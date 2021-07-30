package controllers

import (
	"net/http"
	"sort"

	"github.com/owncast/owncast/core/data"
)

type variantsResponse struct {
	Name               string `json:"name"`
	Index              int    `json:"index"`
	VideoBitrate       int    `json:"-"`
	IsVideoPassthrough bool   `json:"-"`
}

// GetVideoStreamOutputVariants will return the video variants available.
func GetVideoStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	outputVariants := data.GetStreamOutputVariants()

	result := make([]variantsResponse, len(outputVariants))
	for i, variant := range outputVariants {
		variantResponse := variantsResponse{
			Index:              i,
			Name:               variant.GetName(),
			VideoBitrate:       variant.VideoBitrate,
			IsVideoPassthrough: variant.IsVideoPassthrough,
		}
		result[i] = variantResponse
	}

	sort.Slice(result, func(i, j int) bool {
		if outputVariants[i].IsVideoPassthrough && !outputVariants[j].IsVideoPassthrough {
			return true
		}

		if !outputVariants[i].IsVideoPassthrough && outputVariants[j].IsVideoPassthrough {
			return false
		}

		return outputVariants[i].VideoBitrate > outputVariants[j].VideoBitrate
	})

	WriteResponse(w, result)
}
