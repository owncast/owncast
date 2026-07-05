#!/bin/bash

set -e

# shellcheck disable=SC1091  # tools.sh is sourced at runtime; not available to the linter
source ../tools.sh

# Install the node test framework
npm install --quiet --no-progress

install_ffmpeg

start_owncast

start_stream

sleep 10

# Run tests against a fresh install with no settings.
npm test
