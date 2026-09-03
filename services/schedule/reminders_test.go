package schedule

import (
	"strings"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/scheduleeventsrepository"
)

type reminderScheduleRepo struct {
	scheduleeventsrepository.ScheduleEventsRepository
	events map[int][]models.ScheduledEvent
	sentAt map[int]map[string]time.Time
}

func (r *reminderScheduleRepo) GetEventsNeedingReminder(_ time.Time, _ time.Time, reminderNumber int) ([]models.ScheduledEvent, error) {
	var pending []models.ScheduledEvent
	for _, event := range r.events[reminderNumber] {
		if _, sent := r.sentAt[reminderNumber][event.ID]; !sent {
			pending = append(pending, event)
		}
	}
	return pending, nil
}

func (r *reminderScheduleRepo) SetEventReminderSentAt(id string, reminderNumber int, sentAt time.Time) error {
	r.sentAt[reminderNumber][id] = sentAt
	return nil
}

func TestScheduledEventRemindersUseOverridesAndFireOnce(t *testing.T) {
	now := time.Date(2030, time.September, 8, 18, 0, 0, 0, time.UTC)
	repo := &reminderScheduleRepo{
		events: map[int][]models.ScheduledEvent{
			scheduleeventsrepository.ReminderFirst:  {{ID: "custom", Name: "Custom Event", ReminderMessage: "Custom reminder", StartTime: now.Add(2 * time.Hour), Timezone: "UTC"}},
			scheduleeventsrepository.ReminderSecond: {{ID: "default", Name: "Default Event", StartTime: now.Add(15 * time.Minute), Timezone: "UTC"}},
		},
		sentAt: map[int]map[string]time.Time{
			scheduleeventsrepository.ReminderFirst:  {},
			scheduleeventsrepository.ReminderSecond: {},
		},
	}
	var messages []string
	service := New(Deps{
		ScheduleEventsRepository:         repo,
		GetScheduleEnabled:               func() bool { return true },
		GetScheduleReminderMessage:       func() string { return "Default reminder" },
		GetScheduleFirstReminderMinutes:  func() int { return 120 },
		GetScheduleSecondReminderMinutes: func() int { return 15 },
		GetServerURL:                     func() string { return "https://owncast.example" },
		NotifyScheduledEvent: func(message string) {
			messages = append(messages, message)
		},
	})

	service.sendScheduledEventReminders(now)
	service.sendScheduledEventReminders(now.Add(time.Minute))

	if len(messages) != 2 {
		t.Fatalf("sent %d reminders, want 2", len(messages))
	}
	if !strings.Contains(messages[0], "Custom reminder") || !strings.Contains(messages[0], "Custom Event") {
		t.Errorf("custom reminder message = %q", messages[0])
	}
	if !strings.Contains(messages[1], "Default reminder") || !strings.Contains(messages[1], "Default Event") {
		t.Errorf("default reminder message = %q", messages[1])
	}
	if _, ok := repo.sentAt[scheduleeventsrepository.ReminderFirst]["custom"]; !ok {
		t.Error("custom reminder was not marked sent")
	}
	if _, ok := repo.sentAt[scheduleeventsrepository.ReminderSecond]["default"]; !ok {
		t.Error("default reminder was not marked sent")
	}
}

func TestScheduledEventRemindersSkipDisabledSlots(t *testing.T) {
	repo := &reminderScheduleRepo{
		events: map[int][]models.ScheduledEvent{scheduleeventsrepository.ReminderFirst: {{ID: "event"}}},
		sentAt: map[int]map[string]time.Time{scheduleeventsrepository.ReminderFirst: {}},
	}
	called := false
	service := New(Deps{
		ScheduleEventsRepository:         repo,
		GetScheduleFirstReminderMinutes:  func() int { return 0 },
		GetScheduleSecondReminderMinutes: func() int { return 0 },
		NotifyScheduledEvent: func(string) {
			called = true
		},
	})

	service.sendScheduledEventReminders(time.Now())
	if called {
		t.Error("disabled reminder slot sent a message")
	}
}

func TestScheduledEventRemindersDoNotSendEarly(t *testing.T) {
	now := time.Date(2030, time.September, 8, 18, 0, 0, 0, time.UTC)
	repo := &reminderScheduleRepo{
		events: map[int][]models.ScheduledEvent{
			scheduleeventsrepository.ReminderFirst: {{ID: "future", Name: "Future Event", StartTime: now.Add(3 * time.Hour), Timezone: "UTC"}},
		},
		sentAt: map[int]map[string]time.Time{
			scheduleeventsrepository.ReminderFirst: {},
		},
	}
	called := false
	service := New(Deps{
		ScheduleEventsRepository:        repo,
		GetScheduleFirstReminderMinutes: func() int { return 120 },
		NotifyScheduledEvent: func(string) {
			called = true
		},
	})

	service.sendScheduledEventReminders(now)
	if called {
		t.Error("reminder was sent before its due time")
	}
}
