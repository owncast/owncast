package middleware

import (
	"net/http/httptest"
	"strings"
	"testing"

	logtest "github.com/sirupsen/logrus/hooks/test"
)

func TestInvalidAccessTokenLogIncludesRequestContext(t *testing.T) {
	hook := logtest.NewGlobal()
	defer hook.Reset()

	req := httptest.NewRequest("POST", "http://owncast.test/api/integrations/chat?ignored=true", nil)
	req.RemoteAddr = "203.0.113.10:1234"
	req.Header.Set("Authorization", "Bearer secret-token")

	logInvalidAccessToken(req, "CAN_SEND_MESSAGES")

	entry := hook.LastEntry()
	if entry == nil {
		t.Fatal("expected an invalid access token log entry")
	}
	for _, value := range []string{"203.0.113.10", "POST", "/api/integrations/chat", "CAN_SEND_MESSAGES"} {
		if !strings.Contains(entry.Message, value) {
			t.Errorf("log entry %q does not include %q", entry.Message, value)
		}
	}
	if strings.Contains(entry.Message, "secret-token") || strings.Contains(entry.Message, "ignored=true") {
		t.Errorf("log entry includes sensitive request data: %q", entry.Message)
	}
}
