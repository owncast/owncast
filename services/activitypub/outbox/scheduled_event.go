package outbox

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/url"
	"path"
	"time"

	"code.superseriousbusiness.org/activity/streams"
	"code.superseriousbusiness.org/activity/streams/vocab"

	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/activitypub/apmodels"
)

type scheduledEventActivity interface {
	vocab.Type
	SetJSONLDId(vocab.JSONLDIdProperty)
	SetActivityStreamsActor(vocab.ActivityStreamsActorProperty)
	SetActivityStreamsTo(vocab.ActivityStreamsToProperty)
	SetActivityStreamsCc(vocab.ActivityStreamsCcProperty)
}

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
	activity := s.builder.CreateCreateActivity(path.Join("event", event.ID, "activity", "create"), s.localActorIRI())
	object := streams.NewActivityStreamsObjectProperty()
	object.AppendActivityStreamsEvent(apEvent)
	activity.SetActivityStreamsObject(object)
	s.addressEventActivity(activity)

	key := "scheduled-event:" + event.ID
	return s.persistAndSendEventActivity(apEvent, activity, key+":create", key, event.FederationVersion, false)
}

// SendScheduledEventUpdate publishes a complete replacement Event.
func (s *Service) SendScheduledEventUpdate(event models.ScheduledEvent) error {
	if !s.configRepository.GetFederationEnabled() {
		return errors.New("federation is disabled")
	}

	apEvent, err := s.builder.MakeScheduledEvent(event)
	if err != nil {
		return err
	}
	activity := streams.NewActivityStreamsUpdate()
	s.setEventActivityIdentity(activity, path.Join("event", event.ID, "activity", "update", fmt.Sprintf("%d", event.FederationVersion)))
	object := streams.NewActivityStreamsObjectProperty()
	object.AppendActivityStreamsEvent(apEvent)
	activity.SetActivityStreamsObject(object)
	s.addressEventActivity(activity)

	key := "scheduled-event:" + event.ID
	return s.persistAndSendEventActivity(apEvent, activity, key+":state", key, event.FederationVersion, event.Status == models.ScheduledEventStatusCancelled)
}

// SendScheduledEventDelete publishes a Delete containing a Tombstone and leaves
// that Tombstone at the Event's stable IRI.
func (s *Service) SendScheduledEventDelete(event models.ScheduledEvent) error {
	if !s.configRepository.GetFederationEnabled() {
		return errors.New("federation is disabled")
	}

	deletedAt := time.Now().UTC()
	eventIRI := s.builder.MakeLocalIRIForResource(path.Join("event", event.ID))
	tombstone := streams.NewActivityStreamsTombstone()
	id := streams.NewJSONLDIdProperty()
	id.Set(eventIRI)
	tombstone.SetJSONLDId(id)
	formerType := streams.NewActivityStreamsFormerTypeProperty()
	formerType.AppendXMLSchemaString("Event")
	tombstone.SetActivityStreamsFormerType(formerType)
	deleted := streams.NewActivityStreamsDeletedProperty()
	deleted.Set(deletedAt)
	tombstone.SetActivityStreamsDeleted(deleted)

	activity := streams.NewActivityStreamsDelete()
	s.setEventActivityIdentity(activity, path.Join("event", event.ID, "activity", "delete"))
	object := streams.NewActivityStreamsObjectProperty()
	object.AppendActivityStreamsTombstone(tombstone)
	activity.SetActivityStreamsObject(object)
	s.addressEventActivity(activity)

	key := "scheduled-event:" + event.ID
	return s.persistAndSendEventActivity(tombstone, activity, key+":state", key, math.MaxInt64, true)
}

// SendScheduledEventReminder publishes one stable Note for a configured reminder slot.
func (s *Service) SendScheduledEventReminder(event models.ScheduledEvent, reminderNumber int, message string) error {
	if !s.configRepository.GetFederationEnabled() {
		return errors.New("federation is disabled")
	}

	actorIRI := s.localActorIRI()
	noteIRI := s.builder.MakeLocalIRIForResource(path.Join("event", event.ID, "reminder", fmt.Sprintf("%d", reminderNumber)))
	note := apmodels.MakeNote(message, noteIRI, actorIRI)
	to, cc := s.getAddressingToFollowers()
	note.SetActivityStreamsTo(to)
	note.SetActivityStreamsCc(cc)
	inReplyTo := streams.NewActivityStreamsInReplyToProperty()
	inReplyTo.AppendIRI(s.builder.MakeLocalIRIForResource(path.Join("event", event.ID)))
	note.SetActivityStreamsInReplyTo(inReplyTo)

	activity := s.builder.CreateCreateActivity(path.Join("event", event.ID, "reminder", fmt.Sprintf("%d", reminderNumber), "activity"), actorIRI)
	object := streams.NewActivityStreamsObjectProperty()
	object.AppendActivityStreamsNote(note)
	activity.SetActivityStreamsObject(object)
	activity.SetActivityStreamsTo(to)
	activity.SetActivityStreamsCc(cc)

	notePayload, err := apmodels.Serialize(note)
	if err != nil {
		return err
	}
	activityPayload, err := apmodels.Serialize(activity)
	if err != nil {
		return err
	}
	return s.persistAndSendEventReminder(event, reminderNumber, noteIRI.String(), note.GetTypeName(), notePayload, activityPayload)
}

func (s *Service) persistAndSendEventReminder(event models.ScheduledEvent, reminderNumber int, noteIRI, noteType string, notePayload, activityPayload []byte) error {
	key := "scheduled-event:" + event.ID
	deliveries, err := s.followerDeliveries(activityPayload, fmt.Sprintf("%s:reminder:%d", key, reminderNumber), key, event.FederationVersion, false)
	if err != nil {
		return err
	}

	return s.workerpool.WithOrderingKey(key, func() error {
		ctx := context.Background()
		ds := s.persistence.Datastore()
		ds.DbLock.Lock()
		defer ds.DbLock.Unlock()
		tx, err := ds.DB.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer func() { _ = tx.Rollback() }()
		queries := db.New(tx)
		scheduled, err := queries.GetStreamEvent(ctx, event.ID)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		if scheduled.Status != models.ScheduledEventStatusScheduled || scheduled.FederationVersion != event.FederationVersion {
			return nil
		}
		current, err := queries.GetOutboxObjectState(ctx, s.builder.MakeLocalIRIForResource(path.Join("event", event.ID)).String())
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		if current.Type != "Event" || current.CoalesceVersion != event.FederationVersion {
			return nil
		}
		if err := queries.UpsertOutboxObject(ctx, db.UpsertOutboxObjectParams{
			Iri:              noteIRI,
			Value:            notePayload,
			Type:             noteType,
			LiveNotification: sql.NullBool{Bool: false, Valid: true},
		}); err != nil {
			return err
		}
		for _, delivery := range deliveries {
			if err := s.workerpool.EnqueueTx(ctx, tx, delivery); err != nil {
				return err
			}
		}
		return tx.Commit()
	})
}

func (s *Service) persistAndSendEventActivity(object, activity vocab.Type, coalesceKey, orderingKey string, coalesceVersion int64, dropReminders bool) error {
	objectIRI, err := apmodels.GetIRIStringFromJSONLDIdProperty(object.GetJSONLDId())
	if err != nil {
		return err
	}
	activityIRI, err := apmodels.GetIRIStringFromJSONLDIdProperty(activity.GetJSONLDId())
	if err != nil {
		return err
	}
	objectPayload, err := apmodels.SerializeEvent(object)
	if err != nil {
		return err
	}
	activityPayload, err := apmodels.SerializeEvent(activity)
	if err != nil {
		return err
	}
	deliveries, err := s.followerDeliveries(activityPayload, coalesceKey, orderingKey, coalesceVersion, true)
	if err != nil {
		return err
	}

	return s.workerpool.WithOrderingKey(orderingKey, func() error {
		ctx := context.Background()
		ds := s.persistence.Datastore()
		ds.DbLock.Lock()
		defer ds.DbLock.Unlock()
		tx, err := ds.DB.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer func() { _ = tx.Rollback() }()
		queries := db.New(tx)
		if dropReminders {
			if err := queries.DeleteActivityPubReminderDeliveries(ctx, sql.NullString{String: orderingKey, Valid: true}); err != nil {
				return err
			}
		}
		if _, err := queries.UpsertVersionedOutboxObject(ctx, db.UpsertVersionedOutboxObjectParams{
			Iri:              objectIRI,
			Value:            objectPayload,
			Type:             object.GetTypeName(),
			LiveNotification: sql.NullBool{Bool: false, Valid: true},
			CoalesceVersion:  coalesceVersion,
		}); errors.Is(err, sql.ErrNoRows) {
			return nil
		} else if err != nil {
			return err
		}
		if err := queries.UpsertOutboxObject(ctx, db.UpsertOutboxObjectParams{
			Iri:              activityIRI,
			Value:            activityPayload,
			Type:             activity.GetTypeName(),
			LiveNotification: sql.NullBool{Bool: false, Valid: true},
		}); err != nil {
			return err
		}
		for _, delivery := range deliveries {
			if err := s.workerpool.EnqueueTx(ctx, tx, delivery); err != nil {
				return err
			}
		}
		return tx.Commit()
	})
}

func (s *Service) localActorIRI() *url.URL {
	return s.builder.MakeLocalIRIForAccount(s.configRepository.GetDefaultFederationUsername())
}

func (s *Service) setEventActivityIdentity(activity scheduledEventActivity, resource string) {
	id := streams.NewJSONLDIdProperty()
	id.Set(s.builder.MakeLocalIRIForResource(resource))
	activity.SetJSONLDId(id)
	actor := streams.NewActivityStreamsActorProperty()
	actor.AppendIRI(s.localActorIRI())
	activity.SetActivityStreamsActor(actor)
}

func (s *Service) addressEventActivity(activity scheduledEventActivity) {
	to, cc := s.getAddressingToFollowers()
	activity.SetActivityStreamsTo(to)
	activity.SetActivityStreamsCc(cc)
}
