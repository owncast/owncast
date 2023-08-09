package replays

import (
	"bytes"
	"fmt"

	"github.com/grafov/m3u8"
)

// MediaPlaylistAllowCacheTag is a custom tag to explicitly state that this
// playlist is allowed to be cached.
type MediaPlaylistAllowCacheTag struct {
	Type string
}

// TagName should return the full tag identifier including the leading
// '#' and trailing ':' if the tag also contains a value or attribute
// list.
func (tag *MediaPlaylistAllowCacheTag) TagName() string {
	return "#EXT-X-ALLOW-CACHE"
}

// Decode decodes the input line. The line will be the entire matched
// line, including the identifier.
func (tag *MediaPlaylistAllowCacheTag) Decode(line string) (m3u8.CustomTag, error) {
	_, err := fmt.Sscanf(line, "#EXT-X-ALLOW-CACHE")

	return tag, err
}

// SegmentTag specifies that this tag is not for segments.
func (tag *MediaPlaylistAllowCacheTag) SegmentTag() bool {
	return false
}

// Encode formats the structure to the text result.
func (tag *MediaPlaylistAllowCacheTag) Encode() *bytes.Buffer {
	buf := new(bytes.Buffer)

	buf.WriteString(tag.TagName())
	buf.WriteString(tag.Type)

	return buf
}

// String implements Stringer interface.
func (tag *MediaPlaylistAllowCacheTag) String() string {
	return tag.Encode().String()
}
