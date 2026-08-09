package admin

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/models"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// replaysAvailable reports whether the replay subsystem is enabled, writing a
// 404 when it isn't so the admin UI can hide the section.
func (a *Admin) replaysAvailable(w http.ResponseWriter) bool {
	if a.replays == nil || !a.configRepository.GetReplayFeaturesEnabled() {
		w.WriteHeader(http.StatusNotFound)
		return false
	}

	return true
}

// GetAdminReplays returns every replay with the disk usage and clip count the
// admin needs to manage it.
func (a *Admin) GetAdminReplays(w http.ResponseWriter, r *http.Request) {
	if !a.replaysAvailable(w) {
		return
	}

	replays, err := a.replays.GetStreams()
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	webutils.WriteResponse(w, replays)
}

// GetAdminClips returns every clip for administration.
func (a *Admin) GetAdminClips(w http.ResponseWriter, r *http.Request) {
	if !a.replaysAvailable(w) {
		return
	}

	clips, err := a.replays.GetAllClips()
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	webutils.WriteResponse(w, clips)
}

// idFromRequest reads the target object's id from a JSON request body.
func idFromRequest(w http.ResponseWriter, r *http.Request) (string, bool) {
	type idRequest struct {
		ID string `json:"id"`
	}

	var request idRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to read request")
		return "", false
	}

	if request.ID == "" {
		webutils.WriteSimpleResponse(w, false, "an id is required")
		return "", false
	}

	return request.ID, true
}

// DeleteReplay deletes a replay along with its clips and recorded video.
func (a *Admin) DeleteReplay(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	if !a.replaysAvailable(w) {
		return
	}

	id, ok := idFromRequest(w, r)
	if !ok {
		return
	}

	if err := a.replays.DeleteReplay(id); err != nil {
		log.Errorln(err)
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "replay deleted")
}

// DeleteClip deletes a single clip. The replay it came from is untouched.
func (a *Admin) DeleteClip(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	if !a.replaysAvailable(w) {
		return
	}

	id, ok := idFromRequest(w, r)
	if !ok {
		return
	}

	if err := a.replays.DeleteClip(id); err != nil {
		log.Errorln(err)
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "clip deleted")
}

// SetReplayFeaturesEnabled enables or disables replay recording and clips.
//
// Recorded video is kept on disk for as long as a replay or clip references
// it, so disk use grows with every recorded stream until replays are deleted.
func (a *Admin) SetReplayFeaturesEnabled(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	enabled, ok := configValue.Value.(bool)
	if !ok {
		webutils.WriteSimpleResponse(w, false, "a boolean value is required")
		return
	}

	if err := a.configRepository.SetReplayFeaturesEnabled(enabled); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "replay features saved")
}

// SetClipsEnabled enables or disables viewer clip creation.
func (a *Admin) SetClipsEnabled(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	enabled, ok := configValue.Value.(bool)
	if !ok {
		webutils.WriteSimpleResponse(w, false, "a boolean value is required")
		return
	}

	if err := a.configRepository.SetClipsEnabled(enabled); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "clip creation saved")
}

// SetMaxClipDuration sets the longest clip a viewer may create.
func (a *Admin) SetMaxClipDuration(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	seconds, ok := configValue.Value.(float64)
	if !ok {
		webutils.WriteSimpleResponse(w, false, "a numeric value is required")
		return
	}

	if seconds < 1 || seconds > models.MaxAllowedClipDurationSeconds {
		webutils.WriteSimpleResponse(w, false, "clip duration must be between 1 and 3600 seconds")
		return
	}

	if err := a.configRepository.SetMaxClipDurationSeconds(int(seconds)); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "maximum clip duration saved")
}

// SetClipPermissions sets who may create clips.
func (a *Admin) SetClipPermissions(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	permissions, ok := configValue.Value.(string)
	if !ok || !models.IsValidClipPermissions(permissions) {
		webutils.WriteSimpleResponse(w, false, "clip permissions must be one of: moderators, authenticated, established")
		return
	}

	if err := a.configRepository.SetClipPermissions(permissions); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "clip permissions saved")
}
