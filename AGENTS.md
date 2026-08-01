# Owncast AI Agent Instructions

This is the single source of truth for agent instructions in this repository. `CLAUDE.md` and `.github/copilot-instructions.md` point here.

Owncast is a self-hosted live streaming server: a Go backend plus a React/Next.js frontend. One frontend project builds two apps, the public viewer page and the admin at `/admin`, and the Go binary serves both. The backend exposes internal APIs for those apps plus a documented public API that third-party clients depend on.

## Quick Reference

| Task                    | Command                             |
| ----------------------- | ----------------------------------- |
| Build backend           | `go build -o owncast .`             |
| Run backend             | `go run main.go`                    |
| Run frontend dev server | `cd web && npm run dev`             |
| Run Go tests            | `go test ./...`                     |
| Run JS tests            | `cd web && npm test`                |
| Lint Go                 | `make lint`                         |
| Format Go               | `make fmt`                          |
| Lint JS/CSS             | `cd web && npm run lint`            |
| Format JS/CSS           | `cd web && npm run format`          |
| Typecheck frontend      | `cd web && npm run typecheck`       |
| Full web CI check       | `cd web && npm run check`           |
| Build Storybook         | `cd web && npm run build-storybook` |
| Generate API code       | `make api-generate`                 |
| Generate DB code        | `make sqlc`                         |
| Install dev tools       | `make install-tools`                |
| Install git hooks       | `make install-hooks`                |
| Start test stream       | `./test/ocTestStream.sh`            |

## Working Agreements

- Commit messages follow Conventional Commits (`type(scope): summary`). The `commit-msg` hook rejects anything else. Allowed types: `build`, `chore`, `ci`, `docs`, `feat`, `fix`, `perf`, `refactor`, `revert`, `style`, `test`.
- No emoji in code, comments, or commit messages.
- Run `make install-hooks` once and lefthook covers the pre-commit gate: gofumpt, golangci-lint, eslint, prettier, stylelint, knip, and the Go and JS tests. Without hooks, run `make fmt`, `make lint`, and `cd web && npm run lint && npm run format` by hand before committing.
- `develop` is the development branch and moves fast (automated bundle and translation commits land often); rebase before pushing. `master` is the release branch.

## Development Flow

- **Backend**: `go build -o owncast .`, `go test ./...`, `make lint`.
- **Frontend**: `cd web && npm run typecheck` for fast feedback. `npm run check` is the full CI-equivalent gate (lint, format, typecheck, knip, tests, build).
- **Never** copy frontend build output into `static/web/` and never commit that directory. It is generated during release builds; local development uses the dev server.
- **UI work**: components stay standalone and reusable with a story beside them, and `npm run build-storybook` must keep passing. Pattern guide: `web/components/_COMPONENT_HOW_TO.md`.
- **APIs**: edit `openapi.yaml` first, then run `make api-generate` (wraps `build/gen-api.sh`, which needs `npm install -g @redocly/cli`). It regenerates `webserver/handlers/generated/`; never hand-edit generated files, implement the stubs in `webserver/handlers/`.
- **Database**: schema changes are new numbered goose migrations in `persistence/migrations/`, and a migration that has shipped is never edited. Queries live in `db/query.sql`; run `make sqlc` after changing either, and don't hand-write SQL in Go. Adding a migration also means updating the pinned assertions in `persistence/migrations/migrations_test.go`. Walkthrough: `db/README.md`.
- **Tests**: Go tests beside the code, JS unit tests via `npm test`, API integration tests in `test/automated/api/`, browser flows in `test/automated/browser/`. Other suites live under `test/automated/` (`activitypub/`, `auth/`, `hls/`, `plugins/`, `upgrades/`).

## Testing the Web Frontend

Run the backend and the frontend dev server together: `go run main.go` in the repo root, `cd web && npm run dev`, then use `http://localhost:3000`.

Do not test against `http://localhost:8080`. That serves the web bundle embedded in the binary at compile time, not your working tree, so changes appear to have no effect.

`./test/ocTestStream.sh` starts a real stream against the local server when live-stream state is needed. Admin dev credentials are `admin` / `abc123` over HTTP Basic Auth; the admin UI is at `/admin`.

## Repository Structure

- Root of the repo: Go source for the backend.
- `web/`: the React frontend (viewer page, admin, embeds).
- `static/web/`: generated web bundle. Never edit or commit it.
- `test/automated/`: integration, browser, and protocol test suites.
- `build/`: build and code generation scripts.
- `docs/`: internal design and reference docs (`backend.md`, `plugins.md`, `product-definition.md`, `Release.md`).

### Backend Structure

- `main.go`: entry point and composition root. Initializes logging, database, config, core services, and metrics, then constructs every migrated service from `services/` and injects them into the HTTP router and other consumers. See "Service migration pattern" below.
- `webserver/router/`: Chi (v5) HTTP router with HTTP/2 support. Routes under `/api/`, `/api/admin/`, federation endpoints, `/hls/`, and `/ws`.
- `webserver/handlers/`: request handlers for web, admin, static files, HLS, moderation, auth.
- `webserver/handlers/generated/`: auto-generated API types and Chi server stubs from the OpenAPI spec. Do not hand-edit.
- `webserver/router/middleware/`: authentication (`RequireAdminAuth()` for admin endpoints via HTTP Basic Auth) and ActivityPub content negotiation.
- `persistence/`: repository implementations (config, user, chat message, webhook, auth, federated servers, notifications) plus goose migrations. All database access goes through these.
- `db/`: sqlc-generated type-safe database code. Queries in `db/query.sql`; the schema comes from the migrations, there is no `schema.sql`. Run `make sqlc` after changes.
- `services/`: constructor-injected services (stream, chat, transcoder, storage, rtmp, activitypub, plugins, webhooks, notifications, geoip, cache, datastore, dispatcher).
- `services/activitypub/`: federation (controllers, inbox/outbox, HTTP signatures, WebFinger, NodeInfo, delivery queue).
- `pluginhost/`: plugin host surface (auth gates, chat bots, plugin events, per-plugin SQLite storage in `sqlite.go`) on top of `services/plugins/`.
- `models/`: shared data models. `config/`: defaults and constants. `auth/`: IndieAuth and Fediverse chat auth. `yp/`: directory registration. `utils/`, `logging/`, `metrics/`, `notifications/`: supporting packages.
- `tools/`: separate `go.mod` for dev tool pins, installed to `./bin/`.

### Service migration pattern

Owncast is being moved off package-level singletons toward constructor-injected services. New code and migrated services follow these rules:

1. **A migrated service lives in `services/<domain>/`** and exposes `New(Deps) *Service`. Services with a lifecycle also expose `Start(ctx) error` / `Stop(ctx) error`. The package has no package-level state and no `Get()` (or equivalent) accessor.

2. **`main.go` is the composition root.** It's the only place that calls `services/*.New(...)`. Other code receives services via constructor injection. New deps appear in a `Deps` struct on each service or consumer.

3. **Existing `<pkg>.Get()` accessors (in `persistence/`) are deprecated** and exist only until their callers migrate. New code never calls `Get()`. PRs may reduce `Get()` callers; PRs must not add them.

When a service publishes events to multiple listeners, model the listeners as injected interfaces (`type StreamLifecycleListener interface { ... }`) and iterate them in the publisher. Don't introduce a generic event bus; explicit calls make stack traces, grep, and tests work the way the rest of the codebase expects.

`services/cache/` is the reference shape for a service, and the `*Handlers` struct in `webserver/handlers/handlers.go` is the reference pattern for consumers of injected services.

### Frontend Structure

- `web/pages/`: Next.js pages router with static export (viewer, `admin/`, `embed/`).
- `web/components/`: React components by domain (`admin/`, `chat/`, `common/`, `layouts/`, `modals/`, `ui/`, `video/`, `action-buttons/`, `theme/`).
- `web/components/stores/`: jotai state and the client config/status stores.
- `web/components/theme/`: Ant Design token theming (`AntdProvider.tsx`).
- `web/interfaces/`: TypeScript interfaces.
- `web/i18n/`: localization files (next-export-i18n).
- `web/style-definitions/`: design token definitions.

### Key Technical Details

- **Go version**: see `go.mod` (currently 1.26).
- **Database**: SQLite, single file, no external DB server.
- **Frontend**: Next.js 16 (pages router, static export), React 18, TypeScript.
- **UI library**: Ant Design v6. Tokens are themed in `web/components/theme/AntdProvider.tsx` and component CSS is pre-extracted to `web/styles/antd.css` via `npm run antd:extract`.
- **State management**: jotai.
- **Styling**: SCSS with component-scoped files.
- **API spec**: OpenAPI 3.1 in `openapi.yaml`.
- **Go linting**: golangci-lint with cyclop (max complexity 15), dupl (threshold 200), gosec, gocritic, forbidigo (no `fmt.Print*`, `print`, `println`, or `panic`), and others. Full config in `.golangci.yml`.
- **JS linting**: ESLint with airbnb config, Prettier, Stylelint, and knip for unused code detection.
- **CI**: GitHub Actions for Go tests, JS lint/test, browser tests, Storybook builds, Chromatic visual regression, CodeQL.
- **Makefile**: common tasks for building, testing, linting, code generation, and installing dev tools and hooks. `make help` lists them.

## Guidelines

1. All user-facing UI strings must be localizable. Prefer the `Translation` component because it carries default text; use `next-export-i18n`'s `t()` only for strings built at runtime. Test with `?lang=de` (untranslated strings won't visibly change, but previously translated ones must not regress). Reference: https://docs.owncast.dev/web-translations
2. `web/i18n/*/translation.json` files are tab-indented and outside prettier's config. Add keys with a tab-preserving edit; running prettier on them rewrites the entire file.
3. UI components use Ant Design v6. When you use an antd component that isn't already in the build, add it to the manifest in `web/build-scripts/extract-antd-styles.js` and run `npm run antd:extract`.
4. All database access goes through the `persistence/` repositories.
5. Substantial Federation or ActivityPub changes must be reflected in `FEDERATION.md`, which documents only the peer-visible protocol contract, not internal implementation.
6. Public API changes (config, status, HLS, directory, auth, payload shapes) affect downstream native clients (`owncast/owncasts-ios`, `owncast/roku`). Call that out in the PR.

## Pull Requests

- UI changes: include before and after screenshots, plus the Chromatic Storybook link from the PR's Chromatic job.
- API changes: include before and after example responses.
- Backend changes: include before and after log output demonstrating the change.
- Take screenshots from Storybook for standalone components (hide unrelated controls), or from the dev server at `http://localhost:3000` for anything that isn't in Storybook, which is most of the admin. Embed them inline in the PR, not as attachments.
- Keep screenshot scripts and images in `/tmp`. Never commit them.

## References

- Operator documentation: https://owncast.online/docs/
- Component library and design docs (Storybook): https://owncast.online/components
- Development guide: https://docs.owncast.dev/development
- Frontend component guide: `web/components/_COMPONENT_HOW_TO.md` and https://docs.owncast.dev/develop-frontend-components
- API and web routing: https://docs.owncast.dev/api-web-routing
- Contributor guide: https://docs.owncast.dev/contributor-guide
- Public API reference: https://owncast.online/docs/api/ (the spec of record is `openapi.yaml`)
- Federation protocol: `FEDERATION.md`
- Database workflow: `db/README.md`
