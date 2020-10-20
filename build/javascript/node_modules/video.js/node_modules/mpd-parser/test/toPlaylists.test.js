import {
  toPlaylists
} from '../src/toPlaylists';
import QUnit from 'qunit';

QUnit.module('toPlaylists');

QUnit.test('no representations', function(assert) {
  assert.deepEqual(toPlaylists([]), []);
});

QUnit.test('pretty simple', function(assert) {
  const representations = [{
    attributes: { baseUrl: 'http://example.com/', periodIndex: 0, sourceDuration: 2 },
    segmentInfo: {
      template: { }
    }
  }];

  const playlists = [{
    attributes: {
      baseUrl: 'http://example.com/',
      periodIndex: 0,
      sourceDuration: 2,
      duration: 2
    },
    segments: [{
      uri: '',
      timeline: 0,
      duration: 2,
      resolvedUri: 'http://example.com/',
      map: {
        uri: '',
        resolvedUri: 'http://example.com/'
      },
      number: 1
    }]
  }];

  assert.deepEqual(toPlaylists(representations), playlists);
});

QUnit.test('segment base', function(assert) {
  const representations = [{
    attributes: { baseUrl: 'http://example.com/', periodIndex: 0, sourceDuration: 2 },
    segmentInfo: {
      base: true
    }
  }];

  const playlists = [{
    attributes: {
      baseUrl: 'http://example.com/',
      periodIndex: 0,
      sourceDuration: 2,
      duration: 2
    },
    segments: [{
      map: {
        resolvedUri: 'http://example.com/',
        uri: ''
      },
      resolvedUri: 'http://example.com/',
      uri: 'http://example.com/',
      timeline: 0,
      duration: 2,
      number: 0
    }]
  }];

  assert.deepEqual(toPlaylists(representations), playlists);
});

QUnit.test('segment base with sidx', function(assert) {
  const representations = [{
    attributes: {
      baseUrl: 'http://example.com/',
      periodIndex: 0,
      sourceDuration: 2,
      indexRange: '10-19'
    },
    segmentInfo: {
      base: true
    }
  }];

  const playlists = [{
    attributes: {
      baseUrl: 'http://example.com/',
      periodIndex: 0,
      sourceDuration: 2,
      duration: 2,
      indexRange: '10-19'
    },
    segments: [],
    sidx: {
      map: {
        resolvedUri: 'http://example.com/',
        uri: ''
      },
      resolvedUri: 'http://example.com/',
      uri: 'http://example.com/',
      byterange: {
        offset: 10,
        length: 10
      },
      timeline: 0,
      duration: 2,
      number: 0
    }
  }];

  assert.deepEqual(toPlaylists(representations), playlists);
});

QUnit.test('segment list', function(assert) {
  const representations = [{
    attributes: {
      baseUrl: 'http://example.com/',
      duration: 10,
      sourceDuration: 11,
      periodIndex: 0
    },
    segmentInfo: {
      list: {
        segmentUrls: [{
          media: '1.fmp4'
        }, {
          media: '2.fmp4'
        }]
      }
    }
  }];

  const playlists = [{
    attributes: {
      baseUrl: 'http://example.com/',
      duration: 10,
      sourceDuration: 11,
      segmentUrls: [{
        media: '1.fmp4'
      }, {
        media: '2.fmp4'
      }],
      periodIndex: 0
    },
    segments: [{
      duration: 10,
      map: {
        resolvedUri: 'http://example.com/',
        uri: ''
      },
      resolvedUri: 'http://example.com/1.fmp4',
      timeline: 0,
      uri: '1.fmp4',
      number: 1
    }, {
      duration: 1,
      map: {
        resolvedUri: 'http://example.com/',
        uri: ''
      },
      resolvedUri: 'http://example.com/2.fmp4',
      timeline: 0,
      uri: '2.fmp4',
      number: 2
    }]
  }];

  assert.deepEqual(toPlaylists(representations), playlists);
});
