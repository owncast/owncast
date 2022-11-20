#!/bin/sh

ffmpeg -hide_banner -loglevel panic -re -f lavfi \
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
	-f flv "rtmp://127.0.0.1/live/abc123" >/dev/null
