package outbox

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/go-fed/activity/streams"
	"github.com/owncast/owncast/activitypub/models"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
	"github.com/teris-io/shortid"
)

func SendLive() {
	liveMessage := fmt.Sprintf("%s is live!", data.GetServerName())
	streamDescription := data.GetStreamTitle()
	if streamDescription == "" {
		streamDescription = data.GetServerSummary()
	}

	textContent := fmt.Sprintf("%s\n\n%s", liveMessage, streamDescription)
	textContent = utils.RenderSimpleMarkdown(textContent)

	localActor := models.MakeLocalIRIForAccount(data.GetDefaultFederationUsername())
	noteId := models.MakeLocalIRIForResource(data.GetDefaultFederationUsername() + "/" + shortid.MustGenerate())

	activity := models.CreateCreateActivity(localActor)

	object := streams.NewActivityStreamsObjectProperty()
	note := models.MakeNote(textContent, noteId, localActor)
	object.AppendActivityStreamsNote(note)

	activity.SetActivityStreamsObject(object)

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

	// var tagStrings []string
	// reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	// if err != nil {
	// 	fmt.Println(err)
	// }
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
		// tagWithoutSpecialCharacters := reg.ReplaceAllString(tag, "")
		// tagStrings = append(tagStrings, "#"+tagWithoutSpecialCharacters)
	}
	activity.SetActivityStreamsTag(tagProp)

	b, err := models.Serialize(activity)
	if err != nil {
		panic(err)
	}
	SendToFollowers(b)
}

// SendPublicMessage will send a public message to all followers.
func SendPublicMessage(textContent string) {
	localActor := models.MakeLocalIRIForAccount(data.GetDefaultFederationUsername())
	message := models.CreateMessageActivity(textContent, localActor)

	b, err := models.Serialize(message)
	if err != nil {
		panic(err)
	}
	SendToFollowers(b)
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
