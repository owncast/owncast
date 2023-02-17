#!/bin/bash

set -e

source ../tools.sh

# Install the node test framework
npm install --silent >/dev/null

install_ffmpeg

start_owncast

start_stream

sleep 10

# Run tests against a fresh install with no settings.
npm test

# Kill the stream.
kill_with_kids "$STREAM_PID"

# Remove comments when #2615 is fixed.

# Determine if we should continue testing with S3 configuration.
# if [[ -z "${S3_BUCKET}" ]]; then
# 	echo "No S3 configuration is set. Skipping S3 tests!"
# 	exit 0
# fi

# # Update the server config to use S3 for storage.
# update_storage_config

# start_stream

# sleep 10

# # Re-run the HLS test against the external storage configuration.
# npm test
