#!/bin/bash
set -e

# Generate mkcert certificates for the test domains.
# These are placed in a container-local path so they don't leak into the
# mounted source tree.
export CERT_DIR="/tmp/test-certs"
mkdir -p "${CERT_DIR}"
mkcert -cert-file "${CERT_DIR}/cert.pem" \
       -key-file "${CERT_DIR}/key.pem" \
       owncast.local snac.local localhost 127.0.0.1

# Change CWD away from the mounted repo root so Owncast doesn't pick up the
# host's (possibly wrong-architecture) ffmpeg binary via ./ffmpeg detection.
# Stay inside the git repo so `git rev-parse --show-toplevel` still works.
cd /owncast/test/automated/activitypub

# Run the specified test script (default: test-federation.sh)
TEST_SCRIPT="${1:-test-federation.sh}"
exec "/owncast/test/automated/activitypub/${TEST_SCRIPT}"
