package inbox

import (
	"fmt"

	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/internal/activitypub/crypto"
	"github.com/owncast/owncast/internal/activitypub/follower"
	"github.com/owncast/owncast/internal/activitypub/persistence"
)

func New(p *persistence.Service) (*Service, error) {
	s := &Service{}

	s.persistence = p
	s.data = p.Data

	cryptoService, err := crypto.New(p.Data)
	if err != nil {
		return nil, fmt.Errorf("initializing crypto service: %v", err)
	}

	s.crypto = cryptoService

	return s, nil
}

type Service struct {
	data        *data.Service
	persistence *persistence.Service
	crypto      *crypto.Service
	follower    *follower.Service
	chat        *chat.Service
}
