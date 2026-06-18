package outbox

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-fed/activity/streams"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/persistence/federatedserversrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/utils"
	"github.com/owncast/owncast/services/activitypub/webfinger"
)

const (
	maxFederatedServerNameLen = 200
	maxFederatedServerURLLen  = 2048
)

// clampFederatedMetadata bounds remote-supplied strings before persisting.
func clampFederatedMetadata(s string, max int) string {
	if len(s) > max {
		return s[:max]
	}
	return s
}

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

// SendUnfollowRequestToOwncastServerURL sends an Undo of our Follow to the
// given Owncast server so it stops treating us as a follower. Used when an
// admin unfeatures a server. Best-effort: the local record is removed
// regardless, so callers should treat a returned error as non-fatal.
func (s *Service) SendUnfollowRequestToOwncastServerURL(targetServerURL string) error {
	if !s.configRepository.GetFederationEnabled() {
		return nil
	}

	parsedURL, err := url.Parse(targetServerURL)
	if err != nil {
		return fmt.Errorf("failed to parse target URL: %w", err)
	}

	nodeinfo, err := utils.FetchNodeInfo(targetServerURL)
	if err != nil {
		return fmt.Errorf("failed to fetch nodeinfo from %s: %w", targetServerURL, err)
	}

	targetUsername, err := utils.ExtractFederationUsername(nodeinfo)
	if err != nil {
		return fmt.Errorf("failed to extract federation username: %w", err)
	}

	account := fmt.Sprintf("%s@%s", targetUsername, parsedURL.Host)
	links, err := webfinger.GetWebfingerLinks(account)
	if err != nil {
		return fmt.Errorf("failed to get webfinger links for %s: %w", account, err)
	}

	user := apmodels.MakeWebFingerRequestResponseFromData(links)
	if user.Self == "" {
		return fmt.Errorf("no actor IRI found in webfinger response for %s", account)
	}

	actor, err := s.resolver.GetResolvedActorFromIRI(user.Self)
	if err != nil {
		return fmt.Errorf("failed to resolve actor from IRI %s: %w", user.Self, err)
	}

	if !strings.EqualFold(actor.ActorIri.Host, parsedURL.Host) {
		return fmt.Errorf("resolved actor host %q does not match server host %q", actor.ActorIri.Host, parsedURL.Host)
	}
	if actor.Inbox == nil {
		return fmt.Errorf("no inbox found for actor %s", actor.ActorIri.String())
	}

	return s.sendUndoFollow(actor.ActorIri, actor.Inbox)
}

// sendUndoFollow builds and queues an Undo of a Follow whose actor is this
// server and whose object is the given remote actor.
func (s *Service) sendUndoFollow(targetActorIRI, inbox *url.URL) error {
	localUsername := s.configRepository.GetFederationUsername()
	localActorIRI := s.builder.MakeLocalIRIForAccount(localUsername)

	undoID := shortid.MustGenerate()
	undoIRI := s.builder.MakeLocalIRIForResource(fmt.Sprintf("undo/%s", undoID))

	// Reconstruct the Follow being undone (actor = us, object = remote actor).
	follow := streams.NewActivityStreamsFollow()
	followActor := streams.NewActivityStreamsActorProperty()
	followActor.AppendIRI(localActorIRI)
	follow.SetActivityStreamsActor(followActor)
	followObject := streams.NewActivityStreamsObjectProperty()
	followObject.AppendIRI(targetActorIRI)
	follow.SetActivityStreamsObject(followObject)

	undo := streams.NewActivityStreamsUndo()
	undoIDProperty := streams.NewJSONLDIdProperty()
	undoIDProperty.SetIRI(undoIRI)
	undo.SetJSONLDId(undoIDProperty)

	undoActor := streams.NewActivityStreamsActorProperty()
	undoActor.AppendIRI(localActorIRI)
	undo.SetActivityStreamsActor(undoActor)

	undoObject := streams.NewActivityStreamsObjectProperty()
	undoObject.AppendActivityStreamsFollow(follow)
	undo.SetActivityStreamsObject(undoObject)

	undoTo := streams.NewActivityStreamsToProperty()
	undoTo.AppendIRI(targetActorIRI)
	undo.SetActivityStreamsTo(undoTo)

	jsonData, err := apmodels.Serialize(undo)
	if err != nil {
		return fmt.Errorf("failed to serialize undo activity: %w", err)
	}

	req, err := s.signer.CreateSignedRequest(jsonData, inbox, localActorIRI)
	if err != nil {
		return fmt.Errorf("failed to create signed request: %w", err)
	}

	s.workerpool.AddToOutboundQueue(req)
	log.Infof("Sent unfollow (Undo) to %s", targetActorIRI.String())
	return nil
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

	// Pin the resolved actor to the host of the account we were asked to
	// follow. The featured-streams flow always targets a server URL the admin
	// entered, and an Owncast server's actor lives on that same host. If
	// webfinger points us at a different host, the directory entry would link
	// to one server while actually representing another, so refuse it. (This
	// path is only used by featured streams; general fediverse follows, which
	// may legitimately delegate the actor to another host, do not go through
	// here.)
	if at := strings.LastIndex(targetActorAccount, "@"); at != -1 {
		expectedHost := targetActorAccount[at+1:]
		if !strings.EqualFold(actor.ActorIri.Host, expectedHost) {
			return fmt.Errorf("resolved actor host %q does not match requested server host %q", actor.ActorIri.Host, expectedHost)
		}
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
	// The federated_servers record is keyed by the base server URL
	// (scheme://host) -- the same key the admin handler created it with and the
	// Accept handler later looks it up by. The targetServerURL argument is only
	// the host, so derive the full key here for repo lookups/updates.
	serverKey := fmt.Sprintf("%s://%s", targetIRI.Scheme, targetIRI.Host)
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
		_ = repo.RemoveFederatedServerByIRI(serverKey)
		return fmt.Errorf("failed to resolve target actor: %w", err)
	}

	var inboxURL string
	if actorResponse.Inbox != nil {
		inboxURL = actorResponse.Inbox.String()
	} else {
		_ = repo.RemoveFederatedServerByIRI(serverKey)
		return fmt.Errorf("no inbox URL found for target actor")
	}

	// Populate the pending record with the remote server's name and logo from
	// its resolved actor, so the directory shows the server immediately instead
	// of a blank row while the follow awaits acceptance. The actor is public, so
	// this works before the follow is accepted. The Accept handler refreshes
	// these values later.
	name := clampFederatedMetadata(actorResponse.Username, maxFederatedServerNameLen)
	displayName := clampFederatedMetadata(actorResponse.Name, maxFederatedServerNameLen)
	var logoURL string
	if actorResponse.Image != nil {
		logoURL = clampFederatedMetadata(actorResponse.Image.String(), maxFederatedServerURLLen)
	}
	if err := repo.UpdateServerMetadata(serverKey, name, displayName, displayName, logoURL); err != nil {
		log.Warnf("Failed to set initial metadata for featured server %s: %v", serverKey, err)
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
