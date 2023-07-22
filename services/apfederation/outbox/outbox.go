package outbox

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
<<<<<<< HEAD:activitypub/outbox/outbox.go
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/activitypub/webfinger"
	"github.com/owncast/owncast/activitypub/workerpool"
	"github.com/owncast/owncast/core/data"
=======

>>>>>>> 4f9fbfba1 (WIP):services/apfederation/outbox/outbox.go
	"github.com/owncast/owncast/storage/configrepository"
	"github.com/owncast/owncast/storage/federationrepository"
	"github.com/pkg/errors"

	"github.com/owncast/owncast/services/apfederation/apmodels"
	"github.com/owncast/owncast/services/apfederation/crypto"
	"github.com/owncast/owncast/services/apfederation/requests"
	"github.com/owncast/owncast/services/apfederation/resolvers"
	"github.com/owncast/owncast/services/apfederation/webfinger"
	"github.com/owncast/owncast/services/apfederation/workerpool"
	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

type APOutbox struct {
	configRepository     configrepository.ConfigRepository
	federationRepository *federationrepository.FederationRepository
	resolvers            *resolvers.APResolvers
	workerpool           *workerpool.WorkerPool
	requests             *requests.Requests
}

func New() *APOutbox {
	return &APOutbox{
		configRepository:     configrepository.Get(),
		federationRepository: federationrepository.Get(),
		resolvers:            resolvers.Get(),
		workerpool:           workerpool.Get(),
		requests:             requests.Get(),
	}
}

var temporaryGlobalInstance *APOutbox

func Get() *APOutbox {
	if temporaryGlobalInstance == nil {
		temporaryGlobalInstance = New()
	}
	return temporaryGlobalInstance
}

// SendLive will send all followers the message saying you started a live stream.
func (apo *APOutbox) SendLive() error {
	textContent := apo.configRepository.GetFederationGoLiveMessage()

	// If the message is empty then do not send it.
	if textContent == "" {
		return nil
	}

	tagStrings := []string{}
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")

	tagProp := streams.NewActivityStreamsTagProperty()
	for _, tagString := range apo.configRepository.GetServerMetadataTags() {
		tagWithoutSpecialCharacters := reg.ReplaceAllString(tagString, "")
		hashtag := apmodels.MakeHashtag(tagWithoutSpecialCharacters)
		tagProp.AppendTootHashtag(hashtag)
		tagString := apo.getHashtagLinkHTMLFromTagString(tagWithoutSpecialCharacters)
		tagStrings = append(tagStrings, tagString)
	}

	// Manually add Owncast hashtag if it doesn't already exist so it shows up
	// in Owncast search results.
	// We can remove this down the road, but it'll be nice for now.
	if _, exists := utils.FindInSlice(tagStrings, "owncast"); !exists {
		hashtag := apmodels.MakeHashtag("owncast")
		tagProp.AppendTootHashtag(hashtag)
	}

	tagsString := strings.Join(tagStrings, " ")

	var streamTitle string
	if title := apo.configRepository.GetStreamTitle(); title != "" {
		streamTitle = fmt.Sprintf("<p>%s</p>", title)
	}
<<<<<<< HEAD:activitypub/outbox/outbox.go
	textContent = fmt.Sprintf("<p>%s</p>%s<p>%s</p><p><a href=\"%s\">%s</a></p>", textContent, streamTitle, tagsString, data.GetServerURL(), data.GetServerURL())
=======
	textContent = fmt.Sprintf("<p>%s</p>%s<p>%s</p><a href=\"%s\">%s</a>", textContent, streamTitle, tagsString, apo.configRepository.GetServerURL(), apo.configRepository.GetServerURL())
>>>>>>> 4f9fbfba1 (WIP):services/apfederation/outbox/outbox.go

	activity, _, note, noteID := apo.createBaseOutboundMessage(textContent)

	// To the public if we're not treating ActivityPub as "private".
	if !apo.configRepository.GetFederationIsPrivate() {
		note = apmodels.MakeNotePublic(note)
		activity = apmodels.MakeActivityPublic(activity)
	}

	note.SetActivityStreamsTag(tagProp)

	// Attach an image along with the Federated message.
	previewURL, err := url.Parse(apo.configRepository.GetServerURL())
	if err == nil {
		var imageToAttach string
		var mediaType string
		c := config.Get()
		previewGif := filepath.Join(c.TempDir, "preview.gif")
		thumbnailJpg := filepath.Join(c.TempDir, "thumbnail.jpg")
		uniquenessString := shortid.MustGenerate()
		if utils.DoesFileExists(previewGif) {
			imageToAttach = "preview.gif"
			mediaType = "image/gif"
		} else if utils.DoesFileExists(thumbnailJpg) {
			imageToAttach = "thumbnail.jpg"
			mediaType = "image/jpeg"
		}
		if imageToAttach != "" {
			previewURL.Path = imageToAttach
			previewURL.RawQuery = "us=" + uniquenessString
			apmodels.AddImageAttachmentToNote(note, previewURL.String(), mediaType)
		}
	}

	if apo.configRepository.GetNSFW() {
		// Mark content as sensitive.
		sensitive := streams.NewActivityStreamsSensitiveProperty()
		sensitive.AppendXMLSchemaBoolean(true)
		note.SetActivityStreamsSensitive(sensitive)
	}

	b, err := apmodels.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize go live message activity", err)
		return errors.New("unable to serialize go live message activity " + err.Error())
	}

	if err := apo.SendToFollowers(b); err != nil {
		return err
	}

	if err := apo.Add(note, noteID, true); err != nil {
		return err
	}

	return nil
}

// SendDirectMessageToAccount will send a direct message to a single account.
func (apo *APOutbox) SendDirectMessageToAccount(textContent, account string) error {
	wf := webfinger.Get()

	links, err := wf.GetWebfingerLinks(account)
	if err != nil {
		return errors.Wrap(err, "unable to get webfinger links when sending private message")
	}
	user := apmodels.MakeWebFingerRequestResponseFromData(links)

	iri := user.Self
	actor, err := apo.resolvers.GetResolvedActorFromIRI(iri)
	if err != nil {
		return errors.Wrap(err, "unable to resolve actor to send message to")
	}

	activity, _, note, _ := apo.createBaseOutboundMessage(textContent)

	// Set direct message visibility
	activity = apmodels.MakeActivityDirect(activity, actor.ActorIri)
	note = apmodels.MakeNoteDirect(note, actor.ActorIri)
	object := activity.GetActivityStreamsObject()
	object.SetActivityStreamsNote(0, note)

	b, err := apmodels.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize custom fediverse message activity", err)
		return errors.Wrap(err, "unable to serialize custom fediverse message activity")
	}

	return apo.SendToUser(actor.Inbox, b)
}

// SendPublicMessage will send a public message to all followers.
func (apo *APOutbox) SendPublicMessage(textContent string) error {
	originalContent := textContent
	textContent = utils.RenderSimpleMarkdown(textContent)

	tagProp := streams.NewActivityStreamsTagProperty()

	hashtagStrings := utils.GetHashtagsFromText(originalContent)

	for _, hashtag := range hashtagStrings {
		tagWithoutHashtag := strings.TrimPrefix(hashtag, "#")

		// Replace the instances of the tag with a link to the tag page.
		tagHTML := apo.getHashtagLinkHTMLFromTagString(tagWithoutHashtag)
		textContent = strings.ReplaceAll(textContent, hashtag, tagHTML)

		// Create Hashtag object for the tag.
		hashtag := apmodels.MakeHashtag(tagWithoutHashtag)
		tagProp.AppendTootHashtag(hashtag)
	}

	activity, _, note, noteID := apo.createBaseOutboundMessage(textContent)
	note.SetActivityStreamsTag(tagProp)

	if !apo.configRepository.GetFederationIsPrivate() {
		note = apmodels.MakeNotePublic(note)
		activity = apmodels.MakeActivityPublic(activity)
	}

	b, err := apmodels.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize custom fediverse message activity", err)
		return errors.New("unable to serialize custom fediverse message activity " + err.Error())
	}

	if err := apo.SendToFollowers(b); err != nil {
		return err
	}

	if err := apo.Add(note, noteID, false); err != nil {
		return err
	}

	return nil
}

// nolint: unparam
func (apo *APOutbox) createBaseOutboundMessage(textContent string) (vocab.ActivityStreamsCreate, string, vocab.ActivityStreamsNote, string) {
	localActor := apmodels.MakeLocalIRIForAccount(apo.configRepository.GetDefaultFederationUsername())
	noteID := shortid.MustGenerate()
	noteIRI := apmodels.MakeLocalIRIForResource(noteID)
	id := shortid.MustGenerate()
	activity := apmodels.CreateCreateActivity(id, localActor)
	object := streams.NewActivityStreamsObjectProperty()
	activity.SetActivityStreamsObject(object)

	note := apmodels.MakeNote(textContent, noteIRI, localActor)
	object.AppendActivityStreamsNote(note)

	return activity, id, note, noteID
}

// Get Hashtag HTML link for a given tag (without # prefix).
func (apo *APOutbox) getHashtagLinkHTMLFromTagString(baseHashtag string) string {
	return fmt.Sprintf("<a class=\"hashtag\" href=\"https://directory.owncast.online/tags/%s\">#%s</a>", baseHashtag, baseHashtag)
}

// SendToFollowers will send an arbitrary payload to all follower inboxes.
func (apo *APOutbox) SendToFollowers(payload []byte) error {
	localActor := apmodels.MakeLocalIRIForAccount(apo.configRepository.GetDefaultFederationUsername())

	followers, _, err := apo.federationRepository.GetFederationFollowers(-1, 0)
	if err != nil {
		log.Errorln("unable to fetch followers to send to", err)
		return errors.New("unable to fetch followers to send payload to")
	}

	for _, follower := range followers {
		inbox, _ := url.Parse(follower.Inbox)
		req, err := crypto.CreateSignedRequest(payload, inbox, localActor)
		if err != nil {
			log.Errorln("unable to create outbox request", follower.Inbox, err)
			return errors.New("unable to create outbox request: " + follower.Inbox)
		}

		apo.workerpool.AddToOutboundQueue(req)
	}
	return nil
}

// SendToUser will send a payload to a single specific inbox.
func (apo *APOutbox) SendToUser(inbox *url.URL, payload []byte) error {
	localActor := apmodels.MakeLocalIRIForAccount(apo.configRepository.GetDefaultFederationUsername())

	req, err := apo.requests.CreateSignedRequest(payload, inbox, localActor)
	if err != nil {
		return errors.Wrap(err, "unable to create outbox request")
	}

	apo.workerpool.AddToOutboundQueue(req)

	return nil
}

// UpdateFollowersWithAccountUpdates will send an update to all followers alerting of a profile update.
func (apo *APOutbox) UpdateFollowersWithAccountUpdates() error {
	// Don't do anything if federation is disabled.
	if !apo.configRepository.GetFederationEnabled() {
		return nil
	}

	id := shortid.MustGenerate()
	objectID := apmodels.MakeLocalIRIForResource(id)
	activity := apmodels.MakeUpdateActivity(objectID)

	actor := streams.NewActivityStreamsPerson()
	actorID := apmodels.MakeLocalIRIForAccount(apo.configRepository.GetDefaultFederationUsername())
	actorIDProperty := streams.NewJSONLDIdProperty()
	actorIDProperty.Set(actorID)
	actor.SetJSONLDId(actorIDProperty)

	actorProperty := streams.NewActivityStreamsActorProperty()
	actorProperty.AppendActivityStreamsPerson(actor)
	activity.SetActivityStreamsActor(actorProperty)

	obj := streams.NewActivityStreamsObjectProperty()
	obj.AppendIRI(actorID)
	activity.SetActivityStreamsObject(obj)

	b, err := apmodels.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize send update actor activity", err)
		return errors.New("unable to serialize send update actor activity")
	}
	return apo.SendToFollowers(b)
}

// Add will save an ActivityPub object to the datastore.
func (apo *APOutbox) Add(item vocab.Type, id string, isLiveNotification bool) error {
	iri := item.GetJSONLDId().GetIRI().String()
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

	return apo.federationRepository.AddToOutbox(iri, b, typeString, isLiveNotification)
}
