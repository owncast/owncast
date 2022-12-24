#!/bin/sh

# Docker build
# Must authenticate first: https://docs.github.com/en/packages/using-github-packages-with-your-projects-ecosystem/configuring-docker-for-use-with-github-packages#authenticating-to-github-packages
EARTHLY_IMAGE_NAME="owncast"
BUILD_TAG=${EARTHLY_BUILD_TAG:-webv2}
DATE=$(date +"%Y%m%d")
VERSION="${DATE}-${BUILD_TAG}"


echo "Building Docker image ${EARTHLY_IMAGE_NAME}:${BUILD_TAG}..."

# Change to the root directory of the repository
cd "$(git rev-parse --show-toplevel)" || exit
git checkout "${EARTHLY_BUILD_BRANCH:-webv2}"

earthly --ci --push +docker-all --image="ghcr.io/owncast/${EARTHLY_IMAGE_NAME}" --tag=${BUILD_TAG} --version="${VERSION}"
