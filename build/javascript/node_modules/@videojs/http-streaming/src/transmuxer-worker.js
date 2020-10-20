/* global self */
/**
 * @file transmuxer-worker.js
 */

/**
 * videojs-contrib-media-sources
 *
 * Copyright (c) 2015 Brightcove
 * All rights reserved.
 *
 * Handles communication between the browser-world and the mux.js
 * transmuxer running inside of a WebWorker by exposing a simple
 * message-based interface to a Transmuxer object.
 */

import {Transmuxer as FullMux} from 'mux.js/lib/mp4/transmuxer';
import PartialMux from 'mux.js/lib/partial/transmuxer';
import CaptionParser from 'mux.js/lib/mp4/caption-parser';
import {
  secondsToVideoTs,
  videoTsToSeconds
} from 'mux.js/lib/utils/clock';

const typeFromStreamString = (streamString) => {
  if (streamString === 'AudioSegmentStream') {
    return 'audio';
  }
  return streamString === 'VideoSegmentStream' ? 'video' : '';
};

/**
 * Re-emits transmuxer events by converting them into messages to the
 * world outside the worker.
 *
 * @param {Object} transmuxer the transmuxer to wire events on
 * @private
 */
const wireFullTransmuxerEvents = function(self, transmuxer) {
  transmuxer.on('data', function(segment) {
    // transfer ownership of the underlying ArrayBuffer
    // instead of doing a copy to save memory
    // ArrayBuffers are transferable but generic TypedArrays are not
    // @link https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API/Using_web_workers#Passing_data_by_transferring_ownership_(transferable_objects)
    const initArray = segment.initSegment;

    segment.initSegment = {
      data: initArray.buffer,
      byteOffset: initArray.byteOffset,
      byteLength: initArray.byteLength
    };

    const typedArray = segment.data;

    segment.data = typedArray.buffer;
    self.postMessage({
      action: 'data',
      segment,
      byteOffset: typedArray.byteOffset,
      byteLength: typedArray.byteLength
    }, [segment.data]);
  });

  transmuxer.on('done', function(data) {
    self.postMessage({ action: 'done' });
  });

  transmuxer.on('gopInfo', function(gopInfo) {
    self.postMessage({
      action: 'gopInfo',
      gopInfo
    });
  });

  transmuxer.on('videoSegmentTimingInfo', function(timingInfo) {
    const videoSegmentTimingInfo = {
      start: {
        decode: videoTsToSeconds(timingInfo.start.dts),
        presentation: videoTsToSeconds(timingInfo.start.pts)
      },
      end: {
        decode: videoTsToSeconds(timingInfo.end.dts),
        presentation: videoTsToSeconds(timingInfo.end.pts)
      },
      baseMediaDecodeTime: videoTsToSeconds(timingInfo.baseMediaDecodeTime)
    };

    if (timingInfo.prependedContentDuration) {
      videoSegmentTimingInfo.prependedContentDuration = videoTsToSeconds(timingInfo.prependedContentDuration);
    }
    self.postMessage({
      action: 'videoSegmentTimingInfo',
      videoSegmentTimingInfo
    });
  });

  transmuxer.on('id3Frame', function(id3Frame) {
    self.postMessage({
      action: 'id3Frame',
      id3Frame
    });
  });

  transmuxer.on('caption', function(caption) {
    self.postMessage({
      action: 'caption',
      caption
    });
  });

  transmuxer.on('trackinfo', function(trackInfo) {
    self.postMessage({
      action: 'trackinfo',
      trackInfo
    });
  });

  transmuxer.on('audioTimingInfo', function(audioTimingInfo) {
    // convert to video TS since we prioritize video time over audio
    self.postMessage({
      action: 'audioTimingInfo',
      audioTimingInfo: {
        start: videoTsToSeconds(audioTimingInfo.start),
        end: videoTsToSeconds(audioTimingInfo.end)
      }
    });
  });

  transmuxer.on('videoTimingInfo', function(videoTimingInfo) {
    self.postMessage({
      action: 'videoTimingInfo',
      videoTimingInfo: {
        start: videoTsToSeconds(videoTimingInfo.start),
        end: videoTsToSeconds(videoTimingInfo.end)
      }
    });
  });

};

const wirePartialTransmuxerEvents = function(self, transmuxer) {
  transmuxer.on('data', function(event) {
    // transfer ownership of the underlying ArrayBuffer
    // instead of doing a copy to save memory
    // ArrayBuffers are transferable but generic TypedArrays are not
    // @link https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API/Using_web_workers#Passing_data_by_transferring_ownership_(transferable_objects)

    const initSegment = {
      data: event.data.track.initSegment.buffer,
      byteOffset: event.data.track.initSegment.byteOffset,
      byteLength: event.data.track.initSegment.byteLength
    };
    const boxes = {
      data: event.data.boxes.buffer,
      byteOffset: event.data.boxes.byteOffset,
      byteLength: event.data.boxes.byteLength
    };
    const segment = {
      boxes,
      initSegment,
      type: event.type,
      sequence: event.data.sequence
    };

    if (typeof event.data.videoFrameDts !== 'undefined') {
      segment.videoFrameDtsTime = videoTsToSeconds(event.data.videoFrameDts);
    }

    if (typeof event.data.videoFramePts !== 'undefined') {
      segment.videoFramePtsTime = videoTsToSeconds(event.data.videoFramePts);
    }

    self.postMessage({
      action: 'data',
      segment
    }, [ segment.boxes.data, segment.initSegment.data ]);
  });

  transmuxer.on('id3Frame', function(id3Frame) {
    self.postMessage({
      action: 'id3Frame',
      id3Frame
    });
  });

  transmuxer.on('caption', function(caption) {
    self.postMessage({
      action: 'caption',
      caption
    });
  });

  transmuxer.on('done', function(data) {
    self.postMessage({
      action: 'done',
      type: typeFromStreamString(data)
    });
  });

  transmuxer.on('partialdone', function(data) {
    self.postMessage({
      action: 'partialdone',
      type: typeFromStreamString(data)
    });
  });

  transmuxer.on('endedsegment', function(data) {
    self.postMessage({
      action: 'endedSegment',
      type: typeFromStreamString(data)
    });
  });

  transmuxer.on('trackinfo', function(trackInfo) {
    self.postMessage({ action: 'trackinfo', trackInfo });
  });

  transmuxer.on('audioTimingInfo', function(audioTimingInfo) {
    // This can happen if flush is called when no
    // audio has been processed. This should be an
    // unusual case, but if it does occur should not
    // result in valid data being returned
    if (audioTimingInfo.start === null) {
      self.postMessage({
        action: 'audioTimingInfo',
        audioTimingInfo
      });
      return;
    }

    // convert to video TS since we prioritize video time over audio
    const timingInfoInSeconds = {
      start: videoTsToSeconds(audioTimingInfo.start)
    };

    if (audioTimingInfo.end) {
      timingInfoInSeconds.end = videoTsToSeconds(audioTimingInfo.end);
    }

    self.postMessage({
      action: 'audioTimingInfo',
      audioTimingInfo: timingInfoInSeconds
    });
  });

  transmuxer.on('videoTimingInfo', function(videoTimingInfo) {
    const timingInfoInSeconds = {
      start: videoTsToSeconds(videoTimingInfo.start)
    };

    if (videoTimingInfo.end) {
      timingInfoInSeconds.end = videoTsToSeconds(videoTimingInfo.end);
    }

    self.postMessage({
      action: 'videoTimingInfo',
      videoTimingInfo: timingInfoInSeconds
    });
  });
};

/**
 * All incoming messages route through this hash. If no function exists
 * to handle an incoming message, then we ignore the message.
 *
 * @class MessageHandlers
 * @param {Object} options the options to initialize with
 */
class MessageHandlers {
  constructor(self, options) {
    this.options = options || {};
    this.self = self;
    this.init();
  }

  /**
   * initialize our web worker and wire all the events.
   */
  init() {
    if (this.transmuxer) {
      this.transmuxer.dispose();
    }
    this.transmuxer = this.options.handlePartialData ?
      new PartialMux(this.options) :
      new FullMux(this.options);

    if (this.options.handlePartialData) {
      wirePartialTransmuxerEvents(this.self, this.transmuxer);
    } else {
      wireFullTransmuxerEvents(this.self, this.transmuxer);
    }

  }

  pushMp4Captions(data) {
    if (!this.captionParser) {
      this.captionParser = new CaptionParser();
      this.captionParser.init();
    }
    const segment = new Uint8Array(data.data, data.byteOffset, data.byteLength);
    const parsed = this.captionParser.parse(
      segment,
      data.trackIds,
      data.timescales
    );

    this.self.postMessage({
      action: 'mp4Captions',
      captions: parsed && parsed.captions || [],
      data: segment.buffer
    }, [segment.buffer]);
  }

  clearAllMp4Captions() {
    if (this.captionParser) {
      this.captionParser.clearAllCaptions();
    }
  }

  clearParsedMp4Captions() {
    if (this.captionParser) {
      this.captionParser.clearParsedCaptions();
    }
  }

  /**
   * Adds data (a ts segment) to the start of the transmuxer pipeline for
   * processing.
   *
   * @param {ArrayBuffer} data data to push into the muxer
   */
  push(data) {
    // Cast array buffer to correct type for transmuxer
    const segment = new Uint8Array(data.data, data.byteOffset, data.byteLength);

    this.transmuxer.push(segment);
  }

  /**
   * Recreate the transmuxer so that the next segment added via `push`
   * start with a fresh transmuxer.
   */
  reset() {
    this.transmuxer.reset();
  }

  /**
   * Set the value that will be used as the `baseMediaDecodeTime` time for the
   * next segment pushed in. Subsequent segments will have their `baseMediaDecodeTime`
   * set relative to the first based on the PTS values.
   *
   * @param {Object} data used to set the timestamp offset in the muxer
   */
  setTimestampOffset(data) {
    const timestampOffset = data.timestampOffset || 0;

    this.transmuxer.setBaseMediaDecodeTime(Math.round(secondsToVideoTs(timestampOffset)));
  }

  setAudioAppendStart(data) {
    this.transmuxer.setAudioAppendStart(Math.ceil(secondsToVideoTs(data.appendStart)));
  }

  setRemux(data) {
    this.transmuxer.setRemux(data.remux);
  }

  /**
   * Forces the pipeline to finish processing the last segment and emit it's
   * results.
   *
   * @param {Object} data event data, not really used
   */
  flush(data) {
    this.transmuxer.flush();
    // transmuxed done action is fired after both audio/video pipelines are flushed
    self.postMessage({
      action: 'done',
      type: 'transmuxed'
    });
  }

  partialFlush(data) {
    this.transmuxer.partialFlush();
    // transmuxed partialdone action is fired after both audio/video pipelines are flushed
    self.postMessage({
      action: 'partialdone',
      type: 'transmuxed'
    });
  }

  endTimeline() {
    this.transmuxer.endTimeline();
    // transmuxed endedtimeline action is fired after both audio/video pipelines end their
    // timelines
    self.postMessage({
      action: 'endedtimeline',
      type: 'transmuxed'
    });
  }

  alignGopsWith(data) {
    this.transmuxer.alignGopsWith(data.gopsToAlignWith.slice());
  }
}

/**
 * Our web worker interface so that things can talk to mux.js
 * that will be running in a web worker. the scope is passed to this by
 * webworkify.
 *
 * @param {Object} self the scope for the web worker
 */
const TransmuxerWorker = function(self) {
  self.onmessage = function(event) {
    if (event.data.action === 'init' && event.data.options) {
      this.messageHandlers = new MessageHandlers(self, event.data.options);
      return;
    }

    if (!this.messageHandlers) {
      this.messageHandlers = new MessageHandlers(self);
    }

    if (event.data && event.data.action && event.data.action !== 'init') {
      if (this.messageHandlers[event.data.action]) {
        this.messageHandlers[event.data.action](event.data);
      }
    }
  };
};

export default new TransmuxerWorker(self);
