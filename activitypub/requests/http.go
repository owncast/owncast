package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

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
