package requests

import (
	"encoding/json"
	"net/http"
	"net/url"

	"code.superseriousbusiness.org/activity/streams"
	"code.superseriousbusiness.org/activity/streams/vocab"

	"github.com/owncast/owncast/services/activitypub/crypto"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// WriteStreamResponse will write a ActivityPub object to the provided ResponseWriter and sign with the provided key.
func WriteStreamResponse(item vocab.Type, w http.ResponseWriter, publicKey crypto.PublicKey, signer *crypto.Signer) error {
	var jsonmap map[string]interface{}
	jsonmap, _ = streams.Serialize(item)
	b, err := json.Marshal(jsonmap)
	if err != nil {
		return err
	}

	return WriteResponse(b, w, publicKey, signer)
}

// WritePayloadResponse will write any arbitrary object to the provided ResponseWriter and sign with the provided key.
func WritePayloadResponse(payload interface{}, w http.ResponseWriter, publicKey crypto.PublicKey, signer *crypto.Signer) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return WriteResponse(b, w, publicKey, signer)
}

// WriteResponse will write any arbitrary payload to the provided ResponseWriter and sign with the provided key.
func WriteResponse(payload []byte, w http.ResponseWriter, publicKey crypto.PublicKey, signer *crypto.Signer) error {
	w.Header().Set("Content-Type", "application/activity+json")

	if err := signer.SignResponse(w, payload, publicKey); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln("unable to sign response", err)
		return err
	}

	if _, err := w.Write(payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}

// validateRemoteInbox rejects malformed or non-HTTPS inbox URLs. The durable
// queue performs public-address validation before storing the delivery.
func validateRemoteInbox(inbox *url.URL) error {
	if inbox == nil || inbox.Scheme != "https" || inbox.Hostname() == "" {
		return errors.Errorf("rejecting invalid inbox URL: %s", inbox)
	}
	return nil
}
