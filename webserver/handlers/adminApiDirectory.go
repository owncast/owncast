package handlers

import (
	"net/http"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/webserver/responses"
	log "github.com/sirupsen/logrus"
)

// ResetYPRegistration will clear the YP protocol registration key.
func (h *Handlers) ResetYPRegistration(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Resetting YP registration key")
	if err := data.SetDirectoryRegistrationKey(""); err != nil {
		log.Errorln(err)
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}
	responses.WriteSimpleResponse(w, true, "reset")
}
