# Docker build
# Must authenticate first: https://docs.github.com/en/packages/using-github-packages-with-your-projects-ecosystem/configuring-docker-for-use-with-github-packages#authenticating-to-github-packages
DOCKER_IMAGE="owncast"
DATE=$(date +"%Y%m%d")
VERSION="${DATE}-nightly"
GIT_COMMIT=$(git rev-list -1 HEAD)

# Create production build of Tailwind CSS
pushd ../../build/javascript >> /dev/null
# Install the tailwind & postcss CLIs
npm install --quiet --no-progress
# Run the tailwind CLI and pipe it to postcss for minification.
# Save it to a temp directory that we will reference below.
NODE_ENV="production" ./node_modules/.bin/tailwind build | ./node_modules/.bin/postcss >  "../../webroot/js/web_modules/tailwindcss/dist/tailwind.min.css"
popd

echo "Building Docker image ${DOCKER_IMAGE}..."

# Change to the root directory of the repository
cd $(git rev-parse --show-toplevel)

# Docker build
docker build --build-arg NAME=docker --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$GIT_COMMIT -t ghcr.io/owncast/${DOCKER_IMAGE}:nightly .

# Dockerhub
# You must be authenticated via `docker login` with your Dockerhub credentials first.
# docker push gabekangas/owncast:nightly

docker push ghcr.io/owncast/${DOCKER_IMAGE}:nightly