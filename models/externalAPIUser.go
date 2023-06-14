package models

import (
	"time"
)

// ExternalAPIUser represents a single 3rd party integration that uses an access token.
// This struct mostly matches the User struct so they can be used interchangeably.
type ExternalAPIUser struct {
	CreatedAt    time.Time  `json:"createdAt"`
	LastUsedAt   *time.Time `json:"lastUsedAt,omitempty"`
	ID           string     `json:"id"`
	AccessToken  string     `json:"accessToken"`
	DisplayName  string     `json:"displayName"`
	Type         string     `json:"type,omitempty"` // Should be API
	Scopes       []string   `json:"scopes"`
	DisplayColor int        `json:"displayColor"`
	IsBot        bool       `json:"isBot"`
}
