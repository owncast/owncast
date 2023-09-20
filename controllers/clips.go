package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/replays"
	log "github.com/sirupsen/logrus"
)

// GetAllClips will return all clips that have been previously created.
func GetAllClips(w http.ResponseWriter, r *http.Request) {
	if !config.EnableReplayFeatures {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	clips, err := replays.GetAllClips()
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteResponse(w, clips)
}

// AddClip will create a new clip for a given stream and time window.
func AddClip(w http.ResponseWriter, r *http.Request) {
	if !config.EnableReplayFeatures {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	type addClipRequest struct {
		StreamId                 string  `json:"streamId"`
		ClipTitle                string  `json:"clipTitle"`
		RelativeStartTimeSeconds float32 `json:"relativeStartTimeSeconds"`
		RelativeEndTimeSeconds   float32 `json:"relativeEndTimeSeconds"`
	}

	if r.Method != http.MethodPost {
		BadRequestHandler(w, nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request addClipRequest

	if request.RelativeEndTimeSeconds < request.RelativeStartTimeSeconds {
		BadRequestHandler(w, errors.New("end time must be after start time"))
		return
	}

	if err := decoder.Decode(&request); err != nil {
		log.Errorln(err)
		WriteSimpleResponse(w, false, "unable to create clip")
		return
	}

	streamId := request.StreamId
	clipTitle := request.ClipTitle
	startTime := request.RelativeStartTimeSeconds
	endTime := request.RelativeEndTimeSeconds

	// Some validation
	playlistGenerator := replays.NewPlaylistGenerator()

	stream, err := playlistGenerator.GetStream(streamId)
	if err != nil {
		BadRequestHandler(w, errors.New("stream not found"))
		return
	}

	if stream.StartTime.IsZero() {
		BadRequestHandler(w, errors.New("stream start time not found"))
		return
	}

	// Make sure the proposed clip start time and end time are within
	// the start and end time of the stream.
	finalSegment, err := replays.GetFinalSegmentForStream(streamId)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	if finalSegment.RelativeTimestamp < startTime {
		BadRequestHandler(w, errors.New("start time is after the known end of the stream"))
		return
	}

	clipId, duration, err := replays.AddClipForStream(streamId, clipTitle, "", startTime, endTime)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteSimpleResponse(w, true, "clip "+clipId+" created with duration of "+fmt.Sprint(duration)+" seconds")
}

// GetClip will return playable content for a given clip Id.
func GetClip(w http.ResponseWriter, r *http.Request) {
	if !config.EnableReplayFeatures {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pathComponents := strings.Split(r.URL.Path, "/")
	if len(pathComponents) == 3 {
		// Return the master playlist for the requested stream
		clipId := pathComponents[2]
		getClipMasterPlaylist(clipId, w)
		return
	} else if len(pathComponents) == 4 {
		// Return the media playlist for the requested stream and output config
		clipId := pathComponents[2]
		outputConfigId := pathComponents[3]
		getClipMediaPlaylist(clipId, outputConfigId, w)
		return
	}

	BadRequestHandler(w, nil)
}

// getReplayMasterPlaylist will return a complete replay of a stream
// as a HLS playlist.
func getClipMasterPlaylist(clipId string, w http.ResponseWriter) {
	playlistGenerator := replays.NewPlaylistGenerator()
	playlist, err := playlistGenerator.GenerateMasterPlaylistForClip(clipId)
	if err != nil {
		log.Println(err)
	}

	if playlist == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/x-mpegURL")
	if _, err := w.Write(playlist.Encode().Bytes()); err != nil {
		log.Errorln(err)
		return
	}
}

// getClipMediaPlaylist will return media playlist for a given clip
// and stream output configuration.
func getClipMediaPlaylist(clipId, outputConfigId string, w http.ResponseWriter) {
	playlistGenerator := replays.NewPlaylistGenerator()
	playlist, err := playlistGenerator.GenerateMediaPlaylistForClipAndConfiguration(clipId, outputConfigId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/x-mpegURL")
	if _, err := w.Write(playlist.Encode().Bytes()); err != nil {
		log.Errorln(err)
		return
	}
}
