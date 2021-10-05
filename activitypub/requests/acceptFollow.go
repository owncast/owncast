package requests

import (
	"encoding/json"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/teris-io/shortid"

	log "github.com/sirupsen/logrus"
)

// SendFollowAccept will send an accept activity to a follow request from a specified local user.
func SendFollowAccept(followRequest apmodels.ActivityPubActor, fromLocalAccountName string) error {
	log.Println("SendFollowAccept 1")

	followAccept := makeAcceptFollow(followRequest, fromLocalAccountName)
	localAccountIRI := apmodels.MakeLocalIRIForAccount(fromLocalAccountName)

	var jsonmap map[string]interface{}
	jsonmap, _ = streams.Serialize(followAccept)
	b, _ := json.Marshal(jsonmap)
	log.Println("SendFollowAccept 2")

	_, err := PostSignedRequest(b, followRequest.Inbox, localAccountIRI)
	log.Println("SendFollowAccept 3")

	if err != nil {
		return err
	}

	return nil
}

func makeAcceptFollow(follow apmodels.ActivityPubActor, fromAccountName string) vocab.ActivityStreamsAccept {
	acceptIDString := shortid.MustGenerate()
	acceptID := apmodels.MakeLocalIRIForResource(acceptIDString)
	actorID := apmodels.MakeLocalIRIForAccount(fromAccountName)

	accept := streams.NewActivityStreamsAccept()
	idProperty := streams.NewJSONLDIdProperty()
	idProperty.SetIRI(acceptID)
	accept.SetJSONLDId(idProperty)

	actor := apmodels.MakeActorPropertyWithID(actorID)
	accept.SetActivityStreamsActor(actor)

	object := streams.NewActivityStreamsObjectProperty()
	object.AppendIRI(follow.FollowIri)
	accept.SetActivityStreamsObject(object)

	return accept
}
