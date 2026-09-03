package configrepository

import (
	"testing"
)

func TestScheduleReminderMinutesDefaults(t *testing.T) {
	repo := newAutoplayTestRepo(t)
	if got := repo.GetScheduleFirstReminderMinutes(); got != DefaultScheduleFirstReminderMinutes {
		t.Errorf("first reminder default = %d, want %d", got, DefaultScheduleFirstReminderMinutes)
	}
	if got := repo.GetScheduleSecondReminderMinutes(); got != DefaultScheduleSecondReminderMinutes {
		t.Errorf("second reminder default = %d, want %d", got, DefaultScheduleSecondReminderMinutes)
	}
}

func TestScheduleReminderMinutesRoundTrip(t *testing.T) {
	repo := newAutoplayTestRepo(t)
	if err := repo.SetScheduleFirstReminderMinutes(60); err != nil {
		t.Fatalf("SetScheduleFirstReminderMinutes: %v", err)
	}
	if err := repo.SetScheduleSecondReminderMinutes(15); err != nil {
		t.Fatalf("SetScheduleSecondReminderMinutes: %v", err)
	}
	if got := repo.GetScheduleFirstReminderMinutes(); got != 60 {
		t.Errorf("first reminder = %d, want 60", got)
	}
	if got := repo.GetScheduleSecondReminderMinutes(); got != 15 {
		t.Errorf("second reminder = %d, want 15", got)
	}
}
