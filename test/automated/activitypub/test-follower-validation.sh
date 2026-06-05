#!/bin/bash
# shellcheck disable=SC2317  # cleanup() is invoked via trap, not direct call

# Follower Validation Test
#
# This test verifies that the follower validation job correctly:
# 1. Removes invalid/fake followers that have been failing for > 7 days
# 2. Keeps valid real followers
#
# The test:
# 1. Builds and starts Owncast with a 20-second validation interval
# 2. Configures federation via admin API
# 3. Inserts 5 fake followers and 5 real followers directly into the database
# 4. Sets first_validation_failure_at to 8+ days ago for fake followers
# 5. Waits for the validation job to run
# 6. Verifies fake followers were removed and real followers remain

set -e

_SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(git rev-parse --show-toplevel)"

# Configuration
OWNCAST_PORT="${OWNCAST_PORT:-8080}"
ADMIN_USER="admin"
ADMIN_PASS="abc123"
FEDERATION_USERNAME="streamer"
VALIDATION_INTERVAL=20  # seconds
MAX_TEST_DURATION=600   # Maximum test duration in seconds (10 minutes)
TEST_START_TIME=""

# Directories
TEMP_DIR=""
OWNCAST_BIN=""
OWNCAST_DB=""

# PIDs
OWNCAST_PID=""

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

# Check if test has exceeded maximum duration
check_timeout() {
    if [[ -z "${TEST_START_TIME}" ]]; then
        return 0
    fi
    local now
    now=$(date +%s)
    local elapsed=$((now - TEST_START_TIME))
    if [[ ${elapsed} -ge ${MAX_TEST_DURATION} ]]; then
        log_error "Test exceeded maximum duration of ${MAX_TEST_DURATION} seconds"
        exit 1
    fi
}

# Real ActivityPub accounts (these should resolve successfully)
REAL_ACCOUNTS=(
    # Original 5
    "https://social.gabekangas.com/users/gabek"
    "https://social.owncast.online/users/owncast"
    "https://sfba.social/users/dnalounge"
    "https://mastodon.social/users/docpop"
    "https://mastodon.social/users/Gargron"
    # Additional 20 real accounts from various instances
    "https://mastodon.social/users/Mastodon"
    "https://mastodon.social/users/joinmastodon"
    "https://mastodon.social/users/thunderbird"
    "https://mastodon.social/users/mozilla"
    "https://mastodon.social/users/pixelfed"
    "https://fosstodon.org/users/fosstodon"
    "https://fosstodon.org/users/kde"
    "https://fosstodon.org/users/gnome"
    "https://hachyderm.io/users/hachyderm"
    "https://hachyderm.io/users/nova"
    "https://infosec.exchange/users/jerry"
    "https://mas.to/users/Mer__edith"
    "https://mastodon.online/users/mastodonusercount"
    "https://mstdn.social/users/stux"
    "https://techhub.social/users/techhub"
    "https://universeodon.com/users/Universeodon"
    "https://mastodon.world/users/mastodonworld"
    "https://social.vivaldi.net/users/user_vivaldi"
    "https://c.im/users/publicvoit"
    "https://mastodon.social/users/w3c"
)

# Fake follower IRIs (these will fail validation)
FAKE_ACCOUNTS=(
    # Original 5
    "https://fake-server-1.invalid/users/fakeuser1"
    "https://fake-server-2.invalid/users/fakeuser2"
    "https://fake-server-3.invalid/users/fakeuser3"
    "https://nonexistent-instance.fake/users/nobody"
    "https://totally-fake.example/actors/testuser"
    # Additional 20 fake accounts
    "https://fake-server-4.invalid/users/fakeuser4"
    "https://fake-server-5.invalid/users/fakeuser5"
    "https://fake-server-6.invalid/users/fakeuser6"
    "https://fake-server-7.invalid/users/fakeuser7"
    "https://fake-server-8.invalid/users/fakeuser8"
    "https://fake-server-9.invalid/users/fakeuser9"
    "https://fake-server-10.invalid/users/fakeuser10"
    "https://nonexistent-1.fake/users/ghost1"
    "https://nonexistent-2.fake/users/ghost2"
    "https://nonexistent-3.fake/users/ghost3"
    "https://nonexistent-4.fake/users/ghost4"
    "https://nonexistent-5.fake/users/ghost5"
    "https://imaginary-instance.test/actors/phantom1"
    "https://imaginary-instance.test/actors/phantom2"
    "https://imaginary-instance.test/actors/phantom3"
    "https://madeup-server.invalid/users/invented1"
    "https://madeup-server.invalid/users/invented2"
    "https://madeup-server.invalid/users/invented3"
    "https://fictional-fedi.fake/users/notreal1"
    "https://fictional-fedi.fake/users/notreal2"
)

cleanup() {
    log_info "Cleaning up..."

    if [[ -n "${OWNCAST_PID}" ]] && kill -0 "${OWNCAST_PID}" 2>/dev/null; then
        kill "${OWNCAST_PID}" 2>/dev/null || true
        wait "${OWNCAST_PID}" 2>/dev/null || true
    fi

    if [[ -n "${TEMP_DIR}" ]] && [[ -d "${TEMP_DIR}" ]]; then
        rm -rf "${TEMP_DIR}"
    fi

    log_info "Cleanup complete."
}

trap cleanup EXIT

setup_temp_dir() {
    TEMP_DIR=$(mktemp -d)
    OWNCAST_DB="${TEMP_DIR}/owncast.db"
    log_info "Temp directory: ${TEMP_DIR}"
    log_info "Database: ${OWNCAST_DB}"
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
    log_info "Starting Owncast with ${VALIDATION_INTERVAL}s validation interval..."

    # Start Owncast with test configuration
    OWNCAST_ALLOW_INTERNAL_FEDERATION=true \
    OWNCAST_INSECURE_SKIP_VERIFY=true \
    "${OWNCAST_BIN}" \
        -database "${OWNCAST_DB}" \
        -followervalidationinterval "${VALIDATION_INTERVAL}" \
        -enableVerboseLogging &
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

stop_owncast() {
    log_info "Stopping Owncast..."
    if [[ -n "${OWNCAST_PID}" ]] && kill -0 "${OWNCAST_PID}" 2>/dev/null; then
        kill "${OWNCAST_PID}" 2>/dev/null || true
        wait "${OWNCAST_PID}" 2>/dev/null || true
        OWNCAST_PID=""
    fi
    # Give the database time to fully close
    sleep 2
    log_info "Owncast stopped"
}

configure_owncast() {
    log_info "Configuring Owncast federation..."

    local base_url="http://localhost:${OWNCAST_PORT}"
    local auth
    auth=$(echo -n "${ADMIN_USER}:${ADMIN_PASS}" | base64)

    # Set server URL
    curl -s -X POST "${base_url}/api/admin/config/serverurl" \
        -H "Authorization: Basic ${auth}" \
        -H "Content-Type: application/json" \
        -d '{"value": "http://localhost:8080"}' > /dev/null

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

    log_info "Owncast federation configured"
}

# Get the domain from an actor IRI
get_domain_from_iri() {
    local iri="$1"
    echo "${iri}" | sed -E 's|https?://([^/]+)/.*|\1|'
}

# Get the username from an actor IRI
get_username_from_iri() {
    local iri="$1"
    echo "${iri}" | sed -E 's|.*/([^/]+)$|\1|'
}

insert_followers() {
    log_info "Inserting test followers into database..."

    local now
    now=$(date -u +"%Y-%m-%d %H:%M:%S")

    # 8 days ago - past the 7-day threshold for removal
    local eight_days_ago
    eight_days_ago=$(date -u -d "8 days ago" +"%Y-%m-%d %H:%M:%S" 2>/dev/null || date -u -v-8d +"%Y-%m-%d %H:%M:%S")

    # Insert real followers (these should validate successfully and remain)
    log_info "Inserting ${#REAL_ACCOUNTS[@]} real followers..."
    for iri in "${REAL_ACCOUNTS[@]}"; do
        local domain
        domain=$(get_domain_from_iri "${iri}")
        local username
        username=$(get_username_from_iri "${iri}")
        local inbox="https://${domain}/inbox"
        local shared_inbox="https://${domain}/inbox"
        local request_iri="https://localhost:8080/federation/follow-request-${username}"

        sqlite3 "${OWNCAST_DB}" <<EOF
INSERT INTO ap_followers (iri, inbox, shared_inbox, name, username, image, request, request_object, created_at, approved_at)
VALUES (
    '${iri}',
    '${inbox}',
    '${shared_inbox}',
    '${username}',
    '${username}@${domain}',
    '',
    '${request_iri}',
    '{}',
    '${now}',
    '${now}'
);
EOF
        log_info "  Inserted real follower: ${iri}"
    done

    # Insert fake followers with first_validation_failure_at set to 8 days ago
    # These should be removed after validation fails
    log_info "Inserting ${#FAKE_ACCOUNTS[@]} fake followers (with 8-day-old failure timestamp)..."
    for iri in "${FAKE_ACCOUNTS[@]}"; do
        local domain
        domain=$(get_domain_from_iri "${iri}")
        local username
        username=$(get_username_from_iri "${iri}")
        local inbox="https://${domain}/inbox"
        local request_iri="https://localhost:8080/federation/follow-request-${username}"

        sqlite3 "${OWNCAST_DB}" <<EOF
INSERT INTO ap_followers (iri, inbox, name, username, image, request, request_object, created_at, approved_at, first_validation_failure_at)
VALUES (
    '${iri}',
    '${inbox}',
    '${username}',
    '${username}@${domain}',
    '',
    '${request_iri}',
    '{}',
    '${now}',
    '${now}',
    '${eight_days_ago}'
);
EOF
        log_info "  Inserted fake follower: ${iri}"
    done

    log_info "Inserted ${#REAL_ACCOUNTS[@]} real and ${#FAKE_ACCOUNTS[@]} fake followers"
}

# Get all follower IRIs via the admin API
# Returns a newline-separated list of follower IRIs
get_followers_via_api() {
    local base_url="http://localhost:${OWNCAST_PORT}"
    local auth
    auth=$(echo -n "${ADMIN_USER}:${ADMIN_PASS}" | base64)

    # Fetch all followers (using high limit to get all at once)
    local response
    response=$(curl -s "${base_url}/api/admin/followers?limit=1000" \
        -H "Authorization: Basic ${auth}" 2>/dev/null || echo '{"results":[]}')

    # Extract IRIs from the JSON response
    echo "${response}" | jq -r '.results[]?.link // empty' 2>/dev/null || echo ""
}

get_follower_count() {
    local followers
    followers=$(get_followers_via_api)
    if [[ -z "${followers}" ]]; then
        echo "0"
    else
        echo "${followers}" | wc -l | tr -d ' '
    fi
}

get_real_follower_count() {
    # Count followers that match real account patterns
    local followers
    followers=$(get_followers_via_api)
    local count=0
    for iri in "${REAL_ACCOUNTS[@]}"; do
        if echo "${followers}" | grep -qF "${iri}"; then
            count=$((count + 1))
        fi
    done
    echo "${count}"
}

get_fake_follower_count() {
    # Count followers that match fake account patterns
    local followers
    followers=$(get_followers_via_api)
    local count=0
    for iri in "${FAKE_ACCOUNTS[@]}"; do
        if echo "${followers}" | grep -qF "${iri}"; then
            count=$((count + 1))
        fi
    done
    echo "${count}"
}

wait_for_validation() {
    log_info "Waiting for follower validation job to run..."

    local _initial_fake_count
    _initial_fake_count=$(get_fake_follower_count)

    # Calculate wait time based on number of followers
    # The job validates 5 followers per run (FollowersPerRun)
    # With 2-second delay between followers and VALIDATION_INTERVAL between runs
    local total_followers=$((${#REAL_ACCOUNTS[@]} + ${#FAKE_ACCOUNTS[@]}))
    local cycles_needed=$(( (total_followers + 4) / 5 ))  # Round up
    local wait_time=$(( (VALIDATION_INTERVAL + 12) * cycles_needed + 30 ))  # Extra buffer for network delays

    log_info "Total followers: ${total_followers}, estimated cycles needed: ${cycles_needed}"
    log_info "Waiting up to ${wait_time} seconds for validation cycles..."

    local elapsed=0
    local check_interval=10

    while [[ ${elapsed} -lt ${wait_time} ]]; do
        # Check if we've exceeded the maximum test duration
        check_timeout

        sleep ${check_interval}
        elapsed=$((elapsed + check_interval))

        local current_fake_count
        current_fake_count=$(get_fake_follower_count)
        local current_real_count
        current_real_count=$(get_real_follower_count)

        log_info "  ${elapsed}s: Real followers: ${current_real_count}/${#REAL_ACCOUNTS[@]}, Fake followers: ${current_fake_count}/${#FAKE_ACCOUNTS[@]}"

        # If all fake followers are gone and all real remain, we can stop early
        if [[ ${current_fake_count} -eq 0 ]] && [[ ${current_real_count} -eq ${#REAL_ACCOUNTS[@]} ]]; then
            log_info "Validation complete early - all fake followers removed"
            return 0
        fi
    done

    return 0
}

verify_results() {
    log_info "Verifying validation results..."

    local total_count
    total_count=$(get_follower_count)
    local real_count
    real_count=$(get_real_follower_count)
    local fake_count
    fake_count=$(get_fake_follower_count)

    log_test "Total followers: ${total_count}"
    log_test "Real followers remaining: ${real_count}/${#REAL_ACCOUNTS[@]}"
    log_test "Fake followers remaining: ${fake_count}/${#FAKE_ACCOUNTS[@]}"

    # List remaining followers for debugging
    log_info "Remaining followers:"
    get_followers_via_api | while read -r iri; do
        if [[ -n "${iri}" ]]; then
            log_info "  - ${iri}"
        fi
    done

    echo ""
    echo "========================================"
    echo "Follower Validation Test Results"
    echo "========================================"
    echo "Expected real followers:  ${#REAL_ACCOUNTS[@]}"
    echo "Actual real followers:    ${real_count}"
    echo "Expected fake followers:  0"
    echo "Actual fake followers:    ${fake_count}"
    echo "========================================"

    local passed=true

    # Check that all fake followers were removed
    if [[ ${fake_count} -ne 0 ]]; then
        log_error "FAIL: ${fake_count} fake followers should have been removed"
        passed=false
    else
        log_test "PASS: All fake followers were removed"
    fi

    # Check that all real followers remain
    if [[ ${real_count} -ne ${#REAL_ACCOUNTS[@]} ]]; then
        log_error "FAIL: Only ${real_count}/${#REAL_ACCOUNTS[@]} real followers remain"
        passed=false
    else
        log_test "PASS: All real followers remain"
    fi

    echo ""
    if [[ "${passed}" == "true" ]]; then
        echo -e "${GREEN}TEST PASSED${NC}"
        return 0
    else
        echo -e "${RED}TEST FAILED${NC}"
        return 1
    fi
}

main() {
    # Record test start time for timeout checking
    TEST_START_TIME=$(date +%s)

    echo ""
    echo "========================================"
    echo "Follower Validation Test"
    echo "========================================"
    echo ""

    # Setup
    setup_temp_dir
    build_owncast
    start_owncast
    configure_owncast

    # Stop Owncast before inserting test data to avoid database contention
    # The database schema is now created, so we can safely insert data
    stop_owncast

    # Insert test data while Owncast is stopped
    insert_followers

    # Restart Owncast to run the validation job
    start_owncast

    # Verify initial state
    log_info "Initial state:"
    log_info "  Total followers: $(get_follower_count)"
    log_info "  Real followers: $(get_real_follower_count)"
    log_info "  Fake followers: $(get_fake_follower_count)"

    # Wait for validation
    wait_for_validation

    # Verify results
    if verify_results; then
        exit 0
    else
        exit 1
    fi
}

main "$@"
