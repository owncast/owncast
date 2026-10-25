package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/owncast/owncast/services/activitypub/inbox"
)

// TestInboxBodySizeCap verifies that the inbox handler enforces a 1 MiB
// request body limit: an oversized POST must return 413, and a normal-sized
// POST must still return 202.
func TestInboxBodySizeCap(t *testing.T) {
	// Enable federation and set a username so the handler doesn't short-circuit.
	if err := testControllers.configRepository.SetFederationEnabled(true); err != nil {
		t.Fatalf("SetFederationEnabled: %v", err)
	}
	defer testControllers.configRepository.SetFederationEnabled(false) //nolint

	const fedUsername = "testuser"
	if err := testControllers.configRepository.SetFederationUsername(fedUsername); err != nil {
		t.Fatalf("SetFederationUsername: %v", err)
	}

	// Provide a started inbox service so AddToQueue doesn't block forever.
	// Workers have a recover() so nil deps won't crash the test process.
	inboxSvc := inbox.New(inbox.Deps{})
	inboxSvc.Start()
	testControllers.inbox = inboxSvc

	path := "/federation/user/" + fedUsername + "/inbox"

	t.Run("oversized body returns 413", func(t *testing.T) {
		body := bytes.Repeat([]byte("x"), (1<<20)+1) // 1 MiB + 1 byte
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/activity+json")
		rr := httptest.NewRecorder()
		testControllers.InboxHandler(rr, req)
		if rr.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("oversized body: got status %d, want 413", rr.Code)
		}
	})

	t.Run("body within limit returns 202", func(t *testing.T) {
		body := strings.Repeat("x", 512) // well under 1 MiB
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/activity+json")
		rr := httptest.NewRecorder()
		testControllers.InboxHandler(rr, req)
		if rr.Code != http.StatusAccepted {
			t.Errorf("small body: got status %d, want 202", rr.Code)
		}
	})
}
