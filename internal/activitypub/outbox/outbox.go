package outbox

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/internal/activitypub/apmodels"
	"github.com/owncast/owncast/internal/activitypub/crypto"
	"github.com/owncast/owncast/internal/activitypub/follower"
	"github.com/owncast/owncast/internal/activitypub/persistence"
	"github.com/owncast/owncast/internal/activitypub/webfinger"
	"github.com/owncast/owncast/internal/activitypub/workerpool"
	"github.com/owncast/owncast/utils"
)

func New(
	p *persistence.Service,
	c *crypto.Service,
	m *apmodels.Service,
	f *follower.Service,
	w *workerpool.Service,
) (*Service, error) {
	s := &Service{}

	s.persistence = p
	s.crypto = c
	s.models = m
	s.follower = f
	s.worker = w

	return s, nil
}

type Service struct {
	persistence *persistence.Service
	crypto      *crypto.Service
	models      *apmodels.Service
	follower    *follower.Service
	worker      *workerpool.Service
}

// SendLive will send all followers the message saying you started a live stream.
func (s *Service) SendLive() error {
	textContent := s.persistence.Data.GetFederationGoLiveMessage()

	// If the message is empty then do not send it.
	if textContent == "" {
		return nil
	}

	tagStrings := []string{}
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")

	tagProp := streams.NewActivityStreamsTagProperty()
	for _, tagString := range s.persistence.Data.GetServerMetadataTags() {
		tagWithoutSpecialCharacters := reg.ReplaceAllString(tagString, "")
		hashtag := s.models.MakeHashtag(tagWithoutSpecialCharacters)
		tagProp.AppendTootHashtag(hashtag)
		tagString := s.getHashtagLinkHTMLFromTagString(tagWithoutSpecialCharacters)
		tagStrings = append(tagStrings, tagString)
	}

	// Manually add Owncast hashtag if it doesn't already exist so it shows up
	// in Owncast search results.
	// We can remove this down the road, but it'll be nice for now.
	if _, exists := utils.FindInSlice(tagStrings, "owncast"); !exists {
		hashtag := s.models.MakeHashtag("owncast")
		tagProp.AppendTootHashtag(hashtag)
	}

	tagsString := strings.Join(tagStrings, " ")

	var streamTitle string
	if title := s.persistence.Data.GetStreamTitle(); title != "" {
		streamTitle = fmt.Sprintf("<p>%s</p>", title)
	}
	textContent = fmt.Sprintf("<p>%s</p>%s<p>%s</p><a href=\"%s\">%s</a>", textContent, streamTitle, tagsString, s.persistence.Data.GetServerURL(), s.persistence.Data.GetServerURL())

	activity, _, note, noteID := s.createBaseOutboundMessage(textContent)

	// To the public if we're not treating ActivityPub as "private".
	if !s.persistence.Data.GetFederationIsPrivate() {
		note = s.models.MakeNotePublic(note)
		activity = s.models.MakeActivityPublic(activity)
	}

	note.SetActivityStreamsTag(tagProp)

	// Attach an image along with the Federated message.
	previewURL, err := url.Parse(s.persistence.Data.GetServerURL())
	if err == nil {
		var imageToAttach string
		var mediaType string
		previewGif := filepath.Join(config.TempDir, "preview.gif")
		thumbnailJpg := filepath.Join(config.TempDir, "thumbnail.jpg")
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
			s.models.AddImageAttachmentToNote(note, previewURL.String(), mediaType)
		}
	}

	if s.persistence.Data.GetNSFW() {
		// Mark content as sensitive.
		sensitive := streams.NewActivityStreamsSensitiveProperty()
		sensitive.AppendXMLSchemaBoolean(true)
		note.SetActivityStreamsSensitive(sensitive)
	}

	b, err := s.models.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize go live message activity", err)
		return errors.New("unable to serialize go live message activity " + err.Error())
	}

	if err := s.SendToFollowers(b); err != nil {
		return err
	}

	if err := s.Add(note, noteID, true); err != nil {
		return err
	}

	return nil
}

// SendDirectMessageToAccount will send a direct message to a single account.
func (s *Service) SendDirectMessageToAccount(textContent, account string) error {
	links, err := webfinger.GetWebfingerLinks(account)
	if err != nil {
		return errors.Wrap(err, "unable to get webfinger links when sending private message")
	}
	user := s.models.MakeWebFingerRequestResponseFromData(links)

	iri := user.Self
	actor, err := s.follower.GetResolvedActorFromIRI(iri)
	if err != nil {
		return errors.Wrap(err, "unable to resolve actor to send message to")
	}

	activity, _, note, _ := s.createBaseOutboundMessage(textContent)

	// Set direct message visibility
	activity = s.models.MakeActivityDirect(activity, actor.ActorIri)
	note = s.models.MakeNoteDirect(note, actor.ActorIri)
	object := activity.GetActivityStreamsObject()
	object.SetActivityStreamsNote(0, note)

	b, err := s.models.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize custom fediverse message activity", err)
		return errors.Wrap(err, "unable to serialize custom fediverse message activity")
	}

	return s.SendToUser(actor.Inbox, b)
}

// SendPublicMessage will send a public message to all followers.
func (s *Service) SendPublicMessage(textContent string, p *persistence.Service) error {
	originalContent := textContent
	textContent = utils.RenderSimpleMarkdown(textContent)

	tagProp := streams.NewActivityStreamsTagProperty()

	hashtagStrings := utils.GetHashtagsFromText(originalContent)

	for _, hashtag := range hashtagStrings {
		tagWithoutHashtag := strings.TrimPrefix(hashtag, "#")

		// Replace the instances of the tag with a link to the tag page.
		tagHTML := s.getHashtagLinkHTMLFromTagString(tagWithoutHashtag)
		textContent = strings.ReplaceAll(textContent, hashtag, tagHTML)

		// Create Hashtag object for the tag.
		hashtag := s.models.MakeHashtag(tagWithoutHashtag)
		tagProp.AppendTootHashtag(hashtag)
	}

	activity, _, note, noteID := s.createBaseOutboundMessage(textContent)
	note.SetActivityStreamsTag(tagProp)

	if !s.persistence.Data.GetFederationIsPrivate() {
		note = s.models.MakeNotePublic(note)
		activity = s.models.MakeActivityPublic(activity)
	}

	b, err := s.models.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize custom fediverse message activity", err)
		return errors.New("unable to serialize custom fediverse message activity " + err.Error())
	}

	if err := s.SendToFollowers(b); err != nil {
		return err
	}

	if err := s.Add(note, noteID, false); err != nil {
		return err
	}

	return nil
}

// nolint: unparam
func (s *Service) createBaseOutboundMessage(textContent string) (vocab.ActivityStreamsCreate, string, vocab.ActivityStreamsNote, string) {
	localActor := s.models.MakeLocalIRIForAccount(s.persistence.Data.GetDefaultFederationUsername())
	noteID := shortid.MustGenerate()
	noteIRI := s.models.MakeLocalIRIForResource(noteID)
	id := shortid.MustGenerate()
	activity := s.models.CreateCreateActivity(id, localActor)
	object := streams.NewActivityStreamsObjectProperty()
	activity.SetActivityStreamsObject(object)

	note := s.models.MakeNote(textContent, noteIRI, localActor)
	object.AppendActivityStreamsNote(note)

	return activity, id, note, noteID
}

// Get Hashtag HTML link for a given tag (without # prefix).
func (s *Service) getHashtagLinkHTMLFromTagString(baseHashtag string) string {
	return fmt.Sprintf("<a class=\"hashtag\" href=\"https://directory.owncast.online/tags/%s\">#%s</a>", baseHashtag, baseHashtag)
}

// SendToFollowers will send an arbitrary payload to all follower inboxes.
func (s *Service) SendToFollowers(payload []byte) error {
	localActor := s.models.MakeLocalIRIForAccount(s.persistence.Data.GetDefaultFederationUsername())

	followers, _, err := s.follower.GetFederationFollowers(-1, 0)
	if err != nil {
		log.Errorln("unable to fetch followers to send to", err)
		return errors.New("unable to fetch followers to send payload to")
	}

	for _, follower := range followers {
		inbox, _ := url.Parse(follower.Inbox)
		req, err := s.crypto.CreateSignedRequest(payload, inbox, localActor)
		if err != nil {
			log.Errorln("unable to create outbox request", follower.Inbox, err)
			return errors.New("unable to create outbox request: " + follower.Inbox)
		}

		s.worker.AddToOutboundQueue(req)
	}
	return nil
}

// SendToUser will send a payload to a single specific inbox.
func (s *Service) SendToUser(inbox *url.URL, payload []byte) error {
	localActor := s.models.MakeLocalIRIForAccount(s.persistence.Data.GetDefaultFederationUsername())

	req, err := follower.CreateSignedRequest(payload, inbox, localActor, s.crypto)
	if err != nil {
		return errors.Wrap(err, "unable to create outbox request")
	}

	s.worker.AddToOutboundQueue(req)

	return nil
}

// UpdateFollowersWithAccountUpdates will send an update to all followers alerting of a profile update.
func (s *Service) UpdateFollowersWithAccountUpdates(p *persistence.Service) error {
	// Don't do anything if federation is disabled.
	if !s.persistence.Data.GetFederationEnabled() {
		return nil
	}

	id := shortid.MustGenerate()
	objectID := s.models.MakeLocalIRIForResource(id)
	activity := s.models.MakeUpdateActivity(objectID)

	actor := streams.NewActivityStreamsPerson()
	actorID := s.models.MakeLocalIRIForAccount(s.persistence.Data.GetDefaultFederationUsername())
	actorIDProperty := streams.NewJSONLDIdProperty()
	actorIDProperty.Set(actorID)
	actor.SetJSONLDId(actorIDProperty)

	actorProperty := streams.NewActivityStreamsActorProperty()
	actorProperty.AppendActivityStreamsPerson(actor)
	activity.SetActivityStreamsActor(actorProperty)

	obj := streams.NewActivityStreamsObjectProperty()
	obj.AppendIRI(actorID)
	activity.SetActivityStreamsObject(obj)

	b, err := s.models.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize send update actor activity", err)
		return errors.New("unable to serialize send update actor activity")
	}

	return s.SendToFollowers(b)
}

// Add will save an ActivityPub object to the datastore.
func (s *Service) Add(item vocab.Type, id string, isLiveNotification bool) error {
	iri := item.GetJSONLDId().GetIRI().String()
	typeString := item.GetTypeName()

	if iri == "" {
		log.Errorln("Unable to get iri from item")
		return errors.New("Unable to get iri from item " + id)
	}

	b, err := s.models.Serialize(item)
	if err != nil {
		log.Errorln("unable to serialize model when saving to outbox", err)
		return err
	}

	return s.persistence.AddToOutbox(iri, b, typeString, isLiveNotification)
}
