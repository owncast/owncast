package requests

import (
	"encoding/json"
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"

	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/crypto"
	"github.com/owncast/owncast/services/activitypub/workerpool"
	"github.com/owncast/owncast/utils"

	"github.com/teris-io/shortid"
)

// SendFollowAccept will send an accept activity to a follow request
// from a specified local user, queuing the outbound delivery on the
// provided workerpool. The Accept carries this server's current Owncast
// stream status so a newly accepted featured-streams follower reflects our
// live state immediately rather than waiting for the next periodic ping.
func SendFollowAccept(wp *workerpool.Service, inbox *url.URL, originalFollowActivity vocab.ActivityStreamsFollow, fromLocalAccountName string, builder *apmodels.Builder, signer *crypto.Signer, configRepository configrepository.ConfigRepository, streamActive bool) error {
	// SSRF protection: reject non-HTTPS schemes and internal/loopback hosts.
	if inbox.Scheme != "https" {
		return errors.Errorf("rejecting non-HTTPS inbox URL for SSRF protection: %s", inbox.String())
	}
	if utils.IsHostnameInternal(inbox.Hostname()) {
		return errors.Errorf("rejecting internal/loopback inbox URL for SSRF protection: %s", inbox.String())
	}

	followAccept := makeAcceptFollow(originalFollowActivity, fromLocalAccountName, builder, configRepository, streamActive)
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

	// Attach our current stream status/metadata so a featured-streams follower
	// can show our live state the moment the follow is accepted.
	apmodels.SetOwncastMetadata(accept.GetUnknownProperties(), configRepository, streamActive)

	return accept
}
