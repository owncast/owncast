#!/bin/bash

TEMP_DB=$(mktemp)

# Install the node test framework
npm install --quiet --no-progress

# Download a specific version of ffmpeg
if [ ! -d "ffmpeg" ]; then
	mkdir ffmpeg
	pushd ffmpeg >/dev/null || exit
	curl -sL https://github.com/vot/ffbinaries-prebuilt/releases/download/v4.2.1/ffmpeg-4.2.1-linux-64.zip --output ffmpeg.zip >/dev/null
	unzip -o ffmpeg.zip >/dev/null
	PATH=$(pwd):$PATH
	popd >/dev/null || exit
fi

pushd ../../.. >/dev/null || exit

# Build and run owncast from source
go build -o owncast main.go
./owncast -database "$TEMP_DB" -enableVerboseLogging &
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
