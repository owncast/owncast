package schedule

import (
	"sync"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/scheduleeventsrepository"
	log "github.com/sirupsen/logrus"
)

const (
	// MaterializationHorizon is how far ahead occurrence rows exist for
	// recurring series. Far enough for a month-view calendar.
	MaterializationHorizon = 30 * 24 * time.Hour

	// chatOpenLeadTime is how long before a scheduled event chat is
	// treated as open for early arrivals.
	chatOpenLeadTime = 10 * time.Minute

	// tickInterval is how often the scheduler materializes occurrences.
	// Each tick is idempotent and driven purely from table state, so a
	// missed or repeated tick never corrupts anything.
	tickInterval = 1 * time.Minute

	// upcomingCacheTTL bounds how stale the next-event answer served to
	// the 5s status poll may be. Admin mutations bypass it via Refresh.
	upcomingCacheTTL = 10 * time.Second
)

var (
	ticker      *time.Ticker
	tickerDone  chan bool
	tickerMutex sync.Mutex

	upcomingMutex     sync.Mutex
	upcomingEvent     *models.ScheduledEvent
	upcomingFetchedAt time.Time
)

// Start begins the scheduled streams background loop: materializing
// occurrence rows from recurring series. Runs even while the feature toggle
// is off so the schedule is current the moment an admin enables it; the
// HTTP layer owns hiding disabled state from the public.
func Start() {
	tickerMutex.Lock()
	defer tickerMutex.Unlock()

	if ticker != nil {
		log.Debugln("Schedule service already running")
		return
	}

	ticker = time.NewTicker(tickInterval)
	tickerDone = make(chan bool)

	done := tickerDone
	t := ticker

	go func() {
		tick()

		for {
			select {
			case <-t.C:
				tick()
			case <-done:
				return
			}
		}
	}()

	log.Infof("Started schedule service (%s interval, %s materialization horizon)", tickInterval, MaterializationHorizon)
}

// Stop halts the schedule service if it is running.
func Stop() {
	tickerMutex.Lock()
	defer tickerMutex.Unlock()

	if ticker != nil {
		ticker.Stop()
		close(tickerDone)
		ticker = nil
		tickerDone = nil
		log.Infoln("Stopped schedule service")
	}
}

func tick() {
	repo := scheduleeventsrepository.Get()
	if repo == nil {
		return
	}

	inserted, err := MaterializeAllSeries(repo, time.Now(), MaterializationHorizon)
	if err != nil {
		log.Errorf("unable to materialize scheduled stream events: %v", err)
	}
	if inserted > 0 {
		log.Debugf("Materialized %d scheduled stream event(s)", inserted)
	}

	Refresh()
}

// Refresh invalidates the cached next-event answer. Called after admin
// mutations so the status API reflects changes immediately.
func Refresh() {
	upcomingMutex.Lock()
	defer upcomingMutex.Unlock()
	upcomingFetchedAt = time.Time{}
}

// GetUpcomingEvent returns the current-or-next scheduled (not cancelled)
// occurrence: an event still inside its running window (start through
// start+duration) counts, so chat stays open when a stream starts late. The
// end-time predicate lives in SQL, so one row is always the exact answer.
// Returns nil when nothing is current or upcoming. Backed by a short-lived
// cache: the status endpoint is polled every few seconds by every viewer.
func GetUpcomingEvent() *models.ScheduledEvent {
	upcomingMutex.Lock()
	defer upcomingMutex.Unlock()

	if time.Since(upcomingFetchedAt) < upcomingCacheTTL {
		return upcomingEvent
	}

	repo := scheduleeventsrepository.Get()
	if repo == nil {
		return nil
	}

	now := time.Now()
	events, err := repo.GetCurrentOrUpcomingEvents(now, 1)
	if err != nil {
		log.Errorf("unable to fetch the next scheduled stream event: %v", err)
		return nil
	}

	upcomingFetchedAt = now
	upcomingEvent = nil
	if len(events) > 0 {
		upcomingEvent = &events[0]
	}
	return upcomingEvent
}

// IsChatOpenForEvent reports whether the pre-event chat window is open: the
// event starts within chatOpenLeadTime, or has started and not yet run past
// its duration.
func IsChatOpenForEvent(event *models.ScheduledEvent, now time.Time) bool {
	if event == nil {
		return false
	}
	windowStart := event.StartTime.Add(-chatOpenLeadTime)
	windowEnd := event.StartTime.Add(time.Duration(event.DurationMinutes) * time.Minute)
	return !now.Before(windowStart) && now.Before(windowEnd)
}
