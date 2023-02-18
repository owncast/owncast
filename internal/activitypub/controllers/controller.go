package controllers

import (
	"fmt"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/internal/activitypub/apmodels"
	"github.com/owncast/owncast/internal/activitypub/crypto"
	"github.com/owncast/owncast/internal/activitypub/follower"
	"github.com/owncast/owncast/internal/activitypub/inbox"
	"github.com/owncast/owncast/internal/activitypub/persistence"
)

func New(p *persistence.Service, m *apmodels.Service, f *follower.Service) (*Service, error) {
	c := &Service{}

	cryptoService, err := crypto.New(p.Data)
	if err != nil {
		return nil, fmt.Errorf("initializing crypto service: %v", err)
	}

	c.data = p.Data
	c.crypto = cryptoService
	c.models = m
	c.follower = f

	return c, nil
}

type Service struct {
	data        *data.Service
	persistence *persistence.Service
	crypto      *crypto.Service
	models      *apmodels.Service
	follower    *follower.Service
	inbox       *inbox.Service
}
