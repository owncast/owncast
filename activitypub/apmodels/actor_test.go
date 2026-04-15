package apmodels

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/persistence/configrepository"
)

func makeFakeService() vocab.ActivityStreamsService {
	iri, _ := url.Parse("https://fake.fediverse.server/user/mrfoo")
	name := "Mr Foo"
	username := "foodawg"
	inbox, _ := url.Parse("https://fake.fediverse.server/user/mrfoo/inbox")
	userAvatarURL, _ := url.Parse("https://fake.fediverse.server/user/mrfoo/avatar.png")

	service := streams.NewActivityStreamsService()

	id := streams.NewJSONLDIdProperty()
	id.Set(iri)
	service.SetJSONLDId(id)

	nameProperty := streams.NewActivityStreamsNameProperty()
	nameProperty.AppendXMLSchemaString(name)
	service.SetActivityStreamsName(nameProperty)

	preferredUsernameProperty := streams.NewActivityStreamsPreferredUsernameProperty()
	preferredUsernameProperty.SetXMLSchemaString(username)
	service.SetActivityStreamsPreferredUsername(preferredUsernameProperty)

	inboxProp := streams.NewActivityStreamsInboxProperty()
	inboxProp.SetIRI(inbox)
	service.SetActivityStreamsInbox(inboxProp)

	image := streams.NewActivityStreamsImage()
	imgProp := streams.NewActivityStreamsUrlProperty()
	imgProp.AppendIRI(userAvatarURL)
	image.SetActivityStreamsUrl(imgProp)
	icon := streams.NewActivityStreamsIconProperty()
	icon.AppendActivityStreamsImage(image)
	service.SetActivityStreamsIcon(icon)

	publicKeyProperty := streams.NewW3IDSecurityV1PublicKeyProperty()
	publicKeyType := streams.NewW3IDSecurityV1PublicKey()
	publicKeyProperty.AppendW3IDSecurityV1PublicKey(publicKeyType)
	service.SetW3IDSecurityV1PublicKey(publicKeyProperty)

	return service
}

func TestMain(m *testing.M) {
	dbFile, err := ioutil.TempFile(os.TempDir(), "owncast-test-db.db")
	if err != nil {
		panic(err)
	}
	data.SetupPersistence(dbFile.Name())

	configRepository := configrepository.Get()

	configRepository.SetServerURL("https://my.cool.site.biz")

	m.Run()
}

func TestMakeActorPropertyWithID(t *testing.T) {
	iri, _ := url.Parse("https://fake.fediverse.server/user/mrfoo")
	actor := MakeActorPropertyWithID(iri)

	if actor.Begin().GetIRI() != iri {
		t.Errorf("actor.IRI = %v, want %v", actor.Begin().GetIRI(), iri)
	}
}

func TestGetFullUsernameFromPerson(t *testing.T) {
	expected := "foodawg@fake.fediverse.server"
	person := makeFakeService()
	username := GetFullUsernameFromExternalEntity(person)

	if username != expected {
		t.Errorf("actor.Username = %v, want %v", username, expected)
	}
}

func TestAddMetadataLinkToProfile(t *testing.T) {
	person := makeFakeService()
	addMetadataLinkToProfile(person, "my site", "https://my.cool.site.biz")
	attchment := person.GetActivityStreamsAttachment().At(0)

	nameValue := attchment.GetActivityStreamsObject().GetActivityStreamsName().At(0).GetXMLSchemaString()
	expected := "my site"
	if nameValue != expected {
		t.Errorf("attachment name = %v, want %v", nameValue, expected)
	}

	propertyValue := attchment.GetActivityStreamsObject().GetUnknownProperties()["value"]
	expected = `<a href="https://my.cool.site.biz" rel="me nofollow noopener noreferrer" target="_blank">https://my.cool.site.biz</a>`
	if propertyValue != expected {
		t.Errorf("attachment value = %v, want %v", propertyValue, expected)
	}
}

func TestMakeServiceForAccount(t *testing.T) {
	person := MakeServiceForAccount("accountname")
	expectedIRI := "https://my.cool.site.biz/federation/user/accountname"
	if person.GetJSONLDId().Get().String() != expectedIRI {
		t.Errorf("actor.IRI = %v, want %v", person.GetJSONLDId().Get().String(), expectedIRI)
	}

	if person.GetActivityStreamsPreferredUsername().GetXMLSchemaString() != "accountname" {
		t.Errorf("actor.PreferredUsername = %v, want %v", person.GetActivityStreamsPreferredUsername().GetXMLSchemaString(), expectedIRI)
	}

	expectedInbox := "https://my.cool.site.biz/federation/user/accountname/inbox"
	if person.GetActivityStreamsInbox().GetIRI().String() != expectedInbox {
		t.Errorf("actor.Inbox = %v, want %v", person.GetActivityStreamsInbox().GetIRI().String(), expectedInbox)
	}

	expectedOutbox := "https://my.cool.site.biz/federation/user/accountname/outbox"
	if person.GetActivityStreamsOutbox().GetIRI().String() != expectedOutbox {
		t.Errorf("actor.Outbox = %v, want %v", person.GetActivityStreamsOutbox().GetIRI().String(), expectedOutbox)
	}

	expectedFollowers := "https://my.cool.site.biz/federation/user/accountname/followers"
	if person.GetActivityStreamsFollowers().GetIRI().String() != expectedFollowers {
		t.Errorf("actor.Followers = %v, want %v", person.GetActivityStreamsFollowers().GetIRI().String(), expectedFollowers)
	}

	expectedName := "New Owncast Server"
	if person.GetActivityStreamsName().Begin().GetXMLSchemaString() != expectedName {
		t.Errorf("actor.Name = %v, want %v", person.GetActivityStreamsName().Begin().GetXMLSchemaString(), expectedName)
	}

	expectedAvatar := "https://my.cool.site.biz/logo/external"
	u, err := url.Parse(person.GetActivityStreamsIcon().At(0).GetActivityStreamsImage().GetActivityStreamsUrl().Begin().GetIRI().String())
	if err != nil {
		t.Error(err)
	}
	u.RawQuery = ""
	if u.String() != expectedAvatar {
		t.Errorf("actor.Avatar = %v, want %v", person.GetActivityStreamsIcon().At(0).GetActivityStreamsImage().GetActivityStreamsUrl().Begin().GetIRI().String(), expectedAvatar)
	}

	expectedSummary := "This is a new live video streaming server powered by Owncast."
	if person.GetActivityStreamsSummary().At(0).GetXMLSchemaString() != expectedSummary {
		t.Errorf("actor.Summary = %v, want %v", person.GetActivityStreamsSummary().At(0).GetXMLSchemaString(), expectedSummary)
	}

	if person.GetActivityStreamsUrl().At(0).GetIRI().String() != expectedIRI {
		t.Errorf("actor.URL = %v, want %v", person.GetActivityStreamsUrl().At(0).GetIRI().String(), expectedIRI)
	}
}

func TestMakeServiceForAccountWithIDNServerURL(t *testing.T) {
	configRepository := configrepository.Get()
	configRepository.SetServerURL("https://live.retrospection.みんな")
	t.Cleanup(func() {
		configRepository.SetServerURL("https://my.cool.site.biz")
	})

	person := MakeServiceForAccount("retrots3m")
	payload, err := Serialize(person)
	if err != nil {
		t.Fatal(err)
	}

	var actor map[string]interface{}
	if err := json.Unmarshal(payload, &actor); err != nil {
		t.Fatal(err)
	}

	expectedActorURL := "https://live.retrospection.xn--q9jyb4c/federation/user/retrots3m"
	if actor["id"] != expectedActorURL {
		t.Errorf("actor id = %v, want %v", actor["id"], expectedActorURL)
	}
	if actor["url"] != expectedActorURL {
		t.Errorf("actor url = %v, want %v", actor["url"], expectedActorURL)
	}
	if actor["inbox"] != expectedActorURL+"/inbox" {
		t.Errorf("actor inbox = %v, want %v", actor["inbox"], expectedActorURL+"/inbox")
	}
	if actor["outbox"] != expectedActorURL+"/outbox" {
		t.Errorf("actor outbox = %v, want %v", actor["outbox"], expectedActorURL+"/outbox")
	}
	if actor["followers"] != expectedActorURL+"/followers" {
		t.Errorf("actor followers = %v, want %v", actor["followers"], expectedActorURL+"/followers")
	}

	publicKey := actor["publicKey"].(map[string]interface{})
	if publicKey["id"] != expectedActorURL+"#main-key" {
		t.Errorf("public key id = %v, want %v", publicKey["id"], expectedActorURL+"#main-key")
	}
	if publicKey["owner"] != expectedActorURL {
		t.Errorf("public key owner = %v, want %v", publicKey["owner"], expectedActorURL)
	}

	attachments := actor["attachment"].([]interface{})
	streamAttachment := attachments[0].(map[string]interface{})
	expectedDisplayValue := `<a href="https://live.retrospection.みんな" rel="me nofollow noopener noreferrer" target="_blank">https://live.retrospection.みんな</a>`
	if streamAttachment["value"] != expectedDisplayValue {
		t.Errorf("stream attachment = %v, want %v", streamAttachment["value"], expectedDisplayValue)
	}
}

// Tests for nil-safe accessor methods

func TestActorIriStringWithNilValue(t *testing.T) {
	actor := ActivityPubActor{}
	result := actor.ActorIriString()
	if result != "" {
		t.Errorf("ActorIriString() with nil ActorIri = %v, want empty string", result)
	}
}

func TestActorIriStringWithValue(t *testing.T) {
	iri, _ := url.Parse("https://example.com/user/test")
	actor := ActivityPubActor{ActorIri: iri}
	result := actor.ActorIriString()
	if result != "https://example.com/user/test" {
		t.Errorf("ActorIriString() = %v, want %v", result, "https://example.com/user/test")
	}
}

func TestInboxStringWithNilValue(t *testing.T) {
	actor := ActivityPubActor{}
	result := actor.InboxString()
	if result != "" {
		t.Errorf("InboxString() with nil Inbox = %v, want empty string", result)
	}
}

func TestInboxStringWithValue(t *testing.T) {
	inbox, _ := url.Parse("https://example.com/user/test/inbox")
	actor := ActivityPubActor{Inbox: inbox}
	result := actor.InboxString()
	if result != "https://example.com/user/test/inbox" {
		t.Errorf("InboxString() = %v, want %v", result, "https://example.com/user/test/inbox")
	}
}

func TestImageStringWithNilValue(t *testing.T) {
	actor := ActivityPubActor{}
	result := actor.ImageString()
	if result != "" {
		t.Errorf("ImageString() with nil Image = %v, want empty string", result)
	}
}

func TestImageStringWithValue(t *testing.T) {
	image, _ := url.Parse("https://example.com/avatar.png")
	actor := ActivityPubActor{Image: image}
	result := actor.ImageString()
	if result != "https://example.com/avatar.png" {
		t.Errorf("ImageString() = %v, want %v", result, "https://example.com/avatar.png")
	}
}

func TestFollowRequestIriStringWithNilValue(t *testing.T) {
	actor := ActivityPubActor{}
	result := actor.FollowRequestIriString()
	if result != "" {
		t.Errorf("FollowRequestIriString() with nil FollowRequestIri = %v, want empty string", result)
	}
}

func TestFollowRequestIriStringWithValue(t *testing.T) {
	followIri, _ := url.Parse("https://example.com/follow/123")
	actor := ActivityPubActor{FollowRequestIri: followIri}
	result := actor.FollowRequestIriString()
	if result != "https://example.com/follow/123" {
		t.Errorf("FollowRequestIriString() = %v, want %v", result, "https://example.com/follow/123")
	}
}

func TestActorIriHostnameWithNilValue(t *testing.T) {
	actor := ActivityPubActor{}
	result := actor.ActorIriHostname()
	if result != "" {
		t.Errorf("ActorIriHostname() with nil ActorIri = %v, want empty string", result)
	}
}

func TestActorIriHostnameWithValue(t *testing.T) {
	iri, _ := url.Parse("https://example.com/user/test")
	actor := ActivityPubActor{ActorIri: iri}
	result := actor.ActorIriHostname()
	if result != "example.com" {
		t.Errorf("ActorIriHostname() = %v, want %v", result, "example.com")
	}
}

// Tests for Validate() and IsValid() methods

func TestValidateWithAllNilFields(t *testing.T) {
	actor := ActivityPubActor{}
	err := actor.Validate()
	if err == nil {
		t.Error("Validate() with all nil fields should return error")
	}
	if !errors.Is(err, ErrActorMissingRequiredField) {
		t.Errorf("Validate() error = %v, want ErrActorMissingRequiredField", err)
	}
}

func TestValidateWithNilActorIri(t *testing.T) {
	inbox, _ := url.Parse("https://example.com/inbox")
	actor := ActivityPubActor{Inbox: inbox}
	err := actor.Validate()
	if err == nil {
		t.Error("Validate() with nil ActorIri should return error")
	}
	if !errors.Is(err, ErrActorMissingRequiredField) {
		t.Errorf("Validate() error = %v, want ErrActorMissingRequiredField", err)
	}
}

func TestValidateWithNilInbox(t *testing.T) {
	iri, _ := url.Parse("https://example.com/user/test")
	actor := ActivityPubActor{ActorIri: iri}
	err := actor.Validate()
	if err == nil {
		t.Error("Validate() with nil Inbox should return error")
	}
	if !errors.Is(err, ErrActorMissingRequiredField) {
		t.Errorf("Validate() error = %v, want ErrActorMissingRequiredField", err)
	}
}

func TestValidateWithRequiredFields(t *testing.T) {
	iri, _ := url.Parse("https://example.com/user/test")
	inbox, _ := url.Parse("https://example.com/inbox")
	actor := ActivityPubActor{ActorIri: iri, Inbox: inbox}
	err := actor.Validate()
	if err != nil {
		t.Errorf("Validate() with required fields should not return error, got %v", err)
	}
}

func TestIsValidWithInvalidActor(t *testing.T) {
	actor := ActivityPubActor{}
	if actor.IsValid() {
		t.Error("IsValid() with invalid actor should return false")
	}
}

func TestIsValidWithValidActor(t *testing.T) {
	iri, _ := url.Parse("https://example.com/user/test")
	inbox, _ := url.Parse("https://example.com/inbox")
	actor := ActivityPubActor{ActorIri: iri, Inbox: inbox}
	if !actor.IsValid() {
		t.Error("IsValid() with valid actor should return true")
	}
}

// Tests for NewActivityPubActor constructor

func TestNewActivityPubActorWithNilActorIri(t *testing.T) {
	inbox, _ := url.Parse("https://example.com/inbox")
	_, err := NewActivityPubActor(nil, inbox)
	if err == nil {
		t.Error("NewActivityPubActor with nil actorIri should return error")
	}
	if !errors.Is(err, ErrActorMissingRequiredField) {
		t.Errorf("NewActivityPubActor error = %v, want ErrActorMissingRequiredField", err)
	}
}

func TestNewActivityPubActorWithNilInbox(t *testing.T) {
	iri, _ := url.Parse("https://example.com/user/test")
	_, err := NewActivityPubActor(iri, nil)
	if err == nil {
		t.Error("NewActivityPubActor with nil inbox should return error")
	}
	if !errors.Is(err, ErrActorMissingRequiredField) {
		t.Errorf("NewActivityPubActor error = %v, want ErrActorMissingRequiredField", err)
	}
}

func TestNewActivityPubActorWithBothNil(t *testing.T) {
	_, err := NewActivityPubActor(nil, nil)
	if err == nil {
		t.Error("NewActivityPubActor with both nil should return error")
	}
	if !errors.Is(err, ErrActorMissingRequiredField) {
		t.Errorf("NewActivityPubActor error = %v, want ErrActorMissingRequiredField", err)
	}
}

func TestNewActivityPubActorWithValidArgs(t *testing.T) {
	iri, _ := url.Parse("https://example.com/user/test")
	inbox, _ := url.Parse("https://example.com/inbox")
	actor, err := NewActivityPubActor(iri, inbox)
	if err != nil {
		t.Errorf("NewActivityPubActor with valid args should not return error, got %v", err)
	}
	if actor == nil {
		t.Error("NewActivityPubActor with valid args should return non-nil actor")
	}
	if actor.ActorIri != iri {
		t.Errorf("actor.ActorIri = %v, want %v", actor.ActorIri, iri)
	}
	if actor.Inbox != inbox {
		t.Errorf("actor.Inbox = %v, want %v", actor.Inbox, inbox)
	}
}

// Tests for NewActivityPubActorFromEntity with invalid entities

func makeFakeServiceWithoutUsername() vocab.ActivityStreamsService {
	iri, _ := url.Parse("https://fake.fediverse.server/user/mrfoo")
	inbox, _ := url.Parse("https://fake.fediverse.server/user/mrfoo/inbox")

	service := streams.NewActivityStreamsService()

	id := streams.NewJSONLDIdProperty()
	id.Set(iri)
	service.SetJSONLDId(id)

	inboxProp := streams.NewActivityStreamsInboxProperty()
	inboxProp.SetIRI(inbox)
	service.SetActivityStreamsInbox(inboxProp)

	publicKeyProperty := streams.NewW3IDSecurityV1PublicKeyProperty()
	service.SetW3IDSecurityV1PublicKey(publicKeyProperty)

	return service
}

func makeFakeServiceWithoutPublicKey() vocab.ActivityStreamsService {
	iri, _ := url.Parse("https://fake.fediverse.server/user/mrfoo")
	inbox, _ := url.Parse("https://fake.fediverse.server/user/mrfoo/inbox")
	username := "foodawg"

	service := streams.NewActivityStreamsService()

	id := streams.NewJSONLDIdProperty()
	id.Set(iri)
	service.SetJSONLDId(id)

	preferredUsernameProperty := streams.NewActivityStreamsPreferredUsernameProperty()
	preferredUsernameProperty.SetXMLSchemaString(username)
	service.SetActivityStreamsPreferredUsername(preferredUsernameProperty)

	inboxProp := streams.NewActivityStreamsInboxProperty()
	inboxProp.SetIRI(inbox)
	service.SetActivityStreamsInbox(inboxProp)

	return service
}

func makeFakeServiceWithEmptyPublicKey() vocab.ActivityStreamsService {
	iri, _ := url.Parse("https://fake.fediverse.server/user/mrfoo")
	inbox, _ := url.Parse("https://fake.fediverse.server/user/mrfoo/inbox")
	username := "foodawg"

	service := streams.NewActivityStreamsService()

	id := streams.NewJSONLDIdProperty()
	id.Set(iri)
	service.SetJSONLDId(id)

	preferredUsernameProperty := streams.NewActivityStreamsPreferredUsernameProperty()
	preferredUsernameProperty.SetXMLSchemaString(username)
	service.SetActivityStreamsPreferredUsername(preferredUsernameProperty)

	inboxProp := streams.NewActivityStreamsInboxProperty()
	inboxProp.SetIRI(inbox)
	service.SetActivityStreamsInbox(inboxProp)

	// Set an empty public key property (Len() == 0)
	publicKeyProperty := streams.NewW3IDSecurityV1PublicKeyProperty()
	service.SetW3IDSecurityV1PublicKey(publicKeyProperty)

	return service
}

func makeFakeServiceWithoutId() vocab.ActivityStreamsService {
	inbox, _ := url.Parse("https://fake.fediverse.server/user/mrfoo/inbox")
	username := "foodawg"

	service := streams.NewActivityStreamsService()

	preferredUsernameProperty := streams.NewActivityStreamsPreferredUsernameProperty()
	preferredUsernameProperty.SetXMLSchemaString(username)
	service.SetActivityStreamsPreferredUsername(preferredUsernameProperty)

	inboxProp := streams.NewActivityStreamsInboxProperty()
	inboxProp.SetIRI(inbox)
	service.SetActivityStreamsInbox(inboxProp)

	publicKeyProperty := streams.NewW3IDSecurityV1PublicKeyProperty()
	service.SetW3IDSecurityV1PublicKey(publicKeyProperty)

	return service
}

func makeFakeServiceWithoutInbox() vocab.ActivityStreamsService {
	iri, _ := url.Parse("https://fake.fediverse.server/user/mrfoo")
	username := "foodawg"

	service := streams.NewActivityStreamsService()

	id := streams.NewJSONLDIdProperty()
	id.Set(iri)
	service.SetJSONLDId(id)

	preferredUsernameProperty := streams.NewActivityStreamsPreferredUsernameProperty()
	preferredUsernameProperty.SetXMLSchemaString(username)
	service.SetActivityStreamsPreferredUsername(preferredUsernameProperty)

	publicKeyProperty := streams.NewW3IDSecurityV1PublicKeyProperty()
	service.SetW3IDSecurityV1PublicKey(publicKeyProperty)

	return service
}

func TestNewActivityPubActorFromEntityWithoutUsername(t *testing.T) {
	service := makeFakeServiceWithoutUsername()
	_, err := NewActivityPubActorFromEntity(service)
	if err == nil {
		t.Error("NewActivityPubActorFromEntity without username should return error")
	}
	if !errors.Is(err, ErrActorMissingRequiredField) {
		t.Errorf("NewActivityPubActorFromEntity error = %v, want ErrActorMissingRequiredField", err)
	}
}

func TestNewActivityPubActorFromEntityWithoutPublicKey(t *testing.T) {
	service := makeFakeServiceWithoutPublicKey()
	_, err := NewActivityPubActorFromEntity(service)
	if err == nil {
		t.Error("NewActivityPubActorFromEntity without public key should return error")
	}
	if !errors.Is(err, ErrActorMissingRequiredField) {
		t.Errorf("NewActivityPubActorFromEntity error = %v, want ErrActorMissingRequiredField", err)
	}
}

func TestNewActivityPubActorFromEntityWithEmptyPublicKey(t *testing.T) {
	service := makeFakeServiceWithEmptyPublicKey()
	_, err := NewActivityPubActorFromEntity(service)
	if err == nil {
		t.Error("NewActivityPubActorFromEntity with empty public key should return error")
	}
	if !errors.Is(err, ErrActorMissingRequiredField) {
		t.Errorf("NewActivityPubActorFromEntity error = %v, want ErrActorMissingRequiredField", err)
	}
}

func TestNewActivityPubActorFromEntityWithoutId(t *testing.T) {
	service := makeFakeServiceWithoutId()
	_, err := NewActivityPubActorFromEntity(service)
	if err == nil {
		t.Error("NewActivityPubActorFromEntity without ID should return error")
	}
	if !errors.Is(err, ErrActorMissingRequiredField) {
		t.Errorf("NewActivityPubActorFromEntity error = %v, want ErrActorMissingRequiredField", err)
	}
}

func TestNewActivityPubActorFromEntityWithoutInbox(t *testing.T) {
	service := makeFakeServiceWithoutInbox()
	_, err := NewActivityPubActorFromEntity(service)
	if err == nil {
		t.Error("NewActivityPubActorFromEntity without inbox should return error")
	}
	if !errors.Is(err, ErrActorMissingRequiredField) {
		t.Errorf("NewActivityPubActorFromEntity error = %v, want ErrActorMissingRequiredField", err)
	}
}

func TestNewActivityPubActorFromEntityWithValidEntity(t *testing.T) {
	service := makeFakeService()
	actor, err := NewActivityPubActorFromEntity(service)
	if err != nil {
		t.Errorf("NewActivityPubActorFromEntity with valid entity should not return error, got %v", err)
	}
	if actor == nil {
		t.Fatal("NewActivityPubActorFromEntity with valid entity should return non-nil actor")
	}

	// Verify required fields are non-nil
	if actor.ActorIri == nil {
		t.Error("actor.ActorIri should not be nil")
	}
	if actor.Inbox == nil {
		t.Error("actor.Inbox should not be nil")
	}

	// Verify extracted values match the fake service data
	expectedIri := "https://fake.fediverse.server/user/mrfoo"
	if actor.ActorIriString() != expectedIri {
		t.Errorf("actor.ActorIri = %v, want %v", actor.ActorIriString(), expectedIri)
	}

	expectedInbox := "https://fake.fediverse.server/user/mrfoo/inbox"
	if actor.InboxString() != expectedInbox {
		t.Errorf("actor.Inbox = %v, want %v", actor.InboxString(), expectedInbox)
	}

	expectedName := "Mr Foo"
	if actor.Name != expectedName {
		t.Errorf("actor.Name = %v, want %v", actor.Name, expectedName)
	}

	expectedUsername := "foodawg"
	if actor.Username != expectedUsername {
		t.Errorf("actor.Username = %v, want %v", actor.Username, expectedUsername)
	}

	expectedImage := "https://fake.fediverse.server/user/mrfoo/avatar.png"
	if actor.ImageString() != expectedImage {
		t.Errorf("actor.Image = %v, want %v", actor.ImageString(), expectedImage)
	}
}

// Test that safe accessors don't panic on zero-value struct

func TestSafeAccessorsOnZeroValueStruct(t *testing.T) {
	var actor ActivityPubActor

	// These should not panic
	_ = actor.ActorIriString()
	_ = actor.InboxString()
	_ = actor.ImageString()
	_ = actor.FollowRequestIriString()
	_ = actor.ActorIriHostname()
	_ = actor.Validate()
	_ = actor.IsValid()
}

// Test that safe accessors work correctly with optional nil fields on otherwise valid actor

func TestSafeAccessorsWithOptionalNilFields(t *testing.T) {
	iri, _ := url.Parse("https://example.com/user/test")
	inbox, _ := url.Parse("https://example.com/inbox")
	actor := ActivityPubActor{
		ActorIri: iri,
		Inbox:    inbox,
		// Image and FollowRequestIri are intentionally nil
	}

	// Required fields should return values
	if actor.ActorIriString() != "https://example.com/user/test" {
		t.Errorf("ActorIriString() = %v, want non-empty", actor.ActorIriString())
	}
	if actor.InboxString() != "https://example.com/inbox" {
		t.Errorf("InboxString() = %v, want non-empty", actor.InboxString())
	}

	// Optional nil fields should return empty strings without panicking
	if actor.ImageString() != "" {
		t.Errorf("ImageString() = %v, want empty string", actor.ImageString())
	}
	if actor.FollowRequestIriString() != "" {
		t.Errorf("FollowRequestIriString() = %v, want empty string", actor.FollowRequestIriString())
	}

	// Actor should still be valid (only ActorIri and Inbox are required)
	if !actor.IsValid() {
		t.Error("Actor with ActorIri and Inbox should be valid even with nil optional fields")
	}
}
