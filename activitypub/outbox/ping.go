package outbox

import (
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

// SendStreamPing sends an Offer activity to all followers indicating the stream is still live.
func SendStreamPing() error {
	configRepository := configrepository.Get()

	// Don't send if federation is disabled
	if !configRepository.GetFederationEnabled() {
		return nil
	}

	id := shortid.MustGenerate()
	activityID := apmodels.MakeLocalIRIForResource(id)
	localActor := apmodels.MakeLocalIRIForAccount(configRepository.GetDefaultFederationUsername())
	serverURL := configRepository.GetServerURL()

	// Create the Offer activity
	activity := streams.NewActivityStreamsOffer()

	// Set the activity ID
	idProperty := streams.NewJSONLDIdProperty()
	idProperty.Set(activityID)
	activity.SetJSONLDId(idProperty)

	// Set the actor (the Owncast server)
	actorProperty := streams.NewActivityStreamsActorProperty()
	actorProperty.AppendIRI(localActor)
	activity.SetActivityStreamsActor(actorProperty)

	// Set the object (the server URL - offering the live stream)
	objectProperty := streams.NewActivityStreamsObjectProperty()
	serverIRI, err := url.Parse(serverURL)
	if err != nil {
		return errors.Wrap(err, "unable to parse server URL for Offer activity")
	}
	objectProperty.AppendIRI(serverIRI)
	activity.SetActivityStreamsObject(objectProperty)

	// Add custom Owncast metadata indicating stream is live
	unknownProps := activity.GetUnknownProperties()
	apmodels.SetOwncastMetadata(unknownProps, configRepository, true)

	// Set addressing to followers
	to, cc := getAddressingToFollowers()
	activity.SetActivityStreamsTo(to)
	activity.SetActivityStreamsCc(cc)

	// Serialize and send
	b, err := apmodels.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize stream ping Offer activity", err)
		return errors.Wrap(err, "unable to serialize stream ping Offer activity")
	}

	if err := SendToFollowers(b); err != nil {
		return err
	}

	// Add to outbox
	if err := Add(activity, id, false); err != nil {
		return err
	}

	log.Debugln("Sent stream ping Offer activity to all followers")
	return nil
}
