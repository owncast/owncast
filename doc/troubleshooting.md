## CPU and RAM usage alerts

If your hardware is being maxed out then your video may not be processed and delivered fast enough to keep up with the real-time requirements of live video.

Here are some steps you can try taking to resolve this.

1. You may have too many bitrates defined as separate video quality variants.  Try removing one.
1. Change to a [faster encoder preset](https://github.com/gabek/owncast/blob/master/doc/encoding.md#encoder-preset) in your configuration.  If you're currently using `veryfast`, try `superfast`, for example.
1. Try reducing [the quality of the video you're sending to Owncast from your broadcasting software](https://github.com/gabek/owncast/blob/master/doc/encoding.md#how-you-configure-your-broadcasting-software-matters).
1. If you've gone down to a single bitrate, changed the encoder preset to the fastest, and experimented with different qualities in your broadcasting software, it's possible the server you're running Owncast is just not powerful enough for the task and you might need to try a different environment to run this on.
