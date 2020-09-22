#!/usr/bin/env bash

# A recent version of ffmpeg is required for the loop of the provided videos
# to repeat indefinitely.

if ! ([[ $1 ]] && [[ $2 ]])
then
  echo "ocTestStream is used for sending pre-recorded content to a RTMP server."
  echo "./ocTestStream.sh *.mp4 rtmp://desination.host/live/streamkey"
exit
fi

array=( $@ )
ARGS_LEN=${#array[@]}
CONTENT=${array[@]:0:$ARGS_LEN-1}
DESTINATION_HOST=${@: -1}
FILE_COUNT=$( expr ${#} - 1 )

# Delete the old list of files if it exists
if test -f list.txt; then
  rm list.txt
fi

for file in $CONTENT
do
  echo "file '$file'" >> list.txt
done

echo "Streaming a loop of ${FILE_COUNT} videos to $DESTINATION_HOST...  ctl+c to exit"

ffmpeg -hide_banner -loglevel panic -stream_loop -1 -re -f concat -safe 0 -i list.txt -vcodec libx264 -profile:v main -sc_threshold 0 -b:v 1300k -preset ultrafast -acodec aac -b:a 192k -x264opts keyint=50 -g 25 -pix_fmt yuv420p -f flv $DESTINATION_HOST
