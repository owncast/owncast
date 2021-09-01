package inbox

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/go-fed/httpsig"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/resolvers"

	log "github.com/sirupsen/logrus"
)

var _queue = make(chan apmodels.InboxRequest, 5)

func init() {
	go run()
}

func run() {
	for r := range _queue {
		handle(r)
	}
}

func Add(request apmodels.InboxRequest) {
	_queue <- request
}

func handle(request apmodels.InboxRequest) chan bool {
	c := context.WithValue(context.Background(), "account", request.ForLocalAccount)
	r := make(chan bool)

	if verified, err := verify(request.Request); err != nil || !verified {
		log.Warnln("Unable to verify remote request", err)
		return nil
	}

	createCallback := func(c context.Context, activity vocab.ActivityStreamsCreate) error {
		r <- false
		return nil
	}

	updateCallback := func(c context.Context, activity vocab.ActivityStreamsUpdate) error {
		r <- false
		return nil
	}

	personCallback := func(c context.Context, activity vocab.ActivityStreamsPerson) error {
		r <- false
		return nil
	}

	deleteCallback := func(c context.Context, activity vocab.ActivityStreamsDelete) error {
		r <- false
		return nil
	}

	if err := resolvers.Resolve(request.Body, c, createCallback, deleteCallback, updateCallback, handleFollowInboxRequest, personCallback, handleLikeRequest, handleAnnounceRequest, handleUndoInboxRequest); err != nil {
		log.Errorln("resolver error:", err)
	}

	return r
}

// Verify will verify the http signature of an inbound request.
func verify(request *http.Request) (bool, error) {
	verifier, err := httpsig.NewVerifier(request)
	if err != nil {
		return false, err
	}
	pubKeyId, err := url.Parse(verifier.KeyId())
	if err != nil {
		return false, err
	}

	log.Traceln("Fetching key", pubKeyId)

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

	var actor vocab.ActivityStreamsPerson
	personCallback := func(c context.Context, person vocab.ActivityStreamsPerson) error {
		actor = person
		return nil
	}

	resolvers.ResolveIRI(pubKeyId.String(), context.TODO(), personCallback)

	if actor == nil {
		return false, errors.New("unable to resolve actor to fetch key " + pubKeyId.String())
	}

	key := actor.GetW3IDSecurityV1PublicKey().Begin().Get().GetW3IDSecurityV1PublicKeyPem().Get()
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		log.Errorln("failed to parse PEM block containing the public key")
		return false, errors.New("failed to parse PEM block containing the public key")
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Errorln("failed to parse DER encoded public key: " + err.Error())
		return false, errors.New("failed to parse DER encoded public key: " + err.Error())
	}

	if err != nil {
		return false, err
	}

	var algorithm httpsig.Algorithm

	switch algorithmString {
	case "rsa-sha256":
		algorithm = httpsig.RSA_SHA256
	}

	// The verifier will verify the Digest in addition to the HTTP signature
	if err := verifier.Verify(parsedKey, algorithm); err != nil {
		log.Warnln("verification error for", pubKeyId, err)
		return false, err
	}

	return true, nil
}
