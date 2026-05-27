package pluginhost

import (
	"fmt"
	"sync"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/userrepository"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/utils"
)

// pluginChatbotKeyPrefix namespaces the datastore keys that map a plugin
// to its chatbot user's access token.
const pluginChatbotKeyPrefix = "plugins.chatbot."

// pluginChatbotProvisioner resolves the persistent chatbot user for a plugin,
// creating it on first use. The plugin->token mapping is stored in the
// datastore so the chatbot keeps the same identity (and chat history
// attribution) across restarts. The chatbot is a type="API" user, so it loads
// with IsBot=true and a DisplayName equal to the plugin name.
//
// Chatbot users are intentionally not removed when a plugin is disabled or
// uninstalled, so past messages keep their authorship.
type pluginChatbotProvisioner struct {
	users     userrepository.UserRepository
	datastore *datastore.Datastore

	mu    sync.Mutex
	cache map[string]*models.User
}

func newPluginChatbotProvisioner(users userrepository.UserRepository, ds *datastore.Datastore) *pluginChatbotProvisioner {
	return &pluginChatbotProvisioner{
		users:     users,
		datastore: ds,
		cache:     make(map[string]*models.User),
	}
}

// chatbotUser returns the plugin's chatbot user, creating and persisting one
// the first time it's needed.
func (p *pluginChatbotProvisioner) chatbotUser(pluginName string) (*models.User, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if user, ok := p.cache[pluginName]; ok {
		return user, nil
	}

	key := pluginChatbotKeyPrefix + pluginName
	if token, err := p.datastore.GetString(key); err == nil && token != "" {
		if user := p.users.GetUserByToken(token); user != nil {
			p.cache[pluginName] = user
			return user, nil
		}
	}

	token, err := utils.GenerateAccessToken()
	if err != nil {
		return nil, err
	}
	if err := p.users.InsertExternalAPIUser(token, pluginName, 0, nil); err != nil {
		return nil, err
	}
	if err := p.datastore.SetString(key, token); err != nil {
		return nil, err
	}

	user := p.users.GetUserByToken(token)
	if user == nil {
		return nil, fmt.Errorf("plugin chatbot user %q not found after creation", pluginName)
	}
	p.cache[pluginName] = user
	return user, nil
}
