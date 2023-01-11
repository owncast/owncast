#!/bin/bash

ffmpegInstall(){
    # install a specific version of ffmpeg

    FFMPEG_VER="4.4.1"
    FFMPEG_PATH="$(pwd)/ffmpeg-$FFMPEG_VER"

    if ! [[ -d "$FFMPEG_PATH" ]]; then
        mkdir "$FFMPEG_PATH"
    fi

    pushd "$FFMPEG_PATH" >/dev/null || exit

    if [[ -x "$FFMPEG_PATH/ffmpeg" ]]; then
    
        ffmpeg_version=$("$FFMPEG_PATH/ffmpeg" -version | awk -F 'ffmpeg version' '{print $2}' | awk 'NR==1{print $1}')

        if [[ "$ffmpeg_version" == "$FFMPEG_VER-static" ]]; then
            return 0
        else
            mv "$FFMPEG_PATH/ffmpeg" "$FFMPEG_PATH/ffmpeg.bk" || rm -f "$FFMPEG_PATH/ffmpeg"
        fi
    fi

    rm -f ffmpeg.zip
    curl -sL --fail https://github.com/ffbinaries/ffbinaries-prebuilt/releases/download/v${FFMPEG_VER}/ffmpeg-${FFMPEG_VER}-linux-64.zip --output ffmpeg.zip >/dev/null
    unzip -o ffmpeg.zip >/dev/null && rm -f ffmpeg.zip
    chmod +x ffmpeg
    PATH=$FFMPEG_PATH:$PATH

    popd >/dev/null || exit
}