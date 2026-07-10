package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestParseCMCDQueryParameter(t *testing.T) {
	// Example payload shape from the CTA-5004 spec.
	r := httptest.NewRequest("GET", "/hls/0/stream-abc-42.ts?CMCD="+
		"br%3D1200%2Cbs%2Cd%3D3000%2Cmtp%3D25400%2Cot%3Dv%2C"+
		"sid%3D%226e2fb550-c457-11e9-bb97-0800200c9a66%22%2Ctb%3D6000", nil)

	keys := parseCMCDRequest(r)

	if keys == nil {
		t.Fatal("expected CMCD data to be detected")
	}
	if mtp, _ := cmcdNumber(keys, "mtp"); mtp != 25400 {
		t.Errorf("expected mtp 25400, got %v", keys["mtp"])
	}
	if starved, _ := keys["bs"].(bool); !starved {
		t.Error("expected buffer starvation flag to be set")
	}
	if sid, _ := keys["sid"].(string); sid != "6e2fb550-c457-11e9-bb97-0800200c9a66" {
		t.Errorf("unexpected session id %v", keys["sid"])
	}
	if ot, _ := keys["ot"].(string); ot != "v" {
		t.Errorf("expected token value to parse as string, got %v", keys["ot"])
	}
}

func TestParseCMCDHeaders(t *testing.T) {
	r := httptest.NewRequest("GET", "/hls/0/stream-abc-42.ts", nil)
	r.Header.Set("CMCD-Request", "mtp=48100")
	r.Header.Set("CMCD-Session", `sid="4f5a6b7c",cid="with,comma and \"quote\""`)
	r.Header.Set("CMCD-Status", "bs")

	keys := parseCMCDRequest(r)

	if keys == nil {
		t.Fatal("expected CMCD data to be detected")
	}
	if mtp, _ := cmcdNumber(keys, "mtp"); mtp != 48100 {
		t.Errorf("expected mtp 48100, got %v", keys["mtp"])
	}
	if sid, _ := keys["sid"].(string); sid != "4f5a6b7c" {
		t.Errorf("unexpected session id %v", keys["sid"])
	}
	if cid, _ := keys["cid"].(string); cid != `with,comma and "quote"` {
		t.Errorf("unexpected content id %v", keys["cid"])
	}
	if starved, _ := keys["bs"].(bool); !starved {
		t.Error("expected buffer starvation flag to be set")
	}
}

func TestParseCMCDV2LiveLatency(t *testing.T) {
	// CMCD v2 (CTA-5004-A) request-mode payload with live latency in ms.
	r := httptest.NewRequest("GET", "/hls/0/stream-abc-42.ts?CMCD="+
		"ltc%3D4200%2Cmtp%3D12300%2Cpr%3D1%2Csid%3D%22v2-player%22%2Cv%3D2", nil)

	keys := parseCMCDRequest(r)

	if ltc, _ := cmcdNumber(keys, "ltc"); ltc != 4200 {
		t.Errorf("expected ltc 4200, got %v", keys["ltc"])
	}
	if v, _ := cmcdNumber(keys, "v"); v != 2 {
		t.Errorf("expected version 2, got %v", keys["v"])
	}
	if sid := cmcdClientID(r, keys); sid != "v2-player" {
		t.Errorf("expected sid identity, got %q", sid)
	}
}

func TestParseCMCDAbsent(t *testing.T) {
	r := httptest.NewRequest("GET", "/hls/0/stream-abc-42.ts", nil)
	if keys := parseCMCDRequest(r); keys != nil {
		t.Error("expected no CMCD data on a bare request")
	}
	if id := cmcdClientID(r, map[string]any{}); id == "" {
		t.Error("expected fallback client id when sid is missing")
	}
}

func TestReportCmcdPreflight(t *testing.T) {
	// A cross-origin beacon preflights with OPTIONS; it must be answered
	// with permissive CORS headers before any report processing.
	h := &Handlers{}
	r := httptest.NewRequest("OPTIONS", "/api/metrics/cmcd", nil)
	rec := httptest.NewRecorder()

	h.ReportCmcd(rec, r)

	if rec.Code != 204 {
		t.Errorf("expected 204, got %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected wildcard allow-origin")
	}
	if got := rec.Header().Get("Access-Control-Allow-Headers"); !strings.Contains(got, "Content-Type") {
		t.Errorf("expected Content-Type in allow-headers, got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); !strings.Contains(got, "POST") {
		t.Errorf("expected POST in allow-methods, got %q", got)
	}
}

func TestParseCmcdReportsJSON(t *testing.T) {
	// A batched event-mode POST: an array of report objects.
	body := `[{"e":"ps","ts":1750000000000,"sta":"p","ltc":3100,"sid":"s1","sn":1},` +
		`{"e":"t","mtp":8200,"bs":true,"sid":"s1","sn":2}]`
	r := httptest.NewRequest("POST", "/api/metrics/cmcd", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")

	reports, err := parseCmcdReports(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reports) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(reports))
	}
	if ltc, _ := cmcdNumber(reports[0], "ltc"); ltc != 3100 {
		t.Errorf("expected ltc 3100, got %v", reports[0]["ltc"])
	}
	if starved, _ := reports[1]["bs"].(bool); !starved {
		t.Error("expected bs true in second report")
	}

	// A single report object is also accepted.
	r = httptest.NewRequest("POST", "/api/metrics/cmcd", strings.NewReader(`{"e":"e","ec":"3","sid":"s2"}`))
	r.Header.Set("Content-Type", "application/json")
	reports, err = parseCmcdReports(r)
	if err != nil || len(reports) != 1 {
		t.Fatalf("expected 1 report, got %d (err %v)", len(reports), err)
	}
	if event, _ := reports[0]["e"].(string); event != "e" {
		t.Errorf("expected error event, got %v", reports[0]["e"])
	}
}

func TestParseCmcdReportsQuery(t *testing.T) {
	r := httptest.NewRequest("GET", "/api/metrics/cmcd?CMCD=e%3Dt%2Cmtp%3D5100%2Csid%3D%22s3%22", nil)
	reports, err := parseCmcdReports(r)
	if err != nil || len(reports) != 1 {
		t.Fatalf("expected 1 report, got %d (err %v)", len(reports), err)
	}
	if mtp, _ := cmcdNumber(reports[0], "mtp"); mtp != 5100 {
		t.Errorf("expected mtp 5100, got %v", reports[0]["mtp"])
	}
}

func TestSegmentSpeedSample(t *testing.T) {
	const oneMB = 1024 * 1024

	// A 1MB segment delivered in 1s is 8389 kbps.
	kbps, seconds, ok := segmentSpeedSample(oneMB, time.Second, 3)
	if !ok {
		t.Fatal("expected a valid sample")
	}
	if kbps != 8389 {
		t.Errorf("expected 8389 kbps, got %f", kbps)
	}
	if seconds != 1.0 {
		t.Errorf("expected 1s, got %f", seconds)
	}

	// A struggling-but-watching viewer (slower than realtime, under the
	// pause cutoff) must be kept: that's the population health reporting
	// exists to surface.
	if _, _, ok := segmentSpeedSample(oneMB, 5*time.Second, 3); !ok {
		t.Error("expected a sample slower than realtime to be kept")
	}

	// At or beyond 3x realtime it's indistinguishable from a paused client
	// back-pressuring the socket: discarded, not scored.
	if _, _, ok := segmentSpeedSample(oneMB, 9*time.Second, 3); ok {
		t.Error("expected a >=3x realtime sample to be discarded")
	}

	// Too small to measure anything but kernel buffers.
	if _, _, ok := segmentSpeedSample(10*1024, time.Second, 3); ok {
		t.Error("expected a tiny transfer to be discarded")
	}

	// Sub-millisecond serves are clamped rather than dividing by ~zero.
	if _, seconds, ok := segmentSpeedSample(oneMB, 0, 3); !ok || seconds != 0.001 {
		t.Errorf("expected clamped duration 0.001, got ok=%v seconds=%f", ok, seconds)
	}
}
