#!/bin/bash

TEMP_DB=$(mktemp)

# Install the node test framework
npm install --quiet --no-progress

# Download a specific version of ffmpeg
FFMPEG_PATH=$(mktemp -d)
pushd "$FFMPEG_PATH" >/dev/null || exit
curl -sL --fail https://github.com/ffbinaries/ffbinaries-prebuilt/releases/download/v4.4.1/ffmpeg-4.4.1-linux-64.zip --output ffmpeg.zip >/dev/null
unzip -o ffmpeg.zip >/dev/null
chmod +x ffmpeg
PATH=$FFMPEG_PATH:$PATH
popd >/dev/null || exit

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
	rm "$TEMP_DB"
	kill $SERVER_PID $FFMPEG_PID
}
trap finish EXIT

echo "Waiting..."
sleep 15

# Run the tests against the instance.
npm test
