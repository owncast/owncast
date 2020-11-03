#!/usr/bin/env bash
# shellcheck disable=SC2059

set -o errexit
set -o nounset
set -o pipefail

INSTALL_TEMP_DIRECTORY="$(mktemp -d)"
PROJECT_SOURCE_DIR=$(pwd)
cd $INSTALL_TEMP_DIRECTORY

shutdown () {
  rm -rf "$PROJECT_SOURCE_DIR/admin"
  rm -rf "$INSTALL_TEMP_DIRECTORY"
}
trap shutdown INT TERM ABRT EXIT

echo "Cloning owncast admin into $INSTALL_TEMP_DIRECTORY..."
git clone --depth 1 https://github.com/owncast/owncast-admin 2> /dev/null
cd owncast-admin

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
cp -R ${ADMIN_BUILD_DIR}/out/* .

echo "Bundling admin into owncast codebase..."
~/go/bin/pkger

shutdown
echo "Done."
