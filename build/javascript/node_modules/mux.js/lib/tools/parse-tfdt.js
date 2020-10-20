var toUnsigned = require('../utils/bin').toUnsigned;

var tfdt = function(data) {
  var result = {
    version: data[0],
    flags: new Uint8Array(data.subarray(1, 4)),
    baseMediaDecodeTime: toUnsigned(data[4] << 24 | data[5] << 16 | data[6] << 8 | data[7])
  };
  if (result.version === 1) {
    result.baseMediaDecodeTime *= Math.pow(2, 32);
    result.baseMediaDecodeTime += toUnsigned(data[8] << 24 | data[9] << 16 | data[10] << 8 | data[11]);
  }
  return result;
};

module.exports = tfdt;

