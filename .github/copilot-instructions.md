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
5. For UI component changes, a before and after screenshot of the component should be added to the pull request to help with review. Additionally a link to the PR's Storybook on Chromatic via the PR's Chromatic job should be included to help with review.
6. For API changes a before and after example of the API response should be added to the pull request to help with review.
7. For backend changes, a before and after example of logs to demonstrate the change should be added to the pull request to help with review.
8. For frontend changes run the frontend development server using `npm run dev` in the `web` directory to test changes locally and access it at `http://localhost:3000`. It is not required to build the web project or copy the files to the `static/web` directory for local development.
9. When running the frontend development server, the local backend service must also be running. You can run the backend service using `go run main.go` in the root of the repository.
10. When taking a screenshot of the web frontend or the admin web application, an instance of the Owncast backend service needs to be running locally by running `go run main.go` in the root of the repository as well.
11. The credentials for the backend development backend are username: admin and password: abc123 and uses HTTP Basic Auth. This is used for the admin web application and the admin APIs.
12. The admin is found at `/admin`.
13. If a live stream video is needed to run, you can run `./test/ocTestStream.sh` to start an actual stream that will begin streaming from the local development server.
