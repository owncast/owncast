#!/bin/bash
# shellcheck disable=SC2317  # cleanup() is invoked via trap, not direct call

# ActivityPub Federation Test using snac2
#
# This test:
# 1. Builds/installs snac2
# 2. Creates 100 test users in snac2
# 3. Starts a local HTTPS reverse proxy
# 4. Starts snac2 server (HTTP internally, HTTPS via proxy)
# 5. Starts Owncast (HTTP internally, HTTPS via proxy)
# 6. Has all 100 snac2 users follow Owncast
# 7. Sends a message from Owncast
# 8. Verifies message delivery to snac2 inboxes
#
# Requirements:
# - openssl (for generating certs)
# - curl
# - node (for the HTTPS proxy)
# - Go (for building Owncast)
# - C compiler (for building snac2)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(git rev-parse --show-toplevel)"

# Configuration
USER_COUNT="${USER_COUNT:-100}"
FOLLOW_DELAY="${FOLLOW_DELAY:-0.1}"  # Delay between follow requests (seconds)
CI="${CI:-false}"  # Skip interactive prompts for CI environments
PROXY_PORT="${PROXY_PORT:-8443}"
SNAC_PORT="${SNAC_PORT:-9080}"
SNAC_HOSTNAME="snac.local"
OWNCAST_PORT="${OWNCAST_PORT:-8080}"
OWNCAST_HOSTNAME="owncast.local"
ADMIN_USER="admin"
ADMIN_PASS="abc123"
FEDERATION_USERNAME="streamer"
CLEAR_SHARED_INBOX_PERCENT="${CLEAR_SHARED_INBOX_PERCENT:-30}"  # Percentage of followers to clear shared_inbox (0 = none)

# URLs (HTTPS via proxy)
SNAC_URL="https://${SNAC_HOSTNAME}:${PROXY_PORT}"
OWNCAST_URL="https://${OWNCAST_HOSTNAME}:${PROXY_PORT}"

# Directories
TEMP_DIR=""
SNAC_DATA_DIR=""
SNAC_BIN=""
OWNCAST_BIN=""
OWNCAST_DB=""

# PIDs and state
SNAC_PID=""
OWNCAST_PID=""
PROXY_PID=""
TEST_STREAM_PID=""
SNAC_USERNAMES=()
CONFIRMED_FOLLOWERS=()

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

kill_leftover_processes() {
    # Kill any leftover test processes from previous runs
    local killed=false

    # Kill snac2 instances running from /tmp (test instances only)
    if pkill -f "snac httpd /tmp" 2>/dev/null; then
        killed=true
    fi

    # Kill any local-proxy.js instances
    if pkill -f "local-proxy.js" 2>/dev/null; then
        killed=true
    fi

    # Kill anything on our test ports
    local proxy_pid
    proxy_pid=$(lsof -ti :"${PROXY_PORT}" 2>/dev/null) || true
    if [[ -n "${proxy_pid}" ]]; then
        kill "${proxy_pid}" 2>/dev/null || true
        killed=true
    fi

    local snac_pid
    snac_pid=$(lsof -ti :"${SNAC_PORT}" 2>/dev/null) || true
    if [[ -n "${snac_pid}" ]]; then
        kill "${snac_pid}" 2>/dev/null || true
        killed=true
    fi

    if [[ "${killed}" == "true" ]]; then
        log_info "Killed leftover processes from previous run"
        sleep 1
    fi
}

# shellcheck disable=SC2329  # invoked via trap, not called directly
cleanup() {
    log_info "Cleaning up..."

    if [[ -n "${PROXY_PID}" ]] && kill -0 "${PROXY_PID}" 2>/dev/null; then
        kill "${PROXY_PID}" 2>/dev/null || true
        wait "${PROXY_PID}" 2>/dev/null || true
    fi

    if [[ -n "${OWNCAST_PID}" ]] && kill -0 "${OWNCAST_PID}" 2>/dev/null; then
        kill "${OWNCAST_PID}" 2>/dev/null || true
        wait "${OWNCAST_PID}" 2>/dev/null || true
    fi

    if [[ -n "${SNAC_PID}" ]] && kill -0 "${SNAC_PID}" 2>/dev/null; then
        kill "${SNAC_PID}" 2>/dev/null || true
        wait "${SNAC_PID}" 2>/dev/null || true
    fi

    if [[ -n "${TEST_STREAM_PID}" ]] && kill -0 "${TEST_STREAM_PID}" 2>/dev/null; then
        kill "${TEST_STREAM_PID}" 2>/dev/null || true
        wait "${TEST_STREAM_PID}" 2>/dev/null || true
    fi

    if [[ -n "${TEMP_DIR}" ]] && [[ -d "${TEMP_DIR}" ]]; then
        rm -rf "${TEMP_DIR}"
    fi

    log_info "Cleanup complete."
}

trap cleanup EXIT

setup_temp_dir() {
    TEMP_DIR=$(mktemp -d)
    SNAC_DATA_DIR="${TEMP_DIR}/snac-data"
    OWNCAST_DB="${TEMP_DIR}/owncast.db"

    log_info "Temp directory: ${TEMP_DIR}"
}

check_hosts_entry() {
    # Check if hosts entries exist
    if ! grep -q "${OWNCAST_HOSTNAME}" /etc/hosts 2>/dev/null || ! grep -q "${SNAC_HOSTNAME}" /etc/hosts 2>/dev/null; then
        log_warn "Required /etc/hosts entries not found."
        log_warn "Please add the following to /etc/hosts:"
        echo ""
        echo "    127.0.0.1 ${OWNCAST_HOSTNAME} ${SNAC_HOSTNAME}"
        echo ""
        log_warn "You may need to run: sudo sh -c 'echo \"127.0.0.1 ${OWNCAST_HOSTNAME} ${SNAC_HOSTNAME}\" >> /etc/hosts'"
        exit 1
    fi
    log_info "Hosts entries verified"
}

install_snac2() {
    log_info "Setting up snac2..."

    # Check if snac2 is already installed
    if command -v snac &> /dev/null; then
        SNAC_BIN=$(command -v snac)
        log_info "Using system snac2: ${SNAC_BIN}"
        return
    fi

    # Clone and build snac2
    local snac_src="${TEMP_DIR}/snac2-src"
    log_info "Cloning snac2..."
    git clone --depth 1 https://codeberg.org/grunfink/snac2.git "${snac_src}" 2>/dev/null

    log_info "Building snac2..."
    pushd "${snac_src}" > /dev/null
    make
    SNAC_BIN="${snac_src}/snac"
    popd > /dev/null

    log_info "snac2 built: ${SNAC_BIN}"
}

check_certs() {
    # CERT_DIR may be set by the Docker entrypoint; fall back to local certs/
    CERT_DIR="${CERT_DIR:-${SCRIPT_DIR}/certs}"

    if [[ ! -f "${CERT_DIR}/cert.pem" ]] || [[ ! -f "${CERT_DIR}/key.pem" ]]; then
        log_error "Certificates not found in ${CERT_DIR}"
        log_error "Please run the one-time setup. See README.md for instructions:"
        log_error "  mkcert -install"
        log_error "  mkcert -cert-file certs/cert.pem -key-file certs/key.pem owncast.local snac.local localhost 127.0.0.1"
        exit 1
    fi

    log_info "Using certificates from ${CERT_DIR}"
}

init_snac2() {
    log_info "Initializing snac2..."

    # snac2 advertises URLs via the proxy (HTTPS on PROXY_PORT)
    local snac_host_port="${SNAC_HOSTNAME}:${PROXY_PORT}"

    # Use snac init with piped input: address, port, hostname, prefix, admin email
    printf "127.0.0.1\n%s\n%s\n\ntest@test.local\n" "${SNAC_PORT}" "${snac_host_port}" | \
        "${SNAC_BIN}" init "${SNAC_DATA_DIR}" > /dev/null 2>&1

    if [[ ! -f "${SNAC_DATA_DIR}/server.json" ]]; then
        log_error "snac2 init failed - server.json not created"
        return 1
    fi

    log_info "snac2 initialized"
}

enable_snac_shared_inbox() {
    log_info "Enabling shared_inboxes in snac2 configuration..."

    local server_json="${SNAC_DATA_DIR}/server.json"

    if [[ ! -f "${server_json}" ]]; then
        log_error "server.json not found at ${server_json}"
        return 1
    fi

    # Use jq if available, otherwise use Python as a fallback
    if command -v jq &> /dev/null; then
        local tmp_file
        tmp_file=$(mktemp)
        jq '. + {"shared_inboxes": true}' "${server_json}" > "${tmp_file}" && mv "${tmp_file}" "${server_json}"
    elif command -v python3 &> /dev/null; then
        python3 -c "
import json
with open('${server_json}', 'r') as f:
    data = json.load(f)
data['shared_inboxes'] = True
with open('${server_json}', 'w') as f:
    json.dump(data, f, indent=2)
"
    else
        log_error "Neither jq nor python3 available to modify server.json"
        return 1
    fi

    # Verify the setting was applied
    if grep -q '"shared_inboxes"' "${server_json}"; then
        log_info "shared_inboxes enabled in snac2 configuration"
    else
        log_error "Failed to enable shared_inboxes"
        return 1
    fi
}

create_snac_users() {
    log_info "Creating ${USER_COUNT} users in snac2 using snac adduser..."

    # Generate a unique prefix for this test run to avoid conflicts
    local run_id
    run_id=$(date +%s%N | sha256sum | head -c 8)
    log_info "Test run ID: ${run_id}"

    # Store usernames for later use in follow requests
    SNAC_USERNAMES=()

    local created=0
    local failed=0

    for i in $(seq 1 "${USER_COUNT}"); do
        local username="test${run_id}u${i}"
        local displayname="Test User ${i}"

        # Use snac adduser with piped input (username, display name)
        if printf "%s\n%s\n" "${username}" "${displayname}" | "${SNAC_BIN}" adduser "${SNAC_DATA_DIR}" > /dev/null 2>&1; then
            SNAC_USERNAMES+=("${username}")
            created=$((created + 1))
        else
            failed=$((failed + 1))
        fi

        if [[ $((created + failed)) -gt 0 ]] && [[ $(((created + failed) % 20)) -eq 0 ]]; then
            log_info "Created $((created + failed))/${USER_COUNT} users (successful: ${created})..."
        fi
    done

    log_info "Created ${created} users in snac2 (${failed} failed)"
}

start_proxy() {
    log_info "Starting HTTPS reverse proxy (Caddy)..."

    # Check if caddy is installed
    if ! command -v caddy &> /dev/null; then
        log_error "Caddy is not installed. Please install it: https://caddyserver.com/docs/install"
        return 1
    fi

    # Set environment variables for Caddyfile
    export PROXY_PORT="${PROXY_PORT}"
    export OWNCAST_PORT="${OWNCAST_PORT}"
    export OWNCAST2_PORT="${OWNCAST2_PORT:-8081}"  # Used by Caddyfile; not active in this test
    export SNAC_PORT="${SNAC_PORT}"
    export CERT_FILE="${CERT_DIR}/cert.pem"
    export KEY_FILE="${CERT_DIR}/key.pem"

    # Redirect caddy output to a log file to reduce test noise
    local caddy_log="${TEMP_DIR}/caddy.log"
    caddy run --config "${SCRIPT_DIR}/Caddyfile" --adapter caddyfile > "${caddy_log}" 2>&1 &
    PROXY_PID=$!

    log_info "Caddy started with PID ${PROXY_PID}"
    sleep 2

    # Verify proxy is running
    if ! kill -0 "${PROXY_PID}" 2>/dev/null; then
        log_error "Caddy failed to start"
        return 1
    fi

    # Verify proxy is accepting connections
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

start_snac2() {
    log_info "Starting snac2 server..."

    # mkcert certificates are trusted system-wide, no special env vars needed
    # Redirect snac2 output to a log file to reduce test noise
    local snac_log="${TEMP_DIR}/snac2.log"
    DEBUG=0 "${SNAC_BIN}" httpd "${SNAC_DATA_DIR}" > "${snac_log}" 2>&1 &
    SNAC_PID=$!
    log_info "snac2 logs available at: ${snac_log}"

    log_info "snac2 started with PID ${SNAC_PID}"

    # Wait for snac2 to be ready
    local max_attempts=30
    local attempt=0

    while [[ ${attempt} -lt ${max_attempts} ]]; do
        if curl -s "http://127.0.0.1:${SNAC_PORT}/" > /dev/null 2>&1; then
            log_info "snac2 is ready"
            return 0
        fi
        attempt=$((attempt + 1))
        sleep 1
    done

    log_error "snac2 did not become ready"
    return 1
}

verify_snac_shared_inbox() {
    log_info "Verifying snac2 actors have sharedInbox endpoint..."

    # Get the first test user to check their actor object
    if [[ ${#SNAC_USERNAMES[@]} -eq 0 ]]; then
        log_warn "No snac2 users created yet, skipping sharedInbox verification"
        return 0
    fi

    local test_user="${SNAC_USERNAMES[0]}"
    local actor_url="${SNAC_URL}/${test_user}"

    # Fetch the actor object and check for sharedInbox
    local actor_response
    actor_response=$(curl -s --max-time 10 -H "Accept: application/activity+json" "${actor_url}" 2>&1)

    if echo "${actor_response}" | grep -q '"sharedInbox"'; then
        log_info "sharedInbox endpoint found in snac2 actor objects"
        return 0
    else
        log_error "sharedInbox endpoint NOT found in snac2 actor objects"
        log_error "Actor response: ${actor_response}"
        return 1
    fi
}

build_owncast() {
    log_info "Building Owncast..."

    OWNCAST_BIN="${TEMP_DIR}/owncast"
    pushd "${REPO_ROOT}" > /dev/null
    CGO_ENABLED=1 go build -o "${OWNCAST_BIN}" main.go
    popd > /dev/null

    log_info "Owncast built: ${OWNCAST_BIN}"
}

start_owncast() {
    log_info "Starting Owncast..."

    # Start Owncast with test environment variables and debug flags
    OWNCAST_ALLOW_INTERNAL_FEDERATION=true \
    OWNCAST_INSECURE_SKIP_VERIFY=true \
    "${OWNCAST_BIN}" -database "${OWNCAST_DB}" &
    OWNCAST_PID=$!

    log_info "Owncast started with PID ${OWNCAST_PID}"

    # Wait for Owncast to be ready
    local max_attempts=30
    local attempt=0

    while [[ ${attempt} -lt ${max_attempts} ]]; do
        if curl -s "http://localhost:${OWNCAST_PORT}/api/status" > /dev/null 2>&1; then
            log_info "Owncast is ready"
            return 0
        fi
        attempt=$((attempt + 1))
        sleep 1
    done

    log_error "Owncast did not become ready"
    return 1
}

configure_owncast() {
    log_info "Configuring Owncast with URL: ${OWNCAST_URL}"

    local base_url="http://localhost:${OWNCAST_PORT}"
    local auth
    auth=$(echo -n "${ADMIN_USER}:${ADMIN_PASS}" | base64)

    # Set server URL
    curl -s -X POST "${base_url}/api/admin/config/serverurl" \
        -H "Authorization: Basic ${auth}" \
        -H "Content-Type: application/json" \
        -d "{\"value\": \"${OWNCAST_URL}\"}" > /dev/null

    # Set federation username
    curl -s -X POST "${base_url}/api/admin/config/federation/username" \
        -H "Authorization: Basic ${auth}" \
        -H "Content-Type: application/json" \
        -d "{\"value\": \"${FEDERATION_USERNAME}\"}" > /dev/null

    # Enable federation
    curl -s -X POST "${base_url}/api/admin/config/federation/enable" \
        -H "Authorization: Basic ${auth}" \
        -H "Content-Type: application/json" \
        -d '{"value": true}' > /dev/null

    # Disable private mode (auto-accept follows)
    curl -s -X POST "${base_url}/api/admin/config/federation/private" \
        -H "Authorization: Basic ${auth}" \
        -H "Content-Type: application/json" \
        -d '{"value": false}' > /dev/null

    log_info "Owncast configured"
}

prompt_for_test_stream() {
    # Skip interactive prompts in CI mode
    if [[ "${CI}" == "true" ]]; then
        log_info "CI mode - skipping test stream prompt"
        return
    fi

    echo ""
    echo "----------------------------------------"
    read -p "Would you like to start a test live stream? [y/N] " -n 1 -r
    echo ""

    if [[ ${REPLY} =~ ^[Yy]$ ]]; then
        log_info "Starting test stream..."
        "${REPO_ROOT}/test/ocTestStream.sh" &
        TEST_STREAM_PID=$!

        echo ""
        log_info "Test stream started with PID ${TEST_STREAM_PID}"
        log_info "You can view the stream at: http://localhost:${OWNCAST_PORT}"
        echo ""
        read -p "Press Enter when you are watching the stream and ready to continue..." -r
        echo ""
    fi
}

send_follow_requests() {
    log_info "Sending follow requests from snac2 users..."

    local owncast_actor="${OWNCAST_URL}/federation/user/${FEDERATION_USERNAME}"
    log_info "Following actor: ${owncast_actor}"

    # First verify the actor URL is accessible via curl
    log_info "Verifying actor URL is accessible..."
    local curl_test
    curl_test=$(curl -s --max-time 5 -H "Accept: application/activity+json" "${owncast_actor}" 2>&1)
    if echo "${curl_test}" | grep -q '"type"'; then
        log_info "Actor URL is accessible via curl"
    else
        log_error "Actor URL not accessible: ${curl_test}"
    fi

    local successful=0
    local failed=0
    local total=${#SNAC_USERNAMES[@]}

    for username in "${SNAC_USERNAMES[@]}"; do
        # Use snac follow command and capture output
        # mkcert certificates are trusted system-wide, no special env vars needed
        local follow_output
        follow_output=$("${SNAC_BIN}" follow "${SNAC_DATA_DIR}" "${username}" "${owncast_actor}" 2>&1)
        local follow_exit=$?

        if [[ ${follow_exit} -eq 0 ]] && [[ ! "${follow_output}" =~ "cannot" ]]; then
            successful=$((successful + 1))
        else
            failed=$((failed + 1))
            # Show the first error for debugging
            if [[ ${failed} -eq 1 ]]; then
                log_warn "snac follow error: ${follow_output}"
            fi
        fi

        if [[ $((successful + failed)) -gt 0 ]] && [[ $(((successful + failed) % 20)) -eq 0 ]]; then
            log_info "Follows sent: $((successful + failed))/${total} (successful: ${successful})"
        fi

        # Throttle follow requests to prevent overwhelming snac2
        sleep "${FOLLOW_DELAY}"
    done

    log_test "Follow requests: ${successful} successful, ${failed} failed"

    # Check snac2 queue directory for pending activities
    log_info "Checking snac2 queue..."
    local queue_count
    queue_count=$(find "${SNAC_DATA_DIR}" -name "*.json" -path "*/queue/*" 2>/dev/null | wc -l)
    log_info "snac2 queue has ${queue_count} pending items"

    # Give snac2 background thread time to process all pending follows
    # Each follow requires an HTTP round-trip; scale wait with user count
    local wait_time=$((10 + USER_COUNT / 5))
    log_info "Waiting ${wait_time}s for snac2 to process follow requests..."
    sleep "${wait_time}"
}

verify_followers() {
    log_info "Verifying followers..." >&2

    local auth
    auth=$(echo -n "${ADMIN_USER}:${ADMIN_PASS}" | base64)
    local response

    response=$(curl -s "http://localhost:${OWNCAST_PORT}/api/admin/followers?limit=200" \
        -H "Authorization: Basic ${auth}")

    local count
    count=$(echo "${response}" | grep -o '"total":[0-9]*' | grep -o '[0-9]*' || echo "0")

    log_test "Owncast reports ${count} followers" >&2
    echo "${count}"
}

populate_confirmed_followers() {
    # Query Owncast API for actual registered followers and map their IRIs
    # back to snac2 usernames so we only check those inboxes for delivery.
    CONFIRMED_FOLLOWERS=()
    local auth
    auth=$(echo -n "${ADMIN_USER}:${ADMIN_PASS}" | base64)

    local response
    response=$(curl -s "http://localhost:${OWNCAST_PORT}/api/admin/followers?limit=1000" \
        -H "Authorization: Basic ${auth}" 2>/dev/null || echo '{"results":[]}')

    # Follower IRIs look like https://snac.local:8443/username
    while IFS= read -r iri; do
        if [[ -n "${iri}" ]]; then
            local username="${iri##*/}"
            CONFIRMED_FOLLOWERS+=("${username}")
        fi
    done < <(echo "${response}" | jq -r '.results[]?.link // empty' 2>/dev/null)

    log_info "Confirmed ${#CONFIRMED_FOLLOWERS[@]} follower usernames from Owncast API"
}

send_test_message() {
    log_info "Sending test message from Owncast..."

    local auth
    auth=$(echo -n "${ADMIN_USER}:${ADMIN_PASS}" | base64)
    local message
    message="Test message sent at $(date -u +%Y-%m-%dT%H:%M:%SZ) to ${USER_COUNT} followers"

    curl -s -X POST "http://localhost:${OWNCAST_PORT}/api/admin/federation/send" \
        -H "Authorization: Basic ${auth}" \
        -H "Content-Type: application/json" \
        -d "{\"value\": \"${message}\"}" > /dev/null

    log_info "Message sent: ${message}"
}

check_snac_inboxes_count() {
    # Only check users confirmed as followers by Owncast, not all snac2 users.
    # snac2's shared inbox distributes messages to all users it considers followers,
    # which may include users whose follow hasn't been registered by Owncast yet.
    local usernames_to_check=("${CONFIRMED_FOLLOWERS[@]}")
    if [[ ${#usernames_to_check[@]} -eq 0 ]]; then
        usernames_to_check=("${SNAC_USERNAMES[@]}")
    fi

    local users_with_messages=0
    for username in "${usernames_to_check[@]}"; do
        if user_has_message "${username}"; then
            users_with_messages=$((users_with_messages + 1))
        fi
    done

    echo "${users_with_messages}"
}

user_has_message() {
    local username="$1"
    local user_dir="${SNAC_DATA_DIR}/user/${username}"

    # snac2 stores incoming posts in various places depending on type
    # Check timeline, public, private directories for any .json files
    local count
    for subdir in public private timeline; do
        if [[ -d "${user_dir}/${subdir}" ]]; then
            count=$(find "${user_dir}/${subdir}" -name "*.json" -type f 2>/dev/null | wc -l)
            if [[ "${count}" -gt 0 ]]; then
                return 0
            fi
        fi
    done

    # Also check the global object store for objects addressed to this user
    # snac2 might store activities in the object directory
    return 1
}

check_snac_inboxes() {
    log_info "Checking snac2 user inboxes for delivered messages..." >&2

    local usernames_to_check=("${CONFIRMED_FOLLOWERS[@]}")
    if [[ ${#usernames_to_check[@]} -eq 0 ]]; then
        usernames_to_check=("${SNAC_USERNAMES[@]}")
    fi

    local users_with_messages=0
    local users_without_messages=()
    local total=${#usernames_to_check[@]}

    for username in "${usernames_to_check[@]}"; do
        if user_has_message "${username}"; then
            users_with_messages=$((users_with_messages + 1))
        else
            users_without_messages+=("${username}")
        fi
    done

    log_test "${users_with_messages}/${total} confirmed followers received the message" >&2

    if [[ ${#users_without_messages[@]} -gt 0 ]] && [[ ${#users_without_messages[@]} -le 10 ]]; then
        log_warn "Users missing messages: ${users_without_messages[*]}" >&2
    elif [[ ${#users_without_messages[@]} -gt 10 ]]; then
        log_warn "Users missing messages: ${users_without_messages[*]:0:10} ... and $((${#users_without_messages[@]} - 10)) more" >&2
    fi

    echo "${users_with_messages}"
}

# Count followers with shared_inbox set vs those with individual inboxes only
count_shared_inbox_followers() {
    local shared_count
    shared_count=$(sqlite3 "${OWNCAST_DB}" "SELECT COUNT(*) FROM ap_followers WHERE shared_inbox IS NOT NULL AND shared_inbox != '';" 2>/dev/null || echo "0")
    echo "${shared_count}"
}

count_individual_inbox_followers() {
    local individual_count
    individual_count=$(sqlite3 "${OWNCAST_DB}" "SELECT COUNT(*) FROM ap_followers WHERE shared_inbox IS NULL OR shared_inbox = '';" 2>/dev/null || echo "0")
    echo "${individual_count}"
}

# Clear shared_inbox from a percentage of followers to test mixed delivery
clear_some_shared_inboxes() {
    local percentage=${1:-50}  # Default to 50% of followers

    log_info "Clearing shared_inbox from ${percentage}% of followers to test mixed delivery..."

    # Get total follower count
    local total_followers
    total_followers=$(sqlite3 "${OWNCAST_DB}" "SELECT COUNT(*) FROM ap_followers;" 2>/dev/null || echo "0")

    if [[ "${total_followers}" -eq 0 ]]; then
        log_warn "No followers found in database"
        return 1
    fi

    # Calculate how many to clear
    local clear_count=$((total_followers * percentage / 100))

    log_info "Total followers: ${total_followers}, clearing shared_inbox from ${clear_count} followers"

    # Clear shared_inbox from a random subset of followers
    # SQLite's RANDOM() function helps select random rows
    sqlite3 "${OWNCAST_DB}" "UPDATE ap_followers SET shared_inbox = NULL WHERE iri IN (SELECT iri FROM ap_followers ORDER BY RANDOM() LIMIT ${clear_count});"

    local shared_remaining
    shared_remaining=$(count_shared_inbox_followers)
    local individual_count
    individual_count=$(count_individual_inbox_followers)

    log_info "After clearing: ${shared_remaining} followers with shared_inbox, ${individual_count} with individual inbox only"
}

verify_all_followers_received_message() {
    local followers=$1
    local max_wait=${2:-60}

    log_info "Verifying all ${followers} followers receive the message (max ${max_wait}s)..." >&2

    local start_time
    start_time=$(date +%s)
    local last_count=0
    local delivery_time=0
    local current_time
    local count

    while true; do
        current_time=$(date +%s)
        local elapsed=$((current_time - start_time))

        if [[ ${elapsed} -ge ${max_wait} ]]; then
            log_warn "Timeout after ${max_wait}s" >&2
            delivery_time=${elapsed}
            break
        fi

        count=$(check_snac_inboxes_count)

        if [[ ${count} -ne ${last_count} ]]; then
            local pct=$((count * 100 / followers))
            log_info "Delivery progress: ${count}/${followers} (${pct}%)" >&2
            last_count=${count}
        fi

        if [[ ${count} -ge ${followers} ]]; then
            delivery_time=${elapsed}
            log_test "All ${followers} followers received the message in ${delivery_time}s" >&2
            echo "${delivery_time}"
            return 0
        fi

        sleep 2
    done

    local final_count
    final_count=$(check_snac_inboxes_count)
    echo "${delivery_time}"
    if [[ ${final_count} -lt ${followers} ]]; then
        log_error "Only ${final_count}/${followers} followers received the message" >&2
        return 1
    fi

    return 0
}

print_results() {
    local followers=$1
    local delivered=$2
    local delivery_time=$3

    # Get shared inbox vs individual inbox counts
    local shared_inbox_count
    shared_inbox_count=$(count_shared_inbox_followers)
    local individual_inbox_count
    individual_inbox_count=$(count_individual_inbox_followers)

    echo ""
    echo "========================================"
    echo "ActivityPub Federation Test Results"
    echo "========================================"
    echo "Test Users Created:   ${USER_COUNT}"
    echo "Followers Registered: ${followers}"
    echo "  - Shared Inbox:     ${shared_inbox_count}"
    echo "  - Individual Inbox: ${individual_inbox_count}"
    echo "Messages Delivered:   ${delivered}"
    if [[ -n "${delivery_time}" ]] && [[ "${delivery_time}" -gt 0 ]]; then
        echo "Delivery Time:        ${delivery_time}s"
    fi
    echo ""

    if [[ "${followers}" -gt 0 ]]; then
        local follow_rate=$((followers * 100 / USER_COUNT))
        echo "Follow Success Rate:  ${follow_rate}%"
    fi

    if [[ "${followers}" -gt 0 ]]; then
        local delivery_rate=$((delivered * 100 / followers))
        echo "Delivery Rate:        ${delivery_rate}%"
    fi

    echo ""
    echo "----------------------------------------"
    if [[ "${delivered}" -eq "${followers}" ]] && [[ "${followers}" -gt 0 ]]; then
        echo -e "${GREEN}TEST PASSED${NC}"
        echo "All ${followers} followers received the message."
    elif [[ "${delivered}" -gt 0 ]]; then
        echo -e "${RED}TEST FAILED${NC}"
        echo "Only ${delivered} of ${followers} followers received the message."
        echo "Missing: $((followers - delivered)) followers"
    elif [[ "${followers}" -eq 0 ]]; then
        echo -e "${RED}TEST FAILED${NC}"
        echo "No followers were registered."
    else
        echo -e "${RED}TEST FAILED${NC}"
        echo "No messages were delivered to any followers."
    fi
    echo "========================================"
}

main() {
    # Kill any leftover processes from previous runs
    kill_leftover_processes

    echo ""
    echo "========================================"
    echo "ActivityPub Federation Test"
    echo "========================================"
    log_info "Configuration: ${USER_COUNT} test users"
    log_info "Owncast URL: ${OWNCAST_URL}"
    log_info "snac2 URL: ${SNAC_URL}"
    if [[ "${CLEAR_SHARED_INBOX_PERCENT}" -gt 0 ]]; then
        log_info "Mixed delivery test: ${CLEAR_SHARED_INBOX_PERCENT}% of followers will use individual inbox"
    fi
    echo ""

    # ==========================================
    # STEP 1: Setup snac2 server with test users
    # ==========================================
    echo "----------------------------------------"
    echo "STEP 1: Setup snac2 server with test users"
    echo "----------------------------------------"
    setup_temp_dir
    check_hosts_entry
    install_snac2
    check_certs
    init_snac2
    enable_snac_shared_inbox
    create_snac_users
    start_proxy
    start_snac2
    verify_snac_shared_inbox
    echo ""

    # ==========================================
    # STEP 2: Setup Owncast with federation
    # ==========================================
    echo "----------------------------------------"
    echo "STEP 2: Setup Owncast with federation"
    echo "----------------------------------------"
    build_owncast
    start_owncast
    configure_owncast
    sleep 2

    # Optional: prompt for test stream
    prompt_for_test_stream
    echo ""

    # ==========================================
    # STEP 3: snac2 users follow Owncast
    # ==========================================
    echo "----------------------------------------"
    echo "STEP 3: snac2 users follow Owncast"
    echo "----------------------------------------"
    send_follow_requests

    # Wait for all followers to be registered (with timeout)
    local followers=0
    local max_wait=$((60 + USER_COUNT))  # Scale timeout with user count
    local waited=0
    local check_interval=2

    log_info "Waiting for ${USER_COUNT} followers to be registered..."
    while [[ "${followers}" -lt "${USER_COUNT}" ]] && [[ "${waited}" -lt "${max_wait}" ]]; do
        sleep "${check_interval}"
        waited=$((waited + check_interval))
        followers=$(verify_followers)
        if [[ "${followers}" -lt "${USER_COUNT}" ]]; then
            log_info "Followers registered: ${followers}/${USER_COUNT} (waited ${waited}s)"
        fi
    done

    if [[ "${followers}" -eq 0 ]]; then
        log_error "No followers registered - cannot proceed with message delivery test"
        print_results "${followers}" 0
        exit 1
    fi

    if [[ "${followers}" -lt "${USER_COUNT}" ]]; then
        log_warn "Only ${followers}/${USER_COUNT} followers registered after ${max_wait}s timeout"
    else
        log_info "All ${followers} followers registered"
    fi

    # Optionally clear shared_inbox from some followers to test mixed delivery
    if [[ "${CLEAR_SHARED_INBOX_PERCENT}" -gt 0 ]]; then
        echo ""
        echo "----------------------------------------"
        echo "STEP 3.5: Clear shared_inbox from ${CLEAR_SHARED_INBOX_PERCENT}% of followers"
        echo "----------------------------------------"
        clear_some_shared_inboxes "${CLEAR_SHARED_INBOX_PERCENT}"
    fi
    echo ""

    # ==========================================
    # STEP 4: Owncast sends message to followers
    # ==========================================
    echo "----------------------------------------"
    echo "STEP 4: Owncast sends message to followers"
    echo "----------------------------------------"
    populate_confirmed_followers
    send_test_message
    echo ""

    # ==========================================
    # STEP 5: Verify all followers received message
    # ==========================================
    echo "----------------------------------------"
    echo "STEP 5: Verify all followers received message"
    echo "----------------------------------------"
    local delivery_success=0
    local delivery_time
    if delivery_time=$(verify_all_followers_received_message "${followers}" 60); then
        delivery_success=1
    fi

    local delivered
    delivered=$(check_snac_inboxes)

    print_results "${followers}" "${delivered}" "${delivery_time}"

    if [[ "${KEEP_RUNNING:-}" == "true" ]]; then
        log_info "Keeping servers running (Ctrl+C to stop)..."
        log_info "  - Owncast: http://localhost:${OWNCAST_PORT} (${OWNCAST_URL})"
        log_info "  - snac2: http://localhost:${SNAC_PORT} (${SNAC_URL})"
        wait
    fi

    # Exit with appropriate code
    # All conditions must be met:
    # 1. All follow requests succeeded (followers == USER_COUNT)
    # 2. Message delivery completed successfully
    # 3. All followers received the message (delivered == followers)
    if [[ "${followers}" -ne "${USER_COUNT}" ]]; then
        log_error "Follow count mismatch: expected ${USER_COUNT}, got ${followers}"
        exit 1
    fi

    if [[ "${delivery_success}" -ne 1 ]]; then
        log_error "Message delivery failed"
        exit 1
    fi

    if [[ "${delivered}" -ne "${followers}" ]]; then
        log_error "Delivery count mismatch: expected ${followers}, got ${delivered}"
        exit 1
    fi

    log_info "All tests passed"
    exit 0
}

main "$@"
