package inbox

import (
	"net/url"
	"testing"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/core/data"
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
	data.SetupPersistence(":memory:")
	data.SetServerURL("https://my.cool.site.biz")
	persistence.Setup(data.GetDatastore())
	m.Run()
}

func TestBlockedDomains(t *testing.T) {
	person := makeFakePerson()

	data.SetBlockedFederatedDomains([]string{"freedom.eagle", "guns.life"})

	if len(data.GetBlockedFederatedDomains()) != 2 {
		t.Error("Blocked federated domains is not set correctly")
	}

	for _, domain := range data.GetBlockedFederatedDomains() {
		if domain == person.GetJSONLDId().GetIRI().Host {
			return
		}
	}

	t.Error("Failed to catch blocked domain")
}

func TestBlockedActors(t *testing.T) {
	person := makeFakePerson()
	persistence.AddFollow(apmodels.ActivityPubActor{
		ActorIri:         person.GetJSONLDId().GetIRI(),
		Inbox:            person.GetJSONLDId().GetIRI(),
		FollowRequestIri: person.GetJSONLDId().GetIRI(),
	}, false)
	persistence.BlockOrRejectFollower(person.GetJSONLDId().GetIRI().String())

	blocked, err := isBlockedActor(person.GetJSONLDId().GetIRI())
	if err != nil {
		t.Error(err)
		return
	}

	if !blocked {
		t.Error("Failed to block actor")
	}

	failedBlockIRI, _ := url.Parse("https://freedom.eagle/user/mrbar")
	failedBlock, err := isBlockedActor(failedBlockIRI)

	if failedBlock {
		t.Error("Invalid blocking of unblocked actor IRI")
	}
}
