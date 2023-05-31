#!/bin/bash

set -e
set -o errexit
set -o pipefail

finish() {
	# shellcheck disable=SC2317
	kill_with_kids "$BROWSERSTACK_PID"
	# shellcheck disable=SC2317
	kill_with_kids "$STREAM_PID"
}

rm -rf ./screenshots
mkdir -p ./screenshots

curl -o ./BrowserStackLocal-linux-x64.zip https://www.browserstack.com/browserstack-local/BrowserStackLocal-linux-x64.zip
unzip -o ./BrowserStackLocal-linux-x64.zip
./BrowserStackLocal --key "$BROWSERSTACK_KEY" &
BROWSERSTACK_PID=$!

trap finish EXIT TERM INT

npm install --silent >/dev/null
source ../tools.sh
install_ffmpeg
start_owncast

# Offline screenshots
FILE_SUFFIX="offline" node index.js

# Online screenshots
start_stream
sleep 20

FILE_SUFFIX="online" node index.js

SCREENSHOTS="$(pwd)/screenshots"
echo "$SCREENSHOTS"

# Change to the root directory of the repository
cd "$(git rev-parse --show-toplevel)"

cd web/.storybook/story-assets
rm -rf ./screenshots
mv "$SCREENSHOTS" .
