package outbox

import (
	"fmt"
	"net/url"

	"github.com/go-fed/activity/streams"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/persistence/federatedserversrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/utils"
	"github.com/owncast/owncast/services/activitypub/webfinger"
)

// SendFollowRequestToOwncastServerURL fetches the nodeinfo of the target
// Owncast server, validates it, and sends an ActivityPub Follow request.
// Used by the admin UI to follow another Owncast instance for the
// featured-streams mini-directory.
func (s *Service) SendFollowRequestToOwncastServerURL(targetServerURL string, isStreamConnected bool) error {
	if !s.configRepository.GetFederationEnabled() {
		return fmt.Errorf("federation is not enabled")
	}

	nodeinfo, err := utils.FetchNodeInfo(targetServerURL)
	if err != nil {
		return fmt.Errorf("failed to fetch nodeinfo from %s: %w", targetServerURL, err)
	}

	if err := utils.ValidateOwncastServer(nodeinfo); err != nil {
		return fmt.Errorf("server validation failed: %w", err)
	}

	if err := utils.ValidateFeaturedStreamsSupport(nodeinfo); err != nil {
		return fmt.Errorf("server validation failed: %w", err)
	}

	targetUsername, err := utils.ExtractFederationUsername(nodeinfo)
	if err != nil {
		return fmt.Errorf("failed to extract federation username: %w", err)
	}

	parsedURL, err := url.Parse(targetServerURL)
	if err != nil {
		return fmt.Errorf("failed to parse target URL: %w", err)
	}

	targetActorAccount := fmt.Sprintf("%s@%s", targetUsername, parsedURL.Host)

	return s.SendFollowToAccount(targetActorAccount, isStreamConnected)
}

// SendFollowToAccount sends a Follow activity to a fediverse account
// expressed in webfinger form (user@host).
func (s *Service) SendFollowToAccount(targetActorAccount string, isStreamConnected bool) error {
	links, err := webfinger.GetWebfingerLinks(targetActorAccount)
	if err != nil {
		return fmt.Errorf("failed to get webfinger links for %s: %w", targetActorAccount, err)
	}

	user := apmodels.MakeWebFingerRequestResponseFromData(links)

	actorIRI := user.Self
	if actorIRI == "" {
		return fmt.Errorf("no actor IRI found in webfinger response for %s", targetActorAccount)
	}

	actor, err := s.resolver.GetResolvedActorFromIRI(actorIRI)
	if err != nil {
		return fmt.Errorf("failed to resolve actor from IRI %s: %w", actorIRI, err)
	}

	return s.SendFollowToAccountURI(actor.ActorIri.String(), actor.Username, actor.ActorIri.Host, isStreamConnected)
}

// SendFollowToAccountURI sends a Follow activity to a fully-resolved
// actor IRI. The repo is updated with the resulting follow record so the
// later Accept can mark it accepted.
func (s *Service) SendFollowToAccountURI(targetActorID, targetUsername, targetServerURL string, isStreamConnected bool) error {
	localUsername := s.configRepository.GetFederationUsername()
	localActorIRI := s.builder.MakeLocalIRIForAccount(localUsername)

	followID := shortid.MustGenerate()
	followIRI := s.builder.MakeLocalIRIForResource(fmt.Sprintf("follow/%s", followID))

	followActivity := streams.NewActivityStreamsFollow()

	idProperty := streams.NewJSONLDIdProperty()
	idProperty.SetIRI(followIRI)
	followActivity.SetJSONLDId(idProperty)

	actorProperty := streams.NewActivityStreamsActorProperty()
	actorProperty.AppendIRI(localActorIRI)
	followActivity.SetActivityStreamsActor(actorProperty)

	objectProperty := streams.NewActivityStreamsObjectProperty()
	targetIRI, err := url.Parse(targetActorID)
	if err != nil {
		return fmt.Errorf("failed to parse target actor ID: %w", err)
	}
	objectProperty.AppendIRI(targetIRI)
	followActivity.SetActivityStreamsObject(objectProperty)

	toProperty := streams.NewActivityStreamsToProperty()
	toProperty.AppendIRI(targetIRI)
	followActivity.SetActivityStreamsTo(toProperty)

	unknownProps := followActivity.GetUnknownProperties()
	apmodels.SetBasicOwncastMetadata(unknownProps, s.configRepository, isStreamConnected)

	repo := federatedserversrepository.Get()
	if repo == nil {
		return fmt.Errorf("federated servers repository not initialised")
	}

	actorResponse, err := s.resolver.GetResolvedActorFromIRI(targetActorID)
	if err != nil {
		_ = repo.RemoveFederatedServerByIRI(targetServerURL)
		return fmt.Errorf("failed to resolve target actor: %w", err)
	}

	var inboxURL string
	if actorResponse.Inbox != nil {
		inboxURL = actorResponse.Inbox.String()
	} else {
		_ = repo.RemoveFederatedServerByIRI(targetServerURL)
		return fmt.Errorf("no inbox URL found for target actor")
	}

	jsonData, err := apmodels.Serialize(followActivity)
	if err != nil {
		log.Errorf("Failed to serialize follow activity: %v", err)
		return fmt.Errorf("failed to serialize follow activity: %w", err)
	}

	inboxURLParsed, err := url.Parse(inboxURL)
	if err != nil {
		log.Errorf("Failed to parse inbox URL %s: %v", inboxURL, err)
		return fmt.Errorf("failed to parse inbox URL: %w", err)
	}

	req, err := s.signer.CreateSignedRequest(jsonData, inboxURLParsed, localActorIRI)
	if err != nil {
		log.Errorf("Failed to create signed request: %v", err)
		return fmt.Errorf("failed to create signed request: %w", err)
	}

	s.workerpool.AddToOutboundQueue(req)

	log.Infof("Sent follow request to %s (actor: %s)", targetServerURL, targetActorID)

	return nil
}
