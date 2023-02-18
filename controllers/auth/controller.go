package auth

import (
	"fmt"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/controllers/auth/fediverse"
	"github.com/owncast/owncast/controllers/auth/indieauth"
)

func New(s *controllers.Service) (c *Controller, err error) {
	c = &Controller{
		Service: s,
	}

	if c.IndieAuth, err = indieauth.New(s); err != nil {
		return nil, fmt.Errorf("initializing indieauth controller: %v", err)
	}

	if c.Fediverse, err = fediverse.New(s); err != nil {
		return nil, fmt.Errorf("initializing fediverse controller: %v", err)
	}

	return c, nil
}

type Controller struct {
	*controllers.Service
	IndieAuth *indieauth.Controller
	Fediverse *fediverse.Controller
}
