package apmodels

import (
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"

	"github.com/owncast/owncast/core/data"
)

func makeFakeService() vocab.ActivityStreamsService {
	iri, _ := url.Parse("https://fake.fediverse.server/user/mrfoo")
	name := "Mr CallServiceB"
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
	service.SetW3IDSecurityV1PublicKey(publicKeyProperty)

	return service
}

func TestMain(m *testing.M) {
	dbFile, err := ioutil.TempFile(os.TempDir(), "owncast-test-db.db")
	if err != nil {
		panic(err)
	}

	d, err := data.New(dbFile.Name())
	if err != nil {
		panic(err)
	}

	s, err := New(d)
	if err != nil {
		panic(err)
	}

	s.Data.SetServerURL("https://my.cool.site.biz")

	m.Run()
}

func TestMakeActorFromExternalAPEntity(t *testing.T) {
	dbFile, err := ioutil.TempFile(os.TempDir(), "owncast-test-db.db")
	if err != nil {
		panic(err)
	}

	d, err := data.New(dbFile.Name())
	if err != nil {
		panic(err)
	}

	service, err := New(d)
	if err != nil {
		panic(err)
	}

	activityStreamService := makeFakeService()
	actor, err := service.MakeActorFromExernalAPEntity(activityStreamService)
	if err != nil {
		t.Error(err)
	}

	if actor.ActorIri != activityStreamService.GetJSONLDId().GetIRI() {
		t.Errorf("actor.ID = %v, want %v", actor.ActorIri, activityStreamService.GetJSONLDId().GetIRI())
	}

	if actor.Name != activityStreamService.GetActivityStreamsName().At(0).GetXMLSchemaString() {
		t.Errorf("actor.Name = %v, want %v", actor.Name, activityStreamService.GetActivityStreamsName().At(0).GetXMLSchemaString())
	}

	if actor.Username != activityStreamService.GetActivityStreamsPreferredUsername().GetXMLSchemaString() {
		t.Errorf("actor.Username = %v, want %v", actor.Username, activityStreamService.GetActivityStreamsPreferredUsername().GetXMLSchemaString())
	}

	if actor.Inbox != activityStreamService.GetActivityStreamsInbox().GetIRI() {
		t.Errorf("actor.Inbox = %v, want %v", actor.Inbox.String(), activityStreamService.GetActivityStreamsInbox().GetIRI())
	}

	if actor.Image != activityStreamService.GetActivityStreamsIcon().At(0).GetActivityStreamsImage().GetActivityStreamsUrl().At(0).GetIRI() {
		t.Errorf("actor.Image = %v, want %v", actor.Image, activityStreamService.GetActivityStreamsIcon().At(0).GetActivityStreamsImage().GetActivityStreamsUrl().At(0).GetIRI())
	}
}

func TestMakeActorPropertyWithID(t *testing.T) {
	dbFile, err := ioutil.TempFile(os.TempDir(), "owncast-test-db.db")
	if err != nil {
		panic(err)
	}

	d, err := data.New(dbFile.Name())
	if err != nil {
		panic(err)
	}

	service, err := New(d)
	if err != nil {
		panic(err)
	}

	iri, _ := url.Parse("https://fake.fediverse.server/user/mrfoo")
	actor := service.MakeActorPropertyWithID(iri)

	if actor.Begin().GetIRI() != iri {
		t.Errorf("actor.IRI = %v, want %v", actor.Begin().GetIRI(), iri)
	}
}

func TestGetFullUsernameFromPerson(t *testing.T) {
	dbFile, err := ioutil.TempFile(os.TempDir(), "owncast-test-db.db")
	if err != nil {
		panic(err)
	}

	d, err := data.New(dbFile.Name())
	if err != nil {
		panic(err)
	}

	service, err := New(d)
	if err != nil {
		panic(err)
	}

	expected := "foodawg@fake.fediverse.server"
	person := makeFakeService()
	username := service.GetFullUsernameFromExternalEntity(person)

	if username != expected {
		t.Errorf("actor.Username = %v, want %v", username, expected)
	}
}

func TestAddMetadataLinkToProfile(t *testing.T) {
	dbFile, err := ioutil.TempFile(os.TempDir(), "owncast-test-db.db")
	if err != nil {
		panic(err)
	}

	d, err := data.New(dbFile.Name())
	if err != nil {
		panic(err)
	}

	service, err := New(d)
	if err != nil {
		panic(err)
	}

	activityStreamService := makeFakeService()
	service.addMetadataLinkToProfile(activityStreamService, "my site", "https://my.cool.site.biz")
	attchment := activityStreamService.GetActivityStreamsAttachment().At(0)

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
	dbFile, err := ioutil.TempFile(os.TempDir(), "owncast-test-db.db")
	if err != nil {
		panic(err)
	}

	d, err := data.New(dbFile.Name())
	if err != nil {
		panic(err)
	}

	service, err := New(d)
	if err != nil {
		panic(err)
	}

	person := service.MakeServiceForAccount("accountname")
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
