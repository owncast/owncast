package apmodels

import (
	"fmt"
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

// ActivityPubActor represents a single actor.
type ActivityPubActor struct {
	ActorIri                *url.URL
	FollowIri               *url.URL
	Inbox                   *url.URL
	Name                    string
	Username                string
	FullUsername            string
	Image                   *url.URL
	W3IDSecurityV1PublicKey vocab.W3IDSecurityV1PublicKeyProperty
}

// DeleteRequest represents a request for delete.
type DeleteRequest struct {
	ActorIri string
}

func MakeActorFromPerson(person vocab.ActivityStreamsPerson) ActivityPubActor {
	apActor := ActivityPubActor{
		ActorIri:                person.GetJSONLDId().Get(),
		Inbox:                   person.GetActivityStreamsInbox().GetIRI(),
		Name:                    person.GetActivityStreamsName().Begin().GetXMLSchemaString(),
		Username:                person.GetActivityStreamsPreferredUsername().GetXMLSchemaString(),
		FullUsername:            GetFullUsernameFromPerson(person),
		W3IDSecurityV1PublicKey: person.GetW3IDSecurityV1PublicKey(),
	}

	if person.GetActivityStreamsIcon() != nil && person.GetActivityStreamsIcon().Len() > 0 {
		apActor.Image = person.GetActivityStreamsIcon().At(0).GetActivityStreamsImage().GetActivityStreamsUrl().Begin().GetIRI()
	}

	return apActor
}

func MakeActorFromService(service vocab.ActivityStreamsService) ActivityPubActor {
	apActor := ActivityPubActor{
		ActorIri:                service.GetJSONLDId().Get(),
		Inbox:                   service.GetActivityStreamsInbox().GetIRI(),
		Name:                    service.GetActivityStreamsName().Begin().GetXMLSchemaString(),
		Username:                service.GetActivityStreamsPreferredUsername().GetXMLSchemaString(),
		FullUsername:            GetFullUsernameFromService(service),
		W3IDSecurityV1PublicKey: service.GetW3IDSecurityV1PublicKey(),
	}

	if service.GetActivityStreamsIcon() != nil && service.GetActivityStreamsIcon().Len() > 0 {
		apActor.Image = service.GetActivityStreamsIcon().At(0).GetActivityStreamsImage().GetActivityStreamsUrl().Begin().GetIRI()
	}

	return apActor
}

// MakeActorPropertyWithID will return an actor property filled with the provided IRI.
func MakeActorPropertyWithID(idIRI *url.URL) vocab.ActivityStreamsActorProperty {
	actor := streams.NewActivityStreamsActorProperty()
	actor.AppendIRI(idIRI)
	return actor
}

// MakePersonForAccount will create a new local actor person withe the provided username.
func MakePersonForAccount(accountName string) vocab.ActivityStreamsPerson {
	actorIRI := MakeLocalIRIForAccount(accountName)

	person := streams.NewActivityStreamsPerson()
	nameProperty := streams.NewActivityStreamsNameProperty()
	nameProperty.AppendXMLSchemaString(data.GetServerName())
	person.SetActivityStreamsName(nameProperty)

	preferredUsernameProperty := streams.NewActivityStreamsPreferredUsernameProperty()
	preferredUsernameProperty.SetXMLSchemaString(accountName)
	person.SetActivityStreamsPreferredUsername(preferredUsernameProperty)

	inboxIRI := MakeLocalIRIForResource("/user/" + accountName + "/inbox")

	inboxProp := streams.NewActivityStreamsInboxProperty()
	inboxProp.SetIRI(inboxIRI)
	person.SetActivityStreamsInbox(inboxProp)

	needsFollowApprovalProperty := streams.NewActivityStreamsManuallyApprovesFollowersProperty()
	needsFollowApprovalProperty.Set(data.GetFollowApprovalRequired())
	person.SetActivityStreamsManuallyApprovesFollowers(needsFollowApprovalProperty)

	outboxIRI := MakeLocalIRIForResource("/user/" + accountName + "/outbox")

	outboxProp := streams.NewActivityStreamsOutboxProperty()
	outboxProp.SetIRI(outboxIRI)
	person.SetActivityStreamsOutbox(outboxProp)

	id := streams.NewJSONLDIdProperty()
	id.Set(actorIRI)
	person.SetJSONLDId(id)

	publicKey := crypto.GetPublicKey(actorIRI)

	publicKeyProp := streams.NewW3IDSecurityV1PublicKeyProperty()
	publicKeyType := streams.NewW3IDSecurityV1PublicKey()

	pubKeyIDProp := streams.NewJSONLDIdProperty()
	pubKeyIDProp.Set(publicKey.ID)

	publicKeyType.SetJSONLDId(pubKeyIDProp)

	ownerProp := streams.NewW3IDSecurityV1OwnerProperty()
	ownerProp.SetIRI(publicKey.Owner)
	publicKeyType.SetW3IDSecurityV1Owner(ownerProp)

	publicKeyPemProp := streams.NewW3IDSecurityV1PublicKeyPemProperty()
	publicKeyPemProp.Set(publicKey.PublicKeyPem)
	publicKeyType.SetW3IDSecurityV1PublicKeyPem(publicKeyPemProp)
	publicKeyProp.AppendW3IDSecurityV1PublicKey(publicKeyType)
	person.SetW3IDSecurityV1PublicKey(publicKeyProp)

	if t, err := data.GetServerInitTime(); t != nil {
		publishedDateProp := streams.NewActivityStreamsPublishedProperty()
		publishedDateProp.Set(t.Time)
		person.SetActivityStreamsPublished(publishedDateProp)
	} else {
		log.Errorln("unable to fetch server init time", err)
	}

	// Profile properties

	// Avatar
	userAvatarURLString := data.GetServerURL() + "/logo/external"
	userAvatarURL, err := url.Parse(userAvatarURLString)
	if err != nil {
		log.Errorln("unable to parse user avatar url", userAvatarURLString, err)
	}

	image := streams.NewActivityStreamsImage()
	imgProp := streams.NewActivityStreamsUrlProperty()
	imgProp.AppendIRI(userAvatarURL)
	image.SetActivityStreamsUrl(imgProp)
	icon := streams.NewActivityStreamsIconProperty()
	icon.AppendActivityStreamsImage(image)
	person.SetActivityStreamsIcon(icon)

	// Actor  URL

	urlProperty := streams.NewActivityStreamsUrlProperty()
	urlProperty.AppendIRI(actorIRI)
	person.SetActivityStreamsUrl(urlProperty)

	// Profile header

	headerImage := streams.NewActivityStreamsImage()
	headerImgPropURL := streams.NewActivityStreamsUrlProperty()
	headerImgPropURL.AppendIRI(userAvatarURL)
	headerImage.SetActivityStreamsUrl(headerImgPropURL)
	headerImageProp := streams.NewActivityStreamsImageProperty()
	headerImageProp.AppendActivityStreamsImage(headerImage)
	person.SetActivityStreamsImage(headerImageProp)

	// Profile bio

	summaryProperty := streams.NewActivityStreamsSummaryProperty()
	summaryProperty.AppendXMLSchemaString(data.GetServerSummary())
	person.SetActivityStreamsSummary(summaryProperty)

	// Links
	for _, link := range data.GetSocialHandles() {
		addMetadataLinkToProfile(person, link.Platform, link.URL)
	}

	// Discoverable
	discoverableProperty := streams.NewTootDiscoverableProperty()
	discoverableProperty.Set(true)
	person.SetTootDiscoverable(discoverableProperty)

	// followers
	followersProperty := streams.NewActivityStreamsFollowersProperty()
	followersURL := *actorIRI
	followersURL.Path = actorIRI.Path + "/followers"
	followersProperty.SetIRI(&followersURL)
	person.SetActivityStreamsFollowers(followersProperty)

	// tags
	tagProp := streams.NewActivityStreamsTagProperty()
	for _, tagString := range data.GetServerMetadataTags() {
		hashtag := MakeHashtag(tagString)
		tagProp.AppendTootHashtag(hashtag)
	}

	person.SetActivityStreamsTag(tagProp)

	// Work around an issue where a single attachment will not serialize
	// as an array, so add another item to the mix.
	if len(data.GetSocialHandles()) == 1 {
		addMetadataLinkToProfile(person, "Owncast", "https://owncast.online")
	}

	return person
}

// GetFullUsernameFromPerson will return the user@host.tld formatted user given a person object.
func GetFullUsernameFromPerson(person vocab.ActivityStreamsPerson) string {
	hostname := person.GetJSONLDId().GetIRI().Hostname()
	username := person.GetActivityStreamsPreferredUsername().GetXMLSchemaString()
	fullUsername := fmt.Sprintf("%s@%s", username, hostname)

	return fullUsername
}

// GetFullUsernameFromService will return the user@host.tld formatted user given a service object.
func GetFullUsernameFromService(person vocab.ActivityStreamsService) string {
	hostname := person.GetJSONLDId().GetIRI().Hostname()
	username := person.GetActivityStreamsPreferredUsername().GetXMLSchemaString()
	fullUsername := fmt.Sprintf("%s@%s", username, hostname)

	return fullUsername
}

func addMetadataLinkToProfile(profile vocab.ActivityStreamsPerson, name string, url string) {
	attachments := profile.GetActivityStreamsAttachment()
	if attachments == nil {
		attachments = streams.NewActivityStreamsAttachmentProperty()
	}

	displayName := name
	socialHandle := models.GetSocialHandle(name)
	if socialHandle != nil {
		displayName = socialHandle.Platform
	}

	linkValue := fmt.Sprintf("<a href=\"%s\" rel=\"me nofollow noopener noreferrer\" target=\"_blank\">%s</a>", url, url)

	attachment := streams.NewActivityStreamsObject()
	attachmentProp := streams.NewJSONLDTypeProperty()
	attachmentProp.AppendXMLSchemaString("PropertyValue")
	attachment.SetJSONLDType(attachmentProp)
	attachmentName := streams.NewActivityStreamsNameProperty()
	attachmentName.AppendXMLSchemaString(displayName)
	attachment.SetActivityStreamsName(attachmentName)
	attachment.GetUnknownProperties()["value"] = linkValue

	attachments.AppendActivityStreamsObject(attachment)
	profile.SetActivityStreamsAttachment(attachments)
}
