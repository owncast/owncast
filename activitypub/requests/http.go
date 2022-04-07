package requests

import (
	"encoding/json"
	"net/http"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/crypto"

	log "github.com/sirupsen/logrus"
)

// WriteStreamResponse will write a ActivityPub object to the provided ResponseWriter and sign with the provided key.
func WriteStreamResponse(item vocab.Type, w http.ResponseWriter, publicKey crypto.PublicKey) error {
	var jsonmap map[string]interface{}
	jsonmap, _ = streams.Serialize(item)
	b, err := json.Marshal(jsonmap)
	if err != nil {
		return err
	}

	return WriteResponse(b, w, publicKey)
}

// WritePayloadResponse will write any arbitrary object to the provided ResponseWriter and sign with the provided key.
func WritePayloadResponse(payload interface{}, w http.ResponseWriter, publicKey crypto.PublicKey) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return WriteResponse(b, w, publicKey)
}

// WriteResponse will write any arbitrary payload to the provided ResponseWriter and sign with the provided key.
func WriteResponse(payload []byte, w http.ResponseWriter, publicKey crypto.PublicKey) error {
	w.Header().Set("Content-Type", "application/activity+json")

	if err := crypto.SignResponse(w, payload, publicKey); err != nil {
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
