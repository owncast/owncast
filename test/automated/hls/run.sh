#!/bin/bash

set -e

source ../tools.sh

# Install the node test framework
npm install --silent >/dev/null

install_ffmpeg

pushd ../../.. >/dev/null

start_owncast

popd >/dev/null

start_stream

# Run tests against a fresh install with no settings.
npm test

# Determine if we should continue testing with S3 configuration.
if [[ -z "${S3_BUCKET}" ]]; then
  echo "No S3 configuration is set. Skipping S3 tests!"
  exit 0
fi

# Kill the stream.
kill_with_kids "$STREAM_PID"

# Update the server config to use S3 for storage.
update_storage_config

start_stream

# Re-run the HLS test against the external storage configuration.
npm test
