package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/controllers"

	log "github.com/sirupsen/logrus"
)

// ChangeStreamTags will change the stream key (in memory).
func ChangeStreamTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != controllers.POST {
		controllers.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request changeStreamTagsRequest
	err := decoder.Decode(&request)
	if err != nil {
		log.Errorln(err)
		controllers.WriteSimpleResponse(w, false, "")
		return
	}

	config.Config.InstanceDetails.Tags = request.Tags
	controllers.WriteSimpleResponse(w, true, "changed")
}

type changeStreamTagsRequest struct {
	Tags []string `json:"tags"`
}
