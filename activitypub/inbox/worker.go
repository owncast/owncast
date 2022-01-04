package inbox

import (
	"context"
	"crypto/x509"
	"encoding/pem"
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
		log.Debugln("Error in attempting to verify request", err)
		return
	} else if !verified {
		log.Errorln("Request failed verification")
		return
	}

	// c := context.WithValue(context.Background(), "account", request.ForLocalAccount) //nolint

	if err := resolvers.Resolve(context.Background(), request.Body, handleUpdateRequest, handleFollowInboxRequest, handleLikeRequest, handleAnnounceRequest, handleUndoInboxRequest); err != nil {
		log.Errorln("resolver error:", err)
	}
}

// Verify will Verify the http signature of an inbound request as well as
// check it against the list of blocked domains.
func Verify(request *http.Request) (bool, error) {
	blockedDomains := data.GetBlockedFederatedDomains()

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

	actor, err := resolvers.GetResolvedActorFromIRI(pubKeyID.String())
	if err != nil {
		return false, errors.Wrap(err, "failed to resolve actor from IRI to fetch key")
	}

	// Test to see if the actor is in the list of blocked federated domains.
	for _, blockedDomain := range blockedDomains {
		if strings.Contains(actor.ActorIri.Host, blockedDomain) {
			return false, errors.New("actor domain is blocked: " + blockedDomain)
		}
	}

	// If actor is specifically blocked, then fail validation.
	blockedactor, err := persistence.GetFollower(actor.ActorIri.String())
	if err != nil {
		return false, errors.Wrap(err, "error validating actor against blocked actors")
	}
	if blockedactor != nil && blockedactor.DisabledAt != nil {
		return false, errors.Wrap(err, "remote actor is blocked")
	}

	key := actor.W3IDSecurityV1PublicKey.Begin().Get().GetW3IDSecurityV1PublicKeyPem().Get()
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

	var algorithm httpsig.Algorithm = httpsig.Algorithm(algorithmString)

	// The verifier will verify the Digest in addition to the HTTP signature
	if err := verifier.Verify(parsedKey, algorithm); err != nil {
		log.Warnln("verification error for", pubKeyID, err)
		return false, errors.Wrap(err, "verification error: "+pubKeyID.String())
	}

	return true, nil
}
