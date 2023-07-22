package state

import (
	"github.com/owncast/owncast/services/chat"
	"github.com/owncast/owncast/services/metrics"
	"github.com/owncast/owncast/services/status"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/video/transcoder"
)

type VideoState struct {
	transcoder   *transcoder.Transcoder
	metrics      *metrics.Metrics
	status       *status.Status
	lastNotified utils.NullTime
	chatService  *chat.Chat
}

func New() *VideoState {
	return &VideoState{
		transcoder: transcoder.Get(),
		metrics:    metrics.Get(),
		status:     status.Get(),
	}
}

var temporaryGlobalInstance *VideoState

func Get() *VideoState {
	if temporaryGlobalInstance == nil {
		temporaryGlobalInstance = New()
	}

	return temporaryGlobalInstance
}
