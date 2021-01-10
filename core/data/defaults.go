package data

import (
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

// Determine if the defaults have been inserted into the database
func HasPopulatedDefaults() bool {
	hasPopulated, err := _datastore.GetBool("HAS_POPULATED_DEFAULTS")
	if err != nil {
		log.Errorln(err)
		return false
	}
	return hasPopulated
}

func PopulateDefaults() {
	defaults := config.GetDefaults()

	if HasPopulatedDefaults() {
		return
	}

	SetStreamKey(defaults.StreamKey)
	SetHTTPPortNumber(defaults.WebServerPort)
	SetRTMPPortNumber(defaults.RTMPServerPort)
	SetLogoPath(defaults.InstanceDetails.Logo)
	SetServerMetadataTags([]string{"owncast", "streaming"})
	SetServerSummary("Welcome to your new Owncast server!  This description can be changed in the admin")
	SetServerName("Owncast")
	SetStreamKey(defaults.StreamKey)
	SetExtraPageBodyContent("This is your page's content that can be edited in the admin.")
	SetSocialHandles([]models.SocialHandle{
		{
			Platform: "github",
			URL:      "https://github.com/owncast/owncast",
		},
	})

	_datastore.warmCache()
	_datastore.SetBool("HAS_POPULATED_DEFAULTS", true)

}
