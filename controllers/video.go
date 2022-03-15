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

	streamSortVariants := make([]variantsSort, len(outputVariants))
	for i, variant := range outputVariants {
		variantSort := variantsSort{
			Index:              i,
			Name:               variant.GetName(),
			IsVideoPassthrough: variant.IsVideoPassthrough,
			VideoBitrate:       variant.VideoBitrate,
		}
		streamSortVariants[i] = variantSort
	}

	sort.Slice(streamSortVariants, func(i, j int) bool {
		if streamSortVariants[i].IsVideoPassthrough && !streamSortVariants[j].IsVideoPassthrough {
			return true
		}

		if !streamSortVariants[i].IsVideoPassthrough && streamSortVariants[j].IsVideoPassthrough {
			return false
		}

		return streamSortVariants[i].VideoBitrate > streamSortVariants[j].VideoBitrate
	})

	response := make([]variantsResponse, len(streamSortVariants))
	for i, variant := range streamSortVariants {
		variantResponse := variantsResponse{
			Index: variant.Index,
			Name:  variant.Name,
		}
		response[i] = variantResponse
	}

	WriteResponse(w, response)
}

// GetKeepWholeStream will return the current value of the keep whole stream setting.
func GetKeepWholeStream(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, data.GetKeepWholeStream())
}
