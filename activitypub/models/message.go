package models

import (
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/core/data"
	"github.com/teris-io/shortid"
)

func CreateCreateActivity(localAccountIRI *url.URL) vocab.ActivityStreamsCreate {
	id := shortid.MustGenerate()
	objectId := MakeLocalIRIForResource("/create-" + id)
	message := MakeActivity(objectId)

	actorProp := streams.NewActivityStreamsActorProperty()
	actorProp.AppendIRI(localAccountIRI)
	message.SetActivityStreamsActor(actorProp)

	return message
}

func CreateMessageActivity(content string, localAccountIRI *url.URL) vocab.ActivityStreamsCreate {
	toPublic, _ := url.Parse(PUBLIC)

	id := shortid.MustGenerate()
	objectId := MakeLocalIRIForResource("/create-" + id)
	noteId := MakeLocalIRIForResource(data.GetDefaultFederationUsername() + "/" + id)

	actorProp := streams.NewActivityStreamsActorProperty()
	actorProp.AppendIRI(localAccountIRI)

	message := MakeActivity(objectId)

	object := streams.NewActivityStreamsObjectProperty()

	note := MakeNote(content, noteId, localAccountIRI)
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
