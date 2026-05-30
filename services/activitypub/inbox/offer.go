package inbox

import (
	"context"

	"github.com/go-fed/activity/streams/vocab"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/federatedserversrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
)

func (s *Service) handleOfferInboxRequest(c context.Context, activity vocab.ActivityStreamsOffer) error {
	actorIRI, valid := extractActorIRI(activity)
	if !valid {
		return nil
	}

	if !isValidOwncastStreamOffer(activity) {
		return nil
	}

	log.Debugf("Received Owncast stream ping from %s", actorIRI)

	repo := federatedserversrepository.Get()
	if !shouldProcessOfferFromServer(repo, actorIRI) {
		return nil
	}

	unknownProps := activity.GetUnknownProperties()
	extractedMetadata := apmodels.ParseOwncastMetadata(unknownProps)
	update := buildStreamUpdateFromMetadata(extractedMetadata)

	// Update the server status to online
	err := repo.UpdateServerStatus(actorIRI, true, update)
	if err != nil {
		log.Errorf("Failed to update federated server status from Offer: %v", err)
		return err
	}

	log.Debugf("Updated federated server %s status from stream ping", actorIRI)
	return nil
}

func extractActorIRI(activity vocab.ActivityStreamsOffer) (string, bool) {
	actorProp := activity.GetActivityStreamsActor()
	if actorProp == nil || actorProp.Len() == 0 {
		return "", false
	}

	if actorProp.At(0).GetIRI() != nil {
		return actorProp.At(0).GetIRI().String(), true
	}

	return "", false
}

func isValidOwncastStreamOffer(activity vocab.ActivityStreamsOffer) bool {
	unknownProps := activity.GetUnknownProperties()
	streamStatus, hasStreamStatus := unknownProps[config.APOwncastNamespaceStreamStatus]

	if !hasStreamStatus {
		return false
	}

	statusStr, ok := streamStatus.(string)
	return ok && statusStr == "live"
}

func shouldProcessOfferFromServer(repo federatedserversrepository.FederatedServersRepository, actorIRI string) bool {
	server, err := repo.GetFederatedServer(actorIRI)
	if err != nil || server == nil {
		log.Debugf("Ignoring Offer activity from unfollowed server: %s", actorIRI)
		return false
	}

	if server.Pending || server.FollowStatus == "rejected" || server.FollowStatus == "none" {
		log.Debugf("Ignoring Offer activity from server we're not actively following: %s (status: %s)", actorIRI, server.FollowStatus)
		return false
	}

	return true
}

func buildStreamUpdateFromMetadata(metadata *apmodels.OwncastMetadata) *models.FederatedStreamUpdate {
	var streamTitle, streamDescription, thumbnailURL *string
	var tags []string

	if metadata.StreamTitle != "" {
		streamTitle = &metadata.StreamTitle
	}
	if metadata.StreamDescription != "" {
		streamDescription = &metadata.StreamDescription
	}
	if metadata.ThumbnailURL != "" {
		thumbnailURL = &metadata.ThumbnailURL
	}
	if len(metadata.Tags) > 0 {
		tags = metadata.Tags
	}

	return &models.FederatedStreamUpdate{
		Title:        streamTitle,
		Description:  streamDescription,
		ThumbnailURL: thumbnailURL,
		Tags:         tags,
	}
}
