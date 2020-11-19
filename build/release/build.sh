#!/bin/sh

# Human readable names of binary distributions
DISTRO=(macOS-64bit linux-64bit linux-32bit)
# Operating systems for the respective distributions
OS=(darwin linux linux)
# Architectures for the respective distributions
ARCH=(amd64 amd64 386)

# Version
VERSION=$1
SHOULD_RELEASE=$2

if [[ -z "${VERSION}" ]]; then
  echo "Version must be specified when running build"
  exit
fi

[[ -z "${VERSION}" ]] && VERSION='unknownver' || VERSION="${VERSION}"
GIT_COMMIT=$(git rev-list -1 HEAD)
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
# Change to the root directory of the repository
cd $(git rev-parse --show-toplevel)

echo "Cleaning working directories..."
rm -rf ./webroot/hls/* ./hls/* ./webroot/thumbnail.jpg

echo "Creating version ${VERSION} from commit ${GIT_COMMIT}"

# Create production build of Tailwind CSS
pushd build/javascript >> /dev/null
# Install the tailwind & postcss CLIs
npm install --quiet --no-progress
# Run the tailwind CLI and pipe it to postcss for minification.
# Save it to a temp directory that we will reference below.
NODE_ENV="production" ./node_modules/.bin/tailwind build | ./node_modules/.bin/postcss >  "${TMPDIR}tailwind.min.css"
popd

mkdir -p dist

build() {
  NAME=$1
  OS=$2
  ARCH=$3
  VERSION=$4
  GIT_COMMIT=$5

  echo "Building ${NAME} (${OS}/${ARCH}) release from ${GIT_BRANCH}..."

  mkdir -p dist/${NAME}
  mkdir -p dist/${NAME}/webroot/static
  mkdir -p dist/${NAME}/data

  # Default files
  cp config-default.yaml dist/${NAME}/config.yaml
  cp data/content-example.md dist/${NAME}/data/content.md

  cp -R webroot/ dist/${NAME}/webroot/

  # Copy the production pruned+minified css to the build's directory.
  cp "${TMPDIR}tailwind.min.css" ./dist/${NAME}/webroot/js/web_modules/tailwindcss/dist/tailwind.min.css
  cp -R static/ dist/${NAME}/static
  cp README.md dist/${NAME}

  pushd dist/${NAME} >> /dev/null

  CGO_ENABLED=1 ~/go/bin/xgo --branch ${GIT_BRANCH} -ldflags "-s -w -X main.GitCommit=${GIT_COMMIT} -X main.BuildVersion=${VERSION} -X main.BuildType=${NAME}" -targets "${OS}/${ARCH}" github.com/owncast/owncast
  mv owncast-*-${ARCH} owncast

  zip -r -q -8 ../owncast-$VERSION-$NAME.zip .
  popd >> /dev/null

  rm -rf dist/${NAME}/
}

for i in "${!DISTRO[@]}"; do
  build ${DISTRO[$i]} ${OS[$i]} ${ARCH[$i]} $VERSION $GIT_COMMIT
done

# Use the second argument "release" to create an actual release.
if [ "$SHOULD_RELEASE" != "release" ]; then
  echo "Not creating a release."
  exit
fi

# Create the tag
git tag -a "v${VERSION}" -m "Release build v${VERSION}"

# On macOS open the Github page for new releases so they can be uploaded
if test -f "/usr/bin/open"; then
  open "https://github.com/owncast/owncast/releases/new"
  open dist
fi

# Docker build
# Must authenticate first: https://docs.github.com/en/packages/using-github-packages-with-your-projects-ecosystem/configuring-docker-for-use-with-github-packages#authenticating-to-github-packages
DOCKER_IMAGE="owncast-${VERSION}"
echo "Building Docker image ${DOCKER_IMAGE}..."

# Change to the root directory of the repository
cd $(git rev-parse --show-toplevel)

# Github Packages
docker build --build-arg NAME=docker --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$GIT_COMMIT -t owncast . -f scripts/Dockerfile-build
docker tag $DOCKER_IMAGE docker.pkg.github.com/owncast/owncast/$DOCKER_IMAGE:$VERSION
docker push docker.pkg.github.com/owncast/owncast/$DOCKER_IMAGE:$VERSION
#
# Dockerhub
# You must be authenticated via `docker login` with your Dockerhub credentials first.
docker tag owncast gabekangas/owncast:$VERSION
docker push gabekangas/owncast
