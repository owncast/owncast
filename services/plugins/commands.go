package plugins

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type commandCooldownKey struct {
	command string
	user    string
}

const (
	commandCooldownSweepIntervalMs int64 = 60_000
	commandModeratorScope                = "MODERATOR"
)

type commandTarget struct {
	plugin *Loaded
	event  CommandEvent
}

// DispatchCommands matches an accepted human chat message against every loaded
// plugin's command declarations and delivers chat.command directly to each
// match. Duplicate declarations across plugins intentionally all receive it.
func (d *Dispatcher) DispatchCommands(ctx context.Context, msg HostChatMessage) {
	targets := d.commandTargets(msg)

	var wg sync.WaitGroup
	for _, target := range targets {
		wg.Add(1)
		go func(target commandTarget) {
			defer wg.Done()
			deliverCommand(ctx, target.plugin, target.event)
		}(target)
	}
	wg.Wait()
}

func (d *Dispatcher) commandTargets(msg HostChatMessage) []commandTarget {
	var targets []commandTarget
	for _, p := range d.snapshot() {
		if p == nil || p.Manifest == nil {
			continue
		}
		for _, command := range p.Manifest.Commands {
			event, matched := matchCommand(msg, command)
			if !matched || !commandAllowed(msg, command) || !commandCooldownReady(p, msg, command) {
				continue
			}
			targets = append(targets, commandTarget{plugin: p, event: event})
		}
	}
	return targets
}

func matchCommand(msg HostChatMessage, command CommandInfo) (CommandEvent, bool) {
	prefix := command.Prefix
	if prefix == "" {
		prefix = DefaultCommandPrefix
	}
	if !strings.HasPrefix(msg.Body, prefix) {
		return CommandEvent{}, false
	}

	rest := strings.TrimSpace(strings.TrimPrefix(msg.Body, prefix))
	parts := strings.Fields(rest)
	if len(parts) == 0 {
		return CommandEvent{}, false
	}
	invokedAs := parts[0]
	if !commandNameMatches(invokedAs, command.Name, command.CaseSensitive) {
		matchedAlias := false
		for _, alias := range command.Aliases {
			if commandNameMatches(invokedAs, alias, command.CaseSensitive) {
				matchedAlias = true
				break
			}
		}
		if !matchedAlias {
			return CommandEvent{}, false
		}
	}

	return CommandEvent{
		Message:   msg,
		Command:   command.Name,
		InvokedAs: invokedAs,
		Args:      parts[1:],
		ArgString: strings.TrimSpace(rest[len(invokedAs):]),
	}, true
}

func commandNameMatches(called, declared string, caseSensitive bool) bool {
	if caseSensitive {
		return called == declared
	}
	return strings.EqualFold(called, declared)
}

func commandAllowed(msg HostChatMessage, command CommandInfo) bool {
	if !command.ModOnly || msg.User == nil {
		return !command.ModOnly
	}
	for _, scope := range msg.User.Scopes {
		if scope == commandModeratorScope {
			return true
		}
	}
	return false
}

func commandCooldownReady(p *Loaded, msg HostChatMessage, command CommandInfo) bool {
	if command.CooldownMs <= 0 {
		return true
	}
	now, err := time.Parse(time.RFC3339Nano, msg.Timestamp)
	if err != nil {
		return true
	}

	nowMs := now.UnixMilli()
	key := commandCooldownKey{command: command.Name, user: commandUserKey(msg)}
	p.commandCooldownMu.Lock()
	defer p.commandCooldownMu.Unlock()
	if p.commandCooldowns == nil {
		p.commandCooldowns = make(map[commandCooldownKey]int64)
	}
	if p.commandCooldownSweep == 0 || nowMs-p.commandCooldownSweep >= commandCooldownSweepIntervalMs {
		for existing, expiresAt := range p.commandCooldowns {
			if expiresAt <= nowMs {
				delete(p.commandCooldowns, existing)
			}
		}
		p.commandCooldownSweep = nowMs
	}
	if expiresAt, ok := p.commandCooldowns[key]; ok && nowMs < expiresAt {
		return false
	}
	expiresAt := int64(math.MaxInt64)
	if nowMs <= math.MaxInt64-command.CooldownMs {
		expiresAt = nowMs + command.CooldownMs
	}
	p.commandCooldowns[key] = expiresAt
	return true
}

func commandUserKey(msg HostChatMessage) string {
	if msg.User != nil && msg.User.ID != "" {
		return msg.User.ID
	}
	if msg.ClientID != 0 {
		return "client:" + strconv.FormatUint(msg.ClientID, 10)
	}
	return "anonymous"
}

func deliverCommand(ctx context.Context, p *Loaded, event CommandEvent) {
	envelope, err := json.Marshal(Envelope{EventType: EventChatCommand, Payload: event})
	if err != nil {
		fmt.Fprintf(os.Stderr, "plugin %s: marshal command event failed: %v\n", p.Manifest.Slug, err)
		return
	}
	if err := callOnEvent(ctx, p, envelope); err != nil &&
		!errors.Is(err, errPluginNotLoaded) && !errors.Is(err, errPluginNoSuchExport) {
		fmt.Fprintf(os.Stderr, "plugin %s: on_event(%s) failed: %v\n", p.Manifest.Slug, EventChatCommand, err)
	}
}
