package federation

import (
	"io"
	"net/http"
	"strings"

	"github.com/owncast/owncast/services/apfederation/apmodels"
	"github.com/owncast/owncast/services/apfederation/inbox"
	"github.com/owncast/owncast/storage/configrepository"

	log "github.com/sirupsen/logrus"
)

// InboxHandler handles inbound federated requests.
func InboxHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		acceptInboxRequest(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func acceptInboxRequest(w http.ResponseWriter, r *http.Request) {
	configRepository := configrepository.Get()

	if !configRepository.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	urlPathComponents := strings.Split(r.URL.Path, "/")
	var forLocalAccount string
	if len(urlPathComponents) == 5 {
		forLocalAccount = urlPathComponents[3]
	} else {
		log.Errorln("Unable to determine username from url path")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// The account this request is for must match the account name we have set
	// for federation.
	if forLocalAccount != configRepository.GetFederationUsername() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorln("Unable to read inbox request payload", err)
		return
	}

	ib := inbox.Get()
	inboxRequest := apmodels.InboxRequest{Request: r, ForLocalAccount: forLocalAccount, Body: data}
	ib.AddToQueue(inboxRequest)
	w.WriteHeader(http.StatusAccepted)
}
