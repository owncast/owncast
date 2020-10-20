export const parsedManifest = {
  allowCache: true,
  discontinuityStarts: [],
  duration: 6,
  endList: true,
  mediaGroups: {
    AUDIO: {
      audio: {
        ['en (main)']: {
          autoselect: true,
          default: true,
          language: 'en',
          playlists: [{
            attributes: {
              BANDWIDTH: 125000,
              CODECS: 'mp4a.40.2',
              NAME: '125000',
              ['PROGRAM-ID']: 1
            },
            contentProtection: {
              'com.widevine.alpha': {
                attributes: {
                  schemeIdUri: 'urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed'
                },
                pssh: new Uint8Array([181, 235, 45])
              }
            },
            endList: true,
            mediaSequence: 0,
            targetDuration: 1.984,
            resolvedUri: '',
            segments: [{
              duration: 1.984,
              map: {
                resolvedUri: 'https://www.example.com/125000/init.m4f',
                uri: '125000/init.m4f'
              },
              resolvedUri: 'https://www.example.com/125000/0.m4f',
              timeline: 0,
              uri: '125000/0.m4f',
              number: 0
            }, {
              duration: 1.984,
              map: {
                resolvedUri: 'https://www.example.com/125000/init.m4f',
                uri: '125000/init.m4f'
              },
              resolvedUri: 'https://www.example.com/125000/1.m4f',
              timeline: 0,
              uri: '125000/1.m4f',
              number: 1
            }, {
              duration: 1.984,
              map: {
                resolvedUri: 'https://www.example.com/125000/init.m4f',
                uri: '125000/init.m4f'
              },
              resolvedUri: 'https://www.example.com/125000/2.m4f',
              timeline: 0,
              uri: '125000/2.m4f',
              number: 2
            }, {
              duration: 0.04800000000000004,
              map: {
                resolvedUri: 'https://www.example.com/125000/init.m4f',
                uri: '125000/init.m4f'
              },
              resolvedUri: 'https://www.example.com/125000/3.m4f',
              timeline: 0,
              uri: '125000/3.m4f',
              number: 3
            }],
            timeline: 0,
            uri: ''
          }],
          uri: ''
        },
        ['es']: {
          autoselect: true,
          default: false,
          language: 'es',
          playlists: [{
            attributes: {
              BANDWIDTH: 125000,
              CODECS: 'mp4a.40.2',
              NAME: '125000',
              ['PROGRAM-ID']: 1
            },
            contentProtection: {
              'com.widevine.alpha': {
                attributes: {
                  schemeIdUri: 'urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed'
                },
                pssh: new Uint8Array([181, 235, 45])
              }
            },
            endList: true,
            targetDuration: 1.984,
            mediaSequence: 0,
            resolvedUri: '',
            segments: [{
              duration: 1.984,
              map: {
                resolvedUri: 'https://www.example.com/125000/es/init.m4f',
                uri: '125000/es/init.m4f'
              },
              resolvedUri: 'https://www.example.com/125000/es/0.m4f',
              timeline: 0,
              uri: '125000/es/0.m4f',
              number: 0
            }, {
              duration: 1.984,
              map: {
                resolvedUri: 'https://www.example.com/125000/es/init.m4f',
                uri: '125000/es/init.m4f'
              },
              resolvedUri: 'https://www.example.com/125000/es/1.m4f',
              timeline: 0,
              uri: '125000/es/1.m4f',
              number: 1
            }, {
              duration: 1.984,
              map: {
                resolvedUri: 'https://www.example.com/125000/es/init.m4f',
                uri: '125000/es/init.m4f'
              },
              resolvedUri: 'https://www.example.com/125000/es/2.m4f',
              timeline: 0,
              uri: '125000/es/2.m4f',
              number: 2
            }, {
              duration: 0.04800000000000004,
              map: {
                resolvedUri: 'https://www.example.com/125000/es/init.m4f',
                uri: '125000/es/init.m4f'
              },
              resolvedUri: 'https://www.example.com/125000/es/3.m4f',
              timeline: 0,
              uri: '125000/es/3.m4f',
              number: 3
            }],
            timeline: 0,
            uri: ''
          }],
          uri: ''
        }
      }
    },
    ['CLOSED-CAPTIONS']: {},
    SUBTITLES: {
      subs: {
        en: {
          autoselect: false,
          default: false,
          language: 'en',
          playlists: [{
            attributes: {
              BANDWIDTH: 256,
              NAME: 'en',
              ['PROGRAM-ID']: 1
            },
            mediaSequence: 0,
            endList: true,
            targetDuration: 6,
            resolvedUri: 'https://example.com/en.vtt',
            segments: [{
              duration: 6,
              resolvedUri: 'https://example.com/en.vtt',
              timeline: 0,
              uri: 'https://example.com/en.vtt',
              number: 0
            }],
            timeline: 0,
            uri: ''
          }],
          uri: ''
        },
        es: {
          autoselect: false,
          default: false,
          language: 'es',
          playlists: [{
            attributes: {
              BANDWIDTH: 256,
              NAME: 'es',
              ['PROGRAM-ID']: 1
            },
            endList: true,
            targetDuration: 6,
            mediaSequence: 0,
            resolvedUri: 'https://example.com/es.vtt',
            segments: [{
              duration: 6,
              resolvedUri: 'https://example.com/es.vtt',
              timeline: 0,
              uri: 'https://example.com/es.vtt',
              number: 0
            }],
            timeline: 0,
            uri: ''
          }],
          uri: ''
        }
      }
    },
    VIDEO: {}
  },
  playlists: [{
    attributes: {
      AUDIO: 'audio',
      SUBTITLES: 'subs',
      BANDWIDTH: 449000,
      CODECS: 'avc1.420015',
      NAME: '482',
      ['PROGRAM-ID']: 1,
      RESOLUTION: {
        height: 270,
        width: 482
      }
    },
    contentProtection: {
      'com.widevine.alpha': {
        attributes: {
          schemeIdUri: 'urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed'
        },
        pssh: new Uint8Array([181, 235, 45])
      }
    },
    endList: true,
    targetDuration: 1.9185833333333333,
    mediaSequence: 0,
    resolvedUri: '',
    segments: [{
      duration: 1.9185833333333333,
      map: {
        resolvedUri: 'https://www.example.com/482/init.m4f',
        uri: '482/init.m4f'
      },
      resolvedUri: 'https://www.example.com/482/0.m4f',
      timeline: 0,
      uri: '482/0.m4f',
      number: 0
    }, {
      duration: 1.9185833333333333,
      map: {
        resolvedUri: 'https://www.example.com/482/init.m4f',
        uri: '482/init.m4f'
      },
      resolvedUri: 'https://www.example.com/482/1.m4f',
      timeline: 0,
      uri: '482/1.m4f',
      number: 1
    }, {
      duration: 1.9185833333333333,
      map: {
        resolvedUri: 'https://www.example.com/482/init.m4f',
        uri: '482/init.m4f'
      },
      resolvedUri: 'https://www.example.com/482/2.m4f',
      timeline: 0,
      uri: '482/2.m4f',
      number: 2
    }, {
      duration: 0.24425000000000008,
      map: {
        resolvedUri: 'https://www.example.com/482/init.m4f',
        uri: '482/init.m4f'
      },
      resolvedUri: 'https://www.example.com/482/3.m4f',
      timeline: 0,
      uri: '482/3.m4f',
      number: 3
    }],
    timeline: 0,
    uri: ''
  }, {
    attributes: {
      AUDIO: 'audio',
      SUBTITLES: 'subs',
      BANDWIDTH: 3971000,
      CODECS: 'avc1.64001e',
      NAME: '720',
      ['PROGRAM-ID']: 1,
      RESOLUTION: {
        height: 404,
        width: 720
      }
    },
    contentProtection: {
      'com.widevine.alpha': {
        attributes: {
          schemeIdUri: 'urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed'
        },
        pssh: new Uint8Array([181, 235, 45])
      }
    },
    endList: true,
    targetDuration: 1.9185833333333333,
    mediaSequence: 0,
    resolvedUri: '',
    segments: [{
      duration: 1.9185833333333333,
      map: {
        resolvedUri: 'https://www.example.com/720/init.m4f',
        uri: '720/init.m4f'
      },
      resolvedUri: 'https://www.example.com/720/0.m4f',
      timeline: 0,
      uri: '720/0.m4f',
      number: 0
    }, {
      duration: 1.9185833333333333,
      map: {
        resolvedUri: 'https://www.example.com/720/init.m4f',
        uri: '720/init.m4f'
      },
      resolvedUri: 'https://www.example.com/720/1.m4f',
      timeline: 0,
      uri: '720/1.m4f',
      number: 1
    }, {
      duration: 1.9185833333333333,
      map: {
        resolvedUri: 'https://www.example.com/720/init.m4f',
        uri: '720/init.m4f'
      },
      resolvedUri: 'https://www.example.com/720/2.m4f',
      timeline: 0,
      uri: '720/2.m4f',
      number: 2
    }, {
      duration: 0.24425000000000008,
      map: {
        resolvedUri: 'https://www.example.com/720/init.m4f',
        uri: '720/init.m4f'
      },
      resolvedUri: 'https://www.example.com/720/3.m4f',
      timeline: 0,
      uri: '720/3.m4f',
      number: 3
    }],
    timeline: 0,
    uri: ''
  }],
  segments: [],
  uri: ''
};
