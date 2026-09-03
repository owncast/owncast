#!/bin/bash
# shellcheck disable=SC2317 # cleanup() is invoked via trap, not direct call
#
# Boots the official Gancio v2 beta image beside Owncast, establishes a real
# trusted-follow relationship, publishes an Owncast ActivityPub Event, and
# verifies that Gancio imports it.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(git rev-parse --show-toplevel)"

PROXY_PORT="${PROXY_PORT:-443}"
OWNCAST_PORT="${OWNCAST_PORT:-8080}"
GANCIO_PORT="${GANCIO_PORT:-13120}"
OWNCAST_HOSTNAME="owncast.local"
GANCIO_HOSTNAME="gancio.local"
OWNCAST_FED_USERNAME="streamer"
PROXY_PORT_SUFFIX=""
[[ "${PROXY_PORT}" == "443" ]] || PROXY_PORT_SUFFIX=":${PROXY_PORT}"
OWNCAST_URL="https://${OWNCAST_HOSTNAME}${PROXY_PORT_SUFFIX}"
GANCIO_URL="https://${GANCIO_HOSTNAME}${PROXY_PORT_SUFFIX}"
GANCIO_IMAGE="${GANCIO_IMAGE:-cisti/gancio@sha256:d50ebfde6c3342cac52ba8121b9d06d87bbf1409e63584b45651b0cbd6a14d3d}"
GANCIO_ADMIN_EMAIL="admin@gancio.test"
GANCIO_ADMIN_PASSWORD="owncast-gancio-test-password"

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
[[ -f "${CERT_DIR}/cert.pem" && -f "${CERT_DIR}/key.pem" && -f "${CERT_DIR}/rootCA.pem" ]] || fail "test certificates are missing"
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
    "/api/admin/config/federation/private|{\"value\":false}" \
    "/api/admin/config/schedule/enabled|{\"value\":true}"; do
    endpoint=${payload%%|*}
    body=${payload#*|}
    curl -sf -X POST "http://127.0.0.1:${OWNCAST_PORT}${endpoint}" \
        -H "Authorization: Basic ${auth}" -H 'Content-Type: application/json' -d "${body}" >/dev/null
done

log "Starting ${GANCIO_IMAGE}"
cp "${CERT_DIR}/rootCA.pem" "${GANCIO_DIR}/rootCA.pem"
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
    -v "${HOST_GANCIO_DIR}/rootCA.pem:/tmp/owncast-ap-rootCA.pem:ro" \
    -v "${HOST_GANCIO_DIR}/data:/app/data" \
    -v "${HOST_GANCIO_DIR}/gancio.env:/app/.env:ro" \
    --entrypoint sh \
    "${GANCIO_IMAGE}" \
    -c 'echo "127.0.0.1 owncast.local gancio.local" >> /etc/hosts && exec su node -s /bin/sh -c "exec env NODE_EXTRA_CA_CERTS=/tmp/owncast-ap-rootCA.pem /usr/local/bin/docker-entrypoint.sh ./server/gancio"' \
    >/dev/null
wait_for "Gancio" "curl --cacert '${CERT_DIR}/cert.pem' -sf '${GANCIO_URL}/'" 45

nodeinfo=$(curl --cacert "${CERT_DIR}/cert.pem" -sf "${GANCIO_URL}/.well-known/nodeinfo")
printf '%s' "${nodeinfo}" | jq -e '.links | length > 0' >/dev/null || fail "Gancio NodeInfo discovery response is invalid"

actor=$(curl --cacert "${CERT_DIR}/cert.pem" -sf -H 'Accept: application/activity+json' "${OWNCAST_URL}/federation/user/${OWNCAST_FED_USERNAME}")
printf '%s' "${actor}" | jq -e --arg id "${OWNCAST_URL}/federation/user/${OWNCAST_FED_USERNAME}" '.id == $id and .type == "Service"' >/dev/null || fail "Owncast actor is not discoverable through the shared proxy"

log "Registering Gancio administrator"
curl --cacert "${CERT_DIR}/cert.pem" -sf -X POST "${GANCIO_URL}/api/user/register" \
    -H 'Content-Type: application/json' \
    -d "$(jq -nc --arg email "${GANCIO_ADMIN_EMAIL}" --arg password "${GANCIO_ADMIN_PASSWORD}" '{email: $email, password: $password}')" \
    >/dev/null
gancio_token=$(curl --cacert "${CERT_DIR}/cert.pem" -sf -X POST "${GANCIO_URL}/api/login/token" \
    -H 'Content-Type: application/json' \
    -d "$(jq -nc --arg email "${GANCIO_ADMIN_EMAIL}" --arg password "${GANCIO_ADMIN_PASSWORD}" '{email: $email, password: $password}')" \
    | jq -r '.access_token // empty')
[[ -n "${gancio_token}" ]] || fail "Gancio login did not return an access token"

log "Trusting and following the Owncast actor from Gancio"
trust_response=""
trust_succeeded=false
for _ in $(seq 1 10); do
    if trust_response=$(curl --cacert "${CERT_DIR}/cert.pem" --fail-with-body -sS -X POST "${GANCIO_URL}/api/ap_actors/add_trust" \
        -H "Authorization: Bearer ${gancio_token}" \
        -H 'Content-Type: application/json' \
        -d "$(jq -nc --arg url "${OWNCAST_URL}/federation/user/${OWNCAST_FED_USERNAME}" '{url: $url}')" 2>/dev/null); then
        trust_succeeded=true
        break
    fi
    sleep 1
done
if [[ "${trust_succeeded}" != "true" ]]; then
    docker logs "${GANCIO_CONTAINER}" >&2 || true
    fail "Gancio could not trust the Owncast actor: ${trust_response}"
fi

follower_count=0
for _ in $(seq 1 90); do
    follower_count=$(curl -sf "http://127.0.0.1:${OWNCAST_PORT}/api/admin/followers?limit=10" \
        -H "Authorization: Basic ${auth}" | jq -r '.total // 0')
    [[ "${follower_count}" -ge 1 ]] && break
    sleep 1
done
if [[ "${follower_count}" -lt 1 ]]; then
    docker logs "${GANCIO_CONTAINER}" >&2 || true
    cat "${TEMP_DIR}/owncast.log" >&2 || true
    fail "Gancio did not become an Owncast follower"
fi
event_name="Owncast Gancio interop ${RANDOM}"
event_start=$(date -u -d '+4 hours' '+%Y-%m-%dT%H:%M:%SZ')
log "Publishing scheduled event ${event_name}"
event_response=$(curl -sf -X POST "http://127.0.0.1:${OWNCAST_PORT}/api/admin/schedule/event" \
    -H "Authorization: Basic ${auth}" \
    -H 'Content-Type: application/json' \
    -d "$(jq -nc --arg name "${event_name}" --arg start "${event_start}" \
        '{name: $name, description: "ActivityPub Event interoperability test", start: $start, durationMinutes: 90, timezone: "UTC"}')")
event_id=$(printf '%s' "${event_response}" | jq -r --arg name "${event_name}" 'first(.events[]? | select(.name == $name) | .id) // empty')
[[ -n "${event_id}" ]] || fail "Owncast schedule API did not return the created event"

event_url="${OWNCAST_URL}/schedule/${event_id}"
object_url="${OWNCAST_URL}/federation/event/${event_id}"
event_object=$(curl --cacert "${CERT_DIR}/cert.pem" -sf -H 'Accept: application/activity+json' "${object_url}")
printf '%s' "${event_object}" | jq -e \
    --arg id "${object_url}" --arg name "${event_name}" --arg url "${event_url}" \
    '.type == "Event" and .id == $id and .name == $name and .url == $url and
     .eventStatus == "EventScheduled" and .joinMode == "none" and
     .location.type == "VirtualLocation" and .location.url == $url and
     (.organizers.totalItems == 1)' >/dev/null \
    || fail "Owncast event object does not match the FEP-8a8e contract"

gancio_event=""
for _ in $(seq 1 45); do
    events=$(curl --cacert "${CERT_DIR}/cert.pem" -sfG "${GANCIO_URL}/api/events" \
        --data-urlencode 'show_federated=true' \
        --data-urlencode "query=${event_name}")
    gancio_event=$(printf '%s' "${events}" | jq -c --arg name "${event_name}" 'first(.[]? | select(.title == $name)) // empty')
    [[ -n "${gancio_event}" ]] && break
    sleep 1
done
if [[ -z "${gancio_event}" ]]; then
    docker logs "${GANCIO_CONTAINER}" >&2 || true
    cat "${TEMP_DIR}/owncast.log" >&2 || true
    fail "Gancio did not import the Owncast event"
fi
printf '%s' "${gancio_event}" | jq -e --arg url "${event_url}" \
    '(.online_locations | index($url)) != null' >/dev/null \
    || fail "Gancio imported the event without the Owncast event URL: ${gancio_event}"

digest=$(docker image inspect "${GANCIO_IMAGE}" --format '{{index .RepoDigests 0}}')
[[ -n "${digest}" ]] || fail "could not resolve Gancio image digest"
log "Gancio imported ${event_name}. Resolved image: ${digest}"

if [[ "${KEEP_RUNNING:-}" == "true" ]]; then
    host_proxy_port="${HOST_PROXY_PORT:-8443}"
    log "Keeping services running. Press Ctrl+C to stop."
    log "Owncast admin: http://localhost:${OWNCAST_PORT}/admin"
    log "Gancio: https://${GANCIO_HOSTNAME}:${host_proxy_port}"
    log "If needed, map ${OWNCAST_HOSTNAME} and ${GANCIO_HOSTNAME} to 127.0.0.1 in /etc/hosts."
    wait
fi
