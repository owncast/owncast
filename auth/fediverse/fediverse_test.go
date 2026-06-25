package fediverse

import (
	"strings"
	"testing"

	"github.com/owncast/owncast/utils"
)

const (
	accessToken     = "fake-access-token"
	account         = "blah"
	userID          = "fake-user-id"
	userDisplayName = "fake-user-display-name"
)

func TestOTPFlowValidation(t *testing.T) {
	svc := New()
	r, success, err := svc.RegisterFediverseOTP(accessToken, userID, userDisplayName, account)
	if err != nil {
		t.Error(err)
	}

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

	valid, registration := svc.ValidateFediverseOTP(accessToken, r.Code)
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
	svc := New()
	r1, _, _ := svc.RegisterFediverseOTP(accessToken, userID, userDisplayName, account)
	r2, s2, _ := svc.RegisterFediverseOTP(accessToken, userID, userDisplayName, account)

	if r1.Code != r2.Code {
		t.Error("Only one registration should be permitted.")
	}

	if s2 {
		t.Error("Second registration should not be permitted.")
	}
}

func TestReplacePendingWhenAccountDiffers(t *testing.T) {
	svc := New()
	token := "retry-token"

	first, _, _ := svc.RegisterFediverseOTP(token, userID, userDisplayName, "wrong@example.com")

	// A different account replaces the pending request with a fresh code.
	second, ok, _ := svc.RegisterFediverseOTP(token, userID, userDisplayName, "right@example.com")
	if !ok {
		t.Fatal("a request for a different account should replace the pending one")
	}
	if second.Code == first.Code {
		t.Error("replacement should issue a new code")
	}

	// The same account is still blocked from re-sending.
	if _, ok, _ := svc.RegisterFediverseOTP(token, userID, userDisplayName, "right@example.com"); ok {
		t.Error("a duplicate request for the same account should be blocked")
	}

	// Only the replacement code is valid now.
	if valid, _ := svc.ValidateFediverseOTP(token, first.Code); valid {
		t.Error("the replaced code should no longer validate")
	}
	if valid, _ := svc.ValidateFediverseOTP(token, second.Code); !valid {
		t.Error("the replacement code should validate")
	}
}

func TestAccountCaseInsensitive(t *testing.T) {
	svc := New()
	account := "Account"
	accessToken := "another-fake-access-token"
	r1, _, _ := svc.RegisterFediverseOTP(accessToken, userID, userDisplayName, account)
	_, reg1 := svc.ValidateFediverseOTP(accessToken, r1.Code)

	// Simulate second auth with account in different case
	r2, _, _ := svc.RegisterFediverseOTP(accessToken, userID, userDisplayName, strings.ToUpper(account))
	_, reg2 := svc.ValidateFediverseOTP(accessToken, r2.Code)

	if reg1.Account != reg2.Account {
		t.Errorf("Account names should be case-insensitive: %s %s", reg1.Account, reg2.Account)
	}
}

func TestLimitGlobalPendingRequests(t *testing.T) {
	svc := New()
	for i := 0; i < maxPendingRequests; i++ {
		at, _ := utils.GenerateRandomString(10)
		uid, _ := utils.GenerateRandomString(10)
		account, _ := utils.GenerateRandomString(10)

		_, success, error := svc.RegisterFediverseOTP(at, uid, "userDisplayName", account)
		if !success {
			t.Error("Registration should be permitted.", i, " of ", len(svc.pendingAuthRequests))
		}
		if error != nil {
			t.Error(error)
		}
	}

	// This one should fail
	at, _ := utils.GenerateRandomString(10)
	uid, _ := utils.GenerateRandomString(10)
	account, _ := utils.GenerateRandomString(10)
	_, success, error := svc.RegisterFediverseOTP(at, uid, "userDisplayName", account)
	if success {
		t.Error("Registration should not be permitted.")
	}
	if error == nil {
		t.Error("Error should be returned.")
	}
}
