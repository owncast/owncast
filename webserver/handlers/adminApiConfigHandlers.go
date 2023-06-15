package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/owncast/owncast/activitypub/outbox"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/webhooks"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/requests"
	"github.com/owncast/owncast/webserver/responses"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

// SetTags will handle the web config request to set tags.
func (h *Handlers) SetTags(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValues, success := requests.GetValuesFromRequest(w, r)
	if !success {
		return
	}

	tagStrings := make([]string, 0)
	for _, tag := range configValues {
		tagStrings = append(tagStrings, strings.TrimLeft(tag.Value.(string), "#"))
	}

	if err := data.SetServerMetadataTags(tagStrings); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update Fediverse followers about this change.
	if err := outbox.UpdateFollowersWithAccountUpdates(); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "changed")
}

// SetStreamTitle will handle the web config request to set the current stream title.
func (h *Handlers) SetStreamTitle(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	value := configValue.Value.(string)

	if err := data.SetStreamTitle(value); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}
	if value != "" {
		sendSystemChatAction(fmt.Sprintf("Stream title changed to **%s**", value), true)
		webhookManager := webhooks.GetWebhooks()
		go webhookManager.SendStreamStatusEvent(models.StreamTitleUpdated)
	}
	responses.WriteSimpleResponse(w, true, "changed")
}

// ExternalSetStreamTitle will change the stream title on behalf of an external integration API request.
func (h *Handlers) ExternalSetStreamTitle(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	h.SetStreamTitle(w, r)
}

func sendSystemChatAction(messageText string, ephemeral bool) {
	if err := chat.SendSystemAction(messageText, ephemeral); err != nil {
		log.Errorln(err)
	}
}

// SetServerName will handle the web config request to set the server's name.
func (h *Handlers) SetServerName(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetServerName(configValue.Value.(string)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update Fediverse followers about this change.
	if err := outbox.UpdateFollowersWithAccountUpdates(); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "changed")
}

// SetServerSummary will handle the web config request to set the about/summary text.
func (h *Handlers) SetServerSummary(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetServerSummary(configValue.Value.(string)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update Fediverse followers about this change.
	if err := outbox.UpdateFollowersWithAccountUpdates(); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "changed")
}

// SetCustomOfflineMessage will set a message to display when the server is offline.
func (h *Handlers) SetCustomOfflineMessage(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetCustomOfflineMessage(strings.TrimSpace(configValue.Value.(string))); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "changed")
}

// SetServerWelcomeMessage will handle the web config request to set the welcome message text.
func (h *Handlers) SetServerWelcomeMessage(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetServerWelcomeMessage(strings.TrimSpace(configValue.Value.(string))); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "changed")
}

// SetExtraPageContent will handle the web config request to set the page markdown content.
func (h *Handlers) SetExtraPageContent(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetExtraPageBodyContent(configValue.Value.(string)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "changed")
}

// SetAdminPassword will handle the web config request to set the server admin password.
func (h *Handlers) SetAdminPassword(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetAdminPassword(configValue.Value.(string)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "changed")
}

// SetLogo will handle a new logo image file being uploaded.
func (h *Handlers) SetLogo(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	value, ok := configValue.Value.(string)
	if !ok {
		responses.WriteSimpleResponse(w, false, "unable to find image data")
		return
	}
	bytes, extension, err := utils.DecodeBase64Image(value)
	if err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	imgPath := filepath.Join("data", "logo"+extension)
	if err := os.WriteFile(imgPath, bytes, 0o600); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if err := data.SetLogoPath("logo" + extension); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if err := data.SetLogoUniquenessString(shortid.MustGenerate()); err != nil {
		log.Error("Error saving logo uniqueness string: ", err)
	}

	// Update Fediverse followers about this change.
	if err := outbox.UpdateFollowersWithAccountUpdates(); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "changed")
}

// SetNSFW will handle the web config request to set the NSFW flag.
func (h *Handlers) SetNSFW(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetNSFW(configValue.Value.(bool)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "changed")
}

// SetFfmpegPath will handle the web config request to validate and set an updated copy of ffmpg.
func (h *Handlers) SetFfmpegPath(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	path := configValue.Value.(string)
	if err := utils.VerifyFFMpegPath(path); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if err := data.SetFfmpegPath(configValue.Value.(string)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "changed")
}

// SetWebServerPort will handle the web config request to set the server's HTTP port.
func (h *Handlers) SetWebServerPort(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if port, ok := configValue.Value.(float64); ok {
		if (port < 1) || (port > 65535) {
			responses.WriteSimpleResponse(w, false, "Port number must be between 1 and 65535")
			return
		}
		if err := data.SetHTTPPortNumber(port); err != nil {
			responses.WriteSimpleResponse(w, false, err.Error())
			return
		}

		responses.WriteSimpleResponse(w, true, "HTTP port set")
		return
	}

	responses.WriteSimpleResponse(w, false, "Invalid type or value, port must be a number")
}

// SetWebServerIP will handle the web config request to set the server's HTTP listen address.
func (h *Handlers) SetWebServerIP(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if input, ok := configValue.Value.(string); ok {
		if ip := net.ParseIP(input); ip != nil {
			if err := data.SetHTTPListenAddress(ip.String()); err != nil {
				responses.WriteSimpleResponse(w, false, err.Error())
				return
			}

			responses.WriteSimpleResponse(w, true, "HTTP listen address set")
			return
		}

		responses.WriteSimpleResponse(w, false, "Invalid IP address")
		return
	}
	responses.WriteSimpleResponse(w, false, "Invalid type or value, IP address must be a string")
}

// SetRTMPServerPort will handle the web config request to set the inbound RTMP port.
func (h *Handlers) SetRTMPServerPort(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetRTMPPortNumber(configValue.Value.(float64)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "rtmp port set")
}

// SetServerURL will handle the web config request to set the full server URL.
func (h *Handlers) SetServerURL(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	rawValue, ok := configValue.Value.(string)
	if !ok {
		responses.WriteSimpleResponse(w, false, "could not read server url")
		return
	}

	serverHostString := utils.GetHostnameFromURLString(rawValue)
	if serverHostString == "" {
		responses.WriteSimpleResponse(w, false, "server url value invalid")
		return
	}

	// Trim any trailing slash
	serverURL := strings.TrimRight(rawValue, "/")

	if err := data.SetServerURL(serverURL); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "server url set")
}

// SetSocketHostOverride will set the host override for the websocket.
func (h *Handlers) SetSocketHostOverride(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetWebsocketOverrideHost(configValue.Value.(string)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "websocket host override set")
}

// SetDirectoryEnabled will handle the web config request to enable or disable directory registration.
func (h *Handlers) SetDirectoryEnabled(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetDirectoryEnabled(configValue.Value.(bool)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}
	responses.WriteSimpleResponse(w, true, "directory state changed")
}

// SetStreamLatencyLevel will handle the web config request to set the stream latency level.
func (h *Handlers) SetStreamLatencyLevel(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetStreamLatencyLevel(configValue.Value.(float64)); err != nil {
		responses.WriteSimpleResponse(w, false, "error setting stream latency "+err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "set stream latency")
}

// SetS3Configuration will handle the web config request to set the storage configuration.
func (h *Handlers) SetS3Configuration(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	type s3ConfigurationRequest struct {
		Value models.S3 `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var newS3Config s3ConfigurationRequest
	if err := decoder.Decode(&newS3Config); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update s3 config with provided values")
		return
	}

	if newS3Config.Value.Enabled {
		if newS3Config.Value.Endpoint == "" || !utils.IsValidURL((newS3Config.Value.Endpoint)) {
			responses.WriteSimpleResponse(w, false, "s3 support requires an endpoint")
			return
		}

		if newS3Config.Value.AccessKey == "" || newS3Config.Value.Secret == "" {
			responses.WriteSimpleResponse(w, false, "s3 support requires an access key and secret")
			return
		}

		if newS3Config.Value.Region == "" {
			responses.WriteSimpleResponse(w, false, "s3 support requires a region and endpoint")
			return
		}

		if newS3Config.Value.Bucket == "" {
			responses.WriteSimpleResponse(w, false, "s3 support requires a bucket created for storing public video segments")
			return
		}
	}

	if err := data.SetS3Config(newS3Config.Value); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}
	responses.WriteSimpleResponse(w, true, "storage configuration changed")
}

// SetStreamOutputVariants will handle the web config request to set the video output stream variants.
func (h *Handlers) SetStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	type streamOutputVariantRequest struct {
		Value []models.StreamOutputVariant `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var videoVariants streamOutputVariantRequest
	if err := decoder.Decode(&videoVariants); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update video config with provided values "+err.Error())
		return
	}

	if err := data.SetStreamOutputVariants(videoVariants.Value); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update video config with provided values "+err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "stream output variants updated")
}

// SetSocialHandles will handle the web config request to set the external social profile links.
func (h *Handlers) SetSocialHandles(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	type socialHandlesRequest struct {
		Value []models.SocialHandle `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var socialHandles socialHandlesRequest
	if err := decoder.Decode(&socialHandles); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update social handles with provided values")
		return
	}

	if err := data.SetSocialHandles(socialHandles.Value); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update social handles with provided values")
		return
	}

	// Update Fediverse followers about this change.
	if err := outbox.UpdateFollowersWithAccountUpdates(); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "social handles updated")
}

// SetChatDisabled will disable chat functionality.
func (h *Handlers) SetChatDisabled(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		responses.WriteSimpleResponse(w, false, "unable to update chat disabled")
		return
	}

	if err := data.SetChatDisabled(configValue.Value.(bool)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "chat disabled status updated")
}

// SetVideoCodec will change the codec used for video encoding.
func (h *Handlers) SetVideoCodec(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		responses.WriteSimpleResponse(w, false, "unable to change video codec")
		return
	}

	if err := data.SetVideoCodec(configValue.Value.(string)); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update codec")
		return
	}

	responses.WriteSimpleResponse(w, true, "video codec updated")
}

// SetExternalActions will set the 3rd party actions for the web interface.
func (h *Handlers) SetExternalActions(w http.ResponseWriter, r *http.Request) {
	type externalActionsRequest struct {
		Value []models.ExternalAction `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var actions externalActionsRequest
	if err := decoder.Decode(&actions); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update external actions with provided values")
		return
	}

	if err := data.SetExternalActions(actions.Value); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update external actions with provided values")
		return
	}

	responses.WriteSimpleResponse(w, true, "external actions update")
}

// SetCustomStyles will set the CSS string we insert into the page.
func (h *Handlers) SetCustomStyles(w http.ResponseWriter, r *http.Request) {
	customStyles, success := requests.GetValueFromRequest(w, r)
	if !success {
		responses.WriteSimpleResponse(w, false, "unable to update custom styles")
		return
	}

	if err := data.SetCustomStyles(customStyles.Value.(string)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "custom styles updated")
}

// SetCustomJavascript will set the Javascript string we insert into the page.
func (h *Handlers) SetCustomJavascript(w http.ResponseWriter, r *http.Request) {
	customJavascript, success := requests.GetValueFromRequest(w, r)
	if !success {
		responses.WriteSimpleResponse(w, false, "unable to update custom javascript")
		return
	}

	if err := data.SetCustomJavascript(customJavascript.Value.(string)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "custom styles updated")
}

// SetForbiddenUsernameList will set the list of usernames we do not allow to use.
func (h *Handlers) SetForbiddenUsernameList(w http.ResponseWriter, r *http.Request) {
	type forbiddenUsernameListRequest struct {
		Value []string `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var request forbiddenUsernameListRequest
	if err := decoder.Decode(&request); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update forbidden usernames with provided values")
		return
	}

	if err := data.SetForbiddenUsernameList(request.Value); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "forbidden username list updated")
}

// SetSuggestedUsernameList will set the list of suggested usernames that newly registered users are assigned if it isn't inferred otherwise (i.e. through a proxy).
func (h *Handlers) SetSuggestedUsernameList(w http.ResponseWriter, r *http.Request) {
	type suggestedUsernameListRequest struct {
		Value []string `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var request suggestedUsernameListRequest

	if err := decoder.Decode(&request); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update suggested usernames with provided values")
		return
	}

	if err := data.SetSuggestedUsernamesList(request.Value); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "suggested username list updated")
}

// SetChatJoinMessagesEnabled will enable or disable the chat join messages.
func (h *Handlers) SetChatJoinMessagesEnabled(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		responses.WriteSimpleResponse(w, false, "unable to update chat join messages enabled")
		return
	}

	if err := data.SetChatJoinMessagesEnabled(configValue.Value.(bool)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "chat join message status updated")
}

// SetHideViewerCount will enable or disable hiding the viewer count.
func (h *Handlers) SetHideViewerCount(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		responses.WriteSimpleResponse(w, false, "unable to update hiding viewer count")
		return
	}

	if err := data.SetHideViewerCount(configValue.Value.(bool)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "hide viewer count setting updated")
}

// SetDisableSearchIndexing will set search indexing support.
func (h *Handlers) SetDisableSearchIndexing(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	configValue, success := requests.GetValueFromRequest(w, r)
	if !success {
		responses.WriteSimpleResponse(w, false, "unable to update search indexing")
		return
	}

	if err := data.SetDisableSearchIndexing(configValue.Value.(bool)); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "search indexing support updated")
}

// SetVideoServingEndpoint will save the video serving endpoint.
func (h *Handlers) SetVideoServingEndpoint(w http.ResponseWriter, r *http.Request) {
	endpoint, success := requests.GetValueFromRequest(w, r)
	if !success {
		responses.WriteSimpleResponse(w, false, "unable to update custom video serving endpoint")
		return
	}

	value, ok := endpoint.Value.(string)
	if !ok {
		responses.WriteSimpleResponse(w, false, "unable to update custom video serving endpoint")
		return
	}

	if err := data.SetVideoServingEndpoint(value); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "custom video serving endpoint updated")
}

// SetStreamKeys will set the valid stream keys.
func (h *Handlers) SetStreamKeys(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	type streamKeysRequest struct {
		Value []models.StreamKey `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var streamKeys streamKeysRequest
	if err := decoder.Decode(&streamKeys); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update stream keys with provided values")
		return
	}

	if len(streamKeys.Value) == 0 {
		responses.WriteSimpleResponse(w, false, "must provide at least one valid stream key")
		return
	}

	for _, streamKey := range streamKeys.Value {
		if streamKey.Key == "" {
			responses.WriteSimpleResponse(w, false, "stream key cannot be empty")
			return
		}
	}

	if err := data.SetStreamKeys(streamKeys.Value); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "changed")
}
