package models

import (
	"net/url"
	"time"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
)

type PrivacyAudience = string

const (
	PUBLIC PrivacyAudience = "https://www.w3.org/ns/activitystreams#Public"
)

func MakeActivity(activityID *url.URL) vocab.ActivityStreamsCreate {
	activity := streams.NewActivityStreamsCreate()
	id := streams.NewJSONLDIdProperty()
	id.Set(activityID)
	activity.SetJSONLDId(id)

	return activity
}

func MakeNote(text string, noteIRI *url.URL, attributedToIRI *url.URL) vocab.ActivityStreamsNote {
	note := streams.NewActivityStreamsNote()
	content := streams.NewActivityStreamsContentProperty()
	content.AppendXMLSchemaString(text)
	note.SetActivityStreamsContent(content)
	id := streams.NewJSONLDIdProperty()
	id.Set(noteIRI)
	note.SetJSONLDId(id)

	published := streams.NewActivityStreamsPublishedProperty()
	published.Set(time.Now())
	note.SetActivityStreamsPublished(published)

	attributedTo := attributedToIRI
	attr := streams.NewActivityStreamsAttributedToProperty()
	attr.AppendIRI(attributedTo)
	note.SetActivityStreamsAttributedTo(attr)

	public, _ := url.Parse(PUBLIC)
	to := streams.NewActivityStreamsToProperty()
	to.AppendIRI(public)
	note.SetActivityStreamsTo(to)

	return note
}
