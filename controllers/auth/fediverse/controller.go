package fediverse

import (
	"github.com/owncast/owncast/controllers"
)

func New(s *controllers.Service) (*Controller, error) {
	c := &Controller{
		Service: s,
	}

	return c, nil
}

type Controller struct {
	*controllers.Service
}
