#!/bin/sh

# Human readable names of binary distributions
DISTRO=(macOS linux)
# Operating systems for the respective distributions
OS=(darwin linux)
# Architectures for the respective distributions
ARCH=(amd64 amd64)

# Version
VERSION=$1

if [[ -z "${VERSION}" ]]; then
  echo "Version must be specified when running build"
fi

[[ -z "${VERSION}" ]] && VERSION='unknownver' || VERSION="${VERSION}"
GIT_COMMIT=$(git rev-list -1 HEAD)

# Change to the root directory of the repository
cd $(git rev-parse --show-toplevel)

echo "Cleaning working directories..."
rm -rf ./webroot/hls/* ./hls/* ./webroot/thumbnail.jpg

echo "Creating version ${VERSION} from commit ${GIT_COMMIT}"

mkdir -p dist

build() {
  NAME=$1
  OS=$2
  ARCH=$3
  VERSION=$4
  GIT_COMMIT=$5

  echo "Building ${NAME} (${OS}/${ARCH}) release..."

  mkdir -p dist/${NAME}

  # Default files
  cp config-example.yaml dist/${NAME}/config.yaml
  cp webroot/static/content-example.md dist/${NAME}/webroot/static/content.md

  cp -R webroot/ dist/${NAME}/webroot/
  cp -R doc/ dist/${NAME}/doc/
  cp README.md dist/${NAME}

  env CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -ldflags "-X core.GitCommit=$GIT_COMMIT -X core.BuildVersion=$VERSION -X core.BuildType=$NAME" -a -o dist/$NAME/owncast

  pushd dist/${NAME} >> /dev/null
  zip -r -q -8 ../owncast-$NAME-$VERSION.zip .
  popd >> /dev/null

  rm -rf dist/${NAME}/
}

for i in "${!DISTRO[@]}"; do
  build ${DISTRO[$i]} ${OS[$i]} ${ARCH[$i]} $VERSION $GIT_COMMIT
done

# Create the tag
git tag -a "v${VERSION}" -m "Release build v${VERSION}"

# On macOS open the Github page for new releases so they can be uploaded
if test -f "/usr/bin/open"; then
  open "https://github.com/gabek/owncast/releases/new"
  open dist
fi
