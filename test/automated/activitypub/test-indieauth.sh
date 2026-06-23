#!/bin/bash

# End-to-end IndieAuth authentication test.
#
# Stands up a FAKE IndieAuth provider (served by Caddy at https://indieauth.local)
# and drives Owncast's IndieAuth client flow against it, with no stubs in the
# backend:
#   1. Register an anonymous chat user on Owncast (gets an access token).
#   2. Start the IndieAuth flow pointed at the fake provider; Owncast discovers
#      the provider's authorization endpoint and returns a redirect URL.
#   3. Follow that redirect like a browser would: the fake provider auto-approves
#      and bounces back to Owncast's callback with a code; Owncast exchanges the
#      code with the fake provider and links the verified identity.
#   4. Assert the admin users API now reports the user as authenticated with an
#      "IndieAuth" provider.
#
# Runs in the same Docker harness as the federation tests:
#   ./run.sh test-indieauth.sh
#
# No snac2 is needed. OWNCAST_ALLOW_INTERNAL_FEDERATION lets Owncast accept the
# local fake host; the mkcert CA (installed in the container) makes Owncast's
# own HTTPS calls to the provider trust the test certificate.

set -o pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

PROXY_PORT="${PROXY_PORT:-8443}"
OWNCAST_PORT="${OWNCAST_PORT:-8080}"
OWNCAST_HOSTNAME="owncast.local"
INDIEAUTH_HOSTNAME="indieauth.local"

OWNCAST_URL="https://${OWNCAST_HOSTNAME}:${PROXY_PORT}"
INDIEAUTH_HOST="https://${INDIEAUTH_HOSTNAME}:${PROXY_PORT}/"

ADMIN_USER="admin"
ADMIN_PASS="abc123"

TEMP_DIR=""
OWNCAST_BIN=""
OWNCAST_DB=""
OWNCAST_PID=""
PROXY_PID=""

RED='\033[0;31m'; GREEN='\033[0;32m'; CYAN='\033[0;36m'; NC='\033[0m'
log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_test() { echo -e "${CYAN}[TEST]${NC} $1"; }

# shellcheck disable=SC2329  # invoked via trap
cleanup() {
	log_info "Cleaning up..."
	for pid in "${PROXY_PID}" "${OWNCAST_PID}"; do
		if [[ -n "${pid}" ]] && kill -0 "${pid}" 2>/dev/null; then
			kill "${pid}" 2>/dev/null || true
			wait "${pid}" 2>/dev/null || true
		fi
	done
	[[ -n "${TEMP_DIR}" && -d "${TEMP_DIR}" ]] && rm -rf "${TEMP_DIR}"
	log_info "Cleanup complete."
}
trap cleanup EXIT

setup_temp_dir() {
	TEMP_DIR=$(mktemp -d)
	OWNCAST_DB="${TEMP_DIR}/owncast.db"
	log_info "Temp directory: ${TEMP_DIR}"
}

check_certs() {
	CERT_DIR="${CERT_DIR:-${SCRIPT_DIR}/certs}"
	if [[ ! -f "${CERT_DIR}/cert.pem" || ! -f "${CERT_DIR}/key.pem" ]]; then
		log_error "Certificates not found in ${CERT_DIR} (run inside the Docker harness, or see README.md)"
		exit 1
	fi
	export CERT_FILE="${CERT_DIR}/cert.pem"
	export KEY_FILE="${CERT_DIR}/key.pem"
	log_info "Using certificates from ${CERT_DIR}"
}

build_owncast() {
	log_info "Building Owncast..."
	OWNCAST_BIN="${TEMP_DIR}/owncast"
	(cd "$(git rev-parse --show-toplevel)" && CGO_ENABLED=1 go build -o "${OWNCAST_BIN}" main.go) ||
		{ log_error "Owncast build failed"; exit 1; }
}

start_owncast() {
	log_info "Starting Owncast..."
	OWNCAST_ALLOW_INTERNAL_FEDERATION=true \
		OWNCAST_INSECURE_SKIP_VERIFY=true \
		"${OWNCAST_BIN}" -database "${OWNCAST_DB}" >"${TEMP_DIR}/owncast.log" 2>&1 &
	OWNCAST_PID=$!
	for _ in $(seq 1 30); do
		curl -s "http://localhost:${OWNCAST_PORT}/api/status" >/dev/null 2>&1 && { log_info "Owncast ready"; return 0; }
		sleep 1
	done
	log_error "Owncast did not become ready"; return 1
}

configure_owncast() {
	# IndieAuth requires the Owncast server URL to be set (it becomes the OAuth
	# client_id and the redirect_uri base).
	log_info "Configuring Owncast server URL (${OWNCAST_URL})..."
	curl -s -X POST "http://localhost:${OWNCAST_PORT}/api/admin/config/serverurl" \
		-u "${ADMIN_USER}:${ADMIN_PASS}" -H "Content-Type: application/json" \
		-d "{\"value\": \"${OWNCAST_URL}\"}" >/dev/null
	sleep 1
}

start_proxy() {
	log_info "Starting Caddy (Owncast + fake IndieAuth provider)..."
	export PROXY_PORT OWNCAST_PORT
	caddy run --config "${SCRIPT_DIR}/Caddyfile.indieauth" --adapter caddyfile >"${TEMP_DIR}/caddy.log" 2>&1 &
	PROXY_PID=$!
	sleep 2
	for _ in $(seq 1 10); do
		curl -sk "https://127.0.0.1:${PROXY_PORT}/" >/dev/null 2>&1 && { log_info "Caddy ready"; return 0; }
		sleep 1
	done
	log_error "Caddy did not become ready"; return 1
}

# ---- the actual test ----------------------------------------------------

run_indieauth_test() {
	local base="http://localhost:${OWNCAST_PORT}"

	log_test "Sanity-check: the fake IndieAuth provider advertises its endpoint..."
	if ! curl -sk "${INDIEAUTH_HOST}" | grep -q 'rel="authorization_endpoint"'; then
		log_error "fake provider discovery document missing the authorization_endpoint link"
		return 1
	fi

	log_test "Registering an anonymous chat user..."
	local reg user_id token
	reg=$(curl -s -X POST "${base}/api/chat/register")
	user_id=$(echo "${reg}" | jq -r '.id')
	token=$(echo "${reg}" | jq -r '.accessToken')
	if [[ -z "${token}" || "${token}" == "null" ]]; then
		log_error "chat register failed: ${reg}"; return 1
	fi
	log_info "user_id=${user_id}"

	log_test "Starting the IndieAuth flow against ${INDIEAUTH_HOST}..."
	local start_resp redirect
	start_resp=$(curl -s -X POST "${base}/api/auth/indieauth?accessToken=${token}" \
		-H "Content-Type: application/json" -d "{\"authHost\": \"${INDIEAUTH_HOST}\"}")
	redirect=$(echo "${start_resp}" | jq -r '.redirect')
	if [[ -z "${redirect}" || "${redirect}" == "null" ]]; then
		log_error "StartAuthFlow returned no redirect: ${start_resp}"
		log_error "owncast log tail:"; tail -20 "${TEMP_DIR}/owncast.log" 2>/dev/null
		return 1
	fi
	log_info "authorization redirect: ${redirect}"

	log_test "Following the authorization redirect through the fake provider and back to Owncast..."
	# Like a browser: GET the provider's /auth (302 back to Owncast's callback
	# with a code), then Owncast's callback exchanges the code with the provider
	# and links the identity. -L follows the whole chain; -k because the browser
	# leg hits the mkcert cert by IP/host.
	local final_code
	final_code=$(curl -Lsk -o /dev/null -w "%{http_code}" "${redirect}")
	log_info "final response after callback chain: HTTP ${final_code}"

	log_test "Asserting the admin users API reports the IndieAuth identity..."
	local users
	users=$(curl -s -u "${ADMIN_USER}:${ADMIN_PASS}" "${base}/api/admin/users?limit=200")
	if echo "${users}" | jq -e --arg id "${user_id}" \
		'.results[] | select(.id==$id) | (.authenticated==true and ((.authProviders // []) | index("IndieAuth")))' >/dev/null; then
		log_test "${GREEN}PASS${NC}: user is authenticated and shown with the IndieAuth provider"
		return 0
	fi
	log_error "user ${user_id} is not shown as IndieAuth-authenticated"
	echo "${users}" | jq --arg id "${user_id}" '.results[] | select(.id==$id)' 2>/dev/null
	log_error "owncast log tail:"; tail -20 "${TEMP_DIR}/owncast.log" 2>/dev/null
	return 1
}

main() {
	echo "=================================================="
	echo " IndieAuth authentication — end-to-end test"
	echo "=================================================="
	setup_temp_dir
	check_certs
	build_owncast
	start_proxy
	start_owncast
	configure_owncast

	if run_indieauth_test; then
		echo ""
		log_test "${GREEN}IndieAuth auth test passed.${NC}"
		exit 0
	else
		echo ""
		log_error "IndieAuth auth test FAILED."
		exit 1
	fi
}

main "$@"
