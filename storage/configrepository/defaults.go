package configrepository

import (
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/config"
)

// HasPopulatedDefaults will determine if the defaults have been inserted into the database.
func (cr *SqlConfigRepository) HasPopulatedDefaults() bool {
	hasPopulated, err := cr.datastore.GetBool("HAS_POPULATED_DEFAULTS")
	if err != nil {
		return false
	}
	return hasPopulated
}

func (cr *SqlConfigRepository) hasPopulatedFederationDefaults() bool {
	hasPopulated, err := cr.datastore.GetBool("HAS_POPULATED_FEDERATION_DEFAULTS")
	if err != nil {
		return false
	}
	return hasPopulated
}

// PopulateDefaults will set default values in the database.
func (cr *SqlConfigRepository) PopulateDefaults() {
	defaults := config.GetDefaults()

	if cr.HasPopulatedDefaults() {
		return
	}

	_ = cr.SetAdminPassword(defaults.AdminPassword)
	_ = cr.SetStreamKeys(defaults.StreamKeys)
	_ = cr.SetHTTPPortNumber(float64(defaults.WebServerPort))
	_ = cr.SetRTMPPortNumber(float64(defaults.RTMPServerPort))
	_ = cr.SetLogoPath(defaults.Logo)
	_ = cr.SetServerMetadataTags([]string{"owncast", "streaming"})
	_ = cr.SetServerSummary(defaults.Summary)
	_ = cr.SetServerWelcomeMessage("")
	_ = cr.SetServerName(defaults.Name)
	_ = cr.SetExtraPageBodyContent(defaults.PageBodyContent)
	_ = cr.SetFederationGoLiveMessage(defaults.FederationGoLiveMessage)
	_ = cr.SetSocialHandles([]models.SocialHandle{
		{
			Platform: "github",
			URL:      "https://github.com/owncast/owncast",
		},
	})

	_ = cr.datastore.SetBool("HAS_POPULATED_DEFAULTS", true)
}
