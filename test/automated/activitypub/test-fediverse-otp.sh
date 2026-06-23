#!/bin/bash

# End-to-end Fediverse authentication test.
#
# Exercises the real OTP-over-ActivityPub login flow against a live snac2
# instance, with no stubs:
#   1. Register an anonymous chat user on Owncast (gets an access token).
#   2. Ask Owncast to link a Fediverse account; Owncast DMs a one-time code to
#      that account over real federation (webfinger + signed inbox POST).
#   3. Read the code back out of snac2's stored DM (same approach the
#      federation test uses to confirm message delivery).
#   4. Submit the code to Owncast's verify endpoint.
#   5. Assert the admin users API now reports the user as authenticated with a
#      "Fediverse" provider.
#
# Runs inside the same Docker harness as the federation tests:
#   ./run.sh test-fediverse-otp.sh
#
# Reuses the activitypub harness's certs (mkcert), Caddy reverse proxy, and
# snac2 build. No backend code changes are needed: OWNCAST_ALLOW_INTERNAL_FEDERATION
# and OWNCAST_INSECURE_SKIP_VERIFY let the federated DM flow run against localhost.

set -o pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

PROXY_PORT="${PROXY_PORT:-8443}"
SNAC_PORT="${SNAC_PORT:-9080}"
SNAC_HOSTNAME="snac.local"
OWNCAST_PORT="${OWNCAST_PORT:-8080}"
OWNCAST_HOSTNAME="owncast.local"

OWNCAST_URL="https://${OWNCAST_HOSTNAME}:${PROXY_PORT}"

ADMIN_USER="admin"
ADMIN_PASS="abc123"
FEDERATION_USERNAME="streamer"

# The Fediverse account we'll authenticate against — a user we create in snac2.
SNAC_USERNAME="otptester"
FEDIVERSE_ACCOUNT="${SNAC_USERNAME}@${SNAC_HOSTNAME}:${PROXY_PORT}"

TEMP_DIR=""
SNAC_DATA_DIR=""
SNAC_BIN=""
OWNCAST_BIN=""
OWNCAST_DB=""
SNAC_PID=""
OWNCAST_PID=""
PROXY_PID=""

RED='\033[0;31m'; GREEN='\033[0;32m'; CYAN='\033[0;36m'; NC='\033[0m'
log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_test() { echo -e "${CYAN}[TEST]${NC} $1"; }

# shellcheck disable=SC2329  # invoked via trap
cleanup() {
	log_info "Cleaning up..."
	for pid in "${PROXY_PID}" "${OWNCAST_PID}" "${SNAC_PID}"; do
		if [[ -n "${pid}" ]] && kill -0 "${pid}" 2>/dev/null; then
			kill "${pid}" 2>/dev/null || true
			wait "${pid}" 2>/dev/null || true
		fi
	done
	pkill -f "snac httpd /tmp" 2>/dev/null || true
	[[ -n "${TEMP_DIR}" && -d "${TEMP_DIR}" ]] && rm -rf "${TEMP_DIR}"
	log_info "Cleanup complete."
}
trap cleanup EXIT

setup_temp_dir() {
	TEMP_DIR=$(mktemp -d)
	SNAC_DATA_DIR="${TEMP_DIR}/snac-data"
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

install_snac2() {
	if command -v snac &>/dev/null; then
		SNAC_BIN=$(command -v snac)
		log_info "Using system snac2: ${SNAC_BIN}"
		return
	fi
	local snac_src="${TEMP_DIR}/snac2-src"
	log_info "Building snac2..."
	git clone --depth 1 https://codeberg.org/grunfink/snac2.git "${snac_src}" 2>/dev/null
	(cd "${snac_src}" && make >/dev/null 2>&1)
	SNAC_BIN="${snac_src}/snac"
}

setup_snac() {
	log_info "Initializing snac2 and creating @${SNAC_USERNAME}..."
	printf "127.0.0.1\n%s\n%s\n\ntest@test.local\n" "${SNAC_PORT}" "${SNAC_HOSTNAME}:${PROXY_PORT}" |
		"${SNAC_BIN}" init "${SNAC_DATA_DIR}" >/dev/null 2>&1
	if [[ ! -f "${SNAC_DATA_DIR}/server.json" ]]; then
		log_error "snac2 init failed"
		exit 1
	fi
	printf "%s\n%s\n" "${SNAC_USERNAME}" "OTP Tester" |
		"${SNAC_BIN}" adduser "${SNAC_DATA_DIR}" >/dev/null 2>&1
}

start_snac() {
	log_info "Starting snac2..."
	DEBUG=0 "${SNAC_BIN}" httpd "${SNAC_DATA_DIR}" >"${TEMP_DIR}/snac2.log" 2>&1 &
	SNAC_PID=$!
	for _ in $(seq 1 30); do
		curl -s "http://127.0.0.1:${SNAC_PORT}/" >/dev/null 2>&1 && { log_info "snac2 ready"; return 0; }
		sleep 1
	done
	log_error "snac2 did not become ready"; return 1
}

start_proxy() {
	log_info "Starting Caddy reverse proxy..."
	export PROXY_PORT OWNCAST_PORT SNAC_PORT
	export OWNCAST2_PORT="${OWNCAST2_PORT:-8081}"
	caddy run --config "${SCRIPT_DIR}/Caddyfile" --adapter caddyfile >"${TEMP_DIR}/caddy.log" 2>&1 &
	PROXY_PID=$!
	sleep 2
	for _ in $(seq 1 10); do
		curl -sk "https://127.0.0.1:${PROXY_PORT}/" >/dev/null 2>&1 && { log_info "Caddy ready"; return 0; }
		sleep 1
	done
	log_error "Caddy did not become ready"; return 1
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

owncast_admin_post() {
	local path="$1" body="$2"
	curl -s -X POST "http://localhost:${OWNCAST_PORT}/api/admin/${path}" \
		-u "${ADMIN_USER}:${ADMIN_PASS}" -H "Content-Type: application/json" -d "${body}" >/dev/null
}

configure_owncast() {
	log_info "Configuring Owncast federation (${OWNCAST_URL})..."
	owncast_admin_post "config/serverurl" "{\"value\": \"${OWNCAST_URL}\"}"
	owncast_admin_post "config/federation/username" "{\"value\": \"${FEDERATION_USERNAME}\"}"
	owncast_admin_post "config/federation/enable" '{"value": true}'
	owncast_admin_post "config/federation/private" '{"value": false}'
	sleep 2
}

# ---- the actual test ----------------------------------------------------

run_otp_test() {
	local base="http://localhost:${OWNCAST_PORT}"

	log_test "Registering an anonymous chat user..."
	local reg user_id token display
	reg=$(curl -s -X POST "${base}/api/chat/register")
	user_id=$(echo "${reg}" | jq -r '.id')
	token=$(echo "${reg}" | jq -r '.accessToken')
	display=$(echo "${reg}" | jq -r '.displayName')
	if [[ -z "${token}" || "${token}" == "null" ]]; then
		log_error "chat register failed: ${reg}"; return 1
	fi
	log_info "user_id=${user_id} display=${display}"

	log_test "Requesting a Fediverse OTP for ${FEDIVERSE_ACCOUNT}..."
	local otp_resp
	otp_resp=$(curl -s -X POST "${base}/api/auth/fediverse?accessToken=${token}" \
		-H "Content-Type: application/json" -d "{\"account\": \"${FEDIVERSE_ACCOUNT}\"}")
	if ! echo "${otp_resp}" | grep -q '"success":true'; then
		log_error "OTP registration failed: ${otp_resp}"
		log_error "owncast log tail:"; tail -20 "${TEMP_DIR}/owncast.log" 2>/dev/null
		return 1
	fi

	log_test "Waiting for the OTP to arrive as a Fediverse DM in snac2..."
	local code=""
	for _ in $(seq 1 30); do
		code=$(grep -rl 'One-time code' "${SNAC_DATA_DIR}" 2>/dev/null |
			xargs -r grep -hoE 'One-time code[^0-9]*[0-9]{6}' 2>/dev/null |
			grep -oE '[0-9]{6}' | head -1)
		[[ -n "${code}" ]] && break
		sleep 1
	done
	if [[ -z "${code}" ]]; then
		log_error "OTP code never arrived in snac2's inbox"
		log_error "snac2 log tail:"; tail -20 "${TEMP_DIR}/snac2.log" 2>/dev/null
		return 1
	fi
	log_info "Received OTP via federation: ${code}"

	log_test "Verifying the OTP code..."
	local verify_resp
	verify_resp=$(curl -s -X POST "${base}/api/auth/fediverse/verify?accessToken=${token}" \
		-H "Content-Type: application/json" -d "{\"code\": \"${code}\"}")
	if ! echo "${verify_resp}" | grep -q '"success":true'; then
		log_error "OTP verification failed: ${verify_resp}"; return 1
	fi

	log_test "Asserting the admin users API reports the Fediverse identity..."
	local users
	users=$(curl -s -u "${ADMIN_USER}:${ADMIN_PASS}" "${base}/api/admin/users?limit=200")
	if echo "${users}" | jq -e --arg id "${user_id}" \
		'.results[] | select(.id==$id) | (.authenticated==true and ((.authProviders // []) | index("Fediverse")))' >/dev/null; then
		log_test "${GREEN}PASS${NC}: user is authenticated and shown with the Fediverse provider"
		return 0
	fi
	log_error "user ${user_id} is not shown as Fediverse-authenticated"
	echo "${users}" | jq --arg id "${user_id}" '.results[] | select(.id==$id)' 2>/dev/null
	return 1
}

main() {
	echo "=================================================="
	echo " Fediverse OTP authentication — end-to-end test"
	echo "=================================================="
	setup_temp_dir
	check_certs
	install_snac2
	setup_snac
	build_owncast
	start_snac
	start_proxy
	start_owncast
	configure_owncast

	if run_otp_test; then
		echo ""
		log_test "${GREEN}Fediverse OTP auth test passed.${NC}"
		exit 0
	else
		echo ""
		log_error "Fediverse OTP auth test FAILED."
		exit 1
	fi
}

main "$@"
