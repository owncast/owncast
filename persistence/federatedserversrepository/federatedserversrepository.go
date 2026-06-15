package federatedserversrepository

import (
	"time"

	"github.com/owncast/owncast/models"
)

// FederatedServersRepository defines the interface for federated servers repository operations.
type FederatedServersRepository interface {
	GetFederatedServers() ([]models.FederatedServer, error)
	GetFederatedServer(iri string) (*models.FederatedServer, error)
	AddFederatedServer(iri, name, logoURL string, followedAt time.Time, pending bool, username, followStatus string) error
	UpdateServerStatus(iri string, isOnline bool, metadata *models.FederatedStreamUpdate) error
	RemoveFederatedServer(id int64) error
	RemoveFederatedServerByIRI(iri string) error
	UpdateFollowStatus(iri, followStatus string, pending bool, acceptedAt, rejectedAt *time.Time) error
	UpdateServerMetadata(iri, name, displayName, summary, logoURL string) error
	GetPendingFederatedServers() ([]models.FederatedServer, error)
}
