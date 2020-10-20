import {ParseStream, LineStream, Parser} from '../src';
import QUnit from 'qunit';
import sinon from 'sinon';
import testDataExpected from './dist/test-expected.js';
import testDataManifests from './dist/test-manifests.js';

QUnit.module('LineStream', {
  beforeEach() {
    this.lineStream = new LineStream();
  }
});
QUnit.test('empty inputs produce no tokens', function(assert) {
  let data = false;

  this.lineStream.on('data', function() {
    data = true;
  });
  this.lineStream.push('');
  assert.ok(!data, 'no tokens were produced');
});
QUnit.test('splits on newlines', function(assert) {
  const lines = [];

  this.lineStream.on('data', function(line) {
    lines.push(line);
  });
  this.lineStream.push('#EXTM3U\nmovie.ts\n');

  assert.strictEqual(2, lines.length, 'two lines are ready');
  assert.strictEqual('#EXTM3U', lines.shift(), 'the first line is the first token');
  assert.strictEqual('movie.ts', lines.shift(), 'the second line is the second token');
});
QUnit.test('empty lines become empty strings', function(assert) {
  const lines = [];

  this.lineStream.on('data', function(line) {
    lines.push(line);
  });
  this.lineStream.push('\n\n');

  assert.strictEqual(2, lines.length, 'two lines are ready');
  assert.strictEqual('', lines.shift(), 'the first line is empty');
  assert.strictEqual('', lines.shift(), 'the second line is empty');
});
QUnit.test('handles lines broken across appends', function(assert) {
  const lines = [];

  this.lineStream.on('data', function(line) {
    lines.push(line);
  });
  this.lineStream.push('#EXTM');
  assert.strictEqual(0, lines.length, 'no lines are ready');

  this.lineStream.push('3U\nmovie.ts\n');
  assert.strictEqual(2, lines.length, 'two lines are ready');
  assert.strictEqual('#EXTM3U', lines.shift(), 'the first line is the first token');
  assert.strictEqual('movie.ts', lines.shift(), 'the second line is the second token');
});
QUnit.test('stops sending events after deregistering', function(assert) {
  const temporaryLines = [];
  const temporary = function(line) {
    temporaryLines.push(line);
  };
  const permanentLines = [];
  const permanent = function(line) {
    permanentLines.push(line);
  };

  this.lineStream.on('data', temporary);
  this.lineStream.on('data', permanent);
  this.lineStream.push('line one\n');
  assert.strictEqual(
    temporaryLines.length,
    permanentLines.length,
    'both callbacks receive the event'
  );

  assert.ok(this.lineStream.off('data', temporary), 'a listener was removed');
  this.lineStream.push('line two\n');
  assert.strictEqual(1, temporaryLines.length, 'no new events are received');
  assert.strictEqual(2, permanentLines.length, 'new events are still received');
});

QUnit.module('ParseStream', {
  beforeEach() {
    this.lineStream = new LineStream();
    this.parseStream = new ParseStream();
    this.lineStream.pipe(this.parseStream);
  }
});

QUnit.test('parses custom tags', function(assert) {
  const manifest = '#VOD-STARTTIMESTAMP:1501533337573\n';
  let element;

  this.parseStream.addParser({
    expression: /^#VOD-STARTTIMESTAMP/,
    customType: 'startTimestamp'
  });

  this.parseStream.on('data', function(elem) {
    element = elem;
  });

  this.lineStream.push(manifest);
  assert.ok(element, 'element');
  assert.strictEqual(element.type, 'custom', 'the type of the data is custom');
  assert.strictEqual(
    element.customType,
    'startTimestamp',
    'the customType is startTimestamp'
  );
});

QUnit.test('mapper does not conflict with parser', function(assert) {
  const manifest = '#EXAMPLE\n';
  const commentMapper = sinon.spy(line => '#NEW-COMMENT');
  const commentMapper2 = sinon.spy(line => '#SOMETHING-ELSE');
  const dataCallback = sinon.spy();

  this.parseStream.addTagMapper({
    expression: /^#EXAMPLE/,
    map: commentMapper
  });
  this.parseStream.addTagMapper({
    expression: /^#EXAMPLE/,
    map: commentMapper2
  });

  this.parseStream.addParser({
    expression: /^#EXAMPLE/,
    customType: 'test'
  });
  this.parseStream.addParser({
    expression: /^#NEW-COMMENT/,
    customType: 'test2'
  });

  this.parseStream.on('data', dataCallback);
  this.lineStream.push(manifest);

  assert.ok(commentMapper.called);
  assert.ok(commentMapper2.called);
  assert.strictEqual(dataCallback.callCount, 3);

  assert.deepEqual(dataCallback.getCall(0).args[0], {
    data: '#EXAMPLE',
    type: 'custom',
    customType: 'test',
    segment: undefined
  });
  assert.deepEqual(dataCallback.getCall(1).args[0], {
    data: '#NEW-COMMENT',
    type: 'custom',
    customType: 'test2',
    segment: undefined
  });
  assert.deepEqual(dataCallback.getCall(2).args[0], {
    text: 'SOMETHING-ELSE',
    type: 'comment'
  });
});

QUnit.test('maps custom tags', function(assert) {
  const manifest = '#EXAMPLE\n';
  const commentMapper = sinon.spy(line => '#NEW-COMMENT');
  const dataCallback = sinon.spy();

  this.parseStream.addTagMapper({
    expression: /^#EXAMPLE/,
    map: commentMapper
  });

  this.parseStream.on('data', dataCallback);
  this.lineStream.push(manifest);

  assert.ok(commentMapper.called);
  assert.strictEqual(dataCallback.callCount, 2);

  assert.deepEqual(dataCallback.getCall(0).args[0], {
    text: 'EXAMPLE',
    type: 'comment'
  });
  assert.deepEqual(dataCallback.getCall(1).args[0], {
    text: 'NEW-COMMENT',
    type: 'comment'
  });
});

QUnit.test('maps multiple custom tags', function(assert) {
  const manifest = '#VOD-STARTTIMESTAMP:1501533337573\n';
  const commentMapper = sinon.spy(line => '#NEW-COMMENT');
  const pdtMapper = sinon.spy((line) => {
    const match = /#VOD-STARTTIMESTAMP:(\d+)/g.exec(line)[1];
    const ISOdate = new Date(Number(match)).toISOString();

    return `#EXT-X-PROGRAM-DATE-TIME:${ISOdate}`;
  });
  const dataCallback = sinon.spy();

  this.parseStream.addTagMapper({
    expression: /^#VOD-STARTTIMESTAMP/,
    map: commentMapper
  });
  this.parseStream.addTagMapper({
    expression: /^#VOD-STARTTIMESTAMP/,
    map: pdtMapper
  });

  this.parseStream.on('data', dataCallback);
  this.lineStream.push(manifest);

  assert.ok(commentMapper.called);
  assert.ok(pdtMapper.called);
  assert.strictEqual(dataCallback.callCount, 3);

  assert.deepEqual(dataCallback.getCall(0).args[0], {
    text: 'VOD-STARTTIMESTAMP:1501533337573',
    type: 'comment'
  });

  assert.deepEqual(dataCallback.getCall(1).args[0], {
    text: 'NEW-COMMENT',
    type: 'comment'
  });

  const dateTag = dataCallback.getCall(2).args[0];

  assert.strictEqual(dateTag.dateTimeString, '2017-07-31T20:35:37.573Z');
  assert.strictEqual(dateTag.tagType, 'program-date-time');
  assert.strictEqual(dateTag.type, 'tag');
});

QUnit.test('mapper ignores tags', function(assert) {
  const manifest = '#TAG\n';
  const dataCallback = sinon.spy();

  this.parseStream.addTagMapper({
    expression: /^#NO-MATCH/,
    map(line) {
      return '#MAPPED';
    }
  });

  this.parseStream.on('data', dataCallback);
  this.lineStream.push(manifest);

  assert.strictEqual(dataCallback.callCount, 1);
  assert.deepEqual(dataCallback.getCall(0).args[0], {
    text: 'TAG',
    type: 'comment'
  });
});

QUnit.test('parses comment lines', function(assert) {
  const manifest = '# a line that starts with a hash mark without "EXT" is a comment\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'comment', 'the type is comment');
  assert.strictEqual(
    element.text,
    manifest.slice(1, manifest.length - 1),
    'the comment text is parsed'
  );
});
QUnit.test('parses uri lines', function(assert) {
  const manifest = 'any non-blank line that does not start with a hash-mark is a URI\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'uri', 'the type is uri');
  assert.strictEqual(
    element.uri,
    manifest.substring(0, manifest.length - 1),
    'the uri text is parsed'
  );
});
QUnit.test('parses unknown tag types', function(assert) {
  const manifest = '#EXT-X-EXAMPLE-TAG:some,additional,stuff\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the type is tag');
  assert.strictEqual(
    element.data,
    manifest.slice(4, manifest.length - 1),
    'unknown tag data is preserved'
  );
});

// #EXTM3U
QUnit.test('parses #EXTM3U tags', function(assert) {
  const manifest = '#EXTM3U\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'm3u', 'the tag type is m3u');
});

// #EXTINF
QUnit.test('parses minimal #EXTINF tags', function(assert) {
  const manifest = '#EXTINF\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'inf', 'the tag type is inf');
});
QUnit.test('parses #EXTINF tags with durations', function(assert) {
  let manifest = '#EXTINF:15\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'inf', 'the tag type is inf');
  assert.strictEqual(element.duration, 15, 'the duration is parsed');
  assert.ok(!('title' in element), 'no title is parsed');

  manifest = '#EXTINF:21,\n';
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'inf', 'the tag type is inf');
  assert.strictEqual(element.duration, 21, 'the duration is parsed');
  assert.ok(!('title' in element), 'no title is parsed');
});
QUnit.test('parses #EXTINF tags with a duration and title', function(assert) {
  const manifest = '#EXTINF:13,Does anyone really use the title attribute?\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'inf', 'the tag type is inf');
  assert.strictEqual(element.duration, 13, 'the duration is parsed');
  assert.strictEqual(
    element.title,
    manifest.substring(manifest.indexOf(',') + 1, manifest.length - 1),
    'the title is parsed'
  );
});
QUnit.test('parses #EXTINF tags with carriage returns', function(assert) {
  const manifest = '#EXTINF:13,Does anyone really use the title attribute?\r\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'inf', 'the tag type is inf');
  assert.strictEqual(element.duration, 13, 'the duration is parsed');
  assert.strictEqual(
    element.title,
    manifest.substring(manifest.indexOf(',') + 1, manifest.length - 2),
    'the title is parsed'
  );
});

// #EXT-X-TARGETDURATION
QUnit.test('parses minimal #EXT-X-TARGETDURATION tags', function(assert) {
  const manifest = '#EXT-X-TARGETDURATION\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'targetduration', 'the tag type is targetduration');
  assert.ok(!('duration' in element), 'no duration is parsed');
});
QUnit.test('parses #EXT-X-TARGETDURATION with duration', function(assert) {
  const manifest = '#EXT-X-TARGETDURATION:47\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'targetduration', 'the tag type is targetduration');
  assert.strictEqual(element.duration, 47, 'the duration is parsed');
});

// #EXT-X-VERSION
QUnit.test('parses minimal #EXT-X-VERSION tags', function(assert) {
  const manifest = '#EXT-X-VERSION:\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'version', 'the tag type is version');
  assert.ok(!('version' in element), 'no version is present');
});
QUnit.test('parses #EXT-X-VERSION with a version', function(assert) {
  const manifest = '#EXT-X-VERSION:99\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'version', 'the tag type is version');
  assert.strictEqual(element.version, 99, 'the version is parsed');
});

// #EXT-X-MEDIA-SEQUENCE
QUnit.test('parses minimal #EXT-X-MEDIA-SEQUENCE tags', function(assert) {
  const manifest = '#EXT-X-MEDIA-SEQUENCE\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'media-sequence', 'the tag type is media-sequence');
  assert.ok(!('number' in element), 'no number is present');
});
QUnit.test('parses #EXT-X-MEDIA-SEQUENCE with sequence numbers', function(assert) {
  const manifest = '#EXT-X-MEDIA-SEQUENCE:109\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'media-sequence', 'the tag type is media-sequence');
  assert.ok(element.number, 109, 'the number is parsed');
});

// #EXT-X-PLAYLIST-TYPE
QUnit.test('parses minimal #EXT-X-PLAYLIST-TYPE tags', function(assert) {
  const manifest = '#EXT-X-PLAYLIST-TYPE:\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'playlist-type', 'the tag type is playlist-type');
  assert.ok(!('playlistType' in element), 'no playlist type is present');
});
QUnit.test('parses #EXT-X-PLAYLIST-TYPE with mutability info', function(assert) {
  let manifest = '#EXT-X-PLAYLIST-TYPE:EVENT\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'playlist-type', 'the tag type is playlist-type');
  assert.strictEqual(element.playlistType, 'EVENT', 'the playlist type is EVENT');

  manifest = '#EXT-X-PLAYLIST-TYPE:VOD\n';
  this.lineStream.push(manifest);
  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'playlist-type', 'the tag type is playlist-type');
  assert.strictEqual(element.playlistType, 'VOD', 'the playlist type is VOD');

  manifest = '#EXT-X-PLAYLIST-TYPE:nonsense\n';
  this.lineStream.push(manifest);
  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'playlist-type', 'the tag type is playlist-type');
  assert.strictEqual(element.playlistType, 'nonsense', 'the playlist type is parsed');
});

// #EXT-X-BYTERANGE
QUnit.test('parses minimal #EXT-X-BYTERANGE tags', function(assert) {
  const manifest = '#EXT-X-BYTERANGE\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'byterange', 'the tag type is byterange');
  assert.ok(!('length' in element), 'no length is present');
  assert.ok(!('offset' in element), 'no offset is present');
});
QUnit.test('parses #EXT-X-BYTERANGE with length and offset', function(assert) {
  let manifest = '#EXT-X-BYTERANGE:45\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'byterange', 'the tag type is byterange');
  assert.strictEqual(element.length, 45, 'length is parsed');
  assert.ok(!('offset' in element), 'no offset is present');

  manifest = '#EXT-X-BYTERANGE:108@16\n';
  this.lineStream.push(manifest);
  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'byterange', 'the tag type is byterange');
  assert.strictEqual(element.length, 108, 'length is parsed');
  assert.strictEqual(element.offset, 16, 'offset is parsed');
});

// #EXT-X-ALLOW-CACHE
QUnit.test('parses minimal #EXT-X-ALLOW-CACHE tags', function(assert) {
  const manifest = '#EXT-X-ALLOW-CACHE:\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'allow-cache', 'the tag type is allow-cache');
  assert.ok(!('allowed' in element), 'no allowed is present');
});
QUnit.test('parses valid #EXT-X-ALLOW-CACHE tags', function(assert) {
  let manifest = '#EXT-X-ALLOW-CACHE:YES\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'allow-cache', 'the tag type is allow-cache');
  assert.ok(element.allowed, 'allowed is parsed');

  manifest = '#EXT-X-ALLOW-CACHE:NO\n';
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'allow-cache', 'the tag type is allow-cache');
  assert.ok(!element.allowed, 'allowed is parsed');
});
// #EXT-X-MAP
QUnit.test('parses minimal #EXT-X-MAP tags', function(assert) {
  const manifest = '#EXT-X-MAP:URI="init.m4s"\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });

  this.lineStream.push(manifest);
  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'map', 'the tag type is map');
  assert.strictEqual(element.uri, 'init.m4s', 'parsed the uri');
});
QUnit.test('parses #EXT-X-MAP tags with a byterange', function(assert) {
  const manifest = '#EXT-X-MAP:URI="0.m4s", BYTERANGE="1000@23"\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });

  this.lineStream.push(manifest);
  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.uri, '0.m4s', 'parsed the uri');
  assert.strictEqual(
    element.byterange.length,
    1000,
    'parsed the byterange length'
  );
  assert.strictEqual(
    element.byterange.offset,
    23,
    'parsed the byterange offset'
  );
});
QUnit.test('parses #EXT-X-MAP tags with arbitrary attributes', function(assert) {
  const manifest = '#EXT-X-MAP:URI="init.mp4", SOMETHING=YES,BYTERANGE="720@0"\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });

  this.lineStream.push(manifest);
  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.uri, 'init.mp4', 'parsed the uri');
  assert.strictEqual(
    element.byterange.length,
    720,
    'parsed the byterange length'
  );
  assert.strictEqual(
    element.byterange.offset,
    0,
    'parsed the byterange offset'
  );
});
// #EXT-X-STREAM-INF
QUnit.test('parses minimal #EXT-X-STREAM-INF tags', function(assert) {
  const manifest = '#EXT-X-STREAM-INF\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'stream-inf', 'the tag type is stream-inf');
  assert.ok(!('attributes' in element), 'no attributes are present');
});
// #EXT-X-PROGRAM-DATE-TIME
QUnit.test('parses minimal EXT-X-PROGRAM-DATE-TIME tags', function(assert) {
  const manifest = '#EXT-X-PROGRAM-DATE-TIME\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'program-date-time', 'the tag type is date-time');
  assert.ok(!('dateTimeString' in element), 'no dateTime is present');
});
QUnit.test(
  'parses EXT-X-PROGRAM-DATE-TIME tags with valid date-time formats',
  function(assert) {
    let manifest = '#EXT-X-PROGRAM-DATE-TIME:2016-06-22T09:20:16.166-04:00\n';
    let element;

    this.parseStream.on('data', function(elem) {
      element = elem;
    });
    this.lineStream.push(manifest);

    assert.ok(element, 'an event was triggered');
    assert.strictEqual(element.type, 'tag', 'the line type is tag');
    assert.strictEqual(element.tagType, 'program-date-time', 'the tag type is date-time');
    assert.strictEqual(
      element.dateTimeString, '2016-06-22T09:20:16.166-04:00',
      'dateTimeString is parsed'
    );
    assert.deepEqual(
      element.dateTimeObject, new Date('2016-06-22T09:20:16.166-04:00'),
      'dateTimeObject is parsed'
    );

    manifest = '#EXT-X-PROGRAM-DATE-TIME:2016-06-22T09:20:16.16389Z\n';
    this.lineStream.push(manifest);

    assert.ok(element, 'an event was triggered');
    assert.strictEqual(element.type, 'tag', 'the line type is tag');
    assert.strictEqual(element.tagType, 'program-date-time', 'the tag type is date-time');
    assert.strictEqual(
      element.dateTimeString, '2016-06-22T09:20:16.16389Z',
      'dateTimeString is parsed'
    );
    assert.deepEqual(
      element.dateTimeObject, new Date('2016-06-22T09:20:16.16389Z'),
      'dateTimeObject is parsed'
    );
  }
);
QUnit.test('parses #EXT-X-STREAM-INF with common attributes', function(assert) {
  let manifest = '#EXT-X-STREAM-INF:BANDWIDTH=14400\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'stream-inf', 'the tag type is stream-inf');
  assert.strictEqual(element.attributes.BANDWIDTH, 14400, 'bandwidth is parsed');

  manifest = '#EXT-X-STREAM-INF:PROGRAM-ID=7\n';
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'stream-inf', 'the tag type is stream-inf');
  assert.strictEqual(element.attributes['PROGRAM-ID'], 7, 'program-id is parsed');

  manifest = '#EXT-X-STREAM-INF:RESOLUTION=396x224\n';
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'stream-inf', 'the tag type is stream-inf');
  assert.strictEqual(element.attributes.RESOLUTION.width, 396, 'width is parsed');
  assert.strictEqual(element.attributes.RESOLUTION.height, 224, 'heigth is parsed');

  manifest = '#EXT-X-STREAM-INF:CODECS="avc1.4d400d, mp4a.40.2"\n';
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'stream-inf', 'the tag type is stream-inf');
  assert.strictEqual(
    element.attributes.CODECS,
    'avc1.4d400d, mp4a.40.2',
    'codecs are parsed'
  );
});
QUnit.test('parses #EXT-X-STREAM-INF with arbitrary attributes', function(assert) {
  const manifest = '#EXT-X-STREAM-INF:NUMERIC=24,ALPHA=Value,MIXED=123abc\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'stream-inf', 'the tag type is stream-inf');
  assert.strictEqual(element.attributes.NUMERIC, '24', 'numeric attributes are parsed');
  assert.strictEqual(
    element.attributes.ALPHA,
    'Value',
    'alphabetic attributes are parsed'
  );
  assert.strictEqual(element.attributes.MIXED, '123abc', 'mixed attributes are parsed');
});
// #EXT-X-ENDLIST
QUnit.test('parses #EXT-X-ENDLIST tags', function(assert) {
  const manifest = '#EXT-X-ENDLIST\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'endlist', 'the tag type is stream-inf');
});

// #EXT-X-KEY
QUnit.test('parses valid #EXT-X-KEY tags', function(assert) {
  let manifest =
    '#EXT-X-KEY:METHOD=AES-128,URI="https://priv.example.com/key.php?r=52"\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.deepEqual(element, {
    type: 'tag',
    tagType: 'key',
    attributes: {
      METHOD: 'AES-128',
      URI: 'https://priv.example.com/key.php?r=52'
    }
  }, 'parsed a valid key');

  manifest = '#EXT-X-KEY:URI="https://example.com/key#1",METHOD=FutureType-1024\n';
  this.lineStream.push(manifest);
  assert.ok(element, 'an event was triggered');
  assert.deepEqual(element, {
    type: 'tag',
    tagType: 'key',
    attributes: {
      METHOD: 'FutureType-1024',
      URI: 'https://example.com/key#1'
    }
  }, 'parsed the attribute list independent of order');

  manifest = '#EXT-X-KEY:IV=1234567890abcdef1234567890abcdef\n';
  this.lineStream.push(manifest);
  assert.ok(element.attributes.IV, 'detected an IV attribute');
  assert.deepEqual(element.attributes.IV, new Uint32Array([
    0x12345678,
    0x90abcdef,
    0x12345678,
    0x90abcdef
  ]), 'parsed an IV value');
});

QUnit.test('parses minimal #EXT-X-KEY tags', function(assert) {
  const manifest = '#EXT-X-KEY:\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.deepEqual(element, {
    type: 'tag',
    tagType: 'key'
  }, 'parsed a minimal key tag');
});

QUnit.test('parses lightly-broken #EXT-X-KEY tags', function(assert) {
  let manifest = '#EXT-X-KEY:URI=\'https://example.com/single-quote\',METHOD=AES-128\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.strictEqual(
    element.attributes.URI,
    'https://example.com/single-quote',
    'parsed a single-quoted uri'
  );

  element = null;
  manifest = '#EXT-X-KEYURI="https://example.com/key",METHOD=AES-128\n';
  this.lineStream.push(manifest);
  assert.strictEqual(element.tagType, 'key', 'parsed the tag type');
  assert.strictEqual(
    element.attributes.URI,
    'https://example.com/key',
    'inferred a colon after the tag type'
  );

  element = null;
  manifest = '#EXT-X-KEY:  URI =  "https://example.com/key",METHOD=AES-128\n';
  this.lineStream.push(manifest);
  assert.strictEqual(
    element.attributes.URI,
    'https://example.com/key',
    'trims and removes quotes around the URI'
  );
});

QUnit.test('parses prefixed with 0x or 0X #EXT-X-KEY:IV tags', function(assert) {
  let manifest;
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });

  manifest = '#EXT-X-KEY:IV=0x1234567890abcdef1234567890abcdef\n';
  this.lineStream.push(manifest);
  assert.ok(element.attributes.IV, 'detected an IV attribute');
  assert.deepEqual(element.attributes.IV, new Uint32Array([
    0x12345678,
    0x90abcdef,
    0x12345678,
    0x90abcdef
  ]), 'parsed an IV value with 0x');

  manifest = '#EXT-X-KEY:IV=0X1234567890abcdef1234567890abcdef\n';
  this.lineStream.push(manifest);
  assert.ok(element.attributes.IV, 'detected an IV attribute');
  assert.deepEqual(element.attributes.IV, new Uint32Array([
    0x12345678,
    0x90abcdef,
    0x12345678,
    0x90abcdef
  ]), 'parsed an IV value with 0X');
});

// #EXT-X-START
QUnit.test('parses EXT-X-START tags', function(assert) {
  const manifest = '#EXT-X-START:TIME-OFFSET=1.1\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'start', 'the tag type is start');
  assert.strictEqual(element.attributes['TIME-OFFSET'], 1.1, 'parses time offset');
  assert.strictEqual(element.attributes.PRECISE, false, 'precise defaults to false');
});
QUnit.test('parses EXT-X-START PRECISE attribute', function(assert) {
  const manifest = '#EXT-X-START:TIME-OFFSET=1.4,PRECISE=YES\n';
  let element;

  this.parseStream.on('data', function(elem) {
    element = elem;
  });
  this.lineStream.push(manifest);

  assert.ok(element, 'an event was triggered');
  assert.strictEqual(element.type, 'tag', 'the line type is tag');
  assert.strictEqual(element.tagType, 'start', 'the tag type is start');
  assert.strictEqual(element.attributes['TIME-OFFSET'], 1.4, 'parses time offset');
  assert.strictEqual(element.attributes.PRECISE, true, 'parses precise attribute');
});
QUnit.test('flags missing EXT-X-START TIME-OFFSET attribute', function(assert) {
  const parser = new Parser();

  const manifest = [
    '#EXT-X-VERSION:3',
    '#EXT-X-TARGETDURATION:10',
    '#EXT-X-START:PRECISE=YES',
    '#EXTINF:10,',
    'media-00001.ts',
    '#EXT-X-ENDLIST'
  ].join('\n');
  let warning;

  parser.on('warn', function(warn) {
    warning = warn;
  });
  parser.push(manifest);

  assert.ok(warning, 'a warning was triggered');
  assert.ok((/ignoring start/).test(warning.message), 'message is about start tag');
  assert.strictEqual(typeof parser.manifest.start, 'undefined', 'does not parse start');
});

QUnit.test('ignores empty lines', function(assert) {
  const manifest = '\n';
  let event = false;

  this.parseStream.on('data', function() {
    event = true;
  });
  this.lineStream.push(manifest);

  assert.ok(!event, 'no event is triggered');
});

QUnit.module('m3u8 parser');

QUnit.test('can be constructed', function(assert) {
  assert.notStrictEqual(typeof new Parser(), 'undefined', 'parser is defined');
});

QUnit.test('can set custom parsers', function(assert) {
  const parser = new Parser();
  const manifest = [
    '#EXTM3U',
    '#EXT-X-VERSION:3',
    '#EXT-X-TARGETDURATION:10',
    '#EXT-X-MEDIA-SEQUENCE:0',
    '#EXT-X-PROGRAM-DATE-TIME:2017-07-31T20:35:35.053+00:00',
    '#VOD-STARTTIMESTAMP:1501533337573',
    '#VOD-TOTALDELETEDDURATION:0.0',
    '#VOD-FRAMERATE:29.97',
    ''
  ].join('\n');

  parser.addParser({
    expression: /^#VOD-STARTTIMESTAMP/,
    customType: 'startTimestamp'
  });
  parser.addParser({
    expression: /^#VOD-TOTALDELETEDDURATION/,
    customType: 'totalDeleteDuration'
  });
  parser.addParser({
    expression: /^#VOD-FRAMERATE/,
    customType: 'framerate',
    dataParser: (line) => (line.split(':')[1])
  });

  parser.push(manifest);
  assert.strictEqual(
    parser.manifest.custom.startTimestamp,
    '#VOD-STARTTIMESTAMP:1501533337573',
    'sets custom timestamp line'
  );

  assert.strictEqual(
    parser.manifest.custom.totalDeleteDuration,
    '#VOD-TOTALDELETEDDURATION:0.0',
    'sets custom delete duration'
  );

  assert.strictEqual(parser.manifest.custom.framerate, '29.97', 'sets framerate');
});

QUnit.test('segment level custom data', function(assert) {
  const parser = new Parser();

  const manifest = [
    '#EXTM3U',
    '#VOD-TIMING:1511816599485',
    '#COMMENT',
    '#EXTINF:8.0,',
    'ex1.ts',
    '#VOD-TIMING',
    '#EXTINF:8.0,',
    'ex2.ts',
    '#VOD-TIMING:1511816615485',
    '#EXT-UNKNOWN',
    '#EXTINF:8.0,',
    'ex3.ts',
    '#VOD-TIMING:1511816623485',
    '#EXTINF:8.0,',
    'ex3.ts',
    '#EXT-X-ENDLIST'
  ].join('\n');

  parser.addParser({
    expression: /^#VOD-TIMING/,
    customType: 'vodTiming',
    segment: true
  });

  parser.push(manifest);
  assert.equal(
    parser.manifest.segments[0].custom.vodTiming,
    '#VOD-TIMING:1511816599485',
    'parser attached segment level custom data'
  );
  assert.equal(
    parser.manifest.segments[1].custom.vodTiming,
    '#VOD-TIMING',
    'parser got segment level custom data without :'
  );
});

QUnit.test('attaches cue-out data to segment', function(assert) {
  const parser = new Parser();

  const manifest = [
    '#EXTM3U',
    '#EXTINF:5,',
    '#COMMENT',
    'ex1.ts',
    '#EXT-X-CUE-OUT:10',
    '#EXTINF:5,',
    'ex2.ts',
    '#EXT-X-CUE-OUT15',
    '#EXT-UKNOWN-TAG',
    '#EXTINF:5,',
    'ex3.ts',
    '#EXT-X-CUE-OUT',
    '#EXTINF:5,',
    'ex3.ts',
    '#EXT-X-ENDLIST'
  ].join('\n');

  parser.push(manifest);

  assert.equal(parser.manifest.segments[1].cueOut, '10', 'parser attached cue out tag');
  assert.equal(parser.manifest.segments[2].cueOut, '15', 'cue out without : seperator');
  assert.equal(parser.manifest.segments[3].cueOut, '', 'cue out without data');
});

QUnit.test('attaches cue-out-cont data to segment', function(assert) {
  const parser = new Parser();

  const manifest = [
    '#EXTM3U',
    '#EXTINF:5,',
    '#COMMENT',
    'ex1.ts',
    '#EXT-X-CUE-OUT-CONT:10/60',
    '#EXTINF:5,',
    'ex2.ts',
    '#EXT-X-CUE-OUT-CONT15/30',
    '#EXT-UKNOWN-TAG',
    '#EXTINF:5,',
    'ex3.ts',
    '#EXT-X-CUE-OUT-CONT',
    '#EXTINF:5,',
    'ex3.ts',
    '#EXT-X-ENDLIST'
  ].join('\n');

  parser.push(manifest);

  assert.equal(
    parser.manifest.segments[1].cueOutCont, '10/60',
    'parser attached cue out cont tag'
  );
  assert.equal(
    parser.manifest.segments[2].cueOutCont, '15/30',
    'cue out cont without : seperator'
  );
  assert.equal(parser.manifest.segments[3].cueOutCont, '', 'cue out cont without data');
});

QUnit.test('attaches cue-in data to segment', function(assert) {
  const parser = new Parser();

  const manifest = [
    '#EXTM3U',
    '#EXTINF:5,',
    '#COMMENT',
    'ex1.ts',
    '#EXT-X-CUE-IN',
    '#EXTINF:5,',
    'ex2.ts',
    '#EXT-X-CUE-IN:15',
    '#EXT-UKNOWN-TAG',
    '#EXTINF:5,',
    'ex3.ts',
    '#EXT-X-CUE-IN=abc',
    '#EXTINF:5,',
    'ex3.ts',
    '#EXT-X-ENDLIST'
  ].join('\n');

  parser.push(manifest);

  assert.equal(parser.manifest.segments[1].cueIn, '', 'parser attached cue in tag');
  assert.equal(parser.manifest.segments[2].cueIn, '15', 'cue in with data');
  assert.equal(
    parser.manifest.segments[3].cueIn, '=abc',
    'cue in without colon seperator'
  );
});

QUnit.test('parses characteristics attribute', function(assert) {
  const parser = new Parser();

  const manifest = [
    '#EXTM3U',
    '#EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID="subs",CHARACTERISTICS="char",NAME="test"',
    '#EXT-X-STREAM-INF:BANDWIDTH=1,CODECS="mp4a.40.2, avc1.4d400d",SUBTITLES="subs"',
    'index.m3u8'
  ].join('\n');

  parser.push(manifest);

  assert.equal(
    parser.manifest.mediaGroups.SUBTITLES.subs.test.characteristics,
    'char',
    'parsed CHARACTERISTICS attribute'
  );
});

QUnit.test('parses FORCED attribute', function(assert) {
  const parser = new Parser();

  const manifest = [
    '#EXTM3U',
    '#EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID="subs",CHARACTERISTICS="char",NAME="test",FORCED=YES',
    '#EXT-X-STREAM-INF:BANDWIDTH=1,CODECS="mp4a.40.2, avc1.4d400d",SUBTITLES="subs"',
    'index.m3u8'
  ].join('\n');

  parser.push(manifest);

  assert.ok(
    parser.manifest.mediaGroups.SUBTITLES.subs.test.forced,
    'parsed FORCED attribute'
  );
});

QUnit.test('parses Widevine #EXT-X-KEY attributes and attaches to manifest', function(assert) {
  const parser = new Parser();

  const manifest = [
    '#EXTM3U',
    '#EXT-X-KEY:METHOD=SAMPLE-AES-CTR,' +
      'URI="data:text/plain;base64,AAAAPnBzc2gAAAAA7e+LqXnWSs6jyCfc1R0h7QAAAB4iFnN' +
      'oYWthX2NlYzJmNjRhYTc4OTBhMTFI49yVmwY=",KEYID=0x800AACAA522958AE888062B5695DB6BF,' +
      'KEYFORMATVERSIONS="1",KEYFORMAT="urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed"',
    '#EXTINF:5,',
    'ex1.ts',
    '#EXT-X-ENDLIST'
  ].join('\n');

  parser.push(manifest);

  assert.ok(parser.manifest.contentProtection, 'contentProtection property added');
  assert.equal(
    parser.manifest.contentProtection['com.widevine.alpha'].attributes.schemeIdUri,
    'urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed',
    'schemeIdUri set correctly'
  );
  assert.equal(
    parser.manifest.contentProtection['com.widevine.alpha'].attributes.keyId,
    '800AACAA522958AE888062B5695DB6BF',
    'keyId set correctly'
  );
  assert.equal(
    parser.manifest.contentProtection['com.widevine.alpha'].pssh.byteLength,
    62,
    'base64 URI decoded to TypedArray'
  );
});

QUnit.test('Widevine #EXT-X-KEY attributes not attached to manifest if METHOD is invalid', function(assert) {
  const parser = new Parser();

  const manifest = [
    '#EXTM3U',
    '#EXT-X-KEY:METHOD=NONE,' +
      'URI="data:text/plain;base64,AAAAPnBzc2gAAAAA7e+LqXnWSs6jyCfc1R0h7QAAAB4iFnN' +
      'oYWthX2NlYzJmNjRhYTc4OTBhMTFI49yVmwY=",KEYID=0x800AACAA522958AE888062B5695DB6BF,' +
      'KEYFORMATVERSIONS="1",KEYFORMAT="urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed"',
    '#EXTINF:5,',
    'ex1.ts',
    '#EXT-X-ENDLIST'
  ].join('\n');

  parser.push(manifest);

  assert.notOk(parser.manifest.contentProtection, 'contentProtection not added');
});

QUnit.test('Widevine #EXT-X-KEY attributes not attached to manifest if URI is invalid', function(assert) {
  const parser = new Parser();

  const manifest = [
    '#EXTM3U',
    '#EXT-X-KEY:METHOD=SAMPLE-AES-CTR,' +
      'URI="AAAAPnBzc2gAAAAA7e+LqXnWSs6jyCfc1R0h7QAAAB4iFnN' +
      'oYWthX2NlYzJmNjRhYTc4OTBhMTFI49yVmwY=",KEYID=0x800AACAA522958AE888062B5695DB6BF,' +
      'KEYFORMATVERSIONS="1",KEYFORMAT="urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed"',
    '#EXTINF:5,',
    'ex1.ts',
    '#EXT-X-ENDLIST'
  ].join('\n');

  parser.push(manifest);

  assert.notOk(parser.manifest.contentProtection, 'contentProtection not added');
});

QUnit.test('Widevine #EXT-X-KEY attributes not attached to manifest if KEYID is invalid', function(assert) {
  const parser = new Parser();

  const manifest = [
    '#EXTM3U',
    '#EXT-X-KEY:METHOD=SAMPLE-AES-CTR,' +
      'URI="data:text/plain;base64,AAAAPnBzc2gAAAAA7e+LqXnWSs6jyCfc1R0h7QAAAB4iFnN' +
      'oYWthX2NlYzJmNjRhYTc4OTBhMTFI49yVmwY=",KEYID=800AACAA522958AE888062B5695DB6BF,' +
      'KEYFORMATVERSIONS="1",KEYFORMAT="urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed"',
    '#EXTINF:5,',
    'ex1.ts',
    '#EXT-X-ENDLIST'
  ].join('\n');

  parser.push(manifest);

  assert.notOk(parser.manifest.contentProtection, 'contentProtection not added');
});

QUnit.test('Widevine #EXT-X-KEY attributes not attached to manifest if KEYFORMAT is not Widevine UUID', function(assert) {
  const parser = new Parser();

  const manifest = [
    '#EXTM3U',
    '#EXT-X-KEY:METHOD=SAMPLE-AES-CTR,' +
      'URI="data:text/plain;base64,AAAAPnBzc2gAAAAA7e+LqXnWSs6jyCfc1R0h7QAAAB4iFnN' +
      'oYWthX2NlYzJmNjRhYTc4OTBhMTFI49yVmwY=",KEYID=0x800AACAA522958AE888062B5695DB6BF,' +
      'KEYFORMATVERSIONS="1",KEYFORMAT="invalid-keyformat"',
    '#EXTINF:5,',
    'ex1.ts',
    '#EXT-X-ENDLIST'
  ].join('\n');

  parser.push(manifest);

  assert.notOk(parser.manifest.contentProtection, 'contentProtection not added');
});

QUnit.module('m3u8s');

QUnit.test('parses static manifests as expected', function(assert) {
  let key;

  for (key in testDataManifests) {
    if (testDataExpected[key]) {
      const parser = new Parser();

      parser.push(testDataManifests[key]);
      assert.deepEqual(
        parser.manifest,
        testDataExpected[key],
        key + '.m3u8 was parsed correctly'
      );
    }
  }
});
