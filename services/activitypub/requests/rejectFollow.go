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

// SendFollowReject sends a Reject of a previously received Follow to that
// follower's inbox, queuing delivery on the provided workerpool. Owncast uses
// this to revoke a follower it had accepted: when the operator removes a
// directory that is listing this server, the Reject tells the directory its
// follow is no longer accepted so it drops the listing, rather than leaving the
// server showing offline there forever. It is the same signal the wider
// Fediverse sends when an account removes one of its followers.
func SendFollowReject(wp *workerpool.Service, inbox *url.URL, originalFollowActivity vocab.ActivityStreamsFollow, fromLocalAccountName string, builder *apmodels.Builder, signer *crypto.Signer) error {
	// SSRF protection: reject non-HTTPS schemes and internal/loopback hosts.
	if inbox.Scheme != "https" {
		return errors.Errorf("rejecting non-HTTPS inbox URL for SSRF protection: %s", inbox.String())
	}
	if utils.IsHostnameInternal(inbox.Hostname()) {
		return errors.Errorf("rejecting internal/loopback inbox URL for SSRF protection: %s", inbox.String())
	}

	followReject := makeRejectFollow(originalFollowActivity, fromLocalAccountName, builder)
	localAccountIRI := builder.MakeLocalIRIForAccount(fromLocalAccountName)

	var jsonmap map[string]interface{}
	jsonmap, _ = streams.Serialize(followReject)
	b, _ := json.Marshal(jsonmap)
	req, err := signer.CreateSignedRequest(b, inbox, localAccountIRI)
	if err != nil {
		return err
	}

	wp.AddToOutboundQueue(req)

	return nil
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
