package inbox

import (
	"net/url"
	"testing"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"

	"github.com/owncast/owncast/internal/activitypub/apmodels"
	"github.com/owncast/owncast/internal/activitypub/persistence"
)

func (s *Service) makeFakePerson() vocab.ActivityStreamsPerson {
	iri, _ := url.Parse("https://freedom.eagle/user/mrfoo")
	name := "Mr CallServiceB"
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

func (s *Service) TestBlockedDomains(t *testing.T) {
	person := s.makeFakePerson()

	s.Data.SetBlockedFederatedDomains([]string{"freedom.eagle", "guns.life"})

	if len(s.Data.GetBlockedFederatedDomains()) != 2 {
		t.Error("Blocked federated domains is not set correctly")
	}

	for _, domain := range s.Data.GetBlockedFederatedDomains() {
		if domain == person.GetJSONLDId().GetIRI().Host {
			return
		}
	}

	t.Error("Failed to catch blocked domain")
}

func (s *Service) TestBlockedActors(t *testing.T) {
	s.Data.Init(":memory:")
	s.Data.SetServerURL("https://my.cool.site.biz")

	p, err := persistence.New(s.Data)
	if err != nil {
		t.Fatalf("creating persistence service: %v", err)
	}

	person := s.makeFakePerson()
	fakeRequest := streams.NewActivityStreamsFollow()
	p.AddFollow(apmodels.ActivityPubActor{
		ActorIri:         person.GetJSONLDId().GetIRI(),
		Inbox:            person.GetJSONLDId().GetIRI(),
		FollowRequestIri: person.GetJSONLDId().GetIRI(),
		RequestObject:    fakeRequest,
	}, false)

	p.BlockOrRejectFollower(person.GetJSONLDId().GetIRI().String())

	blocked, err := s.isBlockedActor(person.GetJSONLDId().GetIRI())
	if err != nil {
		t.Error(err)
		return
	}

	if !blocked {
		t.Error("Failed to block actor")
	}

	failedBlockIRI, _ := url.Parse("https://freedom.eagle/user/mrbar")
	failedBlock, err := s.isBlockedActor(failedBlockIRI)

	if failedBlock {
		t.Error("Invalid blocking of unblocked actor IRI")
	}
}
