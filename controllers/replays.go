package controllers

import (
	"net/http"
	"strings"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/replays"
	log "github.com/sirupsen/logrus"
)

// GetReplays will return a list of all available replays.
func GetReplays(w http.ResponseWriter, r *http.Request) {
	if !config.EnableReplayFeatures {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	streams, err := replays.GetStreams()
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteResponse(w, streams)
}

// GetReplay will return a playable content for a given stream Id.
func GetReplay(w http.ResponseWriter, r *http.Request) {
	pathComponents := strings.Split(r.URL.Path, "/")
	if len(pathComponents) == 3 {
		// Return the master playlist for the requested stream
		streamId := pathComponents[2]
		getReplayMasterPlaylist(streamId, w)
		return
	} else if len(pathComponents) == 4 {
		// Return the media playlist for the requested stream and output config
		streamId := pathComponents[2]
		outputConfigId := pathComponents[3]
		getReplayMediaPlaylist(streamId, outputConfigId, w)
		return
	}

	BadRequestHandler(w, nil)
}

// getReplayMasterPlaylist will return a complete replay of a stream as a HLS playlist.
// /api/replay/{streamId}.
func getReplayMasterPlaylist(streamId string, w http.ResponseWriter) {
	playlistGenerator := replays.NewPlaylistGenerator()
	playlist, err := playlistGenerator.GenerateMasterPlaylistForStream(streamId)
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

// getReplayMediaPlaylist will return a media playlist for a given stream.
// /api/replay/{streamId}/{outputConfigId}.
func getReplayMediaPlaylist(streamId, outputConfigId string, w http.ResponseWriter) {
	playlistGenerator := replays.NewPlaylistGenerator()
	playlist, err := playlistGenerator.GenerateMediaPlaylistForStreamAndConfiguration(streamId, outputConfigId)
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
