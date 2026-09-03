package schedule

import (
	"strings"
	"testing"
	"time"
)

func TestScheduledEventRemindersUseRepositoryAndPersistEachSlot(t *testing.T) {
	nowFirst := time.Date(2045, time.January, 1, 12, 0, 0, 0, time.UTC)
	nowSecond := time.Date(2045, time.January, 2, 12, 0, 0, 0, time.UTC)
	firstID, err := testRepo.AddOneOffEvent("first integration reminder", "", "First", nowFirst.Add(2*time.Hour), 60, "UTC")
	if err != nil {
		t.Fatalf("AddOneOffEvent(first): %v", err)
	}
	secondID, err := testRepo.AddOneOffEvent("second integration reminder", "", "Second", nowSecond.Add(15*time.Minute), 60, "UTC")
	if err != nil {
		t.Fatalf("AddOneOffEvent(second): %v", err)
	}
	t.Cleanup(func() {
		if err := testRepo.DeleteEvent(firstID); err != nil {
			t.Errorf("DeleteEvent(first): %v", err)
		}
		if err := testRepo.DeleteEvent(secondID); err != nil {
			t.Errorf("DeleteEvent(second): %v", err)
		}
	})

	var messages []string
	firstService := New(Deps{
		ScheduleEventsRepository:        testRepo,
		GetScheduleFirstReminderMinutes: func() int { return 120 },
		NotifyScheduledEvent: func(message string) {
			messages = append(messages, message)
		},
	})
	firstService.sendScheduledEventReminders(nowFirst)

	secondService := New(Deps{
		ScheduleEventsRepository:         testRepo,
		GetScheduleSecondReminderMinutes: func() int { return 15 },
		NotifyScheduledEvent: func(message string) {
			messages = append(messages, message)
		},
	})
	secondService.sendScheduledEventReminders(nowSecond)

	if len(messages) != 2 {
		t.Fatalf("sent %d reminders, want 2", len(messages))
	}
	if !strings.Contains(messages[0], "First") || !strings.Contains(messages[0], "first integration reminder") {
		t.Errorf("first reminder = %q", messages[0])
	}
	if !strings.Contains(messages[1], "Second") || !strings.Contains(messages[1], "second integration reminder") {
		t.Errorf("second reminder = %q", messages[1])
	}

	first, err := testRepo.GetEvent(firstID)
	if err != nil {
		t.Fatalf("GetEvent(first): %v", err)
	}
	if first.Reminder1SentAt == nil || first.Reminder2SentAt != nil {
		t.Errorf("first markers = (%v, %v), want only first marker", first.Reminder1SentAt, first.Reminder2SentAt)
	}
	second, err := testRepo.GetEvent(secondID)
	if err != nil {
		t.Fatalf("GetEvent(second): %v", err)
	}
	if second.Reminder1SentAt != nil || second.Reminder2SentAt == nil {
		t.Errorf("second markers = (%v, %v), want only second marker", second.Reminder1SentAt, second.Reminder2SentAt)
	}

	firstService.sendScheduledEventReminders(nowFirst)
	secondService.sendScheduledEventReminders(nowSecond)
	if len(messages) != 2 {
		t.Fatalf("sent %d reminders after repeat, want 2", len(messages))
	}
}
