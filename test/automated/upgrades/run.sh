#!/bin/bash

# This script will download all the releases of Owncast and test them
# to ensure that upgrades work as expected.  It will also test the
# development branch as the final test.
# It is hard coded to run under 64bit intel linux.

# set -o errexit
set -o pipefail

releases=(
	"https://github.com/owncast/owncast/releases/download/v0.0.6/owncast-0.0.6-linux-64bit.zip"
	"https://github.com/owncast/owncast/releases/download/v0.0.7/owncast-0.0.7-linux-64bit.zip"
	"https://github.com/owncast/owncast/releases/download/v0.0.8/owncast-0.0.8-linux-64bit.zip"
	"https://github.com/owncast/owncast/releases/download/v0.0.9/owncast-0.0.9-linux-64bit.zip"
	"https://github.com/owncast/owncast/releases/download/v0.0.10/owncast-0.0.10-linux-64bit.zip"
	"https://github.com/owncast/owncast/releases/download/v0.0.11/owncast-0.0.11-linux-64bit.zip"
	"https://github.com/owncast/owncast/releases/download/v0.0.12/owncast-0.0.12-linux-64bit.zip"
	"https://github.com/owncast/owncast/releases/download/v0.0.13/owncast-0.0.13-linux-64bit.zip"
)

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
	git clone https://github.com/owncast/owncast
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
