package fediverse

import "testing"

const (
	accessToken     = "fake-access-token"
	account         = "blah"
	userID          = "fake-user-id"
	userDisplayName = "fake-user-display-name"
)

func TestOTPFlowValidation(t *testing.T) {
	r := RegisterFediverseOTP(accessToken, userID, userDisplayName, account)

	if r.Code == "" {
		t.Error("Code is empty")
	}

	if r.Account != account {
		t.Error("Account is not set correctly")
	}

	if r.Timestamp.IsZero() {
		t.Error("Timestamp is empty")
	}

	valid, registration := ValidateFediverseOTP(accessToken, r.Code)
	if !valid {
		t.Error("Code is not valid")
	}

	if registration.Account != account {
		t.Error("Account is not set correctly")
	}

	if registration.UserID != userID {
		t.Error("UserID is not set correctly")
	}

	if registration.UserDisplayName != userDisplayName {
		t.Error("UserDisplayName is not set correctly")
	}
}
