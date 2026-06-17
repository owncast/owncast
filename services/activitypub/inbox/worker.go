package inbox

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-fed/httpsig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/services/activitypub/apmodels"
)

func (s *Service) handle(request apmodels.InboxRequest) {
	keyOwner, err := s.Verify(request.Request)
	if err != nil {
		log.Debugln("Error in attempting to verify request", err)
		return
	}
	if keyOwner == nil {
		log.Debugln("Request failed verification")
		return
	}

	// Bind the activity to its signer. A validly signed request may only carry
	// activities authored by the verified key's owner; otherwise a server could
	// sign with its own key while claiming, in the body, to be a different actor
	// (e.g. forging a featured-stream Offer/Leave/Accept as another server). We
	// compare origins (hostname), the standard fediverse binding, which is
	// robust to id/fragment formatting differences. Fail closed: if no actor
	// IRI can be determined the activity is rejected.
	actorIRI, err := actorIRIFromActivity(request.Body)
	if err != nil {
		log.Warnln("rejecting inbound activity: unable to bind actor to signing key:", err)
		return
	}
	if !sameActorOrigin(actorIRI, keyOwner) {
		log.Warnf("rejecting inbound activity: actor %q does not match signing key owner %q", actorIRI, keyOwner.String())
		return
	}

	if err := s.resolver.Resolve(context.Background(), request.Body, s.handleUpdateRequest, s.handleFollowInboxRequest, s.handleLikeRequest, s.handleAnnounceRequest, s.handleUndoInboxRequest, s.handleCreateRequest, s.handleOfferInboxRequest, s.handleAcceptInboxRequest, s.handleRejectInboxRequest, s.handleLeaveInboxRequest); err != nil {
		log.Debugln("resolver error:", err)
	}
}

// actorIRIFromActivity extracts the top-level "actor" IRI from a raw inbound
// ActivityPub activity. The actor may be a string IRI, an object with an "id",
// or an array of either; the first usable IRI is returned.
func actorIRIFromActivity(body []byte) (string, error) {
	var envelope struct {
		Actor json.RawMessage `json:"actor"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return "", errors.Wrap(err, "unable to parse activity body")
	}
	return iriFromRawActor(envelope.Actor)
}

func iriFromRawActor(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", errors.New("activity has no actor")
	}

	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		if asString == "" {
			return "", errors.New("actor IRI is empty")
		}
		return asString, nil
	}

	var asObject struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(raw, &asObject); err == nil && asObject.ID != "" {
		return asObject.ID, nil
	}

	var asArray []json.RawMessage
	if err := json.Unmarshal(raw, &asArray); err == nil && len(asArray) > 0 {
		return iriFromRawActor(asArray[0])
	}

	return "", errors.New("unable to determine actor IRI from activity")
}

// sameActorOrigin reports whether the activity actor IRI and the verified key
// owner share the same host. Cross-host mismatches are rejected, which is what
// prevents one server from forging activities as another.
func sameActorOrigin(actorIRI string, keyOwner *url.URL) bool {
	actor, err := url.Parse(actorIRI)
	if err != nil {
		return false
	}
	return actor.Hostname() != "" && strings.EqualFold(actor.Hostname(), keyOwner.Hostname())
}

// maxRequestDateAge / maxRequestDateSkew bound how far an inbound request's
// signed Date may be from now. The window is generous to tolerate clock skew
// and delivery/queue delay while still limiting how long a captured,
// validly-signed request can be replayed.
const (
	maxRequestDateAge  = 1 * time.Hour
	maxRequestDateSkew = 1 * time.Hour
)

// requestDateWithinTolerance checks that a present, parseable Date header is
// recent. It is intentionally lenient: a missing or unparseable Date is not an
// error (to avoid rejecting otherwise-valid federation from senders that omit
// or format it unusually). Senders that sign a standard Date — Owncast,
// Mastodon, etc. — get replay bounding because the signed Date can't be altered
// without invalidating the signature.
func requestDateWithinTolerance(request *http.Request) error {
	dateHeader := request.Header.Get("Date")
	if dateHeader == "" {
		return nil
	}

	sent, err := http.ParseTime(dateHeader)
	if err != nil {
		log.Debugln("inbound request has an unparseable Date header, skipping freshness check:", dateHeader)
		return nil
	}

	now := time.Now()
	if now.Sub(sent) > maxRequestDateAge {
		return fmt.Errorf("request Date is too old (possible replay): %s", dateHeader)
	}
	if sent.Sub(now) > maxRequestDateSkew {
		return fmt.Errorf("request Date is too far in the future: %s", dateHeader)
	}

	return nil
}

// Verify validates the HTTP signature of an inbound request against the
// signing actor's public key and checks it against blocked domains/actors. On
// success it returns the verified key owner's actor IRI (used by the caller to
// bind the activity's actor to its signer); on failure it returns a nil IRI
// and an error.
// nolint: cyclop
func (s *Service) Verify(request *http.Request) (*url.URL, error) {
	verifier, err := httpsig.NewVerifier(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create key verifier for request")
	}
	pubKeyID, err := url.Parse(verifier.KeyId())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse key to get key ID")
	}

	// Force federation only via servers using https.
	if pubKeyID.Scheme != "https" {
		return nil, errors.New("federated servers must use https: " + pubKeyID.String())
	}

	signature := request.Header.Get("signature")
	if signature == "" {
		return nil, errors.New("http signature header not found in request")
	}

	var algorithmString string
	signatureComponents := strings.Split(signature, ",")
	for _, component := range signatureComponents {
		kv := strings.Split(component, "=")
		if kv[0] == "algorithm" {
			algorithmString = kv[1]
			break
		}
	}

	algorithmString = strings.Trim(algorithmString, "\"")
	if algorithmString == "" {
		return nil, errors.New("Unable to determine algorithm to verify request")
	}

	publicKey, err := s.resolver.GetResolvedPublicKeyFromIRI(pubKeyID.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve actor from IRI to fetch key")
	}

	var publicKeyActorIRI *url.URL
	if ownerProp := publicKey.GetW3IDSecurityV1Owner(); ownerProp != nil {
		publicKeyActorIRI = ownerProp.Get()
	}

	if publicKeyActorIRI == nil {
		return nil, errors.New("public key owner IRI is empty")
	}

	// Test to see if the actor is in the list of blocked federated domains.
	if s.isBlockedDomain(publicKeyActorIRI.Hostname()) {
		return nil, errors.New("domain is blocked")
	}

	// If actor is specifically blocked, then fail validation.
	if blocked, err := s.isBlockedActor(publicKeyActorIRI); err != nil || blocked {
		return nil, err
	}

	key, err := apmodels.GetPublicKeyPem(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get public key PEM")
	}
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		log.Errorln("failed to parse PEM block containing the public key")
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Errorln("failed to parse DER encoded public key: " + err.Error())
		return nil, errors.Wrap(err, "failed to parse DER encoded public key")
	}

	algos := []httpsig.Algorithm{
		httpsig.Algorithm(algorithmString), // try stated algorithm first then other common algorithms
		httpsig.RSA_SHA256,                 // <- used by almost all fedi software
		httpsig.RSA_SHA512,
	}

	// The verifier will verify the Digest in addition to the HTTP signature
	triedAlgos := make(map[httpsig.Algorithm]error)
	for _, algorithm := range algos {
		if _, tried := triedAlgos[algorithm]; !tried {
			err := verifier.Verify(parsedKey, algorithm)
			if err == nil {
				// Bound replay of captured, validly-signed requests: a replayed
				// capture carries its original (signed) Date, so reject ones
				// outside a tolerance window.
				if dateErr := requestDateWithinTolerance(request); dateErr != nil {
					return nil, dateErr
				}
				return publicKeyActorIRI, nil
			}
			triedAlgos[algorithm] = err
		}
	}

	return nil, fmt.Errorf("http signature verification error(s) for: %s: %+v", pubKeyID.String(), triedAlgos)
}

func (s *Service) isBlockedDomain(domain string) bool {
	blockedDomains := s.configRepository.GetBlockedFederatedDomains()

	for _, blockedDomain := range blockedDomains {
		if strings.Contains(domain, blockedDomain) {
			return true
		}
	}

	return false
}

func (s *Service) isBlockedActor(actorIRI *url.URL) (bool, error) {
	blockedactor, err := s.followers.GetByIRI(actorIRI.String())

	if blockedactor != nil && blockedactor.DisabledAt != nil {
		return true, errors.Wrap(err, "remote actor is blocked")
	}

	return false, nil
}
