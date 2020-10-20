import QUnit from 'qunit';
import {
  bytesToString,
  stringToBytes,
  toUint8,
  concatTypedArrays
} from '../src/byte-helpers.js';
import window from 'global/window';

const arrayNames = [];

[
  'Array',
  'Int8Array',
  'Uint8Array',
  'Uint8ClampedArray',
  'Int16Array',
  'Uint16Array',
  'Int32Array',
  'Uint32Array',
  'Float32Array',
  'Float64Array'
].forEach(function(name) {
  if (window[name]) {
    arrayNames.push(name);
  }
});

QUnit.module('bytesToString');

const testString = 'hello竜';
const testBytes = [
  // h
  0x68,
  // e
  0x65,
  // l
  0x6c,
  // l
  0x6c,
  // o
  0x6f,
  // 竜
  0xe7, 0xab, 0x9c
];

const rawBytes = [0x47, 0x40, 0x00, 0x10, 0x00, 0x00, 0xb0, 0x0d, 0x00, 0x01];

QUnit.test('should function as expected', function(assert) {
  arrayNames.forEach(function(name) {
    const testObj = name === 'Array' ? testBytes : new window[name](testBytes);

    assert.equal(bytesToString(testObj), testString, `testString work as a string arg with ${name}`);
    assert.equal(bytesToString(new window[name]()), '', `empty ${name} returns empty string`);
  });

  assert.equal(bytesToString(), '', 'undefined returns empty string');
  assert.equal(bytesToString(null), '', 'null returns empty string');
  assert.equal(bytesToString(stringToBytes(testString)), testString, 'stringToBytes -> bytesToString works');
});

QUnit.module('stringToBytes');

QUnit.test('should function as expected', function(assert) {
  assert.deepEqual(stringToBytes(testString), testBytes, 'returns an array of bytes');
  assert.deepEqual(stringToBytes(), [], 'empty array for undefined');
  assert.deepEqual(stringToBytes(null), [], 'empty array for null');
  assert.deepEqual(stringToBytes(''), [], 'empty array for empty string');
  assert.deepEqual(stringToBytes(10), [0x31, 0x30], 'converts numbers to strings');
  assert.deepEqual(stringToBytes(bytesToString(testBytes)), testBytes, 'bytesToString -> stringToBytes works');
  assert.deepEqual(stringToBytes(bytesToString(rawBytes), true), rawBytes, 'equal to original with raw bytes mode');
  assert.notDeepEqual(stringToBytes(bytesToString(rawBytes)), rawBytes, 'without raw byte mode works, not equal');
});

QUnit.module('toUint8');

QUnit.test('should function as expected', function(assert) {
  const undef = toUint8();

  assert.ok(undef instanceof Uint8Array && undef.length === 0, 'undef is a blank Uint8Array');

  const nul = toUint8(null);

  assert.ok(nul instanceof Uint8Array && nul.length === 0, 'undef is a blank Uint8Array');

  arrayNames.forEach(function(name) {
    const testObj = name === 'Array' ? testBytes : new window[name](testBytes);
    const uint = toUint8(testObj);

    assert.ok(uint instanceof Uint8Array && uint.length > 0, `converted ${name} to Uint8Array`);
  });
});

QUnit.module('concatTypedArrays');

QUnit.test('should function as expected', function(assert) {
  const tests = {
    undef: {
      data: concatTypedArrays(),
      expected: toUint8([])
    },
    empty: {
      data: concatTypedArrays(toUint8([])),
      expected: toUint8([])
    },
    single: {
      data: concatTypedArrays([0x01]),
      expected: toUint8([0x01])
    },
    array: {
      data: concatTypedArrays([0x01], [0x02]),
      expected: toUint8([0x01, 0x02])
    },
    uint: {
      data: concatTypedArrays(toUint8([0x01]), toUint8([0x02])),
      expected: toUint8([0x01, 0x02])
    },
    buffer: {
      data: concatTypedArrays(toUint8([0x01]).buffer, toUint8([0x02]).buffer),
      expected: toUint8([0x01, 0x02])
    },
    manyarray: {
      data: concatTypedArrays([0x01], [0x02], [0x03], [0x04]),
      expected: toUint8([0x01, 0x02, 0x03, 0x04])
    },
    manyuint: {
      data: concatTypedArrays(toUint8([0x01]), toUint8([0x02]), toUint8([0x03]), toUint8([0x04])),
      expected: toUint8([0x01, 0x02, 0x03, 0x04])
    }
  };

  Object.keys(tests).forEach(function(name) {
    const {data, expected} = tests[name];

    assert.ok(data instanceof Uint8Array, `obj is a Uint8Array for ${name}`);
    assert.deepEqual(data, expected, `data is as expected for ${name}`);
  });
});
