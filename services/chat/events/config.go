package events

import "github.com/owncast/owncast/persistence/configrepository"

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
