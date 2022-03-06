VERSION 0.6

FROM bdwyertech/go-crosscompile

ARG version=dev
WORKDIR /build

code:
  # COPY . /build
  GIT CLONE --branch=develop git@github.com:owncast/owncast.git /build

build:
  FROM +code
  RUN echo $EARTHLY_GIT_HASH
  ARG NAME
  ARG GOOS
  ARG GOARCH
  ARG CC=$CC
  ARG CXX=$CXX
  ENV CGO_ENABLED=1
  ENV GOOS=$GOOS
  ENV GOARCH=$GOARCH
  RUN go build -a -installsuffix cgo -ldflags "-linkmode external -extldflags -static -s -w -X github.com/owncast/owncast/config.GitCommit=$EARTHLY_GIT_HASH -X github.com/owncast/owncast/config.VersionNumber=$version -X github.com/owncast/owncast/config.BuildPlatform=$NAME" -o owncast main.go
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
  SAVE ARTIFACT /build/dist/owncast.zip dist/$ZIPNAME AS LOCAL dist/$ZIPNAME

build-all:
  BUILD +linux-arm64
  BUILD +linux-arm7
  BUILD +linux-amd64
  BUILD +linux-i386
  BUILD +darwin-amd64

linux-arm64:
  ARG NAME=linux-arm64
  FROM +build --platform=linux/arm64 --CC=o64-clang --CXX=o64-clang++ --GOOS=linux --GOARCH=arm64
  FROM +package --NAME=$NAME

linux-arm7:
  ARG NAME=linux-arm7
  FROM +build --platform=linux/arm/v7--CC=o64-clang --CXX=o64-clang++ --GOOS=linux --GOARCH=arm/v7
  FROM +package --NAME=$NAME

linux-amd64:
  ARG NAME=linux-64bit
  FROM +build --GOOS=linux --GOARCH=amd64
  FROM +package --NAME=$NAME

linux-i386:
  ARG NAME=linux-32bit
  FROM +build --GOOS=linux --GOARCH=i386
  FROM +package --NAME=$NAME

darwin-amd64:
  ARG NAME=macOS-64bit
  FROM +build --platform=darwin/amd64 --CC=o64-clang --CXX=o64-clang++ --GOOS=darwin --GOARCH=amd64
  FROM +package --NAME=$NAME
