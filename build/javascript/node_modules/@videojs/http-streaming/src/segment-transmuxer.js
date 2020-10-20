const transmuxQueue = [];
let currentTransmux;

export const handleData_ = (event, transmuxedData, callback) => {
  const {
    type,
    initSegment,
    captions,
    captionStreams,
    metadata,
    videoFrameDtsTime,
    videoFramePtsTime
  } = event.data.segment;

  transmuxedData.buffer.push({
    captions,
    captionStreams,
    metadata
  });

  // right now, boxes will come back from partial transmuxer, data from full
  const boxes = event.data.segment.boxes || {
    data: event.data.segment.data
  };

  const result = {
    type,
    // cast ArrayBuffer to TypedArray
    data: new Uint8Array(
      boxes.data,
      boxes.data.byteOffset,
      boxes.data.byteLength
    ),
    initSegment: new Uint8Array(
      initSegment.data,
      initSegment.byteOffset,
      initSegment.byteLength
    )
  };

  if (typeof videoFrameDtsTime !== 'undefined') {
    result.videoFrameDtsTime = videoFrameDtsTime;
  }

  if (typeof videoFramePtsTime !== 'undefined') {
    result.videoFramePtsTime = videoFramePtsTime;
  }

  callback(result);
};

export const handleDone_ = ({
  transmuxedData,
  callback
}) => {
  // Previously we only returned data on data events,
  // not on done events. Clear out the buffer to keep that consistent.
  transmuxedData.buffer = [];

  // all buffers should have been flushed from the muxer, so start processing anything we
  // have received
  callback(transmuxedData);
};

export const handleGopInfo_ = (event, transmuxedData) => {
  transmuxedData.gopInfo = event.data.gopInfo;
};

export const processTransmux = ({
  transmuxer,
  bytes,
  audioAppendStart,
  gopsToAlignWith,
  isPartial,
  remux,
  onData,
  onTrackInfo,
  onAudioTimingInfo,
  onVideoTimingInfo,
  onVideoSegmentTimingInfo,
  onId3,
  onCaptions,
  onDone
}) => {
  const transmuxedData = {
    isPartial,
    buffer: []
  };

  const handleMessage = (event) => {
    if (!currentTransmux) {
      // disposed
      return;
    }

    if (event.data.action === 'data') {
      handleData_(event, transmuxedData, onData);
    }
    if (event.data.action === 'trackinfo') {
      onTrackInfo(event.data.trackInfo);
    }
    if (event.data.action === 'gopInfo') {
      handleGopInfo_(event, transmuxedData);
    }
    if (event.data.action === 'audioTimingInfo') {
      onAudioTimingInfo(event.data.audioTimingInfo);
    }
    if (event.data.action === 'videoTimingInfo') {
      onVideoTimingInfo(event.data.videoTimingInfo);
    }
    if (event.data.action === 'videoSegmentTimingInfo') {
      onVideoSegmentTimingInfo(event.data.videoSegmentTimingInfo);
    }
    if (event.data.action === 'id3Frame') {
      onId3([event.data.id3Frame], event.data.id3Frame.dispatchType);
    }
    if (event.data.action === 'caption') {
      onCaptions(event.data.caption);
    }

    // wait for the transmuxed event since we may have audio and video
    if (event.data.type !== 'transmuxed') {
      return;
    }

    transmuxer.onmessage = null;
    handleDone_({
      transmuxedData,
      callback: onDone
    });

    /* eslint-disable no-use-before-define */
    dequeue();
    /* eslint-enable */
  };

  transmuxer.onmessage = handleMessage;

  if (audioAppendStart) {
    transmuxer.postMessage({
      action: 'setAudioAppendStart',
      appendStart: audioAppendStart
    });
  }

  // allow empty arrays to be passed to clear out GOPs
  if (Array.isArray(gopsToAlignWith)) {
    transmuxer.postMessage({
      action: 'alignGopsWith',
      gopsToAlignWith
    });
  }

  if (typeof remux !== 'undefined') {
    transmuxer.postMessage({
      action: 'setRemux',
      remux
    });
  }

  if (bytes.byteLength) {
    const buffer = bytes instanceof ArrayBuffer ? bytes : bytes.buffer;
    const byteOffset = bytes instanceof ArrayBuffer ? 0 : bytes.byteOffset;

    transmuxer.postMessage(
      {
        action: 'push',
        // Send the typed-array of data as an ArrayBuffer so that
        // it can be sent as a "Transferable" and avoid the costly
        // memory copy
        data: buffer,
        // To recreate the original typed-array, we need information
        // about what portion of the ArrayBuffer it was a view into
        byteOffset,
        byteLength: bytes.byteLength
      },
      [ buffer ]
    );
  }

  // even if we didn't push any bytes, we have to make sure we flush in case we reached
  // the end of the segment
  transmuxer.postMessage({ action: isPartial ? 'partialFlush' : 'flush' });
};

export const dequeue = () => {
  currentTransmux = null;
  if (transmuxQueue.length) {
    currentTransmux = transmuxQueue.shift();
    if (typeof currentTransmux === 'function') {
      currentTransmux();
    } else {
      processTransmux(currentTransmux);
    }
  }
};

export const processAction = (transmuxer, action) => {
  transmuxer.postMessage({ action });
  dequeue();
};

export const enqueueAction = (action, transmuxer) => {
  if (!currentTransmux) {
    currentTransmux = action;
    processAction(transmuxer, action);
    return;
  }
  transmuxQueue.push(processAction.bind(null, transmuxer, action));
};

export const reset = (transmuxer) => {
  enqueueAction('reset', transmuxer);
};

export const endTimeline = (transmuxer) => {
  enqueueAction('endTimeline', transmuxer);
};

export const transmux = (options) => {
  if (!currentTransmux) {
    currentTransmux = options;
    processTransmux(options);
    return;
  }
  transmuxQueue.push(options);
};

export const dispose = () => {
  // clear out module-level references
  currentTransmux = null;
  transmuxQueue.length = 0;
};

export default {
  reset,
  dispose,
  endTimeline,
  transmux
};
