package schedule

import (
	"fmt"
	"sync"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/scheduleeventsrepository"
	"github.com/owncast/owncast/services/webhooks"
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
	// ScheduledEventWarningLeadTime is how long before an event its webhook fires.
	ScheduledEventWarningLeadTime = 10 * time.Minute

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
	repo                             scheduleeventsrepository.ScheduleEventsRepository
	getStatus                        func() models.Status
	getChatOpenMinutes               func() int
	onMissedEventWarning             func(*models.ScheduledEvent)
	getScheduleEnabled               func() bool
	getFederationEnabled             func() bool
	federateScheduledEvent           func(models.ScheduledEvent) error
	getScheduleReminderMessage       func() string
	getScheduleFirstReminderMinutes  func() int
	getScheduleSecondReminderMinutes func() int
	getServerURL                     func() string
	notifyScheduledEvent             func(string)
	setStreamTitle                   func(string) error
	webhooks                         *webhooks.Service
	chatWindowStates                 map[string]chatWindowState

	tickerMutex sync.Mutex
	ticker      *time.Ticker
	tickerDone  chan bool

	upcomingMutex     sync.Mutex
	upcomingEvent     *models.ScheduledEvent
	upcomingFetchedAt time.Time
}

// Deps lists everything a schedule Service consumes.
type Deps struct {
	ScheduleEventsRepository         scheduleeventsrepository.ScheduleEventsRepository
	GetStatus                        func() models.Status
	GetChatOpenMinutes               func() int
	OnMissedEventWarning             func(*models.ScheduledEvent)
	GetScheduleEnabled               func() bool
	GetFederationEnabled             func() bool
	FederateScheduledEvent           func(models.ScheduledEvent) error
	GetScheduleReminderMessage       func() string
	GetScheduleFirstReminderMinutes  func() int
	GetScheduleSecondReminderMinutes func() int
	GetServerURL                     func() string
	NotifyScheduledEvent             func(string)
	SetStreamTitle                   func(string) error
	Webhooks                         *webhooks.Service
}

// New constructs the schedule service.
func New(deps Deps) *Service {
	return &Service{
		repo:                             deps.ScheduleEventsRepository,
		getStatus:                        deps.GetStatus,
		getChatOpenMinutes:               deps.GetChatOpenMinutes,
		onMissedEventWarning:             deps.OnMissedEventWarning,
		getScheduleEnabled:               deps.GetScheduleEnabled,
		getFederationEnabled:             deps.GetFederationEnabled,
		federateScheduledEvent:           deps.FederateScheduledEvent,
		getScheduleReminderMessage:       deps.GetScheduleReminderMessage,
		getScheduleFirstReminderMinutes:  deps.GetScheduleFirstReminderMinutes,
		getScheduleSecondReminderMinutes: deps.GetScheduleSecondReminderMinutes,
		getServerURL:                     deps.GetServerURL,
		notifyScheduledEvent:             deps.NotifyScheduledEvent,
		setStreamTitle:                   deps.SetStreamTitle,
		webhooks:                         deps.Webhooks,
		chatWindowStates:                 make(map[string]chatWindowState),
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

	s.publishPendingEvents(now)
	s.Refresh()
	s.updateChatWindow(now)
	s.sendScheduledEventReminders(now)

	s.sendScheduledEventWebhooks(now)
}

// PublishPendingEvents immediately queues new occurrences for ActivityPub
// delivery. Admin mutations call this instead of waiting for the next tick.
func (s *Service) PublishPendingEvents() {
	s.publishPendingEvents(time.Now())
}

func (s *Service) publishPendingEvents(now time.Time) {
	if s.repo == nil || s.federateScheduledEvent == nil ||
		(s.getScheduleEnabled != nil && !s.getScheduleEnabled()) ||
		(s.getFederationEnabled != nil && !s.getFederationEnabled()) {
		return
	}

	events, err := s.repo.GetEventsToFederate(now)
	if err != nil {
		log.Errorf("unable to fetch scheduled stream events for federation: %v", err)
		return
	}
	for _, event := range events {
		if err := s.federateScheduledEvent(event); err != nil {
			log.Errorf("unable to federate scheduled stream event %s: %v", event.ID, err)
			continue
		}
		if err := s.repo.SetEventFederatedAt(event.ID, now); err != nil {
			log.Errorf("unable to mark scheduled stream event %s as federated: %v", event.ID, err)
		}
	}
}

func (s *Service) sendScheduledEventReminders(now time.Time) {
	if s.repo == nil || s.notifyScheduledEvent == nil || (s.getScheduleEnabled != nil && !s.getScheduleEnabled()) {
		return
	}

	defaultMessage := ""
	if s.getScheduleReminderMessage != nil {
		defaultMessage = s.getScheduleReminderMessage()
	}
	serverURL := ""
	if s.getServerURL != nil {
		serverURL = s.getServerURL()
	}
	firstReminderMinutes := 0
	if s.getScheduleFirstReminderMinutes != nil {
		firstReminderMinutes = s.getScheduleFirstReminderMinutes()
	}
	secondReminderMinutes := 0
	if s.getScheduleSecondReminderMinutes != nil {
		secondReminderMinutes = s.getScheduleSecondReminderMinutes()
	}
	slots := []struct {
		number  int
		minutes int
	}{
		{number: scheduleeventsrepository.ReminderFirst, minutes: firstReminderMinutes},
		{number: scheduleeventsrepository.ReminderSecond, minutes: secondReminderMinutes},
	}
	seenOffsets := make(map[int]struct{}, len(slots))
	for _, slot := range slots {
		if slot.minutes <= 0 {
			continue
		}
		if _, seen := seenOffsets[slot.minutes]; seen {
			continue
		}
		seenOffsets[slot.minutes] = struct{}{}
		events, err := s.repo.GetEventsNeedingReminder(now, now.Add(time.Duration(slot.minutes)*time.Minute), slot.number)
		if err != nil {
			log.Errorf("unable to fetch scheduled stream reminders: %v", err)
			continue
		}
		for _, event := range events {
			if event.StartTime.Add(-time.Duration(slot.minutes) * time.Minute).After(now) {
				continue
			}
			s.sendScheduledEventReminder(event, slot.number, now, defaultMessage, serverURL)
		}
	}
}

func (s *Service) sendScheduledEventReminder(event models.ScheduledEvent, reminderNumber int, sentAt time.Time, defaultMessage, serverURL string) {
	message := event.ReminderMessage
	if message == "" {
		message = defaultMessage
	}
	if message == "" {
		return
	}
	message = formatScheduledEventReminder(message, event, serverURL)
	s.notifyScheduledEvent(message)
	if err := s.repo.SetEventReminderSentAt(event.ID, reminderNumber, sentAt); err != nil {
		log.Errorf("unable to mark scheduled stream reminder %s as sent: %v", event.ID, err)
	}
}

func formatScheduledEventReminder(message string, event models.ScheduledEvent, serverURL string) string {
	location := time.UTC
	if event.Timezone != "" {
		if loaded, err := time.LoadLocation(event.Timezone); err == nil {
			location = loaded
		}
	}
	start := event.StartTime.In(location).Format("Monday, January 2 at 3:04 PM MST")
	if serverURL == "" {
		return fmt.Sprintf("%s\n\n%s\n%s", message, event.Name, start)
	}
	return fmt.Sprintf("%s\n\n%s\n%s\n%s", message, event.Name, start, serverURL)
}

func (s *Service) sendScheduledEventWebhooks(now time.Time) {
	if s.webhooks == nil || s.repo == nil || (s.getScheduleEnabled != nil && !s.getScheduleEnabled()) {
		return
	}

	warnings, err := s.repo.GetEventsNeedingWebhookWarning(now, now.Add(ScheduledEventWarningLeadTime))
	if err != nil {
		log.Errorf("unable to fetch scheduled stream event warnings: %v", err)
	} else {
		for _, event := range warnings {
			s.webhooks.SendScheduledEvent(event, models.ScheduledEventWarning)
			if err := s.repo.SetEventWebhookWarningSentAt(event.ID, now); err != nil {
				log.Errorf("unable to mark scheduled stream event warning %s as sent: %v", event.ID, err)
			}
		}
	}

	s.sendScheduledEventStartedWebhooks(now)

	ended, err := s.repo.GetEventsNeedingWebhookEnd(now)
	if err != nil {
		log.Errorf("unable to fetch ended scheduled stream events: %v", err)
	} else {
		for _, event := range ended {
			s.webhooks.SendScheduledEvent(event, models.ScheduledEventEnded)
			if err := s.repo.SetEventWebhookEndedSentAt(event.ID, now); err != nil {
				log.Errorf("unable to mark scheduled stream event end %s as sent: %v", event.ID, err)
			}
		}
	}
}

func (s *Service) sendScheduledEventStartedWebhooks(now time.Time) {
	started, err := s.repo.GetEventsNeedingWebhookStart(now)
	if err != nil {
		log.Errorf("unable to fetch started scheduled stream events: %v", err)
		return
	}
	for _, event := range started {
		if s.setStreamTitle != nil {
			if err := s.setStreamTitle(event.Name); err != nil {
				log.Errorf("unable to set scheduled stream title %q: %v", event.Name, err)
				continue
			}
		}
		s.webhooks.SendScheduledEvent(event, models.ScheduledEventStarted)
		if err := s.repo.SetEventWebhookStartedSentAt(event.ID, now); err != nil {
			log.Errorf("unable to mark scheduled stream event start %s as sent: %v", event.ID, err)
		}
	}
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
