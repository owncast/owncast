# Browser tests

Cypress end-to-end tests for the Owncast web UI. They run against a real built
web bundle served by a real Owncast binary, with a real RTMP test stream for
the online specs. Nothing is mocked.

## Running

```bash
./run.sh
```

This builds the web project into the server, installs npm dependencies, builds
and starts Owncast on port 8080 with a throwaway database, and runs every test
group. A full local run takes about 2.5 minutes after the web build.

Useful variations:

```bash
SKIP_BUILD=1 ./run.sh                # reuse the existing web bundle and node_modules
npx cypress run --browser chrome \
  --env tags=desktop \
  --spec "cypress/e2e/offline/*.cy.js"   # one group, against an already-running server
npx cypress open --env tags=desktop      # interactive runner
```

Cypress needs the `tags` env value (`desktop` or `mobile`) to decide which
tag-filtered tests run. Without it, tests wrapped in `filterTests()` are
skipped entirely.

Results are recorded to Cypress Cloud only when `CI` is set. Local runs upload
nothing and work offline.

## Test groups

`run.sh` runs five groups, in this order:

| Group           | Specs                  | Viewport  | Server state |
|-----------------|------------------------|-----------|--------------|
| desktop-offline | `cypress/e2e/offline/` | default   | offline      |
| mobile-offline  | `cypress/e2e/offline/` | 375x667   | offline      |
| desktop-admin   | `cypress/e2e/admin/`   | default   | offline      |
| desktop-online  | `cypress/e2e/online/`  | default   | live stream  |
| mobile-online   | `cypress/e2e/online/`  | 375x667   | live stream  |

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
- **Chat**: visibility toggling, the mobile chat modal, name changes round-
  tripped over the websocket.
- **Fediverse viewer UI** (`offline/07`): the offline banner account string,
  followers collection, and the follow modal end to end. The modal test
  submits `someone@remote.invalid` to `/api/remotefollow`, which fails
  WebFinger deterministically (`.invalid` never resolves, RFC 2606) and must
  surface the error alert. Known issue: at the mobile viewport the
  `#follow-button` fails Cypress's visibility check, so the modal tests are
  desktop-only. This may be a real mobile layout bug worth investigating.
- **Admin UI** (`admin/01-admin_smoke`): authenticated visit, sidebar
  navigation sections, federation-gated menu items, status indicator, the
  General configuration page and its tabs, and the followers table. The admin
  is the most Ant Design-heavy surface in the project, so this group is the
  early-warning system for antd upgrades.

Federation protocol behavior (signed activities, follows from real fediverse
servers, delivery) is tested separately in `test/automated/activitypub/`. API
behavior is tested in `test/automated/api/`. Don't duplicate either here:
this suite is for what a browser renders and what a user clicks.

## Conventions

- Specs set server config through `cy.setConfig(path, value)` and call other
  admin endpoints through `cy.adminRequest(method, url, body)`, both defined
  in `cypress/support/commands.js`. These chain into the Cypress command
  queue, so a following `cy.visit()` sees the new config. Don't use raw
  `fetch()` for setup, it races the test.
- `setup()` from `cypress/support/setup.js` spies on `console.error` and
  fails any test that logged one. Call it at the top of every spec.
- `testIsolation` is off (`cypress.config.js`): specs visit a page once per
  `describe` and assert across multiple `it()` blocks. Keep new specs
  compatible with that pattern, and remember state carries between tests.
- Failed tests retry up to 3 times. Prefer an assertion that waits for the
  real effect (`cy.get('#user-menu').should('contain', ...)`) over
  `cy.wait(ms)`.
- Admin credentials are `admin` / `abc123`, passed to `cy.visit` as
  `{ auth: { username, password } }`.
- The `before()` hook in `cypress/support/e2e.js` sets the server URL to
  `https://testing.biz` for every spec file. Several assertions depend on
  that exact value.

## CI

`.github/workflows/browser-testing.yml` runs `./run.sh` on pushes and pull
requests that touch `web/`, Go code, or this directory, with 3 retry
attempts and a 20 minute timeout per attempt. In CI the groups are recorded
to Cypress Cloud under one build ID (`--record --parallel --group`), so all
five groups appear as a single run in the dashboard.

## Versions

Cypress 15.x, pinned by `package-lock.json`. Video recording is disabled,
failure screenshots land in `cypress/screenshots/`. The Lighthouse
performance specs that used to live here were removed along with the
`cypress-audit` plugin when the suite moved off Cypress 10. If performance
budgets come back, run them as a separate lighthouse-ci job rather than
inside this suite.
