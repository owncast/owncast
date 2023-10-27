package replays

import (
	"context"
	"strings"

	"github.com/grafov/m3u8"
	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/models"
	"github.com/pkg/errors"
)

// GenerateMasterPlaylistForClip returns a master playlist for a given clip Id.
// It includes references to the media playlists for each output configuration.
func (p *PlaylistGenerator) GenerateMasterPlaylistForClip(clipId string) (*m3u8.MasterPlaylist, error) {
	clip, err := p.datastore.GetQueries().GetClip(context.Background(), clipId)
	if err != nil {
		return nil, errors.Wrap(err, "unable to fetch requested clip")
	}

	streamId := clip.StreamID
	configs, err := p.GetConfigurationsForStream(streamId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get configurations for stream")
	}

	// Create the master playlist that will hold the different media playlists.
	masterPlaylist := p.createNewMasterPlaylist()

	// Create the media playlists for each output configuration.
	for _, config := range configs {
		// Verify the validity of the configuration.
		if err := config.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid output configuration")
		}

		mediaPlaylist, err := p.GenerateMediaPlaylistForClipAndConfiguration(clipId, config.ID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create clip media playlist")
		}

		// Append the media playlist to the master playlist.
		params := p.getMediaPlaylistParamsForConfig(config)

		// Add the media playlist to the master playlist.
		publicPlaylistPath := strings.Join([]string{"/clip", clipId, config.ID}, "/")
		masterPlaylist.Append(publicPlaylistPath, mediaPlaylist, params)
	}

	// Return the final master playlist that contains all the media playlists.
	return masterPlaylist, nil
}

// GenerateMediaPlaylistForClipAndConfiguration returns a media playlist for a
// given clip Id and output configuration.
func (p *PlaylistGenerator) GenerateMediaPlaylistForClipAndConfiguration(clipId, outputConfigurationId string) (*m3u8.MediaPlaylist, error) {
	clip, err := p.GetClip(clipId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stream")
	}

	config, err := p.GetOutputConfig(outputConfigurationId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get output configuration")
	}

	clipStartSeconds := clip.RelativeStartTime
	clipEndSeconds := clip.RelativeEndTime

	// Fetch all the segments for this configuration.
	segments, err := p.GetAllSegmentsForOutputConfigurationAndWindow(outputConfigurationId, clipStartSeconds, clipEndSeconds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all clip segments for output configuration")
	}

	// Create the media playlist for this configuration and add the segments.
	mediaPlaylist, err := p.createMediaPlaylistForConfigurationAndSegments(config, clip.Timestamp, false, segments)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create clip media playlist")
	}

	return mediaPlaylist, nil
}

// GetClip returns a clip by its ID.
func (p *PlaylistGenerator) GetClip(clipId string) (*Clip, error) {
	clip, err := p.datastore.GetQueries().GetClip(context.Background(), clipId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get clip")
	}

	if clip.ClipID == "" {
		return nil, errors.Wrap(err, "failed to get clip")
	}

	if !clip.RelativeEndTime.Valid {
		return nil, errors.Wrap(err, "failed to get clip")
	}

	timestamp, err := models.FlexibleDateParse(clip.ClipTimestamp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse clip timestamp")
	}

	c := Clip{
		ID:                clip.ClipID,
		StreamId:          clip.StreamID,
		ClipTitle:         clip.ClipTitle.String,
		RelativeStartTime: float32(clip.RelativeStartTime.Float64),
		RelativeEndTime:   float32(clip.RelativeEndTime.Float64),
		Timestamp:         timestamp,
	}

	return &c, nil
}

// GetAllSegmentsForOutputConfigurationAndWindow returns all the segments for a
// given output config and time window.
func (p *PlaylistGenerator) GetAllSegmentsForOutputConfigurationAndWindow(configId string, startSeconds, endSeconds float32) ([]HLSSegment, error) {
	segmentRows, err := p.datastore.GetQueries().GetSegmentsForOutputIdAndWindow(context.Background(), db.GetSegmentsForOutputIdAndWindowParams{
		OutputConfigurationID: configId,
		StartSeconds:          startSeconds,
		EndSeconds:            endSeconds,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get clip segments for output config")
	}

	segments := []HLSSegment{}
	for _, row := range segmentRows {
		segment := HLSSegment{
			ID:                    row.ID,
			StreamID:              row.StreamID,
			OutputConfigurationID: row.OutputConfigurationID,
			Timestamp:             row.Timestamp.Time,
			Path:                  row.Path,
		}
		segments = append(segments, segment)
	}

	return segments, nil
}
