package inbox

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/go-fed/httpsig"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/core/data"

	log "github.com/sirupsen/logrus"
)

func handle(request apmodels.InboxRequest) {
	if verified, err := Verify(request.Request); err != nil {
		log.Errorln("Error in attempting to verify request", err)
		return
	} else if !verified {
		log.Debugln("Request failed verification", err)
		return
	}

	if err := resolvers.Resolve(context.Background(), request.Body, handleUpdateRequest, handleFollowInboxRequest, handleLikeRequest, handleAnnounceRequest, handleUndoInboxRequest, handleCreateRequest); err != nil {
		log.Debugln("resolver error:", err)
	}
}

// Verify will Verify the http signature of an inbound request as well as
// check it against the list of blocked domains.
// nolint: cyclop
func Verify(request *http.Request) (bool, error) {
	verifier, err := httpsig.NewVerifier(request)
	if err != nil {
		return false, errors.Wrap(err, "failed to create key verifier for request")
	}
	pubKeyID, err := url.Parse(verifier.KeyId())
	if err != nil {
		return false, errors.Wrap(err, "failed to parse key to get key ID")
	}

	// Force federation only via servers using https.
	if pubKeyID.Scheme != "https" {
		return false, errors.New("federated servers must use https: " + pubKeyID.String())
	}

	signature := request.Header.Get("signature")
	if signature == "" {
		return false, errors.New("http signature header not found in request")
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
		return false, errors.New("Unable to determine algorithm to verify request")
	}

	publicKey, err := resolvers.GetResolvedPublicKeyFromIRI(pubKeyID.String())
	if err != nil {
		return false, errors.Wrap(err, "failed to resolve actor from IRI to fetch key")
	}

	var publicKeyActorIRI *url.URL
	if ownerProp := publicKey.GetW3IDSecurityV1Owner(); ownerProp != nil {
		publicKeyActorIRI = ownerProp.Get()
	}

	if publicKeyActorIRI == nil {
		return false, errors.New("public key owner IRI is empty")
	}

	// Test to see if the actor is in the list of blocked federated domains.
	if isBlockedDomain(publicKeyActorIRI.Hostname()) {
		return false, errors.New("domain is blocked")
	}

	// If actor is specifically blocked, then fail validation.
	if blocked, err := isBlockedActor(publicKeyActorIRI); err != nil || blocked {
		return false, err
	}

	key := publicKey.GetW3IDSecurityV1PublicKeyPem().Get()
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		log.Errorln("failed to parse PEM block containing the public key")
		return false, errors.New("failed to parse PEM block containing the public key")
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Errorln("failed to parse DER encoded public key: " + err.Error())
		return false, errors.Wrap(err, "failed to parse DER encoded public key")
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
				return true, nil
			}
			triedAlgos[algorithm] = err
		}
	}

	return false, fmt.Errorf("http signature verification error(s) for: %s: %+v", pubKeyID.String(), triedAlgos)
}

func isBlockedDomain(domain string) bool {
	blockedDomains := data.GetBlockedFederatedDomains()

	for _, blockedDomain := range blockedDomains {
		if strings.Contains(domain, blockedDomain) {
			return true
		}
	}

	return false
}

func isBlockedActor(actorIRI *url.URL) (bool, error) {
	blockedactor, err := persistence.GetFollower(actorIRI.String())

	if blockedactor != nil && blockedactor.DisabledAt != nil {
		return true, errors.Wrap(err, "remote actor is blocked")
	}

	return false, nil
}
