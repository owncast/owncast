// Package controllers hosts the HTTP request handlers for the
// ActivityPub federation surface: actor, inbox, outbox, followers,
// nodeinfo, webfinger and host-meta endpoints. Construct *Controllers
// in main.go (via activitypub.New) with the back-end services and
// register the methods on the router.
package controllers

import (
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/crypto"
	"github.com/owncast/owncast/services/activitypub/inbox"
	"github.com/owncast/owncast/services/activitypub/outbox"
	"github.com/owncast/owncast/services/activitypub/persistence"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
)

// Controllers is the dependency-bearing HTTP handler set for
// ActivityPub. Construct with New(Deps); register methods on the
// router. Pure helpers (createPageURL, writeResponse) and stateless
// handlers (ActorHandler, WebfingerHandler, HostMetaController) remain
// package-level free functions.
type Controllers struct {
	persistence      *persistence.Service
	outbox           *outbox.Service
	inbox            *inbox.Service
	followers        followersrepository.FollowersRepository
	configRepository configrepository.ConfigRepository
	builder          *apmodels.Builder
	signer           *crypto.Signer
}

// Deps lists every service the *Controllers consume.
type Deps struct {
	Persistence      *persistence.Service
	Outbox           *outbox.Service
	Inbox            *inbox.Service
	Followers        followersrepository.FollowersRepository
	ConfigRepository configrepository.ConfigRepository
	Builder          *apmodels.Builder
	Signer           *crypto.Signer
}

// New constructs the controllers set with explicit dependencies.
func New(deps Deps) *Controllers {
	return &Controllers{
		persistence:      deps.Persistence,
		outbox:           deps.Outbox,
		inbox:            deps.Inbox,
		followers:        deps.Followers,
		configRepository: deps.ConfigRepository,
		builder:          deps.Builder,
		signer:           deps.Signer,
	}
}
