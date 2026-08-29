package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDisabledReasonIsAdminOnly(t *testing.T) {
	user := &User{ID: "viewer", DisabledReason: "Repeated spam"}

	publicPayload, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(publicPayload), "disabledReason") || strings.Contains(string(publicPayload), user.DisabledReason) {
		t.Fatalf("public user payload exposed disabled reason: %s", publicPayload)
	}

	adminPayload, err := json.Marshal(UserWithDisabledReasonFrom(user))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(adminPayload), `"disabledReason":"Repeated spam"`) {
		t.Fatalf("admin user payload omitted disabled reason: %s", adminPayload)
	}
}
