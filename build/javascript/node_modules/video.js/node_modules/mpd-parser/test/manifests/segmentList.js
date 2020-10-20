export const parsedManifest = {
  allowCache: true,
  discontinuityStarts: [],
  duration: 6,
  endList: true,
  minimumUpdatePeriod: 0,
  mediaGroups: {
    'AUDIO': {},
    'CLOSED-CAPTIONS': {},
    'SUBTITLES': {},
    'VIDEO': {}
  },
  playlists: [
    {
      attributes: {
        'AUDIO': 'audio',
        'BANDWIDTH': 449000,
        'CODECS': 'avc1.420015',
        'NAME': '482',
        'PROGRAM-ID': 1,
        'RESOLUTION': {
          height: 270,
          width: 482
        },
        'SUBTITLES': 'subs'
      },
      endList: true,
      mediaSequence: 1,
      targetDuration: 1,
      resolvedUri: '',
      segments: [
        {
          duration: 1,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/low/segment-1.ts',
          timeline: 0,
          uri: 'low/segment-1.ts',
          number: 1
        },
        {
          duration: 1,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/low/segment-2.ts',
          timeline: 0,
          uri: 'low/segment-2.ts',
          number: 2
        },
        {
          duration: 1,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/low/segment-3.ts',
          timeline: 0,
          uri: 'low/segment-3.ts',
          number: 3
        },
        {
          duration: 1,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/low/segment-4.ts',
          timeline: 0,
          uri: 'low/segment-4.ts',
          number: 4
        },
        {
          duration: 1,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/low/segment-5.ts',
          timeline: 0,
          uri: 'low/segment-5.ts',
          number: 5
        },
        {
          duration: 1,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/low/segment-6.ts',
          timeline: 0,
          uri: 'low/segment-6.ts',
          number: 6
        }
      ],
      timeline: 0,
      uri: ''
    },
    {
      attributes: {
        'AUDIO': 'audio',
        'BANDWIDTH': 3971000,
        'CODECS': 'avc1.420015',
        'NAME': '720',
        'PROGRAM-ID': 1,
        'RESOLUTION': {
          height: 404,
          width: 720
        },
        'SUBTITLES': 'subs'
      },
      endList: true,
      resolvedUri: '',
      mediaSequence: 1,
      targetDuration: 60,
      segments: [
        {
          duration: 60,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/high/segment-1.ts',
          timeline: 0,
          uri: 'high/segment-1.ts',
          number: 1
        },
        {
          duration: 60,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/high/segment-2.ts',
          timeline: 0,
          uri: 'high/segment-2.ts',
          number: 2
        },
        {
          duration: 60,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/high/segment-3.ts',
          timeline: 0,
          uri: 'high/segment-3.ts',
          number: 3
        },
        {
          duration: 60,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/high/segment-4.ts',
          timeline: 0,
          uri: 'high/segment-4.ts',
          number: 4
        },
        {
          duration: 60,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/high/segment-5.ts',
          timeline: 0,
          uri: 'high/segment-5.ts',
          number: 5
        },
        {
          duration: 60,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/high/segment-6.ts',
          timeline: 0,
          uri: 'high/segment-6.ts',
          number: 6
        },
        {
          duration: 60,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/high/segment-7.ts',
          timeline: 0,
          uri: 'high/segment-7.ts',
          number: 7
        },
        {
          duration: 60,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/high/segment-8.ts',
          timeline: 0,
          uri: 'high/segment-8.ts',
          number: 8
        },
        {
          duration: 60,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/high/segment-9.ts',
          timeline: 0,
          uri: 'high/segment-9.ts',
          number: 9
        },
        {
          duration: 60,
          map: {
            uri: '',
            resolvedUri: 'https://www.example.com/base'
          },
          resolvedUri: 'https://www.example.com/high/segment-10.ts',
          timeline: 0,
          uri: 'high/segment-10.ts',
          number: 10
        }
      ],
      timeline: 0,
      uri: ''
    }
  ],
  segments: [],
  uri: ''
};
