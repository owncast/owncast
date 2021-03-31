# Docker build
# Must authenticate first: https://docs.github.com/en/packages/using-github-packages-with-your-projects-ecosystem/configuring-docker-for-use-with-github-packages#authenticating-to-github-packages
DOCKER_IMAGE="owncast"
DATE=$(date +"%Y%m%d")
VERSION="${DATE}-nightly"
GIT_COMMIT=$(git rev-list -1 HEAD)

echo "Building Docker image ${DOCKER_IMAGE}..."

# Change to the root directory of the repository
cd $(git rev-parse --show-toplevel)

# Docker build
docker build --build-arg NAME=docker --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$GIT_COMMIT -t ghcr.io/owncast/${DOCKER_IMAGE}:nightly . -f build/release/Dockerfile-build

# Dockerhub
# You must be authenticated via `docker login` with your Dockerhub credentials first.
# docker push gabekangas/owncast:nightly

docker push ghcr.io/owncast/${DOCKER_IMAGE}:nightly