import window from 'global/window';
import QUnit from 'qunit';
import {
  mapLegacyAvcCodecs,
  translateLegacyCodecs,
  parseCodecs,
  codecsFromDefault,
  isVideoCodec,
  isAudioCodec,
  muxerSupportsCodec,
  browserSupportsCodec,
  getMimeForCodec
} from '../src/codecs';

const supportedMuxerCodecs = [
  'mp4a',
  'avc1'
];

const unsupportedMuxerCodecs = [
  'hvc1',
  'ac-3',
  'ec-3',
  'mp3'
];

QUnit.module('Legacy Codecs');

QUnit.test('maps legacy AVC codecs', function(assert) {
  assert.equal(
    mapLegacyAvcCodecs('avc1.deadbeef'),
    'avc1.deadbeef',
    'does nothing for non legacy pattern'
  );
  assert.equal(
    mapLegacyAvcCodecs('avc1.dead.beef, mp4a.something'),
    'avc1.dead.beef, mp4a.something',
    'does nothing for non legacy pattern'
  );
  assert.equal(
    mapLegacyAvcCodecs('avc1.dead.beef,mp4a.something'),
    'avc1.dead.beef,mp4a.something',
    'does nothing for non legacy pattern'
  );
  assert.equal(
    mapLegacyAvcCodecs('mp4a.something,avc1.dead.beef'),
    'mp4a.something,avc1.dead.beef',
    'does nothing for non legacy pattern'
  );
  assert.equal(
    mapLegacyAvcCodecs('mp4a.something, avc1.dead.beef'),
    'mp4a.something, avc1.dead.beef',
    'does nothing for non legacy pattern'
  );
  assert.equal(
    mapLegacyAvcCodecs('avc1.42001e'),
    'avc1.42001e',
    'does nothing for non legacy pattern'
  );
  assert.equal(
    mapLegacyAvcCodecs('avc1.4d0020,mp4a.40.2'),
    'avc1.4d0020,mp4a.40.2',
    'does nothing for non legacy pattern'
  );
  assert.equal(
    mapLegacyAvcCodecs('mp4a.40.2,avc1.4d0020'),
    'mp4a.40.2,avc1.4d0020',
    'does nothing for non legacy pattern'
  );
  assert.equal(
    mapLegacyAvcCodecs('mp4a.40.40'),
    'mp4a.40.40',
    'does nothing for non video codecs'
  );

  assert.equal(
    mapLegacyAvcCodecs('avc1.66.30'),
    'avc1.42001e',
    'translates legacy video codec alone'
  );
  assert.equal(
    mapLegacyAvcCodecs('avc1.66.30, mp4a.40.2'),
    'avc1.42001e, mp4a.40.2',
    'translates legacy video codec when paired with audio'
  );
  assert.equal(
    mapLegacyAvcCodecs('mp4a.40.2, avc1.66.30'),
    'mp4a.40.2, avc1.42001e',
    'translates video codec when specified second'
  );
});

QUnit.test('translates legacy codecs', function(assert) {
  assert.deepEqual(
    translateLegacyCodecs(['avc1.66.30', 'avc1.66.30']),
    ['avc1.42001e', 'avc1.42001e'],
    'translates legacy avc1.66.30 codec'
  );

  assert.deepEqual(
    translateLegacyCodecs(['avc1.42C01E', 'avc1.42C01E']),
    ['avc1.42C01E', 'avc1.42C01E'],
    'does not translate modern codecs'
  );

  assert.deepEqual(
    translateLegacyCodecs(['avc1.42C01E', 'avc1.66.30']),
    ['avc1.42C01E', 'avc1.42001e'],
    'only translates legacy codecs when mixed'
  );

  assert.deepEqual(
    translateLegacyCodecs(['avc1.4d0020', 'avc1.100.41', 'avc1.77.41',
      'avc1.77.32', 'avc1.77.31', 'avc1.77.30',
      'avc1.66.30', 'avc1.66.21', 'avc1.42C01e']),
    ['avc1.4d0020', 'avc1.640029', 'avc1.4d0029',
      'avc1.4d0020', 'avc1.4d001f', 'avc1.4d001e',
      'avc1.42001e', 'avc1.420015', 'avc1.42C01e'],
    'translates a whole bunch'
  );
});

QUnit.module('parseCodecs');

QUnit.test('parses video only codec string', function(assert) {
  assert.deepEqual(
    parseCodecs('avc1.42001e'),
    {video: {type: 'avc1', details: '.42001e'}},
    'parsed video only codec string'
  );
});

QUnit.test('parses audio only codec string', function(assert) {
  assert.deepEqual(
    parseCodecs('mp4a.40.2'),
    {audio: {type: 'mp4a', details: '.40.2'}},
    'parsed audio only codec string'
  );
});

QUnit.test('parses video and audio codec string', function(assert) {
  assert.deepEqual(
    parseCodecs('avc1.42001e, mp4a.40.2'),
    {
      video: {type: 'avc1', details: '.42001e'},
      audio: {type: 'mp4a', details: '.40.2'}
    },
    'parsed video and audio codec string'
  );
});

QUnit.test('parses video and audio codec with mixed case', function(assert) {
  assert.deepEqual(
    parseCodecs('AvC1.42001E, Mp4A.40.E'),
    {
      video: {type: 'AvC1', details: '.42001E'},
      audio: {type: 'Mp4A', details: '.40.E'}
    },
    'parsed video and audio codec string'
  );
});

QUnit.module('codecsFromDefault');

QUnit.test('returns falsey when no audio group ID', function(assert) {
  assert.notOk(
    codecsFromDefault(
      { mediaGroups: { AUDIO: {} } },
      '',
    ),
    'returns falsey when no audio group ID'
  );
});

QUnit.test('returns falsey when no matching audio group', function(assert) {
  assert.notOk(
    codecsFromDefault(
      {
        mediaGroups: {
          AUDIO: {
            au1: {
              en: {
                default: false,
                playlists: [{
                  attributes: { CODECS: 'mp4a.40.2' }
                }]
              },
              es: {
                default: true,
                playlists: [{
                  attributes: { CODECS: 'mp4a.40.5' }
                }]
              }
            }
          }
        }
      },
      'au2'
    ),
    'returned falsey when no matching audio group'
  );
});

QUnit.test('returns falsey when no default for audio group', function(assert) {
  assert.notOk(
    codecsFromDefault(
      {
        mediaGroups: {
          AUDIO: {
            au1: {
              en: {
                default: false,
                playlists: [{
                  attributes: { CODECS: 'mp4a.40.2' }
                }]
              },
              es: {
                default: false,
                playlists: [{
                  attributes: { CODECS: 'mp4a.40.5' }
                }]
              }
            }
          }
        }
      },
      'au1'
    ),
    'returned falsey when no default for audio group'
  );
});

QUnit.test('returns audio profile for default in audio group', function(assert) {
  assert.deepEqual(
    codecsFromDefault(
      {
        mediaGroups: {
          AUDIO: {
            au1: {
              en: {
                default: false,
                playlists: [{
                  attributes: { CODECS: 'mp4a.40.2' }
                }]
              },
              es: {
                default: true,
                playlists: [{
                  attributes: { CODECS: 'mp4a.40.5' }
                }]
              }
            }
          }
        }
      },
      'au1'
    ),
    {audio: {type: 'mp4a', details: '.40.5'}},
    'returned parsed codec audio profile'
  );
});

QUnit.module('isVideoCodec');
QUnit.test('works as expected', function(assert) {
  [
    'av1',
    'avc01',
    'avc1',
    'avc02',
    'avc2',
    'vp09',
    'vp9',
    'vp8',
    'vp08',
    'hvc1',
    'hev1',
    'theora',
    'mp4v'
  ].forEach(function(codec) {
    assert.ok(isVideoCodec(codec), `"${codec}" is seen as a video codec`);
    assert.ok(isVideoCodec(` ${codec} `), `" ${codec} " is seen as video codec`);
    assert.ok(isVideoCodec(codec.toUpperCase()), `"${codec.toUpperCase()}" is seen as video codec`);
  });

  ['invalid', 'foo', 'mp4a', 'opus', 'vorbis'].forEach(function(codec) {
    assert.notOk(isVideoCodec(codec), `${codec} is not a video codec`);
  });

});

QUnit.module('isAudioCodec');
QUnit.test('works as expected', function(assert) {
  [
    'mp4a',
    'flac',
    'vorbis',
    'opus',
    'ac-3',
    'ac-4',
    'ec-3',
    'alac'
  ].forEach(function(codec) {
    assert.ok(isAudioCodec(codec), `"${codec}" is seen as an audio codec`);
    assert.ok(isAudioCodec(` ${codec} `), `" ${codec} " is seen as an audio codec`);
    assert.ok(isAudioCodec(codec.toUpperCase()), `"${codec.toUpperCase()}" is seen as an audio codec`);
  });

  ['invalid', 'foo', 'bar', 'avc1', 'av1'].forEach(function(codec) {
    assert.notOk(isAudioCodec(codec), `${codec} is not an audio codec`);
  });
});

QUnit.module('muxerSupportsCodec');
QUnit.test('works as expected', function(assert) {
  const validMuxerCodecs = [];
  const invalidMuxerCodecs = [];

  unsupportedMuxerCodecs.forEach(function(badCodec) {
    invalidMuxerCodecs.push(badCodec);
    supportedMuxerCodecs.forEach(function(goodCodec) {
      invalidMuxerCodecs.push(`${goodCodec}, ${badCodec}`);
    });
  });

  // generate all combinations of valid codecs
  supportedMuxerCodecs.forEach(function(codec, i) {
    validMuxerCodecs.push(codec);

    supportedMuxerCodecs.forEach(function(_codec, z) {
      if (z === i) {
        return;
      }
      validMuxerCodecs.push(`${codec}, ${_codec}`);
      validMuxerCodecs.push(`${codec},${_codec}`);
    });
  });

  validMuxerCodecs.forEach(function(codec) {
    assert.ok(muxerSupportsCodec(codec), `"${codec}" is supported`);
    assert.ok(muxerSupportsCodec(` ${codec} `), `" ${codec} " is supported`);
    assert.ok(muxerSupportsCodec(codec.toUpperCase()), `"${codec.toUpperCase()}" is supported`);
  });

  invalidMuxerCodecs.forEach(function(codec) {
    assert.notOk(muxerSupportsCodec(codec), `${codec} not supported`);
  });
});

QUnit.module('browserSupportsCodec', {
  beforeEach() {
    this.oldMediaSource = window.MediaSource;
  },
  afterEach() {
    window.MediaSource = this.oldMediaSource;
  }
});

QUnit.test('works as expected', function(assert) {
  window.MediaSource = {isTypeSupported: () => true};
  assert.ok(browserSupportsCodec('test'), 'isTypeSupported true, browser does support codec');

  window.MediaSource = {isTypeSupported: () => false};
  assert.notOk(browserSupportsCodec('test'), 'isTypeSupported false, browser does not support codec');

  window.MediaSource = null;
  assert.notOk(browserSupportsCodec('test'), 'no MediaSource, browser does not support codec');

  window.MediaSource = {isTypeSupported: null};
  assert.notOk(browserSupportsCodec('test'), 'no isTypeSupported, browser does not support codec');
});

QUnit.module('getMimeForCodec');

QUnit.test('works as expected', function(assert) {
  // mp4
  assert.equal(getMimeForCodec('vp9,mp4a'), 'video/mp4;codecs="vp9,mp4a"', 'mp4 video/audio works');
  assert.equal(getMimeForCodec('vp9'), 'video/mp4;codecs="vp9"', 'mp4 video works');
  assert.equal(getMimeForCodec('mp4a'), 'audio/mp4;codecs="mp4a"', 'mp4 audio works');

  // webm
  assert.equal(getMimeForCodec('vp8,opus'), 'video/webm;codecs="vp8,opus"', 'webm video/audio works');
  assert.equal(getMimeForCodec('vp8'), 'video/webm;codecs="vp8"', 'webm video works');
  assert.equal(getMimeForCodec('vorbis'), 'audio/webm;codecs="vorbis"', 'webm audio works');

  // ogg
  assert.equal(getMimeForCodec('theora,vorbis'), 'video/ogg;codecs="theora,vorbis"', 'ogg video/audio works');
  assert.equal(getMimeForCodec('theora'), 'video/ogg;codecs="theora"', 'ogg video works');
  // ogg will never be selected for audio only

  // mixed
  assert.equal(getMimeForCodec('opus'), 'audio/mp4;codecs="opus"', 'mp4 takes priority over everything');
  assert.equal(getMimeForCodec('vorbis'), 'audio/webm;codecs="vorbis"', 'webm takes priority over ogg');
  assert.equal(getMimeForCodec('foo'), 'video/mp4;codecs="foo"', 'mp4 is the default');

  assert.notOk(getMimeForCodec(), 'invalid codec returns undefined');

  assert.equal(getMimeForCodec('Mp4A.40.2,AvC1.42001E'), 'video/mp4;codecs="Mp4A.40.2,AvC1.42001E"', 'case is preserved');
});
