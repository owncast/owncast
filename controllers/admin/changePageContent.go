package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/controllers"

	log "github.com/sirupsen/logrus"
)

// ChangeExtraPageContent will change the optional page content
func ChangeExtraPageContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		controllers.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request changeExtraPageContentRequest
	err := decoder.Decode(&request)
	if err != nil {
		log.Errorln(err)
		controllers.WriteSimpleResponse(w, false, "")
		return
	}

	config.Config.InstanceDetails.ExtraPageContent = request.Key
	controllers.WriteSimpleResponse(w, true, "changed")
}

type changeExtraPageContentRequest struct {
	Key string `json:"content"`
}
