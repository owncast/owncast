package admin

import (
	"net/http"

	"github.com/owncast/owncast/core/data"
	webutils "github.com/owncast/owncast/webserver/utils"
	log "github.com/sirupsen/logrus"
)

// ResetYPRegistration will clear the YP protocol registration key.
func ResetYPRegistration(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Resetting YP registration key")
	if err := data.SetDirectoryRegistrationKey(""); err != nil {
		log.Errorln(err)
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}
	webutils.WriteSimpleResponse(w, true, "reset")
}
