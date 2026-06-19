package outbox

import (
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/services/activitypub/apmodels"
)

// streamStatusActivity is satisfied by the ActivityStreams activity types used
// to advertise stream status to followers: an Offer when going live and a
// Leave when going offline. It captures just the setters the shared builder
// needs so both can travel one code path.
type streamStatusActivity interface {
	vocab.Type
	GetUnknownProperties() map[string]interface{}
	SetActivityStreamsActor(vocab.ActivityStreamsActorProperty)
	SetActivityStreamsObject(vocab.ActivityStreamsObjectProperty)
	SetActivityStreamsTo(vocab.ActivityStreamsToProperty)
	SetActivityStreamsCc(vocab.ActivityStreamsCcProperty)
}

// sendStreamStatusToFollowers stamps the given activity with this server's
// identity, attaches its Owncast stream metadata (live or offline), and
// delivers it to directory followers only (the servers listing this stream).
// Fan followers learn about going live from the Create/Note post instead, so
// they are deliberately excluded here. logLabel is used only for log and error
// messages. It is the shared implementation behind SendStreamPing (Offer/live)
// and SendStreamGoingOffline (Leave/offline).
func (s *Service) sendStreamStatusToFollowers(activity streamStatusActivity, isLive bool, logLabel string) error {
	id := shortid.MustGenerate()
	activityID := s.builder.MakeLocalIRIForResource(id)
	localActor := s.builder.MakeLocalIRIForAccount(s.configRepository.GetDefaultFederationUsername())
	serverURL := s.configRepository.GetServerURL()

	idProperty := streams.NewJSONLDIdProperty()
	idProperty.Set(activityID)
	activity.SetJSONLDId(idProperty)

	actorProperty := streams.NewActivityStreamsActorProperty()
	actorProperty.AppendIRI(localActor)
	activity.SetActivityStreamsActor(actorProperty)

	// The object is our server URL: we're offering (or leaving) the live
	// directory.
	objectProperty := streams.NewActivityStreamsObjectProperty()
	serverIRI, err := url.Parse(serverURL)
	if err != nil {
		return errors.Wrapf(err, "unable to parse server URL for %s", logLabel)
	}
	objectProperty.AppendIRI(serverIRI)
	activity.SetActivityStreamsObject(objectProperty)

	// Attach Owncast metadata so receivers can populate their
	// federated_servers table from this activity alone.
	apmodels.SetOwncastMetadata(activity.GetUnknownProperties(), s.configRepository, isLive)

	to, cc := s.getAddressingToFollowers()
	activity.SetActivityStreamsTo(to)
	activity.SetActivityStreamsCc(cc)

	b, err := apmodels.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize "+logLabel, err)
		return errors.Wrapf(err, "unable to serialize %s", logLabel)
	}

	if err := s.SendToDirectoryFollowers(b); err != nil {
		return err
	}

	if err := s.Add(activity, id, false); err != nil {
		return err
	}

	log.Debugln("Sent " + logLabel + " to directory followers")
	return nil
}

// SendStreamPing sends an Offer activity to directory followers indicating the
// stream is live. Used by the featured-streams flow (at go-live and on a
// timer) so directories can keep their list of live streams fresh without
// polling.
func (s *Service) SendStreamPing() error {
	if !s.configRepository.GetFederationEnabled() {
		return nil
	}

	return s.sendStreamStatusToFollowers(streams.NewActivityStreamsOffer(), true, "stream ping Offer activity")
}

// SendStreamGoingOffline sends a Leave activity to directory followers
// indicating the stream has ended. This is the offline counterpart to
// SendStreamPing: it lets directories drop this server from the live section
// of their list immediately, rather than waiting for the staleness sweep to
// time the entry out.
func (s *Service) SendStreamGoingOffline() error {
	if !s.configRepository.GetFederationEnabled() {
		return nil
	}

	return s.sendStreamStatusToFollowers(streams.NewActivityStreamsLeave(), false, "stream-offline Leave activity")
}
