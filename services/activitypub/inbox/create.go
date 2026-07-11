package inbox

import (
	"context"
	"database/sql"
	"net/url"
	"sort"
	"time"

	"code.superseriousbusiness.org/activity/streams/vocab"
	"github.com/pkg/errors"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/events"
	"github.com/owncast/owncast/utils"
)

func (s *Service) handleCreateRequest(c context.Context, activity vocab.ActivityStreamsCreate) error {
	validated, err := validateCreateNote(activity)
	if err != nil || validated == nil {
		return err
	}

	eventType, inReplyTo, relevant, err := s.classifyCreateNote(activity, validated.note)
	if err != nil || !relevant {
		return err
	}

	actor, err := s.resolver.GetResolvedActorFromActorProperty(activity.GetActivityStreamsActor())
	if err != nil {
		return errors.Wrap(err, "unable to resolve create actor")
	}
	if actor.ActorIriString() != validated.actorIRI {
		return errors.New("resolved create actor does not match activity actor")
	}

	payload := createPostPayload(validated.note, validated.noteIRI, inReplyTo, actor)
	claimed, err := s.persistence.ClaimInboundFediverseActivity(validated.createIRI, validated.actorIRI, activity.GetTypeName(), time.Now())
	if err != nil {
		return errors.Wrap(err, "unable to claim inbound create activity")
	}
	if !claimed {
		return nil
	}
	s.publishFediverseEvent(c, eventType, payload)
	return nil
}

type validatedCreate struct {
	createIRI string
	actorIRI  string
	noteIRI   string
	note      vocab.ActivityStreamsNote
}

func validateCreateNote(activity vocab.ActivityStreamsCreate) (*validatedCreate, error) {
	createIRI, err := apmodels.GetIRIStringFromJSONLDIdProperty(activity.GetJSONLDId())
	if err != nil {
		return nil, errors.Wrap(err, "create activity is missing IRI")
	}

	actorIRI, err := activityActorIRI(activity.GetActivityStreamsActor())
	if err != nil {
		return nil, errors.Wrap(err, "create activity has invalid actor")
	}

	object := activity.GetActivityStreamsObject()
	if object == nil || object.Len() != 1 || object.At(0) == nil || !object.At(0).IsActivityStreamsNote() {
		return nil, nil
	}
	note := object.At(0).GetActivityStreamsNote()

	noteIRI, err := apmodels.GetIRIStringFromJSONLDIdProperty(note.GetJSONLDId())
	if err != nil {
		return nil, errors.Wrap(err, "create note is missing IRI")
	}
	if err := validateNoteAttribution(note.GetActivityStreamsAttributedTo(), actorIRI); err != nil {
		return nil, err
	}

	return &validatedCreate{
		createIRI: createIRI,
		actorIRI:  actorIRI,
		noteIRI:   noteIRI,
		note:      note,
	}, nil
}

func (s *Service) classifyCreateNote(activity vocab.ActivityStreamsCreate, note vocab.ActivityStreamsNote) (models.EventType, string, bool, error) {
	inReplyTo := firstInReplyToIRI(note.GetActivityStreamsInReplyTo())
	if inReplyTo != "" {
		_, _, _, err := s.persistence.GetNoteByIRI(inReplyTo)
		if err == nil {
			return models.FediverseReply, inReplyTo, true, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return "", "", false, errors.Wrap(err, "unable to look up create note parent")
		}
	}

	localActor := s.builder.MakeLocalIRIForAccount(s.configRepository.GetFederationUsername())
	if localActor == nil || !createOrNoteAddresses(activity, note, localActor.String()) {
		return "", "", false, nil
	}
	return models.FediverseMention, inReplyTo, true, nil
}

func createPostPayload(note vocab.ActivityStreamsNote, noteIRI, inReplyTo string, actor apmodels.ActivityPubActor) *events.FediverseInboundPostEvent {
	content, language := noteContent(note.GetActivityStreamsContent())
	postedAt := time.Now().UTC()
	if published := note.GetActivityStreamsPublished(); published != nil && published.IsXMLSchemaDateTime() {
		postedAt = published.Get().UTC()
	}

	postURL := noteIRI
	if urlProperty := note.GetActivityStreamsUrl(); urlProperty != nil && urlProperty.Len() > 0 {
		if first := urlProperty.At(0); first != nil && first.IsIRI() && isSafeHTTPURL(first.GetIRI()) {
			postURL = first.GetIRI().String()
		}
	}

	return &events.FediverseInboundPostEvent{
		Actor:       fediverseActorFromResolvedActor(actor),
		Content:     content,
		ContentText: utils.StripHTML(content),
		URL:         postURL,
		PostedAt:    postedAt.Format(time.RFC3339),
		InReplyTo:   inReplyTo,
		Attachments: noteAttachments(note.GetActivityStreamsAttachment()),
		Language:    language,
	}
}

func activityActorIRI(actor vocab.ActivityStreamsActorProperty) (string, error) {
	if actor == nil || actor.Len() != 1 || actor.At(0) == nil {
		return "", errors.New("actor must contain exactly one value")
	}
	value := actor.At(0)
	if value.IsIRI() && value.GetIRI() != nil {
		return value.GetIRI().String(), nil
	}
	if object := value.GetType(); object != nil {
		return apmodels.GetIRIStringFromJSONLDIdProperty(object.GetJSONLDId())
	}
	return "", errors.New("actor is missing IRI")
}

func validateNoteAttribution(attributedTo vocab.ActivityStreamsAttributedToProperty, actorIRI string) error {
	if attributedTo == nil || attributedTo.Len() != 1 || attributedTo.At(0) == nil {
		return errors.New("create note must have exactly one attributedTo")
	}
	value := attributedTo.At(0)
	var attributedToIRI string
	if value.IsIRI() && value.GetIRI() != nil {
		attributedToIRI = value.GetIRI().String()
	} else if object := value.GetType(); object != nil {
		var err error
		attributedToIRI, err = apmodels.GetIRIStringFromJSONLDIdProperty(object.GetJSONLDId())
		if err != nil {
			return errors.Wrap(err, "create note attributedTo is missing IRI")
		}
	}
	if attributedToIRI == "" || attributedToIRI != actorIRI {
		return errors.New("create note attributedTo does not match activity actor")
	}
	return nil
}

func firstInReplyToIRI(inReplyTo vocab.ActivityStreamsInReplyToProperty) string {
	if inReplyTo == nil || inReplyTo.Len() == 0 {
		return ""
	}
	first := inReplyTo.At(0)
	if first == nil {
		return ""
	}
	if first.IsIRI() && first.GetIRI() != nil {
		return first.GetIRI().String()
	}
	if object := first.GetType(); object != nil {
		iri, _ := apmodels.GetIRIStringFromJSONLDIdProperty(object.GetJSONLDId())
		return iri
	}
	return ""
}

func createOrNoteAddresses(activity vocab.ActivityStreamsCreate, note vocab.ActivityStreamsNote, localActorIRI string) bool {
	return toContainsIRI(activity.GetActivityStreamsTo(), localActorIRI) ||
		ccContainsIRI(activity.GetActivityStreamsCc(), localActorIRI) ||
		toContainsIRI(note.GetActivityStreamsTo(), localActorIRI) ||
		ccContainsIRI(note.GetActivityStreamsCc(), localActorIRI)
}

func toContainsIRI(property vocab.ActivityStreamsToProperty, target string) bool {
	if property == nil {
		return false
	}
	for iterator := property.Begin(); iterator != property.End(); iterator = iterator.Next() {
		if iterator.IsIRI() && iterator.GetIRI() != nil && iterator.GetIRI().String() == target {
			return true
		}
	}
	return false
}

func ccContainsIRI(property vocab.ActivityStreamsCcProperty, target string) bool {
	if property == nil {
		return false
	}
	for iterator := property.Begin(); iterator != property.End(); iterator = iterator.Next() {
		if iterator.IsIRI() && iterator.GetIRI() != nil && iterator.GetIRI().String() == target {
			return true
		}
	}
	return false
}

func noteContent(content vocab.ActivityStreamsContentProperty) (string, string) {
	if content == nil || content.Len() == 0 {
		return "", ""
	}
	first := content.At(0)
	if first == nil {
		return "", ""
	}
	if first.IsXMLSchemaString() {
		return first.GetXMLSchemaString(), ""
	}
	if first.IsRDFLangString() {
		values := first.GetRDFLangString()
		languages := make([]string, 0, len(values))
		for language := range values {
			languages = append(languages, language)
		}
		sort.Strings(languages)
		if len(languages) > 0 {
			return values[languages[0]], languages[0]
		}
	}
	return "", ""
}

type attachmentObject interface {
	GetActivityStreamsUrl() vocab.ActivityStreamsUrlProperty
	GetActivityStreamsMediaType() vocab.ActivityStreamsMediaTypeProperty
	GetActivityStreamsName() vocab.ActivityStreamsNameProperty
}

func noteAttachments(property vocab.ActivityStreamsAttachmentProperty) []events.FediverseAttachment {
	if property == nil {
		return nil
	}
	var attachments []events.FediverseAttachment
	for iterator := property.Begin(); iterator != property.End(); iterator = iterator.Next() {
		var object attachmentObject
		if iterator.IsActivityStreamsImage() {
			object = iterator.GetActivityStreamsImage()
		} else if iterator.IsActivityStreamsDocument() {
			object = iterator.GetActivityStreamsDocument()
		} else {
			continue
		}
		if attachment, ok := makeAttachment(object); ok {
			attachments = append(attachments, attachment)
		}
	}
	return attachments
}

func isSafeHTTPURL(iri *url.URL) bool {
	return iri != nil && (iri.Scheme == "http" || iri.Scheme == "https") && iri.Host != ""
}

func makeAttachment(object attachmentObject) (events.FediverseAttachment, bool) {
	urlProperty := object.GetActivityStreamsUrl()
	if urlProperty == nil || urlProperty.Len() == 0 || urlProperty.At(0) == nil || !urlProperty.At(0).IsIRI() {
		return events.FediverseAttachment{}, false
	}
	attachmentIRI := urlProperty.At(0).GetIRI()
	if !isSafeHTTPURL(attachmentIRI) {
		return events.FediverseAttachment{}, false
	}

	attachment := events.FediverseAttachment{URL: attachmentIRI.String()}
	if mediaType := object.GetActivityStreamsMediaType(); mediaType != nil && mediaType.IsRFCRfc2045() {
		attachment.MediaType = mediaType.Get()
	}
	if name := object.GetActivityStreamsName(); name != nil && name.Len() > 0 && name.At(0) != nil && name.At(0).IsXMLSchemaString() {
		attachment.Alt = utils.StripHTML(name.At(0).GetXMLSchemaString())
	}
	return attachment, true
}
