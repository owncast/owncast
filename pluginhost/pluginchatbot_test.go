package pluginhost

import (
	"testing"

	"github.com/owncast/owncast/persistence/userrepository"
)

func TestPluginChatbotProvisioner(t *testing.T) {
	ds := newTestDatastore(t)
	users := userrepository.New(ds)
	prov := newPluginChatbotProvisioner(users, ds)

	chatbot, err := prov.chatbotUser("uptime-bot")
	if err != nil {
		t.Fatalf("chatbotUser: %v", err)
	}
	if chatbot == nil {
		t.Fatal("expected a chatbot user")
	}
	if chatbot.DisplayName != "uptime-bot" {
		t.Errorf("display name = %q want uptime-bot", chatbot.DisplayName)
	}
	if !chatbot.IsBot {
		t.Error("chatbot user should have IsBot=true")
	}

	// Second call returns the same identity (cache).
	again, err := prov.chatbotUser("uptime-bot")
	if err != nil {
		t.Fatalf("chatbotUser (cached): %v", err)
	}
	if again.ID != chatbot.ID {
		t.Errorf("cached lookup returned a different user: %q vs %q", again.ID, chatbot.ID)
	}

	// A different plugin gets a distinct identity.
	other, err := prov.chatbotUser("welcome-bot")
	if err != nil {
		t.Fatalf("chatbotUser other: %v", err)
	}
	if other.ID == chatbot.ID {
		t.Error("different plugins must get different chatbot users")
	}

	// A fresh provisioner (cold cache) resolves the same persisted identity,
	// so the chatbot survives restarts.
	prov2 := newPluginChatbotProvisioner(users, ds)
	resolved, err := prov2.chatbotUser("uptime-bot")
	if err != nil {
		t.Fatalf("chatbotUser after restart: %v", err)
	}
	if resolved.ID != chatbot.ID {
		t.Errorf("persisted chatbot identity not reused: %q vs %q", resolved.ID, chatbot.ID)
	}
}
