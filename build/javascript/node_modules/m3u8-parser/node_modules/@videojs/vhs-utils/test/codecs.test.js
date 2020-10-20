import QUnit from 'qunit';
import {
  mapLegacyAvcCodecs,
  translateLegacyCodecs,
  parseCodecs,
  audioProfileFromDefault
} from '../src/codecs';

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
    {
      codecCount: 1,
      videoCodec: 'avc1',
      videoObjectTypeIndicator: '.42001e',
      audioProfile: null
    },
    'parsed video only codec string'
  );
});

QUnit.test('parses audio only codec string', function(assert) {
  assert.deepEqual(
    parseCodecs('mp4a.40.2'),
    {
      codecCount: 1,
      audioProfile: '2'
    },
    'parsed audio only codec string'
  );
});

QUnit.test('parses video and audio codec string', function(assert) {
  assert.deepEqual(
    parseCodecs('avc1.42001e, mp4a.40.2'),
    {
      codecCount: 2,
      videoCodec: 'avc1',
      videoObjectTypeIndicator: '.42001e',
      audioProfile: '2'
    },
    'parsed video and audio codec string'
  );
});

QUnit.module('audioProfileFromDefault');

QUnit.test('returns falsey when no audio group ID', function(assert) {
  assert.notOk(
    audioProfileFromDefault(
      { mediaGroups: { AUDIO: {} } },
      '',
    ),
    'returns falsey when no audio group ID'
  );
});

QUnit.test('returns falsey when no matching audio group', function(assert) {
  assert.notOk(
    audioProfileFromDefault(
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
    audioProfileFromDefault(
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
    audioProfileFromDefault(
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
    '5',
    'returned parsed codec audio profile'
  );
});
