VERSION 0.6

FROM alpine:latest
ARG version=develop
WORKDIR /build

crosscompiler:
  # This image is missing a few platforms, so we'll add them locally
  FROM bdwyertech/go-crosscompile
  RUN curl -sfL "https://musl.cc/armv7l-linux-musleabihf-cross.tgz" | tar zxf - -C /usr/ --strip-components=1
  RUN curl -sfL "https://musl.cc/i686-linux-musl-cross.tgz" | tar zxf - -C /usr/ --strip-components=1
  RUN curl -sfL "https://musl.cc/x86_64-linux-musl-cross.tgz" | tar zxf - -C /usr/ --strip-components=1

code:
  FROM +crosscompiler
  # COPY . /build
  GIT CLONE --branch=$version git@github.com:owncast/owncast.git /build

build:
  FROM +code
  RUN echo $EARTHLY_GIT_HASH
  ARG NAME
  ARG GOOS
  ARG GOARCH
  ARG GOARM
  ARG CC=clang
  ARG CXX=clang++
  ENV CGO_ENABLED=1
  ENV GOOS=$GOOS
  ENV GOARCH=$GOARCH
  ENV GOARM=$GOARM
  ENV CC=$CC
  ENV CXX=$CXX

  WORKDIR /build
  # MacOSX disallows static executables, so we omit the static flag on this platform
  RUN go build -a -installsuffix cgo -ldflags "$([ $GOOS != darwin ] && echo "-linkmode external -extldflags -static ") -s -w -X github.com/owncast/owncast/config.GitCommit=$EARTHLY_GIT_HASH -X github.com/owncast/owncast/config.VersionNumber=$version -X github.com/owncast/owncast/config.BuildPlatform=$NAME" -o owncast main.go
  SAVE ARTIFACT owncast owncast
  SAVE ARTIFACT webroot webroot
  SAVE ARTIFACT README.md README.md

tailwind:
  FROM +code
  WORKDIR /build/build/javascript
  RUN apk add --update --no-cache npm >> /dev/null
  ENV NODE_ENV=production
  RUN cd /build/build/javascript && npm install --quiet --no-progress >> /dev/null && npm install -g cssnano postcss postcss-cli --quiet --no-progress --save-dev >> /dev/null && ./node_modules/.bin/tailwind build > /build/tailwind.min.css
  RUN npx postcss /build/tailwind.min.css > /build/prod-tailwind.min.css
  SAVE ARTIFACT /build/prod-tailwind.min.css

package:
  RUN apk add --update --no-cache zip >> /dev/null

  ARG NAME

  COPY +build/webroot /build/dist/webroot
  COPY +build/owncast /build/dist/owncast
  COPY +tailwind/prod-tailwind.min.css /build/dist/webroot/js/web_modules/tailwindcss/dist/tailwind.min.css
  COPY +build/README.md /build/dist/README.md
  ENV ZIPNAME owncast-$version-$NAME.zip
  RUN cd /build/dist && zip -r -q -8 /build/dist/owncast.zip .
  SAVE ARTIFACT /build/dist/owncast.zip owncast.zip AS LOCAL dist/$ZIPNAME

build-all:
  BUILD +linux-arm64
  BUILD +linux-arm7
  BUILD +linux-amd64
  BUILD +linux-i386
  BUILD +darwin-amd64

linux-arm64:
  ARG NAME=linux-arm64
  BUILD +package --NAME=$NAME --GOOS=linux --GOARCH=arm64 --CC=aarch64-linux-musl-gcc --CXX=aarch64-linux-musl-g++

linux-arm7:
  ARG NAME=linux-arm7
  BUILD +package --NAME=$NAME --GOOS=linux --GOARCH=arm --GOARM=7 --CC=armv7l-linux-musleabihf-gcc --CXX=armv7l-linux-musleabihf-g++

linux-amd64:
  ARG NAME=linux-64bit
  BUILD +package --NAME=$NAME --GOOS=linux --GOARCH=amd64 --CC=x86_64-linux-musl-gcc --CXX=x86_64-linux-musl-g++

linux-i386:
  ARG NAME=linux-32bit
  BUILD +package --NAME=$NAME --GOOS=linux --GOARCH=386 --CC=i686-linux-musl-gcc --CXX=i686-linux-musl-g++

darwin-amd64:
  ARG NAME=macOS-64bit
  BUILD +package --NAME=$NAME --GOOS=darwin --GOARCH=amd64 --CC=o64-clang --CXX=o64-clang++
