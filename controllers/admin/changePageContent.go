package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/data"

	log "github.com/sirupsen/logrus"
)

// ChangeExtraPageContent will change the optional page content.
func ChangeExtraPageContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != controllers.POST {
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

	if err := data.SetExtraPageBodyContent(request.Content); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "changed")
}

type changeExtraPageContentRequest struct {
	Content string `json:"content"`
}
