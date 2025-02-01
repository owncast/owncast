#!/bin/bash

# This script will make your local Owncast server available at a public URL.
# It's particularly useful for testing on mobile devices or want to test
# activitypub integration.
# Pass a custom domain as an argument if you have previously set it up at
# localhost.run. Otherwise, a random hostname will be generated.
# DOMAIN=me.example.com ./test/test-local.sh
# Pass a port number as an argument if you are running Owncast on a different port.
# By default, it will use port 8080.
# PORT=8080 ./test/test-local.sh

PORT="${PORT:=8080}"
HOST="localhost"

# Using nc (netcat) to check if the port is open
if nc -zv $HOST $PORT 2>/dev/null; then
	echo "Your web server is running on port $PORT. Good."
else
	echo "Please make sure your Owncast server is running on port $PORT."
	exit 1
fi

if [ -n "$DOMAIN" ]; then
	echo "Attempting to use custom domain: $DOMAIN"
	ssh -R $DOMAIN:80:localhost:$PORT localhost.run

else
	echo "Using auto-generated hostname for tunnel."
	ssh -R 80:localhost:$PORT localhost.run
fi
