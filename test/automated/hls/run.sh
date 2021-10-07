#!/bin/bash

set -e

function start_stream() {
  # Start streaming the test file over RTMP to
  # the local owncast instance.
  ffmpeg -hide_banner -loglevel panic -stream_loop -1 -re -i ../test.mp4 -vcodec libx264 -profile:v main -sc_threshold 0 -b:v 1300k -acodec copy -f flv rtmp://127.0.0.1/live/abc123 &
  STREAMING_CLIENT=$!
}

function update_storage_config() {
  echo "Configuring external storage to use ${S3_BUCKET}..."

  # Hard coded to admin:abc123 for auth
  curl 'http://localhost:8080/api/admin/config/s3' \
  -H 'Authorization: Basic YWRtaW46YWJjMTIz' \
  --data-raw "{\"value\":{\"accessKey\":\"${S3_ACCESS_KEY}\",\"acl\":\"\",\"bucket\":\"${S3_BUCKET}\",\"enabled\":true,\"endpoint\":\"${S3_ENDPOINT}\",\"region\":\"${S3_REGION}\",\"secret\":\"${S3_SECRET}\",\"servingEndpoint\":\"\"}}"
}

TEMP_DB=$(mktemp)

# Install the node test framework
npm install --silent >/dev/null

# Download a specific version of ffmpeg
if [ ! -d "ffmpeg" ]; then
  mkdir ffmpeg
  pushd ffmpeg >/dev/null
  curl -sL https://github.com/vot/ffbinaries-prebuilt/releases/download/v4.2.1/ffmpeg-4.2.1-linux-64.zip --output ffmpeg.zip >/dev/null
  unzip -o ffmpeg.zip >/dev/null
  PATH=$PATH:$(pwd)
  popd >/dev/null
fi

pushd ../../.. >/dev/null

# Build and run owncast from source
go build -o owncast main.go pkged.go
./owncast -database $TEMP_DB &
SERVER_PID=$!

function finish {
  echo "Cleaning up..."
  rm $TEMP_DB
  kill $SERVER_PID $STREAMING_CLIENT
}
trap finish EXIT

popd >/dev/null
sleep 5

# Start the stream.
start_stream

echo "Waiting..."
sleep 13

# Run tests against a fresh install with no settings.
npm test

# Determine if we should continue testing with S3 configuration.
if [[ -z "${S3_BUCKET}" ]]; then
  echo "No S3 configuration set"
  exit 0
fi

# Kill the stream.
kill $STREAMING_CLIENT
sleep 5

# Update the server config to use S3 for storage.
update_storage_config

# start the stream.
start_stream
echo "Waiting..."
sleep 13

# Re-run the HLS test against the external storage configuration.
npm test
