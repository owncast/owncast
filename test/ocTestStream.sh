#!/usr/bin/env bash

# A recent version of ffmpeg is required for the loop of the provided videos
# to repeat indefinitely.
# Example: ./test/ocTestStream.sh ~/Downloads/*.mp4 rtmp://localhost/live/abc123

if ! ([[ $1 ]])
then
  echo "ocTestStream is used for sending pre-recorded content to a RTMP server."
  echo "Will default to localhost with the stream key of abc123 if one isn't provided."
  echo "./ocTestStream.sh *.mp4 [RTMPDESINATION]"
exit
fi

# Make the destination optional and point to localhost with default key
if [[ ${@: -1} == *"rtmp://"* ]]; then
  echo "RTMP is specified"
  DESTINATION_HOST=${@: -1}
  array=( $@ )
  ARGS_LEN=${#array[@]}
  CONTENT=${array[@]:0:$ARGS_LEN-1}
  DESTINATION_HOST=${@: -1}
  FILE_COUNT=$( expr ${#} - 1 )
else
  DESTINATION_HOST="rtmp://localhost/live/abc123"
  array=( $@ )
  ARGS_LEN=${#array[@]}
  CONTENT=${array[@]:0:$ARGS_LEN}
  FILE_COUNT=$( expr ${#} )
fi

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

echo "Streaming a loop of ${FILE_COUNT} videos to $DESTINATION_HOST.  Warning: If these files differ greatly in formats transitioning from one to another may not always work correctly...  ctl+c to exit"

ffmpeg -hide_banner -loglevel panic -stream_loop -1 -re -f concat -safe 0 -i list.txt -vcodec libx264 -profile:v main -sc_threshold 0 -b:v 1300k -preset veryfast -acodec copy -f flv $DESTINATION_HOST
