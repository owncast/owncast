package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/webserver/handlers/generated"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// maxScheduleRangeDays bounds how much schedule a single public request may
// ask for. The viewer calendar asks month-by-month.
const maxScheduleRangeDays = 366

// GetSchedule returns the public schedule of stream events, past and
// upcoming. Always an empty list while the feature is disabled: the schedule
// surface stays inert without leaking whether data exists.
func (h *Handlers) GetSchedule(w http.ResponseWriter, r *http.Request, params generated.GetScheduleParams) {
	if !h.configRepository.GetScheduleEnabled() {
		webutils.WriteResponse(w, []models.ScheduledEvent{})
		return
	}

	now := time.Now()
	from := now.AddDate(0, 0, -7)
	to := now.AddDate(0, 0, 90)
	if params.From != nil {
		from = *params.From
	}
	if params.To != nil {
		to = *params.To
	}

	if !to.After(from) {
		webutils.BadRequestHandler(w, errors.New("to must be after from"))
		return
	}
	if to.Sub(from) > maxScheduleRangeDays*24*time.Hour {
		webutils.BadRequestHandler(w, errors.New("requested schedule range is too large"))
		return
	}

	events, err := h.scheduleEventsRepository.GetEventsInRange(from, to)
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	webutils.WriteResponse(w, events)
}
