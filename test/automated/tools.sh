#!/bin/bash

set -e

function install_ffmpeg() {
	# install a specific version of ffmpeg

	FFMPEG_VER="8.0"
	FFMPEG_BUILD_VERSION="20260223192056"
	FFMPEG_PATH="$(pwd)"
	PATH=$FFMPEG_PATH:$PATH

	case "$OSTYPE" in
		linux-*) ffmpeg_os="linux" ;;
		darwin*) ffmpeg_os="darwin" ;;
		*)
			echo "Unsupported platform: $OSTYPE"
			exit 1
			;;
	esac

	case "$(uname -m)" in
		x86_64 | amd64) ffmpeg_arch="amd64" ;;
		aarch64 | arm64) ffmpeg_arch="arm64" ;;
		*)
			echo "Unsupported architecture: $(uname -m)"
			exit 1
			;;
	esac

	if [[ "$ffmpeg_os" == "linux" ]]; then
		ffmpeg_asset="ffmpeg${FFMPEG_VER}-linux-${ffmpeg_arch}-static.tar.gz"
	else
		ffmpeg_asset="ffmpeg${FFMPEG_VER}-darwin-${ffmpeg_arch}.tar.gz"
	fi

	if [[ -x "$FFMPEG_PATH/ffmpeg" ]]; then

		ffmpeg_version=$("$FFMPEG_PATH/ffmpeg" -version | awk -F 'ffmpeg version' '{print $2}' | awk 'NR==1{print $1}')

		# Linux static builds report "8.0-static"; macOS builds report "8.0" or "8.0.1".
		if [[ "$ffmpeg_version" == "$FFMPEG_VER" || "$ffmpeg_version" == "$FFMPEG_VER".* || "$ffmpeg_version" == "$FFMPEG_VER"-* ]]; then
			return 0
		else
			mv "$FFMPEG_PATH/ffmpeg" "$FFMPEG_PATH/ffmpeg.bk" || rm -f "$FFMPEG_PATH/ffmpeg"
		fi
	fi

	echo "Downloading ffmpeg v${FFMPEG_VER} release ${FFMPEG_BUILD_VERSION} for ${ffmpeg_os}/${ffmpeg_arch}"
	rm -rf ffmpeg.tar.gz
	curl -sL --fail "https://github.com/owncast/ffmpeg-builds/releases/download/${FFMPEG_BUILD_VERSION}/${ffmpeg_asset}" --output ffmpeg.tar.gz >/dev/null
	tar -xzf ffmpeg.tar.gz
	rm -f ffmpeg.tar.gz
	chmod +x ffmpeg
	PATH=$FFMPEG_PATH:$PATH
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
		while IFS= read -r child_pid; do
			[[ -n "$child_pid" ]] || continue
			kill "$child_pid" &>/dev/null || true
			wait "$child_pid" &>/dev/null || true
		done < <(pgrep -P "$1" 2>/dev/null || true)
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
