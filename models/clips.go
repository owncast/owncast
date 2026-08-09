package models

import "time"

const (
	// DefaultMaxClipDurationSeconds is the longest clip a viewer may create
	// when the operator has not configured a limit.
	DefaultMaxClipDurationSeconds = 120

	// MaxAllowedClipDurationSeconds caps what an operator may configure as
	// the maximum clip length.
	MaxAllowedClipDurationSeconds = 3600
)

// Clip permission levels: who may create clips. Every level requires a
// registered chat identity; moderators always qualify.
const (
	// ClipPermissionsModerators restricts clip creation to moderators.
	ClipPermissionsModerators = "moderators"
	// ClipPermissionsAuthenticated restricts clip creation to viewers who
	// authenticated via an auth provider (IndieAuth/FediAuth).
	ClipPermissionsAuthenticated = "authenticated"
	// ClipPermissionsEstablished allows any chat identity that has existed
	// for at least MinClipperAccountAge, keeping drive-by fresh accounts out.
	ClipPermissionsEstablished = "established"

	// DefaultClipPermissions applies when the operator has not chosen.
	DefaultClipPermissions = ClipPermissionsEstablished
)

// MinClipperAccountAge is how old a chat identity must be before it may
// create clips under the "established" permission level.
const MinClipperAccountAge = time.Hour

// IsValidClipPermissions reports whether the value names a known clip
// permission level.
func IsValidClipPermissions(value string) bool {
	switch value {
	case ClipPermissionsModerators, ClipPermissionsAuthenticated, ClipPermissionsEstablished:
		return true
	}
	return false
}

// ClipCreatedResponse is returned when a clip is successfully created so the
// client can link to or play the new clip immediately.
type ClipCreatedResponse struct {
	Success         bool   `json:"success"`
	Message         string `json:"message"`
	ID              string `json:"id"`
	DurationSeconds int    `json:"durationSeconds"`
}
