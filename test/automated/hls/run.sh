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
