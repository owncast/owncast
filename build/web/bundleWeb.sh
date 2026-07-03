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

if [ ! "$OFFLINE" ]; then
	echo "Installing npm modules for the owncast web..."
	# Keep stderr: discarding it turns real failures (wrong node version,
	# missing npm, registry errors) into a silent exit.
	npm --silent install
fi

echo "Building owncast web..."
rm -rf .next out
# No output filtering: piping through grep hid real build errors, and under
# pipefail a successful build with no matching lines would fail the script.
node_modules/.bin/next build

# Guard against a build that "succeeds" without producing the static export.
# Known failure mode: unsupported node versions (e.g. v23) make next build
# exit 0 after seconds having compiled nothing. Check before touching the
# existing bundle so a bad build never deletes a working static/web.
if [ ! -d ./out ]; then
	echo "ERROR: next build produced no ./out directory despite exiting successfully." >&2
	echo "This usually means the node version in use ($(node -v)) is unsupported by this Next.js version." >&2
	exit 1
fi

echo "Copying web project to dist directory..."

# Remove the old one
rm -rf ../static/web

# Copy over the new one
mv ./out ../static/web

echo "Done."
