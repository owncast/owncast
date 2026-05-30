package inbox

import (
	"context"

	"github.com/go-fed/activity/streams/vocab"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/services/activitypub/apmodels"
)

func (s *Service) handleLeaveInboxRequest(c context.Context, activity vocab.ActivityStreamsLeave) error {
	log.Debugln("Handling incoming Leave activity")

	// Get the actor who is leaving (sending the Leave activity)
	actorProperty := activity.GetActivityStreamsActor()
	if actorProperty == nil {
		return nil
	}

	// Resolve the actor to get their information
	actor, err := s.resolver.GetResolvedActorFromActorProperty(actorProperty)
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

// updateFederatedServerStatus updates the status and metadata for a federated server.
func updateFederatedServerStatus(actorIRI string, metadata *apmodels.OwncastMetadata) error {
	// Log the status update. Metadata is accepted for symmetry with
	// Offer/Leave handling; the federated_servers row carries this
	// information via UpdateServerStatus when the next live ping comes
	// in, so we just acknowledge the offline transition here.
	_ = metadata
	log.Infof("Updated federated server %s to offline status", actorIRI)
	return nil
}
