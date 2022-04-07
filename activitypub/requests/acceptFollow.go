package requests

import (
	"encoding/json"
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/workerpool"

	"github.com/teris-io/shortid"
)

// SendFollowAccept will send an accept activity to a follow request from a specified local user.
func SendFollowAccept(inbox *url.URL, originalFollowActivity vocab.ActivityStreamsFollow, fromLocalAccountName string) error {
	followAccept := makeAcceptFollow(originalFollowActivity, fromLocalAccountName)
	localAccountIRI := apmodels.MakeLocalIRIForAccount(fromLocalAccountName)

	var jsonmap map[string]interface{}
	jsonmap, _ = streams.Serialize(followAccept)
	b, _ := json.Marshal(jsonmap)
	req, err := crypto.CreateSignedRequest(b, inbox, localAccountIRI)
	if err != nil {
		return err
	}

	workerpool.AddToOutboundQueue(req)

	return nil
}

func makeAcceptFollow(originalFollowActivity vocab.ActivityStreamsFollow, fromAccountName string) vocab.ActivityStreamsAccept {
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
	object.AppendActivityStreamsFollow(originalFollowActivity)
	accept.SetActivityStreamsObject(object)

	return accept
}
