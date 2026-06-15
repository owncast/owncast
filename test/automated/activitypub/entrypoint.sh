#!/bin/bash
# shellcheck disable=SC2317  # cleanup() is invoked via trap, not direct call
set -e

# Generate mkcert certificates for the test domains.
# These are placed in a container-local path so they don't leak into the
# mounted source tree.
export CERT_DIR="/tmp/test-certs"
mkdir -p "${CERT_DIR}"
mkcert -cert-file "${CERT_DIR}/cert.pem" \
       -key-file "${CERT_DIR}/key.pem" \
       owncast.local snac.local localhost 127.0.0.1

# This container runs as root and the repo is bind-mounted from the host, so
# anything the test writes into the tree (Owncast's ./data dir, created mode
# 0700, and any built owncast binary) lands on the host owned by root and the
# invoking user can't remove it. Hand those artifacts back on exit. HOST_UID/
# HOST_GID are set by run.sh; when unset (entrypoint run directly) this is a
# no-op.
# shellcheck disable=SC2329  # invoked via trap, not called directly
cleanup() {
    if [[ -n "${HOST_UID}" && -n "${HOST_GID}" ]]; then
        chown -R "${HOST_UID}:${HOST_GID}" /owncast/test/automated/activitypub 2>/dev/null || true
    fi
}
trap cleanup EXIT

# Change CWD away from the mounted repo root so Owncast doesn't pick up the
# host's (possibly wrong-architecture) ffmpeg binary via ./ffmpeg detection.
# Stay inside the git repo so `git rev-parse --show-toplevel` still works.
cd /owncast/test/automated/activitypub

# Run the specified test script (default: test-federation.sh). Not exec'd, so
# the cleanup trap above still runs afterward; the test's exit code is preserved.
TEST_SCRIPT="${1:-test-federation.sh}"
set +e
"/owncast/test/automated/activitypub/${TEST_SCRIPT}"
exit_code=$?
exit "${exit_code}"
