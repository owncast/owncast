package apmodels

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/models"
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
	// SharedInbox is the shared inbox URL of the remote server (optional)
	SharedInbox *url.URL
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

// SharedInboxString returns the string representation of SharedInbox, or empty string if nil.
func (a *ActivityPubActor) SharedInboxString() string {
	if a.SharedInbox == nil {
		return ""
	}
	return a.SharedInbox.String()
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

// validateEntityRequiredFields checks that all required fields are present on the entity.
func validateEntityRequiredFields(entity ExternalEntity) error {
	if entity.GetJSONLDId() == nil || entity.GetJSONLDId().Get() == nil {
		return fmt.Errorf("%w: entity is missing actor IRI", ErrActorMissingRequiredField)
	}
	if entity.GetActivityStreamsInbox() == nil || entity.GetActivityStreamsInbox().GetIRI() == nil {
		return fmt.Errorf("%w: entity is missing inbox", ErrActorMissingRequiredField)
	}
	if entity.GetActivityStreamsPreferredUsername() == nil || entity.GetActivityStreamsPreferredUsername().GetXMLSchemaString() == "" {
		return fmt.Errorf("%w: entity is missing preferred username", ErrActorMissingRequiredField)
	}
	if entity.GetW3IDSecurityV1PublicKey() == nil || entity.GetW3IDSecurityV1PublicKey().Len() == 0 {
		return fmt.Errorf("%w: entity is missing public key", ErrActorMissingRequiredField)
	}
	return nil
}

// getNameFromEntity extracts the optional name from an entity.
func getNameFromEntity(entity ExternalEntity) string {
	nameProp := entity.GetActivityStreamsName()
	if nameProp == nil || nameProp.Empty() {
		return ""
	}
	return nameProp.At(0).GetXMLSchemaString()
}

// getSharedInboxFromEntity extracts the optional shared inbox URL from an entity.
func getSharedInboxFromEntity(entity ExternalEntity) *url.URL {
	endpointsProp := entity.GetActivityStreamsEndpoints()
	if endpointsProp == nil || !endpointsProp.IsActivityStreamsEndpoints() {
		return nil
	}

	endpoints := endpointsProp.Get()
	if endpoints == nil {
		return nil
	}

	sharedInboxProp := endpoints.GetActivityStreamsSharedInbox()
	if sharedInboxProp == nil || !sharedInboxProp.HasAny() {
		return nil
	}

	return sharedInboxProp.Get()
}

// NewActivityPubActorFromEntity creates a new ActivityPubActor from an external entity
// with validation of required fields.
func NewActivityPubActorFromEntity(entity ExternalEntity) (*ActivityPubActor, error) {
	if err := validateEntityRequiredFields(entity); err != nil {
		return nil, err
	}

	apActor := &ActivityPubActor{
		ActorIri:                entity.GetJSONLDId().Get(),
		Inbox:                   entity.GetActivityStreamsInbox().GetIRI(),
		SharedInbox:             getSharedInboxFromEntity(entity),
		Name:                    getNameFromEntity(entity),
		Username:                entity.GetActivityStreamsPreferredUsername().GetXMLSchemaString(),
		FullUsername:            GetFullUsernameFromExternalEntity(entity),
		W3IDSecurityV1PublicKey: entity.GetW3IDSecurityV1PublicKey(),
		Image:                   GetImageFromIcon(entity.GetActivityStreamsIcon()),
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
	GetActivityStreamsEndpoints() vocab.ActivityStreamsEndpointsProperty
}

// MakeActorPropertyWithID will return an actor property filled with the provided IRI.
func MakeActorPropertyWithID(idIRI *url.URL) vocab.ActivityStreamsActorProperty {
	actor := streams.NewActivityStreamsActorProperty()
	actor.AppendIRI(idIRI)
	return actor
}

// MakeServiceForAccount will create a new local actor service with the the provided username.
func (b *Builder) MakeServiceForAccount(accountName string) vocab.ActivityStreamsService {
	actorIRI := b.MakeLocalIRIForAccount(accountName)

	person := streams.NewActivityStreamsService()
	nameProperty := streams.NewActivityStreamsNameProperty()
	nameProperty.AppendXMLSchemaString(b.configRepository.GetServerName())
	person.SetActivityStreamsName(nameProperty)

	preferredUsernameProperty := streams.NewActivityStreamsPreferredUsernameProperty()
	preferredUsernameProperty.SetXMLSchemaString(accountName)
	person.SetActivityStreamsPreferredUsername(preferredUsernameProperty)

	inboxIRI := b.MakeLocalIRIForResource("/user/" + accountName + "/inbox")

	inboxProp := streams.NewActivityStreamsInboxProperty()
	inboxProp.SetIRI(inboxIRI)
	person.SetActivityStreamsInbox(inboxProp)

	needsFollowApprovalProperty := streams.NewActivityStreamsManuallyApprovesFollowersProperty()
	needsFollowApprovalProperty.Set(b.configRepository.GetFederationIsPrivate())
	person.SetActivityStreamsManuallyApprovesFollowers(needsFollowApprovalProperty)

	outboxIRI := b.MakeLocalIRIForResource("/user/" + accountName + "/outbox")

	outboxProp := streams.NewActivityStreamsOutboxProperty()
	outboxProp.SetIRI(outboxIRI)
	person.SetActivityStreamsOutbox(outboxProp)

	id := streams.NewJSONLDIdProperty()
	id.Set(actorIRI)
	person.SetJSONLDId(id)

	publicKey := b.signer.GetPublicKey(actorIRI)

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

	if t, err := b.configRepository.GetServerInitTime(); t != nil {
		publishedDateProp := streams.NewActivityStreamsPublishedProperty()
		publishedDateProp.Set(t.Time)
		person.SetActivityStreamsPublished(publishedDateProp)
	} else {
		log.Errorln("unable to fetch server init time", err)
	}

	// Profile properties

	// Avatar
	uniquenessString := b.configRepository.GetLogoUniquenessString()
	userAvatarURLString := b.configRepository.GetServerURL() + "/logo/external"
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
	summaryProperty.AppendXMLSchemaString(b.configRepository.GetServerSummary())
	person.SetActivityStreamsSummary(summaryProperty)

	// Links
	if serverURL := b.configRepository.GetServerURL(); serverURL != "" {
		addMetadataLinkToProfile(person, "Stream", serverURL)
	}
	for _, link := range b.configRepository.GetSocialHandles() {
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
	for _, tagString := range b.configRepository.GetServerMetadataTags() {
		hashtag := MakeHashtag(tagString)
		tagProp.AppendTootHashtag(hashtag)
	}

	person.SetActivityStreamsTag(tagProp)

	// Work around an issue where a single attachment will not serialize
	// as an array, so add another item to the mix.
	if len(b.configRepository.GetSocialHandles()) == 1 {
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
