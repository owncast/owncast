#!/bin/bash

set -e

source ../tools.sh

TEMP_DB=$(mktemp)

# Install the node test framework
npm install --silent >/dev/null

ffmpegInstall

pushd ../../.. >/dev/null

# Build and run owncast from source
go build -o owncast main.go
./owncast -database "$TEMP_DB" &
SERVER_PID=$!

popd >/dev/null
sleep 5

# Start the stream.
../../ocTestStream.sh &
STREAM_PID=$!

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
kill_with_kids "$STREAM_PID"
sleep 5

# Update the server config to use S3 for storage.
update_storage_config

# start the stream.
../../ocTestStream.sh &
STREAM_PID=$!

echo "Waiting..."
sleep 13

# Re-run the HLS test against the external storage configuration.
npm test
