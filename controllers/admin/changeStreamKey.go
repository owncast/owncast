package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/controllers"

	log "github.com/sirupsen/logrus"
)

// ChangeStreamKey will change the stream key (in memory).
func ChangeStreamKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		controllers.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request changeStreamKeyRequest
	err := decoder.Decode(&request)
	if err != nil {
		log.Errorln(err)
		controllers.WriteSimpleResponse(w, false, "")
		return
	}

	config.Config.VideoSettings.StreamingKey = request.Key
	controllers.WriteSimpleResponse(w, true, "changed")
}

type changeStreamKeyRequest struct {
	Key string `json:"key"`
}
