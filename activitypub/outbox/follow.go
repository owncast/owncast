package outbox

import (
	"fmt"
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/activitypub/utils"
	"github.com/owncast/owncast/activitypub/webfinger"
	"github.com/owncast/owncast/activitypub/workerpool"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/federatedserversrepository"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

// SendFollowRequestToOwncastServerURL sends a follow request to another Owncast server
func SendFollowRequestToOwncastServerURL(targetServerURL string, isStreamConnected bool) error {
	configRepository := configrepository.Get()

	if !configRepository.GetFederationEnabled() {
		return fmt.Errorf("federation is not enabled")
	}

	// Fetch nodeinfo to get the federation username
	nodeinfo, err := utils.FetchNodeInfo(targetServerURL)
	if err != nil {
		return fmt.Errorf("failed to fetch nodeinfo from %s: %w", targetServerURL, err)
	}

	// Validate it's an Owncast server with ActivityPub enabled
	if err := utils.ValidateOwncastServer(nodeinfo); err != nil {
		return fmt.Errorf("server validation failed: %w", err)
	}

	// Extract the federation username
	targetUsername, err := utils.ExtractFederationUsername(nodeinfo)
	if err != nil {
		return fmt.Errorf("failed to extract federation username: %w", err)
	}

	// Parse the target URL to get the host
	parsedURL, err := url.Parse(targetServerURL)
	if err != nil {
		return fmt.Errorf("failed to parse target URL: %w", err)
	}

	// Construct the target actor account in webfinger format
	targetActorAccount := fmt.Sprintf("%s@%s", targetUsername, parsedURL.Host)

	return SendFollowToAccount(targetActorAccount, isStreamConnected)
}

func SendFollowToAccount(targetActorAccount string, isStreamConnected bool) error {
	// Use webfinger to get the actor links
	links, err := webfinger.GetWebfingerLinks(targetActorAccount)
	if err != nil {
		return fmt.Errorf("failed to get webfinger links for %s: %w", targetActorAccount, err)
	}

	// Create user from webfinger data
	user := apmodels.MakeWebFingerRequestResponseFromData(links)

	// Get the actor IRI from webfinger response
	actorIRI := user.Self
	if actorIRI == "" {
		return fmt.Errorf("no actor IRI found in webfinger response for %s", targetActorAccount)
	}

	// Resolve the actual actor to verify it exists
	actor, err := resolvers.GetResolvedActorFromIRI(actorIRI)
	if err != nil {
		return fmt.Errorf("failed to resolve actor from IRI %s: %w", actorIRI, err)
	}

	// Use the resolved actor's IRI to send the follow request
	return SendFollowToAccountURI(actor.ActorIri.String(), actor.Username, actor.ActorIri.Host, isStreamConnected)
}

func SendFollowToAccountURI(targetActorID, targetUsername, targetServerURL string, isStreamConnected bool) error {
	configRepository := configrepository.Get()

	// Get our local actor information
	localUsername := configRepository.GetFederationUsername()
	localActorIRI := apmodels.MakeLocalIRIForAccount(localUsername)

	// Create the Follow activity using go-fed methods
	followID := shortid.MustGenerate()
	followIRI := apmodels.MakeLocalIRIForResource(fmt.Sprintf("follow/%s", followID))

	// Create the Follow activity
	followActivity := streams.NewActivityStreamsFollow()

	// Set the activity ID
	idProperty := streams.NewJSONLDIdProperty()
	idProperty.SetIRI(followIRI)
	followActivity.SetJSONLDId(idProperty)

	// Set the actor (the local server)
	actorProperty := streams.NewActivityStreamsActorProperty()
	actorProperty.AppendIRI(localActorIRI)
	followActivity.SetActivityStreamsActor(actorProperty)

	// Set the object (the target server/actor)
	objectProperty := streams.NewActivityStreamsObjectProperty()
	targetIRI, err := url.Parse(targetActorID)
	if err != nil {
		return fmt.Errorf("failed to parse target actor ID: %w", err)
	}
	objectProperty.AppendIRI(targetIRI)
	followActivity.SetActivityStreamsObject(objectProperty)

	// Set the "to" field
	toProperty := streams.NewActivityStreamsToProperty()
	toProperty.AppendIRI(targetIRI)
	followActivity.SetActivityStreamsTo(toProperty)

	// Add Owncast metadata to the follow activity
	unknownProps := followActivity.GetUnknownProperties()
	apmodels.SetBasicOwncastMetadata(unknownProps, configRepository, isStreamConnected)

	// Store the follow request in the database with pending status
	repo := federatedserversrepository.Get()
	// err = repo.AddFederatedServer(
	// 	targetServerURL,
	// 	"",             // name will be populated when we receive the Accept
	// 	"",             // logo_url will be populated when we receive the Accept
	// 	time.Now(),     // followed_at
	// 	true,           // pending
	// 	targetUsername, // username
	// 	"pending",      // follow_status
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to store follow request: %w", err)
	// }

	// Resolve the target actor to get their inbox
	actorResponse, err := resolvers.GetResolvedActorFromIRI(targetActorID)
	if err != nil {
		// If we can't resolve the actor, remove the pending follow
		_ = repo.RemoveFederatedServerByIRI(targetServerURL)
		return fmt.Errorf("failed to resolve target actor: %w", err)
	}

	// Get the inbox URL
	var inboxURL string
	if actorResponse.Inbox != nil {
		inboxURL = actorResponse.Inbox.String()
	} else {
		// If we can't find an inbox, remove the pending follow
		_ = repo.RemoveFederatedServerByIRI(targetServerURL)
		return fmt.Errorf("no inbox URL found for target actor")
	}

	// Serialize the Follow activity
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

	req, err := crypto.CreateSignedRequest(jsonData, inboxURLParsed, localActorIRI)
	if err != nil {
		log.Errorf("Failed to create signed request: %v", err)
		return fmt.Errorf("failed to create signed request: %w", err)
	}

	workerpool.AddToOutboundQueue(req)

	log.Infof("Sent follow request to %s (actor: %s)", targetServerURL, targetActorID)

	return nil
}
