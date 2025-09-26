package requests

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/workerpool"
	"github.com/owncast/owncast/persistence/configrepository"

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

	// Add Owncast metadata to the accept activity
	configRepository := configrepository.Get()
	unknownProps := accept.GetUnknownProperties()

	// Always include server identification
	unknownProps["https://owncast.online/ns#serverName"] = configRepository.GetServerName()
	unknownProps["https://owncast.online/ns#streamDescription"] = configRepository.GetServerSummary()

	// Add logo if available
	if logoPath := configRepository.GetLogoPath(); logoPath != "" {
		logoURL := fmt.Sprintf("%s/%s", configRepository.GetServerURL(), logoPath)
		unknownProps["https://owncast.online/ns#logoUrl"] = logoURL
	}

	// Add tags if available
	if tags := configRepository.GetServerMetadataTags(); len(tags) > 0 {
		unknownProps["https://owncast.online/ns#streamTags"] = tags
	}

	// Add current stream status - this allows the follower to immediately know our status
	// We don't include the full stream metadata here since it's just an accept,
	// but we can include basic server info
	if configRepository.GetStreamTitle() != "" {
		unknownProps["https://owncast.online/ns#streamTitle"] = configRepository.GetStreamTitle()
		// If there's a stream title, we might be live or have a standing title
		// For now, don't assume live status in accept, let the regular Note activities handle that
	}

	return accept
}
