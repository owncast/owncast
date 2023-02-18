package controllers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/rtmp"
	"github.com/owncast/owncast/internal/activitypub"
	"github.com/owncast/owncast/internal/activitypub/follower"
	"github.com/owncast/owncast/internal/activitypub/outbox"
	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/notifications"
	"github.com/owncast/owncast/static"
)

type j map[string]interface{}

func New(
	ap *activitypub.Service,
	c *core.Service,
	m *metrics.Service,
	n *notifications.Notifier,
	r *rtmp.Service,
	f *follower.Service,
	o *outbox.Service,
) (*Service, error) {
	s := &Service{
		ActivityPub:   ap,
		Data:          ap.Persistence.Data,
		Core:          c,
		Metrics:       m,
		Notifications: n,
		Rtmp:          r,
		Follower:      f,
		Outbox:        o,
	}

	s.static = http.FileServer(http.FS(static.GetWeb()))

	return s, nil
}

type Service struct {
	Data          *data.Service
	ActivityPub   *activitypub.Service
	Core          *core.Service
	Metrics       *metrics.Service
	static        http.Handler
	Notifications *notifications.Notifier
	Rtmp          *rtmp.Service
	controllers   []interface{ Register(*Service) }
	Follower      *follower.Service
	Outbox        *outbox.Service
}

// InternalErrorHandler will return an error message as an HTTP response.
func (s *Service) InternalErrorHandler(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(j{"error": err.Error()}); err != nil {
		s.InternalErrorHandler(w, err)
	}
}

// BadRequestHandler will return an HTTP 500 as an HTTP response.
func (s *Service) BadRequestHandler(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	log.Debugln(err)

	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(j{"error": err.Error()}); err != nil {
		s.InternalErrorHandler(w, err)
	}
}

// WriteSimpleResponse will return a message as a response.
func (s *Service) WriteSimpleResponse(w http.ResponseWriter, success bool, message string) {
	response := models.BaseAPIResponse{
		Success: success,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")

	if success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.InternalErrorHandler(w, err)
	}
}

// WriteResponse will return an object as a JSON encoded response.
func (s *Service) WriteResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.InternalErrorHandler(w, err)
	}
}

// WriteString will return a basic string and a status code to the client.
func (s *Service) WriteString(w http.ResponseWriter, text string, status int) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	_, err := w.Write([]byte(text))
	return err
}
