package admin

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
)

// ConfigValue is a container object that holds a value, is encoded, and saved to the database.
type ConfigValue struct {
	Value interface{} `json:"value"`
}

// SetTags will handle the web config request to set tags.
func (c *Controller) SetTags(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValues, success := c.getValuesFromRequest(w, r)
	if !success {
		return
	}

	tagStrings := make([]string, 0)
	for _, tag := range configValues {
		tagStrings = append(tagStrings, strings.TrimLeft(tag.Value.(string), "#"))
	}

	if err := c.Data.SetServerMetadataTags(tagStrings); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update Fediverse followers about this change.
	if err := c.Outbox.UpdateFollowersWithAccountUpdates(c.ActivityPub.Persistence); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "changed")
}

// SetStreamTitle will handle the web config request to set the current stream title.
func (c *Controller) SetStreamTitle(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	value := configValue.Value.(string)

	if err := c.Data.SetStreamTitle(value); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}
	if value != "" {
		c.sendSystemChatAction(fmt.Sprintf("Stream title changed to **%s**", value), true)
	}
	c.WriteSimpleResponse(w, true, "changed")
}

// ExternalSetStreamTitle will change the stream title on behalf of an external integration API request.
func (c *Controller) ExternalSetStreamTitle(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	c.SetStreamTitle(w, r)
}

func (c *Controller) sendSystemChatAction(messageText string, ephemeral bool) {
	if err := c.Core.Chat.SendSystemAction(messageText, ephemeral); err != nil {
		log.Errorln(err)
	}
}

// SetServerName will handle the web config request to set the server's name.
func (c *Controller) SetServerName(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetServerName(configValue.Value.(string)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update Fediverse followers about this change.
	if err := c.Outbox.UpdateFollowersWithAccountUpdates(c.ActivityPub.Persistence); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "changed")
}

// SetServerSummary will handle the web config request to set the about/summary text.
func (c *Controller) SetServerSummary(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetServerSummary(configValue.Value.(string)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update Fediverse followers about this change.
	if err := c.Outbox.UpdateFollowersWithAccountUpdates(c.ActivityPub.Persistence); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "changed")
}

// SetCustomOfflineMessage will set a message to display when the server is offline.
func (c *Controller) SetCustomOfflineMessage(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetCustomOfflineMessage(strings.TrimSpace(configValue.Value.(string))); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "changed")
}

// SetServerWelcomeMessage will handle the web config request to set the welcome message text.
func (c *Controller) SetServerWelcomeMessage(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetServerWelcomeMessage(strings.TrimSpace(configValue.Value.(string))); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "changed")
}

// SetExtraPageContent will handle the web config request to set the page markdown content.
func (c *Controller) SetExtraPageContent(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetExtraPageBodyContent(configValue.Value.(string)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "changed")
}

// SetAdminPassword will handle the web config request to set the server admin password.
func (c *Controller) SetAdminPassword(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetAdminPassword(configValue.Value.(string)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "changed")
}

// SetLogo will handle a new logo image file being uploaded.
func (c *Controller) SetLogo(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	value, ok := configValue.Value.(string)
	if !ok {
		c.WriteSimpleResponse(w, false, "unable to find image data")
		return
	}
	bytes, extension, err := utils.DecodeBase64Image(value)
	if err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	imgPath := filepath.Join("data", "logo"+extension)
	if err := os.WriteFile(imgPath, bytes, 0o600); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if err := c.Data.SetLogoPath("logo" + extension); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if err := c.Data.SetLogoUniquenessString(shortid.MustGenerate()); err != nil {
		log.Error("Error saving logo uniqueness string: ", err)
	}

	// Update Fediverse followers about this change.
	if err := c.Outbox.UpdateFollowersWithAccountUpdates(c.ActivityPub.Persistence); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "changed")
}

// SetNSFW will handle the web config request to set the NSFW flag.
func (c *Controller) SetNSFW(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetNSFW(configValue.Value.(bool)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "changed")
}

// SetFfmpegPath will handle the web config request to validate and set an updated copy of ffmpg.
func (c *Controller) SetFfmpegPath(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	path := configValue.Value.(string)
	if err := utils.VerifyFFMpegPath(path); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if err := c.Data.SetFfmpegPath(configValue.Value.(string)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "changed")
}

// SetWebServerPort will handle the web config request to set the server's HTTP port.
func (c *Controller) SetWebServerPort(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if port, ok := configValue.Value.(float64); ok {
		if (port < 1) || (port > 65535) {
			c.WriteSimpleResponse(w, false, "Port number must be between 1 and 65535")
			return
		}
		if err := c.Data.SetHTTPPortNumber(port); err != nil {
			c.WriteSimpleResponse(w, false, err.Error())
			return
		}

		c.WriteSimpleResponse(w, true, "HTTP port set")
		return
	}

	c.WriteSimpleResponse(w, false, "Invalid type or value, port must be a number")
}

// SetWebServerIP will handle the web config request to set the server's HTTP listen address.
func (c *Controller) SetWebServerIP(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if input, ok := configValue.Value.(string); ok {
		if ip := net.ParseIP(input); ip != nil {
			if err := c.Data.SetHTTPListenAddress(ip.String()); err != nil {
				c.WriteSimpleResponse(w, false, err.Error())
				return
			}

			c.WriteSimpleResponse(w, true, "HTTP listen address set")
			return
		}

		c.WriteSimpleResponse(w, false, "Invalid IP address")
		return
	}
	c.WriteSimpleResponse(w, false, "Invalid type or value, IP address must be a string")
}

// SetRTMPServerPort will handle the web config request to set the inbound RTMP port.
func (c *Controller) SetRTMPServerPort(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetRTMPPortNumber(configValue.Value.(float64)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "rtmp port set")
}

// SetServerURL will handle the web config request to set the full server URL.
func (c *Controller) SetServerURL(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	rawValue, ok := configValue.Value.(string)
	if !ok {
		c.WriteSimpleResponse(w, false, "could not read server url")
		return
	}

	serverHostString := utils.GetHostnameFromURLString(rawValue)
	if serverHostString == "" {
		c.WriteSimpleResponse(w, false, "server url value invalid")
		return
	}

	// Trim any trailing slash
	serverURL := strings.TrimRight(rawValue, "/")

	if err := c.Data.SetServerURL(serverURL); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "server url set")
}

// SetSocketHostOverride will set the host override for the websocket.
func (c *Controller) SetSocketHostOverride(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetWebsocketOverrideHost(configValue.Value.(string)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "websocket host override set")
}

// SetDirectoryEnabled will handle the web config request to enable or disable directory registration.
func (c *Controller) SetDirectoryEnabled(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetDirectoryEnabled(configValue.Value.(bool)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}
	c.WriteSimpleResponse(w, true, "directory state changed")
}

// SetStreamLatencyLevel will handle the web config request to set the stream latency level.
func (c *Controller) SetStreamLatencyLevel(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := c.Data.SetStreamLatencyLevel(configValue.Value.(float64)); err != nil {
		c.WriteSimpleResponse(w, false, "error setting stream latency "+err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "set stream latency")
}

// SetS3Configuration will handle the web config request to set the storage configuration.
func (c *Controller) SetS3Configuration(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	type s3ConfigurationRequest struct {
		Value models.S3 `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var newS3Config s3ConfigurationRequest
	if err := decoder.Decode(&newS3Config); err != nil {
		c.WriteSimpleResponse(w, false, "unable to update s3 config with provided values")
		return
	}

	if newS3Config.Value.Enabled {
		if newS3Config.Value.Endpoint == "" || !utils.IsValidURL((newS3Config.Value.Endpoint)) {
			c.WriteSimpleResponse(w, false, "s3 support requires an endpoint")
			return
		}

		if newS3Config.Value.AccessKey == "" || newS3Config.Value.Secret == "" {
			c.WriteSimpleResponse(w, false, "s3 support requires an access key and secret")
			return
		}

		if newS3Config.Value.Region == "" {
			c.WriteSimpleResponse(w, false, "s3 support requires a region and endpoint")
			return
		}

		if newS3Config.Value.Bucket == "" {
			c.WriteSimpleResponse(w, false, "s3 support requires a bucket created for storing public video segments")
			return
		}
	}

	if err := c.Data.SetS3Config(newS3Config.Value); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}
	c.WriteSimpleResponse(w, true, "storage configuration changed")
}

// SetStreamOutputVariants will handle the web config request to set the video output stream variants.
func (c *Controller) SetStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	type streamOutputVariantRequest struct {
		Value []models.StreamOutputVariant `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var videoVariants streamOutputVariantRequest
	if err := decoder.Decode(&videoVariants); err != nil {
		c.WriteSimpleResponse(w, false, "unable to update video config with provided values "+err.Error())
		return
	}

	if err := c.Data.SetStreamOutputVariants(videoVariants.Value); err != nil {
		c.WriteSimpleResponse(w, false, "unable to update video config with provided values "+err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "stream output variants updated")
}

// SetSocialHandles will handle the web config request to set the external social profile links.
func (c *Controller) SetSocialHandles(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	type socialHandlesRequest struct {
		Value []models.SocialHandle `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var socialHandles socialHandlesRequest
	if err := decoder.Decode(&socialHandles); err != nil {
		c.WriteSimpleResponse(w, false, "unable to update social handles with provided values")
		return
	}

	if err := c.Data.SetSocialHandles(socialHandles.Value); err != nil {
		c.WriteSimpleResponse(w, false, "unable to update social handles with provided values")
		return
	}

	// Update Fediverse followers about this change.
	if err := c.Outbox.UpdateFollowersWithAccountUpdates(c.ActivityPub.Persistence); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "social handles updated")
}

// SetChatDisabled will disable chat functionality.
func (c *Controller) SetChatDisabled(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		c.WriteSimpleResponse(w, false, "unable to update chat disabled")
		return
	}

	if err := c.Data.SetChatDisabled(configValue.Value.(bool)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "chat disabled status updated")
}

// SetVideoCodec will change the codec used for video encoding.
func (c *Controller) SetVideoCodec(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		c.WriteSimpleResponse(w, false, "unable to change video codec")
		return
	}

	if err := c.Data.SetVideoCodec(configValue.Value.(string)); err != nil {
		c.WriteSimpleResponse(w, false, "unable to update codec")
		return
	}

	c.WriteSimpleResponse(w, true, "video codec updated")
}

// SetExternalActions will set the 3rd party actions for the web interface.
func (c *Controller) SetExternalActions(w http.ResponseWriter, r *http.Request) {
	type externalActionsRequest struct {
		Value []models.ExternalAction `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var actions externalActionsRequest
	if err := decoder.Decode(&actions); err != nil {
		c.WriteSimpleResponse(w, false, "unable to update external actions with provided values")
		return
	}

	if err := c.Data.SetExternalActions(actions.Value); err != nil {
		c.WriteSimpleResponse(w, false, "unable to update external actions with provided values")
		return
	}

	c.WriteSimpleResponse(w, true, "external actions update")
}

// SetCustomStyles will set the CSS string we insert into the page.
func (c *Controller) SetCustomStyles(w http.ResponseWriter, r *http.Request) {
	customStyles, success := c.getValueFromRequest(w, r)
	if !success {
		c.WriteSimpleResponse(w, false, "unable to update custom styles")
		return
	}

	if err := c.Data.SetCustomStyles(customStyles.Value.(string)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "custom styles updated")
}

// SetCustomJavascript will set the Javascript string we insert into the page.
func (c *Controller) SetCustomJavascript(w http.ResponseWriter, r *http.Request) {
	customJavascript, success := c.getValueFromRequest(w, r)
	if !success {
		c.WriteSimpleResponse(w, false, "unable to update custom javascript")
		return
	}

	if err := c.Data.SetCustomJavascript(customJavascript.Value.(string)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "custom styles updated")
}

// SetForbiddenUsernameList will set the list of usernames we do not allow to use.
func (c *Controller) SetForbiddenUsernameList(w http.ResponseWriter, r *http.Request) {
	type forbiddenUsernameListRequest struct {
		Value []string `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var request forbiddenUsernameListRequest
	if err := decoder.Decode(&request); err != nil {
		c.WriteSimpleResponse(w, false, "unable to update forbidden usernames with provided values")
		return
	}

	if err := c.Data.SetForbiddenUsernameList(request.Value); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "forbidden username list updated")
}

// SetSuggestedUsernameList will set the list of suggested usernames that newly registered users are assigned if it isn't inferred otherwise (i.e. through a proxy).
func (c *Controller) SetSuggestedUsernameList(w http.ResponseWriter, r *http.Request) {
	type suggestedUsernameListRequest struct {
		Value []string `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var request suggestedUsernameListRequest

	if err := decoder.Decode(&request); err != nil {
		c.WriteSimpleResponse(w, false, "unable to update suggested usernames with provided values")
		return
	}

	if err := c.Data.SetSuggestedUsernamesList(request.Value); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "suggested username list updated")
}

// SetChatJoinMessagesEnabled will enable or disable the chat join messages.
func (c *Controller) SetChatJoinMessagesEnabled(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		c.WriteSimpleResponse(w, false, "unable to update chat join messages enabled")
		return
	}

	if err := c.Data.SetChatJoinMessagesEnabled(configValue.Value.(bool)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "chat join message status updated")
}

// SetHideViewerCount will enable or disable hiding the viewer count.
func (c *Controller) SetHideViewerCount(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	configValue, success := c.getValueFromRequest(w, r)
	if !success {
		c.WriteSimpleResponse(w, false, "unable to update hiding viewer count")
		return
	}

	if err := c.Data.SetHideViewerCount(configValue.Value.(bool)); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "hide viewer count setting updated")
}

func (c *Controller) requirePOST(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != controllers.POST {
		c.WriteSimpleResponse(w, false, r.Method+" not supported")
		return false
	}

	return true
}

func (c *Controller) getValueFromRequest(w http.ResponseWriter, r *http.Request) (ConfigValue, bool) {
	decoder := json.NewDecoder(r.Body)
	var configValue ConfigValue
	if err := decoder.Decode(&configValue); err != nil {
		log.Warnln(err)
		c.WriteSimpleResponse(w, false, "unable to parse new value")
		return configValue, false
	}

	return configValue, true
}

func (c *Controller) getValuesFromRequest(w http.ResponseWriter, r *http.Request) ([]ConfigValue, bool) {
	var values []ConfigValue

	decoder := json.NewDecoder(r.Body)
	var configValue ConfigValue
	if err := decoder.Decode(&configValue); err != nil {
		c.WriteSimpleResponse(w, false, "unable to parse array of values")
		return values, false
	}

	object := reflect.ValueOf(configValue.Value)

	for i := 0; i < object.Len(); i++ {
		values = append(values, ConfigValue{Value: object.Index(i).Interface()})
	}

	return values, true
}

// SetStreamKeys will set the valid stream keys.
func (c *Controller) SetStreamKeys(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	type streamKeysRequest struct {
		Value []models.StreamKey `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var streamKeys streamKeysRequest
	if err := decoder.Decode(&streamKeys); err != nil {
		c.WriteSimpleResponse(w, false, "unable to update stream keys with provided values")
		return
	}

	if err := c.Data.SetStreamKeys(streamKeys.Value); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "changed")
}
