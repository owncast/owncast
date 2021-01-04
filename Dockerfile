# BUILD STEP

FROM golang:alpine as build

RUN set -ex && apk add --no-cache gcc build-base linux-headers

RUN mkdir /build
WORKDIR /build

COPY go.mod /build
COPY go.sum /build
RUN CGO_ENABLED=1 GOOS=linux go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-s -w -extldflags "-static"' -o owncast .

# COMMAND STEP

FROM alpine:latest as final
EXPOSE 8080 1935

RUN set -ex && apk add --no-cache ffmpeg ffmpeg-libs
	
RUN mkdir /app
WORKDIR /app
COPY --from=build /build/owncast ./owncast
COPY . .

CMD ["/app/owncast"]
