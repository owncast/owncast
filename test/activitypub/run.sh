#!/bin/sh

# configure a tunnel for snac using localhost.run
start_tunnel() {
	TUNNEL_PIPE="/tmp/localhostrun"
	mkfifo "$TUNNEL_PIPE"

	echo "Using username $USERNAME"
	ssh -T "$USERNAME"@srv.us -R 1:localhost:8001 >"$TUNNEL_PIPE" &
	while read line; do
		echo $line
		DOMAIN=$(echo $line | grep -o 'https://[^ ]*' | tail -n 1)
		if [ -n "$DOMAIN" ]; then
			break
		fi
	done <"$TUNNEL_PIPE"
	echo "Found domain: $DOMAIN"
}

# set SNAC_DATA_PATH externally, must be shared across scripts
configure_snac() {
	echo -e "\n\n$DOMAIN\n\n" | snac init "$SNAC_DATA_PATH"
}

run_snac() {
	snac httpd "$SNAC_DATA_PATH"
}

if [ -z "$SNAC_DATA_PATH" ]; then
	echo "SNAC_DATA_PATH not set"
	exit 1
fi

start_tunnel

if [ -z "$DOMAIN" ]; then
	echo "tunnel failed to get domain"
	exit 1
fi
echo $DOMAIN
configure_snac
run_snac # daemon
