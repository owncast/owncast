#!/bin/bash

# Install the node test framework
npm install --silent > /dev/null

# Download a specific version of ffmpeg
if [ ! -d "ffmpeg" ]; then
  mkdir ffmpeg
  pushd ffmpeg > /dev/null
  curl -sL https://github.com/vot/ffbinaries-prebuilt/releases/download/v4.2.1/ffmpeg-4.2.1-linux-64.zip --output ffmpeg.zip > /dev/null
  unzip -o ffmpeg.zip > /dev/null
  PATH=$PATH:$(pwd)
  popd > /dev/null
fi

pushd ../.. > /dev/null

# If there is no existing config file then copy the default.
if [ ! -f "config.yaml" ]; then
  cp config-example.yaml config.yaml
fi

# Build and run owncast from source
go build -o owncast main.go pkged.go
./owncast &
SERVER_PID=$!

popd > /dev/null
sleep 5

# Start streaming the test file over RTMP to
# the local owncast instance.
ffmpeg -hide_banner -loglevel panic -re -i test.mp4 -vcodec libx264 -profile:v main -sc_threshold 0 -b:v 1300k -acodec copy -f flv rtmp://127.0.0.1/live/abc123 &
FFMPEG_PID=$!

function finish {
  kill $SERVER_PID $FFMPEG_PID
}
trap finish EXIT

sleep 15

# Run the tests against the instance.
npm test