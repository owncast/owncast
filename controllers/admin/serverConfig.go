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
			Name:             data.GetServerName(),
			Summary:          data.GetServerSummary(),
			Tags:             data.GetServerMetadataTags(),
			ExtraPageContent: data.GetExtraPageBodyContent(),
			StreamTitle:      data.GetStreamTitle(),
			WelcomeMessage:   data.GetServerWelcomeMessage(),
			Logo:             data.GetLogoPath(),
			SocialHandles:    data.GetSocialHandles(),
			NSFW:             data.GetNSFW(),
			CustomStyles:     data.GetCustomStyles(),
		},
		FFmpegPath:              ffmpeg,
		StreamKey:               data.GetStreamKey(),
		WebServerPort:           config.WebServerPort,
		WebServerIP:             config.WebServerIP,
		RTMPServerPort:          data.GetRTMPPortNumber(),
		ChatDisabled:            data.GetChatDisabled(),
		ChatJoinMessagesEnabled: data.GetChatJoinMessagesEnabled(),
		SocketHostOverride:      data.GetWebsocketOverrideHost(),
		ChatEstablishedUserMode: data.GetChatEstbalishedUsersOnlyMode(),
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
			Twilio:  data.GetTwilioConfig(),
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
	FFmpegPath              string                      `json:"ffmpegPath"`
	StreamKey               string                      `json:"streamKey"`
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
}

type videoSettings struct {
	VideoQualityVariants []models.StreamOutputVariant `json:"videoQualityVariants"`
	LatencyLevel         int                          `json:"latencyLevel"`
}

type webConfigResponse struct {
	Name             string                      `json:"name"`
	Summary          string                      `json:"summary"`
	WelcomeMessage   string                      `json:"welcomeMessage"`
	Logo             string                      `json:"logo"`
	Tags             []string                    `json:"tags"`
	Version          string                      `json:"version"`
	NSFW             bool                        `json:"nsfw"`
	ExtraPageContent string                      `json:"extraPageContent"`
	StreamTitle      string                      `json:"streamTitle"` // What's going on with the current stream
	SocialHandles    []models.SocialHandle       `json:"socialHandles"`
	CustomStyles     string                      `json:"customStyles"`
	Notifications    notificationsConfigResponse `json:"notifications"`
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
	Twilio  models.TwilioConfiguration              `json:"twilio"`
	Discord models.DiscordConfiguration             `json:"discord"`
}
