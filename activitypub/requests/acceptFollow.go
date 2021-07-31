package requests

import (
	"encoding/json"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/models"
	"github.com/teris-io/shortid"
)

func SendFollowAccept(followRequest models.ActivityPubActor, fromLocalAccountName string) error {
	followAccept := makeAcceptFollow(followRequest, fromLocalAccountName)
	localAccountIRI := models.MakeLocalIRIForAccount(fromLocalAccountName)

	var jsonmap map[string]interface{}
	jsonmap, _ = streams.Serialize(followAccept)
	b, _ := json.Marshal(jsonmap)

	_, err := PostSignedRequest(b, followRequest.Inbox, localAccountIRI)
	if err != nil {
		return err
	}

	return nil
}

func makeAcceptFollow(follow models.ActivityPubActor, fromAccountName string) vocab.ActivityStreamsAccept {
	acceptIdString := shortid.MustGenerate()
	acceptId := models.MakeLocalIRIForResource(acceptIdString)

	actorId := models.MakeLocalIRIForAccount(fromAccountName)

	accept := streams.NewActivityStreamsAccept()
	idProperty := streams.NewJSONLDIdProperty()
	idProperty.SetIRI(acceptId)
	accept.SetJSONLDId(idProperty)

	actor := models.MakeActorPropertyWithId(actorId)
	accept.SetActivityStreamsActor(actor)

	object := streams.NewActivityStreamsObjectProperty()
	object.AppendIRI(follow.FollowIri)
	accept.SetActivityStreamsObject(object)

	return accept
}
