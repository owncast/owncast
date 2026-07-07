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

// Service owns the scheduled streams background loop: materializing
// occurrence rows from recurring series and answering "what's next" for the
// status endpoint. Construct with New in main.go and inject into consumers.
type Service struct {
	repo scheduleeventsrepository.ScheduleEventsRepository

	tickerMutex sync.Mutex
	ticker      *time.Ticker
	tickerDone  chan bool

	upcomingMutex     sync.Mutex
	upcomingEvent     *models.ScheduledEvent
	upcomingFetchedAt time.Time
}

// Deps lists everything a schedule Service consumes.
type Deps struct {
	ScheduleEventsRepository scheduleeventsrepository.ScheduleEventsRepository
}

// New constructs the schedule service.
func New(deps Deps) *Service {
	return &Service{
		repo: deps.ScheduleEventsRepository,
	}
}

// Start begins the background loop. Runs even while the feature toggle is
// off so the schedule is current the moment an admin enables it; the HTTP
// layer owns hiding disabled state from the public.
func (s *Service) Start() {
	s.tickerMutex.Lock()
	defer s.tickerMutex.Unlock()

	if s.ticker != nil {
		log.Debugln("Schedule service already running")
		return
	}

	s.ticker = time.NewTicker(tickInterval)
	s.tickerDone = make(chan bool)

	done := s.tickerDone
	ticker := s.ticker

	go func() {
		s.tick()

		for {
			select {
			case <-ticker.C:
				s.tick()
			case <-done:
				return
			}
		}
	}()

	log.Infof("Started schedule service (%s interval, %s materialization horizon)", tickInterval, MaterializationHorizon)
}

// Stop halts the schedule service if it is running.
func (s *Service) Stop() {
	s.tickerMutex.Lock()
	defer s.tickerMutex.Unlock()

	if s.ticker != nil {
		s.ticker.Stop()
		close(s.tickerDone)
		s.ticker = nil
		s.tickerDone = nil
		log.Infoln("Stopped schedule service")
	}
}

func (s *Service) tick() {
	inserted, err := MaterializeAllSeries(s.repo, time.Now(), MaterializationHorizon)
	if err != nil {
		log.Errorf("unable to materialize scheduled stream events: %v", err)
	}
	if inserted > 0 {
		log.Debugf("Materialized %d scheduled stream event(s)", inserted)
	}

	s.Refresh()
}

// Refresh invalidates the cached next-event answer. Called after admin
// mutations so the status API reflects changes immediately.
func (s *Service) Refresh() {
	s.upcomingMutex.Lock()
	defer s.upcomingMutex.Unlock()
	s.upcomingFetchedAt = time.Time{}
}

// GetUpcomingEvent returns the current-or-next scheduled (not cancelled)
// occurrence: an event still inside its running window (start through
// start+duration) counts, so chat stays open when a stream starts late. The
// end-time predicate lives in SQL, so one row is always the exact answer.
// Returns nil when nothing is current or upcoming. Backed by a short-lived
// cache: the status endpoint is polled every few seconds by every viewer.
func (s *Service) GetUpcomingEvent() *models.ScheduledEvent {
	s.upcomingMutex.Lock()
	defer s.upcomingMutex.Unlock()

	if time.Since(s.upcomingFetchedAt) < upcomingCacheTTL {
		return s.upcomingEvent
	}

	now := time.Now()
	events, err := s.repo.GetCurrentOrUpcomingEvents(now, 1)

	// Cache the answer, including a failed one: the status endpoint polls
	// every few seconds, and a transient database error should back off
	// for the TTL instead of re-querying and logging on every poll.
	s.upcomingFetchedAt = now
	s.upcomingEvent = nil
	if err != nil {
		log.Errorf("unable to fetch the next scheduled stream event: %v", err)
		return nil
	}
	if len(events) > 0 {
		s.upcomingEvent = &events[0]
	}
	return s.upcomingEvent
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
