# Owncast AI Agent Instructions

This is a repository consisting of a Go backend and a React frontend. It supports both an internal web application that is public facing, an internal admin web application, internal-only APIs for powering these web interfaces, and a set of APIs for third parties to take advantage of.

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
| Full web CI check       | `cd web && npm run check`           |
| Build Storybook         | `cd web && npm run build-storybook` |
| Generate API code       | `make api-generate`                 |
| Generate DB code        | `make sqlc`                         |
| Install dev tools       | `make install-tools`                |
| Install git hooks       | `make install-hooks`                |
| Start test stream       | `./test/ocTestStream.sh`            |

## Code Standards

- UI component standards can be found in https://docs.owncast.dev/develop-frontend-components

### Required Before Each Commit

- Ensure all code is formatted using `golangci-lint` for Go, `stylelint` for CSS, and `prettier` for JavaScript/React.
- Fix any formatting or linting errors.

### Development Flow

- Build: Ensure the code builds successfully using `go build` for Go and `npm run build` for React in the `web` directory.
- Tests: Write unit tests for Go code and Javascript logic. Use `go test` and `npm run test`.
- UI components: Components should be standalone and reusable and show up in Storybook to enable testing, prototyping, and iteration.
- API: All APIs should be documented using OpenAPI specifications. Use `build/gen-api.sh` to generate the stubs, and then fill them in with your actual API implementation.
- If you made new UI components or made changes ensure that Storybook still builds by running `npm run build-storybook`.

## Repository Structure

- `web/`: The web source code for the React frontend.
- Root of the repo: Go source code for the backend.
- `static/web/`: Generated static files for the web application. Do not edit or commit files in this directory directly; they are generated from the `web` directory. Ignore this.
- `test/automated/api/`: A series of automated integration tests for API endpoints written in Javascript.
- `test/automated/browser/`: A series of automated browser UI tests for actual real-world browser interaction.
- `build/`: Build and code generation scripts.

### Backend Structure

- `main.go`: Entry point. Initializes logging, database, config, core services, metrics, and the HTTP router.
- `webserver/router/`: Chi (v5) HTTP router with HTTP/2 support. Routes organized under `/api/`, `/api/admin/`, federation endpoints, `/hls/`, and `/ws`.
- `webserver/handlers/`: Request handlers for web, admin, static files, HLS streaming.
- `webserver/handlers/generated/`: Auto-generated API types and Chi server stubs from OpenAPI spec. Do not hand-edit.
- `webserver/router/middleware/`: Authentication (`RequireAdminAuth()` for admin endpoints via HTTP Basic Auth) and ActivityPub content negotiation middleware.
- `persistence/`: Repository pattern implementations (ConfigRepository, UserRepository, ChatMessageRepository, WebhookRepository, AuthRepository). All database access goes through these.
- `db/`: sqlc-generated type-safe database code. Schema in `db/schema.sql`, queries in `db/query.sql`. Run `make sqlc` after modifying.
- `core/`: Core streaming logic including transcoder, chat, webhooks, storage providers, playlist management.
- `core/storageproviders/`: Video storage backends (local filesystem and S3).
- `core/transcoder/`: FFmpeg-based video transcoding.
- `activitypub/`: Federation support (controllers, inbox/outbox, HTTP signatures, WebFinger, NodeInfo, worker pool).
- `models/`: Shared data models.
- `config/`: Configuration package with defaults and constants.
- `tools/`: Separate `go.mod` for development tool dependencies installed to `./bin/`.

### Frontend Structure

- `web/pages/`: Next.js 14 pages with static export.
- `web/components/`: React components organized by domain (`admin/`, `chat/`, `common/`, `layouts/`, `modals/`, `ui/`, `video/`, `action-buttons/`).
- `web/components/stores/`: Recoil state management stores.
- `web/interfaces/`: TypeScript interfaces.
- `web/i18n/`: Localization files (next-export-i18n).
- `web/style-definitions/`: Design token definitions.

### Key Technical Details

- **Go version**: 1.24+ (see `go.mod`)
- **Database**: SQLite (single file, no external DB server needed)
- **Frontend framework**: Next.js 14, React 18, TypeScript, static export
- **UI library**: Ant Design v4 with Less theming
- **State management**: Recoil
- **Styling**: SCSS with component-scoped files, Less for Ant Design
- **API spec**: OpenAPI 3.1 at `openapi.yaml`
- **Commit style**: Conventional Commits (`type(scope): description`)
- **Branches**: `develop` is the development branch; `master` is the production/release branch
- **CI**: GitHub Actions for Go tests, JS linting/testing, browser tests, Storybook builds, Chromatic visual regression, CodeQL security scanning
- **Go linting**: golangci-lint with cyclop (max complexity 15), dupl (threshold 200), gosec, gocritic, forbidigo (no `fmt.Print*`, `print`, `println`, or `panic`), and others. Full config in `.golangci.yml`.
- **JS linting**: ESLint with airbnb config, Prettier, Stylelint, knip for unused code detection
- **Makefile**: For common tasks like building, testing, linting, code generation, and installing dev tools/hooks. Run `make help` for a list of functions.

## Testing Web Frontend

Testing the web frontend requires running the frontend development server and the backend service together. The frontend development server must be started using `npm run dev` in the `web` directory, which will serve the web application at `http://localhost:3000`. The backend service can be started using `go run main.go` in the root of the repository.

Do not access the web application via `http://localhost:8080` or build the web project to copy files to the `static/web` directory for local development.

## Key Guidelines

1. All APIs are to be documented using OpenAPI specifications and code is to be generated using `build/gen-api.sh`. Additional details can be found at https://docs.owncast.dev/api-web-routing.
2. Write API tests for all new endpoints in the `test/automated/api` directory.
3. Use the `test/automated/browser` directory for browser-based tests for new functionality that simulate user interactions.
4. All user-facing frontend UI strings need to support localization. Use the `Translation` component to show most displayable strings. To create dynamic translated strings use the `next-export-i18n` library and the `t()` function. But use the `Translation` component unless there is a reason not to as it allows you to set default text. Read https://docs.owncast.dev/web-translations for more details. Test localization by adding "?lang=XX" with XX being a country code, such as "de" for German. Strings that have not yet been translated will not show as changed, but it's good to test anyway to make sure that previously translated strings have not been broken or regressed in any way.
5. A link to the PR's Storybook on Chromatic via the PR's Chromatic job should be included to help with review.
6. For API changes a before and after example of the API response should be added to the pull request to help with review.
7. For backend changes, a before and after example of logs to demonstrate the change should be added to the pull request to help with review.
8. Do not build the web project or copy the files to the `static/web` directory for local development.
9. When running the frontend development server, the local backend service must also be running. You can run the backend service using `go run main.go` in the root of the repository.
10. The credentials for the backend development backend are username: admin and password: abc123 and uses HTTP Basic Auth. This is used for the admin web application and the admin APIs.
11. The admin is found at `/admin`.
12. If a live stream video is needed to run, you can run `./test/ocTestStream.sh` to start an actual stream that will begin streaming from the local development server.
13. You should never commit the `static/web` directory to the repository. It is generated from the `web` directory and should be ignored in your commits.
14. Don't use emoji in code comments or commit messages. That's lame.
15. User interface components should use Ant Design v4 components, not Ant Design v5.
16. Database access should be through the repository patterns.
17. Any substantial new features or changes around Federation or ActivityPub should be documented in FEDERATION.md.

## External Documentation

- Project documentation: https://docs.owncast.dev
- Frontend component guide: https://docs.owncast.dev/develop-frontend-components
- API and routing: https://docs.owncast.dev/api-web-routing
- Localization: https://docs.owncast.dev/web-translations
- Contributing guide: https://docs.owncast.dev/contributor-guide
- Federation protocol: `FEDERATION.md`
