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

# Run the tests against the instance. Also write machine-readable results for
# CI artifact upload; jest keeps human-readable progress on stderr.
mkdir -p results
npm test -- --json --outputFile=results/jest-results.json
