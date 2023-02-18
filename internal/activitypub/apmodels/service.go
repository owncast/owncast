package apmodels

import (
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/internal/activitypub/crypto"
)

func New(d *data.Service, c *crypto.Service) (*Service, error) {
	s := &Service{}

	s.crypto = c

	return s, nil
}

type Service struct {
	data   *data.Service
	crypto *crypto.Service
}
