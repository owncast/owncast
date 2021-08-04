package controllers

import (
	"net/http"
	"sort"

	"github.com/owncast/owncast/core/data"
)

type variantsSort struct {
	Index              int
	Name               string
	IsVideoPassthrough bool
	VideoBitrate       int
}

type variantsResponse struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
}

// GetVideoStreamOutputVariants will return the video variants available.
func GetVideoStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	outputVariants := data.GetStreamOutputVariants()

	StreamSortVariants := make([]variantsSort, len(outputVariants))
	for i, variant := range outputVariants {
		variantSort := variantsSort{
			Index:              i,
			Name:               variant.GetName(),
			IsVideoPassthrough: variant.IsVideoPassthrough,
			VideoBitrate:       variant.VideoBitrate,
		}
		StreamSortVariants[i] = variantSort
	}

	sort.Slice(StreamSortVariants, func(i, j int) bool {
		if StreamSortVariants[i].IsVideoPassthrough && !StreamSortVariants[j].IsVideoPassthrough {
			return true
		}

		if !StreamSortVariants[i].IsVideoPassthrough && StreamSortVariants[j].IsVideoPassthrough {
			return false
		}

		return StreamSortVariants[i].VideoBitrate > StreamSortVariants[j].VideoBitrate
	})

	response := make([]variantsResponse, len(StreamSortVariants))
	for i, variant := range StreamSortVariants {
		variantResponse := variantsResponse{
			Index: variant.Index,
			Name:  variant.Name,
		}
		response[i] = variantResponse
	}

	WriteResponse(w, response)
}
