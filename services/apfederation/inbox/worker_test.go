package inbox

import (
	"net/url"
	"testing"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/services/apfederation/apmodels"
	"github.com/owncast/owncast/storage/configrepository"
	"github.com/owncast/owncast/storage/data"
	"github.com/owncast/owncast/storage/federationrepository"
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

func TestMain(m *testing.M) {
	ds, err := data.NewStore(":memory:")
	if err != nil {
		panic(err)
	}

	_ = federationrepository.New(ds)
	configRepository := configrepository.New(ds)
	configRepository.PopulateDefaults()
	configRepository.SetServerURL("https://my.cool.site.biz")

	m.Run()
}

func TestBlockedDomains(t *testing.T) {
	cr := configrepository.Get()

	person := makeFakePerson()

	cr.SetBlockedFederatedDomains([]string{"freedom.eagle", "guns.life"})

	if len(cr.GetBlockedFederatedDomains()) != 2 {
		t.Error("Blocked federated domains is not set correctly")
	}

	for _, domain := range cr.GetBlockedFederatedDomains() {
		if domain == person.GetJSONLDId().GetIRI().Host {
			return
		}
	}

	t.Error("Failed to catch blocked domain")
}

func TestBlockedActors(t *testing.T) {
	federationRepository := federationrepository.Get()
	person := makeFakePerson()
	fakeRequest := streams.NewActivityStreamsFollow()
	ib := Get()

	federationRepository.AddFollow(apmodels.ActivityPubActor{
		ActorIri:         person.GetJSONLDId().GetIRI(),
		Inbox:            person.GetJSONLDId().GetIRI(),
		FollowRequestIri: person.GetJSONLDId().GetIRI(),
		RequestObject:    fakeRequest,
	}, false)
	federationRepository.BlockOrRejectFollower(person.GetJSONLDId().GetIRI().String())

	blocked, err := ib.isBlockedActor(person.GetJSONLDId().GetIRI())
	if err != nil {
		t.Error(err)
		return
	}

	if !blocked {
		t.Error("Failed to block actor")
	}

	failedBlockIRI, _ := url.Parse("https://freedom.eagle/user/mrbar")
	failedBlock, err := ib.isBlockedActor(failedBlockIRI)

	if failedBlock {
		t.Error("Invalid blocking of unblocked actor IRI")
	}
}
