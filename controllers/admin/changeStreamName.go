package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/controllers"

	log "github.com/sirupsen/logrus"
)

// ChangeStreamName will change the stream key (in memory).
func ChangeStreamName(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		controllers.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request changeStreamNameRequest
	err := decoder.Decode(&request)
	if err != nil {
		log.Errorln(err)
		controllers.WriteSimpleResponse(w, false, "")
		return
	}

	config.Config.InstanceDetails.Name = request.Name
	controllers.WriteSimpleResponse(w, true, "changed")
}

type changeStreamNameRequest struct {
	Name string `json:"name"`
}
