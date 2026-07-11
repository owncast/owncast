package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/router/middleware"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// ReportCmcd is the CMCD v2 (CTA-5004-A) collector endpoint. It accepts
// event and response mode reports either as a JSON body (a single report
// object or an array of batched reports, keyed by CMCD key names) or as a
// CMCD query parameter in dictionary payload syntax. Players embedded on
// other origins beacon here, so CORS is fully enabled.
func (h *Handlers) ReportCmcd(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)

	if r.Method == http.MethodOptions {
		// Cross-origin JSON POSTs preflight with OPTIONS; allow them.
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, CMCD-Request, CMCD-Object, CMCD-Status, CMCD-Session")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Reports are tiny; bound the unauthenticated body so a large batched
	// POST can't force unbounded allocation.
	r.Body = http.MaxBytesReader(w, r.Body, 64*1024)

	reports, err := parseCmcdReports(r)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	registered := false
	for _, keys := range reports {
		if len(keys) == 0 {
			continue
		}
		h.registerCMCDKeys(cmcdClientID(r, keys), keys)
		registered = true
	}

	if !registered {
		webutils.WriteSimpleResponse(w, false, "no CMCD report found in request")
		return
	}

	// A beaconing player's media requests come from the same client, so
	// suppress the lower-fidelity server-derived observation for it even
	// if it doesn't decorate its media requests. Only a request that
	// actually carried a report earns the suppression.
	h.metrics.RegisterSelfReportingClient(utils.GenerateClientIDFromRequest(r))

	w.WriteHeader(http.StatusOK)
}

// parseCmcdReports extracts CMCD reports from a collector request: a JSON
// object, a JSON array of objects, or a CMCD query parameter.
func parseCmcdReports(r *http.Request) ([]map[string]any, error) {
	if r.Method == http.MethodPost {
		var body any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return nil, err
		}

		switch payload := body.(type) {
		case []any:
			reports := make([]map[string]any, 0, len(payload))
			for _, entry := range payload {
				if keys, ok := entry.(map[string]any); ok {
					reports = append(reports, keys)
				}
			}
			return reports, nil
		case map[string]any:
			return []map[string]any{payload}, nil
		default:
			return nil, nil
		}
	}

	if payload := r.URL.Query().Get("CMCD"); payload != "" {
		keys := map[string]any{}
		parseCMCDPayload(payload, keys)
		return []map[string]any{keys}, nil
	}

	return nil, nil
}
