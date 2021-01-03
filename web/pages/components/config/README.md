# Config


TODO: explain how to use <Form> and how the custom `form-xxxx` components work together.


## Misc notes
- `instanceDetails` needs to be filled out before `yp.enabled` can be turned on.



## Config data structure (with default values)
```
{
  streamKey: '',
  instanceDetails: {
    tags: [],
    nsfw: false,
  },
  yp: {
    enabled: false,
    instance: '',
  },
  videoSettings: {
    videoQualityVariants: [
      {
        audioPassthrough: false,
        videoPassthrough: false,
        videoBitrate: 0,
        audioBitrate: 0,
        framerate: 0,
      },
    ],
  }
};
```