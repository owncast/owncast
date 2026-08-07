package models

const (
	// DefaultMaxClipDurationSeconds is the longest clip a viewer may create
	// when the operator has not configured a limit.
	DefaultMaxClipDurationSeconds = 120

	// MaxAllowedClipDurationSeconds caps what an operator may configure as
	// the maximum clip length.
	MaxAllowedClipDurationSeconds = 3600

	// DefaultClipDurationSeconds is the trailing window captured when a
	// viewer clips a live stream without specifying a duration.
	DefaultClipDurationSeconds = 30
)

// ClipCreatedResponse is returned when a clip is successfully created so the
// client can link to or play the new clip immediately.
type ClipCreatedResponse struct {
	Success         bool   `json:"success"`
	Message         string `json:"message"`
	ID              string `json:"id"`
	DurationSeconds int    `json:"durationSeconds"`
}
