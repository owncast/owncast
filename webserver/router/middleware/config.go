package middleware

import (
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/userrepository"
)

// configRepository is a TRANSITIONAL package-level handle installed by
// main.go via SetConfigRepository. Follow-up: replace with explicit
// parameter passing so the package has no globals (option B from the
// configrepository migration).
var configRepository configrepository.ConfigRepository

// SetConfigRepository installs the package-level ConfigRepository handle.
// TRANSITIONAL — see configRepository doc above.
func SetConfigRepository(repo configrepository.ConfigRepository) {
	configRepository = repo
}

// userRepository is a TRANSITIONAL package-level handle installed by
// main.go via SetUserRepository. Follow-up: replace with explicit
// parameter passing so the package has no globals.
var userRepository userrepository.UserRepository

// SetUserRepository installs the package-level UserRepository handle.
// TRANSITIONAL — see userRepository doc above.
func SetUserRepository(repo userrepository.UserRepository) {
	userRepository = repo
}
