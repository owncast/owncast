package models

import "fmt"

// AutoplayMode controls whether viewers must press play before playback starts.
type AutoplayMode string

const (
	AutoplayOff       AutoplayMode = "off"
	AutoplayAlways    AutoplayMode = "always"
	AutoplaySoundOnly AutoplayMode = "sound-only"
)

// Valid reports whether the mode is supported by Owncast.
func (m AutoplayMode) Valid() bool {
	switch m {
	case AutoplayOff, AutoplayAlways, AutoplaySoundOnly:
		return true
	default:
		return false
	}
}

// Validate returns an error when the mode is unsupported.
func (m AutoplayMode) Validate() error {
	if !m.Valid() {
		return fmt.Errorf("invalid autoplay value %q", m)
	}
	return nil
}
