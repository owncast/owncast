package models

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// ScopeCanSendUserMessages will allow sending chat messages as users.
	ScopeCanSendUserMessages = "CAN_SEND_MESSAGES"
	// ScopeCanSendSystemMessages will allow sending chat messages as the system.
	ScopeCanSendSystemMessages = "CAN_SEND_SYSTEM_MESSAGES"
	// ScopeHasAdminAccess will allow performing administrative actions on the server.
	ScopeHasAdminAccess = "HAS_ADMIN_ACCESS"
)

// For a scope to be seen as "valid" it must live in this slice.
var validAccessTokenScopes = []string{
	ScopeCanSendUserMessages,
	ScopeCanSendSystemMessages,
	ScopeHasAdminAccess,
}

// AccessToken gives access to 3rd party code to access specific Owncast APIs.
type AccessToken struct {
	Token     string     `json:"token"`
	Name      string     `json:"name"`
	Scopes    []string   `json:"scopes"`
	Timestamp time.Time  `json:"timestamp"`
	LastUsed  *time.Time `json:"lastUsed"`
}

// HasValidScopes will verify that all the scopes provided are valid.
// This is not a efficient method.
func HasValidScopes(scopes []string) bool {
	for _, scope := range scopes {
		if !findItemInSlice(validAccessTokenScopes, scope) {
			log.Errorln("Invalid scope", scope)
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
