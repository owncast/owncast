package outbox

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/activitypub/utils"
	"github.com/owncast/owncast/activitypub/workerpool"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/federatedserversrepository"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

// SendFollowRequest sends a follow request to another Owncast server
func SendFollowRequest(targetServerURL string) error {
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

	// Construct the target actor ID
	targetActorID := fmt.Sprintf("%s://%s/federation/user/%s", parsedURL.Scheme, parsedURL.Host, targetUsername)

	// Get our local actor information
	localUsername := configRepository.GetFederationUsername()
	localServerURL := configRepository.GetServerURL()
	localActorIRI := apmodels.MakeLocalIRIForAccount(localUsername)

	// Create the Follow activity
	followID, _ := shortid.Generate()
	followActivity := map[string]interface{}{
		"@context": "https://www.w3.org/ns/activitystreams",
		"id":       fmt.Sprintf("%s/federation/follow/%s", localServerURL, followID),
		"type":     "Follow",
		"actor":    localActorIRI.String(),
		"object":   targetActorID,
		"to":       []string{targetActorID},
	}

	// Store the follow request in the database with pending status
	repo := federatedserversrepository.Get()
	err = repo.AddFederatedServer(
		targetServerURL,
		"",             // name will be populated when we receive the Accept
		"",             // logo_url will be populated when we receive the Accept
		time.Now(),     // followed_at
		true,           // pending
		targetUsername, // username
		"pending",      // follow_status
	)
	if err != nil {
		return fmt.Errorf("failed to store follow request: %w", err)
	}

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

	// Send the Follow activity to the target's inbox
	jsonData, err := json.Marshal(followActivity)
	if err != nil {
		log.Errorf("Failed to marshal follow activity: %v", err)
		return fmt.Errorf("failed to marshal follow activity: %w", err)
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
