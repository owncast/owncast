package admin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/owncast/owncast/persistence/configrepository"
)

func makeNumericConfigValueBody(value int) string {
	return fmt.Sprintf(`{"value": %d}`, value)
}

func TestSetScheduleReminderMinutes(t *testing.T) {
	repo := configrepository.New(testDatastore)
	for _, value := range []int{0, 15, 30, 60, 120, 1440} {
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/admin/config/schedule/firstreminderminutes",
			strings.NewReader(makeNumericConfigValueBody(value)),
		)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		testAdmin.SetScheduleFirstReminderMinutes(w, req)

		resp := parseResponse(t, w)
		if !resp.Success {
			t.Fatalf("value %d: expected success, got error: %s", value, resp.Message)
		}
		if got := repo.GetScheduleFirstReminderMinutes(); got != value {
			t.Errorf("value %d: persisted first reminder = %d", value, got)
		}
	}
}

func TestSetScheduleReminderMinutesRejectsInvalidValue(t *testing.T) {
	repo := configrepository.New(testDatastore)
	if err := repo.SetScheduleSecondReminderMinutes(15); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/admin/config/schedule/secondreminderminutes",
		strings.NewReader(makeNumericConfigValueBody(10)),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testAdmin.SetScheduleSecondReminderMinutes(w, req)

	resp := parseResponse(t, w)
	if resp.Success {
		t.Fatal("expected failure for an invalid reminder value")
	}
	if got := repo.GetScheduleSecondReminderMinutes(); got != 15 {
		t.Errorf("invalid set changed second reminder to %d", got)
	}
}
