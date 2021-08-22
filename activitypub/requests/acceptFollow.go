package requests

import (
	"encoding/json"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/teris-io/shortid"
)

func SendFollowAccept(followRequest apmodels.ActivityPubActor, fromLocalAccountName string) error {
	followAccept := makeAcceptFollow(followRequest, fromLocalAccountName)
	localAccountIRI := apmodels.MakeLocalIRIForAccount(fromLocalAccountName)

	var jsonmap map[string]interface{}
	jsonmap, _ = streams.Serialize(followAccept)
	b, _ := json.Marshal(jsonmap)

	_, err := PostSignedRequest(b, followRequest.Inbox, localAccountIRI)
	if err != nil {
		return err
	}

	return nil
}

func makeAcceptFollow(follow apmodels.ActivityPubActor, fromAccountName string) vocab.ActivityStreamsAccept {
	acceptIdString := shortid.MustGenerate()
	acceptId := apmodels.MakeLocalIRIForResource(acceptIdString)

	actorId := apmodels.MakeLocalIRIForAccount(fromAccountName)

	accept := streams.NewActivityStreamsAccept()
	idProperty := streams.NewJSONLDIdProperty()
	idProperty.SetIRI(acceptId)
	accept.SetJSONLDId(idProperty)

	actor := apmodels.MakeActorPropertyWithId(actorId)
	accept.SetActivityStreamsActor(actor)

	object := streams.NewActivityStreamsObjectProperty()
	object.AppendIRI(follow.FollowIri)
	accept.SetActivityStreamsObject(object)

	return accept
}
