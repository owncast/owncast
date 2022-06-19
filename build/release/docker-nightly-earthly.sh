#!/bin/sh

# Docker build
# Must authenticate first: https://docs.github.com/en/packages/using-github-packages-with-your-projects-ecosystem/configuring-docker-for-use-with-github-packages#authenticating-to-github-packages
DOCKER_IMAGE="owncast-earthly"
DATE=$(date +"%Y%m%d")
VERSION="${DATE}-nightly"

echo "Building Docker image ${DOCKER_IMAGE}..."

earthly --ci --push +docker --image="ghcr.io/owncast/${DOCKER_IMAGE}" --tag=nightly --version=${VERSION}

# docker push ghcr.io/owncast/${DOCKER_IMAGE}:nightly
