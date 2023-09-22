#!/usr/bin/env bash
# shellcheck disable=SC2059

set -o errexit
set -o nounset
set -o pipefail

OFFLINE=
while [[ $# -gt 0 ]]; do
  case $1 in
    --offline)
      OFFLINE=1
      ;;
  esac
  shift
done

# Change to the root directory of the repository
cd "$(git rev-parse --show-toplevel)"

cd web

if [ ! "$OFFLINE" ]
then
  echo "Installing npm modules for the owncast web..."
  npm --silent install 2>/dev/null
fi

echo "Building owncast web..."
rm -rf .next
(node_modules/.bin/next build && node_modules/.bin/next export) | grep info

echo "Copying web project to dist directory..."

# Remove the old one
rm -rf ../static/web

# Copy over the new one
mv ./out ../static/web

echo "Done."
