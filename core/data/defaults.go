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

func hasPopulatedFederationDefaults() bool {
	hasPopulated, err := _datastore.GetBool("HAS_POPULATED_FEDERATION_DEFAULTS")
	if err != nil {
		return false
	}
	return hasPopulated
}

// PopulateDefaults will set default values in the database.
func PopulateDefaults() {
	_datastore.warmCache()

	defaults := config.GetDefaults()

	if HasPopulatedDefaults() {
		return
	}

	_ = SetAdminPassword(defaults.AdminPassword)
	_ = SetStreamKeys(defaults.StreamKeys)
	_ = SetHTTPPortNumber(float64(defaults.WebServerPort))
	_ = SetRTMPPortNumber(float64(defaults.RTMPServerPort))
	_ = SetLogoPath(defaults.Logo)
	_ = SetServerMetadataTags([]string{"owncast", "streaming"})
	_ = SetServerSummary(defaults.Summary)
	_ = SetServerWelcomeMessage("")
	_ = SetServerName(defaults.Name)
	_ = SetExtraPageBodyContent(defaults.PageBodyContent)
	_ = SetFederationGoLiveMessage(defaults.FederationGoLiveMessage)
	_ = SetSocialHandles([]models.SocialHandle{
		{
			Platform: "github",
			URL:      "https://github.com/owncast/owncast",
		},
	})

	_ = _datastore.SetBool("HAS_POPULATED_DEFAULTS", true)
}
