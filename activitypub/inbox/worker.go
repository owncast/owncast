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
	"github.com/owncast/owncast/activitypub/models"
	"github.com/owncast/owncast/activitypub/resolvers"

	log "github.com/sirupsen/logrus"
)

var _queue = make(chan models.InboxRequest, 5)

func init() {
	go run()
}

func run() {
	for r := range _queue {
		fmt.Println("run...")
		handle(r)
	}
}

func Add(request models.InboxRequest) {
	fmt.Println("Adding AP Payload...")
	_queue <- request
}

func handle(request models.InboxRequest) chan bool {
	c := context.WithValue(context.Background(), "account", request.ForLocalAccount)
	r := make(chan bool)

	if verified, err := verify(request.Request); err != nil || !verified {
		log.Errorln("Unable to verify remote request", err)
		return nil
	}

	fmt.Println("Handling payload via worker...")

	createCallback := func(c context.Context, activity vocab.ActivityStreamsCreate) error {
		fmt.Println("createCallback fired!")
		fmt.Println(activity)
		r <- false
		return nil
	}

	updateCallback := func(c context.Context, activity vocab.ActivityStreamsUpdate) error {
		fmt.Println("updateCallback fired!")

		fmt.Println(activity)
		r <- false
		return nil
	}

	personCallback := func(c context.Context, activity vocab.ActivityStreamsPerson) error {
		fmt.Println("personCallback fired!")
		r <- false
		return nil
	}

	if err := resolvers.Resolve(request.Body, c, createCallback, updateCallback, handleFollowInboxRequest, personCallback, handleUndoInboxRequest); err != nil {
		panic(err)
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

	fmt.Println("Fetching key", pubKeyId)

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
		fmt.Println("personCallback fired!")
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
		panic("failed to parse PEM block containing the public key")
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
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
	verificationError := verifier.Verify(parsedKey, algorithm)
	fmt.Println(verificationError)

	return verificationError == nil, verificationError
}
