package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// GetReplays will return a list of all available stream replays.
func (h *Handlers) GetReplays(w http.ResponseWriter, r *http.Request) {
	if !config.EnableReplayFeatures || h.replays == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	streams, err := h.replays.GetStreams()
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	webutils.WriteResponse(w, streams)
}

// GetReplay will return playable HLS content for a given stream Id.
// /replay/{streamId} returns the master playlist.
// /replay/{streamId}/{outputConfigId} returns a media playlist.
func (h *Handlers) GetReplay(w http.ResponseWriter, r *http.Request) {
	if !config.EnableReplayFeatures || h.replays == nil {
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
	playlist, err := playlistGenerator.GenerateMasterPlaylistForStream(streamID)
	if err != nil {
		log.Debugln(err)
	}

	if playlist == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/x-mpegURL")
	// #nosec G705 -- m3u8 playlist output, not HTML; Content-Type is application/x-mpegURL.
	if _, err := w.Write(playlist.Encode().Bytes()); err != nil {
		log.Errorln(err)
		return
	}
}

func (h *Handlers) getReplayMediaPlaylist(streamID, outputConfigID string, w http.ResponseWriter) {
	playlistGenerator := h.replays.NewPlaylistGenerator()
	playlist, err := playlistGenerator.GenerateMediaPlaylistForStreamAndConfiguration(streamID, outputConfigID)
	if err != nil {
		log.Debugln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/x-mpegURL")
	// #nosec G705 -- m3u8 playlist output, not HTML; Content-Type is application/x-mpegURL.
	if _, err := w.Write(playlist.Encode().Bytes()); err != nil {
		log.Errorln(err)
		return
	}
}

// GetAllClips will return all clips that have been previously created.
func (h *Handlers) GetAllClips(w http.ResponseWriter, r *http.Request) {
	if !config.EnableReplayFeatures || h.replays == nil {
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

// AddClip will create a new clip for a given stream and time window.
func (h *Handlers) AddClip(w http.ResponseWriter, r *http.Request) {
	if !config.EnableReplayFeatures || h.replays == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		webutils.BadRequestHandler(w, errors.New("a POST request is required to create a clip"))
		return
	}

	type addClipRequest struct {
		StreamID                 string  `json:"streamId"`
		ClipTitle                string  `json:"clipTitle"`
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

	if request.RelativeEndTimeSeconds < request.RelativeStartTimeSeconds {
		webutils.BadRequestHandler(w, errors.New("end time must be after start time"))
		return
	}

	streamID := request.StreamID
	startTime := request.RelativeStartTimeSeconds
	endTime := request.RelativeEndTimeSeconds

	// Some validation. Verify the stream exists and the requested window is
	// within the known bounds of the stream.
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

	finalSegment, err := h.replays.GetFinalSegmentForStream(streamID)
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	if float32(finalSegment.RelativeTimestamp) < startTime {
		webutils.BadRequestHandler(w, errors.New("start time is after the known end of the stream"))
		return
	}

	clipID, duration, err := h.replays.AddClipForStream(streamID, request.ClipTitle, "", startTime, endTime)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	webutils.WriteSimpleResponse(w, true, "clip "+clipID+" created with duration of "+fmt.Sprint(duration)+" seconds")
}

// GetClip will return playable HLS content for a given clip Id.
// /clip/{clipId} returns the master playlist.
// /clip/{clipId}/{outputConfigId} returns a media playlist.
func (h *Handlers) GetClip(w http.ResponseWriter, r *http.Request) {
	if !config.EnableReplayFeatures || h.replays == nil {
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
	}

	if playlist == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/x-mpegURL")
	// #nosec G705 -- m3u8 playlist output, not HTML; Content-Type is application/x-mpegURL.
	if _, err := w.Write(playlist.Encode().Bytes()); err != nil {
		log.Errorln(err)
		return
	}
}

func (h *Handlers) getClipMediaPlaylist(clipID, outputConfigID string, w http.ResponseWriter) {
	playlistGenerator := h.replays.NewPlaylistGenerator()
	playlist, err := playlistGenerator.GenerateMediaPlaylistForClipAndConfiguration(clipID, outputConfigID)
	if err != nil {
		log.Debugln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/x-mpegURL")
	// #nosec G705 -- m3u8 playlist output, not HTML; Content-Type is application/x-mpegURL.
	if _, err := w.Write(playlist.Encode().Bytes()); err != nil {
		log.Errorln(err)
		return
	}
}
