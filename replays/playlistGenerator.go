package replays

import (
	"context"
	"fmt"
	"strings"
	"time"

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
		if config.VideoBitrate == 0 {
			return nil, errors.New("video bitrate is unavailable")
		}

		if config.Framerate == 0 {
			return nil, errors.New("video framerate is unavailable")
		}

		mediaPlaylist, err := p.GenerateMediaPlaylistForStreamAndConfiguration(streamId, config.ID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create media playlist")
		}

		// Append the media playlist to the master playlist.
		params := m3u8.VariantParams{
			ProgramId: 1,
			Name:      config.Name,
			FrameRate: float64(config.Framerate),
			Bandwidth: uint32(config.VideoBitrate * 1000),
			// Match what is generated in our live playlists.
			Codecs: "avc1.64001f,mp4a.40.2",
		}

		// If both the width and height are set then we can set that as
		// the resolution in the media playlist.
		if config.ScaledHeight > 0 && config.ScaledWidth > 0 {
			params.Resolution = fmt.Sprintf("%dx%d", config.ScaledWidth, config.ScaledHeight)
		}

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
		StartTime:  stream.StartTime,
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

// GetConfigurationsForStream returns the output configurations for a given stream.
func (p *PlaylistGenerator) GetConfigurationsForStream(streamId string) ([]*HLSOutputConfiguration, error) {
	outputConfigRows, err := p.datastore.GetQueries().GetOutputConfigurationsForStreamId(context.Background(), streamId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get output configurations for stream")
	}

	outputConfigs := []*HLSOutputConfiguration{}
	for _, row := range outputConfigRows {
		config := &HLSOutputConfiguration{
			ID:              row.ID,
			StreamId:        streamId,
			VariantId:       row.VariantID,
			Name:            row.Name,
			VideoBitrate:    int(row.Bitrate),
			Framerate:       int(row.Framerate),
			ScaledHeight:    int(row.ResolutionWidth.Int32),
			ScaledWidth:     int(row.ResolutionHeight.Int32),
			SegmentDuration: float64(row.SegmentDuration),
		}
		outputConfigs = append(outputConfigs, config)
	}

	return outputConfigs, nil
}

// GetAllSegmentsForOutputConfiguration returns all the segments for a given output config.
func (p *PlaylistGenerator) GetAllSegmentsForOutputConfiguration(outputId string) ([]HLSSegment, error) {
	segmentRows, err := p.datastore.GetQueries().GetSegmentsForOutputId(context.Background(), outputId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get segments for output config")
	}

	segments := []HLSSegment{}
	for _, row := range segmentRows {
		segment := HLSSegment{
			ID:                    row.ID,
			StreamID:              row.StreamID,
			OutputConfigurationID: row.OutputConfigurationID,
			Timestamp:             row.Timestamp,
			Path:                  row.Path,
		}
		segments = append(segments, segment)
	}

	return segments, nil
}

func (p *PlaylistGenerator) createMediaPlaylistForConfigurationAndSegments(configuration *HLSOutputConfiguration, startTime time.Time, inProgress bool, segments []HLSSegment) (*m3u8.MediaPlaylist, error) {
	playlistSize := len(segments)
	segmentDuration := configuration.SegmentDuration
	playlist, err := m3u8.NewMediaPlaylist(0, uint(playlistSize))

	playlist.TargetDuration = configuration.SegmentDuration

	if !inProgress {
		playlist.MediaType = m3u8.VOD
	} else {
		playlist.MediaType = m3u8.EVENT
	}

	// Add the segments to the playlist.
	for index, segment := range segments {
		mediaSegment := m3u8.MediaSegment{
			URI:             "/" + segment.Path,
			Duration:        segmentDuration,
			SeqId:           uint64(index),
			ProgramDateTime: segment.Timestamp,
		}
		if err := playlist.AppendSegment(&mediaSegment); err != nil {
			return nil, errors.Wrap(err, "failed to append segment to recording playlist")
		}
	}

	if err != nil {
		return nil, err
	}

	// Configure the properties of this media playlist.
	if err := playlist.SetProgramDateTime(startTime); err != nil {
		return nil, errors.Wrap(err, "failed to set media playlist program date time")
	}

	// Our live output is specified as v6, so let's match it to be as close as
	// possible to what we're doing for live streams.
	playlist.SetVersion(6)

	if !inProgress {
		// Specify explicitly that the playlist content is allowed to be cached.
		// However, if in-progress recordings are supported this should not be enabled
		// in order for the playlist to be updated with new segments. inProgress is
		// determined by seeing if the stream has an endTime or not.
		playlist.SetCustomTag(&MediaPlaylistAllowCacheTag{})

		// Set the ENDLIST tag and close the playlist for writing if the stream is
		// not still in progress.
		playlist.Close()
	}

	return playlist, nil
}

func (p *PlaylistGenerator) createNewMasterPlaylist() *m3u8.MasterPlaylist {
	playlist := m3u8.NewMasterPlaylist()
	playlist.SetIndependentSegments(true)
	playlist.SetVersion(6)

	return playlist
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
