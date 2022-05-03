package apmodels

import (
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
)

// CreateCreateActivity will create a new Create Activity model with the provided ID and IRI.
func CreateCreateActivity(id string, localAccountIRI *url.URL) vocab.ActivityStreamsCreate {
	objectID := MakeLocalIRIForResource(id)
	message := MakeCreateActivity(objectID)

	actorProp := streams.NewActivityStreamsActorProperty()
	actorProp.AppendIRI(localAccountIRI)
	message.SetActivityStreamsActor(actorProp)

	return message
}

// AddImageAttachmentToNote will add the provided image URL to the provided note object.
func AddImageAttachmentToNote(note vocab.ActivityStreamsNote, image, mediaType string) {
	imageURL, err := url.Parse(image)
	if err != nil {
		return
	}

	attachments := note.GetActivityStreamsAttachment()
	if attachments == nil {
		attachments = streams.NewActivityStreamsAttachmentProperty()
	}

	urlProp := streams.NewActivityStreamsUrlProperty()
	urlProp.AppendIRI(imageURL)

	apImage := streams.NewActivityStreamsImage()
	apImage.SetActivityStreamsUrl(urlProp)

	imageProp := streams.NewActivityStreamsImageProperty()
	imageProp.AppendActivityStreamsImage(apImage)

	imageDescription := streams.NewActivityStreamsNameProperty()
	imageDescription.AppendXMLSchemaString("Live stream preview")
	apImage.SetActivityStreamsName(imageDescription)

	mediaTypeProperty := streams.NewActivityStreamsMediaTypeProperty()
	mediaTypeProperty.Set(mediaType)
	apImage.SetActivityStreamsMediaType(mediaTypeProperty)

	attachments.AppendActivityStreamsImage(apImage)

	note.SetActivityStreamsAttachment(attachments)
}
