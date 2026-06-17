package inbox

import (
	"net/url"
	"testing"
)

func TestActorIRIFromActivity(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    string
		wantErr bool
	}{
		{
			name: "string actor",
			body: `{"type":"Offer","actor":"https://remote.example/federation/user/streamer"}`,
			want: "https://remote.example/federation/user/streamer",
		},
		{
			name: "object actor with id",
			body: `{"type":"Follow","actor":{"id":"https://remote.example/federation/user/streamer","type":"Service"}}`,
			want: "https://remote.example/federation/user/streamer",
		},
		{
			name: "array actor",
			body: `{"type":"Leave","actor":["https://remote.example/federation/user/streamer"]}`,
			want: "https://remote.example/federation/user/streamer",
		},
		{
			name:    "missing actor",
			body:    `{"type":"Offer"}`,
			wantErr: true,
		},
		{
			name:    "empty actor",
			body:    `{"type":"Offer","actor":""}`,
			wantErr: true,
		},
		{
			name:    "invalid json",
			body:    `{not json`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := actorIRIFromActivity([]byte(tt.body))
			if (err != nil) != tt.wantErr {
				t.Fatalf("actorIRIFromActivity() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("actorIRIFromActivity() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSameActorOrigin(t *testing.T) {
	keyOwner, _ := url.Parse("https://remote.example/federation/user/streamer")

	tests := []struct {
		name     string
		actorIRI string
		want     bool
	}{
		{
			name:     "same host exact actor",
			actorIRI: "https://remote.example/federation/user/streamer",
			want:     true,
		},
		{
			name:     "same host different path",
			actorIRI: "https://remote.example/federation/user/someoneelse",
			want:     true,
		},
		{
			name:     "same host case-insensitive",
			actorIRI: "https://REMOTE.example/federation/user/streamer",
			want:     true,
		},
		{
			// The core spoofing case: a different host claiming to be us.
			name:     "different host is rejected",
			actorIRI: "https://attacker.evil/federation/user/streamer",
			want:     false,
		},
		{
			name:     "empty actor is rejected",
			actorIRI: "",
			want:     false,
		},
		{
			name:     "non-http scheme different host is rejected",
			actorIRI: "javascript:alert(1)",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sameActorOrigin(tt.actorIRI, keyOwner); got != tt.want {
				t.Errorf("sameActorOrigin(%q) = %v, want %v", tt.actorIRI, got, tt.want)
			}
		})
	}
}
