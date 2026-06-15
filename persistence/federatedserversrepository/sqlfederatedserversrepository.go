package federatedserversrepository

import (
	"context"
	"time"

	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/datastore"
)

// SqlFederatedServersRepository is a SQL implementation of the FederatedServersRepository interface.
type SqlFederatedServersRepository struct {
	datastore *datastore.Datastore
}

// temporaryGlobalInstance is set once during application startup so
// helper code that has not yet been migrated to the dependency-injection
// pattern can still reach the federated-servers repository. Get returns
// nil until SetGlobalInstance has been called.
var temporaryGlobalInstance FederatedServersRepository

// SetGlobalInstance registers the application's single
// FederatedServersRepository for Get to return. Called from main.go
// after constructing the repository.
func SetGlobalInstance(r FederatedServersRepository) {
	temporaryGlobalInstance = r
}

// Get returns the global FederatedServersRepository registered with
// SetGlobalInstance. Returns nil until startup has wired one in.
func Get() FederatedServersRepository {
	return temporaryGlobalInstance
}

// New will create a new instance of the FederatedServersRepository.
func New(datastore *datastore.Datastore) FederatedServersRepository {
	return &SqlFederatedServersRepository{
		datastore: datastore,
	}
}

// GetFederatedServers returns all federated servers we are following.
func (r *SqlFederatedServersRepository) GetFederatedServers() ([]models.FederatedServer, error) {
	queries := db.New(r.datastore.DB)
	dbServers, err := queries.GetFederatedServers(context.Background())
	if err != nil {
		return nil, err
	}

	var servers []models.FederatedServer
	for _, dbServer := range dbServers {
		var apiServer models.FederatedServer
		apiServer.FromDatabaseModel(dbServer)
		servers = append(servers, apiServer)
	}

	return servers, nil
}

// GetFederatedServer returns a specific federated server by IRI.
func (r *SqlFederatedServersRepository) GetFederatedServer(iri string) (*models.FederatedServer, error) {
	queries := db.New(r.datastore.DB)
	server, err := queries.GetFederatedServer(context.Background(), iri)
	if err != nil {
		return nil, err
	}

	var apiServer models.FederatedServer
	apiServer.FromDatabaseModel(server)
	return &apiServer, nil
}

// AddFederatedServer adds a new federated server to follow.
func (r *SqlFederatedServersRepository) AddFederatedServer(iri, name, logoURL string, followedAt time.Time, pending bool, username, followStatus string) error {
	queries := db.New(r.datastore.DB)

	params := db.AddFederatedServerParams{
		Iri:          iri,
		Name:         models.PointerToNullString(&name),
		LogoUrl:      models.PointerToNullString(&logoURL),
		FollowedAt:   models.TimeToNullTime(followedAt),
		Pending:      models.BoolToNullBool(pending),
		Username:     models.PointerToNullString(&username),
		FollowStatus: models.PointerToNullString(&followStatus),
	}

	return queries.AddFederatedServer(context.Background(), params)
}

// UpdateServerStatus updates a federated server's online status and metadata.
func (r *SqlFederatedServersRepository) UpdateServerStatus(iri string, isOnline bool, metadata *models.FederatedStreamUpdate) error {
	queries := db.New(r.datastore.DB)
	now := time.Now()

	if isOnline && metadata != nil {
		// Server came online with stream metadata
		params := db.UpdateFederatedServerStatusParams{
			IsOnline:          models.BoolToNullBool(isOnline),
			StreamTitle:       models.PointerToNullString(metadata.Title),
			StreamDescription: models.PointerToNullString(metadata.Description),
			StreamTags:        models.StringSliceToNullString(metadata.Tags),
			ThumbnailUrl:      models.PointerToNullString(metadata.ThumbnailURL),
			LastStatusUpdate:  models.TimeToNullTime(now),
			Iri:               iri,
		}
		return queries.UpdateFederatedServerStatus(context.Background(), params)
	} else {
		// Server went offline or just status update without metadata
		var lastSeenOnline time.Time
		if isOnline {
			lastSeenOnline = now
		} else {
			// Don't update last seen online when going offline
			lastSeenOnline = time.Time{}
		}

		params := db.UpdateFederatedServerOnlineStatusParams{
			IsOnline:         models.BoolToNullBool(isOnline),
			LastSeenOnline:   models.TimeToNullTime(lastSeenOnline),
			LastStatusUpdate: models.TimeToNullTime(now),
			Iri:              iri,
		}
		return queries.UpdateFederatedServerOnlineStatus(context.Background(), params)
	}
}

// RemoveFederatedServer removes a federated server by ID.
func (r *SqlFederatedServersRepository) RemoveFederatedServer(id int64) error {
	queries := db.New(r.datastore.DB)
	return queries.RemoveFederatedServer(context.Background(), id)
}

// RemoveFederatedServerByIRI removes a federated server by IRI.
func (r *SqlFederatedServersRepository) RemoveFederatedServerByIRI(iri string) error {
	// First get the server to find its ID
	server, err := r.GetFederatedServer(iri)
	if err != nil {
		return err
	}
	if server == nil {
		return nil // Server doesn't exist, nothing to remove
	}
	return r.RemoveFederatedServer(server.ID)
}

// UpdateFollowStatus updates the follow status of a federated server.
func (r *SqlFederatedServersRepository) UpdateFollowStatus(iri, followStatus string, pending bool, acceptedAt, rejectedAt *time.Time) error {
	queries := db.New(r.datastore.DB)

	params := db.UpdateFederatedServerFollowStatusParams{
		FollowStatus: models.PointerToNullString(&followStatus),
		Pending:      models.BoolToNullBool(pending),
		AcceptedAt:   models.PointerToNullTime(acceptedAt),
		RejectedAt:   models.PointerToNullTime(rejectedAt),
		Iri:          iri,
	}

	return queries.UpdateFederatedServerFollowStatus(context.Background(), params)
}

// UpdateServerMetadata updates the metadata of a federated server.
func (r *SqlFederatedServersRepository) UpdateServerMetadata(iri, name, displayName, summary, logoURL string) error {
	queries := db.New(r.datastore.DB)

	params := db.UpdateFederatedServerMetadataParams{
		Name:        models.PointerToNullString(&name),
		DisplayName: models.PointerToNullString(&displayName),
		Summary:     models.PointerToNullString(&summary),
		LogoUrl:     models.PointerToNullString(&logoURL),
		Iri:         iri,
	}

	return queries.UpdateFederatedServerMetadata(context.Background(), params)
}

// GetPendingFederatedServers returns all federated servers with pending follow status.
func (r *SqlFederatedServersRepository) GetPendingFederatedServers() ([]models.FederatedServer, error) {
	queries := db.New(r.datastore.DB)
	dbServers, err := queries.GetPendingFederatedServers(context.Background())
	if err != nil {
		return nil, err
	}

	var servers []models.FederatedServer
	for _, dbServer := range dbServers {
		var apiServer models.FederatedServer
		apiServer.FromDatabaseModel(dbServer)
		servers = append(servers, apiServer)
	}

	return servers, nil
}
