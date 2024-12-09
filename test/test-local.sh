#!/bin/bash

# This script will make your local Owncast server available at a public URL.
# It's particularly useful for testing on mobile devices or want to test
# activitypub integration.

PORT=8080
HOST="localhost"

# Using nc (netcat) to check if the port is open
if nc -zv $HOST $PORT 2>/dev/null; then
	echo "Your web server is running on port $PORT. Good."
else
	echo "Please make sure your Owncast server is running on port $PORT."
	exit 1
fi

ssh -R 80:localhost:$PORT localhost.run
