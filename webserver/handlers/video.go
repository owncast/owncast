package handlers

import (
	"net/http"
	"sort"

	webutils "github.com/owncast/owncast/webserver/utils"
)

type variantsSort struct {
	Name               string
	Index              int
	VideoBitrate       int
	IsVideoPassthrough bool
}

type variantsResponse struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
}

// GetVideoStreamOutputVariants will return the video variants available.
func (h *Handlers) GetVideoStreamOutputVariants(w http.ResponseWriter, r *http.Request) {
	outputVariants := h.configRepository.GetStreamOutputVariants()

	streamSortVariants := make([]variantsSort, 0, len(outputVariants))
	enabledIndex := 0
	for _, variant := range outputVariants {
		if !variant.Enabled {
			continue
		}
		variantSort := variantsSort{
			Index:              enabledIndex,
			Name:               variant.GetName(),
			IsVideoPassthrough: variant.IsVideoPassthrough,
			VideoBitrate:       variant.VideoBitrate,
		}
		streamSortVariants = append(streamSortVariants, variantSort)
		enabledIndex++
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

	webutils.WriteResponse(w, response)
}
