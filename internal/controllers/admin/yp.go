package admin

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/internal/controllers"
	"github.com/owncast/owncast/internal/core/data"
)

// ResetYPRegistration will clear the YP protocol registration key.
func ResetYPRegistration(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Resetting YP registration key")
	if err := data.SetDirectoryRegistrationKey(""); err != nil {
		log.Errorln(err)
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}
	controllers.WriteSimpleResponse(w, true, "reset")
}
