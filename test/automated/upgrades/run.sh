#!/bin/bash

# Upgrade safety test.
#
# Proves that a data directory created by the latest published Owncast
# release survives an upgrade to the code in this checkout:
#
#   1. Download the latest GitHub release asset for the current platform.
#   2. Run it in a scratch directory on non-default ports and write
#      recognizable config values through the admin API.
#   3. Stop it, then start a build of this checkout on the same data dir.
#   4. Assert readiness, that the reported version changed to the dev
#      build's, that the previously written config survived, and that the
#      database migrated to the newest goose schema version.
#
# Requires curl, jq, unzip, sqlite3, go, and ffmpeg in PATH.
# Exits nonzero on the first failed assertion. Server logs are kept in
# ./scratch for post-mortem (CI uploads them on failure).

set -euo pipefail

cd "$(dirname "$0")" || exit 1
REPO_ROOT=$(cd ../../.. && pwd)
SCRATCH="$PWD/scratch"
INSTANCE="$SCRATCH/instance"

# Non-default ports: 8080/1935 are used by dev instances and sibling tests.
WEB_PORT="${UPGRADE_TEST_WEB_PORT:-8098}"
RTMP_PORT="${UPGRADE_TEST_RTMP_PORT:-1998}"
BASE_URL="http://127.0.0.1:${WEB_PORT}"
ADMIN_AUTH="admin:abc123" # Owncast's default admin credentials.

TEST_NAME="upgrade-test-$$"
TEST_SUMMARY="Written by the pre-upgrade release binary ($$)."

SERVER_PID=""
PASS_COUNT=0

pass() {
	PASS_COUNT=$((PASS_COUNT + 1))
	echo "PASS: $1"
}

fail() {
	echo "FAIL: $1" >&2
	exit 1
}

# assert_eq <description> <expected> <actual>
assert_eq() {
	if [ "$2" = "$3" ]; then
		pass "$1: '$3'"
	else
		fail "$1: expected '$2', got '$3'"
	fi
}

cleanup() {
	status=$?
	if [ -n "$SERVER_PID" ] && kill -0 "$SERVER_PID" 2>/dev/null; then
		kill "$SERVER_PID" 2>/dev/null || true
		wait "$SERVER_PID" 2>/dev/null || true
	fi
	echo "--------------------------------------------"
	if [ "$status" -eq 0 ]; then
		echo "RESULT: PASS (${PASS_COUNT} assertions)"
	else
		echo "RESULT: FAIL after ${PASS_COUNT} passing assertions. Server logs: $SCRATCH" >&2
	fi
	exit "$status"
}
trap cleanup EXIT

require_command() {
	command -v "$1" >/dev/null 2>&1 || fail "'$1' is required but was not found in PATH"
}

require_port_free() {
	if (exec 3<>"/dev/tcp/127.0.0.1/$1") 2>/dev/null; then
		fail "port $1 is already in use; refusing to run"
	fi
}

# Authenticated GitHub API GET when GITHUB_TOKEN is set (CI runners share
# the anonymous rate limit). The token is only sent to api.github.com,
# never to the redirected asset download host. Retries cover transient
# GitHub API outages (5xx), which otherwise fail the run before the first
# assertion.
gh_api() {
	if [ -n "${GITHUB_TOKEN:-}" ]; then
		curl -sSL --fail --retry 5 --retry-all-errors -H "Authorization: Bearer ${GITHUB_TOKEN}" -H "Accept: application/vnd.github+json" "$1" ||
			fail "GitHub API request failed: $1"
	else
		curl -sSL --fail --retry 5 --retry-all-errors -H "Accept: application/vnd.github+json" "$1" ||
			fail "GitHub API request failed: $1"
	fi
}

# start_server <binary> <logfile>
# Runs the binary with $INSTANCE as its working directory so both the
# release and the dev build operate on the same ./data directory.
start_server() {
	(cd "$INSTANCE" && exec "$1" -webserverport "$WEB_PORT" -rtmpport "$RTMP_PORT") >"$2" 2>&1 &
	SERVER_PID=$!
	# Keep bash's job-control "Terminated" notice out of the test output.
	disown "$SERVER_PID"
}

# wait_ready <description>
wait_ready() {
	i=0
	while [ "$i" -lt 90 ]; do
		kill -0 "$SERVER_PID" 2>/dev/null || fail "$1 exited before becoming ready"
		if curl -sf "$BASE_URL/api/status" >/dev/null 2>&1; then
			pass "$1 is serving /api/status"
			return 0
		fi
		sleep 1
		i=$((i + 1))
	done
	fail "$1 did not become ready within 90s"
}

# stop_server <description>
stop_server() {
	kill "$SERVER_PID"
	i=0
	while [ "$i" -lt 30 ]; do
		if ! kill -0 "$SERVER_PID" 2>/dev/null; then
			SERVER_PID=""
			return 0
		fi
		sleep 1
		i=$((i + 1))
	done
	kill -9 "$SERVER_PID" 2>/dev/null || true
	fail "$1 did not exit within 30s of SIGTERM"
}

# admin_post <path> <json payload>
admin_post() {
	response=$(curl -sf -u "$ADMIN_AUTH" -X POST -H 'Content-Type: application/json' --data "$2" "$BASE_URL$1") ||
		fail "POST $1 failed"
	jq -e '.success == true' >/dev/null <<<"$response" || fail "POST $1 rejected: $response"
}

# api_get <path> <jq filter>
api_get() {
	curl -sf "$BASE_URL$1" | jq -r "$2"
}

echo "--------------------------------------------"
echo "Owncast upgrade test: latest release -> local checkout"
echo "--------------------------------------------"

for cmd in curl jq unzip sqlite3 go ffmpeg; do
	require_command "$cmd"
done
require_port_free "$WEB_PORT"
require_port_free "$RTMP_PORT"

rm -rf "$SCRATCH"
mkdir -p "$INSTANCE"

case "$(uname -s)/$(uname -m)" in
	Linux/x86_64) ASSET_SUFFIX="linux-64bit" ;;
	Linux/aarch64 | Linux/arm64) ASSET_SUFFIX="linux-arm64" ;;
	Darwin/arm64) ASSET_SUFFIX="macOS-arm64" ;;
	Darwin/x86_64) ASSET_SUFFIX="macOS-64bit" ;;
	*) fail "unsupported platform: $(uname -s)/$(uname -m)" ;;
esac

echo "Building the local checkout..."
(cd "$REPO_ROOT" && CGO_ENABLED=1 go build -o "$SCRATCH/owncast-dev" .)

echo "Looking up the latest published release..."
RELEASE_JSON=$(gh_api "https://api.github.com/repos/owncast/owncast/releases/latest")
RELEASE_TAG=$(jq -r '.tag_name' <<<"$RELEASE_JSON")
RELEASE_VERSION="${RELEASE_TAG#v}"
ASSET_URL=$(jq -r --arg suffix "-${ASSET_SUFFIX}.zip" \
	'[.assets[] | select(.name | endswith($suffix)) | .browser_download_url][0] // empty' <<<"$RELEASE_JSON")
[ -n "$ASSET_URL" ] || fail "release ${RELEASE_TAG} has no asset ending in -${ASSET_SUFFIX}.zip"
pass "release ${RELEASE_TAG} has asset $(basename "$ASSET_URL")"

echo "Downloading $(basename "$ASSET_URL")..."
curl -sSL --fail --retry 5 --retry-all-errors "$ASSET_URL" --output "$SCRATCH/release.zip" ||
	fail "release asset download failed: $ASSET_URL"
unzip -oq "$SCRATCH/release.zip" -d "$INSTANCE"
[ -x "$INSTANCE/owncast" ] || fail "release zip did not contain an executable owncast binary"

echo "Starting release ${RELEASE_TAG} on ports ${WEB_PORT}/${RTMP_PORT}..."
start_server "$INSTANCE/owncast" "$SCRATCH/release.log"
wait_ready "release ${RELEASE_TAG}"
assert_eq "release reports its own version" "$RELEASE_VERSION" "$(api_get /api/status .versionNumber)"

echo "Writing config through the admin API..."
admin_post /api/admin/config/name "{\"value\": \"$TEST_NAME\"}"
admin_post /api/admin/config/serversummary "{\"value\": \"$TEST_SUMMARY\"}"
assert_eq "release persisted the server name" "$TEST_NAME" "$(api_get /api/config .name)"
assert_eq "release persisted the server summary" "$TEST_SUMMARY" "$(api_get /api/config .summary)"

echo "Stopping the release binary..."
stop_server "release ${RELEASE_TAG}"
[ -f "$INSTANCE/data/owncast.db" ] || fail "release did not create data/owncast.db"

echo "Starting the dev build on the release's data directory..."
start_server "$SCRATCH/owncast-dev" "$SCRATCH/dev.log"
wait_ready "dev build (upgraded server)"

DEV_EXPECTED_VERSION=$(sed -n 's/.*StaticVersionNumber = "\([^"]*\)".*/\1/p' "$REPO_ROOT/config/constants.go")
[ -n "$DEV_EXPECTED_VERSION" ] || fail "could not read StaticVersionNumber from config/constants.go"
DEV_VERSION=$(api_get /api/status .versionNumber)
assert_eq "upgraded server reports the dev version" "$DEV_EXPECTED_VERSION" "$DEV_VERSION"
if [ "$DEV_VERSION" = "$RELEASE_VERSION" ]; then
	fail "version did not change across the upgrade (still ${RELEASE_VERSION})"
fi
pass "version changed across the upgrade: ${RELEASE_VERSION} -> ${DEV_VERSION}"

assert_eq "server name survived the upgrade" "$TEST_NAME" "$(api_get /api/config .name)"
assert_eq "server summary survived the upgrade" "$TEST_SUMMARY" "$(api_get /api/config .summary)"
assert_eq "admin API works after the upgrade" "$TEST_NAME" \
	"$(curl -sf -u "$ADMIN_AUTH" "$BASE_URL/api/admin/serverconfig" | jq -r '.instanceDetails.name')"

echo "Stopping the dev build..."
stop_server "dev build"

# The goose migration table is the queryable record of schema state: its max
# version_id must equal the newest migration shipped in this checkout.
EXPECTED_SCHEMA=$(find "$REPO_ROOT/persistence/migrations" -name '[0-9]*.sql' |
	sed 's|.*/0*\([0-9][0-9]*\)_.*|\1|' | sort -n | tail -n 1)
ACTUAL_SCHEMA=$(sqlite3 "$INSTANCE/data/owncast.db" 'SELECT COALESCE(MAX(version_id), 0) FROM goose_db_version;')
assert_eq "database migrated to the newest goose schema version" "$EXPECTED_SCHEMA" "$ACTUAL_SCHEMA"
