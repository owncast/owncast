package inbox

import (
	"context"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/federatedserversrepository"
	log "github.com/sirupsen/logrus"
)

func handleCreateRequest(c context.Context, activity vocab.ActivityStreamsCreate) error {
	// Get the object being created
	objectProp := activity.GetActivityStreamsObject()
	if objectProp == nil {
		return nil
	}

	// Check if it's a Note
	for i := 0; i < objectProp.Len(); i++ {
		if objectProp.At(i).GetActivityStreamsNote() != nil {
			note := objectProp.At(i).GetActivityStreamsNote()
			return handleNoteActivity(c, activity, note)
		}
	}

	return nil
}

func handleNoteActivity(_ context.Context, activity vocab.ActivityStreamsCreate, note vocab.ActivityStreamsNote) error {
	actorIRI, statusStr, unknownProps, ok := extractNoteStreamStatus(activity, note)
	if !ok {
		return nil
	}

	log.Debugf("Received Owncast stream status update from %s: %s", actorIRI, statusStr)

	repo := federatedserversrepository.Get()
	server, err := repo.GetFederatedServer(actorIRI)
	if err != nil || server == nil {
		log.Debugf("Ignoring Note activity from unfollowed server: %s", actorIRI)
		return nil
	}

	if server.Pending || server.FollowStatus == "rejected" || server.FollowStatus == "none" {
		log.Debugf("Ignoring Note activity from server we're not actively following: %s (status: %s)", actorIRI, server.FollowStatus)
		return nil
	}

	now := time.Now()
	server.LastStatusUpdate = &now

	switch statusStr {
	case config.APStreamStatusLive:
		server.IsOnline = true
		server.LastSeenOnline = &now
		applyLiveMetadataToServer(server, unknownProps)
		log.Infof("Federated server %s is now online", actorIRI)

	case config.APStreamStatusOffline:
		server.IsOnline = false
		server.StreamTitle = nil
		server.StreamDescription = nil
		server.ThumbnailURL = nil
		log.Infof("Federated server %s is now offline", actorIRI)
	}

	update := &models.FederatedStreamUpdate{
		Title:        server.StreamTitle,
		Description:  server.StreamDescription,
		ThumbnailURL: server.ThumbnailURL,
		Tags:         server.Tags,
	}

	err = repo.UpdateServerStatus(server.IRI, server.IsOnline, update)
	if err != nil {
		log.Errorf("Failed to update federated server status %s: %v", actorIRI, err)
	}

	return nil
}

// extractNoteStreamStatus extracts the actor IRI and stream status from a Note activity.
// Returns the actorIRI, status string, unknown properties map, and whether extraction succeeded.
func extractNoteStreamStatus(activity vocab.ActivityStreamsCreate, note vocab.ActivityStreamsNote) (string, string, map[string]interface{}, bool) {
	actorProp := activity.GetActivityStreamsActor()
	if actorProp == nil || actorProp.Len() == 0 {
		return "", "", nil, false
	}

	if actorProp.At(0).GetIRI() == nil {
		return "", "", nil, false
	}

	actorIRI := actorProp.At(0).GetIRI().String()

	unknownProps := note.GetUnknownProperties()
	streamStatus, hasStreamStatus := unknownProps[config.APOwncastNamespaceStreamStatus]
	if !hasStreamStatus {
		return "", "", nil, false
	}

	statusStr, ok := streamStatus.(string)
	if !ok {
		return "", "", nil, false
	}

	return actorIRI, statusStr, unknownProps, true
}

// applyLiveMetadataToServer applies parsed Owncast metadata to a server model when going live.
func applyLiveMetadataToServer(server *models.FederatedServer, unknownProps map[string]interface{}) {
	extractedMetadata := apmodels.ParseOwncastMetadata(unknownProps)
	if extractedMetadata.StreamTitle != "" {
		server.StreamTitle = &extractedMetadata.StreamTitle
	}
	if extractedMetadata.StreamDescription != "" {
		server.StreamDescription = &extractedMetadata.StreamDescription
	}
	if extractedMetadata.ThumbnailURL != "" {
		server.ThumbnailURL = &extractedMetadata.ThumbnailURL
	}
	if extractedMetadata.LogoURL != "" {
		server.LogoURL = &extractedMetadata.LogoURL
	}
	if len(extractedMetadata.Tags) > 0 {
		server.Tags = extractedMetadata.Tags
	}
	if extractedMetadata.ServerName != "" {
		server.Name = &extractedMetadata.ServerName
	}
}
