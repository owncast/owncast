// Package apmodels owns the ActivityPub data-model builders and
// helpers: actor/service profiles, notes, activities, webfinger
// responses, IRI construction. The free-function nil-safety accessors
// in utils.go and the data-model types in actor.go, activity.go,
// hashtag.go, inboxRequest.go, message.go, webfinger.go remain
// package-level since they don't read configuration.
//
// The builders that DO read configuration (MakeServiceForAccount,
// MakeLocalIRIForAccount, MakeNote, MakeCreateActivity,
// GetCanonicalServerURL, GetLogoType, MakeActivityPublic,
// MakeUpdateActivity, MakeAddressingToFollowers, MakeWebfingerResponse,
// CreateCreateActivity, MakeLocalIRIForResource,
// MakeLocalIRIForStreamURL, MakeLocalIRIforLogo, MakeLocalURLForPath)
// are methods on *Builder. Construct in main.go with New(Deps); pass
// to consumers via their own Deps structs.
package apmodels

import (
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/crypto"
)

// Builder is the dependency-bearing handle for the ActivityPub
// model-construction helpers. It reads from configRepository for
// server URL, federation flags, logo, etc., and delegates to the
// crypto.Signer for actor public-key lookup.
type Builder struct {
	configRepository configrepository.ConfigRepository
	signer           *crypto.Signer
}

// Deps is the explicit dependency contract for Builder.
type Deps struct {
	ConfigRepository configrepository.ConfigRepository
	Signer           *crypto.Signer
}

// New constructs a Builder with the provided dependencies.
func New(deps Deps) *Builder {
	return &Builder{
		configRepository: deps.ConfigRepository,
		signer:           deps.Signer,
	}
}
