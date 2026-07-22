package requests

import (
	"encoding/json"
	"fmt"
	"net/url"

	"code.superseriousbusiness.org/activity/streams"
	"code.superseriousbusiness.org/activity/streams/vocab"

	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/workerpool"

	"github.com/teris-io/shortid"
)

// SendFollowReject queues a Reject for a previously received Follow. Owncast
// uses this to revoke a directory follower so the remote server drops its
// listing instead of leaving it offline forever.
func SendFollowReject(wp *workerpool.Service, inbox *url.URL, originalFollowActivity vocab.ActivityStreamsFollow, fromLocalAccountName string, builder *apmodels.Builder) error {
	delivery, err := MakeFollowRejectDelivery(inbox, originalFollowActivity, fromLocalAccountName, builder)
	if err != nil {
		return err
	}
	return wp.Enqueue(delivery)
}

// MakeFollowRejectDelivery builds the durable delivery for a Follow Reject.
func MakeFollowRejectDelivery(inbox *url.URL, originalFollowActivity vocab.ActivityStreamsFollow, fromLocalAccountName string, builder *apmodels.Builder) (workerpool.Delivery, error) {
	if err := validateRemoteInbox(inbox); err != nil {
		return workerpool.Delivery{}, err
	}

	followReject := makeRejectFollow(originalFollowActivity, fromLocalAccountName, builder)
	localAccountIRI := builder.MakeLocalIRIForAccount(fromLocalAccountName)

	jsonMap, err := streams.Serialize(followReject)
	if err != nil {
		return workerpool.Delivery{}, fmt.Errorf("serializing Follow Reject: %w", err)
	}
	payload, err := json.Marshal(jsonMap)
	if err != nil {
		return workerpool.Delivery{}, fmt.Errorf("marshalling Follow Reject: %w", err)
	}
	remoteActorIRI, err := apmodels.GetIRIStringFromActorProperty(originalFollowActivity.GetActivityStreamsActor())
	if err != nil {
		return workerpool.Delivery{}, fmt.Errorf("reading Follow actor: %w", err)
	}

	return workerpool.Delivery{
		Inbox:        inbox,
		Payload:      payload,
		ActorIRI:     localAccountIRI,
		ActivityType: "Reject",
		CoalesceKey:  "follow-response:" + remoteActorIRI,
	}, nil
}

func makeRejectFollow(originalFollowActivity vocab.ActivityStreamsFollow, fromAccountName string, builder *apmodels.Builder) vocab.ActivityStreamsReject {
	rejectIDString := shortid.MustGenerate()
	rejectID := builder.MakeLocalIRIForResource(rejectIDString)
	actorID := builder.MakeLocalIRIForAccount(fromAccountName)

	reject := streams.NewActivityStreamsReject()
	idProperty := streams.NewJSONLDIdProperty()
	idProperty.SetIRI(rejectID)
	reject.SetJSONLDId(idProperty)

	actor := apmodels.MakeActorPropertyWithID(actorID)
	reject.SetActivityStreamsActor(actor)

	object := streams.NewActivityStreamsObjectProperty()
	object.AppendActivityStreamsFollow(originalFollowActivity)
	reject.SetActivityStreamsObject(object)

	return reject
}
