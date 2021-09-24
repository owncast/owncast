package inbox

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
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

// Add will place an InboxRequest into the worker queue to be processed.
func Add(request apmodels.InboxRequest) {
	_queue <- request
}

func handle(request apmodels.InboxRequest) chan bool {
	c := context.WithValue(context.Background(), "account", request.ForLocalAccount) //nolint
	r := make(chan bool)

	createCallback := func(c context.Context, activity vocab.ActivityStreamsCreate) error {
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

	if err := resolvers.Resolve(c, request.Body, createCallback, deleteCallback, handleUpdateRequest, handleFollowInboxRequest, personCallback, handleLikeRequest, handleAnnounceRequest, handleUndoInboxRequest); err != nil {
		log.Errorln("resolver error:", err)
	}

	return r
}

// Verify will Verify the http signature of an inbound request.
func Verify(request *http.Request) (bool, error) {
	verifier, err := httpsig.NewVerifier(request)
	if err != nil {
		return false, err
	}
	pubKeyID, err := url.Parse(verifier.KeyId())
	if err != nil {
		return false, err
	}

	log.Println("Fetching key", pubKeyID)

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
		return false, errors.New("unable to determine algorithm to verify request")
	}

	var actor vocab.ActivityStreamsPerson
	var application vocab.ActivityStreamsApplication
	personCallback := func(c context.Context, person vocab.ActivityStreamsPerson) error {
		actor = person
		return nil
	}

	applicationCallback := func(c context.Context, a vocab.ActivityStreamsApplication) error {
		application = a
		return nil
	}

	if err := resolvers.ResolveIRI(context.TODO(), pubKeyID.String(), personCallback, applicationCallback); err != nil {
		return false, fmt.Errorf("%v %v", err, pubKeyID)
	}

	var pk vocab.W3IDSecurityV1PublicKeyProperty
	if actor != nil {
		pk = actor.GetW3IDSecurityV1PublicKey()
	} else if application != nil {
		pk = application.GetW3IDSecurityV1PublicKey()
	} else {
		return false, errors.New("unable to resolve actor as either a person or an application to fetch public key")
	}

	key := pk.Begin().Get().GetW3IDSecurityV1PublicKeyPem().Get()
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
		log.Warnln("verification error for", pubKeyID, err)
		return false, err
	}

	return true, nil
}
