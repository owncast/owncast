package replays

import (
	"context"
	"fmt"
	"time"

	"github.com/grafov/m3u8"
	"github.com/pkg/errors"
)

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
			Timestamp:             row.Timestamp.Time,
			Path:                  row.Path,
		}
		segments = append(segments, segment)
	}

	return segments, nil
}

func (p *PlaylistGenerator) getMediaPlaylistParamsForConfig(config *HLSOutputConfiguration) m3u8.VariantParams {
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

	return params
}
