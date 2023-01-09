#!/bin/bash

source ../tools.sh

TEMP_DB=$(mktemp)

# Install the node test framework
npm install --quiet --no-progress

ffmpegInstall

pushd ../../.. >/dev/null || exit

# Build and run owncast from source
go build -o owncast main.go
./owncast -database "$TEMP_DB" &
SERVER_PID=$!

popd >/dev/null || exit
sleep 5

# Start streaming the test file over RTMP to
# the local owncast instance.
../../ocTestStream.sh &
FFMPEG_PID=$!

function finish {
	kill $SERVER_PID $FFMPEG_PID
	rm -fr "$TEMP_DB" "$FFMPEG_PATH"
}
trap finish EXIT

echo "Waiting..."
sleep 15

# Run the tests against the instance.
npm test
