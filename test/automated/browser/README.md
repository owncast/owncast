# Browser tests

Cypress end-to-end tests for the Owncast web UI. They run against a real built
web bundle served by a real Owncast binary, with a real RTMP test stream for
the online specs and a real (local) remote ActivityPub server for the
federation specs. Nothing is mocked inside Owncast.

## Running

```bash
./run.sh
```

This builds the web project into the server, installs npm dependencies, builds
and starts Owncast on port 8080 with a throwaway database, and runs every test
group. A failed group does not stop the run: all groups execute so a single
run reports every failure.

Useful variations:

```bash
SKIP_BUILD=1 ./run.sh                 # reuse the existing web bundle and node_modules
./run.sh desktop-federation           # run only one group (see table below)
SKIP_BUILD=1 ./run.sh desktop-online  # fastest iteration on one group
npx cypress open --env tags=desktop   # interactive runner, against a running server
```

Cypress needs the `tags` env value (`desktop` or `mobile`) to decide which
tag-filtered tests run. Without it, tests wrapped in `filterTests()` are
skipped entirely.

Results are recorded to Cypress Cloud only when `CI` is set. Local runs upload
nothing and work offline.

## Reading failures

The suite is built so a failing run tells you exactly what broke:

- `cypress/results/failures.json` is written across all groups of a run and
  printed by `run.sh` at the end. Each entry has the spec, the full test
  title, the complete error text, and the failure screenshot paths.
- Any `console.error` logged by the application fails the current test, and
  the assertion message contains the logged error text itself (see
  `cypress/support/setup.js`).
- Failure screenshots land in `cypress/screenshots/`.
- Federation task failures (`cy.task('fediverse:*')`) time out with a message
  describing what the fake remote server received instead and what to check.

The intended loop for automated agents: make a change, `./run.sh`, read
`cypress/results/failures.json`, fix, re-run the single affected group with
`SKIP_BUILD=1 ./run.sh <group>` (after re-running `build/web/bundleWeb.sh`
from the repo root if web code changed).

## Test groups

`run.sh` runs six groups, in this order:

| Group              | Specs                     | Viewport | Server state |
|--------------------|---------------------------|----------|--------------|
| desktop-offline    | `cypress/e2e/offline/`    | default  | offline      |
| mobile-offline     | `cypress/e2e/offline/`    | 375x667  | offline      |
| desktop-admin      | `cypress/e2e/admin/`      | default  | offline      |
| desktop-online     | `cypress/e2e/online/`     | default  | live stream  |
| mobile-online      | `cypress/e2e/online/`     | 375x667  | live stream  |
| desktop-federation | `cypress/e2e/federation/` | default  | live stream  |

The offline and online directories run twice, once per viewport. Tests that
only make sense at one size are wrapped in `filterTests(['desktop'])` or
`filterTests(['mobile'])`. Between the admin and online groups, `run.sh`
starts `ocTestStream.sh` and polls `/api/status` until the server reports
`"online": true`.

## What is covered

- **Viewer page** (`offline/`, `online/`): header, offline banner, tags,
  video player presence, the followers tab, and the CSS identifier contract
  that theme authors rely on (`05-*_identifier_check`). If you rename one of
  those identifiers you must update the customization docs.
- **Embeds**: `/embed/video`, `/embed/chat/readwrite`, `/embed/chat/readonly`.
- **Chat**: visibility toggling, the mobile chat modal, name changes and
  message sends round-tripped over the real websocket (main page, mobile
  modal, and the readwrite embed).
- **Fediverse viewer UI** (`offline/07`): the offline banner account string,
  followers collection, and the follow modal end to end. The modal test
  submits `someone@remote.invalid` to `/api/remotefollow`, which fails
  WebFinger deterministically (`.invalid` never resolves, RFC 2606) and must
  surface the error alert. Known issue: at the mobile viewport the
  `#follow-button` fails Cypress's visibility check, so the modal tests are
  desktop-only. This may be a real mobile layout bug worth investigating.
- **Admin UI** (`admin/`): authenticated navigation smoke tests
  (`01-admin_smoke`) plus a render check of every admin page
  (`02-admin_pages`), each as its own named test. The admin is the most Ant
  Design-heavy surface in the project, so this group is the early-warning
  system for antd upgrades: a page that crashes or logs a React/antd error
  fails its own test, naming the page. `/admin/upgrade` is excluded because
  it fetches from the GitHub API.
- **Federation flows** (`federation/`): full end-to-end flows against the
  fake remote fediverse server (below). A signed inbound Follow must produce
  a federated Accept, a live engagement card in chat ("followed this
  stream"), the follower in the Followers tab, and a row in the admin
  followers table. Fediverse chat authentication runs the whole OTP journey:
  the account form, a real WebFinger lookup and signed DM delivery of the
  code, code entry, and the authenticated badge on a subsequent chat message.

Federation protocol behavior against real fediverse software (snac2) is
tested separately in `test/automated/activitypub/`. API behavior is tested in
`test/automated/api/`. Don't duplicate either here: this suite is for what a
browser renders and what a user clicks.

## The fake remote fediverse server

`cypress/support/remote-fediverse-server.js` runs inside the Cypress Node
process and is exposed to specs as `cy.task()`s:

- `fediverse:createActor` starts the server (idempotent) and returns a fresh
  remote actor `{ name, account, iri }` with its own RSA keypair and inbox.
- `fediverse:sendFollow` delivers a signed Follow to Owncast's inbox and
  waits for the Accept to arrive back.
- `fediverse:waitForOTPCode` waits for the one-time-code DM and extracts the
  code.

It serves real WebFinger, actor documents, and an inbox over self-signed
HTTPS on `https://127.0.0.1:9443`, and signs outbound activities with
draft-cavage HTTP signatures. Owncast treats it exactly like a remote
fediverse server. This requires `OWNCAST_ALLOW_INTERNAL_FEDERATION=true` and
`OWNCAST_INSECURE_SKIP_VERIFY=true` in Owncast's environment at process
start; `run.sh` exports both before starting the server.

Flow specs create a fresh actor per test attempt because Owncast
deduplicates repeat follows and pending OTP requests per account; a fresh
actor keeps Cypress retries meaningful.

## Conventions

- Specs set server config through `cy.setConfig(path, value)` and call other
  admin endpoints through `cy.adminRequest(method, url, body)`, both defined
  in `cypress/support/commands.js`. These chain into the Cypress command
  queue, so a following `cy.visit()` sees the new config. Don't use raw
  `fetch()` for setup, it races the test.
- `setup()` from `cypress/support/setup.js` fails any test whose page logged
  a `console.error`, and includes the logged text in the failure. Call it at
  the top of every spec.
- `testIsolation` is off (`cypress.config.js`): most specs visit a page once
  per `describe` and assert across multiple `it()` blocks. Keep new specs
  compatible with that pattern, and remember state carries between tests.
  Multi-step journeys (like the federation flows) should instead run as one
  `it()` so a retry replays the whole journey.
- Failed tests retry up to 3 times. Prefer an assertion that waits for the
  real effect (`cy.get('#user-menu').should('contain', ...)`) over
  `cy.wait(ms)`.
- Admin credentials are `admin` / `abc123`, passed to `cy.visit` as
  `{ auth: { username, password } }`.
- The `before()` hook in `cypress/support/e2e.js` sets the server URL to
  `https://testing.biz` for every spec file. Several assertions depend on
  that exact value, including the Owncast actor IRI used by the federation
  specs (`https://testing.biz/federation/user/streamer`).

## CI

`.github/workflows/browser-testing.yml` runs `./run.sh` on pushes and pull
requests that touch `web/`, Go code, or this directory, with 3 retry
attempts and a 20 minute timeout per attempt. In CI the groups are recorded
to Cypress Cloud under one build ID (`--record --parallel --group`), so all
groups appear as a single run in the dashboard.

## Versions

Cypress 15.x, pinned by `package-lock.json`. Video recording is disabled,
failure screenshots land in `cypress/screenshots/`. The Lighthouse
performance specs that used to live here were removed along with the
`cypress-audit` plugin when the suite moved off Cypress 10. If performance
budgets come back, run them as a separate lighthouse-ci job rather than
inside this suite.
