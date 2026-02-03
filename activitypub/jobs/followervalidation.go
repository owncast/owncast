package jobs

import (
	"time"

	"github.com/owncast/owncast/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	log "github.com/sirupsen/logrus"
)

const (
	// ValidationInterval is how often the validation job runs.
	ValidationInterval = 1 * time.Hour
	// FollowersPerRun is how many followers to validate per job run.
	FollowersPerRun = 5
	// FailureDurationThreshold is how long a follower must be unreachable before removal.
	FailureDurationThreshold = 7 * 24 * time.Hour // 7 days
	// DelayBetweenFollowers is the delay between validating individual followers.
	DelayBetweenFollowers = 2 * time.Second
)

// GetValidationInterval returns the configured validation interval, or the default.
func GetValidationInterval() time.Duration {
	if config.FollowerValidationInterval > 0 {
		return config.FollowerValidationInterval
	}
	return ValidationInterval
}

// StartFollowerValidationJob starts the background job that periodically validates followers.
func StartFollowerValidationJob() {
	interval := GetValidationInterval()
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			runFollowerValidation()
		}
	}()
	log.Debugf("Follower validation job scheduled with interval %v", interval)
}

func runFollowerValidation() {
	configRepo := configrepository.Get()
	if !configRepo.GetFederationEnabled() {
		return
	}

	followersRepo := followersrepository.Get()
	followers, err := followersRepo.GetFollowersToValidate(FollowersPerRun)
	if err != nil {
		log.Errorln("Failed to get followers for validation:", err)
		return
	}

	for _, follower := range followers {
		validateAndUpdateFollower(followersRepo, follower)
		time.Sleep(DelayBetweenFollowers)
	}
}

func validateAndUpdateFollower(repo followersrepository.FollowersRepository, follower models.Follower) {
	resolvedActor, err := resolvers.GetResolvedActorFromIRI(follower.ActorIRI)
	if err != nil {
		handleValidationFailure(repo, follower, err)
		return
	}

	// Success - clear failure timestamp and update data
	if err := repo.UpdateFollowerValidationSuccess(follower.ActorIRI); err != nil {
		log.Errorln("Failed to update validation success:", err)
	}

	// Update follower data
	if err := repo.Update(
		resolvedActor.ActorIriString(),
		resolvedActor.InboxString(),
		resolvedActor.SharedInboxString(),
		resolvedActor.Name,
		resolvedActor.FullUsername,
		resolvedActor.ImageString(),
	); err != nil {
		log.Errorln("Failed to update follower data:", err)
	}
}

func handleValidationFailure(repo followersrepository.FollowersRepository, follower models.Follower, resolveErr error) {
	log.Debugf("Follower validation failed for %s: %v", follower.ActorIRI, resolveErr)

	// Check removal eligibility BEFORE updating failure timestamp.
	// We use the existing FirstValidationFailureAt from the database to determine
	// if this follower has been failing long enough to be removed.
	// If FirstValidationFailureAt is not set, this is a new failure and we just record it.
	shouldRemove := false
	if follower.FirstValidationFailureAt.Valid {
		failureDuration := time.Since(follower.FirstValidationFailureAt.Time)
		shouldRemove = failureDuration >= FailureDurationThreshold
	}

	// Update failure timestamp (sets first_validation_failure_at if not already set)
	if err := repo.UpdateFollowerValidationFailure(follower.ActorIRI); err != nil {
		log.Errorln("Failed to update validation failure:", err)
		return
	}

	// Remove follower if they've exceeded the failure threshold
	if shouldRemove {
		failureDuration := time.Since(follower.FirstValidationFailureAt.Time)
		log.Infof("Removing follower %s after %v of consecutive failures",
			follower.ActorIRI, failureDuration.Round(time.Hour))
		if err := repo.RemoveByIRI(follower.ActorIRI); err != nil {
			log.Errorln("Failed to remove invalid follower:", err)
		}
	}
}
