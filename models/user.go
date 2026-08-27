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

// AdminUser includes moderation-only fields that must not appear in public
// chat or webhook payloads.
type AdminUser struct {
	*User
	DisabledReason string `json:"disabledReason,omitempty"`
}

// AdminUserFrom exposes a user's moderation-only fields to authenticated
// admin and moderator endpoints.
func AdminUserFrom(user *User) *AdminUser {
	if user == nil {
		return nil
	}
	return &AdminUser{User: user, DisabledReason: user.DisabledReason}
}

// AdminUsersFrom converts a collection for authenticated admin endpoints.
func AdminUsersFrom(users []*User) []*AdminUser {
	result := make([]*AdminUser, 0, len(users))
	for _, user := range users {
		result = append(result, AdminUserFrom(user))
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
