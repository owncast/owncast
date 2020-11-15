# Builder
FROM golang:alpine AS builder

RUN apk update && \
    apk add --no-cache gcc build-base linux-headers

WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o owncast .

# Actual image
FROM alpine

RUN apk update && \
    apk add --no-cache ffmpeg ffmpeg-libs

COPY --from=builder /app /app

EXPOSE 8080 1935
WORKDIR /app
ENTRYPOINT /app/owncast
