package data

import (
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
)

// HasPopulatedDefaults will determine if the defaults have been inserted into the database.
func HasPopulatedDefaults() bool {
	hasPopulated, err := _datastore.GetBool("HAS_POPULATED_DEFAULTS")
	if err != nil {
		return false
	}
	return hasPopulated
}

// PopulateDefaults will set default values in the database.
func PopulateDefaults() {
	defaults := config.GetDefaults()

	if HasPopulatedDefaults() {
		return
	}

	_ = SetStreamKey(defaults.StreamKey)
	_ = SetHTTPPortNumber(float64(defaults.WebServerPort))
	_ = SetRTMPPortNumber(float64(defaults.RTMPServerPort))
	_ = SetLogoPath(defaults.Logo)
	_ = SetServerMetadataTags([]string{"owncast", "streaming"})
	_ = SetServerSummary("Welcome to your new Owncast server!  This description can be changed in the admin")
	_ = SetServerName("Owncast")
	_ = SetStreamKey(defaults.StreamKey)
	_ = SetExtraPageBodyContent("This is your page's content that can be edited in the admin.")
	_ = SetSocialHandles([]models.SocialHandle{
		{
			Platform: "github",
			URL:      "https://github.com/owncast/owncast",
		},
	})

	_datastore.warmCache()
	_ = _datastore.SetBool("HAS_POPULATED_DEFAULTS", true)
}
