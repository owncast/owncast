package apmodels

import (
	"encoding/json"
	"net/url"
	"testing"
)

// asPublicIRI is the audience IRI Mastodon's InteractionPolicyParser accepts
// as "anyone may quote". Asserted as a literal because it is the wire
// contract, not an implementation detail.
const asPublicIRI = "https://www.w3.org/ns/activitystreams#Public"

// jsonStringValue unwraps a JSON-LD property that may serialize as either a
// bare string or a single-element string array. Both forms are equivalent on
// the wire and remote parsers accept either.
func jsonStringValue(t *testing.T, key string, v interface{}) string {
	t.Helper()
	switch val := v.(type) {
	case string:
		return val
	case []interface{}:
		if len(val) == 1 {
			if s, ok := val[0].(string); ok {
				return s
			}
		}
	}
	t.Fatalf("%s: expected string or single-element string array, got %T: %v", key, v, v)
	return ""
}

// jsonObjectValue unwraps a JSON-LD property that may serialize as either an
// object or a single-element array holding an object.
func jsonObjectValue(t *testing.T, key string, v interface{}) map[string]interface{} {
	t.Helper()
	switch val := v.(type) {
	case map[string]interface{}:
		return val
	case []interface{}:
		if len(val) == 1 {
			if obj, ok := val[0].(map[string]interface{}); ok {
				return obj
			}
		}
	}
	t.Fatalf("%s: expected object or single-element object array, got %T: %v", key, v, v)
	return nil
}

// contextContains reports whether a serialized @context (a bare string or a
// mixed array of strings and alias objects) includes the given URI.
func contextContains(context interface{}, uri string) bool {
	switch v := context.(type) {
	case string:
		return v == uri
	case []interface{}:
		for _, elem := range v {
			if s, ok := elem.(string); ok && s == uri {
				return true
			}
		}
	}
	return false
}

// TestMakeQuoteAuthorizationSerialization pins the FEP-044f QuoteAuthorization
// stamp wire format that Mastodon's VerifyQuoteService checks when fetching
// the stamp: an unprefixed "QuoteAuthorization" type, the ActivityStreams
// context URI, and id/attributedTo/interactingObject/interactionTarget IRIs.
func TestMakeQuoteAuthorizationSerialization(t *testing.T) {
	const (
		stampIRI      = "https://live.example/federation/stamp/x1"
		actorIRI      = "https://live.example/federation/user/streamer"
		quotePostIRI  = "https://remote.example/users/alice/statuses/99"
		quotedPostIRI = "https://live.example/federation/abc"
	)
	mustParse := func(s string) *url.URL {
		u, err := url.Parse(s)
		if err != nil {
			t.Fatalf("bad fixture URL %q: %v", s, err)
		}
		return u
	}

	stamp := MakeQuoteAuthorization(mustParse(stampIRI), mustParse(actorIRI), mustParse(quotePostIRI), mustParse(quotedPostIRI))
	b, err := Serialize(stamp)
	if err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("serialized output is not valid JSON: %v", err)
	}

	// Mastodon matches the raw string "QuoteAuthorization". A vocabulary
	// prefix (e.g. "gts:QuoteAuthorization") would fail its type check.
	if got := jsonStringValue(t, "type", m["type"]); got != "QuoteAuthorization" {
		t.Errorf("type = %q, want exactly \"QuoteAuthorization\"", got)
	}

	if !contextContains(m["@context"], "https://www.w3.org/ns/activitystreams") {
		t.Errorf("@context %v is missing https://www.w3.org/ns/activitystreams, remote servers will reject the document", m["@context"])
	}

	for key, want := range map[string]string{
		"id":                stampIRI,
		"attributedTo":      actorIRI,
		"interactingObject": quotePostIRI,
		"interactionTarget": quotedPostIRI,
	} {
		if got := jsonStringValue(t, key, m[key]); got != want {
			t.Errorf("%s = %q, want %q", key, got, want)
		}
	}
}

// TestMakeNoteQuotableSerialization pins the FEP-044f interaction policy wire
// format on outbound notes: interactionPolicy.canQuote.automaticApproval must
// contain the ActivityStreams Public IRI, which is exactly the path and value
// Mastodon's InteractionPolicyParser digs out to decide a post is quotable.
func TestMakeNoteQuotableSerialization(t *testing.T) {
	noteIRI, _ := url.Parse("https://live.example/federation/abc")
	actorIRI, _ := url.Parse("https://live.example/federation/user/streamer")

	note := MakeNoteQuotable(MakeNote("text", noteIRI, actorIRI))
	b, err := Serialize(note)
	if err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("serialized output is not valid JSON: %v", err)
	}

	policy, ok := m["interactionPolicy"]
	if !ok {
		t.Fatalf("serialized note has no interactionPolicy: %s", b)
	}
	canQuote := jsonObjectValue(t, "interactionPolicy", policy)["canQuote"]
	if canQuote == nil {
		t.Fatalf("interactionPolicy has no canQuote: %s", b)
	}
	approval := jsonObjectValue(t, "canQuote", canQuote)["automaticApproval"]
	if approval == nil {
		t.Fatalf("canQuote has no automaticApproval: %s", b)
	}

	if got := jsonStringValue(t, "automaticApproval", approval); got != asPublicIRI {
		t.Errorf("automaticApproval = %q, want %q", got, asPublicIRI)
	}
}
