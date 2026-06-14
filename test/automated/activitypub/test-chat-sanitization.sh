#!/bin/bash
# shellcheck disable=SC2317,SC2329  # cleanup() is invoked via trap, not direct call
# shellcheck disable=SC2034  # SNAC_URL is unused but kept for consistency with other AP tests

# Chat Sanitization Test for Fediverse Engagement Events
#
# This test verifies that malicious HTML and markdown in ActivityPub actor
# display names is sanitized before being rendered in chat messages.
#
# The test:
# 1. Starts snac2 + Owncast + HTTPS proxy (same as federation test)
# 2. Creates snac2 users with malicious display names (HTML, markdown)
# 3. Enables engagement display in chat
# 4. Has the malicious users follow Owncast
# 5. Queries the chat messages from the database
# 6. Verifies no HTML tags or markdown artifacts appear in stored messages
#
# Prerequisites:
#   Add to /etc/hosts: 127.0.0.1 owncast.local snac.local

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(git rev-parse --show-toplevel)"

# Configuration
FOLLOW_DELAY="${FOLLOW_DELAY:-0.3}"
CI="${CI:-false}"
PROXY_PORT="${PROXY_PORT:-8443}"
SNAC_PORT="${SNAC_PORT:-9080}"
SNAC_HOSTNAME="snac.local"
OWNCAST_PORT="${OWNCAST_PORT:-8080}"
OWNCAST_HOSTNAME="owncast.local"
ADMIN_USER="admin"
ADMIN_PASS="abc123"
FEDERATION_USERNAME="streamer"

# URLs
SNAC_URL="https://${SNAC_HOSTNAME}:${PROXY_PORT}"
OWNCAST_URL="https://${OWNCAST_HOSTNAME}:${PROXY_PORT}"

# Directories
TEMP_DIR=""
SNAC_DATA_DIR=""
SNAC_BIN=""
OWNCAST_DB=""

# PIDs
SNAC_PID=""
OWNCAST_PID=""
PROXY_PID=""

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

# Test users: username -> malicious display name
# These simulate real-world attack payloads from the security advisory
declare -A MALICIOUS_USERS
MALICIOUS_USERS=(
    ["htmlscript"]='<script>alert("xss")</script>'
    ["htmliframe"]='<iframe src="https://evil.com"></iframe>'
    ["htmlimg"]='<img src="https://evil.com/track.png">'
    ["htmlform"]='<form action="https://evil.com"><input type="password"></form>'
    ["htmlmeta"]='<meta http-equiv="refresh" content="0;url=https://evil.com">'
    ["htmlstyle"]='<style>body{display:none}</style>Visible'
    ["mdimage"]='![xss](https://evil.com/track.png)'
    ["mdlink"]='[Click me](https://evil.com)'
    ["mixedhtml"]='Alice <script>alert(1)</script> Bob'
    ["cleanuser"]='Legitimate User'
)

# Patterns that must NOT appear in stored chat message bodies
FORBIDDEN_PATTERNS=(
    '<script'
    '<iframe'
    '<img'
    '<form'
    '<meta'
    '<style'
    '<input'
    'src="https://evil'
    'action="https://evil'
    'onerror='
    'onload='
)

kill_leftover_processes() {
    local killed=false
    if pkill -f "snac httpd /tmp" 2>/dev/null; then killed=true; fi
    if pkill -f "local-proxy.js" 2>/dev/null; then killed=true; fi

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
    if ! grep -q "${OWNCAST_HOSTNAME}" /etc/hosts 2>/dev/null || ! grep -q "${SNAC_HOSTNAME}" /etc/hosts 2>/dev/null; then
        log_warn "Required /etc/hosts entries not found."
        log_warn "Please add: 127.0.0.1 ${OWNCAST_HOSTNAME} ${SNAC_HOSTNAME}"
        exit 1
    fi
    log_info "Hosts entries verified"
}

install_snac2() {
    log_info "Setting up snac2..."
    if command -v snac &> /dev/null; then
        SNAC_BIN=$(command -v snac)
        log_info "Using system snac2: ${SNAC_BIN}"
        return
    fi

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
        log_error "Certificates not found in ${CERT_DIR}. See README.md for setup."
        exit 1
    fi
    log_info "Using certificates from ${CERT_DIR}"
}

init_snac2() {
    log_info "Initializing snac2..."
    local snac_host_port="${SNAC_HOSTNAME}:${PROXY_PORT}"
    printf "127.0.0.1\n%s\n%s\n\ntest@test.local\n" "${SNAC_PORT}" "${snac_host_port}" | \
        "${SNAC_BIN}" init "${SNAC_DATA_DIR}" > /dev/null 2>&1
    log_info "snac2 initialized"
}

create_malicious_users() {
    log_info "Creating snac2 users with malicious display names..."

    local run_id
    run_id=$(date +%s%N | sha256sum | head -c 8)

    local created=0
    for username in "${!MALICIOUS_USERS[@]}"; do
        local displayname="${MALICIOUS_USERS[$username]}"
        local full_username="${username}${run_id}"

        if printf "%s\n%s\n" "${full_username}" "${displayname}" | "${SNAC_BIN}" adduser "${SNAC_DATA_DIR}" > /dev/null 2>&1; then
            # Update the key to include the run_id so we can find them later
            MALICIOUS_USERS["${username}"]="${displayname}"
            # Store the full username for follow requests
            SNAC_FULL_USERNAMES["${username}"]="${full_username}"
            created=$((created + 1))
            log_info "  Created user '${full_username}' with display name: ${displayname}"
        else
            log_error "  Failed to create user '${full_username}'"
        fi
    done

    log_info "Created ${created} users with malicious display names"
}

start_proxy() {
    log_info "Starting HTTPS reverse proxy (Caddy)..."
    if ! command -v caddy &> /dev/null; then
        log_error "Caddy is not installed."
        return 1
    fi

    export PROXY_PORT OWNCAST_PORT SNAC_PORT
    export CERT_FILE="${CERT_DIR}/cert.pem"
    export KEY_FILE="${CERT_DIR}/key.pem"

    local caddy_log="${TEMP_DIR}/caddy.log"
    caddy run --config "${SCRIPT_DIR}/Caddyfile" --adapter caddyfile > "${caddy_log}" 2>&1 &
    PROXY_PID=$!
    sleep 2

    if ! kill -0 "${PROXY_PID}" 2>/dev/null; then
        log_error "Caddy failed to start"
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

start_snac2() {
    log_info "Starting snac2 server..."
    local snac_log="${TEMP_DIR}/snac2.log"
    DEBUG=0 "${SNAC_BIN}" httpd "${SNAC_DATA_DIR}" > "${snac_log}" 2>&1 &
    SNAC_PID=$!

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

build_owncast() {
    log_info "Building Owncast..."
    pushd "${REPO_ROOT}" > /dev/null
    CGO_ENABLED=1 go build -o owncast main.go
    popd > /dev/null
    log_info "Owncast built"
}

start_owncast() {
    log_info "Starting Owncast..."
    OWNCAST_ALLOW_INTERNAL_FEDERATION=true \
    OWNCAST_INSECURE_SKIP_VERIFY=true \
    "${REPO_ROOT}/owncast" -database "${OWNCAST_DB}" &
    OWNCAST_PID=$!

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
    log_info "Configuring Owncast..."
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

    # Disable private mode
    curl -s -X POST "${base_url}/api/admin/config/federation/private" \
        -H "Authorization: Basic ${auth}" \
        -H "Content-Type: application/json" \
        -d '{"value": false}' > /dev/null

    # Enable engagement display in chat
    curl -s -X POST "${base_url}/api/admin/config/federation/showengagement" \
        -H "Authorization: Basic ${auth}" \
        -H "Content-Type: application/json" \
        -d '{"value": true}' > /dev/null

    log_info "Owncast configured (engagement display enabled)"
}

send_follow_requests() {
    log_info "Sending follow requests from malicious users..."
    local owncast_actor="${OWNCAST_URL}/federation/user/${FEDERATION_USERNAME}"
    local successful=0

    for username in "${!MALICIOUS_USERS[@]}"; do
        local full_username="${SNAC_FULL_USERNAMES[$username]}"
        local follow_output
        follow_output=$("${SNAC_BIN}" follow "${SNAC_DATA_DIR}" "${full_username}" "${owncast_actor}" 2>&1)
        local follow_exit=$?

        if [[ ${follow_exit} -eq 0 ]] && [[ ! "${follow_output}" =~ "cannot" ]]; then
            successful=$((successful + 1))
            log_info "  ${full_username} followed Owncast"
        else
            log_warn "  ${full_username} follow failed: ${follow_output}"
        fi

        sleep "${FOLLOW_DELAY}"
    done

    log_test "${successful}/${#MALICIOUS_USERS[@]} follow requests sent"

    # Wait for follow requests to be processed
    local user_count=${#MALICIOUS_USERS[@]}
    local wait_time=$((15 + user_count))
    log_info "Waiting ${wait_time}s for engagement events to be processed..."
    sleep "${wait_time}"
}

verify_chat_sanitization() {
    log_info "Verifying chat message sanitization..."

    local passed=true
    local tests_run=0
    local tests_passed=0

    # Query fediverse engagement messages from the database
    local messages
    messages=$(sqlite3 "${OWNCAST_DB}" \
        "SELECT body FROM messages WHERE eventType IN ('FEDIVERSE_ENGAGEMENT_FOLLOW', 'FEDIVERSE_ENGAGEMENT_LIKE', 'FEDIVERSE_ENGAGEMENT_REPOST');" 2>/dev/null)

    if [[ -z "${messages}" ]]; then
        log_error "No fediverse engagement messages found in database"
        return 1
    fi

    local message_count
    message_count=$(echo "${messages}" | wc -l | tr -d ' ')
    log_info "Found ${message_count} engagement messages in chat"

    # Test 1: No forbidden HTML patterns in any message
    for pattern in "${FORBIDDEN_PATTERNS[@]}"; do
        tests_run=$((tests_run + 1))
        if echo "${messages}" | grep -qi "${pattern}"; then
            log_error "FAIL: Found forbidden pattern '${pattern}' in chat messages:"
            echo "${messages}" | grep -i "${pattern}" | while read -r line; do
                log_error "  Body: ${line}"
            done
            passed=false
        else
            tests_passed=$((tests_passed + 1))
            log_test "PASS: No '${pattern}' found in messages"
        fi
    done

    # Test 2: Verify the clean user's message is present
    # Note: snac2 may not serve the display name in actor objects, so we check
    # for either the display name or the username fallback.
    tests_run=$((tests_run + 1))
    local clean_username="${SNAC_FULL_USERNAMES[cleanuser]}"
    if echo "${messages}" | grep -q "Legitimate User"; then
        tests_passed=$((tests_passed + 1))
        log_test "PASS: Clean display name 'Legitimate User' preserved correctly"
    elif echo "${messages}" | grep -q "${clean_username}"; then
        tests_passed=$((tests_passed + 1))
        log_test "PASS: Clean user present via username fallback '${clean_username}'"
    else
        log_error "FAIL: Clean user message not found in chat"
        passed=false
    fi

    # Test 3: Verify that messages with stripped HTML fell back to expected content
    # Messages from users whose display names were entirely HTML should show
    # the follow action text but not the HTML
    tests_run=$((tests_run + 1))
    if echo "${messages}" | grep -q "followed this stream"; then
        tests_passed=$((tests_passed + 1))
        log_test "PASS: Follow action text present in messages"
    else
        log_error "FAIL: No 'followed this stream' text found in any message"
        passed=false
    fi

    # Test 4: No markdown image syntax rendered as HTML img tags
    tests_run=$((tests_run + 1))
    if echo "${messages}" | grep -qi '<img.*src=.*evil'; then
        log_error "FAIL: Markdown image syntax was rendered as HTML img tag"
        passed=false
    else
        tests_passed=$((tests_passed + 1))
        log_test "PASS: No markdown-rendered img tags found"
    fi

    # Test 5: No markdown link syntax rendered as HTML anchor tags with evil URLs
    tests_run=$((tests_run + 1))
    if echo "${messages}" | grep -qi '<a.*href=.*evil'; then
        log_error "FAIL: Markdown link syntax was rendered as HTML anchor tag"
        passed=false
    else
        tests_passed=$((tests_passed + 1))
        log_test "PASS: No markdown-rendered anchor tags with evil URLs found"
    fi

    # Print all stored messages for inspection
    echo ""
    log_info "All stored engagement messages:"
    echo "${messages}" | while read -r line; do
        log_info "  ${line}"
    done

    echo ""
    echo "========================================"
    echo "Chat Sanitization Test Results"
    echo "========================================"
    echo "Tests run:    ${tests_run}"
    echo "Tests passed: ${tests_passed}"
    echo "========================================"

    if [[ "${passed}" == "true" ]]; then
        echo -e "${GREEN}TEST PASSED${NC}"
        return 0
    else
        echo -e "${RED}TEST FAILED${NC}"
        return 1
    fi
}

# We need an associative array for the full usernames
declare -A SNAC_FULL_USERNAMES

main() {
    kill_leftover_processes

    echo ""
    echo "========================================"
    echo "Chat Sanitization Test"
    echo "========================================"
    echo ""

    # Setup infrastructure
    setup_temp_dir
    check_hosts_entry
    install_snac2
    check_certs
    init_snac2
    create_malicious_users
    start_proxy
    start_snac2
    build_owncast
    start_owncast
    configure_owncast
    sleep 2

    # Send follows from malicious users
    send_follow_requests

    # Verify results
    if verify_chat_sanitization; then
        exit 0
    else
        exit 1
    fi
}

main "$@"
