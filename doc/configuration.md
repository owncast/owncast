# Configuration


## Video Quality

Owncast supports HLS [Adaptive bitrate streaming](https://en.wikipedia.org/wiki/Adaptive_bitrate_streaming), or in other words, different video qualities that can be used for different network conditions.

You can edit the `config.yaml` file and add as many stream _variants_ as you like under the `videoSettings` block, like so:

```
  streamQualities:
    - low:
      bitrate: 700
      scaledWidth: 600
      encoderPreset: superfast
    - medium:
      bitrate: 900
```

You must have at least one bitrate specified.

### Important caveats

#### CPU Usage

Each bitrate variant adds significant CPU usage and slows down the overall generation of video segments.  If you have a slow server running Owncast you should probably only have one bitrate variant in play.  If you add more and you notice that playback becomes choppy it's likely that everything is running too slowly for consistent playback.  Consider removing the additional variants and tweaking your single variant so it supports a wider variety of network conditions.

#### Disk Usage

More stream quality variants requires more disk space, since it's another copy of the video on disk.  If you're serving video locally and you have enough disk space then it's probably no big deal and files will rather quickly get rotated and cleaned up.  If you're using something like [S3 for storage](S3.md) then files won't get cleaned up until some point in the future, so you'll have more remote storage use in play.

## Customization

Editing `config.yaml` allows you to change the name, description, logo and links to external sites.

`socialHandles` currently supports the following services by name:

* `facebook`
* `twitter`
* `instagram`
* `snapchat`
* `tiktok`
* `soundcloud`
* `bandcamp`
* `patreon`
* `youtube`
* `spotify`
* `twitch`
* `paypal`
* `github`
* `linkedin`
* `discord`
* `mastadon`

Update your `tags` in the config to display the topics type of content you want to call attention to.