package webhooks

import (
	"github.com/owncast/owncast/core/data"
)

func New(d *data.Service) (*Service, error) {
	s := &Service{}

	s.data = d

	return s, nil
}

type Service struct {
	data *data.Service
}
