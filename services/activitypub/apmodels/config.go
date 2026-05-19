package apmodels

import "github.com/owncast/owncast/persistence/configrepository"

// configRepository is a TRANSITIONAL package-level handle. The plan is to
// remove this and have callers pass the specific config values each
// helper needs as parameters (option B from the configrepository
// migration). Until that follow-up lands, main.go installs the repo
// here via SetConfigRepository at startup.
var configRepository configrepository.ConfigRepository

// SetConfigRepository installs the package-level ConfigRepository handle.
// TRANSITIONAL — see configRepository doc above.
func SetConfigRepository(repo configrepository.ConfigRepository) {
	configRepository = repo
}
