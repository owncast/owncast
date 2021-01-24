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

TODO:
- fill out readme for how to use form fields and about data flow w/ ant and react

- cleanup 
  - more consitent constants
  - cleanup types
  - cleanup style sheets..? make style module for each config page? (but what about ant deisgn overrides?)
- redesign
  - label+form layout - put them into a table, table of rows
  - change Encpder preset into slider
  - page headers - diff color? 
  - fix social handles icon in table
  - consolidate things into 1 page?
  - things could use smaller font?