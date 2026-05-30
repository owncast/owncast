#!/bin/bash

# Standalone end-to-end plugin integration test.
#
# Unlike the Go unit tests, this builds plugins from source using the Owncast
# plugin SDK (so no WASM is bundled in this repo), installs them into a real
# Owncast instance, and verifies they actually run by exercising chat and HTTP
# through them. Run by CI via the Earthly `+plugin-tests` target.
#
# The plugins built here are the SDK's own examples — the suite tracks them
# rather than maintaining plugin source in this repo.
#
# Env overrides:
#   PLUGIN_SDK_DIR    path to an existing SDK checkout; if set, no clone happens
#                     (use this for local runs: PLUGIN_SDK_DIR=~/src/plugin-sdk ./run.sh)
#   PLUGIN_SDK_REPO   git URL of the plugin SDK (default: owncast/plugin-sdk)
#   PLUGIN_SDK_REF    branch/tag/sha to build from (default: main)

set -e

source ../tools.sh

REPO_ROOT="$(git rev-parse --show-toplevel)"
PLUGIN_SDK_DIR="${PLUGIN_SDK_DIR:-}"
PLUGIN_SDK_REPO="${PLUGIN_SDK_REPO:-https://github.com/owncast/plugin-sdk}"
PLUGIN_SDK_REF="${PLUGIN_SDK_REF:-main}"

# SDK example plugins exercised by the tests in this directory.
PLUGIN_NAMES=(profanity-filter echo-bot overlay styles-demo scripts-demo viewer-gate page-content-demo tabs-demo)

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
	# Never block on an interactive credential prompt — fail clearly instead.
	export GIT_TERMINAL_PROMPT=0
	git clone --depth 1 --branch "$PLUGIN_SDK_REF" "$PLUGIN_SDK_REPO" "$SDK_DIR"
fi

# Tear down the installed plugins (and the SDK clone, but never a local SDK
# checkout the caller pointed us at) on top of tools.sh's server/database
# cleanup, which `finish` performs.
plugins_finish() {
	rm -rf "$PLUGIN_DIR"
	if [[ "$CLONED_SDK" == "1" ]]; then
		rm -rf "$SDK_DIR"
	fi
	finish
}
trap plugins_finish EXIT

# Install the SDK package's own dependencies. The build CLI (owncast-plugin)
# runs from sdks/js/bin and requires esbuild/jszip resolvable from there, and
# its postinstall fetches extism-js + binaryen. build-plugin.sh only installs
# each example's deps, not the SDK package's, so do it here once up front.
echo "Installing plugin SDK build dependencies..."
(cd "${SDK_DIR}/sdks/js" && npm install --no-audit --no-fund)

echo "Building example plugins..."
rm -rf "$PLUGIN_DIR"
mkdir -p "$PLUGIN_DIR"
for name in "${PLUGIN_NAMES[@]}"; do
	echo "  building ${name}..."
	# build-plugin.sh runs `npm install` (whose postinstall fetches extism-js
	# and friends) then `npm run build`, emitting an .ocpkg into the SDK's
	# plugins/ directory.
	(cd "$SDK_DIR" && ./tools/build-plugin.sh "examples/js/${name}")
	cp "${SDK_DIR}/plugins/${name}.ocpkg" "${PLUGIN_DIR}/"
done

# Install the JS test framework for this suite.
npm install --quiet --no-progress

install_ffmpeg
start_owncast

# Run the tests against the running instance.
npm test
