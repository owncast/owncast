package outbox

import (
	"errors"

	"code.superseriousbusiness.org/activity/streams"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/activitypub/apmodels"
)

// SendScheduledEvent publishes a Create containing an ActivityPub Event to all
// approved follower inboxes.
func (s *Service) SendScheduledEvent(event models.ScheduledEvent) error {
	if !s.configRepository.GetFederationEnabled() {
		return errors.New("federation is disabled")
	}

	apEvent, err := s.builder.MakeScheduledEvent(event)
	if err != nil {
		return err
	}
	objectIRI, err := apmodels.GetIRIStringFromJSONLDIdProperty(apEvent.GetJSONLDId())
	if err != nil {
		return err
	}

	activity := s.builder.CreateCreateActivity("event/"+event.ID+"/"+shortid.MustGenerate(), s.builder.MakeLocalIRIForAccount(s.configRepository.GetDefaultFederationUsername()))
	object := streams.NewActivityStreamsObjectProperty()
	object.AppendActivityStreamsEvent(apEvent)
	activity.SetActivityStreamsObject(object)

	to, cc := s.getAddressingToFollowers()
	activity.SetActivityStreamsTo(to)
	activity.SetActivityStreamsCc(cc)

	objectPayload, err := apmodels.SerializeEvent(apEvent)
	if err != nil {
		return err
	}
	activityPayload, err := apmodels.SerializeEvent(activity)
	if err != nil {
		return err
	}
	if err := s.persistence.UpsertOutboxObject(objectIRI, objectPayload, apEvent.GetTypeName()); err != nil {
		return err
	}
	return s.SendToFollowers(activityPayload)
}
