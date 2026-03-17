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

## Testing web frontend

Testing the web frontend requires running the frontend development server and the backend service together. The frontend development server must be started using `npm run dev` in the `web` directory, which will serve the web application at `http://localhost:3000`. The backend service can be started using `go run main.go` in the root of the repository.

Do not access the web application via http://localhost:8080 or build the web project to copy files to the `static/web` directory for local development.

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
