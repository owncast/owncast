VERSION --new-platform 0.6

FROM --platform=linux/amd64 alpine:3.15.5
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
  RUN apk add --update --no-cache tar gzip upx >> /dev/null
  RUN curl -sfL "https://owncast-infra.nyc3.cdn.digitaloceanspaces.com/build/armv7l-linux-musleabihf-cross.tgz" | tar zxf - -C /usr/ --strip-components=1
  RUN curl -sfL "https://owncast-infra.nyc3.cdn.digitaloceanspaces.com/build/i686-linux-musl-cross.tgz" | tar zxf - -C /usr/ --strip-components=1
  RUN curl -sfL "https://owncast-infra.nyc3.cdn.digitaloceanspaces.com/build/x86_64-linux-musl-cross.tgz" | tar zxf - -C /usr/ --strip-components=1

code:
  FROM --platform=linux/amd64 +crosscompiler
  COPY . /build

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
  RUN go build -a -installsuffix cgo -ldflags "$([ "$GOOS"z != darwinz ] && echo "-linkmode external -extldflags -static ") -s -w -X github.com/owncast/owncast/config.GitCommit=$EARTHLY_GIT_HASH -X github.com/owncast/owncast/config.VersionNumber=$version -X github.com/owncast/owncast/config.BuildPlatform=$NAME" -tags sqlite_omit_load_extension -o owncast main.go

	# Decrease the size of the shipped binary
	RUN upx --best --lzma owncast
	# Test the binary
	RUN upx -t owncast

  SAVE ARTIFACT owncast owncast

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

  COPY (+build/owncast --platform $TARGETPLATFORM) /build/dist/owncast
  ENV ZIPNAME owncast-$version-$NAME.zip
  RUN cd /build/dist && zip -r -q -8 /build/dist/owncast.zip .
  SAVE ARTIFACT /build/dist/owncast.zip owncast.zip AS LOCAL dist/$ZIPNAME

docker:
	# Multiple image names can be tagged at once. They should all be passed
	# in as space separated strings using the full account/repo:tag format.
	# https://github.com/earthly/earthly/blob/aea38448fa9c0064b1b70d61be717ae740689fb9/docs/earthfile/earthfile.md#assigning-multiple-image-names
  ARG TARGETPLATFORM
  FROM --platform=$TARGETPLATFORM alpine:3.15.5
  RUN apk update && apk add --no-cache ffmpeg ffmpeg-libs ca-certificates unzip && update-ca-certificates
  RUN addgroup -g 101 -S owncast && adduser -u 101 -S owncast -G owncast
  WORKDIR /app
  COPY --platform=$TARGETPLATFORM +package/owncast.zip /app
  RUN unzip -x owncast.zip && mkdir data

	# temporarily disable until we figure out how to move forward
  # RUN chown -R owncast:owncast /app
  # USER owncast

  ENTRYPOINT ["/app/owncast"]
  EXPOSE 8080 1935

  ARG images=ghcr.io/owncast/owncast:testing
	RUN echo "Saving images: ${images}"

	# Tag this image with the list of names
	# passed along.
	FOR --no-cache i IN ${images}
		SAVE IMAGE --push "${i}"
	END

dockerfile:
  FROM DOCKERFILE -f Dockerfile .

testing:
	ARG images
	FOR i IN ${images}
		RUN echo "Testing ${i}"
	END

unit-tests:
  FROM --platform=linux/amd64 bdwyertech/go-crosscompile
  COPY . /build
	WORKDIR /build
	RUN go test ./...

api-tests:
	FROM --platform=linux/amd64 bdwyertech/go-crosscompile
	RUN apk add npm font-noto && fc-cache -f
  COPY . /build
	WORKDIR /build/test/automated/api
	RUN npm install
	RUN ./run.sh

hls-tests:
	FROM --platform=linux/amd64 bdwyertech/go-crosscompile
	RUN apk add npm font-noto && fc-cache -f
  COPY . /build
	WORKDIR /build/test/automated/hls
	RUN npm install
	RUN ./run.sh
