package follower

import (
	"encoding/json"
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/teris-io/shortid"
)

// SendFollowAccept will send an accept activity to a follow request from a specified local user.
func (s *Service) SendFollowAccept(inbox *url.URL, originalFollowActivity vocab.ActivityStreamsFollow, fromLocalAccountName string) error {
	followAccept := s.makeAcceptFollow(originalFollowActivity, fromLocalAccountName)
	localAccountIRI := s.models.MakeLocalIRIForAccount(fromLocalAccountName)

	var jsonmap map[string]interface{}
	jsonmap, _ = streams.Serialize(followAccept)
	b, _ := json.Marshal(jsonmap)
	req, err := s.crypto.CreateSignedRequest(b, inbox, localAccountIRI)
	if err != nil {
		return err
	}

	s.worker.AddToOutboundQueue(req)

	return nil
}

func (s *Service) makeAcceptFollow(originalFollowActivity vocab.ActivityStreamsFollow, fromAccountName string) vocab.ActivityStreamsAccept {
	acceptIDString := shortid.MustGenerate()
	acceptID := s.models.MakeLocalIRIForResource(acceptIDString)
	actorID := s.models.MakeLocalIRIForAccount(fromAccountName)

	accept := streams.NewActivityStreamsAccept()
	idProperty := streams.NewJSONLDIdProperty()
	idProperty.SetIRI(acceptID)
	accept.SetJSONLDId(idProperty)

	actor := s.models.MakeActorPropertyWithID(actorID)
	accept.SetActivityStreamsActor(actor)

	object := streams.NewActivityStreamsObjectProperty()
	object.AppendActivityStreamsFollow(originalFollowActivity)
	accept.SetActivityStreamsObject(object)

	return accept
}
