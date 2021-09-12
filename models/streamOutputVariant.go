package models

import (
	"encoding/json"
	"fmt"
	"math"
)

// StreamOutputVariant defines the output specifics of a single HLS stream variant.
type StreamOutputVariant struct {
	// Name is an optional human-readable label for this stream output.
	Name string `json:"name"`

	// Enable passthrough to copy the video and/or audio directly from the
	// incoming stream and disable any transcoding.  It will ignore any of
	// the below settings.
	IsVideoPassthrough bool `yaml:"videoPassthrough" json:"videoPassthrough"`
	IsAudioPassthrough bool `yaml:"audioPassthrough" json:"audioPassthrough"`

	VideoBitrate int `yaml:"videoBitrate" json:"videoBitrate"`
	AudioBitrate int `yaml:"audioBitrate" json:"audioBitrate"`

	// Set only one of these in order to keep your current aspect ratio.
	// Or set neither to not scale the video.
	ScaledWidth  int `yaml:"scaledWidth" json:"scaledWidth,omitempty"`
	ScaledHeight int `yaml:"scaledHeight" json:"scaledHeight,omitempty"`

	Framerate int `yaml:"framerate" json:"framerate"`
	// CPUUsageLevel represents a codec preset to configure CPU usage.
	CPUUsageLevel int `json:"cpuUsageLevel"`
}

// GetFramerate returns the framerate or default.
func (q *StreamOutputVariant) GetFramerate() int {
	if q.IsVideoPassthrough {
		return 0
	}

	if q.Framerate > 0 {
		return q.Framerate
	}

	return 24
}

// GetIsAudioPassthrough will return if this variant audio is passthrough.
func (q *StreamOutputVariant) GetIsAudioPassthrough() bool {
	if q.IsAudioPassthrough {
		return true
	}

	if q.AudioBitrate == 0 {
		return true
	}

	return false
}

// GetName will return the human readable name for this stream output.
func (q *StreamOutputVariant) GetName() string {
	bitrate := getBitrateString(q.VideoBitrate)

	if q.Name != "" {
		return q.Name
	} else if q.IsVideoPassthrough {
		return "Source"
	} else if q.ScaledHeight == 720 && q.ScaledWidth == 1080 {
		return fmt.Sprintf("720p @%s", bitrate)
	} else if q.ScaledHeight == 1080 && q.ScaledWidth == 1920 {
		return fmt.Sprintf("1080p @%s", bitrate)
	} else if q.ScaledHeight != 0 {
		return fmt.Sprintf("%dh", q.ScaledHeight)
	} else if q.ScaledWidth != 0 {
		return fmt.Sprintf("%dw", q.ScaledWidth)
	} else {
		return fmt.Sprintf("%s@%dfps", bitrate, q.Framerate)
	}
}

func getBitrateString(bitrate int) string {
	if bitrate == 0 {
		return ""
	} else if bitrate < 1000 {
		return fmt.Sprintf("%dKbps", bitrate)
	} else if bitrate >= 1000 {
		if math.Mod(float64(bitrate), 1000) == 0 {
			return fmt.Sprintf("%dMbps", bitrate/1000.0)
		}

		return fmt.Sprintf("%.1fMbps", float32(bitrate)/1000.0)
	}

	return ""
}

// MarshalJSON is a custom JSON marshal function for video stream qualities.
func (q *StreamOutputVariant) MarshalJSON() ([]byte, error) {
	type Alias StreamOutputVariant
	return json.Marshal(&struct {
		Framerate int `json:"framerate"`
		*Alias
	}{
		Framerate: q.GetFramerate(),
		Alias:     (*Alias)(q),
	})
}
