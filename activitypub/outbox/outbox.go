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
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/activitypub/webfinger"
	"github.com/owncast/owncast/activitypub/workerpool"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/pkg/errors"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

// SendLive will send all followers the message saying you started a live stream.
func SendLive() error {
	configRepository := configrepository.Get()

	textContent := configRepository.GetFederationGoLiveMessage()

	// If the message is empty then do not send it.
	if textContent == "" {
		return nil
	}

	tagStrings := []string{}
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")

	tagProp := streams.NewActivityStreamsTagProperty()
	for _, tagString := range configRepository.GetServerMetadataTags() {
		tagWithoutSpecialCharacters := reg.ReplaceAllString(tagString, "")
		hashtag := apmodels.MakeHashtag(tagWithoutSpecialCharacters)
		tagProp.AppendTootHashtag(hashtag)
		tagString := getHashtagLinkHTMLFromTagString(tagWithoutSpecialCharacters)
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
	if title := configRepository.GetStreamTitle(); title != "" {
		streamTitle = fmt.Sprintf("<p>%s</p>", title)
	}
	textContent = fmt.Sprintf("<p>%s</p>%s<p>%s</p><p><a href=\"%s\">%s</a></p>", textContent, streamTitle, tagsString, configRepository.GetServerURL(), configRepository.GetServerURL())

	activity, _, note, noteID := createBaseOutboundMessage(textContent)

	to, cc := getAddressingToFollowers()
	note.SetActivityStreamsTo(to)
	note.SetActivityStreamsCc(cc)
	activity.SetActivityStreamsTo(to)
	activity.SetActivityStreamsCc(cc)

	note.SetActivityStreamsTag(tagProp)

	// Attach an image along with the Federated message.
	previewURL, err := url.Parse(configRepository.GetServerURL())
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
			apmodels.AddImageAttachmentToNote(note, previewURL.String(), mediaType)
		}
	}

	if configRepository.GetNSFW() {
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

	if err := SendToFollowers(b); err != nil {
		return err
	}

	if err := Add(note, noteID, true); err != nil {
		return err
	}

	return nil
}

// SendDirectMessageToAccount will send a direct message to a single account.
func SendDirectMessageToAccount(textContent, account string) error {
	links, err := webfinger.GetWebfingerLinks(account)
	if err != nil {
		return errors.Wrap(err, "unable to get webfinger links when sending private message")
	}
	user := apmodels.MakeWebFingerRequestResponseFromData(links)

	iri := user.Self
	actor, err := resolvers.GetResolvedActorFromIRI(iri)
	if err != nil {
		return errors.Wrap(err, "unable to resolve actor to send message to")
	}

	activity, _, note, _ := createBaseOutboundMessage(textContent)

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

	return SendToUser(actor.Inbox, b)
}

// SendPublicMessage will send a public message to all followers.
func SendPublicMessage(textContent string) error {
	originalContent := textContent
	textContent = utils.RenderSimpleMarkdown(textContent)

	tagProp := streams.NewActivityStreamsTagProperty()

	hashtagStrings := utils.GetHashtagsFromText(originalContent)

	for _, hashtag := range hashtagStrings {
		tagWithoutHashtag := strings.TrimPrefix(hashtag, "#")

		// Replace the instances of the tag with a link to the tag page.
		tagHTML := getHashtagLinkHTMLFromTagString(tagWithoutHashtag)
		textContent = strings.ReplaceAll(textContent, hashtag, tagHTML)

		// Create Hashtag object for the tag.
		hashtag := apmodels.MakeHashtag(tagWithoutHashtag)
		tagProp.AppendTootHashtag(hashtag)
	}

	activity, _, note, noteID := createBaseOutboundMessage(textContent)
	note.SetActivityStreamsTag(tagProp)

	to, cc := getAddressingToFollowers()
	note.SetActivityStreamsTo(to)
	note.SetActivityStreamsCc(cc)
	activity.SetActivityStreamsTo(to)
	activity.SetActivityStreamsCc(cc)

	b, err := apmodels.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize custom fediverse message activity", err)
		return errors.New("unable to serialize custom fediverse message activity " + err.Error())
	}

	if err := SendToFollowers(b); err != nil {
		return err
	}

	if err := Add(note, noteID, false); err != nil {
		return err
	}

	return nil
}

// if public, cc the followers and to the Public uri, else private, address followers directly.
func getAddressingToFollowers() (vocab.ActivityStreamsToProperty, vocab.ActivityStreamsCcProperty) {
	configRepository := configrepository.Get()
	server_url := configRepository.GetServerURL()
	followers_iri, _ := url.Parse(server_url)
	username := configRepository.GetDefaultFederationUsername()

	followers_iri = followers_iri.JoinPath("federation", "user", username, "followers")

	return apmodels.MakeAddressingToFollowers(followers_iri, !configRepository.GetFederationIsPrivate())
}

// nolint: unparam
func createBaseOutboundMessage(textContent string) (vocab.ActivityStreamsCreate, string, vocab.ActivityStreamsNote, string) {
	configRepository := configrepository.Get()
	localActor := apmodels.MakeLocalIRIForAccount(configRepository.GetDefaultFederationUsername())
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
func getHashtagLinkHTMLFromTagString(baseHashtag string) string {
	return fmt.Sprintf("<a class=\"hashtag\" href=\"https://owncast.directory/tags/%s\">#%s</a>", baseHashtag, baseHashtag)
}

// SendToFollowers will send an arbitrary payload to all follower inboxes.
// It uses shared inboxes when available to reduce the number of outbound requests.
func SendToFollowers(payload []byte) error {
	configRepository := configrepository.Get()
	followersRepo := followersrepository.Get()
	localActor := apmodels.MakeLocalIRIForAccount(configRepository.GetDefaultFederationUsername())

	// Get unique delivery inboxes (prefers shared inboxes over individual inboxes)
	inboxes, err := followersRepo.GetUniqueDeliveryInboxes()
	if err != nil {
		log.Errorln("unable to fetch delivery inboxes", err)
		return errors.New("unable to fetch delivery inboxes to send payload to")
	}

	// Batch size and delay to prevent resource exhaustion during delivery.
	// This spreads CPU load from cryptographic signing over time.
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

		// SSRF protection: reject non-HTTPS schemes and internal/loopback hosts.
		// A malicious remote actor could set their inbox to an internal address
		// to trick this server into making requests to internal services.
		if inbox.Scheme != "https" {
			log.Warnln("rejecting non-HTTPS inbox URL for SSRF protection:", inboxURL)
			continue
		}
		if utils.IsHostnameInternal(inbox.Hostname()) {
			log.Warnln("rejecting internal/loopback inbox URL for SSRF protection:", inboxURL)
			continue
		}

		// Pre-check circuit breaker BEFORE expensive cryptographic signing.
		// This saves CPU cycles for domains we know are failing.
		if workerpool.ShouldSkipDomain(inbox.Host) {
			skipped++
			continue
		}

		req, err := crypto.CreateSignedRequest(payload, inbox, localActor)
		if err != nil {
			log.Errorln("unable to create outbox request", inboxURL, err)
			continue
		}

		workerpool.AddToOutboundQueue(req)
		queued++

		// Add a small delay between batches to spread out CPU and network load.
		// This helps prevent ActivityPub delivery from competing with video encoding.
		// Use queued count (not loop index) to ensure consistent rate limiting
		// even when followers are skipped due to circuit breaker or parse errors.
		if queued%batchSize == 0 && i+1 < len(inboxes) {
			time.Sleep(batchDelay)
		}
	}

	if skipped > 0 {
		log.Debugf("Skipped %d followers due to circuit breaker, queued %d", skipped, queued)
	}

	return nil
}

// SendToUser will send a payload to a single specific inbox.
func SendToUser(inbox *url.URL, payload []byte) error {
	// SSRF protection: reject non-HTTPS schemes and internal/loopback hosts.
	if inbox.Scheme != "https" {
		return errors.Errorf("rejecting non-HTTPS inbox URL for SSRF protection: %s", inbox.String())
	}
	if utils.IsHostnameInternal(inbox.Hostname()) {
		return errors.Errorf("rejecting internal/loopback inbox URL for SSRF protection: %s", inbox.String())
	}

	configRepository := configrepository.Get()
	localActor := apmodels.MakeLocalIRIForAccount(configRepository.GetDefaultFederationUsername())

	req, err := requests.CreateSignedRequest(payload, inbox, localActor)
	if err != nil {
		return errors.Wrap(err, "unable to create outbox request")
	}

	workerpool.AddToOutboundQueue(req)

	return nil
}

// UpdateFollowersWithAccountUpdates will send an update to all followers alerting of a profile update.
func UpdateFollowersWithAccountUpdates() error {
	configRepository := configrepository.Get()

	// Don't do anything if federation is disabled.
	if !configRepository.GetFederationEnabled() {
		return nil
	}

	id := shortid.MustGenerate()
	objectID := apmodels.MakeLocalIRIForResource(id)
	activity := apmodels.MakeUpdateActivity(objectID)

	actor := streams.NewActivityStreamsPerson()
	actorID := apmodels.MakeLocalIRIForAccount(configRepository.GetDefaultFederationUsername())
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
	return SendToFollowers(b)
}

// Add will save an ActivityPub object to the datastore.
func Add(item vocab.Type, id string, isLiveNotification bool) error {
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

	return persistence.AddToOutbox(iri, b, typeString, isLiveNotification)
}
