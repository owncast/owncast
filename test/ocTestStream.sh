#!/bin/bash

# Requirements:
#   ffmpeg (a recent version with loop video support)
#   a Sans family font (for overlay text)
#   awk
#   readlink

# Example: ./test/ocTestStream.sh ~/Downloads/*.mp4 rtmp://127.0.0.1/live/abc123


ffmpeg_execs=( 'ffmpeg' 'ffmpeg.exe' )
ffmpeg_paths=( './' '../' '' )

for _ffmpeg_exec in "${ffmpeg_execs[@]}"; do
  for _ffmpeg_path in "${ffmpeg_paths[@]}"; do
    if [[ -x "$(command -v "${_ffmpeg_path}${_ffmpeg_exec}")" ]]; then
      ffmpeg_exec="${_ffmpeg_path}${_ffmpeg_exec}"
      break
    fi
  done
done

if [[ ${*: -1} == "--help" ]]; then
  echo "ocTestStream is used for sending pre-recorded or internal test content to an RTMP server."
  echo "Usage: ./ocTestStream.sh [VIDEO_FILES] [RTMP_DESINATION]"
  echo "VIDEO_FILES: path to one or multiple videos for sending to the RTMP server (optional)"
  echo "RTMP_DESINATION: URL of RTMP server with key (optional; default: rtmp://127.0.0.1/live/abc123)"
  exit
elif [[ ${*: -1} == *"rtmp://"* ]]; then
  echo "RTMP server is specified"
  DESTINATION_HOST=${*: -1}
  FILE_COUNT=$(( ${#} - 1 ))
else
  echo "RTMP server is not specified"
  DESTINATION_HOST="rtmp://127.0.0.1/live/abc123"
  FILE_COUNT=${#}
fi

if [[ -z "$ffmpeg_exec" ]]; then
	echo "ERROR: ffmpeg was not found in path or in the current directory! Please install ffmpeg before using this script."
  exit 1
else
  ffmpeg_version=$("$ffmpeg_exec" -version | awk -F 'ffmpeg version' '{print $2}' | awk 'NR==1{print $1}')
  echo "ffmpeg executable: $ffmpeg_exec ($ffmpeg_version)"
  echo "ffmpeg path: $(readlink -f "$(which "$ffmpeg_exec")")" 
fi

if [[ ${FILE_COUNT} -eq 0 ]]; then
  echo "Streaming internal test video loop to $DESTINATION_HOST"
  echo "...press ctrl+c to exit"

  command "${ffmpeg_exec}" -hide_banner -loglevel panic -nostdin -re -f lavfi \
    -i "testsrc=size=1280x720:rate=60[out0];sine=frequency=400:sample_rate=48000[out1]" \
    -vf "[in]drawtext=fontsize=96: box=1: boxcolor=black@0.75: boxborderw=5: fontcolor=white: x=(w-text_w)/2: y=((h-text_h)/2)+((h-text_h)/-2): text='Owncast Test Stream', drawtext=fontsize=96: box=1: boxcolor=black@0.75: boxborderw=5: fontcolor=white: x=(w-text_w)/2: y=((h-text_h)/2)+((h-text_h)/2): text='%{gmtime\:%H\\\\\:%M\\\\\:%S} UTC'[out]" \
    -nal-hrd cbr \
    -metadata:s:v encoder=test \
    -vcodec libx264 \
    -acodec aac \
    -preset veryfast \
    -profile:v baseline \
    -tune zerolatency \
    -bf 0 \
    -g 0 \
    -b:v 6320k \
    -b:a 160k \
    -ac 2 \
    -ar 48000 \
    -minrate 6320k \
    -maxrate 6320k \
    -bufsize 6320k \
    -muxrate 6320k \
    -r 60 \
    -pix_fmt yuv420p \
    -color_range 1 -colorspace 1 -color_primaries 1 -color_trc 1 \
    -flags:v +global_header \
    -bsf:v dump_extra \
    -x264-params "nal-hrd=cbr:min-keyint=2:keyint=2:scenecut=0:bframes=0" \
    -f flv "$DESTINATION_HOST"

else

  CONTENT=${*:1:${FILE_COUNT}}

  rm -f list.txt
  for file in $CONTENT
  do
    if [[ -f "$file" ]]; then
      echo "file '$file'" >> list.txt
    else
      echo "ERROR: File not found: $file"
      exit 1
    fi
    
  done

  function finish {
    rm list.txt
  }
  trap finish EXIT

  echo "Streaming a loop of ${FILE_COUNT} video(s) to $DESTINATION_HOST"
  if [[ ${FILE_COUNT} -gt 1 ]]; then
    echo "Warning: If these files differ greatly in formats, transitioning from one to another may not always work correctly."
  fi
  echo "$CONTENT"
  echo "...press ctrl+c to exit"

  command "${ffmpeg_exec}" -hide_banner -loglevel panic -nostdin -stream_loop -1 -re -f concat \
    -safe 0 \
    -i list.txt \
    -vcodec libx264 \
    -profile:v high \
    -g 48 \
    -r 24 \
    -sc_threshold 0 \
    -b:v 1300k \
    -preset veryfast \
    -acodec copy \
    -vf drawtext="fontsize=96: box=1: boxcolor=black@0.75: boxborderw=5: fontcolor=white: x=(w-text_w)/2: y=((h-text_h)/2)+((h-text_h)/4): text='%{gmtime\:%H\\\\\:%M\\\\\:%S}'" \
    -f flv "$DESTINATION_HOST"
fi
