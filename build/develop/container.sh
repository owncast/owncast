#!/bin/sh
set -e

# Development container builder
#
# Must authenticate first: https://docs.github.com/en/packages/using-github-packages-with-your-projects-ecosystem/configuring-docker-for-use-with-github-packages#authenticating-to-github-packages
# env vars:
#   $EARTHLY_BUILD_BRANCH: git branch to checkout
#   $EARTHLY_BUILD_TAG: tag for container image

EARTHLY_IMAGE_NAME="owncast"
BUILD_TAG=${EARTHLY_BUILD_TAG:-develop}
DATE=$(date +"%Y%m%d")
VERSION="${DATE}-${BUILD_TAG}"

echo "Building container image ${EARTHLY_IMAGE_NAME}:${BUILD_TAG} ..."

# Change to the root directory of the repository
cd "$(git rev-parse --show-toplevel)" || exit
if [ -n "${EARTHLY_BUILD_BRANCH}" ]; then
	git checkout "${EARTHLY_BUILD_BRANCH}" || exit
fi

earthly --ci +docker-all --images="ghcr.io/owncast/${EARTHLY_IMAGE_NAME}:${BUILD_TAG}" --version="${VERSION}"
