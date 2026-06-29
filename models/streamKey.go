package models

// StreamKey is a valid RTMP stream key plus an optional human-readable
// comment. It is the domain type used throughout config storage and the
// streaming engine. The web layer has its own generated wire type with an
// identical JSON shape; convert at the HTTP edge.
type StreamKey struct {
	Comment *string `json:"comment,omitempty"`
	Key     *string `json:"key,omitempty"`
}
