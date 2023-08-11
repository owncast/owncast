package replays

import (
	"context"
	"fmt"
	"time"

	"github.com/owncast/owncast/core/data"
	"github.com/pkg/errors"
)

type Stream struct {
	ID         string    `json:"id"`
	Title      string    `json:"title,omitempty"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime,omitempty"`
	InProgress bool      `json:"inProgress,omitempty"`
	Manifest   string    `json:"manifest,omitempty"`
}

// GetStreams will return all streams that have been recorded.
func GetStreams() ([]*Stream, error) {
	streams, err := data.GetDatastore().GetQueries().GetStreams(context.Background())
	if err != nil {
		return nil, errors.WithMessage(err, "failure to get streams")
	}

	response := []*Stream{}
	for _, stream := range streams {
		s := Stream{
			ID:         stream.ID,
			Title:      stream.StreamTitle.String,
			StartTime:  stream.StartTime.Time,
			EndTime:    stream.EndTime.Time,
			InProgress: !stream.EndTime.Valid,
			Manifest:   fmt.Sprintf("/replay/%s", stream.ID),
		}
		response = append(response, &s)
	}
	return response, nil
}
