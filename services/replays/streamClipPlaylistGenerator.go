package replays

import (
	"context"
	"math"
	"strings"

	"github.com/grafov/m3u8"
	"github.com/pkg/errors"

	"github.com/owncast/owncast/db"
)

// GenerateMasterPlaylistForClip returns a master playlist for a given clip ID.
// It includes references to the media playlists for each output configuration.
func (p *PlaylistGenerator) GenerateMasterPlaylistForClip(clipID string) (*m3u8.MasterPlaylist, error) {
	clip, err := p.datastore.GetQueries().GetClip(context.Background(), clipID)
	if err != nil {
		return nil, errors.Wrap(err, "unable to fetch requested clip")
	}

	streamID := clip.StreamID
	configs, err := p.GetConfigurationsForStream(streamID)
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

		mediaPlaylist, err := p.GenerateMediaPlaylistForClipAndConfiguration(clipID, config.ID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create clip media playlist")
		}

		// Append the media playlist to the master playlist.
		params := p.getMediaPlaylistParamsForConfig(config)

		// Add the media playlist to the master playlist.
		publicPlaylistPath := strings.Join([]string{"/clip", clipID, config.ID}, "/")
		masterPlaylist.Append(publicPlaylistPath, mediaPlaylist, params)
	}

	// Return the final master playlist that contains all the media playlists.
	return masterPlaylist, nil
}

// GenerateMediaPlaylistForClipAndConfiguration returns a media playlist for a
// given clip ID and output configuration.
func (p *PlaylistGenerator) GenerateMediaPlaylistForClipAndConfiguration(clipID, outputConfigurationID string) (*m3u8.MediaPlaylist, error) {
	clip, err := p.GetClip(clipID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get clip")
	}

	config, err := p.GetOutputConfig(outputConfigurationID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get output configuration")
	}

	clipStartSeconds := clip.RelativeStartTime
	clipEndSeconds := clip.RelativeEndTime

	// Fetch all the segments for this configuration within the clip window.
	segments, err := p.GetAllSegmentsForOutputConfigurationAndWindow(outputConfigurationID, clipStartSeconds, clipEndSeconds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all clip segments for output configuration")
	}

	// Create the media playlist for this configuration and add the segments.
	mediaPlaylist, err := p.createMediaPlaylistForConfigurationAndSegments(config, false, segments)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create clip media playlist")
	}

	return mediaPlaylist, nil
}

// GetClip returns a clip by its ID.
func (p *PlaylistGenerator) GetClip(clipID string) (*Clip, error) {
	clip, err := p.datastore.GetQueries().GetClip(context.Background(), clipID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get clip")
	}

	if clip.ClipID == "" {
		return nil, errors.New("clip not found")
	}

	if !clip.RelativeEndTime.Valid {
		return nil, errors.New("clip has no end time")
	}

	c := Clip{
		ID:                clip.ClipID,
		StreamID:          clip.StreamID,
		ClipTitle:         clip.ClipTitle.String,
		StreamTitle:       clip.StreamTitle.String,
		ClippedBy:         clip.ClippedBy.String,
		RelativeStartTime: float32(clip.RelativeStartTime.Float64),
		RelativeEndTime:   float32(clip.RelativeEndTime.Float64),
		DurationSeconds:   int(math.Round(clip.RelativeEndTime.Float64 - clip.RelativeStartTime.Float64)),
		Timestamp:         clip.ClipTimestamp.Time,
		Manifest:          "/clip/" + clip.ClipID,
		Thumbnail:         clipThumbnailURL(clip.ClipID),
	}

	return &c, nil
}

// GetAllSegmentsForOutputConfigurationAndWindow returns all the segments for a
// given output config and time window.
func (p *PlaylistGenerator) GetAllSegmentsForOutputConfigurationAndWindow(configID string, startSeconds, endSeconds float32) ([]HLSSegment, error) {
	segmentRows, err := p.datastore.GetQueries().GetSegmentsForOutputIdAndWindow(context.Background(), db.GetSegmentsForOutputIdAndWindowParams{
		OutputConfigurationID: configID,
		StartSeconds:          float64(startSeconds),
		EndSeconds:            float64(endSeconds),
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
			Duration:              row.Duration.Float64,
			MediaOffset:           row.MediaOffset.Float64,
			Timestamp:             row.Timestamp.Time,
			Path:                  row.Path,
		}
		segments = append(segments, segment)
	}

	return segments, nil
}
