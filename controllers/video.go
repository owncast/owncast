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

// Len returns the number of variants.
func (v variants) Len() int { return len(v) }

// Less is less than..
func (v variants) Less(i, j int) bool { return v[i].VideoBitrate < v[j].VideoBitrate }

// Swap will swap two values.
func (v variants) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

// GetVideoStreamOutputVariants will return the video variants available
func GetVideoStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	outputVariants := data.GetStreamOutputVariants()
	result := make([]variantsResponse, len(outputVariants))

	for i, variant := range outputVariants {
		variantResponse := variantsResponse{
			Index: i,
			Name:  variant.GetName(),
		}
		result[i] = variantResponse
	}

	sort.Slice(result, func(i, j int) bool {
		return outputVariants[i].VideoBitrate > outputVariants[j].VideoBitrate || !outputVariants[i].IsVideoPassthrough
	})

	WriteResponse(w, result)
}
