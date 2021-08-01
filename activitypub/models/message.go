package models

import (
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/core/data"
	"github.com/teris-io/shortid"
)

func CreateMessage(content string, localAccountIRI *url.URL) vocab.ActivityStreamsCreate {
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

	addImageAttachmentToNote(note, "https://watch.owncast.online/preview.gif")

	message.SetActivityStreamsObject(object)

	return message
}

func addImageAttachmentToNote(note vocab.ActivityStreamsNote, image string) {
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
	attachments.AppendActivityStreamsImage(apImage)

	note.SetActivityStreamsAttachment(attachments)
}
