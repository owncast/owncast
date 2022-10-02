#!/bin/sh

# Docker build
# Must authenticate first: https://docs.github.com/en/packages/using-github-packages-with-your-projects-ecosystem/configuring-docker-for-use-with-github-packages#authenticating-to-github-packages
DOCKER_IMAGE="owncast"
DATE=$(date +"%Y%m%d")
TAG="webv2"
VERSION="${DATE}-${TAG}"
echo "Building Docker image ${DOCKER_IMAGE}..."

# Change to the root directory of the repository
cd $(git rev-parse --show-toplevel)
git checkout webv2

earthly --ci --push +docker-all --image="ghcr.io/owncast/${DOCKER_IMAGE}" --tag=${TAG} --version="${VERSION}"
