package crypto

import (
	"github.com/owncast/owncast/core/data"
)

func New(ds *data.Service) (*Service, error) {
	s := &Service{
		Data: ds,
	}

	return s, nil
}

type Service struct {
	Data *data.Service
}
