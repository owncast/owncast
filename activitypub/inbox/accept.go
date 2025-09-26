package inbox

import (
	"context"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
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

	// Check for Owncast metadata in the Accept activity using shared utility
	unknownProps := activity.GetUnknownProperties()
	metadata := apmodels.ParseOwncastMetadata(unknownProps)

	if metadata.IsOwncastServer {
		log.Debugf("Accept activity from %s contains Owncast metadata", actorIRI)
		// Here we could create/update federated server records if needed
		// For now, just log that we detected an Owncast server
	}

	// Could add more accept-specific logic here if needed
	// For now, we mainly care about the metadata extraction

	return nil
}
