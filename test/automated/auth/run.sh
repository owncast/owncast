#!/bin/bash

# Standalone end-to-end viewer-authentication ("auth gate") integration tests.
#
# Builds two auth.gate plugins and installs them into a real Owncast instance:
#   - the SDK's `basic-auth` example (a shared-password gate), and
#   - `oauth-gate-test` (this repo's fixture: a standard OAuth2 client built
#     from scratch against the SDK), which authenticates against a real,
#     certified OpenID Connect provider (node-oidc-provider) started below.
#
# The suites verify, for each gate, that viewer endpoints are blocked (302 to
# the login screen) before auth while the admin surface and the gate plugin's
# own routes stay reachable, that authenticating issues a signed
# `owncast_session` cookie that unlocks the gated endpoints (and a tampered one
# does not), and that logout clears the cookie.
#
# Env overrides (same contract as test/automated/plugins/run.sh):
#   PLUGIN_SDK_DIR   path to an existing SDK checkout (skips the clone)
#   PLUGIN_SDK_REPO  git URL of the plugin SDK (default: owncast/plugin-sdk)
#   PLUGIN_SDK_REF   branch/tag/sha to build from (default: main)

set -e

# shellcheck disable=SC1091  # tools.sh is sourced at runtime; not available to the linter
source ../tools.sh

REPO_ROOT="$(git rev-parse --show-toplevel)"
PLUGIN_SDK_DIR="${PLUGIN_SDK_DIR:-}"
PLUGIN_SDK_REPO="${PLUGIN_SDK_REPO:-https://github.com/owncast/plugin-sdk}"
PLUGIN_SDK_REF="${PLUGIN_SDK_REF:-main}"

PLUGIN_NAME="basic-auth"
OAUTH_PLUGIN_NAME="oauth-gate-test"
PLUGIN_DIR="${REPO_ROOT}/data/plugins"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

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
	if [[ -n "${FAKE_OAUTH_PID}" ]] && kill -0 "${FAKE_OAUTH_PID}" 2>/dev/null; then
		kill "${FAKE_OAUTH_PID}" 2>/dev/null || true
	fi
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

echo "Building the ${OAUTH_PLUGIN_NAME} plugin (this repo's fixture, against the SDK)..."
cp -R "${SCRIPT_DIR}/${OAUTH_PLUGIN_NAME}" "${SDK_DIR}/examples/js/${OAUTH_PLUGIN_NAME}"
(cd "$SDK_DIR" && ./tools/build-plugin.sh "examples/js/${OAUTH_PLUGIN_NAME}")
cp "${SDK_DIR}/plugins/${OAUTH_PLUGIN_NAME}.ocpkg" "${PLUGIN_DIR}/"

echo "Installing the JS test framework..."
npm install --quiet --no-progress

install_ffmpeg
start_owncast

# Start the internal OAuth2 / OIDC provider the oauth-gate-test plugin
# authenticates against (node-oidc-provider), then wait for it to be ready.
echo "Starting the internal OAuth2 provider (node-oidc-provider)..."
node "${SCRIPT_DIR}/fake-oauth-server.mjs" &
FAKE_OAUTH_PID=$!
for _ in $(seq 1 15); do
	if curl -s "http://localhost:9876/.well-known/openid-configuration" >/dev/null 2>&1; then
		break
	fi
	sleep 1
done

# The plugins are installed (discovered) but disabled; each suite enables its
# gate in beforeAll and disables it in afterAll so a gate is only live for the
# duration of that suite.
npm test
