#!/bin/bash

# This script will make your local Owncast server available at a public URL.
# It's particularly useful for testing on mobile devices.

PORT=8080
USERNAME=$(whoami)
BRANCH=$(git symbolic-ref --short -q HEAD | sed 's/[^a-z A-Z 0-9]//g')
HOSTNAME="oc-$USERNAME-$BRANCH"

echo "Please wait. Making the local server port $PORT available at https://$HOSTNAME.serveo.net"
ssh -R "$HOSTNAME".serveo.net:80:localhost:$PORT serveo.net
