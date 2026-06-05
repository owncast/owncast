// Package resolvers translates raw ActivityPub payloads and IRIs into
// our internal actor model: fetching, signing, resolving callbacks, and
// constructing follow/unfollow request representations. Construct a
// *Resolver in main.go with New(Deps) and pass it to consumers
// (inbox.*Service, jobs.*Service, outbox.*Service, persistence.*Service)
// via their own Deps structs.
package resolvers

import (
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/crypto"
)

// Resolver is the dependency-bearing handle for the ActivityPub
// resolver helpers. It reads federation defaults from
// configRepository, signs outbound requests with the *crypto.Signer,
// and uses the *apmodels.Builder for local IRI construction.
type Resolver struct {
	configRepository configrepository.ConfigRepository
	builder          *apmodels.Builder
	signer           *crypto.Signer
}

// Deps is the explicit dependency contract for Resolver.
type Deps struct {
	ConfigRepository configrepository.ConfigRepository
	Builder          *apmodels.Builder
	Signer           *crypto.Signer
}

// New constructs a Resolver with the provided dependencies.
func New(deps Deps) *Resolver {
	return &Resolver{
		configRepository: deps.ConfigRepository,
		builder:          deps.Builder,
		signer:           deps.Signer,
	}
}
