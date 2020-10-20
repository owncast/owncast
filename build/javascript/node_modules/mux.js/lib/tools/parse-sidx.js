var parseSidx = function(data) {
  var view = new DataView(data.buffer, data.byteOffset, data.byteLength),
    result = {
      version: data[0],
      flags: new Uint8Array(data.subarray(1, 4)),
      references: [],
      referenceId: view.getUint32(4),
      timescale: view.getUint32(8),
      earliestPresentationTime: view.getUint32(12),
      firstOffset: view.getUint32(16)
    },
    referenceCount = view.getUint16(22),
    i;

  for (i = 24; referenceCount; i += 12, referenceCount--) {
    result.references.push({
      referenceType: (data[i] & 0x80) >>> 7,
      referencedSize: view.getUint32(i) & 0x7FFFFFFF,
      subsegmentDuration: view.getUint32(i + 4),
      startsWithSap: !!(data[i + 8] & 0x80),
      sapType: (data[i + 8] & 0x70) >>> 4,
      sapDeltaTime: view.getUint32(i + 8) & 0x0FFFFFFF
    });
  }

  return result;
};

module.exports = parseSidx;
