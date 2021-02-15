package controllers

import (
	"fmt"
	"math"
	"net/http"
	"sort"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
)

type variants []models.StreamOutputVariant

type variantsResponse struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
}

// Len returns the number of variants.
func (v variants) Len() int { return len(v) }

// Less does less.
func (v variants) Less(i, j int) bool { return v[i].VideoBitrate < v[j].VideoBitrate }

// Swap will swap.
func (v variants) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

// GetVideoStreamOutputVariants will return the video variants available
func GetVideoStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	outputVariants := data.GetStreamOutputVariants()
	result := make([]variantsResponse, len(outputVariants))

	for i, variant := range outputVariants {
		var name string
		bitrate := getBitrateString(variant.VideoBitrate)

		if variant.IsVideoPassthrough {
			name = "Source"
		} else if variant.ScaledHeight == 720 && variant.ScaledWidth == 1080 {
			name = fmt.Sprintf("720p @%s", bitrate)
		} else if variant.ScaledHeight == 1080 && variant.ScaledWidth == 1920 {
			name = fmt.Sprintf("1080p @%s", bitrate)
		} else if variant.ScaledHeight != 0 {
			name = fmt.Sprintf("%dh", variant.ScaledHeight)
		} else if variant.ScaledWidth != 0 {
			name = fmt.Sprintf("%dw", variant.ScaledWidth)
		} else {
			name = fmt.Sprintf("%s@%dfps", bitrate, variant.Framerate)
		}

		variantResponse := variantsResponse{
			Index: i,
			Name:  name,
		}
		result[i] = variantResponse
	}

	sort.Slice(result, func(i, j int) bool {
		// Sort video passthrough to the front of the list
		if outputVariants[i].IsVideoPassthrough {
			return true
		}

		return outputVariants[i].VideoBitrate > outputVariants[j].VideoBitrate
	})

	WriteResponse(w, result)
}

func getBitrateString(bitrate int) string {
	if bitrate == 0 {
		return ""
	} else if bitrate < 1000 {
		return fmt.Sprintf("%dKbps", bitrate)
	} else if bitrate >= 1000 {
		if math.Mod(float64(bitrate), 1000) == 0 {
			return fmt.Sprintf("%dMbps", bitrate/1000.0)
		} else {
			return fmt.Sprintf("%.1fMbps", float32(bitrate)/1000.0)
		}
	}

	return ""
}
