package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"code.superseriousbusiness.org/activity/streams"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	apcrypto "github.com/owncast/owncast/services/activitypub/crypto"
	"github.com/owncast/owncast/services/activitypub/persistence"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/services/activitypub/resolvers"
	"github.com/owncast/owncast/services/activitypub/workerpool"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/services/dispatcher"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("OWNCAST_ALLOW_INTERNAL_FEDERATION", "true")
	_ = os.Setenv("OWNCAST_INSECURE_SKIP_VERIFY", "true")
	os.Exit(m.Run())
}

type outboundFixture struct {
	service *Service
	ds      *datastore.Datastore
	cfg     configrepository.ConfigRepository
	events  []json.RawMessage
}

func newOutboundFixture(t *testing.T) *outboundFixture {
	t.Helper()
	directory := t.TempDir()
	ds, err := datastore.SetupPersistence(filepath.Join(directory, "owncast.db"), directory)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ds.DB.Close() })
	cfg := configrepository.New(ds)
	privateKey, publicKey, err := apcrypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}
	for _, err := range []error{
		cfg.SetServerURL("https://owncast.example"), cfg.SetFederationEnabled(true),
		cfg.SetPrivateKey(string(privateKey)), cfg.SetPublicKey(string(publicKey)),
		cfg.SetFederationGoLiveMessage("We are live"),
	} {
		if err != nil {
			t.Fatal(err)
		}
	}
	builder := apmodels.New(apmodels.Deps{ConfigRepository: cfg})
	signer := apcrypto.New(apcrypto.Deps{ConfigRepository: cfg})
	resolver := resolvers.New(resolvers.Deps{ConfigRepository: cfg, Builder: builder, Signer: signer})
	fixture := &outboundFixture{ds: ds, cfg: cfg}
	events := dispatcher.New()
	events.AddListener(func(_ context.Context, event dispatcher.Event) {
		if event.Type == models.FediverseOutboundActivity {
			fixture.events = append(fixture.events, event.Payload.(json.RawMessage))
		}
	})
	fixture.service = New(Deps{
		Persistence: persistence.New(ds, resolver),
		Workerpool:  workerpool.New(workerpool.Deps{Datastore: ds, Signer: signer, WorkerPoolSize: 1}),
		Followers:   followersrepository.New(ds), ConfigRepository: cfg,
		Builder: builder, Resolver: resolver, Config: &config.Config{TempDir: directory}, Events: events,
	})
	return fixture
}

func (f *outboundFixture) addFollower(t *testing.T, actorIRI, inboxIRI string) {
	t.Helper()
	actor, _ := url.Parse(actorIRI)
	inbox, _ := url.Parse(inboxIRI)
	follow := streams.NewActivityStreamsFollow()
	id := streams.NewJSONLDIdProperty()
	id.SetIRI(actor.JoinPath("follow"))
	follow.SetJSONLDId(id)
	actorProperty := streams.NewActivityStreamsActorProperty()
	actorProperty.AppendIRI(actor)
	follow.SetActivityStreamsActor(actorProperty)
	object := streams.NewActivityStreamsObjectProperty()
	object.AppendIRI(f.service.builder.MakeLocalIRIForAccount(f.cfg.GetFederationUsername()))
	follow.SetActivityStreamsObject(object)
	if err := f.service.followers.Add(apmodels.ActivityPubActor{
		ActorIri: actor, Inbox: inbox, Username: "follower", RequestObject: follow,
	}, true); err != nil {
		t.Fatal(err)
	}
}

func TestOutboundPostCommitsBeforeNotification(t *testing.T) {
	for _, followers := range []int{0, 2} {
		t.Run(fmt.Sprintf("%d followers", followers), func(t *testing.T) {
			f := newOutboundFixture(t)
			for i := range followers {
				f.addFollower(t, fmt.Sprintf("https://remote.example/users/%d", i), fmt.Sprintf("https://remote.example/inbox/%d", i))
			}
			f.service.events.AddListener(func(_ context.Context, event dispatcher.Event) {
				var activity struct {
					Object struct {
						ID string `json:"id"`
					} `json:"object"`
				}
				if err := json.Unmarshal(event.Payload.(json.RawMessage), &activity); err != nil {
					t.Fatal(err)
				}
				if _, _, _, err := f.service.persistence.GetNoteByIRI(activity.Object.ID); err != nil {
					t.Fatalf("observer could not read committed post: %v", err)
				}
				var queued int
				if err := f.ds.DB.QueryRow(`SELECT count(*) FROM ap_delivery_queue`).Scan(&queued); err != nil || queued != followers {
					t.Fatalf("queued = %d, want %d: %v", queued, followers, err)
				}
			})
			if err := f.service.SendPublicMessage("Hello outbound"); err != nil {
				t.Fatal(err)
			}
			if len(f.events) != 1 || !strings.Contains(string(f.events[0]), "Hello outbound") {
				t.Fatalf("outbound post events = %s", f.events)
			}
		})
	}
}

func TestOutboundFailedPersistenceDoesNotNotify(t *testing.T) {
	for _, table := range []string{"ap_outbox", "ap_delivery_queue"} {
		t.Run(table, func(t *testing.T) {
			f := newOutboundFixture(t)
			f.addFollower(t, "https://remote.example/users/alice", "https://remote.example/inbox")
			if _, err := f.ds.DB.Exec(`CREATE TRIGGER reject_insert BEFORE INSERT ON ` + table + ` BEGIN SELECT RAISE(ABORT, 'forced persistence failure'); END`); err != nil {
				t.Fatal(err)
			}
			if err := f.service.SendPublicMessage("Must not crosspost"); err == nil || !strings.Contains(err.Error(), "forced persistence failure") {
				t.Fatalf("publication error = %v", err)
			}
			if len(f.events) != 0 {
				t.Fatalf("failed publication notified plugins: %s", f.events)
			}
		})
	}
}

func remoteOutboundActor(t *testing.T, f *outboundFixture) *httptest.Server {
	t.Helper()
	var remote *httptest.Server
	remote = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/activity+json")
		if r.URL.Path == "/.well-known/webfinger" {
			_ = json.NewEncoder(w).Encode(map[string]any{"links": []map[string]string{{"rel": "self", "type": "application/activity+json", "href": remote.URL + "/users/alice"}}})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"@context": []string{"https://www.w3.org/ns/activitystreams", "https://w3id.org/security/v1"},
			"id":       remote.URL + "/users/alice", "type": "Person", "preferredUsername": "alice",
			"inbox": remote.URL + "/inbox", "publicKey": map[string]string{
				"id": remote.URL + "/users/alice#main-key", "owner": remote.URL + "/users/alice", "publicKeyPem": f.cfg.GetPublicKey(),
			},
		})
	}))
	t.Cleanup(remote.Close)
	return remote
}

func TestOutboundDirectMessagesExcludeAuthentication(t *testing.T) {
	f := newOutboundFixture(t)
	remote := remoteOutboundActor(t, f)
	account := "alice@" + strings.TrimPrefix(remote.URL, "https://")
	if err := f.service.SendAuthenticationMessageToAccount("One-time code: 123456", account); err != nil {
		t.Fatal(err)
	}
	if len(f.events) != 0 {
		t.Fatalf("authentication code reached plugins: %s", f.events)
	}
	if err := f.service.SendDirectMessageToAccount("Ordinary private message", account); err != nil {
		t.Fatal(err)
	}
	if len(f.events) != 1 || !strings.Contains(string(f.events[0]), "Ordinary private message") || strings.Contains(string(f.events[0]), "123456") {
		t.Fatalf("private activity events = %s", f.events)
	}
	var queued int
	if err := f.ds.DB.QueryRow(`SELECT count(*) FROM ap_delivery_queue`).Scan(&queued); err != nil || queued != 2 {
		t.Fatalf("queued = %d, want both messages: %v", queued, err)
	}
}

func TestOutboundFollowAndUndoAfterPersistence(t *testing.T) {
	f := newOutboundFixture(t)
	remote := remoteOutboundActor(t, f)
	actor, _ := url.Parse(remote.URL + "/users/alice")
	inbox, _ := url.Parse(remote.URL + "/inbox")
	f.service.events.AddListener(func(_ context.Context, _ dispatcher.Event) {
		if !f.ds.DbLock.TryLock() {
			t.Error("outbound observer invoked under datastore lock")
			return
		}
		f.ds.DbLock.Unlock()
		var status string
		if err := f.ds.DB.QueryRow(`SELECT follow_status FROM federated_servers WHERE iri = ?`, remote.URL).Scan(&status); err != nil || status != "pending" {
			t.Errorf("committed follow status = %q: %v", status, err)
		}
	})
	if err := f.service.SendFollowToAccountURI(actor.String(), "alice", actor.Host, false); err != nil {
		t.Fatal(err)
	}
	if err := f.service.sendUndoFollow(actor, inbox); err != nil {
		t.Fatal(err)
	}
	for i, expected := range []string{"Follow", "Undo"} {
		if len(f.events) <= i || !strings.Contains(string(f.events[i]), `"type":"`+expected+`"`) {
			t.Fatalf("event %d should be %s: %s", i, expected, f.events)
		}
	}
	if len(f.events) != 2 {
		t.Fatalf("outbound follow events = %s", f.events)
	}
}

func TestOutboundRetryDoesNotNotifyAgain(t *testing.T) {
	f := newOutboundFixture(t)
	var attempts atomic.Int32
	remote := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if attempts.Add(1) == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer remote.Close()
	f.addFollower(t, remote.URL+"/users/alice", remote.URL+"/inbox")
	if err := f.service.SendPublicMessage("Retry this delivery"); err != nil {
		t.Fatal(err)
	}
	f.service.workerpool.Start()
	defer f.service.workerpool.Stop()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		var remaining int
		if err := f.ds.DB.QueryRow(`SELECT count(*) FROM ap_delivery_queue`).Scan(&remaining); err != nil {
			t.Fatal(err)
		}
		if attempts.Load() >= 2 && remaining == 0 {
			if len(f.events) != 1 {
				t.Fatalf("retry notified plugins again: %s", f.events)
			}
			return
		}
		if attempts.Load() == 1 {
			if _, err := f.ds.DB.Exec(`UPDATE ap_delivery_queue SET next_attempt_at = ? WHERE attempts > 0`, time.Now().Add(-time.Second)); err != nil {
				t.Fatal(err)
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("delivery did not retry and drain, attempts = %d", attempts.Load())
}
