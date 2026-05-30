#!/bin/bash

# This script will download all the releases of Owncast and test them
# to ensure that upgrades work as expected.  It will also test the
# development branch as the final test.
# The release list is fetched dynamically from the GitHub API so it stays
# current as new versions are published.
# It is hard coded to run under 64bit intel linux, and requires curl and jq.

# set -o errexit
set -o pipefail

if ! command -v jq >/dev/null 2>&1; then
	echo "Error: 'jq' is required to fetch the list of releases but was not found." >&2
	exit 1
fi

# Fetch the list of linux 64-bit release download URLs from the GitHub API,
# sorted oldest to newest. This keeps the test against the full set of
# published releases without hard-coding versions that go out of date.
fetch_release_urls() {
	curl -sL -H "Accept: application/vnd.github+json" \
		"https://api.github.com/repos/owncast/owncast/releases?per_page=100" |
		jq -r '.[]
			| select(.prerelease == false and .draft == false)
			| .assets[]
			| select(.name | test("linux-64bit\\.zip$"))
			| .browser_download_url' |
		sort -V
}

echo "Fetching the list of Owncast releases from GitHub..."
mapfile -t releases < <(fetch_release_urls)

if [ ${#releases[@]} -eq 0 ]; then
	echo "Error: failed to fetch any releases from the GitHub API." >&2
	exit 1
fi

echo "--------------------------------------------"
echo "Owncast releases upgrade test."
echo "Will download ${#releases[@]} releases plus the development branch."
echo "Please wait, as this will take a while."
printf "\n"

rm -rf releases
rm -rf owncast
rm -rf src

mkdir -p releases
mkdir -p src

download_release() {
	url=$1

	echo "--------------------------------------------"
	echo "Downloading $url"

	zipfile="releases/$(basename "$url")"
	curl -sL "${url}" --output "${zipfile}"
}

test_release() {
	pushd ./owncast >>/dev/null || exit
	timeout --preserve-status 10 ./owncast
	popd >>/dev/null || exit
}

build_development() {
	echo "Building test release from current development branch..."
	cd src || exit
	git clone --depth 1 https://github.com/owncast/owncast
	cd owncast || exit
	earthly +package --platform="linux/amd64"
	mv dist/owncast-develop-linux-64bit.zip ../../releases/owncast-develop-linux-64bit.zip
	cd ../..
}

unzip_release() {
	zipfile="releases/$(basename "$1")"
	unzip -o "${zipfile}" -d "owncast" >>/dev/null
}

# Test all the releases in a row
for release in "${releases[@]}"; do
	if [ ! -f "releases/$(basename "$1")" ]; then
		download_release "$release"
	fi

	unzip_release "$(basename "$release")"
	test_release
done

# # Build and run the latest release
build_development
unzip_release owncast-develop-linux-64bit.zip
test_release

# Test jumping from the first release to the development release
rm -rf owncast

if [ ! -f "releases/$(basename "${releases[0]}")" ]; then
	download_release "${releases[0]}"
fi
unzip_release "$(basename "${releases[0]}")"
test_release
echo "--------------------------------------------"
echo "Testing upgrade from the first release to the development branch."
printf "\n"

if [ ! -f "releases/owncast-develop-linux-64bit.zip" ]; then
	build_development
fi
unzip_release owncast-develop-linux-64bit.zip
test_release

echo "Done."
