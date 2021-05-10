FROM golang:alpine AS build
RUN apk add --no-cache gcc build-base linux-headers

WORKDIR /build
COPY . /build
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o owncast .


FROM alpine
RUN apk add --no-cache ffmpeg ffmpeg-libs

WORKDIR /app
COPY webroot /app/webroot
COPY static /app/static
COPY --from=build /build/owncast /app/owncast

EXPOSE 8080 1935
CMD ["/app/owncast"]
