/**
 * @file xhr.js
 */

/**
 * A wrapper for videojs.xhr that tracks bandwidth.
 *
 * @param {Object} options options for the XHR
 * @param {Function} callback the callback to call when done
 * @return {Request} the xhr request that is going to be made
 */
import videojs from 'video.js';

const {
  xhr: videojsXHR,
  mergeOptions
} = videojs;

const callbackWrapper = function(request, error, response, callback) {
  const reqResponse = request.responseType === 'arraybuffer' ? request.response : request.responseText;

  if (!error && reqResponse) {
    request.responseTime = Date.now();
    request.roundTripTime = request.responseTime - request.requestTime;
    request.bytesReceived = reqResponse.byteLength || reqResponse.length;
    if (!request.bandwidth) {
      request.bandwidth =
        Math.floor((request.bytesReceived / request.roundTripTime) * 8 * 1000);
    }
  }

  if (response.headers) {
    request.responseHeaders = response.headers;
  }

  // videojs.xhr now uses a specific code on the error
  // object to signal that a request has timed out instead
  // of setting a boolean on the request object
  if (error && error.code === 'ETIMEDOUT') {
    request.timedout = true;
  }

  // videojs.xhr no longer considers status codes outside of 200 and 0
  // (for file uris) to be errors, but the old XHR did, so emulate that
  // behavior. Status 206 may be used in response to byterange requests.
  if (!error &&
    !request.aborted &&
    response.statusCode !== 200 &&
    response.statusCode !== 206 &&
    response.statusCode !== 0) {
    error = new Error('XHR Failed with a response of: ' +
                      (request && (reqResponse || request.responseText)));
  }

  callback(error, request);
};

const xhrFactory = function() {
  const xhr = function XhrFunction(options, callback) {
    // Add a default timeout
    options = mergeOptions({
      timeout: 45e3
    }, options);

    // Allow an optional user-specified function to modify the option
    // object before we construct the xhr request
    const beforeRequest = XhrFunction.beforeRequest || videojs.Vhs.xhr.beforeRequest;

    if (beforeRequest && typeof beforeRequest === 'function') {
      const newOptions = beforeRequest(options);

      if (newOptions) {
        options = newOptions;
      }
    }

    const request = videojsXHR(options, function(error, response) {
      return callbackWrapper(request, error, response, callback);
    });
    const originalAbort = request.abort;

    request.abort = function() {
      request.aborted = true;
      return originalAbort.apply(request, arguments);
    };
    request.uri = options.uri;
    request.requestTime = Date.now();
    return request;
  };

  return xhr;
};

/**
 * Turns segment byterange into a string suitable for use in
 * HTTP Range requests
 *
 * @param {Object} byterange - an object with two values defining the start and end
 *                             of a byte-range
 */
const byterangeStr = function(byterange) {
  // `byterangeEnd` is one less than `offset + length` because the HTTP range
  // header uses inclusive ranges
  const byterangeEnd = byterange.offset + byterange.length - 1;
  const byterangeStart = byterange.offset;

  return 'bytes=' + byterangeStart + '-' + byterangeEnd;
};

/**
 * Defines headers for use in the xhr request for a particular segment.
 *
 * @param {Object} segment - a simplified copy of the segmentInfo object
 *                           from SegmentLoader
 */
const segmentXhrHeaders = function(segment) {
  const headers = {};

  if (segment.byterange) {
    headers.Range = byterangeStr(segment.byterange);
  }
  return headers;
};

export {segmentXhrHeaders, callbackWrapper, xhrFactory};

export default xhrFactory;
