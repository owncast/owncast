package data

import (
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
)

// Determine if the defaults have been inserted into the database
func HasPopulatedDefaults() bool {
	hasPopulated, err := _datastore.GetBool("HAS_POPULATED_DEFAULTS")
	if err != nil {
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
	SetHTTPPortNumber(float64(defaults.WebServerPort))
	SetRTMPPortNumber(float64(defaults.RTMPServerPort))
	SetLogoPath(defaults.Logo)
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
