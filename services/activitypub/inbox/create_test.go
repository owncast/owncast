package inbox

import (
	"context"
	"net/url"
	"sync"
	"testing"
	"time"

	"code.superseriousbusiness.org/activity/streams"
	"code.superseriousbusiness.org/activity/streams/vocab"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/events"
	"github.com/owncast/owncast/services/activitypub/persistence"
	"github.com/owncast/owncast/services/activitypub/resolvers"
	"github.com/owncast/owncast/services/dispatcher"
)

type createTestFixture struct {
	service    *Service
	activity   vocab.ActivityStreamsCreate
	note       vocab.ActivityStreamsNote
	actorIRI   *url.URL
	localActor *url.URL
	events     []dispatcher.Event
	eventsMu   sync.Mutex
}

func newCreateTestFixture(t *testing.T, suffix string) *createTestFixture {
	t.Helper()

	configRepository := configrepository.New(testDatastore)
	builder := apmodels.New(apmodels.Deps{ConfigRepository: configRepository})
	eventDispatcher := dispatcher.New()
	fixture := &createTestFixture{
		actorIRI:   mustCreateTestURL(t, "https://remote.example/users/alice"),
		localActor: builder.MakeLocalIRIForAccount(configRepository.GetFederationUsername()),
	}
	eventDispatcher.AddListener(func(_ context.Context, event dispatcher.Event) {
		fixture.eventsMu.Lock()
		defer fixture.eventsMu.Unlock()
		fixture.events = append(fixture.events, event)
	})
	fixture.service = New(Deps{
		Persistence:      persistence.New(testDatastore, nil),
		ConfigRepository: configRepository,
		Builder:          builder,
		Resolver: resolvers.New(resolvers.Deps{
			ConfigRepository: configRepository,
			Builder:          builder,
		}),
		Events: eventDispatcher,
	})

	person := streams.NewActivityStreamsPerson()
	personID := streams.NewJSONLDIdProperty()
	personID.SetIRI(fixture.actorIRI)
	person.SetJSONLDId(personID)
	inbox := streams.NewActivityStreamsInboxProperty()
	inbox.SetIRI(mustCreateTestURL(t, "https://remote.example/users/alice/inbox"))
	person.SetActivityStreamsInbox(inbox)
	name := streams.NewActivityStreamsNameProperty()
	name.AppendXMLSchemaString("Alice")
	person.SetActivityStreamsName(name)
	username := streams.NewActivityStreamsPreferredUsernameProperty()
	username.SetXMLSchemaString("alice")
	person.SetActivityStreamsPreferredUsername(username)
	publicKey := streams.NewW3IDSecurityV1PublicKeyProperty()
	publicKey.AppendW3IDSecurityV1PublicKey(streams.NewW3IDSecurityV1PublicKey())
	person.SetW3IDSecurityV1PublicKey(publicKey)

	fixture.activity = streams.NewActivityStreamsCreate()
	createID := streams.NewJSONLDIdProperty()
	createID.SetIRI(mustCreateTestURL(t, "https://remote.example/activities/"+suffix))
	fixture.activity.SetJSONLDId(createID)
	actor := streams.NewActivityStreamsActorProperty()
	actor.AppendActivityStreamsPerson(person)
	fixture.activity.SetActivityStreamsActor(actor)

	fixture.note = streams.NewActivityStreamsNote()
	noteID := streams.NewJSONLDIdProperty()
	noteID.SetIRI(mustCreateTestURL(t, "https://remote.example/posts/"+suffix))
	fixture.note.SetJSONLDId(noteID)
	attributedTo := streams.NewActivityStreamsAttributedToProperty()
	attributedTo.AppendIRI(fixture.actorIRI)
	fixture.note.SetActivityStreamsAttributedTo(attributedTo)
	object := streams.NewActivityStreamsObjectProperty()
	object.AppendActivityStreamsNote(fixture.note)
	fixture.activity.SetActivityStreamsObject(object)

	return fixture
}

func TestHandleCreateMentionExtractsPost(t *testing.T) {
	fixture := newCreateTestFixture(t, "mention")
	cc := streams.NewActivityStreamsCcProperty()
	cc.AppendIRI(fixture.localActor)
	fixture.activity.SetActivityStreamsCc(cc)

	content := streams.NewActivityStreamsContentProperty()
	content.AppendXMLSchemaString("<p>Hello <strong>world</strong></p>")
	fixture.note.SetActivityStreamsContent(content)
	postURL := streams.NewActivityStreamsUrlProperty()
	postURL.AppendIRI(mustCreateTestURL(t, "https://remote.example/@alice/mention"))
	fixture.note.SetActivityStreamsUrl(postURL)
	publishedAt := time.Date(2026, time.July, 10, 12, 34, 56, 0, time.FixedZone("offset", 2*60*60))
	published := streams.NewActivityStreamsPublishedProperty()
	published.Set(publishedAt)
	fixture.note.SetActivityStreamsPublished(published)

	attachments := streams.NewActivityStreamsAttachmentProperty()
	image := streams.NewActivityStreamsImage()
	imageURL := streams.NewActivityStreamsUrlProperty()
	imageURL.AppendIRI(mustCreateTestURL(t, "https://cdn.remote.example/image.jpg"))
	image.SetActivityStreamsUrl(imageURL)
	imageType := streams.NewActivityStreamsMediaTypeProperty()
	imageType.Set("image/jpeg")
	image.SetActivityStreamsMediaType(imageType)
	imageName := streams.NewActivityStreamsNameProperty()
	imageName.AppendXMLSchemaString("<b>Alt text</b>")
	image.SetActivityStreamsName(imageName)
	attachments.AppendActivityStreamsImage(image)
	document := streams.NewActivityStreamsDocument()
	documentURL := streams.NewActivityStreamsUrlProperty()
	documentURL.AppendIRI(mustCreateTestURL(t, "https://cdn.remote.example/caption.vtt"))
	document.SetActivityStreamsUrl(documentURL)
	documentType := streams.NewActivityStreamsMediaTypeProperty()
	documentType.Set("text/vtt")
	document.SetActivityStreamsMediaType(documentType)
	attachments.AppendActivityStreamsDocument(document)

	unsafeDocument := streams.NewActivityStreamsDocument()
	unsafeURL := streams.NewActivityStreamsUrlProperty()
	unsafeURL.AppendIRI(mustCreateTestURL(t, "file:///etc/passwd"))
	unsafeDocument.SetActivityStreamsUrl(unsafeURL)
	attachments.AppendActivityStreamsDocument(unsafeDocument)
	fixture.note.SetActivityStreamsAttachment(attachments)

	if err := fixture.service.handleCreateRequest(context.Background(), fixture.activity); err != nil {
		t.Fatalf("handle create mention: %v", err)
	}
	post := requireCreatePostEvent(t, fixture, models.FediverseMention)
	if post.Content != "<p>Hello <strong>world</strong></p>" || post.ContentText != "Hello world" {
		t.Fatalf("unexpected content: HTML %q text %q", post.Content, post.ContentText)
	}
	if post.URL != "https://remote.example/@alice/mention" {
		t.Fatalf("unexpected post URL %q", post.URL)
	}
	if post.PostedAt != "2026-07-10T10:34:56Z" {
		t.Fatalf("unexpected posted time %q", post.PostedAt)
	}
	if post.Actor.Handle != "alice@remote.example" || post.Actor.URL != fixture.actorIRI.String() {
		t.Fatalf("unexpected actor: %+v", post.Actor)
	}
	if len(post.Attachments) != 2 {
		t.Fatalf("expected safe embedded attachments only, got %+v", post.Attachments)
	}
	if got := post.Attachments[0]; got.URL != "https://cdn.remote.example/image.jpg" || got.MediaType != "image/jpeg" || got.Alt != "Alt text" {
		t.Fatalf("unexpected image attachment: %+v", got)
	}
	if got := post.Attachments[1]; got.URL != "https://cdn.remote.example/caption.vtt" || got.MediaType != "text/vtt" {
		t.Fatalf("unexpected document attachment: %+v", got)
	}
}

func TestHandleCreateLocalReplyTakesPrecedence(t *testing.T) {
	fixture := newCreateTestFixture(t, "reply")
	parentIRI := "https://my.cool.site.biz/federation/posts/local-parent"
	if err := fixture.service.persistence.AddToOutbox(parentIRI, []byte(`{"type":"Note"}`), "Note", false); err != nil {
		t.Fatalf("save local parent: %v", err)
	}
	inReplyTo := streams.NewActivityStreamsInReplyToProperty()
	inReplyTo.AppendIRI(mustCreateTestURL(t, parentIRI))
	fixture.note.SetActivityStreamsInReplyTo(inReplyTo)

	if err := fixture.service.handleCreateRequest(context.Background(), fixture.activity); err != nil {
		t.Fatalf("handle create reply: %v", err)
	}
	post := requireCreatePostEvent(t, fixture, models.FediverseReply)
	if post.InReplyTo != parentIRI {
		t.Fatalf("unexpected parent %q", post.InReplyTo)
	}
}

func TestHandleCreateUnknownParentIsMention(t *testing.T) {
	fixture := newCreateTestFixture(t, "unknown-parent")
	inReplyTo := streams.NewActivityStreamsInReplyToProperty()
	inReplyTo.AppendIRI(mustCreateTestURL(t, "https://elsewhere.example/posts/unknown"))
	fixture.note.SetActivityStreamsInReplyTo(inReplyTo)
	to := streams.NewActivityStreamsToProperty()
	to.AppendIRI(fixture.localActor)
	fixture.note.SetActivityStreamsTo(to)

	if err := fixture.service.handleCreateRequest(context.Background(), fixture.activity); err != nil {
		t.Fatalf("handle create mention: %v", err)
	}
	post := requireCreatePostEvent(t, fixture, models.FediverseMention)
	if post.InReplyTo != "https://elsewhere.example/posts/unknown" {
		t.Fatalf("unexpected parent %q", post.InReplyTo)
	}
}

func TestHandleCreateUnrelatedAndIRIObjectIgnored(t *testing.T) {
	fixture := newCreateTestFixture(t, "unrelated")
	actor := streams.NewActivityStreamsActorProperty()
	actor.AppendIRI(fixture.actorIRI)
	fixture.activity.SetActivityStreamsActor(actor)
	attributedTo := streams.NewActivityStreamsAttributedToProperty()
	attributedTo.AppendIRI(fixture.actorIRI)
	fixture.note.SetActivityStreamsAttributedTo(attributedTo)

	if err := fixture.service.handleCreateRequest(context.Background(), fixture.activity); err != nil {
		t.Fatalf("unrelated create should be ignored: %v", err)
	}
	if len(fixture.events) != 0 {
		t.Fatalf("unrelated create published %d events", len(fixture.events))
	}

	object := streams.NewActivityStreamsObjectProperty()
	object.AppendIRI(mustCreateTestURL(t, "https://remote.example/posts/iri-only"))
	fixture.activity.SetActivityStreamsObject(object)
	if err := fixture.service.handleCreateRequest(context.Background(), fixture.activity); err != nil {
		t.Fatalf("IRI-only object should be ignored: %v", err)
	}
	if len(fixture.events) != 0 {
		t.Fatalf("IRI-only create published %d events", len(fixture.events))
	}
	for _, eventType := range []models.EventType{models.FediverseMention, models.FediverseReply} {
		handled, err := fixture.service.persistence.HasPreviouslyHandledInboundActivity(
			"https://remote.example/activities/unrelated",
			fixture.actorIRI.String(),
			string(eventType),
		)
		if err != nil {
			t.Fatalf("check ignored activity persistence: %v", err)
		}
		if handled {
			t.Fatalf("ignored create was persisted as %q", eventType)
		}
	}
}

func TestHandleCreateRejectsActorAttributionMismatch(t *testing.T) {
	fixture := newCreateTestFixture(t, "mismatch")
	attributedTo := streams.NewActivityStreamsAttributedToProperty()
	attributedTo.AppendIRI(mustCreateTestURL(t, "https://remote.example/users/mallory"))
	fixture.note.SetActivityStreamsAttributedTo(attributedTo)
	to := streams.NewActivityStreamsToProperty()
	to.AppendIRI(fixture.localActor)
	fixture.note.SetActivityStreamsTo(to)

	if err := fixture.service.handleCreateRequest(context.Background(), fixture.activity); err == nil {
		t.Fatal("expected actor/attributedTo mismatch to be rejected")
	}
	if len(fixture.events) != 0 {
		t.Fatalf("mismatched create published %d events", len(fixture.events))
	}
}

func TestHandleCreateDuplicateSuppressed(t *testing.T) {
	fixture := newCreateTestFixture(t, "duplicate")
	to := streams.NewActivityStreamsToProperty()
	to.AppendIRI(fixture.localActor)
	fixture.note.SetActivityStreamsTo(to)

	if err := fixture.service.handleCreateRequest(context.Background(), fixture.activity); err != nil {
		t.Fatalf("first delivery: %v", err)
	}
	if err := fixture.service.handleCreateRequest(context.Background(), fixture.activity); err != nil {
		t.Fatalf("duplicate delivery: %v", err)
	}
	if len(fixture.events) != 1 {
		t.Fatalf("expected one event after duplicate delivery, got %d", len(fixture.events))
	}
}

func TestHandleCreateConcurrentDuplicateSuppressed(t *testing.T) {
	testDatastore.DB.SetMaxOpenConns(1)
	fixture := newCreateTestFixture(t, "concurrent-duplicate")
	to := streams.NewActivityStreamsToProperty()
	to.AppendIRI(fixture.localActor)
	fixture.note.SetActivityStreamsTo(to)

	const deliveries = 16
	start := make(chan struct{})
	errs := make(chan error, deliveries)
	var wg sync.WaitGroup
	wg.Add(deliveries)
	for range deliveries {
		go func() {
			defer wg.Done()
			<-start
			if err := fixture.service.handleCreateRequest(context.Background(), fixture.activity); err != nil {
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
	if len(fixture.events) != 1 {
		t.Fatalf("expected one event after concurrent deliveries, got %d", len(fixture.events))
	}
}

func TestHandleCreateReplayClassificationChangeSuppressed(t *testing.T) {
	fixture := newCreateTestFixture(t, "classification-change")
	parentIRI := "https://my.cool.site.biz/federation/posts/delayed-parent"
	inReplyTo := streams.NewActivityStreamsInReplyToProperty()
	inReplyTo.AppendIRI(mustCreateTestURL(t, parentIRI))
	fixture.note.SetActivityStreamsInReplyTo(inReplyTo)
	to := streams.NewActivityStreamsToProperty()
	to.AppendIRI(fixture.localActor)
	fixture.note.SetActivityStreamsTo(to)

	if err := fixture.service.handleCreateRequest(context.Background(), fixture.activity); err != nil {
		t.Fatalf("first delivery: %v", err)
	}
	requireCreatePostEvent(t, fixture, models.FediverseMention)

	if err := fixture.service.persistence.AddToOutbox(parentIRI, []byte(`{"type":"Note"}`), "Note", false); err != nil {
		t.Fatalf("save delayed parent: %v", err)
	}
	if err := fixture.service.handleCreateRequest(context.Background(), fixture.activity); err != nil {
		t.Fatalf("replayed delivery: %v", err)
	}
	if len(fixture.events) != 1 {
		t.Fatalf("classification-changing replay published an event; total = %d, want 1", len(fixture.events))
	}
}

func TestHandleCreatePublishedFallbackAndLanguage(t *testing.T) {
	fixture := newCreateTestFixture(t, "fallback")
	to := streams.NewActivityStreamsToProperty()
	to.AppendIRI(fixture.localActor)
	fixture.note.SetActivityStreamsTo(to)
	content := streams.NewActivityStreamsContentProperty()
	content.AppendRDFLangString(map[string]string{
		"fr": "<p>Bonjour</p>",
		"en": "<p>Hello</p>",
	})
	fixture.note.SetActivityStreamsContent(content)

	before := time.Now().UTC().Add(-time.Second)
	if err := fixture.service.handleCreateRequest(context.Background(), fixture.activity); err != nil {
		t.Fatalf("handle create fallback: %v", err)
	}
	after := time.Now().UTC().Add(time.Second)
	post := requireCreatePostEvent(t, fixture, models.FediverseMention)
	postedAt, err := time.Parse(time.RFC3339, post.PostedAt)
	if err != nil || postedAt.Before(before) || postedAt.After(after) {
		t.Fatalf("postedAt %q is not receipt-time fallback", post.PostedAt)
	}
	if post.Language != "en" || post.Content != "<p>Hello</p>" || post.ContentText != "Hello" {
		t.Fatalf("unexpected deterministic language content: %+v", post)
	}
}

func requireCreatePostEvent(t *testing.T, fixture *createTestFixture, eventType models.EventType) *events.FediverseInboundPostEvent {
	t.Helper()
	if len(fixture.events) != 1 {
		t.Fatalf("expected one event, got %d", len(fixture.events))
	}
	if fixture.events[0].Type != string(eventType) {
		t.Fatalf("expected event %q, got %q", eventType, fixture.events[0].Type)
	}
	post, ok := fixture.events[0].Payload.(*events.FediverseInboundPostEvent)
	if !ok {
		t.Fatalf("unexpected payload type %T", fixture.events[0].Payload)
	}
	return post
}

func mustCreateTestURL(t *testing.T, value string) *url.URL {
	t.Helper()
	parsed, err := url.Parse(value)
	if err != nil {
		t.Fatalf("parse URL %q: %v", value, err)
	}
	return parsed
}
