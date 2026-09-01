package models

import (
	"time"

	"github.com/owncast/owncast/utils"
)

const (
	moderatorScopeKey = "MODERATOR"
)

type User struct {
	CreatedAt       time.Time  `json:"createdAt"`
	DisabledAt      *time.Time `json:"disabledAt,omitempty"`
	DisabledReason  string     `json:"-"`
	NameChangedAt   *time.Time `json:"nameChangedAt,omitempty"`
	AuthenticatedAt *time.Time `json:"-"`
	ID              string     `json:"id"`
	DisplayName     string     `json:"displayName"`
	PreviousNames   []string   `json:"previousNames"`
	Scopes          []string   `json:"scopes,omitempty"`
	AuthProviders   []string   `json:"authProviders,omitempty"`
	DisplayColor    int        `json:"displayColor"`
	IsBot           bool       `json:"isBot"`
	Authenticated   bool       `json:"authenticated"`
}

// UserWithDisabledReason includes the moderation-only reason that must not
// appear in public chat or webhook payloads.
type UserWithDisabledReason struct {
	*User
	DisabledReason string `json:"disabledReason,omitempty"`
}

// UserWithDisabledReasonFrom exposes the moderation-only reason to
// authenticated admin and moderator endpoints.
func UserWithDisabledReasonFrom(user *User) *UserWithDisabledReason {
	if user == nil {
		return nil
	}
	return &UserWithDisabledReason{User: user, DisabledReason: user.DisabledReason}
}

// UsersWithDisabledReasonsFrom converts a collection for authenticated admin
// endpoints.
func UsersWithDisabledReasonsFrom(users []*User) []*UserWithDisabledReason {
	result := make([]*UserWithDisabledReason, 0, len(users))
	for _, user := range users {
		result = append(result, UserWithDisabledReasonFrom(user))
	}
	return result
}

// IsEnabled will return if this single user is enabled.
func (u *User) IsEnabled() bool {
	return u.DisabledAt == nil
}

// IsModerator will return if the user has moderation privileges.
func (u *User) IsModerator() bool {
	_, hasModerationScope := utils.FindInSlice(u.Scopes, moderatorScopeKey)
	return hasModerationScope
}
