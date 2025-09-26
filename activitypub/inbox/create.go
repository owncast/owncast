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

func handleNoteActivity(c context.Context, activity vocab.ActivityStreamsCreate, note vocab.ActivityStreamsNote) error {
	// Get the actor who created the note
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

	// Check for Owncast custom properties in the note
	unknownProps := note.GetUnknownProperties()
	streamStatus, hasStreamStatus := unknownProps[config.APOwncastNamespaceStreamStatus]

	if !hasStreamStatus {
		// Not an Owncast stream status update, ignore
		return nil
	}

	statusStr, ok := streamStatus.(string)
	if !ok {
		return nil
	}

	log.Debugf("Received Owncast stream status update from %s: %s", actorIRI, statusStr)

	// Get or create the federated server
	repo := federatedserversrepository.Get()
	server, err := repo.GetFederatedServer(actorIRI)
	if err != nil || server == nil {
		// Server doesn't exist in our database - we're not following them
		log.Debugf("Ignoring Note activity from unfollowed server: %s", actorIRI)
		return nil
	}

	// Check if we're actually following this server (not pending or rejected)
	if server.Pending || server.FollowStatus == "rejected" || server.FollowStatus == "none" {
		log.Debugf("Ignoring Note activity from server we're not actively following: %s (status: %s)", actorIRI, server.FollowStatus)
		return nil
	}

	now := time.Now()
	server.LastStatusUpdate = &now

	switch statusStr {
	case "live":
		server.IsOnline = true
		server.LastSeenOnline = &now

		// Extract stream metadata from unknown properties using utility function
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

		log.Infof("Federated server %s is now online", actorIRI)

	case "offline":
		server.IsOnline = false
		// Clear stream metadata when offline
		server.StreamTitle = nil
		server.StreamDescription = nil
		server.ThumbnailURL = nil

		log.Infof("Federated server %s is now offline", actorIRI)
	}

	// Update the server status (we know server exists and we're following them)
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
