#!/bin/bash
# Persistent two-instance Owncast federation dev environment.
#
# Brings up two Owncast instances (owncast.local + owncast2.local) behind a
# shared HTTPS Caddy proxy, both configured for federation, and LEAVES THEM
# RUNNING so you can iterate on featured streams locally instead of deploying
# to production. Ctrl-C tears everything down.
#
# This reuses the exact mechanism the CI test (test-featured-streams.sh) relies
# on -- HTTPS via Caddy + mkcert, plus the two env vars that let federation talk
# to internal/self-signed hosts:
#   OWNCAST_ALLOW_INTERNAL_FEDERATION=true  (allow owncast.local / loopback)
#   OWNCAST_INSECURE_SKIP_VERIFY=true       (accept the mkcert self-signed cert)
#
# One-time prerequisites (see setup.sh):
#   - caddy, mkcert, jq, ffmpeg, go installed
#   - mkcert CA trusted:  mkcert -install
#   - certs present (regenerate to include owncast2.local):
#       CAROOT="$HOME/.local/share/mkcert" mkcert \
#         -cert-file certs/cert.pem -key-file certs/key.pem \
#         owncast.local owncast2.local snac.local localhost 127.0.0.1
#   - /etc/hosts entry:
#       127.0.0.1 owncast.local owncast2.local snac.local
#
# Usage:
#   ./dev-up.sh              # build, start both instances, leave running
#   ./dev-up.sh --no-build   # skip 'go build', reuse the existing binary
#   FRESH=true ./dev-up.sh   # wipe instance data dirs first (clean slate)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(git rev-parse --show-toplevel)"

PROXY_PORT="${PROXY_PORT:-443}"
OWNCAST_PORT="${OWNCAST_PORT:-8080}"
OWNCAST2_PORT="${OWNCAST2_PORT:-8081}"
OWNCAST_RTMP_PORT="${OWNCAST_RTMP_PORT:-1935}"
OWNCAST2_RTMP_PORT="${OWNCAST2_RTMP_PORT:-1936}"
SNAC_PORT="${SNAC_PORT:-9080}" # unused; keeps the shared Caddyfile valid

OWNCAST_HOSTNAME="owncast.local"
OWNCAST2_HOSTNAME="owncast2.local"
OWNCAST_FED_USERNAME="streamerone"
OWNCAST2_FED_USERNAME="streamertwo"

ADMIN_USER="admin"
ADMIN_PASS="abc123"

# The admin UI rejects any federation URL whose port isn't 443 (see
# web/components/admin/FederatedServers/FeatureStreamModal.tsx). Running the
# proxy on 443 keeps URLs portless so the real UI accepts them and the setup
# mirrors production. Override PROXY_PORT for a non-privileged port (CLI/curl
# flow only).
if [[ "${PROXY_PORT}" == "443" ]]; then
    OWNCAST_URL="https://${OWNCAST_HOSTNAME}"
    OWNCAST2_URL="https://${OWNCAST2_HOSTNAME}"
else
    OWNCAST_URL="https://${OWNCAST_HOSTNAME}:${PROXY_PORT}"
    OWNCAST2_URL="https://${OWNCAST2_HOSTNAME}:${PROXY_PORT}"
fi

DEVDATA="${SCRIPT_DIR}/devdata"
CERT_FILE="${SCRIPT_DIR}/certs/cert.pem"
KEY_FILE="${SCRIPT_DIR}/certs/key.pem"
OWNCAST_BIN="${DEVDATA}/owncast"

# Web dev servers (one per instance) so the current web UI can be tested on
# both without rebuilding the embedded bundle. Set SKIP_WEB=true for backend
# only.
WEB_DIR="${REPO_ROOT}/web"
WEB_PORT_A="${WEB_PORT_A:-3000}"
WEB_PORT_B="${WEB_PORT_B:-3001}"

GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m'
info()  { echo -e "${GREEN}[dev-up]${NC} $1"; }
warn()  { echo -e "${YELLOW}[dev-up]${NC} $1"; }
err()   { echo -e "${RED}[dev-up]${NC} $1"; }

CADDY_PID=""
OC1_PID=""
OC2_PID=""
WEB_A_PID=""
WEB_B_PID=""

# kill_tree kills a process and all its descendants. npm spawns next, which
# spawns workers; killing only the npm pid would orphan them, so recurse.
kill_tree() {
    local pid=$1
    [[ -z "${pid}" ]] && return
    local child
    for child in $(pgrep -P "${pid}" 2>/dev/null); do
        kill_tree "${child}"
    done
    kill "${pid}" 2>/dev/null || true
}

cleanup() {
    info "Shutting down..."
    kill_tree "${WEB_A_PID}"
    kill_tree "${WEB_B_PID}"
    [[ -n "${OC1_PID}" ]]   && kill "${OC1_PID}"   2>/dev/null || true
    [[ -n "${OC2_PID}" ]]   && kill "${OC2_PID}"   2>/dev/null || true
    [[ -n "${CADDY_PID}" ]] && kill "${CADDY_PID}" 2>/dev/null || true
    wait 2>/dev/null || true
    info "Done."
}
trap cleanup EXIT INT TERM

# --- preflight ---------------------------------------------------------------
for tool in caddy mkcert jq go; do
    command -v "${tool}" >/dev/null 2>&1 || { err "missing required tool: ${tool}"; exit 1; }
done

# Locate npm. nvm-managed node is often not on a non-interactive PATH, so fall
# back to the newest ~/.nvm node. NPM_BIN is set so its dir can be prepended to
# PATH when launching (next needs node on PATH).
NPM_BIN=""
if [[ "${SKIP_WEB:-}" != "true" ]]; then
    if command -v npm >/dev/null 2>&1; then
        NPM_BIN="$(command -v npm)"
    else
        NPM_BIN="$(find "${HOME}/.nvm/versions/node" -maxdepth 3 -name npm 2>/dev/null | sort -V | tail -1)"
    fi
    if [[ -z "${NPM_BIN}" ]]; then
        err "npm not found (and SKIP_WEB!=true). Install node/npm or run with SKIP_WEB=true."
        exit 1
    fi
    [[ -d "${WEB_DIR}/node_modules" ]] || { err "${WEB_DIR}/node_modules missing -- run 'npm install' in web/ first (or SKIP_WEB=true)."; exit 1; }
fi
[[ -f "${CERT_FILE}" && -f "${KEY_FILE}" ]] || { err "certs missing; see header for the mkcert command"; exit 1; }
# Caddy needs permission to bind the privileged port (443 by default) without
# root. Grant the capability once with setcap.
if [[ "${PROXY_PORT}" -lt 1024 ]]; then
    caddy_bin="$(command -v caddy)"
    if ! getcap "${caddy_bin}" 2>/dev/null | grep -q cap_net_bind_service; then
        err "Caddy can't bind port ${PROXY_PORT} without privilege. Grant it once:"
        err "  sudo setcap 'cap_net_bind_service=+ep' ${caddy_bin}"
        err "(or run with a non-privileged port: PROXY_PORT=8443 ./dev-up.sh -- but then the admin UI rejects the URL; use the curl flow.)"
        exit 1
    fi
fi
if ! grep -q "owncast2.local" /etc/hosts; then
    err "owncast2.local not in /etc/hosts. Run:"
    err "  sudo sed -i 's/^127\\.0\\.0\\.1 owncast\\.local snac\\.local/127.0.0.1 owncast.local owncast2.local snac.local/' /etc/hosts"
    exit 1
fi

if [[ "${FRESH:-}" == "true" ]]; then
    warn "FRESH=true -> wiping instance data"
    rm -rf "${DEVDATA}/oc1" "${DEVDATA}/oc2" "${DEVDATA}/oc1.db" "${DEVDATA}/oc2.db"
fi
mkdir -p "${DEVDATA}/oc1" "${DEVDATA}/oc2"

# --- build -------------------------------------------------------------------
if [[ "${1:-}" != "--no-build" ]]; then
    info "Building Owncast (CGO_ENABLED=1)..."
    ( cd "${REPO_ROOT}" && CGO_ENABLED=1 go build -o "${OWNCAST_BIN}" main.go )
else
    [[ -x "${OWNCAST_BIN}" ]] || { err "--no-build but no binary at ${OWNCAST_BIN}"; exit 1; }
    info "Reusing existing binary (--no-build)."
fi

# --- caddy -------------------------------------------------------------------
info "Starting Caddy HTTPS proxy on :${PROXY_PORT}..."
PROXY_PORT="${PROXY_PORT}" CERT_FILE="${CERT_FILE}" KEY_FILE="${KEY_FILE}" \
OWNCAST_PORT="${OWNCAST_PORT}" OWNCAST2_PORT="${OWNCAST2_PORT}" SNAC_PORT="${SNAC_PORT}" \
    caddy run --config "${SCRIPT_DIR}/Caddyfile" --adapter caddyfile \
    > "${DEVDATA}/caddy.log" 2>&1 &
CADDY_PID=$!

for _ in $(seq 1 10); do
    curl -sk "https://127.0.0.1:${PROXY_PORT}/" >/dev/null 2>&1 && break
    sleep 1
done

# start_instance LABEL WORKDIR DB WEB_PORT RTMP_PORT -> sets _PID
start_instance() {
    local label=$1 workdir=$2 db=$3 web_port=$4 rtmp_port=$5
    local log="${DEVDATA}/${label}.log"
    info "Starting '${label}' (web ${web_port}, rtmp ${rtmp_port}, log ${log})..."
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
    _PID=$!
    for _ in $(seq 1 30); do
        curl -s "http://localhost:${web_port}/api/status" >/dev/null 2>&1 && { info "'${label}' ready"; return 0; }
        kill -0 "${_PID}" 2>/dev/null || { err "'${label}' exited early:"; tail -30 "${log}"; return 1; }
        sleep 1
    done
    err "'${label}' did not become ready:"; tail -30 "${log}"; return 1
}

auth() { echo -n "${ADMIN_USER}:${ADMIN_PASS}" | base64; }

# configure_instance WEB_PORT SERVER_URL FED_USERNAME
configure_instance() {
    local web_port=$1 server_url=$2 fed_username=$3
    local base="http://localhost:${web_port}" a; a=$(auth)
    info "Configuring ${base} as ${server_url} (user ${fed_username})"
    for kv in \
        "config/serverurl:{\"value\": \"${server_url}\"}" \
        "config/federation/username:{\"value\": \"${fed_username}\"}" \
        "config/federation/enable:{\"value\": true}" \
        "config/federation/private:{\"value\": false}"; do
        local path="${kv%%:*}" body="${kv#*:}"
        curl -s -X POST "${base}/api/admin/${path}" \
            -H "Authorization: Basic ${a}" -H "Content-Type: application/json" \
            -d "${body}" >/dev/null
    done
}

start_instance oc1 "${DEVDATA}/oc1" "${DEVDATA}/oc1.db" "${OWNCAST_PORT}" "${OWNCAST_RTMP_PORT}"
OC1_PID=$_PID
start_instance oc2 "${DEVDATA}/oc2" "${DEVDATA}/oc2.db" "${OWNCAST2_PORT}" "${OWNCAST2_RTMP_PORT}"
OC2_PID=$_PID

configure_instance "${OWNCAST_PORT}"  "${OWNCAST_URL}"  "${OWNCAST_FED_USERNAME}"
configure_instance "${OWNCAST2_PORT}" "${OWNCAST2_URL}" "${OWNCAST2_FED_USERNAME}"

# start_web LABEL WEB_PORT BACKEND_PORT DISTDIR -> sets _WEBPID
# Runs `next dev` from web/ proxied to the given backend. Each server uses its
# own distDir so the two don't fight over the .next build cache.
start_web() {
    local label=$1 web_port=$2 backend_port=$3 distdir=$4
    local log="${DEVDATA}/${label}.log"
    info "Starting web dev '${label}' on :${web_port} -> backend :${backend_port} (log ${log})..."
    (
        cd "${WEB_DIR}" || exit 1
        # NEXT_PUBLIC_API_HOST=/ forces the admin's API base to be relative so
        # its calls ride THIS dev server's proxy to ITS backend. .env.development
        # hardcodes http://localhost:8080 (absolute -> always backend A), which
        # breaks the second instance's admin. Real env vars win over .env files
        # in Next, so this overrides it without editing the shared file.
        exec env \
            PATH="$(dirname "${NPM_BIN}"):${PATH}" \
            OWNCAST_DEV_BACKEND="http://localhost:${backend_port}" \
            OWNCAST_DEV_DISTDIR="${distdir}" \
            NEXT_PUBLIC_API_HOST="/" \
            "${NPM_BIN}" run dev -- -p "${web_port}"
    ) > "${log}" 2>&1 &
    _WEBPID=$!
    local attempt=0
    while [[ ${attempt} -lt 60 ]]; do
        if curl -s -o /dev/null "http://localhost:${web_port}/" 2>/dev/null; then
            info "web '${label}' ready (http://localhost:${web_port}/admin)"
            return 0
        fi
        kill -0 "${_WEBPID}" 2>/dev/null || { err "web '${label}' exited early:"; tail -30 "${log}"; return 1; }
        attempt=$((attempt + 1))
        sleep 1
    done
    warn "web '${label}' not ready after 60s; check ${log} (continuing anyway)"
}

if [[ "${SKIP_WEB:-}" != "true" ]]; then
    start_web webA "${WEB_PORT_A}" "${OWNCAST_PORT}"  ".next"
    WEB_A_PID=$_WEBPID
    start_web webB "${WEB_PORT_B}" "${OWNCAST2_PORT}" ".next-dev-b"
    WEB_B_PID=$_WEBPID
fi

A=$(auth)
cat <<EOF

$(echo -e "${GREEN}=========================================================${NC}")
 Two-instance federation env is UP. Leave this running.
$(echo -e "${GREEN}=========================================================${NC}")

 Use these LIVE web UIs (current source, hot-reload, no rebuild):
   Instance A admin:  http://localhost:${WEB_PORT_A}/admin   -> backend A (:${OWNCAST_PORT})
   Instance B admin:  http://localhost:${WEB_PORT_B}/admin   -> backend B (:${OWNCAST2_PORT})
 Admin login:         ${ADMIN_USER} / ${ADMIN_PASS}
 Logs:                ${DEVDATA}/{oc1,oc2,webA,webB,caddy(.log)}

 (Backends also serve their own COMPILED bundle at ${OWNCAST_URL}/admin and
  ${OWNCAST2_URL}/admin, but that is stale until rebuilt -- use the dev UIs
  above for web work. Run with SKIP_WEB=true to skip the dev servers.)

 ----- Drive the featured-streams flow from the CLI -----

 # A features B (sends a Follow that B must approve):
 curl -s -X POST http://localhost:${OWNCAST_PORT}/api/admin/federation/servers \\
   -H "Authorization: Basic ${A}" -H "Content-Type: application/json" \\
   -d '{"url": "${OWNCAST2_URL}"}' | jq

 # B sees a pending feature request:
 curl -s http://localhost:${OWNCAST2_PORT}/api/admin/federation/feature-requests \\
   -H "Authorization: Basic ${A}" | jq

 # B approves A (use the actorIRI from the feature-requests list):
 curl -s -X POST http://localhost:${OWNCAST2_PORT}/api/admin/followers/approve \\
   -H "Authorization: Basic ${A}" -H "Content-Type: application/json" \\
   -d '{"actorIRI": "PASTE_IRI", "approved": true}' | jq

 # A's directory (should show B, then flip to accepted + metadata):
 curl -s http://localhost:${OWNCAST_PORT}/api/federation/servers | jq

 # Make B go live so A's entry flips to live:
 STREAM_KEY=abc123 ${REPO_ROOT}/test/ocTestStream.sh rtmp://localhost:${OWNCAST2_RTMP_PORT}/live

 Ctrl-C to stop everything.

EOF

info "Tailing instance logs (Ctrl-C to stop)..."
wait
