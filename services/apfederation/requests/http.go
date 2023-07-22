package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/services/apfederation/crypto"
	"github.com/owncast/owncast/services/config"

	log "github.com/sirupsen/logrus"
)

// WriteStreamResponse will write a ActivityPub object to the provided ResponseWriter and sign with the provided key.
func (r *Requests) WriteStreamResponse(item vocab.Type, w http.ResponseWriter, publicKey crypto.PublicKey) error {
	var jsonmap map[string]interface{}
	jsonmap, _ = streams.Serialize(item)
	b, err := json.Marshal(jsonmap)
	if err != nil {
		return err
	}

	return r.WriteResponse(b, w, publicKey)
}

// WritePayloadResponse will write any arbitrary object to the provided ResponseWriter and sign with the provided key.
func (r *Requests) WritePayloadResponse(payload interface{}, w http.ResponseWriter, publicKey crypto.PublicKey) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return r.WriteResponse(b, w, publicKey)
}

// WriteResponse will write any arbitrary payload to the provided ResponseWriter and sign with the provided key.
func (r *Requests) WriteResponse(payload []byte, w http.ResponseWriter, publicKey crypto.PublicKey) error {
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

// CreateSignedRequest will create a signed POST request of a payload to the provided destination.
func (r *Requests) CreateSignedRequest(payload []byte, url *url.URL, fromActorIRI *url.URL) (*http.Request, error) {
	log.Debugln("Sending", string(payload), "to", url)

	req, _ := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(payload))
	c := config.Get()
	ua := fmt.Sprintf("%s; https://owncast.online", c.GetReleaseString())
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Content-Type", "application/activity+json")

	if err := crypto.SignRequest(req, payload, fromActorIRI); err != nil {
		log.Errorln("error signing request:", err)
		return nil, err
	}

	return req, nil
}
