package models

import "time"

const (
	ScopeCanSendUserMessages   = "CAN_SEND_MESSAGES"
	ScopeCanSendSystemMessages = "CAN_SEND_SYSTEM_MESSAGES"
)

// For a scope to be seen as "valid" it must live in this slice.
var validAccessTokenScopes = []string{
	ScopeCanSendUserMessages,
	ScopeCanSendSystemMessages,
}

type AccessToken struct {
	Token     string    `json:"token"`
	Name      string    `json:"name"`
	Scopes    []string  `json:"scopes"`
	Timestamp time.Time `json:"timestamp"`
	LastUsed  time.Time `json:"lastUsed"`
}

// HasValidScopes will verify that all the scopes provided are valid.
// This is not a efficient method.
func HasValidScopes(scopes []string) bool {
	for _, scope := range scopes {
		if !findItemInSlice(validAccessTokenScopes, scope) {
			return false
		}
	}
	return true
}

func findItemInSlice(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
