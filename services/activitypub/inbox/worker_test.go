package inbox

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"code.superseriousbusiness.org/activity/streams"
	"code.superseriousbusiness.org/activity/streams/vocab"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/persistence"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/services/dispatcher"
)

func makeFakePerson() vocab.ActivityStreamsPerson {
	iri, _ := url.Parse("https://freedom.eagle/user/mrfoo")
	name := "Mr Foo"
	username := "foodawg"
	inbox, _ := url.Parse("https://fake.fediverse.server/user/mrfoo/inbox")
	userAvatarURL, _ := url.Parse("https://fake.fediverse.server/user/mrfoo/avatar.png")

	person := streams.NewActivityStreamsPerson()

	id := streams.NewJSONLDIdProperty()
	id.Set(iri)
	person.SetJSONLDId(id)

	nameProperty := streams.NewActivityStreamsNameProperty()
	nameProperty.AppendXMLSchemaString(name)
	person.SetActivityStreamsName(nameProperty)

	preferredUsernameProperty := streams.NewActivityStreamsPreferredUsernameProperty()
	preferredUsernameProperty.SetXMLSchemaString(username)
	person.SetActivityStreamsPreferredUsername(preferredUsernameProperty)

	inboxProp := streams.NewActivityStreamsInboxProperty()
	inboxProp.SetIRI(inbox)
	person.SetActivityStreamsInbox(inboxProp)

	image := streams.NewActivityStreamsImage()
	imgProp := streams.NewActivityStreamsUrlProperty()
	imgProp.AppendIRI(userAvatarURL)
	image.SetActivityStreamsUrl(imgProp)
	icon := streams.NewActivityStreamsIconProperty()
	icon.AppendActivityStreamsImage(image)
	person.SetActivityStreamsIcon(icon)

	return person
}

// testService is the inbox Service used by tests in this file. It's
// initialized in TestMain with a real in-memory datastore so handler
// methods that hit persistence/followers work.
var (
	testService   *Service
	testDatastore *datastore.Datastore
)

func TestMain(m *testing.M) {
	if err := os.Setenv("OWNCAST_ALLOW_INTERNAL_FEDERATION", "true"); err != nil {
		panic(err)
	}
	if err := os.Setenv("OWNCAST_INSECURE_SKIP_VERIFY", "true"); err != nil {
		panic(err)
	}
	ds, err := datastore.SetupPersistence(":memory:", os.TempDir())
	if err != nil {
		panic(err)
	}
	testDatastore = ds
	configRepository := configrepository.New(testDatastore)
	configRepository.SetServerURL("https://my.cool.site.biz")
	persistenceSvc := persistence.New(testDatastore, nil)
	testService = New(Deps{
		Persistence: persistenceSvc,
		Followers:   followersrepository.New(testDatastore),
	})
	os.Exit(m.Run())
}

func TestBlockedDomains(t *testing.T) {
	configRepository := configrepository.New(testDatastore)

	person := makeFakePerson()

	configRepository.SetBlockedFederatedDomains([]string{"freedom.eagle", "guns.life"})

	if len(configRepository.GetBlockedFederatedDomains()) != 2 {
		t.Error("Blocked federated domains is not set correctly")
	}

	for _, domain := range configRepository.GetBlockedFederatedDomains() {
		if domain == person.GetJSONLDId().GetIRI().Host {
			return
		}
	}

	t.Error("Failed to catch blocked domain")
}

func TestBlockedActors(t *testing.T) {
	person := makeFakePerson()
	fakeRequest := streams.NewActivityStreamsFollow()
	followersRepository := followersrepository.New(testDatastore)
	followersRepository.Add(apmodels.ActivityPubActor{
		ActorIri:         person.GetJSONLDId().GetIRI(),
		Inbox:            person.GetJSONLDId().GetIRI(),
		FollowRequestIri: person.GetJSONLDId().GetIRI(),
		RequestObject:    fakeRequest,
	}, false)
	followersRepository.BlockOrReject(person.GetJSONLDId().GetIRI().String())

	blocked, err := testService.isBlockedActor(person.GetJSONLDId().GetIRI())
	if err != nil {
		t.Error(err)
		return
	}

	if !blocked {
		t.Error("Failed to block actor")
	}

	failedBlockIRI, _ := url.Parse("https://freedom.eagle/user/mrbar")
	failedBlock, err := testService.isBlockedActor(failedBlockIRI)

	if failedBlock {
		t.Error("Invalid blocking of unblocked actor IRI")
	}
}

func TestHandleVerifiedIngress(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		keyOwner  string
		wantEvent bool
	}{
		{
			name:      "matching actor",
			body:      `{"type":"Like","actor":"https://remote.example/users/alice"}`,
			keyOwner:  "https://remote.example/users/alice#main-key",
			wantEvent: true,
		},
		{
			name:      "unknown activity type",
			body:      `{ "type": "FutureActivity", "actor": "https://remote.example/users/alice", "custom": true }`,
			keyOwner:  "https://remote.example/users/alice#main-key",
			wantEvent: true,
		},
		{
			name:     "mismatched actor",
			body:     `{"type":"Like","actor":"https://forged.example/users/alice"}`,
			keyOwner: "https://remote.example/users/alice#main-key",
		},
		{
			name:     "malformed activity",
			body:     `{"type":"Like","actor":`,
			keyOwner: "https://remote.example/users/alice#main-key",
		},
		{
			name:     "activity without actor",
			body:     `{"type":"Like"}`,
			keyOwner: "https://remote.example/users/alice#main-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventDispatcher := dispatcher.New()
			var published []dispatcher.Event
			eventDispatcher.AddListener(func(_ context.Context, event dispatcher.Event) {
				published = append(published, event)
			})
			service := New(Deps{Events: eventDispatcher})
			keyOwner, err := url.Parse(tt.keyOwner)
			if err != nil {
				t.Fatal(err)
			}
			body := []byte(tt.body)

			accepted := service.handleVerifiedIngress(context.Background(), body, keyOwner)

			if accepted != tt.wantEvent {
				t.Fatalf("accepted = %v, want %v", accepted, tt.wantEvent)
			}
			if !tt.wantEvent {
				if len(published) != 0 {
					t.Fatalf("published %d events, want 0", len(published))
				}
				return
			}
			if len(published) != 1 {
				t.Fatalf("published %d events, want 1", len(published))
			}
			if published[0].Type != models.FediverseActivity {
				t.Fatalf("event type = %q, want %q", published[0].Type, models.FediverseActivity)
			}
			payload, ok := published[0].Payload.(json.RawMessage)
			if !ok {
				t.Fatalf("payload type = %T, want json.RawMessage", published[0].Payload)
			}
			body[0] = '!'
			if string(payload) != tt.body {
				t.Fatalf("payload = %q, want exact raw JSON %q", payload, tt.body)
			}
		})
	}
}

func TestVerifyRequestBody(t *testing.T) {
	body := []byte(`{"type":"Like"}`)
	sum := sha256.Sum256(body)
	digest := "SHA-256=" + base64.StdEncoding.EncodeToString(sum[:])

	tests := []struct {
		name    string
		digest  string
		body    []byte
		wantErr bool
	}{
		{name: "matching digest", digest: digest, body: body},
		{name: "whitespace", digest: " SHA-256 = " + base64.StdEncoding.EncodeToString(sum[:]) + " ", body: body},
		{name: "multiple digests", digest: "SHA-512=invalid, " + digest, body: body},
		{name: "changed body", digest: digest, body: []byte(`{"type":"Follow"}`), wantErr: true},
		{name: "missing digest", body: body, wantErr: true},
		{name: "malformed digest", digest: "SHA-256", body: body, wantErr: true},
		{name: "unsupported algorithm", digest: "MD5=abc", body: body, wantErr: true},
		{name: "invalid encoding", digest: "SHA-256=not-base64", body: body, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "https://local.example/inbox", nil)
			if tt.digest != "" {
				request.Header.Set("Digest", tt.digest)
			}
			if err := verifyRequestBody(request, tt.body); (err != nil) != tt.wantErr {
				t.Fatalf("verifyRequestBody() error = %v, want error: %v", err, tt.wantErr)
			}
		})
	}
}

func TestSignedHeaders(t *testing.T) {
	signature := `keyId="https://remote.example/actor#key",headers="(request-target) host date digest",signature="abc"`
	headers, err := signedHeaders(signature)
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"(request-target)", "host", "date", "digest"} {
		if _, ok := headers[required]; !ok {
			t.Errorf("signed headers missing %q", required)
		}
	}
}

func TestVerifyRequiresSignedDigest(t *testing.T) {
	request := httptest.NewRequest("POST", "https://local.example/inbox", nil)
	request.Header.Set("Signature", `keyId="https://remote.example/actor#main-key",algorithm="rsa-sha256",headers="(request-target) host date",signature="abc"`)

	_, err := (&Service{}).Verify(request, []byte(`{"type":"Like"}`))
	if err == nil || !strings.Contains(err.Error(), `does not sign required header "digest"`) {
		t.Fatalf("Verify() error = %v, want missing signed Digest error", err)
	}
}

func TestHandleDoesNotPublishUnverifiedRequest(t *testing.T) {
	eventDispatcher := dispatcher.New()
	published := 0
	eventDispatcher.AddListener(func(_ context.Context, _ dispatcher.Event) {
		published++
	})
	service := New(Deps{Events: eventDispatcher})
	body := []byte(`{"type":"Unknown","actor":"https://remote.example/users/alice"}`)
	request := httptest.NewRequest("POST", "https://local.example/inbox", nil)

	service.handle(apmodels.InboxRequest{Request: request, Body: body})

	if published != 0 {
		t.Fatalf("published %d events for an unverified request, want 0", published)
	}
}
