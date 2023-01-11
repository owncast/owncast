#!/bin/bash

source ../tools.sh

# Install the node test framework
npm install --quiet --no-progress

install_ffmpeg

pushd ../../.. >/dev/null || exit

start_owncast

popd >/dev/null || exit

start_stream

# Run the tests against the instance.
npm test
