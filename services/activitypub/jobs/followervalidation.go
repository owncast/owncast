// Package jobs hosts the background validators and periodic tasks for
// the ActivityPub subsystem. Construct *Service with the followers
// repository and call Start to schedule the recurring follower
// validation pass.
package jobs

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	apresolvers "github.com/owncast/owncast/services/activitypub/resolvers"
)

const (
	// ValidationInterval is how often the validation job runs.
	ValidationInterval = 1 * time.Hour
	// FollowersPerRun is how many followers to validate per job run.
	FollowersPerRun = 5
	// FailureDurationThreshold is how long a follower must be unreachable before removal.
	FailureDurationThreshold = 7 * 24 * time.Hour
	// DelayBetweenFollowers is the delay between validating individual followers.
	DelayBetweenFollowers = 2 * time.Second
)

// GetValidationInterval returns the configured validation interval (from
// the supplied Config) when non-zero, otherwise the default
// ValidationInterval.
func GetValidationInterval(cfg *config.Config) time.Duration {
	if cfg != nil && cfg.FollowerValidationInterval > 0 {
		return cfg.FollowerValidationInterval
	}
	return ValidationInterval
}

// Service runs the recurring AP-side jobs.
type Service struct {
	followers        followersrepository.FollowersRepository
	configRepository configrepository.ConfigRepository
	resolver         *apresolvers.Resolver
	cfg              *config.Config
}

// Deps is the dependency contract for jobs.
type Deps struct {
	Followers        followersrepository.FollowersRepository
	ConfigRepository configrepository.ConfigRepository
	Resolver         *apresolvers.Resolver
	Config           *config.Config
}

// New constructs the jobs Service. Call Start to schedule the
// recurring tasks.
func New(deps Deps) *Service {
	return &Service{
		followers:        deps.Followers,
		configRepository: deps.ConfigRepository,
		resolver:         deps.Resolver,
		cfg:              deps.Config,
	}
}

// Start schedules the follower-validation tick.
func (s *Service) Start() {
	interval := GetValidationInterval(s.cfg)
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			s.runFollowerValidation()
		}
	}()
	log.Debugf("Follower validation job scheduled with interval %v", interval)
}

func (s *Service) runFollowerValidation() {
	if !s.configRepository.GetFederationEnabled() {
		return
	}

	followers, err := s.followers.GetFollowersToValidate(FollowersPerRun)
	if err != nil {
		log.Errorln("Failed to get followers for validation:", err)
		return
	}

	for _, follower := range followers {
		s.validateAndUpdateFollower(follower)
		time.Sleep(DelayBetweenFollowers)
	}
}

func (s *Service) validateAndUpdateFollower(follower models.Follower) {
	resolvedActor, err := s.resolver.GetResolvedActorFromIRI(follower.ActorIRI)
	if err != nil {
		s.handleValidationFailure(follower, err)
		return
	}

	// Success — clear failure timestamp and refresh actor data.
	if err := s.followers.UpdateFollowerValidationSuccess(follower.ActorIRI); err != nil {
		log.Errorln("Failed to update validation success:", err)
	}

	if err := s.followers.Update(
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

func (s *Service) handleValidationFailure(follower models.Follower, resolveErr error) {
	log.Debugf("Follower validation failed for %s: %v", follower.ActorIRI, resolveErr)

	// Check removal eligibility BEFORE updating the failure timestamp.
	// FirstValidationFailureAt being unset means this is a new failure
	// and we just record it.
	shouldRemove := false
	if follower.FirstValidationFailureAt.Valid {
		failureDuration := time.Since(follower.FirstValidationFailureAt.Time)
		shouldRemove = failureDuration >= FailureDurationThreshold
	}

	if err := s.followers.UpdateFollowerValidationFailure(follower.ActorIRI); err != nil {
		log.Errorln("Failed to update validation failure:", err)
		return
	}

	if shouldRemove {
		failureDuration := time.Since(follower.FirstValidationFailureAt.Time)
		log.Infof("Removing follower %s after %v of consecutive failures",
			follower.ActorIRI, failureDuration.Round(time.Hour))
		if err := s.followers.RemoveByIRI(follower.ActorIRI); err != nil {
			log.Errorln("Failed to remove invalid follower:", err)
		}
	}
}
