package apmodels

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	log "github.com/sirupsen/logrus"
)

// ActivityPubActor represents a single actor in handling ActivityPub activity.
type ActivityPubActor struct {
	// RequestObject is the actual follow request object.
	RequestObject vocab.ActivityStreamsFollow
	// W3IDSecurityV1PublicKey is the public key of the actor.
	W3IDSecurityV1PublicKey vocab.W3IDSecurityV1PublicKeyProperty
	// ActorIRI is the IRI of the remote actor.
	ActorIri *url.URL
	// FollowRequestIRI is the unique identifier of the follow request.
	FollowRequestIri *url.URL
	// Inbox is the inbox URL of the remote follower
	Inbox *url.URL
	// Image is the avatar image of the Actor.
	Image *url.URL
	// DisabledAt is the time, if any, this follower was blocked/removed.
	DisabledAt *time.Time
	// Name is the display name of the follower.
	Name string
	// Username is the account username of the remote actor.
	Username string
	// FullUsername is the username@account.tld representation of the user.
	FullUsername string
}

// ErrActorMissingRequiredField is returned when an actor is missing a required field.
var ErrActorMissingRequiredField = errors.New("actor missing required field")

// Validate checks that required fields are present on the actor.
// Returns an error if ActorIri or Inbox are nil.
func (a *ActivityPubActor) Validate() error {
	if a.ActorIri == nil {
		return fmt.Errorf("%w: ActorIri is required", ErrActorMissingRequiredField)
	}
	if a.Inbox == nil {
		return fmt.Errorf("%w: Inbox is required", ErrActorMissingRequiredField)
	}
	return nil
}

// IsValid returns true if the actor has all required fields.
func (a *ActivityPubActor) IsValid() bool {
	return a.Validate() == nil
}

// ActorIriString returns the string representation of ActorIri, or empty string if nil.
func (a *ActivityPubActor) ActorIriString() string {
	if a.ActorIri == nil {
		return ""
	}
	return a.ActorIri.String()
}

// InboxString returns the string representation of Inbox, or empty string if nil.
func (a *ActivityPubActor) InboxString() string {
	if a.Inbox == nil {
		return ""
	}
	return a.Inbox.String()
}

// ImageString returns the string representation of Image, or empty string if nil.
func (a *ActivityPubActor) ImageString() string {
	if a.Image == nil {
		return ""
	}
	return a.Image.String()
}

// FollowRequestIriString returns the string representation of FollowRequestIri, or empty string if nil.
func (a *ActivityPubActor) FollowRequestIriString() string {
	if a.FollowRequestIri == nil {
		return ""
	}
	return a.FollowRequestIri.String()
}

// ActorIriHostname returns the hostname of ActorIri, or empty string if nil.
func (a *ActivityPubActor) ActorIriHostname() string {
	if a.ActorIri == nil {
		return ""
	}
	return a.ActorIri.Hostname()
}

// NewActivityPubActor creates a new ActivityPubActor with required fields.
// Returns an error if actorIri or inbox are nil.
func NewActivityPubActor(actorIri, inbox *url.URL) (*ActivityPubActor, error) {
	if actorIri == nil {
		return nil, fmt.Errorf("%w: actorIri is required", ErrActorMissingRequiredField)
	}
	if inbox == nil {
		return nil, fmt.Errorf("%w: inbox is required", ErrActorMissingRequiredField)
	}
	return &ActivityPubActor{
		ActorIri: actorIri,
		Inbox:    inbox,
	}, nil
}

// NewActivityPubActorFromEntity creates a new ActivityPubActor from an external entity
// with validation of required fields.
func NewActivityPubActorFromEntity(entity ExternalEntity) (*ActivityPubActor, error) {
	// ActorIri is required (must validate before GetFullUsernameFromExternalEntity which uses it)
	if entity.GetJSONLDId() == nil || entity.GetJSONLDId().Get() == nil {
		return nil, fmt.Errorf("%w: entity is missing actor IRI", ErrActorMissingRequiredField)
	}
	actorIri := entity.GetJSONLDId().Get()

	// Inbox is required
	if entity.GetActivityStreamsInbox() == nil || entity.GetActivityStreamsInbox().GetIRI() == nil {
		return nil, fmt.Errorf("%w: entity is missing inbox", ErrActorMissingRequiredField)
	}
	inbox := entity.GetActivityStreamsInbox().GetIRI()

	// Username is required (but not a part of the official ActivityPub spec)
	if entity.GetActivityStreamsPreferredUsername() == nil || entity.GetActivityStreamsPreferredUsername().GetXMLSchemaString() == "" {
		return nil, fmt.Errorf("%w: entity is missing preferred username", ErrActorMissingRequiredField)
	}
	username := GetFullUsernameFromExternalEntity(entity)

	// Key is required
	if entity.GetW3IDSecurityV1PublicKey() == nil || entity.GetW3IDSecurityV1PublicKey().Len() == 0 {
		return nil, fmt.Errorf("%w: entity is missing public key", ErrActorMissingRequiredField)
	}

	// Name is optional
	var name string
	if entity.GetActivityStreamsName() != nil && !entity.GetActivityStreamsName().Empty() {
		name = entity.GetActivityStreamsName().At(0).GetXMLSchemaString()
	}

	// Image is optional
	image := GetImageFromIcon(entity.GetActivityStreamsIcon())

	apActor := &ActivityPubActor{
		ActorIri:                actorIri,
		Inbox:                   inbox,
		Name:                    name,
		Username:                entity.GetActivityStreamsPreferredUsername().GetXMLSchemaString(),
		FullUsername:            username,
		W3IDSecurityV1PublicKey: entity.GetW3IDSecurityV1PublicKey(),
		Image:                   image,
	}

	return apActor, nil
}

// DeleteRequest represents a request for delete.
type DeleteRequest struct {
	ActorIri string
}

// ExternalEntity represents an ActivityPub Person, Service or Application.
type ExternalEntity interface {
	GetJSONLDId() vocab.JSONLDIdProperty
	GetActivityStreamsInbox() vocab.ActivityStreamsInboxProperty
	GetActivityStreamsName() vocab.ActivityStreamsNameProperty
	GetActivityStreamsPreferredUsername() vocab.ActivityStreamsPreferredUsernameProperty
	GetActivityStreamsIcon() vocab.ActivityStreamsIconProperty
	GetW3IDSecurityV1PublicKey() vocab.W3IDSecurityV1PublicKeyProperty
}

// MakeActorPropertyWithID will return an actor property filled with the provided IRI.
func MakeActorPropertyWithID(idIRI *url.URL) vocab.ActivityStreamsActorProperty {
	actor := streams.NewActivityStreamsActorProperty()
	actor.AppendIRI(idIRI)
	return actor
}

// MakeServiceForAccount will create a new local actor service with the the provided username.
func MakeServiceForAccount(accountName string) vocab.ActivityStreamsService {
	configRepository := configrepository.Get()

	actorIRI := MakeLocalIRIForAccount(accountName)

	person := streams.NewActivityStreamsService()
	nameProperty := streams.NewActivityStreamsNameProperty()
	nameProperty.AppendXMLSchemaString(configRepository.GetServerName())
	person.SetActivityStreamsName(nameProperty)

	preferredUsernameProperty := streams.NewActivityStreamsPreferredUsernameProperty()
	preferredUsernameProperty.SetXMLSchemaString(accountName)
	person.SetActivityStreamsPreferredUsername(preferredUsernameProperty)

	inboxIRI := MakeLocalIRIForResource("/user/" + accountName + "/inbox")

	inboxProp := streams.NewActivityStreamsInboxProperty()
	inboxProp.SetIRI(inboxIRI)
	person.SetActivityStreamsInbox(inboxProp)

	needsFollowApprovalProperty := streams.NewActivityStreamsManuallyApprovesFollowersProperty()
	needsFollowApprovalProperty.Set(configRepository.GetFederationIsPrivate())
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

	if t, err := configRepository.GetServerInitTime(); t != nil {
		publishedDateProp := streams.NewActivityStreamsPublishedProperty()
		publishedDateProp.Set(t.Time)
		person.SetActivityStreamsPublished(publishedDateProp)
	} else {
		log.Errorln("unable to fetch server init time", err)
	}

	// Profile properties

	// Avatar
	uniquenessString := configRepository.GetLogoUniquenessString()
	userAvatarURLString := configRepository.GetServerURL() + "/logo/external"
	userAvatarURL, err := url.Parse(userAvatarURLString)
	userAvatarURL.RawQuery = "uc=" + uniquenessString
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
	summaryProperty.AppendXMLSchemaString(configRepository.GetServerSummary())
	person.SetActivityStreamsSummary(summaryProperty)

	// Links
	if serverURL := configRepository.GetServerURL(); serverURL != "" {
		addMetadataLinkToProfile(person, "Stream", serverURL)
	}
	for _, link := range configRepository.GetSocialHandles() {
		addMetadataLinkToProfile(person, link.Platform, link.URL)
	}

	// Discoverable
	discoverableProperty := streams.NewTootDiscoverableProperty()
	discoverableProperty.Set(true)
	person.SetTootDiscoverable(discoverableProperty)

	// Followers
	followersProperty := streams.NewActivityStreamsFollowersProperty()
	followersURL := *actorIRI
	followersURL.Path = actorIRI.Path + "/followers"
	followersProperty.SetIRI(&followersURL)
	person.SetActivityStreamsFollowers(followersProperty)

	// Tags
	tagProp := streams.NewActivityStreamsTagProperty()
	for _, tagString := range configRepository.GetServerMetadataTags() {
		hashtag := MakeHashtag(tagString)
		tagProp.AppendTootHashtag(hashtag)
	}

	person.SetActivityStreamsTag(tagProp)

	// Work around an issue where a single attachment will not serialize
	// as an array, so add another item to the mix.
	if len(configRepository.GetSocialHandles()) == 1 {
		addMetadataLinkToProfile(person, "Owncast", "https://owncast.online")
	}

	return person
}

// GetFullUsernameFromExternalEntity will return the full username from an
// internal representation of an ExternalEntity. Returns user@host.tld.
func GetFullUsernameFromExternalEntity(entity ExternalEntity) string {
	hostname := GetHostnameFromJSONLDId(entity.GetJSONLDId())
	username := entity.GetActivityStreamsPreferredUsername().GetXMLSchemaString()
	fullUsername := fmt.Sprintf("%s@%s", username, hostname)

	return fullUsername
}

func addMetadataLinkToProfile(profile vocab.ActivityStreamsService, name string, url string) {
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
