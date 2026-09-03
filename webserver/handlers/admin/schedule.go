package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/schedule"
	"github.com/owncast/owncast/webserver/handlers/generated"
	webutils "github.com/owncast/owncast/webserver/utils"
)

const defaultEventDurationMinutes = 60

// maxEventDurationMinutes bounds an event to a week. Anything longer is a
// data-entry mistake, not a stream.
const maxEventDurationMinutes = 7 * 24 * 60

type adminScheduleResponse struct {
	Series []models.ScheduledEventSeries `json:"series"`
	Events []models.ScheduledEvent       `json:"events"`
}

// GetAdminSchedule returns every series and the occurrences an admin cares
// about managing.
func (a *Admin) GetAdminSchedule(w http.ResponseWriter, r *http.Request) {
	a.writeAdminSchedule(w)
}

func (a *Admin) writeAdminSchedule(w http.ResponseWriter) {
	series, err := a.scheduleEventsRepository.GetAllSeries()
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	// ponytail: fixed admin window (3 months back, 2 months out) instead of
	// pagination; widen when somebody actually schedules beyond it.
	now := time.Now()
	events, err := a.scheduleEventsRepository.GetEventsInRange(now.AddDate(0, -3, 0), now.AddDate(0, 2, 0))
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	webutils.WriteResponse(w, adminScheduleResponse{Series: series, Events: events})
}

// UpsertScheduledEvent creates or updates a one-off event (start set) or a
// recurring series (recurrence set).
func (a *Admin) UpsertScheduledEvent(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	var request generated.ScheduledEventInput
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		webutils.BadRequestHandler(w, err)
		return
	}

	name := strings.TrimSpace(request.Name)
	if name == "" {
		webutils.BadRequestHandler(w, errors.New("name is required"))
		return
	}

	if request.DurationMinutes != nil && (*request.DurationMinutes <= 0 || *request.DurationMinutes > maxEventDurationMinutes) {
		webutils.BadRequestHandler(w, errors.New("durationMinutes is out of range"))
		return
	}

	recurrence := ""
	if request.Recurrence != nil {
		recurrence = strings.TrimSpace(*request.Recurrence)
	}
	if recurrence != "" && request.Start != nil {
		webutils.BadRequestHandler(w, errors.New("start and recurrence are mutually exclusive"))
		return
	}

	var err error
	if request.Id == nil {
		err = a.createScheduledEvent(w, name, recurrence, request)
	} else {
		err = a.updateScheduledEvent(w, *request.Id, name, recurrence, request)
	}
	if err != nil {
		// The helper already wrote the HTTP error response.
		return
	}

	a.schedule.Refresh()
	a.writeAdminSchedule(w)
}

// createScheduledEvent writes an HTTP error response and returns a non-nil
// error when creation fails.
func (a *Admin) createScheduledEvent(w http.ResponseWriter, name, recurrence string, request generated.ScheduledEventInput) error {
	description := ""
	if request.Description != nil {
		description = *request.Description
	}
	duration := defaultEventDurationMinutes
	if request.DurationMinutes != nil {
		duration = *request.DurationMinutes
	}

	if recurrence != "" {
		return a.createSeries(w, name, description, recurrence, duration)
	}
	return a.createOneOffEvent(w, name, description, duration, request)
}

func (a *Admin) createSeries(w http.ResponseWriter, name, description, recurrence string, duration int) error {
	if _, err := schedule.ParseRecurrence(recurrence); err != nil {
		webutils.BadRequestHandler(w, err)
		return err
	}
	seriesID, err := a.scheduleEventsRepository.AddSeries(name, description, recurrence, duration)
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return err
	}
	series, err := a.scheduleEventsRepository.GetSeries(seriesID)
	if err != nil || series == nil {
		if err == nil {
			err = errors.New("unable to load created series")
		}
		webutils.InternalErrorHandler(w, err)
		return err
	}
	if _, err := schedule.MaterializeSeries(a.scheduleEventsRepository, *series, time.Now(), schedule.MaterializationHorizon); err != nil {
		webutils.InternalErrorHandler(w, err)
		return err
	}
	return nil
}

func (a *Admin) createOneOffEvent(w http.ResponseWriter, name, description string, duration int, request generated.ScheduledEventInput) error {
	if request.Start == nil {
		err := errors.New("start is required for a one-off event")
		webutils.BadRequestHandler(w, err)
		return err
	}
	if request.Start.Before(time.Now()) {
		err := errors.New("start must be in the future")
		webutils.BadRequestHandler(w, err)
		return err
	}

	timezone := "UTC"
	if request.Timezone != nil && *request.Timezone != "" {
		if _, err := time.LoadLocation(*request.Timezone); err != nil {
			webutils.BadRequestHandler(w, errors.New("unknown timezone"))
			return err
		}
		timezone = *request.Timezone
	}

	if _, err := a.scheduleEventsRepository.AddOneOffEvent(name, description, *request.Start, duration, timezone); err != nil {
		webutils.InternalErrorHandler(w, err)
		return err
	}
	return nil
}

// updateScheduledEvent writes an HTTP error response and returns a non-nil
// error when the update fails. Omitted optional fields keep their stored
// values, and field/id-type combinations that cannot apply are rejected
// instead of silently dropped.
func (a *Admin) updateScheduledEvent(w http.ResponseWriter, id, name, recurrence string, request generated.ScheduledEventInput) error {
	if request.Timezone != nil {
		err := errors.New("timezone cannot be changed on an update")
		webutils.BadRequestHandler(w, err)
		return err
	}

	series, err := a.scheduleEventsRepository.GetSeries(id)
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return err
	}
	if series != nil {
		return a.updateSeries(w, *series, name, recurrence, request)
	}

	event, err := a.scheduleEventsRepository.GetEvent(id)
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return err
	}
	if event == nil {
		err := errors.New("no scheduled event or series with that id")
		webutils.BadRequestHandler(w, err)
		return err
	}
	return a.updateOneOffEvent(w, *event, name, recurrence, request)
}

func (a *Admin) updateSeries(w http.ResponseWriter, series models.ScheduledEventSeries, name, recurrence string, request generated.ScheduledEventInput) error {
	if request.Start != nil {
		err := errors.New("a recurring series has no single start, edit its recurrence instead")
		webutils.BadRequestHandler(w, err)
		return err
	}
	if recurrence == "" {
		recurrence = series.Recurrence
	}
	if _, err := schedule.ParseRecurrence(recurrence); err != nil {
		webutils.BadRequestHandler(w, err)
		return err
	}

	description := series.Description
	if request.Description != nil {
		description = *request.Description
	}
	duration := series.DurationMinutes
	if request.DurationMinutes != nil {
		duration = *request.DurationMinutes
	}

	if err := a.scheduleEventsRepository.UpdateSeries(series.ID, name, description, recurrence, duration); err != nil {
		webutils.InternalErrorHandler(w, err)
		return err
	}
	updated, err := a.scheduleEventsRepository.GetSeries(series.ID)
	if err != nil || updated == nil {
		if err == nil {
			err = errors.New("unable to load updated series")
		}
		webutils.InternalErrorHandler(w, err)
		return err
	}
	if _, err := schedule.RegenerateSeries(a.scheduleEventsRepository, *updated, time.Now(), schedule.MaterializationHorizon); err != nil {
		webutils.InternalErrorHandler(w, err)
		return err
	}
	return nil
}

func (a *Admin) updateOneOffEvent(w http.ResponseWriter, event models.ScheduledEvent, name, recurrence string, request generated.ScheduledEventInput) error {
	if recurrence != "" {
		err := errors.New("a one-off event cannot become a recurring series, create a new series instead")
		webutils.BadRequestHandler(w, err)
		return err
	}
	// Validate everything before the first write so a rejected update never
	// half-applies.
	if request.Start != nil && request.Start.Before(time.Now()) {
		err := errors.New("start must be in the future")
		webutils.BadRequestHandler(w, err)
		return err
	}

	description := event.Description
	if request.Description != nil {
		description = *request.Description
	}
	duration := event.DurationMinutes
	if request.DurationMinutes != nil {
		duration = *request.DurationMinutes
	}

	if err := a.scheduleEventsRepository.UpdateEventDetails(event.ID, name, description, duration); err != nil {
		webutils.InternalErrorHandler(w, err)
		return err
	}
	if request.Start != nil {
		if err := a.scheduleEventsRepository.MoveEvent(event.ID, *request.Start); err != nil {
			webutils.InternalErrorHandler(w, err)
			return err
		}
	}
	return nil
}

// DeleteScheduledEvent deletes or cancels an event or series. Cancelling an
// occurrence keeps the row (shown as cancelled and never re-materialized);
// cancelling a series pauses it. Deleting a series also removes its future
// never-federated occurrences while past ones stay for the record.
func (a *Admin) DeleteScheduledEvent(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	var request generated.DeleteScheduledEventJSONBody
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		webutils.BadRequestHandler(w, err)
		return
	}
	cancel := request.Cancel != nil && *request.Cancel

	series, err := a.scheduleEventsRepository.GetSeries(request.Id)
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	if series != nil {
		if cancel {
			err = a.scheduleEventsRepository.SetSeriesActive(request.Id, false)
		} else {
			if err := a.scheduleEventsRepository.DeleteUnfederatedFutureEventsForSeries(request.Id, time.Now()); err != nil {
				webutils.InternalErrorHandler(w, err)
				return
			}
			err = a.scheduleEventsRepository.DeleteSeries(request.Id)
		}
		if err != nil {
			webutils.InternalErrorHandler(w, err)
			return
		}
	} else {
		event, err := a.scheduleEventsRepository.GetEvent(request.Id)
		if err != nil {
			webutils.InternalErrorHandler(w, err)
			return
		}
		if event == nil {
			webutils.BadRequestHandler(w, errors.New("no scheduled event or series with that id"))
			return
		}
		if cancel {
			err = a.scheduleEventsRepository.CancelEvent(request.Id)
		} else {
			err = a.scheduleEventsRepository.DeleteEvent(request.Id)
		}
		if err != nil {
			webutils.InternalErrorHandler(w, err)
			return
		}
	}

	a.schedule.Refresh()
	a.writeAdminSchedule(w)
}

// PreviewScheduleRecurrence expands a recurrence rule server-side so the
// admin form can show exactly what the materializer will produce. One source
// of truth: the same parser and expander the scheduler uses.
func (a *Admin) PreviewScheduleRecurrence(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	var request generated.PreviewScheduleRecurrenceJSONBody
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		webutils.BadRequestHandler(w, err)
		return
	}

	now := time.Now()
	occurrences, err := schedule.ExpandBetween(request.Recurrence, now, now.AddDate(1, 0, 0))
	if err != nil {
		webutils.BadRequestHandler(w, err)
		return
	}

	const previewCount = 5
	if len(occurrences) > previewCount {
		occurrences = occurrences[:previewCount]
	}
	if occurrences == nil {
		occurrences = []time.Time{}
	}

	webutils.WriteResponse(w, map[string][]time.Time{"occurrences": occurrences})
}

// SetScheduleEnabled turns the scheduled streams feature on or off.
func (a *Admin) SetScheduleEnabled(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	enabled, ok := configValue.Value.(bool)
	if !ok {
		webutils.BadRequestHandler(w, errors.New("value must be a boolean"))
		return
	}

	if err := a.configRepository.SetScheduleEnabled(enabled); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "schedule feature updated")
}

// SetScheduleShowCountdown enables or disables the viewer's event countdown.
func (a *Admin) SetScheduleShowCountdown(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	showCountdown, ok := configValue.Value.(bool)
	if !ok {
		webutils.BadRequestHandler(w, errors.New("value must be a boolean"))
		return
	}

	if err := a.configRepository.SetScheduleShowCountdown(showCountdown); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "schedule countdown setting updated")
}

// SetScheduleChatOpenMinutes sets how many minutes before an event chat opens.
func (a *Admin) SetScheduleChatOpenMinutes(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	value, ok := configValue.Value.(float64)
	if !ok || value != float64(int(value)) {
		webutils.BadRequestHandler(w, errors.New("value must be a whole number of minutes"))
		return
	}
	minutes := int(value)
	switch minutes {
	case 0, 5, 10, 30, 60:
	default:
		webutils.BadRequestHandler(w, errors.New("value must be 0, 5, 10, 30, or 60 minutes"))
		return
	}

	if err := a.configRepository.SetScheduleChatOpenMinutes(minutes); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "schedule chat open setting updated")
}

// SetScheduleReminderMessage sets the message posted to the Fediverse before
// a scheduled event starts. Empty disables reminders.
func (a *Admin) SetScheduleReminderMessage(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	message, ok := configValue.Value.(string)
	if !ok {
		webutils.BadRequestHandler(w, errors.New("value must be a string"))
		return
	}

	if err := a.configRepository.SetScheduleReminderMessage(message); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "schedule reminder message updated")
}
