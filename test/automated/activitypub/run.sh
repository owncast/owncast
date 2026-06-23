#!/bin/bash

# Run the ActivityPub federation test inside a Docker container.
#
# Usage:
#   ./run.sh                                    # Run federation test with 100 users
#   ./run.sh test-follower-validation.sh        # Run follower validation test
#   USER_COUNT=50 ./run.sh                      # Run with 50 users
#   KEEP_RUNNING=true ./run.sh                  # Keep servers running after test
#
# Prerequisites:
#   Docker must be installed and running.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(git rev-parse --show-toplevel)"
IMAGE_NAME="owncast-ap-test"

echo "Building Docker image..."
docker build -t "${IMAGE_NAME}" "${SCRIPT_DIR}"

# Collect environment variables to pass through
ENV_ARGS=()
for var in USER_COUNT FOLLOW_DELAY KEEP_RUNNING CI PROXY_PORT SNAC_PORT OWNCAST_PORT OWNCAST2_PORT CLEAR_SHARED_INBOX_PERCENT; do
    if [[ -n "${!var}" ]]; then
        ENV_ARGS+=("-e" "${var}=${!var}")
    fi
done

# Always skip interactive prompts inside the container
ENV_ARGS+=("-e" "CI=true")

# Pass the host user/group so the container can hand any files it writes into
# the bind-mounted repo (notably Owncast's data dir) back to us on exit,
# instead of leaving them owned by root and unremovable.
ENV_ARGS+=("-e" "HOST_UID=$(id -u)" "-e" "HOST_GID=$(id -g)")

# Port-forward when KEEP_RUNNING is set so the user can access the services
EXTRA_ARGS=()
if [[ "${KEEP_RUNNING}" == "true" ]]; then
    OWNCAST_PORT="${OWNCAST_PORT:-8080}"
    OWNCAST2_PORT="${OWNCAST2_PORT:-8081}"
    PROXY_PORT="${PROXY_PORT:-8443}"
    EXTRA_ARGS+=("-p" "${OWNCAST_PORT}:${OWNCAST_PORT}" "-p" "${OWNCAST2_PORT}:${OWNCAST2_PORT}" "-p" "${PROXY_PORT}:${PROXY_PORT}")
fi

echo "Running test in Docker container..."
docker run --rm \
    --add-host owncast.local:127.0.0.1 \
    --add-host owncast2.local:127.0.0.1 \
    --add-host snac.local:127.0.0.1 \
    --add-host indieauth.local:127.0.0.1 \
    -v "${REPO_ROOT}:/owncast" \
    -v owncast-ap-test-gomod:/go/pkg/mod \
    -v owncast-ap-test-gobuild:/root/.cache/go-build \
    "${ENV_ARGS[@]}" \
    "${EXTRA_ARGS[@]}" \
    "${IMAGE_NAME}" \
    "$@"
