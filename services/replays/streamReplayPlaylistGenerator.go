package replays

import (
	"context"
	"strings"

	"github.com/grafov/m3u8"
	"github.com/pkg/errors"

	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/services/datastore"
)

/*
The PlaylistGenerator is responsible for creating the master and media
playlists, in order to replay a stream in whole, or part. It requires detailed
metadata about how the initial live stream was configured, as well as
access to every segment that was created during the live stream.
*/

// PlaylistGenerator reconstructs HLS playlists from recorded stream metadata.
type PlaylistGenerator struct {
	datastore *datastore.Datastore
}

// GenerateMasterPlaylistForStream returns a master playlist referencing a media
// playlist for each output configuration of a recorded stream.
func (p *PlaylistGenerator) GenerateMasterPlaylistForStream(streamID string) (*m3u8.MasterPlaylist, error) {
	// Determine the different output configurations for this stream.
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

		mediaPlaylist, err := p.GenerateMediaPlaylistForStreamAndConfiguration(streamID, config.ID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create media playlist")
		}

		// Append the media playlist to the master playlist.
		params := p.getMediaPlaylistParamsForConfig(config)

		// Add the media playlist to the master playlist.
		publicPlaylistPath := strings.Join([]string{"/replay", streamID, config.ID}, "/")
		masterPlaylist.Append(publicPlaylistPath, mediaPlaylist, params)
	}

	// Return the final master playlist that contains all the media playlists.
	return masterPlaylist, nil
}

// GenerateMediaPlaylistForStreamAndConfiguration returns the media playlist for
// a single output configuration of a recorded stream.
func (p *PlaylistGenerator) GenerateMediaPlaylistForStreamAndConfiguration(streamID, outputConfigurationID string) (*m3u8.MediaPlaylist, error) {
	stream, err := p.GetStream(streamID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stream")
	}

	config, err := p.GetOutputConfig(outputConfigurationID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get output configuration")
	}

	// Fetch all the segments for this configuration.
	segments, err := p.GetAllSegmentsForOutputConfiguration(outputConfigurationID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all segments for output configuration")
	}

	// Create the media playlist for this configuration and add the segments.
	mediaPlaylist, err := p.createMediaPlaylistForConfigurationAndSegments(config, stream.InProgress, segments)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create media playlist")
	}

	return mediaPlaylist, nil
}

// GetStream returns the recorded stream metadata for a stream ID.
func (p *PlaylistGenerator) GetStream(streamID string) (*Stream, error) {
	stream, err := p.datastore.GetQueries().GetStreamById(context.Background(), streamID)
	if err != nil || stream.ID == "" {
		return nil, errors.Wrap(err, "failed to get stream")
	}

	s := Stream{
		ID:         stream.ID,
		Title:      stream.StreamTitle.String,
		StartTime:  stream.StartTime.Time,
		EndTime:    stream.EndTime.Time,
		InProgress: !stream.EndTime.Valid,
	}

	return &s, nil
}

// GetOutputConfig returns a single output configuration by ID.
func (p *PlaylistGenerator) GetOutputConfig(outputConfigID string) (*HLSOutputConfiguration, error) {
	config, err := p.datastore.GetQueries().GetOutputConfigurationForId(context.Background(), outputConfigID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get output configuration")
	}

	return createConfigFromConfigRow(config), nil
}

func createConfigFromConfigRow(row db.GetOutputConfigurationForIdRow) *HLSOutputConfiguration {
	config := HLSOutputConfiguration{
		ID:              row.ID,
		StreamID:        row.StreamID,
		VariantID:       row.VariantID,
		Name:            row.Name,
		VideoBitrate:    int(row.Bitrate),
		Framerate:       int(row.Framerate),
		ScaledWidth:     int(row.ResolutionWidth.Int64),
		ScaledHeight:    int(row.ResolutionHeight.Int64),
		SegmentDuration: float64(row.SegmentDuration),
	}
	return &config
}
