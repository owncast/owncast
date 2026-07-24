package inbox

import (
	"context"
	"time"

	"code.superseriousbusiness.org/activity/streams/vocab"
	"github.com/pkg/errors"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	activityevents "github.com/owncast/owncast/services/activitypub/events"
	"github.com/owncast/owncast/services/activitypub/requests"
)

// handleQuoteRequestInboxRequest handles an inbound QuoteRequest (FEP-044f):
// a remote user asking permission to quote one of our posts. Our public
// posts are quotable by anyone, so any request targeting a post we actually
// authored is accepted and answered with a verifiable QuoteAuthorization
// stamp. Requests for unknown posts, or any request while federation is
// private, are rejected so the pending quote clears on the remote end.
func (s *Service) handleQuoteRequestInboxRequest(c context.Context, activity vocab.GoToSocialQuoteRequest) error {
	activityIRI, err := apmodels.GetIRIStringFromJSONLDIdProperty(activity.GetJSONLDId())
	if err != nil {
		return errors.Wrap(err, "quote request is missing activity IRI")
	}

	actor, err := s.resolver.GetResolvedActorFromActorProperty(activity.GetActivityStreamsActor())
	if err != nil {
		return errors.Wrap(err, "unable to resolve actor of quote request")
	}
	actorIRI := actor.ActorIriString()

	// object = the post being quoted, instrument = the quote post itself.
	quotedPostIRI, err := apmodels.GetIRIFromObjectProperty(activity.GetActivityStreamsObject())
	if err != nil {
		return errors.Wrap(err, "quote request is missing object IRI")
	}

	quotePostIRI, err := apmodels.GetIRIFromInstrumentProperty(activity.GetActivityStreamsInstrument())
	if err != nil {
		return errors.Wrap(err, "quote request is missing instrument IRI")
	}

	localAccountName := s.configRepository.GetDefaultFederationUsername()

	// Only posts (Notes) we authored are quotable, only while our posts are
	// public, and only when the operator has quoting enabled. The outbox store
	// also holds non-post objects like QuoteAuthorization stamps, which must
	// not be quotable themselves. In private federation mode posts are
	// follower-only, so quoting them would leak them to a wider audience.
	if _, _, _, noteErr := s.persistence.GetNoteByIRI(quotedPostIRI.String()); noteErr != nil || s.configRepository.GetFederationIsPrivate() || !s.configRepository.GetFederationEnableQuotes() {
		return requests.SendQuoteRequestReject(s.workerpool, actor.Inbox, activity, localAccountName, s.builder)
	}

	// Store the QuoteAuthorization stamp so other servers can fetch it by IRI
	// to verify the quote was approved.
	stampIRI := s.builder.MakeLocalIRIForResource(shortid.MustGenerate())
	localActorIRI := s.builder.MakeLocalIRIForAccount(localAccountName)
	stamp := apmodels.MakeQuoteAuthorization(stampIRI, localActorIRI, quotePostIRI, quotedPostIRI)

	stampBytes, err := apmodels.Serialize(stamp)
	if err != nil {
		return errors.Wrap(err, "unable to serialize quote authorization")
	}

	if err := s.persistence.AddToOutbox(stampIRI.String(), stampBytes, stamp.GetTypeName(), false); err != nil {
		return errors.Wrap(err, "unable to store quote authorization")
	}

	if err := requests.SendQuoteRequestAccept(s.workerpool, actor.Inbox, activity, stampIRI, localAccountName, s.builder); err != nil {
		return err
	}
	claimed, err := s.persistence.ClaimInboundFediverseActivity(activityIRI, actorIRI, activity.GetTypeName(), time.Now())
	if err != nil {
		return errors.Wrap(err, "unable to claim inbound quote request")
	}
	if !claimed {
		return nil
	}
	s.publishFediverseEvent(c, models.FediverseEngagementQuote, &activityevents.FediverseEngagementEvent{
		Actor:  fediverseActorFromResolvedActor(actor),
		Target: &activityevents.FediverseTarget{URL: quotedPostIRI.String()},
	})
	return nil
}
