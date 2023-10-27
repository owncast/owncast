package replays

import (
	"context"
	"strings"

	"github.com/grafov/m3u8"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/db"
	"github.com/pkg/errors"
)

/*
The PlaylistGenerator is responsible for creating the master and media
playlists, in order to replay a stream in whole, or part. It requires detailed
metadata about how the initial live stream was configured, as well as a
access to every segment that was created during the live stream.
*/

type PlaylistGenerator struct {
	datastore *data.Datastore
}

func NewPlaylistGenerator() *PlaylistGenerator {
	return &PlaylistGenerator{
		datastore: data.GetDatastore(),
	}
}

func (p *PlaylistGenerator) GenerateMasterPlaylistForStream(streamId string) (*m3u8.MasterPlaylist, error) {
	// Determine the different output configurations for this stream.
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

		mediaPlaylist, err := p.GenerateMediaPlaylistForStreamAndConfiguration(streamId, config.ID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create media playlist")
		}

		// Append the media playlist to the master playlist.
		params := p.getMediaPlaylistParamsForConfig(config)

		// Add the media playlist to the master playlist.
		publicPlaylistPath := strings.Join([]string{"/replay", streamId, config.ID}, "/")
		masterPlaylist.Append(publicPlaylistPath, mediaPlaylist, params)
	}

	// Return the final master playlist that contains all the media playlists.
	return masterPlaylist, nil
}

func (p *PlaylistGenerator) GenerateMediaPlaylistForStreamAndConfiguration(streamId, outputConfigurationId string) (*m3u8.MediaPlaylist, error) {
	stream, err := p.GetStream(streamId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stream")
	}

	config, err := p.GetOutputConfig(outputConfigurationId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get output configuration")
	}

	// Fetch all the segments for this configuration.
	segments, err := p.GetAllSegmentsForOutputConfiguration(outputConfigurationId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all segments for output configuration")
	}

	// Create the media playlist for this configuration and add the segments.
	mediaPlaylist, err := p.createMediaPlaylistForConfigurationAndSegments(config, stream.StartTime, stream.InProgress, segments)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create media playlist")
	}

	return mediaPlaylist, nil
}

func (p *PlaylistGenerator) GetStream(streamId string) (*Stream, error) {
	stream, err := p.datastore.GetQueries().GetStreamById(context.Background(), streamId)
	if stream.ID == "" {
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

func (p *PlaylistGenerator) GetOutputConfig(outputConfigId string) (*HLSOutputConfiguration, error) {
	config, err := p.datastore.GetQueries().GetOutputConfigurationForId(context.Background(), outputConfigId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get output configuration")
	}

	return createConfigFromConfigRow(config), nil
}

func createConfigFromConfigRow(row db.GetOutputConfigurationForIdRow) *HLSOutputConfiguration {
	config := HLSOutputConfiguration{
		ID:              row.ID,
		StreamId:        row.StreamID,
		VariantId:       row.VariantID,
		Name:            row.Name,
		VideoBitrate:    int(row.Bitrate),
		Framerate:       int(row.Framerate),
		ScaledHeight:    int(row.ResolutionWidth.Int32),
		ScaledWidth:     int(row.ResolutionHeight.Int32),
		SegmentDuration: float64(row.SegmentDuration),
	}
	return &config
}
