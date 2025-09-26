package federatedserversrepository

import (
	"github.com/owncast/owncast/models"
)

// FederatedServersRepository defines the interface for federated servers repository operations.
type FederatedServersRepository interface {
	GetFederatedServers() ([]models.FederatedServer, error)
	GetFederatedServer(iri string) (*models.FederatedServer, error)
	AddFederatedServer(iri, name, logoURL string) error
	UpdateServerStatus(iri string, isOnline bool, metadata *models.FederatedStreamUpdate) error
	RemoveFederatedServer(id int32) error
}
