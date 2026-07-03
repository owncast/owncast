#!/bin/bash

# Browser test runner. Runs every Cypress group against a real Owncast
# instance, then prints a consolidated failure summary.
#
#   ./run.sh                    # run all groups
#   ./run.sh desktop-online     # run a single group (faster iteration)
#   SKIP_BUILD=1 ./run.sh       # reuse existing web bundle + node_modules
#
# A failed group does not stop the run: every group executes so one run
# reports every failure. The exit code is non-zero if any group failed, and
# cypress/results/failures.json holds machine-readable details per failure.

set -o errexit
set -o pipefail

# shellcheck disable=SC1091  # tools.sh is sourced at runtime; not available to the linter
source ../tools.sh

GROUP_FILTER="${1:-}"

BUILD_ID=$((RANDOM % 7200 + 600))
CYPRESS_RECORD_KEY="e9c8b547-7a8f-452d-8c53-fd7531491e3b"
BROWSER="electron" # Default. Will try to use Google Chrome.

if hash google-chrome 2>/dev/null; then
	BROWSER="chrome"
	echo "Using Google Chrome as a browser."
else
	echo "Google Chrome not found. Using Electron."
fi

# Bundle the updated web code into the server codebase. This is the slowest
# step of a run (npm install + a full Next.js production build), so say so
# up front and keep the noisy build output in a log instead of silently
# swallowing it: a silent multi-minute pause reads as a hang.
if [ -z "$SKIP_BUILD" ]; then
	WEB_BUILD_LOG="/tmp/owncast-bundle-web.log"
	echo "Bundling web code into server. This takes a few minutes with no output;"
	echo "progress is logged to ${WEB_BUILD_LOG}. Use SKIP_BUILD=1 to reuse the current bundle."

	# Change to the root directory of the repository
	pushd "$(git rev-parse --show-toplevel)" >/dev/null

	if ! ./build/web/bundleWeb.sh >"$WEB_BUILD_LOG" 2>&1; then
		echo "ERROR: web bundle build failed. Last lines of ${WEB_BUILD_LOG}:" >&2
		tail -30 "$WEB_BUILD_LOG" >&2
		exit 1
	fi

	popd >/dev/null
else
	echo "Skipping web build..."
fi

# Install the web test framework
if [ -z "$SKIP_BUILD" ]; then
	echo "Installing test dependencies..."
	npm install --quiet --no-progress

else
	echo "Skipping dependencies installation"
fi

set -o nounset

FAILED_GROUPS=()

# Run one cypress group. Cloud recording (and the flags that require it) is
# CI-only: local runs stay fast and work offline. A failing group is recorded
# and the run continues, so a single run surfaces every broken group.
#   run_cypress <group> <tags-env> <spec-glob> [extra cypress args...]
run_cypress() {
	local group=$1 tags=$2 spec=$3
	shift 3

	if [ -n "$GROUP_FILTER" ] && [ "$group" != "$GROUP_FILTER" ]; then
		return 0
	fi

	local args=(--browser "$BROWSER" --env "tags=$tags" --spec "$spec")
	if [ -n "${CI:-}" ]; then
		args+=(--record --key "$CYPRESS_RECORD_KEY" --parallel --ci-build-id "$BUILD_ID" --group "$group" --tag "${group/-/,}")
	fi

	if ! npx cypress run "${args[@]}" "$@"; then
		FAILED_GROUPS+=("$group")
	fi
}

# The federation specs drive real signed ActivityPub exchanges between a fake
# remote fediverse server (cypress/support/remote-fediverse-server.js) and
# this Owncast instance, both on localhost. These env vars relax Owncast's
# SSRF guard and TLS verification for that loopback federation. They are read
# once at first use (sync.Once), so they must be set before the server starts.
export OWNCAST_ALLOW_INTERNAL_FEDERATION=true
export OWNCAST_INSECURE_SKIP_VERIFY=true

install_ffmpeg

start_owncast

# Start each run with a clean failure ledger (appended to by after:spec in
# cypress.config.js across all groups).
rm -f cypress/results/failures.json

run_cypress "desktop-offline" desktop "cypress/e2e/offline/*.cy.js"
run_cypress "mobile-offline" mobile "cypress/e2e/offline/*.cy.js" --config viewportWidth=375,viewportHeight=667

# Admin UI tests: desktop-only smoke coverage of the Ant Design-heavy admin
# interface. Real federation protocol tests live in test/automated/activitypub/.
run_cypress "desktop-admin" desktop "cypress/e2e/admin/*.cy.js"

start_stream

run_cypress "desktop-online" desktop "cypress/e2e/online/*.cy.js"
run_cypress "mobile-online" mobile "cypress/e2e/online/*.cy.js" --config viewportWidth=375,viewportHeight=667

# Full end-to-end federation flows (fediverse chat auth, inbound follows)
# against the fake remote fediverse server. Runs while the stream is live so
# chat is fully available.
run_cypress "desktop-federation" desktop "cypress/e2e/federation/*.cy.js"

echo ""
echo "=================================================="
if [ ${#FAILED_GROUPS[@]} -gt 0 ]; then
	echo "FAILED groups: ${FAILED_GROUPS[*]}"
	if [ -f cypress/results/failures.json ]; then
		echo "Failure details (cypress/results/failures.json):"
		cat cypress/results/failures.json
	fi
	echo "Failure screenshots: cypress/screenshots/"
	exit 1
fi
echo "All browser test groups passed."
