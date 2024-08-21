#!/bin/bash

set -e

function install_ffmpeg() {
	# install a specific version of ffmpeg

	if [[ "$OSTYPE" == "linux-"* ]]; then
		OS="linux"
	elif [[ "$OSTYPE" == "darwin"* ]]; then
		OS="macos"
	else
		echo "Exiting!!!"
		exit 1
	fi

	OS="linux"

	FFMPEG_VER="5.1"
	FFMPEG_PATH="$(pwd)/ffmpeg-$FFMPEG_VER"
	PATH=$FFMPEG_PATH:$PATH

	if ! [[ -d "$FFMPEG_PATH" ]]; then
		mkdir "$FFMPEG_PATH"
	fi

	pushd "$FFMPEG_PATH" >/dev/null

	if [[ -x "$FFMPEG_PATH/ffmpeg" ]]; then

		ffmpeg_version=$("$FFMPEG_PATH/ffmpeg" -version | awk -F 'ffmpeg version' '{print $2}' | awk 'NR==1{print $1}')

		if [[ "$ffmpeg_version" == "$FFMPEG_VER-static" ]]; then
			popd >/dev/null
			return 0
		else
			mv "$FFMPEG_PATH/ffmpeg" "$FFMPEG_PATH/ffmpeg.bk" || rm -f "$FFMPEG_PATH/ffmpeg"
		fi
	fi

	rm -f ffmpeg.zip
	curl -sL --fail https://github.com/ffbinaries/ffbinaries-prebuilt/releases/download/v${FFMPEG_VER}/ffmpeg-${FFMPEG_VER}-${OS}-64.zip --output ffmpeg.zip >/dev/null
	unzip -o ffmpeg.zip >/dev/null && rm -f ffmpeg.zip
	chmod +x ffmpeg
	PATH=$FFMPEG_PATH:$PATH

	popd >/dev/null
}

function start_owncast() {
	# Build and run owncast from source
	echo "Building owncast..."
	pushd "$(git rev-parse --show-toplevel)" >/dev/null
	CGO_ENABLED=1 go build -o owncast main.go

	echo "Running owncast..."
	./owncast -database "$TEMP_DB" &
	SERVER_PID=$!
	popd >/dev/null

	sleep 5

}

function start_stream() {
	# Start streaming the test file over RTMP to the local owncast instance.
	../../ocTestStream.sh &
	STREAM_PID=$!

	echo "Waiting for stream to start..."
	sleep 12
}

function update_storage_config() {
	echo "Configuring external storage to use ${S3_BUCKET}..."

	# Hard-coded to admin:abc123 for auth
	curl --fail 'http://localhost:8080/api/admin/config/s3' \
		-H 'Authorization: Basic YWRtaW46YWJjMTIz' \
		--data-raw "{\"value\":{\"accessKey\":\"${S3_ACCESS_KEY}\",\"acl\":\"\",\"bucket\":\"${S3_BUCKET}\",\"enabled\":true,\"endpoint\":\"${S3_ENDPOINT}\",\"region\":\"${S3_REGION}\",\"secret\":\"${S3_SECRET}\",\"servingEndpoint\":\"\"}}"
}

function kill_with_kids() {
	# kill a process and all its children (by pid)! return no error.

	if [[ -n $1 ]]; then
		mapfile -t CHILDREN_PID_LIST < <(ps --ppid "$1" -o pid= &>/dev/null || true)
		for child_pid in "${CHILDREN_PID_LIST[@]}"; do
			kill "$child_pid" &>/dev/null || true
			wait "$child_pid" &>/dev/null || true
		done
		kill "$1" &>/dev/null || true
		wait "$1" &>/dev/null || true
	fi
}

function finish() {
	echo "Cleaning up..."
	kill_with_kids "$STREAM_PID"
	kill "$SERVER_PID" &>/dev/null || true
	wait "$SERVER_PID" &>/dev/null || true
	rm -fr "$TEMP_DB"
}

trap finish EXIT

TEMP_DB=$(mktemp)
