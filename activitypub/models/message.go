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

	// addImageAttachmentToNote(note, "https://watch.owncast.online/preview.gif")
	// addVideoAttachmentToNote(note, "https://goth.land/hls/stream.m3u8")

	message.SetActivityStreamsObject(object)

	return message
}
