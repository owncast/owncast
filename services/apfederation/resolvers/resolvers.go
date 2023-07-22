package resolvers

import (
	"github.com/owncast/owncast/storage/configrepository"
)

type APResolvers struct {
	configRepository configrepository.ConfigRepository
}

func New() *APResolvers {
	return &APResolvers{
		configRepository: configrepository.Get(),
	}
}

var temporaryGlobalInstance *APResolvers

func Get() *APResolvers {
	if temporaryGlobalInstance == nil {
		temporaryGlobalInstance = New()
	}
	return temporaryGlobalInstance
}
