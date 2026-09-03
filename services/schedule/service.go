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

	// tickInterval is how often the scheduler materializes occurrences.
	// Each tick is idempotent and driven purely from table state, so a
	// missed or repeated tick never corrupts anything.
	tickInterval = 1 * time.Minute

	// upcomingCacheTTL bounds how stale the next-event answer served to
	// the 5s status poll may be. Admin mutations bypass it via Refresh.
	upcomingCacheTTL = 10 * time.Second

	// MissedEventGracePeriod is how long after the scheduled start chat stays
	// open while waiting for a late stream.
	MissedEventGracePeriod = 10 * time.Minute

	// MissedEventWarningLeadTime leaves the chat shutdown message visible.
	MissedEventWarningLeadTime = 1 * time.Minute

	// MissedEventChatMessage is sent when a pre-opened event never starts.
	MissedEventChatMessage = "This scheduled live stream never started, so chat is being turned off. Thanks for coming!"
)

type chatWindowState struct {
	event       *models.ScheduledEvent
	preEnabled  bool
	started     bool
	warningSent bool
	cancelled   bool
	expiresAt   time.Time
}

// Service owns the scheduled streams background loop: materializing
// occurrence rows from recurring series and answering "what's next" for the
// status endpoint. Construct with New in main.go and inject into consumers.
type Service struct {
	repo scheduleeventsrepository.ScheduleEventsRepository

	getStatus            func() models.Status
	getChatOpenMinutes   func() int
	onMissedEventWarning func(*models.ScheduledEvent)
	chatWindowStates     map[string]chatWindowState

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
	GetStatus                func() models.Status
	GetChatOpenMinutes       func() int
	OnMissedEventWarning     func(*models.ScheduledEvent)
}

// New constructs the schedule service.
func New(deps Deps) *Service {
	return &Service{
		repo:                 deps.ScheduleEventsRepository,
		getStatus:            deps.GetStatus,
		getChatOpenMinutes:   deps.GetChatOpenMinutes,
		onMissedEventWarning: deps.OnMissedEventWarning,
		chatWindowStates:     make(map[string]chatWindowState),
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
	now := time.Now()
	inserted, err := MaterializeAllSeries(s.repo, now, MaterializationHorizon)
	if err != nil {
		log.Errorf("unable to materialize scheduled stream events: %v", err)
	}
	if inserted > 0 {
		log.Debugf("Materialized %d scheduled stream event(s)", inserted)
	}

	s.Refresh()
	s.updateChatWindow(now)
}

func (s *Service) updateChatWindow(now time.Time) {
	if s.getStatus == nil || s.getChatOpenMinutes == nil {
		return
	}

	event := s.GetUpcomingEvent()
	s.cleanupChatWindowStates(now)
	status := s.getStatus()
	if event != nil {
		state := s.chatWindowStates[event.ID]
		state.event = event
		leadTime := time.Duration(s.getChatOpenMinutes()) * time.Minute
		state.preEnabled = state.preEnabled || shouldPreEnableChat(event, status, now, leadTime)
		state.started = state.started || eventHasStarted(event, status, now)
		state.expiresAt = event.StartTime.Add(time.Duration(event.DurationMinutes) * time.Minute)
		graceExpiresAt := event.StartTime.Add(MissedEventGracePeriod + tickInterval)
		if state.expiresAt.Before(graceExpiresAt) {
			state.expiresAt = graceExpiresAt
		}
		s.chatWindowStates[event.ID] = state
	}

	s.processChatWindows(now, status)
}

func (s *Service) processChatWindows(now time.Time, status models.Status) {
	for id, state := range s.chatWindowStates {
		event := state.event
		if event == nil || !state.preEnabled || state.started || status.Online {
			continue
		}

		if !state.warningSent && IsEventMissedWarning(event, now) {
			state.warningSent = true
			if s.onMissedEventWarning != nil {
				s.onMissedEventWarning(event)
			}
		}

		if !state.cancelled && IsEventMissed(event, now) && s.repo != nil {
			if err := s.repo.CancelEvent(id); err != nil {
				log.Errorf("unable to cancel missed scheduled stream event %s: %v", id, err)
			} else {
				state.cancelled = true
				s.Refresh()
			}
		}
		s.chatWindowStates[id] = state
	}
}

func (s *Service) cleanupChatWindowStates(now time.Time) {
	for id, state := range s.chatWindowStates {
		if !now.Before(state.expiresAt) {
			delete(s.chatWindowStates, id)
		}
	}
}

func shouldPreEnableChat(event *models.ScheduledEvent, status models.Status, now time.Time, leadTime time.Duration) bool {
	return leadTime > 0 &&
		now.Before(event.StartTime) &&
		!now.Before(event.StartTime.Add(-leadTime)) &&
		!status.Online
}

func eventHasStarted(event *models.ScheduledEvent, status models.Status, now time.Time) bool {
	if now.Before(event.StartTime) {
		return false
	}
	if status.Online {
		return true
	}
	return status.LastConnectTime != nil && !status.LastConnectTime.Time.Before(event.StartTime)
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

// IsChatOpen reports whether a scheduled event currently permits chat while
// the stream is offline.
func (s *Service) IsChatOpen() bool {
	if s.getStatus == nil || s.getChatOpenMinutes == nil {
		return false
	}
	event := s.GetUpcomingEvent()
	if event == nil {
		return false
	}
	now := time.Now()
	return !IsEventMissed(event, now) &&
		IsChatOpenForEvent(event, now, time.Duration(s.getChatOpenMinutes())*time.Minute)
}

// IsEventMissed reports whether the late-start grace period has elapsed.
func IsEventMissed(event *models.ScheduledEvent, now time.Time) bool {
	return event != nil && !now.Before(event.StartTime.Add(MissedEventGracePeriod))
}

// IsEventMissedWarning reports whether it is time to warn before shutdown.
func IsEventMissedWarning(event *models.ScheduledEvent, now time.Time) bool {
	return event != nil && !now.Before(event.StartTime.Add(MissedEventGracePeriod-MissedEventWarningLeadTime))
}

// IsChatOpenForEvent reports whether the pre-event chat window is open:
// the event starts within chatOpenLeadTime, or has started and not yet run
// past its duration.
func IsChatOpenForEvent(event *models.ScheduledEvent, now time.Time, chatOpenLeadTime time.Duration) bool {
	if event == nil {
		return false
	}
	windowStart := event.StartTime.Add(-chatOpenLeadTime)
	windowEnd := event.StartTime.Add(time.Duration(event.DurationMinutes) * time.Minute)
	return !now.Before(windowStart) && now.Before(windowEnd)
}
