package admin

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// ResetYPRegistration will clear the YP protocol registration key.
func (c *Controller) ResetYPRegistration(w http.ResponseWriter, r *http.Request) {
	log.Traceln("Resetting YP registration key")
	if err := c.Data.SetDirectoryRegistrationKey(""); err != nil {
		log.Errorln(err)
		c.Service.WriteSimpleResponse(w, false, err.Error())
		return
	}
	c.Service.WriteSimpleResponse(w, true, "reset")
}
