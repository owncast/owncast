#!/bin/bash

set -o errexit
set -o pipefail

# shellcheck disable=SC1091  # tools.sh is sourced at runtime; not available to the linter
source ../tools.sh

BUILD_ID=$((RANDOM % 7200 + 600))
CYPRESS_RECORD_KEY="e9c8b547-7a8f-452d-8c53-fd7531491e3b"
BROWSER="electron" # Default. Will try to use Google Chrome.

if hash google-chrome 2>/dev/null; then
	BROWSER="chrome"
	echo "Using Google Chrome as a browser."
else
	echo "Google Chrome not found. Using Electron."
fi

# Bundle the updated web code into the server codebase.
if [ -z "$SKIP_BUILD" ]; then
	echo "Bundling web code into server..."

	# Change to the root directory of the repository
	pushd "$(git rev-parse --show-toplevel)"

	./build/web/bundleWeb.sh >/dev/null

	popd
else
	echo "Skipping web build..."
fi

# Install the web test framework
if [ -z "$SKIP_BUILD" ]; then
	echo "Installing test dependencies..."
	npm install --quiet --no-progress

else
	echo "Skipping dependencies installation"
fi

set -o nounset

# Run one cypress group. Cloud recording (and the flags that require it) is
# CI-only: local runs stay fast and work offline.
#   run_cypress <group> <tags-env> <spec-glob> [extra cypress args...]
run_cypress() {
	local group=$1 tags=$2 spec=$3
	shift 3

	local args=(--browser "$BROWSER" --env "tags=$tags" --spec "$spec")
	if [ -n "${CI:-}" ]; then
		args+=(--record --key "$CYPRESS_RECORD_KEY" --parallel --ci-build-id "$BUILD_ID" --group "$group" --tag "${group/-/,}")
	fi

	npx cypress run "${args[@]}" "$@"
}

install_ffmpeg

start_owncast

run_cypress "desktop-offline" desktop "cypress/e2e/offline/*.cy.js"
run_cypress "mobile-offline" mobile "cypress/e2e/offline/*.cy.js" --config viewportWidth=375,viewportHeight=667

# Admin UI tests: desktop-only smoke coverage of the Ant Design-heavy admin
# interface. Real federation protocol tests live in test/automated/activitypub/.
run_cypress "desktop-admin" desktop "cypress/e2e/admin/*.cy.js"

start_stream

run_cypress "desktop-online" desktop "cypress/e2e/online/*.cy.js"
run_cypress "mobile-online" mobile "cypress/e2e/online/*.cy.js" --config viewportWidth=375,viewportHeight=667

