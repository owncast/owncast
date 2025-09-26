package federatedserversrepository

import (
	"context"
	"time"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/models"
)

// SqlFederatedServersRepository is a SQL implementation of the FederatedServersRepository interface.
type SqlFederatedServersRepository struct {
	datastore *data.Datastore
}

// NOTE: This is temporary during the transition period.
var temporaryGlobalInstance FederatedServersRepository

// Get will return the federated servers repository.
func Get() FederatedServersRepository {
	if temporaryGlobalInstance == nil {
		i := New(data.GetDatastore())
		temporaryGlobalInstance = i
	}
	return temporaryGlobalInstance
}

// New will create a new instance of the FederatedServersRepository.
func New(datastore *data.Datastore) FederatedServersRepository {
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
func (r *SqlFederatedServersRepository) AddFederatedServer(iri, name, logoURL string) error {
	queries := db.New(r.datastore.DB)

	now := time.Now()
	followedAt := &now

	params := db.AddFederatedServerParams{
		Iri:        iri,
		Name:       models.PointerToNullString(&name),
		LogoUrl:    models.PointerToNullString(&logoURL),
		FollowedAt: models.PointerToNullTime(followedAt),
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
func (r *SqlFederatedServersRepository) RemoveFederatedServer(id int32) error {
	queries := db.New(r.datastore.DB)
	return queries.RemoveFederatedServer(context.Background(), id)
}
