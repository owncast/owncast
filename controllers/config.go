package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/app/middleware"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
)

type webConfigResponse struct {
	Name                 string                       `json:"name"`
	Summary              string                       `json:"summary"`
	OfflineMessage       string                       `json:"offlineMessage"`
	Logo                 string                       `json:"logo"`
	Tags                 []string                     `json:"tags"`
	Version              string                       `json:"version"`
	NSFW                 bool                         `json:"nsfw"`
	SocketHostOverride   string                       `json:"socketHostOverride,omitempty"`
	ExtraPageContent     string                       `json:"extraPageContent"`
	StreamTitle          string                       `json:"streamTitle,omitempty"` // What's going on with the current stream
	SocialHandles        []models.SocialHandle        `json:"socialHandles"`
	ChatDisabled         bool                         `json:"chatDisabled"`
	ExternalActions      []models.ExternalAction      `json:"externalActions"`
	CustomStyles         string                       `json:"customStyles"`
	AppearanceVariables  map[string]string            `json:"appearanceVariables"`
	MaxSocketPayloadSize int                          `json:"maxSocketPayloadSize"`
	Federation           federationConfigResponse     `json:"federation"`
	Notifications        notificationsConfigResponse  `json:"Notifications"`
	Authentication       authenticationConfigResponse `json:"authentication"`
}

type federationConfigResponse struct {
	Enabled       bool   `json:"enabled"`
	Account       string `json:"account,omitempty"`
	FollowerCount int    `json:"followerCount,omitempty"`
}

type browserNotificationsConfigResponse struct {
	Enabled   bool   `json:"enabled"`
	PublicKey string `json:"publicKey,omitempty"`
}

type notificationsConfigResponse struct {
	Browser browserNotificationsConfigResponse `json:"browser"`
}

type authenticationConfigResponse struct {
	IndieAuthEnabled bool `json:"indieAuthEnabled"`
}

// GetWebConfig gets the status of the server.
func (s *Service) GetWebConfig(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)
	middleware.DisableCache(w)
	w.Header().Set("Content-Type", "application/json")

	configuration := s.getConfigResponse()

	if err := json.NewEncoder(w).Encode(configuration); err != nil {
		s.BadRequestHandler(w, err)
	}
}

func (s *Service) getConfigResponse() webConfigResponse {
	pageContent := utils.RenderPageContentMarkdown(s.Data.GetExtraPageBodyContent())
	socialHandles := s.Data.GetSocialHandles()
	for i, handle := range socialHandles {
		platform := models.GetSocialHandle(handle.Platform)
		if platform != nil {
			handle.Icon = platform.Icon
			socialHandles[i] = handle
		}
	}

	serverSummary := s.Data.GetServerSummary()

	var federationResponse federationConfigResponse
	federationEnabled := s.Data.GetFederationEnabled()

	followerCount, _ := s.ActivityPub.GetFollowerCount()
	if federationEnabled {
		serverURLString := s.Data.GetServerURL()
		serverURL, _ := url.Parse(serverURLString)
		account := fmt.Sprintf("%s@%s", s.Data.GetDefaultFederationUsername(), serverURL.Host)
		federationResponse = federationConfigResponse{
			Enabled:       federationEnabled,
			FollowerCount: int(followerCount),
			Account:       account,
		}
	}

	browserPushEnabled := s.Data.GetBrowserPushConfig().Enabled
	browserPushPublicKey, err := s.Data.GetBrowserPushPublicKey()
	if err != nil {
		log.Errorln("unable to fetch browser push Notifications public key", err)
		browserPushEnabled = false
	}

	notificationsResponse := notificationsConfigResponse{
		Browser: browserNotificationsConfigResponse{
			Enabled:   browserPushEnabled,
			PublicKey: browserPushPublicKey,
		},
	}

	authenticationResponse := authenticationConfigResponse{
		IndieAuthEnabled: s.Data.GetServerURL() != "",
	}

	return webConfigResponse{
		Name:                 s.Data.GetServerName(),
		Summary:              serverSummary,
		OfflineMessage:       s.Data.GetCustomOfflineMessage(),
		Logo:                 "/logo",
		Tags:                 s.Data.GetServerMetadataTags(),
		Version:              config.GetReleaseString(),
		NSFW:                 s.Data.GetNSFW(),
		SocketHostOverride:   s.Data.GetWebsocketOverrideHost(),
		ExtraPageContent:     pageContent,
		StreamTitle:          s.Data.GetStreamTitle(),
		SocialHandles:        socialHandles,
		ChatDisabled:         s.Data.GetChatDisabled(),
		ExternalActions:      s.Data.GetExternalActions(),
		CustomStyles:         s.Data.GetCustomStyles(),
		MaxSocketPayloadSize: config.MaxSocketPayloadSize,
		Federation:           federationResponse,
		Notifications:        notificationsResponse,
		Authentication:       authenticationResponse,
		AppearanceVariables:  s.Data.GetCustomColorVariableValues(),
	}
}

// GetAllSocialPlatforms will return a list of all social platform types.
func (s *Service) GetAllSocialPlatforms(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)
	w.Header().Set("Content-Type", "application/json")

	platforms := models.GetAllSocialHandles()
	if err := json.NewEncoder(w).Encode(platforms); err != nil {
		s.InternalErrorHandler(w, err)
	}
}
