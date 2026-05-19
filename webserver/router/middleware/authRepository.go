package middleware

import "github.com/owncast/owncast/persistence/authrepository"

// authRepository is a TRANSITIONAL package-level handle installed by
// main.go via SetAuthRepository. Follow-up: replace with explicit
// parameter passing so the package has no globals (option B from the
// configrepository migration).
var authRepository authrepository.AuthRepository

// SetAuthRepository installs the package-level AuthRepository handle.
// TRANSITIONAL — see authRepository doc above.
func SetAuthRepository(repository authrepository.AuthRepository) {
	authRepository = repository
}
