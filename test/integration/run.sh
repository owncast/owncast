#!/bin/bash
npm install > /dev/null

git clone https://github.com/owncast/owncast
pushd owncast
cp config-example.yaml config.yaml
go build -o owncast main.go pkged.go

./owncast &
SERVER_PID=$!

popd

sleep 5
ffmpeg -hide_banner -loglevel panic -re -i test.mp4 -vcodec libx264 -profile:v main -sc_threshold 0 -b:v 1300k -acodec copy -f flv rtmp://127.0.0.1/live/abc123 &
FFMPEG_PID=$!

function finish {
  kill $SERVER_PID $FFMPEG_PID
  # rm -rf owncast
}
trap finish EXIT

sleep 15

npm test