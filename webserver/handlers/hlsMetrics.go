package handlers

import (
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/utils"
)

// Server-side playback observation for players that don't self-report
// metrics (Safari/iOS native HLS, VLC, mpv, ffmpeg, etc.) plus parsing of
// CMCD (CTA-5004 v1 and CTA-5004-A v2) data attached to media requests by
// players that support it (hls.js, dash.js, Shaka, ExoPlayer, ...).
// Event and response mode reports POSTed to the collector endpoint reuse
// the same key registration, see cmcdCollector.go.

// cmcdHeaders are the four request headers CMCD data may arrive in when a
// player uses header transmission instead of the CMCD query parameter.
var cmcdHeaders = []string{"CMCD-Request", "CMCD-Object", "CMCD-Status", "CMCD-Session"}

// parseCMCDRequest extracts CMCD keys from a media request. Players send
// them either as a single urlencoded CMCD query parameter or split across
// the CMCD-* request headers. Returns nil when the request carries no CMCD
// data.
func parseCMCDRequest(r *http.Request) map[string]any {
	keys := map[string]any{}
	if payload := r.URL.Query().Get("CMCD"); payload != "" {
		parseCMCDPayload(payload, keys)
	}
	for _, header := range cmcdHeaders {
		if payload := r.Header.Get(header); payload != "" {
			parseCMCDPayload(payload, keys)
		}
	}
	if len(keys) == 0 {
		return nil
	}
	return keys
}

// parseCMCDPayload parses a CMCD dictionary payload into keys: comma
// separated key or key=value pairs, where string values are double-quoted,
// bare values are numbers or tokens, and a value-less key is boolean true.
// The value typing matches what a JSON-decoded CMCD report produces so both
// transports share the registration path.
func parseCMCDPayload(payload string, keys map[string]any) {
	for _, item := range splitOutsideQuotes(payload, ',') {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		key := item
		value := ""
		if idx := indexOutsideQuotes(item, '='); idx != -1 {
			key = item[:idx]
			value = item[idx+1:]
		}

		switch {
		case key == "":
			continue
		case key == item:
			keys[key] = true
		case strings.HasPrefix(value, `"`):
			keys[key] = unquoteCMCDString(value)
		default:
			if number, err := strconv.ParseFloat(value, 64); err == nil {
				keys[key] = number
			} else {
				keys[key] = value
			}
		}
	}
}

// unquoteCMCDString strips the surrounding double quotes of a CMCD string
// value and resolves \" and \\ escapes.
func unquoteCMCDString(value string) string {
	value = strings.TrimPrefix(value, `"`)
	value = strings.TrimSuffix(value, `"`)
	if !strings.Contains(value, `\`) {
		return value
	}
	var b strings.Builder
	b.Grow(len(value))
	for i := 0; i < len(value); i++ {
		if value[i] == '\\' && i+1 < len(value) {
			i++
		}
		b.WriteByte(value[i])
	}
	return b.String()
}

// splitOutsideQuotes splits s on sep, ignoring separators inside
// double-quoted strings ("\"" escapes a quote per structured-field syntax).
func splitOutsideQuotes(s string, sep byte) []string {
	var parts []string
	start := 0
	inQuotes := false
	for i := 0; i < len(s); i++ {
		switch {
		case s[i] == '\\' && inQuotes:
			i++ // skip escaped character
		case s[i] == '"':
			inQuotes = !inQuotes
		case s[i] == sep && !inQuotes:
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	return append(parts, s[start:])
}

// indexOutsideQuotes returns the index of the first sep outside quoted
// strings, or -1.
func indexOutsideQuotes(s string, sep byte) int {
	inQuotes := false
	for i := 0; i < len(s); i++ {
		switch {
		case s[i] == '\\' && inQuotes:
			i++
		case s[i] == '"':
			inQuotes = !inQuotes
		case s[i] == sep && !inQuotes:
			return i
		}
	}
	return -1
}

// cmcdNumber returns the numeric value of a CMCD key, whether it arrived
// via JSON or the dictionary payload syntax.
func cmcdNumber(keys map[string]any, key string) (float64, bool) {
	number, ok := keys[key].(float64)
	return number, ok
}

// cmcdClientID returns the metrics identity for a CMCD report: the CMCD
// session ID when present — it identifies the player better than
// IP+UserAgent (it distinguishes players behind shared NAT) — otherwise
// the request-derived client ID.
func cmcdClientID(r *http.Request, keys map[string]any) string {
	if id, _ := keys["sid"].(string); id != "" {
		return id
	}
	return utils.GenerateClientIDFromRequest(r)
}

// registerCMCDKeys maps the consumed subset of a CMCD report — from any
// transmission or reporting mode — onto the playback metrics pipeline in a
// single locked update, since this runs per media request.
func (h *Handlers) registerCMCDKeys(id string, keys map[string]any) {
	report := metrics.PlaybackReport{
		ClientID: id,
		// A zero error count is registered for every report so healthy
		// clients stay in the health overview's denominator, matching the
		// legacy metrics endpoint's behavior.
		HasErrorCount: true,
	}

	if mtp, ok := cmcdNumber(keys, "mtp"); ok && mtp > 0 {
		report.BandwidthKbps = mtp
	}
	// v2: live latency from the live edge, in milliseconds.
	if ltc, ok := cmcdNumber(keys, "ltc"); ok && ltc > 0 {
		report.LatencySeconds = ltc / 1000
	}
	// Error signals, both reported by the player so unlike server-side
	// inference they can't be confused with a pause: bs is a buffer
	// starvation and e=e is a fatal error event.
	if starved, ok := keys["bs"].(bool); ok && starved {
		report.ErrorCount++
	}
	if event, _ := keys["e"].(string); event == "e" {
		report.ErrorCount++
	}
	// v2 response mode: time to last byte is the client-measured download
	// duration of the reported object, in milliseconds.
	if ttlb, ok := cmcdNumber(keys, "ttlb"); ok && ttlb > 0 {
		report.DownloadSeconds = ttlb / 1000
	}
	// The encoded bitrate being played, observed for variant change
	// detection.
	if br, ok := cmcdNumber(keys, "br"); ok && br > 0 {
		report.BitrateKbps = br
	}

	h.metrics.RegisterPlaybackReport(report)
}

// Segments smaller than this measure kernel socket buffers, not the
// client's connection, so they are not usable speed samples.
const minMeasurableSegmentBytes = 64 * 1024

// segmentSpeedSample converts a completed segment transfer into a speed
// sample. Returns ok=false when the transfer isn't a trustworthy
// measurement: too small, or so slow relative to realtime that it is
// indistinguishable from a paused client back-pressuring the socket.
func segmentSpeedSample(bytes int64, duration time.Duration, segmentSeconds int) (kbps, seconds float64, ok bool) {
	if bytes < minMeasurableSegmentBytes {
		return 0, 0, false
	}
	// A transfer taking >=3x realtime means a paused/backgrounded client
	// stopped draining its socket, not a slow connection. Genuine strugglers
	// cluster just above 1x and are kept.
	if segmentSeconds > 0 && duration.Seconds() >= float64(3*segmentSeconds) {
		return 0, 0, false
	}

	seconds = math.Max(duration.Seconds(), 0.001)
	kbps = math.Round(float64(bytes) * 8 / seconds / 1000)
	return kbps, seconds, true
}

// Serve-rate readings above this aren't measurements of any real viewer
// connection: they mean a local relay (reverse proxy, CDN edge, tunnel
// agent) drained the socket instead of the viewer, which reads as
// multi-gigabit "speed". No honest viewer health signal lives above this
// ceiling — stream bitrates top out orders of magnitude below it.
const maxPlausibleViewerKbps = 50000

// registerServedSegmentMetrics records a server-observed sample for a
// served segment. durationOnly is set for CMCD clients: their reported
// throughput owns the speed metric, but the serve timing still provides
// the download duration.
func (h *Handlers) registerServedSegmentMetrics(r *http.Request, clientID string, bytes int64, duration time.Duration, durationOnly bool) {
	// An aborted transfer isn't a complete download, so it isn't a speed
	// sample. Aborts are normal player behavior (quality switches, seeks,
	// leaving), not errors.
	if r.Context().Err() != nil {
		return
	}

	segmentSeconds := h.configRepository.GetStreamLatencyLevel().SecondsPerSegment
	kbps, seconds, ok := segmentSpeedSample(bytes, duration, segmentSeconds)
	if !ok {
		return
	}

	// A reading above the plausibility ceiling means a local relay drained
	// the socket, not the viewer — neither the speed nor the duration
	// measured anything about the viewer, so drop the whole sample.
	if kbps > maxPlausibleViewerKbps {
		return
	}

	report := metrics.PlaybackReport{
		ClientID:        clientID,
		DownloadSeconds: seconds,
	}
	if !durationOnly {
		report.BandwidthKbps = kbps
	}
	h.metrics.RegisterPlaybackReport(report)
}

// countingResponseWriter counts the bytes written to the client so a
// segment transfer can be turned into a speed sample.
type countingResponseWriter struct {
	http.ResponseWriter
	bytes int64
}

func (c *countingResponseWriter) Write(p []byte) (int, error) {
	n, err := c.ResponseWriter.Write(p)
	c.bytes += int64(n)
	return n, err
}

// ReadFrom preserves the underlying writer's io.ReaderFrom fast path —
// kernel sendfile for file-backed responses like segments — while still
// counting the bytes sent. Without it, wrapping the writer would silently
// downgrade every measured segment serve to userspace copies.
func (c *countingResponseWriter) ReadFrom(src io.Reader) (int64, error) {
	if rf, ok := c.ResponseWriter.(io.ReaderFrom); ok {
		n, err := rf.ReadFrom(src)
		c.bytes += n
		return n, err
	}
	n, err := io.Copy(c.ResponseWriter, src)
	c.bytes += n
	return n, err
}
