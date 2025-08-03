@echo off
setlocal enabledelayedexpansion

REM Requirements:
REM   ffmpeg (a recent version with loop video support)
REM   a Sans family font (for overlay text)

REM Example: ocTestStream.bat "C:\Videos\*.mp4" rtmp://127.0.0.1/live/abc123

set "ffmpeg_exec="
set "DESTINATION_HOST=rtmp://127.0.0.1/live/abc123"
set "FILE_COUNT=0"

REM Try to find ffmpeg executable
for %%e in (ffmpeg.exe ffmpeg) do (
    for %%p in ("." ".." "") do (
        if exist "%%~p%%e" (
            set "ffmpeg_exec=%%~p%%e"
            goto :found_ffmpeg
        )
    )
)

REM Check if ffmpeg is in PATH
where ffmpeg.exe >nul 2>&1
if !ERRORLEVEL! == 0 (
    set "ffmpeg_exec=ffmpeg.exe"
    goto :found_ffmpeg
)

where ffmpeg >nul 2>&1
if !ERRORLEVEL! == 0 (
    set "ffmpeg_exec=ffmpeg"
    goto :found_ffmpeg
)

:found_ffmpeg

REM Check for help argument
if "%~1"=="--help" (
    echo ocTestStream is used for sending pre-recorded or internal test content to an RTMP server.
    echo Usage: ocTestStream.bat [VIDEO_FILES] [RTMP_DESTINATION]
    echo VIDEO_FILES: path to one or multiple videos for sending to the RTMP server ^(optional^)
    echo RTMP_DESTINATION: URL of RTMP server with key ^(optional; default: rtmp://127.0.0.1/live/abc123^)
    exit /b 0
)

REM Parse arguments - check if last argument is RTMP URL
set "args_count=0"
set "last_arg="
for %%i in (%*) do (
    set /a args_count+=1
    set "last_arg=%%i"
)

REM Check if last argument contains rtmp://
echo !last_arg! | findstr /c:"rtmp://" >nul
if !ERRORLEVEL! == 0 (
    echo RTMP server is specified
    set "DESTINATION_HOST=!last_arg!"
    set /a FILE_COUNT=!args_count!-1
) else (
    echo RTMP server is not specified
    set "FILE_COUNT=!args_count!"
)

REM Check if ffmpeg was found
if "!ffmpeg_exec!"=="" (
    echo ERROR: ffmpeg was not found in path or in the current directory! Please install ffmpeg before using this script.
    exit /b 1
) else (
    REM Get ffmpeg version
    for /f "tokens=3" %%i in ('"!ffmpeg_exec!" -version 2^>nul ^| findstr "ffmpeg version"') do (
        set "ffmpeg_version=%%i"
        goto :got_version
    )
    :got_version
    echo ffmpeg executable: !ffmpeg_exec! (!ffmpeg_version!)
    
    REM Get ffmpeg path
    where "!ffmpeg_exec!" >nul 2>&1
    if !ERRORLEVEL! == 0 (
        for /f "delims=" %%i in ('where "!ffmpeg_exec!" 2^>nul') do (
            echo ffmpeg path: %%i
            goto :got_path
        )
    )
    :got_path
)

REM If no files specified, stream test pattern
if !FILE_COUNT! == 0 (
    echo Streaming internal test video loop to !DESTINATION_HOST!
    echo ...press ctrl+c to exit
    
    "!ffmpeg_exec!" -hide_banner -loglevel warning -nostdin -re -f lavfi ^
        -i "testsrc=size=1280x720:rate=60[out0];sine=frequency=400:sample_rate=48000[out1]" ^
        -vf "[in]drawtext=fontsize=96: box=1: boxcolor=black@0.75: boxborderw=5: fontcolor=white: x=(w-text_w)/2: y=((h-text_h)/2)+((h-text_h)/-2): text='Owncast Test Stream', drawtext=fontsize=96: box=1: boxcolor=black@0.75: boxborderw=5: fontcolor=white: x=(w-text_w)/2: y=((h-text_h)/2)+((h-text_h)/2): text='%%{gmtime\:%%H-%%M-%%S} UTC'[out]" ^
        -nal-hrd cbr ^
        -metadata:s:v encoder=test ^
        -vcodec libx264 ^
        -acodec aac ^
        -preset veryfast ^
        -profile:v baseline ^
        -tune zerolatency ^
        -bf 0 ^
        -g 0 ^
        -b:v 6320k ^
        -b:a 160k ^
        -ac 2 ^
        -ar 48000 ^
        -minrate 6320k ^
        -maxrate 6320k ^
        -bufsize 6320k ^
        -muxrate 6320k ^
        -r 60 ^
        -pix_fmt yuv420p ^
        -color_range 1 -colorspace 1 -color_primaries 1 -color_trc 1 ^
        -flags:v +global_header ^
        -bsf:v dump_extra ^
        -x264-params "nal-hrd=cbr:min-keyint=2:keyint=2:scenecut=0:bframes=0" ^
        -f flv "!DESTINATION_HOST!"
) else (
    REM Handle video files
    if exist list.txt del list.txt
    
    set "file_index=0"
    set "valid_files=0"
    
    REM Process each argument (except the last one if it's RTMP URL)
    for %%i in (%*) do (
        set /a file_index+=1
        
        REM Skip last argument if it's RTMP URL
        if !file_index! lss !args_count! (
            if exist "%%i" (
                echo file '%%i' >> list.txt
                set /a valid_files+=1
                echo Found file: %%i
            ) else (
                echo ERROR: File not found: %%i
                if exist list.txt del list.txt
                exit /b 1
            )
        ) else (
            REM This is the last argument
            echo !last_arg! | findstr /c:"rtmp://" >nul
            if !ERRORLEVEL! neq 0 (
                REM Last argument is not RTMP, so it's a file
                if exist "%%i" (
                    echo file '%%i' >> list.txt
                    set /a valid_files+=1
                    echo Found file: %%i
                ) else (
                    echo ERROR: File not found: %%i
                    if exist list.txt del list.txt
                    exit /b 1
                )
            )
        )
    )
    
    if !valid_files! == 0 (
        echo ERROR: No valid video files found
        if exist list.txt del list.txt
        exit /b 1
    )
    
    echo Streaming a loop of !valid_files! video^(s^) to !DESTINATION_HOST!
    if !valid_files! gtr 1 (
        echo Warning: If these files differ greatly in formats, transitioning from one to another may not always work correctly.
    )
    echo ...press ctrl+c to exit
    
    "!ffmpeg_exec!" -hide_banner -loglevel panic -nostdin -stream_loop -1 -re -f concat ^
        -safe 0 ^
        -i list.txt ^
        -vcodec libx264 ^
        -profile:v high ^
        -g 48 ^
        -r 24 ^
        -sc_threshold 0 ^
        -b:v 1300k ^
        -preset veryfast ^
        -acodec copy ^
        -vf drawtext="fontsize=96: box=1: boxcolor=black@0.75: boxborderw=5: fontcolor=white: x=(w-text_w)/2: y=((h-text_h)/2)+((h-text_h)/4): text='%%{gmtime\:%%H-%%M-%%S}'" ^
        -f flv "!DESTINATION_HOST!"
    
    REM Cleanup
    if exist list.txt del list.txt
)