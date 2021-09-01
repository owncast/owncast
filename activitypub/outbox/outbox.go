package outbox

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

func SendLive() {
	textContent := data.GetFederationGoLiveMessage()
	textContent = utils.RenderSimpleMarkdown(textContent)

	localActor := apmodels.MakeLocalIRIForAccount(data.GetDefaultFederationUsername())
	noteID := shortid.MustGenerate()
	noteIRI := apmodels.MakeLocalIRIForResource(noteID)
	id := shortid.MustGenerate()
	activity := apmodels.CreateCreateActivity(id, localActor)

	object := streams.NewActivityStreamsObjectProperty()

	activity.SetActivityStreamsObject(object)

	var tagStrings []string
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")

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
		tagString := fmt.Sprintf("<a class=\"hashtag\" href=\"https://directory.owncast.online/tags/%s\">#%s</a>", tagWithoutSpecialCharacters, tagWithoutSpecialCharacters)
		tagStrings = append(tagStrings, tagString)
	}
	tagsString := strings.Join(tagStrings, " ")

	activity.SetActivityStreamsTag(tagProp)

	var streamTitle string
	if title := data.GetStreamTitle(); title != "" {
		streamTitle = fmt.Sprintf("<p>%s</p>", title)
	}
	textContent = fmt.Sprintf("<p>%s</p><a href=\"%s\">%s</a>%s<p>%s</p>", textContent, data.GetServerURL(), data.GetServerURL(), streamTitle, tagsString)

	log.Println(textContent)

	note := apmodels.MakeNote(textContent, noteIRI, localActor)
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
			apmodels.AddImageAttachmentToNote(note, previewURL.String())
		}
	}

	b, err := apmodels.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize go live message activity", err)
		return
	}

	SendToFollowers(b)
	Add(note, noteID)
}

// SendPublicMessage will send a public message to all followers.
func SendPublicMessage(textContent string) {
	localActor := apmodels.MakeLocalIRIForAccount(data.GetDefaultFederationUsername())
	id := shortid.MustGenerate()
	activity := apmodels.CreateMessageActivity(id, textContent, localActor)
	message := activity.GetActivityStreamsObject().Begin().GetActivityStreamsNote()

	b, err := apmodels.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize send public message activity", err)
		return
	}
	SendToFollowers(b)

	Add(message, message.GetJSONLDId().Get().String())
}

func SendToFollowers(payload []byte) {
	localActor := apmodels.MakeLocalIRIForAccount(data.GetDefaultFederationUsername())

	followers, err := persistence.GetFederationFollowers()
	if err != nil {
		log.Errorln("unable to fetch followers to send to", err)
		return
	}

	for _, follower := range followers {
		if _, err := requests.PostSignedRequest(payload, follower.Inbox, localActor); err != nil {
			log.Errorln("unable to send to follower inbox", follower.Inbox, err)
			return
		}
	}
}

func Add(item vocab.Type, id string) error {
	iri := "/" + item.GetJSONLDId().GetIRI().Path
	typeString := item.GetTypeName()

	if iri == "" {
		log.Errorln("Unable to get iri from item")
		return errors.New("Unable to get iri from item " + id)
	}

	b, err := apmodels.Serialize(item)
	if err != nil {
		log.Errorln("unable to serialize model when saving to outbox", err)
		return err
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
