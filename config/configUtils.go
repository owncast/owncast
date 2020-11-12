package config

import (
	"encoding/json"
	"sort"
)

func findHighestQuality(qualities []StreamQuality) int {
	type IndexedQuality struct {
		index   int
		quality StreamQuality
	}

	if len(qualities) < 2 {
		return 0
	}

	indexedQualities := make([]IndexedQuality, 0)
	for index, quality := range qualities {
		indexedQuality := IndexedQuality{index, quality}
		indexedQualities = append(indexedQualities, indexedQuality)
	}

	sort.Slice(indexedQualities, func(a, b int) bool {
		if indexedQualities[a].quality.IsVideoPassthrough && !indexedQualities[b].quality.IsVideoPassthrough {
			return true
		}

		if !indexedQualities[a].quality.IsVideoPassthrough && indexedQualities[b].quality.IsVideoPassthrough {
			return false
		}

		return indexedQualities[a].quality.VideoBitrate > indexedQualities[b].quality.VideoBitrate
	})

	return indexedQualities[0].index
}

// MarshalJSON is a custom JSON marshal function for video stream qualities.
func (q *StreamQuality) MarshalJSON() ([]byte, error) {
	type Alias StreamQuality
	return json.Marshal(&struct {
		Framerate int `json:"framerate"`
		*Alias
	}{
		Framerate: q.GetFramerate(),
		Alias:     (*Alias)(q),
	})
}
