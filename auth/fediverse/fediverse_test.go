package fediverse

import (
	"strings"
	"testing"
)

const (
	accessToken     = "fake-access-token"
	account         = "blah"
	userID          = "fake-user-id"
	userDisplayName = "fake-user-display-name"
)

func TestOTPFlowValidation(t *testing.T) {
	r, success := RegisterFediverseOTP(accessToken, userID, userDisplayName, account)

	if !success {
		t.Error("Registration should be permitted.")
	}

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

func TestSingleOTPFlowRequest(t *testing.T) {
	r1, _ := RegisterFediverseOTP(accessToken, userID, userDisplayName, account)
	r2, s2 := RegisterFediverseOTP(accessToken, userID, userDisplayName, account)

	if r1.Code != r2.Code {
		t.Error("Only one registration should be permitted.")
	}

	if s2 {
		t.Error("Second registration should not be permitted.")
	}
}

func TestAccountCaseInsensitive(t *testing.T) {
	account := "Account"
	accessToken := "another-fake-access-token"
	r1, _ := RegisterFediverseOTP(accessToken, userID, userDisplayName, account)
	_, reg1 := ValidateFediverseOTP(accessToken, r1.Code)

	// Simulate second auth with account in different case
	r2, _ := RegisterFediverseOTP(accessToken, userID, userDisplayName, strings.ToUpper(account))
	_, reg2 := ValidateFediverseOTP(accessToken, r2.Code)

	if reg1.Account != reg2.Account {
		t.Errorf("Account names should be case-insensitive: %s %s", reg1.Account, reg2.Account)
	}
}
