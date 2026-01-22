#!/bin/bash

# Run the ActivityPub federation test
#
# Usage:
#   ./run.sh                    # Run with 100 users
#   USER_COUNT=50 ./run.sh      # Run with 50 users
#   KEEP_RUNNING=true ./run.sh  # Keep servers running after test
#
# Prerequisites:
#   Add to /etc/hosts: 127.0.0.1 owncast.local snac.local

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exec "${SCRIPT_DIR}/test-federation.sh" "$@"
