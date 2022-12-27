package models

// StreamKey represents a single stream key.
type StreamKeyHashed struct {
	Key     []byte `json:"key"`
	Comment string `json:"comment"`
}
