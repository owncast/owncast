package inbox

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"testing"

	"code.superseriousbusiness.org/activity/streams"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/dispatcher"
)

func TestOutboundQuoteResponsesAfterQueueing(t *testing.T) {
	for _, mode := range []string{"accepted", "rejected", "queue failure"} {
		t.Run(mode, func(t *testing.T) {
			service, persistenceService, _ := newEngagementTestService(t)
			if mode != "rejected" {
				if err := persistenceService.AddToOutbox("https://owncast.example/federation/quoted", []byte(`{"type":"Note"}`), "Note", false); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "queue failure" {
				if _, err := persistenceService.Datastore().DB.Exec(`CREATE TRIGGER reject_delivery BEFORE INSERT ON ap_delivery_queue BEGIN SELECT RAISE(ABORT, 'forced queue failure'); END`); err != nil {
					t.Fatal(err)
				}
			}
			var activities []string
			service.events.AddListener(func(_ context.Context, event dispatcher.Event) {
				if event.Type != models.FediverseOutboundActivity {
					return
				}
				var envelope struct {
					Type string `json:"type"`
				}
				if err := json.Unmarshal(event.Payload.(json.RawMessage), &envelope); err != nil {
					t.Fatal(err)
				}
				var queued int
				if err := persistenceService.Datastore().DB.QueryRow(`SELECT count(*) FROM ap_delivery_queue`).Scan(&queued); err != nil || queued != 1 {
					t.Fatalf("observer queued responses = %d: %v", queued, err)
				}
				activities = append(activities, envelope.Type)
			})
			err := service.handleQuoteRequestInboxRequest(context.Background(), quoteRequest("https://remote.example/inbox"))
			if mode == "queue failure" {
				if err == nil || !strings.Contains(err.Error(), "forced queue failure") || len(activities) != 0 {
					t.Fatalf("failed response = %v, activities = %v", err, activities)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			want := "Accept"
			if mode == "rejected" {
				want = "Reject"
			}
			if len(activities) != 1 || activities[0] != want {
				t.Fatalf("response activities = %v, want %s", activities, want)
			}
		})
	}
}

func TestOutboundFollowAcceptAfterCommit(t *testing.T) {
	for _, mode := range []string{"approved", "pending", "queue failure"} {
		t.Run(mode, func(t *testing.T) {
			service, persistenceService, _ := newEngagementTestService(t)
			ds := persistenceService.Datastore()
			actor, _ := url.Parse("https://remote.example/users/alice")
			inbox, _ := url.Parse("https://remote.example/inbox")
			follow := streams.NewActivityStreamsFollow()
			id := streams.NewJSONLDIdProperty()
			id.SetIRI(actor.JoinPath("follow"))
			follow.SetJSONLDId(id)
			actorProperty := streams.NewActivityStreamsActorProperty()
			actorProperty.AppendIRI(actor)
			follow.SetActivityStreamsActor(actorProperty)
			object := streams.NewActivityStreamsObjectProperty()
			object.AppendIRI(service.builder.MakeLocalIRIForAccount(service.configRepository.GetFederationUsername()))
			follow.SetActivityStreamsObject(object)
			request := apmodels.ActivityPubActor{ActorIri: actor, Inbox: inbox, Username: "alice", RequestObject: follow}
			if mode == "queue failure" {
				if _, err := ds.DB.Exec(`CREATE TRIGGER reject_delivery BEFORE INSERT ON ap_delivery_queue BEGIN SELECT RAISE(ABORT, 'forced queue failure'); END`); err != nil {
					t.Fatal(err)
				}
			}
			var responses int
			service.events.AddListener(func(_ context.Context, event dispatcher.Event) {
				if event.Type != models.FediverseOutboundActivity {
					return
				}
				if !ds.DbLock.TryLock() {
					t.Fatal("observer invoked under datastore lock")
				}
				ds.DbLock.Unlock()
				stored, err := service.followers.GetByIRI(actor.String())
				if err != nil || stored.ApprovedAt == nil {
					t.Fatalf("observer could not read approved follow: %v", err)
				}
				var envelope struct {
					Type string `json:"type"`
				}
				if err := json.Unmarshal(event.Payload.(json.RawMessage), &envelope); err != nil || envelope.Type != "Accept" {
					t.Fatalf("follow response = %+v: %v", envelope, err)
				}
				responses++
			})
			_, err := service.createOrHandleFollow(context.Background(), follow, &request, request, mode != "pending", actor.String(), service.configRepository.GetFederationUsername())
			if mode == "queue failure" {
				if err == nil || !strings.Contains(err.Error(), "forced queue failure") || responses != 0 {
					t.Fatalf("failed follow = %v, responses = %d", err, responses)
				}
				if _, err := service.followers.GetByIRI(actor.String()); err == nil {
					t.Fatal("failed follow transaction committed follower")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			want := 1
			if mode == "pending" {
				want = 0
			}
			if responses != want {
				t.Fatalf("responses = %d, want %d", responses, want)
			}
		})
	}
}
