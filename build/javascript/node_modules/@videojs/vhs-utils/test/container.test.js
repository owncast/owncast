import QUnit from 'qunit';
import {detectContainerForBytes, isLikelyFmp4MediaSegment} from '../src/containers.js';
import {stringToBytes} from '../src/byte-helpers.js';

const fillerArray = (size) => Array.apply(null, Array(size)).map(() => 0x00);
const otherMp4Data = [0x00, 0x00, 0x00, 0x00].concat(stringToBytes('stypiso'));
const id3Data = []
  // id3 header is 10 bytes without footer
  // 10th byte is length 0x23 or 35 in decimal
  // so a total length of 45
  .concat(stringToBytes('ID3').concat([0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x23]))
  // add in the id3 content
  .concat(Array.apply(null, Array(35)).map(() => 0x00));

const id3DataWithFooter = []
  // id3 header is 20 bytes with footer
  // "we have a footer" is the sixth byte
  // 10th byte is length of 0x23 or 35 in decimal
  // so a total length of 55
  .concat(stringToBytes('ID3').concat([0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x23]))
  // add in the id3 content
  .concat(Array.apply(null, Array(45)).map(() => 0x00));

const testData = {
  'webm': [0x1A, 0x45, 0xDf, 0xA3],
  'flac': stringToBytes('fLaC'),
  'ogg': stringToBytes('OggS'),
  'aac': [0xFF, 0xF1],
  'mp3': [0xFF, 0xFB],
  '3gp': [0x00, 0x00, 0x00, 0x00].concat(stringToBytes('ftyp3g')),
  'mp4': [0x00, 0x00, 0x00, 0x00].concat(stringToBytes('ftypiso')),
  'ts': [0x47]
};

QUnit.module('detectContainerForBytes');

QUnit.test('should identify known types', function(assert) {
  Object.keys(testData).forEach(function(key) {
    const data = new Uint8Array(testData[key]);

    assert.equal(detectContainerForBytes(testData[key]), key, `found ${key} with Array`);
    assert.equal(detectContainerForBytes(data.buffer), key, `found ${key} with ArrayBuffer`);
    assert.equal(detectContainerForBytes(data), key, `found ${key} with Uint8Array`);
  });

  const mp4Bytes = new Uint8Array([0x00, 0x00, 0x00, 0x00].concat(stringToBytes('styp')));

  assert.equal(detectContainerForBytes(mp4Bytes), 'mp4', 'styp mp4 detected as mp4');

  // mp3 and aac audio can have id3 data before the
  // signature for the file, so we need to handle that.
  ['mp3', 'aac'].forEach(function(type) {
    const dataWithId3 = new Uint8Array([].concat(id3Data).concat(testData[type]));
    const dataWithId3Footer = new Uint8Array([].concat(id3DataWithFooter).concat(testData[type]));

    const recursiveDataWithId3 = new Uint8Array([]
      .concat(id3Data)
      .concat(id3Data)
      .concat(id3Data)
      .concat(testData[type]));
    const recursiveDataWithId3Footer = new Uint8Array([]
      .concat(id3DataWithFooter)
      .concat(id3DataWithFooter)
      .concat(id3DataWithFooter)
      .concat(testData[type]));

    const differentId3Sections = new Uint8Array([]
      .concat(id3DataWithFooter)
      .concat(id3Data)
      .concat(id3DataWithFooter)
      .concat(id3Data)
      .concat(testData[type]));

    assert.equal(detectContainerForBytes(dataWithId3), type, `id3 skipped and ${type} detected`);
    assert.equal(detectContainerForBytes(dataWithId3Footer), type, `id3 + footer skipped and ${type} detected`);
    assert.equal(detectContainerForBytes(recursiveDataWithId3), type, `id3 x3 skipped and ${type} detected`);
    assert.equal(detectContainerForBytes(recursiveDataWithId3Footer), type, `id3 + footer x3 skipped and ${type} detected`);
    assert.equal(detectContainerForBytes(differentId3Sections), type, `id3 with/without footer skipped and ${type} detected`);
  });

  const notTs = []
    .concat(testData.ts)
    .concat(fillerArray(188));
  const longTs = []
    .concat(testData.ts)
    .concat(fillerArray(187))
    .concat(testData.ts);

  const unsyncTs = []
    .concat(fillerArray(187))
    .concat(testData.ts)
    .concat(fillerArray(187))
    .concat(testData.ts);

  const badTs = []
    .concat(fillerArray(188))
    .concat(testData.ts)
    .concat(fillerArray(187))
    .concat(testData.ts);

  assert.equal(detectContainerForBytes(longTs), 'ts', 'long ts data is detected');
  assert.equal(detectContainerForBytes(unsyncTs), 'ts', 'unsynced ts is detected');
  assert.equal(detectContainerForBytes(badTs), '', 'ts without a sync byte in 188 bytes is not detected');
  assert.equal(detectContainerForBytes(notTs), '', 'ts missing 0x47 at 188 is not ts at all');
  assert.equal(detectContainerForBytes(otherMp4Data), 'mp4', 'fmp4 detected as mp4');
  assert.equal(detectContainerForBytes(new Uint8Array()), '', 'no type');
  assert.equal(detectContainerForBytes(), '', 'no type');
});

const createBox = function(type) {
  const size = 0x20;

  // size bytes
  return [0x00, 0x00, 0x00, size]
    // box identfier styp
    .concat(stringToBytes(type))
    // filler data for size minus identfier and size bytes
    .concat(fillerArray(size - 8));
};

QUnit.module('isLikelyFmp4MediaSegment');
QUnit.test('works as expected', function(assert) {
  const fmp4Data = []
    .concat(createBox('styp'))
    .concat(createBox('sidx'))
    .concat(createBox('moof'));

  const mp4Data = []
    .concat(createBox('ftyp'))
    .concat(createBox('sidx'))
    .concat(createBox('moov'));

  const fmp4Fake = []
    .concat(createBox('test'))
    .concat(createBox('moof'))
    .concat(createBox('fooo'))
    .concat(createBox('bar'));

  assert.ok(isLikelyFmp4MediaSegment(fmp4Data), 'fmp4 is recognized as fmp4');
  assert.ok(isLikelyFmp4MediaSegment(fmp4Fake), 'fmp4 with moof and unknown boxes is still fmp4');
  assert.ok(isLikelyFmp4MediaSegment(createBox('moof')), 'moof alone is recognized as fmp4');
  assert.notOk(isLikelyFmp4MediaSegment(mp4Data), 'mp4 is not recognized');
  assert.notOk(isLikelyFmp4MediaSegment([].concat(id3DataWithFooter).concat(testData.mp3)), 'bad data is not recognized');
  assert.notOk(isLikelyFmp4MediaSegment(new Uint8Array()), 'no errors on empty data');
  assert.notOk(isLikelyFmp4MediaSegment(), 'no errors on empty data');
});
