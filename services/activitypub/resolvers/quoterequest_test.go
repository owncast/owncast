package resolvers

import (
	"context"
	"testing"

	"code.superseriousbusiness.org/activity/streams/vocab"

	"github.com/owncast/owncast/services/activitypub/apmodels"
)

// TestResolveQuoteRequestDispatch feeds raw inbound QuoteRequest JSON, shaped
// exactly as real fediverse servers emit it, through Resolve and verifies it
// dispatches to a vocab.GoToSocialQuoteRequest callback with the actor,
// object, and instrument IRIs intact. Resolve is pure JSON dispatch with no
// network I/O, so the dependency-free testResolver is enough.
//
// This is load-bearing for Mastodon interop: Mastodon 4.5 aliases the type to
// https://w3id.org/fep/044f#QuoteRequest and embeds the quote post as an
// object under instrument, while GoToSocial uses its own ns alias and a plain
// IRI instrument. All shapes must reach the handler. Resolve swallows
// unmatched-type errors and returns nil, so the fired flag, not the error, is
// what proves dispatch.
func TestResolveQuoteRequestDispatch(t *testing.T) {
	const (
		actorIRI      = "https://remote.example/users/alice"
		quotedPostIRI = "https://live.example/federation/abc"
		quotePostIRI  = "https://remote.example/users/alice/statuses/99"
	)

	tests := []struct {
		name string
		body string
	}{
		{
			name: "mastodon 4.5: fep-044f type alias, embedded instrument object",
			body: `{
				"@context": ["https://www.w3.org/ns/activitystreams", {"QuoteRequest": "https://w3id.org/fep/044f#QuoteRequest"}],
				"type": "QuoteRequest",
				"id": "https://remote.example/quote_requests/1",
				"actor": "https://remote.example/users/alice",
				"object": "https://live.example/federation/abc",
				"instrument": {"type": "Note", "id": "https://remote.example/users/alice/statuses/99"}
			}`,
		},
		{
			name: "gotosocial: bare ns context, plain IRI instrument",
			body: `{
				"@context": ["https://www.w3.org/ns/activitystreams", "https://gotosocial.org/ns"],
				"type": "QuoteRequest",
				"id": "https://remote.example/quote_requests/2",
				"actor": "https://remote.example/users/alice",
				"object": "https://live.example/federation/abc",
				"instrument": "https://remote.example/users/alice/statuses/99"
			}`,
		},
		{
			name: "gotosocial ns type alias in context object",
			body: `{
				"@context": ["https://www.w3.org/ns/activitystreams", {"QuoteRequest": "https://gotosocial.org/ns#QuoteRequest"}],
				"type": "QuoteRequest",
				"id": "https://remote.example/quote_requests/3",
				"actor": "https://remote.example/users/alice",
				"object": "https://live.example/federation/abc",
				"instrument": "https://remote.example/users/alice/statuses/99"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got vocab.GoToSocialQuoteRequest
			callback := func(_ context.Context, activity vocab.GoToSocialQuoteRequest) error {
				got = activity
				return nil
			}

			if err := testResolver.Resolve(context.Background(), []byte(tt.body), callback); err != nil {
				t.Fatalf("Resolve returned error: %v", err)
			}
			if got == nil {
				t.Fatal("QuoteRequest callback did not fire, payload was not dispatched")
			}

			actor, err := apmodels.GetIRIStringFromActorProperty(got.GetActivityStreamsActor())
			if err != nil {
				t.Fatalf("GetIRIStringFromActorProperty: %v", err)
			}
			if actor != actorIRI {
				t.Errorf("actor = %q, want %q", actor, actorIRI)
			}

			object, err := apmodels.GetIRIFromObjectProperty(got.GetActivityStreamsObject())
			if err != nil {
				t.Fatalf("GetIRIFromObjectProperty: %v", err)
			}
			if object.String() != quotedPostIRI {
				t.Errorf("object = %q, want %q", object, quotedPostIRI)
			}

			instrument, err := apmodels.GetIRIFromInstrumentProperty(got.GetActivityStreamsInstrument())
			if err != nil {
				t.Fatalf("GetIRIFromInstrumentProperty: %v", err)
			}
			if instrument.String() != quotePostIRI {
				t.Errorf("instrument = %q, want %q", instrument, quotePostIRI)
			}
		})
	}
}
