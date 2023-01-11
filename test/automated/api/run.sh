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
STREAM_PID=$!

echo "Waiting..."
sleep 15

# Run the tests against the instance.
npm test
