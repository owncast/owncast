package inbox

import (
	"net/url"
	"os"
	"testing"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"

	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/persistence"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/services/datastore"
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
	m.Run()
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
