This is a repository consisting of a Go backend and a React frontend. It supports both an internal web application that is public facing, an internal admin web application, internal-only APIs for powering these web interfaces, and a set of APIs for third parties to take advantage of.

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

- `web`: The web source code for the React frontend.
- Root of the repo: Go source code for the backend.
- `static/web`: Generated static files for the web application. Do not edit or commit files in this directory directly; they are generated from the `web` directory. Ignore this.
- `test/automated/api`: A series of automated integration tests for API endpoints written in Javascript.
- `test/automated/browser`: A series of automated browser UI tests for actual real-world browser interaction.
- `build/web`: Script to build and bundle the web application.

## Key Guidelines

1. All APIs are to be documented using OpenAPI specifications and code is to be generated using `build/gen-api.sh`. Additional details can be found at https://docs.owncast.dev/api-web-routing.
2. Write API tests for all new endpoints in the `test/automated/api` directory.
3. Use the `test/automated/browser` directory for browser-based tests for new functionality that simulate user interactions.
4. All user-facing frontend UI strings need to support localization. Use the `next-export-i18n` package for wrapping strings to enable this. Read https://docs.owncast.dev/web-translations for more details.
