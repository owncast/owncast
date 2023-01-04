#!/usr/bin/env bash

# A recent version of ffmpeg is required for the loop of the provided videos
# to repeat indefinitely.
# Example: ./test/ocTestStream.sh ~/Downloads/*.mp4 rtmp://localhost/live/abc123

if ! [[ $1 ]]
then
  echo "ocTestStream is used for sending pre-recorded content to a RTMP server."
  echo "Will default to localhost with the stream key of abc123 if one isn't provided."
  echo "./ocTestStream.sh *.mp4 [RTMPDESINATION]"
  exit
fi

# Make the destination optional and point to localhost with default key
if [[ ${*: -1} == *"rtmp://"* ]]; then
  echo "RTMP server is specified"
  DESTINATION_HOST=${*: -1}
  FILE_COUNT=$(( ${#} - 1 ))
else
  echo "RTMP server is not specified"
  DESTINATION_HOST="rtmp://localhost/live/abc123"
  FILE_COUNT=${#}
fi

if [[ FILE_COUNT -eq 0 ]]; then
  echo "ERROR: ocTestStream needs a video file for sending to the RTMP server."
  exit
fi

CONTENT=${*:1:${FILE_COUNT}}

# Delete the old list of files if it exists
if test -f list.txt; then
  rm list.txt
fi

for file in $CONTENT
do
  echo "file '$file'" >> list.txt
done

function finish {
  rm list.txt
}
trap finish EXIT

echo "Streaming a loop of ${FILE_COUNT} videos to $DESTINATION_HOST."
if [[ FILE_COUNT -gt 1 ]]; then
  echo "Warning: If these files differ greatly in formats transitioning from one to another may not always work correctly."
fi
echo "$CONTENT"
echo "...press ctl+c to exit"

ffmpeg -hide_banner -loglevel panic -stream_loop -1 -re -f concat -safe 0 -i list.txt -vcodec libx264 -profile:v high -g 48 -r 24 -sc_threshold 0 -b:v 1300k -preset veryfast -acodec copy -vf drawtext="fontfile=monofonto.ttf: fontsize=96: box=1: boxcolor=black@0.75: boxborderw=5: fontcolor=white: x=(w-text_w)/2: y=((h-text_h)/2)+((h-text_h)/4): text='%{gmtime\:%H\\\\\:%M\\\\\:%S}'" -f flv "$DESTINATION_HOST"
