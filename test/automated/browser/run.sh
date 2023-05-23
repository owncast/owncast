#!/bin/bash

set -o errexit
set -o pipefail

source ../tools.sh

BUILD_ID=$((RANDOM % 7200 + 600))
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

install_ffmpeg

start_owncast

# Run cypress browser tests for desktop
npx cypress run --parallel --browser "$BROWSER" --group "desktop-offline" --env tags=desktop --ci-build-id $BUILD_ID --tag "desktop,offline" --record --key e9c8b547-7a8f-452d-8c53-fd7531491e3b --spec "cypress/e2e/offline/*.cy.js"
# Run cypress browser tests for mobile
npx cypress run --parallel --browser "$BROWSER" --group "mobile-offline" --env tags=mobile --ci-build-id $BUILD_ID --tag "mobile,offline" --record --key e9c8b547-7a8f-452d-8c53-fd7531491e3b --spec "cypress/e2e/offline/*.cy.js" --config viewportWidth=375,viewportHeight=667

start_stream

# Run cypress browser tests for desktop
npx cypress run --parallel --browser "$BROWSER" --group "desktop-online" --env tags=desktop --ci-build-id $BUILD_ID --tag "desktop,online" --record --key e9c8b547-7a8f-452d-8c53-fd7531491e3b --spec "cypress/e2e/online/*.cy.js"
# Run cypress browser tests for mobile
npx cypress run --parallel --browser "$BROWSER" --group "mobile-online" --env tags=mobile --ci-build-id $BUILD_ID --tag "mobile,online" --record --key e9c8b547-7a8f-452d-8c53-fd7531491e3b --spec "cypress/e2e/online/*.cy.js" --config viewportWidth=375,viewportHeight=667
