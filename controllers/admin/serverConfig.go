package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/transcoder"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

// GetServerConfig gets the config details of the server.
func GetServerConfig(w http.ResponseWriter, r *http.Request) {
	ffmpeg := utils.ValidatedFfmpegPath(data.GetFfMpegPath())
	usernameBlocklist := data.GetForbiddenUsernameList()
	usernameSuggestions := data.GetSuggestedUsernamesList()

	videoQualityVariants := make([]models.StreamOutputVariant, 0)
	for _, variant := range data.GetStreamOutputVariants() {
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
			Name:                data.GetServerName(),
			Summary:             data.GetServerSummary(),
			Tags:                data.GetServerMetadataTags(),
			ExtraPageContent:    data.GetExtraPageBodyContent(),
			StreamTitle:         data.GetStreamTitle(),
			WelcomeMessage:      data.GetServerWelcomeMessage(),
			OfflineMessage:      data.GetCustomOfflineMessage(),
			Logo:                data.GetLogoPath(),
			SocialHandles:       data.GetSocialHandles(),
			NSFW:                data.GetNSFW(),
			CustomStyles:        data.GetCustomStyles(),
			CustomJavascript:    data.GetCustomJavascript(),
			AppearanceVariables: data.GetCustomColorVariableValues(),
		},
		FFmpegPath:              ffmpeg,
		AdminPassword:           data.GetAdminPassword(),
		StreamKeys:              data.GetStreamKeys(),
		StreamKeyOverridden:     config.TemporaryStreamKey != "",
		WebServerPort:           config.WebServerPort,
		WebServerIP:             config.WebServerIP,
		RTMPServerPort:          data.GetRTMPPortNumber(),
		ChatDisabled:            data.GetChatDisabled(),
		ChatJoinMessagesEnabled: data.GetChatJoinMessagesEnabled(),
		SocketHostOverride:      data.GetWebsocketOverrideHost(),
		VideoServingEndpoint:    data.GetVideoServingEndpoint(),
		ChatEstablishedUserMode: data.GetChatEstbalishedUsersOnlyMode(),
		HideViewerCount:         data.GetHideViewerCount(),
		DisableSearchIndexing:   data.GetDisableSearchIndexing(),
		VideoSettings: videoSettings{
			VideoQualityVariants: videoQualityVariants,
			LatencyLevel:         data.GetStreamLatencyLevel().Level,
		},
		YP: yp{
			Enabled:     data.GetDirectoryEnabled(),
			InstanceURL: data.GetServerURL(),
		},
		S3:                 data.GetS3Config(),
		ExternalActions:    data.GetExternalActions(),
		SupportedCodecs:    transcoder.GetCodecs(ffmpeg),
		VideoCodec:         data.GetVideoCodec(),
		ForbiddenUsernames: usernameBlocklist,
		SuggestedUsernames: usernameSuggestions,
		Federation: federationConfigResponse{
			Enabled:        data.GetFederationEnabled(),
			IsPrivate:      data.GetFederationIsPrivate(),
			Username:       data.GetFederationUsername(),
			GoLiveMessage:  data.GetFederationGoLiveMessage(),
			ShowEngagement: data.GetFederationShowEngagement(),
			BlockedDomains: data.GetBlockedFederatedDomains(),
		},
		Notifications: notificationsConfigResponse{
			Discord: data.GetDiscordConfig(),
			Browser: data.GetBrowserPushConfig(),
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
	Notifications           notificationsConfigResponse `json:"notifications"`
	YP                      yp                          `json:"yp"`
	FFmpegPath              string                      `json:"ffmpegPath"`
	AdminPassword           string                      `json:"adminPassword"`
	SocketHostOverride      string                      `json:"socketHostOverride,omitempty"`
	WebServerIP             string                      `json:"webServerIP"`
	VideoCodec              string                      `json:"videoCodec"`
	VideoServingEndpoint    string                      `json:"videoServingEndpoint"`
	S3                      models.S3                   `json:"s3"`
	Federation              federationConfigResponse    `json:"federation"`
	SupportedCodecs         []string                    `json:"supportedCodecs"`
	ExternalActions         []models.ExternalAction     `json:"externalActions"`
	ForbiddenUsernames      []string                    `json:"forbiddenUsernames"`
	SuggestedUsernames      []string                    `json:"suggestedUsernames"`
	StreamKeys              []models.StreamKey          `json:"streamKeys"`
	VideoSettings           videoSettings               `json:"videoSettings"`
	RTMPServerPort          int                         `json:"rtmpServerPort"`
	WebServerPort           int                         `json:"webServerPort"`
	ChatDisabled            bool                        `json:"chatDisabled"`
	ChatJoinMessagesEnabled bool                        `json:"chatJoinMessagesEnabled"`
	ChatEstablishedUserMode bool                        `json:"chatEstablishedUserMode"`
	DisableSearchIndexing   bool                        `json:"disableSearchIndexing"`
	StreamKeyOverridden     bool                        `json:"streamKeyOverridden"`
	HideViewerCount         bool                        `json:"hideViewerCount"`
}

type videoSettings struct {
	VideoQualityVariants []models.StreamOutputVariant `json:"videoQualityVariants"`
	LatencyLevel         int                          `json:"latencyLevel"`
}

type webConfigResponse struct {
	AppearanceVariables map[string]string     `json:"appearanceVariables"`
	Version             string                `json:"version"`
	WelcomeMessage      string                `json:"welcomeMessage"`
	OfflineMessage      string                `json:"offlineMessage"`
	Logo                string                `json:"logo"`
	Name                string                `json:"name"`
	ExtraPageContent    string                `json:"extraPageContent"`
	StreamTitle         string                `json:"streamTitle"` // What's going on with the current stream
	CustomStyles        string                `json:"customStyles"`
	CustomJavascript    string                `json:"customJavascript"`
	Summary             string                `json:"summary"`
	Tags                []string              `json:"tags"`
	SocialHandles       []models.SocialHandle `json:"socialHandles"`
	NSFW                bool                  `json:"nsfw"`
}

type yp struct {
	InstanceURL  string `json:"instanceUrl"` // The public URL the directory should link to
	YPServiceURL string `json:"-"`           // The base URL to the YP API to register with (optional)
	Enabled      bool   `json:"enabled"`
}

type federationConfigResponse struct {
	Username       string   `json:"username"`
	GoLiveMessage  string   `json:"goLiveMessage"`
	BlockedDomains []string `json:"blockedDomains"`
	Enabled        bool     `json:"enabled"`
	IsPrivate      bool     `json:"isPrivate"`
	ShowEngagement bool     `json:"showEngagement"`
}

type notificationsConfigResponse struct {
	Browser models.BrowserNotificationConfiguration `json:"browser"`
	Discord models.DiscordConfiguration             `json:"discord"`
}
