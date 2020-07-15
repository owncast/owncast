FROM golang:alpine
EXPOSE 8080 1935
RUN mkdir /app 
ADD . /app
WORKDIR /app
RUN apk add --no-cache ffmpeg ffmpeg-libs
RUN apk update && apk add --no-cache gcc build-base linux-headers
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o owncast .
WORKDIR /app
CMD ["/app/owncast"]