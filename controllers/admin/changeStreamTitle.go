package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/controllers"

	log "github.com/sirupsen/logrus"
)

// ChangeStreamTitle will change the stream key (in memory).
func ChangeStreamTitle(w http.ResponseWriter, r *http.Request) {
	if r.Method != controllers.POST {
		controllers.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request changeStreamTitleRequest
	err := decoder.Decode(&request)
	if err != nil {
		log.Errorln(err)
		controllers.WriteSimpleResponse(w, false, "")
		return
	}

	config.Config.InstanceDetails.Title = request.Title
	controllers.WriteSimpleResponse(w, true, "changed")
}

type changeStreamTitleRequest struct {
	Title string `json:"title"`
}
