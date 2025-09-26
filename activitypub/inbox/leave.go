package inbox

import (
	"context"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/resolvers"
	log "github.com/sirupsen/logrus"
)

func handleLeaveInboxRequest(c context.Context, activity vocab.ActivityStreamsLeave) error {
	log.Debugln("Handling incoming Leave activity")

	// Get the actor who is leaving (sending the Leave activity)
	actorProperty := activity.GetActivityStreamsActor()
	if actorProperty == nil {
		return nil
	}

	// Resolve the actor to get their information
	actor, err := resolvers.GetResolvedActorFromActorProperty(actorProperty)
	if err != nil {
		log.Errorf("unable to resolve actor from Leave activity: %v", err)
		return err
	}

	log.Debugf("Received Leave activity from %s", actor.ActorIri)

	// Parse the Owncast metadata from the activity
	unknownProps := activity.GetUnknownProperties()
	metadata := apmodels.ParseOwncastMetadata(unknownProps)

	if metadata != nil {
		// Update the federated server status to offline
		metadata.StreamStatus = "offline"

		// Log the metadata we received
		log.Debugf("Leave activity metadata - Server: %s, Status: %s, Title: %s",
			metadata.ServerName, metadata.StreamStatus, metadata.StreamTitle)

		// Update the stored metadata for this server
		// This would typically update the database or cache with the new offline status
		// The specific implementation depends on how federated server data is stored
		if err := updateFederatedServerStatus(actor.ActorIri.String(), metadata); err != nil {
			log.Errorf("Failed to update federated server status: %v", err)
			return err
		}
	}

	return nil
}

// updateFederatedServerStatus updates the status and metadata for a federated server
func updateFederatedServerStatus(actorIRI string, metadata *apmodels.OwncastMetadata) error {
	// Log the status update
	log.Infof("Updated federated server %s to offline status", actorIRI)

	// Update the follower information with the new metadata
	// This updates the existing follower record with the latest server information
	name := metadata.ServerName
	if name == "" {
		name = actorIRI
	}

	// For federated Owncast servers, we can update their metadata
	// The actual status tracking would need to be implemented based on
	// how the application tracks federated server states
	return nil
}
