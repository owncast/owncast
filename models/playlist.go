package models

import "path"

//Segment represents a segment of the live stream
type Segment struct {
	VariantIndex       int    // The bitrate variant
	FullDiskPath       string // Where it lives on disk
	RelativeUploadPath string // Path it should have remotely
	RemoteID           string // Used for IPFS
}

//Variant represents a single bitrate variant and the segments that make it up
type Variant struct {
	VariantIndex int
	Segments     []Segment
}

//GetSegmentForFilename gets the segment for the provided filename
func (v *Variant) GetSegmentForFilename(filename string) *Segment {
	for _, segment := range v.Segments {
		if path.Base(segment.FullDiskPath) == filename {
			return &segment
		}
	}

	return nil
}
