package inbox

import (
	"context"
	"net/url"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/persistence"
)

func handleDeleteRequest(c context.Context, activity vocab.ActivityStreamsDelete) error {
	obj := activity.GetActivityStreamsObject()
	if obj != nil && obj.Len() > 0 && obj.At(0) != nil {
		o := obj.At(0)

		var actorIri *url.URL
		switch true {
		case o.IsActivityStreamsPerson():
			actorIri = o.GetActivityStreamsPerson().GetJSONLDId().Get()
		case o.IsActivityStreamsService():
			actorIri = o.GetActivityStreamsService().GetJSONLDId().Get()
		case o.IsActivityStreamsApplication():
			actorIri = o.GetActivityStreamsApplication().GetJSONLDId().Get()
		default:
			// not a type we care about here
			return nil
		}

		return persistence.RemoveFollowIRI(actorIri)
	}

	return nil
}
