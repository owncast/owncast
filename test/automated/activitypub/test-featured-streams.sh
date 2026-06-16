#!/bin/bash
# shellcheck disable=SC2317  # cleanup() is invoked via trap, not direct call

# Featured Streams (mini-directory) Test
#
# This test verifies the Owncast-to-Owncast "featured streams" federation
# flow end to end. It is the only test in this directory that runs TWO
# Owncast instances, because the federated-servers directory is an
# Owncast-to-Owncast feature: one instance follows another, the follow is
# accepted, and the followed server then shows up in the follower's
# directory listing.
#
# Flow:
# 1. Build Owncast once, start two instances (owncast.local, owncast2.local)
#    behind the shared HTTPS proxy, both with federation enabled and public.
# 2. Instance 1 adds Instance 2 via POST /api/admin/federation/servers.
# 3. Verify Instance 2 immediately appears in Instance 1's directory listing
#    as a pending follow (regression guard: the record must be persisted).
# 4. Verify the follow transitions to "accepted" once Instance 2 returns its
#    ActivityPub Accept (regression guard: the Accept must be matched to the
#    stored record).
# 5. Verify the accepted record is actually populated with the remote server's
#    name, display name and logo (regression guard: a green "accepted" with a
#    blank, image-less row is the broken state users hit).
# 6. Verify the listing exposes the documented field names the web consumes
#    (regression guard: the API and web previously drifted onto different
#    field names and the feature shipped broken but green).
# 7. Verify the listing is readable on the PUBLIC (unauthenticated) endpoint.
# 8. Verify the reverse direction (Instance 2 adds Instance 1) also works.
# 9. Stream a real test video into Instance 2 and verify its entry in
#    Instance 1's directory flips to live with the stream title, then stop the
#    stream and verify it flips back to offline promptly (the core of the
#    feature, both directions). Streams automatically under CI; prompts when
#    run interactively.
#
# Requirements:
# - Go, C compiler (for building Owncast)
# - Caddy, mkcert
# - curl, jq
# - ffmpeg and a Sans font (for the live-status test stream; the Docker image
#   installs both, see run.sh)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(git rev-parse --show-toplevel)"

# Configuration
PROXY_PORT="${PROXY_PORT:-8443}"

OWNCAST_PORT="${OWNCAST_PORT:-8080}"
OWNCAST_HOSTNAME="owncast.local"
OWNCAST_FED_USERNAME="streamerone"

OWNCAST2_PORT="${OWNCAST2_PORT:-8081}"
OWNCAST2_HOSTNAME="owncast2.local"
OWNCAST2_FED_USERNAME="streamertwo"

# RTMP ports must differ so both instances can bind their stream servers.
OWNCAST_RTMP_PORT="${OWNCAST_RTMP_PORT:-1935}"
OWNCAST2_RTMP_PORT="${OWNCAST2_RTMP_PORT:-1936}"

# snac2 isn't used by this test, but the shared Caddyfile references
# SNAC_PORT, so give it a value to keep the config valid.
SNAC_PORT="${SNAC_PORT:-9080}"

ADMIN_USER="admin"
ADMIN_PASS="abc123"
# Default stream key for a fresh Owncast instance (config/defaults.go).
STREAM_KEY="${STREAM_KEY:-abc123}"
CI="${CI:-false}"

# URLs (HTTPS via proxy)
OWNCAST_URL="https://${OWNCAST_HOSTNAME}:${PROXY_PORT}"
OWNCAST2_URL="https://${OWNCAST2_HOSTNAME}:${PROXY_PORT}"

# Directories and state
TEMP_DIR=""
OWNCAST_BIN=""
PROXY_PID=""
OC1_PID=""
OC2_PID=""
OC_LAST_PID=""
TEST_STREAM_PID=""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_test() { echo -e "${CYAN}[TEST]${NC} $1"; }

# shellcheck disable=SC2329  # invoked via trap, not called directly
cleanup() {
    log_info "Cleaning up..."

    # Stop any test stream first so ffmpeg releases the RTMP connection before
    # the Owncast instances are torn down. Kill the ffmpeg child directly too,
    # since the ocTestStream wrapper does not forward signals to it.
    pkill -f "rtmp://127.0.0.1:${OWNCAST2_RTMP_PORT}/live/" 2>/dev/null || true
    if [[ -n "${TEST_STREAM_PID}" ]] && kill -0 "${TEST_STREAM_PID}" 2>/dev/null; then
        kill "${TEST_STREAM_PID}" 2>/dev/null || true
        wait "${TEST_STREAM_PID}" 2>/dev/null || true
    fi

    for pid_var in PROXY_PID OC2_PID OC1_PID; do
        local pid="${!pid_var}"
        if [[ -n "${pid}" ]] && kill -0 "${pid}" 2>/dev/null; then
            kill "${pid}" 2>/dev/null || true
            wait "${pid}" 2>/dev/null || true
        fi
    done

    if [[ -n "${TEMP_DIR}" ]] && [[ -d "${TEMP_DIR}" ]]; then
        rm -rf "${TEMP_DIR}"
    fi

    log_info "Cleanup complete."
}

trap cleanup EXIT

# ==========================
# Setup
# ==========================

setup_temp_dir() {
    TEMP_DIR=$(mktemp -d)
    log_info "Temp directory: ${TEMP_DIR}"
}

check_hosts_entry() {
    local missing=0
    for host in "${OWNCAST_HOSTNAME}" "${OWNCAST2_HOSTNAME}"; do
        if ! grep -q "${host}" /etc/hosts 2>/dev/null; then
            log_error "Missing /etc/hosts entry for ${host}"
            missing=1
        fi
    done
    if [[ ${missing} -ne 0 ]]; then
        log_error "Please add: 127.0.0.1 ${OWNCAST_HOSTNAME} ${OWNCAST2_HOSTNAME}"
        exit 1
    fi
    log_info "Hosts entries verified"
}

check_certs() {
    CERT_DIR="${CERT_DIR:-${SCRIPT_DIR}/certs}"
    if [[ ! -f "${CERT_DIR}/cert.pem" ]] || [[ ! -f "${CERT_DIR}/key.pem" ]]; then
        log_error "Certificates not found in ${CERT_DIR}"
        exit 1
    fi
    log_info "Using certificates from ${CERT_DIR}"
}

start_proxy() {
    log_info "Starting HTTPS reverse proxy (Caddy)..."

    if ! command -v caddy &> /dev/null; then
        log_error "Caddy is not installed."
        return 1
    fi

    export PROXY_PORT="${PROXY_PORT}"
    export OWNCAST_PORT="${OWNCAST_PORT}"
    export OWNCAST2_PORT="${OWNCAST2_PORT}"
    export SNAC_PORT="${SNAC_PORT}"
    export CERT_FILE="${CERT_DIR}/cert.pem"
    export KEY_FILE="${CERT_DIR}/key.pem"

    local caddy_log="${TEMP_DIR}/caddy.log"
    caddy run --config "${SCRIPT_DIR}/Caddyfile" --adapter caddyfile > "${caddy_log}" 2>&1 &
    PROXY_PID=$!

    log_info "Caddy started with PID ${PROXY_PID}"
    sleep 2

    if ! kill -0 "${PROXY_PID}" 2>/dev/null; then
        log_error "Caddy failed to start. Check ${caddy_log}"
        cat "${caddy_log}"
        return 1
    fi

    local max_attempts=10
    local attempt=0
    while [[ ${attempt} -lt ${max_attempts} ]]; do
        if curl -sk "https://127.0.0.1:${PROXY_PORT}/" > /dev/null 2>&1; then
            log_info "Caddy proxy is ready"
            return 0
        fi
        attempt=$((attempt + 1))
        sleep 1
    done

    log_error "Caddy proxy did not become ready"
    return 1
}

build_owncast() {
    log_info "Building Owncast..."
    OWNCAST_BIN="${TEMP_DIR}/owncast"
    pushd "${REPO_ROOT}" > /dev/null
    CGO_ENABLED=1 go build -o "${OWNCAST_BIN}" main.go
    popd > /dev/null
    log_info "Owncast built: ${OWNCAST_BIN}"
}

# start_owncast_instance LABEL WEB_PORT RTMP_PORT
#
# Each instance runs from its own working directory so it gets a private
# ./data directory (the data dir path is a compile-time constant relative to
# the process CWD and can't be overridden by a flag). The database is kept in
# the shared temp dir and selected explicitly. The resulting PID is returned
# in the global OC_LAST_PID.
start_owncast_instance() {
    local label=$1
    local web_port=$2
    local rtmp_port=$3
    local workdir="${TEMP_DIR}/${label}"
    local db="${TEMP_DIR}/${label}.db"
    local log="${TEMP_DIR}/${label}.log"

    mkdir -p "${workdir}"

    log_info "Starting Owncast '${label}' (web ${web_port}, rtmp ${rtmp_port})..."
    (
        cd "${workdir}" || exit 1
        exec env \
            OWNCAST_ALLOW_INTERNAL_FEDERATION=true \
            OWNCAST_INSECURE_SKIP_VERIFY=true \
            "${OWNCAST_BIN}" \
            -database "${db}" \
            -webserverport "${web_port}" \
            -rtmpport "${rtmp_port}" \
            -enableVerboseLogging
    ) > "${log}" 2>&1 &
    OC_LAST_PID=$!
    log_info "Owncast '${label}' PID ${OC_LAST_PID} (log: ${log})"

    local max_attempts=30
    local attempt=0
    while [[ ${attempt} -lt ${max_attempts} ]]; do
        if curl -s "http://localhost:${web_port}/api/status" > /dev/null 2>&1; then
            log_info "Owncast '${label}' is ready"
            return 0
        fi
        if ! kill -0 "${OC_LAST_PID}" 2>/dev/null; then
            log_error "Owncast '${label}' exited early. Log:"
            cat "${log}"
            return 1
        fi
        attempt=$((attempt + 1))
        sleep 1
    done

    log_error "Owncast '${label}' did not become ready. Log tail:"
    tail -50 "${log}"
    return 1
}

get_admin_auth() {
    echo -n "${ADMIN_USER}:${ADMIN_PASS}" | base64
}

# configure_owncast WEB_PORT SERVER_URL FED_USERNAME
configure_owncast() {
    local web_port=$1
    local server_url=$2
    local fed_username=$3
    local base_url="http://localhost:${web_port}"
    local auth
    auth=$(get_admin_auth)

    log_info "Configuring Owncast at ${base_url} as ${server_url}"

    curl -s -X POST "${base_url}/api/admin/config/serverurl" \
        -H "Authorization: Basic ${auth}" -H "Content-Type: application/json" \
        -d "{\"value\": \"${server_url}\"}" > /dev/null

    curl -s -X POST "${base_url}/api/admin/config/federation/username" \
        -H "Authorization: Basic ${auth}" -H "Content-Type: application/json" \
        -d "{\"value\": \"${fed_username}\"}" > /dev/null

    curl -s -X POST "${base_url}/api/admin/config/federation/enable" \
        -H "Authorization: Basic ${auth}" -H "Content-Type: application/json" \
        -d '{"value": true}' > /dev/null

    # Public federation so inbound follows are auto-accepted (the directory
    # flow relies on the Accept being returned without manual approval).
    curl -s -X POST "${base_url}/api/admin/config/federation/private" \
        -H "Authorization: Basic ${auth}" -H "Content-Type: application/json" \
        -d '{"value": false}' > /dev/null
}

# ==========================
# API helpers
# ==========================

# add_featured_server WEB_PORT TARGET_URL -> prints JSON response body
add_featured_server() {
    local web_port=$1
    local target_url=$2
    local auth
    auth=$(get_admin_auth)

    curl -s -X POST "http://localhost:${web_port}/api/admin/federation/servers" \
        -H "Authorization: Basic ${auth}" \
        -H "Content-Type: application/json" \
        -d "{\"url\": \"${target_url}\"}"
}

# get_featured_servers WEB_PORT [auth] -> prints JSON body
# When the second argument is "admin" the request is authenticated; otherwise
# it hits the public, unauthenticated endpoint.
get_featured_servers() {
    local web_port=$1
    local mode="${2:-public}"

    if [[ "${mode}" == "admin" ]]; then
        local auth
        auth=$(get_admin_auth)
        curl -s "http://localhost:${web_port}/api/federation/servers" \
            -H "Authorization: Basic ${auth}"
    else
        curl -s "http://localhost:${web_port}/api/federation/servers"
    fi
}

# server_has_iri JSON IRI -> exit 0 if a server with that iri is present
server_has_iri() {
    echo "$1" | jq -e --arg iri "$2" '.servers[]? | select(.iri == $iri)' > /dev/null 2>&1
}

# server_follow_status JSON IRI -> prints followStatus (or empty)
server_follow_status() {
    echo "$1" | jq -r --arg iri "$2" '.servers[]? | select(.iri == $iri) | .followStatus' 2>/dev/null
}

# server_object JSON IRI -> prints the compact server object (or empty)
server_object() {
    echo "$1" | jq -c --arg iri "$2" '.servers[]? | select(.iri == $iri)' 2>/dev/null
}

# server_field JSON IRI FIELD -> prints the value of FIELD (or empty).
# Note: we deliberately avoid `// empty` here because jq treats a literal
# false as "absent", which would turn isOnline=false into an empty string and
# make an offline server indistinguishable from a missing field. Map only
# null/missing to empty, preserving false/true and string values.
server_field() {
    echo "$1" | jq -r --arg iri "$2" --arg field "$3" \
        '.servers[]? | select(.iri == $iri) | .[$field] | if . == null then empty else . end' 2>/dev/null
}

# wait_for_follow_status WEB_PORT IRI EXPECTED_STATUS TIMEOUT_SECONDS
wait_for_follow_status() {
    local web_port=$1
    local iri=$2
    local expected=$3
    local timeout=${4:-30}
    local waited=0

    while [[ ${waited} -lt ${timeout} ]]; do
        local json status
        json=$(get_featured_servers "${web_port}" admin)
        status=$(server_follow_status "${json}" "${iri}")
        if [[ "${status}" == "${expected}" ]]; then
            return 0
        fi
        sleep 2
        waited=$((waited + 2))
    done
    return 1
}

# ==========================
# Test scenarios
# ==========================

test_add_persists_pending_record() {
    log_test "TEST 1: Adding a featured server persists a directory record"

    local response success message
    response=$(add_featured_server "${OWNCAST_PORT}" "${OWNCAST2_URL}")
    log_info "Add response: ${response}"

    success=$(echo "${response}" | jq -r '.success // false' 2>/dev/null)
    if [[ "${success}" != "true" ]]; then
        message=$(echo "${response}" | jq -r '.message // ""' 2>/dev/null)
        log_error "TEST 1 FAILED: add request was not accepted: ${message}"
        return 1
    fi

    # The record must exist in the listing immediately, before any Accept
    # round-trips back. This is the core regression guard: previously nothing
    # was ever written to the federated_servers table, so the listing stayed
    # empty even though the Follow had been sent.
    local json
    json=$(get_featured_servers "${OWNCAST_PORT}" admin)
    if ! server_has_iri "${json}" "${OWNCAST2_URL}"; then
        log_error "TEST 1 FAILED: ${OWNCAST2_URL} not present in directory after add"
        log_error "Listing: ${json}"
        return 1
    fi

    log_test "TEST 1 PASSED: directory record persisted for ${OWNCAST2_URL}"
    return 0
}

test_follow_is_accepted() {
    log_test "TEST 2: Follow transitions to accepted after remote Accept"

    if wait_for_follow_status "${OWNCAST_PORT}" "${OWNCAST2_URL}" "accepted" 40; then
        log_test "TEST 2 PASSED: follow to ${OWNCAST2_URL} was accepted"
        return 0
    fi

    local json status
    json=$(get_featured_servers "${OWNCAST_PORT}" admin)
    status=$(server_follow_status "${json}" "${OWNCAST2_URL}")
    log_error "TEST 2 FAILED: follow status is '${status}', expected 'accepted'"
    log_error "Listing: ${json}"
    return 1
}

test_metadata_is_populated() {
    log_test "TEST 3: Remote server metadata (name/logo) is populated after accept"

    # Once the Accept comes back, the follower resolves the remote actor and
    # stores its username, display name and logo. If this step is skipped the
    # directory entry is a blank row with no name and no image -- exactly the
    # broken state users reported. A green "follow accepted" is not enough;
    # the row has to actually carry the data the UI renders.
    local json name displayName logoUrl
    json=$(get_featured_servers "${OWNCAST_PORT}" admin)

    name=$(server_field "${json}" "${OWNCAST2_URL}" name)
    displayName=$(server_field "${json}" "${OWNCAST2_URL}" displayName)
    logoUrl=$(server_field "${json}" "${OWNCAST2_URL}" logoUrl)

    local failed=0
    if [[ -z "${name}" ]]; then
        log_error "TEST 3 FAILED: 'name' is empty; the Accept handler did not store the remote username"
        failed=1
    fi
    if [[ -z "${displayName}" ]]; then
        log_error "TEST 3 FAILED: 'displayName' is empty; the remote server's display name was not stored"
        failed=1
    fi
    if [[ -z "${logoUrl}" ]]; then
        log_error "TEST 3 FAILED: 'logoUrl' is empty; the featured stream would render with no image"
        failed=1
    fi

    if [[ ${failed} -ne 0 ]]; then
        log_error "Server record: $(server_object "${json}" "${OWNCAST2_URL}")"
        return 1
    fi

    log_test "TEST 3 PASSED: name='${name}', displayName='${displayName}', logoUrl='${logoUrl}'"
    return 0
}

test_listing_field_contract() {
    log_test "TEST 4: Directory listing exposes the documented field contract"

    # The web reads these field names verbatim (web/hooks/useFederatedServers.tsx).
    # The whole feature shipped broken once because the API and the web had
    # drifted onto different names, so guard the contract here: the documented
    # names must be present and the legacy names must never come back.
    local json server
    json=$(get_featured_servers "${OWNCAST_PORT}" public)
    server=$(server_object "${json}" "${OWNCAST2_URL}")

    if [[ -z "${server}" ]]; then
        log_error "TEST 4 FAILED: ${OWNCAST2_URL} not present in public listing"
        log_error "Listing: ${json}"
        return 1
    fi

    local failed=0 field
    local required=(iri name displayName logoUrl isOnline addedAt followStatus)
    for field in "${required[@]}"; do
        if ! echo "${server}" | jq -e "has(\"${field}\")" > /dev/null 2>&1; then
            log_error "TEST 4 FAILED: response is missing documented field '${field}'"
            failed=1
        fi
    done

    local forbidden=(url logo thumbnail lastChecked)
    for field in "${forbidden[@]}"; do
        if echo "${server}" | jq -e "has(\"${field}\")" > /dev/null 2>&1; then
            log_error "TEST 4 FAILED: response contains legacy field '${field}' the web no longer reads"
            failed=1
        fi
    done

    if [[ ${failed} -ne 0 ]]; then
        log_error "Server object: ${server}"
        return 1
    fi

    log_test "TEST 4 PASSED: listing exposes the documented field contract"
    return 0
}

test_public_listing_is_readable() {
    log_test "TEST 5: Directory listing is readable without admin auth"

    # The public watch page fetches this endpoint on load; it must not be
    # gated behind admin basic auth.
    local code
    code=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:${OWNCAST_PORT}/api/federation/servers")
    if [[ "${code}" != "200" ]]; then
        log_error "TEST 5 FAILED: public listing returned HTTP ${code}, expected 200"
        return 1
    fi

    local json
    json=$(get_featured_servers "${OWNCAST_PORT}" public)
    if ! server_has_iri "${json}" "${OWNCAST2_URL}"; then
        log_error "TEST 5 FAILED: ${OWNCAST2_URL} not present in public listing"
        log_error "Listing: ${json}"
        return 1
    fi

    log_test "TEST 5 PASSED: public listing returns the featured server"
    return 0
}

test_reverse_direction() {
    log_test "TEST 6: Reverse direction (instance 2 adds instance 1)"

    local response success message
    response=$(add_featured_server "${OWNCAST2_PORT}" "${OWNCAST_URL}")
    log_info "Add response: ${response}"

    success=$(echo "${response}" | jq -r '.success // false' 2>/dev/null)
    if [[ "${success}" != "true" ]]; then
        message=$(echo "${response}" | jq -r '.message // ""' 2>/dev/null)
        log_error "TEST 6 FAILED: add request was not accepted: ${message}"
        return 1
    fi

    if wait_for_follow_status "${OWNCAST2_PORT}" "${OWNCAST_URL}" "accepted" 40; then
        log_test "TEST 6 PASSED: instance 2 followed and accepted instance 1"
        return 0
    fi

    local json status
    json=$(get_featured_servers "${OWNCAST2_PORT}" admin)
    status=$(server_follow_status "${json}" "${OWNCAST_URL}")
    log_error "TEST 6 FAILED: follow status is '${status}', expected 'accepted'"
    log_error "Listing: ${json}"
    return 1
}

# stop_test_stream stops the running test stream, if any.
#
# ocTestStream.sh runs ffmpeg in the foreground without forwarding signals or
# trapping EXIT (for the internal test-pattern path), so killing the wrapper
# alone orphans ffmpeg and the RTMP stream keeps running -- which would leave
# the source server live and no Leave would ever be sent. Kill the ffmpeg
# pushing to instance 2's RTMP endpoint directly as well.
stop_test_stream() {
    pkill -f "rtmp://127.0.0.1:${OWNCAST2_RTMP_PORT}/live/" 2>/dev/null || true
    if [[ -n "${TEST_STREAM_PID}" ]] && kill -0 "${TEST_STREAM_PID}" 2>/dev/null; then
        kill "${TEST_STREAM_PID}" 2>/dev/null || true
        wait "${TEST_STREAM_PID}" 2>/dev/null || true
    fi
    TEST_STREAM_PID=""
}

test_live_status_flip() {
    log_test "TEST 7: Featured stream flips to live when the remote goes online"

    # This is the heart of the feature: a featured server must show as live,
    # with its stream metadata, while it is actually streaming. Instance 1
    # already follows instance 2 (TEST 2), so when instance 2 goes live it
    # sends an immediate Offer ping that should flip instance 2's row in
    # instance 1's directory to online.
    #
    # In CI we stream automatically; run interactively and we ask first so a
    # developer can decline (or watch it happen).
    if [[ "${CI}" != "true" ]]; then
        echo ""
        read -p "Start a test stream on instance 2 to verify live-status propagation? [Y/n] " -n 1 -r
        echo ""
        if [[ "${REPLY}" =~ ^[Nn]$ ]]; then
            log_warn "TEST 7 SKIPPED: declined to start a test stream"
            return 0
        fi
    fi

    if [[ ! -x "${REPO_ROOT}/test/ocTestStream.sh" ]]; then
        log_error "TEST 7 FAILED: ${REPO_ROOT}/test/ocTestStream.sh not found or not executable"
        return 1
    fi

    # Give instance 2 a known stream title so we can assert it propagates with
    # the live status, rather than just trusting the boolean flag.
    local expected_title="Featured Streams Live Test"
    local auth
    auth=$(get_admin_auth)
    curl -s -X POST "http://localhost:${OWNCAST2_PORT}/api/admin/config/streamtitle" \
        -H "Authorization: Basic ${auth}" -H "Content-Type: application/json" \
        -d "{\"value\": \"${expected_title}\"}" > /dev/null

    log_info "Starting test stream into instance 2 (rtmp port ${OWNCAST2_RTMP_PORT})..."
    "${REPO_ROOT}/test/ocTestStream.sh" "rtmp://127.0.0.1:${OWNCAST2_RTMP_PORT}/live/${STREAM_KEY}" \
        > "${TEMP_DIR}/teststream.log" 2>&1 &
    TEST_STREAM_PID=$!

    # Wait for instance 1 to observe instance 2 as online. The Offer fires on
    # RTMP connect, but allow generous time for ffmpeg startup and delivery.
    local timeout=90 waited=0 online="" json
    while [[ ${waited} -lt ${timeout} ]]; do
        if ! kill -0 "${TEST_STREAM_PID}" 2>/dev/null; then
            log_error "TEST 7 FAILED: test stream process exited early. ffmpeg log:"
            cat "${TEMP_DIR}/teststream.log"
            TEST_STREAM_PID=""
            return 1
        fi
        json=$(get_featured_servers "${OWNCAST_PORT}" admin)
        online=$(server_field "${json}" "${OWNCAST2_URL}" isOnline)
        if [[ "${online}" == "true" ]]; then
            break
        fi
        sleep 3
        waited=$((waited + 3))
    done

    if [[ "${online}" != "true" ]]; then
        log_error "TEST 7 FAILED: instance 2 never showed as online on instance 1 within ${timeout}s"
        log_error "Server record: $(server_object "${json}" "${OWNCAST2_URL}")"
        stop_test_stream
        return 1
    fi

    # The Offer carries stream metadata; the title we set must have propagated.
    local title
    title=$(server_field "$(get_featured_servers "${OWNCAST_PORT}" admin)" "${OWNCAST2_URL}" streamTitle)
    if [[ "${title}" != "${expected_title}" ]]; then
        log_error "TEST 7 FAILED: streamTitle is '${title}', expected '${expected_title}'"
        stop_test_stream
        return 1
    fi

    log_info "Instance 2 is live on instance 1 with title '${title}'; stopping the stream..."
    stop_test_stream

    # Ending the stream makes instance 2 send a Leave activity; instance 1 must
    # flip the entry back to offline promptly, well before the 20-minute
    # staleness sweep would otherwise time it out. The latency floor here is
    # how long instance 2 takes to detect the RTMP disconnect and transition
    # itself offline, so allow generous headroom and report the measured time.
    local offline_timeout=180 offline_waited=0
    online="true"
    while [[ ${offline_waited} -lt ${offline_timeout} ]]; do
        json=$(get_featured_servers "${OWNCAST_PORT}" admin)
        online=$(server_field "${json}" "${OWNCAST2_URL}" isOnline)
        if [[ "${online}" == "false" ]]; then
            break
        fi
        sleep 3
        offline_waited=$((offline_waited + 3))
    done

    if [[ "${online}" != "false" ]]; then
        log_error "TEST 7 FAILED: instance 2 still shows online on instance 1 ${offline_timeout}s after the stream stopped"
        log_error "Server record: $(server_object "${json}" "${OWNCAST2_URL}")"

        # Diagnostics: did instance 2 itself go offline, and did the Leave flow?
        local oc2_self
        oc2_self=$(curl -s "http://localhost:${OWNCAST2_PORT}/api/status" 2>/dev/null | jq -c '{online}' 2>/dev/null)
        log_error "Instance 2 self-reported status: ${oc2_self:-<none>}"
        log_error "--- instance 2 log (disconnect / leave / ping) ---"
        grep -iE "disconnect|leave|offer|ping|offline|transcoder complet|federat" "${TEMP_DIR}/owncast2.log" 2>/dev/null | tail -25 >&2 || true
        log_error "--- instance 1 log (inbox / leave) ---"
        grep -iE "leave|offer|offline|federated server|inbox" "${TEMP_DIR}/owncast1.log" 2>/dev/null | tail -25 >&2 || true
        return 1
    fi

    log_test "TEST 7 PASSED: went live with title '${title}', then back offline ~${offline_waited}s after the stream ended"
    return 0
}

# ==========================
# Results
# ==========================

print_results() {
    local passed=$1
    local failed=$2
    local total=$((passed + failed))

    echo ""
    echo "========================================"
    echo "Featured Streams Test Results"
    echo "========================================"
    echo "Tests Run:    ${total}"
    echo "Tests Passed: ${passed}"
    echo "Tests Failed: ${failed}"
    echo ""

    if [[ "${failed}" -eq 0 ]]; then
        echo -e "${GREEN}ALL TESTS PASSED${NC}"
    else
        echo -e "${RED}${failed} TEST(S) FAILED${NC}"
    fi
    echo "========================================"
}

# ==========================
# Main
# ==========================

main() {
    echo ""
    echo "========================================"
    echo "Featured Streams (mini-directory) Test"
    echo "========================================"
    echo "Verifies the Owncast-to-Owncast featured"
    echo "streams directory flow across two servers."
    echo ""

    local passed=0
    local failed=0

    echo "----------------------------------------"
    echo "STEP 1: Setup proxy and certificates"
    echo "----------------------------------------"
    setup_temp_dir
    check_hosts_entry
    check_certs
    start_proxy
    echo ""

    echo "----------------------------------------"
    echo "STEP 2: Build and start two Owncast instances"
    echo "----------------------------------------"
    build_owncast

    # OC1_PID/OC2_PID are read indirectly in cleanup() via ${!pid_var},
    # which shellcheck cannot trace, hence the SC2034 suppressions.
    start_owncast_instance "owncast1" "${OWNCAST_PORT}" "${OWNCAST_RTMP_PORT}"
    # shellcheck disable=SC2034
    OC1_PID="${OC_LAST_PID}"
    start_owncast_instance "owncast2" "${OWNCAST2_PORT}" "${OWNCAST2_RTMP_PORT}"
    # shellcheck disable=SC2034
    OC2_PID="${OC_LAST_PID}"

    configure_owncast "${OWNCAST_PORT}" "${OWNCAST_URL}" "${OWNCAST_FED_USERNAME}"
    configure_owncast "${OWNCAST2_PORT}" "${OWNCAST2_URL}" "${OWNCAST2_FED_USERNAME}"
    sleep 2
    echo ""

    echo "----------------------------------------"
    echo "STEP 3: Run test scenarios"
    echo "----------------------------------------"

    if test_add_persists_pending_record; then passed=$((passed + 1)); else failed=$((failed + 1)); fi
    echo ""
    if test_follow_is_accepted; then passed=$((passed + 1)); else failed=$((failed + 1)); fi
    echo ""
    if test_metadata_is_populated; then passed=$((passed + 1)); else failed=$((failed + 1)); fi
    echo ""
    if test_listing_field_contract; then passed=$((passed + 1)); else failed=$((failed + 1)); fi
    echo ""
    if test_public_listing_is_readable; then passed=$((passed + 1)); else failed=$((failed + 1)); fi
    echo ""
    if test_reverse_direction; then passed=$((passed + 1)); else failed=$((failed + 1)); fi
    echo ""
    if test_live_status_flip; then passed=$((passed + 1)); else failed=$((failed + 1)); fi
    echo ""

    print_results "${passed}" "${failed}"

    if [[ "${KEEP_RUNNING:-}" == "true" ]]; then
        log_info "Keeping servers running (Ctrl+C to stop)..."
        log_info "  Owncast 1: http://localhost:${OWNCAST_PORT} (${OWNCAST_URL})"
        log_info "  Owncast 2: http://localhost:${OWNCAST2_PORT} (${OWNCAST2_URL})"
        wait
    fi

    if [[ "${failed}" -gt 0 ]]; then
        exit 1
    fi

    exit 0
}

main "$@"
