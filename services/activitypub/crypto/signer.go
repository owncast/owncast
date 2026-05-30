// Package crypto holds the ActivityPub signing/verification surface:
// keypair lookup, request signing, response signing, and signed-request
// construction. Construct a *Signer in main.go with New(Deps) and pass
// it to consumers (apmodels.Builder, apresolvers.Resolver, outbox,
// controllers, etc.) via their own Deps structs.
package crypto

import "github.com/owncast/owncast/persistence/configrepository"

// Signer is the dependency-bearing handle for the crypto helpers.
// It reads the local server's PEM-encoded keypair from the
// configRepository each time keys are needed. Construct with New(Deps).
type Signer struct {
	configRepository configrepository.ConfigRepository
}

// Deps is the explicit dependency contract for Signer.
type Deps struct {
	ConfigRepository configrepository.ConfigRepository
}

// New constructs a Signer with the provided dependencies.
func New(deps Deps) *Signer {
	return &Signer{
		configRepository: deps.ConfigRepository,
	}
}
