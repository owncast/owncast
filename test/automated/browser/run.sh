#!/bin/bash

set -o errexit
set -o pipefail

source ../tools.sh

TEMP_DB=$(mktemp)
BUILD_ID=$((RANDOM % 7200 + 600))
BROWSER="electron" # Default. Will try to use Google Chrome.

if hash google-chrome 2>/dev/null; then
	BROWSER="chrome"
	echo "Using Google Chrome as a browser."
else
	echo "Google Chrome not found. Using Electron."
fi

# Change to the root directory of the repository
pushd "$(git rev-parse --show-toplevel)"

# Bundle the updated web code into the server codebase.
if [ -z "$SKIP_BUILD" ]; then
	echo "Bundling web code into server..."
	./build/web/bundleWeb.sh >/dev/null
else
	echo "Skipping web build..."
fi

# Install the web test framework
if [ -z "$SKIP_BUILD" ]; then
	echo "Installing test dependencies..."
	pushd test/automated/browser
	npm install --quiet --no-progress

	popd
else
	echo "Skipping dependencies installation"
fi

set -o nounset

ffmpegInstall

# Build and run owncast from source
echo "Building owncast..."
go build -o owncast main.go

echo "Running owncast..."
./owncast -database "$TEMP_DB" &
SERVER_PID=$!

pushd test/automated/browser

# Run cypress browser tests for desktop
npx cypress run --browser "$BROWSER" --group "desktop-offline" --env tags=desktop --ci-build-id $BUILD_ID --tag "desktop,offline" --record --key e9c8b547-7a8f-452d-8c53-fd7531491e3b --spec "cypress/e2e/offline/*.cy.js"
# Run cypress browser tests for mobile
npx cypress run --browser "$BROWSER" --group "mobile-offline" --ci-build-id $BUILD_ID --tag "mobile,offline" --record --key e9c8b547-7a8f-452d-8c53-fd7531491e3b --spec "cypress/e2e/offline/*.cy.js" --config viewportWidth=375,viewportHeight=667

# Start streaming the test file over RTMP to
# the local owncast instance.
echo "Waiting for stream to start..."
../../ocTestStream.sh &
STREAMING_CLIENT=$!

function finish {
	echo "Cleaning up..."
	kill $SERVER_PID $STREAMING_CLIENT
	rm -fr "$TEMP_DB" "$FFMPEG_PATH"
}
trap finish EXIT SIGHUP SIGINT SIGTERM SIGQUIT SIGABRT SIGTERM

sleep 20

# Run cypress browser tests for desktop
npx cypress run --browser "$BROWSER" --group "desktop-online" --env tags=desktop --ci-build-id $BUILD_ID --tag "desktop,online" --record --key e9c8b547-7a8f-452d-8c53-fd7531491e3b --spec "cypress/e2e/online/*.cy.js"
# Run cypress browser tests for mobile
npx cypress run --browser "$BROWSER" --group "mobile-online" --ci-build-id $BUILD_ID --tag "mobile,online" --record --key e9c8b547-7a8f-452d-8c53-fd7531491e3b --spec "cypress/e2e/online/*.cy.js" --config viewportWidth=375,viewportHeight=667
