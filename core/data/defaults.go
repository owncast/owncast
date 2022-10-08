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

	const defaultPageContent = `
This is a live stream powered by [Owncast](https://owncast.online), a free and open source livestreaming server.
<hr/>

### Why you might be interested in Owncast

- It can be self-hosted on your own server.
- It's easy to install and configure. You can be up and running in minutes.
- It's not controlled by a corporation and is built by people like [you](https://owncast.online/help).
- Viewers can start chatting immediately, no account required.
- No ads, tracking or data collection.

<hr/>

To discover more streams, visit [Owncast's directory](https://directory.owncast.online).
	`

	_ = SetStreamKey(defaults.StreamKey)
	_ = SetHTTPPortNumber(float64(defaults.WebServerPort))
	_ = SetRTMPPortNumber(float64(defaults.RTMPServerPort))
	_ = SetLogoPath(defaults.Logo)
	_ = SetServerMetadataTags([]string{"owncast", "streaming"})
	_ = SetServerSummary("Welcome to your new Owncast server! This description can be changed in the admin. Visit https://owncast.online/docs/configuration/ to learn more.")
	_ = SetServerWelcomeMessage("")
	_ = SetServerName("Owncast")
	_ = SetStreamKey(defaults.StreamKey)
	_ = SetExtraPageBodyContent(defaultPageContent)
	_ = SetFederationGoLiveMessage(defaults.FederationGoLiveMessage)
	_ = SetSocialHandles([]models.SocialHandle{
		{
			Platform: "github",
			URL:      "https://github.com/owncast/owncast",
		},
	})

	_ = _datastore.SetBool("HAS_POPULATED_DEFAULTS", true)
}
