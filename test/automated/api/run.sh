#!/bin/bash

set -e

source ../tools.sh

# Install the node test framework
npm install --quiet --no-progress

install_ffmpeg

start_owncast

start_stream

sleep 10

# Run the tests against the instance.
npm test
