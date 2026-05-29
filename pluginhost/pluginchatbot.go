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
// to its chatbot user's access token. The key suffix is the plugin's
// slug so the same chatbot user (and message history) sticks with the
// plugin across reinstalls; the display-name half of the manifest can
// change without losing the bot's identity.
const pluginChatbotKeyPrefix = "plugins.chatbot."

// pluginChatbotProvisioner resolves the persistent chatbot user for a plugin,
// creating it on first use. The plugin->token mapping is stored in the
// datastore so the chatbot keeps the same identity (and chat history
// attribution) across restarts. The chatbot is a type="API" user, so it loads
// with IsBot=true and the DisplayName the manifest declared at creation
// time (Manifest.Bot.DisplayName, or Manifest.Name when Bot.DisplayName is
// empty: ChatDisplayName() picks).
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
// the first time it's needed. Two identity inputs: pluginSlug is the
// stable cache + datastore key (so the bot identity survives manifest
// edits to display name), and displayName is what chat viewers see for
// the bot. displayName is only used on first provisioning; once the
// user record exists, later calls return the existing record and
// don't rename it.
func (p *pluginChatbotProvisioner) chatbotUser(pluginSlug, displayName string) (*models.User, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if user, ok := p.cache[pluginSlug]; ok {
		return user, nil
	}

	key := pluginChatbotKeyPrefix + pluginSlug
	if token, err := p.datastore.GetString(key); err == nil && token != "" {
		if user := p.users.GetUserByToken(token); user != nil {
			p.cache[pluginSlug] = user
			return user, nil
		}
	}

	// New chatbot user: pick a sensible display name. Falling back to
	// the slug here is defensive: the host fn closure already resolves
	// ChatDisplayName(), so displayName should never be empty in
	// practice.
	chosenName := displayName
	if chosenName == "" {
		chosenName = pluginSlug
	}

	token, err := utils.GenerateAccessToken()
	if err != nil {
		return nil, err
	}
	if err := p.users.InsertExternalAPIUser(token, chosenName, 0, nil); err != nil {
		return nil, err
	}
	if err := p.datastore.SetString(key, token); err != nil {
		return nil, err
	}

	user := p.users.GetUserByToken(token)
	if user == nil {
		return nil, fmt.Errorf("plugin chatbot user %q not found after creation", pluginSlug)
	}
	p.cache[pluginSlug] = user
	return user, nil
}
