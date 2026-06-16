package inbox

import (
	"context"
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/persistence/federatedserversrepository"
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

// updateFederatedServerStatus marks the federated server identified by the
// given actor IRI as offline. The Leave activity is how a peer announces its
// stream has ended, so it is the mirror of the Offer (go-live) handler: it
// flips the directory entry offline immediately rather than waiting for the
// staleness sweep. Metadata is accepted for symmetry with the Offer handler
// but an offline transition carries no stream fields to store.
func updateFederatedServerStatus(actorIRI string, metadata *apmodels.OwncastMetadata) error {
	_ = metadata

	// Federated server records are keyed by the base server URL
	// (scheme://host); the Leave's actor is the full actor IRI.
	serverURL := serverURLFromActorIRI(actorIRI)
	if serverURL == "" {
		log.Debugf("Could not derive server URL from actor IRI %s, ignoring Leave", actorIRI)
		return nil
	}

	repo := federatedserversrepository.Get()
	if repo == nil {
		// No repository wired up (e.g. in unit tests); nothing to update.
		return nil
	}

	// Only act on servers we are actively following, matching the Offer path.
	if !shouldProcessOfferFromServer(repo, serverURL) {
		return nil
	}

	if err := repo.UpdateServerStatus(serverURL, false, nil); err != nil {
		return fmt.Errorf("failed to mark federated server %s offline: %w", serverURL, err)
	}

	log.Debugf("Marked federated server %s offline from Leave activity", serverURL)
	return nil
}
