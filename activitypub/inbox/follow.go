package inbox

import (
	"context"
	"fmt"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/webhooks"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/federatedserversrepository"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

func handleFollowInboxRequest(c context.Context, activity vocab.ActivityStreamsFollow) error {
	configRepository := configrepository.Get()
	followersRepo := followersrepository.Get()

	follow, err := resolvers.MakeFollowRequest(c, activity)
	if err != nil {
		log.Errorln("unable to create follow inbox request", err)
		return err
	}

	if follow == nil {
		return fmt.Errorf("unable to handle request")
	}

	approved := !configRepository.GetFederationIsPrivate()

	followRequest := *follow

	if err := followersRepo.Add(followRequest, approved); err != nil {
		log.Errorln("unable to save follow request", err)
		return err
	}

	// If this is an Owncast server, also update/create the federated server record
	if followRequest.IsOwncastServer {
		if err := handleOwncastServerFollow(activity, followRequest); err != nil {
			log.Errorf("Failed to update federated server from follow request: %v", err)
			// Don't return error as the follow was successful, just log the federated server update failure
		}
	}

	localAccountName := configRepository.GetDefaultFederationUsername()

	objectIRI, err := apmodels.GetIRIStringFromObjectProperty(activity.GetActivityStreamsObject())
	if err != nil {
		return errors.Wrap(err, "follow activity is missing object IRI")
	}

	actorIRI, err := apmodels.GetIRIStringFromActorProperty(activity.GetActivityStreamsActor())
	if err != nil {
		return errors.Wrap(err, "follow activity is missing actor IRI")
	}

	actorReference := activity.GetActivityStreamsActor()

	if approved {
		if err := requests.SendFollowAccept(follow.Inbox, activity, localAccountName, false); err != nil {
			log.Errorln("unable to send follow accept", err)
			return err
		}
		go webhooks.SendFediverseEngagementFollowEvent(actorIRI)
	}

	// If this request is approved and we have not previously sent an action to
	// chat due to a previous follow request, then do so.
	hasPreviouslyhandled := true // Default so we don't send anything if it fails.
	if approved {
		hasPreviouslyhandled, err = persistence.HasPreviouslyHandledInboundActivity(objectIRI, actorIRI, events.FediverseEngagementFollow)
		if err != nil {
			log.Errorln("error checking for previously handled follow activity", err)
		}
	}

	// Save this follow action to our activities table.
	if err := persistence.SaveInboundFediverseActivity(objectIRI, actorIRI, events.FediverseEngagementFollow, time.Now()); err != nil {
		return errors.Wrap(err, "unable to save inbound share/re-post activity")
	}

	// Send action to chat if it has not been previously handled.
	if !hasPreviouslyhandled {
		return handleEngagementActivity(events.FediverseEngagementFollow, false, actorReference, events.FediverseEngagementFollow)
	}

	return nil
}

func handleUnfollowRequest(c context.Context, activity vocab.ActivityStreamsUndo) error {
	request := resolvers.MakeUnFollowRequest(c, activity)
	if request == nil {
		log.Errorf("unable to handle unfollow request")
		return errors.New("unable to handle unfollow request")
	}

	unfollowRequest := *request
	log.Traceln("unfollow request:", unfollowRequest)

	followersRepo := followersrepository.Get()
	return followersRepo.Remove(unfollowRequest)
}

func handleOwncastServerFollow(activity vocab.ActivityStreamsFollow, followRequest apmodels.ActivityPubActor) error {
	// Extract Owncast metadata from the follow request using shared utility.
	unknownProps := activity.GetUnknownProperties()
	metadata := apmodels.ParseOwncastMetadata(unknownProps)

	if !metadata.IsOwncastServer {
		return nil // Not an Owncast server, nothing to do.
	}

	repo := federatedserversrepository.Get()
	actorIRI := followRequest.ActorIri.String()

	// Get or create the federated server record.
	server, err := getOrCreateFederatedServer(repo, actorIRI, metadata, &followRequest)
	if err != nil {
		return err
	}

	// Update server with metadata from follow request.
	now := time.Now()
	server.LastStatusUpdate = &now
	server.FollowedAt = &now

	streamUpdate, hasUpdates := buildFollowStreamUpdate(metadata, server)

	isOnline := metadata.StreamStatus == config.APStreamStatusLive
	if isOnline {
		server.LastSeenOnline = &now
	}

	if hasUpdates || isOnline {
		err = repo.UpdateServerStatus(actorIRI, isOnline, streamUpdate)
		if err != nil {
			log.Errorf("Failed to update federated server status: %v", err)
		}
	}

	log.Infof("Updated federated server %s from follow request (online: %v)", actorIRI, isOnline)
	return nil
}

// getOrCreateFederatedServer retrieves an existing server or creates a new one.
func getOrCreateFederatedServer(repo federatedserversrepository.FederatedServersRepository, actorIRI string, metadata *apmodels.OwncastMetadata, followRequest *apmodels.ActivityPubActor) (*models.FederatedServer, error) {
	server, err := repo.GetFederatedServer(actorIRI)
	if err == nil {
		return server, nil
	}

	// Server doesn't exist, create it.
	var name string
	if metadata.ServerName != "" {
		name = metadata.ServerName
	} else if followRequest.Name != "" {
		name = followRequest.Name
	}

	var logoURL string
	if metadata.LogoURL != "" {
		logoURL = metadata.LogoURL
	} else if followRequest.Image != nil {
		logoURL = followRequest.Image.String()
	}

	err = repo.AddFederatedServer(actorIRI, name, logoURL, time.Now(), false, "", "accepted")
	if err != nil {
		return nil, err
	}

	return repo.GetFederatedServer(actorIRI)
}

// buildFollowStreamUpdate creates a FederatedStreamUpdate from parsed metadata and updates the server model.
func buildFollowStreamUpdate(metadata *apmodels.OwncastMetadata, server *models.FederatedServer) (*models.FederatedStreamUpdate, bool) {
	streamUpdate := &models.FederatedStreamUpdate{}
	hasUpdates := false

	if metadata.StreamTitle != "" {
		streamUpdate.Title = &metadata.StreamTitle
		hasUpdates = true
	}
	if metadata.StreamDescription != "" {
		streamUpdate.Description = &metadata.StreamDescription
		hasUpdates = true
	}
	if metadata.ThumbnailURL != "" {
		streamUpdate.ThumbnailURL = &metadata.ThumbnailURL
		hasUpdates = true
	}
	if len(metadata.Tags) > 0 {
		streamUpdate.Tags = metadata.Tags
		hasUpdates = true
	}
	if metadata.ServerName != "" {
		server.Name = &metadata.ServerName
	}
	if metadata.LogoURL != "" {
		server.LogoURL = &metadata.LogoURL
	}

	return streamUpdate, hasUpdates
}
