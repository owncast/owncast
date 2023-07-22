package cmd

import log "github.com/sirupsen/logrus"

func (app *Application) showStartupMessage() {
	log.Infoln(app.configservice.GetReleaseString())
}
