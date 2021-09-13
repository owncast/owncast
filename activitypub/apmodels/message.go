package apmodels

import (
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
)

// CreateCreateActivity will create a new Create Activity model with the provided ID and IRI.
func CreateCreateActivity(id string, localAccountIRI *url.URL) vocab.ActivityStreamsCreate {
	objectID := MakeLocalIRIForResource("/create-" + id)
	message := MakeCreateActivity(objectID)

	actorProp := streams.NewActivityStreamsActorProperty()
	actorProp.AppendIRI(localAccountIRI)
	message.SetActivityStreamsActor(actorProp)

	return message
}

// CreateMessageActivity will create a new Create Activity model with a note object.
func CreateMessageActivity(id string, content string, localAccountIRI *url.URL) vocab.ActivityStreamsCreate {
	toPublic, _ := url.Parse(PUBLIC)

	objectID := MakeLocalIRIForResource(id)
	noteID := localAccountIRI
	noteID.Path = noteID.Path + "/note-" + id

	actorProp := streams.NewActivityStreamsActorProperty()
	actorProp.AppendIRI(localAccountIRI)

	message := MakeCreateActivity(objectID)

	object := streams.NewActivityStreamsObjectProperty()

	note := MakeNote(content, noteID, localAccountIRI)
	object.AppendActivityStreamsNote(note)

	message.SetActivityStreamsActor(actorProp)

	to := streams.NewActivityStreamsToProperty()
	to.AppendIRI(toPublic)
	note.SetActivityStreamsTo(to)

	// TODO: Add CW for NSFW streams.
	// Content warning is done via the summary property.
	// summary := NewActivityStreamsSummaryProperty()
	// summary.AppendXMLSchemaString("NSFW")

	message.SetActivityStreamsObject(object)

	return message
}

// AddImageAttachmentToNote will add the provided image URL to the provided note object.
func AddImageAttachmentToNote(note vocab.ActivityStreamsNote, image string) {
	imageURL, err := url.Parse(image)
	if err != nil {
		return
	}

	var attachments = note.GetActivityStreamsAttachment()
	if attachments == nil {
		attachments = streams.NewActivityStreamsAttachmentProperty()
	}

	urlProp := streams.NewActivityStreamsUrlProperty()
	urlProp.AppendIRI(imageURL)

	apImage := streams.NewActivityStreamsImage()
	apImage.SetActivityStreamsUrl(urlProp)

	imageProp := streams.NewActivityStreamsImageProperty()
	imageProp.AppendActivityStreamsImage(apImage)

	imageDescription := streams.NewActivityStreamsContentProperty()
	imageDescription.AppendXMLSchemaString("Live stream preview")
	apImage.SetActivityStreamsContent(imageDescription)

	attachments.AppendActivityStreamsImage(apImage)

	note.SetActivityStreamsAttachment(attachments)
}
