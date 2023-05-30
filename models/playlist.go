package models

// Segment represents a segment of the live stream.
type Segment struct {
	FullDiskPath       string // Where it lives on disk
	RelativeUploadPath string // Path it should have remotely
	RemoteURL          string
	VariantIndex       int // The bitrate variant
}

// Variant represents a single video variant and the segments that make it up.
type Variant struct {
	Segments     map[string]*Segment
	VariantIndex int
}

// GetSegmentForFilename gets the segment for the provided filename.
func (v *Variant) GetSegmentForFilename(filename string) *Segment {
	return v.Segments[filename]
}
