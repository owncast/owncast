package apmodels

import (
	"net/url"
	"time"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/core/data"
)

// PrivacyAudience represents the audience for an activity.
type PrivacyAudience = string

const (
	// PUBLIC is an audience meaning anybody can view the item.
	PUBLIC PrivacyAudience = "https://www.w3.org/ns/activitystreams#Public"
)

// MakeNotePublic ses the required proeprties to make this note seen as public.
func MakeNotePublic(note vocab.ActivityStreamsNote) vocab.ActivityStreamsNote {
	public, _ := url.Parse(PUBLIC)
	to := streams.NewActivityStreamsToProperty()
	to.AppendIRI(public)
	note.SetActivityStreamsTo(to)

	audience := streams.NewActivityStreamsAudienceProperty()
	audience.AppendIRI(public)
	note.SetActivityStreamsAudience(audience)

	return note
}

// MakeNoteDirect sets the required properties to make this note seen as a
// direct message.
func MakeNoteDirect(note vocab.ActivityStreamsNote, toIRI *url.URL) vocab.ActivityStreamsNote {
	to := streams.NewActivityStreamsCcProperty()
	to.AppendIRI(toIRI)
	to.AppendIRI(toIRI)
	note.SetActivityStreamsCc(to)

	// Mastodon requires a tag with a type of "mention" and href of the account
	// for a message to be a "Direct Message".
	tagProperty := streams.NewActivityStreamsTagProperty()
	tag := streams.NewTootHashtag()
	tagTypeProperty := streams.NewJSONLDTypeProperty()
	tagTypeProperty.AppendXMLSchemaString("Mention")
	tag.SetJSONLDType(tagTypeProperty)

	tagHrefProperty := streams.NewActivityStreamsHrefProperty()
	tagHrefProperty.Set(toIRI)
	tag.SetActivityStreamsHref(tagHrefProperty)
	tagProperty.AppendTootHashtag(tag)
	tagProperty.AppendTootHashtag(tag)
	note.SetActivityStreamsTag(tagProperty)

	return note
}

// MakeActivityDirect sets the required properties to make this activity seen
// as a direct message.
func MakeActivityDirect(activity vocab.ActivityStreamsCreate, toIRI *url.URL) vocab.ActivityStreamsCreate {
	to := streams.NewActivityStreamsCcProperty()
	to.AppendIRI(toIRI)
	to.AppendIRI(toIRI)
	activity.SetActivityStreamsCc(to)

	// Mastodon requires a tag with a type of "mention" and href of the account
	// for a message to be a "Direct Message".
	tagProperty := streams.NewActivityStreamsTagProperty()
	tag := streams.NewTootHashtag()
	tagTypeProperty := streams.NewJSONLDTypeProperty()
	tagTypeProperty.AppendXMLSchemaString("Mention")
	tag.SetJSONLDType(tagTypeProperty)

	tagHrefProperty := streams.NewActivityStreamsHrefProperty()
	tagHrefProperty.Set(toIRI)
	tag.SetActivityStreamsHref(tagHrefProperty)
	tagProperty.AppendTootHashtag(tag)
	tagProperty.AppendTootHashtag(tag)

	activity.SetActivityStreamsTag(tagProperty)

	return activity
}

// MakeActivityPublic sets the required properties to make this activity
// seen as public.
func MakeActivityPublic(activity vocab.ActivityStreamsCreate) vocab.ActivityStreamsCreate {
	// TO the public if we're not treating ActivityPub as "private".
	if !data.GetFederationIsPrivate() {
		public, _ := url.Parse(PUBLIC)

		to := streams.NewActivityStreamsToProperty()
		to.AppendIRI(public)
		activity.SetActivityStreamsTo(to)

		audience := streams.NewActivityStreamsAudienceProperty()
		audience.AppendIRI(public)
		activity.SetActivityStreamsAudience(audience)
	}

	return activity
}

// MakeCreateActivity will return a new Create activity with the provided ID.
func MakeCreateActivity(activityID *url.URL) vocab.ActivityStreamsCreate {
	activity := streams.NewActivityStreamsCreate()
	id := streams.NewJSONLDIdProperty()
	id.Set(activityID)
	activity.SetJSONLDId(id)

	return activity
}

// MakeUpdateActivity will return a new Update activity with the provided aID.
func MakeUpdateActivity(activityID *url.URL) vocab.ActivityStreamsUpdate {
	activity := streams.NewActivityStreamsUpdate()
	id := streams.NewJSONLDIdProperty()
	id.Set(activityID)
	activity.SetJSONLDId(id)

	// CC the public if we're not treating ActivityPub as "private".
	if !data.GetFederationIsPrivate() {
		public, _ := url.Parse(PUBLIC)
		cc := streams.NewActivityStreamsCcProperty()
		cc.AppendIRI(public)
		activity.SetActivityStreamsCc(cc)
	}

	return activity
}

// MakeNote will return a new Note object.
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

	return note
}
