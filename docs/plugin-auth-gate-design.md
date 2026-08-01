# Plugin-based Viewer Authentication ("Auth Gate"): Design

Status: implemented. See `pluginhost/authgate.go` for the middleware,
`pluginhost/pluginhost.go` for the settings store and admin endpoint, and
`yp/` for directory behavior.

## Goal

Let operators require viewers to authenticate before they can access an Owncast
instance, where the *authentication method* is supplied entirely by a plugin
(OAuth, Discord, x.com, magic links, SAML, anything reachable over HTTP). This
replaces the high-friction "put Vouch/oauth2-proxy in front of Owncast" pattern
with a first-class plugin capability.

The plugin is the **identity provider**; Owncast core is the **gatekeeper and
session authority**.

## Scope decision: gate the web surface, with an operator-chosen edge

The gate is opt-out. Every request is challenged unless a named exemption
applies, so the viewer page (`/`), `/api/config`, chat, and the rest of the
public API are behind it.

The operator chooses one cumulative access mode. Native players such as VLC,
QuickTime, mobile apps, and restreamers cannot complete a browser login or carry
a session cookie. Any mode that gates `/hls/*` blocks those players.

The access mode is stored per gate plugin and selected on the plugin's
**Authentication** tab:

| Access mode | Effect |
|---|---|
| Website only (default) | The web UI requires sign-in. `/hls/*`, `/api/status`, and directory listing stay public. |
| Website, video players, and other resources | Also gates `/hls/*`. Native players without a browser session are blocked. `/api/status` and directory listing stay public. |
| Website, video players, and server status requests | Gates the web UI, `/hls/*`, and `/api/status`. Directory listing is disabled. |

The modes are cumulative. There is no status-only mode that hides
`/api/status` while leaving HLS public.

The default is deliberately the least disruptive. Enabling a gate on a running
instance protects the website without breaking players or monitoring that
already work. Operators who need the video itself private select one of the
stream-protection modes and accept losing native players.

### External storage (S3/CDN) caveat: document, do not block

Under S3 storage, `rewritePlaylistLocations` rewrites playlists to **absolute
remote URLs**, so segments are fetched directly from S3/CDN and never touch the
Owncast server. The gate cannot see those requests.

- Gating `/hls/stream.m3u8` still prevents an unauthenticated visitor from
  *discovering* the current segment list.
- But under S3 the segment URLs themselves are world-readable, so a *leaked or
  shared* segment URL remains fetchable.

Therefore: **stream protection + local storage = airtight. Stream protection +
S3 = good friction, not airtight.** This is a documentation item for operators
and plugin authors, not a hard block.

## Architecture

### The plugin is out of the per-request hot path

Once auth is on, every non-exempt request hits the gate. Under stream
protection that includes each HLS segment (a live viewer pulls a new segment
every ~2-4s × N viewers). Calling into the plugin's embedded JS/Python engine
per request would melt the server. So:

- **The plugin runs the login flow only** (infrequent: ~once per viewer session),
  over its existing `onHttpRequest` routes under `/plugins/<slug>/*`.
- **Core mints and checks a signed session cookie.** The per-request gate check
  is pure Go (verify signature + expiry), no DB lookup and no engine call.
- **The access policy is memoized.** Settings are read from the plugin's KV
  namespace once and cached in an atomic snapshot, so the hot path does not
  touch storage either. A read failure is not cached: it falls back to the
  defaults for that request and retries on the next one.

### Three tiers of checking

| Tier | When | Cost | What happens |
|------|------|------|--------------|
| 1 | every non-exempt request | Go-only: verify cookie signature + expiry | valid passes. Invalid or absent gets a 302 for GET/HEAD, 401 otherwise |
| 2 | `/` (index navigation) only, with a valid session | optional engine call: `onAuthCheck` | re-validate against provider, returning `ok`/`refresh`/`deny` |
| 3 | n/a | none | no per-request denylist, revocation is Tier 2 |

Which requests are exempt is decided by a single named list in
`pluginhost/authgate.go`. Adding a bypass means adding a named, individually
testable rule there, in the open:

| Exemption | Why | Conditional on |
|---|---|---|
| `admin` | `/admin/*` and `/api/admin/*` carry their own credential and have their own gate, so an operator can always disable a broken gate | never |
| `static-assets` | files that actually exist in the embedded web build, which the admin app needs to render. HTML entry points are excluded, so the viewer UI stays gated | never |
| `active-gate-plugin` | the login screen and OAuth callback must be reachable while signed out. Only the *active* gate's namespace, not every plugin | never |
| `external-api` | `/api/integrations/*` validates its own bearer token per route | never |
| `directory-api` | `/api/yp` is anonymous public metadata by design, and the YP handler itself refuses when listing is off | never |
| `third-party-player` | `/hls/*` for players that cannot log in | off when "protect the stream" is on |
| `stream-status` | `/api/status` for uptime monitors | off when "block stream status" is on |

A valid session is checked before any exemption, so a signed-in viewer is never
bounced off a path the policy closed to anonymous callers.

### The session cookie = a signed carrier for the existing access token

We do **not** invent a new identity primitive. Owncast already identifies users
by **access token** (`GetUserByToken`) and already links external identities to
users via `AddAuth(userID, authToken, authType)` / `GetUserByAuth` — this is how
IndieAuth and Fediverse auth work today. **A plugin auth gate is just a new
`AuthType`** feeding the same machinery.

The gate cookie is a **signed envelope carrying that user's existing access
token** plus a gate-session expiry:

```
cookie = sign({ accessToken, exp }, coreSecret)
```

Two readers, one cookie, both keyed on the access token you already understand:

| Path | Reader | Mechanism |
|------|--------|-----------|
| `/`, gated `/hls/*`, page assets | gate | verify signature + expiry → pass (no DB; this is why it's signed) |
| `/ws`, chat REST | chat | extract `accessToken` → `GetUserByToken(accessToken)` → existing chat identity path |

Why signed rather than a raw token in a cookie:
1. **Hot path** — the gate trusts `sig+exp` without a `GetUserByToken` per segment.
2. **Session expiry** — Owncast access tokens are long-lived; the envelope's
   `exp` gives the *gate session* its own sliding lifetime without changing the
   token model.

Backward compatible: the existing `localStorage` + `?accessToken=` path still
works on ungated servers; on a gated server, chat gains a "no query param? read
the access token from the gate cookie" fallback.

## Authoring surface

### New permissions

- **`auth.gate`** — be the gate plugin; use `grantSession` / `endSession`.
- **`users.register`** — create/link an *authenticated* user (separate so a
  non-gating chat-auth plugin, e.g. an IndieAuth-style "verified member" badge,
  can use identity without gating). Host fn: `owncast.users.register` (plural, to
  match existing `users.read` / `users.list`).

### New host functions

```ts
// Identity — find-or-create the user for an external identity, link + authenticate it.
// Pass the raw external id. Core stores your slug in the auth row's provider
// column and scopes lookups to (provider, authId), so plugins cannot collide
// or spoof each other's users. Do not prefix authId yourself.
owncast.users.register({ authId: string, displayName?: string, scopes?: string[] })
  : { userId: string }

// Session — issue the signed gate cookie carrying that user's access token.
// Only valid inside onHttpRequest (needs a live response to attach Set-Cookie).
owncast.auth.grantSession({ userId: string, ttl?: number }): void

// Self-logout — clear the current viewer's gate cookie on this response.
// Only valid inside onHttpRequest. Plugin still controls the redirect (and may
// bounce to the provider's logout for single-logout).
owncast.auth.endSession(): void
```

Core owns the cookie end to end: it reserves the cookie name, mints + signs, and
attaches `Set-Cookie` to the in-flight `onHttpRequest` response. The plugin never
sees or sets the signed token, so it can't forge or leak it.

### New optional hook

```ts
// Fires ONLY on '/' navigation (Tier 2). Optional. Core-driven (not a host fn).
// Lets the plugin re-validate against the provider and refresh/deny.
onAuthCheck(input: { user: { userId, displayName, scopes, authId } }):
  | { action: "ok" }                                      // pass, cookie unchanged
  | { action: "refresh", displayName?, scopes?, ttl? }    // still good; re-mint cookie
  | { action: "deny", reason?: string }                   // clear cookie, bounce to login

// Error / timeout in onAuthCheck → fail closed (treat as deny for that request).
```

Revocation model: a plugin that wants to kick a user (deleted upstream) returns
`deny` from `onAuthCheck` on that user's next `/` load. **Accepted limitation:** a
revoked viewer with an open tab keeps pulling segments until they reload `/` or
the cookie expires — so the cookie TTL is the hard backstop (see TTL below).
There is intentionally **no per-request denylist and no cross-user `revoke()`**.

## Lifecycle & enforcement

### Designation and arming

- A plugin declares `auth.gate`. **Declaring it does nothing on its own.**
- **Arming = enabling the plugin** via the existing enable/disable lifecycle
  (`services/plugins` Manager). There is **no separate "Require viewer auth"
  toggle** — that framing would falsely imply built-in Owncast auth. The plugin's
  admin page states plainly: *"While enabled, all viewers must authenticate
  through this plugin before accessing the site."*
- **`Manager.Enable` refuses to enable a second `auth.gate` plugin** while one is
  already enabled ("disable X first"). Centralized enablement means two can never
  be live at once — no load-order tie-break needed.
- Configure-before-live and fast-off-switch fall out for free: configure while
  discovered-but-disabled, enable to go live; disable to drop the gate instantly.

### Control loop (Model 1: plugin = web app)

```
GET /                              gate: no cookie → 302 /plugins/<slug>/?return_to=%2F
  (gate-plugin namespace is exempt)
GET /plugins/<slug>/               plugin renders login screen ("Sign in with GitHub")
  → 302 to provider (github.com/login/oauth/authorize?...&state=<rand>)
GET /plugins/<slug>/callback?code&state
  plugin: validate state (KV), exchange code (fetch), fetch user, enforce org
  → owncast.users.register({ authId, displayName })  → { userId }
  → owncast.auth.grantSession({ userId })            → core attaches Set-Cookie
  → return 302 → <sanitized return_to>
GET /                              gate: cookie valid → Tier 2 onAuthCheck → ok → render
GET /hls/3.ts                      gate: cookie valid → serve segment
WS  /ws                            chat: cookie → accessToken → GetUserByToken → "octocat"
```

- **Entry path by convention:** `/plugins/<slug>/`. No manifest field.
- **`return_to`** is appended by core and **sanitized to same-origin absolute
  paths** (reject `//host`, `https://…`) to avoid an open redirect.

### Fail-closed

The gate's posture is **decoupled from plugin runtime health**:

- armed + plugin healthy → normal flow
- armed + plugin unavailable (crashed, failed to load, errored, hit the
  auto-disable strike threshold) → **deny all viewer traffic**, serve a static
  core-owned "Authentication temporarily unavailable" page. **Never open.**
- **admin always bypasses** (existing Basic Auth) to fix config or disable.
- **already-valid sessions survive** an outage (cookie check needs no plugin).

The existing auto-disable-on-strikes must resolve to **closed + loud admin
alert**, never to a silently open site.

### Where the access policy lives

The access mode is host-owned, not plugin config. A gate plugin cannot widen or
narrow its own reach. It is keyed per plugin under the reserved KV key
`owncast.auth-gate-settings` in that plugin's namespace, so each gate keeps its
own policy. Switching gates does not carry one provider's choice over to
another.

The admin UI shows an **Authentication** tab for any installed plugin that
declares `auth.gate`, whether or not it is currently enabled, so an operator can
set the mode before arming the gate. It is backed by
`GET`/`POST /api/admin/plugins/<slug>/auth-settings`. The endpoint accepts
`website-only`, `website-and-stream`, or `website-stream-and-status` as the
`accessMode`. It returns 404 for plugins that do not declare the permission and
rejects unknown fields or unsupported modes.

### Directory interaction

`/api/yp` is never challenged by the gate. The Owncast Directory fetches it
anonymously: the registration key travels in the server-to-directory ping, not
in the directory's read of this endpoint, so there is no credential to check.
The payload is public metadata by design.

The `website-stream-and-status` mode does not gate `/api/yp`. It switches the
directory off. `Host.DirectoryAvailable()` reports the policy, and the YP
service consults it in two places: `GetYPResponse` answers 404, and the outbound
ping is skipped so the instance stops announcing itself to a directory that can
no longer read it. A listed stream cannot function without a readable status
endpoint, so status protection and directory availability move together.

### Session credential details

- **Stateless signed cookie**, `HttpOnly`, `Secure`, `SameSite=Lax` (Lax, not
  Strict — the provider callback is a cross-site top-level redirect), `Path=/`.
- **Signing secret is core's responsibility**: auto-generated on first use,
  persisted in config, core-owned rotation (rotating = invalidate all sessions =
  panic button). Mechanically required: core verifies on the hot path without
  calling the plugin, so core must hold the secret. Plugin authors never touch it.
  (The OAuth *client* secret is a separate, plugin-config concern.)
- **TTL is plugin-requested** via `grantSession({ ttl })`, defaulting to
  `DefaultSessionTTL` (24h) and capped at `MaxSessionTTL` (30 days). It is a
  real security knob because it is the revocation backstop: `onAuthCheck` only
  runs on a page load, so an idle revoked viewer lingers until the cookie
  expires.

## The chat / identity bridge

A gate login produces an authenticated chat identity automatically, because
`users.register` creates/links a real Owncast `User` (`Authenticated=true`,
display name seeded from the provider) and the cookie carries that user's access
token. Required core change: **`/ws` connect (and chat REST) gain a fallback —
no `?accessToken=` query param → read the access token from the gate cookie →
`GetUserByToken`.** No token is ever shuttled into the browser's localStorage.

## Worked example: "Sign in with GitHub" plugin

`manifest.json` permissions:
`["auth.gate", "users.register", "http.serve", "network.fetch", "storage.kv", "server.read"]`
with `network.allowedHosts: ["github.com", "api.github.com"]` and config fields
`clientId`, `clientSecret`, `allowedOrg`.

Handlers: `onHttpRequest` (routes `/`, `/callback`, `/logout`) + optional
`onAuthCheck` (re-verify org membership on each `/`).

- The plugin learns its own public base URL from `owncast.server.info()` (errors
  if the operator hasn't set the server URL) to build `redirect_uri` and to show
  the admin the exact callback URL to register at GitHub:
  `<serverURL>/plugins/github-auth/callback`.
- CSRF `state` stored in `storage.kv` with a short TTL, keyed to `return_to`.

## Still deferred

- **Hard per-session revocation** (KV denylist). Revocation stays at page-load
  granularity via `onAuthCheck`, with the cookie TTL as the backstop.
- **Signing secret rotation UI.** The secret is generated and persisted
  automatically, but rotating it (which invalidates every session) has no admin
  control yet.
- **Cross-gate policy migration.** Access settings are per plugin, so switching
  gate providers means setting the policy again on the new one. That is the
  intended behavior, not a gap, but an operator may be surprised by it.

## Where the code lives

**Owncast core:**
- `pluginhost/authgate.go`: middleware, exemption list, gate decision,
  `onAuthCheck` dispatch, fail-closed page.
- `pluginhost/pluginhost.go`: `AuthGateSettings`, the memoized settings store,
  `DirectoryAvailable`, and the `auth-settings` admin endpoint.
- `services/plugins/authsession.go`: cookie mint/verify, TTL clamping, the
  request-scoped grant sink.
- `services/plugins/manager.go`: `ActiveAuthGate`, `IsAuthGate`, and the
  single-gate guard in `Enable`.
- `yp/`: directory ping and `/api/yp`, both gated on `DirectoryAvailable`.
- `web/components/admin/plugins/AuthGateSettingsForm.tsx`: the Authentication
  tab.

**Plugin SDK** (`owncast-plugin-sdk`, JS and Python in parallel): the
`onAuthCheck` handler type, facades for `users.register`, `auth.grantSession`,
and `auth.endSession`, permission docs, and the `github-auth` / `basic-auth`
examples.

## Test coverage

| Area | Where |
|---|---|
| Cookie sign/verify, expiry, tampering | `services/plugins/authsession_test.go` |
| Single-gate guard, `ActiveAuthGate`, `IsAuthGate` | `services/plugins/manager_authgate_test.go` |
| Exemption list and gate decision, full policy matrix | `pluginhost/authgate_test.go` |
| Middleware: credential rejection, 302/401 split, identity propagation, fail-closed, revocation clearing the cookie | `pluginhost/authgate_middleware_test.go` |
| Settings endpoint, persistence, corrupt-value fallback | `pluginhost/authgate_settings_test.go` |
| Directory 404 and ping suppression | `yp/api_test.go` |
