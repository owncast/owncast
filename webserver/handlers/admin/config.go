package admin

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/handlers/generated"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// ConfigValue is a container object that holds a value, is encoded, and saved to the database.
type ConfigValue struct {
	Value interface{} `json:"value"`
}

// SetTags will handle the web config request to set tags.
func (a *Admin) SetTags(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValues, success := getValuesFromRequest(w, r)
	if !success {
		return
	}

	tagStrings := make([]string, 0)
	for _, tag := range configValues {
		tagStrings = append(tagStrings, strings.TrimLeft(tag.Value.(string), "#"))
	}

	configRepository := a.configRepository
	if err := configRepository.SetServerMetadataTags(tagStrings); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update Fediverse followers about this change.
	if err := a.activitypub.UpdateFollowersWithAccountUpdates(); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "changed")
}

// SetStreamTitle will handle the web config request to set the current stream title.
func (a *Admin) SetStreamTitle(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	value := configValue.Value.(string)
	configRepository := a.configRepository

	if err := configRepository.SetStreamTitle(value); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}
	if value != "" {
		a.sendSystemChatAction(fmt.Sprintf("Stream title changed to **%s**", value), true)
		go a.webhooks.SendStreamStatusEvent(models.StreamTitleUpdated)
	}
	webutils.WriteSimpleResponse(w, true, "changed")
}

// ExternalSetStreamTitle will change the stream title on behalf of an external integration API request.
func (a *Admin) ExternalSetStreamTitle(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	a.SetStreamTitle(w, r)
}

// ExternalGetStatus will return the status of the server.
func (a *Admin) ExternalGetStatus(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	status := a.stream.GetStatus()
	webutils.WriteResponse(w, status)
}

func (a *Admin) sendSystemChatAction(messageText string, ephemeral bool) {
	if err := a.chat.SendSystemAction(messageText, ephemeral); err != nil {
		log.Errorln(err)
	}
}

// SetServerName will handle the web config request to set the server's name.
func (a *Admin) SetServerName(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetServerName(configValue.Value.(string)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update Fediverse followers about this change.
	if err := a.activitypub.UpdateFollowersWithAccountUpdates(); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "changed")
}

// SetServerSummary will handle the web config request to set the about/summary text.
func (a *Admin) SetServerSummary(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetServerSummary(configValue.Value.(string)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update Fediverse followers about this change.
	if err := a.activitypub.UpdateFollowersWithAccountUpdates(); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "changed")
}

// SetCustomOfflineMessage will set a message to display when the server is offline.
func (a *Admin) SetCustomOfflineMessage(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetCustomOfflineMessage(strings.TrimSpace(configValue.Value.(string))); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "changed")
}

// SetServerWelcomeMessage will handle the web config request to set the welcome message text.
func (a *Admin) SetServerWelcomeMessage(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetServerWelcomeMessage(strings.TrimSpace(configValue.Value.(string))); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "changed")
}

// SetExtraPageContent will handle the web config request to set the page markdown content.
func (a *Admin) SetExtraPageContent(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetExtraPageBodyContent(configValue.Value.(string)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "changed")
}

// SetAdminPassword will handle the web config request to set the server admin password.
func (a *Admin) SetAdminPassword(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetAdminPassword(configValue.Value.(string)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "changed")
}

// SetLogo will handle a new logo image file being uploaded.
func (a *Admin) SetLogo(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	value, ok := configValue.Value.(string)
	if !ok {
		webutils.WriteSimpleResponse(w, false, "unable to find image data")
		return
	}
	bytes, extension, err := utils.DecodeBase64Image(value)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	imgPath := filepath.Join("data", "logo"+extension)
	if err := os.WriteFile(imgPath, bytes, 0o600); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	configRepository := a.configRepository

	if err := configRepository.SetLogoPath("logo" + extension); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if err := configRepository.SetLogoUniquenessString(shortid.MustGenerate()); err != nil {
		log.Error("Error saving logo uniqueness string: ", err)
	}

	// Update Fediverse followers about this change.
	if err := a.activitypub.UpdateFollowersWithAccountUpdates(); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "changed")
}

// SetFavicon will handle a new favicon image being set via base64 data.
func (a *Admin) SetFavicon(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	// Limit request body size to prevent large allocations before the decoded
	// size check runs. 200KB decoded = ~267KB base64 + JSON overhead ≈ 300KB.
	r.Body = http.MaxBytesReader(w, r.Body, 300*1024)

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	value, ok := configValue.Value.(string)
	if !ok {
		webutils.WriteSimpleResponse(w, false, "unable to find image data")
		return
	}

	bytes, extension, err := utils.DecodeBase64Image(value)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Only allow PNG and ICO formats for favicons.
	if extension != ".png" && extension != ".ico" {
		webutils.WriteSimpleResponse(w, false, "favicon must be PNG or ICO format")
		return
	}

	// Enforce 200KB size limit.
	const maxFaviconSize = 200 * 1024
	if len(bytes) > maxFaviconSize {
		webutils.WriteSimpleResponse(w, false, "file too large, max 200KB")
		return
	}

	imgPath := filepath.Join("data", "favicon"+extension)
	if err := os.WriteFile(imgPath, bytes, 0o600); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	configRepository := a.configRepository

	if err := configRepository.SetFaviconPath("favicon" + extension); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "favicon updated")
}

// ResetFavicon will reset the favicon to the default by removing the custom one.
func (a *Admin) ResetFavicon(w http.ResponseWriter, r *http.Request) {
	configRepository := a.configRepository

	// Get the current favicon path before clearing it
	currentFavicon := configRepository.GetFaviconPath()

	// Clear the favicon path in the database
	if err := configRepository.SetFaviconPath(""); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Delete the favicon file if it exists
	if currentFavicon != "" {
		faviconPath := filepath.Join("data", currentFavicon)
		if err := os.Remove(faviconPath); err != nil && !os.IsNotExist(err) {
			log.Debugln("error removing favicon file:", err)
		}
	}

	// Also try to remove any favicon files that might exist with different extensions
	for _, ext := range []string{".ico", ".png"} {
		faviconPath := filepath.Join("data", "favicon"+ext)
		if err := os.Remove(faviconPath); err != nil && !os.IsNotExist(err) {
			log.Debugln("error removing favicon file:", err)
		}
	}

	webutils.WriteSimpleResponse(w, true, "favicon reset to default")
}

// SetNSFW will handle the web config request to set the NSFW flag.
func (a *Admin) SetNSFW(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetNSFW(configValue.Value.(bool)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "changed")
}

// SetFfmpegPath will handle the web config request to validate and set an updated copy of ffmpg.
func (a *Admin) SetFfmpegPath(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	path := configValue.Value.(string)
	if err := utils.VerifyFFMpegPath(path); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetFfmpegPath(configValue.Value.(string)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "changed")
}

// SetWebServerPort will handle the web config request to set the server's HTTP port.
func (a *Admin) SetWebServerPort(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository
	if port, ok := configValue.Value.(float64); ok {
		if (port < 1) || (port > 65535) {
			webutils.WriteSimpleResponse(w, false, "Port number must be between 1 and 65535")
			return
		}
		if err := configRepository.SetHTTPPortNumber(port); err != nil {
			webutils.WriteSimpleResponse(w, false, err.Error())
			return
		}

		webutils.WriteSimpleResponse(w, true, "HTTP port set")
		return
	}

	webutils.WriteSimpleResponse(w, false, "Invalid type or value, port must be a number")
}

// SetWebServerIP will handle the web config request to set the server's HTTP listen address.
func (a *Admin) SetWebServerIP(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository
	if input, ok := configValue.Value.(string); ok {
		if ip := net.ParseIP(input); ip != nil {
			if err := configRepository.SetHTTPListenAddress(ip.String()); err != nil {
				webutils.WriteSimpleResponse(w, false, err.Error())
				return
			}

			webutils.WriteSimpleResponse(w, true, "HTTP listen address set")
			return
		}

		webutils.WriteSimpleResponse(w, false, "Invalid IP address")
		return
	}
	webutils.WriteSimpleResponse(w, false, "Invalid type or value, IP address must be a string")
}

// SetRTMPServerBindAddress will handle the web config request to set the inbound RTMP server ip.
func (a *Admin) SetRTMPServerBindAddress(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetRTMPBindAddress(configValue.Value.(string)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "rtmp address set")
}

// SetRTMPServerPort will handle the web config request to set the inbound RTMP port.
func (a *Admin) SetRTMPServerPort(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetRTMPPortNumber(configValue.Value.(float64)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "rtmp port set")
}

// SetServerURL will handle the web config request to set the full server URL.
func (a *Admin) SetServerURL(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	rawValue, ok := configValue.Value.(string)
	if !ok {
		webutils.WriteSimpleResponse(w, false, "could not read server url")
		return
	}

	serverHostString := utils.GetHostnameFromURLString(rawValue)
	if serverHostString == "" {
		webutils.WriteSimpleResponse(w, false, "server url value invalid")
		return
	}

	// Block Private IP URLs
	ipAddr, ipErr := netip.ParseAddr(utils.GetHostnameWithoutPortFromURLString(rawValue))

	if ipErr == nil && ipAddr.IsPrivate() {
		webutils.WriteSimpleResponse(w, false, "Server URL cannot be private")
		return
	}

	configRepository := a.configRepository

	// Trim any trailing slash
	serverURL := strings.TrimRight(rawValue, "/")

	if err := configRepository.SetServerURL(serverURL); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "server url set")
}

// SetSocketHostOverride will set the host override for the websocket.
func (a *Admin) SetSocketHostOverride(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository

	if err := configRepository.SetWebsocketOverrideHost(configValue.Value.(string)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "websocket host override set")
}

// SetDirectoryEnabled will handle the web config request to enable or disable directory registration.
func (a *Admin) SetDirectoryEnabled(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository

	if err := configRepository.SetDirectoryEnabled(configValue.Value.(bool)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}
	webutils.WriteSimpleResponse(w, true, "directory state changed")
}

// SetStreamLatencyLevel will handle the web config request to set the stream latency level.
func (a *Admin) SetStreamLatencyLevel(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository

	if err := configRepository.SetStreamLatencyLevel(configValue.Value.(float64)); err != nil {
		webutils.WriteSimpleResponse(w, false, "error setting stream latency "+err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "set stream latency")
}

// SetS3Configuration will handle the web config request to set the storage configuration.
func (a *Admin) SetS3Configuration(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	type s3ConfigurationRequest struct {
		Value models.S3 `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var newS3Config s3ConfigurationRequest
	if err := decoder.Decode(&newS3Config); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to update s3 config with provided values")
		return
	}

	if newS3Config.Value.Enabled {
		if newS3Config.Value.Endpoint == "" || !utils.IsValidURL((newS3Config.Value.Endpoint)) {
			webutils.WriteSimpleResponse(w, false, "s3 support requires an endpoint")
			return
		}

		if newS3Config.Value.AccessKey == "" || newS3Config.Value.Secret == "" {
			webutils.WriteSimpleResponse(w, false, "s3 support requires an access key and secret")
			return
		}

		if newS3Config.Value.Region == "" {
			webutils.WriteSimpleResponse(w, false, "s3 support requires a region and endpoint")
			return
		}

		if newS3Config.Value.Bucket == "" {
			webutils.WriteSimpleResponse(w, false, "s3 support requires a bucket created for storing public video segments")
			return
		}
	}

	configRepository := a.configRepository
	if err := configRepository.SetS3Config(newS3Config.Value); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}
	webutils.WriteSimpleResponse(w, true, "storage configuration changed")
}

// SetStreamOutputVariants will handle the web config request to set the video output stream variants.
func (a *Admin) SetStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	type streamOutputVariantRequest struct {
		Value []models.StreamOutputVariant `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var videoVariants streamOutputVariantRequest
	if err := decoder.Decode(&videoVariants); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to update video config with provided values "+err.Error())
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetStreamOutputVariants(videoVariants.Value); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to update video config with provided values "+err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "stream output variants updated")
}

// SetSocialHandles will handle the web config request to set the external social profile links.
func (a *Admin) SetSocialHandles(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	type socialHandlesRequest struct {
		Value []models.SocialHandle `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var socialHandles socialHandlesRequest
	if err := decoder.Decode(&socialHandles); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to update social handles with provided values")
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetSocialHandles(socialHandles.Value); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to update social handles with provided values")
		return
	}

	// Update Fediverse followers about this change.
	if err := a.activitypub.UpdateFollowersWithAccountUpdates(); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "social handles updated")
}

// SetChatDisabled will disable chat functionality.
func (a *Admin) SetChatDisabled(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		webutils.WriteSimpleResponse(w, false, "unable to update chat disabled")
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetChatDisabled(configValue.Value.(bool)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "chat disabled status updated")
}

// SetVideoCodec will change the codec used for video encoding.
func (a *Admin) SetVideoCodec(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		webutils.WriteSimpleResponse(w, false, "unable to change video codec")
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetVideoCodec(configValue.Value.(string)); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to update codec")
		return
	}

	webutils.WriteSimpleResponse(w, true, "video codec updated")
}

// SetExternalActions will set the 3rd party actions for the web interface.
func (a *Admin) SetExternalActions(w http.ResponseWriter, r *http.Request) {
	type externalActionsRequest struct {
		Value []models.ExternalAction `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var actions externalActionsRequest
	if err := decoder.Decode(&actions); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to update external actions with provided values")
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetExternalActions(actions.Value); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to update external actions with provided values")
		return
	}

	webutils.WriteSimpleResponse(w, true, "external actions update")
}

// SetCustomStyles will set the CSS string we insert into the page.
func (a *Admin) SetCustomStyles(w http.ResponseWriter, r *http.Request) {
	customStyles, success := getValueFromRequest(w, r)
	if !success {
		webutils.WriteSimpleResponse(w, false, "unable to update custom styles")
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetCustomStyles(customStyles.Value.(string)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "custom styles updated")
}

// SetCustomJavascript will set the Javascript string we insert into the page.
func (a *Admin) SetCustomJavascript(w http.ResponseWriter, r *http.Request) {
	customJavascript, success := getValueFromRequest(w, r)
	if !success {
		webutils.WriteSimpleResponse(w, false, "unable to update custom javascript")
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetCustomJavascript(customJavascript.Value.(string)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "custom styles updated")
}

// SetForbiddenUsernameList will set the list of usernames we do not allow to use.
func (a *Admin) SetForbiddenUsernameList(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request generated.SetForbiddenUsernameListJSONBody

	if err := decoder.Decode(&request); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to update forbidden usernames with provided values")
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetForbiddenUsernameList(*request.Value); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "forbidden username list updated")
}

// SetSuggestedUsernameList will set the list of suggested usernames that newly registered users are assigned if it isn't inferred otherwise (i.e. through a proxy).
func (a *Admin) SetSuggestedUsernameList(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request generated.SetSuggestedUsernameListJSONBody

	if err := decoder.Decode(&request); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to update suggested usernames with provided values")
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetSuggestedUsernamesList(*request.Value); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "suggested username list updated")
}

// SetChatJoinMessagesEnabled will enable or disable the chat join messages.
func (a *Admin) SetChatJoinMessagesEnabled(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		webutils.WriteSimpleResponse(w, false, "unable to update chat join messages enabled")
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetChatJoinMessagesEnabled(configValue.Value.(bool)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "chat join message status updated")
}

// SetHideViewerCount will enable or disable hiding the viewer count.
func (a *Admin) SetHideViewerCount(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		webutils.WriteSimpleResponse(w, false, "unable to update hiding viewer count")
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetHideViewerCount(configValue.Value.(bool)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "hide viewer count setting updated")
}

// SetDisableSearchIndexing will set search indexing support.
func (a *Admin) SetDisableSearchIndexing(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		webutils.WriteSimpleResponse(w, false, "unable to update search indexing")
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetDisableSearchIndexing(configValue.Value.(bool)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "search indexing support updated")
}

// SetVideoServingEndpoint will save the video serving endpoint.
func (a *Admin) SetVideoServingEndpoint(w http.ResponseWriter, r *http.Request) {
	endpoint, success := getValueFromRequest(w, r)
	if !success {
		webutils.WriteSimpleResponse(w, false, "unable to update custom video serving endpoint")
		return
	}

	value, ok := endpoint.Value.(string)
	if !ok {
		webutils.WriteSimpleResponse(w, false, "unable to update custom video serving endpoint")
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetVideoServingEndpoint(value); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "custom video serving endpoint updated")
}

// SetChatSpamProtectionEnabled will enable or disable the chat spam protection.
func (a *Admin) SetChatSpamProtectionEnabled(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetChatSpamProtectionEnabled(configValue.Value.(bool)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}
	webutils.WriteSimpleResponse(w, true, "chat spam protection changed")
}

// SetChatSlurFilterEnabled will enable or disable the chat slur filter.
func (a *Admin) SetChatSlurFilterEnabled(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetChatSlurFilterEnabled(configValue.Value.(bool)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}
	webutils.WriteSimpleResponse(w, true, "chat message slur filter changed")
}

// SetChatRequireAuthentication will enable or disable requiring authentication for chat.
func (a *Admin) SetChatRequireAuthentication(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	configRepository := a.configRepository
	if err := configRepository.SetChatRequireAuthentication(configValue.Value.(bool)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}
	webutils.WriteSimpleResponse(w, true, "chat authentication requirement changed")
}

func requirePOST(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		webutils.WriteSimpleResponse(w, false, r.Method+" not supported")
		return false
	}

	return true
}

func getValueFromRequest(w http.ResponseWriter, r *http.Request) (ConfigValue, bool) {
	decoder := json.NewDecoder(r.Body)
	var configValue ConfigValue
	if err := decoder.Decode(&configValue); err != nil {
		log.Warnln(err)
		webutils.WriteSimpleResponse(w, false, "unable to parse new value")
		return configValue, false
	}

	return configValue, true
}

func getValuesFromRequest(w http.ResponseWriter, r *http.Request) ([]ConfigValue, bool) {
	var values []ConfigValue

	decoder := json.NewDecoder(r.Body)
	var configValue ConfigValue
	if err := decoder.Decode(&configValue); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to parse array of values")
		return values, false
	}

	object := reflect.ValueOf(configValue.Value)

	for i := 0; i < object.Len(); i++ {
		values = append(values, ConfigValue{Value: object.Index(i).Interface()})
	}

	return values, true
}

// SetStreamKeys will set the valid stream keys.
func (a *Admin) SetStreamKeys(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var streamKeys generated.SetStreamKeysJSONRequestBody
	if err := decoder.Decode(&streamKeys); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to update stream keys with provided values")
		return
	}

	if streamKeys.Value == nil || len(*streamKeys.Value) == 0 {
		webutils.WriteSimpleResponse(w, false, "must provide at least one valid stream key")
		return
	}

	for _, streamKey := range *streamKeys.Value {
		if *streamKey.Key == "" {
			webutils.WriteSimpleResponse(w, false, "stream key cannot be empty")
			return
		}
	}

	configRepository := a.configRepository
	if err := configRepository.SetStreamKeys(*streamKeys.Value); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "changed")
}
