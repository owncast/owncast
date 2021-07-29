package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/models"
)

func WriteStreamResponse(item vocab.Type, w http.ResponseWriter, publicKey models.PublicKey) error {
	var jsonmap map[string]interface{}
	jsonmap, _ = streams.Serialize(item)
	b, err := json.Marshal(jsonmap)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	if err := crypto.SignResponse(w, b, publicKey); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return err
	}

	fmt.Println(string(b))
	if _, err := w.Write(b); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println(string(b))

	return nil
}

func PostSignedRequest(payload []byte, url *url.URL, fromActorIRI *url.URL) ([]byte, error) {
	fmt.Println("Sending", string(payload), "to", url)

	req, _ := http.NewRequest("POST", url.String(), bytes.NewBuffer(payload))
	if err := crypto.SignRequest(req, payload, fromActorIRI); err != nil {
		fmt.Println("error signing request:", err)
		return nil, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println("Response: ", response.StatusCode, string(body))
	return body, nil
}
