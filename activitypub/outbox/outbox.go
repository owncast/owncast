package outbox

import (
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/models"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
	"github.com/teris-io/shortid"
)

func SendLive() {
	textContent := data.GetFederationGoLiveMessage()
	textContent = utils.RenderSimpleMarkdown(textContent)

	localActor := models.MakeLocalIRIForAccount(data.GetDefaultFederationUsername())
	noteId := models.MakeLocalIRIForResource(data.GetDefaultFederationUsername() + "/" + shortid.MustGenerate())
	id := shortid.MustGenerate()
	activity := models.CreateCreateActivity(id, localActor)

	object := streams.NewActivityStreamsObjectProperty()

	activity.SetActivityStreamsObject(object)

	var tagStrings []string
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")

	// TODO: Need to assign to a `hashtag` property, not `tag`
	tagProp := streams.NewActivityStreamsTagProperty()
	for _, tagString := range data.GetServerMetadataTags() {
		mention := streams.NewActivityStreamsMention()
		mentionName := streams.NewActivityStreamsNameProperty()
		mentionName.AppendXMLSchemaString("#" + tagString)
		mention.SetActivityStreamsName(mentionName)

		name := streams.NewActivityStreamsNameProperty()
		name.AppendXMLSchemaString("#" + tagString)

		tagProp.AppendActivityStreamsMention(mention)

		// TODO: Do we want to display tags or just assign them?
		tagWithoutSpecialCharacters := reg.ReplaceAllString(tagString, "")
		tagStrings = append(tagStrings, "#"+tagWithoutSpecialCharacters)
	}
	tagsString := strings.Join(tagStrings, " ")

	activity.SetActivityStreamsTag(tagProp)

	textContent = textContent + "\n\n" + tagsString

	note := models.MakeNote(textContent, noteId, localActor)
	object.AppendActivityStreamsNote(note)

	// Attach an image along with the Federated message.
	previewURL, err := url.Parse(data.GetServerURL())
	if err == nil {
		var imageToAttach string
		previewGif := filepath.Join(config.WebRoot, "preview.gif")
		thumbnailJpg := filepath.Join(config.WebRoot, "thumbnail.jpg")

		if utils.DoesFileExists(previewGif) {
			imageToAttach = "preview.gif"
		} else if utils.DoesFileExists(thumbnailJpg) {
			imageToAttach = "thumbnail.jpg"
		}
		if imageToAttach != "" {
			previewURL.Path = imageToAttach
			models.AddImageAttachmentToNote(note, previewURL.String())
		}
	}

	b, err := models.Serialize(activity)
	if err != nil {
		panic(err)
	}
	SendToFollowers(b)
	Add(activity, id)
}

// SendPublicMessage will send a public message to all followers.
func SendPublicMessage(textContent string) {
	localActor := models.MakeLocalIRIForAccount(data.GetDefaultFederationUsername())
	id := shortid.MustGenerate()
	message := models.CreateMessageActivity(id, textContent, localActor)

	b, err := models.Serialize(message)
	if err != nil {
		panic(err)
	}
	SendToFollowers(b)

	Add(message, id)
}

func SendToFollowers(payload []byte) {
	localActor := models.MakeLocalIRIForAccount(data.GetDefaultFederationUsername())

	followers, err := persistence.GetFederationFollowers()
	if err != nil {
		panic(err)
	}

	for _, follower := range followers {
		if _, err := requests.PostSignedRequest(payload, follower.Inbox, localActor); err != nil {
			panic(err)
		}
	}
}

func Add(item vocab.Type, id string) error {
	iri := item.GetJSONLDId().GetIRI().String()
	typeString := item.GetTypeName()

	if iri == "" {
		panic("Unable to get iri from item")
	}

	b, err := models.Serialize(item)
	if err != nil {
		panic(err)
	}

	return persistence.AddToOutbox(id, iri, b, typeString)
}

func Get() vocab.ActivityStreamsOrderedCollectionPage {
	orderedCollection, _ := persistence.GetOutbox()
	return orderedCollection
}

// func AddPayloadToOutbox(id string, payload []byte) {
// 	persistence.AddToOutbox(id, payload)
// }
