package requests

import (
	"encoding/json"
	"fmt"
	"net/url"

	"code.superseriousbusiness.org/activity/streams"
	"code.superseriousbusiness.org/activity/streams/vocab"

	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/workerpool"

	"github.com/teris-io/shortid"
)

// SendFollowAccept queues an Accept for a received Follow. The Accept carries
// this server's current stream status so a newly accepted featured-streams
// follower reflects our live state immediately.
func SendFollowAccept(wp *workerpool.Service, inbox *url.URL, originalFollowActivity vocab.ActivityStreamsFollow, fromLocalAccountName string, builder *apmodels.Builder, configRepository configrepository.ConfigRepository, streamActive bool) error {
	delivery, err := MakeFollowAcceptDelivery(inbox, originalFollowActivity, fromLocalAccountName, builder, configRepository, streamActive)
	if err != nil {
		return err
	}
	return wp.Enqueue(delivery)
}

// MakeFollowAcceptDelivery builds the durable delivery for a Follow Accept.
func MakeFollowAcceptDelivery(inbox *url.URL, originalFollowActivity vocab.ActivityStreamsFollow, fromLocalAccountName string, builder *apmodels.Builder, configRepository configrepository.ConfigRepository, streamActive bool) (workerpool.Delivery, error) {
	if err := validateRemoteInbox(inbox); err != nil {
		return workerpool.Delivery{}, err
	}

	followAccept := makeAcceptFollow(originalFollowActivity, fromLocalAccountName, builder, configRepository, streamActive)
	localAccountIRI := builder.MakeLocalIRIForAccount(fromLocalAccountName)

	jsonMap, err := streams.Serialize(followAccept)
	if err != nil {
		return workerpool.Delivery{}, fmt.Errorf("serializing Follow Accept: %w", err)
	}
	payload, err := json.Marshal(jsonMap)
	if err != nil {
		return workerpool.Delivery{}, fmt.Errorf("marshalling Follow Accept: %w", err)
	}
	remoteActorIRI, err := apmodels.GetIRIStringFromActorProperty(originalFollowActivity.GetActivityStreamsActor())
	if err != nil {
		return workerpool.Delivery{}, fmt.Errorf("reading Follow actor: %w", err)
	}

	return workerpool.Delivery{
		Inbox:        inbox,
		Payload:      payload,
		ActorIRI:     localAccountIRI,
		ActivityType: "Accept",
		CoalesceKey:  "follow-response:" + remoteActorIRI,
	}, nil
}

func makeAcceptFollow(originalFollowActivity vocab.ActivityStreamsFollow, fromAccountName string, builder *apmodels.Builder, configRepository configrepository.ConfigRepository, streamActive bool) vocab.ActivityStreamsAccept {
	acceptIDString := shortid.MustGenerate()
	acceptID := builder.MakeLocalIRIForResource(acceptIDString)
	actorID := builder.MakeLocalIRIForAccount(fromAccountName)

	accept := streams.NewActivityStreamsAccept()
	idProperty := streams.NewJSONLDIdProperty()
	idProperty.SetIRI(acceptID)
	accept.SetJSONLDId(idProperty)

	actor := apmodels.MakeActorPropertyWithID(actorID)
	accept.SetActivityStreamsActor(actor)

	object := streams.NewActivityStreamsObjectProperty()
	object.AppendActivityStreamsFollow(originalFollowActivity)
	accept.SetActivityStreamsObject(object)

	apmodels.SetOwncastMetadata(accept.GetUnknownProperties(), configRepository, streamActive)

	return accept
}
