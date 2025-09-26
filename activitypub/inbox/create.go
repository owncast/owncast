package inbox

import (
	"context"
	"net/url"
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
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

	// Parse actor IRI to get host
	actorURL, err := url.Parse(actorIRI)
	if err != nil {
		return nil
	}

	// Check for Owncast custom properties in the note
	unknownProps := note.GetUnknownProperties()
	streamStatus, hasStreamStatus := unknownProps["https://owncast.online/ns#streamStatus"]

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
	if err != nil {
		// Server doesn't exist, create it
		server = &models.FederatedServer{
			IRI:      actorIRI,
			IsOnline: false,
			AddedAt:  time.Now(),
		}

		// Try to get server name from actor
		if nameProp := note.GetActivityStreamsAttributedTo(); nameProp != nil {
			// This is a simplified approach - in reality you'd need to resolve the actor
			// to get the name, but for now we'll use the hostname
			server.Name = &actorURL.Host
		}
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

	// Save or update the server
	if server.ID == 0 {
		var name string
		if server.Name != nil {
			name = *server.Name
		}
		err = repo.AddFederatedServer(server.IRI, name, "")
		if err != nil {
			log.Errorf("Failed to add federated server %s: %v", actorIRI, err)
			return nil
		}
		// Reload to get the ID
		server, _ = repo.GetFederatedServer(actorIRI)
	}

	if server != nil {
		update := &models.FederatedStreamUpdate{
			Title:        server.StreamTitle,
			Description:  server.StreamDescription,
			ThumbnailURL: server.ThumbnailURL,
			Tags:         server.Tags,
		}

		err = repo.UpdateServerStatus(server.IRI, server.IsOnline, update)
		if err != nil {
			log.Errorf("Failed to update federated server %s: %v", actorIRI, err)
		}
	}

	return nil
}
