import {bytesToString, toUint8} from './byte-helpers.js';

export const id3Size = function(bytes, offset = 0) {
  bytes = toUint8(bytes);

  const returnSize = (bytes[offset + 6] << 21) |
    (bytes[offset + 7] << 14) |
    (bytes[offset + 8] << 7) |
    (bytes[offset + 9]);
  const flags = bytes[offset + 5];
  const footerPresent = (flags & 16) >> 4;

  if (footerPresent) {
    return returnSize + 20;
  }

  return returnSize + 10;
};

export const getId3Offset = function(bytes, offset = 0) {
  bytes = toUint8(bytes);

  if ((bytes.length - offset) < 10 || bytesToString(bytes.subarray(offset, offset + 3)) !== 'ID3') {
    return offset;
  }

  offset += id3Size(bytes, offset);

  // recursive check for id3 tags as some files
  // have multiple ID3 tag sections even though
  // they should not.
  return getId3Offset(bytes, offset);
};

export const isLikely = {
  aac(bytes) {
    const offset = getId3Offset(bytes);

    return bytes.length >= offset + 2 &&
      (bytes[offset] & 0xFF) === 0xFF &&
      (bytes[offset + 1] & 0xE0) === 0xE0 &&
      (bytes[offset + 1] & 0x16) === 0x10;
  },

  mp3(bytes) {
    const offset = getId3Offset(bytes);

    return bytes.length >= offset + 2 &&
      (bytes[offset] & 0xFF) === 0xFF &&
      (bytes[offset + 1] & 0xE0) === 0xE0 &&
      (bytes[offset + 1] & 0x06) === 0x02;
  },

  webm(bytes) {
    return bytes.length >= 4 &&
      (bytes[0] & 0xFF) === 0x1A &&
      (bytes[1] & 0xFF) === 0x45 &&
      (bytes[2] & 0xFF) === 0xDF &&
      (bytes[3] & 0xFF) === 0xA3;
  },

  mp4(bytes) {
    return bytes.length >= 8 &&
      (/^(f|s)typ$/).test(bytesToString(bytes.subarray(4, 8))) &&
      // not 3gp data
      !(/^ftyp3g$/).test(bytesToString(bytes.subarray(4, 10)));
  },

  '3gp'(bytes) {
    return bytes.length >= 10 &&
      (/^ftyp3g$/).test(bytesToString(bytes.subarray(4, 10)));
  },

  ts(bytes) {
    if (bytes.length < 189 && bytes.length >= 1) {
      return bytes[0] === 0x47;
    }

    let i = 0;

    // check the first 376 bytes for two matching sync bytes
    while (i + 188 < bytes.length && i < 188) {
      if (bytes[i] === 0x47 && bytes[i + 188] === 0x47) {
        return true;
      }
      i += 1;
    }

    return false;
  },
  flac(bytes) {
    return bytes.length >= 4 &&
      (/^fLaC$/).test(bytesToString(bytes.subarray(0, 4)));
  },
  ogg(bytes) {
    return bytes.length >= 4 &&
      (/^OggS$/).test(bytesToString(bytes.subarray(0, 4)));
  }
};

// get all the isLikely functions
// but make sure 'ts' is at the bottom
// as it is the least specific
const isLikelyTypes = Object.keys(isLikely)
  // remove ts
  .filter((t) => t !== 'ts')
  // add it back to the bottom
  .concat('ts');

// make sure we are dealing with uint8 data.
isLikelyTypes.forEach(function(type) {
  const isLikelyFn = isLikely[type];

  isLikely[type] = (bytes) => isLikelyFn(toUint8(bytes));
});

// A useful list of file signatures can be found here
// https://en.wikipedia.org/wiki/List_of_file_signatures
export const detectContainerForBytes = (bytes) => {
  bytes = toUint8(bytes);

  for (let i = 0; i < isLikelyTypes.length; i++) {
    const type = isLikelyTypes[i];

    if (isLikely[type](bytes)) {
      return type;
    }
  }

  return '';
};

// fmp4 is not a container
export const isLikelyFmp4MediaSegment = (bytes) => {
  bytes = toUint8(bytes);
  let i = 0;

  while (i < bytes.length) {
    const size = (bytes[i] << 24 | bytes[i + 1] << 16 | bytes[i + 2] << 8 | bytes[i + 3]) >>> 0;
    const type = bytesToString(bytes.subarray(i + 4, i + 8));

    if (type === 'moof') {
      return true;
    }

    if (size === 0 || (size + i) > bytes.length) {
      i = bytes.length;
    } else {
      i += size;
    }
  }

  return false;
};
