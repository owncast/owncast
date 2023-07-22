package status

import (
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/storage/configrepository"
)

type Status struct {
	models.Stats
	models.Status

	broadcast       *models.CurrentBroadcast
	broadcaster     *models.Broadcaster
	StreamConnected bool

	VersionNumber         string `json:"versionNumber"`
	StreamTitle           string `json:"streamTitle"`
	ViewerCount           int    `json:"viewerCount"`
	OverallMaxViewerCount int    `json:"overallMaxViewerCount"`
	SessionMaxViewerCount int    `json:"sessionMaxViewerCount"`

	Online bool `json:"online"`
}

var temporaryGlobalInstance *Status

func New() *Status {
	configRepository := configrepository.Get()

	return &Status{
		StreamTitle: configRepository.GetStreamTitle(),
	}
}

// Get will return the global instance of the status service.
func Get() *Status {
	if temporaryGlobalInstance == nil {
		temporaryGlobalInstance = &Status{}
	}

	return temporaryGlobalInstance
}

// GetCurrentBroadcast will return the currently active broadcast.
func (s *Status) GetCurrentBroadcast() *models.CurrentBroadcast {
	return s.broadcast
}

func (s *Status) SetCurrentBroadcast(broadcast *models.CurrentBroadcast) {
	s.broadcast = broadcast
}

// SetBroadcaster will store the current inbound broadcasting details.
func (s *Status) SetBroadcaster(broadcaster *models.Broadcaster) {
	s.broadcaster = broadcaster
}

// GetBroadcaster will return the details of the currently active broadcaster.
func (s *Status) GetBroadcaster() *models.Broadcaster {
	return s.broadcaster
}
