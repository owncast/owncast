VERSION --new-platform 0.6

FROM --platform=linux/amd64 alpine:latest
ARG version=develop

WORKDIR /build

build-all:
  BUILD --platform=linux/amd64 --platform=linux/386 --platform=linux/arm64 --platform=linux/arm/v7 --platform=darwin/amd64 +build

package-all:
  BUILD --platform=linux/amd64 --platform=linux/386 --platform=linux/arm64 --platform=linux/arm/v7 --platform=darwin/amd64 +package

docker-all:
  BUILD --platform=linux/amd64 --platform=linux/386 --platform=linux/arm64 --platform=linux/arm/v7 +docker

crosscompiler:
  # This image is missing a few platforms, so we'll add them locally
  FROM --platform=linux/amd64 bdwyertech/go-crosscompile
  RUN curl -sfL "https://musl.cc/armv7l-linux-musleabihf-cross.tgz" | tar zxf - -C /usr/ --strip-components=1
  RUN curl -sfL "https://musl.cc/i686-linux-musl-cross.tgz" | tar zxf - -C /usr/ --strip-components=1
  RUN curl -sfL "https://musl.cc/x86_64-linux-musl-cross.tgz" | tar zxf - -C /usr/ --strip-components=1

code:
  FROM --platform=linux/amd64 +crosscompiler
  COPY . /build
  # GIT CLONE --branch=$version git@github.com:owncast/owncast.git /build

build:
  ARG EARTHLY_GIT_HASH # provided by Earthly
  ARG TARGETPLATFORM   # provided by Earthly
  ARG TARGETOS         # provided by Earthly
  ARG TARGETARCH       # provided by Earthly
  ARG GOOS=$TARGETOS
  ARG GOARCH=$TARGETARCH

  FROM --platform=linux/amd64 +code

  RUN echo $EARTHLY_GIT_HASH
  RUN echo "Finding CC configuration for $TARGETPLATFORM"
  IF [ "$TARGETPLATFORM" = "linux/amd64" ]
    ARG NAME=linux-64bit
    ARG CC=x86_64-linux-musl-gcc
    ARG CXX=x86_64-linux-musl-g++
  ELSE IF [ "$TARGETPLATFORM" = "linux/386" ]
    ARG NAME=linux-32bit
    ARG CC=i686-linux-musl-gcc
    ARG CXX=i686-linux-musl-g++
  ELSE IF [ "$TARGETPLATFORM" = "linux/arm64" ]
    ARG NAME=linux-arm64
    ARG CC=aarch64-linux-musl-gcc
    ARG CXX=aarch64-linux-musl-g++
  ELSE IF [ "$TARGETPLATFORM" = "linux/arm/v7" ]
    ARG NAME=linux-arm7
    ARG CC=armv7l-linux-musleabihf-gcc
    ARG CXX=armv7l-linux-musleabihf-g++
    ARG GOARM=7
  ELSE IF [ "$TARGETPLATFORM" = "darwin/amd64" ]
    ARG NAME=macOS-64bit
    ARG CC=o64-clang
    ARG CXX=o64-clang++
  ELSE
    RUN echo "Failed to find CC configuration for $TARGETPLATFORM"
    ARG --required CC
    ARG --required CXX
  END

  ENV CGO_ENABLED=1
  ENV GOOS=$GOOS
  ENV GOARCH=$GOARCH
  ENV GOARM=$GOARM
  ENV CC=$CC
  ENV CXX=$CXX

  WORKDIR /build
  # MacOSX disallows static executables, so we omit the static flag on this platform
  RUN go build -a -installsuffix cgo -ldflags "$([ "$GOOS"z != darwinz ] && echo "-linkmode external -extldflags -static ") -s -w -X github.com/owncast/owncast/config.GitCommit=$EARTHLY_GIT_HASH -X github.com/owncast/owncast/config.VersionNumber=$version -X github.com/owncast/owncast/config.BuildPlatform=$NAME" -o owncast main.go
  COPY +tailwind/prod-tailwind.min.css /build/dist/webroot/js/web_modules/tailwindcss/dist/tailwind.min.css

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
  SAVE ARTIFACT /build/prod-tailwind.min.css prod-tailwind.min.css

package:
  RUN apk add --update --no-cache zip >> /dev/null

  ARG TARGETPLATFORM   # provided by Earthly
  IF [ "$TARGETPLATFORM" = "linux/amd64" ]
    ARG NAME=linux-64bit
  ELSE IF [ "$TARGETPLATFORM" = "linux/386" ]
    ARG NAME=linux-32bit
  ELSE IF [ "$TARGETPLATFORM" = "linux/arm64" ]
    ARG NAME=linux-arm64
  ELSE IF [ "$TARGETPLATFORM" = "linux/arm/v7" ]
    ARG NAME=linux-arm7
  ELSE IF [ "$TARGETPLATFORM" = "darwin/amd64" ]
    ARG NAME=macOS-64bit
  ELSE
    ARG NAME=custom
  END

  COPY (+build/webroot --platform $TARGETPLATFORM) /build/dist/webroot
  COPY (+build/owncast --platform $TARGETPLATFORM) /build/dist/owncast
  COPY (+build/README.md --platform $TARGETPLATFORM) /build/dist/README.md
  ENV ZIPNAME owncast-$version-$NAME.zip
  RUN cd /build/dist && zip -r -q -8 /build/dist/owncast.zip .
  SAVE ARTIFACT /build/dist/owncast.zip owncast.zip AS LOCAL dist/$ZIPNAME

docker:
  ARG image=ghcr.io/owncast/owncast
  ARG tag=develop
  ARG TARGETPLATFORM
  FROM --platform=$TARGETPLATFORM alpine:latest
  RUN apk update && apk add --no-cache ffmpeg ffmpeg-libs ca-certificates unzip && update-ca-certificates
  WORKDIR /app
  COPY --platform=$TARGETPLATFORM +package/owncast.zip /app
  RUN unzip -x owncast.zip && mkdir data
  ENTRYPOINT ["/app/owncast"]
  EXPOSE 8080 1935
  SAVE IMAGE --push $image:$tag
