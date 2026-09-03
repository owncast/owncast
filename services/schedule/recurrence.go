// Package schedule owns scheduled-stream recurrence logic: parsing RFC 5545
// recurrence values and materializing concrete occurrence rows from them.
// The scheduler service loop (ticker) is added by a later milestone.
package schedule

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/teambition/rrule-go"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/scheduleeventsrepository"

	log "github.com/sirupsen/logrus"
)

// A stored recurrence value is the iCalendar RFC 5545 pair of lines:
//
//	DTSTART;TZID=America/Los_Angeles:20260706T180000
//	RRULE:FREQ=WEEKLY;BYDAY=MO
//
// The TZID travels inside the value, which is what makes expansion
// DST-correct: "Monday 18:00" stays 18:00 wall time in that zone while the
// UTC instant shifts across transitions.

// Guards against admin-supplied rules that are valid RFC 5545 but would
// break the system.
const (
	// maxExpansionIterations bounds the total rule-walking work one
	// expansion may do, so a dense rule with an old DTSTART cannot burn
	// unbounded CPU before ever reaching the requested window.
	maxExpansionIterations = 500000
	// maxOccurrencesPerExpansion bounds how many occurrences one window may
	// produce (a 30-day horizon of hourly occurrences is 720; anything past
	// this is a rule nobody means for a stream schedule).
	maxOccurrencesPerExpansion = 1000
)

// ParseRecurrence parses and validates a stored recurrence value. It rejects
// empty values, values without an RRULE or DTSTART, and unknown timezones.
func ParseRecurrence(recurrence string) (*rrule.Set, error) {
	trimmed := strings.TrimSpace(recurrence)
	if trimmed == "" {
		return nil, errors.New("recurrence value is empty")
	}

	set, err := rrule.StrToRRuleSet(trimmed)
	if err != nil {
		return nil, fmt.Errorf("unable to parse recurrence value: %w", err)
	}

	if set.GetRRule() == nil {
		return nil, errors.New("recurrence value has no RRULE")
	}

	// Without an explicit DTSTART, rrule-go silently defaults the rule's
	// start to time.Now() at every parse, which makes each expansion of the
	// same value produce different instants and breaks materialization
	// idempotency. Refuse it at the door.
	if set.GetDTStart().IsZero() {
		return nil, errors.New("recurrence value has no DTSTART")
	}

	return set, nil
}

// ExpandBetween returns the occurrence start instants of a recurrence value
// inside the half-open interval [from, to). The returned times are in the
// rule's own timezone.
func ExpandBetween(recurrence string, from, to time.Time) ([]time.Time, error) {
	set, err := ParseRecurrence(recurrence)
	if err != nil {
		return nil, err
	}

	// Walk the rule's iterator instead of Set.Between so both guard caps
	// apply before a pathological rule can allocate millions of instants.
	next := set.Iterator()
	var result []time.Time
	for iterations := 0; ; iterations++ {
		if iterations > maxExpansionIterations {
			return nil, errors.New("recurrence rule is too dense to expand")
		}
		occurrence, ok := next()
		if !ok {
			break
		}
		// Occurrences arrive in chronological order: [from, to) half-open.
		if !occurrence.Before(to) {
			break
		}
		if occurrence.Before(from) {
			continue
		}
		result = append(result, occurrence)
		if len(result) > maxOccurrencesPerExpansion {
			return nil, errors.New("recurrence rule produces too many occurrences")
		}
	}
	return result, nil
}

// RecurrenceTimezone reports the IANA zone name a recurrence value's DTSTART
// carries, falling back to UTC.
func RecurrenceTimezone(set *rrule.Set) string {
	return set.GetDTStart().Location().String()
}

// MaterializeSeries expands one series over [now, now+horizon) and inserts
// any missing occurrence rows. Idempotent: existing (series, original start)
// slots are skipped by the unique index, so cancelled or moved occurrences
// are never re-created. Returns the number of rows inserted.
func MaterializeSeries(repo scheduleeventsrepository.ScheduleEventsRepository, series models.ScheduledEventSeries, now time.Time, horizon time.Duration) (int, error) {
	set, err := ParseRecurrence(series.Recurrence)
	if err != nil {
		return 0, fmt.Errorf("series %s: %w", series.ID, err)
	}

	timezone := RecurrenceTimezone(set)
	occurrences, err := ExpandBetween(series.Recurrence, now, now.Add(horizon))
	if err != nil {
		return 0, fmt.Errorf("series %s: %w", series.ID, err)
	}

	inserted := 0
	for _, occurrence := range occurrences {
		wasInserted, err := repo.AddOccurrence(series.ID, occurrence, series.Name, series.Description, series.ReminderMessage, occurrence, series.DurationMinutes, timezone)
		if err != nil {
			return inserted, fmt.Errorf("series %s: %w", series.ID, err)
		}
		if wasInserted {
			inserted++
		}
	}
	return inserted, nil
}

// MaterializeAllSeries expands every active series up to the horizon. One
// unparseable or failing series is logged and skipped, never stopping the
// loop for the others. Returns the total number of rows inserted.
func MaterializeAllSeries(repo scheduleeventsrepository.ScheduleEventsRepository, now time.Time, horizon time.Duration) (int, error) {
	series, err := repo.GetActiveSeries()
	if err != nil {
		return 0, err
	}

	inserted := 0
	for _, s := range series {
		count, err := MaterializeSeries(repo, s, now, horizon)
		inserted += count
		if err != nil {
			log.Errorf("unable to materialize scheduled stream events: %v", err)
			continue
		}
	}
	return inserted, nil
}

// RegenerateSeries removes the future, never-federated occurrences of an
// edited series and re-materializes them from the updated rule. Occurrences
// already announced to followers are left alone; the caller owns sending
// Update/Delete activities for those.
func RegenerateSeries(repo scheduleeventsrepository.ScheduleEventsRepository, series models.ScheduledEventSeries, now time.Time, horizon time.Duration) (int, error) {
	if err := repo.DeleteUnfederatedFutureEventsForSeries(series.ID, now); err != nil {
		return 0, err
	}
	return MaterializeSeries(repo, series, now, horizon)
}
