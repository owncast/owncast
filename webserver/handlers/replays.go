package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
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

// replayFeaturesEnabled reports whether the replay subsystem is available.
func (h *Handlers) replayFeaturesEnabled() bool {
	return h.replays != nil && h.configRepository != nil && h.configRepository.GetReplayFeaturesEnabled()
}

// GetReplay will return playable HLS content for a given stream Id.
// /replay/{streamId} returns the master playlist.
// /replay/{streamId}/{outputConfigId} returns a media playlist.
//
// Full-stream replay playback is admin-only for now; viewers only get clips.
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

// GetAllClips will return all clips that have been previously created.
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

// clipWindowRequest is the caller's requested clip window: either a trailing
// duration, or an explicit media-time range.
type clipWindowRequest struct {
	DurationSeconds float32
	StartSeconds    float32
	EndSeconds      float32
}

// resolveClipWindow turns a clip request into a concrete media-time window,
// clamped to what has actually been recorded and to the operator's maximum
// clip length.
//
// A trailing duration is resolved against the newest recorded media, which is
// what a viewer clipping a live stream sends: the server owns this arithmetic
// so the client never has to know where the media timeline currently sits.
func resolveClipWindow(request clipWindowRequest, mediaEnd float64, maxDuration float32) (float32, float32, error) {
	if request.DurationSeconds > 0 {
		duration := min(request.DurationSeconds, maxDuration)
		end := float32(mediaEnd)
		return max(end-duration, 0), end, nil
	}

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

	// Clamp to what actually exists, then enforce the operator's limit.
	end = float32(min(float64(end), mediaEnd))
	if end-start > maxDuration {
		return 0, 0, fmt.Errorf("clips may be at most %.0f seconds long", maxDuration)
	}

	return start, end, nil
}

// AddClip will create a new clip for a given stream.
//
// A clip is either an explicit media-time window
// (relativeStartTimeSeconds + relativeEndTimeSeconds) or, for a viewer
// clipping what just happened on a live stream, a trailing window of
// durationSeconds ending at the newest recorded media.
func (h *Handlers) AddClip(w http.ResponseWriter, r *http.Request) {
	if !h.replayFeaturesEnabled() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !h.configRepository.GetClipsEnabled() {
		webutils.WriteSimpleResponse(w, false, "clip creation is disabled")
		return
	}

	if !allowClipFromIP(utils.GetIPAddressFromRequest(r)) {
		w.WriteHeader(http.StatusTooManyRequests)
		webutils.WriteSimpleResponse(w, false, "too many clips created, please wait before creating another")
		return
	}

	type addClipRequest struct {
		StreamID  string `json:"streamId"`
		ClipTitle string `json:"clipTitle"`
		// DurationSeconds requests a trailing clip of this length ending at
		// the newest recorded media for the stream.
		DurationSeconds float32 `json:"durationSeconds"`
		// An explicit media-time window, used instead of DurationSeconds.
		RelativeStartTimeSeconds float32 `json:"relativeStartTimeSeconds"`
		RelativeEndTimeSeconds   float32 `json:"relativeEndTimeSeconds"`
	}

	decoder := json.NewDecoder(r.Body)
	var request addClipRequest
	if err := decoder.Decode(&request); err != nil {
		log.Errorln(err)
		webutils.WriteSimpleResponse(w, false, "unable to create clip")
		return
	}

	streamID := request.StreamID

	// Verify the stream exists and has a known media timeline.
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

	startTime, endTime, err := resolveClipWindow(clipWindowRequest{
		DurationSeconds: request.DurationSeconds,
		StartSeconds:    request.RelativeStartTimeSeconds,
		EndSeconds:      request.RelativeEndTimeSeconds,
	}, mediaEnd, float32(h.configRepository.GetMaxClipDurationSeconds()))
	if err != nil {
		webutils.BadRequestHandler(w, err)
		return
	}

	// Attribute the clip when the requester has a chat identity. Clipping
	// does not require one.
	clippedBy := ""
	if token := utils.ChatAccessTokenFromRequest(r); token != "" && h.userRepository != nil {
		if user := h.userRepository.GetUserByToken(token); user != nil {
			clippedBy = user.DisplayName
		}
	}

	clipID, duration, err := h.replays.AddClipForStream(streamID, request.ClipTitle, clippedBy, startTime, endTime)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Generate a poster image for the clip so it has a thumbnail in listings
	// and in link previews. Best effort: a clip without one still plays.
	h.replays.GenerateClipThumbnail(clipID)

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
