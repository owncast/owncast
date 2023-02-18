package data

import (
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
)

// HasPopulatedDefaults will determine if the defaults have been inserted into the database.
func (s *Service) HasPopulatedDefaults() bool {
	hasPopulated, err := s.Store.GetBool("HAS_POPULATED_DEFAULTS")
	if err != nil {
		return false
	}
	return hasPopulated
}

func (s *Service) hasPopulatedFederationDefaults() bool {
	hasPopulated, err := s.Store.GetBool("HAS_POPULATED_FEDERATION_DEFAULTS")
	if err != nil {
		return false
	}
	return hasPopulated
}

// PopulateDefaults will set default values in the database.
func (s *Service) PopulateDefaults() {
	s.Store.warmCache()

	defaults := config.GetDefaults()

	if s.HasPopulatedDefaults() {
		return
	}

	_ = s.SetAdminPassword(defaults.AdminPassword)
	_ = s.SetStreamKeys(defaults.StreamKeys)
	_ = s.SetHTTPPortNumber(float64(defaults.WebServerPort))
	_ = s.SetRTMPPortNumber(float64(defaults.RTMPServerPort))
	_ = s.SetLogoPath(defaults.Logo)
	_ = s.SetServerMetadataTags([]string{"owncast", "streaming"})
	_ = s.SetServerSummary(defaults.Summary)
	_ = s.SetServerWelcomeMessage("")
	_ = s.SetServerName(defaults.Name)
	_ = s.SetExtraPageBodyContent(defaults.PageBodyContent)
	_ = s.SetFederationGoLiveMessage(defaults.FederationGoLiveMessage)
	_ = s.SetSocialHandles([]models.SocialHandle{
		{
			Platform: "github",
			URL:      "https://github.com/owncast/owncast",
		},
	})

	_ = s.Store.SetBool("HAS_POPULATED_DEFAULTS", true)
}
