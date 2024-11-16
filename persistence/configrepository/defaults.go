package configrepository

import (
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

// PopulateDefaults will set default values in the database.
func (r *SqlConfigRepository) PopulateDefaults() {
	key := "HAS_POPULATED_DEFAULTS"

	r.datastore.WarmCache()

	defaults := config.GetDefaults()

	_ = r.SetAdminPassword(defaults.AdminPassword)
	_ = r.SetStreamKeys(defaults.StreamKeys)
	_ = r.SetHTTPPortNumber(float64(defaults.WebServerPort))
	_ = r.SetRTMPPortNumber(float64(defaults.RTMPServerPort))
	_ = r.SetLogoPath(defaults.Logo)
	_ = r.SetServerMetadataTags([]string{"owncast", "streaming"})
	_ = r.SetServerSummary(defaults.Summary)
	_ = r.SetServerWelcomeMessage("")
	_ = r.SetServerName(defaults.Name)
	_ = r.SetExtraPageBodyContent(defaults.PageBodyContent)
	_ = r.SetFederationGoLiveMessage(defaults.FederationGoLiveMessage)
	_ = r.SetSocialHandles([]models.SocialHandle{
		{
			Platform: "github",
			URL:      "https://github.com/owncast/owncast",
		},
	})

	if !r.HasPopulatedFederationDefaults() {
		if err := r.SetFederationGoLiveMessage(config.GetDefaults().FederationGoLiveMessage); err != nil {
			log.Errorln(err)
		}
		if err := r.datastore.SetBool("HAS_POPULATED_FEDERATION_DEFAULTS", true); err != nil {
			log.Errorln(err)
		}
	}

	_ = r.datastore.SetBool(key, true)
}

// HasPopulatedDefaults will determine if the defaults have been inserted into the database.
func (r *SqlConfigRepository) HasPopulatedDefaults() bool {
	hasPopulated, err := r.datastore.GetBool("HAS_POPULATED_DEFAULTS")
	if err != nil {
		return false
	}
	return hasPopulated
}

func (r *SqlConfigRepository) HasPopulatedFederationDefaults() bool {
	hasPopulated, err := r.datastore.GetBool("HAS_POPULATED_FEDERATION_DEFAULTS")
	if err != nil {
		return false
	}
	return hasPopulated
}
