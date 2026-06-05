package admin

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	webutils "github.com/owncast/owncast/webserver/utils"
)

// ResetYPRegistration will clear the YP protocol registration key.
func (a *Admin) ResetYPRegistration(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Resetting YP registration key")
	if err := a.configRepository.SetDirectoryRegistrationKey(""); err != nil {
		log.Errorln(err)
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}
	webutils.WriteSimpleResponse(w, true, "reset")
}
