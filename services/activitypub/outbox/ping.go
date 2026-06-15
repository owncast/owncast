package outbox

import (
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/services/activitypub/apmodels"
)

// SendStreamPing sends an Offer activity to all followers indicating
// the stream is still live. Used by the featured-streams flow so
// remote Owncast servers can keep their mini-directory of live
// streams fresh without polling.
func (s *Service) SendStreamPing() error {
	if !s.configRepository.GetFederationEnabled() {
		return nil
	}

	id := shortid.MustGenerate()
	activityID := s.builder.MakeLocalIRIForResource(id)
	localActor := s.builder.MakeLocalIRIForAccount(s.configRepository.GetDefaultFederationUsername())
	serverURL := s.configRepository.GetServerURL()

	// Create the Offer activity.
	activity := streams.NewActivityStreamsOffer()

	idProperty := streams.NewJSONLDIdProperty()
	idProperty.Set(activityID)
	activity.SetJSONLDId(idProperty)

	actorProperty := streams.NewActivityStreamsActorProperty()
	actorProperty.AppendIRI(localActor)
	activity.SetActivityStreamsActor(actorProperty)

	// The object of the Offer is our server URL (we're offering the
	// live stream).
	objectProperty := streams.NewActivityStreamsObjectProperty()
	serverIRI, err := url.Parse(serverURL)
	if err != nil {
		return errors.Wrap(err, "unable to parse server URL for Offer activity")
	}
	objectProperty.AppendIRI(serverIRI)
	activity.SetActivityStreamsObject(objectProperty)

	// Attach Owncast metadata so receivers can populate their
	// federated_servers table from the ping alone.
	unknownProps := activity.GetUnknownProperties()
	apmodels.SetOwncastMetadata(unknownProps, s.configRepository, true)

	to, cc := s.getAddressingToFollowers()
	activity.SetActivityStreamsTo(to)
	activity.SetActivityStreamsCc(cc)

	b, err := apmodels.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize stream ping Offer activity", err)
		return errors.Wrap(err, "unable to serialize stream ping Offer activity")
	}

	if err := s.SendToFollowers(b); err != nil {
		return err
	}

	if err := s.Add(activity, id, false); err != nil {
		return err
	}

	log.Debugln("Sent stream ping Offer activity to all followers")
	return nil
}
