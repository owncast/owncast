# Plugin host integration

How Owncast runs plugins and how the plugin system is wired into the server.
This is a contributor-facing overview of the **host side**; it is not a guide
to writing plugins.

## What plugins are

A plugin is a sandboxed WebAssembly module (authors write JavaScript/TypeScript
and compile it to wasm). Owncast loads plugins and runs them in-process via
[Extism](https://extism.org), which uses [Wazero](https://wazero.io) — a pure-Go
wasm runtime, so there is no CGo and no external process. A plugin can only do
what its manifest declares as **permissions**, and it can only reach the host
through a fixed set of **host functions** Owncast provides.

Plugins are optional infrastructure: if the runtime or a plugin fails to start,
Owncast logs and continues rather than aborting startup.

## Where the code lives

| Path | Responsibility |
|---|---|
| `services/plugins/` | The plugin **runtime** — discovery/enable lifecycle (`manager.go`), event + filter fan-out (`dispatcher.go`), the HTTP handler for `/plugins/<name>/*` (`server.go`), the host functions and their types/permissions (`hostfns.go`), the KV store interface (`kv/`), and the SSE hub (`sse.go`). |
| `pluginhost/` | The **Owncast integration layer** — builds the runtime's `HostEnv` from Owncast services, starts the manager, and exposes the HTTP handlers. This is the only package that knows about both the runtime and Owncast's internals. |

`services/plugins/` is a **vendored copy** of the upstream plugin SDK's
`host-runtime/plugin` package (`github.com/owncast/plugin-sdk`). See
[Keeping the runtime in sync](#keeping-the-runtime-in-sync-with-the-sdk) below —
this is the most important thing to understand before changing it.

## How a plugin call reaches Owncast

The runtime defines host functions in `services/plugins/hostfns.go`, but those
functions don't know anything about Owncast. Each one just calls a
function-pointer field on a `HostEnv` struct, e.g.:

```go
// services/plugins/hostfns.go — runtime, Owncast-agnostic
func hostVideoConfigRead(env *HostEnv) extism.HostFunction { /* calls env.VideoConfig() */ }

type HostEnv struct {
    VideoConfig func() VideoConfig // wired by the embedding host
    // ...
}
```

`pluginhost/pluginhost.go` is where Owncast **fills in** those pointers with
real services. `New(ctx, Deps)` constructs the `HostEnv` and `wirePluginHostEnv`
sets every field, grouped into helpers:

```go
// pluginhost/pluginhost.go — Owncast wiring
func wireVideoConfigHostFns(env *plugins.HostEnv, deps Deps) {
    cfg := deps.ConfigRepository
    env.VideoConfig = func() plugins.VideoConfig {
        variants := cfg.GetStreamOutputVariants()
        // map Owncast's models.StreamOutputVariant -> the plugin's StreamVariant
        // ...
    }
    env.WriteVideoConfig = func(pluginName string, u plugins.VideoConfigUpdate) error {
        // persist via cfg.SetStreamLatencyLevel / SetVideoCodec / SetStreamOutputVariants
    }
}
```

Two layers, two responsibilities:

- **`hostfns.go`** defines the plugin-facing API: the host-function names, the
  permissions that gate them, and the data **types** plugins receive
  (`VideoConfig`, `StreamBroadcaster`, …). This must match the SDK.
- **`pluginhost.go`** is the adapter: it reads Owncast's real state
  (`ConfigRepository`, `stream.Service`, chat, …) and **maps Owncast's internal
  `models.*` types into the plugin types**. This is Owncast-specific and is
  expected to differ from the SDK.

The plugin types are deliberately separate from `models.*`: they're a curated,
stable contract shared with plugins and the TypeScript SDK, so internal model
refactors don't silently change what plugins see. The mapping is the price of
that decoupling, and it lives at the boundary on purpose.

## What's wired

`wirePluginHostEnv` wires these subsystems (each behind its declared permission):

- **Chat send** (`chat.send`) — posts under the plugin's own provisioned bot
  identity (`newPluginChatbotProvisioner`), plus system/action messages and DMs.
- **Chat read / moderation** (`chat.history`, `chat.moderate`) — history,
  connected clients, delete message, kick client.
- **Server reads** (`server.read`) — stream status, server info, socials,
  federation, metadata tags, and read-only inbound-broadcast telemetry
  (`stream.Service.GetBroadcaster()`).
- **Video config** (`videoconfig.read` / `videoconfig.write`) — read and
  partially update output variants, codec, and latency via `ConfigRepository`.
  Writes **persist only**; they take effect on the next stream start (the admin
  UI already shows this), so the host deliberately does **not** restart a live
  transcoder.
- **Users** (`users.read` / `users.moderate`) — list/get users, set enabled,
  ban IP.
- **Notifications** (`notifications.send`) — Discord, browser push, fediverse.
- **Storage** — KV backed by Owncast's datastore (`newDatastoreKVStore`),
  `storage.upload` for browser-accessible plugin assets (written under
  `public/`), and `storage.fs` for a private, sandboxed filesystem
  (`data/plugin-data/<slug>/`) the plugin can read/write/list/delete within.
  The host rejects any `storage.fs` path that escapes the plugin's own
  directory, and the bytes are never served over HTTP.
- **HTTP** (`http.serve`, `http.sse`) — `req.authenticated` / `req.user` from
  Owncast's auth middleware, plus the SSE hub for server-pushed events.

## Events and filters

Owncast and plugins communicate through a shared dispatcher
(`services/dispatcher`), injected as `Deps.Events`:

- Owncast **publishes** events (chat messages, stream lifecycle, moderation, …)
  to that dispatcher. `pluginhost` registers a listener
  (`newPluginEventListener`) that forwards them to plugins' `on_event` handlers.
- For chat, `pluginhost` registers a **filter** (`newPluginChatFilter`) so
  plugins' `filterChatMessage` handlers run on inbound messages before they're
  broadcast (redact, drop, rewrite).
- Plugins can also **emit** custom events to each other; `env.Emit` is wired to
  the runtime's live dispatcher after the manager starts.

## Lifecycle, persistence, and HTTP

- **Discovery / enable:** the manager scans `<data>/plugins/` and tracks plugins
  as *discovered* vs *enabled*. Files are never auto-loaded; an admin enables
  them. The enabled set persists through `configEnabledStore` (Owncast's
  datastore, key `plugins.enabled`) rather than a JSON file.
- **HTTP:** `Host.Handler()` serves `/plugins/<name>/*` (mounted in the chi
  router). `Host.AdminHandler()` serves the admin API: `GET /api/admin/plugins`,
  `POST /api/admin/plugins/<name>/{enable,disable,reload}`, and the public
  `GET /api/plugins/actions`.

## Keeping the runtime in sync with the SDK

`services/plugins/` is a copy of the SDK's runtime. The **implementation** may
diverge for integration (e.g. Owncast adds a datastore-backed enabled-set store),
but the **plugin API** — the host-function names, permission identifiers, and
data-type shapes in `hostfns.go` — must stay identical to the SDK. If it
drifts, plugins built against the SDK can fail to load or receive malformed data
(this is exactly how the `videoconfig` functions were once missing here).

`services/plugins/contract_test.go` guards this. It re-derives the API surface
from this repo's `hostfns.go` and compares it to `services/plugins/plugin-contract.json`
(a snapshot copied verbatim from the SDK). If they differ, the test fails with
instructions. It checks **only** the API surface, not the `pluginhost.go`
wiring, so Owncast-specific adapter code is free to differ.

### Adding or changing a host function

Do it in the SDK first, then bring it here:

1. In the SDK (`github.com/owncast/plugin-sdk`), add the host function + type +
   permission to `host-runtime/plugin/hostfns.go`, expose it in the JS SDK, and
   regenerate the snapshot: `UPDATE_CONTRACT=1 go test ./plugin/ -run TestPluginContract`.
2. Copy the changed runtime files **and** `plugin-contract.json` into
   `services/plugins/` here.
3. Wire the new `HostEnv` field to real Owncast data in `pluginhost/pluginhost.go`.
4. `go test ./services/plugins/... ./pluginhost/...` — the contract test passes
   once the two `hostfns.go` copies agree.

For the SDK side of this system, see the SDK repo's `docs/ARCHITECTURE.md`.

## Tests

- Unit tests live in `services/plugins/` and `pluginhost/`.
- End-to-end coverage is in `test/automated/plugins/`: it builds the SDK's
  example plugins from source, installs them into a running Owncast, and
  exercises chat, filtering, and HTTP through them.
