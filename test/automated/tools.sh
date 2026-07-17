#!/bin/bash

set -e

# When a GitHub token is available (e.g. GITHUB_TOKEN in CI), send it on
# requests to github.com so release-asset downloads use the authenticated
# 5,000/hour rate limit instead of the anonymous, shared-per-IP 60/hour limit
# that frequently fails on CI runners. Local runs without a token are
# unaffected and keep downloading anonymously. curl drops the Authorization
# header on the cross-host redirect to the asset CDN, so the token is not
# leaked to objects.githubusercontent.com.
gh_curl_auth=()
if [[ -n "${GITHUB_TOKEN:-}" ]]; then
	gh_curl_auth=(-H "Authorization: Bearer ${GITHUB_TOKEN}")
fi

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
	curl -sL --fail "${gh_curl_auth[@]}" "https://github.com/owncast/ffmpeg-builds/releases/download/${FFMPEG_BUILD_VERSION}/${ffmpeg_asset}" --output ffmpeg.tar.gz >/dev/null
	tar -xzf ffmpeg.tar.gz
	rm -f ffmpeg.tar.gz
	chmod +x ffmpeg
	PATH=$FFMPEG_PATH:$PATH
}

function start_owncast() {
	# Refuse to run if something already answers on the test port. The
	# readiness check below polls http://localhost:8080, so a pre-existing
	# server (most commonly a `go run main.go` dev backend) would answer it
	# and the whole suite would silently run against that instance and its
	# database instead of the throwaway one started here.
	if curl -fsS -o /dev/null "http://localhost:8080/api/status" 2>/dev/null; then
		echo "ERROR: something is already serving http://localhost:8080." >&2
		echo "Stop it before running this suite (it would be tested, and its database modified, instead of the throwaway test instance)." >&2
		exit 1
	fi

	# Build and run owncast from source. -cover instruments the binary so e2e
	# runs produce Go coverage data in GOCOVERDIR (reported at teardown).
	echo "Building owncast..."
	pushd "$REPO_ROOT" >/dev/null
	CGO_ENABLED=1 go build -cover -o owncast main.go
	# Remove any log left by a previous run so the artifact collected at
	# teardown belongs to this run.
	rm -f data/logs/owncast.log
	echo "Running owncast..."
	GOCOVERDIR="$GOCOVERDIR" ./owncast -database "$TEMP_DB" &
	SERVER_PID=$!
	popd >/dev/null

	# Wait for the API to actually come up instead of a blind sleep. The server
	# runs in the background, so `set -e` can't see it crash; without this, a
	# failed start (most commonly port 8080/1935 already in use by another
	# Owncast) leaves the suite running against a stale instance and producing
	# confusing, unrelated test failures. Fail loudly and early instead.
	echo "Waiting for owncast to be ready on :8080..."
	for _ in $(seq 1 30); do
		if ! kill -0 "$SERVER_PID" 2>/dev/null; then
			echo "ERROR: owncast exited during startup. Is port 8080 or 1935 already in use by another instance?" >&2
			exit 1
		fi
		if curl -fsS -o /dev/null "http://localhost:8080/api/status" 2>/dev/null; then
			return 0
		fi
		sleep 1
	done
	echo "ERROR: owncast did not become ready on :8080 within 30s." >&2
	exit 1
}

function start_stream() {
	# Start streaming the test file over RTMP to the local owncast instance.
	../../ocTestStream.sh &
	STREAM_PID=$!

	# Wait for the stream to actually go live instead of a blind sleep: the
	# online player tests assert on UI that only renders once the API reports
	# the stream online.
	echo "Waiting for stream to go live..."
	for _ in $(seq 1 30); do
		if curl -fsS "http://localhost:8080/api/status" 2>/dev/null | grep -q '"online":true'; then
			return 0
		fi
		sleep 1
	done
	echo "ERROR: stream did not go live within 30s." >&2
	exit 1
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

function stop_owncast() {
	# Coverage counter data is only written on a clean exit (return from main
	# or os.Exit); a plain signal death loses it because owncast installs no
	# signal handler. Ask the server to exit through the authenticated
	# forcequit admin endpoint (which calls os.Exit(0)), then fall back to
	# SIGTERM if the request or the wait fails.
	if [[ -n "${SERVER_PID:-}" ]] && kill -0 "$SERVER_PID" 2>/dev/null; then
		curl -m 5 -s -o /dev/null -u admin:abc123 \
			"http://localhost:8080/api/admin/update/forcequit" || true
		for _ in $(seq 1 20); do
			kill -0 "$SERVER_PID" 2>/dev/null || break
			sleep 0.5
		done
	fi
	kill "${SERVER_PID:-}" &>/dev/null || true
	wait "${SERVER_PID:-}" &>/dev/null || true
}

function collect_artifacts() {
	# Preserve the server log and the throwaway database so CI can upload
	# them on failure. Cheap and unconditional.
	mkdir -p "$SUITE_DIR/test-artifacts"
	cp "$REPO_ROOT/data/logs/owncast.log" "$SUITE_DIR/test-artifacts/" 2>/dev/null || true
	cp "$REPO_ROOT/data/logs/transcoder.log" "$SUITE_DIR/test-artifacts/" 2>/dev/null || true
	# Owncast runs sqlite in WAL mode and forcequit exits without a clean
	# close, so most data lives in the -wal sidecar. Copy it (and -shm) with
	# matching names so the artifact opens as a normal sqlite database.
	cp "$TEMP_DB" "$SUITE_DIR/test-artifacts/owncast.db" 2>/dev/null || true
	cp "$TEMP_DB-wal" "$SUITE_DIR/test-artifacts/owncast.db-wal" 2>/dev/null || true
	cp "$TEMP_DB-shm" "$SUITE_DIR/test-artifacts/owncast.db-shm" 2>/dev/null || true
}

function report_coverage() {
	if ! compgen -G "$GOCOVERDIR/covcounters.*" >/dev/null; then
		echo "WARNING: no coverage counter data in GOCOVERDIR; owncast likely did not exit cleanly. Skipping coverage report." >&2
		return 0
	fi
	mkdir -p "$SUITE_DIR/e2e-coverage"
	go tool covdata percent -i "$GOCOVERDIR" >"$SUITE_DIR/e2e-coverage/percent.txt" || true
	go tool covdata textfmt -i "$GOCOVERDIR" -o "$SUITE_DIR/e2e-coverage/profile.txt" || true
	local total
	total=$(go tool cover -func="$SUITE_DIR/e2e-coverage/profile.txt" 2>/dev/null | tail -n 1) || true
	if [[ -n "$total" ]]; then
		echo "$total" >>"$SUITE_DIR/e2e-coverage/percent.txt"
		echo "E2E Go coverage $total"
	fi
}

function finish() {
	echo "Cleaning up..."
	kill_with_kids "${STREAM_PID:-}"
	stop_owncast
	collect_artifacts
	report_coverage
	rm -f "$TEMP_DB" "$TEMP_DB-wal" "$TEMP_DB-shm"
	rm -fr "$GOCOVERDIR"
}

trap finish EXIT

# tools.sh is sourced from the suite's own directory; remember it so teardown
# can drop artifacts next to the suite regardless of the cwd at exit time.
SUITE_DIR="$(pwd)"
REPO_ROOT="$(git rev-parse --show-toplevel)"
TEMP_DB=$(mktemp)
GOCOVERDIR=$(mktemp -d)
