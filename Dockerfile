FROM golang:alpine
EXPOSE 8080 1935
WORKDIR /app
ADD . .
RUN set -ex && \
    apk add --no-cache ffmpeg ffmpeg-libs && \
    apk add --no-cache gcc build-base linux-headers && \ 
    CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o owncast .
CMD ["/app/owncast"]
