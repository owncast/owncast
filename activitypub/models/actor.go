package models

import (
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
)

type ActivityPubActor struct {
	ActorIri  *url.URL
	FollowIri *url.URL
	Inbox     *url.URL
}

type DeleteRequest struct {
	ActorIri string
}

func MakeActorPropertyWithId(idIRI *url.URL) vocab.ActivityStreamsActorProperty {
	actor := streams.NewActivityStreamsActorProperty()
	actor.AppendIRI(idIRI)
	return actor
}
