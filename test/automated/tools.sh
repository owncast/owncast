#!/bin/bash

ffmpegInstall(){
    # install a specific version of ffmpeg

    FFMPEG_VER="4.4.1"
    FFMPEG_PATH="$(pwd)/ffmpeg-$FFMPEG_VER"

    if ! [[ -d "$FFMPEG_PATH" ]]; then
        mkdir "$FFMPEG_PATH"
    fi

    pushd "$FFMPEG_PATH" >/dev/null || exit

    if [[ -x "./ffmpeg" ]]; then
    
        ffmpeg_version=$(./ffmpeg -version | awk -F 'ffmpeg version' '{print $2}' | awk 'NR==1{print $1}')

        if [[ "$ffmpeg_version" == "$FFMPEG_VER-static" ]]; then
            return 0
        else
            mv ./ffmpeg ./ffmpeg.bk || rm -f ./ffmpeg
        fi
    fi

    curl -sL --fail https://github.com/ffbinaries/ffbinaries-prebuilt/releases/download/v${FFMPEG_VER}/ffmpeg-${FFMPEG_VER}-linux-64.zip --output ffmpeg.zip >/dev/null
    unzip -o ffmpeg.zip >/dev/null && rm -f ffmpeg.zip
    chmod +x ffmpeg
    PATH=$FFMPEG_PATH:$PATH

    popd >/dev/null || exit
}