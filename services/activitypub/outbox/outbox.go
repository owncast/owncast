// Package outbox is the federated-message producer for the ActivityPub
// subsystem: builds Create activities, addresses them to followers, and
// hands the signed bytes to the outbound worker pool. Construct via
// New(Deps) with persistence + workerpool services; all entry points
// are methods on *Service.
package outbox

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	apcrypto "github.com/owncast/owncast/services/activitypub/crypto"
	"github.com/owncast/owncast/services/activitypub/persistence"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/services/activitypub/requests"
	apresolvers "github.com/owncast/owncast/services/activitypub/resolvers"
	"github.com/owncast/owncast/services/activitypub/webfinger"
	"github.com/owncast/owncast/services/activitypub/workerpool"
	"github.com/owncast/owncast/utils"
)

// Service is the outbound side of federation. It composes a *persistence.Service
// for storing what we send and a *workerpool.Service for actually
// delivering it.
type Service struct {
	persistence      *persistence.Service
	workerpool       *workerpool.Service
	followers        followersrepository.FollowersRepository
	configRepository configrepository.ConfigRepository
	builder          *apmodels.Builder
	signer           *apcrypto.Signer
	resolver         *apresolvers.Resolver
	cfg              *config.Config
}

// Deps is the explicit dependency contract for outbox.
type Deps struct {
	Persistence      *persistence.Service
	Workerpool       *workerpool.Service
	Followers        followersrepository.FollowersRepository
	ConfigRepository configrepository.ConfigRepository
	Builder          *apmodels.Builder
	Signer           *apcrypto.Signer
	Resolver         *apresolvers.Resolver
	Config           *config.Config
}

// New constructs an outbox Service. All deps are required.
func New(deps Deps) *Service {
	return &Service{
		persistence:      deps.Persistence,
		workerpool:       deps.Workerpool,
		followers:        deps.Followers,
		configRepository: deps.ConfigRepository,
		builder:          deps.Builder,
		signer:           deps.Signer,
		resolver:         deps.Resolver,
		cfg:              deps.Config,
	}
}

// SendLive notifies all followers that the stream has gone live.
func (s *Service) SendLive() error {
	textContent := s.configRepository.GetFederationGoLiveMessage()

	// If the message is empty then do not send it.
	if textContent == "" {
		return nil
	}

	tagStrings := []string{}
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")

	tagProp := streams.NewActivityStreamsTagProperty()
	for _, tagString := range s.configRepository.GetServerMetadataTags() {
		tagWithoutSpecialCharacters := reg.ReplaceAllString(tagString, "")
		hashtag := apmodels.MakeHashtag(tagWithoutSpecialCharacters)
		tagProp.AppendTootHashtag(hashtag)
		tagString := getHashtagLinkHTMLFromTagString(tagWithoutSpecialCharacters)
		tagStrings = append(tagStrings, tagString)
	}

	// Manually add Owncast hashtag if it doesn't already exist so it
	// shows up in Owncast search results.
	if _, exists := utils.FindInSlice(tagStrings, "owncast"); !exists {
		hashtag := apmodels.MakeHashtag("owncast")
		tagProp.AppendTootHashtag(hashtag)
	}

	tagsString := strings.Join(tagStrings, " ")

	var streamTitle string
	if title := s.configRepository.GetStreamTitle(); title != "" {
		streamTitle = fmt.Sprintf("<p>%s</p>", title)
	}
	textContent = fmt.Sprintf("<p>%s</p>%s<p>%s</p><p><a href=\"%s\">%s</a></p>", textContent, streamTitle, tagsString, s.configRepository.GetServerURL(), s.configRepository.GetServerURL())

	activity, _, note, noteID := s.createBaseOutboundMessage(textContent)

	to, cc := s.getAddressingToFollowers()
	note.SetActivityStreamsTo(to)
	note.SetActivityStreamsCc(cc)
	activity.SetActivityStreamsTo(to)
	activity.SetActivityStreamsCc(cc)

	note.SetActivityStreamsTag(tagProp)

	// Attach an image along with the Federated message.
	previewURL, err := url.Parse(s.configRepository.GetServerURL())
	if err == nil {
		var imageToAttach string
		var mediaType string
		previewGif := filepath.Join(s.cfg.TempDir, "preview.gif")
		thumbnailJpg := filepath.Join(s.cfg.TempDir, "thumbnail.jpg")
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

	if s.configRepository.GetNSFW() {
		sensitive := streams.NewActivityStreamsSensitiveProperty()
		sensitive.AppendXMLSchemaBoolean(true)
		note.SetActivityStreamsSensitive(sensitive)
	}

	b, err := apmodels.Serialize(activity)
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

// SendDirectMessageToAccount sends a direct message to a single account.
func (s *Service) SendDirectMessageToAccount(textContent, account string) error {
	links, err := webfinger.GetWebfingerLinks(account)
	if err != nil {
		return errors.Wrap(err, "unable to get webfinger links when sending private message")
	}
	user := apmodels.MakeWebFingerRequestResponseFromData(links)

	iri := user.Self
	actor, err := s.resolver.GetResolvedActorFromIRI(iri)
	if err != nil {
		return errors.Wrap(err, "unable to resolve actor to send message to")
	}

	activity, _, note, _ := s.createBaseOutboundMessage(textContent)

	// Set direct message visibility.
	activity = apmodels.MakeActivityDirect(activity, actor.ActorIri)
	note = apmodels.MakeNoteDirect(note, actor.ActorIri)
	object := activity.GetActivityStreamsObject()
	object.SetActivityStreamsNote(0, note)

	b, err := apmodels.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize custom fediverse message activity", err)
		return errors.Wrap(err, "unable to serialize custom fediverse message activity")
	}

	return s.SendToUser(actor.Inbox, b)
}

// SendPublicMessage sends a public message to all followers.
func (s *Service) SendPublicMessage(textContent string) error {
	originalContent := textContent
	textContent = utils.RenderSimpleMarkdown(textContent)

	tagProp := streams.NewActivityStreamsTagProperty()

	hashtagStrings := utils.GetHashtagsFromText(originalContent)

	for _, hashtag := range hashtagStrings {
		tagWithoutHashtag := strings.TrimPrefix(hashtag, "#")

		// Replace each instance with a link to the tag page.
		tagHTML := getHashtagLinkHTMLFromTagString(tagWithoutHashtag)
		textContent = strings.ReplaceAll(textContent, hashtag, tagHTML)

		hashtag := apmodels.MakeHashtag(tagWithoutHashtag)
		tagProp.AppendTootHashtag(hashtag)
	}

	activity, _, note, noteID := s.createBaseOutboundMessage(textContent)
	note.SetActivityStreamsTag(tagProp)

	to, cc := s.getAddressingToFollowers()
	note.SetActivityStreamsTo(to)
	note.SetActivityStreamsCc(cc)
	activity.SetActivityStreamsTo(to)
	activity.SetActivityStreamsCc(cc)

	b, err := apmodels.Serialize(activity)
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

// getAddressingToFollowers builds the to/cc properties for a follower-addressed
// activity: public → cc:followers, to:public; private → to:followers.
func (s *Service) getAddressingToFollowers() (vocab.ActivityStreamsToProperty, vocab.ActivityStreamsCcProperty) {
	username := s.configRepository.GetDefaultFederationUsername()

	followersIRI := s.builder.MakeLocalIRIForAccount(username)
	followersIRI = followersIRI.JoinPath("followers")

	return apmodels.MakeAddressingToFollowers(followersIRI, !s.configRepository.GetFederationIsPrivate())
}

// nolint: unparam
func (s *Service) createBaseOutboundMessage(textContent string) (vocab.ActivityStreamsCreate, string, vocab.ActivityStreamsNote, string) {
	localActor := s.builder.MakeLocalIRIForAccount(s.configRepository.GetDefaultFederationUsername())
	noteID := shortid.MustGenerate()
	noteIRI := s.builder.MakeLocalIRIForResource(noteID)
	id := shortid.MustGenerate()
	activity := s.builder.CreateCreateActivity(id, localActor)
	object := streams.NewActivityStreamsObjectProperty()
	activity.SetActivityStreamsObject(object)

	note := apmodels.MakeNote(textContent, noteIRI, localActor)
	object.AppendActivityStreamsNote(note)

	return activity, id, note, noteID
}

// getHashtagLinkHTMLFromTagString returns the HTML link for a tag (without the # prefix).
func getHashtagLinkHTMLFromTagString(baseHashtag string) string {
	return fmt.Sprintf("<a class=\"hashtag\" href=\"https://owncast.directory/tags/%s\">#%s</a>", baseHashtag, baseHashtag)
}

// SendToFollowers sends an arbitrary payload to all follower inboxes,
// preferring shared inboxes to reduce outbound request count.
func (s *Service) SendToFollowers(payload []byte) error {
	localActor := s.builder.MakeLocalIRIForAccount(s.configRepository.GetDefaultFederationUsername())

	// Prefer shared inboxes over individual inboxes.
	inboxes, err := s.followers.GetUniqueDeliveryInboxes()
	if err != nil {
		log.Errorln("unable to fetch delivery inboxes", err)
		return errors.New("unable to fetch delivery inboxes to send payload to")
	}

	// Batch size and delay to prevent resource exhaustion during
	// delivery; spreads CPU load from cryptographic signing over time.
	const batchSize = 50
	const batchDelay = 100 * time.Millisecond

	queued := 0
	skipped := 0

	for i, inboxURL := range inboxes {
		inbox, err := url.Parse(inboxURL)
		if err != nil {
			log.Warnln("unable to parse inbox URL", inboxURL, err)
			continue
		}

		// SSRF protection: reject non-HTTPS schemes and
		// internal/loopback hosts. A malicious remote actor could set
		// their inbox to an internal address to trick this server into
		// making requests to internal services.
		if inbox.Scheme != "https" {
			log.Warnln("rejecting non-HTTPS inbox URL for SSRF protection:", inboxURL)
			continue
		}
		if utils.IsHostnameInternal(inbox.Hostname()) {
			log.Warnln("rejecting internal/loopback inbox URL for SSRF protection:", inboxURL)
			continue
		}

		// Pre-check circuit breaker BEFORE expensive cryptographic
		// signing. This saves CPU cycles for domains we know are
		// failing.
		if s.workerpool.ShouldSkipDomain(inbox.Host) {
			skipped++
			continue
		}

		req, err := s.signer.CreateSignedRequest(payload, inbox, localActor)
		if err != nil {
			log.Errorln("unable to create outbox request", inboxURL, err)
			continue
		}

		s.workerpool.AddToOutboundQueue(req)
		queued++

		// Spread CPU and network load across batches so ActivityPub
		// delivery doesn't compete with video encoding. Use queued count
		// (not loop index) so rate limiting is consistent even when
		// followers are skipped due to circuit breaker or parse errors.
		if queued%batchSize == 0 && i+1 < len(inboxes) {
			time.Sleep(batchDelay)
		}
	}

	if skipped > 0 {
		log.Debugf("Skipped %d followers due to circuit breaker, queued %d", skipped, queued)
	}

	return nil
}

// SendToUser sends a payload to a single specific inbox.
func (s *Service) SendToUser(inbox *url.URL, payload []byte) error {
	// SSRF protection: reject non-HTTPS schemes and internal/loopback hosts.
	if inbox.Scheme != "https" {
		return errors.Errorf("rejecting non-HTTPS inbox URL for SSRF protection: %s", inbox.String())
	}
	if utils.IsHostnameInternal(inbox.Hostname()) {
		return errors.Errorf("rejecting internal/loopback inbox URL for SSRF protection: %s", inbox.String())
	}

	localActor := s.builder.MakeLocalIRIForAccount(s.configRepository.GetDefaultFederationUsername())

	req, err := requests.CreateSignedRequest(payload, inbox, localActor, s.signer)
	if err != nil {
		return errors.Wrap(err, "unable to create outbox request")
	}

	s.workerpool.AddToOutboundQueue(req)

	return nil
}

// UpdateFollowersWithAccountUpdates broadcasts a profile-update Activity
// to all followers.
func (s *Service) UpdateFollowersWithAccountUpdates() error {
	if !s.configRepository.GetFederationEnabled() {
		return nil
	}

	id := shortid.MustGenerate()
	objectID := s.builder.MakeLocalIRIForResource(id)
	activity := s.builder.MakeUpdateActivity(objectID)

	actor := streams.NewActivityStreamsPerson()
	actorID := s.builder.MakeLocalIRIForAccount(s.configRepository.GetDefaultFederationUsername())
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
	return s.SendToFollowers(b)
}

// Add saves an ActivityPub object to the datastore.
func (s *Service) Add(item vocab.Type, id string, isLiveNotification bool) error {
	iri, err := apmodels.GetIRIStringFromJSONLDIdProperty(item.GetJSONLDId())
	if err != nil {
		log.Errorln("Unable to get iri from item:", err)
		return errors.Wrap(err, "unable to get iri from item "+id)
	}
	typeString := item.GetTypeName()

	b, err := apmodels.Serialize(item)
	if err != nil {
		log.Errorln("unable to serialize model when saving to outbox", err)
		return err
	}

	return s.persistence.AddToOutbox(iri, b, typeString, isLiveNotification)
}
