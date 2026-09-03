#!/bin/bash
# shellcheck disable=SC2317 # cleanup() is invoked via trap, not direct call
#
# Boots the official Gancio v2 beta image beside an Owncast instance inside the
# existing ActivityPub harness. It proves the cross-instance topology, TLS
# routing, and Gancio's basic federation discovery surface before scheduled
# Event delivery exists in Owncast.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(git rev-parse --show-toplevel)"

PROXY_PORT="${PROXY_PORT:-8443}"
OWNCAST_PORT="${OWNCAST_PORT:-8080}"
GANCIO_PORT="${GANCIO_PORT:-13120}"
OWNCAST_HOSTNAME="owncast.local"
GANCIO_HOSTNAME="gancio.local"
OWNCAST_FED_USERNAME="streamer"
OWNCAST_URL="https://${OWNCAST_HOSTNAME}:${PROXY_PORT}"
GANCIO_URL="https://${GANCIO_HOSTNAME}:${PROXY_PORT}"
GANCIO_IMAGE="${GANCIO_IMAGE:-cisti/gancio@sha256:d50ebfde6c3342cac52ba8121b9d06d87bbf1409e63584b45651b0cbd6a14d3d}"

TEMP_DIR=""
OWNCAST_PID=""
PROXY_PID=""
GANCIO_CONTAINER=""
OWNCAST_BIN=""
GANCIO_DIR=""
HOST_GANCIO_DIR=""
CERT_DIR="${CERT_DIR:-${SCRIPT_DIR}/certs}"
log() { printf '[gancio-v2] %s\n' "$*"; }
fail() { printf '[gancio-v2] ERROR: %s\n' "$*" >&2; exit 1; }

cleanup() {
    if [[ -n "${GANCIO_CONTAINER}" ]]; then
        docker rm -f "${GANCIO_CONTAINER}" >/dev/null 2>&1 || true
    fi
    if [[ -n "${PROXY_PID}" ]] && kill -0 "${PROXY_PID}" 2>/dev/null; then
        kill "${PROXY_PID}" 2>/dev/null || true
        wait "${PROXY_PID}" 2>/dev/null || true
    fi
    if [[ -n "${OWNCAST_PID}" ]] && kill -0 "${OWNCAST_PID}" 2>/dev/null; then
        kill "${OWNCAST_PID}" 2>/dev/null || true
        wait "${OWNCAST_PID}" 2>/dev/null || true
    fi
    [[ -z "${TEMP_DIR}" ]] || rm -rf "${TEMP_DIR}"
    [[ -z "${GANCIO_DIR}" ]] || rm -rf "${GANCIO_DIR}"
}
trap cleanup EXIT

wait_for() {
    local description=$1
    local command=$2
    local attempts=${3:-30}

    for _ in $(seq 1 "${attempts}"); do
        if eval "${command}" >/dev/null 2>&1; then
            return 0
        fi
        sleep 1
    done
    if [[ "${description}" == "Gancio" && -n "${GANCIO_CONTAINER}" ]]; then
        docker logs "${GANCIO_CONTAINER}" >&2 || true
    fi
    fail "${description} did not become ready"
}

[[ -S /var/run/docker.sock ]] || fail "Gancio test requires the Docker socket. Run it through ./run.sh test-gancio-v2.sh"
[[ -n "${HOST_REPO_ROOT:-}" ]] || fail "HOST_REPO_ROOT is missing. Run this test through ./run.sh"
[[ -f "${CERT_DIR}/cert.pem" && -f "${CERT_DIR}/key.pem" ]] || fail "test certificates are missing"
grep -q "${GANCIO_HOSTNAME}" /etc/hosts || fail "${GANCIO_HOSTNAME} is missing from /etc/hosts"

TEMP_DIR=$(mktemp -d)
gancio_dir_name=".gancio-v2-${RANDOM}-${RANDOM}"
GANCIO_DIR="${SCRIPT_DIR}/${gancio_dir_name}"
HOST_GANCIO_DIR="${HOST_REPO_ROOT}/test/automated/activitypub/${gancio_dir_name}"
mkdir -p "${GANCIO_DIR}/data"
chown -R 1000:1000 "${GANCIO_DIR}"
OWNCAST_BIN="${TEMP_DIR}/bin/owncast"
mkdir -p "$(dirname "${OWNCAST_BIN}")"

log "Building Owncast"
(
    cd "${REPO_ROOT}"
    CGO_ENABLED=1 go build -o "${OWNCAST_BIN}" main.go
)

log "Starting Caddy"
PROXY_PORT="${PROXY_PORT}" OWNCAST_PORT="${OWNCAST_PORT}" \
OWNCAST2_PORT="${OWNCAST2_PORT:-8081}" SNAC_PORT="${SNAC_PORT:-9080}" \
GANCIO_PORT="${GANCIO_PORT}" CERT_FILE="${CERT_DIR}/cert.pem" KEY_FILE="${CERT_DIR}/key.pem" \
    caddy run --config "${SCRIPT_DIR}/Caddyfile" --adapter caddyfile \
    >"${TEMP_DIR}/caddy.log" 2>&1 &
PROXY_PID=$!

log "Starting Owncast"
mkdir -p "${TEMP_DIR}/owncast"
(
    cd "${TEMP_DIR}/owncast"
    exec env OWNCAST_ALLOW_INTERNAL_FEDERATION=true OWNCAST_INSECURE_SKIP_VERIFY=true \
        "${OWNCAST_BIN}" -database "${TEMP_DIR}/owncast.db" -webserverport "${OWNCAST_PORT}" -enableVerboseLogging
) >"${TEMP_DIR}/owncast.log" 2>&1 &
OWNCAST_PID=$!
wait_for "Owncast" "curl -sf http://127.0.0.1:${OWNCAST_PORT}/api/status"

auth=$(printf '%s' 'admin:abc123' | base64)
for payload in \
    "/api/admin/config/serverurl|{\"value\":\"${OWNCAST_URL}\"}" \
    "/api/admin/config/federation/username|{\"value\":\"${OWNCAST_FED_USERNAME}\"}" \
    "/api/admin/config/federation/enable|{\"value\":true}" \
    "/api/admin/config/federation/private|{\"value\":false}"; do
    endpoint=${payload%%|*}
    body=${payload#*|}
    curl -sf -X POST "http://127.0.0.1:${OWNCAST_PORT}${endpoint}" \
        -H "Authorization: Basic ${auth}" -H 'Content-Type: application/json' -d "${body}" >/dev/null
done

log "Starting ${GANCIO_IMAGE}"
cp "${CERT_DIR}/cert.pem" "${GANCIO_DIR}/cert.pem"
cat >"${GANCIO_DIR}/gancio.env" <<EOF
BASEURL=${GANCIO_URL}
HOST=0.0.0.0
PORT=${GANCIO_PORT}
DB_DIALECT=sqlite
DB_STORAGE=/app/data/gancio.sqlite
DB_LOGGING=false
UPLOAD_PATH=/app/data/uploads
LOG_PATH=/app/data/logs
NUXT_SESSION_PASSWORD=owncast-gancio-v2-test-session-password
EOF
GANCIO_CONTAINER="owncast-ap-gancio-${RANDOM}-${RANDOM}"
docker run -d --name "${GANCIO_CONTAINER}" \
    --network "container:${HOSTNAME}" \
    --user root \
    -e NODE_EXTRA_CA_CERTS=/tmp/owncast-ap-cert.pem \
    -v "${HOST_GANCIO_DIR}/cert.pem:/tmp/owncast-ap-cert.pem:ro" \
    -v "${HOST_GANCIO_DIR}/data:/app/data" \
    -v "${HOST_GANCIO_DIR}/gancio.env:/app/.env:ro" \
    --entrypoint sh \
    "${GANCIO_IMAGE}" \
    -c 'echo "127.0.0.1 owncast.local gancio.local" >> /etc/hosts && exec su node -s /bin/sh -c "exec /usr/local/bin/docker-entrypoint.sh ./server/gancio"' \
    >/dev/null
wait_for "Gancio" "curl --cacert '${CERT_DIR}/cert.pem' -sf '${GANCIO_URL}/'" 45

nodeinfo=$(curl --cacert "${CERT_DIR}/cert.pem" -sf "${GANCIO_URL}/.well-known/nodeinfo")
printf '%s' "${nodeinfo}" | jq -e '.links | length > 0' >/dev/null || fail "Gancio NodeInfo discovery response is invalid"

actor=$(curl --cacert "${CERT_DIR}/cert.pem" -sf -H 'Accept: application/activity+json' "${OWNCAST_URL}/federation/user/${OWNCAST_FED_USERNAME}")
printf '%s' "${actor}" | jq -e --arg id "${OWNCAST_URL}/federation/user/${OWNCAST_FED_USERNAME}" '.id == $id and .type == "Service"' >/dev/null || fail "Owncast actor is not discoverable through the shared proxy"

digest=$(docker image inspect "${GANCIO_IMAGE}" --format '{{index .RepoDigests 0}}')
[[ -n "${digest}" ]] || fail "could not resolve Gancio image digest"
log "Gancio v2 sidecar is reachable. Resolved image: ${digest}"
