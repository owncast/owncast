#!/bin/bash

# Standalone end-to-end viewer-authentication ("auth gate") integration test.
#
# Builds the SDK's `basic-auth` example (a shared-password auth.gate plugin),
# installs it into a real Owncast instance, then verifies via HTTP that:
#   - viewer endpoints are blocked (302 to the login screen) before auth,
#   - the admin API and the gate plugin's own routes stay reachable,
#   - a correct password issues a signed `owncast_session` cookie,
#   - that cookie unlocks the gated endpoints, and a tampered one does not,
#   - logout clears the cookie.
#
# Env overrides (same contract as test/automated/plugins/run.sh):
#   PLUGIN_SDK_DIR   path to an existing SDK checkout (skips the clone)
#   PLUGIN_SDK_REPO  git URL of the plugin SDK (default: owncast/plugin-sdk)
#   PLUGIN_SDK_REF   branch/tag/sha to build from (default: main)

set -e

source ../tools.sh

REPO_ROOT="$(git rev-parse --show-toplevel)"
PLUGIN_SDK_DIR="${PLUGIN_SDK_DIR:-}"
PLUGIN_SDK_REPO="${PLUGIN_SDK_REPO:-https://github.com/owncast/plugin-sdk}"
PLUGIN_SDK_REF="${PLUGIN_SDK_REF:-main}"

PLUGIN_NAME="basic-auth"
PLUGIN_DIR="${REPO_ROOT}/data/plugins"

# Resolve the SDK source: an existing local checkout, or a fresh clone.
if [[ -n "$PLUGIN_SDK_DIR" ]]; then
	echo "Using local plugin SDK at ${PLUGIN_SDK_DIR}"
	SDK_DIR="$PLUGIN_SDK_DIR"
	CLONED_SDK=0
else
	echo "Cloning the plugin SDK (${PLUGIN_SDK_REPO}@${PLUGIN_SDK_REF})..."
	SDK_DIR="$(mktemp -d)"
	CLONED_SDK=1
	export GIT_TERMINAL_PROMPT=0
	git clone "$PLUGIN_SDK_REPO" "$SDK_DIR"
	(cd "$SDK_DIR" && git checkout "$PLUGIN_SDK_REF")
fi

# Tear down installed plugins (and the SDK clone, never a caller-supplied local
# checkout) on top of tools.sh's server/database cleanup.
auth_finish() {
	rm -rf "$PLUGIN_DIR"
	if [[ "$CLONED_SDK" == "1" ]]; then
		rm -rf "$SDK_DIR"
	fi
	finish
}
trap auth_finish EXIT

echo "Installing plugin SDK build dependencies..."
(cd "${SDK_DIR}/sdks/js" && npm install --no-audit --no-fund)

echo "Building the ${PLUGIN_NAME} plugin..."
rm -rf "$PLUGIN_DIR"
mkdir -p "$PLUGIN_DIR"
(cd "$SDK_DIR" && ./tools/build-plugin.sh "examples/js/${PLUGIN_NAME}")
cp "${SDK_DIR}/plugins/${PLUGIN_NAME}.ocpkg" "${PLUGIN_DIR}/"

echo "Installing the JS test framework..."
npm install --quiet --no-progress

install_ffmpeg
start_owncast

# The plugin is installed (discovered) but disabled; the test enables it in
# its beforeAll and disables it in afterAll so the gate is only live for the
# duration of this suite.
npm test
