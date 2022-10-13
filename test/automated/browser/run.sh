#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
set -e

TEMP_DB=$(mktemp)

# Change to the root directory of the repository
cd "$(git rev-parse --show-toplevel)"

# Bundle the updated web code into the server codebase.
echo "Bundling web code into server..."
./build/web/bundleWeb.sh >/dev/null

# Install the web test framework
echo "Installing test dependencies..."
pushd test/automated/browser
npm install --silent >/dev/null

popd

# Download a specific version of ffmpeg
if [ ! -d "ffmpeg" ]; then
	echo "Downloading ffmpeg..."
	mkdir -p /tmp/ffmpeg
	pushd /tmp/ffmpeg >/dev/null
	curl -sL https://github.com/vot/ffbinaries-prebuilt/releases/download/v4.2.1/ffmpeg-4.2.1-linux-64.zip --output ffmpeg.zip >/dev/null
	unzip -o ffmpeg.zip >/dev/null
	PATH=$PATH:$(pwd)
	popd >/dev/null
fi

# Build and run owncast from source
echo "Building owncast..."
go build -o owncast main.go
echo "Running owncast..."
./owncast -database $TEMP_DB &
SERVER_PID=$!

function finish {
	echo "Cleaning up..."
	rm $TEMP_DB
	kill $SERVER_PID $STREAMING_CLIENT
}
trap finish EXIT

pushd test/automated/browser

# Run cypress browser tests
npx cypress run --spec "cypress/e2e/offline/*.cy.js"

# Start streaming the test file over RTMP to
# the local owncast instance.
echo "Waiting for stream to start..."
ffmpeg -hide_banner -loglevel panic -stream_loop -1 -re -i ../test.mp4 -vcodec libx264 -profile:v main -sc_threshold 0 -b:v 1300k -acodec copy -f flv rtmp://127.0.0.1/live/abc123 &
STREAMING_CLIENT=$!

sleep 20

# npx cypress run --spec "cypress/e2e/*.cy.js"
npx cypress run --spec "cypress/e2e/online/*.cy.js"
