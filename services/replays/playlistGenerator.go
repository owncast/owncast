package replays

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/grafov/m3u8"
	"github.com/pkg/errors"
)

// GetConfigurationsForStream returns the output configurations for a given stream.
func (p *PlaylistGenerator) GetConfigurationsForStream(streamID string) ([]*HLSOutputConfiguration, error) {
	outputConfigRows, err := p.datastore.GetQueries().GetOutputConfigurationsForStreamId(context.Background(), streamID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get output configurations for stream")
	}

	outputConfigs := []*HLSOutputConfiguration{}
	for _, row := range outputConfigRows {
		config := &HLSOutputConfiguration{
			ID:              row.ID,
			StreamID:        streamID,
			VariantID:       row.VariantID,
			Name:            row.Name,
			VideoBitrate:    int(row.Bitrate),
			Framerate:       int(row.Framerate),
			ScaledWidth:     int(row.ResolutionWidth.Int64),
			ScaledHeight:    int(row.ResolutionHeight.Int64),
			SegmentDuration: float64(row.SegmentDuration),
		}
		outputConfigs = append(outputConfigs, config)
	}

	return outputConfigs, nil
}

func (p *PlaylistGenerator) createMediaPlaylistForConfigurationAndSegments(configuration *HLSOutputConfiguration, startTime time.Time, inProgress bool, segments []HLSSegment) (*m3u8.MediaPlaylist, error) {
	playlistSize := len(segments)
	segmentDuration := configuration.SegmentDuration
	playlist, err := m3u8.NewMediaPlaylist(0, uint(playlistSize))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create media playlist")
	}

	playlist.TargetDuration = configuration.SegmentDuration

	if !inProgress {
		playlist.MediaType = m3u8.VOD
	} else {
		playlist.MediaType = m3u8.EVENT
	}

	// Add the segments to the playlist.
	for index, segment := range segments {
		// If it's a URL leave it as is, if it's a local path then append a slash.
		path := segment.Path
		if !strings.HasPrefix(path, "http") {
			path = "/" + path
		}

		mediaSegment := m3u8.MediaSegment{
			URI:             path,
			Duration:        segmentDuration,
			SeqId:           uint64(index),
			ProgramDateTime: segment.Timestamp,
		}
		if err := playlist.AppendSegment(&mediaSegment); err != nil {
			return nil, errors.Wrap(err, "failed to append segment to recording playlist")
		}
	}

	// Configure the properties of this media playlist. SetProgramDateTime
	// attaches to the first segment, so it's only meaningful (and only valid)
	// when the playlist actually has segments. A stream that has just started
	// may legitimately have none yet.
	if len(segments) > 0 {
		if err := playlist.SetProgramDateTime(startTime); err != nil {
			return nil, errors.Wrap(err, "failed to set media playlist program date time")
		}
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

// GetAllSegmentsForOutputConfiguration returns all the segments for a given output config.
func (p *PlaylistGenerator) GetAllSegmentsForOutputConfiguration(outputID string) ([]HLSSegment, error) {
	segmentRows, err := p.datastore.GetQueries().GetSegmentsForOutputId(context.Background(), outputID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get segments for output config")
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

func (p *PlaylistGenerator) getMediaPlaylistParamsForConfig(config *HLSOutputConfiguration) m3u8.VariantParams {
	// Clamp to the uint32 range the playlist library requires.
	bandwidth := max(config.VideoBitrate, 0) * 1000
	if bandwidth > math.MaxUint32 {
		bandwidth = math.MaxUint32
	}
	params := m3u8.VariantParams{
		ProgramId: 1,
		Name:      config.Name,
		FrameRate: float64(config.Framerate),
		Bandwidth: uint32(bandwidth),
		// Match what is generated in our live playlists.
		Codecs: "avc1.64001f,mp4a.40.2",
	}

	// If both the width and height are set then we can set that as
	// the resolution in the media playlist.
	if config.ScaledHeight > 0 && config.ScaledWidth > 0 {
		params.Resolution = fmt.Sprintf("%dx%d", config.ScaledWidth, config.ScaledHeight)
	}

	return params
}
