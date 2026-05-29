package pluginhost

import (
	"testing"

	"github.com/owncast/owncast/persistence/userrepository"
)

func TestPluginChatbotProvisioner(t *testing.T) {
	ds := newTestDatastore(t)
	users := userrepository.New(ds)
	prov := newPluginChatbotProvisioner(users, ds)

	chatbot, err := prov.chatbotUser("uptime-bot", "Uptime Bot")
	if err != nil {
		t.Fatalf("chatbotUser: %v", err)
	}
	if chatbot == nil {
		t.Fatal("expected a chatbot user")
	}
	if chatbot.DisplayName != "Uptime Bot" {
		t.Errorf("display name = %q want %q", chatbot.DisplayName, "Uptime Bot")
	}
	if !chatbot.IsBot {
		t.Error("chatbot user should have IsBot=true")
	}

	// Second call returns the same identity (cache). Passing a different
	// display name on subsequent calls is a no-op: the user record was
	// created with the first display name.
	again, err := prov.chatbotUser("uptime-bot", "Uptime Bot")
	if err != nil {
		t.Fatalf("chatbotUser (cached): %v", err)
	}
	if again.ID != chatbot.ID {
		t.Errorf("cached lookup returned a different user: %q vs %q", again.ID, chatbot.ID)
	}

	// A different plugin slug gets a distinct identity.
	other, err := prov.chatbotUser("welcome-bot", "Welcome Bot")
	if err != nil {
		t.Fatalf("chatbotUser other: %v", err)
	}
	if other.ID == chatbot.ID {
		t.Error("different plugins must get different chatbot users")
	}

	// A fresh provisioner (cold cache) resolves the same persisted identity,
	// so the chatbot survives restarts.
	prov2 := newPluginChatbotProvisioner(users, ds)
	resolved, err := prov2.chatbotUser("uptime-bot", "Uptime Bot")
	if err != nil {
		t.Fatalf("chatbotUser after restart: %v", err)
	}
	if resolved.ID != chatbot.ID {
		t.Errorf("persisted chatbot identity not reused: %q vs %q", resolved.ID, chatbot.ID)
	}
}
