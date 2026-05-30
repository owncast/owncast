package inbox

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/persistence/federatedserversrepository"
	log "github.com/sirupsen/logrus"
)

func (s *Service) handleRejectInboxRequest(c context.Context, activity vocab.ActivityStreamsReject) error {
	// Get the actor who sent the reject
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

	log.Debugf("Received Reject activity from %s", actorIRI)

	// Check what object is being rejected (usually our Follow request)
	objectProp := activity.GetActivityStreamsObject()
	if objectProp == nil || objectProp.Len() == 0 {
		log.Debugf("Reject activity has no object, ignoring")
		return nil
	}

	// Check if this is rejecting a Follow
	for iter := objectProp.Begin(); iter != objectProp.End(); iter = iter.Next() {
		if iter.IsActivityStreamsFollow() {
			// This is a Reject of a Follow request
			log.Debugf("Received Reject for Follow request from %s", actorIRI)

			// Extract the server URL from the actor IRI
			parsedIRI, err := url.Parse(actorIRI)
			if err != nil {
				log.Errorf("Failed to parse actor IRI %s: %v", actorIRI, err)
				return nil
			}

			// Construct the server URL (base URL without the federation path)
			serverURL := fmt.Sprintf("%s://%s", parsedIRI.Scheme, parsedIRI.Host)

			// Update the follow status in the database
			repo := federatedserversrepository.Get()

			// Get the existing server record
			server, err := repo.GetFederatedServer(serverURL)
			if err != nil || server == nil {
				log.Debugf("No pending follow found for %s", serverURL)
				return nil
			}

			// Update follow status to rejected
			rejectedAt := time.Now()
			err = repo.UpdateFollowStatus(serverURL, "rejected", false, nil, &rejectedAt)
			if err != nil {
				log.Errorf("Failed to update follow status for %s: %v", serverURL, err)
				return nil
			}

			log.Infof("Follow request to %s has been rejected", serverURL)
		}
	}

	return nil
}
