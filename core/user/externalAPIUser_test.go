package user

import (
	"testing"

	"github.com/owncast/owncast/core/data"
)

const (
	tokenName = "test token name"
	token     = "test-token-123"
)

var testScopes = []string{"test-scope"}

func TestMain(m *testing.M) {
	if err := data.SetupPersistence(":memory:"); err != nil {
		panic(err)
	}

	SetupUsers()

	m.Run()
}

func TestCreateExternalAPIUser(t *testing.T) {
	if err := InsertExternalAPIUser(token, tokenName, 0, testScopes); err != nil {
		t.Fatal(err)
	}

	user := GetUserByToken(token)
	if user == nil {
		t.Fatal("api user not found after creating")
	}

	if user.DisplayName != tokenName {
		t.Errorf("expected display name %q, got %q", tokenName, user.DisplayName)
	}

	if user.Scopes[0] != testScopes[0] {
		t.Errorf("expected scopes %q, got %q", testScopes, user.Scopes)
	}
}

func TestDeleteExternalAPIUser(t *testing.T) {
	if err := DeleteExternalAPIUser(token); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyTokenDisabled(t *testing.T) {
	users, err := GetExternalAPIUser()
	if err != nil {
		t.Fatal(err)
	}

	if len(users) > 0 {
		t.Fatal("disabled user returned in list of all API users")
	}
}

func TestVerifyGetUserTokenDisabled(t *testing.T) {
	user := GetUserByToken(token)
	if user == nil {
		t.Fatal("user not returned in GetUserByToken after disabling")
	}

	if user.DisabledAt == nil {
		t.Fatal("user returned in GetUserByToken after disabling")
	}
}

func TestVerifyGetExternalAPIUserForAccessTokenAndScopeTokenDisabled(t *testing.T) {
	user, _ := GetExternalAPIUserForAccessTokenAndScope(token, testScopes[0])

	if user != nil {
		t.Fatal("user returned in GetExternalAPIUserForAccessTokenAndScope after disabling")
	}
}

func TestCreateAdditionalAPIUser(t *testing.T) {
	if err := InsertExternalAPIUser("ignore-me", "token-to-be-ignored", 0, testScopes); err != nil {
		t.Fatal(err)
	}
}

func TestAgainVerifyGetExternalAPIUserForAccessTokenAndScopeTokenDisabled(t *testing.T) {
	user, _ := GetExternalAPIUserForAccessTokenAndScope(token, testScopes[0])

	if user != nil {
		t.Fatal("user returned in TestAgainVerifyGetExternalAPIUserForAccessTokenAndScopeTokenDisabled after disabling")
	}
}
