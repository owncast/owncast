#!/bin/bash

# Requirements:
#   ffmpeg (a recent version with loop video support)
#   a Sans family font (for overlay text)
#   awk
#   readlink

# Example: ./test/ocTestStream.sh ~/Downloads/*.mp4 rtmp://127.0.0.1/live/abc123

exec node ./ocTestStream.js
