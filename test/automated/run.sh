#!/bin/bash

npm install --silent > /dev/null

git clone https://github.com/owncast/owncast
pushd owncast > /dev/null

cp config-default.yaml config.yaml
go build -o owncast main.go pkged.go

curl -sL https://github.com/vot/ffbinaries-prebuilt/releases/download/v4.2.1/ffmpeg-4.2.1-linux-64.zip --output ffmpeg.zip >> /dev/null
unzip -f ffmpeg.zip

./owncast &
SERVER_PID=$!

popd > /dev/null
sleep 5

./owncast/ffmpeg -hide_banner -loglevel panic -re -i test.mp4 -vcodec libx264 -profile:v main -sc_threshold 0 -b:v 1300k -acodec copy -f flv rtmp://127.0.0.1/live/abc123 &
FFMPEG_PID=$!

function finish {
  kill $SERVER_PID $FFMPEG_PID
}
trap finish EXIT
echo "Streaming test.mp4..."

sleep 15

npm test