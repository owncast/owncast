package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/models"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/core/data"
)

func ActorHandler(w http.ResponseWriter, r *http.Request) {
	pathComponents := strings.Split(r.URL.Path, "/")
	accountName := pathComponents[3]

	if _, valid := data.GetFederatedInboxMap()[accountName]; !valid {
		// User is not valid
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// If this request is for an actor's inbox then pass
	// the request to the inbox controller.
	if len(pathComponents) == 5 && pathComponents[4] == "inbox" {
		InboxHandler(w, r)
		return
	} else if len(pathComponents) == 5 && pathComponents[4] == "outbox" {
		OutboxHandler(w, r)
		return
	} else if len(pathComponents) == 5 {
		ActorObjectHandler(w, r)
	} else {
		// dunno
	}
	actorIRI := models.MakeLocalIRIForAccount(accountName)

	person := streams.NewActivityStreamsPerson()
	nameProperty := streams.NewActivityStreamsNameProperty()
	nameProperty.AppendXMLSchemaString(data.GetServerName())
	person.SetActivityStreamsName(nameProperty)

	preferredUsernameProperty := streams.NewActivityStreamsPreferredUsernameProperty()
	preferredUsernameProperty.SetXMLSchemaString(accountName)
	person.SetActivityStreamsPreferredUsername(preferredUsernameProperty)

	inboxIRI := models.MakeLocalIRIForResource("/user/" + accountName + "/inbox")

	inboxProp := streams.NewActivityStreamsInboxProperty()
	inboxProp.SetIRI(inboxIRI)
	person.SetActivityStreamsInbox(inboxProp)

	outboxIRI := models.MakeLocalIRIForResource("/user/" + accountName + "/outbox")

	outboxProp := streams.NewActivityStreamsOutboxProperty()
	outboxProp.SetIRI(outboxIRI)
	person.SetActivityStreamsOutbox(outboxProp)

	id := streams.NewJSONLDIdProperty()
	id.Set(actorIRI)
	person.SetJSONLDId(id)

	publicKey := crypto.GetPublicKey(actorIRI)

	publicKeyProp := streams.NewW3IDSecurityV1PublicKeyProperty()
	publicKeyType := streams.NewW3IDSecurityV1PublicKey()

	pubKeyIdProp := streams.NewJSONLDIdProperty()
	pubKeyIdProp.Set(publicKey.Id)

	publicKeyType.SetJSONLDId(pubKeyIdProp)

	ownerProp := streams.NewW3IDSecurityV1OwnerProperty()
	ownerProp.SetIRI(publicKey.Owner)
	publicKeyType.SetW3IDSecurityV1Owner(ownerProp)

	publicKeyPemProp := streams.NewW3IDSecurityV1PublicKeyPemProperty()
	publicKeyPemProp.Set(publicKey.PublicKeyPem)
	publicKeyType.SetW3IDSecurityV1PublicKeyPem(publicKeyPemProp)
	publicKeyProp.AppendW3IDSecurityV1PublicKey(publicKeyType)
	person.SetW3IDSecurityV1PublicKey(publicKeyProp)

	if t, err := data.GetServerInitTime(); t.Valid {
		publishedDateProp := streams.NewActivityStreamsPublishedProperty()
		publishedDateProp.Set(t.Time)
		person.SetActivityStreamsPublished(publishedDateProp)
		log.Errorln(err)
	}
	// Profile properties

	// Avatar
	userAvatarUrlString := data.GetServerURL() + "/logo"
	userAvatarUrl, err := url.Parse(userAvatarUrlString)
	if err != nil {
		panic(err)
	}

	image := streams.NewActivityStreamsImage()
	imgProp := streams.NewActivityStreamsUrlProperty()
	imgProp.AppendIRI(userAvatarUrl)
	image.SetActivityStreamsUrl(imgProp)
	icon := streams.NewActivityStreamsIconProperty()
	icon.AppendActivityStreamsImage(image)
	person.SetActivityStreamsIcon(icon)

	// Site URL

	siteURL, err := url.Parse(data.GetServerURL())
	if err != nil {
		panic(err)
	}

	link := streams.NewActivityStreamsLink()
	hrefLink := streams.NewActivityStreamsHrefProperty()
	hrefLink.Set(siteURL)
	link.SetActivityStreamsHref(hrefLink)
	urlProperty := streams.NewActivityStreamsUrlProperty()
	urlProperty.AppendActivityStreamsLink(link)
	person.SetActivityStreamsUrl(urlProperty)

	// Profile header

	headerImage := streams.NewActivityStreamsImage()
	headerImgPropUrl := streams.NewActivityStreamsUrlProperty()
	headerImgPropUrl.AppendIRI(userAvatarUrl)
	headerImage.SetActivityStreamsUrl(headerImgPropUrl)
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

	// Work around an issue where a single attachment will not serialize
	// as an array, so add another item to the mix.
	if len(data.GetSocialHandles()) == 1 {
		addMetadataLinkToProfile(person, "Owncast", "https://owncast.online")
	}

	if err := requests.WriteStreamResponse(person, w, publicKey); err != nil {
		fmt.Println(err)
	}
}

func addMetadataLinkToProfile(profile vocab.ActivityStreamsPerson, name string, url string) {
	var attachments = profile.GetActivityStreamsAttachment()
	if attachments == nil {
		attachments = streams.NewActivityStreamsAttachmentProperty()
	}

	linkValue := fmt.Sprintf("<a href=\"%s\" rel=\"me nofollow noopener noreferrer\" target=\"_blank\">%s</a>", url, url)

	attachment := streams.NewActivityStreamsObject()
	attachmentProp := streams.NewJSONLDTypeProperty()
	attachmentProp.AppendXMLSchemaString("PropertyValue")
	attachment.SetJSONLDType(attachmentProp)
	attachmentName := streams.NewActivityStreamsNameProperty()
	attachmentName.AppendXMLSchemaString(name)
	attachment.SetActivityStreamsName(attachmentName)
	attachment.GetUnknownProperties()["value"] = linkValue

	attachments.AppendActivityStreamsObject(attachment)
	profile.SetActivityStreamsAttachment(attachments)
}
