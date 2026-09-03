#!/usr/bin/env bash
# Build and run every plugin SDK example against this Owncast checkout.
#
# PLUGIN_SDK_DIR may point at a local SDK checkout. The default is the sibling
# checkout used by local Owncast development.
set -euo pipefail

CORE_DIR="$(git rev-parse --show-toplevel)"
SDK_DIR="${PLUGIN_SDK_DIR:-${CORE_DIR%/owncast}/owncast-plugin-sdk}"
HOST_BIN="${OWNCAST_PLUGIN_HOST_BIN:-${SDK_DIR}/.compatibility/owncast-plugin-test}"

if [[ ! -f "${SDK_DIR}/host-runtime/go.mod" ]]; then
	echo "plugin SDK checkout not found: ${SDK_DIR}" >&2
	exit 1
fi

mkdir -p "$(dirname "$HOST_BIN")"
HOST_WORK_DIR="$(mktemp -d)"
cleanup() {
	rm -rf "$HOST_WORK_DIR"
}
trap cleanup EXIT

# Use a workspace so the host binary imports this exact core checkout instead
# of the version pinned in the SDK's host-runtime/go.mod.
(
	cd "$HOST_WORK_DIR"
	go work init "$CORE_DIR" "${SDK_DIR}/host-runtime"
)
(
	cd "${SDK_DIR}/host-runtime"
	GOWORK="${HOST_WORK_DIR}/go.work" go build -o "$HOST_BIN" ./cmd/owncast-plugin-test
)

run_host_check() {
	local project="$1"
	if [[ -d "${project%/}/__tests__" ]]; then
		"$HOST_BIN" "$project"
	else
		"$HOST_BIN" --load-only "$project"
	fi
}

run_js() {
	local project slug
	(
		cd "${SDK_DIR}/sdks/js"
		if [[ -f package-lock.json ]]; then
			npm ci --no-audit --no-fund --ignore-scripts
		else
			npm install --no-audit --no-fund --ignore-scripts
		fi
	)
	for project in "${SDK_DIR}"/examples/js/*/; do
		[[ -f "$project/plugin.manifest.json" ]] || continue
		[[ -d "$project" ]] || continue
		slug="$(basename "$project")"
		echo "[js] ${slug}"
		(
			cd "$project"
			if [[ -f package-lock.json ]]; then
				npm ci --no-audit --no-fund --ignore-scripts
			else
				npm install --no-audit --no-fund --ignore-scripts
			fi
			npx --no-install owncast-plugin build
			run_host_check .
		)
	done
}

run_python() {
	local project slug
	for project in "${SDK_DIR}"/examples/python/*/; do
		[[ -f "$project/plugin.manifest.json" ]] || continue
		[[ -d "$project" ]] || continue
		slug="$(basename "$project")"
		echo "[python] ${slug}"
		python3 "${SDK_DIR}/sdks/python/owncast_plugin_build.py" "$project"
		run_host_check "$project"
	done
}

case "${PLUGIN_LANGUAGE:-all}" in
javascript)
	run_js
	;;
python)
	run_python
	;;
all)
	run_js
	run_python
	;;
*)
	echo "PLUGIN_LANGUAGE must be javascript, python, or all" >&2
	exit 2
	;;
esac

echo "plugin SDK compatibility checks passed"
