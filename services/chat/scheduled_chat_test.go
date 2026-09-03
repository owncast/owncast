package chat

import (
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
)

func TestShouldRejectOfflineMessage(t *testing.T) {
	now := time.Date(2030, time.September, 8, 18, 0, 0, 0, time.UTC)
	status := models.Status{LastDisconnectTime: &utils.NullTime{Time: now.Add(-6 * time.Minute)}}

	if shouldRejectOfflineMessage(status, true, now) {
		t.Fatal("scheduled chat message was rejected")
	}
	if !shouldRejectOfflineMessage(status, false, now) {
		t.Fatal("offline message was accepted after the disconnect grace period")
	}
}
