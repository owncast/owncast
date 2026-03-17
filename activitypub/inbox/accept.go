package inbox

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/persistence/federatedserversrepository"
	log "github.com/sirupsen/logrus"
)

func handleAcceptInboxRequest(c context.Context, activity vocab.ActivityStreamsAccept) error {
	// Get the actor who sent the accept
	actorProp := activity.GetActivityStreamsActor()
	if actorProp == nil || actorProp.Len() == 0 {
		return nil
	}

	var actorIRI string
	if actorProp.At(0).GetIRI() != nil {
		actorIRI = actorProp.At(0).GetIRI().String()
	} else {
		return nil
	}

	log.Debugf("Received Accept activity from %s", actorIRI)

	// Check what object is being accepted (usually our Follow request)
	objectProp := activity.GetActivityStreamsObject()
	if objectProp == nil || objectProp.Len() == 0 {
		log.Debugf("Accept activity has no object, ignoring")
		return nil
	}

	// Check if this is accepting a Follow
	for iter := objectProp.Begin(); iter != objectProp.End(); iter = iter.Next() {
		if iter.IsActivityStreamsFollow() {
			handleFollowAccepted(actorIRI)
		}
	}

	// Check for Owncast metadata in the Accept activity using shared utility
	unknownProps := activity.GetUnknownProperties()
	metadata := apmodels.ParseOwncastMetadata(unknownProps)

	if metadata.IsOwncastServer {
		log.Debugf("Accept activity from %s contains Owncast metadata", actorIRI)
	}

	return nil
}

func handleFollowAccepted(actorIRI string) {
	log.Debugf("Received Accept for Follow request from %s", actorIRI)

	// Extract the server URL from the actor IRI
	parsedIRI, err := url.Parse(actorIRI)
	if err != nil {
		log.Errorf("Failed to parse actor IRI %s: %v", actorIRI, err)
		return
	}

	// Construct the server URL (base URL without the federation path)
	serverURL := fmt.Sprintf("%s://%s", parsedIRI.Scheme, parsedIRI.Host)

	// Update the follow status in the database
	repo := federatedserversrepository.Get()

	// Get the existing server record
	server, err := repo.GetFederatedServer(serverURL)
	if err != nil || server == nil {
		log.Debugf("No pending follow found for %s", serverURL)
		return
	}

	// Update follow status to accepted
	acceptedAt := time.Now()
	err = repo.UpdateFollowStatus(serverURL, "accepted", false, &acceptedAt, nil)
	if err != nil {
		log.Errorf("Failed to update follow status for %s: %v", serverURL, err)
		return
	}

	// Try to fetch and update server metadata from the actor
	updateServerMetadataFromActor(repo, serverURL, actorIRI)

	log.Infof("Follow request to %s has been accepted", serverURL)
}

func updateServerMetadataFromActor(repo federatedserversrepository.FederatedServersRepository, serverURL, actorIRI string) {
	actorData, err := resolvers.GetResolvedActorFromIRI(actorIRI)
	if err != nil {
		return
	}

	displayName := actorData.Name
	name := actorData.Username
	// For summary, use the display name as a fallback.
	summary := actorData.Name

	var logoURL string
	if actorData.Image != nil {
		logoURL = actorData.Image.String()
	}

	// Update server metadata
	err = repo.UpdateServerMetadata(serverURL, name, displayName, summary, logoURL)
	if err != nil {
		log.Errorf("Failed to update server metadata for %s: %v", serverURL, err)
	}
}
