package admin

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/app/middleware"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/transcoder"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
)

// GetServerConfig gets the config details of the server.
func (c *Controller) GetServerConfig(w http.ResponseWriter, r *http.Request) {
	ffmpeg := utils.ValidatedFfmpegPath(c.Data.GetFfMpegPath())
	usernameBlocklist := c.Data.GetForbiddenUsernameList()
	usernameSuggestions := c.Data.GetSuggestedUsernamesList()

	videoQualityVariants := make([]models.StreamOutputVariant, 0)
	for _, variant := range c.Data.GetStreamOutputVariants() {
		videoQualityVariants = append(videoQualityVariants, models.StreamOutputVariant{
			Name:               variant.GetName(),
			IsAudioPassthrough: variant.GetIsAudioPassthrough(),
			IsVideoPassthrough: variant.IsVideoPassthrough,
			Framerate:          variant.GetFramerate(),
			VideoBitrate:       variant.VideoBitrate,
			AudioBitrate:       variant.AudioBitrate,
			CPUUsageLevel:      variant.CPUUsageLevel,
			ScaledWidth:        variant.ScaledWidth,
			ScaledHeight:       variant.ScaledHeight,
		})
	}
	response := serverConfigAdminResponse{
		InstanceDetails: webConfigResponse{
			Name:                c.Data.GetServerName(),
			Summary:             c.Data.GetServerSummary(),
			Tags:                c.Data.GetServerMetadataTags(),
			ExtraPageContent:    c.Data.GetExtraPageBodyContent(),
			StreamTitle:         c.Data.GetStreamTitle(),
			WelcomeMessage:      c.Data.GetServerWelcomeMessage(),
			OfflineMessage:      c.Data.GetCustomOfflineMessage(),
			Logo:                c.Data.GetLogoPath(),
			SocialHandles:       c.Data.GetSocialHandles(),
			NSFW:                c.Data.GetNSFW(),
			CustomStyles:        c.Data.GetCustomStyles(),
			CustomJavascript:    c.Data.GetCustomJavascript(),
			AppearanceVariables: c.Data.GetCustomColorVariableValues(),
		},
		FFmpegPath:              ffmpeg,
		AdminPassword:           c.Data.GetAdminPassword(),
		StreamKeys:              c.Data.GetStreamKeys(),
		WebServerPort:           config.WebServerPort,
		WebServerIP:             config.WebServerIP,
		RTMPServerPort:          c.Data.GetRTMPPortNumber(),
		ChatDisabled:            c.Data.GetChatDisabled(),
		ChatJoinMessagesEnabled: c.Data.GetChatJoinMessagesEnabled(),
		SocketHostOverride:      c.Data.GetWebsocketOverrideHost(),
		ChatEstablishedUserMode: c.Data.GetChatEstbalishedUsersOnlyMode(),
		HideViewerCount:         c.Data.GetHideViewerCount(),
		VideoSettings: videoSettings{
			VideoQualityVariants: videoQualityVariants,
			LatencyLevel:         c.Data.GetStreamLatencyLevel().Level,
		},
		YP: yp{
			Enabled:     c.Data.GetDirectoryEnabled(),
			InstanceURL: c.Data.GetServerURL(),
		},
		S3:                 c.Data.GetS3Config(),
		ExternalActions:    c.Data.GetExternalActions(),
		SupportedCodecs:    transcoder.GetCodecs(ffmpeg),
		VideoCodec:         c.Data.GetVideoCodec(),
		ForbiddenUsernames: usernameBlocklist,
		SuggestedUsernames: usernameSuggestions,
		Federation: federationConfigResponse{
			Enabled:        c.Data.GetFederationEnabled(),
			IsPrivate:      c.Data.GetFederationIsPrivate(),
			Username:       c.Data.GetFederationUsername(),
			GoLiveMessage:  c.Data.GetFederationGoLiveMessage(),
			ShowEngagement: c.Data.GetFederationShowEngagement(),
			BlockedDomains: c.Data.GetBlockedFederatedDomains(),
		},
		Notifications: notificationsConfigResponse{
			Discord: c.Data.GetDiscordConfig(),
			Browser: c.Data.GetBrowserPushConfig(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	middleware.DisableCache(w)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Errorln(err)
	}
}

type serverConfigAdminResponse struct {
	InstanceDetails         webConfigResponse           `json:"instanceDetails"`
	FFmpegPath              string                      `json:"ffmpegPath"`
	AdminPassword           string                      `json:"adminPassword"`
	StreamKeys              []models.StreamKey          `json:"streamKeys"`
	WebServerPort           int                         `json:"webServerPort"`
	WebServerIP             string                      `json:"webServerIP"`
	RTMPServerPort          int                         `json:"rtmpServerPort"`
	S3                      models.S3                   `json:"s3"`
	VideoSettings           videoSettings               `json:"videoSettings"`
	YP                      yp                          `json:"yp"`
	ChatDisabled            bool                        `json:"chatDisabled"`
	ChatJoinMessagesEnabled bool                        `json:"chatJoinMessagesEnabled"`
	ChatEstablishedUserMode bool                        `json:"chatEstablishedUserMode"`
	ExternalActions         []models.ExternalAction     `json:"externalActions"`
	SupportedCodecs         []string                    `json:"supportedCodecs"`
	VideoCodec              string                      `json:"videoCodec"`
	ForbiddenUsernames      []string                    `json:"forbiddenUsernames"`
	Federation              federationConfigResponse    `json:"federation"`
	SuggestedUsernames      []string                    `json:"suggestedUsernames"`
	SocketHostOverride      string                      `json:"socketHostOverride,omitempty"`
	Notifications           notificationsConfigResponse `json:"notifications"`
	HideViewerCount         bool                        `json:"hideViewerCount"`
}

type videoSettings struct {
	VideoQualityVariants []models.StreamOutputVariant `json:"videoQualityVariants"`
	LatencyLevel         int                          `json:"latencyLevel"`
}

type webConfigResponse struct {
	Name                string                `json:"name"`
	Summary             string                `json:"summary"`
	WelcomeMessage      string                `json:"welcomeMessage"`
	OfflineMessage      string                `json:"offlineMessage"`
	Logo                string                `json:"logo"`
	Tags                []string              `json:"tags"`
	Version             string                `json:"version"`
	NSFW                bool                  `json:"nsfw"`
	ExtraPageContent    string                `json:"extraPageContent"`
	StreamTitle         string                `json:"streamTitle"` // What's going on with the current stream
	SocialHandles       []models.SocialHandle `json:"socialHandles"`
	CustomStyles        string                `json:"customStyles"`
	CustomJavascript    string                `json:"customJavascript"`
	AppearanceVariables map[string]string     `json:"appearanceVariables"`
}

type yp struct {
	Enabled      bool   `json:"enabled"`
	InstanceURL  string `json:"instanceUrl"` // The public URL the directory should link to
	YPServiceURL string `json:"-"`           // The base URL to the YP API to register with (optional)
}

type federationConfigResponse struct {
	Enabled        bool     `json:"enabled"`
	IsPrivate      bool     `json:"isPrivate"`
	Username       string   `json:"username"`
	GoLiveMessage  string   `json:"goLiveMessage"`
	ShowEngagement bool     `json:"showEngagement"`
	BlockedDomains []string `json:"blockedDomains"`
}

type notificationsConfigResponse struct {
	Browser models.BrowserNotificationConfiguration `json:"browser"`
	Discord models.DiscordConfiguration             `json:"discord"`
}
