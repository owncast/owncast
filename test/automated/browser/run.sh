#!/bin/bash

set -e

TEMP_DB=$(mktemp)

# Change to the root directory of the repository
cd "$(git rev-parse --show-toplevel)"

# Bundle the updated web code into the server codebase.
./build/web/bundleWeb.sh

# Install the web test framework
pushd web
echo "Installing web dependencies..."
npm install --silent >/dev/null

cd "$(git rev-parse --show-toplevel)"

# Download a specific version of ffmpeg
if [ ! -d "ffmpeg" ]; then
	echo "Downloading ffmpeg..."
	mkdir ffmpeg
	pushd ffmpeg >/dev/null
	curl -sL https://github.com/vot/ffbinaries-prebuilt/releases/download/v4.2.1/ffmpeg-4.2.1-linux-64.zip --output ffmpeg.zip >/dev/null
	unzip -o ffmpeg.zip >/dev/null
	PATH=$PATH:$(pwd)
	popd >/dev/null
fi

# Build and run owncast from source
echo "Building owncast..."
go build -o owncast main.go
./owncast -database $TEMP_DB &
SERVER_PID=$!

function finish {
	echo "Cleaning up..."
	rm $TEMP_DB
	kill $SERVER_PID $STREAMING_CLIENT
}
trap finish EXIT

# Run cypress browser tests
cd web
npx cypress run
