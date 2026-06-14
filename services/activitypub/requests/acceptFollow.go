package requests

import (
	"encoding/json"
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"

	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/crypto"
	"github.com/owncast/owncast/services/activitypub/workerpool"
	"github.com/owncast/owncast/utils"

	"github.com/teris-io/shortid"
)

// SendFollowAccept will send an accept activity to a follow request
// from a specified local user, queuing the outbound delivery on the
// provided workerpool.
func SendFollowAccept(wp *workerpool.Service, inbox *url.URL, originalFollowActivity vocab.ActivityStreamsFollow, fromLocalAccountName string, builder *apmodels.Builder, signer *crypto.Signer) error {
	// SSRF protection: reject non-HTTPS schemes and internal/loopback hosts.
	if inbox.Scheme != "https" {
		return errors.Errorf("rejecting non-HTTPS inbox URL for SSRF protection: %s", inbox.String())
	}
	if utils.IsHostnameInternal(inbox.Hostname()) {
		return errors.Errorf("rejecting internal/loopback inbox URL for SSRF protection: %s", inbox.String())
	}

	followAccept := makeAcceptFollow(originalFollowActivity, fromLocalAccountName, builder)
	localAccountIRI := builder.MakeLocalIRIForAccount(fromLocalAccountName)

	var jsonmap map[string]interface{}
	jsonmap, _ = streams.Serialize(followAccept)
	b, _ := json.Marshal(jsonmap)
	req, err := signer.CreateSignedRequest(b, inbox, localAccountIRI)
	if err != nil {
		return err
	}

	wp.AddToOutboundQueue(req)

	return nil
}

func makeAcceptFollow(originalFollowActivity vocab.ActivityStreamsFollow, fromAccountName string, builder *apmodels.Builder) vocab.ActivityStreamsAccept {
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

	return accept
}
