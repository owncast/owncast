# Plugin-based Viewer Authentication ("Auth Gate") — Design

Status: draft for review (no code yet)
Branch: `plugin-auth-gate` (worktree off `develop`)

## Goal

Let operators require viewers to authenticate before they can access an Owncast
instance, where the *authentication method* is supplied entirely by a plugin
(OAuth, Discord, x.com, magic links, SAML, anything reachable over HTTP). This
replaces the high-friction "put Vouch/oauth2-proxy in front of Owncast" pattern
with a first-class plugin capability.

The plugin is the **identity provider**; Owncast core is the **gatekeeper and
session authority**.

## Scope decision: gate the whole web server

When auth is on, the **entire web server is gated** — not just the HTML page.
The viewer page (`/`), the video (`/hls/*`), chat (`/ws`), and the API are all
behind the gate. Consequence, accepted: **the stream is only viewable through
Owncast's own web UI** — VLC/QuickTime/other native players can't carry the
session and therefore can't play a gated stream.

Gating only the HTML page was rejected: the page embeds nothing secret, and the
HLS URL would remain openly fetchable, so it would be privacy theater.

### External storage (S3/CDN) caveat — DOCUMENT, don't block

Under S3 storage, `rewritePlaylistLocations` rewrites playlists to **absolute
remote URLs**, so segments are fetched directly from S3/CDN and never touch the
Owncast server. The gate cannot see those requests.

- Gating `/hls/stream.m3u8` still prevents an unauthenticated visitor from
  *discovering* the current segment list.
- But under S3 the segment URLs themselves are world-readable, so a *leaked or
  shared* segment URL remains fetchable.

Therefore: **gate + local storage = airtight; gate + S3 = good friction, not
airtight.** This is a documentation item for operators and plugin authors, not a
hard block.

## Architecture

### The plugin is out of the per-request hot path

Once auth is on, every request hits the gate, including each HLS segment (a live
viewer pulls a new segment every ~2-4s × N viewers). Calling into the plugin's
embedded JS/Python engine per request would melt the server. So:

- **The plugin runs the login flow only** (infrequent: ~once per viewer session),
  over its existing `onHttpRequest` routes under `/plugins/<slug>/*`.
- **Core mints and checks a signed session cookie.** The per-request gate check
  is pure Go (verify signature + expiry) — no DB lookup, no engine call.

### Three tiers of checking

| Tier | When | Cost | What happens |
|------|------|------|--------------|
| 1 | every request (`/hls/*`, images, `/ws`, API) | Go-only: verify cookie signature + expiry | valid → pass; invalid/absent → 302 to login |
| 2 | `/` (index navigation) only | optional engine call: `onAuthCheck` | re-validate against provider; `ok`/`refresh`/`deny` |
| 3 | n/a | — | (no per-request denylist; revocation is Tier 2) |

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
| `/`, `/hls/*`, images | gate | verify signature + expiry → pass (no DB; this is why it's signed) |
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
// Core namespaces authId by plugin slug internally (github-auth:github:583231),
// so plugins cannot collide or spoof each other's users.
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

### Exemption set (bypass the gate)

Principle: **the gate covers only the otherwise-public surface; any route that
already enforces its own credential bypasses it.**

- the designated gate plugin's namespace `/plugins/<slug>/*` + its static assets
  (bootstrap)
- `/admin/*` — already `RequireAdminAuth`
- all external-API-token routes — already `RequireExternalAPIAccessToken`
  (a valid Bearer client carries no cookie; gating would 302 it into login)
- health/liveness endpoint (if present)

**Everything else is gated, including `/api/status` and embeds** — leaking
live-status / viewer count to anonymous visitors would undercut the privacy goal.
A future "advertise live status even when gated" option would be an explicit
operator opt-in, not a default bypass.

### Session credential details

- **Stateless signed cookie**, `HttpOnly`, `Secure`, `SameSite=Lax` (Lax, not
  Strict — the provider callback is a cross-site top-level redirect), `Path=/`.
- **Signing secret is core's responsibility**: auto-generated on first use,
  persisted in config, core-owned rotation (rotating = invalidate all sessions =
  panic button). Mechanically required: core verifies on the hot path without
  calling the plugin, so core must hold the secret. Plugin authors never touch it.
  (The OAuth *client* secret is a separate, plugin-config concern.)
- **TTL: operator-configurable, default 24h**, with **sliding refresh** on each
  `/` load. TTL is a real security knob because it's the revocation backstop.

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

## Open / deferred items

- **`onAuthCheck` precise input fields** (above is the proposed shape).
- **Admin UX**: callback-URL display + "configure before enabling" guidance +
  the "enabling gates the whole site" warning. Author-built page; core supplies
  `server.info()`.
- **Signing secret config key + rotation UI** in core admin.
- **Static maintenance page** for the fail-closed state.
- **Hard per-session revocation** (KV denylist) — explicitly deferred; revocation
  is page-load granularity via `onAuthCheck`.

## Implementation surface (where code lands)

**Owncast core** (`services/plugins` is the single source of truth; runtime
features go here):
- Gate middleware in the chi chain (`webserver/router`), with the exemption set.
- Signed session cookie mint/verify + persisted signing secret in config.
- Host functions `owncast.users.register`, `owncast.auth.grantSession`,
  `owncast.auth.endSession`; new plugin `AuthType`; per-plugin `authId` namespacing.
- `auth.gate` + `users.register` permission enforcement.
- `Manager.Enable` single-gate guard; fail-closed wiring + auto-disable → closed.
- `/ws` + chat REST cookie fallback (`HandleClientConnection`).
- `onAuthCheck` dispatch on the index path; fail-closed on error.
- Static "auth unavailable" page.

**Plugin SDK** (`owncast-plugin-sdk`, JS + Python in parallel):
- Type/handler for `onAuthCheck`; facades for `users.register`,
  `auth.grantSession`, `auth.endSession`.
- `auth.gate` / `users.register` permission docs; manifest validation.
- Wire-protocol doc updates; author guide section.
- Example plugin (`github-auth`) in `examples/js` **and** `examples/python`
  (kept in sync, README + INSTRUCTIONS).
```
