package inbox

import (
	"github.com/owncast/owncast/services/apfederation/requests"
	"github.com/owncast/owncast/services/apfederation/resolvers"
	"github.com/owncast/owncast/services/chat"

	"github.com/owncast/owncast/storage/configrepository"
	"github.com/owncast/owncast/storage/federationrepository"
)

type APInbox struct {
	configRepository     configrepository.ConfigRepository
	federationRepository *federationrepository.FederationRepository
	resolvers            *resolvers.APResolvers
	requests             *requests.Requests
	chatService          *chat.Chat
}

func New() *APInbox {
	return &APInbox{
		configRepository:     configrepository.Get(),
		federationRepository: federationrepository.Get(),
		resolvers:            resolvers.Get(),
		requests:             requests.Get(),
		chatService:          chat.Get(),
	}
}

var temporaryGlobalInstance *APInbox

func Get() *APInbox {
	if temporaryGlobalInstance == nil {
		temporaryGlobalInstance = New()
	}
	return temporaryGlobalInstance
}
