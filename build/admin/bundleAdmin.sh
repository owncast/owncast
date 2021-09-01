#!/usr/bin/env bash
# shellcheck disable=SC2059

set -o errexit
set -o nounset
set -o pipefail

INSTALL_TEMP_DIRECTORY="$(mktemp -d)"
PROJECT_SOURCE_DIR=$(pwd)
cd $INSTALL_TEMP_DIRECTORY

shutdown () {
  rm -rf "$INSTALL_TEMP_DIRECTORY"
}
trap shutdown INT TERM ABRT EXIT

echo "Cloning owncast admin into $INSTALL_TEMP_DIRECTORY..."
git clone https://github.com/owncast/owncast-admin 2> /dev/null
cd owncast-admin
git checkout gek/activity-pub-1

echo "Installing npm modules for the owncast admin..."
npm --silent install 2> /dev/null

echo "Building owncast admin..."
rm -rf .next
(node_modules/.bin/next build && node_modules/.bin/next export) | grep info

echo "Copying admin to project directory..."
ADMIN_BUILD_DIR=$(pwd)
cd $PROJECT_SOURCE_DIR
mkdir -p admin 2> /dev/null
cd admin

# Remove the old one
rm -rf $PROJECT_SOURCE_DIR/static/admin

# Copy over the new one
mv ${ADMIN_BUILD_DIR}/out $PROJECT_SOURCE_DIR/static/admin

shutdown
echo "Done."
