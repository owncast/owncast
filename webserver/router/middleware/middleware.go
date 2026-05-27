package middleware

import (
	"github.com/owncast/owncast/persistence/authrepository"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/userrepository"
)

// Middleware bundles the dependencies needed by the HTTP middleware
// helpers in this package (admin basic-auth, external-API token auth,
// user-access-token auth, federation gating, etc.). Construct once in
// main() and inject into the components that wire HTTP routes.
type Middleware struct {
	configRepository configrepository.ConfigRepository
	authRepository   authrepository.AuthRepository
	userRepository   userrepository.UserRepository
}

// Deps lists the services *Middleware consumes.
type Deps struct {
	ConfigRepository configrepository.ConfigRepository
	AuthRepository   authrepository.AuthRepository
	UserRepository   userrepository.UserRepository
}

// New constructs a *Middleware.
func New(deps Deps) *Middleware {
	return &Middleware{
		configRepository: deps.ConfigRepository,
		authRepository:   deps.AuthRepository,
		userRepository:   deps.UserRepository,
	}
}
