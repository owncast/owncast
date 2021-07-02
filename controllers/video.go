package controllers

import (
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

// GetVideoStreamOutputVariants will return the video variants available.
func GetVideoStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	outputVariants := data.GetStreamOutputVariants()

	sort.Slice(outputVariants, func(i, j int) bool {
		return outputVariants[j].VideoBitrate < outputVariants[i].VideoBitrate
	})

	result := make([]variantsResponse, len(outputVariants))

	for i, variant := range outputVariants {
		variantResponse := variantsResponse{
			Index: i,
			Name:  variant.GetName(),
		}
		result[i] = variantResponse
	}

	WriteResponse(w, result)
}
