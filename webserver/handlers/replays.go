package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// clipRateLimiters throttles clip creation per client IP. Clip creation is
// open to any viewer, so the limiter is what keeps it from being spammed.
var (
	clipRateLimitersMu sync.Mutex
	clipRateLimiters   = map[string]*clipRateLimiter{}
)

type clipRateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

const (
	// clipsPerMinute is the sustained clip creation rate allowed per IP.
	clipsPerMinute = 6
	// clipBurst is how many clips may be created back to back.
	clipBurst = 3
	// clipRateLimiterTTL is how long an idle limiter is retained before it is
	// swept, so the map doesn't grow with every visiting IP.
	clipRateLimiterTTL = 30 * time.Minute
)

// allowClipFromIP reports whether the given client may create a clip right now.
func allowClipFromIP(ip string) bool {
	clipRateLimitersMu.Lock()
	defer clipRateLimitersMu.Unlock()

	now := time.Now()
	for key, entry := range clipRateLimiters {
		if now.Sub(entry.lastSeen) > clipRateLimiterTTL {
			delete(clipRateLimiters, key)
		}
	}

	entry, found := clipRateLimiters[ip]
	if !found {
		entry = &clipRateLimiter{limiter: rate.NewLimiter(rate.Limit(clipsPerMinute)/60, clipBurst)}
		clipRateLimiters[ip] = entry
	}
	entry.lastSeen = now

	return entry.limiter.Allow()
}

// clipPermissionDenialReason reports why the given chat identity may not
// create clips under the configured permission level, or "" when it may.
// Every level requires a registered, enabled identity; moderators always
// qualify. The client hides the clip button using the same rules, so viewers
// normally never see these messages.
func clipPermissionDenialReason(user *models.User, permissions string) string {
	if user == nil {
		return "clip creation requires joining chat first"
	}
	if !user.IsEnabled() {
		return "clip creation is not available"
	}
	if user.IsModerator() {
		return ""
	}

	switch permissions {
	case models.ClipPermissionsModerators:
		return "clip creation is limited to moderators"
	case models.ClipPermissionsAuthenticated:
		if !user.Authenticated {
			return "clip creation is limited to authenticated viewers"
		}
	default:
		// Established: any identity old enough.
		if time.Since(user.CreatedAt) < models.MinClipperAccountAge {
			return "your account is too new to create clips"
		}
	}
	return ""
}

// replayFeaturesEnabled reports whether the replay subsystem is available.
func (h *Handlers) replayFeaturesEnabled() bool {
	return h.replays != nil && h.configRepository != nil && h.configRepository.GetReplayFeaturesEnabled()
}

// GetReplay will return playable HLS content for a given stream Id.
// /replay/{streamId} returns the master playlist.
// /replay/{streamId}/{outputConfigId} returns a media playlist.
//
// The route requires admin auth: full-stream replay playback is an admin
// capability, while viewers are served clips.
func (h *Handlers) GetReplay(w http.ResponseWriter, r *http.Request) {
	if !h.replayFeaturesEnabled() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pathComponents := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	switch len(pathComponents) {
	case 2:
		// /replay/{streamId} -> master playlist
		h.getReplayMasterPlaylist(pathComponents[1], w)
	case 3:
		// /replay/{streamId}/{outputConfigId} -> media playlist
		h.getReplayMediaPlaylist(pathComponents[1], pathComponents[2], w)
	default:
		webutils.BadRequestHandler(w, errors.New("invalid replay path"))
	}
}

func (h *Handlers) getReplayMasterPlaylist(streamID string, w http.ResponseWriter) {
	playlistGenerator := h.replays.NewPlaylistGenerator()

	stream, err := playlistGenerator.GetStream(streamID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	playlist, err := playlistGenerator.GenerateMasterPlaylistForStream(streamID)
	if err != nil {
		log.Debugln(err)
		webutils.InternalErrorHandler(w, errors.New("unable to generate replay playlist"))
		return
	}

	if playlist == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// While the stream is still recording, its variant list can still grow.
	writePlaylistResponse(w, playlist.Encode().Bytes(), !stream.InProgress)
}

func (h *Handlers) getReplayMediaPlaylist(streamID, outputConfigID string, w http.ResponseWriter) {
	playlistGenerator := h.replays.NewPlaylistGenerator()

	stream, err := playlistGenerator.GetStream(streamID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	playlist, err := playlistGenerator.GenerateMediaPlaylistForStreamAndConfiguration(streamID, outputConfigID)
	if err != nil {
		log.Debugln(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// An in-progress recording keeps growing, so its playlist must not be
	// cached; a finished one is immutable.
	writePlaylistResponse(w, playlist.Encode().Bytes(), !stream.InProgress)
}

// GetAllClips will return every clip that exists.
func (h *Handlers) GetAllClips(w http.ResponseWriter, r *http.Request) {
	if !h.replayFeaturesEnabled() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	clips, err := h.replays.GetAllClips()
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	webutils.WriteResponse(w, clips)
}

// GetClipDetails returns a single clip, so a clip page can render its title,
// creation time and source stream without fetching the whole list.
func (h *Handlers) GetClipDetails(w http.ResponseWriter, r *http.Request, clipID string) {
	if !h.replayFeaturesEnabled() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	clip, err := h.replays.NewPlaylistGenerator().GetClip(clipID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	webutils.WriteResponse(w, clip)
}

// clipWindowRequest is the caller's explicit media-time clip window.
type clipWindowRequest struct {
	StartSeconds float32
	EndSeconds   float32
}

// resolveClipWindow validates the requested media-time window against recorded
// media and the operator's maximum clip duration.
func resolveClipWindow(request clipWindowRequest, mediaEnd float64, maxDuration float32) (float32, float32, error) {
	start := request.StartSeconds
	end := request.EndSeconds

	if start < 0 {
		return 0, 0, errors.New("start time must not be negative")
	}

	if end <= start {
		return 0, 0, errors.New("end time must be after start time")
	}

	if float64(start) >= mediaEnd {
		return 0, 0, errors.New("start time is after the known end of the stream")
	}

	// Preserve the newest part of an over-long request. Playhead samples and
	// wall-clock countdowns drift independently, so rejecting it would lose a
	// capture at confirmation time.
	end = float32(min(float64(end), mediaEnd))
	if end-start > maxDuration {
		start = end - maxDuration
	}

	return start, end, nil
}

// AddClip creates a clip from the explicit media-time window supplied by the
// viewer after the capture ends.
func (h *Handlers) AddClip(w http.ResponseWriter, r *http.Request) {
	if !h.replayFeaturesEnabled() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !h.configRepository.GetClipsEnabled() {
		webutils.WriteSimpleResponse(w, false, "clip creation is disabled")
		return
	}

	type addClipRequest struct {
		StreamID                 string  `json:"streamId"`
		ClipTitle                string  `json:"clipTitle"`
		RelativeStartTimeSeconds float32 `json:"relativeStartTimeSeconds"`
		RelativeEndTimeSeconds   float32 `json:"relativeEndTimeSeconds"`
	}

	var request addClipRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Errorln(err)
		webutils.WriteSimpleResponse(w, false, "unable to create clip")
		return
	}

	streamID := request.StreamID
	playlistGenerator := h.replays.NewPlaylistGenerator()
	stream, err := playlistGenerator.GetStream(streamID)
	if err != nil {
		webutils.BadRequestHandler(w, errors.New("stream not found"))
		return
	}

	if stream.StartTime.IsZero() {
		webutils.BadRequestHandler(w, errors.New("stream start time not found"))
		return
	}

	mediaEnd, err := h.replays.GetStreamMediaEnd(streamID)
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}
	if mediaEnd <= 0 {
		webutils.BadRequestHandler(w, errors.New("stream has no recorded video to clip yet"))
		return
	}

	startTime, endTime, err := resolveClipWindow(
		clipWindowRequest{
			StartSeconds: request.RelativeStartTimeSeconds,
			EndSeconds:   request.RelativeEndTimeSeconds,
		},
		mediaEnd,
		float32(h.configRepository.GetMaxClipDurationSeconds()),
	)
	if err != nil {
		webutils.BadRequestHandler(w, err)
		return
	}

	var user *models.User
	if token := utils.ChatAccessTokenFromRequest(r); token != "" && h.userRepository != nil {
		user = h.userRepository.GetUserByToken(token)
	}
	if reason := clipPermissionDenialReason(user, h.configRepository.GetClipPermissions()); reason != "" {
		w.WriteHeader(http.StatusForbidden)
		webutils.WriteSimpleResponse(w, false, reason)
		return
	}

	if !allowClipFromIP(utils.GetIPAddressFromRequest(r)) {
		w.WriteHeader(http.StatusTooManyRequests)
		webutils.WriteSimpleResponse(w, false, "too many clips created, please wait before creating another")
		return
	}

	clippedBy := ""
	if user != nil {
		clippedBy = user.DisplayName
	}
	clipID, duration, err := h.replays.AddClipForStream(streamID, request.ClipTitle, clippedBy, startTime, endTime)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go h.replays.GenerateClipThumbnail(clipID)

	webutils.WriteResponse(w, models.ClipCreatedResponse{
		Success:         true,
		Message:         "clip created",
		ID:              clipID,
		DurationSeconds: duration,
	})
}

// GetClip will return playable HLS content for a given clip Id.
// /clip/{clipId} returns the master playlist.
// /clip/{clipId}/{outputConfigId} returns a media playlist.
func (h *Handlers) GetClip(w http.ResponseWriter, r *http.Request) {
	if !h.replayFeaturesEnabled() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pathComponents := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	switch len(pathComponents) {
	case 2:
		// /clip/{clipId} -> master playlist
		h.getClipMasterPlaylist(pathComponents[1], w)
	case 3:
		// /clip/{clipId}/{outputConfigId} -> media playlist
		h.getClipMediaPlaylist(pathComponents[1], pathComponents[2], w)
	default:
		webutils.BadRequestHandler(w, errors.New("invalid clip path"))
	}
}

func (h *Handlers) getClipMasterPlaylist(clipID string, w http.ResponseWriter) {
	playlistGenerator := h.replays.NewPlaylistGenerator()
	playlist, err := playlistGenerator.GenerateMasterPlaylistForClip(clipID)
	if err != nil {
		log.Debugln(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if playlist == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	writePlaylistResponse(w, playlist.Encode().Bytes(), true)
}

func (h *Handlers) getClipMediaPlaylist(clipID, outputConfigID string, w http.ResponseWriter) {
	playlistGenerator := h.replays.NewPlaylistGenerator()
	playlist, err := playlistGenerator.GenerateMediaPlaylistForClipAndConfiguration(clipID, outputConfigID)
	if err != nil {
		log.Debugln(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	writePlaylistResponse(w, playlist.Encode().Bytes(), true)
}

// writePlaylistResponse writes an HLS playlist with the CORS and caching
// headers a player needs. Immutable (finished) playlists may be cached;
// in-progress ones must not be, or playback stalls at a stale segment list.
func writePlaylistResponse(w http.ResponseWriter, playlist []byte, cacheable bool) {
	w.Header().Set("Content-Type", "application/x-mpegURL")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if cacheable {
		w.Header().Set("Cache-Control", "public, max-age=600")
	} else {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	}

	// #nosec G705 -- m3u8 playlist output, not HTML; Content-Type is application/x-mpegURL.
	if _, err := w.Write(playlist); err != nil {
		log.Errorln(err)
	}
}
