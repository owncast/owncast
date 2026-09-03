package apmodels

import (
	"encoding/json"
	"errors"
	"path"
	"time"

	"code.superseriousbusiness.org/activity/streams"
	"code.superseriousbusiness.org/activity/streams/vocab"

	"github.com/owncast/owncast/models"
)

const (
	activityStreamsContext = "https://www.w3.org/ns/activitystreams"
	eventContext           = "https://w3id.org/fep/8a8e"
	schemaContext          = "https://schema.org"
)

// MakeScheduledEvent builds the public ActivityPub representation of one
// scheduled stream occurrence.
func (b *Builder) MakeScheduledEvent(event models.ScheduledEvent) (vocab.ActivityStreamsEvent, error) {
	actorIRI := b.MakeLocalIRIForAccount(b.configRepository.GetDefaultFederationUsername())
	eventIRI := b.MakeLocalIRIForResource(path.Join("event", event.ID))
	pageURL := b.MakeLocalURLForPath(path.Join("schedule", event.ID))
	if actorIRI == nil || eventIRI == nil || pageURL == nil {
		return nil, errors.New("unable to build scheduled event URLs")
	}

	result := streams.NewActivityStreamsEvent()

	id := streams.NewJSONLDIdProperty()
	id.Set(eventIRI)
	result.SetJSONLDId(id)

	name := streams.NewActivityStreamsNameProperty()
	name.AppendXMLSchemaString(event.Name)
	result.SetActivityStreamsName(name)

	if event.Description != "" {
		content := streams.NewActivityStreamsContentProperty()
		content.AppendXMLSchemaString(event.Description)
		result.SetActivityStreamsContent(content)
	}

	startTime := streams.NewActivityStreamsStartTimeProperty()
	startTime.Set(event.StartTime)
	result.SetActivityStreamsStartTime(startTime)

	endTime := streams.NewActivityStreamsEndTimeProperty()
	endTime.Set(event.StartTime.Add(time.Duration(event.DurationMinutes) * time.Minute))
	result.SetActivityStreamsEndTime(endTime)

	duration := streams.NewActivityStreamsDurationProperty()
	duration.Set(time.Duration(event.DurationMinutes) * time.Minute)
	result.SetActivityStreamsDuration(duration)

	attributedTo := streams.NewActivityStreamsAttributedToProperty()
	attributedTo.AppendIRI(actorIRI)
	result.SetActivityStreamsAttributedTo(attributedTo)

	url := streams.NewActivityStreamsUrlProperty()
	url.AppendIRI(pageURL)
	result.SetActivityStreamsUrl(url)

	if event.CreatedAt != nil {
		published := streams.NewActivityStreamsPublishedProperty()
		published.Set(*event.CreatedAt)
		result.SetActivityStreamsPublished(published)
	}
	if event.UpdatedAt != nil {
		updated := streams.NewActivityStreamsUpdatedProperty()
		updated.Set(*event.UpdatedAt)
		result.SetActivityStreamsUpdated(updated)
	}

	followersIRI := actorIRI.JoinPath("followers")
	to, cc := MakeAddressingToFollowers(followersIRI, !b.configRepository.GetFederationIsPrivate())
	result.SetActivityStreamsTo(to)
	result.SetActivityStreamsCc(cc)

	serverName := b.configRepository.GetServerName()
	if serverName == "" {
		serverName = "Owncast"
	}
	properties := result.GetUnknownProperties()
	properties["displayEndTime"] = true
	properties["eventStatus"] = "EventScheduled"
	properties["joinMode"] = "none"
	properties["location"] = map[string]interface{}{
		"type": "VirtualLocation",
		"name": serverName,
		"url":  pageURL.String(),
	}
	properties["organizers"] = map[string]interface{}{
		"type":       "OrganizersCollection",
		"totalItems": 1,
		"items":      []string{actorIRI.String()},
	}
	if event.Timezone != "" {
		properties["timezone"] = event.Timezone
	}

	logoURL := b.MakeLocalIRIforLogo()
	if logoURL != nil {
		attachments := streams.NewActivityStreamsAttachmentProperty()
		image := streams.NewActivityStreamsImage()
		imageURL := streams.NewActivityStreamsUrlProperty()
		imageURL.AppendIRI(logoURL)
		image.SetActivityStreamsUrl(imageURL)
		name := streams.NewActivityStreamsNameProperty()
		name.AppendXMLSchemaString(serverName + " logo")
		image.SetActivityStreamsName(name)
		mediaType := streams.NewActivityStreamsMediaTypeProperty()
		mediaType.Set(b.GetLogoType())
		image.SetActivityStreamsMediaType(mediaType)
		attachments.AppendActivityStreamsImage(image)
		result.SetActivityStreamsAttachment(attachments)
	}

	return result, nil
}

// SerializeEvent adds the extension contexts required by FEP-8a8e and
// schema.org to an Event or an activity containing one.
func SerializeEvent(value vocab.Type) ([]byte, error) {
	payload, err := streams.Serialize(value)
	if err != nil {
		return nil, err
	}
	payload["@context"] = []string{activityStreamsContext, eventContext, schemaContext}
	normalizeEventRecipients(payload)
	return json.Marshal(payload)
}

// Gancio expects Event recipients as arrays even when ActivityStreams JSON-LD
// compaction would emit one recipient as a scalar.
func normalizeEventRecipients(payload map[string]interface{}) {
	event := payload
	if object, ok := payload["object"].(map[string]interface{}); ok {
		event = object
	}
	if event["type"] != "Event" {
		return
	}
	for _, property := range []string{"to", "cc", "attachment"} {
		value, ok := event[property]
		if !ok {
			continue
		}
		if _, ok := value.([]interface{}); !ok {
			event[property] = []interface{}{value}
		}
	}
}
