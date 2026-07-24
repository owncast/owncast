package inbox

import (
	"context"
	"sync"
	"testing"

	"code.superseriousbusiness.org/activity/streams"
	"code.superseriousbusiness.org/activity/streams/vocab"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	apcrypto "github.com/owncast/owncast/services/activitypub/crypto"
	activityevents "github.com/owncast/owncast/services/activitypub/events"
	"github.com/owncast/owncast/services/activitypub/persistence"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	apresolvers "github.com/owncast/owncast/services/activitypub/resolvers"
	"github.com/owncast/owncast/services/activitypub/workerpool"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/services/dispatcher"
)

func newEngagementTestService(t *testing.T) (*Service, *persistence.Service, *[]dispatcher.Event) {
	t.Helper()

	ds, err := datastore.SetupPersistence(":memory:", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ds.DB.Close() })

	configRepository := configrepository.New(ds)
	if err := configRepository.SetServerURL("https://owncast.example"); err != nil {
		t.Fatal(err)
	}
	if err := configRepository.SetChatDisabled(true); err != nil {
		t.Fatal(err)
	}
	if err := configRepository.SetFederationShowEngagement(false); err != nil {
		t.Fatal(err)
	}
	if err := configRepository.SetFederationEnableQuotes(true); err != nil {
		t.Fatal(err)
	}

	privateKey, publicKey, err := apcrypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}
	if err := configRepository.SetPrivateKey(string(privateKey)); err != nil {
		t.Fatal(err)
	}
	if err := configRepository.SetPublicKey(string(publicKey)); err != nil {
		t.Fatal(err)
	}

	signer := apcrypto.New(apcrypto.Deps{ConfigRepository: configRepository})
	builder := apmodels.New(apmodels.Deps{ConfigRepository: configRepository, Signer: signer})
	resolver := apresolvers.New(apresolvers.Deps{ConfigRepository: configRepository, Builder: builder, Signer: signer})
	followers := followersrepository.New(ds)
	outbound := workerpool.New(workerpool.Deps{Datastore: ds, Signer: signer})
	outbound.Start()
	t.Cleanup(outbound.Stop)
	eventDispatcher := dispatcher.New()
	var publishedMu sync.Mutex
	published := []dispatcher.Event{}
	eventDispatcher.AddListener(func(_ context.Context, event dispatcher.Event) {
		publishedMu.Lock()
		defer publishedMu.Unlock()
		published = append(published, event)
	})
	persistenceService := persistence.New(ds, nil)

	return New(Deps{
		Persistence:      persistenceService,
		Workerpool:       outbound,
		ConfigRepository: configRepository,
		Followers:        followers,
		Builder:          builder,
		Resolver:         resolver,
		Events:           eventDispatcher,
	}), persistenceService, &published
}

func actorProperty(person vocab.ActivityStreamsPerson) vocab.ActivityStreamsActorProperty {
	publicKey := streams.NewW3IDSecurityV1PublicKeyProperty()
	publicKey.AppendW3IDSecurityV1PublicKey(streams.NewW3IDSecurityV1PublicKey())
	person.SetW3IDSecurityV1PublicKey(publicKey)

	property := streams.NewActivityStreamsActorProperty()
	property.AppendActivityStreamsPerson(person)
	return property
}

func objectProperty(iri string) vocab.ActivityStreamsObjectProperty {
	property := streams.NewActivityStreamsObjectProperty()
	property.AppendIRI(mustParseURL(iri))
	return property
}

func assertEngagementEvent(t *testing.T, published []dispatcher.Event, eventType models.EventType, target string) {
	t.Helper()
	if len(published) != 1 {
		t.Fatalf("published %d events, want 1", len(published))
	}
	if published[0].Type != eventType {
		t.Fatalf("event type = %q, want %q", published[0].Type, eventType)
	}
	payload, ok := published[0].Payload.(*activityevents.FediverseEngagementEvent)
	if !ok || payload == nil {
		t.Fatalf("payload type = %T, want *events.FediverseEngagementEvent", published[0].Payload)
	}
	if payload.Actor.Name != "Mr Foo" || payload.Actor.Handle != "foodawg@freedom.eagle" || payload.Actor.URL != "https://freedom.eagle/user/mrfoo" || payload.Actor.Image != "https://fake.fediverse.server/user/mrfoo/avatar.png" {
		t.Fatalf("unexpected actor payload: %+v", payload.Actor)
	}
	if payload.Target == nil || payload.Target.URL != target {
		t.Fatalf("target = %+v, want %q", payload.Target, target)
	}
}

func TestEngagementPublishesWithChatDisplayDisabledAndSkipsDuplicates(t *testing.T) {
	tests := []struct {
		name      string
		eventType models.EventType
		handle    func(*Service, string) error
	}{
		{
			name:      "like",
			eventType: models.FediverseEngagementLike,
			handle: func(service *Service, target string) error {
				activity := streams.NewActivityStreamsLike()
				activity.SetActivityStreamsActor(actorProperty(makeFakePerson()))
				activity.SetActivityStreamsObject(objectProperty(target))
				return service.handleLikeRequest(context.Background(), activity)
			},
		},
		{
			name:      "repost",
			eventType: models.FediverseEngagementRepost,
			handle: func(service *Service, target string) error {
				activity := streams.NewActivityStreamsAnnounce()
				activity.SetActivityStreamsActor(actorProperty(makeFakePerson()))
				activity.SetActivityStreamsObject(objectProperty(target))
				return service.handleAnnounceRequest(context.Background(), activity)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, persistenceService, published := newEngagementTestService(t)
			target := "https://owncast.example/federation/" + tt.name
			if err := persistenceService.AddToOutbox(target, []byte(`{"type":"Note"}`), "Note", false); err != nil {
				t.Fatal(err)
			}

			if err := tt.handle(service, target); err != nil {
				t.Fatalf("first activity failed: %v", err)
			}
			assertEngagementEvent(t, *published, tt.eventType, target)

			_ = tt.handle(service, target)
			if len(*published) != 1 {
				t.Fatalf("duplicate published an event; total = %d, want 1", len(*published))
			}
		})
	}
}

func TestConcurrentEngagementDeliveriesPublishOnce(t *testing.T) {
	service, persistenceService, published := newEngagementTestService(t)
	persistenceService.Datastore().DB.SetMaxOpenConns(1)
	const target = "https://owncast.example/federation/concurrent-like"
	if err := persistenceService.AddToOutbox(target, []byte(`{"type":"Note"}`), "Note", false); err != nil {
		t.Fatal(err)
	}

	activity := streams.NewActivityStreamsLike()
	activity.SetActivityStreamsActor(actorProperty(makeFakePerson()))
	activity.SetActivityStreamsObject(objectProperty(target))

	const deliveries = 16
	start := make(chan struct{})
	errs := make(chan error, deliveries)
	var wg sync.WaitGroup
	wg.Add(deliveries)
	for range deliveries {
		go func() {
			defer wg.Done()
			<-start
			if err := service.handleLikeRequest(context.Background(), activity); err != nil {
				errs <- err
			}
		}()
	}
	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		t.Fatalf("concurrent delivery: %v", err)
	}
	assertEngagementEvent(t, *published, models.FediverseEngagementLike, target)
}

func quoteRequest(inbox string) vocab.GoToSocialQuoteRequest {
	person := makeFakePerson()
	inboxProperty := streams.NewActivityStreamsInboxProperty()
	inboxProperty.SetIRI(mustParseURL(inbox))
	person.SetActivityStreamsInbox(inboxProperty)

	activity := streams.NewGoToSocialQuoteRequest()
	id := streams.NewJSONLDIdProperty()
	id.SetIRI(mustParseURL("https://remote.example/quote_requests/1"))
	activity.SetJSONLDId(id)
	activity.SetActivityStreamsActor(actorProperty(person))
	activity.SetActivityStreamsObject(objectProperty("https://owncast.example/federation/quoted"))
	instrument := streams.NewActivityStreamsInstrumentProperty()
	instrument.AppendIRI(mustParseURL("https://remote.example/posts/quote"))
	activity.SetActivityStreamsInstrument(instrument)
	return activity
}

func TestQuotePublishesOnlyAfterSuccessfulAccept(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, persistenceService, published := newEngagementTestService(t)
		const target = "https://owncast.example/federation/quoted"
		if err := persistenceService.AddToOutbox(target, []byte(`{"type":"Note"}`), "Note", false); err != nil {
			t.Fatal(err)
		}

		if err := service.handleQuoteRequestInboxRequest(context.Background(), quoteRequest("https://remote.example/inbox")); err != nil {
			t.Fatal(err)
		}
		assertEngagementEvent(t, *published, models.FediverseEngagementQuote, target)
	})

	t.Run("replay retries accept without republishing", func(t *testing.T) {
		service, persistenceService, published := newEngagementTestService(t)
		const target = "https://owncast.example/federation/quoted"
		if err := persistenceService.AddToOutbox(target, []byte(`{"type":"Note"}`), "Note", false); err != nil {
			t.Fatal(err)
		}

		activity := quoteRequest("https://remote.example/inbox")
		if err := service.handleQuoteRequestInboxRequest(context.Background(), activity); err != nil {
			t.Fatalf("first delivery: %v", err)
		}
		if err := service.handleQuoteRequestInboxRequest(context.Background(), activity); err != nil {
			t.Fatalf("replayed delivery: %v", err)
		}
		if len(*published) != 1 {
			t.Fatalf("replayed quote published %d events, want 1", len(*published))
		}

		var outboxItems int
		if err := persistenceService.Datastore().DB.QueryRow("SELECT COUNT(*) FROM ap_outbox").Scan(&outboxItems); err != nil {
			t.Fatalf("count outbox items: %v", err)
		}
		if outboxItems != 3 {
			t.Fatalf("outbox has %d items, want target plus two authorization stamps", outboxItems)
		}
	})

	t.Run("rejected", func(t *testing.T) {
		service, _, published := newEngagementTestService(t)
		if err := service.handleQuoteRequestInboxRequest(context.Background(), quoteRequest("https://remote.example/inbox")); err != nil {
			t.Fatal(err)
		}
		if len(*published) != 0 {
			t.Fatalf("rejected quote published %d events, want 0", len(*published))
		}
	})

	t.Run("accept failed", func(t *testing.T) {
		service, persistenceService, published := newEngagementTestService(t)
		const target = "https://owncast.example/federation/quoted"
		if err := persistenceService.AddToOutbox(target, []byte(`{"type":"Note"}`), "Note", false); err != nil {
			t.Fatal(err)
		}

		if err := service.handleQuoteRequestInboxRequest(context.Background(), quoteRequest("http://remote.example/inbox")); err == nil {
			t.Fatal("quote with insecure inbox returned nil error")
		}
		if len(*published) != 0 {
			t.Fatalf("failed quote published %d events, want 0", len(*published))
		}
	})
}

func TestConcurrentQuoteDeliveriesPublishOnce(t *testing.T) {
	_ = config.GetReleaseString()
	service, persistenceService, published := newEngagementTestService(t)
	persistenceService.Datastore().DB.SetMaxOpenConns(1)
	const target = "https://owncast.example/federation/quoted"
	if err := persistenceService.AddToOutbox(target, []byte(`{"type":"Note"}`), "Note", false); err != nil {
		t.Fatal(err)
	}

	const deliveries = 8
	activity := quoteRequest("https://remote.example/inbox")
	start := make(chan struct{})
	errs := make(chan error, deliveries)
	var wg sync.WaitGroup
	wg.Add(deliveries)
	for range deliveries {
		go func() {
			defer wg.Done()
			<-start
			if err := service.handleQuoteRequestInboxRequest(context.Background(), activity); err != nil {
				errs <- err
			}
		}()
	}
	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		t.Fatalf("concurrent delivery: %v", err)
	}
	if len(*published) != 1 {
		t.Fatalf("concurrent quote deliveries published %d events, want 1", len(*published))
	}
}

func TestDuplicateApprovedDirectoryFollowResendsAccept(t *testing.T) {
	service, persistenceService, _ := newEngagementTestService(t)
	person := makeFakePerson()
	follow := streams.NewActivityStreamsFollow()
	id := streams.NewJSONLDIdProperty()
	id.Set(mustParseURL("https://freedom.eagle/follow/directory"))
	follow.SetJSONLDId(id)
	follow.SetActivityStreamsActor(actorProperty(person))
	follow.SetActivityStreamsObject(objectProperty("https://owncast.example/federation/user/streamer"))
	follow.GetUnknownProperties()[config.APOwncastNamespaceDirectory] = true

	if err := service.handleFollowInboxRequest(context.Background(), follow); err != nil {
		t.Fatalf("store pending directory follow: %v", err)
	}
	actorIRI := person.GetJSONLDId().GetIRI().String()
	if err := service.followers.ApprovePreviousRequest(actorIRI); err != nil {
		t.Fatalf("approve directory follow: %v", err)
	}
	if err := service.handleFollowInboxRequest(context.Background(), follow); err != nil {
		t.Fatalf("handle first duplicate follow: %v", err)
	}
	if err := service.handleFollowInboxRequest(context.Background(), follow); err != nil {
		t.Fatalf("handle second duplicate follow: %v", err)
	}

	var followerCount, deliveryCount, revision int
	if err := persistenceService.Datastore().DB.QueryRow(`SELECT count(*) FROM ap_followers`).Scan(&followerCount); err != nil {
		t.Fatal(err)
	}
	if err := persistenceService.Datastore().DB.QueryRow(`SELECT count(*), max(revision) FROM ap_delivery_queue`).Scan(&deliveryCount, &revision); err != nil {
		t.Fatal(err)
	}
	if followerCount != 1 || deliveryCount != 1 || revision < 1 {
		t.Fatalf("followers=%d deliveries=%d revision=%d, want one follower and one refreshed Accept", followerCount, deliveryCount, revision)
	}
}
