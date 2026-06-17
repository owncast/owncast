package inbox

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/persistence/federatedserversrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
)

func (s *Service) handleAcceptInboxRequest(c context.Context, activity vocab.ActivityStreamsAccept) error {
	actorProp := activity.GetActivityStreamsActor()
	if actorProp == nil || actorProp.Len() == 0 {
		return nil
	}

	if actorProp.At(0).GetIRI() == nil {
		return nil
	}
	actorIRI := actorProp.At(0).GetIRI().String()

	log.Debugf("Received Accept activity from %s", actorIRI)

	objectProp := activity.GetActivityStreamsObject()
	if objectProp == nil || objectProp.Len() == 0 {
		log.Debugf("Accept activity has no object, ignoring")
		return nil
	}

	for iter := objectProp.Begin(); iter != objectProp.End(); iter = iter.Next() {
		if iter.IsActivityStreamsFollow() {
			s.markFederatedServerAccepted(actorIRI)
		}
	}

	unknownProps := activity.GetUnknownProperties()
	metadata := apmodels.ParseOwncastMetadata(unknownProps)
	if metadata.IsOwncastServer {
		log.Debugf("Accept activity from %s contains Owncast metadata", actorIRI)
	}

	return nil
}

// markFederatedServerAccepted handles the bookkeeping side of receiving
// an Accept-of-Follow: it transitions our pending follow record for the
// remote Owncast server into the accepted state and tops up the cached
// metadata from the resolved actor.
func (s *Service) markFederatedServerAccepted(actorIRI string) {
	log.Debugf("Received Accept for Follow request from %s", actorIRI)

	parsedIRI, err := url.Parse(actorIRI)
	if err != nil {
		log.Errorf("Failed to parse actor IRI %s: %v", actorIRI, err)
		return
	}

	serverURL := fmt.Sprintf("%s://%s", parsedIRI.Scheme, parsedIRI.Host)

	repo := federatedserversrepository.Get()
	if repo == nil {
		log.Errorln("Federated servers repository not initialised; cannot mark Accept")
		return
	}

	server, err := repo.GetFederatedServer(serverURL)
	if err != nil || server == nil {
		log.Debugf("No pending follow found for %s", serverURL)
		return
	}

	acceptedAt := time.Now()
	if err := repo.UpdateFollowStatus(serverURL, "accepted", false, &acceptedAt, nil); err != nil {
		log.Errorf("Failed to update follow status for %s: %v", serverURL, err)
		return
	}

	if actorData, err := s.resolver.GetResolvedActorFromIRI(actorIRI); err == nil {
		var logoURL string
		if actorData.Image != nil {
			logoURL = truncateMetadata(actorData.Image.String(), maxMetadataURLLen)
		}
		// Clamp the attacker-controlled remote actor fields before storing.
		// Summary falls back to the display name; Owncast actors don't
		// expose a separate summary field on the Person object.
		name := truncateMetadata(actorData.Username, maxServerNameLen)
		displayName := truncateMetadata(actorData.Name, maxServerNameLen)
		if err := repo.UpdateServerMetadata(serverURL, name, displayName, displayName, logoURL); err != nil {
			log.Errorf("Failed to update server metadata for %s: %v", serverURL, err)
		}
	}

	log.Infof("Follow request to %s has been accepted", serverURL)
}
