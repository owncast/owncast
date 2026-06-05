package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/router/middleware"
	webutils "github.com/owncast/owncast/webserver/utils"
)

type webConfigResponse struct {
	AppearanceVariables map[string]string `json:"appearanceVariables"`
	Name                string            `json:"name"`
	// CustomStyles is the admin-configured CSS plus the concatenated
	// content of every loaded plugin's manifest.styles entries (each
	// preceded by a `/* plugin: <slug> ... */` delimiter for
	// devtools attribution). The viewer renders this as one inline
	// <style> block so plugins can theme the page without each plugin
	// needing its own <link> tag.
	CustomStyles       string                  `json:"customStyles"`
	StreamTitle        string                  `json:"streamTitle,omitempty"` // What's going on with the current stream
	OfflineMessage     string                  `json:"offlineMessage"`
	Logo               string                  `json:"logo"`
	Version            string                  `json:"version"`
	SocketHostOverride string                  `json:"socketHostOverride,omitempty"`
	ExtraPageContent   string                  `json:"extraPageContent"`
	Summary            string                  `json:"summary"`
	Tags               []string                `json:"tags"`
	SocialHandles      []models.SocialHandle   `json:"socialHandles"`
	ExternalActions    []models.ExternalAction `json:"externalActions"`
	// PluginTabs is the list of viewer-page tabs contributed by
	// loaded plugins via manifest.tabs. The viewer page renders one
	// tab per entry alongside the built-in tabs (Followers, About).
	PluginTabs                 []models.PluginTab           `json:"pluginTabs"`
	Notifications              notificationsConfigResponse  `json:"notifications"`
	Federation                 federationConfigResponse     `json:"federation"`
	MaxSocketPayloadSize       int                          `json:"maxSocketPayloadSize"`
	HideViewerCount            bool                         `json:"hideViewerCount"`
	ChatDisabled               bool                         `json:"chatDisabled"`
	ChatSpamProtectionDisabled bool                         `json:"chatSpamProtectionDisabled"`
	ChatRequireAuthentication  bool                         `json:"chatRequireAuthentication"`
	NSFW                       bool                         `json:"nsfw"`
	Authentication             authenticationConfigResponse `json:"authentication"`
}

type federationConfigResponse struct {
	Account       string `json:"account,omitempty"`
	FollowerCount int    `json:"followerCount,omitempty"`
	Enabled       bool   `json:"enabled"`
}

type browserNotificationsConfigResponse struct {
	PublicKey string `json:"publicKey,omitempty"`
	Enabled   bool   `json:"enabled"`
}

type notificationsConfigResponse struct {
	Browser browserNotificationsConfigResponse `json:"browser"`
}

type authenticationConfigResponse struct {
	IndieAuthEnabled bool `json:"indieAuthEnabled"`
}

// GetWebConfig gets the status of the server.
func (h *Handlers) GetWebConfig(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)
	middleware.DisableCache(w)
	w.Header().Set("Content-Type", "application/json")

	configuration := h.getConfigResponse(r)

	if err := json.NewEncoder(w).Encode(configuration); err != nil {
		webutils.BadRequestHandler(w, err)
	}
}

func (h *Handlers) getConfigResponse(r *http.Request) webConfigResponse {
	configRepository := h.configRepository
	pageContent := utils.RenderPageContentMarkdown(configRepository.GetExtraPageBodyContent())
	pageContent = prependPluginPageContent(pageContent, r, h.pluginPageContent)
	offlineMessage := utils.RenderSimpleMarkdown(configRepository.GetCustomOfflineMessage())
	socialHandles := configRepository.GetSocialHandles()
	for i, handle := range socialHandles {
		platform := models.GetSocialHandle(handle.Platform)
		if platform != nil {
			handle.Icon = platform.Icon
			socialHandles[i] = handle
		}
	}

	serverSummary := configRepository.GetServerSummary()

	var federationResponse federationConfigResponse
	federationEnabled := configRepository.GetFederationEnabled()

	followerCount, _ := h.activitypub.GetFollowerCount()
	if federationEnabled {
		serverURLString := configRepository.GetServerURL()
		serverURL, _ := url.Parse(serverURLString)
		account := fmt.Sprintf("%s@%s", configRepository.GetDefaultFederationUsername(), serverURL.Host)
		federationResponse = federationConfigResponse{
			Enabled:       federationEnabled,
			FollowerCount: int(followerCount),
			Account:       account,
		}
	}

	browserPushEnabled := configRepository.GetBrowserPushConfig().Enabled
	browserPushPublicKey, err := configRepository.GetBrowserPushPublicKey()
	if err != nil {
		log.Errorln("unable to fetch browser push notifications public key", err)
		browserPushEnabled = false
	}

	notificationsResponse := notificationsConfigResponse{
		Browser: browserNotificationsConfigResponse{
			Enabled:   browserPushEnabled,
			PublicKey: browserPushPublicKey,
		},
	}

	authenticationResponse := authenticationConfigResponse{
		IndieAuthEnabled: configRepository.GetServerURL() != "",
	}

	return webConfigResponse{
		Name:                       configRepository.GetServerName(),
		Summary:                    serverSummary,
		OfflineMessage:             offlineMessage,
		Logo:                       "/logo",
		Tags:                       configRepository.GetServerMetadataTags(),
		Version:                    config.GetReleaseString(),
		NSFW:                       configRepository.GetNSFW(),
		SocketHostOverride:         configRepository.GetWebsocketOverrideHost(),
		ExtraPageContent:           pageContent,
		StreamTitle:                configRepository.GetStreamTitle(),
		SocialHandles:              socialHandles,
		ChatDisabled:               configRepository.GetChatDisabled(),
		ChatSpamProtectionDisabled: configRepository.GetChatSpamProtectionEnabled(),
		ChatRequireAuthentication:  configRepository.GetChatRequireAuthentication(),
		ExternalActions:            mergePluginActions(configRepository.GetExternalActions(), h.pluginActions),
		PluginTabs:                 pluginTabsOrEmpty(r, h.pluginTabs),
		CustomStyles:               mergePluginCSS(configRepository.GetCustomStyles(), h.pluginCSSContent),
		MaxSocketPayloadSize:       config.MaxSocketPayloadSize,
		Federation:                 federationResponse,
		Notifications:              notificationsResponse,
		Authentication:             authenticationResponse,
		AppearanceVariables:        configRepository.GetCustomColorVariableValues(),
		HideViewerCount:            configRepository.GetHideViewerCount(),
	}
}

// pluginTabsOrEmpty returns the host's plugin-tab list, or an empty
// (non-nil) slice when the getter is unset (no plugin host) or
// returns nothing. The empty-slice contract keeps the JSON wire
// shape stable: `pluginTabs: []` rather than `null`, so the viewer
// doesn't need a defensive nil-check before iterating.
func pluginTabsOrEmpty(r *http.Request, getter func(*http.Request) []models.PluginTab) []models.PluginTab {
	if getter == nil {
		return []models.PluginTab{}
	}
	if t := getter(r); len(t) > 0 {
		return t
	}
	return []models.PluginTab{}
}

// prependPluginPageContent puts plugin-contributed HTML in front of
// the admin's rendered extraPageContent. Plugin HTML lands at the top
// of the extra-content block so plugins can announce themselves above
// the admin's prose. A nil getter or empty contribution leaves the
// admin's value untouched. A newline separates the two sources so
// the trailing markup of one can't run into the next.
func prependPluginPageContent(admin string, r *http.Request, pluginHTML func(*http.Request) []byte) string {
	if pluginHTML == nil {
		return admin
	}
	bytes := pluginHTML(r)
	if len(bytes) == 0 {
		return admin
	}
	prefix := string(bytes)
	if admin == "" {
		return prefix
	}
	if prefix[len(prefix)-1] != '\n' {
		prefix += "\n"
	}
	return prefix + admin
}

// mergePluginCSS appends plugin-contributed CSS bytes to the admin's
// custom CSS so the viewer renders one inline <style> block covering
// both. A nil getter (no plugin host) or empty contribution leaves
// the admin's value untouched. A newline separates the two sources
// so a stylesheet that doesn't terminate in one can't trail into the
// plugin block.
func mergePluginCSS(admin string, pluginCSS func() []byte) string {
	if pluginCSS == nil {
		return admin
	}
	bytes := pluginCSS()
	if len(bytes) == 0 {
		return admin
	}
	if admin == "" {
		return string(bytes)
	}
	if admin[len(admin)-1] != '\n' {
		admin += "\n"
	}
	return admin + string(bytes)
}

// mergePluginActions appends plugin-contributed action buttons to the
// admin-configured list so the viewer sees both in one externalActions
// array. Admin entries stay first (and so render first in the UI) — a
// plugin can extend the action set but can't reorder or replace what
// the admin defined. A nil getter (no plugin host) makes this a no-op.
func mergePluginActions(
	configured []models.ExternalAction,
	pluginActions func() []models.ExternalAction,
) []models.ExternalAction {
	if pluginActions == nil {
		return configured
	}
	contributed := pluginActions()
	if len(contributed) == 0 {
		return configured
	}
	out := make([]models.ExternalAction, 0, len(configured)+len(contributed))
	out = append(out, configured...)
	out = append(out, contributed...)
	return out
}

// GetAllSocialPlatforms will return a list of all social platform types.
func GetAllSocialPlatforms(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)
	w.Header().Set("Content-Type", "application/json")

	platforms := models.GetAllSocialHandles()
	if err := json.NewEncoder(w).Encode(platforms); err != nil {
		webutils.InternalErrorHandler(w, err)
	}
}
