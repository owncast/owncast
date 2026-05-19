package crypto

import "github.com/owncast/owncast/persistence/configrepository"

// configRepository is a TRANSITIONAL package-level handle. Follow-up: have
// GetPublicKey/GetPrivateKey take the PEM string as a parameter so this
// package has no globals.
var configRepository configrepository.ConfigRepository

// SetConfigRepository installs the package-level ConfigRepository handle.
// TRANSITIONAL — see configRepository doc above.
func SetConfigRepository(repo configrepository.ConfigRepository) {
	configRepository = repo
}
