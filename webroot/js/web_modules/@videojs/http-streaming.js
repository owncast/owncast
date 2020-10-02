import { w as window_1, u as urlToolkit, i as inheritsLoose, a as assertThisInitialized, _ as _extends_1, d as domParser, v as videojs$1, b as document_1 } from '../common/video.es-915ea113.js';
import { c as createCommonjsModule } from '../common/_commonjsHelpers-798ad6a7.js';

function _interopDefault (ex) { return (ex && (typeof ex === 'object') && 'default' in ex) ? ex['default'] : ex; }

var URLToolkit = _interopDefault(urlToolkit);
var window$1 = _interopDefault(window_1);

var resolveUrl = function resolveUrl(baseUrl, relativeUrl) {
  // return early if we don't need to resolve
  if (/^[a-z]+:/i.test(relativeUrl)) {
    return relativeUrl;
  } // if the base URL is relative then combine with the current location


  if (!/\/\//i.test(baseUrl)) {
    baseUrl = URLToolkit.buildAbsoluteURL(window$1.location && window$1.location.href || '', baseUrl);
  }

  return URLToolkit.buildAbsoluteURL(baseUrl, relativeUrl);
};

var resolveUrl_1 = resolveUrl;

/*! @name @videojs/vhs-utils @version 1.3.0 @license MIT */

/**
 * @file stream.js
 */

/**
 * A lightweight readable stream implemention that handles event dispatching.
 *
 * @class Stream
 */
var Stream =
/*#__PURE__*/
function () {
  function Stream() {
    this.listeners = {};
  }
  /**
   * Add a listener for a specified event type.
   *
   * @param {string} type the event name
   * @param {Function} listener the callback to be invoked when an event of
   * the specified type occurs
   */


  var _proto = Stream.prototype;

  _proto.on = function on(type, listener) {
    if (!this.listeners[type]) {
      this.listeners[type] = [];
    }

    this.listeners[type].push(listener);
  }
  /**
   * Remove a listener for a specified event type.
   *
   * @param {string} type the event name
   * @param {Function} listener  a function previously registered for this
   * type of event through `on`
   * @return {boolean} if we could turn it off or not
   */
  ;

  _proto.off = function off(type, listener) {
    if (!this.listeners[type]) {
      return false;
    }

    var index = this.listeners[type].indexOf(listener); // TODO: which is better?
    // In Video.js we slice listener functions
    // on trigger so that it does not mess up the order
    // while we loop through.
    //
    // Here we slice on off so that the loop in trigger
    // can continue using it's old reference to loop without
    // messing up the order.

    this.listeners[type] = this.listeners[type].slice(0);
    this.listeners[type].splice(index, 1);
    return index > -1;
  }
  /**
   * Trigger an event of the specified type on this stream. Any additional
   * arguments to this function are passed as parameters to event listeners.
   *
   * @param {string} type the event name
   */
  ;

  _proto.trigger = function trigger(type) {
    var callbacks = this.listeners[type];

    if (!callbacks) {
      return;
    } // Slicing the arguments on every invocation of this method
    // can add a significant amount of overhead. Avoid the
    // intermediate object creation for the common case of a
    // single callback argument


    if (arguments.length === 2) {
      var length = callbacks.length;

      for (var i = 0; i < length; ++i) {
        callbacks[i].call(this, arguments[1]);
      }
    } else {
      var args = Array.prototype.slice.call(arguments, 1);
      var _length = callbacks.length;

      for (var _i = 0; _i < _length; ++_i) {
        callbacks[_i].apply(this, args);
      }
    }
  }
  /**
   * Destroys the stream and cleans up.
   */
  ;

  _proto.dispose = function dispose() {
    this.listeners = {};
  }
  /**
   * Forwards all `data` events on this stream to the destination stream. The
   * destination stream should provide a method `push` to receive the data
   * events as they arrive.
   *
   * @param {Stream} destination the stream that will receive all `data` events
   * @see http://nodejs.org/api/stream.html#stream_readable_pipe_destination_options
   */
  ;

  _proto.pipe = function pipe(destination) {
    this.on('data', function (data) {
      destination.push(data);
    });
  };

  return Stream;
}();

var stream = Stream;

function _interopDefault$1 (ex) { return (ex && (typeof ex === 'object') && 'default' in ex) ? ex['default'] : ex; }

var window$2 = _interopDefault$1(window_1);

var atob = function atob(s) {
  return window$2.atob ? window$2.atob(s) : Buffer.from(s, 'base64').toString('binary');
};

function decodeB64ToUint8Array(b64Text) {
  var decodedString = atob(b64Text);
  var array = new Uint8Array(decodedString.length);

  for (var i = 0; i < decodedString.length; i++) {
    array[i] = decodedString.charCodeAt(i);
  }

  return array;
}

var decodeB64ToUint8Array_1 = decodeB64ToUint8Array;

/*! @name m3u8-parser @version 4.4.3 @license Apache-2.0 */

/**
 * A stream that buffers string input and generates a `data` event for each
 * line.
 *
 * @class LineStream
 * @extends Stream
 */

var LineStream =
/*#__PURE__*/
function (_Stream) {
  inheritsLoose(LineStream, _Stream);

  function LineStream() {
    var _this;

    _this = _Stream.call(this) || this;
    _this.buffer = '';
    return _this;
  }
  /**
   * Add new data to be parsed.
   *
   * @param {string} data the text to process
   */


  var _proto = LineStream.prototype;

  _proto.push = function push(data) {
    var nextNewline;
    this.buffer += data;
    nextNewline = this.buffer.indexOf('\n');

    for (; nextNewline > -1; nextNewline = this.buffer.indexOf('\n')) {
      this.trigger('data', this.buffer.substring(0, nextNewline));
      this.buffer = this.buffer.substring(nextNewline + 1);
    }
  };

  return LineStream;
}(stream);

/**
 * "forgiving" attribute list psuedo-grammar:
 * attributes -> keyvalue (',' keyvalue)*
 * keyvalue   -> key '=' value
 * key        -> [^=]*
 * value      -> '"' [^"]* '"' | [^,]*
 */

var attributeSeparator = function attributeSeparator() {
  var key = '[^=]*';
  var value = '"[^"]*"|[^,]*';
  var keyvalue = '(?:' + key + ')=(?:' + value + ')';
  return new RegExp('(?:^|,)(' + keyvalue + ')');
};
/**
 * Parse attributes from a line given the separator
 *
 * @param {string} attributes the attribute line to parse
 */


var parseAttributes = function parseAttributes(attributes) {
  // split the string using attributes as the separator
  var attrs = attributes.split(attributeSeparator());
  var result = {};
  var i = attrs.length;
  var attr;

  while (i--) {
    // filter out unmatched portions of the string
    if (attrs[i] === '') {
      continue;
    } // split the key and value


    attr = /([^=]*)=(.*)/.exec(attrs[i]).slice(1); // trim whitespace and remove optional quotes around the value

    attr[0] = attr[0].replace(/^\s+|\s+$/g, '');
    attr[1] = attr[1].replace(/^\s+|\s+$/g, '');
    attr[1] = attr[1].replace(/^['"](.*)['"]$/g, '$1');
    result[attr[0]] = attr[1];
  }

  return result;
};
/**
 * A line-level M3U8 parser event stream. It expects to receive input one
 * line at a time and performs a context-free parse of its contents. A stream
 * interpretation of a manifest can be useful if the manifest is expected to
 * be too large to fit comfortably into memory or the entirety of the input
 * is not immediately available. Otherwise, it's probably much easier to work
 * with a regular `Parser` object.
 *
 * Produces `data` events with an object that captures the parser's
 * interpretation of the input. That object has a property `tag` that is one
 * of `uri`, `comment`, or `tag`. URIs only have a single additional
 * property, `line`, which captures the entirety of the input without
 * interpretation. Comments similarly have a single additional property
 * `text` which is the input without the leading `#`.
 *
 * Tags always have a property `tagType` which is the lower-cased version of
 * the M3U8 directive without the `#EXT` or `#EXT-X-` prefix. For instance,
 * `#EXT-X-MEDIA-SEQUENCE` becomes `media-sequence` when parsed. Unrecognized
 * tags are given the tag type `unknown` and a single additional property
 * `data` with the remainder of the input.
 *
 * @class ParseStream
 * @extends Stream
 */


var ParseStream =
/*#__PURE__*/
function (_Stream) {
  inheritsLoose(ParseStream, _Stream);

  function ParseStream() {
    var _this;

    _this = _Stream.call(this) || this;
    _this.customParsers = [];
    _this.tagMappers = [];
    return _this;
  }
  /**
   * Parses an additional line of input.
   *
   * @param {string} line a single line of an M3U8 file to parse
   */


  var _proto = ParseStream.prototype;

  _proto.push = function push(line) {
    var _this2 = this;

    var match;
    var event; // strip whitespace

    line = line.trim();

    if (line.length === 0) {
      // ignore empty lines
      return;
    } // URIs


    if (line[0] !== '#') {
      this.trigger('data', {
        type: 'uri',
        uri: line
      });
      return;
    } // map tags


    var newLines = this.tagMappers.reduce(function (acc, mapper) {
      var mappedLine = mapper(line); // skip if unchanged

      if (mappedLine === line) {
        return acc;
      }

      return acc.concat([mappedLine]);
    }, [line]);
    newLines.forEach(function (newLine) {
      for (var i = 0; i < _this2.customParsers.length; i++) {
        if (_this2.customParsers[i].call(_this2, newLine)) {
          return;
        }
      } // Comments


      if (newLine.indexOf('#EXT') !== 0) {
        _this2.trigger('data', {
          type: 'comment',
          text: newLine.slice(1)
        });

        return;
      } // strip off any carriage returns here so the regex matching
      // doesn't have to account for them.


      newLine = newLine.replace('\r', ''); // Tags

      match = /^#EXTM3U/.exec(newLine);

      if (match) {
        _this2.trigger('data', {
          type: 'tag',
          tagType: 'm3u'
        });

        return;
      }

      match = /^#EXTINF:?([0-9\.]*)?,?(.*)?$/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'inf'
        };

        if (match[1]) {
          event.duration = parseFloat(match[1]);
        }

        if (match[2]) {
          event.title = match[2];
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-TARGETDURATION:?([0-9.]*)?/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'targetduration'
        };

        if (match[1]) {
          event.duration = parseInt(match[1], 10);
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#ZEN-TOTAL-DURATION:?([0-9.]*)?/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'totalduration'
        };

        if (match[1]) {
          event.duration = parseInt(match[1], 10);
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-VERSION:?([0-9.]*)?/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'version'
        };

        if (match[1]) {
          event.version = parseInt(match[1], 10);
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-MEDIA-SEQUENCE:?(\-?[0-9.]*)?/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'media-sequence'
        };

        if (match[1]) {
          event.number = parseInt(match[1], 10);
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-DISCONTINUITY-SEQUENCE:?(\-?[0-9.]*)?/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'discontinuity-sequence'
        };

        if (match[1]) {
          event.number = parseInt(match[1], 10);
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-PLAYLIST-TYPE:?(.*)?$/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'playlist-type'
        };

        if (match[1]) {
          event.playlistType = match[1];
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-BYTERANGE:?([0-9.]*)?@?([0-9.]*)?/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'byterange'
        };

        if (match[1]) {
          event.length = parseInt(match[1], 10);
        }

        if (match[2]) {
          event.offset = parseInt(match[2], 10);
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-ALLOW-CACHE:?(YES|NO)?/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'allow-cache'
        };

        if (match[1]) {
          event.allowed = !/NO/.test(match[1]);
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-MAP:?(.*)$/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'map'
        };

        if (match[1]) {
          var attributes = parseAttributes(match[1]);

          if (attributes.URI) {
            event.uri = attributes.URI;
          }

          if (attributes.BYTERANGE) {
            var _attributes$BYTERANGE = attributes.BYTERANGE.split('@'),
                length = _attributes$BYTERANGE[0],
                offset = _attributes$BYTERANGE[1];

            event.byterange = {};

            if (length) {
              event.byterange.length = parseInt(length, 10);
            }

            if (offset) {
              event.byterange.offset = parseInt(offset, 10);
            }
          }
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-STREAM-INF:?(.*)$/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'stream-inf'
        };

        if (match[1]) {
          event.attributes = parseAttributes(match[1]);

          if (event.attributes.RESOLUTION) {
            var split = event.attributes.RESOLUTION.split('x');
            var resolution = {};

            if (split[0]) {
              resolution.width = parseInt(split[0], 10);
            }

            if (split[1]) {
              resolution.height = parseInt(split[1], 10);
            }

            event.attributes.RESOLUTION = resolution;
          }

          if (event.attributes.BANDWIDTH) {
            event.attributes.BANDWIDTH = parseInt(event.attributes.BANDWIDTH, 10);
          }

          if (event.attributes['PROGRAM-ID']) {
            event.attributes['PROGRAM-ID'] = parseInt(event.attributes['PROGRAM-ID'], 10);
          }
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-MEDIA:?(.*)$/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'media'
        };

        if (match[1]) {
          event.attributes = parseAttributes(match[1]);
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-ENDLIST/.exec(newLine);

      if (match) {
        _this2.trigger('data', {
          type: 'tag',
          tagType: 'endlist'
        });

        return;
      }

      match = /^#EXT-X-DISCONTINUITY/.exec(newLine);

      if (match) {
        _this2.trigger('data', {
          type: 'tag',
          tagType: 'discontinuity'
        });

        return;
      }

      match = /^#EXT-X-PROGRAM-DATE-TIME:?(.*)$/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'program-date-time'
        };

        if (match[1]) {
          event.dateTimeString = match[1];
          event.dateTimeObject = new Date(match[1]);
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-KEY:?(.*)$/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'key'
        };

        if (match[1]) {
          event.attributes = parseAttributes(match[1]); // parse the IV string into a Uint32Array

          if (event.attributes.IV) {
            if (event.attributes.IV.substring(0, 2).toLowerCase() === '0x') {
              event.attributes.IV = event.attributes.IV.substring(2);
            }

            event.attributes.IV = event.attributes.IV.match(/.{8}/g);
            event.attributes.IV[0] = parseInt(event.attributes.IV[0], 16);
            event.attributes.IV[1] = parseInt(event.attributes.IV[1], 16);
            event.attributes.IV[2] = parseInt(event.attributes.IV[2], 16);
            event.attributes.IV[3] = parseInt(event.attributes.IV[3], 16);
            event.attributes.IV = new Uint32Array(event.attributes.IV);
          }
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-START:?(.*)$/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'start'
        };

        if (match[1]) {
          event.attributes = parseAttributes(match[1]);
          event.attributes['TIME-OFFSET'] = parseFloat(event.attributes['TIME-OFFSET']);
          event.attributes.PRECISE = /YES/.test(event.attributes.PRECISE);
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-CUE-OUT-CONT:?(.*)?$/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'cue-out-cont'
        };

        if (match[1]) {
          event.data = match[1];
        } else {
          event.data = '';
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-CUE-OUT:?(.*)?$/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'cue-out'
        };

        if (match[1]) {
          event.data = match[1];
        } else {
          event.data = '';
        }

        _this2.trigger('data', event);

        return;
      }

      match = /^#EXT-X-CUE-IN:?(.*)?$/.exec(newLine);

      if (match) {
        event = {
          type: 'tag',
          tagType: 'cue-in'
        };

        if (match[1]) {
          event.data = match[1];
        } else {
          event.data = '';
        }

        _this2.trigger('data', event);

        return;
      } // unknown tag type


      _this2.trigger('data', {
        type: 'tag',
        data: newLine.slice(4)
      });
    });
  }
  /**
   * Add a parser for custom headers
   *
   * @param {Object}   options              a map of options for the added parser
   * @param {RegExp}   options.expression   a regular expression to match the custom header
   * @param {string}   options.customType   the custom type to register to the output
   * @param {Function} [options.dataParser] function to parse the line into an object
   * @param {boolean}  [options.segment]    should tag data be attached to the segment object
   */
  ;

  _proto.addParser = function addParser(_ref) {
    var _this3 = this;

    var expression = _ref.expression,
        customType = _ref.customType,
        dataParser = _ref.dataParser,
        segment = _ref.segment;

    if (typeof dataParser !== 'function') {
      dataParser = function dataParser(line) {
        return line;
      };
    }

    this.customParsers.push(function (line) {
      var match = expression.exec(line);

      if (match) {
        _this3.trigger('data', {
          type: 'custom',
          data: dataParser(line),
          customType: customType,
          segment: segment
        });

        return true;
      }
    });
  }
  /**
   * Add a custom header mapper
   *
   * @param {Object}   options
   * @param {RegExp}   options.expression   a regular expression to match the custom header
   * @param {Function} options.map          function to translate tag into a different tag
   */
  ;

  _proto.addTagMapper = function addTagMapper(_ref2) {
    var expression = _ref2.expression,
        map = _ref2.map;

    var mapFn = function mapFn(line) {
      if (expression.test(line)) {
        return map(line);
      }

      return line;
    };

    this.tagMappers.push(mapFn);
  };

  return ParseStream;
}(stream);

/**
 * A parser for M3U8 files. The current interpretation of the input is
 * exposed as a property `manifest` on parser objects. It's just two lines to
 * create and parse a manifest once you have the contents available as a string:
 *
 * ```js
 * var parser = new m3u8.Parser();
 * parser.push(xhr.responseText);
 * ```
 *
 * New input can later be applied to update the manifest object by calling
 * `push` again.
 *
 * The parser attempts to create a usable manifest object even if the
 * underlying input is somewhat nonsensical. It emits `info` and `warning`
 * events during the parse if it encounters input that seems invalid or
 * requires some property of the manifest object to be defaulted.
 *
 * @class Parser
 * @extends Stream
 */

var Parser =
/*#__PURE__*/
function (_Stream) {
  inheritsLoose(Parser, _Stream);

  function Parser() {
    var _this;

    _this = _Stream.call(this) || this;
    _this.lineStream = new LineStream();
    _this.parseStream = new ParseStream();

    _this.lineStream.pipe(_this.parseStream);
    /* eslint-disable consistent-this */


    var self = assertThisInitialized(_this);
    /* eslint-enable consistent-this */


    var uris = [];
    var currentUri = {}; // if specified, the active EXT-X-MAP definition

    var currentMap; // if specified, the active decryption key

    var _key;

    var noop = function noop() {};

    var defaultMediaGroups = {
      'AUDIO': {},
      'VIDEO': {},
      'CLOSED-CAPTIONS': {},
      'SUBTITLES': {}
    }; // This is the Widevine UUID from DASH IF IOP. The same exact string is
    // used in MPDs with Widevine encrypted streams.

    var widevineUuid = 'urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed'; // group segments into numbered timelines delineated by discontinuities

    var currentTimeline = 0; // the manifest is empty until the parse stream begins delivering data

    _this.manifest = {
      allowCache: true,
      discontinuityStarts: [],
      segments: []
    }; // keep track of the last seen segment's byte range end, as segments are not required
    // to provide the offset, in which case it defaults to the next byte after the
    // previous segment

    var lastByterangeEnd = 0; // update the manifest with the m3u8 entry from the parse stream

    _this.parseStream.on('data', function (entry) {
      var mediaGroup;
      var rendition;
      ({
        tag: function tag() {
          // switch based on the tag type
          (({
            'allow-cache': function allowCache() {
              this.manifest.allowCache = entry.allowed;

              if (!('allowed' in entry)) {
                this.trigger('info', {
                  message: 'defaulting allowCache to YES'
                });
                this.manifest.allowCache = true;
              }
            },
            byterange: function byterange() {
              var byterange = {};

              if ('length' in entry) {
                currentUri.byterange = byterange;
                byterange.length = entry.length;

                if (!('offset' in entry)) {
                  /*
                   * From the latest spec (as of this writing):
                   * https://tools.ietf.org/html/draft-pantos-http-live-streaming-23#section-4.3.2.2
                   *
                   * Same text since EXT-X-BYTERANGE's introduction in draft 7:
                   * https://tools.ietf.org/html/draft-pantos-http-live-streaming-07#section-3.3.1)
                   *
                   * "If o [offset] is not present, the sub-range begins at the next byte
                   * following the sub-range of the previous media segment."
                   */
                  entry.offset = lastByterangeEnd;
                }
              }

              if ('offset' in entry) {
                currentUri.byterange = byterange;
                byterange.offset = entry.offset;
              }

              lastByterangeEnd = byterange.offset + byterange.length;
            },
            endlist: function endlist() {
              this.manifest.endList = true;
            },
            inf: function inf() {
              if (!('mediaSequence' in this.manifest)) {
                this.manifest.mediaSequence = 0;
                this.trigger('info', {
                  message: 'defaulting media sequence to zero'
                });
              }

              if (!('discontinuitySequence' in this.manifest)) {
                this.manifest.discontinuitySequence = 0;
                this.trigger('info', {
                  message: 'defaulting discontinuity sequence to zero'
                });
              }

              if (entry.duration > 0) {
                currentUri.duration = entry.duration;
              }

              if (entry.duration === 0) {
                currentUri.duration = 0.01;
                this.trigger('info', {
                  message: 'updating zero segment duration to a small value'
                });
              }

              this.manifest.segments = uris;
            },
            key: function key() {
              if (!entry.attributes) {
                this.trigger('warn', {
                  message: 'ignoring key declaration without attribute list'
                });
                return;
              } // clear the active encryption key


              if (entry.attributes.METHOD === 'NONE') {
                _key = null;
                return;
              }

              if (!entry.attributes.URI) {
                this.trigger('warn', {
                  message: 'ignoring key declaration without URI'
                });
                return;
              } // check if the content is encrypted for Widevine
              // Widevine/HLS spec: https://storage.googleapis.com/wvdocs/Widevine_DRM_HLS.pdf


              if (entry.attributes.KEYFORMAT === widevineUuid) {
                var VALID_METHODS = ['SAMPLE-AES', 'SAMPLE-AES-CTR', 'SAMPLE-AES-CENC'];

                if (VALID_METHODS.indexOf(entry.attributes.METHOD) === -1) {
                  this.trigger('warn', {
                    message: 'invalid key method provided for Widevine'
                  });
                  return;
                }

                if (entry.attributes.METHOD === 'SAMPLE-AES-CENC') {
                  this.trigger('warn', {
                    message: 'SAMPLE-AES-CENC is deprecated, please use SAMPLE-AES-CTR instead'
                  });
                }

                if (entry.attributes.URI.substring(0, 23) !== 'data:text/plain;base64,') {
                  this.trigger('warn', {
                    message: 'invalid key URI provided for Widevine'
                  });
                  return;
                }

                if (!(entry.attributes.KEYID && entry.attributes.KEYID.substring(0, 2) === '0x')) {
                  this.trigger('warn', {
                    message: 'invalid key ID provided for Widevine'
                  });
                  return;
                } // if Widevine key attributes are valid, store them as `contentProtection`
                // on the manifest to emulate Widevine tag structure in a DASH mpd


                this.manifest.contentProtection = {
                  'com.widevine.alpha': {
                    attributes: {
                      schemeIdUri: entry.attributes.KEYFORMAT,
                      // remove '0x' from the key id string
                      keyId: entry.attributes.KEYID.substring(2)
                    },
                    // decode the base64-encoded PSSH box
                    pssh: decodeB64ToUint8Array_1(entry.attributes.URI.split(',')[1])
                  }
                };
                return;
              }

              if (!entry.attributes.METHOD) {
                this.trigger('warn', {
                  message: 'defaulting key method to AES-128'
                });
              } // setup an encryption key for upcoming segments


              _key = {
                method: entry.attributes.METHOD || 'AES-128',
                uri: entry.attributes.URI
              };

              if (typeof entry.attributes.IV !== 'undefined') {
                _key.iv = entry.attributes.IV;
              }
            },
            'media-sequence': function mediaSequence() {
              if (!isFinite(entry.number)) {
                this.trigger('warn', {
                  message: 'ignoring invalid media sequence: ' + entry.number
                });
                return;
              }

              this.manifest.mediaSequence = entry.number;
            },
            'discontinuity-sequence': function discontinuitySequence() {
              if (!isFinite(entry.number)) {
                this.trigger('warn', {
                  message: 'ignoring invalid discontinuity sequence: ' + entry.number
                });
                return;
              }

              this.manifest.discontinuitySequence = entry.number;
              currentTimeline = entry.number;
            },
            'playlist-type': function playlistType() {
              if (!/VOD|EVENT/.test(entry.playlistType)) {
                this.trigger('warn', {
                  message: 'ignoring unknown playlist type: ' + entry.playlist
                });
                return;
              }

              this.manifest.playlistType = entry.playlistType;
            },
            map: function map() {
              currentMap = {};

              if (entry.uri) {
                currentMap.uri = entry.uri;
              }

              if (entry.byterange) {
                currentMap.byterange = entry.byterange;
              }
            },
            'stream-inf': function streamInf() {
              this.manifest.playlists = uris;
              this.manifest.mediaGroups = this.manifest.mediaGroups || defaultMediaGroups;

              if (!entry.attributes) {
                this.trigger('warn', {
                  message: 'ignoring empty stream-inf attributes'
                });
                return;
              }

              if (!currentUri.attributes) {
                currentUri.attributes = {};
              }

              _extends_1(currentUri.attributes, entry.attributes);
            },
            media: function media() {
              this.manifest.mediaGroups = this.manifest.mediaGroups || defaultMediaGroups;

              if (!(entry.attributes && entry.attributes.TYPE && entry.attributes['GROUP-ID'] && entry.attributes.NAME)) {
                this.trigger('warn', {
                  message: 'ignoring incomplete or missing media group'
                });
                return;
              } // find the media group, creating defaults as necessary


              var mediaGroupType = this.manifest.mediaGroups[entry.attributes.TYPE];
              mediaGroupType[entry.attributes['GROUP-ID']] = mediaGroupType[entry.attributes['GROUP-ID']] || {};
              mediaGroup = mediaGroupType[entry.attributes['GROUP-ID']]; // collect the rendition metadata

              rendition = {
                default: /yes/i.test(entry.attributes.DEFAULT)
              };

              if (rendition.default) {
                rendition.autoselect = true;
              } else {
                rendition.autoselect = /yes/i.test(entry.attributes.AUTOSELECT);
              }

              if (entry.attributes.LANGUAGE) {
                rendition.language = entry.attributes.LANGUAGE;
              }

              if (entry.attributes.URI) {
                rendition.uri = entry.attributes.URI;
              }

              if (entry.attributes['INSTREAM-ID']) {
                rendition.instreamId = entry.attributes['INSTREAM-ID'];
              }

              if (entry.attributes.CHARACTERISTICS) {
                rendition.characteristics = entry.attributes.CHARACTERISTICS;
              }

              if (entry.attributes.FORCED) {
                rendition.forced = /yes/i.test(entry.attributes.FORCED);
              } // insert the new rendition


              mediaGroup[entry.attributes.NAME] = rendition;
            },
            discontinuity: function discontinuity() {
              currentTimeline += 1;
              currentUri.discontinuity = true;
              this.manifest.discontinuityStarts.push(uris.length);
            },
            'program-date-time': function programDateTime() {
              if (typeof this.manifest.dateTimeString === 'undefined') {
                // PROGRAM-DATE-TIME is a media-segment tag, but for backwards
                // compatibility, we add the first occurence of the PROGRAM-DATE-TIME tag
                // to the manifest object
                // TODO: Consider removing this in future major version
                this.manifest.dateTimeString = entry.dateTimeString;
                this.manifest.dateTimeObject = entry.dateTimeObject;
              }

              currentUri.dateTimeString = entry.dateTimeString;
              currentUri.dateTimeObject = entry.dateTimeObject;
            },
            targetduration: function targetduration() {
              if (!isFinite(entry.duration) || entry.duration < 0) {
                this.trigger('warn', {
                  message: 'ignoring invalid target duration: ' + entry.duration
                });
                return;
              }

              this.manifest.targetDuration = entry.duration;
            },
            totalduration: function totalduration() {
              if (!isFinite(entry.duration) || entry.duration < 0) {
                this.trigger('warn', {
                  message: 'ignoring invalid total duration: ' + entry.duration
                });
                return;
              }

              this.manifest.totalDuration = entry.duration;
            },
            start: function start() {
              if (!entry.attributes || isNaN(entry.attributes['TIME-OFFSET'])) {
                this.trigger('warn', {
                  message: 'ignoring start declaration without appropriate attribute list'
                });
                return;
              }

              this.manifest.start = {
                timeOffset: entry.attributes['TIME-OFFSET'],
                precise: entry.attributes.PRECISE
              };
            },
            'cue-out': function cueOut() {
              currentUri.cueOut = entry.data;
            },
            'cue-out-cont': function cueOutCont() {
              currentUri.cueOutCont = entry.data;
            },
            'cue-in': function cueIn() {
              currentUri.cueIn = entry.data;
            }
          })[entry.tagType] || noop).call(self);
        },
        uri: function uri() {
          currentUri.uri = entry.uri;
          uris.push(currentUri); // if no explicit duration was declared, use the target duration

          if (this.manifest.targetDuration && !('duration' in currentUri)) {
            this.trigger('warn', {
              message: 'defaulting segment duration to the target duration'
            });
            currentUri.duration = this.manifest.targetDuration;
          } // annotate with encryption information, if necessary


          if (_key) {
            currentUri.key = _key;
          }

          currentUri.timeline = currentTimeline; // annotate with initialization segment information, if necessary

          if (currentMap) {
            currentUri.map = currentMap;
          } // prepare for the next URI


          currentUri = {};
        },
        comment: function comment() {// comments are not important for playback
        },
        custom: function custom() {
          // if this is segment-level data attach the output to the segment
          if (entry.segment) {
            currentUri.custom = currentUri.custom || {};
            currentUri.custom[entry.customType] = entry.data; // if this is manifest-level data attach to the top level manifest object
          } else {
            this.manifest.custom = this.manifest.custom || {};
            this.manifest.custom[entry.customType] = entry.data;
          }
        }
      })[entry.type].call(self);
    });

    return _this;
  }
  /**
   * Parse the input string and update the manifest object.
   *
   * @param {string} chunk a potentially incomplete portion of the manifest
   */


  var _proto = Parser.prototype;

  _proto.push = function push(chunk) {
    this.lineStream.push(chunk);
  }
  /**
   * Flush any remaining input. This can be handy if the last line of an M3U8
   * manifest did not contain a trailing newline but the file has been
   * completely received.
   */
  ;

  _proto.end = function end() {
    // flush any buffered input
    this.lineStream.push('\n');
  }
  /**
   * Add an additional parser for non-standard tags
   *
   * @param {Object}   options              a map of options for the added parser
   * @param {RegExp}   options.expression   a regular expression to match the custom header
   * @param {string}   options.type         the type to register to the output
   * @param {Function} [options.dataParser] function to parse the line into an object
   * @param {boolean}  [options.segment]    should tag data be attached to the segment object
   */
  ;

  _proto.addParser = function addParser(options) {
    this.parseStream.addParser(options);
  }
  /**
   * Add a custom header mapper
   *
   * @param {Object}   options
   * @param {RegExp}   options.expression   a regular expression to match the custom header
   * @param {Function} options.map          function to translate tag into a different tag
   */
  ;

  _proto.addTagMapper = function addTagMapper(options) {
    this.parseStream.addTagMapper(options);
  };

  return Parser;
}(stream);

var mediaTypes = createCommonjsModule(function (module, exports) {

Object.defineProperty(exports, '__esModule', { value: true });

var MPEGURL_REGEX = /^(audio|video|application)\/(x-|vnd\.apple\.)?mpegurl/i;
var DASH_REGEX = /^application\/dash\+xml/i;
/**
 * Returns a string that describes the type of source based on a video source object's
 * media type.
 *
 * @see {@link https://dev.w3.org/html5/pf-summary/video.html#dom-source-type|Source Type}
 *
 * @param {string} type
 *        Video source object media type
 * @return {('hls'|'dash'|'vhs-json'|null)}
 *         VHS source type string
 */

var simpleTypeFromSourceType = function simpleTypeFromSourceType(type) {
  if (MPEGURL_REGEX.test(type)) {
    return 'hls';
  }

  if (DASH_REGEX.test(type)) {
    return 'dash';
  } // Denotes the special case of a manifest object passed to http-streaming instead of a
  // source URL.
  //
  // See https://en.wikipedia.org/wiki/Media_type for details on specifying media types.
  //
  // In this case, vnd stands for vendor, video.js for the organization, VHS for this
  // project, and the +json suffix identifies the structure of the media type.


  if (type === 'application/vnd.videojs.vhs+json') {
    return 'vhs-json';
  }

  return null;
};

exports.simpleTypeFromSourceType = simpleTypeFromSourceType;
});

function _interopDefault$2 (ex) { return (ex && (typeof ex === 'object') && 'default' in ex) ? ex['default'] : ex; }

var URLToolkit$1 = _interopDefault$2(urlToolkit);
var window$3 = _interopDefault$2(window_1);

var resolveUrl$1 = function resolveUrl(baseUrl, relativeUrl) {
  // return early if we don't need to resolve
  if (/^[a-z]+:/i.test(relativeUrl)) {
    return relativeUrl;
  } // if the base URL is relative then combine with the current location


  if (!/\/\//i.test(baseUrl)) {
    baseUrl = URLToolkit$1.buildAbsoluteURL(window$3.location && window$3.location.href || '', baseUrl);
  }

  return URLToolkit$1.buildAbsoluteURL(baseUrl, relativeUrl);
};

var resolveUrl_1$1 = resolveUrl$1;

function _interopDefault$3 (ex) { return (ex && (typeof ex === 'object') && 'default' in ex) ? ex['default'] : ex; }

var window$4 = _interopDefault$3(window_1);

var atob$1 = function atob(s) {
  return window$4.atob ? window$4.atob(s) : Buffer.from(s, 'base64').toString('binary');
};

function decodeB64ToUint8Array$1(b64Text) {
  var decodedString = atob$1(b64Text);
  var array = new Uint8Array(decodedString.length);

  for (var i = 0; i < decodedString.length; i++) {
    array[i] = decodedString.charCodeAt(i);
  }

  return array;
}

var decodeB64ToUint8Array_1$1 = decodeB64ToUint8Array$1;

/*! @name mpd-parser @version 0.12.0 @license Apache-2.0 */

var isObject = function isObject(obj) {
  return !!obj && typeof obj === 'object';
};

var merge = function merge() {
  for (var _len = arguments.length, objects = new Array(_len), _key = 0; _key < _len; _key++) {
    objects[_key] = arguments[_key];
  }

  return objects.reduce(function (result, source) {
    Object.keys(source).forEach(function (key) {
      if (Array.isArray(result[key]) && Array.isArray(source[key])) {
        result[key] = result[key].concat(source[key]);
      } else if (isObject(result[key]) && isObject(source[key])) {
        result[key] = merge(result[key], source[key]);
      } else {
        result[key] = source[key];
      }
    });
    return result;
  }, {});
};
var values = function values(o) {
  return Object.keys(o).map(function (k) {
    return o[k];
  });
};

var range = function range(start, end) {
  var result = [];

  for (var i = start; i < end; i++) {
    result.push(i);
  }

  return result;
};
var flatten = function flatten(lists) {
  return lists.reduce(function (x, y) {
    return x.concat(y);
  }, []);
};
var from = function from(list) {
  if (!list.length) {
    return [];
  }

  var result = [];

  for (var i = 0; i < list.length; i++) {
    result.push(list[i]);
  }

  return result;
};
var findIndexes = function findIndexes(l, key) {
  return l.reduce(function (a, e, i) {
    if (e[key]) {
      a.push(i);
    }

    return a;
  }, []);
};

var errors = {
  INVALID_NUMBER_OF_PERIOD: 'INVALID_NUMBER_OF_PERIOD',
  DASH_EMPTY_MANIFEST: 'DASH_EMPTY_MANIFEST',
  DASH_INVALID_XML: 'DASH_INVALID_XML',
  NO_BASE_URL: 'NO_BASE_URL',
  MISSING_SEGMENT_INFORMATION: 'MISSING_SEGMENT_INFORMATION',
  SEGMENT_TIME_UNSPECIFIED: 'SEGMENT_TIME_UNSPECIFIED',
  UNSUPPORTED_UTC_TIMING_SCHEME: 'UNSUPPORTED_UTC_TIMING_SCHEME'
};

/**
 * @typedef {Object} SingleUri
 * @property {string} uri - relative location of segment
 * @property {string} resolvedUri - resolved location of segment
 * @property {Object} byterange - Object containing information on how to make byte range
 *   requests following byte-range-spec per RFC2616.
 * @property {String} byterange.length - length of range request
 * @property {String} byterange.offset - byte offset of range request
 *
 * @see https://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.35.1
 */

/**
 * Converts a URLType node (5.3.9.2.3 Table 13) to a segment object
 * that conforms to how m3u8-parser is structured
 *
 * @see https://github.com/videojs/m3u8-parser
 *
 * @param {string} baseUrl - baseUrl provided by <BaseUrl> nodes
 * @param {string} source - source url for segment
 * @param {string} range - optional range used for range calls,
 *   follows  RFC 2616, Clause 14.35.1
 * @return {SingleUri} full segment information transformed into a format similar
 *   to m3u8-parser
 */

var urlTypeToSegment = function urlTypeToSegment(_ref) {
  var _ref$baseUrl = _ref.baseUrl,
      baseUrl = _ref$baseUrl === void 0 ? '' : _ref$baseUrl,
      _ref$source = _ref.source,
      source = _ref$source === void 0 ? '' : _ref$source,
      _ref$range = _ref.range,
      range = _ref$range === void 0 ? '' : _ref$range,
      _ref$indexRange = _ref.indexRange,
      indexRange = _ref$indexRange === void 0 ? '' : _ref$indexRange;
  var segment = {
    uri: source,
    resolvedUri: resolveUrl_1$1(baseUrl || '', source)
  };

  if (range || indexRange) {
    var rangeStr = range ? range : indexRange;
    var ranges = rangeStr.split('-');
    var startRange = parseInt(ranges[0], 10);
    var endRange = parseInt(ranges[1], 10); // byterange should be inclusive according to
    // RFC 2616, Clause 14.35.1

    segment.byterange = {
      length: endRange - startRange + 1,
      offset: startRange
    };
  }

  return segment;
};
var byteRangeToString = function byteRangeToString(byterange) {
  // `endRange` is one less than `offset + length` because the HTTP range
  // header uses inclusive ranges
  var endRange = byterange.offset + byterange.length - 1;
  return byterange.offset + "-" + endRange;
};

/**
 * Functions for calculating the range of available segments in static and dynamic
 * manifests.
 */

var segmentRange = {
  /**
   * Returns the entire range of available segments for a static MPD
   *
   * @param {Object} attributes
   *        Inheritied MPD attributes
   * @return {{ start: number, end: number }}
   *         The start and end numbers for available segments
   */
  static: function _static(attributes) {
    var duration = attributes.duration,
        _attributes$timescale = attributes.timescale,
        timescale = _attributes$timescale === void 0 ? 1 : _attributes$timescale,
        sourceDuration = attributes.sourceDuration;
    return {
      start: 0,
      end: Math.ceil(sourceDuration / (duration / timescale))
    };
  },

  /**
   * Returns the current live window range of available segments for a dynamic MPD
   *
   * @param {Object} attributes
   *        Inheritied MPD attributes
   * @return {{ start: number, end: number }}
   *         The start and end numbers for available segments
   */
  dynamic: function dynamic(attributes) {
    var NOW = attributes.NOW,
        clientOffset = attributes.clientOffset,
        availabilityStartTime = attributes.availabilityStartTime,
        _attributes$timescale2 = attributes.timescale,
        timescale = _attributes$timescale2 === void 0 ? 1 : _attributes$timescale2,
        duration = attributes.duration,
        _attributes$start = attributes.start,
        start = _attributes$start === void 0 ? 0 : _attributes$start,
        _attributes$minimumUp = attributes.minimumUpdatePeriod,
        minimumUpdatePeriod = _attributes$minimumUp === void 0 ? 0 : _attributes$minimumUp,
        _attributes$timeShift = attributes.timeShiftBufferDepth,
        timeShiftBufferDepth = _attributes$timeShift === void 0 ? Infinity : _attributes$timeShift;
    var now = (NOW + clientOffset) / 1000;
    var periodStartWC = availabilityStartTime + start;
    var periodEndWC = now + minimumUpdatePeriod;
    var periodDuration = periodEndWC - periodStartWC;
    var segmentCount = Math.ceil(periodDuration * timescale / duration);
    var availableStart = Math.floor((now - periodStartWC - timeShiftBufferDepth) * timescale / duration);
    var availableEnd = Math.floor((now - periodStartWC) * timescale / duration);
    return {
      start: Math.max(0, availableStart),
      end: Math.min(segmentCount, availableEnd)
    };
  }
};
/**
 * Maps a range of numbers to objects with information needed to build the corresponding
 * segment list
 *
 * @name toSegmentsCallback
 * @function
 * @param {number} number
 *        Number of the segment
 * @param {number} index
 *        Index of the number in the range list
 * @return {{ number: Number, duration: Number, timeline: Number, time: Number }}
 *         Object with segment timing and duration info
 */

/**
 * Returns a callback for Array.prototype.map for mapping a range of numbers to
 * information needed to build the segment list.
 *
 * @param {Object} attributes
 *        Inherited MPD attributes
 * @return {toSegmentsCallback}
 *         Callback map function
 */

var toSegments = function toSegments(attributes) {
  return function (number, index) {
    var duration = attributes.duration,
        _attributes$timescale3 = attributes.timescale,
        timescale = _attributes$timescale3 === void 0 ? 1 : _attributes$timescale3,
        periodIndex = attributes.periodIndex,
        _attributes$startNumb = attributes.startNumber,
        startNumber = _attributes$startNumb === void 0 ? 1 : _attributes$startNumb;
    return {
      number: startNumber + number,
      duration: duration / timescale,
      timeline: periodIndex,
      time: index * duration
    };
  };
};
/**
 * Returns a list of objects containing segment timing and duration info used for
 * building the list of segments. This uses the @duration attribute specified
 * in the MPD manifest to derive the range of segments.
 *
 * @param {Object} attributes
 *        Inherited MPD attributes
 * @return {{number: number, duration: number, time: number, timeline: number}[]}
 *         List of Objects with segment timing and duration info
 */

var parseByDuration = function parseByDuration(attributes) {
  var _attributes$type = attributes.type,
      type = _attributes$type === void 0 ? 'static' : _attributes$type,
      duration = attributes.duration,
      _attributes$timescale4 = attributes.timescale,
      timescale = _attributes$timescale4 === void 0 ? 1 : _attributes$timescale4,
      sourceDuration = attributes.sourceDuration;

  var _segmentRange$type = segmentRange[type](attributes),
      start = _segmentRange$type.start,
      end = _segmentRange$type.end;

  var segments = range(start, end).map(toSegments(attributes));

  if (type === 'static') {
    var index = segments.length - 1; // final segment may be less than full segment duration

    segments[index].duration = sourceDuration - duration / timescale * index;
  }

  return segments;
};

/**
 * Translates SegmentBase into a set of segments.
 * (DASH SPEC Section 5.3.9.3.2) contains a set of <SegmentURL> nodes.  Each
 * node should be translated into segment.
 *
 * @param {Object} attributes
 *   Object containing all inherited attributes from parent elements with attribute
 *   names as keys
 * @return {Object.<Array>} list of segments
 */

var segmentsFromBase = function segmentsFromBase(attributes) {
  var baseUrl = attributes.baseUrl,
      _attributes$initializ = attributes.initialization,
      initialization = _attributes$initializ === void 0 ? {} : _attributes$initializ,
      sourceDuration = attributes.sourceDuration,
      _attributes$indexRang = attributes.indexRange,
      indexRange = _attributes$indexRang === void 0 ? '' : _attributes$indexRang,
      duration = attributes.duration; // base url is required for SegmentBase to work, per spec (Section 5.3.9.2.1)

  if (!baseUrl) {
    throw new Error(errors.NO_BASE_URL);
  }

  var initSegment = urlTypeToSegment({
    baseUrl: baseUrl,
    source: initialization.sourceURL,
    range: initialization.range
  });
  var segment = urlTypeToSegment({
    baseUrl: baseUrl,
    source: baseUrl,
    indexRange: indexRange
  });
  segment.map = initSegment; // If there is a duration, use it, otherwise use the given duration of the source
  // (since SegmentBase is only for one total segment)

  if (duration) {
    var segmentTimeInfo = parseByDuration(attributes);

    if (segmentTimeInfo.length) {
      segment.duration = segmentTimeInfo[0].duration;
      segment.timeline = segmentTimeInfo[0].timeline;
    }
  } else if (sourceDuration) {
    segment.duration = sourceDuration;
    segment.timeline = 0;
  } // This is used for mediaSequence


  segment.number = 0;
  return [segment];
};
/**
 * Given a playlist, a sidx box, and a baseUrl, update the segment list of the playlist
 * according to the sidx information given.
 *
 * playlist.sidx has metadadata about the sidx where-as the sidx param
 * is the parsed sidx box itself.
 *
 * @param {Object} playlist the playlist to update the sidx information for
 * @param {Object} sidx the parsed sidx box
 * @return {Object} the playlist object with the updated sidx information
 */

var addSegmentsToPlaylist = function addSegmentsToPlaylist(playlist, sidx, baseUrl) {
  // Retain init segment information
  var initSegment = playlist.sidx.map ? playlist.sidx.map : null; // Retain source duration from initial master manifest parsing

  var sourceDuration = playlist.sidx.duration; // Retain source timeline

  var timeline = playlist.timeline || 0;
  var sidxByteRange = playlist.sidx.byterange;
  var sidxEnd = sidxByteRange.offset + sidxByteRange.length; // Retain timescale of the parsed sidx

  var timescale = sidx.timescale; // referenceType 1 refers to other sidx boxes

  var mediaReferences = sidx.references.filter(function (r) {
    return r.referenceType !== 1;
  });
  var segments = []; // firstOffset is the offset from the end of the sidx box

  var startIndex = sidxEnd + sidx.firstOffset;

  for (var i = 0; i < mediaReferences.length; i++) {
    var reference = sidx.references[i]; // size of the referenced (sub)segment

    var size = reference.referencedSize; // duration of the referenced (sub)segment, in  the  timescale
    // this will be converted to seconds when generating segments

    var duration = reference.subsegmentDuration; // should be an inclusive range

    var endIndex = startIndex + size - 1;
    var indexRange = startIndex + "-" + endIndex;
    var attributes = {
      baseUrl: baseUrl,
      timescale: timescale,
      timeline: timeline,
      // this is used in parseByDuration
      periodIndex: timeline,
      duration: duration,
      sourceDuration: sourceDuration,
      indexRange: indexRange
    };
    var segment = segmentsFromBase(attributes)[0];

    if (initSegment) {
      segment.map = initSegment;
    }

    segments.push(segment);
    startIndex += size;
  }

  playlist.segments = segments;
  return playlist;
};

var mergeDiscontiguousPlaylists = function mergeDiscontiguousPlaylists(playlists) {
  var mergedPlaylists = values(playlists.reduce(function (acc, playlist) {
    // assuming playlist IDs are the same across periods
    // TODO: handle multiperiod where representation sets are not the same
    // across periods
    var name = playlist.attributes.id + (playlist.attributes.lang || ''); // Periods after first

    if (acc[name]) {
      var _acc$name$segments;

      // first segment of subsequent periods signal a discontinuity
      if (playlist.segments[0]) {
        playlist.segments[0].discontinuity = true;
      }

      (_acc$name$segments = acc[name].segments).push.apply(_acc$name$segments, playlist.segments); // bubble up contentProtection, this assumes all DRM content
      // has the same contentProtection


      if (playlist.attributes.contentProtection) {
        acc[name].attributes.contentProtection = playlist.attributes.contentProtection;
      }
    } else {
      // first Period
      acc[name] = playlist;
    }

    return acc;
  }, {}));
  return mergedPlaylists.map(function (playlist) {
    playlist.discontinuityStarts = findIndexes(playlist.segments, 'discontinuity');
    return playlist;
  });
};

var addSegmentInfoFromSidx = function addSegmentInfoFromSidx(playlists, sidxMapping) {
  if (sidxMapping === void 0) {
    sidxMapping = {};
  }

  if (!Object.keys(sidxMapping).length) {
    return playlists;
  }

  for (var i in playlists) {
    var playlist = playlists[i];

    if (!playlist.sidx) {
      continue;
    }

    var sidxKey = playlist.sidx.uri + '-' + byteRangeToString(playlist.sidx.byterange);
    var sidxMatch = sidxMapping[sidxKey] && sidxMapping[sidxKey].sidx;

    if (playlist.sidx && sidxMatch) {
      addSegmentsToPlaylist(playlist, sidxMatch, playlist.sidx.resolvedUri);
    }
  }

  return playlists;
};

var formatAudioPlaylist = function formatAudioPlaylist(_ref) {
  var _attributes;

  var attributes = _ref.attributes,
      segments = _ref.segments,
      sidx = _ref.sidx;
  var playlist = {
    attributes: (_attributes = {
      NAME: attributes.id,
      BANDWIDTH: attributes.bandwidth,
      CODECS: attributes.codecs
    }, _attributes['PROGRAM-ID'] = 1, _attributes),
    uri: '',
    endList: (attributes.type || 'static') === 'static',
    timeline: attributes.periodIndex,
    resolvedUri: '',
    targetDuration: attributes.duration,
    segments: segments,
    mediaSequence: segments.length ? segments[0].number : 1
  };

  if (attributes.contentProtection) {
    playlist.contentProtection = attributes.contentProtection;
  }

  if (sidx) {
    playlist.sidx = sidx;
  }

  return playlist;
};
var formatVttPlaylist = function formatVttPlaylist(_ref2) {
  var _attributes2;

  var attributes = _ref2.attributes,
      segments = _ref2.segments;

  if (typeof segments === 'undefined') {
    // vtt tracks may use single file in BaseURL
    segments = [{
      uri: attributes.baseUrl,
      timeline: attributes.periodIndex,
      resolvedUri: attributes.baseUrl || '',
      duration: attributes.sourceDuration,
      number: 0
    }]; // targetDuration should be the same duration as the only segment

    attributes.duration = attributes.sourceDuration;
  }

  return {
    attributes: (_attributes2 = {
      NAME: attributes.id,
      BANDWIDTH: attributes.bandwidth
    }, _attributes2['PROGRAM-ID'] = 1, _attributes2),
    uri: '',
    endList: (attributes.type || 'static') === 'static',
    timeline: attributes.periodIndex,
    resolvedUri: attributes.baseUrl || '',
    targetDuration: attributes.duration,
    segments: segments,
    mediaSequence: segments.length ? segments[0].number : 1
  };
};
var organizeAudioPlaylists = function organizeAudioPlaylists(playlists, sidxMapping) {
  if (sidxMapping === void 0) {
    sidxMapping = {};
  }

  var mainPlaylist;
  var formattedPlaylists = playlists.reduce(function (a, playlist) {
    var role = playlist.attributes.role && playlist.attributes.role.value || '';
    var language = playlist.attributes.lang || '';
    var label = 'main';

    if (language) {
      var roleLabel = role ? " (" + role + ")" : '';
      label = "" + playlist.attributes.lang + roleLabel;
    } // skip if we already have the highest quality audio for a language


    if (a[label] && a[label].playlists[0].attributes.BANDWIDTH > playlist.attributes.bandwidth) {
      return a;
    }

    a[label] = {
      language: language,
      autoselect: true,
      default: role === 'main',
      playlists: addSegmentInfoFromSidx([formatAudioPlaylist(playlist)], sidxMapping),
      uri: ''
    };

    if (typeof mainPlaylist === 'undefined' && role === 'main') {
      mainPlaylist = playlist;
      mainPlaylist.default = true;
    }

    return a;
  }, {}); // if no playlists have role "main", mark the first as main

  if (!mainPlaylist) {
    var firstLabel = Object.keys(formattedPlaylists)[0];
    formattedPlaylists[firstLabel].default = true;
  }

  return formattedPlaylists;
};
var organizeVttPlaylists = function organizeVttPlaylists(playlists, sidxMapping) {
  if (sidxMapping === void 0) {
    sidxMapping = {};
  }

  return playlists.reduce(function (a, playlist) {
    var label = playlist.attributes.lang || 'text'; // skip if we already have subtitles

    if (a[label]) {
      return a;
    }

    a[label] = {
      language: label,
      default: false,
      autoselect: false,
      playlists: addSegmentInfoFromSidx([formatVttPlaylist(playlist)], sidxMapping),
      uri: ''
    };
    return a;
  }, {});
};
var formatVideoPlaylist = function formatVideoPlaylist(_ref3) {
  var _attributes3;

  var attributes = _ref3.attributes,
      segments = _ref3.segments,
      sidx = _ref3.sidx;
  var playlist = {
    attributes: (_attributes3 = {
      NAME: attributes.id,
      AUDIO: 'audio',
      SUBTITLES: 'subs',
      RESOLUTION: {
        width: attributes.width,
        height: attributes.height
      },
      CODECS: attributes.codecs,
      BANDWIDTH: attributes.bandwidth
    }, _attributes3['PROGRAM-ID'] = 1, _attributes3),
    uri: '',
    endList: (attributes.type || 'static') === 'static',
    timeline: attributes.periodIndex,
    resolvedUri: '',
    targetDuration: attributes.duration,
    segments: segments,
    mediaSequence: segments.length ? segments[0].number : 1
  };

  if (attributes.contentProtection) {
    playlist.contentProtection = attributes.contentProtection;
  }

  if (sidx) {
    playlist.sidx = sidx;
  }

  return playlist;
};
var toM3u8 = function toM3u8(dashPlaylists, locations, sidxMapping) {
  var _mediaGroups;

  if (sidxMapping === void 0) {
    sidxMapping = {};
  }

  if (!dashPlaylists.length) {
    return {};
  } // grab all master attributes


  var _dashPlaylists$0$attr = dashPlaylists[0].attributes,
      duration = _dashPlaylists$0$attr.sourceDuration,
      _dashPlaylists$0$attr2 = _dashPlaylists$0$attr.type,
      type = _dashPlaylists$0$attr2 === void 0 ? 'static' : _dashPlaylists$0$attr2,
      suggestedPresentationDelay = _dashPlaylists$0$attr.suggestedPresentationDelay,
      minimumUpdatePeriod = _dashPlaylists$0$attr.minimumUpdatePeriod;

  var videoOnly = function videoOnly(_ref4) {
    var attributes = _ref4.attributes;
    return attributes.mimeType === 'video/mp4' || attributes.contentType === 'video';
  };

  var audioOnly = function audioOnly(_ref5) {
    var attributes = _ref5.attributes;
    return attributes.mimeType === 'audio/mp4' || attributes.contentType === 'audio';
  };

  var vttOnly = function vttOnly(_ref6) {
    var attributes = _ref6.attributes;
    return attributes.mimeType === 'text/vtt' || attributes.contentType === 'text';
  };

  var videoPlaylists = mergeDiscontiguousPlaylists(dashPlaylists.filter(videoOnly)).map(formatVideoPlaylist);
  var audioPlaylists = mergeDiscontiguousPlaylists(dashPlaylists.filter(audioOnly));
  var vttPlaylists = dashPlaylists.filter(vttOnly);
  var master = {
    allowCache: true,
    discontinuityStarts: [],
    segments: [],
    endList: true,
    mediaGroups: (_mediaGroups = {
      AUDIO: {},
      VIDEO: {}
    }, _mediaGroups['CLOSED-CAPTIONS'] = {}, _mediaGroups.SUBTITLES = {}, _mediaGroups),
    uri: '',
    duration: duration,
    playlists: addSegmentInfoFromSidx(videoPlaylists, sidxMapping)
  };

  if (minimumUpdatePeriod >= 0) {
    master.minimumUpdatePeriod = minimumUpdatePeriod * 1000;
  }

  if (locations) {
    master.locations = locations;
  }

  if (type === 'dynamic') {
    master.suggestedPresentationDelay = suggestedPresentationDelay;
  }

  if (audioPlaylists.length) {
    master.mediaGroups.AUDIO.audio = organizeAudioPlaylists(audioPlaylists, sidxMapping);
  }

  if (vttPlaylists.length) {
    master.mediaGroups.SUBTITLES.subs = organizeVttPlaylists(vttPlaylists, sidxMapping);
  }

  return master;
};

/**
 * Calculates the R (repetition) value for a live stream (for the final segment
 * in a manifest where the r value is negative 1)
 *
 * @param {Object} attributes
 *        Object containing all inherited attributes from parent elements with attribute
 *        names as keys
 * @param {number} time
 *        current time (typically the total time up until the final segment)
 * @param {number} duration
 *        duration property for the given <S />
 *
 * @return {number}
 *        R value to reach the end of the given period
 */
var getLiveRValue = function getLiveRValue(attributes, time, duration) {
  var NOW = attributes.NOW,
      clientOffset = attributes.clientOffset,
      availabilityStartTime = attributes.availabilityStartTime,
      _attributes$timescale = attributes.timescale,
      timescale = _attributes$timescale === void 0 ? 1 : _attributes$timescale,
      _attributes$start = attributes.start,
      start = _attributes$start === void 0 ? 0 : _attributes$start,
      _attributes$minimumUp = attributes.minimumUpdatePeriod,
      minimumUpdatePeriod = _attributes$minimumUp === void 0 ? 0 : _attributes$minimumUp;
  var now = (NOW + clientOffset) / 1000;
  var periodStartWC = availabilityStartTime + start;
  var periodEndWC = now + minimumUpdatePeriod;
  var periodDuration = periodEndWC - periodStartWC;
  return Math.ceil((periodDuration * timescale - time) / duration);
};
/**
 * Uses information provided by SegmentTemplate.SegmentTimeline to determine segment
 * timing and duration
 *
 * @param {Object} attributes
 *        Object containing all inherited attributes from parent elements with attribute
 *        names as keys
 * @param {Object[]} segmentTimeline
 *        List of objects representing the attributes of each S element contained within
 *
 * @return {{number: number, duration: number, time: number, timeline: number}[]}
 *         List of Objects with segment timing and duration info
 */


var parseByTimeline = function parseByTimeline(attributes, segmentTimeline) {
  var _attributes$type = attributes.type,
      type = _attributes$type === void 0 ? 'static' : _attributes$type,
      _attributes$minimumUp2 = attributes.minimumUpdatePeriod,
      minimumUpdatePeriod = _attributes$minimumUp2 === void 0 ? 0 : _attributes$minimumUp2,
      _attributes$media = attributes.media,
      media = _attributes$media === void 0 ? '' : _attributes$media,
      sourceDuration = attributes.sourceDuration,
      _attributes$timescale2 = attributes.timescale,
      timescale = _attributes$timescale2 === void 0 ? 1 : _attributes$timescale2,
      _attributes$startNumb = attributes.startNumber,
      startNumber = _attributes$startNumb === void 0 ? 1 : _attributes$startNumb,
      timeline = attributes.periodIndex;
  var segments = [];
  var time = -1;

  for (var sIndex = 0; sIndex < segmentTimeline.length; sIndex++) {
    var S = segmentTimeline[sIndex];
    var duration = S.d;
    var repeat = S.r || 0;
    var segmentTime = S.t || 0;

    if (time < 0) {
      // first segment
      time = segmentTime;
    }

    if (segmentTime && segmentTime > time) {
      // discontinuity
      // TODO: How to handle this type of discontinuity
      // timeline++ here would treat it like HLS discontuity and content would
      // get appended without gap
      // E.G.
      //  <S t="0" d="1" />
      //  <S d="1" />
      //  <S d="1" />
      //  <S t="5" d="1" />
      // would have $Time$ values of [0, 1, 2, 5]
      // should this be appened at time positions [0, 1, 2, 3],(#EXT-X-DISCONTINUITY)
      // or [0, 1, 2, gap, gap, 5]? (#EXT-X-GAP)
      // does the value of sourceDuration consider this when calculating arbitrary
      // negative @r repeat value?
      // E.G. Same elements as above with this added at the end
      //  <S d="1" r="-1" />
      //  with a sourceDuration of 10
      // Would the 2 gaps be included in the time duration calculations resulting in
      // 8 segments with $Time$ values of [0, 1, 2, 5, 6, 7, 8, 9] or 10 segments
      // with $Time$ values of [0, 1, 2, 5, 6, 7, 8, 9, 10, 11] ?
      time = segmentTime;
    }

    var count = void 0;

    if (repeat < 0) {
      var nextS = sIndex + 1;

      if (nextS === segmentTimeline.length) {
        // last segment
        if (type === 'dynamic' && minimumUpdatePeriod > 0 && media.indexOf('$Number$') > 0) {
          count = getLiveRValue(attributes, time, duration);
        } else {
          // TODO: This may be incorrect depending on conclusion of TODO above
          count = (sourceDuration * timescale - time) / duration;
        }
      } else {
        count = (segmentTimeline[nextS].t - time) / duration;
      }
    } else {
      count = repeat + 1;
    }

    var end = startNumber + segments.length + count;
    var number = startNumber + segments.length;

    while (number < end) {
      segments.push({
        number: number,
        duration: duration / timescale,
        time: time,
        timeline: timeline
      });
      time += duration;
      number++;
    }
  }

  return segments;
};

var identifierPattern = /\$([A-z]*)(?:(%0)([0-9]+)d)?\$/g;
/**
 * Replaces template identifiers with corresponding values. To be used as the callback
 * for String.prototype.replace
 *
 * @name replaceCallback
 * @function
 * @param {string} match
 *        Entire match of identifier
 * @param {string} identifier
 *        Name of matched identifier
 * @param {string} format
 *        Format tag string. Its presence indicates that padding is expected
 * @param {string} width
 *        Desired length of the replaced value. Values less than this width shall be left
 *        zero padded
 * @return {string}
 *         Replacement for the matched identifier
 */

/**
 * Returns a function to be used as a callback for String.prototype.replace to replace
 * template identifiers
 *
 * @param {Obect} values
 *        Object containing values that shall be used to replace known identifiers
 * @param {number} values.RepresentationID
 *        Value of the Representation@id attribute
 * @param {number} values.Number
 *        Number of the corresponding segment
 * @param {number} values.Bandwidth
 *        Value of the Representation@bandwidth attribute.
 * @param {number} values.Time
 *        Timestamp value of the corresponding segment
 * @return {replaceCallback}
 *         Callback to be used with String.prototype.replace to replace identifiers
 */

var identifierReplacement = function identifierReplacement(values) {
  return function (match, identifier, format, width) {
    if (match === '$$') {
      // escape sequence
      return '$';
    }

    if (typeof values[identifier] === 'undefined') {
      return match;
    }

    var value = '' + values[identifier];

    if (identifier === 'RepresentationID') {
      // Format tag shall not be present with RepresentationID
      return value;
    }

    if (!format) {
      width = 1;
    } else {
      width = parseInt(width, 10);
    }

    if (value.length >= width) {
      return value;
    }

    return "" + new Array(width - value.length + 1).join('0') + value;
  };
};
/**
 * Constructs a segment url from a template string
 *
 * @param {string} url
 *        Template string to construct url from
 * @param {Obect} values
 *        Object containing values that shall be used to replace known identifiers
 * @param {number} values.RepresentationID
 *        Value of the Representation@id attribute
 * @param {number} values.Number
 *        Number of the corresponding segment
 * @param {number} values.Bandwidth
 *        Value of the Representation@bandwidth attribute.
 * @param {number} values.Time
 *        Timestamp value of the corresponding segment
 * @return {string}
 *         Segment url with identifiers replaced
 */

var constructTemplateUrl = function constructTemplateUrl(url, values) {
  return url.replace(identifierPattern, identifierReplacement(values));
};
/**
 * Generates a list of objects containing timing and duration information about each
 * segment needed to generate segment uris and the complete segment object
 *
 * @param {Object} attributes
 *        Object containing all inherited attributes from parent elements with attribute
 *        names as keys
 * @param {Object[]|undefined} segmentTimeline
 *        List of objects representing the attributes of each S element contained within
 *        the SegmentTimeline element
 * @return {{number: number, duration: number, time: number, timeline: number}[]}
 *         List of Objects with segment timing and duration info
 */

var parseTemplateInfo = function parseTemplateInfo(attributes, segmentTimeline) {
  if (!attributes.duration && !segmentTimeline) {
    // if neither @duration or SegmentTimeline are present, then there shall be exactly
    // one media segment
    return [{
      number: attributes.startNumber || 1,
      duration: attributes.sourceDuration,
      time: 0,
      timeline: attributes.periodIndex
    }];
  }

  if (attributes.duration) {
    return parseByDuration(attributes);
  }

  return parseByTimeline(attributes, segmentTimeline);
};
/**
 * Generates a list of segments using information provided by the SegmentTemplate element
 *
 * @param {Object} attributes
 *        Object containing all inherited attributes from parent elements with attribute
 *        names as keys
 * @param {Object[]|undefined} segmentTimeline
 *        List of objects representing the attributes of each S element contained within
 *        the SegmentTimeline element
 * @return {Object[]}
 *         List of segment objects
 */

var segmentsFromTemplate = function segmentsFromTemplate(attributes, segmentTimeline) {
  var templateValues = {
    RepresentationID: attributes.id,
    Bandwidth: attributes.bandwidth || 0
  };
  var _attributes$initializ = attributes.initialization,
      initialization = _attributes$initializ === void 0 ? {
    sourceURL: '',
    range: ''
  } : _attributes$initializ;
  var mapSegment = urlTypeToSegment({
    baseUrl: attributes.baseUrl,
    source: constructTemplateUrl(initialization.sourceURL, templateValues),
    range: initialization.range
  });
  var segments = parseTemplateInfo(attributes, segmentTimeline);
  return segments.map(function (segment) {
    templateValues.Number = segment.number;
    templateValues.Time = segment.time;
    var uri = constructTemplateUrl(attributes.media || '', templateValues);
    return {
      uri: uri,
      timeline: segment.timeline,
      duration: segment.duration,
      resolvedUri: resolveUrl_1$1(attributes.baseUrl || '', uri),
      map: mapSegment,
      number: segment.number
    };
  });
};

/**
 * Converts a <SegmentUrl> (of type URLType from the DASH spec 5.3.9.2 Table 14)
 * to an object that matches the output of a segment in videojs/mpd-parser
 *
 * @param {Object} attributes
 *   Object containing all inherited attributes from parent elements with attribute
 *   names as keys
 * @param {Object} segmentUrl
 *   <SegmentURL> node to translate into a segment object
 * @return {Object} translated segment object
 */

var SegmentURLToSegmentObject = function SegmentURLToSegmentObject(attributes, segmentUrl) {
  var baseUrl = attributes.baseUrl,
      _attributes$initializ = attributes.initialization,
      initialization = _attributes$initializ === void 0 ? {} : _attributes$initializ;
  var initSegment = urlTypeToSegment({
    baseUrl: baseUrl,
    source: initialization.sourceURL,
    range: initialization.range
  });
  var segment = urlTypeToSegment({
    baseUrl: baseUrl,
    source: segmentUrl.media,
    range: segmentUrl.mediaRange
  });
  segment.map = initSegment;
  return segment;
};
/**
 * Generates a list of segments using information provided by the SegmentList element
 * SegmentList (DASH SPEC Section 5.3.9.3.2) contains a set of <SegmentURL> nodes.  Each
 * node should be translated into segment.
 *
 * @param {Object} attributes
 *   Object containing all inherited attributes from parent elements with attribute
 *   names as keys
 * @param {Object[]|undefined} segmentTimeline
 *        List of objects representing the attributes of each S element contained within
 *        the SegmentTimeline element
 * @return {Object.<Array>} list of segments
 */


var segmentsFromList = function segmentsFromList(attributes, segmentTimeline) {
  var duration = attributes.duration,
      _attributes$segmentUr = attributes.segmentUrls,
      segmentUrls = _attributes$segmentUr === void 0 ? [] : _attributes$segmentUr; // Per spec (5.3.9.2.1) no way to determine segment duration OR
  // if both SegmentTimeline and @duration are defined, it is outside of spec.

  if (!duration && !segmentTimeline || duration && segmentTimeline) {
    throw new Error(errors.SEGMENT_TIME_UNSPECIFIED);
  }

  var segmentUrlMap = segmentUrls.map(function (segmentUrlObject) {
    return SegmentURLToSegmentObject(attributes, segmentUrlObject);
  });
  var segmentTimeInfo;

  if (duration) {
    segmentTimeInfo = parseByDuration(attributes);
  }

  if (segmentTimeline) {
    segmentTimeInfo = parseByTimeline(attributes, segmentTimeline);
  }

  var segments = segmentTimeInfo.map(function (segmentTime, index) {
    if (segmentUrlMap[index]) {
      var segment = segmentUrlMap[index];
      segment.timeline = segmentTime.timeline;
      segment.duration = segmentTime.duration;
      segment.number = segmentTime.number;
      return segment;
    } // Since we're mapping we should get rid of any blank segments (in case
    // the given SegmentTimeline is handling for more elements than we have
    // SegmentURLs for).

  }).filter(function (segment) {
    return segment;
  });
  return segments;
};

var generateSegments = function generateSegments(_ref) {
  var attributes = _ref.attributes,
      segmentInfo = _ref.segmentInfo;
  var segmentAttributes;
  var segmentsFn;

  if (segmentInfo.template) {
    segmentsFn = segmentsFromTemplate;
    segmentAttributes = merge(attributes, segmentInfo.template);
  } else if (segmentInfo.base) {
    segmentsFn = segmentsFromBase;
    segmentAttributes = merge(attributes, segmentInfo.base);
  } else if (segmentInfo.list) {
    segmentsFn = segmentsFromList;
    segmentAttributes = merge(attributes, segmentInfo.list);
  }

  var segmentsInfo = {
    attributes: attributes
  };

  if (!segmentsFn) {
    return segmentsInfo;
  }

  var segments = segmentsFn(segmentAttributes, segmentInfo.timeline); // The @duration attribute will be used to determin the playlist's targetDuration which
  // must be in seconds. Since we've generated the segment list, we no longer need
  // @duration to be in @timescale units, so we can convert it here.

  if (segmentAttributes.duration) {
    var _segmentAttributes = segmentAttributes,
        duration = _segmentAttributes.duration,
        _segmentAttributes$ti = _segmentAttributes.timescale,
        timescale = _segmentAttributes$ti === void 0 ? 1 : _segmentAttributes$ti;
    segmentAttributes.duration = duration / timescale;
  } else if (segments.length) {
    // if there is no @duration attribute, use the largest segment duration as
    // as target duration
    segmentAttributes.duration = segments.reduce(function (max, segment) {
      return Math.max(max, Math.ceil(segment.duration));
    }, 0);
  } else {
    segmentAttributes.duration = 0;
  }

  segmentsInfo.attributes = segmentAttributes;
  segmentsInfo.segments = segments; // This is a sidx box without actual segment information

  if (segmentInfo.base && segmentAttributes.indexRange) {
    segmentsInfo.sidx = segments[0];
    segmentsInfo.segments = [];
  }

  return segmentsInfo;
};
var toPlaylists = function toPlaylists(representations) {
  return representations.map(generateSegments);
};

var findChildren = function findChildren(element, name) {
  return from(element.childNodes).filter(function (_ref) {
    var tagName = _ref.tagName;
    return tagName === name;
  });
};
var getContent = function getContent(element) {
  return element.textContent.trim();
};

var parseDuration = function parseDuration(str) {
  var SECONDS_IN_YEAR = 365 * 24 * 60 * 60;
  var SECONDS_IN_MONTH = 30 * 24 * 60 * 60;
  var SECONDS_IN_DAY = 24 * 60 * 60;
  var SECONDS_IN_HOUR = 60 * 60;
  var SECONDS_IN_MIN = 60; // P10Y10M10DT10H10M10.1S

  var durationRegex = /P(?:(\d*)Y)?(?:(\d*)M)?(?:(\d*)D)?(?:T(?:(\d*)H)?(?:(\d*)M)?(?:([\d.]*)S)?)?/;
  var match = durationRegex.exec(str);

  if (!match) {
    return 0;
  }

  var _match$slice = match.slice(1),
      year = _match$slice[0],
      month = _match$slice[1],
      day = _match$slice[2],
      hour = _match$slice[3],
      minute = _match$slice[4],
      second = _match$slice[5];

  return parseFloat(year || 0) * SECONDS_IN_YEAR + parseFloat(month || 0) * SECONDS_IN_MONTH + parseFloat(day || 0) * SECONDS_IN_DAY + parseFloat(hour || 0) * SECONDS_IN_HOUR + parseFloat(minute || 0) * SECONDS_IN_MIN + parseFloat(second || 0);
};
var parseDate = function parseDate(str) {
  // Date format without timezone according to ISO 8601
  // YYY-MM-DDThh:mm:ss.ssssss
  var dateRegex = /^\d+-\d+-\d+T\d+:\d+:\d+(\.\d+)?$/; // If the date string does not specifiy a timezone, we must specifiy UTC. This is
  // expressed by ending with 'Z'

  if (dateRegex.test(str)) {
    str += 'Z';
  }

  return Date.parse(str);
};

var parsers = {
  /**
   * Specifies the duration of the entire Media Presentation. Format is a duration string
   * as specified in ISO 8601
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {number}
   *         The duration in seconds
   */
  mediaPresentationDuration: function mediaPresentationDuration(value) {
    return parseDuration(value);
  },

  /**
   * Specifies the Segment availability start time for all Segments referred to in this
   * MPD. For a dynamic manifest, it specifies the anchor for the earliest availability
   * time. Format is a date string as specified in ISO 8601
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {number}
   *         The date as seconds from unix epoch
   */
  availabilityStartTime: function availabilityStartTime(value) {
    return parseDate(value) / 1000;
  },

  /**
   * Specifies the smallest period between potential changes to the MPD. Format is a
   * duration string as specified in ISO 8601
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {number}
   *         The duration in seconds
   */
  minimumUpdatePeriod: function minimumUpdatePeriod(value) {
    return parseDuration(value);
  },

  /**
   * Specifies the suggested presentation delay. Format is a
   * duration string as specified in ISO 8601
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {number}
   *         The duration in seconds
   */
  suggestedPresentationDelay: function suggestedPresentationDelay(value) {
    return parseDuration(value);
  },

  /**
   * specifices the type of mpd. Can be either "static" or "dynamic"
   *
   * @param {string} value
   *        value of attribute as a string
   *
   * @return {string}
   *         The type as a string
   */
  type: function type(value) {
    return value;
  },

  /**
   * Specifies the duration of the smallest time shifting buffer for any Representation
   * in the MPD. Format is a duration string as specified in ISO 8601
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {number}
   *         The duration in seconds
   */
  timeShiftBufferDepth: function timeShiftBufferDepth(value) {
    return parseDuration(value);
  },

  /**
   * Specifies the PeriodStart time of the Period relative to the availabilityStarttime.
   * Format is a duration string as specified in ISO 8601
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {number}
   *         The duration in seconds
   */
  start: function start(value) {
    return parseDuration(value);
  },

  /**
   * Specifies the width of the visual presentation
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {number}
   *         The parsed width
   */
  width: function width(value) {
    return parseInt(value, 10);
  },

  /**
   * Specifies the height of the visual presentation
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {number}
   *         The parsed height
   */
  height: function height(value) {
    return parseInt(value, 10);
  },

  /**
   * Specifies the bitrate of the representation
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {number}
   *         The parsed bandwidth
   */
  bandwidth: function bandwidth(value) {
    return parseInt(value, 10);
  },

  /**
   * Specifies the number of the first Media Segment in this Representation in the Period
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {number}
   *         The parsed number
   */
  startNumber: function startNumber(value) {
    return parseInt(value, 10);
  },

  /**
   * Specifies the timescale in units per seconds
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {number}
   *         The aprsed timescale
   */
  timescale: function timescale(value) {
    return parseInt(value, 10);
  },

  /**
   * Specifies the constant approximate Segment duration
   * NOTE: The <Period> element also contains an @duration attribute. This duration
   *       specifies the duration of the Period. This attribute is currently not
   *       supported by the rest of the parser, however we still check for it to prevent
   *       errors.
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {number}
   *         The parsed duration
   */
  duration: function duration(value) {
    var parsedValue = parseInt(value, 10);

    if (isNaN(parsedValue)) {
      return parseDuration(value);
    }

    return parsedValue;
  },

  /**
   * Specifies the Segment duration, in units of the value of the @timescale.
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {number}
   *         The parsed duration
   */
  d: function d(value) {
    return parseInt(value, 10);
  },

  /**
   * Specifies the MPD start time, in @timescale units, the first Segment in the series
   * starts relative to the beginning of the Period
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {number}
   *         The parsed time
   */
  t: function t(value) {
    return parseInt(value, 10);
  },

  /**
   * Specifies the repeat count of the number of following contiguous Segments with the
   * same duration expressed by the value of @d
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {number}
   *         The parsed number
   */
  r: function r(value) {
    return parseInt(value, 10);
  },

  /**
   * Default parser for all other attributes. Acts as a no-op and just returns the value
   * as a string
   *
   * @param {string} value
   *        value of attribute as a string
   * @return {string}
   *         Unparsed value
   */
  DEFAULT: function DEFAULT(value) {
    return value;
  }
};
/**
 * Gets all the attributes and values of the provided node, parses attributes with known
 * types, and returns an object with attribute names mapped to values.
 *
 * @param {Node} el
 *        The node to parse attributes from
 * @return {Object}
 *         Object with all attributes of el parsed
 */

var parseAttributes$1 = function parseAttributes(el) {
  if (!(el && el.attributes)) {
    return {};
  }

  return from(el.attributes).reduce(function (a, e) {
    var parseFn = parsers[e.name] || parsers.DEFAULT;
    a[e.name] = parseFn(e.value);
    return a;
  }, {});
};

var keySystemsMap = {
  'urn:uuid:1077efec-c0b2-4d02-ace3-3c1e52e2fb4b': 'org.w3.clearkey',
  'urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed': 'com.widevine.alpha',
  'urn:uuid:9a04f079-9840-4286-ab92-e65be0885f95': 'com.microsoft.playready',
  'urn:uuid:f239e769-efa3-4850-9c16-a903c6932efb': 'com.adobe.primetime'
};
/**
 * Builds a list of urls that is the product of the reference urls and BaseURL values
 *
 * @param {string[]} referenceUrls
 *        List of reference urls to resolve to
 * @param {Node[]} baseUrlElements
 *        List of BaseURL nodes from the mpd
 * @return {string[]}
 *         List of resolved urls
 */

var buildBaseUrls = function buildBaseUrls(referenceUrls, baseUrlElements) {
  if (!baseUrlElements.length) {
    return referenceUrls;
  }

  return flatten(referenceUrls.map(function (reference) {
    return baseUrlElements.map(function (baseUrlElement) {
      return resolveUrl_1$1(reference, getContent(baseUrlElement));
    });
  }));
};
/**
 * Contains all Segment information for its containing AdaptationSet
 *
 * @typedef {Object} SegmentInformation
 * @property {Object|undefined} template
 *           Contains the attributes for the SegmentTemplate node
 * @property {Object[]|undefined} timeline
 *           Contains a list of atrributes for each S node within the SegmentTimeline node
 * @property {Object|undefined} list
 *           Contains the attributes for the SegmentList node
 * @property {Object|undefined} base
 *           Contains the attributes for the SegmentBase node
 */

/**
 * Returns all available Segment information contained within the AdaptationSet node
 *
 * @param {Node} adaptationSet
 *        The AdaptationSet node to get Segment information from
 * @return {SegmentInformation}
 *         The Segment information contained within the provided AdaptationSet
 */

var getSegmentInformation = function getSegmentInformation(adaptationSet) {
  var segmentTemplate = findChildren(adaptationSet, 'SegmentTemplate')[0];
  var segmentList = findChildren(adaptationSet, 'SegmentList')[0];
  var segmentUrls = segmentList && findChildren(segmentList, 'SegmentURL').map(function (s) {
    return merge({
      tag: 'SegmentURL'
    }, parseAttributes$1(s));
  });
  var segmentBase = findChildren(adaptationSet, 'SegmentBase')[0];
  var segmentTimelineParentNode = segmentList || segmentTemplate;
  var segmentTimeline = segmentTimelineParentNode && findChildren(segmentTimelineParentNode, 'SegmentTimeline')[0];
  var segmentInitializationParentNode = segmentList || segmentBase || segmentTemplate;
  var segmentInitialization = segmentInitializationParentNode && findChildren(segmentInitializationParentNode, 'Initialization')[0]; // SegmentTemplate is handled slightly differently, since it can have both
  // @initialization and an <Initialization> node.  @initialization can be templated,
  // while the node can have a url and range specified.  If the <SegmentTemplate> has
  // both @initialization and an <Initialization> subelement we opt to override with
  // the node, as this interaction is not defined in the spec.

  var template = segmentTemplate && parseAttributes$1(segmentTemplate);

  if (template && segmentInitialization) {
    template.initialization = segmentInitialization && parseAttributes$1(segmentInitialization);
  } else if (template && template.initialization) {
    // If it is @initialization we convert it to an object since this is the format that
    // later functions will rely on for the initialization segment.  This is only valid
    // for <SegmentTemplate>
    template.initialization = {
      sourceURL: template.initialization
    };
  }

  var segmentInfo = {
    template: template,
    timeline: segmentTimeline && findChildren(segmentTimeline, 'S').map(function (s) {
      return parseAttributes$1(s);
    }),
    list: segmentList && merge(parseAttributes$1(segmentList), {
      segmentUrls: segmentUrls,
      initialization: parseAttributes$1(segmentInitialization)
    }),
    base: segmentBase && merge(parseAttributes$1(segmentBase), {
      initialization: parseAttributes$1(segmentInitialization)
    })
  };
  Object.keys(segmentInfo).forEach(function (key) {
    if (!segmentInfo[key]) {
      delete segmentInfo[key];
    }
  });
  return segmentInfo;
};
/**
 * Contains Segment information and attributes needed to construct a Playlist object
 * from a Representation
 *
 * @typedef {Object} RepresentationInformation
 * @property {SegmentInformation} segmentInfo
 *           Segment information for this Representation
 * @property {Object} attributes
 *           Inherited attributes for this Representation
 */

/**
 * Maps a Representation node to an object containing Segment information and attributes
 *
 * @name inheritBaseUrlsCallback
 * @function
 * @param {Node} representation
 *        Representation node from the mpd
 * @return {RepresentationInformation}
 *         Representation information needed to construct a Playlist object
 */

/**
 * Returns a callback for Array.prototype.map for mapping Representation nodes to
 * Segment information and attributes using inherited BaseURL nodes.
 *
 * @param {Object} adaptationSetAttributes
 *        Contains attributes inherited by the AdaptationSet
 * @param {string[]} adaptationSetBaseUrls
 *        Contains list of resolved base urls inherited by the AdaptationSet
 * @param {SegmentInformation} adaptationSetSegmentInfo
 *        Contains Segment information for the AdaptationSet
 * @return {inheritBaseUrlsCallback}
 *         Callback map function
 */

var inheritBaseUrls = function inheritBaseUrls(adaptationSetAttributes, adaptationSetBaseUrls, adaptationSetSegmentInfo) {
  return function (representation) {
    var repBaseUrlElements = findChildren(representation, 'BaseURL');
    var repBaseUrls = buildBaseUrls(adaptationSetBaseUrls, repBaseUrlElements);
    var attributes = merge(adaptationSetAttributes, parseAttributes$1(representation));
    var representationSegmentInfo = getSegmentInformation(representation);
    return repBaseUrls.map(function (baseUrl) {
      return {
        segmentInfo: merge(adaptationSetSegmentInfo, representationSegmentInfo),
        attributes: merge(attributes, {
          baseUrl: baseUrl
        })
      };
    });
  };
};
/**
 * Tranforms a series of content protection nodes to
 * an object containing pssh data by key system
 *
 * @param {Node[]} contentProtectionNodes
 *        Content protection nodes
 * @return {Object}
 *        Object containing pssh data by key system
 */

var generateKeySystemInformation = function generateKeySystemInformation(contentProtectionNodes) {
  return contentProtectionNodes.reduce(function (acc, node) {
    var attributes = parseAttributes$1(node);
    var keySystem = keySystemsMap[attributes.schemeIdUri];

    if (keySystem) {
      acc[keySystem] = {
        attributes: attributes
      };
      var psshNode = findChildren(node, 'cenc:pssh')[0];

      if (psshNode) {
        var pssh = getContent(psshNode);
        var psshBuffer = pssh && decodeB64ToUint8Array_1$1(pssh);
        acc[keySystem].pssh = psshBuffer;
      }
    }

    return acc;
  }, {});
};
/**
 * Maps an AdaptationSet node to a list of Representation information objects
 *
 * @name toRepresentationsCallback
 * @function
 * @param {Node} adaptationSet
 *        AdaptationSet node from the mpd
 * @return {RepresentationInformation[]}
 *         List of objects containing Representaion information
 */

/**
 * Returns a callback for Array.prototype.map for mapping AdaptationSet nodes to a list of
 * Representation information objects
 *
 * @param {Object} periodAttributes
 *        Contains attributes inherited by the Period
 * @param {string[]} periodBaseUrls
 *        Contains list of resolved base urls inherited by the Period
 * @param {string[]} periodSegmentInfo
 *        Contains Segment Information at the period level
 * @return {toRepresentationsCallback}
 *         Callback map function
 */


var toRepresentations = function toRepresentations(periodAttributes, periodBaseUrls, periodSegmentInfo) {
  return function (adaptationSet) {
    var adaptationSetAttributes = parseAttributes$1(adaptationSet);
    var adaptationSetBaseUrls = buildBaseUrls(periodBaseUrls, findChildren(adaptationSet, 'BaseURL'));
    var role = findChildren(adaptationSet, 'Role')[0];
    var roleAttributes = {
      role: parseAttributes$1(role)
    };
    var attrs = merge(periodAttributes, adaptationSetAttributes, roleAttributes);
    var contentProtection = generateKeySystemInformation(findChildren(adaptationSet, 'ContentProtection'));

    if (Object.keys(contentProtection).length) {
      attrs = merge(attrs, {
        contentProtection: contentProtection
      });
    }

    var segmentInfo = getSegmentInformation(adaptationSet);
    var representations = findChildren(adaptationSet, 'Representation');
    var adaptationSetSegmentInfo = merge(periodSegmentInfo, segmentInfo);
    return flatten(representations.map(inheritBaseUrls(attrs, adaptationSetBaseUrls, adaptationSetSegmentInfo)));
  };
};
/**
 * Maps an Period node to a list of Representation inforamtion objects for all
 * AdaptationSet nodes contained within the Period
 *
 * @name toAdaptationSetsCallback
 * @function
 * @param {Node} period
 *        Period node from the mpd
 * @param {number} periodIndex
 *        Index of the Period within the mpd
 * @return {RepresentationInformation[]}
 *         List of objects containing Representaion information
 */

/**
 * Returns a callback for Array.prototype.map for mapping Period nodes to a list of
 * Representation information objects
 *
 * @param {Object} mpdAttributes
 *        Contains attributes inherited by the mpd
 * @param {string[]} mpdBaseUrls
 *        Contains list of resolved base urls inherited by the mpd
 * @return {toAdaptationSetsCallback}
 *         Callback map function
 */

var toAdaptationSets = function toAdaptationSets(mpdAttributes, mpdBaseUrls) {
  return function (period, index) {
    var periodBaseUrls = buildBaseUrls(mpdBaseUrls, findChildren(period, 'BaseURL'));
    var periodAtt = parseAttributes$1(period);
    var parsedPeriodId = parseInt(periodAtt.id, 10); // fallback to mapping index if Period@id is not a number

    var periodIndex = window_1.isNaN(parsedPeriodId) ? index : parsedPeriodId;
    var periodAttributes = merge(mpdAttributes, {
      periodIndex: periodIndex
    });
    var adaptationSets = findChildren(period, 'AdaptationSet');
    var periodSegmentInfo = getSegmentInformation(period);
    return flatten(adaptationSets.map(toRepresentations(periodAttributes, periodBaseUrls, periodSegmentInfo)));
  };
};
/**
 * Traverses the mpd xml tree to generate a list of Representation information objects
 * that have inherited attributes from parent nodes
 *
 * @param {Node} mpd
 *        The root node of the mpd
 * @param {Object} options
 *        Available options for inheritAttributes
 * @param {string} options.manifestUri
 *        The uri source of the mpd
 * @param {number} options.NOW
 *        Current time per DASH IOP.  Default is current time in ms since epoch
 * @param {number} options.clientOffset
 *        Client time difference from NOW (in milliseconds)
 * @return {RepresentationInformation[]}
 *         List of objects containing Representation information
 */

var inheritAttributes = function inheritAttributes(mpd, options) {
  if (options === void 0) {
    options = {};
  }

  var _options = options,
      _options$manifestUri = _options.manifestUri,
      manifestUri = _options$manifestUri === void 0 ? '' : _options$manifestUri,
      _options$NOW = _options.NOW,
      NOW = _options$NOW === void 0 ? Date.now() : _options$NOW,
      _options$clientOffset = _options.clientOffset,
      clientOffset = _options$clientOffset === void 0 ? 0 : _options$clientOffset;
  var periods = findChildren(mpd, 'Period');

  if (!periods.length) {
    throw new Error(errors.INVALID_NUMBER_OF_PERIOD);
  }

  var locations = findChildren(mpd, 'Location');
  var mpdAttributes = parseAttributes$1(mpd);
  var mpdBaseUrls = buildBaseUrls([manifestUri], findChildren(mpd, 'BaseURL'));
  mpdAttributes.sourceDuration = mpdAttributes.mediaPresentationDuration || 0;
  mpdAttributes.NOW = NOW;
  mpdAttributes.clientOffset = clientOffset;

  if (locations.length) {
    mpdAttributes.locations = locations.map(getContent);
  }

  return {
    locations: mpdAttributes.locations,
    representationInfo: flatten(periods.map(toAdaptationSets(mpdAttributes, mpdBaseUrls)))
  };
};

var stringToMpdXml = function stringToMpdXml(manifestString) {
  if (manifestString === '') {
    throw new Error(errors.DASH_EMPTY_MANIFEST);
  }

  var parser = new domParser.DOMParser();
  var xml = parser.parseFromString(manifestString, 'application/xml');
  var mpd = xml && xml.documentElement.tagName === 'MPD' ? xml.documentElement : null;

  if (!mpd || mpd && mpd.getElementsByTagName('parsererror').length > 0) {
    throw new Error(errors.DASH_INVALID_XML);
  }

  return mpd;
};

/**
 * Parses the manifest for a UTCTiming node, returning the nodes attributes if found
 *
 * @param {string} mpd
 *        XML string of the MPD manifest
 * @return {Object|null}
 *         Attributes of UTCTiming node specified in the manifest. Null if none found
 */

var parseUTCTimingScheme = function parseUTCTimingScheme(mpd) {
  var UTCTimingNode = findChildren(mpd, 'UTCTiming')[0];

  if (!UTCTimingNode) {
    return null;
  }

  var attributes = parseAttributes$1(UTCTimingNode);

  switch (attributes.schemeIdUri) {
    case 'urn:mpeg:dash:utc:http-head:2014':
    case 'urn:mpeg:dash:utc:http-head:2012':
      attributes.method = 'HEAD';
      break;

    case 'urn:mpeg:dash:utc:http-xsdate:2014':
    case 'urn:mpeg:dash:utc:http-iso:2014':
    case 'urn:mpeg:dash:utc:http-xsdate:2012':
    case 'urn:mpeg:dash:utc:http-iso:2012':
      attributes.method = 'GET';
      break;

    case 'urn:mpeg:dash:utc:direct:2014':
    case 'urn:mpeg:dash:utc:direct:2012':
      attributes.method = 'DIRECT';
      attributes.value = Date.parse(attributes.value);
      break;

    case 'urn:mpeg:dash:utc:http-ntp:2014':
    case 'urn:mpeg:dash:utc:ntp:2014':
    case 'urn:mpeg:dash:utc:sntp:2014':
    default:
      throw new Error(errors.UNSUPPORTED_UTC_TIMING_SCHEME);
  }

  return attributes;
};

var parse = function parse(manifestString, options) {
  if (options === void 0) {
    options = {};
  }

  var parsedManifestInfo = inheritAttributes(stringToMpdXml(manifestString), options);
  var playlists = toPlaylists(parsedManifestInfo.representationInfo);
  return toM3u8(playlists, parsedManifestInfo.locations, options.sidxMapping);
};
/**
 * Parses the manifest for a UTCTiming node, returning the nodes attributes if found
 *
 * @param {string} manifestString
 *        XML string of the MPD manifest
 * @return {Object|null}
 *         Attributes of UTCTiming node specified in the manifest. Null if none found
 */


var parseUTCTiming = function parseUTCTiming(manifestString) {
  return parseUTCTimingScheme(stringToMpdXml(manifestString));
};

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

var parseSidx_1 = parseSidx;

var containers = createCommonjsModule(function (module, exports) {

Object.defineProperty(exports, '__esModule', { value: true });

var toUint8 = function toUint8(bytes) {
  return bytes instanceof Uint8Array ? bytes : new Uint8Array(bytes && bytes.buffer || bytes, bytes && bytes.byteOffset || 0, bytes && bytes.byteLength || 0);
};
var bytesToString = function bytesToString(bytes) {
  if (!bytes) {
    return '';
  }

  bytes = Array.prototype.slice.call(bytes);
  var string = String.fromCharCode.apply(null, toUint8(bytes));

  try {
    return decodeURIComponent(escape(string));
  } catch (e) {// if decodeURIComponent/escape fails, we are dealing with partial
    // or full non string data. Just return the potentially garbled string.
  }

  return string;
};

var getId3Offset = function getId3Offset(bytes) {
  bytes = toUint8(bytes);

  if (bytes.length < 10 || bytesToString(bytes.subarray(0, 3)) !== 'ID3') {
    return 0;
  }

  var returnSize = bytes[6] << 21 | bytes[7] << 14 | bytes[8] << 7 | bytes[9];
  var flags = bytes[5];
  var footerPresent = (flags & 16) >> 4;

  if (footerPresent) {
    return returnSize + 20;
  }

  return returnSize + 10;
};
var isLikely = {
  aac: function aac(bytes) {
    var offset = getId3Offset(bytes);
    return bytes.length >= offset + 2 && (bytes[offset] & 0xFF) === 0xFF && (bytes[offset + 1] & 0xE0) === 0xE0 && (bytes[offset + 1] & 0x16) === 0x10;
  },
  mp3: function mp3(bytes) {
    var offset = getId3Offset(bytes);
    return bytes.length >= offset + 2 && (bytes[offset] & 0xFF) === 0xFF && (bytes[offset + 1] & 0xE0) === 0xE0 && (bytes[offset + 1] & 0x06) === 0x02;
  },
  webm: function webm(bytes) {
    return bytes.length >= 4 && (bytes[0] & 0xFF) === 0x1A && (bytes[1] & 0xFF) === 0x45 && (bytes[2] & 0xFF) === 0xDF && (bytes[3] & 0xFF) === 0xA3;
  },
  mp4: function mp4(bytes) {
    return bytes.length >= 8 && /^(f|s)typ$/.test(bytesToString(bytes.subarray(4, 8))) && // not 3gp data
    !/^ftyp3g$/.test(bytesToString(bytes.subarray(4, 10)));
  },
  '3gp': function gp(bytes) {
    return bytes.length >= 10 && /^ftyp3g$/.test(bytesToString(bytes.subarray(4, 10)));
  },
  ts: function ts(bytes) {
    if (bytes.length < 189 && bytes.length >= 1) {
      return bytes[0] === 0x47;
    }

    var i = 0; // check the first 376 bytes for two matching sync bytes

    while (i + 188 < bytes.length && i < 188) {
      if (bytes[i] === 0x47 && bytes[i + 188] === 0x47) {
        return true;
      }

      i += 1;
    }

    return false;
  },
  flac: function flac(bytes) {
    return bytes.length >= 4 && /^fLaC$/.test(bytesToString(bytes.subarray(0, 4)));
  },
  ogg: function ogg(bytes) {
    return bytes.length >= 4 && /^OggS$/.test(bytesToString(bytes.subarray(0, 4)));
  }
}; // get all the isLikely functions
// but make sure 'ts' is at the bottom
// as it is the least specific

var isLikelyTypes = Object.keys(isLikely) // remove ts
.filter(function (t) {
  return t !== 'ts';
}) // add it back to the bottom
.concat('ts'); // make sure we are dealing with uint8 data.

isLikelyTypes.forEach(function (type) {
  var isLikelyFn = isLikely[type];

  isLikely[type] = function (bytes) {
    return isLikelyFn(toUint8(bytes));
  };
}); // A useful list of file signatures can be found here
// https://en.wikipedia.org/wiki/List_of_file_signatures

var detectContainerForBytes = function detectContainerForBytes(bytes) {
  bytes = toUint8(bytes);

  for (var i = 0; i < isLikelyTypes.length; i++) {
    var type = isLikelyTypes[i];

    if (isLikely[type](bytes)) {
      return type;
    }
  }

  return '';
}; // fmp4 is not a container

var isLikelyFmp4MediaSegment = function isLikelyFmp4MediaSegment(bytes) {
  bytes = toUint8(bytes);
  var i = 0;

  while (i < bytes.length) {
    var size = (bytes[i] << 24 | bytes[i + 1] << 16 | bytes[i + 2] << 8 | bytes[i + 3]) >>> 0;
    var type = bytesToString(bytes.subarray(i + 4, i + 8));

    if (type === 'moof') {
      return true;
    }

    if (size === 0 || size + i > bytes.length) {
      i = bytes.length;
    } else {
      i += size;
    }
  }

  return false;
};

exports.detectContainerForBytes = detectContainerForBytes;
exports.getId3Offset = getId3Offset;
exports.isLikely = isLikely;
exports.isLikelyFmp4MediaSegment = isLikelyFmp4MediaSegment;
});

var byteHelpers = createCommonjsModule(function (module, exports) {

Object.defineProperty(exports, '__esModule', { value: true });

var isTypedArray = function isTypedArray(obj) {
  return ArrayBuffer.isView(obj);
};
var toUint8 = function toUint8(bytes) {
  return bytes instanceof Uint8Array ? bytes : new Uint8Array(bytes && bytes.buffer || bytes, bytes && bytes.byteOffset || 0, bytes && bytes.byteLength || 0);
};
var bytesToString = function bytesToString(bytes) {
  if (!bytes) {
    return '';
  }

  bytes = Array.prototype.slice.call(bytes);
  var string = String.fromCharCode.apply(null, toUint8(bytes));

  try {
    return decodeURIComponent(escape(string));
  } catch (e) {// if decodeURIComponent/escape fails, we are dealing with partial
    // or full non string data. Just return the potentially garbled string.
  }

  return string;
};
var stringToBytes = function stringToBytes(string, stringIsBytes) {
  if (stringIsBytes === void 0) {
    stringIsBytes = false;
  }

  var bytes = [];

  if (typeof string !== 'string' && string && typeof string.toString === 'function') {
    string = string.toString();
  }

  if (typeof string !== 'string') {
    return bytes;
  } // If the string already is bytes, we don't have to do this


  if (!stringIsBytes) {
    string = unescape(encodeURIComponent(string));
  }

  return string.split('').map(function (s) {
    return s.charCodeAt(0) & 0xFF;
  });
};
var concatTypedArrays = function concatTypedArrays() {
  for (var _len = arguments.length, buffers = new Array(_len), _key = 0; _key < _len; _key++) {
    buffers[_key] = arguments[_key];
  }

  var totalLength = buffers.reduce(function (total, buf) {
    var len = buf && (buf.byteLength || buf.length);
    total += len || 0;
    return total;
  }, 0);
  var tempBuffer = new Uint8Array(totalLength);
  var offset = 0;
  buffers.forEach(function (buf) {
    buf = toUint8(buf);
    tempBuffer.set(buf, offset);
    offset += buf.byteLength;
  });
  return tempBuffer;
};

exports.bytesToString = bytesToString;
exports.concatTypedArrays = concatTypedArrays;
exports.isTypedArray = isTypedArray;
exports.stringToBytes = stringToBytes;
exports.toUint8 = toUint8;
});

/**
 * mux.js
 *
 * Copyright (c) Brightcove
 * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
 */

var streamTypes = {
  H264_STREAM_TYPE: 0x1B,
  ADTS_STREAM_TYPE: 0x0F,
  METADATA_STREAM_TYPE: 0x15
};

/**
 * mux.js
 *
 * Copyright (c) Brightcove
 * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
 *
 * A lightweight readable stream implemention that handles event dispatching.
 * Objects that inherit from streams should call init in their constructors.
 */

var Stream$1 = function() {
  this.init = function() {
    var listeners = {};
    /**
     * Add a listener for a specified event type.
     * @param type {string} the event name
     * @param listener {function} the callback to be invoked when an event of
     * the specified type occurs
     */
    this.on = function(type, listener) {
      if (!listeners[type]) {
        listeners[type] = [];
      }
      listeners[type] = listeners[type].concat(listener);
    };
    /**
     * Remove a listener for a specified event type.
     * @param type {string} the event name
     * @param listener {function} a function previously registered for this
     * type of event through `on`
     */
    this.off = function(type, listener) {
      var index;
      if (!listeners[type]) {
        return false;
      }
      index = listeners[type].indexOf(listener);
      listeners[type] = listeners[type].slice();
      listeners[type].splice(index, 1);
      return index > -1;
    };
    /**
     * Trigger an event of the specified type on this stream. Any additional
     * arguments to this function are passed as parameters to event listeners.
     * @param type {string} the event name
     */
    this.trigger = function(type) {
      var callbacks, i, length, args;
      callbacks = listeners[type];
      if (!callbacks) {
        return;
      }
      // Slicing the arguments on every invocation of this method
      // can add a significant amount of overhead. Avoid the
      // intermediate object creation for the common case of a
      // single callback argument
      if (arguments.length === 2) {
        length = callbacks.length;
        for (i = 0; i < length; ++i) {
          callbacks[i].call(this, arguments[1]);
        }
      } else {
        args = [];
        i = arguments.length;
        for (i = 1; i < arguments.length; ++i) {
          args.push(arguments[i]);
        }
        length = callbacks.length;
        for (i = 0; i < length; ++i) {
          callbacks[i].apply(this, args);
        }
      }
    };
    /**
     * Destroys the stream and cleans up.
     */
    this.dispose = function() {
      listeners = {};
    };
  };
};

/**
 * Forwards all `data` events on this stream to the destination stream. The
 * destination stream should provide a method `push` to receive the data
 * events as they arrive.
 * @param destination {stream} the stream that will receive all `data` events
 * @param autoFlush {boolean} if false, we will not call `flush` on the destination
 *                            when the current stream emits a 'done' event
 * @see http://nodejs.org/api/stream.html#stream_readable_pipe_destination_options
 */
Stream$1.prototype.pipe = function(destination) {
  this.on('data', function(data) {
    destination.push(data);
  });

  this.on('done', function(flushSource) {
    destination.flush(flushSource);
  });

  this.on('partialdone', function(flushSource) {
    destination.partialFlush(flushSource);
  });

  this.on('endedtimeline', function(flushSource) {
    destination.endTimeline(flushSource);
  });

  this.on('reset', function(flushSource) {
    destination.reset(flushSource);
  });

  return destination;
};

// Default stream functions that are expected to be overridden to perform
// actual work. These are provided by the prototype as a sort of no-op
// implementation so that we don't have to check for their existence in the
// `pipe` function above.
Stream$1.prototype.push = function(data) {
  this.trigger('data', data);
};

Stream$1.prototype.flush = function(flushSource) {
  this.trigger('done', flushSource);
};

Stream$1.prototype.partialFlush = function(flushSource) {
  this.trigger('partialdone', flushSource);
};

Stream$1.prototype.endTimeline = function(flushSource) {
  this.trigger('endedtimeline', flushSource);
};

Stream$1.prototype.reset = function(flushSource) {
  this.trigger('reset', flushSource);
};

var stream$1 = Stream$1;

var MAX_TS = 8589934592;

var RO_THRESH = 4294967296;

var TYPE_SHARED = 'shared';

var handleRollover = function(value, reference) {
  var direction = 1;

  if (value > reference) {
    // If the current timestamp value is greater than our reference timestamp and we detect a
    // timestamp rollover, this means the roll over is happening in the opposite direction.
    // Example scenario: Enter a long stream/video just after a rollover occurred. The reference
    // point will be set to a small number, e.g. 1. The user then seeks backwards over the
    // rollover point. In loading this segment, the timestamp values will be very large,
    // e.g. 2^33 - 1. Since this comes before the data we loaded previously, we want to adjust
    // the time stamp to be `value - 2^33`.
    direction = -1;
  }

  // Note: A seek forwards or back that is greater than the RO_THRESH (2^32, ~13 hours) will
  // cause an incorrect adjustment.
  while (Math.abs(reference - value) > RO_THRESH) {
    value += (direction * MAX_TS);
  }

  return value;
};

var TimestampRolloverStream = function(type) {
  var lastDTS, referenceDTS;

  TimestampRolloverStream.prototype.init.call(this);

  // The "shared" type is used in cases where a stream will contain muxed
  // video and audio. We could use `undefined` here, but having a string
  // makes debugging a little clearer.
  this.type_ = type || TYPE_SHARED;

  this.push = function(data) {

    // Any "shared" rollover streams will accept _all_ data. Otherwise,
    // streams will only accept data that matches their type.
    if (this.type_ !== TYPE_SHARED && data.type !== this.type_) {
      return;
    }

    if (referenceDTS === undefined) {
      referenceDTS = data.dts;
    }

    data.dts = handleRollover(data.dts, referenceDTS);
    data.pts = handleRollover(data.pts, referenceDTS);

    lastDTS = data.dts;

    this.trigger('data', data);
  };

  this.flush = function() {
    referenceDTS = lastDTS;
    this.trigger('done');
  };

  this.endTimeline = function() {
    this.flush();
    this.trigger('endedtimeline');
  };

  this.discontinuity = function() {
    referenceDTS = void 0;
    lastDTS = void 0;
  };

  this.reset = function() {
    this.discontinuity();
    this.trigger('reset');
  };
};

TimestampRolloverStream.prototype = new stream$1();

var timestampRolloverStream = {
  TimestampRolloverStream: TimestampRolloverStream,
  handleRollover: handleRollover
};

var parsePid = function(packet) {
  var pid = packet[1] & 0x1f;
  pid <<= 8;
  pid |= packet[2];
  return pid;
};

var parsePayloadUnitStartIndicator = function(packet) {
  return !!(packet[1] & 0x40);
};

var parseAdaptionField = function(packet) {
  var offset = 0;
  // if an adaption field is present, its length is specified by the
  // fifth byte of the TS packet header. The adaptation field is
  // used to add stuffing to PES packets that don't fill a complete
  // TS packet, and to specify some forms of timing and control data
  // that we do not currently use.
  if (((packet[3] & 0x30) >>> 4) > 0x01) {
    offset += packet[4] + 1;
  }
  return offset;
};

var parseType = function(packet, pmtPid) {
  var pid = parsePid(packet);
  if (pid === 0) {
    return 'pat';
  } else if (pid === pmtPid) {
    return 'pmt';
  } else if (pmtPid) {
    return 'pes';
  }
  return null;
};

var parsePat = function(packet) {
  var pusi = parsePayloadUnitStartIndicator(packet);
  var offset = 4 + parseAdaptionField(packet);

  if (pusi) {
    offset += packet[offset] + 1;
  }

  return (packet[offset + 10] & 0x1f) << 8 | packet[offset + 11];
};

var parsePmt = function(packet) {
  var programMapTable = {};
  var pusi = parsePayloadUnitStartIndicator(packet);
  var payloadOffset = 4 + parseAdaptionField(packet);

  if (pusi) {
    payloadOffset += packet[payloadOffset] + 1;
  }

  // PMTs can be sent ahead of the time when they should actually
  // take effect. We don't believe this should ever be the case
  // for HLS but we'll ignore "forward" PMT declarations if we see
  // them. Future PMT declarations have the current_next_indicator
  // set to zero.
  if (!(packet[payloadOffset + 5] & 0x01)) {
    return;
  }

  var sectionLength, tableEnd, programInfoLength;
  // the mapping table ends at the end of the current section
  sectionLength = (packet[payloadOffset + 1] & 0x0f) << 8 | packet[payloadOffset + 2];
  tableEnd = 3 + sectionLength - 4;

  // to determine where the table is, we have to figure out how
  // long the program info descriptors are
  programInfoLength = (packet[payloadOffset + 10] & 0x0f) << 8 | packet[payloadOffset + 11];

  // advance the offset to the first entry in the mapping table
  var offset = 12 + programInfoLength;
  while (offset < tableEnd) {
    var i = payloadOffset + offset;
    // add an entry that maps the elementary_pid to the stream_type
    programMapTable[(packet[i + 1] & 0x1F) << 8 | packet[i + 2]] = packet[i];

    // move to the next table entry
    // skip past the elementary stream descriptors, if present
    offset += ((packet[i + 3] & 0x0F) << 8 | packet[i + 4]) + 5;
  }
  return programMapTable;
};

var parsePesType = function(packet, programMapTable) {
  var pid = parsePid(packet);
  var type = programMapTable[pid];
  switch (type) {
    case streamTypes.H264_STREAM_TYPE:
      return 'video';
    case streamTypes.ADTS_STREAM_TYPE:
      return 'audio';
    case streamTypes.METADATA_STREAM_TYPE:
      return 'timed-metadata';
    default:
      return null;
  }
};

var parsePesTime = function(packet) {
  var pusi = parsePayloadUnitStartIndicator(packet);
  if (!pusi) {
    return null;
  }

  var offset = 4 + parseAdaptionField(packet);

  if (offset >= packet.byteLength) {
    // From the H 222.0 MPEG-TS spec
    // "For transport stream packets carrying PES packets, stuffing is needed when there
    //  is insufficient PES packet data to completely fill the transport stream packet
    //  payload bytes. Stuffing is accomplished by defining an adaptation field longer than
    //  the sum of the lengths of the data elements in it, so that the payload bytes
    //  remaining after the adaptation field exactly accommodates the available PES packet
    //  data."
    //
    // If the offset is >= the length of the packet, then the packet contains no data
    // and instead is just adaption field stuffing bytes
    return null;
  }

  var pes = null;
  var ptsDtsFlags;

  // PES packets may be annotated with a PTS value, or a PTS value
  // and a DTS value. Determine what combination of values is
  // available to work with.
  ptsDtsFlags = packet[offset + 7];

  // PTS and DTS are normally stored as a 33-bit number.  Javascript
  // performs all bitwise operations on 32-bit integers but javascript
  // supports a much greater range (52-bits) of integer using standard
  // mathematical operations.
  // We construct a 31-bit value using bitwise operators over the 31
  // most significant bits and then multiply by 4 (equal to a left-shift
  // of 2) before we add the final 2 least significant bits of the
  // timestamp (equal to an OR.)
  if (ptsDtsFlags & 0xC0) {
    pes = {};
    // the PTS and DTS are not written out directly. For information
    // on how they are encoded, see
    // http://dvd.sourceforge.net/dvdinfo/pes-hdr.html
    pes.pts = (packet[offset + 9] & 0x0E) << 27 |
      (packet[offset + 10] & 0xFF) << 20 |
      (packet[offset + 11] & 0xFE) << 12 |
      (packet[offset + 12] & 0xFF) <<  5 |
      (packet[offset + 13] & 0xFE) >>>  3;
    pes.pts *= 4; // Left shift by 2
    pes.pts += (packet[offset + 13] & 0x06) >>> 1; // OR by the two LSBs
    pes.dts = pes.pts;
    if (ptsDtsFlags & 0x40) {
      pes.dts = (packet[offset + 14] & 0x0E) << 27 |
        (packet[offset + 15] & 0xFF) << 20 |
        (packet[offset + 16] & 0xFE) << 12 |
        (packet[offset + 17] & 0xFF) << 5 |
        (packet[offset + 18] & 0xFE) >>> 3;
      pes.dts *= 4; // Left shift by 2
      pes.dts += (packet[offset + 18] & 0x06) >>> 1; // OR by the two LSBs
    }
  }
  return pes;
};

var parseNalUnitType = function(type) {
  switch (type) {
    case 0x05:
      return 'slice_layer_without_partitioning_rbsp_idr';
    case 0x06:
      return 'sei_rbsp';
    case 0x07:
      return 'seq_parameter_set_rbsp';
    case 0x08:
      return 'pic_parameter_set_rbsp';
    case 0x09:
      return 'access_unit_delimiter_rbsp';
    default:
      return null;
  }
};

var videoPacketContainsKeyFrame = function(packet) {
  var offset = 4 + parseAdaptionField(packet);
  var frameBuffer = packet.subarray(offset);
  var frameI = 0;
  var frameSyncPoint = 0;
  var foundKeyFrame = false;
  var nalType;

  // advance the sync point to a NAL start, if necessary
  for (; frameSyncPoint < frameBuffer.byteLength - 3; frameSyncPoint++) {
    if (frameBuffer[frameSyncPoint + 2] === 1) {
      // the sync point is properly aligned
      frameI = frameSyncPoint + 5;
      break;
    }
  }

  while (frameI < frameBuffer.byteLength) {
    // look at the current byte to determine if we've hit the end of
    // a NAL unit boundary
    switch (frameBuffer[frameI]) {
    case 0:
      // skip past non-sync sequences
      if (frameBuffer[frameI - 1] !== 0) {
        frameI += 2;
        break;
      } else if (frameBuffer[frameI - 2] !== 0) {
        frameI++;
        break;
      }

      if (frameSyncPoint + 3 !== frameI - 2) {
        nalType = parseNalUnitType(frameBuffer[frameSyncPoint + 3] & 0x1f);
        if (nalType === 'slice_layer_without_partitioning_rbsp_idr') {
          foundKeyFrame = true;
        }
      }

      // drop trailing zeroes
      do {
        frameI++;
      } while (frameBuffer[frameI] !== 1 && frameI < frameBuffer.length);
      frameSyncPoint = frameI - 2;
      frameI += 3;
      break;
    case 1:
      // skip past non-sync sequences
      if (frameBuffer[frameI - 1] !== 0 ||
          frameBuffer[frameI - 2] !== 0) {
        frameI += 3;
        break;
      }

      nalType = parseNalUnitType(frameBuffer[frameSyncPoint + 3] & 0x1f);
      if (nalType === 'slice_layer_without_partitioning_rbsp_idr') {
        foundKeyFrame = true;
      }
      frameSyncPoint = frameI - 2;
      frameI += 3;
      break;
    default:
      // the current byte isn't a one or zero, so it cannot be part
      // of a sync sequence
      frameI += 3;
      break;
    }
  }
  frameBuffer = frameBuffer.subarray(frameSyncPoint);
  frameI -= frameSyncPoint;
  frameSyncPoint = 0;
  // parse the final nal
  if (frameBuffer && frameBuffer.byteLength > 3) {
    nalType = parseNalUnitType(frameBuffer[frameSyncPoint + 3] & 0x1f);
    if (nalType === 'slice_layer_without_partitioning_rbsp_idr') {
      foundKeyFrame = true;
    }
  }

  return foundKeyFrame;
};


var probe = {
  parseType: parseType,
  parsePat: parsePat,
  parsePmt: parsePmt,
  parsePayloadUnitStartIndicator: parsePayloadUnitStartIndicator,
  parsePesType: parsePesType,
  parsePesTime: parsePesTime,
  videoPacketContainsKeyFrame: videoPacketContainsKeyFrame
};

/**
 * mux.js
 *
 * Copyright (c) Brightcove
 * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
 *
 * Utilities to detect basic properties and metadata about Aac data.
 */

var ADTS_SAMPLING_FREQUENCIES = [
  96000,
  88200,
  64000,
  48000,
  44100,
  32000,
  24000,
  22050,
  16000,
  12000,
  11025,
  8000,
  7350
];

var parseId3TagSize = function(header, byteIndex) {
  var
    returnSize = (header[byteIndex + 6] << 21) |
                 (header[byteIndex + 7] << 14) |
                 (header[byteIndex + 8] << 7) |
                 (header[byteIndex + 9]),
    flags = header[byteIndex + 5],
    footerPresent = (flags & 16) >> 4;

  if (footerPresent) {
    return returnSize + 20;
  }
  return returnSize + 10;
};

// TODO: use vhs-utils
var isLikelyAacData = function(data) {
  var offset = 0;

  if (data.length > 10 && (data[0] === 'I'.charCodeAt(0)) &&
      (data[1] === 'D'.charCodeAt(0)) &&
      (data[2] === '3'.charCodeAt(0))) {
    offset = parseId3TagSize(data, 0);
  }

  return data.length >= offset + 2 &&
    (data[offset] & 0xFF) === 0xFF &&
    (data[offset + 1] & 0xF0) === 0xF0 &&
    // verify that the 2 layer bits are 0, aka this
    // is not mp3 data but aac data.
    (data[offset + 1] & 0x16) === 0x10;
};

var parseSyncSafeInteger = function(data) {
  return (data[0] << 21) |
          (data[1] << 14) |
          (data[2] << 7) |
          (data[3]);
};

// return a percent-encoded representation of the specified byte range
// @see http://en.wikipedia.org/wiki/Percent-encoding
var percentEncode = function(bytes, start, end) {
  var i, result = '';
  for (i = start; i < end; i++) {
    result += '%' + ('00' + bytes[i].toString(16)).slice(-2);
  }
  return result;
};

// return the string representation of the specified byte range,
// interpreted as ISO-8859-1.
var parseIso88591 = function(bytes, start, end) {
  return unescape(percentEncode(bytes, start, end)); // jshint ignore:line
};

var parseAdtsSize = function(header, byteIndex) {
  var
    lowThree = (header[byteIndex + 5] & 0xE0) >> 5,
    middle = header[byteIndex + 4] << 3,
    highTwo = header[byteIndex + 3] & 0x3 << 11;

  return (highTwo | middle) | lowThree;
};

var parseType$1 = function(header, byteIndex) {
  if ((header[byteIndex] === 'I'.charCodeAt(0)) &&
      (header[byteIndex + 1] === 'D'.charCodeAt(0)) &&
      (header[byteIndex + 2] === '3'.charCodeAt(0))) {
    return 'timed-metadata';
  } else if ((header[byteIndex] & 0xff === 0xff) &&
             ((header[byteIndex + 1] & 0xf0) === 0xf0)) {
    return 'audio';
  }
  return null;
};

var parseSampleRate = function(packet) {
  var i = 0;

  while (i + 5 < packet.length) {
    if (packet[i] !== 0xFF || (packet[i + 1] & 0xF6) !== 0xF0) {
      // If a valid header was not found,  jump one forward and attempt to
      // find a valid ADTS header starting at the next byte
      i++;
      continue;
    }
    return ADTS_SAMPLING_FREQUENCIES[(packet[i + 2] & 0x3c) >>> 2];
  }

  return null;
};

var parseAacTimestamp = function(packet) {
  var frameStart, frameSize, frame, frameHeader;

  // find the start of the first frame and the end of the tag
  frameStart = 10;
  if (packet[5] & 0x40) {
    // advance the frame start past the extended header
    frameStart += 4; // header size field
    frameStart += parseSyncSafeInteger(packet.subarray(10, 14));
  }

  // parse one or more ID3 frames
  // http://id3.org/id3v2.3.0#ID3v2_frame_overview
  do {
    // determine the number of bytes in this frame
    frameSize = parseSyncSafeInteger(packet.subarray(frameStart + 4, frameStart + 8));
    if (frameSize < 1) {
      return null;
    }
    frameHeader = String.fromCharCode(packet[frameStart],
                                      packet[frameStart + 1],
                                      packet[frameStart + 2],
                                      packet[frameStart + 3]);

    if (frameHeader === 'PRIV') {
      frame = packet.subarray(frameStart + 10, frameStart + frameSize + 10);

      for (var i = 0; i < frame.byteLength; i++) {
        if (frame[i] === 0) {
          var owner = parseIso88591(frame, 0, i);
          if (owner === 'com.apple.streaming.transportStreamTimestamp') {
            var d = frame.subarray(i + 1);
            var size = ((d[3] & 0x01)  << 30) |
                       (d[4]  << 22) |
                       (d[5] << 14) |
                       (d[6] << 6) |
                       (d[7] >>> 2);
            size *= 4;
            size += d[7] & 0x03;

            return size;
          }
          break;
        }
      }
    }

    frameStart += 10; // advance past the frame header
    frameStart += frameSize; // advance past the frame body
  } while (frameStart < packet.byteLength);
  return null;
};

var utils = {
  isLikelyAacData: isLikelyAacData,
  parseId3TagSize: parseId3TagSize,
  parseAdtsSize: parseAdtsSize,
  parseType: parseType$1,
  parseSampleRate: parseSampleRate,
  parseAacTimestamp: parseAacTimestamp
};

/**
 * mux.js
 *
 * Copyright (c) Brightcove
 * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
 */
var
  ONE_SECOND_IN_TS = 90000, // 90kHz clock
  secondsToVideoTs,
  secondsToAudioTs,
  videoTsToSeconds,
  audioTsToSeconds,
  audioTsToVideoTs,
  videoTsToAudioTs,
  metadataTsToSeconds;

secondsToVideoTs = function(seconds) {
  return seconds * ONE_SECOND_IN_TS;
};

secondsToAudioTs = function(seconds, sampleRate) {
  return seconds * sampleRate;
};

videoTsToSeconds = function(timestamp) {
  return timestamp / ONE_SECOND_IN_TS;
};

audioTsToSeconds = function(timestamp, sampleRate) {
  return timestamp / sampleRate;
};

audioTsToVideoTs = function(timestamp, sampleRate) {
  return secondsToVideoTs(audioTsToSeconds(timestamp, sampleRate));
};

videoTsToAudioTs = function(timestamp, sampleRate) {
  return secondsToAudioTs(videoTsToSeconds(timestamp), sampleRate);
};

/**
 * Adjust ID3 tag or caption timing information by the timeline pts values
 * (if keepOriginalTimestamps is false) and convert to seconds
 */
metadataTsToSeconds = function(timestamp, timelineStartPts, keepOriginalTimestamps) {
  return videoTsToSeconds(keepOriginalTimestamps ? timestamp : timestamp - timelineStartPts);
};

var clock = {
  ONE_SECOND_IN_TS: ONE_SECOND_IN_TS,
  secondsToVideoTs: secondsToVideoTs,
  secondsToAudioTs: secondsToAudioTs,
  videoTsToSeconds: videoTsToSeconds,
  audioTsToSeconds: audioTsToSeconds,
  audioTsToVideoTs: audioTsToVideoTs,
  videoTsToAudioTs: videoTsToAudioTs,
  metadataTsToSeconds: metadataTsToSeconds
};

var handleRollover$1 = timestampRolloverStream.handleRollover;
var probe$1 = {};
probe$1.ts = probe;
probe$1.aac = utils;
var ONE_SECOND_IN_TS$1 = clock.ONE_SECOND_IN_TS;

var
  MP2T_PACKET_LENGTH = 188, // bytes
  SYNC_BYTE = 0x47;

/**
 * walks through segment data looking for pat and pmt packets to parse out
 * program map table information
 */
var parsePsi_ = function(bytes, pmt) {
  var
    startIndex = 0,
    endIndex = MP2T_PACKET_LENGTH,
    packet, type;

  while (endIndex < bytes.byteLength) {
    // Look for a pair of start and end sync bytes in the data..
    if (bytes[startIndex] === SYNC_BYTE && bytes[endIndex] === SYNC_BYTE) {
      // We found a packet
      packet = bytes.subarray(startIndex, endIndex);
      type = probe$1.ts.parseType(packet, pmt.pid);

      switch (type) {
        case 'pat':
          if (!pmt.pid) {
            pmt.pid = probe$1.ts.parsePat(packet);
          }
          break;
        case 'pmt':
          if (!pmt.table) {
            pmt.table = probe$1.ts.parsePmt(packet);
          }
          break;
      }

      // Found the pat and pmt, we can stop walking the segment
      if (pmt.pid && pmt.table) {
        return;
      }

      startIndex += MP2T_PACKET_LENGTH;
      endIndex += MP2T_PACKET_LENGTH;
      continue;
    }

    // If we get here, we have somehow become de-synchronized and we need to step
    // forward one byte at a time until we find a pair of sync bytes that denote
    // a packet
    startIndex++;
    endIndex++;
  }
};

/**
 * walks through the segment data from the start and end to get timing information
 * for the first and last audio pes packets
 */
var parseAudioPes_ = function(bytes, pmt, result) {
  var
    startIndex = 0,
    endIndex = MP2T_PACKET_LENGTH,
    packet, type, pesType, pusi, parsed;

  var endLoop = false;

  // Start walking from start of segment to get first audio packet
  while (endIndex <= bytes.byteLength) {
    // Look for a pair of start and end sync bytes in the data..
    if (bytes[startIndex] === SYNC_BYTE &&
        (bytes[endIndex] === SYNC_BYTE || endIndex === bytes.byteLength)) {
      // We found a packet
      packet = bytes.subarray(startIndex, endIndex);
      type = probe$1.ts.parseType(packet, pmt.pid);

      switch (type) {
        case 'pes':
          pesType = probe$1.ts.parsePesType(packet, pmt.table);
          pusi = probe$1.ts.parsePayloadUnitStartIndicator(packet);
          if (pesType === 'audio' && pusi) {
            parsed = probe$1.ts.parsePesTime(packet);
            if (parsed) {
              parsed.type = 'audio';
              result.audio.push(parsed);
              endLoop = true;
            }
          }
          break;
      }

      if (endLoop) {
        break;
      }

      startIndex += MP2T_PACKET_LENGTH;
      endIndex += MP2T_PACKET_LENGTH;
      continue;
    }

    // If we get here, we have somehow become de-synchronized and we need to step
    // forward one byte at a time until we find a pair of sync bytes that denote
    // a packet
    startIndex++;
    endIndex++;
  }

  // Start walking from end of segment to get last audio packet
  endIndex = bytes.byteLength;
  startIndex = endIndex - MP2T_PACKET_LENGTH;
  endLoop = false;
  while (startIndex >= 0) {
    // Look for a pair of start and end sync bytes in the data..
    if (bytes[startIndex] === SYNC_BYTE &&
        (bytes[endIndex] === SYNC_BYTE || endIndex === bytes.byteLength)) {
      // We found a packet
      packet = bytes.subarray(startIndex, endIndex);
      type = probe$1.ts.parseType(packet, pmt.pid);

      switch (type) {
        case 'pes':
          pesType = probe$1.ts.parsePesType(packet, pmt.table);
          pusi = probe$1.ts.parsePayloadUnitStartIndicator(packet);
          if (pesType === 'audio' && pusi) {
            parsed = probe$1.ts.parsePesTime(packet);
            if (parsed) {
              parsed.type = 'audio';
              result.audio.push(parsed);
              endLoop = true;
            }
          }
          break;
      }

      if (endLoop) {
        break;
      }

      startIndex -= MP2T_PACKET_LENGTH;
      endIndex -= MP2T_PACKET_LENGTH;
      continue;
    }

    // If we get here, we have somehow become de-synchronized and we need to step
    // forward one byte at a time until we find a pair of sync bytes that denote
    // a packet
    startIndex--;
    endIndex--;
  }
};

/**
 * walks through the segment data from the start and end to get timing information
 * for the first and last video pes packets as well as timing information for the first
 * key frame.
 */
var parseVideoPes_ = function(bytes, pmt, result) {
  var
    startIndex = 0,
    endIndex = MP2T_PACKET_LENGTH,
    packet, type, pesType, pusi, parsed, frame, i, pes;

  var endLoop = false;

  var currentFrame = {
    data: [],
    size: 0
  };

  // Start walking from start of segment to get first video packet
  while (endIndex < bytes.byteLength) {
    // Look for a pair of start and end sync bytes in the data..
    if (bytes[startIndex] === SYNC_BYTE && bytes[endIndex] === SYNC_BYTE) {
      // We found a packet
      packet = bytes.subarray(startIndex, endIndex);
      type = probe$1.ts.parseType(packet, pmt.pid);

      switch (type) {
        case 'pes':
          pesType = probe$1.ts.parsePesType(packet, pmt.table);
          pusi = probe$1.ts.parsePayloadUnitStartIndicator(packet);
          if (pesType === 'video') {
            if (pusi && !endLoop) {
              parsed = probe$1.ts.parsePesTime(packet);
              if (parsed) {
                parsed.type = 'video';
                result.video.push(parsed);
                endLoop = true;
              }
            }
            if (!result.firstKeyFrame) {
              if (pusi) {
                if (currentFrame.size !== 0) {
                  frame = new Uint8Array(currentFrame.size);
                  i = 0;
                  while (currentFrame.data.length) {
                    pes = currentFrame.data.shift();
                    frame.set(pes, i);
                    i += pes.byteLength;
                  }
                  if (probe$1.ts.videoPacketContainsKeyFrame(frame)) {
                    var firstKeyFrame = probe$1.ts.parsePesTime(frame);

                    // PTS/DTS may not be available. Simply *not* setting
                    // the keyframe seems to work fine with HLS playback
                    // and definitely preferable to a crash with TypeError...
                    if (firstKeyFrame) {
                      result.firstKeyFrame = firstKeyFrame;
                      result.firstKeyFrame.type = 'video';
                    } else {
                      // eslint-disable-next-line
                      console.warn(
                        'Failed to extract PTS/DTS from PES at first keyframe. ' +
                        'This could be an unusual TS segment, or else mux.js did not ' +
                        'parse your TS segment correctly. If you know your TS ' +
                        'segments do contain PTS/DTS on keyframes please file a bug ' +
                        'report! You can try ffprobe to double check for yourself.'
                      );
                    }
                  }
                  currentFrame.size = 0;
                }
              }
              currentFrame.data.push(packet);
              currentFrame.size += packet.byteLength;
            }
          }
          break;
      }

      if (endLoop && result.firstKeyFrame) {
        break;
      }

      startIndex += MP2T_PACKET_LENGTH;
      endIndex += MP2T_PACKET_LENGTH;
      continue;
    }

    // If we get here, we have somehow become de-synchronized and we need to step
    // forward one byte at a time until we find a pair of sync bytes that denote
    // a packet
    startIndex++;
    endIndex++;
  }

  // Start walking from end of segment to get last video packet
  endIndex = bytes.byteLength;
  startIndex = endIndex - MP2T_PACKET_LENGTH;
  endLoop = false;
  while (startIndex >= 0) {
    // Look for a pair of start and end sync bytes in the data..
    if (bytes[startIndex] === SYNC_BYTE && bytes[endIndex] === SYNC_BYTE) {
      // We found a packet
      packet = bytes.subarray(startIndex, endIndex);
      type = probe$1.ts.parseType(packet, pmt.pid);

      switch (type) {
        case 'pes':
          pesType = probe$1.ts.parsePesType(packet, pmt.table);
          pusi = probe$1.ts.parsePayloadUnitStartIndicator(packet);
          if (pesType === 'video' && pusi) {
              parsed = probe$1.ts.parsePesTime(packet);
              if (parsed) {
                parsed.type = 'video';
                result.video.push(parsed);
                endLoop = true;
              }
          }
          break;
      }

      if (endLoop) {
        break;
      }

      startIndex -= MP2T_PACKET_LENGTH;
      endIndex -= MP2T_PACKET_LENGTH;
      continue;
    }

    // If we get here, we have somehow become de-synchronized and we need to step
    // forward one byte at a time until we find a pair of sync bytes that denote
    // a packet
    startIndex--;
    endIndex--;
  }
};

/**
 * Adjusts the timestamp information for the segment to account for
 * rollover and convert to seconds based on pes packet timescale (90khz clock)
 */
var adjustTimestamp_ = function(segmentInfo, baseTimestamp) {
  if (segmentInfo.audio && segmentInfo.audio.length) {
    var audioBaseTimestamp = baseTimestamp;
    if (typeof audioBaseTimestamp === 'undefined') {
      audioBaseTimestamp = segmentInfo.audio[0].dts;
    }
    segmentInfo.audio.forEach(function(info) {
      info.dts = handleRollover$1(info.dts, audioBaseTimestamp);
      info.pts = handleRollover$1(info.pts, audioBaseTimestamp);
      // time in seconds
      info.dtsTime = info.dts / ONE_SECOND_IN_TS$1;
      info.ptsTime = info.pts / ONE_SECOND_IN_TS$1;
    });
  }

  if (segmentInfo.video && segmentInfo.video.length) {
    var videoBaseTimestamp = baseTimestamp;
    if (typeof videoBaseTimestamp === 'undefined') {
      videoBaseTimestamp = segmentInfo.video[0].dts;
    }
    segmentInfo.video.forEach(function(info) {
      info.dts = handleRollover$1(info.dts, videoBaseTimestamp);
      info.pts = handleRollover$1(info.pts, videoBaseTimestamp);
      // time in seconds
      info.dtsTime = info.dts / ONE_SECOND_IN_TS$1;
      info.ptsTime = info.pts / ONE_SECOND_IN_TS$1;
    });
    if (segmentInfo.firstKeyFrame) {
      var frame = segmentInfo.firstKeyFrame;
      frame.dts = handleRollover$1(frame.dts, videoBaseTimestamp);
      frame.pts = handleRollover$1(frame.pts, videoBaseTimestamp);
      // time in seconds
      frame.dtsTime = frame.dts / ONE_SECOND_IN_TS$1;
      frame.ptsTime = frame.dts / ONE_SECOND_IN_TS$1;
    }
  }
};

/**
 * inspects the aac data stream for start and end time information
 */
var inspectAac_ = function(bytes) {
  var
    endLoop = false,
    audioCount = 0,
    sampleRate = null,
    timestamp = null,
    frameSize = 0,
    byteIndex = 0,
    packet;

  while (bytes.length - byteIndex >= 3) {
    var type = probe$1.aac.parseType(bytes, byteIndex);
    switch (type) {
      case 'timed-metadata':
        // Exit early because we don't have enough to parse
        // the ID3 tag header
        if (bytes.length - byteIndex < 10) {
          endLoop = true;
          break;
        }

        frameSize = probe$1.aac.parseId3TagSize(bytes, byteIndex);

        // Exit early if we don't have enough in the buffer
        // to emit a full packet
        if (frameSize > bytes.length) {
          endLoop = true;
          break;
        }
        if (timestamp === null) {
          packet = bytes.subarray(byteIndex, byteIndex + frameSize);
          timestamp = probe$1.aac.parseAacTimestamp(packet);
        }
        byteIndex += frameSize;
        break;
      case 'audio':
        // Exit early because we don't have enough to parse
        // the ADTS frame header
        if (bytes.length - byteIndex < 7) {
          endLoop = true;
          break;
        }

        frameSize = probe$1.aac.parseAdtsSize(bytes, byteIndex);

        // Exit early if we don't have enough in the buffer
        // to emit a full packet
        if (frameSize > bytes.length) {
          endLoop = true;
          break;
        }
        if (sampleRate === null) {
          packet = bytes.subarray(byteIndex, byteIndex + frameSize);
          sampleRate = probe$1.aac.parseSampleRate(packet);
        }
        audioCount++;
        byteIndex += frameSize;
        break;
      default:
        byteIndex++;
        break;
    }
    if (endLoop) {
      return null;
    }
  }
  if (sampleRate === null || timestamp === null) {
    return null;
  }

  var audioTimescale = ONE_SECOND_IN_TS$1 / sampleRate;

  var result = {
    audio: [
      {
        type: 'audio',
        dts: timestamp,
        pts: timestamp
      },
      {
        type: 'audio',
        dts: timestamp + (audioCount * 1024 * audioTimescale),
        pts: timestamp + (audioCount * 1024 * audioTimescale)
      }
    ]
  };

  return result;
};

/**
 * inspects the transport stream segment data for start and end time information
 * of the audio and video tracks (when present) as well as the first key frame's
 * start time.
 */
var inspectTs_ = function(bytes) {
  var pmt = {
    pid: null,
    table: null
  };

  var result = {};

  parsePsi_(bytes, pmt);

  for (var pid in pmt.table) {
    if (pmt.table.hasOwnProperty(pid)) {
      var type = pmt.table[pid];
      switch (type) {
        case streamTypes.H264_STREAM_TYPE:
          result.video = [];
          parseVideoPes_(bytes, pmt, result);
          if (result.video.length === 0) {
            delete result.video;
          }
          break;
        case streamTypes.ADTS_STREAM_TYPE:
          result.audio = [];
          parseAudioPes_(bytes, pmt, result);
          if (result.audio.length === 0) {
            delete result.audio;
          }
          break;
      }
    }
  }
  return result;
};

/**
 * Inspects segment byte data and returns an object with start and end timing information
 *
 * @param {Uint8Array} bytes The segment byte data
 * @param {Number} baseTimestamp Relative reference timestamp used when adjusting frame
 *  timestamps for rollover. This value must be in 90khz clock.
 * @return {Object} Object containing start and end frame timing info of segment.
 */
var inspect = function(bytes, baseTimestamp) {
  var isAacData = probe$1.aac.isLikelyAacData(bytes);

  var result;

  if (isAacData) {
    result = inspectAac_(bytes);
  } else {
    result = inspectTs_(bytes);
  }

  if (!result || (!result.audio && !result.video)) {
    return null;
  }

  adjustTimestamp_(result, baseTimestamp);

  return result;
};

var tsInspector = {
  inspect: inspect,
  parseAudioPes_: parseAudioPes_
};

/**
 * mux.js
 *
 * Copyright (c) Brightcove
 * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
 */
var toUnsigned = function(value) {
  return value >>> 0;
};

var toHexString = function(value) {
  return ('00' + value.toString(16)).slice(-2);
};

var bin = {
  toUnsigned: toUnsigned,
  toHexString: toHexString
};

var parseType$2 = function(buffer) {
  var result = '';
  result += String.fromCharCode(buffer[0]);
  result += String.fromCharCode(buffer[1]);
  result += String.fromCharCode(buffer[2]);
  result += String.fromCharCode(buffer[3]);
  return result;
};


var parseType_1 = parseType$2;

var toUnsigned$1 = bin.toUnsigned;


var findBox = function(data, path) {
  var results = [],
    i, size, type, end, subresults;

  if (!path.length) {
    // short-circuit the search for empty paths
    return null;
  }

  for (i = 0; i < data.byteLength;) {
    size = toUnsigned$1(data[i]     << 24 |
      data[i + 1] << 16 |
      data[i + 2] <<  8 |
      data[i + 3]);

    type = parseType_1(data.subarray(i + 4, i + 8));

    end = size > 1 ? i + size : data.byteLength;

    if (type === path[0]) {
      if (path.length === 1) {
        // this is the end of the path and we've found the box we were
        // looking for
        results.push(data.subarray(i + 8, end));
      } else {
        // recursively search for the next box along the path
        subresults = findBox(data.subarray(i + 8, end), path.slice(1));
        if (subresults.length) {
          results = results.concat(subresults);
        }
      }
    }
    i = end;
  }

  // we've finished searching all of data
  return results;
};

var findBox_1 = findBox;

var tfhd = function(data) {
  var
  view = new DataView(data.buffer, data.byteOffset, data.byteLength),
    result = {
      version: data[0],
      flags: new Uint8Array(data.subarray(1, 4)),
      trackId: view.getUint32(4)
    },
    baseDataOffsetPresent = result.flags[2] & 0x01,
    sampleDescriptionIndexPresent = result.flags[2] & 0x02,
    defaultSampleDurationPresent = result.flags[2] & 0x08,
    defaultSampleSizePresent = result.flags[2] & 0x10,
    defaultSampleFlagsPresent = result.flags[2] & 0x20,
    durationIsEmpty = result.flags[0] & 0x010000,
    defaultBaseIsMoof =  result.flags[0] & 0x020000,
    i;

  i = 8;
  if (baseDataOffsetPresent) {
    i += 4; // truncate top 4 bytes
    // FIXME: should we read the full 64 bits?
    result.baseDataOffset = view.getUint32(12);
    i += 4;
  }
  if (sampleDescriptionIndexPresent) {
    result.sampleDescriptionIndex = view.getUint32(i);
    i += 4;
  }
  if (defaultSampleDurationPresent) {
    result.defaultSampleDuration = view.getUint32(i);
    i += 4;
  }
  if (defaultSampleSizePresent) {
    result.defaultSampleSize = view.getUint32(i);
    i += 4;
  }
  if (defaultSampleFlagsPresent) {
    result.defaultSampleFlags = view.getUint32(i);
  }
  if (durationIsEmpty) {
    result.durationIsEmpty = true;
  }
  if (!baseDataOffsetPresent && defaultBaseIsMoof) {
    result.baseDataOffsetIsMoof = true;
  }
  return result;
};

var parseTfhd = tfhd;

var parseSampleFlags = function(flags) {
  return {
    isLeading: (flags[0] & 0x0c) >>> 2,
    dependsOn: flags[0] & 0x03,
    isDependedOn: (flags[1] & 0xc0) >>> 6,
    hasRedundancy: (flags[1] & 0x30) >>> 4,
    paddingValue: (flags[1] & 0x0e) >>> 1,
    isNonSyncSample: flags[1] & 0x01,
    degradationPriority: (flags[2] << 8) | flags[3]
  };
};

var parseSampleFlags_1 = parseSampleFlags;

var trun = function(data) {
  var
  result = {
    version: data[0],
    flags: new Uint8Array(data.subarray(1, 4)),
    samples: []
  },
    view = new DataView(data.buffer, data.byteOffset, data.byteLength),
    // Flag interpretation
    dataOffsetPresent = result.flags[2] & 0x01, // compare with 2nd byte of 0x1
    firstSampleFlagsPresent = result.flags[2] & 0x04, // compare with 2nd byte of 0x4
    sampleDurationPresent = result.flags[1] & 0x01, // compare with 2nd byte of 0x100
    sampleSizePresent = result.flags[1] & 0x02, // compare with 2nd byte of 0x200
    sampleFlagsPresent = result.flags[1] & 0x04, // compare with 2nd byte of 0x400
    sampleCompositionTimeOffsetPresent = result.flags[1] & 0x08, // compare with 2nd byte of 0x800
    sampleCount = view.getUint32(4),
    offset = 8,
    sample;

  if (dataOffsetPresent) {
    // 32 bit signed integer
    result.dataOffset = view.getInt32(offset);
    offset += 4;
  }

  // Overrides the flags for the first sample only. The order of
  // optional values will be: duration, size, compositionTimeOffset
  if (firstSampleFlagsPresent && sampleCount) {
    sample = {
      flags: parseSampleFlags_1(data.subarray(offset, offset + 4))
    };
    offset += 4;
    if (sampleDurationPresent) {
      sample.duration = view.getUint32(offset);
      offset += 4;
    }
    if (sampleSizePresent) {
      sample.size = view.getUint32(offset);
      offset += 4;
    }
    if (sampleCompositionTimeOffsetPresent) {
      if (result.version === 1) {
        sample.compositionTimeOffset = view.getInt32(offset);
      } else {
        sample.compositionTimeOffset = view.getUint32(offset);
      }
      offset += 4;
    }
    result.samples.push(sample);
    sampleCount--;
  }

  while (sampleCount--) {
    sample = {};
    if (sampleDurationPresent) {
      sample.duration = view.getUint32(offset);
      offset += 4;
    }
    if (sampleSizePresent) {
      sample.size = view.getUint32(offset);
      offset += 4;
    }
    if (sampleFlagsPresent) {
      sample.flags = parseSampleFlags_1(data.subarray(offset, offset + 4));
      offset += 4;
    }
    if (sampleCompositionTimeOffsetPresent) {
      if (result.version === 1) {
        sample.compositionTimeOffset = view.getInt32(offset);
      } else {
        sample.compositionTimeOffset = view.getUint32(offset);
      }
      offset += 4;
    }
    result.samples.push(sample);
  }
  return result;
};

var parseTrun = trun;

var toUnsigned$2 = bin.toUnsigned;

var tfdt = function(data) {
  var result = {
    version: data[0],
    flags: new Uint8Array(data.subarray(1, 4)),
    baseMediaDecodeTime: toUnsigned$2(data[4] << 24 | data[5] << 16 | data[6] << 8 | data[7])
  };
  if (result.version === 1) {
    result.baseMediaDecodeTime *= Math.pow(2, 32);
    result.baseMediaDecodeTime += toUnsigned$2(data[8] << 24 | data[9] << 16 | data[10] << 8 | data[11]);
  }
  return result;
};

var parseTfdt = tfdt;

var toUnsigned$3 = bin.toUnsigned;
var toHexString$1 = bin.toHexString;





var timescale, startTime, compositionStartTime, getVideoTrackIds, getTracks,
  getTimescaleFromMediaHeader;

/**
 * Parses an MP4 initialization segment and extracts the timescale
 * values for any declared tracks. Timescale values indicate the
 * number of clock ticks per second to assume for time-based values
 * elsewhere in the MP4.
 *
 * To determine the start time of an MP4, you need two pieces of
 * information: the timescale unit and the earliest base media decode
 * time. Multiple timescales can be specified within an MP4 but the
 * base media decode time is always expressed in the timescale from
 * the media header box for the track:
 * ```
 * moov > trak > mdia > mdhd.timescale
 * ```
 * @param init {Uint8Array} the bytes of the init segment
 * @return {object} a hash of track ids to timescale values or null if
 * the init segment is malformed.
 */
timescale = function(init) {
  var
    result = {},
    traks = findBox_1(init, ['moov', 'trak']);

  // mdhd timescale
  return traks.reduce(function(result, trak) {
    var tkhd, version, index, id, mdhd;

    tkhd = findBox_1(trak, ['tkhd'])[0];
    if (!tkhd) {
      return null;
    }
    version = tkhd[0];
    index = version === 0 ? 12 : 20;
    id = toUnsigned$3(tkhd[index]     << 24 |
                    tkhd[index + 1] << 16 |
                    tkhd[index + 2] <<  8 |
                    tkhd[index + 3]);

    mdhd = findBox_1(trak, ['mdia', 'mdhd'])[0];
    if (!mdhd) {
      return null;
    }
    version = mdhd[0];
    index = version === 0 ? 12 : 20;
    result[id] = toUnsigned$3(mdhd[index]     << 24 |
                            mdhd[index + 1] << 16 |
                            mdhd[index + 2] <<  8 |
                            mdhd[index + 3]);
    return result;
  }, result);
};

/**
 * Determine the base media decode start time, in seconds, for an MP4
 * fragment. If multiple fragments are specified, the earliest time is
 * returned.
 *
 * The base media decode time can be parsed from track fragment
 * metadata:
 * ```
 * moof > traf > tfdt.baseMediaDecodeTime
 * ```
 * It requires the timescale value from the mdhd to interpret.
 *
 * @param timescale {object} a hash of track ids to timescale values.
 * @return {number} the earliest base media decode start time for the
 * fragment, in seconds
 */
startTime = function(timescale, fragment) {
  var trafs, baseTimes, result;

  // we need info from two childrend of each track fragment box
  trafs = findBox_1(fragment, ['moof', 'traf']);

  // determine the start times for each track
  baseTimes = [].concat.apply([], trafs.map(function(traf) {
    return findBox_1(traf, ['tfhd']).map(function(tfhd) {
      var id, scale, baseTime;

      // get the track id from the tfhd
      id = toUnsigned$3(tfhd[4] << 24 |
                      tfhd[5] << 16 |
                      tfhd[6] <<  8 |
                      tfhd[7]);
      // assume a 90kHz clock if no timescale was specified
      scale = timescale[id] || 90e3;

      // get the base media decode time from the tfdt
      baseTime = findBox_1(traf, ['tfdt']).map(function(tfdt) {
        var version, result;

        version = tfdt[0];
        result = toUnsigned$3(tfdt[4] << 24 |
                            tfdt[5] << 16 |
                            tfdt[6] <<  8 |
                            tfdt[7]);
        if (version ===  1) {
          result *= Math.pow(2, 32);
          result += toUnsigned$3(tfdt[8]  << 24 |
                               tfdt[9]  << 16 |
                               tfdt[10] <<  8 |
                               tfdt[11]);
        }
        return result;
      })[0];
      baseTime = baseTime || Infinity;

      // convert base time to seconds
      return baseTime / scale;
    });
  }));

  // return the minimum
  result = Math.min.apply(null, baseTimes);
  return isFinite(result) ? result : 0;
};

/**
 * Determine the composition start, in seconds, for an MP4
 * fragment.
 *
 * The composition start time of a fragment can be calculated using the base
 * media decode time, composition time offset, and timescale, as follows:
 *
 * compositionStartTime = (baseMediaDecodeTime + compositionTimeOffset) / timescale
 *
 * All of the aforementioned information is contained within a media fragment's
 * `traf` box, except for timescale info, which comes from the initialization
 * segment, so a track id (also contained within a `traf`) is also necessary to
 * associate it with a timescale
 *
 *
 * @param timescales {object} - a hash of track ids to timescale values.
 * @param fragment {Unit8Array} - the bytes of a media segment
 * @return {number} the composition start time for the fragment, in seconds
 **/
compositionStartTime = function(timescales, fragment) {
  var trafBoxes = findBox_1(fragment, ['moof', 'traf']);
  var baseMediaDecodeTime = 0;
  var compositionTimeOffset = 0;
  var trackId;

  if (trafBoxes && trafBoxes.length) {
    // The spec states that track run samples contained within a `traf` box are contiguous, but
    // it does not explicitly state whether the `traf` boxes themselves are contiguous.
    // We will assume that they are, so we only need the first to calculate start time.
    var tfhd = findBox_1(trafBoxes[0], ['tfhd'])[0];
    var trun = findBox_1(trafBoxes[0], ['trun'])[0];
    var tfdt = findBox_1(trafBoxes[0], ['tfdt'])[0];

    if (tfhd) {
      var parsedTfhd = parseTfhd(tfhd);

      trackId = parsedTfhd.trackId;
    }

    if (tfdt) {
      var parsedTfdt = parseTfdt(tfdt);

      baseMediaDecodeTime = parsedTfdt.baseMediaDecodeTime;
    }

    if (trun) {
      var parsedTrun = parseTrun(trun);

      if (parsedTrun.samples && parsedTrun.samples.length) {
        compositionTimeOffset = parsedTrun.samples[0].compositionTimeOffset || 0;
      }
    }
  }

  // Get timescale for this specific track. Assume a 90kHz clock if no timescale was
  // specified.
  var timescale = timescales[trackId] || 90e3;

  // return the composition start time, in seconds
  return (baseMediaDecodeTime + compositionTimeOffset) / timescale;
};

/**
  * Find the trackIds of the video tracks in this source.
  * Found by parsing the Handler Reference and Track Header Boxes:
  *   moov > trak > mdia > hdlr
  *   moov > trak > tkhd
  *
  * @param {Uint8Array} init - The bytes of the init segment for this source
  * @return {Number[]} A list of trackIds
  *
  * @see ISO-BMFF-12/2015, Section 8.4.3
 **/
getVideoTrackIds = function(init) {
  var traks = findBox_1(init, ['moov', 'trak']);
  var videoTrackIds = [];

  traks.forEach(function(trak) {
    var hdlrs = findBox_1(trak, ['mdia', 'hdlr']);
    var tkhds = findBox_1(trak, ['tkhd']);

    hdlrs.forEach(function(hdlr, index) {
      var handlerType = parseType_1(hdlr.subarray(8, 12));
      var tkhd = tkhds[index];
      var view;
      var version;
      var trackId;

      if (handlerType === 'vide') {
        view = new DataView(tkhd.buffer, tkhd.byteOffset, tkhd.byteLength);
        version = view.getUint8(0);
        trackId = (version === 0) ? view.getUint32(12) : view.getUint32(20);

        videoTrackIds.push(trackId);
      }
    });
  });

  return videoTrackIds;
};

getTimescaleFromMediaHeader = function(mdhd) {
  // mdhd is a FullBox, meaning it will have its own version as the first byte
  var version = mdhd[0];
  var index = version === 0 ? 12 : 20;

  return toUnsigned$3(
    mdhd[index]     << 24 |
    mdhd[index + 1] << 16 |
    mdhd[index + 2] <<  8 |
    mdhd[index + 3]
  );
};

/**
 * Get all the video, audio, and hint tracks from a non fragmented
 * mp4 segment
 */
getTracks = function(init) {
  var traks = findBox_1(init, ['moov', 'trak']);
  var tracks = [];

  traks.forEach(function(trak) {
    var track = {};
    var tkhd = findBox_1(trak, ['tkhd'])[0];
    var view, tkhdVersion;

    // id
    if (tkhd) {
      view = new DataView(tkhd.buffer, tkhd.byteOffset, tkhd.byteLength);
      tkhdVersion = view.getUint8(0);

      track.id = (tkhdVersion === 0) ? view.getUint32(12) : view.getUint32(20);
    }

    var hdlr = findBox_1(trak, ['mdia', 'hdlr'])[0];

    // type
    if (hdlr) {
      var type = parseType_1(hdlr.subarray(8, 12));

      if (type === 'vide') {
        track.type = 'video';
      } else if (type === 'soun') {
        track.type = 'audio';
      } else {
        track.type = type;
      }
    }


    // codec
    var stsd = findBox_1(trak, ['mdia', 'minf', 'stbl', 'stsd'])[0];

    if (stsd) {
      var sampleDescriptions = stsd.subarray(8);
      // gives the codec type string
      track.codec = parseType_1(sampleDescriptions.subarray(4, 8));

      var codecBox = findBox_1(sampleDescriptions, [track.codec])[0];
      var codecConfig, codecConfigType;

      if (codecBox) {
        // https://tools.ietf.org/html/rfc6381#section-3.3
        if ((/^[a-z]vc[1-9]$/i).test(track.codec)) {
          // we don't need anything but the "config" parameter of the
          // avc1 codecBox
          codecConfig = codecBox.subarray(78);
          codecConfigType = parseType_1(codecConfig.subarray(4, 8));

          if (codecConfigType === 'avcC' && codecConfig.length > 11) {
            track.codec += '.';

            // left padded with zeroes for single digit hex
            // profile idc
            track.codec +=  toHexString$1(codecConfig[9]);
            // the byte containing the constraint_set flags
            track.codec += toHexString$1(codecConfig[10]);
            // level idc
            track.codec += toHexString$1(codecConfig[11]);
          } else {
            // TODO: show a warning that we couldn't parse the codec
            // and are using the default
            track.codec = 'avc1.4d400d';
          }
        } else if ((/^mp4[a,v]$/i).test(track.codec)) {
          // we do not need anything but the streamDescriptor of the mp4a codecBox
          codecConfig = codecBox.subarray(28);
          codecConfigType = parseType_1(codecConfig.subarray(4, 8));

          if (codecConfigType === 'esds' && codecConfig.length > 20 && codecConfig[19] !== 0) {
            track.codec += '.' + toHexString$1(codecConfig[19]);
            // this value is only a single digit
            track.codec += '.' + toHexString$1((codecConfig[20] >>> 2) & 0x3f).replace(/^0/, '');
          } else {
            // TODO: show a warning that we couldn't parse the codec
            // and are using the default
            track.codec = 'mp4a.40.2';
          }
        } else ;
      }
    }

    var mdhd = findBox_1(trak, ['mdia', 'mdhd'])[0];

    if (mdhd) {
      track.timescale = getTimescaleFromMediaHeader(mdhd);
    }

    tracks.push(track);
  });

  return tracks;
};

var probe$2 = {
  // export mp4 inspector's findBox and parseType for backwards compatibility
  findBox: findBox_1,
  parseType: parseType_1,
  timescale: timescale,
  startTime: startTime,
  compositionStartTime: compositionStartTime,
  videoTrackIds: getVideoTrackIds,
  tracks: getTracks,
  getTimescaleFromMediaHeader: getTimescaleFromMediaHeader
};

var codecs = createCommonjsModule(function (module, exports) {

Object.defineProperty(exports, '__esModule', { value: true });

function _interopDefault (ex) { return (ex && (typeof ex === 'object') && 'default' in ex) ? ex['default'] : ex; }

var window = _interopDefault(window_1);

var regexs = {
  // to determine mime types
  mp4: /^(av0?1|avc0?[1234]|vp0?9|flac|opus|mp3|mp4a|mp4v)/,
  webm: /^(vp0?[89]|av0?1|opus|vorbis)/,
  ogg: /^(vp0?[89]|theora|flac|opus|vorbis)/,
  // to determine if a codec is audio or video
  video: /^(av0?1|avc0?[1234]|vp0?[89]|hvc1|hev1|theora|mp4v)/,
  audio: /^(mp4a|flac|vorbis|opus|ac-[34]|ec-3|alac|mp3)/,
  // mux.js support regex
  muxerVideo: /^(avc0?1)/,
  muxerAudio: /^(mp4a)/
};
/**
 * Replace the old apple-style `avc1.<dd>.<dd>` codec string with the standard
 * `avc1.<hhhhhh>`
 *
 * @param {string} codec
 *        Codec string to translate
 * @return {string}
 *         The translated codec string
 */

var translateLegacyCodec = function translateLegacyCodec(codec) {
  if (!codec) {
    return codec;
  }

  return codec.replace(/avc1\.(\d+)\.(\d+)/i, function (orig, profile, avcLevel) {
    var profileHex = ('00' + Number(profile).toString(16)).slice(-2);
    var avcLevelHex = ('00' + Number(avcLevel).toString(16)).slice(-2);
    return 'avc1.' + profileHex + '00' + avcLevelHex;
  });
};
/**
 * Replace the old apple-style `avc1.<dd>.<dd>` codec strings with the standard
 * `avc1.<hhhhhh>`
 *
 * @param {string[]} codecs
 *        An array of codec strings to translate
 * @return {string[]}
 *         The translated array of codec strings
 */

var translateLegacyCodecs = function translateLegacyCodecs(codecs) {
  return codecs.map(translateLegacyCodec);
};
/**
 * Replace codecs in the codec string with the old apple-style `avc1.<dd>.<dd>` to the
 * standard `avc1.<hhhhhh>`.
 *
 * @param {string} codecString
 *        The codec string
 * @return {string}
 *         The codec string with old apple-style codecs replaced
 *
 * @private
 */

var mapLegacyAvcCodecs = function mapLegacyAvcCodecs(codecString) {
  return codecString.replace(/avc1\.(\d+)\.(\d+)/i, function (match) {
    return translateLegacyCodecs([match])[0];
  });
};
/**
 * @typedef {Object} ParsedCodecInfo
 * @property {number} codecCount
 *           Number of codecs parsed
 * @property {string} [videoCodec]
 *           Parsed video codec (if found)
 * @property {string} [videoObjectTypeIndicator]
 *           Video object type indicator (if found)
 * @property {string|null} audioProfile
 *           Audio profile
 */

/**
 * Parses a codec string to retrieve the number of codecs specified, the video codec and
 * object type indicator, and the audio profile.
 *
 * @param {string} [codecString]
 *        The codec string to parse
 * @return {ParsedCodecInfo}
 *         Parsed codec info
 */

var parseCodecs = function parseCodecs(codecString) {
  if (codecString === void 0) {
    codecString = '';
  }

  var codecs = codecString.split(',');
  var result = {};
  codecs.forEach(function (codec) {
    codec = codec.trim();
    ['video', 'audio'].forEach(function (name) {
      var match = regexs[name].exec(codec.toLowerCase());

      if (!match || match.length <= 1) {
        return;
      } // maintain codec case


      var type = codec.substring(0, match[1].length);
      var details = codec.replace(type, '');
      result[name] = {
        type: type,
        details: details
      };
    });
  });
  return result;
};
/**
 * Returns a ParsedCodecInfo object for the default alternate audio playlist if there is
 * a default alternate audio playlist for the provided audio group.
 *
 * @param {Object} master
 *        The master playlist
 * @param {string} audioGroupId
 *        ID of the audio group for which to find the default codec info
 * @return {ParsedCodecInfo}
 *         Parsed codec info
 */

var codecsFromDefault = function codecsFromDefault(master, audioGroupId) {
  if (!master.mediaGroups.AUDIO || !audioGroupId) {
    return null;
  }

  var audioGroup = master.mediaGroups.AUDIO[audioGroupId];

  if (!audioGroup) {
    return null;
  }

  for (var name in audioGroup) {
    var audioType = audioGroup[name];

    if (audioType.default && audioType.playlists) {
      // codec should be the same for all playlists within the audio type
      return parseCodecs(audioType.playlists[0].attributes.CODECS);
    }
  }

  return null;
};
var isVideoCodec = function isVideoCodec(codec) {
  if (codec === void 0) {
    codec = '';
  }

  return regexs.video.test(codec.trim().toLowerCase());
};
var isAudioCodec = function isAudioCodec(codec) {
  if (codec === void 0) {
    codec = '';
  }

  return regexs.audio.test(codec.trim().toLowerCase());
};
var getMimeForCodec = function getMimeForCodec(codecString) {
  if (!codecString || typeof codecString !== 'string') {
    return;
  }

  var codecs = codecString.toLowerCase().split(',').map(function (c) {
    return translateLegacyCodec(c.trim());
  }); // default to video type

  var type = 'video'; // only change to audio type if the only codec we have is
  // audio

  if (codecs.length === 1 && isAudioCodec(codecs[0])) {
    type = 'audio';
  } // default the container to mp4


  var container = 'mp4'; // every codec must be able to go into the container
  // for that container to be the correct one

  if (codecs.every(function (c) {
    return regexs.mp4.test(c);
  })) {
    container = 'mp4';
  } else if (codecs.every(function (c) {
    return regexs.webm.test(c);
  })) {
    container = 'webm';
  } else if (codecs.every(function (c) {
    return regexs.ogg.test(c);
  })) {
    container = 'ogg';
  }

  return type + "/" + container + ";codecs=\"" + codecString + "\"";
};
var browserSupportsCodec = function browserSupportsCodec(codecString) {
  if (codecString === void 0) {
    codecString = '';
  }

  return window.MediaSource && window.MediaSource.isTypeSupported && window.MediaSource.isTypeSupported(getMimeForCodec(codecString)) || false;
};
var muxerSupportsCodec = function muxerSupportsCodec(codecString) {
  if (codecString === void 0) {
    codecString = '';
  }

  return codecString.toLowerCase().split(',').every(function (codec) {
    codec = codec.trim();
    return regexs.muxerVideo.test(codec) || regexs.muxerAudio.test(codec);
  });
};
var DEFAULT_AUDIO_CODEC = 'mp4a.40.2';
var DEFAULT_VIDEO_CODEC = 'avc1.4d400d';

exports.DEFAULT_AUDIO_CODEC = DEFAULT_AUDIO_CODEC;
exports.DEFAULT_VIDEO_CODEC = DEFAULT_VIDEO_CODEC;
exports.browserSupportsCodec = browserSupportsCodec;
exports.codecsFromDefault = codecsFromDefault;
exports.getMimeForCodec = getMimeForCodec;
exports.isAudioCodec = isAudioCodec;
exports.isVideoCodec = isVideoCodec;
exports.mapLegacyAvcCodecs = mapLegacyAvcCodecs;
exports.muxerSupportsCodec = muxerSupportsCodec;
exports.parseCodecs = parseCodecs;
exports.translateLegacyCodec = translateLegacyCodec;
exports.translateLegacyCodecs = translateLegacyCodecs;
});

const version = "5.6.6";

const version$1 = "0.12.0";

const version$2 = "4.4.3";

const version$3 = "3.0.2";

/*! @name @videojs/http-streaming @version 2.2.0 @license Apache-2.0 */

/**
 * @file resolve-url.js - Handling how URLs are resolved and manipulated
 */
var resolveUrl$2 = resolveUrl_1;
/**
 * Checks whether xhr request was redirected and returns correct url depending
 * on `handleManifestRedirects` option
 *
 * @api private
 *
 * @param  {string} url - an url being requested
 * @param  {XMLHttpRequest} req - xhr request result
 *
 * @return {string}
 */

var resolveManifestRedirect = function resolveManifestRedirect(handleManifestRedirect, url, req) {
  // To understand how the responseURL below is set and generated:
  // - https://fetch.spec.whatwg.org/#concept-response-url
  // - https://fetch.spec.whatwg.org/#atomic-http-redirect-handling
  if (handleManifestRedirect && req && req.responseURL && url !== req.responseURL) {
    return req.responseURL;
  }

  return url;
};

var log = videojs$1.log;
var createPlaylistID = function createPlaylistID(index, uri) {
  return index + "-" + uri;
};
/**
 * Parses a given m3u8 playlist
 *
 * @param {string} manifestString
 *        The downloaded manifest string
 * @param {Object[]} [customTagParsers]
 *        An array of custom tag parsers for the m3u8-parser instance
 * @param {Object[]} [customTagMappers]
 *         An array of custom tag mappers for the m3u8-parser instance
 * @return {Object}
 *         The manifest object
 */

var parseManifest = function parseManifest(_ref) {
  var manifestString = _ref.manifestString,
      _ref$customTagParsers = _ref.customTagParsers,
      customTagParsers = _ref$customTagParsers === void 0 ? [] : _ref$customTagParsers,
      _ref$customTagMappers = _ref.customTagMappers,
      customTagMappers = _ref$customTagMappers === void 0 ? [] : _ref$customTagMappers;
  var parser = new Parser();
  customTagParsers.forEach(function (customParser) {
    return parser.addParser(customParser);
  });
  customTagMappers.forEach(function (mapper) {
    return parser.addTagMapper(mapper);
  });
  parser.push(manifestString);
  parser.end();
  return parser.manifest;
};
/**
 * Loops through all supported media groups in master and calls the provided
 * callback for each group
 *
 * @param {Object} master
 *        The parsed master manifest object
 * @param {Function} callback
 *        Callback to call for each media group
 */

var forEachMediaGroup = function forEachMediaGroup(master, callback) {
  ['AUDIO', 'SUBTITLES'].forEach(function (mediaType) {
    for (var groupKey in master.mediaGroups[mediaType]) {
      for (var labelKey in master.mediaGroups[mediaType][groupKey]) {
        var mediaProperties = master.mediaGroups[mediaType][groupKey][labelKey];
        callback(mediaProperties, mediaType, groupKey, labelKey);
      }
    }
  });
};
/**
 * Adds properties and attributes to the playlist to keep consistent functionality for
 * playlists throughout VHS.
 *
 * @param {Object} config
 *        Arguments object
 * @param {Object} config.playlist
 *        The media playlist
 * @param {string} [config.uri]
 *        The uri to the media playlist (if media playlist is not from within a master
 *        playlist)
 * @param {string} id
 *        ID to use for the playlist
 */

var setupMediaPlaylist = function setupMediaPlaylist(_ref2) {
  var playlist = _ref2.playlist,
      uri = _ref2.uri,
      id = _ref2.id;
  playlist.id = id;

  if (uri) {
    // For media playlists, m3u8-parser does not have access to a URI, as HLS media
    // playlists do not contain their own source URI, but one is needed for consistency in
    // VHS.
    playlist.uri = uri;
  } // For HLS master playlists, even though certain attributes MUST be defined, the
  // stream may still be played without them.
  // For HLS media playlists, m3u8-parser does not attach an attributes object to the
  // manifest.
  //
  // To avoid undefined reference errors through the project, and make the code easier
  // to write/read, add an empty attributes object for these cases.


  playlist.attributes = playlist.attributes || {};
};
/**
 * Adds ID, resolvedUri, and attributes properties to each playlist of the master, where
 * necessary. In addition, creates playlist IDs for each playlist and adds playlist ID to
 * playlist references to the playlists array.
 *
 * @param {Object} master
 *        The master playlist
 */

var setupMediaPlaylists = function setupMediaPlaylists(master) {
  var i = master.playlists.length;

  while (i--) {
    var playlist = master.playlists[i];
    setupMediaPlaylist({
      playlist: playlist,
      id: createPlaylistID(i, playlist.uri)
    });
    playlist.resolvedUri = resolveUrl$2(master.uri, playlist.uri);
    master.playlists[playlist.id] = playlist; // URI reference added for backwards compatibility

    master.playlists[playlist.uri] = playlist; // Although the spec states an #EXT-X-STREAM-INF tag MUST have a BANDWIDTH attribute,
    // the stream can be played without it. Although an attributes property may have been
    // added to the playlist to prevent undefined references, issue a warning to fix the
    // manifest.

    if (!playlist.attributes.BANDWIDTH) {
      log.warn('Invalid playlist STREAM-INF detected. Missing BANDWIDTH attribute.');
    }
  }
};
/**
 * Adds resolvedUri properties to each media group.
 *
 * @param {Object} master
 *        The master playlist
 */

var resolveMediaGroupUris = function resolveMediaGroupUris(master) {
  forEachMediaGroup(master, function (properties) {
    if (properties.uri) {
      properties.resolvedUri = resolveUrl$2(master.uri, properties.uri);
    }
  });
};
/**
 * Creates a master playlist wrapper to insert a sole media playlist into.
 *
 * @param {Object} media
 *        Media playlist
 * @param {string} uri
 *        The media URI
 *
 * @return {Object}
 *         Master playlist
 */

var masterForMedia = function masterForMedia(media, uri) {
  var id = createPlaylistID(0, uri);
  var master = {
    mediaGroups: {
      'AUDIO': {},
      'VIDEO': {},
      'CLOSED-CAPTIONS': {},
      'SUBTITLES': {}
    },
    uri: window_1.location.href,
    resolvedUri: window_1.location.href,
    playlists: [{
      uri: uri,
      id: id,
      resolvedUri: uri,
      // m3u8-parser does not attach an attributes property to media playlists so make
      // sure that the property is attached to avoid undefined reference errors
      attributes: {}
    }]
  }; // set up ID reference

  master.playlists[id] = master.playlists[0]; // URI reference added for backwards compatibility

  master.playlists[uri] = master.playlists[0];
  return master;
};
/**
 * Does an in-place update of the master manifest to add updated playlist URI references
 * as well as other properties needed by VHS that aren't included by the parser.
 *
 * @param {Object} master
 *        Master manifest object
 * @param {string} uri
 *        The source URI
 */

var addPropertiesToMaster = function addPropertiesToMaster(master, uri) {
  master.uri = uri;

  for (var i = 0; i < master.playlists.length; i++) {
    if (!master.playlists[i].uri) {
      // Set up phony URIs for the playlists since playlists are referenced by their URIs
      // throughout VHS, but some formats (e.g., DASH) don't have external URIs
      // TODO: consider adding dummy URIs in mpd-parser
      var phonyUri = "placeholder-uri-" + i;
      master.playlists[i].uri = phonyUri;
    }
  }

  forEachMediaGroup(master, function (properties, mediaType, groupKey, labelKey) {
    if (!properties.playlists || !properties.playlists.length || properties.playlists[0].uri) {
      return;
    } // Set up phony URIs for the media group playlists since playlists are referenced by
    // their URIs throughout VHS, but some formats (e.g., DASH) don't have external URIs


    var phonyUri = "placeholder-uri-" + mediaType + "-" + groupKey + "-" + labelKey;
    var id = createPlaylistID(0, phonyUri);
    properties.playlists[0].uri = phonyUri;
    properties.playlists[0].id = id; // setup ID and URI references (URI for backwards compatibility)

    master.playlists[id] = properties.playlists[0];
    master.playlists[phonyUri] = properties.playlists[0];
  });
  setupMediaPlaylists(master);
  resolveMediaGroupUris(master);
};

var mergeOptions = videojs$1.mergeOptions,
    EventTarget = videojs$1.EventTarget;
/**
  * Returns a new array of segments that is the result of merging
  * properties from an older list of segments onto an updated
  * list. No properties on the updated playlist will be overridden.
  *
  * @param {Array} original the outdated list of segments
  * @param {Array} update the updated list of segments
  * @param {number=} offset the index of the first update
  * segment in the original segment list. For non-live playlists,
  * this should always be zero and does not need to be
  * specified. For live playlists, it should be the difference
  * between the media sequence numbers in the original and updated
  * playlists.
  * @return a list of merged segment objects
  */

var updateSegments = function updateSegments(original, update, offset) {
  var result = update.slice();
  offset = offset || 0;
  var length = Math.min(original.length, update.length + offset);

  for (var i = offset; i < length; i++) {
    result[i - offset] = mergeOptions(original[i], result[i - offset]);
  }

  return result;
};
var resolveSegmentUris = function resolveSegmentUris(segment, baseUri) {
  if (!segment.resolvedUri) {
    segment.resolvedUri = resolveUrl$2(baseUri, segment.uri);
  }

  if (segment.key && !segment.key.resolvedUri) {
    segment.key.resolvedUri = resolveUrl$2(baseUri, segment.key.uri);
  }

  if (segment.map && !segment.map.resolvedUri) {
    segment.map.resolvedUri = resolveUrl$2(baseUri, segment.map.uri);
  }
};
/**
  * Returns a new master playlist that is the result of merging an
  * updated media playlist into the original version. If the
  * updated media playlist does not match any of the playlist
  * entries in the original master playlist, null is returned.
  *
  * @param {Object} master a parsed master M3U8 object
  * @param {Object} media a parsed media M3U8 object
  * @return {Object} a new object that represents the original
  * master playlist with the updated media playlist merged in, or
  * null if the merge produced no change.
  */

var updateMaster = function updateMaster(master, media) {
  var result = mergeOptions(master, {});
  var playlist = result.playlists[media.id];

  if (!playlist) {
    return null;
  } // consider the playlist unchanged if the number of segments is equal, the media
  // sequence number is unchanged, and this playlist hasn't become the end of the playlist


  if (playlist.segments && media.segments && playlist.segments.length === media.segments.length && playlist.endList === media.endList && playlist.mediaSequence === media.mediaSequence) {
    return null;
  }

  var mergedPlaylist = mergeOptions(playlist, media); // if the update could overlap existing segment information, merge the two segment lists

  if (playlist.segments) {
    mergedPlaylist.segments = updateSegments(playlist.segments, media.segments, media.mediaSequence - playlist.mediaSequence);
  } // resolve any segment URIs to prevent us from having to do it later


  mergedPlaylist.segments.forEach(function (segment) {
    resolveSegmentUris(segment, mergedPlaylist.resolvedUri);
  }); // TODO Right now in the playlists array there are two references to each playlist, one
  // that is referenced by index, and one by URI. The index reference may no longer be
  // necessary.

  for (var i = 0; i < result.playlists.length; i++) {
    if (result.playlists[i].id === media.id) {
      result.playlists[i] = mergedPlaylist;
    }
  }

  result.playlists[media.id] = mergedPlaylist; // URI reference added for backwards compatibility

  result.playlists[media.uri] = mergedPlaylist;
  return result;
};
/**
 * Calculates the time to wait before refreshing a live playlist
 *
 * @param {Object} media
 *        The current media
 * @param {boolean} update
 *        True if there were any updates from the last refresh, false otherwise
 * @return {number}
 *         The time in ms to wait before refreshing the live playlist
 */

var refreshDelay = function refreshDelay(media, update) {
  var lastSegment = media.segments[media.segments.length - 1];
  var delay;

  if (update && lastSegment && lastSegment.duration) {
    delay = lastSegment.duration * 1000;
  } else {
    // if the playlist is unchanged since the last reload or last segment duration
    // cannot be determined, try again after half the target duration
    delay = (media.targetDuration || 10) * 500;
  }

  return delay;
};
/**
 * Load a playlist from a remote location
 *
 * @class PlaylistLoader
 * @extends Stream
 * @param {string|Object} src url or object of manifest
 * @param {boolean} withCredentials the withCredentials xhr option
 * @class
 */

var PlaylistLoader = /*#__PURE__*/function (_EventTarget) {
  inheritsLoose(PlaylistLoader, _EventTarget);

  function PlaylistLoader(src, vhs, options) {
    var _this;

    if (options === void 0) {
      options = {};
    }

    _this = _EventTarget.call(this) || this;

    if (!src) {
      throw new Error('A non-empty playlist URL or object is required');
    }

    var _options = options,
        _options$withCredenti = _options.withCredentials,
        withCredentials = _options$withCredenti === void 0 ? false : _options$withCredenti,
        _options$handleManife = _options.handleManifestRedirects,
        handleManifestRedirects = _options$handleManife === void 0 ? false : _options$handleManife;
    _this.src = src;
    _this.vhs_ = vhs;
    _this.withCredentials = withCredentials;
    _this.handleManifestRedirects = handleManifestRedirects;
    var vhsOptions = vhs.options_;
    _this.customTagParsers = vhsOptions && vhsOptions.customTagParsers || [];
    _this.customTagMappers = vhsOptions && vhsOptions.customTagMappers || []; // initialize the loader state

    _this.state = 'HAVE_NOTHING'; // live playlist staleness timeout

    _this.on('mediaupdatetimeout', function () {
      if (_this.state !== 'HAVE_METADATA') {
        // only refresh the media playlist if no other activity is going on
        return;
      }

      _this.state = 'HAVE_CURRENT_METADATA';
      _this.request = _this.vhs_.xhr({
        uri: resolveUrl$2(_this.master.uri, _this.media().uri),
        withCredentials: _this.withCredentials
      }, function (error, req) {
        // disposed
        if (!_this.request) {
          return;
        }

        if (error) {
          return _this.playlistRequestError(_this.request, _this.media(), 'HAVE_METADATA');
        }

        _this.haveMetadata({
          playlistString: _this.request.responseText,
          url: _this.media().uri,
          id: _this.media().id
        });
      });
    });

    return _this;
  }

  var _proto = PlaylistLoader.prototype;

  _proto.playlistRequestError = function playlistRequestError(xhr, playlist, startingState) {
    var uri = playlist.uri,
        id = playlist.id; // any in-flight request is now finished

    this.request = null;

    if (startingState) {
      this.state = startingState;
    }

    this.error = {
      playlist: this.master.playlists[id],
      status: xhr.status,
      message: "HLS playlist request error at URL: " + uri + ".",
      responseText: xhr.responseText,
      code: xhr.status >= 500 ? 4 : 2
    };
    this.trigger('error');
  }
  /**
   * Update the playlist loader's state in response to a new or updated playlist.
   *
   * @param {string} [playlistString]
   *        Playlist string (if playlistObject is not provided)
   * @param {Object} [playlistObject]
   *        Playlist object (if playlistString is not provided)
   * @param {string} url
   *        URL of playlist
   * @param {string} id
   *        ID to use for playlist
   */
  ;

  _proto.haveMetadata = function haveMetadata(_ref) {
    var _this2 = this;

    var playlistString = _ref.playlistString,
        playlistObject = _ref.playlistObject,
        url = _ref.url,
        id = _ref.id;
    // any in-flight request is now finished
    this.request = null;
    this.state = 'HAVE_METADATA';
    var playlist = playlistObject || parseManifest({
      manifestString: playlistString,
      customTagParsers: this.customTagParsers,
      customTagMappers: this.customTagMappers
    });
    setupMediaPlaylist({
      playlist: playlist,
      uri: url,
      id: id
    }); // merge this playlist into the master

    var update = updateMaster(this.master, playlist);
    this.targetDuration = playlist.targetDuration;

    if (update) {
      this.master = update;
      this.media_ = this.master.playlists[id];
    } else {
      this.trigger('playlistunchanged');
    } // refresh live playlists after a target duration passes


    if (!this.media().endList) {
      window_1.clearTimeout(this.mediaUpdateTimeout);
      this.mediaUpdateTimeout = window_1.setTimeout(function () {
        _this2.trigger('mediaupdatetimeout');
      }, refreshDelay(this.media(), !!update));
    }

    this.trigger('loadedplaylist');
  }
  /**
    * Abort any outstanding work and clean up.
    */
  ;

  _proto.dispose = function dispose() {
    this.trigger('dispose');
    this.stopRequest();
    window_1.clearTimeout(this.mediaUpdateTimeout);
    window_1.clearTimeout(this.finalRenditionTimeout);
    this.off();
  };

  _proto.stopRequest = function stopRequest() {
    if (this.request) {
      var oldRequest = this.request;
      this.request = null;
      oldRequest.onreadystatechange = null;
      oldRequest.abort();
    }
  }
  /**
    * When called without any arguments, returns the currently
    * active media playlist. When called with a single argument,
    * triggers the playlist loader to asynchronously switch to the
    * specified media playlist. Calling this method while the
    * loader is in the HAVE_NOTHING causes an error to be emitted
    * but otherwise has no effect.
    *
    * @param {Object=} playlist the parsed media playlist
    * object to switch to
    * @param {boolean=} is this the last available playlist
    *
    * @return {Playlist} the current loaded media
    */
  ;

  _proto.media = function media(playlist, isFinalRendition) {
    var _this3 = this;

    // getter
    if (!playlist) {
      return this.media_;
    } // setter


    if (this.state === 'HAVE_NOTHING') {
      throw new Error('Cannot switch media playlist from ' + this.state);
    } // find the playlist object if the target playlist has been
    // specified by URI


    if (typeof playlist === 'string') {
      if (!this.master.playlists[playlist]) {
        throw new Error('Unknown playlist URI: ' + playlist);
      }

      playlist = this.master.playlists[playlist];
    }

    window_1.clearTimeout(this.finalRenditionTimeout);

    if (isFinalRendition) {
      var delay = playlist.targetDuration / 2 * 1000 || 5 * 1000;
      this.finalRenditionTimeout = window_1.setTimeout(this.media.bind(this, playlist, false), delay);
      return;
    }

    var startingState = this.state;
    var mediaChange = !this.media_ || playlist.id !== this.media_.id; // switch to fully loaded playlists immediately

    if (this.master.playlists[playlist.id].endList || // handle the case of a playlist object (e.g., if using vhs-json with a resolved
    // media playlist or, for the case of demuxed audio, a resolved audio media group)
    playlist.endList && playlist.segments.length) {
      // abort outstanding playlist requests
      if (this.request) {
        this.request.onreadystatechange = null;
        this.request.abort();
        this.request = null;
      }

      this.state = 'HAVE_METADATA';
      this.media_ = playlist; // trigger media change if the active media has been updated

      if (mediaChange) {
        this.trigger('mediachanging');

        if (startingState === 'HAVE_MASTER') {
          // The initial playlist was a master manifest, and the first media selected was
          // also provided (in the form of a resolved playlist object) as part of the
          // source object (rather than just a URL). Therefore, since the media playlist
          // doesn't need to be requested, loadedmetadata won't trigger as part of the
          // normal flow, and needs an explicit trigger here.
          this.trigger('loadedmetadata');
        } else {
          this.trigger('mediachange');
        }
      }

      return;
    } // switching to the active playlist is a no-op


    if (!mediaChange) {
      return;
    }

    this.state = 'SWITCHING_MEDIA'; // there is already an outstanding playlist request

    if (this.request) {
      if (playlist.resolvedUri === this.request.url) {
        // requesting to switch to the same playlist multiple times
        // has no effect after the first
        return;
      }

      this.request.onreadystatechange = null;
      this.request.abort();
      this.request = null;
    } // request the new playlist


    if (this.media_) {
      this.trigger('mediachanging');
    }

    this.request = this.vhs_.xhr({
      uri: playlist.resolvedUri,
      withCredentials: this.withCredentials
    }, function (error, req) {
      // disposed
      if (!_this3.request) {
        return;
      }

      playlist.resolvedUri = resolveManifestRedirect(_this3.handleManifestRedirects, playlist.resolvedUri, req);

      if (error) {
        return _this3.playlistRequestError(_this3.request, playlist, startingState);
      }

      _this3.haveMetadata({
        playlistString: req.responseText,
        url: playlist.uri,
        id: playlist.id
      }); // fire loadedmetadata the first time a media playlist is loaded


      if (startingState === 'HAVE_MASTER') {
        _this3.trigger('loadedmetadata');
      } else {
        _this3.trigger('mediachange');
      }
    });
  }
  /**
   * pause loading of the playlist
   */
  ;

  _proto.pause = function pause() {
    this.stopRequest();
    window_1.clearTimeout(this.mediaUpdateTimeout);

    if (this.state === 'HAVE_NOTHING') {
      // If we pause the loader before any data has been retrieved, its as if we never
      // started, so reset to an unstarted state.
      this.started = false;
    } // Need to restore state now that no activity is happening


    if (this.state === 'SWITCHING_MEDIA') {
      // if the loader was in the process of switching media, it should either return to
      // HAVE_MASTER or HAVE_METADATA depending on if the loader has loaded a media
      // playlist yet. This is determined by the existence of loader.media_
      if (this.media_) {
        this.state = 'HAVE_METADATA';
      } else {
        this.state = 'HAVE_MASTER';
      }
    } else if (this.state === 'HAVE_CURRENT_METADATA') {
      this.state = 'HAVE_METADATA';
    }
  }
  /**
   * start loading of the playlist
   */
  ;

  _proto.load = function load(isFinalRendition) {
    var _this4 = this;

    window_1.clearTimeout(this.mediaUpdateTimeout);
    var media = this.media();

    if (isFinalRendition) {
      var delay = media ? media.targetDuration / 2 * 1000 : 5 * 1000;
      this.mediaUpdateTimeout = window_1.setTimeout(function () {
        return _this4.load();
      }, delay);
      return;
    }

    if (!this.started) {
      this.start();
      return;
    }

    if (media && !media.endList) {
      this.trigger('mediaupdatetimeout');
    } else {
      this.trigger('loadedplaylist');
    }
  }
  /**
   * start loading of the playlist
   */
  ;

  _proto.start = function start() {
    var _this5 = this;

    this.started = true;

    if (typeof this.src === 'object') {
      // in the case of an entirely constructed manifest object (meaning there's no actual
      // manifest on a server), default the uri to the page's href
      if (!this.src.uri) {
        this.src.uri = window_1.location.href;
      } // resolvedUri is added on internally after the initial request. Since there's no
      // request for pre-resolved manifests, add on resolvedUri here.


      this.src.resolvedUri = this.src.uri; // Since a manifest object was passed in as the source (instead of a URL), the first
      // request can be skipped (since the top level of the manifest, at a minimum, is
      // already available as a parsed manifest object). However, if the manifest object
      // represents a master playlist, some media playlists may need to be resolved before
      // the starting segment list is available. Therefore, go directly to setup of the
      // initial playlist, and let the normal flow continue from there.
      //
      // Note that the call to setup is asynchronous, as other sections of VHS may assume
      // that the first request is asynchronous.

      setTimeout(function () {
        _this5.setupInitialPlaylist(_this5.src);
      }, 0);
      return;
    } // request the specified URL


    this.request = this.vhs_.xhr({
      uri: this.src,
      withCredentials: this.withCredentials
    }, function (error, req) {
      // disposed
      if (!_this5.request) {
        return;
      } // clear the loader's request reference


      _this5.request = null;

      if (error) {
        _this5.error = {
          status: req.status,
          message: "HLS playlist request error at URL: " + _this5.src + ".",
          responseText: req.responseText,
          // MEDIA_ERR_NETWORK
          code: 2
        };

        if (_this5.state === 'HAVE_NOTHING') {
          _this5.started = false;
        }

        return _this5.trigger('error');
      }

      _this5.src = resolveManifestRedirect(_this5.handleManifestRedirects, _this5.src, req);
      var manifest = parseManifest({
        manifestString: req.responseText,
        customTagParsers: _this5.customTagParsers,
        customTagMappers: _this5.customTagMappers
      });

      _this5.setupInitialPlaylist(manifest);
    });
  };

  _proto.srcUri = function srcUri() {
    return typeof this.src === 'string' ? this.src : this.src.uri;
  }
  /**
   * Given a manifest object that's either a master or media playlist, trigger the proper
   * events and set the state of the playlist loader.
   *
   * If the manifest object represents a master playlist, `loadedplaylist` will be
   * triggered to allow listeners to select a playlist. If none is selected, the loader
   * will default to the first one in the playlists array.
   *
   * If the manifest object represents a media playlist, `loadedplaylist` will be
   * triggered followed by `loadedmetadata`, as the only available playlist is loaded.
   *
   * In the case of a media playlist, a master playlist object wrapper with one playlist
   * will be created so that all logic can handle playlists in the same fashion (as an
   * assumed manifest object schema).
   *
   * @param {Object} manifest
   *        The parsed manifest object
   */
  ;

  _proto.setupInitialPlaylist = function setupInitialPlaylist(manifest) {
    this.state = 'HAVE_MASTER';

    if (manifest.playlists) {
      this.master = manifest;
      addPropertiesToMaster(this.master, this.srcUri()); // If the initial master playlist has playlists wtih segments already resolved,
      // then resolve URIs in advance, as they are usually done after a playlist request,
      // which may not happen if the playlist is resolved.

      manifest.playlists.forEach(function (playlist) {
        if (playlist.segments) {
          playlist.segments.forEach(function (segment) {
            resolveSegmentUris(segment, playlist.resolvedUri);
          });
        }
      });
      this.trigger('loadedplaylist');

      if (!this.request) {
        // no media playlist was specifically selected so start
        // from the first listed one
        this.media(this.master.playlists[0]);
      }

      return;
    } // In order to support media playlists passed in as vhs-json, the case where the uri
    // is not provided as part of the manifest should be considered, and an appropriate
    // default used.


    var uri = this.srcUri() || window_1.location.href;
    this.master = masterForMedia(manifest, uri);
    this.haveMetadata({
      playlistObject: manifest,
      url: uri,
      id: this.master.playlists[0].id
    });
    this.trigger('loadedmetadata');
  };

  return PlaylistLoader;
}(EventTarget);

/**
 * ranges
 *
 * Utilities for working with TimeRanges.
 *
 */

var TIME_FUDGE_FACTOR = 1 / 30; // Comparisons between time values such as current time and the end of the buffered range
// can be misleading because of precision differences or when the current media has poorly
// aligned audio and video, which can cause values to be slightly off from what you would
// expect. This value is what we consider to be safe to use in such comparisons to account
// for these scenarios.

var SAFE_TIME_DELTA = TIME_FUDGE_FACTOR * 3;

var filterRanges = function filterRanges(timeRanges, predicate) {
  var results = [];
  var i;

  if (timeRanges && timeRanges.length) {
    // Search for ranges that match the predicate
    for (i = 0; i < timeRanges.length; i++) {
      if (predicate(timeRanges.start(i), timeRanges.end(i))) {
        results.push([timeRanges.start(i), timeRanges.end(i)]);
      }
    }
  }

  return videojs$1.createTimeRanges(results);
};
/**
 * Attempts to find the buffered TimeRange that contains the specified
 * time.
 *
 * @param {TimeRanges} buffered - the TimeRanges object to query
 * @param {number} time  - the time to filter on.
 * @return {TimeRanges} a new TimeRanges object
 */


var findRange = function findRange(buffered, time) {
  return filterRanges(buffered, function (start, end) {
    return start - SAFE_TIME_DELTA <= time && end + SAFE_TIME_DELTA >= time;
  });
};
/**
 * Returns the TimeRanges that begin later than the specified time.
 *
 * @param {TimeRanges} timeRanges - the TimeRanges object to query
 * @param {number} time - the time to filter on.
 * @return {TimeRanges} a new TimeRanges object.
 */

var findNextRange = function findNextRange(timeRanges, time) {
  return filterRanges(timeRanges, function (start) {
    return start - TIME_FUDGE_FACTOR >= time;
  });
};
/**
 * Returns gaps within a list of TimeRanges
 *
 * @param {TimeRanges} buffered - the TimeRanges object
 * @return {TimeRanges} a TimeRanges object of gaps
 */

var findGaps = function findGaps(buffered) {
  if (buffered.length < 2) {
    return videojs$1.createTimeRanges();
  }

  var ranges = [];

  for (var i = 1; i < buffered.length; i++) {
    var start = buffered.end(i - 1);
    var end = buffered.start(i);
    ranges.push([start, end]);
  }

  return videojs$1.createTimeRanges(ranges);
};
/**
 * Calculate the intersection of two TimeRanges
 *
 * @param {TimeRanges} bufferA
 * @param {TimeRanges} bufferB
 * @return {TimeRanges} The interesection of `bufferA` with `bufferB`
 */

var bufferIntersection = function bufferIntersection(bufferA, bufferB) {
  var start = null;
  var end = null;
  var arity = 0;
  var extents = [];
  var ranges = [];

  if (!bufferA || !bufferA.length || !bufferB || !bufferB.length) {
    return videojs$1.createTimeRange();
  } // Handle the case where we have both buffers and create an
  // intersection of the two


  var count = bufferA.length; // A) Gather up all start and end times

  while (count--) {
    extents.push({
      time: bufferA.start(count),
      type: 'start'
    });
    extents.push({
      time: bufferA.end(count),
      type: 'end'
    });
  }

  count = bufferB.length;

  while (count--) {
    extents.push({
      time: bufferB.start(count),
      type: 'start'
    });
    extents.push({
      time: bufferB.end(count),
      type: 'end'
    });
  } // B) Sort them by time


  extents.sort(function (a, b) {
    return a.time - b.time;
  }); // C) Go along one by one incrementing arity for start and decrementing
  //    arity for ends

  for (count = 0; count < extents.length; count++) {
    if (extents[count].type === 'start') {
      arity++; // D) If arity is ever incremented to 2 we are entering an
      //    overlapping range

      if (arity === 2) {
        start = extents[count].time;
      }
    } else if (extents[count].type === 'end') {
      arity--; // E) If arity is ever decremented to 1 we leaving an
      //    overlapping range

      if (arity === 1) {
        end = extents[count].time;
      }
    } // F) Record overlapping ranges


    if (start !== null && end !== null) {
      ranges.push([start, end]);
      start = null;
      end = null;
    }
  }

  return videojs$1.createTimeRanges(ranges);
};
/**
 * Gets a human readable string for a TimeRange
 *
 * @param {TimeRange} range
 * @return {string} a human readable string
 */

var printableRange = function printableRange(range) {
  var strArr = [];

  if (!range || !range.length) {
    return '';
  }

  for (var i = 0; i < range.length; i++) {
    strArr.push(range.start(i) + ' => ' + range.end(i));
  }

  return strArr.join(', ');
};
/**
 * Calculates the amount of time left in seconds until the player hits the end of the
 * buffer and causes a rebuffer
 *
 * @param {TimeRange} buffered
 *        The state of the buffer
 * @param {Numnber} currentTime
 *        The current time of the player
 * @param {number} playbackRate
 *        The current playback rate of the player. Defaults to 1.
 * @return {number}
 *         Time until the player has to start rebuffering in seconds.
 * @function timeUntilRebuffer
 */

var timeUntilRebuffer = function timeUntilRebuffer(buffered, currentTime, playbackRate) {
  if (playbackRate === void 0) {
    playbackRate = 1;
  }

  var bufferedEnd = buffered.length ? buffered.end(buffered.length - 1) : 0;
  return (bufferedEnd - currentTime) / playbackRate;
};
/**
 * Converts a TimeRanges object into an array representation
 *
 * @param {TimeRanges} timeRanges
 * @return {Array}
 */

var timeRangesToArray = function timeRangesToArray(timeRanges) {
  var timeRangesList = [];

  for (var i = 0; i < timeRanges.length; i++) {
    timeRangesList.push({
      start: timeRanges.start(i),
      end: timeRanges.end(i)
    });
  }

  return timeRangesList;
};
/**
 * Determines if two time range objects are different.
 *
 * @param {TimeRange} a
 *        the first time range object to check
 *
 * @param {TimeRange} b
 *        the second time range object to check
 *
 * @return {Boolean}
 *         Whether the time range objects differ
 */

var isRangeDifferent = function isRangeDifferent(a, b) {
  // same object
  if (a === b) {
    return false;
  } // one or the other is undefined


  if (!a && b || !b && a) {
    return true;
  } // length is different


  if (a.length !== b.length) {
    return true;
  } // see if any start/end pair is different


  for (var i = 0; i < a.length; i++) {
    if (a.start(i) !== b.start(i) || a.end(i) !== b.end(i)) {
      return true;
    }
  } // if the length and every pair is the same
  // this is the same time range


  return false;
};

/**
 * @file playlist.js
 *
 * Playlist related utilities.
 */
var createTimeRange = videojs$1.createTimeRange;
/**
 * walk backward until we find a duration we can use
 * or return a failure
 *
 * @param {Playlist} playlist the playlist to walk through
 * @param {Number} endSequence the mediaSequence to stop walking on
 */

var backwardDuration = function backwardDuration(playlist, endSequence) {
  var result = 0;
  var i = endSequence - playlist.mediaSequence; // if a start time is available for segment immediately following
  // the interval, use it

  var segment = playlist.segments[i]; // Walk backward until we find the latest segment with timeline
  // information that is earlier than endSequence

  if (segment) {
    if (typeof segment.start !== 'undefined') {
      return {
        result: segment.start,
        precise: true
      };
    }

    if (typeof segment.end !== 'undefined') {
      return {
        result: segment.end - segment.duration,
        precise: true
      };
    }
  }

  while (i--) {
    segment = playlist.segments[i];

    if (typeof segment.end !== 'undefined') {
      return {
        result: result + segment.end,
        precise: true
      };
    }

    result += segment.duration;

    if (typeof segment.start !== 'undefined') {
      return {
        result: result + segment.start,
        precise: true
      };
    }
  }

  return {
    result: result,
    precise: false
  };
};
/**
 * walk forward until we find a duration we can use
 * or return a failure
 *
 * @param {Playlist} playlist the playlist to walk through
 * @param {number} endSequence the mediaSequence to stop walking on
 */


var forwardDuration = function forwardDuration(playlist, endSequence) {
  var result = 0;
  var segment;
  var i = endSequence - playlist.mediaSequence; // Walk forward until we find the earliest segment with timeline
  // information

  for (; i < playlist.segments.length; i++) {
    segment = playlist.segments[i];

    if (typeof segment.start !== 'undefined') {
      return {
        result: segment.start - result,
        precise: true
      };
    }

    result += segment.duration;

    if (typeof segment.end !== 'undefined') {
      return {
        result: segment.end - result,
        precise: true
      };
    }
  } // indicate we didn't find a useful duration estimate


  return {
    result: -1,
    precise: false
  };
};
/**
  * Calculate the media duration from the segments associated with a
  * playlist. The duration of a subinterval of the available segments
  * may be calculated by specifying an end index.
  *
  * @param {Object} playlist a media playlist object
  * @param {number=} endSequence an exclusive upper boundary
  * for the playlist.  Defaults to playlist length.
  * @param {number} expired the amount of time that has dropped
  * off the front of the playlist in a live scenario
  * @return {number} the duration between the first available segment
  * and end index.
  */


var intervalDuration = function intervalDuration(playlist, endSequence, expired) {
  if (typeof endSequence === 'undefined') {
    endSequence = playlist.mediaSequence + playlist.segments.length;
  }

  if (endSequence < playlist.mediaSequence) {
    return 0;
  } // do a backward walk to estimate the duration


  var backward = backwardDuration(playlist, endSequence);

  if (backward.precise) {
    // if we were able to base our duration estimate on timing
    // information provided directly from the Media Source, return
    // it
    return backward.result;
  } // walk forward to see if a precise duration estimate can be made
  // that way


  var forward = forwardDuration(playlist, endSequence);

  if (forward.precise) {
    // we found a segment that has been buffered and so it's
    // position is known precisely
    return forward.result;
  } // return the less-precise, playlist-based duration estimate


  return backward.result + expired;
};
/**
  * Calculates the duration of a playlist. If a start and end index
  * are specified, the duration will be for the subset of the media
  * timeline between those two indices. The total duration for live
  * playlists is always Infinity.
  *
  * @param {Object} playlist a media playlist object
  * @param {number=} endSequence an exclusive upper
  * boundary for the playlist. Defaults to the playlist media
  * sequence number plus its length.
  * @param {number=} expired the amount of time that has
  * dropped off the front of the playlist in a live scenario
  * @return {number} the duration between the start index and end
  * index.
  */


var duration = function duration(playlist, endSequence, expired) {
  if (!playlist) {
    return 0;
  }

  if (typeof expired !== 'number') {
    expired = 0;
  } // if a slice of the total duration is not requested, use
  // playlist-level duration indicators when they're present


  if (typeof endSequence === 'undefined') {
    // if present, use the duration specified in the playlist
    if (playlist.totalDuration) {
      return playlist.totalDuration;
    } // duration should be Infinity for live playlists


    if (!playlist.endList) {
      return window_1.Infinity;
    }
  } // calculate the total duration based on the segment durations


  return intervalDuration(playlist, endSequence, expired);
};
/**
  * Calculate the time between two indexes in the current playlist
  * neight the start- nor the end-index need to be within the current
  * playlist in which case, the targetDuration of the playlist is used
  * to approximate the durations of the segments
  *
  * @param {Object} playlist a media playlist object
  * @param {number} startIndex
  * @param {number} endIndex
  * @return {number} the number of seconds between startIndex and endIndex
  */

var sumDurations = function sumDurations(playlist, startIndex, endIndex) {
  var durations = 0;

  if (startIndex > endIndex) {
    var _ref = [endIndex, startIndex];
    startIndex = _ref[0];
    endIndex = _ref[1];
  }

  if (startIndex < 0) {
    for (var i = startIndex; i < Math.min(0, endIndex); i++) {
      durations += playlist.targetDuration;
    }

    startIndex = 0;
  }

  for (var _i = startIndex; _i < endIndex; _i++) {
    durations += playlist.segments[_i].duration;
  }

  return durations;
};
/**
 * Determines the media index of the segment corresponding to the safe edge of the live
 * window which is the duration of the last segment plus 2 target durations from the end
 * of the playlist.
 *
 * A liveEdgePadding can be provided which will be used instead of calculating the safe live edge.
 * This corresponds to suggestedPresentationDelay in DASH manifests.
 *
 * @param {Object} playlist
 *        a media playlist object
 * @param {number} [liveEdgePadding]
 *        A number in seconds indicating how far from the end we want to be.
 *        If provided, this value is used instead of calculating the safe live index from the target durations.
 *        Corresponds to suggestedPresentationDelay in DASH manifests.
 * @return {number}
 *         The media index of the segment at the safe live point. 0 if there is no "safe"
 *         point.
 * @function safeLiveIndex
 */

var safeLiveIndex = function safeLiveIndex(playlist, liveEdgePadding) {
  if (!playlist.segments.length) {
    return 0;
  }

  var i = playlist.segments.length;
  var lastSegmentDuration = playlist.segments[i - 1].duration || playlist.targetDuration;
  var safeDistance = typeof liveEdgePadding === 'number' ? liveEdgePadding : lastSegmentDuration + playlist.targetDuration * 2;

  if (safeDistance === 0) {
    return i;
  }

  var distanceFromEnd = 0;

  while (i--) {
    distanceFromEnd += playlist.segments[i].duration;

    if (distanceFromEnd >= safeDistance) {
      break;
    }
  }

  return Math.max(0, i);
};
/**
 * Calculates the playlist end time
 *
 * @param {Object} playlist a media playlist object
 * @param {number=} expired the amount of time that has
 *                  dropped off the front of the playlist in a live scenario
 * @param {boolean|false} useSafeLiveEnd a boolean value indicating whether or not the
 *                        playlist end calculation should consider the safe live end
 *                        (truncate the playlist end by three segments). This is normally
 *                        used for calculating the end of the playlist's seekable range.
 *                        This takes into account the value of liveEdgePadding.
 *                        Setting liveEdgePadding to 0 is equivalent to setting this to false.
 * @param {number} liveEdgePadding a number indicating how far from the end of the playlist we should be in seconds.
 *                 If this is provided, it is used in the safe live end calculation.
 *                 Setting useSafeLiveEnd=false or liveEdgePadding=0 are equivalent.
 *                 Corresponds to suggestedPresentationDelay in DASH manifests.
 * @return {number} the end time of playlist
 * @function playlistEnd
 */

var playlistEnd = function playlistEnd(playlist, expired, useSafeLiveEnd, liveEdgePadding) {
  if (!playlist || !playlist.segments) {
    return null;
  }

  if (playlist.endList) {
    return duration(playlist);
  }

  if (expired === null) {
    return null;
  }

  expired = expired || 0;
  var endSequence = useSafeLiveEnd ? safeLiveIndex(playlist, liveEdgePadding) : playlist.segments.length;
  return intervalDuration(playlist, playlist.mediaSequence + endSequence, expired);
};
/**
  * Calculates the interval of time that is currently seekable in a
  * playlist. The returned time ranges are relative to the earliest
  * moment in the specified playlist that is still available. A full
  * seekable implementation for live streams would need to offset
  * these values by the duration of content that has expired from the
  * stream.
  *
  * @param {Object} playlist a media playlist object
  * dropped off the front of the playlist in a live scenario
  * @param {number=} expired the amount of time that has
  * dropped off the front of the playlist in a live scenario
  * @param {number} liveEdgePadding how far from the end of the playlist we should be in seconds.
  *        Corresponds to suggestedPresentationDelay in DASH manifests.
  * @return {TimeRanges} the periods of time that are valid targets
  * for seeking
  */

var seekable = function seekable(playlist, expired, liveEdgePadding) {
  var useSafeLiveEnd = true;
  var seekableStart = expired || 0;
  var seekableEnd = playlistEnd(playlist, expired, useSafeLiveEnd, liveEdgePadding);

  if (seekableEnd === null) {
    return createTimeRange();
  }

  return createTimeRange(seekableStart, seekableEnd);
};
/**
 * Determine the index and estimated starting time of the segment that
 * contains a specified playback position in a media playlist.
 *
 * @param {Object} playlist the media playlist to query
 * @param {number} currentTime The number of seconds since the earliest
 * possible position to determine the containing segment for
 * @param {number} startIndex
 * @param {number} startTime
 * @return {Object}
 */

var getMediaInfoForTime = function getMediaInfoForTime(playlist, currentTime, startIndex, startTime) {
  var i;
  var segment;
  var numSegments = playlist.segments.length;
  var time = currentTime - startTime;

  if (time < 0) {
    // Walk backward from startIndex in the playlist, adding durations
    // until we find a segment that contains `time` and return it
    if (startIndex > 0) {
      for (i = startIndex - 1; i >= 0; i--) {
        segment = playlist.segments[i];
        time += segment.duration + TIME_FUDGE_FACTOR;

        if (time > 0) {
          return {
            mediaIndex: i,
            startTime: startTime - sumDurations(playlist, startIndex, i)
          };
        }
      }
    } // We were unable to find a good segment within the playlist
    // so select the first segment


    return {
      mediaIndex: 0,
      startTime: currentTime
    };
  } // When startIndex is negative, we first walk forward to first segment
  // adding target durations. If we "run out of time" before getting to
  // the first segment, return the first segment


  if (startIndex < 0) {
    for (i = startIndex; i < 0; i++) {
      time -= playlist.targetDuration;

      if (time < 0) {
        return {
          mediaIndex: 0,
          startTime: currentTime
        };
      }
    }

    startIndex = 0;
  } // Walk forward from startIndex in the playlist, subtracting durations
  // until we find a segment that contains `time` and return it


  for (i = startIndex; i < numSegments; i++) {
    segment = playlist.segments[i];
    time -= segment.duration + TIME_FUDGE_FACTOR;

    if (time < 0) {
      return {
        mediaIndex: i,
        startTime: startTime + sumDurations(playlist, startIndex, i)
      };
    }
  } // We are out of possible candidates so load the last one...


  return {
    mediaIndex: numSegments - 1,
    startTime: currentTime
  };
};
/**
 * Check whether the playlist is blacklisted or not.
 *
 * @param {Object} playlist the media playlist object
 * @return {boolean} whether the playlist is blacklisted or not
 * @function isBlacklisted
 */

var isBlacklisted = function isBlacklisted(playlist) {
  return playlist.excludeUntil && playlist.excludeUntil > Date.now();
};
/**
 * Check whether the playlist is compatible with current playback configuration or has
 * been blacklisted permanently for being incompatible.
 *
 * @param {Object} playlist the media playlist object
 * @return {boolean} whether the playlist is incompatible or not
 * @function isIncompatible
 */

var isIncompatible = function isIncompatible(playlist) {
  return playlist.excludeUntil && playlist.excludeUntil === Infinity;
};
/**
 * Check whether the playlist is enabled or not.
 *
 * @param {Object} playlist the media playlist object
 * @return {boolean} whether the playlist is enabled or not
 * @function isEnabled
 */

var isEnabled = function isEnabled(playlist) {
  var blacklisted = isBlacklisted(playlist);
  return !playlist.disabled && !blacklisted;
};
/**
 * Check whether the playlist has been manually disabled through the representations api.
 *
 * @param {Object} playlist the media playlist object
 * @return {boolean} whether the playlist is disabled manually or not
 * @function isDisabled
 */

var isDisabled = function isDisabled(playlist) {
  return playlist.disabled;
};
/**
 * Returns whether the current playlist is an AES encrypted HLS stream
 *
 * @return {boolean} true if it's an AES encrypted HLS stream
 */

var isAes = function isAes(media) {
  for (var i = 0; i < media.segments.length; i++) {
    if (media.segments[i].key) {
      return true;
    }
  }

  return false;
};
/**
 * Checks if the playlist has a value for the specified attribute
 *
 * @param {string} attr
 *        Attribute to check for
 * @param {Object} playlist
 *        The media playlist object
 * @return {boolean}
 *         Whether the playlist contains a value for the attribute or not
 * @function hasAttribute
 */

var hasAttribute = function hasAttribute(attr, playlist) {
  return playlist.attributes && playlist.attributes[attr];
};
/**
 * Estimates the time required to complete a segment download from the specified playlist
 *
 * @param {number} segmentDuration
 *        Duration of requested segment
 * @param {number} bandwidth
 *        Current measured bandwidth of the player
 * @param {Object} playlist
 *        The media playlist object
 * @param {number=} bytesReceived
 *        Number of bytes already received for the request. Defaults to 0
 * @return {number|NaN}
 *         The estimated time to request the segment. NaN if bandwidth information for
 *         the given playlist is unavailable
 * @function estimateSegmentRequestTime
 */

var estimateSegmentRequestTime = function estimateSegmentRequestTime(segmentDuration, bandwidth, playlist, bytesReceived) {
  if (bytesReceived === void 0) {
    bytesReceived = 0;
  }

  if (!hasAttribute('BANDWIDTH', playlist)) {
    return NaN;
  }

  var size = segmentDuration * playlist.attributes.BANDWIDTH;
  return (size - bytesReceived * 8) / bandwidth;
};
/*
 * Returns whether the current playlist is the lowest rendition
 *
 * @return {Boolean} true if on lowest rendition
 */

var isLowestEnabledRendition = function isLowestEnabledRendition(master, media) {
  if (master.playlists.length === 1) {
    return true;
  }

  var currentBandwidth = media.attributes.BANDWIDTH || Number.MAX_VALUE;
  return master.playlists.filter(function (playlist) {
    if (!isEnabled(playlist)) {
      return false;
    }

    return (playlist.attributes.BANDWIDTH || 0) < currentBandwidth;
  }).length === 0;
}; // exports

var Playlist = {
  duration: duration,
  seekable: seekable,
  safeLiveIndex: safeLiveIndex,
  getMediaInfoForTime: getMediaInfoForTime,
  isEnabled: isEnabled,
  isDisabled: isDisabled,
  isBlacklisted: isBlacklisted,
  isIncompatible: isIncompatible,
  playlistEnd: playlistEnd,
  isAes: isAes,
  hasAttribute: hasAttribute,
  estimateSegmentRequestTime: estimateSegmentRequestTime,
  isLowestEnabledRendition: isLowestEnabledRendition
};

/**
 * @file xhr.js
 */
var videojsXHR = videojs$1.xhr,
    mergeOptions$1 = videojs$1.mergeOptions;

var callbackWrapper = function callbackWrapper(request, error, response, callback) {
  var reqResponse = request.responseType === 'arraybuffer' ? request.response : request.responseText;

  if (!error && reqResponse) {
    request.responseTime = Date.now();
    request.roundTripTime = request.responseTime - request.requestTime;
    request.bytesReceived = reqResponse.byteLength || reqResponse.length;

    if (!request.bandwidth) {
      request.bandwidth = Math.floor(request.bytesReceived / request.roundTripTime * 8 * 1000);
    }
  }

  if (response.headers) {
    request.responseHeaders = response.headers;
  } // videojs.xhr now uses a specific code on the error
  // object to signal that a request has timed out instead
  // of setting a boolean on the request object


  if (error && error.code === 'ETIMEDOUT') {
    request.timedout = true;
  } // videojs.xhr no longer considers status codes outside of 200 and 0
  // (for file uris) to be errors, but the old XHR did, so emulate that
  // behavior. Status 206 may be used in response to byterange requests.


  if (!error && !request.aborted && response.statusCode !== 200 && response.statusCode !== 206 && response.statusCode !== 0) {
    error = new Error('XHR Failed with a response of: ' + (request && (reqResponse || request.responseText)));
  }

  callback(error, request);
};

var xhrFactory = function xhrFactory() {
  var xhr = function XhrFunction(options, callback) {
    // Add a default timeout
    options = mergeOptions$1({
      timeout: 45e3
    }, options); // Allow an optional user-specified function to modify the option
    // object before we construct the xhr request

    var beforeRequest = XhrFunction.beforeRequest || videojs$1.Vhs.xhr.beforeRequest;

    if (beforeRequest && typeof beforeRequest === 'function') {
      var newOptions = beforeRequest(options);

      if (newOptions) {
        options = newOptions;
      }
    }

    var request = videojsXHR(options, function (error, response) {
      return callbackWrapper(request, error, response, callback);
    });
    var originalAbort = request.abort;

    request.abort = function () {
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


var byterangeStr = function byterangeStr(byterange) {
  // `byterangeEnd` is one less than `offset + length` because the HTTP range
  // header uses inclusive ranges
  var byterangeEnd = byterange.offset + byterange.length - 1;
  var byterangeStart = byterange.offset;
  return 'bytes=' + byterangeStart + '-' + byterangeEnd;
};
/**
 * Defines headers for use in the xhr request for a particular segment.
 *
 * @param {Object} segment - a simplified copy of the segmentInfo object
 *                           from SegmentLoader
 */


var segmentXhrHeaders = function segmentXhrHeaders(segment) {
  var headers = {};

  if (segment.byterange) {
    headers.Range = byterangeStr(segment.byterange);
  }

  return headers;
};

/**
 * @file bin-utils.js
 */

/**
 * convert a TimeRange to text
 *
 * @param {TimeRange} range the timerange to use for conversion
 * @param {number} i the iterator on the range to convert
 */
var textRange = function textRange(range, i) {
  return range.start(i) + '-' + range.end(i);
};
/**
 * format a number as hex string
 *
 * @param {number} e The number
 * @param {number} i the iterator
 */


var formatHexString = function formatHexString(e, i) {
  var value = e.toString(16);
  return '00'.substring(0, 2 - value.length) + value + (i % 2 ? ' ' : '');
};

var formatAsciiString = function formatAsciiString(e) {
  if (e >= 0x20 && e < 0x7e) {
    return String.fromCharCode(e);
  }

  return '.';
};
/**
 * Creates an object for sending to a web worker modifying properties that are TypedArrays
 * into a new object with seperated properties for the buffer, byteOffset, and byteLength.
 *
 * @param {Object} message
 *        Object of properties and values to send to the web worker
 * @return {Object}
 *         Modified message with TypedArray values expanded
 * @function createTransferableMessage
 */


var createTransferableMessage = function createTransferableMessage(message) {
  var transferable = {};
  Object.keys(message).forEach(function (key) {
    var value = message[key];

    if (ArrayBuffer.isView(value)) {
      transferable[key] = {
        bytes: value.buffer,
        byteOffset: value.byteOffset,
        byteLength: value.byteLength
      };
    } else {
      transferable[key] = value;
    }
  });
  return transferable;
};
/**
 * Returns a unique string identifier for a media initialization
 * segment.
 */

var initSegmentId = function initSegmentId(initSegment) {
  var byterange = initSegment.byterange || {
    length: Infinity,
    offset: 0
  };
  return [byterange.length, byterange.offset, initSegment.resolvedUri].join(',');
};
/**
 * Returns a unique string identifier for a media segment key.
 */

var segmentKeyId = function segmentKeyId(key) {
  return key.resolvedUri;
};
/**
 * utils to help dump binary data to the console
 */

var hexDump = function hexDump(data) {
  var bytes = Array.prototype.slice.call(data);
  var step = 16;
  var result = '';
  var hex;
  var ascii;

  for (var j = 0; j < bytes.length / step; j++) {
    hex = bytes.slice(j * step, j * step + step).map(formatHexString).join('');
    ascii = bytes.slice(j * step, j * step + step).map(formatAsciiString).join('');
    result += hex + ' ' + ascii + '\n';
  }

  return result;
};
var tagDump = function tagDump(_ref) {
  var bytes = _ref.bytes;
  return hexDump(bytes);
};
var textRanges = function textRanges(ranges) {
  var result = '';
  var i;

  for (i = 0; i < ranges.length; i++) {
    result += textRange(ranges, i) + ' ';
  }

  return result;
};

var utils$1 = /*#__PURE__*/Object.freeze({
  __proto__: null,
  createTransferableMessage: createTransferableMessage,
  initSegmentId: initSegmentId,
  segmentKeyId: segmentKeyId,
  hexDump: hexDump,
  tagDump: tagDump,
  textRanges: textRanges
});

// TODO handle fmp4 case where the timing info is accurate and doesn't involve transmux
// 25% was arbitrarily chosen, and may need to be refined over time.

var SEGMENT_END_FUDGE_PERCENT = 0.25;
/**
 * Converts a player time (any time that can be gotten/set from player.currentTime(),
 * e.g., any time within player.seekable().start(0) to player.seekable().end(0)) to a
 * program time (any time referencing the real world (e.g., EXT-X-PROGRAM-DATE-TIME)).
 *
 * The containing segment is required as the EXT-X-PROGRAM-DATE-TIME serves as an "anchor
 * point" (a point where we have a mapping from program time to player time, with player
 * time being the post transmux start of the segment).
 *
 * For more details, see [this doc](../../docs/program-time-from-player-time.md).
 *
 * @param {number} playerTime the player time
 * @param {Object} segment the segment which contains the player time
 * @return {Date} program time
 */

var playerTimeToProgramTime = function playerTimeToProgramTime(playerTime, segment) {
  if (!segment.dateTimeObject) {
    // Can't convert without an "anchor point" for the program time (i.e., a time that can
    // be used to map the start of a segment with a real world time).
    return null;
  }

  var transmuxerPrependedSeconds = segment.videoTimingInfo.transmuxerPrependedSeconds;
  var transmuxedStart = segment.videoTimingInfo.transmuxedPresentationStart; // get the start of the content from before old content is prepended

  var startOfSegment = transmuxedStart + transmuxerPrependedSeconds;
  var offsetFromSegmentStart = playerTime - startOfSegment;
  return new Date(segment.dateTimeObject.getTime() + offsetFromSegmentStart * 1000);
};
var originalSegmentVideoDuration = function originalSegmentVideoDuration(videoTimingInfo) {
  return videoTimingInfo.transmuxedPresentationEnd - videoTimingInfo.transmuxedPresentationStart - videoTimingInfo.transmuxerPrependedSeconds;
};
/**
 * Finds a segment that contains the time requested given as an ISO-8601 string. The
 * returned segment might be an estimate or an accurate match.
 *
 * @param {string} programTime The ISO-8601 programTime to find a match for
 * @param {Object} playlist A playlist object to search within
 */

var findSegmentForProgramTime = function findSegmentForProgramTime(programTime, playlist) {
  // Assumptions:
  //  - verifyProgramDateTimeTags has already been run
  //  - live streams have been started
  var dateTimeObject;

  try {
    dateTimeObject = new Date(programTime);
  } catch (e) {
    return null;
  }

  if (!playlist || !playlist.segments || playlist.segments.length === 0) {
    return null;
  }

  var segment = playlist.segments[0];

  if (dateTimeObject < segment.dateTimeObject) {
    // Requested time is before stream start.
    return null;
  }

  for (var i = 0; i < playlist.segments.length - 1; i++) {
    segment = playlist.segments[i];
    var nextSegmentStart = playlist.segments[i + 1].dateTimeObject;

    if (dateTimeObject < nextSegmentStart) {
      break;
    }
  }

  var lastSegment = playlist.segments[playlist.segments.length - 1];
  var lastSegmentStart = lastSegment.dateTimeObject;
  var lastSegmentDuration = lastSegment.videoTimingInfo ? originalSegmentVideoDuration(lastSegment.videoTimingInfo) : lastSegment.duration + lastSegment.duration * SEGMENT_END_FUDGE_PERCENT;
  var lastSegmentEnd = new Date(lastSegmentStart.getTime() + lastSegmentDuration * 1000);

  if (dateTimeObject > lastSegmentEnd) {
    // Beyond the end of the stream, or our best guess of the end of the stream.
    return null;
  }

  if (dateTimeObject > lastSegmentStart) {
    segment = lastSegment;
  }

  return {
    segment: segment,
    estimatedStart: segment.videoTimingInfo ? segment.videoTimingInfo.transmuxedPresentationStart : Playlist.duration(playlist, playlist.mediaSequence + playlist.segments.indexOf(segment)),
    // Although, given that all segments have accurate date time objects, the segment
    // selected should be accurate, unless the video has been transmuxed at some point
    // (determined by the presence of the videoTimingInfo object), the segment's "player
    // time" (the start time in the player) can't be considered accurate.
    type: segment.videoTimingInfo ? 'accurate' : 'estimate'
  };
};
/**
 * Finds a segment that contains the given player time(in seconds).
 *
 * @param {number} time The player time to find a match for
 * @param {Object} playlist A playlist object to search within
 */

var findSegmentForPlayerTime = function findSegmentForPlayerTime(time, playlist) {
  // Assumptions:
  // - there will always be a segment.duration
  // - we can start from zero
  // - segments are in time order
  if (!playlist || !playlist.segments || playlist.segments.length === 0) {
    return null;
  }

  var segmentEnd = 0;
  var segment;

  for (var i = 0; i < playlist.segments.length; i++) {
    segment = playlist.segments[i]; // videoTimingInfo is set after the segment is downloaded and transmuxed, and
    // should contain the most accurate values we have for the segment's player times.
    //
    // Use the accurate transmuxedPresentationEnd value if it is available, otherwise fall
    // back to an estimate based on the manifest derived (inaccurate) segment.duration, to
    // calculate an end value.

    segmentEnd = segment.videoTimingInfo ? segment.videoTimingInfo.transmuxedPresentationEnd : segmentEnd + segment.duration;

    if (time <= segmentEnd) {
      break;
    }
  }

  var lastSegment = playlist.segments[playlist.segments.length - 1];

  if (lastSegment.videoTimingInfo && lastSegment.videoTimingInfo.transmuxedPresentationEnd < time) {
    // The time requested is beyond the stream end.
    return null;
  }

  if (time > segmentEnd) {
    // The time is within or beyond the last segment.
    //
    // Check to see if the time is beyond a reasonable guess of the end of the stream.
    if (time > segmentEnd + lastSegment.duration * SEGMENT_END_FUDGE_PERCENT) {
      // Technically, because the duration value is only an estimate, the time may still
      // exist in the last segment, however, there isn't enough information to make even
      // a reasonable estimate.
      return null;
    }

    segment = lastSegment;
  }

  return {
    segment: segment,
    estimatedStart: segment.videoTimingInfo ? segment.videoTimingInfo.transmuxedPresentationStart : segmentEnd - segment.duration,
    // Because videoTimingInfo is only set after transmux, it is the only way to get
    // accurate timing values.
    type: segment.videoTimingInfo ? 'accurate' : 'estimate'
  };
};
/**
 * Gives the offset of the comparisonTimestamp from the programTime timestamp in seconds.
 * If the offset returned is positive, the programTime occurs after the
 * comparisonTimestamp.
 * If the offset is negative, the programTime occurs before the comparisonTimestamp.
 *
 * @param {string} comparisonTimeStamp An ISO-8601 timestamp to compare against
 * @param {string} programTime The programTime as an ISO-8601 string
 * @return {number} offset
 */

var getOffsetFromTimestamp = function getOffsetFromTimestamp(comparisonTimeStamp, programTime) {
  var segmentDateTime;
  var programDateTime;

  try {
    segmentDateTime = new Date(comparisonTimeStamp);
    programDateTime = new Date(programTime);
  } catch (e) {// TODO handle error
  }

  var segmentTimeEpoch = segmentDateTime.getTime();
  var programTimeEpoch = programDateTime.getTime();
  return (programTimeEpoch - segmentTimeEpoch) / 1000;
};
/**
 * Checks that all segments in this playlist have programDateTime tags.
 *
 * @param {Object} playlist A playlist object
 */

var verifyProgramDateTimeTags = function verifyProgramDateTimeTags(playlist) {
  if (!playlist.segments || playlist.segments.length === 0) {
    return false;
  }

  for (var i = 0; i < playlist.segments.length; i++) {
    var segment = playlist.segments[i];

    if (!segment.dateTimeObject) {
      return false;
    }
  }

  return true;
};
/**
 * Returns the programTime of the media given a playlist and a playerTime.
 * The playlist must have programDateTime tags for a programDateTime tag to be returned.
 * If the segments containing the time requested have not been buffered yet, an estimate
 * may be returned to the callback.
 *
 * @param {Object} args
 * @param {Object} args.playlist A playlist object to search within
 * @param {number} time A playerTime in seconds
 * @param {Function} callback(err, programTime)
 * @return {string} err.message A detailed error message
 * @return {Object} programTime
 * @return {number} programTime.mediaSeconds The streamTime in seconds
 * @return {string} programTime.programDateTime The programTime as an ISO-8601 String
 */

var getProgramTime = function getProgramTime(_ref) {
  var playlist = _ref.playlist,
      _ref$time = _ref.time,
      time = _ref$time === void 0 ? undefined : _ref$time,
      callback = _ref.callback;

  if (!callback) {
    throw new Error('getProgramTime: callback must be provided');
  }

  if (!playlist || time === undefined) {
    return callback({
      message: 'getProgramTime: playlist and time must be provided'
    });
  }

  var matchedSegment = findSegmentForPlayerTime(time, playlist);

  if (!matchedSegment) {
    return callback({
      message: 'valid programTime was not found'
    });
  }

  if (matchedSegment.type === 'estimate') {
    return callback({
      message: 'Accurate programTime could not be determined.' + ' Please seek to e.seekTime and try again',
      seekTime: matchedSegment.estimatedStart
    });
  }

  var programTimeObject = {
    mediaSeconds: time
  };
  var programTime = playerTimeToProgramTime(time, matchedSegment.segment);

  if (programTime) {
    programTimeObject.programDateTime = programTime.toISOString();
  }

  return callback(null, programTimeObject);
};
/**
 * Seeks in the player to a time that matches the given programTime ISO-8601 string.
 *
 * @param {Object} args
 * @param {string} args.programTime A programTime to seek to as an ISO-8601 String
 * @param {Object} args.playlist A playlist to look within
 * @param {number} args.retryCount The number of times to try for an accurate seek. Default is 2.
 * @param {Function} args.seekTo A method to perform a seek
 * @param {boolean} args.pauseAfterSeek Whether to end in a paused state after seeking. Default is true.
 * @param {Object} args.tech The tech to seek on
 * @param {Function} args.callback(err, newTime) A callback to return the new time to
 * @return {string} err.message A detailed error message
 * @return {number} newTime The exact time that was seeked to in seconds
 */

var seekToProgramTime = function seekToProgramTime(_ref2) {
  var programTime = _ref2.programTime,
      playlist = _ref2.playlist,
      _ref2$retryCount = _ref2.retryCount,
      retryCount = _ref2$retryCount === void 0 ? 2 : _ref2$retryCount,
      seekTo = _ref2.seekTo,
      _ref2$pauseAfterSeek = _ref2.pauseAfterSeek,
      pauseAfterSeek = _ref2$pauseAfterSeek === void 0 ? true : _ref2$pauseAfterSeek,
      tech = _ref2.tech,
      callback = _ref2.callback;

  if (!callback) {
    throw new Error('seekToProgramTime: callback must be provided');
  }

  if (typeof programTime === 'undefined' || !playlist || !seekTo) {
    return callback({
      message: 'seekToProgramTime: programTime, seekTo and playlist must be provided'
    });
  }

  if (!playlist.endList && !tech.hasStarted_) {
    return callback({
      message: 'player must be playing a live stream to start buffering'
    });
  }

  if (!verifyProgramDateTimeTags(playlist)) {
    return callback({
      message: 'programDateTime tags must be provided in the manifest ' + playlist.resolvedUri
    });
  }

  var matchedSegment = findSegmentForProgramTime(programTime, playlist); // no match

  if (!matchedSegment) {
    return callback({
      message: programTime + " was not found in the stream"
    });
  }

  var segment = matchedSegment.segment;
  var mediaOffset = getOffsetFromTimestamp(segment.dateTimeObject, programTime);

  if (matchedSegment.type === 'estimate') {
    // we've run out of retries
    if (retryCount === 0) {
      return callback({
        message: programTime + " is not buffered yet. Try again"
      });
    }

    seekTo(matchedSegment.estimatedStart + mediaOffset);
    tech.one('seeked', function () {
      seekToProgramTime({
        programTime: programTime,
        playlist: playlist,
        retryCount: retryCount - 1,
        seekTo: seekTo,
        pauseAfterSeek: pauseAfterSeek,
        tech: tech,
        callback: callback
      });
    });
    return;
  } // Since the segment.start value is determined from the buffered end or ending time
  // of the prior segment, the seekToTime doesn't need to account for any transmuxer
  // modifications.


  var seekToTime = segment.start + mediaOffset;

  var seekedCallback = function seekedCallback() {
    return callback(null, tech.currentTime());
  }; // listen for seeked event


  tech.one('seeked', seekedCallback); // pause before seeking as video.js will restore this state

  if (pauseAfterSeek) {
    tech.pause();
  }

  seekTo(seekToTime);
};

// which will only happen if the request is complete.

var callbackOnCompleted = function callbackOnCompleted(request, cb) {
  if (request.readyState === 4) {
    return cb();
  }

  return;
};

var containerRequest = function containerRequest(uri, xhr, cb) {
  var bytes = [];
  var id3Offset;
  var finished = false;

  var endRequestAndCallback = function endRequestAndCallback(err, req, type, _bytes) {
    req.abort();
    finished = true;
    return cb(err, req, type, _bytes);
  };

  var progressListener = function progressListener(error, request) {
    if (finished) {
      return;
    }

    if (error) {
      return endRequestAndCallback(error, request, '', bytes);
    } // grap the new part of content that was just downloaded


    var newPart = request.responseText.substring(bytes && bytes.byteLength || 0, request.responseText.length); // add that onto bytes

    bytes = byteHelpers.concatTypedArrays(bytes, byteHelpers.stringToBytes(newPart, true));
    id3Offset = id3Offset || containers.getId3Offset(bytes); // we need at least 10 bytes to determine a type
    // or we need at least two bytes after an id3Offset

    if (bytes.length < 10 || id3Offset && bytes.length < id3Offset + 2) {
      return callbackOnCompleted(request, function () {
        return endRequestAndCallback(error, request, '', bytes);
      });
    }

    var type = containers.detectContainerForBytes(bytes); // if this looks like a ts segment but we don't have enough data
    // to see the second sync byte, wait until we have enough data
    // before declaring it ts

    if (type === 'ts' && bytes.length < 188) {
      return callbackOnCompleted(request, function () {
        return endRequestAndCallback(error, request, '', bytes);
      });
    } // this may be an unsynced ts segment
    // wait for 376 bytes before detecting no container


    if (!type && bytes.length < 376) {
      return callbackOnCompleted(request, function () {
        return endRequestAndCallback(error, request, '', bytes);
      });
    }

    return endRequestAndCallback(null, request, type, bytes);
  };

  var options = {
    uri: uri,
    beforeSend: function beforeSend(request) {
      // this forces the browser to pass the bytes to us unprocessed
      request.overrideMimeType('text/plain; charset=x-user-defined');
      request.addEventListener('progress', function (_ref) {
        var total = _ref.total,
            loaded = _ref.loaded;
        return callbackWrapper(request, null, {
          statusCode: request.status
        }, progressListener);
      });
    }
  };
  var request = xhr(options, function (error, response) {
    return callbackWrapper(request, error, response, progressListener);
  });
  return request;
};

var EventTarget$1 = videojs$1.EventTarget,
    mergeOptions$2 = videojs$1.mergeOptions;
/**
 * Parses the master XML string and updates playlist URI references.
 *
 * @param {Object} config
 *        Object of arguments
 * @param {string} config.masterXml
 *        The mpd XML
 * @param {string} config.srcUrl
 *        The mpd URL
 * @param {Date} config.clientOffset
 *         A time difference between server and client
 * @param {Object} config.sidxMapping
 *        SIDX mappings for moof/mdat URIs and byte ranges
 * @return {Object}
 *         The parsed mpd manifest object
 */

var parseMasterXml = function parseMasterXml(_ref) {
  var masterXml = _ref.masterXml,
      srcUrl = _ref.srcUrl,
      clientOffset = _ref.clientOffset,
      sidxMapping = _ref.sidxMapping;
  var master = parse(masterXml, {
    manifestUri: srcUrl,
    clientOffset: clientOffset,
    sidxMapping: sidxMapping
  });
  addPropertiesToMaster(master, srcUrl);
  return master;
};
/**
 * Returns a new master manifest that is the result of merging an updated master manifest
 * into the original version.
 *
 * @param {Object} oldMaster
 *        The old parsed mpd object
 * @param {Object} newMaster
 *        The updated parsed mpd object
 * @return {Object}
 *         A new object representing the original master manifest with the updated media
 *         playlists merged in
 */

var updateMaster$1 = function updateMaster$1(oldMaster, newMaster) {
  var noChanges = true;
  var update = mergeOptions$2(oldMaster, {
    // These are top level properties that can be updated
    duration: newMaster.duration,
    minimumUpdatePeriod: newMaster.minimumUpdatePeriod
  }); // First update the playlists in playlist list

  for (var i = 0; i < newMaster.playlists.length; i++) {
    var playlistUpdate = updateMaster(update, newMaster.playlists[i]);

    if (playlistUpdate) {
      update = playlistUpdate;
      noChanges = false;
    }
  } // Then update media group playlists


  forEachMediaGroup(newMaster, function (properties, type, group, label) {
    if (properties.playlists && properties.playlists.length) {
      var id = properties.playlists[0].id;

      var _playlistUpdate = updateMaster(update, properties.playlists[0]);

      if (_playlistUpdate) {
        update = _playlistUpdate; // update the playlist reference within media groups

        update.mediaGroups[type][group][label].playlists[0] = update.playlists[id];
        noChanges = false;
      }
    }
  });

  if (newMaster.minimumUpdatePeriod !== oldMaster.minimumUpdatePeriod) {
    noChanges = false;
  }

  if (noChanges) {
    return null;
  }

  return update;
};
var generateSidxKey = function generateSidxKey(sidxInfo) {
  // should be non-inclusive
  var sidxByteRangeEnd = sidxInfo.byterange.offset + sidxInfo.byterange.length - 1;
  return sidxInfo.uri + '-' + sidxInfo.byterange.offset + '-' + sidxByteRangeEnd;
}; // SIDX should be equivalent if the URI and byteranges of the SIDX match.
// If the SIDXs have maps, the two maps should match,
// both `a` and `b` missing SIDXs is considered matching.
// If `a` or `b` but not both have a map, they aren't matching.

var equivalentSidx = function equivalentSidx(a, b) {
  var neitherMap = Boolean(!a.map && !b.map);
  var equivalentMap = neitherMap || Boolean(a.map && b.map && a.map.byterange.offset === b.map.byterange.offset && a.map.byterange.length === b.map.byterange.length);
  return equivalentMap && a.uri === b.uri && a.byterange.offset === b.byterange.offset && a.byterange.length === b.byterange.length;
}; // exported for testing


var compareSidxEntry = function compareSidxEntry(playlists, oldSidxMapping) {
  var newSidxMapping = {};

  for (var id in playlists) {
    var playlist = playlists[id];
    var currentSidxInfo = playlist.sidx;

    if (currentSidxInfo) {
      var key = generateSidxKey(currentSidxInfo);

      if (!oldSidxMapping[key]) {
        break;
      }

      var savedSidxInfo = oldSidxMapping[key].sidxInfo;

      if (equivalentSidx(savedSidxInfo, currentSidxInfo)) {
        newSidxMapping[key] = oldSidxMapping[key];
      }
    }
  }

  return newSidxMapping;
};
/**
 *  A function that filters out changed items as they need to be requested separately.
 *
 *  The method is exported for testing
 *
 *  @param {Object} masterXml the mpd XML
 *  @param {string} srcUrl the mpd url
 *  @param {Date} clientOffset a time difference between server and client (passed through and not used)
 *  @param {Object} oldSidxMapping the SIDX to compare against
 */

var filterChangedSidxMappings = function filterChangedSidxMappings(masterXml, srcUrl, clientOffset, oldSidxMapping) {
  // Don't pass current sidx mapping
  var master = parse(masterXml, {
    manifestUri: srcUrl,
    clientOffset: clientOffset
  });
  var videoSidx = compareSidxEntry(master.playlists, oldSidxMapping);
  var mediaGroupSidx = videoSidx;
  forEachMediaGroup(master, function (properties, mediaType, groupKey, labelKey) {
    if (properties.playlists && properties.playlists.length) {
      var playlists = properties.playlists;
      mediaGroupSidx = mergeOptions$2(mediaGroupSidx, compareSidxEntry(playlists, oldSidxMapping));
    }
  });
  return mediaGroupSidx;
}; // exported for testing

var requestSidx_ = function requestSidx_(loader, sidxRange, playlist, xhr, options, finishProcessingFn) {
  var sidxInfo = {
    // resolve the segment URL relative to the playlist
    uri: resolveManifestRedirect(options.handleManifestRedirects, sidxRange.resolvedUri),
    // resolvedUri: sidxRange.resolvedUri,
    byterange: sidxRange.byterange,
    // the segment's playlist
    playlist: playlist
  };
  var sidxRequestOptions = videojs$1.mergeOptions(sidxInfo, {
    responseType: 'arraybuffer',
    headers: segmentXhrHeaders(sidxInfo)
  });
  return containerRequest(sidxInfo.uri, xhr, function (err, request, container, bytes) {
    if (err) {
      return finishProcessingFn(err, request);
    }

    if (!container || container !== 'mp4') {
      return finishProcessingFn({
        status: request.status,
        message: "Unsupported " + (container || 'unknown') + " container type for sidx segment at URL: " + sidxInfo.uri,
        // response is just bytes in this case
        // but we really don't want to return that.
        response: '',
        playlist: playlist,
        internal: true,
        blacklistDuration: Infinity,
        // MEDIA_ERR_NETWORK
        code: 2
      }, request);
    } // if we already downloaded the sidx bytes in the container request, use them


    var _sidxInfo$byterange = sidxInfo.byterange,
        offset = _sidxInfo$byterange.offset,
        length = _sidxInfo$byterange.length;

    if (bytes.length >= length + offset) {
      return finishProcessingFn(err, {
        response: bytes.subarray(offset, offset + length),
        status: request.status,
        uri: request.uri
      });
    } // otherwise request sidx bytes


    loader.request = xhr(sidxRequestOptions, finishProcessingFn);
  });
};

var DashPlaylistLoader = /*#__PURE__*/function (_EventTarget) {
  inheritsLoose(DashPlaylistLoader, _EventTarget);

  // DashPlaylistLoader must accept either a src url or a playlist because subsequent
  // playlist loader setups from media groups will expect to be able to pass a playlist
  // (since there aren't external URLs to media playlists with DASH)
  function DashPlaylistLoader(srcUrlOrPlaylist, vhs, options, masterPlaylistLoader) {
    var _this;

    if (options === void 0) {
      options = {};
    }

    _this = _EventTarget.call(this) || this;
    var _options = options,
        _options$withCredenti = _options.withCredentials,
        withCredentials = _options$withCredenti === void 0 ? false : _options$withCredenti,
        _options$handleManife = _options.handleManifestRedirects,
        handleManifestRedirects = _options$handleManife === void 0 ? false : _options$handleManife;
    _this.vhs_ = vhs;
    _this.withCredentials = withCredentials;
    _this.handleManifestRedirects = handleManifestRedirects;

    if (!srcUrlOrPlaylist) {
      throw new Error('A non-empty playlist URL or object is required');
    } // event naming?


    _this.on('minimumUpdatePeriod', function () {
      _this.refreshXml_();
    }); // live playlist staleness timeout


    _this.on('mediaupdatetimeout', function () {
      _this.refreshMedia_(_this.media().id);
    });

    _this.state = 'HAVE_NOTHING';
    _this.loadedPlaylists_ = {}; // initialize the loader state
    // The masterPlaylistLoader will be created with a string

    if (typeof srcUrlOrPlaylist === 'string') {
      _this.srcUrl = srcUrlOrPlaylist; // TODO: reset sidxMapping between period changes
      // once multi-period is refactored

      _this.sidxMapping_ = {};
      return assertThisInitialized(_this);
    }

    _this.setupChildLoader(masterPlaylistLoader, srcUrlOrPlaylist);

    return _this;
  }

  var _proto = DashPlaylistLoader.prototype;

  _proto.setupChildLoader = function setupChildLoader(masterPlaylistLoader, playlist) {
    this.masterPlaylistLoader_ = masterPlaylistLoader;
    this.childPlaylist_ = playlist;
  };

  _proto.dispose = function dispose() {
    this.trigger('dispose');
    this.stopRequest();
    this.loadedPlaylists_ = {};
    window_1.clearTimeout(this.minimumUpdatePeriodTimeout_);
    window_1.clearTimeout(this.mediaRequest_);
    window_1.clearTimeout(this.mediaUpdateTimeout);
    this.off();
  };

  _proto.hasPendingRequest = function hasPendingRequest() {
    return this.request || this.mediaRequest_;
  };

  _proto.stopRequest = function stopRequest() {
    if (this.request) {
      var oldRequest = this.request;
      this.request = null;
      oldRequest.onreadystatechange = null;
      oldRequest.abort();
    }
  };

  _proto.sidxRequestFinished_ = function sidxRequestFinished_(playlist, master, startingState, doneFn) {
    var _this2 = this;

    return function (err, request) {
      // disposed
      if (!_this2.request) {
        return;
      } // pending request is cleared


      _this2.request = null;

      if (err) {
        // use the provided error or create one
        // see requestSidx_ for the container request
        // that can cause this.
        _this2.error = typeof err === 'object' ? err : {
          status: request.status,
          message: 'DASH playlist request error at URL: ' + playlist.uri,
          response: request.response,
          // MEDIA_ERR_NETWORK
          code: 2
        };

        if (startingState) {
          _this2.state = startingState;
        }

        _this2.trigger('error');

        return;
      }

      var bytes = byteHelpers.toUint8(request.response);
      var sidx = parseSidx_1(bytes.subarray(8));
      return doneFn(master, sidx);
    };
  };

  _proto.media = function media(playlist) {
    var _this3 = this;

    // getter
    if (!playlist) {
      return this.media_;
    } // setter


    if (this.state === 'HAVE_NOTHING') {
      throw new Error('Cannot switch media playlist from ' + this.state);
    }

    var startingState = this.state; // find the playlist object if the target playlist has been specified by URI

    if (typeof playlist === 'string') {
      if (!this.master.playlists[playlist]) {
        throw new Error('Unknown playlist URI: ' + playlist);
      }

      playlist = this.master.playlists[playlist];
    }

    var mediaChange = !this.media_ || playlist.id !== this.media_.id; // switch to previously loaded playlists immediately

    if (mediaChange && this.loadedPlaylists_[playlist.id] && this.loadedPlaylists_[playlist.id].endList) {
      this.state = 'HAVE_METADATA';
      this.media_ = playlist; // trigger media change if the active media has been updated

      if (mediaChange) {
        this.trigger('mediachanging');
        this.trigger('mediachange');
      }

      return;
    } // switching to the active playlist is a no-op


    if (!mediaChange) {
      return;
    } // switching from an already loaded playlist


    if (this.media_) {
      this.trigger('mediachanging');
    }

    if (!playlist.sidx) {
      // Continue asynchronously if there is no sidx
      // wait one tick to allow haveMaster to run first on a child loader
      this.mediaRequest_ = window_1.setTimeout(this.haveMetadata.bind(this, {
        startingState: startingState,
        playlist: playlist
      }), 0); // exit early and don't do sidx work

      return;
    } // we have sidx mappings


    var oldMaster;
    var sidxMapping; // sidxMapping is used when parsing the masterXml, so store
    // it on the masterPlaylistLoader

    if (this.masterPlaylistLoader_) {
      oldMaster = this.masterPlaylistLoader_.master;
      sidxMapping = this.masterPlaylistLoader_.sidxMapping_;
    } else {
      oldMaster = this.master;
      sidxMapping = this.sidxMapping_;
    }

    var sidxKey = generateSidxKey(playlist.sidx);
    sidxMapping[sidxKey] = {
      sidxInfo: playlist.sidx
    };
    this.request = requestSidx_(this, playlist.sidx, playlist, this.vhs_.xhr, {
      handleManifestRedirects: this.handleManifestRedirects
    }, this.sidxRequestFinished_(playlist, oldMaster, startingState, function (newMaster, sidx) {
      if (!newMaster || !sidx) {
        throw new Error('failed to request sidx');
      } // update loader's sidxMapping with parsed sidx box


      sidxMapping[sidxKey].sidx = sidx; // everything is ready just continue to haveMetadata

      _this3.haveMetadata({
        startingState: startingState,
        playlist: newMaster.playlists[playlist.id]
      });
    }));
  };

  _proto.haveMetadata = function haveMetadata(_ref2) {
    var startingState = _ref2.startingState,
        playlist = _ref2.playlist;
    this.state = 'HAVE_METADATA';
    this.loadedPlaylists_[playlist.id] = playlist;
    this.mediaRequest_ = null; // This will trigger loadedplaylist

    this.refreshMedia_(playlist.id); // fire loadedmetadata the first time a media playlist is loaded
    // to resolve setup of media groups

    if (startingState === 'HAVE_MASTER') {
      this.trigger('loadedmetadata');
    } else {
      // trigger media change if the active media has been updated
      this.trigger('mediachange');
    }
  };

  _proto.pause = function pause() {
    this.stopRequest();
    window_1.clearTimeout(this.mediaUpdateTimeout);
    window_1.clearTimeout(this.minimumUpdatePeriodTimeout_);

    if (this.state === 'HAVE_NOTHING') {
      // If we pause the loader before any data has been retrieved, its as if we never
      // started, so reset to an unstarted state.
      this.started = false;
    }
  };

  _proto.load = function load(isFinalRendition) {
    var _this4 = this;

    window_1.clearTimeout(this.mediaUpdateTimeout);
    window_1.clearTimeout(this.minimumUpdatePeriodTimeout_);
    var media = this.media();

    if (isFinalRendition) {
      var delay = media ? media.targetDuration / 2 * 1000 : 5 * 1000;
      this.mediaUpdateTimeout = window_1.setTimeout(function () {
        return _this4.load();
      }, delay);
      return;
    } // because the playlists are internal to the manifest, load should either load the
    // main manifest, or do nothing but trigger an event


    if (!this.started) {
      this.start();
      return;
    }

    if (media && !media.endList) {
      this.trigger('mediaupdatetimeout');
    } else {
      this.trigger('loadedplaylist');
    }
  };

  _proto.start = function start() {
    var _this5 = this;

    this.started = true; // We don't need to request the master manifest again
    // Call this asynchronously to match the xhr request behavior below

    if (this.masterPlaylistLoader_) {
      this.mediaRequest_ = window_1.setTimeout(this.haveMaster_.bind(this), 0);
      return;
    } // request the specified URL


    this.request = this.vhs_.xhr({
      uri: this.srcUrl,
      withCredentials: this.withCredentials
    }, function (error, req) {
      // disposed
      if (!_this5.request) {
        return;
      } // clear the loader's request reference


      _this5.request = null;

      if (error) {
        _this5.error = {
          status: req.status,
          message: 'DASH playlist request error at URL: ' + _this5.srcUrl,
          responseText: req.responseText,
          // MEDIA_ERR_NETWORK
          code: 2
        };

        if (_this5.state === 'HAVE_NOTHING') {
          _this5.started = false;
        }

        return _this5.trigger('error');
      }

      _this5.masterXml_ = req.responseText;

      if (req.responseHeaders && req.responseHeaders.date) {
        _this5.masterLoaded_ = Date.parse(req.responseHeaders.date);
      } else {
        _this5.masterLoaded_ = Date.now();
      }

      _this5.srcUrl = resolveManifestRedirect(_this5.handleManifestRedirects, _this5.srcUrl, req);

      _this5.syncClientServerClock_(_this5.onClientServerClockSync_.bind(_this5));
    });
  }
  /**
   * Parses the master xml for UTCTiming node to sync the client clock to the server
   * clock. If the UTCTiming node requires a HEAD or GET request, that request is made.
   *
   * @param {Function} done
   *        Function to call when clock sync has completed
   */
  ;

  _proto.syncClientServerClock_ = function syncClientServerClock_(done) {
    var _this6 = this;

    var utcTiming = parseUTCTiming(this.masterXml_); // No UTCTiming element found in the mpd. Use Date header from mpd request as the
    // server clock

    if (utcTiming === null) {
      this.clientOffset_ = this.masterLoaded_ - Date.now();
      return done();
    }

    if (utcTiming.method === 'DIRECT') {
      this.clientOffset_ = utcTiming.value - Date.now();
      return done();
    }

    this.request = this.vhs_.xhr({
      uri: resolveUrl$2(this.srcUrl, utcTiming.value),
      method: utcTiming.method,
      withCredentials: this.withCredentials
    }, function (error, req) {
      // disposed
      if (!_this6.request) {
        return;
      }

      if (error) {
        // sync request failed, fall back to using date header from mpd
        // TODO: log warning
        _this6.clientOffset_ = _this6.masterLoaded_ - Date.now();
        return done();
      }

      var serverTime;

      if (utcTiming.method === 'HEAD') {
        if (!req.responseHeaders || !req.responseHeaders.date) {
          // expected date header not preset, fall back to using date header from mpd
          // TODO: log warning
          serverTime = _this6.masterLoaded_;
        } else {
          serverTime = Date.parse(req.responseHeaders.date);
        }
      } else {
        serverTime = Date.parse(req.responseText);
      }

      _this6.clientOffset_ = serverTime - Date.now();
      done();
    });
  };

  _proto.haveMaster_ = function haveMaster_() {
    this.state = 'HAVE_MASTER'; // clear media request

    this.mediaRequest_ = null;

    if (!this.masterPlaylistLoader_) {
      this.updateMainManifest_(parseMasterXml({
        masterXml: this.masterXml_,
        srcUrl: this.srcUrl,
        clientOffset: this.clientOffset_,
        sidxMapping: this.sidxMapping_
      })); // We have the master playlist at this point, so
      // trigger this to allow MasterPlaylistController
      // to make an initial playlist selection

      this.trigger('loadedplaylist');
    } else if (!this.media_) {
      // no media playlist was specifically selected so select
      // the one the child playlist loader was created with
      this.media(this.childPlaylist_);
    }
  };

  _proto.updateMinimumUpdatePeriodTimeout_ = function updateMinimumUpdatePeriodTimeout_() {
    var _this7 = this;

    // Clear existing timeout
    window_1.clearTimeout(this.minimumUpdatePeriodTimeout_);

    var createMUPTimeout = function createMUPTimeout(mup) {
      _this7.minimumUpdatePeriodTimeout_ = window_1.setTimeout(function () {
        _this7.trigger('minimumUpdatePeriod');
      }, mup);
    };

    var minimumUpdatePeriod = this.master && this.master.minimumUpdatePeriod;

    if (minimumUpdatePeriod > 0) {
      createMUPTimeout(minimumUpdatePeriod); // If the minimumUpdatePeriod has a value of 0, that indicates that the current
      // MPD has no future validity, so a new one will need to be acquired when new
      // media segments are to be made available. Thus, we use the target duration
      // in this case
    } else if (minimumUpdatePeriod === 0) {
      // If we haven't yet selected a playlist, wait until then so we know the
      // target duration
      if (!this.media()) {
        this.one('loadedplaylist', function () {
          createMUPTimeout(_this7.media().targetDuration * 1000);
        });
      } else {
        createMUPTimeout(this.media().targetDuration * 1000);
      }
    }
  }
  /**
   * Handler for after client/server clock synchronization has happened. Sets up
   * xml refresh timer if specificed by the manifest.
   */
  ;

  _proto.onClientServerClockSync_ = function onClientServerClockSync_() {
    this.haveMaster_();

    if (!this.hasPendingRequest() && !this.media_) {
      this.media(this.master.playlists[0]);
    }

    this.updateMinimumUpdatePeriodTimeout_();
  }
  /**
   * Given a new manifest, update our pointer to it and update the srcUrl based on the location elements of the manifest, if they exist.
   *
   * @param {Object} updatedManifest the manifest to update to
   */
  ;

  _proto.updateMainManifest_ = function updateMainManifest_(updatedManifest) {
    this.master = updatedManifest; // if locations isn't set or is an empty array, exit early

    if (!this.master.locations || !this.master.locations.length) {
      return;
    }

    var location = this.master.locations[0];

    if (location !== this.srcUrl) {
      this.srcUrl = location;
    }
  }
  /**
   * Sends request to refresh the master xml and updates the parsed master manifest
   * TODO: Does the client offset need to be recalculated when the xml is refreshed?
   */
  ;

  _proto.refreshXml_ = function refreshXml_() {
    var _this8 = this;

    // The srcUrl here *may* need to pass through handleManifestsRedirects when
    // sidx is implemented
    this.request = this.vhs_.xhr({
      uri: this.srcUrl,
      withCredentials: this.withCredentials
    }, function (error, req) {
      // disposed
      if (!_this8.request) {
        return;
      } // clear the loader's request reference


      _this8.request = null;

      if (error) {
        _this8.error = {
          status: req.status,
          message: 'DASH playlist request error at URL: ' + _this8.srcUrl,
          responseText: req.responseText,
          // MEDIA_ERR_NETWORK
          code: 2
        };

        if (_this8.state === 'HAVE_NOTHING') {
          _this8.started = false;
        }

        return _this8.trigger('error');
      }

      _this8.masterXml_ = req.responseText; // This will filter out updated sidx info from the mapping

      _this8.sidxMapping_ = filterChangedSidxMappings(_this8.masterXml_, _this8.srcUrl, _this8.clientOffset_, _this8.sidxMapping_);
      var master = parseMasterXml({
        masterXml: _this8.masterXml_,
        srcUrl: _this8.srcUrl,
        clientOffset: _this8.clientOffset_,
        sidxMapping: _this8.sidxMapping_
      });
      var updatedMaster = updateMaster$1(_this8.master, master);

      var currentSidxInfo = _this8.media().sidx;

      if (updatedMaster) {
        if (currentSidxInfo) {
          var sidxKey = generateSidxKey(currentSidxInfo); // the sidx was updated, so the previous mapping was removed

          if (!_this8.sidxMapping_[sidxKey]) {
            var playlist = _this8.media();

            _this8.request = requestSidx_(_this8, playlist.sidx, playlist, _this8.vhs_.xhr, {
              handleManifestRedirects: _this8.handleManifestRedirects
            }, _this8.sidxRequestFinished_(playlist, master, _this8.state, function (newMaster, sidx) {
              if (!newMaster || !sidx) {
                throw new Error('failed to request sidx on minimumUpdatePeriod');
              } // update loader's sidxMapping with parsed sidx box


              _this8.sidxMapping_[sidxKey].sidx = sidx;

              _this8.updateMinimumUpdatePeriodTimeout_(); // TODO: do we need to reload the current playlist?


              _this8.refreshMedia_(_this8.media().id);

              return;
            }));
          }
        } else {
          _this8.updateMainManifest_(updatedMaster);

          if (_this8.media_) {
            _this8.media_ = _this8.master.playlists[_this8.media_.id];
          }
        }
      }

      _this8.updateMinimumUpdatePeriodTimeout_();
    });
  }
  /**
   * Refreshes the media playlist by re-parsing the master xml and updating playlist
   * references. If this is an alternate loader, the updated parsed manifest is retrieved
   * from the master loader.
   */
  ;

  _proto.refreshMedia_ = function refreshMedia_(mediaID) {
    var _this9 = this;

    if (!mediaID) {
      throw new Error('refreshMedia_ must take a media id');
    }

    var oldMaster;
    var newMaster;

    if (this.masterPlaylistLoader_) {
      oldMaster = this.masterPlaylistLoader_.master;
      newMaster = parseMasterXml({
        masterXml: this.masterPlaylistLoader_.masterXml_,
        srcUrl: this.masterPlaylistLoader_.srcUrl,
        clientOffset: this.masterPlaylistLoader_.clientOffset_,
        sidxMapping: this.masterPlaylistLoader_.sidxMapping_
      });
    } else {
      oldMaster = this.master;
      newMaster = parseMasterXml({
        masterXml: this.masterXml_,
        srcUrl: this.srcUrl,
        clientOffset: this.clientOffset_,
        sidxMapping: this.sidxMapping_
      });
    }

    var updatedMaster = updateMaster$1(oldMaster, newMaster);

    if (updatedMaster) {
      if (this.masterPlaylistLoader_) {
        this.masterPlaylistLoader_.master = updatedMaster;
      } else {
        this.master = updatedMaster;
      }

      this.media_ = updatedMaster.playlists[mediaID];
    } else {
      this.media_ = oldMaster.playlists[mediaID];
      this.trigger('playlistunchanged');
    }

    if (!this.media().endList) {
      this.mediaUpdateTimeout = window_1.setTimeout(function () {
        _this9.trigger('mediaupdatetimeout');
      }, refreshDelay(this.media(), !!updatedMaster));
    }

    this.trigger('loadedplaylist');
  };

  return DashPlaylistLoader;
}(EventTarget$1);

var Config = {
  GOAL_BUFFER_LENGTH: 30,
  MAX_GOAL_BUFFER_LENGTH: 60,
  BACK_BUFFER_LENGTH: 30,
  GOAL_BUFFER_LENGTH_RATE: 1,
  // 0.5 MB/s
  INITIAL_BANDWIDTH: 4194304,
  // A fudge factor to apply to advertised playlist bitrates to account for
  // temporary flucations in client bandwidth
  BANDWIDTH_VARIANCE: 1.2,
  // How much of the buffer must be filled before we consider upswitching
  BUFFER_LOW_WATER_LINE: 0,
  MAX_BUFFER_LOW_WATER_LINE: 30,
  BUFFER_LOW_WATER_LINE_RATE: 1
};

var stringToArrayBuffer = function stringToArrayBuffer(string) {
  var view = new Uint8Array(new ArrayBuffer(string.length));

  for (var i = 0; i < string.length; i++) {
    view[i] = string.charCodeAt(i);
  }

  return view.buffer;
};

var transmuxQueue = [];
var currentTransmux;
var handleData_ = function handleData_(event, transmuxedData, callback) {
  var _event$data$segment = event.data.segment,
      type = _event$data$segment.type,
      initSegment = _event$data$segment.initSegment,
      captions = _event$data$segment.captions,
      captionStreams = _event$data$segment.captionStreams,
      metadata = _event$data$segment.metadata,
      videoFrameDtsTime = _event$data$segment.videoFrameDtsTime,
      videoFramePtsTime = _event$data$segment.videoFramePtsTime;
  transmuxedData.buffer.push({
    captions: captions,
    captionStreams: captionStreams,
    metadata: metadata
  }); // right now, boxes will come back from partial transmuxer, data from full

  var boxes = event.data.segment.boxes || {
    data: event.data.segment.data
  };
  var result = {
    type: type,
    // cast ArrayBuffer to TypedArray
    data: new Uint8Array(boxes.data, boxes.data.byteOffset, boxes.data.byteLength),
    initSegment: new Uint8Array(initSegment.data, initSegment.byteOffset, initSegment.byteLength)
  };

  if (typeof videoFrameDtsTime !== 'undefined') {
    result.videoFrameDtsTime = videoFrameDtsTime;
  }

  if (typeof videoFramePtsTime !== 'undefined') {
    result.videoFramePtsTime = videoFramePtsTime;
  }

  callback(result);
};
var handleDone_ = function handleDone_(_ref) {
  var transmuxedData = _ref.transmuxedData,
      callback = _ref.callback;
  // Previously we only returned data on data events,
  // not on done events. Clear out the buffer to keep that consistent.
  transmuxedData.buffer = []; // all buffers should have been flushed from the muxer, so start processing anything we
  // have received

  callback(transmuxedData);
};
var handleGopInfo_ = function handleGopInfo_(event, transmuxedData) {
  transmuxedData.gopInfo = event.data.gopInfo;
};
var processTransmux = function processTransmux(_ref2) {
  var transmuxer = _ref2.transmuxer,
      bytes = _ref2.bytes,
      audioAppendStart = _ref2.audioAppendStart,
      gopsToAlignWith = _ref2.gopsToAlignWith,
      isPartial = _ref2.isPartial,
      remux = _ref2.remux,
      onData = _ref2.onData,
      onTrackInfo = _ref2.onTrackInfo,
      onAudioTimingInfo = _ref2.onAudioTimingInfo,
      onVideoTimingInfo = _ref2.onVideoTimingInfo,
      onVideoSegmentTimingInfo = _ref2.onVideoSegmentTimingInfo,
      onId3 = _ref2.onId3,
      onCaptions = _ref2.onCaptions,
      onDone = _ref2.onDone;
  var transmuxedData = {
    isPartial: isPartial,
    buffer: []
  };

  var handleMessage = function handleMessage(event) {
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
    } // wait for the transmuxed event since we may have audio and video


    if (event.data.type !== 'transmuxed') {
      return;
    }

    transmuxer.onmessage = null;
    handleDone_({
      transmuxedData: transmuxedData,
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
  } // allow empty arrays to be passed to clear out GOPs


  if (Array.isArray(gopsToAlignWith)) {
    transmuxer.postMessage({
      action: 'alignGopsWith',
      gopsToAlignWith: gopsToAlignWith
    });
  }

  if (typeof remux !== 'undefined') {
    transmuxer.postMessage({
      action: 'setRemux',
      remux: remux
    });
  }

  if (bytes.byteLength) {
    var buffer = bytes instanceof ArrayBuffer ? bytes : bytes.buffer;
    var byteOffset = bytes instanceof ArrayBuffer ? 0 : bytes.byteOffset;
    transmuxer.postMessage({
      action: 'push',
      // Send the typed-array of data as an ArrayBuffer so that
      // it can be sent as a "Transferable" and avoid the costly
      // memory copy
      data: buffer,
      // To recreate the original typed-array, we need information
      // about what portion of the ArrayBuffer it was a view into
      byteOffset: byteOffset,
      byteLength: bytes.byteLength
    }, [buffer]);
  } // even if we didn't push any bytes, we have to make sure we flush in case we reached
  // the end of the segment


  transmuxer.postMessage({
    action: isPartial ? 'partialFlush' : 'flush'
  });
};
var dequeue = function dequeue() {
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
var processAction = function processAction(transmuxer, action) {
  transmuxer.postMessage({
    action: action
  });
  dequeue();
};
var enqueueAction = function enqueueAction(action, transmuxer) {
  if (!currentTransmux) {
    currentTransmux = action;
    processAction(transmuxer, action);
    return;
  }

  transmuxQueue.push(processAction.bind(null, transmuxer, action));
};
var reset = function reset(transmuxer) {
  enqueueAction('reset', transmuxer);
};
var endTimeline = function endTimeline(transmuxer) {
  enqueueAction('endTimeline', transmuxer);
};
var transmux = function transmux(options) {
  if (!currentTransmux) {
    currentTransmux = options;
    processTransmux(options);
    return;
  }

  transmuxQueue.push(options);
};
var dispose = function dispose() {
  // clear out module-level references
  currentTransmux = null;
  transmuxQueue.length = 0;
};
var segmentTransmuxer = {
  reset: reset,
  dispose: dispose,
  endTimeline: endTimeline,
  transmux: transmux
};

/**
 * Probe an mpeg2-ts segment to determine the start time of the segment in it's
 * internal "media time," as well as whether it contains video and/or audio.
 *
 * @private
 * @param {Uint8Array} bytes - segment bytes
 * @return {Object} The start time of the current segment in "media time" as well as
 *                  whether it contains video and/or audio
 */

var probeTsSegment = function probeTsSegment(bytes, baseStartTime) {
  var timeInfo = tsInspector.inspect(bytes, baseStartTime * clock.ONE_SECOND_IN_TS);

  if (!timeInfo) {
    return null;
  }

  var result = {
    // each type's time info comes back as an array of 2 times, start and end
    hasVideo: timeInfo.video && timeInfo.video.length === 2 || false,
    hasAudio: timeInfo.audio && timeInfo.audio.length === 2 || false
  };

  if (result.hasVideo) {
    result.videoStart = timeInfo.video[0].ptsTime;
  }

  if (result.hasAudio) {
    result.audioStart = timeInfo.audio[0].ptsTime;
  }

  return result;
};
/**
 * Combine all segments into a single Uint8Array
 *
 * @param {Object} segmentObj
 * @return {Uint8Array} concatenated bytes
 * @private
 */

var concatSegments = function concatSegments(segmentObj) {
  var offset = 0;
  var tempBuffer;

  if (segmentObj.bytes) {
    tempBuffer = new Uint8Array(segmentObj.bytes); // combine the individual segments into one large typed-array

    segmentObj.segments.forEach(function (segment) {
      tempBuffer.set(segment, offset);
      offset += segment.byteLength;
    });
  }

  return tempBuffer;
};

var REQUEST_ERRORS = {
  FAILURE: 2,
  TIMEOUT: -101,
  ABORTED: -102
};
/**
 * Abort all requests
 *
 * @param {Object} activeXhrs - an object that tracks all XHR requests
 */

var abortAll = function abortAll(activeXhrs) {
  activeXhrs.forEach(function (xhr) {
    xhr.abort();
  });
};
/**
 * Gather important bandwidth stats once a request has completed
 *
 * @param {Object} request - the XHR request from which to gather stats
 */


var getRequestStats = function getRequestStats(request) {
  return {
    bandwidth: request.bandwidth,
    bytesReceived: request.bytesReceived || 0,
    roundTripTime: request.roundTripTime || 0
  };
};
/**
 * If possible gather bandwidth stats as a request is in
 * progress
 *
 * @param {Event} progressEvent - an event object from an XHR's progress event
 */


var getProgressStats = function getProgressStats(progressEvent) {
  var request = progressEvent.target;
  var roundTripTime = Date.now() - request.requestTime;
  var stats = {
    bandwidth: Infinity,
    bytesReceived: 0,
    roundTripTime: roundTripTime || 0
  };
  stats.bytesReceived = progressEvent.loaded; // This can result in Infinity if stats.roundTripTime is 0 but that is ok
  // because we should only use bandwidth stats on progress to determine when
  // abort a request early due to insufficient bandwidth

  stats.bandwidth = Math.floor(stats.bytesReceived / stats.roundTripTime * 8 * 1000);
  return stats;
};
/**
 * Handle all error conditions in one place and return an object
 * with all the information
 *
 * @param {Error|null} error - if non-null signals an error occured with the XHR
 * @param {Object} request -  the XHR request that possibly generated the error
 */


var handleErrors = function handleErrors(error, request) {
  if (request.timedout) {
    return {
      status: request.status,
      message: 'HLS request timed-out at URL: ' + request.uri,
      code: REQUEST_ERRORS.TIMEOUT,
      xhr: request
    };
  }

  if (request.aborted) {
    return {
      status: request.status,
      message: 'HLS request aborted at URL: ' + request.uri,
      code: REQUEST_ERRORS.ABORTED,
      xhr: request
    };
  }

  if (error) {
    return {
      status: request.status,
      message: 'HLS request errored at URL: ' + request.uri,
      code: REQUEST_ERRORS.FAILURE,
      xhr: request
    };
  }

  return null;
};
/**
 * Handle responses for key data and convert the key data to the correct format
 * for the decryption step later
 *
 * @param {Object} segment - a simplified copy of the segmentInfo object
 *                           from SegmentLoader
 * @param {Function} finishProcessingFn - a callback to execute to continue processing
 *                                        this request
 */


var handleKeyResponse = function handleKeyResponse(segment, finishProcessingFn) {
  return function (error, request) {
    var response = request.response;
    var errorObj = handleErrors(error, request);

    if (errorObj) {
      return finishProcessingFn(errorObj, segment);
    }

    if (response.byteLength !== 16) {
      return finishProcessingFn({
        status: request.status,
        message: 'Invalid HLS key at URL: ' + request.uri,
        code: REQUEST_ERRORS.FAILURE,
        xhr: request
      }, segment);
    }

    var view = new DataView(response);
    segment.key.bytes = new Uint32Array([view.getUint32(0), view.getUint32(4), view.getUint32(8), view.getUint32(12)]);
    return finishProcessingFn(null, segment);
  };
};
/**
 * Handle init-segment responses
 *
 * @param {Object} segment - a simplified copy of the segmentInfo object
 *                           from SegmentLoader
 * @param {Function} finishProcessingFn - a callback to execute to continue processing
 *                                        this request
 */


var handleInitSegmentResponse = function handleInitSegmentResponse(_ref) {
  var segment = _ref.segment,
      finishProcessingFn = _ref.finishProcessingFn;
  return function (error, request) {
    var response = request.response;
    var errorObj = handleErrors(error, request);

    if (errorObj) {
      return finishProcessingFn(errorObj, segment);
    } // stop processing if received empty content


    if (response.byteLength === 0) {
      return finishProcessingFn({
        status: request.status,
        message: 'Empty HLS segment content at URL: ' + request.uri,
        code: REQUEST_ERRORS.FAILURE,
        xhr: request
      }, segment);
    }

    segment.map.bytes = new Uint8Array(request.response);
    var type = containers.detectContainerForBytes(segment.map.bytes); // TODO: We should also handle ts init segments here, but we
    // only know how to parse mp4 init segments at the moment

    if (type !== 'mp4') {
      return finishProcessingFn({
        status: request.status,
        message: "Found unsupported " + (type || 'unknown') + " container for initialization segment at URL: " + request.uri,
        code: REQUEST_ERRORS.FAILURE,
        internal: true,
        xhr: request
      }, segment);
    }

    var tracks = probe$2.tracks(segment.map.bytes);
    tracks.forEach(function (track) {
      segment.map.tracks = segment.map.tracks || {}; // only support one track of each type for now

      if (segment.map.tracks[track.type]) {
        return;
      }

      segment.map.tracks[track.type] = track;

      if (track.id && track.timescale) {
        segment.map.timescales = segment.map.timescales || {};
        segment.map.timescales[track.id] = track.timescale;
      }
    });
    return finishProcessingFn(null, segment);
  };
};
/**
 * Response handler for segment-requests being sure to set the correct
 * property depending on whether the segment is encryped or not
 * Also records and keeps track of stats that are used for ABR purposes
 *
 * @param {Object} segment - a simplified copy of the segmentInfo object
 *                           from SegmentLoader
 * @param {Function} finishProcessingFn - a callback to execute to continue processing
 *                                        this request
 */


var handleSegmentResponse = function handleSegmentResponse(_ref2) {
  var segment = _ref2.segment,
      finishProcessingFn = _ref2.finishProcessingFn,
      responseType = _ref2.responseType;
  return function (error, request) {
    var response = request.response;
    var errorObj = handleErrors(error, request);

    if (errorObj) {
      return finishProcessingFn(errorObj, segment);
    }

    var newBytes = // although responseText "should" exist, this guard serves to prevent an error being
    // thrown for two primary cases:
    // 1. the mime type override stops working, or is not implemented for a specific
    //    browser
    // 2. when using mock XHR libraries like sinon that do not allow the override behavior
    responseType === 'arraybuffer' || !request.responseText ? request.response : stringToArrayBuffer(request.responseText.substring(segment.lastReachedChar || 0)); // stop processing if received empty content

    if (response.byteLength === 0) {
      return finishProcessingFn({
        status: request.status,
        message: 'Empty HLS segment content at URL: ' + request.uri,
        code: REQUEST_ERRORS.FAILURE,
        xhr: request
      }, segment);
    }

    segment.stats = getRequestStats(request);

    if (segment.key) {
      segment.encryptedBytes = new Uint8Array(newBytes);
    } else {
      segment.bytes = new Uint8Array(newBytes);
    }

    return finishProcessingFn(null, segment);
  };
};

var transmuxAndNotify = function transmuxAndNotify(_ref3) {
  var segment = _ref3.segment,
      bytes = _ref3.bytes,
      isPartial = _ref3.isPartial,
      trackInfoFn = _ref3.trackInfoFn,
      timingInfoFn = _ref3.timingInfoFn,
      videoSegmentTimingInfoFn = _ref3.videoSegmentTimingInfoFn,
      id3Fn = _ref3.id3Fn,
      captionsFn = _ref3.captionsFn,
      dataFn = _ref3.dataFn,
      doneFn = _ref3.doneFn;
  var fmp4Tracks = segment.map && segment.map.tracks || {};
  var isMuxed = Boolean(fmp4Tracks.audio && fmp4Tracks.video); // Keep references to each function so we can null them out after we're done with them.
  // One reason for this is that in the case of full segments, we want to trust start
  // times from the probe, rather than the transmuxer.

  var audioStartFn = timingInfoFn.bind(null, segment, 'audio', 'start');
  var audioEndFn = timingInfoFn.bind(null, segment, 'audio', 'end');
  var videoStartFn = timingInfoFn.bind(null, segment, 'video', 'start');
  var videoEndFn = timingInfoFn.bind(null, segment, 'video', 'end'); // Check to see if we are appending a full segment.

  if (!isPartial && !segment.lastReachedChar) {
    // In the full segment transmuxer, we don't yet have the ability to extract a "proper"
    // start time. Meaning cached frame data may corrupt our notion of where this segment
    // really starts. To get around this, full segment appends should probe for the info
    // needed.
    var probeResult = probeTsSegment(bytes, segment.baseStartTime);

    if (probeResult) {
      trackInfoFn(segment, {
        hasAudio: probeResult.hasAudio,
        hasVideo: probeResult.hasVideo,
        isMuxed: isMuxed
      });
      trackInfoFn = null;

      if (probeResult.hasAudio && !isMuxed) {
        audioStartFn(probeResult.audioStart);
      }

      if (probeResult.hasVideo) {
        videoStartFn(probeResult.videoStart);
      }

      audioStartFn = null;
      videoStartFn = null;
    }
  }

  transmux({
    bytes: bytes,
    transmuxer: segment.transmuxer,
    audioAppendStart: segment.audioAppendStart,
    gopsToAlignWith: segment.gopsToAlignWith,
    isPartial: isPartial,
    remux: isMuxed,
    onData: function onData(result) {
      result.type = result.type === 'combined' ? 'video' : result.type;
      dataFn(segment, result);
    },
    onTrackInfo: function onTrackInfo(trackInfo) {
      if (trackInfoFn) {
        if (isMuxed) {
          trackInfo.isMuxed = true;
        }

        trackInfoFn(segment, trackInfo);
      }
    },
    onAudioTimingInfo: function onAudioTimingInfo(audioTimingInfo) {
      // we only want the first start value we encounter
      if (audioStartFn && typeof audioTimingInfo.start !== 'undefined') {
        audioStartFn(audioTimingInfo.start);
        audioStartFn = null;
      } // we want to continually update the end time


      if (audioEndFn && typeof audioTimingInfo.end !== 'undefined') {
        audioEndFn(audioTimingInfo.end);
      }
    },
    onVideoTimingInfo: function onVideoTimingInfo(videoTimingInfo) {
      // we only want the first start value we encounter
      if (videoStartFn && typeof videoTimingInfo.start !== 'undefined') {
        videoStartFn(videoTimingInfo.start);
        videoStartFn = null;
      } // we want to continually update the end time


      if (videoEndFn && typeof videoTimingInfo.end !== 'undefined') {
        videoEndFn(videoTimingInfo.end);
      }
    },
    onVideoSegmentTimingInfo: function onVideoSegmentTimingInfo(videoSegmentTimingInfo) {
      videoSegmentTimingInfoFn(videoSegmentTimingInfo);
    },
    onId3: function onId3(id3Frames, dispatchType) {
      id3Fn(segment, id3Frames, dispatchType);
    },
    onCaptions: function onCaptions(captions) {
      captionsFn(segment, [captions]);
    },
    onDone: function onDone(result) {
      // To handle partial appends, there won't be a done function passed in (since
      // there's still, potentially, more segment to process), so there's nothing to do.
      if (!doneFn || isPartial) {
        return;
      }

      result.type = result.type === 'combined' ? 'video' : result.type;
      doneFn(null, segment, result);
    }
  });
};

var handleSegmentBytes = function handleSegmentBytes(_ref4) {
  var segment = _ref4.segment,
      bytes = _ref4.bytes,
      isPartial = _ref4.isPartial,
      trackInfoFn = _ref4.trackInfoFn,
      timingInfoFn = _ref4.timingInfoFn,
      videoSegmentTimingInfoFn = _ref4.videoSegmentTimingInfoFn,
      id3Fn = _ref4.id3Fn,
      captionsFn = _ref4.captionsFn,
      dataFn = _ref4.dataFn,
      doneFn = _ref4.doneFn;
  var bytesAsUint8Array = new Uint8Array(bytes); // TODO:
  // We should have a handler that fetches the number of bytes required
  // to check if something is fmp4. This will allow us to save bandwidth
  // because we can only blacklist a playlist and abort requests
  // by codec after trackinfo triggers.

  if (containers.isLikelyFmp4MediaSegment(bytesAsUint8Array)) {
    segment.isFmp4 = true;
    var tracks = segment.map.tracks;
    var trackInfo = {
      isFmp4: true,
      hasVideo: !!tracks.video,
      hasAudio: !!tracks.audio
    }; // if we have a audio track, with a codec that is not set to
    // encrypted audio

    if (tracks.audio && tracks.audio.codec && tracks.audio.codec !== 'enca') {
      trackInfo.audioCodec = tracks.audio.codec;
    } // if we have a video track, with a codec that is not set to
    // encrypted video


    if (tracks.video && tracks.video.codec && tracks.video.codec !== 'encv') {
      trackInfo.videoCodec = tracks.video.codec;
    }

    if (tracks.video && tracks.audio) {
      trackInfo.isMuxed = true;
    } // since we don't support appending fmp4 data on progress, we know we have the full
    // segment here


    trackInfoFn(segment, trackInfo); // The probe doesn't provide the segment end time, so only callback with the start
    // time. The end time can be roughly calculated by the receiver using the duration.
    //
    // Note that the start time returned by the probe reflects the baseMediaDecodeTime, as
    // that is the true start of the segment (where the playback engine should begin
    // decoding).

    var timingInfo = probe$2.startTime(segment.map.timescales, bytesAsUint8Array);

    if (trackInfo.hasAudio && !trackInfo.isMuxed) {
      timingInfoFn(segment, 'audio', 'start', timingInfo);
    }

    if (trackInfo.hasVideo) {
      timingInfoFn(segment, 'video', 'start', timingInfo);
    }

    var finishLoading = function finishLoading(captions) {
      // if the track still has audio at this point it is only possible
      // for it to be audio only. See `tracks.video && tracks.audio` if statement
      // above.
      // we make sure to use segment.bytes here as that
      dataFn(segment, {
        data: bytes,
        type: trackInfo.hasAudio && !trackInfo.isMuxed ? 'audio' : 'video'
      });

      if (captions && captions.length) {
        captionsFn(segment, captions);
      }

      doneFn(null, segment, {});
    }; // Run through the CaptionParser in case there are captions.
    // Initialize CaptionParser if it hasn't been yet


    if (!tracks.video || !bytes.byteLength || !segment.transmuxer) {
      finishLoading();
      return;
    }

    var buffer = bytes instanceof ArrayBuffer ? bytes : bytes.buffer;
    var byteOffset = bytes instanceof ArrayBuffer ? 0 : bytes.byteOffset;

    var listenForCaptions = function listenForCaptions(event) {
      if (event.data.action !== 'mp4Captions') {
        return;
      }

      segment.transmuxer.removeEventListener('message', listenForCaptions);
      var data = event.data.data; // transfer ownership of bytes back to us.

      segment.bytes = bytes = new Uint8Array(data, data.byteOffset || 0, data.byteLength);
      finishLoading(event.data.captions);
    };

    segment.transmuxer.addEventListener('message', listenForCaptions); // transfer ownership of bytes to worker.

    segment.transmuxer.postMessage({
      action: 'pushMp4Captions',
      timescales: segment.map.timescales,
      trackIds: [tracks.video.id],
      data: buffer,
      byteOffset: byteOffset,
      byteLength: bytes.byteLength
    }, [buffer]);
    return;
  } // VTT or other segments that don't need processing


  if (!segment.transmuxer) {
    doneFn(null, segment, {});
    return;
  }

  if (typeof segment.container === 'undefined') {
    segment.container = containers.detectContainerForBytes(bytesAsUint8Array);
  }

  if (segment.container !== 'ts' && segment.container !== 'aac') {
    trackInfoFn(segment, {
      hasAudio: false,
      hasVideo: false
    });
    doneFn(null, segment, {});
    return;
  } // ts or aac


  transmuxAndNotify({
    segment: segment,
    bytes: bytes,
    isPartial: isPartial,
    trackInfoFn: trackInfoFn,
    timingInfoFn: timingInfoFn,
    videoSegmentTimingInfoFn: videoSegmentTimingInfoFn,
    id3Fn: id3Fn,
    captionsFn: captionsFn,
    dataFn: dataFn,
    doneFn: doneFn
  });
};
/**
 * Decrypt the segment via the decryption web worker
 *
 * @param {WebWorker} decryptionWorker - a WebWorker interface to AES-128 decryption
 *                                       routines
 * @param {Object} segment - a simplified copy of the segmentInfo object
 *                           from SegmentLoader
 * @param {Function} trackInfoFn - a callback that receives track info
 * @param {Function} dataFn - a callback that is executed when segment bytes are available
 *                            and ready to use
 * @param {Function} doneFn - a callback that is executed after decryption has completed
 */


var decryptSegment = function decryptSegment(_ref5) {
  var decryptionWorker = _ref5.decryptionWorker,
      segment = _ref5.segment,
      trackInfoFn = _ref5.trackInfoFn,
      timingInfoFn = _ref5.timingInfoFn,
      videoSegmentTimingInfoFn = _ref5.videoSegmentTimingInfoFn,
      id3Fn = _ref5.id3Fn,
      captionsFn = _ref5.captionsFn,
      dataFn = _ref5.dataFn,
      doneFn = _ref5.doneFn;

  var decryptionHandler = function decryptionHandler(event) {
    if (event.data.source === segment.requestId) {
      decryptionWorker.removeEventListener('message', decryptionHandler);
      var decrypted = event.data.decrypted;
      segment.bytes = new Uint8Array(decrypted.bytes, decrypted.byteOffset, decrypted.byteLength);
      handleSegmentBytes({
        segment: segment,
        bytes: segment.bytes,
        isPartial: false,
        trackInfoFn: trackInfoFn,
        timingInfoFn: timingInfoFn,
        videoSegmentTimingInfoFn: videoSegmentTimingInfoFn,
        id3Fn: id3Fn,
        captionsFn: captionsFn,
        dataFn: dataFn,
        doneFn: doneFn
      });
    }
  };

  decryptionWorker.addEventListener('message', decryptionHandler);
  var keyBytes;

  if (segment.key.bytes.slice) {
    keyBytes = segment.key.bytes.slice();
  } else {
    keyBytes = new Uint32Array(Array.prototype.slice.call(segment.key.bytes));
  } // this is an encrypted segment
  // incrementally decrypt the segment


  decryptionWorker.postMessage(createTransferableMessage({
    source: segment.requestId,
    encrypted: segment.encryptedBytes,
    key: keyBytes,
    iv: segment.key.iv
  }), [segment.encryptedBytes.buffer, keyBytes.buffer]);
};
/**
 * This function waits for all XHRs to finish (with either success or failure)
 * before continueing processing via it's callback. The function gathers errors
 * from each request into a single errors array so that the error status for
 * each request can be examined later.
 *
 * @param {Object} activeXhrs - an object that tracks all XHR requests
 * @param {WebWorker} decryptionWorker - a WebWorker interface to AES-128 decryption
 *                                       routines
 * @param {Function} trackInfoFn - a callback that receives track info
 * @param {Function} timingInfoFn - a callback that receives timing info
 * @param {Function} id3Fn - a callback that receives ID3 metadata
 * @param {Function} captionsFn - a callback that receives captions
 * @param {Function} dataFn - a callback that is executed when segment bytes are available
 *                            and ready to use
 * @param {Function} doneFn - a callback that is executed after all resources have been
 *                            downloaded and any decryption completed
 */


var waitForCompletion = function waitForCompletion(_ref6) {
  var activeXhrs = _ref6.activeXhrs,
      decryptionWorker = _ref6.decryptionWorker,
      trackInfoFn = _ref6.trackInfoFn,
      timingInfoFn = _ref6.timingInfoFn,
      videoSegmentTimingInfoFn = _ref6.videoSegmentTimingInfoFn,
      id3Fn = _ref6.id3Fn,
      captionsFn = _ref6.captionsFn,
      dataFn = _ref6.dataFn,
      doneFn = _ref6.doneFn;
  var count = 0;
  var didError = false;
  return function (error, segment) {
    if (didError) {
      return;
    }

    if (error) {
      didError = true; // If there are errors, we have to abort any outstanding requests

      abortAll(activeXhrs); // Even though the requests above are aborted, and in theory we could wait until we
      // handle the aborted events from those requests, there are some cases where we may
      // never get an aborted event. For instance, if the network connection is lost and
      // there were two requests, the first may have triggered an error immediately, while
      // the second request remains unsent. In that case, the aborted algorithm will not
      // trigger an abort: see https://xhr.spec.whatwg.org/#the-abort()-method
      //
      // We also can't rely on the ready state of the XHR, since the request that
      // triggered the connection error may also show as a ready state of 0 (unsent).
      // Therefore, we have to finish this group of requests immediately after the first
      // seen error.

      return doneFn(error, segment);
    }

    count += 1;

    if (count === activeXhrs.length) {
      // Keep track of when *all* of the requests have completed
      segment.endOfAllRequests = Date.now();

      if (segment.encryptedBytes) {
        return decryptSegment({
          decryptionWorker: decryptionWorker,
          segment: segment,
          trackInfoFn: trackInfoFn,
          timingInfoFn: timingInfoFn,
          videoSegmentTimingInfoFn: videoSegmentTimingInfoFn,
          id3Fn: id3Fn,
          captionsFn: captionsFn,
          dataFn: dataFn,
          doneFn: doneFn
        });
      } // Otherwise, everything is ready just continue


      handleSegmentBytes({
        segment: segment,
        bytes: segment.bytes,
        isPartial: false,
        trackInfoFn: trackInfoFn,
        timingInfoFn: timingInfoFn,
        videoSegmentTimingInfoFn: videoSegmentTimingInfoFn,
        id3Fn: id3Fn,
        captionsFn: captionsFn,
        dataFn: dataFn,
        doneFn: doneFn
      });
    }
  };
};
/**
 * Calls the abort callback if any request within the batch was aborted. Will only call
 * the callback once per batch of requests, even if multiple were aborted.
 *
 * @param {Object} loadendState - state to check to see if the abort function was called
 * @param {Function} abortFn - callback to call for abort
 */


var handleLoadEnd = function handleLoadEnd(_ref7) {
  var loadendState = _ref7.loadendState,
      abortFn = _ref7.abortFn;
  return function (event) {
    var request = event.target;

    if (request.aborted && abortFn && !loadendState.calledAbortFn) {
      abortFn();
      loadendState.calledAbortFn = true;
    }
  };
};
/**
 * Simple progress event callback handler that gathers some stats before
 * executing a provided callback with the `segment` object
 *
 * @param {Object} segment - a simplified copy of the segmentInfo object
 *                           from SegmentLoader
 * @param {Function} progressFn - a callback that is executed each time a progress event
 *                                is received
 * @param {Function} trackInfoFn - a callback that receives track info
 * @param {Function} dataFn - a callback that is executed when segment bytes are available
 *                            and ready to use
 * @param {Event} event - the progress event object from XMLHttpRequest
 */


var handleProgress = function handleProgress(_ref8) {
  var segment = _ref8.segment,
      progressFn = _ref8.progressFn,
      trackInfoFn = _ref8.trackInfoFn,
      timingInfoFn = _ref8.timingInfoFn,
      videoSegmentTimingInfoFn = _ref8.videoSegmentTimingInfoFn,
      id3Fn = _ref8.id3Fn,
      captionsFn = _ref8.captionsFn,
      dataFn = _ref8.dataFn,
      handlePartialData = _ref8.handlePartialData;
  return function (event) {
    var request = event.target;

    if (request.aborted) {
      return;
    } // don't support encrypted segments or fmp4 for now


    if (handlePartialData && !segment.key && // although responseText "should" exist, this guard serves to prevent an error being
    // thrown on the next check for two primary cases:
    // 1. the mime type override stops working, or is not implemented for a specific
    //    browser
    // 2. when using mock XHR libraries like sinon that do not allow the override behavior
    request.responseText && // in order to determine if it's an fmp4 we need at least 8 bytes
    request.responseText.length >= 8) {
      var newBytes = stringToArrayBuffer(request.responseText.substring(segment.lastReachedChar || 0));

      if (segment.lastReachedChar || !containers.isLikelyFmp4MediaSegment(new Uint8Array(newBytes))) {
        segment.lastReachedChar = request.responseText.length;
        handleSegmentBytes({
          segment: segment,
          bytes: newBytes,
          isPartial: true,
          trackInfoFn: trackInfoFn,
          timingInfoFn: timingInfoFn,
          videoSegmentTimingInfoFn: videoSegmentTimingInfoFn,
          id3Fn: id3Fn,
          captionsFn: captionsFn,
          dataFn: dataFn
        });
      }
    }

    segment.stats = videojs$1.mergeOptions(segment.stats, getProgressStats(event)); // record the time that we receive the first byte of data

    if (!segment.stats.firstBytesReceivedAt && segment.stats.bytesReceived) {
      segment.stats.firstBytesReceivedAt = Date.now();
    }

    return progressFn(event, segment);
  };
};
/**
 * Load all resources and does any processing necessary for a media-segment
 *
 * Features:
 *   decrypts the media-segment if it has a key uri and an iv
 *   aborts *all* requests if *any* one request fails
 *
 * The segment object, at minimum, has the following format:
 * {
 *   resolvedUri: String,
 *   [transmuxer]: Object,
 *   [byterange]: {
 *     offset: Number,
 *     length: Number
 *   },
 *   [key]: {
 *     resolvedUri: String
 *     [byterange]: {
 *       offset: Number,
 *       length: Number
 *     },
 *     iv: {
 *       bytes: Uint32Array
 *     }
 *   },
 *   [map]: {
 *     resolvedUri: String,
 *     [byterange]: {
 *       offset: Number,
 *       length: Number
 *     },
 *     [bytes]: Uint8Array
 *   }
 * }
 * ...where [name] denotes optional properties
 *
 * @param {Function} xhr - an instance of the xhr wrapper in xhr.js
 * @param {Object} xhrOptions - the base options to provide to all xhr requests
 * @param {WebWorker} decryptionWorker - a WebWorker interface to AES-128
 *                                       decryption routines
 * @param {Object} segment - a simplified copy of the segmentInfo object
 *                           from SegmentLoader
 * @param {Function} abortFn - a callback called (only once) if any piece of a request was
 *                             aborted
 * @param {Function} progressFn - a callback that receives progress events from the main
 *                                segment's xhr request
 * @param {Function} trackInfoFn - a callback that receives track info
 * @param {Function} id3Fn - a callback that receives ID3 metadata
 * @param {Function} captionsFn - a callback that receives captions
 * @param {Function} dataFn - a callback that receives data from the main segment's xhr
 *                            request, transmuxed if needed
 * @param {Function} doneFn - a callback that is executed only once all requests have
 *                            succeeded or failed
 * @return {Function} a function that, when invoked, immediately aborts all
 *                     outstanding requests
 */


var mediaSegmentRequest = function mediaSegmentRequest(_ref9) {
  var xhr = _ref9.xhr,
      xhrOptions = _ref9.xhrOptions,
      decryptionWorker = _ref9.decryptionWorker,
      segment = _ref9.segment,
      abortFn = _ref9.abortFn,
      progressFn = _ref9.progressFn,
      trackInfoFn = _ref9.trackInfoFn,
      timingInfoFn = _ref9.timingInfoFn,
      videoSegmentTimingInfoFn = _ref9.videoSegmentTimingInfoFn,
      id3Fn = _ref9.id3Fn,
      captionsFn = _ref9.captionsFn,
      dataFn = _ref9.dataFn,
      doneFn = _ref9.doneFn,
      handlePartialData = _ref9.handlePartialData;
  var activeXhrs = [];
  var finishProcessingFn = waitForCompletion({
    activeXhrs: activeXhrs,
    decryptionWorker: decryptionWorker,
    trackInfoFn: trackInfoFn,
    timingInfoFn: timingInfoFn,
    videoSegmentTimingInfoFn: videoSegmentTimingInfoFn,
    id3Fn: id3Fn,
    captionsFn: captionsFn,
    dataFn: dataFn,
    doneFn: doneFn
  }); // optionally, request the decryption key

  if (segment.key && !segment.key.bytes) {
    var keyRequestOptions = videojs$1.mergeOptions(xhrOptions, {
      uri: segment.key.resolvedUri,
      responseType: 'arraybuffer'
    });
    var keyRequestCallback = handleKeyResponse(segment, finishProcessingFn);
    var keyXhr = xhr(keyRequestOptions, keyRequestCallback);
    activeXhrs.push(keyXhr);
  } // optionally, request the associated media init segment


  if (segment.map && !segment.map.bytes) {
    var initSegmentOptions = videojs$1.mergeOptions(xhrOptions, {
      uri: segment.map.resolvedUri,
      responseType: 'arraybuffer',
      headers: segmentXhrHeaders(segment.map)
    });
    var initSegmentRequestCallback = handleInitSegmentResponse({
      segment: segment,
      finishProcessingFn: finishProcessingFn
    });
    var initSegmentXhr = xhr(initSegmentOptions, initSegmentRequestCallback);
    activeXhrs.push(initSegmentXhr);
  }

  var segmentRequestOptions = videojs$1.mergeOptions(xhrOptions, {
    uri: segment.resolvedUri,
    responseType: 'arraybuffer',
    headers: segmentXhrHeaders(segment)
  });

  if (handlePartialData) {
    // setting to text is required for partial responses
    // conversion to ArrayBuffer happens later
    segmentRequestOptions.responseType = 'text';

    segmentRequestOptions.beforeSend = function (xhrObject) {
      // XHR binary charset opt by Marcus Granado 2006 [http://mgran.blogspot.com]
      // makes the browser pass through the "text" unparsed
      xhrObject.overrideMimeType('text/plain; charset=x-user-defined');
    };
  }

  var segmentRequestCallback = handleSegmentResponse({
    segment: segment,
    finishProcessingFn: finishProcessingFn,
    responseType: segmentRequestOptions.responseType
  });
  var segmentXhr = xhr(segmentRequestOptions, segmentRequestCallback);
  segmentXhr.addEventListener('progress', handleProgress({
    segment: segment,
    progressFn: progressFn,
    trackInfoFn: trackInfoFn,
    timingInfoFn: timingInfoFn,
    videoSegmentTimingInfoFn: videoSegmentTimingInfoFn,
    id3Fn: id3Fn,
    captionsFn: captionsFn,
    dataFn: dataFn,
    handlePartialData: handlePartialData
  }));
  activeXhrs.push(segmentXhr); // since all parts of the request must be considered, but should not make callbacks
  // multiple times, provide a shared state object

  var loadendState = {};
  activeXhrs.forEach(function (activeXhr) {
    activeXhr.addEventListener('loadend', handleLoadEnd({
      loadendState: loadendState,
      abortFn: abortFn
    }));
  });
  return function () {
    return abortAll(activeXhrs);
  };
};

var win = typeof window !== 'undefined' ? window : {},
    TARGET = typeof Symbol === 'undefined' ? '__target' : Symbol(),
    SCRIPT_TYPE = 'application/javascript',
    BlobBuilder = win.BlobBuilder || win.WebKitBlobBuilder || win.MozBlobBuilder || win.MSBlobBuilder,
    URL = win.URL || win.webkitURL || (URL && URL.msURL),
    Worker = win.Worker;

/**
 * Returns a wrapper around Web Worker code that is constructible.
 *
 * @function shimWorker
 *
 * @param { String }    filename    The name of the file
 * @param { Function }  fn          Function wrapping the code of the worker
 */
function shimWorker (filename, fn) {
    return function ShimWorker (forceFallback) {
        var o = this;

        if (!fn) {
            return new Worker(filename);
        }
        else if (Worker && !forceFallback) {
            // Convert the function's inner code to a string to construct the worker
            var source = fn.toString().replace(/^function.+?{/, '').slice(0, -1),
                objURL = createSourceObject(source);

            this[TARGET] = new Worker(objURL);
            wrapTerminate(this[TARGET], objURL);
            return this[TARGET];
        }
        else {
            var selfShim = {
                    postMessage: function(m) {
                        if (o.onmessage) {
                            setTimeout(function(){ o.onmessage({ data: m, target: selfShim }); });
                        }
                    }
                };

            fn.call(selfShim);
            this.postMessage = function(m) {
                setTimeout(function(){ selfShim.onmessage({ data: m, target: o }); });
            };
            this.isThisThread = true;
        }
    };
}
// Test Worker capabilities
if (Worker) {
    var testWorker,
        objURL = createSourceObject('self.onmessage = function () {}'),
        testArray = new Uint8Array(1);

    try {
        testWorker = new Worker(objURL);

        // Native browser on some Samsung devices throws for transferables, let's detect it
        testWorker.postMessage(testArray, [testArray.buffer]);
    }
    catch (e) {
        Worker = null;
    }
    finally {
        URL.revokeObjectURL(objURL);
        if (testWorker) {
            testWorker.terminate();
        }
    }
}

function createSourceObject(str) {
    try {
        return URL.createObjectURL(new Blob([str], { type: SCRIPT_TYPE }));
    }
    catch (e) {
        var blob = new BlobBuilder();
        blob.append(str);
        return URL.createObjectURL(blob.getBlob(type));
    }
}

function wrapTerminate(worker, objURL){
    if(!worker || !objURL) return;
    var term = worker.terminate;
    worker.objURL = objURL;
    worker.terminate = function(){
        if(worker.objURL)
            URL.revokeObjectURL(worker.objURL);
        term.call(worker);
    };
}

var TransmuxWorker = new shimWorker("./transmuxer-worker.worker.js", function (window, document) {
  var self = this;
  /*! @name @videojs/http-streaming @version 2.2.0 @license Apache-2.0 */

  var transmuxerWorker = function () {
    /**
     * mux.js
     *
     * Copyright (c) Brightcove
     * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
     *
     * A lightweight readable stream implemention that handles event dispatching.
     * Objects that inherit from streams should call init in their constructors.
     */

    var Stream = function Stream() {
      this.init = function () {
        var listeners = {};
        /**
         * Add a listener for a specified event type.
         * @param type {string} the event name
         * @param listener {function} the callback to be invoked when an event of
         * the specified type occurs
         */

        this.on = function (type, listener) {
          if (!listeners[type]) {
            listeners[type] = [];
          }

          listeners[type] = listeners[type].concat(listener);
        };
        /**
         * Remove a listener for a specified event type.
         * @param type {string} the event name
         * @param listener {function} a function previously registered for this
         * type of event through `on`
         */


        this.off = function (type, listener) {
          var index;

          if (!listeners[type]) {
            return false;
          }

          index = listeners[type].indexOf(listener);
          listeners[type] = listeners[type].slice();
          listeners[type].splice(index, 1);
          return index > -1;
        };
        /**
         * Trigger an event of the specified type on this stream. Any additional
         * arguments to this function are passed as parameters to event listeners.
         * @param type {string} the event name
         */


        this.trigger = function (type) {
          var callbacks, i, length, args;
          callbacks = listeners[type];

          if (!callbacks) {
            return;
          } // Slicing the arguments on every invocation of this method
          // can add a significant amount of overhead. Avoid the
          // intermediate object creation for the common case of a
          // single callback argument


          if (arguments.length === 2) {
            length = callbacks.length;

            for (i = 0; i < length; ++i) {
              callbacks[i].call(this, arguments[1]);
            }
          } else {
            args = [];
            i = arguments.length;

            for (i = 1; i < arguments.length; ++i) {
              args.push(arguments[i]);
            }

            length = callbacks.length;

            for (i = 0; i < length; ++i) {
              callbacks[i].apply(this, args);
            }
          }
        };
        /**
         * Destroys the stream and cleans up.
         */


        this.dispose = function () {
          listeners = {};
        };
      };
    };
    /**
     * Forwards all `data` events on this stream to the destination stream. The
     * destination stream should provide a method `push` to receive the data
     * events as they arrive.
     * @param destination {stream} the stream that will receive all `data` events
     * @param autoFlush {boolean} if false, we will not call `flush` on the destination
     *                            when the current stream emits a 'done' event
     * @see http://nodejs.org/api/stream.html#stream_readable_pipe_destination_options
     */


    Stream.prototype.pipe = function (destination) {
      this.on('data', function (data) {
        destination.push(data);
      });
      this.on('done', function (flushSource) {
        destination.flush(flushSource);
      });
      this.on('partialdone', function (flushSource) {
        destination.partialFlush(flushSource);
      });
      this.on('endedtimeline', function (flushSource) {
        destination.endTimeline(flushSource);
      });
      this.on('reset', function (flushSource) {
        destination.reset(flushSource);
      });
      return destination;
    }; // Default stream functions that are expected to be overridden to perform
    // actual work. These are provided by the prototype as a sort of no-op
    // implementation so that we don't have to check for their existence in the
    // `pipe` function above.


    Stream.prototype.push = function (data) {
      this.trigger('data', data);
    };

    Stream.prototype.flush = function (flushSource) {
      this.trigger('done', flushSource);
    };

    Stream.prototype.partialFlush = function (flushSource) {
      this.trigger('partialdone', flushSource);
    };

    Stream.prototype.endTimeline = function (flushSource) {
      this.trigger('endedtimeline', flushSource);
    };

    Stream.prototype.reset = function (flushSource) {
      this.trigger('reset', flushSource);
    };

    var stream = Stream;
    /**
     * mux.js
     *
     * Copyright (c) Brightcove
     * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
     *
     * Functions that generate fragmented MP4s suitable for use with Media
     * Source Extensions.
     */

    var UINT32_MAX = Math.pow(2, 32) - 1;
    var box, dinf, esds, ftyp, mdat, mfhd, minf, moof, moov, mvex, mvhd, trak, tkhd, mdia, mdhd, hdlr, sdtp, stbl, stsd, traf, trex, trun, types, MAJOR_BRAND, MINOR_VERSION, AVC1_BRAND, VIDEO_HDLR, AUDIO_HDLR, HDLR_TYPES, VMHD, SMHD, DREF, STCO, STSC, STSZ, STTS; // pre-calculate constants

    (function () {
      var i;
      types = {
        avc1: [],
        // codingname
        avcC: [],
        btrt: [],
        dinf: [],
        dref: [],
        esds: [],
        ftyp: [],
        hdlr: [],
        mdat: [],
        mdhd: [],
        mdia: [],
        mfhd: [],
        minf: [],
        moof: [],
        moov: [],
        mp4a: [],
        // codingname
        mvex: [],
        mvhd: [],
        pasp: [],
        sdtp: [],
        smhd: [],
        stbl: [],
        stco: [],
        stsc: [],
        stsd: [],
        stsz: [],
        stts: [],
        styp: [],
        tfdt: [],
        tfhd: [],
        traf: [],
        trak: [],
        trun: [],
        trex: [],
        tkhd: [],
        vmhd: []
      }; // In environments where Uint8Array is undefined (e.g., IE8), skip set up so that we
      // don't throw an error

      if (typeof Uint8Array === 'undefined') {
        return;
      }

      for (i in types) {
        if (types.hasOwnProperty(i)) {
          types[i] = [i.charCodeAt(0), i.charCodeAt(1), i.charCodeAt(2), i.charCodeAt(3)];
        }
      }

      MAJOR_BRAND = new Uint8Array(['i'.charCodeAt(0), 's'.charCodeAt(0), 'o'.charCodeAt(0), 'm'.charCodeAt(0)]);
      AVC1_BRAND = new Uint8Array(['a'.charCodeAt(0), 'v'.charCodeAt(0), 'c'.charCodeAt(0), '1'.charCodeAt(0)]);
      MINOR_VERSION = new Uint8Array([0, 0, 0, 1]);
      VIDEO_HDLR = new Uint8Array([0x00, // version 0
      0x00, 0x00, 0x00, // flags
      0x00, 0x00, 0x00, 0x00, // pre_defined
      0x76, 0x69, 0x64, 0x65, // handler_type: 'vide'
      0x00, 0x00, 0x00, 0x00, // reserved
      0x00, 0x00, 0x00, 0x00, // reserved
      0x00, 0x00, 0x00, 0x00, // reserved
      0x56, 0x69, 0x64, 0x65, 0x6f, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x00 // name: 'VideoHandler'
      ]);
      AUDIO_HDLR = new Uint8Array([0x00, // version 0
      0x00, 0x00, 0x00, // flags
      0x00, 0x00, 0x00, 0x00, // pre_defined
      0x73, 0x6f, 0x75, 0x6e, // handler_type: 'soun'
      0x00, 0x00, 0x00, 0x00, // reserved
      0x00, 0x00, 0x00, 0x00, // reserved
      0x00, 0x00, 0x00, 0x00, // reserved
      0x53, 0x6f, 0x75, 0x6e, 0x64, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x00 // name: 'SoundHandler'
      ]);
      HDLR_TYPES = {
        video: VIDEO_HDLR,
        audio: AUDIO_HDLR
      };
      DREF = new Uint8Array([0x00, // version 0
      0x00, 0x00, 0x00, // flags
      0x00, 0x00, 0x00, 0x01, // entry_count
      0x00, 0x00, 0x00, 0x0c, // entry_size
      0x75, 0x72, 0x6c, 0x20, // 'url' type
      0x00, // version 0
      0x00, 0x00, 0x01 // entry_flags
      ]);
      SMHD = new Uint8Array([0x00, // version
      0x00, 0x00, 0x00, // flags
      0x00, 0x00, // balance, 0 means centered
      0x00, 0x00 // reserved
      ]);
      STCO = new Uint8Array([0x00, // version
      0x00, 0x00, 0x00, // flags
      0x00, 0x00, 0x00, 0x00 // entry_count
      ]);
      STSC = STCO;
      STSZ = new Uint8Array([0x00, // version
      0x00, 0x00, 0x00, // flags
      0x00, 0x00, 0x00, 0x00, // sample_size
      0x00, 0x00, 0x00, 0x00 // sample_count
      ]);
      STTS = STCO;
      VMHD = new Uint8Array([0x00, // version
      0x00, 0x00, 0x01, // flags
      0x00, 0x00, // graphicsmode
      0x00, 0x00, 0x00, 0x00, 0x00, 0x00 // opcolor
      ]);
    })();

    box = function box(type) {
      var payload = [],
          size = 0,
          i,
          result,
          view;

      for (i = 1; i < arguments.length; i++) {
        payload.push(arguments[i]);
      }

      i = payload.length; // calculate the total size we need to allocate

      while (i--) {
        size += payload[i].byteLength;
      }

      result = new Uint8Array(size + 8);
      view = new DataView(result.buffer, result.byteOffset, result.byteLength);
      view.setUint32(0, result.byteLength);
      result.set(type, 4); // copy the payload into the result

      for (i = 0, size = 8; i < payload.length; i++) {
        result.set(payload[i], size);
        size += payload[i].byteLength;
      }

      return result;
    };

    dinf = function dinf() {
      return box(types.dinf, box(types.dref, DREF));
    };

    esds = function esds(track) {
      return box(types.esds, new Uint8Array([0x00, // version
      0x00, 0x00, 0x00, // flags
      // ES_Descriptor
      0x03, // tag, ES_DescrTag
      0x19, // length
      0x00, 0x00, // ES_ID
      0x00, // streamDependenceFlag, URL_flag, reserved, streamPriority
      // DecoderConfigDescriptor
      0x04, // tag, DecoderConfigDescrTag
      0x11, // length
      0x40, // object type
      0x15, // streamType
      0x00, 0x06, 0x00, // bufferSizeDB
      0x00, 0x00, 0xda, 0xc0, // maxBitrate
      0x00, 0x00, 0xda, 0xc0, // avgBitrate
      // DecoderSpecificInfo
      0x05, // tag, DecoderSpecificInfoTag
      0x02, // length
      // ISO/IEC 14496-3, AudioSpecificConfig
      // for samplingFrequencyIndex see ISO/IEC 13818-7:2006, 8.1.3.2.2, Table 35
      track.audioobjecttype << 3 | track.samplingfrequencyindex >>> 1, track.samplingfrequencyindex << 7 | track.channelcount << 3, 0x06, 0x01, 0x02 // GASpecificConfig
      ]));
    };

    ftyp = function ftyp() {
      return box(types.ftyp, MAJOR_BRAND, MINOR_VERSION, MAJOR_BRAND, AVC1_BRAND);
    };

    hdlr = function hdlr(type) {
      return box(types.hdlr, HDLR_TYPES[type]);
    };

    mdat = function mdat(data) {
      return box(types.mdat, data);
    };

    mdhd = function mdhd(track) {
      var result = new Uint8Array([0x00, // version 0
      0x00, 0x00, 0x00, // flags
      0x00, 0x00, 0x00, 0x02, // creation_time
      0x00, 0x00, 0x00, 0x03, // modification_time
      0x00, 0x01, 0x5f, 0x90, // timescale, 90,000 "ticks" per second
      track.duration >>> 24 & 0xFF, track.duration >>> 16 & 0xFF, track.duration >>> 8 & 0xFF, track.duration & 0xFF, // duration
      0x55, 0xc4, // 'und' language (undetermined)
      0x00, 0x00]); // Use the sample rate from the track metadata, when it is
      // defined. The sample rate can be parsed out of an ADTS header, for
      // instance.

      if (track.samplerate) {
        result[12] = track.samplerate >>> 24 & 0xFF;
        result[13] = track.samplerate >>> 16 & 0xFF;
        result[14] = track.samplerate >>> 8 & 0xFF;
        result[15] = track.samplerate & 0xFF;
      }

      return box(types.mdhd, result);
    };

    mdia = function mdia(track) {
      return box(types.mdia, mdhd(track), hdlr(track.type), minf(track));
    };

    mfhd = function mfhd(sequenceNumber) {
      return box(types.mfhd, new Uint8Array([0x00, 0x00, 0x00, 0x00, // flags
      (sequenceNumber & 0xFF000000) >> 24, (sequenceNumber & 0xFF0000) >> 16, (sequenceNumber & 0xFF00) >> 8, sequenceNumber & 0xFF // sequence_number
      ]));
    };

    minf = function minf(track) {
      return box(types.minf, track.type === 'video' ? box(types.vmhd, VMHD) : box(types.smhd, SMHD), dinf(), stbl(track));
    };

    moof = function moof(sequenceNumber, tracks) {
      var trackFragments = [],
          i = tracks.length; // build traf boxes for each track fragment

      while (i--) {
        trackFragments[i] = traf(tracks[i]);
      }

      return box.apply(null, [types.moof, mfhd(sequenceNumber)].concat(trackFragments));
    };
    /**
     * Returns a movie box.
     * @param tracks {array} the tracks associated with this movie
     * @see ISO/IEC 14496-12:2012(E), section 8.2.1
     */


    moov = function moov(tracks) {
      var i = tracks.length,
          boxes = [];

      while (i--) {
        boxes[i] = trak(tracks[i]);
      }

      return box.apply(null, [types.moov, mvhd(0xffffffff)].concat(boxes).concat(mvex(tracks)));
    };

    mvex = function mvex(tracks) {
      var i = tracks.length,
          boxes = [];

      while (i--) {
        boxes[i] = trex(tracks[i]);
      }

      return box.apply(null, [types.mvex].concat(boxes));
    };

    mvhd = function mvhd(duration) {
      var bytes = new Uint8Array([0x00, // version 0
      0x00, 0x00, 0x00, // flags
      0x00, 0x00, 0x00, 0x01, // creation_time
      0x00, 0x00, 0x00, 0x02, // modification_time
      0x00, 0x01, 0x5f, 0x90, // timescale, 90,000 "ticks" per second
      (duration & 0xFF000000) >> 24, (duration & 0xFF0000) >> 16, (duration & 0xFF00) >> 8, duration & 0xFF, // duration
      0x00, 0x01, 0x00, 0x00, // 1.0 rate
      0x01, 0x00, // 1.0 volume
      0x00, 0x00, // reserved
      0x00, 0x00, 0x00, 0x00, // reserved
      0x00, 0x00, 0x00, 0x00, // reserved
      0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, // transformation: unity matrix
      0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // pre_defined
      0xff, 0xff, 0xff, 0xff // next_track_ID
      ]);
      return box(types.mvhd, bytes);
    };

    sdtp = function sdtp(track) {
      var samples = track.samples || [],
          bytes = new Uint8Array(4 + samples.length),
          flags,
          i; // leave the full box header (4 bytes) all zero
      // write the sample table

      for (i = 0; i < samples.length; i++) {
        flags = samples[i].flags;
        bytes[i + 4] = flags.dependsOn << 4 | flags.isDependedOn << 2 | flags.hasRedundancy;
      }

      return box(types.sdtp, bytes);
    };

    stbl = function stbl(track) {
      return box(types.stbl, stsd(track), box(types.stts, STTS), box(types.stsc, STSC), box(types.stsz, STSZ), box(types.stco, STCO));
    };

    (function () {
      var videoSample, audioSample;

      stsd = function stsd(track) {
        return box(types.stsd, new Uint8Array([0x00, // version 0
        0x00, 0x00, 0x00, // flags
        0x00, 0x00, 0x00, 0x01]), track.type === 'video' ? videoSample(track) : audioSample(track));
      };

      videoSample = function videoSample(track) {
        var sps = track.sps || [],
            pps = track.pps || [],
            sequenceParameterSets = [],
            pictureParameterSets = [],
            i,
            avc1Box; // assemble the SPSs

        for (i = 0; i < sps.length; i++) {
          sequenceParameterSets.push((sps[i].byteLength & 0xFF00) >>> 8);
          sequenceParameterSets.push(sps[i].byteLength & 0xFF); // sequenceParameterSetLength

          sequenceParameterSets = sequenceParameterSets.concat(Array.prototype.slice.call(sps[i])); // SPS
        } // assemble the PPSs


        for (i = 0; i < pps.length; i++) {
          pictureParameterSets.push((pps[i].byteLength & 0xFF00) >>> 8);
          pictureParameterSets.push(pps[i].byteLength & 0xFF);
          pictureParameterSets = pictureParameterSets.concat(Array.prototype.slice.call(pps[i]));
        }

        avc1Box = [types.avc1, new Uint8Array([0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved
        0x00, 0x01, // data_reference_index
        0x00, 0x00, // pre_defined
        0x00, 0x00, // reserved
        0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // pre_defined
        (track.width & 0xff00) >> 8, track.width & 0xff, // width
        (track.height & 0xff00) >> 8, track.height & 0xff, // height
        0x00, 0x48, 0x00, 0x00, // horizresolution
        0x00, 0x48, 0x00, 0x00, // vertresolution
        0x00, 0x00, 0x00, 0x00, // reserved
        0x00, 0x01, // frame_count
        0x13, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x6a, 0x73, 0x2d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x69, 0x62, 0x2d, 0x68, 0x6c, 0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // compressorname
        0x00, 0x18, // depth = 24
        0x11, 0x11 // pre_defined = -1
        ]), box(types.avcC, new Uint8Array([0x01, // configurationVersion
        track.profileIdc, // AVCProfileIndication
        track.profileCompatibility, // profile_compatibility
        track.levelIdc, // AVCLevelIndication
        0xff // lengthSizeMinusOne, hard-coded to 4 bytes
        ].concat([sps.length], // numOfSequenceParameterSets
        sequenceParameterSets, // "SPS"
        [pps.length], // numOfPictureParameterSets
        pictureParameterSets // "PPS"
        ))), box(types.btrt, new Uint8Array([0x00, 0x1c, 0x9c, 0x80, // bufferSizeDB
        0x00, 0x2d, 0xc6, 0xc0, // maxBitrate
        0x00, 0x2d, 0xc6, 0xc0 // avgBitrate
        ]))];

        if (track.sarRatio) {
          var hSpacing = track.sarRatio[0],
              vSpacing = track.sarRatio[1];
          avc1Box.push(box(types.pasp, new Uint8Array([(hSpacing & 0xFF000000) >> 24, (hSpacing & 0xFF0000) >> 16, (hSpacing & 0xFF00) >> 8, hSpacing & 0xFF, (vSpacing & 0xFF000000) >> 24, (vSpacing & 0xFF0000) >> 16, (vSpacing & 0xFF00) >> 8, vSpacing & 0xFF])));
        }

        return box.apply(null, avc1Box);
      };

      audioSample = function audioSample(track) {
        return box(types.mp4a, new Uint8Array([// SampleEntry, ISO/IEC 14496-12
        0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved
        0x00, 0x01, // data_reference_index
        // AudioSampleEntry, ISO/IEC 14496-12
        0x00, 0x00, 0x00, 0x00, // reserved
        0x00, 0x00, 0x00, 0x00, // reserved
        (track.channelcount & 0xff00) >> 8, track.channelcount & 0xff, // channelcount
        (track.samplesize & 0xff00) >> 8, track.samplesize & 0xff, // samplesize
        0x00, 0x00, // pre_defined
        0x00, 0x00, // reserved
        (track.samplerate & 0xff00) >> 8, track.samplerate & 0xff, 0x00, 0x00 // samplerate, 16.16
        // MP4AudioSampleEntry, ISO/IEC 14496-14
        ]), esds(track));
      };
    })();

    tkhd = function tkhd(track) {
      var result = new Uint8Array([0x00, // version 0
      0x00, 0x00, 0x07, // flags
      0x00, 0x00, 0x00, 0x00, // creation_time
      0x00, 0x00, 0x00, 0x00, // modification_time
      (track.id & 0xFF000000) >> 24, (track.id & 0xFF0000) >> 16, (track.id & 0xFF00) >> 8, track.id & 0xFF, // track_ID
      0x00, 0x00, 0x00, 0x00, // reserved
      (track.duration & 0xFF000000) >> 24, (track.duration & 0xFF0000) >> 16, (track.duration & 0xFF00) >> 8, track.duration & 0xFF, // duration
      0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved
      0x00, 0x00, // layer
      0x00, 0x00, // alternate_group
      0x01, 0x00, // non-audio track volume
      0x00, 0x00, // reserved
      0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, // transformation: unity matrix
      (track.width & 0xFF00) >> 8, track.width & 0xFF, 0x00, 0x00, // width
      (track.height & 0xFF00) >> 8, track.height & 0xFF, 0x00, 0x00 // height
      ]);
      return box(types.tkhd, result);
    };
    /**
     * Generate a track fragment (traf) box. A traf box collects metadata
     * about tracks in a movie fragment (moof) box.
     */


    traf = function traf(track) {
      var trackFragmentHeader, trackFragmentDecodeTime, trackFragmentRun, sampleDependencyTable, dataOffset, upperWordBaseMediaDecodeTime, lowerWordBaseMediaDecodeTime;
      trackFragmentHeader = box(types.tfhd, new Uint8Array([0x00, // version 0
      0x00, 0x00, 0x3a, // flags
      (track.id & 0xFF000000) >> 24, (track.id & 0xFF0000) >> 16, (track.id & 0xFF00) >> 8, track.id & 0xFF, // track_ID
      0x00, 0x00, 0x00, 0x01, // sample_description_index
      0x00, 0x00, 0x00, 0x00, // default_sample_duration
      0x00, 0x00, 0x00, 0x00, // default_sample_size
      0x00, 0x00, 0x00, 0x00 // default_sample_flags
      ]));
      upperWordBaseMediaDecodeTime = Math.floor(track.baseMediaDecodeTime / (UINT32_MAX + 1));
      lowerWordBaseMediaDecodeTime = Math.floor(track.baseMediaDecodeTime % (UINT32_MAX + 1));
      trackFragmentDecodeTime = box(types.tfdt, new Uint8Array([0x01, // version 1
      0x00, 0x00, 0x00, // flags
      // baseMediaDecodeTime
      upperWordBaseMediaDecodeTime >>> 24 & 0xFF, upperWordBaseMediaDecodeTime >>> 16 & 0xFF, upperWordBaseMediaDecodeTime >>> 8 & 0xFF, upperWordBaseMediaDecodeTime & 0xFF, lowerWordBaseMediaDecodeTime >>> 24 & 0xFF, lowerWordBaseMediaDecodeTime >>> 16 & 0xFF, lowerWordBaseMediaDecodeTime >>> 8 & 0xFF, lowerWordBaseMediaDecodeTime & 0xFF])); // the data offset specifies the number of bytes from the start of
      // the containing moof to the first payload byte of the associated
      // mdat

      dataOffset = 32 + // tfhd
      20 + // tfdt
      8 + // traf header
      16 + // mfhd
      8 + // moof header
      8; // mdat header
      // audio tracks require less metadata

      if (track.type === 'audio') {
        trackFragmentRun = trun(track, dataOffset);
        return box(types.traf, trackFragmentHeader, trackFragmentDecodeTime, trackFragmentRun);
      } // video tracks should contain an independent and disposable samples
      // box (sdtp)
      // generate one and adjust offsets to match


      sampleDependencyTable = sdtp(track);
      trackFragmentRun = trun(track, sampleDependencyTable.length + dataOffset);
      return box(types.traf, trackFragmentHeader, trackFragmentDecodeTime, trackFragmentRun, sampleDependencyTable);
    };
    /**
     * Generate a track box.
     * @param track {object} a track definition
     * @return {Uint8Array} the track box
     */


    trak = function trak(track) {
      track.duration = track.duration || 0xffffffff;
      return box(types.trak, tkhd(track), mdia(track));
    };

    trex = function trex(track) {
      var result = new Uint8Array([0x00, // version 0
      0x00, 0x00, 0x00, // flags
      (track.id & 0xFF000000) >> 24, (track.id & 0xFF0000) >> 16, (track.id & 0xFF00) >> 8, track.id & 0xFF, // track_ID
      0x00, 0x00, 0x00, 0x01, // default_sample_description_index
      0x00, 0x00, 0x00, 0x00, // default_sample_duration
      0x00, 0x00, 0x00, 0x00, // default_sample_size
      0x00, 0x01, 0x00, 0x01 // default_sample_flags
      ]); // the last two bytes of default_sample_flags is the sample
      // degradation priority, a hint about the importance of this sample
      // relative to others. Lower the degradation priority for all sample
      // types other than video.

      if (track.type !== 'video') {
        result[result.length - 1] = 0x00;
      }

      return box(types.trex, result);
    };

    (function () {
      var audioTrun, videoTrun, trunHeader; // This method assumes all samples are uniform. That is, if a
      // duration is present for the first sample, it will be present for
      // all subsequent samples.
      // see ISO/IEC 14496-12:2012, Section 8.8.8.1

      trunHeader = function trunHeader(samples, offset) {
        var durationPresent = 0,
            sizePresent = 0,
            flagsPresent = 0,
            compositionTimeOffset = 0; // trun flag constants

        if (samples.length) {
          if (samples[0].duration !== undefined) {
            durationPresent = 0x1;
          }

          if (samples[0].size !== undefined) {
            sizePresent = 0x2;
          }

          if (samples[0].flags !== undefined) {
            flagsPresent = 0x4;
          }

          if (samples[0].compositionTimeOffset !== undefined) {
            compositionTimeOffset = 0x8;
          }
        }

        return [0x00, // version 0
        0x00, durationPresent | sizePresent | flagsPresent | compositionTimeOffset, 0x01, // flags
        (samples.length & 0xFF000000) >>> 24, (samples.length & 0xFF0000) >>> 16, (samples.length & 0xFF00) >>> 8, samples.length & 0xFF, // sample_count
        (offset & 0xFF000000) >>> 24, (offset & 0xFF0000) >>> 16, (offset & 0xFF00) >>> 8, offset & 0xFF // data_offset
        ];
      };

      videoTrun = function videoTrun(track, offset) {
        var bytesOffest, bytes, header, samples, sample, i;
        samples = track.samples || [];
        offset += 8 + 12 + 16 * samples.length;
        header = trunHeader(samples, offset);
        bytes = new Uint8Array(header.length + samples.length * 16);
        bytes.set(header);
        bytesOffest = header.length;

        for (i = 0; i < samples.length; i++) {
          sample = samples[i];
          bytes[bytesOffest++] = (sample.duration & 0xFF000000) >>> 24;
          bytes[bytesOffest++] = (sample.duration & 0xFF0000) >>> 16;
          bytes[bytesOffest++] = (sample.duration & 0xFF00) >>> 8;
          bytes[bytesOffest++] = sample.duration & 0xFF; // sample_duration

          bytes[bytesOffest++] = (sample.size & 0xFF000000) >>> 24;
          bytes[bytesOffest++] = (sample.size & 0xFF0000) >>> 16;
          bytes[bytesOffest++] = (sample.size & 0xFF00) >>> 8;
          bytes[bytesOffest++] = sample.size & 0xFF; // sample_size

          bytes[bytesOffest++] = sample.flags.isLeading << 2 | sample.flags.dependsOn;
          bytes[bytesOffest++] = sample.flags.isDependedOn << 6 | sample.flags.hasRedundancy << 4 | sample.flags.paddingValue << 1 | sample.flags.isNonSyncSample;
          bytes[bytesOffest++] = sample.flags.degradationPriority & 0xF0 << 8;
          bytes[bytesOffest++] = sample.flags.degradationPriority & 0x0F; // sample_flags

          bytes[bytesOffest++] = (sample.compositionTimeOffset & 0xFF000000) >>> 24;
          bytes[bytesOffest++] = (sample.compositionTimeOffset & 0xFF0000) >>> 16;
          bytes[bytesOffest++] = (sample.compositionTimeOffset & 0xFF00) >>> 8;
          bytes[bytesOffest++] = sample.compositionTimeOffset & 0xFF; // sample_composition_time_offset
        }

        return box(types.trun, bytes);
      };

      audioTrun = function audioTrun(track, offset) {
        var bytes, bytesOffest, header, samples, sample, i;
        samples = track.samples || [];
        offset += 8 + 12 + 8 * samples.length;
        header = trunHeader(samples, offset);
        bytes = new Uint8Array(header.length + samples.length * 8);
        bytes.set(header);
        bytesOffest = header.length;

        for (i = 0; i < samples.length; i++) {
          sample = samples[i];
          bytes[bytesOffest++] = (sample.duration & 0xFF000000) >>> 24;
          bytes[bytesOffest++] = (sample.duration & 0xFF0000) >>> 16;
          bytes[bytesOffest++] = (sample.duration & 0xFF00) >>> 8;
          bytes[bytesOffest++] = sample.duration & 0xFF; // sample_duration

          bytes[bytesOffest++] = (sample.size & 0xFF000000) >>> 24;
          bytes[bytesOffest++] = (sample.size & 0xFF0000) >>> 16;
          bytes[bytesOffest++] = (sample.size & 0xFF00) >>> 8;
          bytes[bytesOffest++] = sample.size & 0xFF; // sample_size
        }

        return box(types.trun, bytes);
      };

      trun = function trun(track, offset) {
        if (track.type === 'audio') {
          return audioTrun(track, offset);
        }

        return videoTrun(track, offset);
      };
    })();

    var mp4Generator = {
      ftyp: ftyp,
      mdat: mdat,
      moof: moof,
      moov: moov,
      initSegment: function initSegment(tracks) {
        var fileType = ftyp(),
            movie = moov(tracks),
            result;
        result = new Uint8Array(fileType.byteLength + movie.byteLength);
        result.set(fileType);
        result.set(movie, fileType.byteLength);
        return result;
      }
    };
    /**
     * mux.js
     *
     * Copyright (c) Brightcove
     * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
     */
    // Convert an array of nal units into an array of frames with each frame being
    // composed of the nal units that make up that frame
    // Also keep track of cummulative data about the frame from the nal units such
    // as the frame duration, starting pts, etc.

    var groupNalsIntoFrames = function groupNalsIntoFrames(nalUnits) {
      var i,
          currentNal,
          currentFrame = [],
          frames = []; // TODO added for LHLS, make sure this is OK

      frames.byteLength = 0;
      frames.nalCount = 0;
      frames.duration = 0;
      currentFrame.byteLength = 0;

      for (i = 0; i < nalUnits.length; i++) {
        currentNal = nalUnits[i]; // Split on 'aud'-type nal units

        if (currentNal.nalUnitType === 'access_unit_delimiter_rbsp') {
          // Since the very first nal unit is expected to be an AUD
          // only push to the frames array when currentFrame is not empty
          if (currentFrame.length) {
            currentFrame.duration = currentNal.dts - currentFrame.dts; // TODO added for LHLS, make sure this is OK

            frames.byteLength += currentFrame.byteLength;
            frames.nalCount += currentFrame.length;
            frames.duration += currentFrame.duration;
            frames.push(currentFrame);
          }

          currentFrame = [currentNal];
          currentFrame.byteLength = currentNal.data.byteLength;
          currentFrame.pts = currentNal.pts;
          currentFrame.dts = currentNal.dts;
        } else {
          // Specifically flag key frames for ease of use later
          if (currentNal.nalUnitType === 'slice_layer_without_partitioning_rbsp_idr') {
            currentFrame.keyFrame = true;
          }

          currentFrame.duration = currentNal.dts - currentFrame.dts;
          currentFrame.byteLength += currentNal.data.byteLength;
          currentFrame.push(currentNal);
        }
      } // For the last frame, use the duration of the previous frame if we
      // have nothing better to go on


      if (frames.length && (!currentFrame.duration || currentFrame.duration <= 0)) {
        currentFrame.duration = frames[frames.length - 1].duration;
      } // Push the final frame
      // TODO added for LHLS, make sure this is OK


      frames.byteLength += currentFrame.byteLength;
      frames.nalCount += currentFrame.length;
      frames.duration += currentFrame.duration;
      frames.push(currentFrame);
      return frames;
    }; // Convert an array of frames into an array of Gop with each Gop being composed
    // of the frames that make up that Gop
    // Also keep track of cummulative data about the Gop from the frames such as the
    // Gop duration, starting pts, etc.


    var groupFramesIntoGops = function groupFramesIntoGops(frames) {
      var i,
          currentFrame,
          currentGop = [],
          gops = []; // We must pre-set some of the values on the Gop since we
      // keep running totals of these values

      currentGop.byteLength = 0;
      currentGop.nalCount = 0;
      currentGop.duration = 0;
      currentGop.pts = frames[0].pts;
      currentGop.dts = frames[0].dts; // store some metadata about all the Gops

      gops.byteLength = 0;
      gops.nalCount = 0;
      gops.duration = 0;
      gops.pts = frames[0].pts;
      gops.dts = frames[0].dts;

      for (i = 0; i < frames.length; i++) {
        currentFrame = frames[i];

        if (currentFrame.keyFrame) {
          // Since the very first frame is expected to be an keyframe
          // only push to the gops array when currentGop is not empty
          if (currentGop.length) {
            gops.push(currentGop);
            gops.byteLength += currentGop.byteLength;
            gops.nalCount += currentGop.nalCount;
            gops.duration += currentGop.duration;
          }

          currentGop = [currentFrame];
          currentGop.nalCount = currentFrame.length;
          currentGop.byteLength = currentFrame.byteLength;
          currentGop.pts = currentFrame.pts;
          currentGop.dts = currentFrame.dts;
          currentGop.duration = currentFrame.duration;
        } else {
          currentGop.duration += currentFrame.duration;
          currentGop.nalCount += currentFrame.length;
          currentGop.byteLength += currentFrame.byteLength;
          currentGop.push(currentFrame);
        }
      }

      if (gops.length && currentGop.duration <= 0) {
        currentGop.duration = gops[gops.length - 1].duration;
      }

      gops.byteLength += currentGop.byteLength;
      gops.nalCount += currentGop.nalCount;
      gops.duration += currentGop.duration; // push the final Gop

      gops.push(currentGop);
      return gops;
    };
    /*
     * Search for the first keyframe in the GOPs and throw away all frames
     * until that keyframe. Then extend the duration of the pulled keyframe
     * and pull the PTS and DTS of the keyframe so that it covers the time
     * range of the frames that were disposed.
     *
     * @param {Array} gops video GOPs
     * @returns {Array} modified video GOPs
     */


    var extendFirstKeyFrame = function extendFirstKeyFrame(gops) {
      var currentGop;

      if (!gops[0][0].keyFrame && gops.length > 1) {
        // Remove the first GOP
        currentGop = gops.shift();
        gops.byteLength -= currentGop.byteLength;
        gops.nalCount -= currentGop.nalCount; // Extend the first frame of what is now the
        // first gop to cover the time period of the
        // frames we just removed

        gops[0][0].dts = currentGop.dts;
        gops[0][0].pts = currentGop.pts;
        gops[0][0].duration += currentGop.duration;
      }

      return gops;
    };
    /**
     * Default sample object
     * see ISO/IEC 14496-12:2012, section 8.6.4.3
     */


    var createDefaultSample = function createDefaultSample() {
      return {
        size: 0,
        flags: {
          isLeading: 0,
          dependsOn: 1,
          isDependedOn: 0,
          hasRedundancy: 0,
          degradationPriority: 0,
          isNonSyncSample: 1
        }
      };
    };
    /*
     * Collates information from a video frame into an object for eventual
     * entry into an MP4 sample table.
     *
     * @param {Object} frame the video frame
     * @param {Number} dataOffset the byte offset to position the sample
     * @return {Object} object containing sample table info for a frame
     */


    var sampleForFrame = function sampleForFrame(frame, dataOffset) {
      var sample = createDefaultSample();
      sample.dataOffset = dataOffset;
      sample.compositionTimeOffset = frame.pts - frame.dts;
      sample.duration = frame.duration;
      sample.size = 4 * frame.length; // Space for nal unit size

      sample.size += frame.byteLength;

      if (frame.keyFrame) {
        sample.flags.dependsOn = 2;
        sample.flags.isNonSyncSample = 0;
      }

      return sample;
    }; // generate the track's sample table from an array of gops


    var generateSampleTable = function generateSampleTable(gops, baseDataOffset) {
      var h,
          i,
          sample,
          currentGop,
          currentFrame,
          dataOffset = baseDataOffset || 0,
          samples = [];

      for (h = 0; h < gops.length; h++) {
        currentGop = gops[h];

        for (i = 0; i < currentGop.length; i++) {
          currentFrame = currentGop[i];
          sample = sampleForFrame(currentFrame, dataOffset);
          dataOffset += sample.size;
          samples.push(sample);
        }
      }

      return samples;
    }; // generate the track's raw mdat data from an array of gops


    var concatenateNalData = function concatenateNalData(gops) {
      var h,
          i,
          j,
          currentGop,
          currentFrame,
          currentNal,
          dataOffset = 0,
          nalsByteLength = gops.byteLength,
          numberOfNals = gops.nalCount,
          totalByteLength = nalsByteLength + 4 * numberOfNals,
          data = new Uint8Array(totalByteLength),
          view = new DataView(data.buffer); // For each Gop..

      for (h = 0; h < gops.length; h++) {
        currentGop = gops[h]; // For each Frame..

        for (i = 0; i < currentGop.length; i++) {
          currentFrame = currentGop[i]; // For each NAL..

          for (j = 0; j < currentFrame.length; j++) {
            currentNal = currentFrame[j];
            view.setUint32(dataOffset, currentNal.data.byteLength);
            dataOffset += 4;
            data.set(currentNal.data, dataOffset);
            dataOffset += currentNal.data.byteLength;
          }
        }
      }

      return data;
    }; // generate the track's sample table from a frame


    var generateSampleTableForFrame = function generateSampleTableForFrame(frame, baseDataOffset) {
      var sample,
          dataOffset = baseDataOffset || 0,
          samples = [];
      sample = sampleForFrame(frame, dataOffset);
      samples.push(sample);
      return samples;
    }; // generate the track's raw mdat data from a frame


    var concatenateNalDataForFrame = function concatenateNalDataForFrame(frame) {
      var i,
          currentNal,
          dataOffset = 0,
          nalsByteLength = frame.byteLength,
          numberOfNals = frame.length,
          totalByteLength = nalsByteLength + 4 * numberOfNals,
          data = new Uint8Array(totalByteLength),
          view = new DataView(data.buffer); // For each NAL..

      for (i = 0; i < frame.length; i++) {
        currentNal = frame[i];
        view.setUint32(dataOffset, currentNal.data.byteLength);
        dataOffset += 4;
        data.set(currentNal.data, dataOffset);
        dataOffset += currentNal.data.byteLength;
      }

      return data;
    };

    var frameUtils = {
      groupNalsIntoFrames: groupNalsIntoFrames,
      groupFramesIntoGops: groupFramesIntoGops,
      extendFirstKeyFrame: extendFirstKeyFrame,
      generateSampleTable: generateSampleTable,
      concatenateNalData: concatenateNalData,
      generateSampleTableForFrame: generateSampleTableForFrame,
      concatenateNalDataForFrame: concatenateNalDataForFrame
    };
    /**
     * mux.js
     *
     * Copyright (c) Brightcove
     * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
     */

    var highPrefix = [33, 16, 5, 32, 164, 27];
    var lowPrefix = [33, 65, 108, 84, 1, 2, 4, 8, 168, 2, 4, 8, 17, 191, 252];

    var zeroFill = function zeroFill(count) {
      var a = [];

      while (count--) {
        a.push(0);
      }

      return a;
    };

    var makeTable = function makeTable(metaTable) {
      return Object.keys(metaTable).reduce(function (obj, key) {
        obj[key] = new Uint8Array(metaTable[key].reduce(function (arr, part) {
          return arr.concat(part);
        }, []));
        return obj;
      }, {});
    };

    var silence;

    var silence_1 = function silence_1() {
      if (!silence) {
        // Frames-of-silence to use for filling in missing AAC frames
        var coneOfSilence = {
          96000: [highPrefix, [227, 64], zeroFill(154), [56]],
          88200: [highPrefix, [231], zeroFill(170), [56]],
          64000: [highPrefix, [248, 192], zeroFill(240), [56]],
          48000: [highPrefix, [255, 192], zeroFill(268), [55, 148, 128], zeroFill(54), [112]],
          44100: [highPrefix, [255, 192], zeroFill(268), [55, 163, 128], zeroFill(84), [112]],
          32000: [highPrefix, [255, 192], zeroFill(268), [55, 234], zeroFill(226), [112]],
          24000: [highPrefix, [255, 192], zeroFill(268), [55, 255, 128], zeroFill(268), [111, 112], zeroFill(126), [224]],
          16000: [highPrefix, [255, 192], zeroFill(268), [55, 255, 128], zeroFill(268), [111, 255], zeroFill(269), [223, 108], zeroFill(195), [1, 192]],
          12000: [lowPrefix, zeroFill(268), [3, 127, 248], zeroFill(268), [6, 255, 240], zeroFill(268), [13, 255, 224], zeroFill(268), [27, 253, 128], zeroFill(259), [56]],
          11025: [lowPrefix, zeroFill(268), [3, 127, 248], zeroFill(268), [6, 255, 240], zeroFill(268), [13, 255, 224], zeroFill(268), [27, 255, 192], zeroFill(268), [55, 175, 128], zeroFill(108), [112]],
          8000: [lowPrefix, zeroFill(268), [3, 121, 16], zeroFill(47), [7]]
        };
        silence = makeTable(coneOfSilence);
      }

      return silence;
    };
    /**
     * mux.js
     *
     * Copyright (c) Brightcove
     * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
     */


    var ONE_SECOND_IN_TS = 90000,
        // 90kHz clock
    secondsToVideoTs,
        secondsToAudioTs,
        videoTsToSeconds,
        audioTsToSeconds,
        audioTsToVideoTs,
        videoTsToAudioTs,
        metadataTsToSeconds;

    secondsToVideoTs = function secondsToVideoTs(seconds) {
      return seconds * ONE_SECOND_IN_TS;
    };

    secondsToAudioTs = function secondsToAudioTs(seconds, sampleRate) {
      return seconds * sampleRate;
    };

    videoTsToSeconds = function videoTsToSeconds(timestamp) {
      return timestamp / ONE_SECOND_IN_TS;
    };

    audioTsToSeconds = function audioTsToSeconds(timestamp, sampleRate) {
      return timestamp / sampleRate;
    };

    audioTsToVideoTs = function audioTsToVideoTs(timestamp, sampleRate) {
      return secondsToVideoTs(audioTsToSeconds(timestamp, sampleRate));
    };

    videoTsToAudioTs = function videoTsToAudioTs(timestamp, sampleRate) {
      return secondsToAudioTs(videoTsToSeconds(timestamp), sampleRate);
    };
    /**
     * Adjust ID3 tag or caption timing information by the timeline pts values
     * (if keepOriginalTimestamps is false) and convert to seconds
     */


    metadataTsToSeconds = function metadataTsToSeconds(timestamp, timelineStartPts, keepOriginalTimestamps) {
      return videoTsToSeconds(keepOriginalTimestamps ? timestamp : timestamp - timelineStartPts);
    };

    var clock = {
      ONE_SECOND_IN_TS: ONE_SECOND_IN_TS,
      secondsToVideoTs: secondsToVideoTs,
      secondsToAudioTs: secondsToAudioTs,
      videoTsToSeconds: videoTsToSeconds,
      audioTsToSeconds: audioTsToSeconds,
      audioTsToVideoTs: audioTsToVideoTs,
      videoTsToAudioTs: videoTsToAudioTs,
      metadataTsToSeconds: metadataTsToSeconds
    };
    var clock_2 = clock.secondsToVideoTs;
    var clock_4 = clock.videoTsToSeconds;
    /**
     * mux.js
     *
     * Copyright (c) Brightcove
     * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
     */

    /**
     * Sum the `byteLength` properties of the data in each AAC frame
     */

    var sumFrameByteLengths = function sumFrameByteLengths(array) {
      var i,
          currentObj,
          sum = 0; // sum the byteLength's all each nal unit in the frame

      for (i = 0; i < array.length; i++) {
        currentObj = array[i];
        sum += currentObj.data.byteLength;
      }

      return sum;
    }; // Possibly pad (prefix) the audio track with silence if appending this track
    // would lead to the introduction of a gap in the audio buffer


    var prefixWithSilence = function prefixWithSilence(track, frames, audioAppendStartTs, videoBaseMediaDecodeTime) {
      var baseMediaDecodeTimeTs,
          frameDuration = 0,
          audioGapDuration = 0,
          audioFillFrameCount = 0,
          audioFillDuration = 0,
          silentFrame,
          i,
          firstFrame;

      if (!frames.length) {
        return;
      }

      baseMediaDecodeTimeTs = clock.audioTsToVideoTs(track.baseMediaDecodeTime, track.samplerate); // determine frame clock duration based on sample rate, round up to avoid overfills

      frameDuration = Math.ceil(clock.ONE_SECOND_IN_TS / (track.samplerate / 1024));

      if (audioAppendStartTs && videoBaseMediaDecodeTime) {
        // insert the shortest possible amount (audio gap or audio to video gap)
        audioGapDuration = baseMediaDecodeTimeTs - Math.max(audioAppendStartTs, videoBaseMediaDecodeTime); // number of full frames in the audio gap

        audioFillFrameCount = Math.floor(audioGapDuration / frameDuration);
        audioFillDuration = audioFillFrameCount * frameDuration;
      } // don't attempt to fill gaps smaller than a single frame or larger
      // than a half second


      if (audioFillFrameCount < 1 || audioFillDuration > clock.ONE_SECOND_IN_TS / 2) {
        return;
      }

      silentFrame = silence_1()[track.samplerate];

      if (!silentFrame) {
        // we don't have a silent frame pregenerated for the sample rate, so use a frame
        // from the content instead
        silentFrame = frames[0].data;
      }

      for (i = 0; i < audioFillFrameCount; i++) {
        firstFrame = frames[0];
        frames.splice(0, 0, {
          data: silentFrame,
          dts: firstFrame.dts - frameDuration,
          pts: firstFrame.pts - frameDuration
        });
      }

      track.baseMediaDecodeTime -= Math.floor(clock.videoTsToAudioTs(audioFillDuration, track.samplerate));
    }; // If the audio segment extends before the earliest allowed dts
    // value, remove AAC frames until starts at or after the earliest
    // allowed DTS so that we don't end up with a negative baseMedia-
    // DecodeTime for the audio track


    var trimAdtsFramesByEarliestDts = function trimAdtsFramesByEarliestDts(adtsFrames, track, earliestAllowedDts) {
      if (track.minSegmentDts >= earliestAllowedDts) {
        return adtsFrames;
      } // We will need to recalculate the earliest segment Dts


      track.minSegmentDts = Infinity;
      return adtsFrames.filter(function (currentFrame) {
        // If this is an allowed frame, keep it and record it's Dts
        if (currentFrame.dts >= earliestAllowedDts) {
          track.minSegmentDts = Math.min(track.minSegmentDts, currentFrame.dts);
          track.minSegmentPts = track.minSegmentDts;
          return true;
        } // Otherwise, discard it


        return false;
      });
    }; // generate the track's raw mdat data from an array of frames


    var generateSampleTable$1 = function generateSampleTable$1(frames) {
      var i,
          currentFrame,
          samples = [];

      for (i = 0; i < frames.length; i++) {
        currentFrame = frames[i];
        samples.push({
          size: currentFrame.data.byteLength,
          duration: 1024 // For AAC audio, all samples contain 1024 samples

        });
      }

      return samples;
    }; // generate the track's sample table from an array of frames


    var concatenateFrameData = function concatenateFrameData(frames) {
      var i,
          currentFrame,
          dataOffset = 0,
          data = new Uint8Array(sumFrameByteLengths(frames));

      for (i = 0; i < frames.length; i++) {
        currentFrame = frames[i];
        data.set(currentFrame.data, dataOffset);
        dataOffset += currentFrame.data.byteLength;
      }

      return data;
    };

    var audioFrameUtils = {
      prefixWithSilence: prefixWithSilence,
      trimAdtsFramesByEarliestDts: trimAdtsFramesByEarliestDts,
      generateSampleTable: generateSampleTable$1,
      concatenateFrameData: concatenateFrameData
    };
    /**
     * mux.js
     *
     * Copyright (c) Brightcove
     * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
     */

    var ONE_SECOND_IN_TS$1 = clock.ONE_SECOND_IN_TS;
    /**
     * Store information about the start and end of the track and the
     * duration for each frame/sample we process in order to calculate
     * the baseMediaDecodeTime
     */

    var collectDtsInfo = function collectDtsInfo(track, data) {
      if (typeof data.pts === 'number') {
        if (track.timelineStartInfo.pts === undefined) {
          track.timelineStartInfo.pts = data.pts;
        }

        if (track.minSegmentPts === undefined) {
          track.minSegmentPts = data.pts;
        } else {
          track.minSegmentPts = Math.min(track.minSegmentPts, data.pts);
        }

        if (track.maxSegmentPts === undefined) {
          track.maxSegmentPts = data.pts;
        } else {
          track.maxSegmentPts = Math.max(track.maxSegmentPts, data.pts);
        }
      }

      if (typeof data.dts === 'number') {
        if (track.timelineStartInfo.dts === undefined) {
          track.timelineStartInfo.dts = data.dts;
        }

        if (track.minSegmentDts === undefined) {
          track.minSegmentDts = data.dts;
        } else {
          track.minSegmentDts = Math.min(track.minSegmentDts, data.dts);
        }

        if (track.maxSegmentDts === undefined) {
          track.maxSegmentDts = data.dts;
        } else {
          track.maxSegmentDts = Math.max(track.maxSegmentDts, data.dts);
        }
      }
    };
    /**
     * Clear values used to calculate the baseMediaDecodeTime between
     * tracks
     */


    var clearDtsInfo = function clearDtsInfo(track) {
      delete track.minSegmentDts;
      delete track.maxSegmentDts;
      delete track.minSegmentPts;
      delete track.maxSegmentPts;
    };
    /**
     * Calculate the track's baseMediaDecodeTime based on the earliest
     * DTS the transmuxer has ever seen and the minimum DTS for the
     * current track
     * @param track {object} track metadata configuration
     * @param keepOriginalTimestamps {boolean} If true, keep the timestamps
     *        in the source; false to adjust the first segment to start at 0.
     */


    var calculateTrackBaseMediaDecodeTime = function calculateTrackBaseMediaDecodeTime(track, keepOriginalTimestamps) {
      var baseMediaDecodeTime,
          scale,
          minSegmentDts = track.minSegmentDts; // Optionally adjust the time so the first segment starts at zero.

      if (!keepOriginalTimestamps) {
        minSegmentDts -= track.timelineStartInfo.dts;
      } // track.timelineStartInfo.baseMediaDecodeTime is the location, in time, where
      // we want the start of the first segment to be placed


      baseMediaDecodeTime = track.timelineStartInfo.baseMediaDecodeTime; // Add to that the distance this segment is from the very first

      baseMediaDecodeTime += minSegmentDts; // baseMediaDecodeTime must not become negative

      baseMediaDecodeTime = Math.max(0, baseMediaDecodeTime);

      if (track.type === 'audio') {
        // Audio has a different clock equal to the sampling_rate so we need to
        // scale the PTS values into the clock rate of the track
        scale = track.samplerate / ONE_SECOND_IN_TS$1;
        baseMediaDecodeTime *= scale;
        baseMediaDecodeTime = Math.floor(baseMediaDecodeTime);
      }

      return baseMediaDecodeTime;
    };

    var trackDecodeInfo = {
      clearDtsInfo: clearDtsInfo,
      calculateTrackBaseMediaDecodeTime: calculateTrackBaseMediaDecodeTime,
      collectDtsInfo: collectDtsInfo
    };
    /**
     * mux.js
     *
     * Copyright (c) Brightcove
     * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
     *
     * Reads in-band caption information from a video elementary
     * stream. Captions must follow the CEA-708 standard for injection
     * into an MPEG-2 transport streams.
     * @see https://en.wikipedia.org/wiki/CEA-708
     * @see https://www.gpo.gov/fdsys/pkg/CFR-2007-title47-vol1/pdf/CFR-2007-title47-vol1-sec15-119.pdf
     */
    // Supplemental enhancement information (SEI) NAL units have a
    // payload type field to indicate how they are to be
    // interpreted. CEAS-708 caption content is always transmitted with
    // payload type 0x04.

    var USER_DATA_REGISTERED_ITU_T_T35 = 4,
        RBSP_TRAILING_BITS = 128;
    /**
      * Parse a supplemental enhancement information (SEI) NAL unit.
      * Stops parsing once a message of type ITU T T35 has been found.
      *
      * @param bytes {Uint8Array} the bytes of a SEI NAL unit
      * @return {object} the parsed SEI payload
      * @see Rec. ITU-T H.264, 7.3.2.3.1
      */

    var parseSei = function parseSei(bytes) {
      var i = 0,
          result = {
        payloadType: -1,
        payloadSize: 0
      },
          payloadType = 0,
          payloadSize = 0; // go through the sei_rbsp parsing each each individual sei_message

      while (i < bytes.byteLength) {
        // stop once we have hit the end of the sei_rbsp
        if (bytes[i] === RBSP_TRAILING_BITS) {
          break;
        } // Parse payload type


        while (bytes[i] === 0xFF) {
          payloadType += 255;
          i++;
        }

        payloadType += bytes[i++]; // Parse payload size

        while (bytes[i] === 0xFF) {
          payloadSize += 255;
          i++;
        }

        payloadSize += bytes[i++]; // this sei_message is a 608/708 caption so save it and break
        // there can only ever be one caption message in a frame's sei

        if (!result.payload && payloadType === USER_DATA_REGISTERED_ITU_T_T35) {
          var userIdentifier = String.fromCharCode(bytes[i + 3], bytes[i + 4], bytes[i + 5], bytes[i + 6]);

          if (userIdentifier === 'GA94') {
            result.payloadType = payloadType;
            result.payloadSize = payloadSize;
            result.payload = bytes.subarray(i, i + payloadSize);
            break;
          } else {
            result.payload = void 0;
          }
        } // skip the payload and parse the next message


        i += payloadSize;
        payloadType = 0;
        payloadSize = 0;
      }

      return result;
    }; // see ANSI/SCTE 128-1 (2013), section 8.1


    var parseUserData = function parseUserData(sei) {
      // itu_t_t35_contry_code must be 181 (United States) for
      // captions
      if (sei.payload[0] !== 181) {
        return null;
      } // itu_t_t35_provider_code should be 49 (ATSC) for captions


      if ((sei.payload[1] << 8 | sei.payload[2]) !== 49) {
        return null;
      } // the user_identifier should be "GA94" to indicate ATSC1 data


      if (String.fromCharCode(sei.payload[3], sei.payload[4], sei.payload[5], sei.payload[6]) !== 'GA94') {
        return null;
      } // finally, user_data_type_code should be 0x03 for caption data


      if (sei.payload[7] !== 0x03) {
        return null;
      } // return the user_data_type_structure and strip the trailing
      // marker bits


      return sei.payload.subarray(8, sei.payload.length - 1);
    }; // see CEA-708-D, section 4.4


    var parseCaptionPackets = function parseCaptionPackets(pts, userData) {
      var results = [],
          i,
          count,
          offset,
          data; // if this is just filler, return immediately

      if (!(userData[0] & 0x40)) {
        return results;
      } // parse out the cc_data_1 and cc_data_2 fields


      count = userData[0] & 0x1f;

      for (i = 0; i < count; i++) {
        offset = i * 3;
        data = {
          type: userData[offset + 2] & 0x03,
          pts: pts
        }; // capture cc data when cc_valid is 1

        if (userData[offset + 2] & 0x04) {
          data.ccData = userData[offset + 3] << 8 | userData[offset + 4];
          results.push(data);
        }
      }

      return results;
    };

    var discardEmulationPreventionBytes = function discardEmulationPreventionBytes(data) {
      var length = data.byteLength,
          emulationPreventionBytesPositions = [],
          i = 1,
          newLength,
          newData; // Find all `Emulation Prevention Bytes`

      while (i < length - 2) {
        if (data[i] === 0 && data[i + 1] === 0 && data[i + 2] === 0x03) {
          emulationPreventionBytesPositions.push(i + 2);
          i += 2;
        } else {
          i++;
        }
      } // If no Emulation Prevention Bytes were found just return the original
      // array


      if (emulationPreventionBytesPositions.length === 0) {
        return data;
      } // Create a new array to hold the NAL unit data


      newLength = length - emulationPreventionBytesPositions.length;
      newData = new Uint8Array(newLength);
      var sourceIndex = 0;

      for (i = 0; i < newLength; sourceIndex++, i++) {
        if (sourceIndex === emulationPreventionBytesPositions[0]) {
          // Skip this byte
          sourceIndex++; // Remove this position index

          emulationPreventionBytesPositions.shift();
        }

        newData[i] = data[sourceIndex];
      }

      return newData;
    }; // exports


    var captionPacketParser = {
      parseSei: parseSei,
      parseUserData: parseUserData,
      parseCaptionPackets: parseCaptionPackets,
      discardEmulationPreventionBytes: discardEmulationPreventionBytes,
      USER_DATA_REGISTERED_ITU_T_T35: USER_DATA_REGISTERED_ITU_T_T35
    }; // -----------------
    // Link To Transport
    // -----------------

    var CaptionStream = function CaptionStream() {
      CaptionStream.prototype.init.call(this);
      this.captionPackets_ = [];
      this.ccStreams_ = [new Cea608Stream(0, 0), // eslint-disable-line no-use-before-define
      new Cea608Stream(0, 1), // eslint-disable-line no-use-before-define
      new Cea608Stream(1, 0), // eslint-disable-line no-use-before-define
      new Cea608Stream(1, 1) // eslint-disable-line no-use-before-define
      ];
      this.reset(); // forward data and done events from CCs to this CaptionStream

      this.ccStreams_.forEach(function (cc) {
        cc.on('data', this.trigger.bind(this, 'data'));
        cc.on('partialdone', this.trigger.bind(this, 'partialdone'));
        cc.on('done', this.trigger.bind(this, 'done'));
      }, this);
    };

    CaptionStream.prototype = new stream();

    CaptionStream.prototype.push = function (event) {
      var sei, userData, newCaptionPackets; // only examine SEI NALs

      if (event.nalUnitType !== 'sei_rbsp') {
        return;
      } // parse the sei


      sei = captionPacketParser.parseSei(event.escapedRBSP); // ignore everything but user_data_registered_itu_t_t35

      if (sei.payloadType !== captionPacketParser.USER_DATA_REGISTERED_ITU_T_T35) {
        return;
      } // parse out the user data payload


      userData = captionPacketParser.parseUserData(sei); // ignore unrecognized userData

      if (!userData) {
        return;
      } // Sometimes, the same segment # will be downloaded twice. To stop the
      // caption data from being processed twice, we track the latest dts we've
      // received and ignore everything with a dts before that. However, since
      // data for a specific dts can be split across packets on either side of
      // a segment boundary, we need to make sure we *don't* ignore the packets
      // from the *next* segment that have dts === this.latestDts_. By constantly
      // tracking the number of packets received with dts === this.latestDts_, we
      // know how many should be ignored once we start receiving duplicates.


      if (event.dts < this.latestDts_) {
        // We've started getting older data, so set the flag.
        this.ignoreNextEqualDts_ = true;
        return;
      } else if (event.dts === this.latestDts_ && this.ignoreNextEqualDts_) {
        this.numSameDts_--;

        if (!this.numSameDts_) {
          // We've received the last duplicate packet, time to start processing again
          this.ignoreNextEqualDts_ = false;
        }

        return;
      } // parse out CC data packets and save them for later


      newCaptionPackets = captionPacketParser.parseCaptionPackets(event.pts, userData);
      this.captionPackets_ = this.captionPackets_.concat(newCaptionPackets);

      if (this.latestDts_ !== event.dts) {
        this.numSameDts_ = 0;
      }

      this.numSameDts_++;
      this.latestDts_ = event.dts;
    };

    CaptionStream.prototype.flushCCStreams = function (flushType) {
      this.ccStreams_.forEach(function (cc) {
        return flushType === 'flush' ? cc.flush() : cc.partialFlush();
      }, this);
    };

    CaptionStream.prototype.flushStream = function (flushType) {
      // make sure we actually parsed captions before proceeding
      if (!this.captionPackets_.length) {
        this.flushCCStreams(flushType);
        return;
      } // In Chrome, the Array#sort function is not stable so add a
      // presortIndex that we can use to ensure we get a stable-sort


      this.captionPackets_.forEach(function (elem, idx) {
        elem.presortIndex = idx;
      }); // sort caption byte-pairs based on their PTS values

      this.captionPackets_.sort(function (a, b) {
        if (a.pts === b.pts) {
          return a.presortIndex - b.presortIndex;
        }

        return a.pts - b.pts;
      });
      this.captionPackets_.forEach(function (packet) {
        if (packet.type < 2) {
          // Dispatch packet to the right Cea608Stream
          this.dispatchCea608Packet(packet);
        } // this is where an 'else' would go for a dispatching packets
        // to a theoretical Cea708Stream that handles SERVICEn data

      }, this);
      this.captionPackets_.length = 0;
      this.flushCCStreams(flushType);
    };

    CaptionStream.prototype.flush = function () {
      return this.flushStream('flush');
    }; // Only called if handling partial data


    CaptionStream.prototype.partialFlush = function () {
      return this.flushStream('partialFlush');
    };

    CaptionStream.prototype.reset = function () {
      this.latestDts_ = null;
      this.ignoreNextEqualDts_ = false;
      this.numSameDts_ = 0;
      this.activeCea608Channel_ = [null, null];
      this.ccStreams_.forEach(function (ccStream) {
        ccStream.reset();
      });
    }; // From the CEA-608 spec:

    /*
     * When XDS sub-packets are interleaved with other services, the end of each sub-packet shall be followed
     * by a control pair to change to a different service. When any of the control codes from 0x10 to 0x1F is
     * used to begin a control code pair, it indicates the return to captioning or Text data. The control code pair
     * and subsequent data should then be processed according to the FCC rules. It may be necessary for the
     * line 21 data encoder to automatically insert a control code pair (i.e. RCL, RU2, RU3, RU4, RDC, or RTD)
     * to switch to captioning or Text.
    */
    // With that in mind, we ignore any data between an XDS control code and a
    // subsequent closed-captioning control code.


    CaptionStream.prototype.dispatchCea608Packet = function (packet) {
      // NOTE: packet.type is the CEA608 field
      if (this.setsTextOrXDSActive(packet)) {
        this.activeCea608Channel_[packet.type] = null;
      } else if (this.setsChannel1Active(packet)) {
        this.activeCea608Channel_[packet.type] = 0;
      } else if (this.setsChannel2Active(packet)) {
        this.activeCea608Channel_[packet.type] = 1;
      }

      if (this.activeCea608Channel_[packet.type] === null) {
        // If we haven't received anything to set the active channel, or the
        // packets are Text/XDS data, discard the data; we don't want jumbled
        // captions
        return;
      }

      this.ccStreams_[(packet.type << 1) + this.activeCea608Channel_[packet.type]].push(packet);
    };

    CaptionStream.prototype.setsChannel1Active = function (packet) {
      return (packet.ccData & 0x7800) === 0x1000;
    };

    CaptionStream.prototype.setsChannel2Active = function (packet) {
      return (packet.ccData & 0x7800) === 0x1800;
    };

    CaptionStream.prototype.setsTextOrXDSActive = function (packet) {
      return (packet.ccData & 0x7100) === 0x0100 || (packet.ccData & 0x78fe) === 0x102a || (packet.ccData & 0x78fe) === 0x182a;
    }; // ----------------------
    // Session to Application
    // ----------------------
    // This hash maps non-ASCII, special, and extended character codes to their
    // proper Unicode equivalent. The first keys that are only a single byte
    // are the non-standard ASCII characters, which simply map the CEA608 byte
    // to the standard ASCII/Unicode. The two-byte keys that follow are the CEA608
    // character codes, but have their MSB bitmasked with 0x03 so that a lookup
    // can be performed regardless of the field and data channel on which the
    // character code was received.


    var CHARACTER_TRANSLATION = {
      0x2a: 0xe1,
      // á
      0x5c: 0xe9,
      // é
      0x5e: 0xed,
      // í
      0x5f: 0xf3,
      // ó
      0x60: 0xfa,
      // ú
      0x7b: 0xe7,
      // ç
      0x7c: 0xf7,
      // ÷
      0x7d: 0xd1,
      // Ñ
      0x7e: 0xf1,
      // ñ
      0x7f: 0x2588,
      // █
      0x0130: 0xae,
      // ®
      0x0131: 0xb0,
      // °
      0x0132: 0xbd,
      // ½
      0x0133: 0xbf,
      // ¿
      0x0134: 0x2122,
      // ™
      0x0135: 0xa2,
      // ¢
      0x0136: 0xa3,
      // £
      0x0137: 0x266a,
      // ♪
      0x0138: 0xe0,
      // à
      0x0139: 0xa0,
      //
      0x013a: 0xe8,
      // è
      0x013b: 0xe2,
      // â
      0x013c: 0xea,
      // ê
      0x013d: 0xee,
      // î
      0x013e: 0xf4,
      // ô
      0x013f: 0xfb,
      // û
      0x0220: 0xc1,
      // Á
      0x0221: 0xc9,
      // É
      0x0222: 0xd3,
      // Ó
      0x0223: 0xda,
      // Ú
      0x0224: 0xdc,
      // Ü
      0x0225: 0xfc,
      // ü
      0x0226: 0x2018,
      // ‘
      0x0227: 0xa1,
      // ¡
      0x0228: 0x2a,
      // *
      0x0229: 0x27,
      // '
      0x022a: 0x2014,
      // —
      0x022b: 0xa9,
      // ©
      0x022c: 0x2120,
      // ℠
      0x022d: 0x2022,
      // •
      0x022e: 0x201c,
      // “
      0x022f: 0x201d,
      // ”
      0x0230: 0xc0,
      // À
      0x0231: 0xc2,
      // Â
      0x0232: 0xc7,
      // Ç
      0x0233: 0xc8,
      // È
      0x0234: 0xca,
      // Ê
      0x0235: 0xcb,
      // Ë
      0x0236: 0xeb,
      // ë
      0x0237: 0xce,
      // Î
      0x0238: 0xcf,
      // Ï
      0x0239: 0xef,
      // ï
      0x023a: 0xd4,
      // Ô
      0x023b: 0xd9,
      // Ù
      0x023c: 0xf9,
      // ù
      0x023d: 0xdb,
      // Û
      0x023e: 0xab,
      // «
      0x023f: 0xbb,
      // »
      0x0320: 0xc3,
      // Ã
      0x0321: 0xe3,
      // ã
      0x0322: 0xcd,
      // Í
      0x0323: 0xcc,
      // Ì
      0x0324: 0xec,
      // ì
      0x0325: 0xd2,
      // Ò
      0x0326: 0xf2,
      // ò
      0x0327: 0xd5,
      // Õ
      0x0328: 0xf5,
      // õ
      0x0329: 0x7b,
      // {
      0x032a: 0x7d,
      // }
      0x032b: 0x5c,
      // \
      0x032c: 0x5e,
      // ^
      0x032d: 0x5f,
      // _
      0x032e: 0x7c,
      // |
      0x032f: 0x7e,
      // ~
      0x0330: 0xc4,
      // Ä
      0x0331: 0xe4,
      // ä
      0x0332: 0xd6,
      // Ö
      0x0333: 0xf6,
      // ö
      0x0334: 0xdf,
      // ß
      0x0335: 0xa5,
      // ¥
      0x0336: 0xa4,
      // ¤
      0x0337: 0x2502,
      // │
      0x0338: 0xc5,
      // Å
      0x0339: 0xe5,
      // å
      0x033a: 0xd8,
      // Ø
      0x033b: 0xf8,
      // ø
      0x033c: 0x250c,
      // ┌
      0x033d: 0x2510,
      // ┐
      0x033e: 0x2514,
      // └
      0x033f: 0x2518 // ┘

    };

    var getCharFromCode = function getCharFromCode(code) {
      if (code === null) {
        return '';
      }

      code = CHARACTER_TRANSLATION[code] || code;
      return String.fromCharCode(code);
    }; // the index of the last row in a CEA-608 display buffer


    var BOTTOM_ROW = 14; // This array is used for mapping PACs -> row #, since there's no way of
    // getting it through bit logic.

    var ROWS = [0x1100, 0x1120, 0x1200, 0x1220, 0x1500, 0x1520, 0x1600, 0x1620, 0x1700, 0x1720, 0x1000, 0x1300, 0x1320, 0x1400, 0x1420]; // CEA-608 captions are rendered onto a 34x15 matrix of character
    // cells. The "bottom" row is the last element in the outer array.

    var createDisplayBuffer = function createDisplayBuffer() {
      var result = [],
          i = BOTTOM_ROW + 1;

      while (i--) {
        result.push('');
      }

      return result;
    };

    var Cea608Stream = function Cea608Stream(field, dataChannel) {
      Cea608Stream.prototype.init.call(this);
      this.field_ = field || 0;
      this.dataChannel_ = dataChannel || 0;
      this.name_ = 'CC' + ((this.field_ << 1 | this.dataChannel_) + 1);
      this.setConstants();
      this.reset();

      this.push = function (packet) {
        var data, swap, char0, char1, text; // remove the parity bits

        data = packet.ccData & 0x7f7f; // ignore duplicate control codes; the spec demands they're sent twice

        if (data === this.lastControlCode_) {
          this.lastControlCode_ = null;
          return;
        } // Store control codes


        if ((data & 0xf000) === 0x1000) {
          this.lastControlCode_ = data;
        } else if (data !== this.PADDING_) {
          this.lastControlCode_ = null;
        }

        char0 = data >>> 8;
        char1 = data & 0xff;

        if (data === this.PADDING_) {
          return;
        } else if (data === this.RESUME_CAPTION_LOADING_) {
          this.mode_ = 'popOn';
        } else if (data === this.END_OF_CAPTION_) {
          // If an EOC is received while in paint-on mode, the displayed caption
          // text should be swapped to non-displayed memory as if it was a pop-on
          // caption. Because of that, we should explicitly switch back to pop-on
          // mode
          this.mode_ = 'popOn';
          this.clearFormatting(packet.pts); // if a caption was being displayed, it's gone now

          this.flushDisplayed(packet.pts); // flip memory

          swap = this.displayed_;
          this.displayed_ = this.nonDisplayed_;
          this.nonDisplayed_ = swap; // start measuring the time to display the caption

          this.startPts_ = packet.pts;
        } else if (data === this.ROLL_UP_2_ROWS_) {
          this.rollUpRows_ = 2;
          this.setRollUp(packet.pts);
        } else if (data === this.ROLL_UP_3_ROWS_) {
          this.rollUpRows_ = 3;
          this.setRollUp(packet.pts);
        } else if (data === this.ROLL_UP_4_ROWS_) {
          this.rollUpRows_ = 4;
          this.setRollUp(packet.pts);
        } else if (data === this.CARRIAGE_RETURN_) {
          this.clearFormatting(packet.pts);
          this.flushDisplayed(packet.pts);
          this.shiftRowsUp_();
          this.startPts_ = packet.pts;
        } else if (data === this.BACKSPACE_) {
          if (this.mode_ === 'popOn') {
            this.nonDisplayed_[this.row_] = this.nonDisplayed_[this.row_].slice(0, -1);
          } else {
            this.displayed_[this.row_] = this.displayed_[this.row_].slice(0, -1);
          }
        } else if (data === this.ERASE_DISPLAYED_MEMORY_) {
          this.flushDisplayed(packet.pts);
          this.displayed_ = createDisplayBuffer();
        } else if (data === this.ERASE_NON_DISPLAYED_MEMORY_) {
          this.nonDisplayed_ = createDisplayBuffer();
        } else if (data === this.RESUME_DIRECT_CAPTIONING_) {
          if (this.mode_ !== 'paintOn') {
            // NOTE: This should be removed when proper caption positioning is
            // implemented
            this.flushDisplayed(packet.pts);
            this.displayed_ = createDisplayBuffer();
          }

          this.mode_ = 'paintOn';
          this.startPts_ = packet.pts; // Append special characters to caption text
        } else if (this.isSpecialCharacter(char0, char1)) {
          // Bitmask char0 so that we can apply character transformations
          // regardless of field and data channel.
          // Then byte-shift to the left and OR with char1 so we can pass the
          // entire character code to `getCharFromCode`.
          char0 = (char0 & 0x03) << 8;
          text = getCharFromCode(char0 | char1);
          this[this.mode_](packet.pts, text);
          this.column_++; // Append extended characters to caption text
        } else if (this.isExtCharacter(char0, char1)) {
          // Extended characters always follow their "non-extended" equivalents.
          // IE if a "è" is desired, you'll always receive "eè"; non-compliant
          // decoders are supposed to drop the "è", while compliant decoders
          // backspace the "e" and insert "è".
          // Delete the previous character
          if (this.mode_ === 'popOn') {
            this.nonDisplayed_[this.row_] = this.nonDisplayed_[this.row_].slice(0, -1);
          } else {
            this.displayed_[this.row_] = this.displayed_[this.row_].slice(0, -1);
          } // Bitmask char0 so that we can apply character transformations
          // regardless of field and data channel.
          // Then byte-shift to the left and OR with char1 so we can pass the
          // entire character code to `getCharFromCode`.


          char0 = (char0 & 0x03) << 8;
          text = getCharFromCode(char0 | char1);
          this[this.mode_](packet.pts, text);
          this.column_++; // Process mid-row codes
        } else if (this.isMidRowCode(char0, char1)) {
          // Attributes are not additive, so clear all formatting
          this.clearFormatting(packet.pts); // According to the standard, mid-row codes
          // should be replaced with spaces, so add one now

          this[this.mode_](packet.pts, ' ');
          this.column_++;

          if ((char1 & 0xe) === 0xe) {
            this.addFormatting(packet.pts, ['i']);
          }

          if ((char1 & 0x1) === 0x1) {
            this.addFormatting(packet.pts, ['u']);
          } // Detect offset control codes and adjust cursor

        } else if (this.isOffsetControlCode(char0, char1)) {
          // Cursor position is set by indent PAC (see below) in 4-column
          // increments, with an additional offset code of 1-3 to reach any
          // of the 32 columns specified by CEA-608. So all we need to do
          // here is increment the column cursor by the given offset.
          this.column_ += char1 & 0x03; // Detect PACs (Preamble Address Codes)
        } else if (this.isPAC(char0, char1)) {
          // There's no logic for PAC -> row mapping, so we have to just
          // find the row code in an array and use its index :(
          var row = ROWS.indexOf(data & 0x1f20); // Configure the caption window if we're in roll-up mode

          if (this.mode_ === 'rollUp') {
            // This implies that the base row is incorrectly set.
            // As per the recommendation in CEA-608(Base Row Implementation), defer to the number
            // of roll-up rows set.
            if (row - this.rollUpRows_ + 1 < 0) {
              row = this.rollUpRows_ - 1;
            }

            this.setRollUp(packet.pts, row);
          }

          if (row !== this.row_) {
            // formatting is only persistent for current row
            this.clearFormatting(packet.pts);
            this.row_ = row;
          } // All PACs can apply underline, so detect and apply
          // (All odd-numbered second bytes set underline)


          if (char1 & 0x1 && this.formatting_.indexOf('u') === -1) {
            this.addFormatting(packet.pts, ['u']);
          }

          if ((data & 0x10) === 0x10) {
            // We've got an indent level code. Each successive even number
            // increments the column cursor by 4, so we can get the desired
            // column position by bit-shifting to the right (to get n/2)
            // and multiplying by 4.
            this.column_ = ((data & 0xe) >> 1) * 4;
          }

          if (this.isColorPAC(char1)) {
            // it's a color code, though we only support white, which
            // can be either normal or italicized. white italics can be
            // either 0x4e or 0x6e depending on the row, so we just
            // bitwise-and with 0xe to see if italics should be turned on
            if ((char1 & 0xe) === 0xe) {
              this.addFormatting(packet.pts, ['i']);
            }
          } // We have a normal character in char0, and possibly one in char1

        } else if (this.isNormalChar(char0)) {
          if (char1 === 0x00) {
            char1 = null;
          }

          text = getCharFromCode(char0);
          text += getCharFromCode(char1);
          this[this.mode_](packet.pts, text);
          this.column_ += text.length;
        } // finish data processing

      };
    };

    Cea608Stream.prototype = new stream(); // Trigger a cue point that captures the current state of the
    // display buffer

    Cea608Stream.prototype.flushDisplayed = function (pts) {
      var content = this.displayed_ // remove spaces from the start and end of the string
      .map(function (row) {
        try {
          return row.trim();
        } catch (e) {
          // Ordinarily, this shouldn't happen. However, caption
          // parsing errors should not throw exceptions and
          // break playback.
          // eslint-disable-next-line no-console
          console.error('Skipping malformed caption.');
          return '';
        }
      }) // combine all text rows to display in one cue
      .join('\n') // and remove blank rows from the start and end, but not the middle
      .replace(/^\n+|\n+$/g, '');

      if (content.length) {
        this.trigger('data', {
          startPts: this.startPts_,
          endPts: pts,
          text: content,
          stream: this.name_
        });
      }
    };
    /**
     * Zero out the data, used for startup and on seek
     */


    Cea608Stream.prototype.reset = function () {
      this.mode_ = 'popOn'; // When in roll-up mode, the index of the last row that will
      // actually display captions. If a caption is shifted to a row
      // with a lower index than this, it is cleared from the display
      // buffer

      this.topRow_ = 0;
      this.startPts_ = 0;
      this.displayed_ = createDisplayBuffer();
      this.nonDisplayed_ = createDisplayBuffer();
      this.lastControlCode_ = null; // Track row and column for proper line-breaking and spacing

      this.column_ = 0;
      this.row_ = BOTTOM_ROW;
      this.rollUpRows_ = 2; // This variable holds currently-applied formatting

      this.formatting_ = [];
    };
    /**
     * Sets up control code and related constants for this instance
     */


    Cea608Stream.prototype.setConstants = function () {
      // The following attributes have these uses:
      // ext_ :    char0 for mid-row codes, and the base for extended
      //           chars (ext_+0, ext_+1, and ext_+2 are char0s for
      //           extended codes)
      // control_: char0 for control codes, except byte-shifted to the
      //           left so that we can do this.control_ | CONTROL_CODE
      // offset_:  char0 for tab offset codes
      //
      // It's also worth noting that control codes, and _only_ control codes,
      // differ between field 1 and field2. Field 2 control codes are always
      // their field 1 value plus 1. That's why there's the "| field" on the
      // control value.
      if (this.dataChannel_ === 0) {
        this.BASE_ = 0x10;
        this.EXT_ = 0x11;
        this.CONTROL_ = (0x14 | this.field_) << 8;
        this.OFFSET_ = 0x17;
      } else if (this.dataChannel_ === 1) {
        this.BASE_ = 0x18;
        this.EXT_ = 0x19;
        this.CONTROL_ = (0x1c | this.field_) << 8;
        this.OFFSET_ = 0x1f;
      } // Constants for the LSByte command codes recognized by Cea608Stream. This
      // list is not exhaustive. For a more comprehensive listing and semantics see
      // http://www.gpo.gov/fdsys/pkg/CFR-2010-title47-vol1/pdf/CFR-2010-title47-vol1-sec15-119.pdf
      // Padding


      this.PADDING_ = 0x0000; // Pop-on Mode

      this.RESUME_CAPTION_LOADING_ = this.CONTROL_ | 0x20;
      this.END_OF_CAPTION_ = this.CONTROL_ | 0x2f; // Roll-up Mode

      this.ROLL_UP_2_ROWS_ = this.CONTROL_ | 0x25;
      this.ROLL_UP_3_ROWS_ = this.CONTROL_ | 0x26;
      this.ROLL_UP_4_ROWS_ = this.CONTROL_ | 0x27;
      this.CARRIAGE_RETURN_ = this.CONTROL_ | 0x2d; // paint-on mode

      this.RESUME_DIRECT_CAPTIONING_ = this.CONTROL_ | 0x29; // Erasure

      this.BACKSPACE_ = this.CONTROL_ | 0x21;
      this.ERASE_DISPLAYED_MEMORY_ = this.CONTROL_ | 0x2c;
      this.ERASE_NON_DISPLAYED_MEMORY_ = this.CONTROL_ | 0x2e;
    };
    /**
     * Detects if the 2-byte packet data is a special character
     *
     * Special characters have a second byte in the range 0x30 to 0x3f,
     * with the first byte being 0x11 (for data channel 1) or 0x19 (for
     * data channel 2).
     *
     * @param  {Integer} char0 The first byte
     * @param  {Integer} char1 The second byte
     * @return {Boolean}       Whether the 2 bytes are an special character
     */


    Cea608Stream.prototype.isSpecialCharacter = function (char0, char1) {
      return char0 === this.EXT_ && char1 >= 0x30 && char1 <= 0x3f;
    };
    /**
     * Detects if the 2-byte packet data is an extended character
     *
     * Extended characters have a second byte in the range 0x20 to 0x3f,
     * with the first byte being 0x12 or 0x13 (for data channel 1) or
     * 0x1a or 0x1b (for data channel 2).
     *
     * @param  {Integer} char0 The first byte
     * @param  {Integer} char1 The second byte
     * @return {Boolean}       Whether the 2 bytes are an extended character
     */


    Cea608Stream.prototype.isExtCharacter = function (char0, char1) {
      return (char0 === this.EXT_ + 1 || char0 === this.EXT_ + 2) && char1 >= 0x20 && char1 <= 0x3f;
    };
    /**
     * Detects if the 2-byte packet is a mid-row code
     *
     * Mid-row codes have a second byte in the range 0x20 to 0x2f, with
     * the first byte being 0x11 (for data channel 1) or 0x19 (for data
     * channel 2).
     *
     * @param  {Integer} char0 The first byte
     * @param  {Integer} char1 The second byte
     * @return {Boolean}       Whether the 2 bytes are a mid-row code
     */


    Cea608Stream.prototype.isMidRowCode = function (char0, char1) {
      return char0 === this.EXT_ && char1 >= 0x20 && char1 <= 0x2f;
    };
    /**
     * Detects if the 2-byte packet is an offset control code
     *
     * Offset control codes have a second byte in the range 0x21 to 0x23,
     * with the first byte being 0x17 (for data channel 1) or 0x1f (for
     * data channel 2).
     *
     * @param  {Integer} char0 The first byte
     * @param  {Integer} char1 The second byte
     * @return {Boolean}       Whether the 2 bytes are an offset control code
     */


    Cea608Stream.prototype.isOffsetControlCode = function (char0, char1) {
      return char0 === this.OFFSET_ && char1 >= 0x21 && char1 <= 0x23;
    };
    /**
     * Detects if the 2-byte packet is a Preamble Address Code
     *
     * PACs have a first byte in the range 0x10 to 0x17 (for data channel 1)
     * or 0x18 to 0x1f (for data channel 2), with the second byte in the
     * range 0x40 to 0x7f.
     *
     * @param  {Integer} char0 The first byte
     * @param  {Integer} char1 The second byte
     * @return {Boolean}       Whether the 2 bytes are a PAC
     */


    Cea608Stream.prototype.isPAC = function (char0, char1) {
      return char0 >= this.BASE_ && char0 < this.BASE_ + 8 && char1 >= 0x40 && char1 <= 0x7f;
    };
    /**
     * Detects if a packet's second byte is in the range of a PAC color code
     *
     * PAC color codes have the second byte be in the range 0x40 to 0x4f, or
     * 0x60 to 0x6f.
     *
     * @param  {Integer} char1 The second byte
     * @return {Boolean}       Whether the byte is a color PAC
     */


    Cea608Stream.prototype.isColorPAC = function (char1) {
      return char1 >= 0x40 && char1 <= 0x4f || char1 >= 0x60 && char1 <= 0x7f;
    };
    /**
     * Detects if a single byte is in the range of a normal character
     *
     * Normal text bytes are in the range 0x20 to 0x7f.
     *
     * @param  {Integer} char  The byte
     * @return {Boolean}       Whether the byte is a normal character
     */


    Cea608Stream.prototype.isNormalChar = function (char) {
      return char >= 0x20 && char <= 0x7f;
    };
    /**
     * Configures roll-up
     *
     * @param  {Integer} pts         Current PTS
     * @param  {Integer} newBaseRow  Used by PACs to slide the current window to
     *                               a new position
     */


    Cea608Stream.prototype.setRollUp = function (pts, newBaseRow) {
      // Reset the base row to the bottom row when switching modes
      if (this.mode_ !== 'rollUp') {
        this.row_ = BOTTOM_ROW;
        this.mode_ = 'rollUp'; // Spec says to wipe memories when switching to roll-up

        this.flushDisplayed(pts);
        this.nonDisplayed_ = createDisplayBuffer();
        this.displayed_ = createDisplayBuffer();
      }

      if (newBaseRow !== undefined && newBaseRow !== this.row_) {
        // move currently displayed captions (up or down) to the new base row
        for (var i = 0; i < this.rollUpRows_; i++) {
          this.displayed_[newBaseRow - i] = this.displayed_[this.row_ - i];
          this.displayed_[this.row_ - i] = '';
        }
      }

      if (newBaseRow === undefined) {
        newBaseRow = this.row_;
      }

      this.topRow_ = newBaseRow - this.rollUpRows_ + 1;
    }; // Adds the opening HTML tag for the passed character to the caption text,
    // and keeps track of it for later closing


    Cea608Stream.prototype.addFormatting = function (pts, format) {
      this.formatting_ = this.formatting_.concat(format);
      var text = format.reduce(function (text, format) {
        return text + '<' + format + '>';
      }, '');
      this[this.mode_](pts, text);
    }; // Adds HTML closing tags for current formatting to caption text and
    // clears remembered formatting


    Cea608Stream.prototype.clearFormatting = function (pts) {
      if (!this.formatting_.length) {
        return;
      }

      var text = this.formatting_.reverse().reduce(function (text, format) {
        return text + '</' + format + '>';
      }, '');
      this.formatting_ = [];
      this[this.mode_](pts, text);
    }; // Mode Implementations


    Cea608Stream.prototype.popOn = function (pts, text) {
      var baseRow = this.nonDisplayed_[this.row_]; // buffer characters

      baseRow += text;
      this.nonDisplayed_[this.row_] = baseRow;
    };

    Cea608Stream.prototype.rollUp = function (pts, text) {
      var baseRow = this.displayed_[this.row_];
      baseRow += text;
      this.displayed_[this.row_] = baseRow;
    };

    Cea608Stream.prototype.shiftRowsUp_ = function () {
      var i; // clear out inactive rows

      for (i = 0; i < this.topRow_; i++) {
        this.displayed_[i] = '';
      }

      for (i = this.row_ + 1; i < BOTTOM_ROW + 1; i++) {
        this.displayed_[i] = '';
      } // shift displayed rows up


      for (i = this.topRow_; i < this.row_; i++) {
        this.displayed_[i] = this.displayed_[i + 1];
      } // clear out the bottom row


      this.displayed_[this.row_] = '';
    };

    Cea608Stream.prototype.paintOn = function (pts, text) {
      var baseRow = this.displayed_[this.row_];
      baseRow += text;
      this.displayed_[this.row_] = baseRow;
    }; // exports


    var captionStream = {
      CaptionStream: CaptionStream,
      Cea608Stream: Cea608Stream
    };
    /**
     * mux.js
     *
     * Copyright (c) Brightcove
     * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
     */

    var streamTypes = {
      H264_STREAM_TYPE: 0x1B,
      ADTS_STREAM_TYPE: 0x0F,
      METADATA_STREAM_TYPE: 0x15
    };
    var MAX_TS = 8589934592;
    var RO_THRESH = 4294967296;
    var TYPE_SHARED = 'shared';

    var handleRollover = function handleRollover(value, reference) {
      var direction = 1;

      if (value > reference) {
        // If the current timestamp value is greater than our reference timestamp and we detect a
        // timestamp rollover, this means the roll over is happening in the opposite direction.
        // Example scenario: Enter a long stream/video just after a rollover occurred. The reference
        // point will be set to a small number, e.g. 1. The user then seeks backwards over the
        // rollover point. In loading this segment, the timestamp values will be very large,
        // e.g. 2^33 - 1. Since this comes before the data we loaded previously, we want to adjust
        // the time stamp to be `value - 2^33`.
        direction = -1;
      } // Note: A seek forwards or back that is greater than the RO_THRESH (2^32, ~13 hours) will
      // cause an incorrect adjustment.


      while (Math.abs(reference - value) > RO_THRESH) {
        value += direction * MAX_TS;
      }

      return value;
    };

    var TimestampRolloverStream = function TimestampRolloverStream(type) {
      var lastDTS, referenceDTS;
      TimestampRolloverStream.prototype.init.call(this); // The "shared" type is used in cases where a stream will contain muxed
      // video and audio. We could use `undefined` here, but having a string
      // makes debugging a little clearer.

      this.type_ = type || TYPE_SHARED;

      this.push = function (data) {
        // Any "shared" rollover streams will accept _all_ data. Otherwise,
        // streams will only accept data that matches their type.
        if (this.type_ !== TYPE_SHARED && data.type !== this.type_) {
          return;
        }

        if (referenceDTS === undefined) {
          referenceDTS = data.dts;
        }

        data.dts = handleRollover(data.dts, referenceDTS);
        data.pts = handleRollover(data.pts, referenceDTS);
        lastDTS = data.dts;
        this.trigger('data', data);
      };

      this.flush = function () {
        referenceDTS = lastDTS;
        this.trigger('done');
      };

      this.endTimeline = function () {
        this.flush();
        this.trigger('endedtimeline');
      };

      this.discontinuity = function () {
        referenceDTS = void 0;
        lastDTS = void 0;
      };

      this.reset = function () {
        this.discontinuity();
        this.trigger('reset');
      };
    };

    TimestampRolloverStream.prototype = new stream();
    var timestampRolloverStream = {
      TimestampRolloverStream: TimestampRolloverStream,
      handleRollover: handleRollover
    };

    var percentEncode = function percentEncode(bytes, start, end) {
      var i,
          result = '';

      for (i = start; i < end; i++) {
        result += '%' + ('00' + bytes[i].toString(16)).slice(-2);
      }

      return result;
    },
        // return the string representation of the specified byte range,
    // interpreted as UTf-8.
    parseUtf8 = function parseUtf8(bytes, start, end) {
      return decodeURIComponent(percentEncode(bytes, start, end));
    },
        // return the string representation of the specified byte range,
    // interpreted as ISO-8859-1.
    parseIso88591 = function parseIso88591(bytes, start, end) {
      return unescape(percentEncode(bytes, start, end)); // jshint ignore:line
    },
        parseSyncSafeInteger = function parseSyncSafeInteger(data) {
      return data[0] << 21 | data[1] << 14 | data[2] << 7 | data[3];
    },
        tagParsers = {
      TXXX: function TXXX(tag) {
        var i;

        if (tag.data[0] !== 3) {
          // ignore frames with unrecognized character encodings
          return;
        }

        for (i = 1; i < tag.data.length; i++) {
          if (tag.data[i] === 0) {
            // parse the text fields
            tag.description = parseUtf8(tag.data, 1, i); // do not include the null terminator in the tag value

            tag.value = parseUtf8(tag.data, i + 1, tag.data.length).replace(/\0*$/, '');
            break;
          }
        }

        tag.data = tag.value;
      },
      WXXX: function WXXX(tag) {
        var i;

        if (tag.data[0] !== 3) {
          // ignore frames with unrecognized character encodings
          return;
        }

        for (i = 1; i < tag.data.length; i++) {
          if (tag.data[i] === 0) {
            // parse the description and URL fields
            tag.description = parseUtf8(tag.data, 1, i);
            tag.url = parseUtf8(tag.data, i + 1, tag.data.length);
            break;
          }
        }
      },
      PRIV: function PRIV(tag) {
        var i;

        for (i = 0; i < tag.data.length; i++) {
          if (tag.data[i] === 0) {
            // parse the description and URL fields
            tag.owner = parseIso88591(tag.data, 0, i);
            break;
          }
        }

        tag.privateData = tag.data.subarray(i + 1);
        tag.data = tag.privateData;
      }
    },
        _MetadataStream;

    _MetadataStream = function MetadataStream(options) {
      var settings = {
        debug: !!(options && options.debug),
        // the bytes of the program-level descriptor field in MP2T
        // see ISO/IEC 13818-1:2013 (E), section 2.6 "Program and
        // program element descriptors"
        descriptor: options && options.descriptor
      },
          // the total size in bytes of the ID3 tag being parsed
      tagSize = 0,
          // tag data that is not complete enough to be parsed
      buffer = [],
          // the total number of bytes currently in the buffer
      bufferSize = 0,
          i;

      _MetadataStream.prototype.init.call(this); // calculate the text track in-band metadata track dispatch type
      // https://html.spec.whatwg.org/multipage/embedded-content.html#steps-to-expose-a-media-resource-specific-text-track


      this.dispatchType = streamTypes.METADATA_STREAM_TYPE.toString(16);

      if (settings.descriptor) {
        for (i = 0; i < settings.descriptor.length; i++) {
          this.dispatchType += ('00' + settings.descriptor[i].toString(16)).slice(-2);
        }
      }

      this.push = function (chunk) {
        var tag, frameStart, frameSize, frame, i, frameHeader;

        if (chunk.type !== 'timed-metadata') {
          return;
        } // if data_alignment_indicator is set in the PES header,
        // we must have the start of a new ID3 tag. Assume anything
        // remaining in the buffer was malformed and throw it out


        if (chunk.dataAlignmentIndicator) {
          bufferSize = 0;
          buffer.length = 0;
        } // ignore events that don't look like ID3 data


        if (buffer.length === 0 && (chunk.data.length < 10 || chunk.data[0] !== 'I'.charCodeAt(0) || chunk.data[1] !== 'D'.charCodeAt(0) || chunk.data[2] !== '3'.charCodeAt(0))) {
          if (settings.debug) {
            // eslint-disable-next-line no-console
            console.log('Skipping unrecognized metadata packet');
          }

          return;
        } // add this chunk to the data we've collected so far


        buffer.push(chunk);
        bufferSize += chunk.data.byteLength; // grab the size of the entire frame from the ID3 header

        if (buffer.length === 1) {
          // the frame size is transmitted as a 28-bit integer in the
          // last four bytes of the ID3 header.
          // The most significant bit of each byte is dropped and the
          // results concatenated to recover the actual value.
          tagSize = parseSyncSafeInteger(chunk.data.subarray(6, 10)); // ID3 reports the tag size excluding the header but it's more
          // convenient for our comparisons to include it

          tagSize += 10;
        } // if the entire frame has not arrived, wait for more data


        if (bufferSize < tagSize) {
          return;
        } // collect the entire frame so it can be parsed


        tag = {
          data: new Uint8Array(tagSize),
          frames: [],
          pts: buffer[0].pts,
          dts: buffer[0].dts
        };

        for (i = 0; i < tagSize;) {
          tag.data.set(buffer[0].data.subarray(0, tagSize - i), i);
          i += buffer[0].data.byteLength;
          bufferSize -= buffer[0].data.byteLength;
          buffer.shift();
        } // find the start of the first frame and the end of the tag


        frameStart = 10;

        if (tag.data[5] & 0x40) {
          // advance the frame start past the extended header
          frameStart += 4; // header size field

          frameStart += parseSyncSafeInteger(tag.data.subarray(10, 14)); // clip any padding off the end

          tagSize -= parseSyncSafeInteger(tag.data.subarray(16, 20));
        } // parse one or more ID3 frames
        // http://id3.org/id3v2.3.0#ID3v2_frame_overview


        do {
          // determine the number of bytes in this frame
          frameSize = parseSyncSafeInteger(tag.data.subarray(frameStart + 4, frameStart + 8));

          if (frameSize < 1) {
            // eslint-disable-next-line no-console
            return console.log('Malformed ID3 frame encountered. Skipping metadata parsing.');
          }

          frameHeader = String.fromCharCode(tag.data[frameStart], tag.data[frameStart + 1], tag.data[frameStart + 2], tag.data[frameStart + 3]);
          frame = {
            id: frameHeader,
            data: tag.data.subarray(frameStart + 10, frameStart + frameSize + 10)
          };
          frame.key = frame.id;

          if (tagParsers[frame.id]) {
            tagParsers[frame.id](frame); // handle the special PRIV frame used to indicate the start
            // time for raw AAC data

            if (frame.owner === 'com.apple.streaming.transportStreamTimestamp') {
              var d = frame.data,
                  size = (d[3] & 0x01) << 30 | d[4] << 22 | d[5] << 14 | d[6] << 6 | d[7] >>> 2;
              size *= 4;
              size += d[7] & 0x03;
              frame.timeStamp = size; // in raw AAC, all subsequent data will be timestamped based
              // on the value of this frame
              // we couldn't have known the appropriate pts and dts before
              // parsing this ID3 tag so set those values now

              if (tag.pts === undefined && tag.dts === undefined) {
                tag.pts = frame.timeStamp;
                tag.dts = frame.timeStamp;
              }

              this.trigger('timestamp', frame);
            }
          }

          tag.frames.push(frame);
          frameStart += 10; // advance past the frame header

          frameStart += frameSize; // advance past the frame body
        } while (frameStart < tagSize);

        this.trigger('data', tag);
      };
    };

    _MetadataStream.prototype = new stream();
    var metadataStream = _MetadataStream;
    var TimestampRolloverStream$1 = timestampRolloverStream.TimestampRolloverStream; // object types

    var _TransportPacketStream, _TransportParseStream, _ElementaryStream; // constants


    var MP2T_PACKET_LENGTH = 188,
        // bytes
    SYNC_BYTE = 0x47;
    /**
     * Splits an incoming stream of binary data into MPEG-2 Transport
     * Stream packets.
     */

    _TransportPacketStream = function TransportPacketStream() {
      var buffer = new Uint8Array(MP2T_PACKET_LENGTH),
          bytesInBuffer = 0;

      _TransportPacketStream.prototype.init.call(this); // Deliver new bytes to the stream.

      /**
       * Split a stream of data into M2TS packets
      **/


      this.push = function (bytes) {
        var startIndex = 0,
            endIndex = MP2T_PACKET_LENGTH,
            everything; // If there are bytes remaining from the last segment, prepend them to the
        // bytes that were pushed in

        if (bytesInBuffer) {
          everything = new Uint8Array(bytes.byteLength + bytesInBuffer);
          everything.set(buffer.subarray(0, bytesInBuffer));
          everything.set(bytes, bytesInBuffer);
          bytesInBuffer = 0;
        } else {
          everything = bytes;
        } // While we have enough data for a packet


        while (endIndex < everything.byteLength) {
          // Look for a pair of start and end sync bytes in the data..
          if (everything[startIndex] === SYNC_BYTE && everything[endIndex] === SYNC_BYTE) {
            // We found a packet so emit it and jump one whole packet forward in
            // the stream
            this.trigger('data', everything.subarray(startIndex, endIndex));
            startIndex += MP2T_PACKET_LENGTH;
            endIndex += MP2T_PACKET_LENGTH;
            continue;
          } // If we get here, we have somehow become de-synchronized and we need to step
          // forward one byte at a time until we find a pair of sync bytes that denote
          // a packet


          startIndex++;
          endIndex++;
        } // If there was some data left over at the end of the segment that couldn't
        // possibly be a whole packet, keep it because it might be the start of a packet
        // that continues in the next segment


        if (startIndex < everything.byteLength) {
          buffer.set(everything.subarray(startIndex), 0);
          bytesInBuffer = everything.byteLength - startIndex;
        }
      };
      /**
       * Passes identified M2TS packets to the TransportParseStream to be parsed
      **/


      this.flush = function () {
        // If the buffer contains a whole packet when we are being flushed, emit it
        // and empty the buffer. Otherwise hold onto the data because it may be
        // important for decoding the next segment
        if (bytesInBuffer === MP2T_PACKET_LENGTH && buffer[0] === SYNC_BYTE) {
          this.trigger('data', buffer);
          bytesInBuffer = 0;
        }

        this.trigger('done');
      };

      this.endTimeline = function () {
        this.flush();
        this.trigger('endedtimeline');
      };

      this.reset = function () {
        bytesInBuffer = 0;
        this.trigger('reset');
      };
    };

    _TransportPacketStream.prototype = new stream();
    /**
     * Accepts an MP2T TransportPacketStream and emits data events with parsed
     * forms of the individual transport stream packets.
     */

    _TransportParseStream = function TransportParseStream() {
      var parsePsi, parsePat, parsePmt, self;

      _TransportParseStream.prototype.init.call(this);

      self = this;
      this.packetsWaitingForPmt = [];
      this.programMapTable = undefined;

      parsePsi = function parsePsi(payload, psi) {
        var offset = 0; // PSI packets may be split into multiple sections and those
        // sections may be split into multiple packets. If a PSI
        // section starts in this packet, the payload_unit_start_indicator
        // will be true and the first byte of the payload will indicate
        // the offset from the current position to the start of the
        // section.

        if (psi.payloadUnitStartIndicator) {
          offset += payload[offset] + 1;
        }

        if (psi.type === 'pat') {
          parsePat(payload.subarray(offset), psi);
        } else {
          parsePmt(payload.subarray(offset), psi);
        }
      };

      parsePat = function parsePat(payload, pat) {
        pat.section_number = payload[7]; // eslint-disable-line camelcase

        pat.last_section_number = payload[8]; // eslint-disable-line camelcase
        // skip the PSI header and parse the first PMT entry

        self.pmtPid = (payload[10] & 0x1F) << 8 | payload[11];
        pat.pmtPid = self.pmtPid;
      };
      /**
       * Parse out the relevant fields of a Program Map Table (PMT).
       * @param payload {Uint8Array} the PMT-specific portion of an MP2T
       * packet. The first byte in this array should be the table_id
       * field.
       * @param pmt {object} the object that should be decorated with
       * fields parsed from the PMT.
       */


      parsePmt = function parsePmt(payload, pmt) {
        var sectionLength, tableEnd, programInfoLength, offset; // PMTs can be sent ahead of the time when they should actually
        // take effect. We don't believe this should ever be the case
        // for HLS but we'll ignore "forward" PMT declarations if we see
        // them. Future PMT declarations have the current_next_indicator
        // set to zero.

        if (!(payload[5] & 0x01)) {
          return;
        } // overwrite any existing program map table


        self.programMapTable = {
          video: null,
          audio: null,
          'timed-metadata': {}
        }; // the mapping table ends at the end of the current section

        sectionLength = (payload[1] & 0x0f) << 8 | payload[2];
        tableEnd = 3 + sectionLength - 4; // to determine where the table is, we have to figure out how
        // long the program info descriptors are

        programInfoLength = (payload[10] & 0x0f) << 8 | payload[11]; // advance the offset to the first entry in the mapping table

        offset = 12 + programInfoLength;

        while (offset < tableEnd) {
          var streamType = payload[offset];
          var pid = (payload[offset + 1] & 0x1F) << 8 | payload[offset + 2]; // only map a single elementary_pid for audio and video stream types
          // TODO: should this be done for metadata too? for now maintain behavior of
          //       multiple metadata streams

          if (streamType === streamTypes.H264_STREAM_TYPE && self.programMapTable.video === null) {
            self.programMapTable.video = pid;
          } else if (streamType === streamTypes.ADTS_STREAM_TYPE && self.programMapTable.audio === null) {
            self.programMapTable.audio = pid;
          } else if (streamType === streamTypes.METADATA_STREAM_TYPE) {
            // map pid to stream type for metadata streams
            self.programMapTable['timed-metadata'][pid] = streamType;
          } // move to the next table entry
          // skip past the elementary stream descriptors, if present


          offset += ((payload[offset + 3] & 0x0F) << 8 | payload[offset + 4]) + 5;
        } // record the map on the packet as well


        pmt.programMapTable = self.programMapTable;
      };
      /**
       * Deliver a new MP2T packet to the next stream in the pipeline.
       */


      this.push = function (packet) {
        var result = {},
            offset = 4;
        result.payloadUnitStartIndicator = !!(packet[1] & 0x40); // pid is a 13-bit field starting at the last bit of packet[1]

        result.pid = packet[1] & 0x1f;
        result.pid <<= 8;
        result.pid |= packet[2]; // if an adaption field is present, its length is specified by the
        // fifth byte of the TS packet header. The adaptation field is
        // used to add stuffing to PES packets that don't fill a complete
        // TS packet, and to specify some forms of timing and control data
        // that we do not currently use.

        if ((packet[3] & 0x30) >>> 4 > 0x01) {
          offset += packet[offset] + 1;
        } // parse the rest of the packet based on the type


        if (result.pid === 0) {
          result.type = 'pat';
          parsePsi(packet.subarray(offset), result);
          this.trigger('data', result);
        } else if (result.pid === this.pmtPid) {
          result.type = 'pmt';
          parsePsi(packet.subarray(offset), result);
          this.trigger('data', result); // if there are any packets waiting for a PMT to be found, process them now

          while (this.packetsWaitingForPmt.length) {
            this.processPes_.apply(this, this.packetsWaitingForPmt.shift());
          }
        } else if (this.programMapTable === undefined) {
          // When we have not seen a PMT yet, defer further processing of
          // PES packets until one has been parsed
          this.packetsWaitingForPmt.push([packet, offset, result]);
        } else {
          this.processPes_(packet, offset, result);
        }
      };

      this.processPes_ = function (packet, offset, result) {
        // set the appropriate stream type
        if (result.pid === this.programMapTable.video) {
          result.streamType = streamTypes.H264_STREAM_TYPE;
        } else if (result.pid === this.programMapTable.audio) {
          result.streamType = streamTypes.ADTS_STREAM_TYPE;
        } else {
          // if not video or audio, it is timed-metadata or unknown
          // if unknown, streamType will be undefined
          result.streamType = this.programMapTable['timed-metadata'][result.pid];
        }

        result.type = 'pes';
        result.data = packet.subarray(offset);
        this.trigger('data', result);
      };
    };

    _TransportParseStream.prototype = new stream();
    _TransportParseStream.STREAM_TYPES = {
      h264: 0x1b,
      adts: 0x0f
    };
    /**
     * Reconsistutes program elementary stream (PES) packets from parsed
     * transport stream packets. That is, if you pipe an
     * mp2t.TransportParseStream into a mp2t.ElementaryStream, the output
     * events will be events which capture the bytes for individual PES
     * packets plus relevant metadata that has been extracted from the
     * container.
     */

    _ElementaryStream = function ElementaryStream() {
      var self = this,
          // PES packet fragments
      video = {
        data: [],
        size: 0
      },
          audio = {
        data: [],
        size: 0
      },
          timedMetadata = {
        data: [],
        size: 0
      },
          programMapTable,
          parsePes = function parsePes(payload, pes) {
        var ptsDtsFlags; // get the packet length, this will be 0 for video

        pes.packetLength = 6 + (payload[4] << 8 | payload[5]); // find out if this packets starts a new keyframe

        pes.dataAlignmentIndicator = (payload[6] & 0x04) !== 0; // PES packets may be annotated with a PTS value, or a PTS value
        // and a DTS value. Determine what combination of values is
        // available to work with.

        ptsDtsFlags = payload[7]; // PTS and DTS are normally stored as a 33-bit number.  Javascript
        // performs all bitwise operations on 32-bit integers but javascript
        // supports a much greater range (52-bits) of integer using standard
        // mathematical operations.
        // We construct a 31-bit value using bitwise operators over the 31
        // most significant bits and then multiply by 4 (equal to a left-shift
        // of 2) before we add the final 2 least significant bits of the
        // timestamp (equal to an OR.)

        if (ptsDtsFlags & 0xC0) {
          // the PTS and DTS are not written out directly. For information
          // on how they are encoded, see
          // http://dvd.sourceforge.net/dvdinfo/pes-hdr.html
          pes.pts = (payload[9] & 0x0E) << 27 | (payload[10] & 0xFF) << 20 | (payload[11] & 0xFE) << 12 | (payload[12] & 0xFF) << 5 | (payload[13] & 0xFE) >>> 3;
          pes.pts *= 4; // Left shift by 2

          pes.pts += (payload[13] & 0x06) >>> 1; // OR by the two LSBs

          pes.dts = pes.pts;

          if (ptsDtsFlags & 0x40) {
            pes.dts = (payload[14] & 0x0E) << 27 | (payload[15] & 0xFF) << 20 | (payload[16] & 0xFE) << 12 | (payload[17] & 0xFF) << 5 | (payload[18] & 0xFE) >>> 3;
            pes.dts *= 4; // Left shift by 2

            pes.dts += (payload[18] & 0x06) >>> 1; // OR by the two LSBs
          }
        } // the data section starts immediately after the PES header.
        // pes_header_data_length specifies the number of header bytes
        // that follow the last byte of the field.


        pes.data = payload.subarray(9 + payload[8]);
      },

      /**
        * Pass completely parsed PES packets to the next stream in the pipeline
       **/
      flushStream = function flushStream(stream, type, forceFlush) {
        var packetData = new Uint8Array(stream.size),
            event = {
          type: type
        },
            i = 0,
            offset = 0,
            packetFlushable = false,
            fragment; // do nothing if there is not enough buffered data for a complete
        // PES header

        if (!stream.data.length || stream.size < 9) {
          return;
        }

        event.trackId = stream.data[0].pid; // reassemble the packet

        for (i = 0; i < stream.data.length; i++) {
          fragment = stream.data[i];
          packetData.set(fragment.data, offset);
          offset += fragment.data.byteLength;
        } // parse assembled packet's PES header


        parsePes(packetData, event); // non-video PES packets MUST have a non-zero PES_packet_length
        // check that there is enough stream data to fill the packet

        packetFlushable = type === 'video' || event.packetLength <= stream.size; // flush pending packets if the conditions are right

        if (forceFlush || packetFlushable) {
          stream.size = 0;
          stream.data.length = 0;
        } // only emit packets that are complete. this is to avoid assembling
        // incomplete PES packets due to poor segmentation


        if (packetFlushable) {
          self.trigger('data', event);
        }
      };

      _ElementaryStream.prototype.init.call(this);
      /**
       * Identifies M2TS packet types and parses PES packets using metadata
       * parsed from the PMT
       **/


      this.push = function (data) {
        ({
          pat: function pat() {// we have to wait for the PMT to arrive as well before we
            // have any meaningful metadata
          },
          pes: function pes() {
            var stream, streamType;

            switch (data.streamType) {
              case streamTypes.H264_STREAM_TYPE:
                stream = video;
                streamType = 'video';
                break;

              case streamTypes.ADTS_STREAM_TYPE:
                stream = audio;
                streamType = 'audio';
                break;

              case streamTypes.METADATA_STREAM_TYPE:
                stream = timedMetadata;
                streamType = 'timed-metadata';
                break;

              default:
                // ignore unknown stream types
                return;
            } // if a new packet is starting, we can flush the completed
            // packet


            if (data.payloadUnitStartIndicator) {
              flushStream(stream, streamType, true);
            } // buffer this fragment until we are sure we've received the
            // complete payload


            stream.data.push(data);
            stream.size += data.data.byteLength;
          },
          pmt: function pmt() {
            var event = {
              type: 'metadata',
              tracks: []
            };
            programMapTable = data.programMapTable; // translate audio and video streams to tracks

            if (programMapTable.video !== null) {
              event.tracks.push({
                timelineStartInfo: {
                  baseMediaDecodeTime: 0
                },
                id: +programMapTable.video,
                codec: 'avc',
                type: 'video'
              });
            }

            if (programMapTable.audio !== null) {
              event.tracks.push({
                timelineStartInfo: {
                  baseMediaDecodeTime: 0
                },
                id: +programMapTable.audio,
                codec: 'adts',
                type: 'audio'
              });
            }

            self.trigger('data', event);
          }
        })[data.type]();
      };

      this.reset = function () {
        video.size = 0;
        video.data.length = 0;
        audio.size = 0;
        audio.data.length = 0;
        this.trigger('reset');
      };
      /**
       * Flush any remaining input. Video PES packets may be of variable
       * length. Normally, the start of a new video packet can trigger the
       * finalization of the previous packet. That is not possible if no
       * more video is forthcoming, however. In that case, some other
       * mechanism (like the end of the file) has to be employed. When it is
       * clear that no additional data is forthcoming, calling this method
       * will flush the buffered packets.
       */


      this.flushStreams_ = function () {
        // !!THIS ORDER IS IMPORTANT!!
        // video first then audio
        flushStream(video, 'video');
        flushStream(audio, 'audio');
        flushStream(timedMetadata, 'timed-metadata');
      };

      this.flush = function () {
        this.flushStreams_();
        this.trigger('done');
      };
    };

    _ElementaryStream.prototype = new stream();
    var m2ts = {
      PAT_PID: 0x0000,
      MP2T_PACKET_LENGTH: MP2T_PACKET_LENGTH,
      TransportPacketStream: _TransportPacketStream,
      TransportParseStream: _TransportParseStream,
      ElementaryStream: _ElementaryStream,
      TimestampRolloverStream: TimestampRolloverStream$1,
      CaptionStream: captionStream.CaptionStream,
      Cea608Stream: captionStream.Cea608Stream,
      MetadataStream: metadataStream
    };

    for (var type in streamTypes) {
      if (streamTypes.hasOwnProperty(type)) {
        m2ts[type] = streamTypes[type];
      }
    }

    var m2ts_1 = m2ts;
    var ONE_SECOND_IN_TS$2 = clock.ONE_SECOND_IN_TS;

    var _AdtsStream;

    var ADTS_SAMPLING_FREQUENCIES = [96000, 88200, 64000, 48000, 44100, 32000, 24000, 22050, 16000, 12000, 11025, 8000, 7350];
    /*
     * Accepts a ElementaryStream and emits data events with parsed
     * AAC Audio Frames of the individual packets. Input audio in ADTS
     * format is unpacked and re-emitted as AAC frames.
     *
     * @see http://wiki.multimedia.cx/index.php?title=ADTS
     * @see http://wiki.multimedia.cx/?title=Understanding_AAC
     */

    _AdtsStream = function AdtsStream(handlePartialSegments) {
      var buffer,
          frameNum = 0;

      _AdtsStream.prototype.init.call(this);

      this.push = function (packet) {
        var i = 0,
            frameLength,
            protectionSkipBytes,
            frameEnd,
            oldBuffer,
            sampleCount,
            adtsFrameDuration;

        if (!handlePartialSegments) {
          frameNum = 0;
        }

        if (packet.type !== 'audio') {
          // ignore non-audio data
          return;
        } // Prepend any data in the buffer to the input data so that we can parse
        // aac frames the cross a PES packet boundary


        if (buffer) {
          oldBuffer = buffer;
          buffer = new Uint8Array(oldBuffer.byteLength + packet.data.byteLength);
          buffer.set(oldBuffer);
          buffer.set(packet.data, oldBuffer.byteLength);
        } else {
          buffer = packet.data;
        } // unpack any ADTS frames which have been fully received
        // for details on the ADTS header, see http://wiki.multimedia.cx/index.php?title=ADTS


        while (i + 5 < buffer.length) {
          // Look for the start of an ADTS header..
          if (buffer[i] !== 0xFF || (buffer[i + 1] & 0xF6) !== 0xF0) {
            // If a valid header was not found,  jump one forward and attempt to
            // find a valid ADTS header starting at the next byte
            i++;
            continue;
          } // The protection skip bit tells us if we have 2 bytes of CRC data at the
          // end of the ADTS header


          protectionSkipBytes = (~buffer[i + 1] & 0x01) * 2; // Frame length is a 13 bit integer starting 16 bits from the
          // end of the sync sequence

          frameLength = (buffer[i + 3] & 0x03) << 11 | buffer[i + 4] << 3 | (buffer[i + 5] & 0xe0) >> 5;
          sampleCount = ((buffer[i + 6] & 0x03) + 1) * 1024;
          adtsFrameDuration = sampleCount * ONE_SECOND_IN_TS$2 / ADTS_SAMPLING_FREQUENCIES[(buffer[i + 2] & 0x3c) >>> 2];
          frameEnd = i + frameLength; // If we don't have enough data to actually finish this ADTS frame, return
          // and wait for more data

          if (buffer.byteLength < frameEnd) {
            return;
          } // Otherwise, deliver the complete AAC frame


          this.trigger('data', {
            pts: packet.pts + frameNum * adtsFrameDuration,
            dts: packet.dts + frameNum * adtsFrameDuration,
            sampleCount: sampleCount,
            audioobjecttype: (buffer[i + 2] >>> 6 & 0x03) + 1,
            channelcount: (buffer[i + 2] & 1) << 2 | (buffer[i + 3] & 0xc0) >>> 6,
            samplerate: ADTS_SAMPLING_FREQUENCIES[(buffer[i + 2] & 0x3c) >>> 2],
            samplingfrequencyindex: (buffer[i + 2] & 0x3c) >>> 2,
            // assume ISO/IEC 14496-12 AudioSampleEntry default of 16
            samplesize: 16,
            data: buffer.subarray(i + 7 + protectionSkipBytes, frameEnd)
          });
          frameNum++; // If the buffer is empty, clear it and return

          if (buffer.byteLength === frameEnd) {
            buffer = undefined;
            return;
          } // Remove the finished frame from the buffer and start the process again


          buffer = buffer.subarray(frameEnd);
        }
      };

      this.flush = function () {
        frameNum = 0;
        this.trigger('done');
      };

      this.reset = function () {
        buffer = void 0;
        this.trigger('reset');
      };

      this.endTimeline = function () {
        buffer = void 0;
        this.trigger('endedtimeline');
      };
    };

    _AdtsStream.prototype = new stream();
    var adts = _AdtsStream;
    /**
     * mux.js
     *
     * Copyright (c) Brightcove
     * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
     */

    var ExpGolomb;
    /**
     * Parser for exponential Golomb codes, a variable-bitwidth number encoding
     * scheme used by h264.
     */

    ExpGolomb = function ExpGolomb(workingData) {
      var // the number of bytes left to examine in workingData
      workingBytesAvailable = workingData.byteLength,
          // the current word being examined
      workingWord = 0,
          // :uint
      // the number of bits left to examine in the current word
      workingBitsAvailable = 0; // :uint;
      // ():uint

      this.length = function () {
        return 8 * workingBytesAvailable;
      }; // ():uint


      this.bitsAvailable = function () {
        return 8 * workingBytesAvailable + workingBitsAvailable;
      }; // ():void


      this.loadWord = function () {
        var position = workingData.byteLength - workingBytesAvailable,
            workingBytes = new Uint8Array(4),
            availableBytes = Math.min(4, workingBytesAvailable);

        if (availableBytes === 0) {
          throw new Error('no bytes available');
        }

        workingBytes.set(workingData.subarray(position, position + availableBytes));
        workingWord = new DataView(workingBytes.buffer).getUint32(0); // track the amount of workingData that has been processed

        workingBitsAvailable = availableBytes * 8;
        workingBytesAvailable -= availableBytes;
      }; // (count:int):void


      this.skipBits = function (count) {
        var skipBytes; // :int

        if (workingBitsAvailable > count) {
          workingWord <<= count;
          workingBitsAvailable -= count;
        } else {
          count -= workingBitsAvailable;
          skipBytes = Math.floor(count / 8);
          count -= skipBytes * 8;
          workingBytesAvailable -= skipBytes;
          this.loadWord();
          workingWord <<= count;
          workingBitsAvailable -= count;
        }
      }; // (size:int):uint


      this.readBits = function (size) {
        var bits = Math.min(workingBitsAvailable, size),
            // :uint
        valu = workingWord >>> 32 - bits; // :uint
        // if size > 31, handle error

        workingBitsAvailable -= bits;

        if (workingBitsAvailable > 0) {
          workingWord <<= bits;
        } else if (workingBytesAvailable > 0) {
          this.loadWord();
        }

        bits = size - bits;

        if (bits > 0) {
          return valu << bits | this.readBits(bits);
        }

        return valu;
      }; // ():uint


      this.skipLeadingZeros = function () {
        var leadingZeroCount; // :uint

        for (leadingZeroCount = 0; leadingZeroCount < workingBitsAvailable; ++leadingZeroCount) {
          if ((workingWord & 0x80000000 >>> leadingZeroCount) !== 0) {
            // the first bit of working word is 1
            workingWord <<= leadingZeroCount;
            workingBitsAvailable -= leadingZeroCount;
            return leadingZeroCount;
          }
        } // we exhausted workingWord and still have not found a 1


        this.loadWord();
        return leadingZeroCount + this.skipLeadingZeros();
      }; // ():void


      this.skipUnsignedExpGolomb = function () {
        this.skipBits(1 + this.skipLeadingZeros());
      }; // ():void


      this.skipExpGolomb = function () {
        this.skipBits(1 + this.skipLeadingZeros());
      }; // ():uint


      this.readUnsignedExpGolomb = function () {
        var clz = this.skipLeadingZeros(); // :uint

        return this.readBits(clz + 1) - 1;
      }; // ():int


      this.readExpGolomb = function () {
        var valu = this.readUnsignedExpGolomb(); // :int

        if (0x01 & valu) {
          // the number is odd if the low order bit is set
          return 1 + valu >>> 1; // add 1 to make it even, and divide by 2
        }

        return -1 * (valu >>> 1); // divide by two then make it negative
      }; // Some convenience functions
      // :Boolean


      this.readBoolean = function () {
        return this.readBits(1) === 1;
      }; // ():int


      this.readUnsignedByte = function () {
        return this.readBits(8);
      };

      this.loadWord();
    };

    var expGolomb = ExpGolomb;

    var _H264Stream, _NalByteStream;

    var PROFILES_WITH_OPTIONAL_SPS_DATA;
    /**
     * Accepts a NAL unit byte stream and unpacks the embedded NAL units.
     */

    _NalByteStream = function NalByteStream() {
      var syncPoint = 0,
          i,
          buffer;

      _NalByteStream.prototype.init.call(this);
      /*
       * Scans a byte stream and triggers a data event with the NAL units found.
       * @param {Object} data Event received from H264Stream
       * @param {Uint8Array} data.data The h264 byte stream to be scanned
       *
       * @see H264Stream.push
       */


      this.push = function (data) {
        var swapBuffer;

        if (!buffer) {
          buffer = data.data;
        } else {
          swapBuffer = new Uint8Array(buffer.byteLength + data.data.byteLength);
          swapBuffer.set(buffer);
          swapBuffer.set(data.data, buffer.byteLength);
          buffer = swapBuffer;
        }

        var len = buffer.byteLength; // Rec. ITU-T H.264, Annex B
        // scan for NAL unit boundaries
        // a match looks like this:
        // 0 0 1 .. NAL .. 0 0 1
        // ^ sync point        ^ i
        // or this:
        // 0 0 1 .. NAL .. 0 0 0
        // ^ sync point        ^ i
        // advance the sync point to a NAL start, if necessary

        for (; syncPoint < len - 3; syncPoint++) {
          if (buffer[syncPoint + 2] === 1) {
            // the sync point is properly aligned
            i = syncPoint + 5;
            break;
          }
        }

        while (i < len) {
          // look at the current byte to determine if we've hit the end of
          // a NAL unit boundary
          switch (buffer[i]) {
            case 0:
              // skip past non-sync sequences
              if (buffer[i - 1] !== 0) {
                i += 2;
                break;
              } else if (buffer[i - 2] !== 0) {
                i++;
                break;
              } // deliver the NAL unit if it isn't empty


              if (syncPoint + 3 !== i - 2) {
                this.trigger('data', buffer.subarray(syncPoint + 3, i - 2));
              } // drop trailing zeroes


              do {
                i++;
              } while (buffer[i] !== 1 && i < len);

              syncPoint = i - 2;
              i += 3;
              break;

            case 1:
              // skip past non-sync sequences
              if (buffer[i - 1] !== 0 || buffer[i - 2] !== 0) {
                i += 3;
                break;
              } // deliver the NAL unit


              this.trigger('data', buffer.subarray(syncPoint + 3, i - 2));
              syncPoint = i - 2;
              i += 3;
              break;

            default:
              // the current byte isn't a one or zero, so it cannot be part
              // of a sync sequence
              i += 3;
              break;
          }
        } // filter out the NAL units that were delivered


        buffer = buffer.subarray(syncPoint);
        i -= syncPoint;
        syncPoint = 0;
      };

      this.reset = function () {
        buffer = null;
        syncPoint = 0;
        this.trigger('reset');
      };

      this.flush = function () {
        // deliver the last buffered NAL unit
        if (buffer && buffer.byteLength > 3) {
          this.trigger('data', buffer.subarray(syncPoint + 3));
        } // reset the stream state


        buffer = null;
        syncPoint = 0;
        this.trigger('done');
      };

      this.endTimeline = function () {
        this.flush();
        this.trigger('endedtimeline');
      };
    };

    _NalByteStream.prototype = new stream(); // values of profile_idc that indicate additional fields are included in the SPS
    // see Recommendation ITU-T H.264 (4/2013),
    // 7.3.2.1.1 Sequence parameter set data syntax

    PROFILES_WITH_OPTIONAL_SPS_DATA = {
      100: true,
      110: true,
      122: true,
      244: true,
      44: true,
      83: true,
      86: true,
      118: true,
      128: true,
      138: true,
      139: true,
      134: true
    };
    /**
     * Accepts input from a ElementaryStream and produces H.264 NAL unit data
     * events.
     */

    _H264Stream = function H264Stream() {
      var nalByteStream = new _NalByteStream(),
          self,
          trackId,
          currentPts,
          currentDts,
          discardEmulationPreventionBytes,
          readSequenceParameterSet,
          skipScalingList;

      _H264Stream.prototype.init.call(this);

      self = this;
      /*
       * Pushes a packet from a stream onto the NalByteStream
       *
       * @param {Object} packet - A packet received from a stream
       * @param {Uint8Array} packet.data - The raw bytes of the packet
       * @param {Number} packet.dts - Decode timestamp of the packet
       * @param {Number} packet.pts - Presentation timestamp of the packet
       * @param {Number} packet.trackId - The id of the h264 track this packet came from
       * @param {('video'|'audio')} packet.type - The type of packet
       *
       */

      this.push = function (packet) {
        if (packet.type !== 'video') {
          return;
        }

        trackId = packet.trackId;
        currentPts = packet.pts;
        currentDts = packet.dts;
        nalByteStream.push(packet);
      };
      /*
       * Identify NAL unit types and pass on the NALU, trackId, presentation and decode timestamps
       * for the NALUs to the next stream component.
       * Also, preprocess caption and sequence parameter NALUs.
       *
       * @param {Uint8Array} data - A NAL unit identified by `NalByteStream.push`
       * @see NalByteStream.push
       */


      nalByteStream.on('data', function (data) {
        var event = {
          trackId: trackId,
          pts: currentPts,
          dts: currentDts,
          data: data
        };

        switch (data[0] & 0x1f) {
          case 0x05:
            event.nalUnitType = 'slice_layer_without_partitioning_rbsp_idr';
            break;

          case 0x06:
            event.nalUnitType = 'sei_rbsp';
            event.escapedRBSP = discardEmulationPreventionBytes(data.subarray(1));
            break;

          case 0x07:
            event.nalUnitType = 'seq_parameter_set_rbsp';
            event.escapedRBSP = discardEmulationPreventionBytes(data.subarray(1));
            event.config = readSequenceParameterSet(event.escapedRBSP);
            break;

          case 0x08:
            event.nalUnitType = 'pic_parameter_set_rbsp';
            break;

          case 0x09:
            event.nalUnitType = 'access_unit_delimiter_rbsp';
            break;
        } // This triggers data on the H264Stream


        self.trigger('data', event);
      });
      nalByteStream.on('done', function () {
        self.trigger('done');
      });
      nalByteStream.on('partialdone', function () {
        self.trigger('partialdone');
      });
      nalByteStream.on('reset', function () {
        self.trigger('reset');
      });
      nalByteStream.on('endedtimeline', function () {
        self.trigger('endedtimeline');
      });

      this.flush = function () {
        nalByteStream.flush();
      };

      this.partialFlush = function () {
        nalByteStream.partialFlush();
      };

      this.reset = function () {
        nalByteStream.reset();
      };

      this.endTimeline = function () {
        nalByteStream.endTimeline();
      };
      /**
       * Advance the ExpGolomb decoder past a scaling list. The scaling
       * list is optionally transmitted as part of a sequence parameter
       * set and is not relevant to transmuxing.
       * @param count {number} the number of entries in this scaling list
       * @param expGolombDecoder {object} an ExpGolomb pointed to the
       * start of a scaling list
       * @see Recommendation ITU-T H.264, Section 7.3.2.1.1.1
       */


      skipScalingList = function skipScalingList(count, expGolombDecoder) {
        var lastScale = 8,
            nextScale = 8,
            j,
            deltaScale;

        for (j = 0; j < count; j++) {
          if (nextScale !== 0) {
            deltaScale = expGolombDecoder.readExpGolomb();
            nextScale = (lastScale + deltaScale + 256) % 256;
          }

          lastScale = nextScale === 0 ? lastScale : nextScale;
        }
      };
      /**
       * Expunge any "Emulation Prevention" bytes from a "Raw Byte
       * Sequence Payload"
       * @param data {Uint8Array} the bytes of a RBSP from a NAL
       * unit
       * @return {Uint8Array} the RBSP without any Emulation
       * Prevention Bytes
       */


      discardEmulationPreventionBytes = function discardEmulationPreventionBytes(data) {
        var length = data.byteLength,
            emulationPreventionBytesPositions = [],
            i = 1,
            newLength,
            newData; // Find all `Emulation Prevention Bytes`

        while (i < length - 2) {
          if (data[i] === 0 && data[i + 1] === 0 && data[i + 2] === 0x03) {
            emulationPreventionBytesPositions.push(i + 2);
            i += 2;
          } else {
            i++;
          }
        } // If no Emulation Prevention Bytes were found just return the original
        // array


        if (emulationPreventionBytesPositions.length === 0) {
          return data;
        } // Create a new array to hold the NAL unit data


        newLength = length - emulationPreventionBytesPositions.length;
        newData = new Uint8Array(newLength);
        var sourceIndex = 0;

        for (i = 0; i < newLength; sourceIndex++, i++) {
          if (sourceIndex === emulationPreventionBytesPositions[0]) {
            // Skip this byte
            sourceIndex++; // Remove this position index

            emulationPreventionBytesPositions.shift();
          }

          newData[i] = data[sourceIndex];
        }

        return newData;
      };
      /**
       * Read a sequence parameter set and return some interesting video
       * properties. A sequence parameter set is the H264 metadata that
       * describes the properties of upcoming video frames.
       * @param data {Uint8Array} the bytes of a sequence parameter set
       * @return {object} an object with configuration parsed from the
       * sequence parameter set, including the dimensions of the
       * associated video frames.
       */


      readSequenceParameterSet = function readSequenceParameterSet(data) {
        var frameCropLeftOffset = 0,
            frameCropRightOffset = 0,
            frameCropTopOffset = 0,
            frameCropBottomOffset = 0,
            sarScale = 1,
            expGolombDecoder,
            profileIdc,
            levelIdc,
            profileCompatibility,
            chromaFormatIdc,
            picOrderCntType,
            numRefFramesInPicOrderCntCycle,
            picWidthInMbsMinus1,
            picHeightInMapUnitsMinus1,
            frameMbsOnlyFlag,
            scalingListCount,
            sarRatio,
            aspectRatioIdc,
            i;
        expGolombDecoder = new expGolomb(data);
        profileIdc = expGolombDecoder.readUnsignedByte(); // profile_idc

        profileCompatibility = expGolombDecoder.readUnsignedByte(); // constraint_set[0-5]_flag

        levelIdc = expGolombDecoder.readUnsignedByte(); // level_idc u(8)

        expGolombDecoder.skipUnsignedExpGolomb(); // seq_parameter_set_id
        // some profiles have more optional data we don't need

        if (PROFILES_WITH_OPTIONAL_SPS_DATA[profileIdc]) {
          chromaFormatIdc = expGolombDecoder.readUnsignedExpGolomb();

          if (chromaFormatIdc === 3) {
            expGolombDecoder.skipBits(1); // separate_colour_plane_flag
          }

          expGolombDecoder.skipUnsignedExpGolomb(); // bit_depth_luma_minus8

          expGolombDecoder.skipUnsignedExpGolomb(); // bit_depth_chroma_minus8

          expGolombDecoder.skipBits(1); // qpprime_y_zero_transform_bypass_flag

          if (expGolombDecoder.readBoolean()) {
            // seq_scaling_matrix_present_flag
            scalingListCount = chromaFormatIdc !== 3 ? 8 : 12;

            for (i = 0; i < scalingListCount; i++) {
              if (expGolombDecoder.readBoolean()) {
                // seq_scaling_list_present_flag[ i ]
                if (i < 6) {
                  skipScalingList(16, expGolombDecoder);
                } else {
                  skipScalingList(64, expGolombDecoder);
                }
              }
            }
          }
        }

        expGolombDecoder.skipUnsignedExpGolomb(); // log2_max_frame_num_minus4

        picOrderCntType = expGolombDecoder.readUnsignedExpGolomb();

        if (picOrderCntType === 0) {
          expGolombDecoder.readUnsignedExpGolomb(); // log2_max_pic_order_cnt_lsb_minus4
        } else if (picOrderCntType === 1) {
          expGolombDecoder.skipBits(1); // delta_pic_order_always_zero_flag

          expGolombDecoder.skipExpGolomb(); // offset_for_non_ref_pic

          expGolombDecoder.skipExpGolomb(); // offset_for_top_to_bottom_field

          numRefFramesInPicOrderCntCycle = expGolombDecoder.readUnsignedExpGolomb();

          for (i = 0; i < numRefFramesInPicOrderCntCycle; i++) {
            expGolombDecoder.skipExpGolomb(); // offset_for_ref_frame[ i ]
          }
        }

        expGolombDecoder.skipUnsignedExpGolomb(); // max_num_ref_frames

        expGolombDecoder.skipBits(1); // gaps_in_frame_num_value_allowed_flag

        picWidthInMbsMinus1 = expGolombDecoder.readUnsignedExpGolomb();
        picHeightInMapUnitsMinus1 = expGolombDecoder.readUnsignedExpGolomb();
        frameMbsOnlyFlag = expGolombDecoder.readBits(1);

        if (frameMbsOnlyFlag === 0) {
          expGolombDecoder.skipBits(1); // mb_adaptive_frame_field_flag
        }

        expGolombDecoder.skipBits(1); // direct_8x8_inference_flag

        if (expGolombDecoder.readBoolean()) {
          // frame_cropping_flag
          frameCropLeftOffset = expGolombDecoder.readUnsignedExpGolomb();
          frameCropRightOffset = expGolombDecoder.readUnsignedExpGolomb();
          frameCropTopOffset = expGolombDecoder.readUnsignedExpGolomb();
          frameCropBottomOffset = expGolombDecoder.readUnsignedExpGolomb();
        }

        if (expGolombDecoder.readBoolean()) {
          // vui_parameters_present_flag
          if (expGolombDecoder.readBoolean()) {
            // aspect_ratio_info_present_flag
            aspectRatioIdc = expGolombDecoder.readUnsignedByte();

            switch (aspectRatioIdc) {
              case 1:
                sarRatio = [1, 1];
                break;

              case 2:
                sarRatio = [12, 11];
                break;

              case 3:
                sarRatio = [10, 11];
                break;

              case 4:
                sarRatio = [16, 11];
                break;

              case 5:
                sarRatio = [40, 33];
                break;

              case 6:
                sarRatio = [24, 11];
                break;

              case 7:
                sarRatio = [20, 11];
                break;

              case 8:
                sarRatio = [32, 11];
                break;

              case 9:
                sarRatio = [80, 33];
                break;

              case 10:
                sarRatio = [18, 11];
                break;

              case 11:
                sarRatio = [15, 11];
                break;

              case 12:
                sarRatio = [64, 33];
                break;

              case 13:
                sarRatio = [160, 99];
                break;

              case 14:
                sarRatio = [4, 3];
                break;

              case 15:
                sarRatio = [3, 2];
                break;

              case 16:
                sarRatio = [2, 1];
                break;

              case 255:
                {
                  sarRatio = [expGolombDecoder.readUnsignedByte() << 8 | expGolombDecoder.readUnsignedByte(), expGolombDecoder.readUnsignedByte() << 8 | expGolombDecoder.readUnsignedByte()];
                  break;
                }
            }

            if (sarRatio) {
              sarScale = sarRatio[0] / sarRatio[1];
            }
          }
        }

        return {
          profileIdc: profileIdc,
          levelIdc: levelIdc,
          profileCompatibility: profileCompatibility,
          width: Math.ceil(((picWidthInMbsMinus1 + 1) * 16 - frameCropLeftOffset * 2 - frameCropRightOffset * 2) * sarScale),
          height: (2 - frameMbsOnlyFlag) * (picHeightInMapUnitsMinus1 + 1) * 16 - frameCropTopOffset * 2 - frameCropBottomOffset * 2,
          sarRatio: sarRatio
        };
      };
    };

    _H264Stream.prototype = new stream();
    var h264 = {
      H264Stream: _H264Stream,
      NalByteStream: _NalByteStream
    };
    /**
     * mux.js
     *
     * Copyright (c) Brightcove
     * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
     *
     * Utilities to detect basic properties and metadata about Aac data.
     */

    var ADTS_SAMPLING_FREQUENCIES$1 = [96000, 88200, 64000, 48000, 44100, 32000, 24000, 22050, 16000, 12000, 11025, 8000, 7350];

    var parseId3TagSize = function parseId3TagSize(header, byteIndex) {
      var returnSize = header[byteIndex + 6] << 21 | header[byteIndex + 7] << 14 | header[byteIndex + 8] << 7 | header[byteIndex + 9],
          flags = header[byteIndex + 5],
          footerPresent = (flags & 16) >> 4;

      if (footerPresent) {
        return returnSize + 20;
      }

      return returnSize + 10;
    }; // TODO: use vhs-utils


    var isLikelyAacData = function isLikelyAacData(data) {
      var offset = 0;

      if (data.length > 10 && data[0] === 'I'.charCodeAt(0) && data[1] === 'D'.charCodeAt(0) && data[2] === '3'.charCodeAt(0)) {
        offset = parseId3TagSize(data, 0);
      }

      return data.length >= offset + 2 && (data[offset] & 0xFF) === 0xFF && (data[offset + 1] & 0xF0) === 0xF0 && // verify that the 2 layer bits are 0, aka this
      // is not mp3 data but aac data.
      (data[offset + 1] & 0x16) === 0x10;
    };

    var parseSyncSafeInteger$1 = function parseSyncSafeInteger$1(data) {
      return data[0] << 21 | data[1] << 14 | data[2] << 7 | data[3];
    }; // return a percent-encoded representation of the specified byte range
    // @see http://en.wikipedia.org/wiki/Percent-encoding


    var percentEncode$1 = function percentEncode$1(bytes, start, end) {
      var i,
          result = '';

      for (i = start; i < end; i++) {
        result += '%' + ('00' + bytes[i].toString(16)).slice(-2);
      }

      return result;
    }; // return the string representation of the specified byte range,
    // interpreted as ISO-8859-1.


    var parseIso88591$1 = function parseIso88591$1(bytes, start, end) {
      return unescape(percentEncode$1(bytes, start, end)); // jshint ignore:line
    };

    var parseAdtsSize = function parseAdtsSize(header, byteIndex) {
      var lowThree = (header[byteIndex + 5] & 0xE0) >> 5,
          middle = header[byteIndex + 4] << 3,
          highTwo = header[byteIndex + 3] & 0x3 << 11;
      return highTwo | middle | lowThree;
    };

    var parseType = function parseType(header, byteIndex) {
      if (header[byteIndex] === 'I'.charCodeAt(0) && header[byteIndex + 1] === 'D'.charCodeAt(0) && header[byteIndex + 2] === '3'.charCodeAt(0)) {
        return 'timed-metadata';
      } else if (header[byteIndex] & 0xff === 0xff && (header[byteIndex + 1] & 0xf0) === 0xf0) {
        return 'audio';
      }

      return null;
    };

    var parseSampleRate = function parseSampleRate(packet) {
      var i = 0;

      while (i + 5 < packet.length) {
        if (packet[i] !== 0xFF || (packet[i + 1] & 0xF6) !== 0xF0) {
          // If a valid header was not found,  jump one forward and attempt to
          // find a valid ADTS header starting at the next byte
          i++;
          continue;
        }

        return ADTS_SAMPLING_FREQUENCIES$1[(packet[i + 2] & 0x3c) >>> 2];
      }

      return null;
    };

    var parseAacTimestamp = function parseAacTimestamp(packet) {
      var frameStart, frameSize, frame, frameHeader; // find the start of the first frame and the end of the tag

      frameStart = 10;

      if (packet[5] & 0x40) {
        // advance the frame start past the extended header
        frameStart += 4; // header size field

        frameStart += parseSyncSafeInteger$1(packet.subarray(10, 14));
      } // parse one or more ID3 frames
      // http://id3.org/id3v2.3.0#ID3v2_frame_overview


      do {
        // determine the number of bytes in this frame
        frameSize = parseSyncSafeInteger$1(packet.subarray(frameStart + 4, frameStart + 8));

        if (frameSize < 1) {
          return null;
        }

        frameHeader = String.fromCharCode(packet[frameStart], packet[frameStart + 1], packet[frameStart + 2], packet[frameStart + 3]);

        if (frameHeader === 'PRIV') {
          frame = packet.subarray(frameStart + 10, frameStart + frameSize + 10);

          for (var i = 0; i < frame.byteLength; i++) {
            if (frame[i] === 0) {
              var owner = parseIso88591$1(frame, 0, i);

              if (owner === 'com.apple.streaming.transportStreamTimestamp') {
                var d = frame.subarray(i + 1);
                var size = (d[3] & 0x01) << 30 | d[4] << 22 | d[5] << 14 | d[6] << 6 | d[7] >>> 2;
                size *= 4;
                size += d[7] & 0x03;
                return size;
              }

              break;
            }
          }
        }

        frameStart += 10; // advance past the frame header

        frameStart += frameSize; // advance past the frame body
      } while (frameStart < packet.byteLength);

      return null;
    };

    var utils = {
      isLikelyAacData: isLikelyAacData,
      parseId3TagSize: parseId3TagSize,
      parseAdtsSize: parseAdtsSize,
      parseType: parseType,
      parseSampleRate: parseSampleRate,
      parseAacTimestamp: parseAacTimestamp
    }; // Constants

    var _AacStream;
    /**
     * Splits an incoming stream of binary data into ADTS and ID3 Frames.
     */


    _AacStream = function AacStream() {
      var everything = new Uint8Array(),
          timeStamp = 0;

      _AacStream.prototype.init.call(this);

      this.setTimestamp = function (timestamp) {
        timeStamp = timestamp;
      };

      this.push = function (bytes) {
        var frameSize = 0,
            byteIndex = 0,
            bytesLeft,
            chunk,
            packet,
            tempLength; // If there are bytes remaining from the last segment, prepend them to the
        // bytes that were pushed in

        if (everything.length) {
          tempLength = everything.length;
          everything = new Uint8Array(bytes.byteLength + tempLength);
          everything.set(everything.subarray(0, tempLength));
          everything.set(bytes, tempLength);
        } else {
          everything = bytes;
        }

        while (everything.length - byteIndex >= 3) {
          if (everything[byteIndex] === 'I'.charCodeAt(0) && everything[byteIndex + 1] === 'D'.charCodeAt(0) && everything[byteIndex + 2] === '3'.charCodeAt(0)) {
            // Exit early because we don't have enough to parse
            // the ID3 tag header
            if (everything.length - byteIndex < 10) {
              break;
            } // check framesize


            frameSize = utils.parseId3TagSize(everything, byteIndex); // Exit early if we don't have enough in the buffer
            // to emit a full packet
            // Add to byteIndex to support multiple ID3 tags in sequence

            if (byteIndex + frameSize > everything.length) {
              break;
            }

            chunk = {
              type: 'timed-metadata',
              data: everything.subarray(byteIndex, byteIndex + frameSize)
            };
            this.trigger('data', chunk);
            byteIndex += frameSize;
            continue;
          } else if ((everything[byteIndex] & 0xff) === 0xff && (everything[byteIndex + 1] & 0xf0) === 0xf0) {
            // Exit early because we don't have enough to parse
            // the ADTS frame header
            if (everything.length - byteIndex < 7) {
              break;
            }

            frameSize = utils.parseAdtsSize(everything, byteIndex); // Exit early if we don't have enough in the buffer
            // to emit a full packet

            if (byteIndex + frameSize > everything.length) {
              break;
            }

            packet = {
              type: 'audio',
              data: everything.subarray(byteIndex, byteIndex + frameSize),
              pts: timeStamp,
              dts: timeStamp
            };
            this.trigger('data', packet);
            byteIndex += frameSize;
            continue;
          }

          byteIndex++;
        }

        bytesLeft = everything.length - byteIndex;

        if (bytesLeft > 0) {
          everything = everything.subarray(byteIndex);
        } else {
          everything = new Uint8Array();
        }
      };

      this.reset = function () {
        everything = new Uint8Array();
        this.trigger('reset');
      };

      this.endTimeline = function () {
        everything = new Uint8Array();
        this.trigger('endedtimeline');
      };
    };

    _AacStream.prototype = new stream();
    var aac = _AacStream; // constants

    var AUDIO_PROPERTIES = ['audioobjecttype', 'channelcount', 'samplerate', 'samplingfrequencyindex', 'samplesize'];
    var audioProperties = AUDIO_PROPERTIES;
    var VIDEO_PROPERTIES = ['width', 'height', 'profileIdc', 'levelIdc', 'profileCompatibility', 'sarRatio'];
    var videoProperties = VIDEO_PROPERTIES;
    var H264Stream$1 = h264.H264Stream;
    var isLikelyAacData$1 = utils.isLikelyAacData;
    var ONE_SECOND_IN_TS$3 = clock.ONE_SECOND_IN_TS; // object types

    var _VideoSegmentStream, _AudioSegmentStream, _Transmuxer, _CoalesceStream;
    /**
     * Compare two arrays (even typed) for same-ness
     */


    var arrayEquals = function arrayEquals(a, b) {
      var i;

      if (a.length !== b.length) {
        return false;
      } // compare the value of each element in the array


      for (i = 0; i < a.length; i++) {
        if (a[i] !== b[i]) {
          return false;
        }
      }

      return true;
    };

    var generateVideoSegmentTimingInfo = function generateVideoSegmentTimingInfo(baseMediaDecodeTime, startDts, startPts, endDts, endPts, prependedContentDuration) {
      var ptsOffsetFromDts = startPts - startDts,
          decodeDuration = endDts - startDts,
          presentationDuration = endPts - startPts; // The PTS and DTS values are based on the actual stream times from the segment,
      // however, the player time values will reflect a start from the baseMediaDecodeTime.
      // In order to provide relevant values for the player times, base timing info on the
      // baseMediaDecodeTime and the DTS and PTS durations of the segment.

      return {
        start: {
          dts: baseMediaDecodeTime,
          pts: baseMediaDecodeTime + ptsOffsetFromDts
        },
        end: {
          dts: baseMediaDecodeTime + decodeDuration,
          pts: baseMediaDecodeTime + presentationDuration
        },
        prependedContentDuration: prependedContentDuration,
        baseMediaDecodeTime: baseMediaDecodeTime
      };
    };
    /**
     * Constructs a single-track, ISO BMFF media segment from AAC data
     * events. The output of this stream can be fed to a SourceBuffer
     * configured with a suitable initialization segment.
     * @param track {object} track metadata configuration
     * @param options {object} transmuxer options object
     * @param options.keepOriginalTimestamps {boolean} If true, keep the timestamps
     *        in the source; false to adjust the first segment to start at 0.
     */


    _AudioSegmentStream = function AudioSegmentStream(track, options) {
      var adtsFrames = [],
          sequenceNumber = 0,
          earliestAllowedDts = 0,
          audioAppendStartTs = 0,
          videoBaseMediaDecodeTime = Infinity;
      options = options || {};

      _AudioSegmentStream.prototype.init.call(this);

      this.push = function (data) {
        trackDecodeInfo.collectDtsInfo(track, data);

        if (track) {
          audioProperties.forEach(function (prop) {
            track[prop] = data[prop];
          });
        } // buffer audio data until end() is called


        adtsFrames.push(data);
      };

      this.setEarliestDts = function (earliestDts) {
        earliestAllowedDts = earliestDts;
      };

      this.setVideoBaseMediaDecodeTime = function (baseMediaDecodeTime) {
        videoBaseMediaDecodeTime = baseMediaDecodeTime;
      };

      this.setAudioAppendStart = function (timestamp) {
        audioAppendStartTs = timestamp;
      };

      this.flush = function () {
        var frames, moof, mdat, boxes, frameDuration; // return early if no audio data has been observed

        if (adtsFrames.length === 0) {
          this.trigger('done', 'AudioSegmentStream');
          return;
        }

        frames = audioFrameUtils.trimAdtsFramesByEarliestDts(adtsFrames, track, earliestAllowedDts);
        track.baseMediaDecodeTime = trackDecodeInfo.calculateTrackBaseMediaDecodeTime(track, options.keepOriginalTimestamps);
        audioFrameUtils.prefixWithSilence(track, frames, audioAppendStartTs, videoBaseMediaDecodeTime); // we have to build the index from byte locations to
        // samples (that is, adts frames) in the audio data

        track.samples = audioFrameUtils.generateSampleTable(frames); // concatenate the audio data to constuct the mdat

        mdat = mp4Generator.mdat(audioFrameUtils.concatenateFrameData(frames));
        adtsFrames = [];
        moof = mp4Generator.moof(sequenceNumber, [track]);
        boxes = new Uint8Array(moof.byteLength + mdat.byteLength); // bump the sequence number for next time

        sequenceNumber++;
        boxes.set(moof);
        boxes.set(mdat, moof.byteLength);
        trackDecodeInfo.clearDtsInfo(track);
        frameDuration = Math.ceil(ONE_SECOND_IN_TS$3 * 1024 / track.samplerate); // TODO this check was added to maintain backwards compatibility (particularly with
        // tests) on adding the timingInfo event. However, it seems unlikely that there's a
        // valid use-case where an init segment/data should be triggered without associated
        // frames. Leaving for now, but should be looked into.

        if (frames.length) {
          this.trigger('timingInfo', {
            start: frames[0].pts,
            end: frames[0].pts + frames.length * frameDuration
          });
        }

        this.trigger('data', {
          track: track,
          boxes: boxes
        });
        this.trigger('done', 'AudioSegmentStream');
      };

      this.reset = function () {
        trackDecodeInfo.clearDtsInfo(track);
        adtsFrames = [];
        this.trigger('reset');
      };
    };

    _AudioSegmentStream.prototype = new stream();
    /**
     * Constructs a single-track, ISO BMFF media segment from H264 data
     * events. The output of this stream can be fed to a SourceBuffer
     * configured with a suitable initialization segment.
     * @param track {object} track metadata configuration
     * @param options {object} transmuxer options object
     * @param options.alignGopsAtEnd {boolean} If true, start from the end of the
     *        gopsToAlignWith list when attempting to align gop pts
     * @param options.keepOriginalTimestamps {boolean} If true, keep the timestamps
     *        in the source; false to adjust the first segment to start at 0.
     */

    _VideoSegmentStream = function VideoSegmentStream(track, options) {
      var sequenceNumber = 0,
          nalUnits = [],
          gopsToAlignWith = [],
          config,
          pps;
      options = options || {};

      _VideoSegmentStream.prototype.init.call(this);

      delete track.minPTS;
      this.gopCache_ = [];
      /**
        * Constructs a ISO BMFF segment given H264 nalUnits
        * @param {Object} nalUnit A data event representing a nalUnit
        * @param {String} nalUnit.nalUnitType
        * @param {Object} nalUnit.config Properties for a mp4 track
        * @param {Uint8Array} nalUnit.data The nalUnit bytes
        * @see lib/codecs/h264.js
       **/

      this.push = function (nalUnit) {
        trackDecodeInfo.collectDtsInfo(track, nalUnit); // record the track config

        if (nalUnit.nalUnitType === 'seq_parameter_set_rbsp' && !config) {
          config = nalUnit.config;
          track.sps = [nalUnit.data];
          videoProperties.forEach(function (prop) {
            track[prop] = config[prop];
          }, this);
        }

        if (nalUnit.nalUnitType === 'pic_parameter_set_rbsp' && !pps) {
          pps = nalUnit.data;
          track.pps = [nalUnit.data];
        } // buffer video until flush() is called


        nalUnits.push(nalUnit);
      };
      /**
        * Pass constructed ISO BMFF track and boxes on to the
        * next stream in the pipeline
       **/


      this.flush = function () {
        var frames,
            gopForFusion,
            gops,
            moof,
            mdat,
            boxes,
            prependedContentDuration = 0,
            firstGop,
            lastGop; // Throw away nalUnits at the start of the byte stream until
        // we find the first AUD

        while (nalUnits.length) {
          if (nalUnits[0].nalUnitType === 'access_unit_delimiter_rbsp') {
            break;
          }

          nalUnits.shift();
        } // Return early if no video data has been observed


        if (nalUnits.length === 0) {
          this.resetStream_();
          this.trigger('done', 'VideoSegmentStream');
          return;
        } // Organize the raw nal-units into arrays that represent
        // higher-level constructs such as frames and gops
        // (group-of-pictures)


        frames = frameUtils.groupNalsIntoFrames(nalUnits);
        gops = frameUtils.groupFramesIntoGops(frames); // If the first frame of this fragment is not a keyframe we have
        // a problem since MSE (on Chrome) requires a leading keyframe.
        //
        // We have two approaches to repairing this situation:
        // 1) GOP-FUSION:
        //    This is where we keep track of the GOPS (group-of-pictures)
        //    from previous fragments and attempt to find one that we can
        //    prepend to the current fragment in order to create a valid
        //    fragment.
        // 2) KEYFRAME-PULLING:
        //    Here we search for the first keyframe in the fragment and
        //    throw away all the frames between the start of the fragment
        //    and that keyframe. We then extend the duration and pull the
        //    PTS of the keyframe forward so that it covers the time range
        //    of the frames that were disposed of.
        //
        // #1 is far prefereable over #2 which can cause "stuttering" but
        // requires more things to be just right.

        if (!gops[0][0].keyFrame) {
          // Search for a gop for fusion from our gopCache
          gopForFusion = this.getGopForFusion_(nalUnits[0], track);

          if (gopForFusion) {
            // in order to provide more accurate timing information about the segment, save
            // the number of seconds prepended to the original segment due to GOP fusion
            prependedContentDuration = gopForFusion.duration;
            gops.unshift(gopForFusion); // Adjust Gops' metadata to account for the inclusion of the
            // new gop at the beginning

            gops.byteLength += gopForFusion.byteLength;
            gops.nalCount += gopForFusion.nalCount;
            gops.pts = gopForFusion.pts;
            gops.dts = gopForFusion.dts;
            gops.duration += gopForFusion.duration;
          } else {
            // If we didn't find a candidate gop fall back to keyframe-pulling
            gops = frameUtils.extendFirstKeyFrame(gops);
          }
        } // Trim gops to align with gopsToAlignWith


        if (gopsToAlignWith.length) {
          var alignedGops;

          if (options.alignGopsAtEnd) {
            alignedGops = this.alignGopsAtEnd_(gops);
          } else {
            alignedGops = this.alignGopsAtStart_(gops);
          }

          if (!alignedGops) {
            // save all the nals in the last GOP into the gop cache
            this.gopCache_.unshift({
              gop: gops.pop(),
              pps: track.pps,
              sps: track.sps
            }); // Keep a maximum of 6 GOPs in the cache

            this.gopCache_.length = Math.min(6, this.gopCache_.length); // Clear nalUnits

            nalUnits = []; // return early no gops can be aligned with desired gopsToAlignWith

            this.resetStream_();
            this.trigger('done', 'VideoSegmentStream');
            return;
          } // Some gops were trimmed. clear dts info so minSegmentDts and pts are correct
          // when recalculated before sending off to CoalesceStream


          trackDecodeInfo.clearDtsInfo(track);
          gops = alignedGops;
        }

        trackDecodeInfo.collectDtsInfo(track, gops); // First, we have to build the index from byte locations to
        // samples (that is, frames) in the video data

        track.samples = frameUtils.generateSampleTable(gops); // Concatenate the video data and construct the mdat

        mdat = mp4Generator.mdat(frameUtils.concatenateNalData(gops));
        track.baseMediaDecodeTime = trackDecodeInfo.calculateTrackBaseMediaDecodeTime(track, options.keepOriginalTimestamps);
        this.trigger('processedGopsInfo', gops.map(function (gop) {
          return {
            pts: gop.pts,
            dts: gop.dts,
            byteLength: gop.byteLength
          };
        }));
        firstGop = gops[0];
        lastGop = gops[gops.length - 1];
        this.trigger('segmentTimingInfo', generateVideoSegmentTimingInfo(track.baseMediaDecodeTime, firstGop.dts, firstGop.pts, lastGop.dts + lastGop.duration, lastGop.pts + lastGop.duration, prependedContentDuration));
        this.trigger('timingInfo', {
          start: gops[0].pts,
          end: gops[gops.length - 1].pts + gops[gops.length - 1].duration
        }); // save all the nals in the last GOP into the gop cache

        this.gopCache_.unshift({
          gop: gops.pop(),
          pps: track.pps,
          sps: track.sps
        }); // Keep a maximum of 6 GOPs in the cache

        this.gopCache_.length = Math.min(6, this.gopCache_.length); // Clear nalUnits

        nalUnits = [];
        this.trigger('baseMediaDecodeTime', track.baseMediaDecodeTime);
        this.trigger('timelineStartInfo', track.timelineStartInfo);
        moof = mp4Generator.moof(sequenceNumber, [track]); // it would be great to allocate this array up front instead of
        // throwing away hundreds of media segment fragments

        boxes = new Uint8Array(moof.byteLength + mdat.byteLength); // Bump the sequence number for next time

        sequenceNumber++;
        boxes.set(moof);
        boxes.set(mdat, moof.byteLength);
        this.trigger('data', {
          track: track,
          boxes: boxes
        });
        this.resetStream_(); // Continue with the flush process now

        this.trigger('done', 'VideoSegmentStream');
      };

      this.reset = function () {
        this.resetStream_();
        nalUnits = [];
        this.gopCache_.length = 0;
        gopsToAlignWith.length = 0;
        this.trigger('reset');
      };

      this.resetStream_ = function () {
        trackDecodeInfo.clearDtsInfo(track); // reset config and pps because they may differ across segments
        // for instance, when we are rendition switching

        config = undefined;
        pps = undefined;
      }; // Search for a candidate Gop for gop-fusion from the gop cache and
      // return it or return null if no good candidate was found


      this.getGopForFusion_ = function (nalUnit) {
        var halfSecond = 45000,
            // Half-a-second in a 90khz clock
        allowableOverlap = 10000,
            // About 3 frames @ 30fps
        nearestDistance = Infinity,
            dtsDistance,
            nearestGopObj,
            currentGop,
            currentGopObj,
            i; // Search for the GOP nearest to the beginning of this nal unit

        for (i = 0; i < this.gopCache_.length; i++) {
          currentGopObj = this.gopCache_[i];
          currentGop = currentGopObj.gop; // Reject Gops with different SPS or PPS

          if (!(track.pps && arrayEquals(track.pps[0], currentGopObj.pps[0])) || !(track.sps && arrayEquals(track.sps[0], currentGopObj.sps[0]))) {
            continue;
          } // Reject Gops that would require a negative baseMediaDecodeTime


          if (currentGop.dts < track.timelineStartInfo.dts) {
            continue;
          } // The distance between the end of the gop and the start of the nalUnit


          dtsDistance = nalUnit.dts - currentGop.dts - currentGop.duration; // Only consider GOPS that start before the nal unit and end within
          // a half-second of the nal unit

          if (dtsDistance >= -allowableOverlap && dtsDistance <= halfSecond) {
            // Always use the closest GOP we found if there is more than
            // one candidate
            if (!nearestGopObj || nearestDistance > dtsDistance) {
              nearestGopObj = currentGopObj;
              nearestDistance = dtsDistance;
            }
          }
        }

        if (nearestGopObj) {
          return nearestGopObj.gop;
        }

        return null;
      }; // trim gop list to the first gop found that has a matching pts with a gop in the list
      // of gopsToAlignWith starting from the START of the list


      this.alignGopsAtStart_ = function (gops) {
        var alignIndex, gopIndex, align, gop, byteLength, nalCount, duration, alignedGops;
        byteLength = gops.byteLength;
        nalCount = gops.nalCount;
        duration = gops.duration;
        alignIndex = gopIndex = 0;

        while (alignIndex < gopsToAlignWith.length && gopIndex < gops.length) {
          align = gopsToAlignWith[alignIndex];
          gop = gops[gopIndex];

          if (align.pts === gop.pts) {
            break;
          }

          if (gop.pts > align.pts) {
            // this current gop starts after the current gop we want to align on, so increment
            // align index
            alignIndex++;
            continue;
          } // current gop starts before the current gop we want to align on. so increment gop
          // index


          gopIndex++;
          byteLength -= gop.byteLength;
          nalCount -= gop.nalCount;
          duration -= gop.duration;
        }

        if (gopIndex === 0) {
          // no gops to trim
          return gops;
        }

        if (gopIndex === gops.length) {
          // all gops trimmed, skip appending all gops
          return null;
        }

        alignedGops = gops.slice(gopIndex);
        alignedGops.byteLength = byteLength;
        alignedGops.duration = duration;
        alignedGops.nalCount = nalCount;
        alignedGops.pts = alignedGops[0].pts;
        alignedGops.dts = alignedGops[0].dts;
        return alignedGops;
      }; // trim gop list to the first gop found that has a matching pts with a gop in the list
      // of gopsToAlignWith starting from the END of the list


      this.alignGopsAtEnd_ = function (gops) {
        var alignIndex, gopIndex, align, gop, alignEndIndex, matchFound;
        alignIndex = gopsToAlignWith.length - 1;
        gopIndex = gops.length - 1;
        alignEndIndex = null;
        matchFound = false;

        while (alignIndex >= 0 && gopIndex >= 0) {
          align = gopsToAlignWith[alignIndex];
          gop = gops[gopIndex];

          if (align.pts === gop.pts) {
            matchFound = true;
            break;
          }

          if (align.pts > gop.pts) {
            alignIndex--;
            continue;
          }

          if (alignIndex === gopsToAlignWith.length - 1) {
            // gop.pts is greater than the last alignment candidate. If no match is found
            // by the end of this loop, we still want to append gops that come after this
            // point
            alignEndIndex = gopIndex;
          }

          gopIndex--;
        }

        if (!matchFound && alignEndIndex === null) {
          return null;
        }

        var trimIndex;

        if (matchFound) {
          trimIndex = gopIndex;
        } else {
          trimIndex = alignEndIndex;
        }

        if (trimIndex === 0) {
          return gops;
        }

        var alignedGops = gops.slice(trimIndex);
        var metadata = alignedGops.reduce(function (total, gop) {
          total.byteLength += gop.byteLength;
          total.duration += gop.duration;
          total.nalCount += gop.nalCount;
          return total;
        }, {
          byteLength: 0,
          duration: 0,
          nalCount: 0
        });
        alignedGops.byteLength = metadata.byteLength;
        alignedGops.duration = metadata.duration;
        alignedGops.nalCount = metadata.nalCount;
        alignedGops.pts = alignedGops[0].pts;
        alignedGops.dts = alignedGops[0].dts;
        return alignedGops;
      };

      this.alignGopsWith = function (newGopsToAlignWith) {
        gopsToAlignWith = newGopsToAlignWith;
      };
    };

    _VideoSegmentStream.prototype = new stream();
    /**
     * A Stream that can combine multiple streams (ie. audio & video)
     * into a single output segment for MSE. Also supports audio-only
     * and video-only streams.
     * @param options {object} transmuxer options object
     * @param options.keepOriginalTimestamps {boolean} If true, keep the timestamps
     *        in the source; false to adjust the first segment to start at media timeline start.
     */

    _CoalesceStream = function CoalesceStream(options, metadataStream) {
      // Number of Tracks per output segment
      // If greater than 1, we combine multiple
      // tracks into a single segment
      this.numberOfTracks = 0;
      this.metadataStream = metadataStream;
      options = options || {};

      if (typeof options.remux !== 'undefined') {
        this.remuxTracks = !!options.remux;
      } else {
        this.remuxTracks = true;
      }

      if (typeof options.keepOriginalTimestamps === 'boolean') {
        this.keepOriginalTimestamps = options.keepOriginalTimestamps;
      } else {
        this.keepOriginalTimestamps = false;
      }

      this.pendingTracks = [];
      this.videoTrack = null;
      this.pendingBoxes = [];
      this.pendingCaptions = [];
      this.pendingMetadata = [];
      this.pendingBytes = 0;
      this.emittedTracks = 0;

      _CoalesceStream.prototype.init.call(this); // Take output from multiple


      this.push = function (output) {
        // buffer incoming captions until the associated video segment
        // finishes
        if (output.text) {
          return this.pendingCaptions.push(output);
        } // buffer incoming id3 tags until the final flush


        if (output.frames) {
          return this.pendingMetadata.push(output);
        } // Add this track to the list of pending tracks and store
        // important information required for the construction of
        // the final segment


        this.pendingTracks.push(output.track);
        this.pendingBytes += output.boxes.byteLength; // TODO: is there an issue for this against chrome?
        // We unshift audio and push video because
        // as of Chrome 75 when switching from
        // one init segment to another if the video
        // mdat does not appear after the audio mdat
        // only audio will play for the duration of our transmux.

        if (output.track.type === 'video') {
          this.videoTrack = output.track;
          this.pendingBoxes.push(output.boxes);
        }

        if (output.track.type === 'audio') {
          this.audioTrack = output.track;
          this.pendingBoxes.unshift(output.boxes);
        }
      };
    };

    _CoalesceStream.prototype = new stream();

    _CoalesceStream.prototype.flush = function (flushSource) {
      var offset = 0,
          event = {
        captions: [],
        captionStreams: {},
        metadata: [],
        info: {}
      },
          caption,
          id3,
          initSegment,
          timelineStartPts = 0,
          i;

      if (this.pendingTracks.length < this.numberOfTracks) {
        if (flushSource !== 'VideoSegmentStream' && flushSource !== 'AudioSegmentStream') {
          // Return because we haven't received a flush from a data-generating
          // portion of the segment (meaning that we have only recieved meta-data
          // or captions.)
          return;
        } else if (this.remuxTracks) {
          // Return until we have enough tracks from the pipeline to remux (if we
          // are remuxing audio and video into a single MP4)
          return;
        } else if (this.pendingTracks.length === 0) {
          // In the case where we receive a flush without any data having been
          // received we consider it an emitted track for the purposes of coalescing
          // `done` events.
          // We do this for the case where there is an audio and video track in the
          // segment but no audio data. (seen in several playlists with alternate
          // audio tracks and no audio present in the main TS segments.)
          this.emittedTracks++;

          if (this.emittedTracks >= this.numberOfTracks) {
            this.trigger('done');
            this.emittedTracks = 0;
          }

          return;
        }
      }

      if (this.videoTrack) {
        timelineStartPts = this.videoTrack.timelineStartInfo.pts;
        videoProperties.forEach(function (prop) {
          event.info[prop] = this.videoTrack[prop];
        }, this);
      } else if (this.audioTrack) {
        timelineStartPts = this.audioTrack.timelineStartInfo.pts;
        audioProperties.forEach(function (prop) {
          event.info[prop] = this.audioTrack[prop];
        }, this);
      }

      if (this.videoTrack || this.audioTrack) {
        if (this.pendingTracks.length === 1) {
          event.type = this.pendingTracks[0].type;
        } else {
          event.type = 'combined';
        }

        this.emittedTracks += this.pendingTracks.length;
        initSegment = mp4Generator.initSegment(this.pendingTracks); // Create a new typed array to hold the init segment

        event.initSegment = new Uint8Array(initSegment.byteLength); // Create an init segment containing a moov
        // and track definitions

        event.initSegment.set(initSegment); // Create a new typed array to hold the moof+mdats

        event.data = new Uint8Array(this.pendingBytes); // Append each moof+mdat (one per track) together

        for (i = 0; i < this.pendingBoxes.length; i++) {
          event.data.set(this.pendingBoxes[i], offset);
          offset += this.pendingBoxes[i].byteLength;
        } // Translate caption PTS times into second offsets to match the
        // video timeline for the segment, and add track info


        for (i = 0; i < this.pendingCaptions.length; i++) {
          caption = this.pendingCaptions[i];
          caption.startTime = clock.metadataTsToSeconds(caption.startPts, timelineStartPts, this.keepOriginalTimestamps);
          caption.endTime = clock.metadataTsToSeconds(caption.endPts, timelineStartPts, this.keepOriginalTimestamps);
          event.captionStreams[caption.stream] = true;
          event.captions.push(caption);
        } // Translate ID3 frame PTS times into second offsets to match the
        // video timeline for the segment


        for (i = 0; i < this.pendingMetadata.length; i++) {
          id3 = this.pendingMetadata[i];
          id3.cueTime = clock.metadataTsToSeconds(id3.pts, timelineStartPts, this.keepOriginalTimestamps);
          event.metadata.push(id3);
        } // We add this to every single emitted segment even though we only need
        // it for the first


        event.metadata.dispatchType = this.metadataStream.dispatchType; // Reset stream state

        this.pendingTracks.length = 0;
        this.videoTrack = null;
        this.pendingBoxes.length = 0;
        this.pendingCaptions.length = 0;
        this.pendingBytes = 0;
        this.pendingMetadata.length = 0; // Emit the built segment
        // We include captions and ID3 tags for backwards compatibility,
        // ideally we should send only video and audio in the data event

        this.trigger('data', event); // Emit each caption to the outside world
        // Ideally, this would happen immediately on parsing captions,
        // but we need to ensure that video data is sent back first
        // so that caption timing can be adjusted to match video timing

        for (i = 0; i < event.captions.length; i++) {
          caption = event.captions[i];
          this.trigger('caption', caption);
        } // Emit each id3 tag to the outside world
        // Ideally, this would happen immediately on parsing the tag,
        // but we need to ensure that video data is sent back first
        // so that ID3 frame timing can be adjusted to match video timing


        for (i = 0; i < event.metadata.length; i++) {
          id3 = event.metadata[i];
          this.trigger('id3Frame', id3);
        }
      } // Only emit `done` if all tracks have been flushed and emitted


      if (this.emittedTracks >= this.numberOfTracks) {
        this.trigger('done');
        this.emittedTracks = 0;
      }
    };

    _CoalesceStream.prototype.setRemux = function (val) {
      this.remuxTracks = val;
    };
    /**
     * A Stream that expects MP2T binary data as input and produces
     * corresponding media segments, suitable for use with Media Source
     * Extension (MSE) implementations that support the ISO BMFF byte
     * stream format, like Chrome.
     */


    _Transmuxer = function Transmuxer(options) {
      var self = this,
          hasFlushed = true,
          videoTrack,
          audioTrack;

      _Transmuxer.prototype.init.call(this);

      options = options || {};
      this.baseMediaDecodeTime = options.baseMediaDecodeTime || 0;
      this.transmuxPipeline_ = {};

      this.setupAacPipeline = function () {
        var pipeline = {};
        this.transmuxPipeline_ = pipeline;
        pipeline.type = 'aac';
        pipeline.metadataStream = new m2ts_1.MetadataStream(); // set up the parsing pipeline

        pipeline.aacStream = new aac();
        pipeline.audioTimestampRolloverStream = new m2ts_1.TimestampRolloverStream('audio');
        pipeline.timedMetadataTimestampRolloverStream = new m2ts_1.TimestampRolloverStream('timed-metadata');
        pipeline.adtsStream = new adts();
        pipeline.coalesceStream = new _CoalesceStream(options, pipeline.metadataStream);
        pipeline.headOfPipeline = pipeline.aacStream;
        pipeline.aacStream.pipe(pipeline.audioTimestampRolloverStream).pipe(pipeline.adtsStream);
        pipeline.aacStream.pipe(pipeline.timedMetadataTimestampRolloverStream).pipe(pipeline.metadataStream).pipe(pipeline.coalesceStream);
        pipeline.metadataStream.on('timestamp', function (frame) {
          pipeline.aacStream.setTimestamp(frame.timeStamp);
        });
        pipeline.aacStream.on('data', function (data) {
          if (data.type !== 'timed-metadata' && data.type !== 'audio' || pipeline.audioSegmentStream) {
            return;
          }

          audioTrack = audioTrack || {
            timelineStartInfo: {
              baseMediaDecodeTime: self.baseMediaDecodeTime
            },
            codec: 'adts',
            type: 'audio'
          }; // hook up the audio segment stream to the first track with aac data

          pipeline.coalesceStream.numberOfTracks++;
          pipeline.audioSegmentStream = new _AudioSegmentStream(audioTrack, options);
          pipeline.audioSegmentStream.on('timingInfo', self.trigger.bind(self, 'audioTimingInfo')); // Set up the final part of the audio pipeline

          pipeline.adtsStream.pipe(pipeline.audioSegmentStream).pipe(pipeline.coalesceStream); // emit pmt info

          self.trigger('trackinfo', {
            hasAudio: !!audioTrack,
            hasVideo: !!videoTrack
          });
        }); // Re-emit any data coming from the coalesce stream to the outside world

        pipeline.coalesceStream.on('data', this.trigger.bind(this, 'data')); // Let the consumer know we have finished flushing the entire pipeline

        pipeline.coalesceStream.on('done', this.trigger.bind(this, 'done'));
      };

      this.setupTsPipeline = function () {
        var pipeline = {};
        this.transmuxPipeline_ = pipeline;
        pipeline.type = 'ts';
        pipeline.metadataStream = new m2ts_1.MetadataStream(); // set up the parsing pipeline

        pipeline.packetStream = new m2ts_1.TransportPacketStream();
        pipeline.parseStream = new m2ts_1.TransportParseStream();
        pipeline.elementaryStream = new m2ts_1.ElementaryStream();
        pipeline.timestampRolloverStream = new m2ts_1.TimestampRolloverStream();
        pipeline.adtsStream = new adts();
        pipeline.h264Stream = new H264Stream$1();
        pipeline.captionStream = new m2ts_1.CaptionStream();
        pipeline.coalesceStream = new _CoalesceStream(options, pipeline.metadataStream);
        pipeline.headOfPipeline = pipeline.packetStream; // disassemble MPEG2-TS packets into elementary streams

        pipeline.packetStream.pipe(pipeline.parseStream).pipe(pipeline.elementaryStream).pipe(pipeline.timestampRolloverStream); // !!THIS ORDER IS IMPORTANT!!
        // demux the streams

        pipeline.timestampRolloverStream.pipe(pipeline.h264Stream);
        pipeline.timestampRolloverStream.pipe(pipeline.adtsStream);
        pipeline.timestampRolloverStream.pipe(pipeline.metadataStream).pipe(pipeline.coalesceStream); // Hook up CEA-608/708 caption stream

        pipeline.h264Stream.pipe(pipeline.captionStream).pipe(pipeline.coalesceStream);
        pipeline.elementaryStream.on('data', function (data) {
          var i;

          if (data.type === 'metadata') {
            i = data.tracks.length; // scan the tracks listed in the metadata

            while (i--) {
              if (!videoTrack && data.tracks[i].type === 'video') {
                videoTrack = data.tracks[i];
                videoTrack.timelineStartInfo.baseMediaDecodeTime = self.baseMediaDecodeTime;
              } else if (!audioTrack && data.tracks[i].type === 'audio') {
                audioTrack = data.tracks[i];
                audioTrack.timelineStartInfo.baseMediaDecodeTime = self.baseMediaDecodeTime;
              }
            } // hook up the video segment stream to the first track with h264 data


            if (videoTrack && !pipeline.videoSegmentStream) {
              pipeline.coalesceStream.numberOfTracks++;
              pipeline.videoSegmentStream = new _VideoSegmentStream(videoTrack, options);
              pipeline.videoSegmentStream.on('timelineStartInfo', function (timelineStartInfo) {
                // When video emits timelineStartInfo data after a flush, we forward that
                // info to the AudioSegmentStream, if it exists, because video timeline
                // data takes precedence.  Do not do this if keepOriginalTimestamps is set,
                // because this is a particularly subtle form of timestamp alteration.
                if (audioTrack && !options.keepOriginalTimestamps) {
                  audioTrack.timelineStartInfo = timelineStartInfo; // On the first segment we trim AAC frames that exist before the
                  // very earliest DTS we have seen in video because Chrome will
                  // interpret any video track with a baseMediaDecodeTime that is
                  // non-zero as a gap.

                  pipeline.audioSegmentStream.setEarliestDts(timelineStartInfo.dts - self.baseMediaDecodeTime);
                }
              });
              pipeline.videoSegmentStream.on('processedGopsInfo', self.trigger.bind(self, 'gopInfo'));
              pipeline.videoSegmentStream.on('segmentTimingInfo', self.trigger.bind(self, 'videoSegmentTimingInfo'));
              pipeline.videoSegmentStream.on('baseMediaDecodeTime', function (baseMediaDecodeTime) {
                if (audioTrack) {
                  pipeline.audioSegmentStream.setVideoBaseMediaDecodeTime(baseMediaDecodeTime);
                }
              });
              pipeline.videoSegmentStream.on('timingInfo', self.trigger.bind(self, 'videoTimingInfo')); // Set up the final part of the video pipeline

              pipeline.h264Stream.pipe(pipeline.videoSegmentStream).pipe(pipeline.coalesceStream);
            }

            if (audioTrack && !pipeline.audioSegmentStream) {
              // hook up the audio segment stream to the first track with aac data
              pipeline.coalesceStream.numberOfTracks++;
              pipeline.audioSegmentStream = new _AudioSegmentStream(audioTrack, options);
              pipeline.audioSegmentStream.on('timingInfo', self.trigger.bind(self, 'audioTimingInfo')); // Set up the final part of the audio pipeline

              pipeline.adtsStream.pipe(pipeline.audioSegmentStream).pipe(pipeline.coalesceStream);
            } // emit pmt info


            self.trigger('trackinfo', {
              hasAudio: !!audioTrack,
              hasVideo: !!videoTrack
            });
          }
        }); // Re-emit any data coming from the coalesce stream to the outside world

        pipeline.coalesceStream.on('data', this.trigger.bind(this, 'data'));
        pipeline.coalesceStream.on('id3Frame', function (id3Frame) {
          id3Frame.dispatchType = pipeline.metadataStream.dispatchType;
          self.trigger('id3Frame', id3Frame);
        });
        pipeline.coalesceStream.on('caption', this.trigger.bind(this, 'caption')); // Let the consumer know we have finished flushing the entire pipeline

        pipeline.coalesceStream.on('done', this.trigger.bind(this, 'done'));
      }; // hook up the segment streams once track metadata is delivered


      this.setBaseMediaDecodeTime = function (baseMediaDecodeTime) {
        var pipeline = this.transmuxPipeline_;

        if (!options.keepOriginalTimestamps) {
          this.baseMediaDecodeTime = baseMediaDecodeTime;
        }

        if (audioTrack) {
          audioTrack.timelineStartInfo.dts = undefined;
          audioTrack.timelineStartInfo.pts = undefined;
          trackDecodeInfo.clearDtsInfo(audioTrack);

          if (pipeline.audioTimestampRolloverStream) {
            pipeline.audioTimestampRolloverStream.discontinuity();
          }
        }

        if (videoTrack) {
          if (pipeline.videoSegmentStream) {
            pipeline.videoSegmentStream.gopCache_ = [];
          }

          videoTrack.timelineStartInfo.dts = undefined;
          videoTrack.timelineStartInfo.pts = undefined;
          trackDecodeInfo.clearDtsInfo(videoTrack);
          pipeline.captionStream.reset();
        }

        if (pipeline.timestampRolloverStream) {
          pipeline.timestampRolloverStream.discontinuity();
        }
      };

      this.setAudioAppendStart = function (timestamp) {
        if (audioTrack) {
          this.transmuxPipeline_.audioSegmentStream.setAudioAppendStart(timestamp);
        }
      };

      this.setRemux = function (val) {
        var pipeline = this.transmuxPipeline_;
        options.remux = val;

        if (pipeline && pipeline.coalesceStream) {
          pipeline.coalesceStream.setRemux(val);
        }
      };

      this.alignGopsWith = function (gopsToAlignWith) {
        if (videoTrack && this.transmuxPipeline_.videoSegmentStream) {
          this.transmuxPipeline_.videoSegmentStream.alignGopsWith(gopsToAlignWith);
        }
      }; // feed incoming data to the front of the parsing pipeline


      this.push = function (data) {
        if (hasFlushed) {
          var isAac = isLikelyAacData$1(data);

          if (isAac && this.transmuxPipeline_.type !== 'aac') {
            this.setupAacPipeline();
          } else if (!isAac && this.transmuxPipeline_.type !== 'ts') {
            this.setupTsPipeline();
          }

          hasFlushed = false;
        }

        this.transmuxPipeline_.headOfPipeline.push(data);
      }; // flush any buffered data


      this.flush = function () {
        hasFlushed = true; // Start at the top of the pipeline and flush all pending work

        this.transmuxPipeline_.headOfPipeline.flush();
      };

      this.endTimeline = function () {
        this.transmuxPipeline_.headOfPipeline.endTimeline();
      };

      this.reset = function () {
        if (this.transmuxPipeline_.headOfPipeline) {
          this.transmuxPipeline_.headOfPipeline.reset();
        }
      }; // Caption data has to be reset when seeking outside buffered range


      this.resetCaptions = function () {
        if (this.transmuxPipeline_.captionStream) {
          this.transmuxPipeline_.captionStream.reset();
        }
      };
    };

    _Transmuxer.prototype = new stream();
    var transmuxer = {
      Transmuxer: _Transmuxer,
      VideoSegmentStream: _VideoSegmentStream,
      AudioSegmentStream: _AudioSegmentStream,
      AUDIO_PROPERTIES: audioProperties,
      VIDEO_PROPERTIES: videoProperties,
      // exported for testing
      generateVideoSegmentTimingInfo: generateVideoSegmentTimingInfo
    };
    var transmuxer_1 = transmuxer.Transmuxer;
    /**
     * mux.js
     *
     * Copyright (c) Brightcove
     * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
     */

    var codecs = {
      Adts: adts,
      h264: h264
    };
    var ONE_SECOND_IN_TS$4 = clock.ONE_SECOND_IN_TS;
    /**
     * Constructs a single-track, ISO BMFF media segment from AAC data
     * events. The output of this stream can be fed to a SourceBuffer
     * configured with a suitable initialization segment.
     */

    var AudioSegmentStream$1 = function AudioSegmentStream$1(track, options) {
      var adtsFrames = [],
          sequenceNumber = 0,
          earliestAllowedDts = 0,
          audioAppendStartTs = 0,
          videoBaseMediaDecodeTime = Infinity,
          segmentStartPts = null,
          segmentEndPts = null;
      options = options || {};
      AudioSegmentStream$1.prototype.init.call(this);

      this.push = function (data) {
        trackDecodeInfo.collectDtsInfo(track, data);

        if (track) {
          audioProperties.forEach(function (prop) {
            track[prop] = data[prop];
          });
        } // buffer audio data until end() is called


        adtsFrames.push(data);
      };

      this.setEarliestDts = function (earliestDts) {
        earliestAllowedDts = earliestDts;
      };

      this.setVideoBaseMediaDecodeTime = function (baseMediaDecodeTime) {
        videoBaseMediaDecodeTime = baseMediaDecodeTime;
      };

      this.setAudioAppendStart = function (timestamp) {
        audioAppendStartTs = timestamp;
      };

      this.processFrames_ = function () {
        var frames, moof, mdat, boxes, timingInfo; // return early if no audio data has been observed

        if (adtsFrames.length === 0) {
          return;
        }

        frames = audioFrameUtils.trimAdtsFramesByEarliestDts(adtsFrames, track, earliestAllowedDts);

        if (frames.length === 0) {
          // return early if the frames are all after the earliest allowed DTS
          // TODO should we clear the adtsFrames?
          return;
        }

        track.baseMediaDecodeTime = trackDecodeInfo.calculateTrackBaseMediaDecodeTime(track, options.keepOriginalTimestamps);
        audioFrameUtils.prefixWithSilence(track, frames, audioAppendStartTs, videoBaseMediaDecodeTime); // we have to build the index from byte locations to
        // samples (that is, adts frames) in the audio data

        track.samples = audioFrameUtils.generateSampleTable(frames); // concatenate the audio data to constuct the mdat

        mdat = mp4Generator.mdat(audioFrameUtils.concatenateFrameData(frames));
        adtsFrames = [];
        moof = mp4Generator.moof(sequenceNumber, [track]); // bump the sequence number for next time

        sequenceNumber++;
        track.initSegment = mp4Generator.initSegment([track]); // it would be great to allocate this array up front instead of
        // throwing away hundreds of media segment fragments

        boxes = new Uint8Array(moof.byteLength + mdat.byteLength);
        boxes.set(moof);
        boxes.set(mdat, moof.byteLength);
        trackDecodeInfo.clearDtsInfo(track);

        if (segmentStartPts === null) {
          segmentEndPts = segmentStartPts = frames[0].pts;
        }

        segmentEndPts += frames.length * (ONE_SECOND_IN_TS$4 * 1024 / track.samplerate);
        timingInfo = {
          start: segmentStartPts
        };
        this.trigger('timingInfo', timingInfo);
        this.trigger('data', {
          track: track,
          boxes: boxes
        });
      };

      this.flush = function () {
        this.processFrames_(); // trigger final timing info

        this.trigger('timingInfo', {
          start: segmentStartPts,
          end: segmentEndPts
        });
        this.resetTiming_();
        this.trigger('done', 'AudioSegmentStream');
      };

      this.partialFlush = function () {
        this.processFrames_();
        this.trigger('partialdone', 'AudioSegmentStream');
      };

      this.endTimeline = function () {
        this.flush();
        this.trigger('endedtimeline', 'AudioSegmentStream');
      };

      this.resetTiming_ = function () {
        trackDecodeInfo.clearDtsInfo(track);
        segmentStartPts = null;
        segmentEndPts = null;
      };

      this.reset = function () {
        this.resetTiming_();
        adtsFrames = [];
        this.trigger('reset');
      };
    };

    AudioSegmentStream$1.prototype = new stream();
    var audioSegmentStream = AudioSegmentStream$1;

    var VideoSegmentStream$1 = function VideoSegmentStream$1(track, options) {
      var sequenceNumber = 0,
          nalUnits = [],
          frameCache = [],
          // gopsToAlignWith = [],
      config,
          pps,
          segmentStartPts = null,
          segmentEndPts = null,
          gops,
          ensureNextFrameIsKeyFrame = true;
      options = options || {};
      VideoSegmentStream$1.prototype.init.call(this);

      this.push = function (nalUnit) {
        trackDecodeInfo.collectDtsInfo(track, nalUnit);

        if (typeof track.timelineStartInfo.dts === 'undefined') {
          track.timelineStartInfo.dts = nalUnit.dts;
        } // record the track config


        if (nalUnit.nalUnitType === 'seq_parameter_set_rbsp' && !config) {
          config = nalUnit.config;
          track.sps = [nalUnit.data];
          videoProperties.forEach(function (prop) {
            track[prop] = config[prop];
          }, this);
        }

        if (nalUnit.nalUnitType === 'pic_parameter_set_rbsp' && !pps) {
          pps = nalUnit.data;
          track.pps = [nalUnit.data];
        } // buffer video until flush() is called


        nalUnits.push(nalUnit);
      };

      this.processNals_ = function (cacheLastFrame) {
        var i;
        nalUnits = frameCache.concat(nalUnits); // Throw away nalUnits at the start of the byte stream until
        // we find the first AUD

        while (nalUnits.length) {
          if (nalUnits[0].nalUnitType === 'access_unit_delimiter_rbsp') {
            break;
          }

          nalUnits.shift();
        } // Return early if no video data has been observed


        if (nalUnits.length === 0) {
          return;
        }

        var frames = frameUtils.groupNalsIntoFrames(nalUnits);

        if (!frames.length) {
          return;
        } // note that the frame cache may also protect us from cases where we haven't
        // pushed data for the entire first or last frame yet


        frameCache = frames[frames.length - 1];

        if (cacheLastFrame) {
          frames.pop();
          frames.duration -= frameCache.duration;
          frames.nalCount -= frameCache.length;
          frames.byteLength -= frameCache.byteLength;
        }

        if (!frames.length) {
          nalUnits = [];
          return;
        }

        this.trigger('timelineStartInfo', track.timelineStartInfo);

        if (ensureNextFrameIsKeyFrame) {
          gops = frameUtils.groupFramesIntoGops(frames);

          if (!gops[0][0].keyFrame) {
            gops = frameUtils.extendFirstKeyFrame(gops);

            if (!gops[0][0].keyFrame) {
              // we haven't yet gotten a key frame, so reset nal units to wait for more nal
              // units
              nalUnits = [].concat.apply([], frames).concat(frameCache);
              frameCache = [];
              return;
            }

            frames = [].concat.apply([], gops);
            frames.duration = gops.duration;
          }

          ensureNextFrameIsKeyFrame = false;
        }

        if (segmentStartPts === null) {
          segmentStartPts = frames[0].pts;
          segmentEndPts = segmentStartPts;
        }

        segmentEndPts += frames.duration;
        this.trigger('timingInfo', {
          start: segmentStartPts,
          end: segmentEndPts
        });

        for (i = 0; i < frames.length; i++) {
          var frame = frames[i];
          track.samples = frameUtils.generateSampleTableForFrame(frame);
          var mdat = mp4Generator.mdat(frameUtils.concatenateNalDataForFrame(frame));
          trackDecodeInfo.clearDtsInfo(track);
          trackDecodeInfo.collectDtsInfo(track, frame);
          track.baseMediaDecodeTime = trackDecodeInfo.calculateTrackBaseMediaDecodeTime(track, options.keepOriginalTimestamps);
          var moof = mp4Generator.moof(sequenceNumber, [track]);
          sequenceNumber++;
          track.initSegment = mp4Generator.initSegment([track]);
          var boxes = new Uint8Array(moof.byteLength + mdat.byteLength);
          boxes.set(moof);
          boxes.set(mdat, moof.byteLength);
          this.trigger('data', {
            track: track,
            boxes: boxes,
            sequence: sequenceNumber,
            videoFrameDts: frame.dts,
            videoFramePts: frame.pts
          });
        }

        nalUnits = [];
      };

      this.resetTimingAndConfig_ = function () {
        config = undefined;
        pps = undefined;
        segmentStartPts = null;
        segmentEndPts = null;
      };

      this.partialFlush = function () {
        this.processNals_(true);
        this.trigger('partialdone', 'VideoSegmentStream');
      };

      this.flush = function () {
        this.processNals_(false); // reset config and pps because they may differ across segments
        // for instance, when we are rendition switching

        this.resetTimingAndConfig_();
        this.trigger('done', 'VideoSegmentStream');
      };

      this.endTimeline = function () {
        this.flush();
        this.trigger('endedtimeline', 'VideoSegmentStream');
      };

      this.reset = function () {
        this.resetTimingAndConfig_();
        frameCache = [];
        nalUnits = [];
        ensureNextFrameIsKeyFrame = true;
        this.trigger('reset');
      };
    };

    VideoSegmentStream$1.prototype = new stream();
    var videoSegmentStream = VideoSegmentStream$1;
    var isLikelyAacData$2 = utils.isLikelyAacData;

    var createPipeline = function createPipeline(object) {
      object.prototype = new stream();
      object.prototype.init.call(object);
      return object;
    };

    var tsPipeline = function tsPipeline(options) {
      var pipeline = {
        type: 'ts',
        tracks: {
          audio: null,
          video: null
        },
        packet: new m2ts_1.TransportPacketStream(),
        parse: new m2ts_1.TransportParseStream(),
        elementary: new m2ts_1.ElementaryStream(),
        timestampRollover: new m2ts_1.TimestampRolloverStream(),
        adts: new codecs.Adts(),
        h264: new codecs.h264.H264Stream(),
        captionStream: new m2ts_1.CaptionStream(),
        metadataStream: new m2ts_1.MetadataStream()
      };
      pipeline.headOfPipeline = pipeline.packet; // Transport Stream

      pipeline.packet.pipe(pipeline.parse).pipe(pipeline.elementary).pipe(pipeline.timestampRollover); // H264

      pipeline.timestampRollover.pipe(pipeline.h264); // Hook up CEA-608/708 caption stream

      pipeline.h264.pipe(pipeline.captionStream);
      pipeline.timestampRollover.pipe(pipeline.metadataStream); // ADTS

      pipeline.timestampRollover.pipe(pipeline.adts);
      pipeline.elementary.on('data', function (data) {
        if (data.type !== 'metadata') {
          return;
        }

        for (var i = 0; i < data.tracks.length; i++) {
          if (!pipeline.tracks[data.tracks[i].type]) {
            pipeline.tracks[data.tracks[i].type] = data.tracks[i];
            pipeline.tracks[data.tracks[i].type].timelineStartInfo.baseMediaDecodeTime = options.baseMediaDecodeTime;
          }
        }

        if (pipeline.tracks.video && !pipeline.videoSegmentStream) {
          pipeline.videoSegmentStream = new videoSegmentStream(pipeline.tracks.video, options);
          pipeline.videoSegmentStream.on('timelineStartInfo', function (timelineStartInfo) {
            if (pipeline.tracks.audio && !options.keepOriginalTimestamps) {
              pipeline.audioSegmentStream.setEarliestDts(timelineStartInfo.dts - options.baseMediaDecodeTime);
            }
          });
          pipeline.videoSegmentStream.on('timingInfo', pipeline.trigger.bind(pipeline, 'videoTimingInfo'));
          pipeline.videoSegmentStream.on('data', function (data) {
            pipeline.trigger('data', {
              type: 'video',
              data: data
            });
          });
          pipeline.videoSegmentStream.on('done', pipeline.trigger.bind(pipeline, 'done'));
          pipeline.videoSegmentStream.on('partialdone', pipeline.trigger.bind(pipeline, 'partialdone'));
          pipeline.videoSegmentStream.on('endedtimeline', pipeline.trigger.bind(pipeline, 'endedtimeline'));
          pipeline.h264.pipe(pipeline.videoSegmentStream);
        }

        if (pipeline.tracks.audio && !pipeline.audioSegmentStream) {
          pipeline.audioSegmentStream = new audioSegmentStream(pipeline.tracks.audio, options);
          pipeline.audioSegmentStream.on('data', function (data) {
            pipeline.trigger('data', {
              type: 'audio',
              data: data
            });
          });
          pipeline.audioSegmentStream.on('done', pipeline.trigger.bind(pipeline, 'done'));
          pipeline.audioSegmentStream.on('partialdone', pipeline.trigger.bind(pipeline, 'partialdone'));
          pipeline.audioSegmentStream.on('endedtimeline', pipeline.trigger.bind(pipeline, 'endedtimeline'));
          pipeline.audioSegmentStream.on('timingInfo', pipeline.trigger.bind(pipeline, 'audioTimingInfo'));
          pipeline.adts.pipe(pipeline.audioSegmentStream);
        } // emit pmt info


        pipeline.trigger('trackinfo', {
          hasAudio: !!pipeline.tracks.audio,
          hasVideo: !!pipeline.tracks.video
        });
      });
      pipeline.captionStream.on('data', function (caption) {
        var timelineStartPts;

        if (pipeline.tracks.video) {
          timelineStartPts = pipeline.tracks.video.timelineStartInfo.pts || 0;
        } else {
          // This will only happen if we encounter caption packets before
          // video data in a segment. This is an unusual/unlikely scenario,
          // so we assume the timeline starts at zero for now.
          timelineStartPts = 0;
        } // Translate caption PTS times into second offsets into the
        // video timeline for the segment


        caption.startTime = clock.metadataTsToSeconds(caption.startPts, timelineStartPts, options.keepOriginalTimestamps);
        caption.endTime = clock.metadataTsToSeconds(caption.endPts, timelineStartPts, options.keepOriginalTimestamps);
        pipeline.trigger('caption', caption);
      });
      pipeline = createPipeline(pipeline);
      pipeline.metadataStream.on('data', pipeline.trigger.bind(pipeline, 'id3Frame'));
      return pipeline;
    };

    var aacPipeline = function aacPipeline(options) {
      var pipeline = {
        type: 'aac',
        tracks: {
          audio: null
        },
        metadataStream: new m2ts_1.MetadataStream(),
        aacStream: new aac(),
        audioRollover: new m2ts_1.TimestampRolloverStream('audio'),
        timedMetadataRollover: new m2ts_1.TimestampRolloverStream('timed-metadata'),
        adtsStream: new adts(true)
      }; // set up the parsing pipeline

      pipeline.headOfPipeline = pipeline.aacStream;
      pipeline.aacStream.pipe(pipeline.audioRollover).pipe(pipeline.adtsStream);
      pipeline.aacStream.pipe(pipeline.timedMetadataRollover).pipe(pipeline.metadataStream);
      pipeline.metadataStream.on('timestamp', function (frame) {
        pipeline.aacStream.setTimestamp(frame.timeStamp);
      });
      pipeline.aacStream.on('data', function (data) {
        if (data.type !== 'timed-metadata' && data.type !== 'audio' || pipeline.audioSegmentStream) {
          return;
        }

        pipeline.tracks.audio = pipeline.tracks.audio || {
          timelineStartInfo: {
            baseMediaDecodeTime: options.baseMediaDecodeTime
          },
          codec: 'adts',
          type: 'audio'
        }; // hook up the audio segment stream to the first track with aac data

        pipeline.audioSegmentStream = new audioSegmentStream(pipeline.tracks.audio, options);
        pipeline.audioSegmentStream.on('data', function (data) {
          pipeline.trigger('data', {
            type: 'audio',
            data: data
          });
        });
        pipeline.audioSegmentStream.on('partialdone', pipeline.trigger.bind(pipeline, 'partialdone'));
        pipeline.audioSegmentStream.on('done', pipeline.trigger.bind(pipeline, 'done'));
        pipeline.audioSegmentStream.on('endedtimeline', pipeline.trigger.bind(pipeline, 'endedtimeline'));
        pipeline.audioSegmentStream.on('timingInfo', pipeline.trigger.bind(pipeline, 'audioTimingInfo')); // Set up the final part of the audio pipeline

        pipeline.adtsStream.pipe(pipeline.audioSegmentStream);
        pipeline.trigger('trackinfo', {
          hasAudio: !!pipeline.tracks.audio,
          hasVideo: !!pipeline.tracks.video
        });
      }); // set the pipeline up as a stream before binding to get access to the trigger function

      pipeline = createPipeline(pipeline);
      pipeline.metadataStream.on('data', pipeline.trigger.bind(pipeline, 'id3Frame'));
      return pipeline;
    };

    var setupPipelineListeners = function setupPipelineListeners(pipeline, transmuxer) {
      pipeline.on('data', transmuxer.trigger.bind(transmuxer, 'data'));
      pipeline.on('done', transmuxer.trigger.bind(transmuxer, 'done'));
      pipeline.on('partialdone', transmuxer.trigger.bind(transmuxer, 'partialdone'));
      pipeline.on('endedtimeline', transmuxer.trigger.bind(transmuxer, 'endedtimeline'));
      pipeline.on('audioTimingInfo', transmuxer.trigger.bind(transmuxer, 'audioTimingInfo'));
      pipeline.on('videoTimingInfo', transmuxer.trigger.bind(transmuxer, 'videoTimingInfo'));
      pipeline.on('trackinfo', transmuxer.trigger.bind(transmuxer, 'trackinfo'));
      pipeline.on('id3Frame', function (event) {
        // add this to every single emitted segment even though it's only needed for the first
        event.dispatchType = pipeline.metadataStream.dispatchType; // keep original time, can be adjusted if needed at a higher level

        event.cueTime = clock.videoTsToSeconds(event.pts);
        transmuxer.trigger('id3Frame', event);
      });
      pipeline.on('caption', function (event) {
        transmuxer.trigger('caption', event);
      });
    };

    var Transmuxer$1 = function Transmuxer$1(options) {
      var pipeline = null,
          hasFlushed = true;
      options = options || {};
      Transmuxer$1.prototype.init.call(this);
      options.baseMediaDecodeTime = options.baseMediaDecodeTime || 0;

      this.push = function (bytes) {
        if (hasFlushed) {
          var isAac = isLikelyAacData$2(bytes);

          if (isAac && (!pipeline || pipeline.type !== 'aac')) {
            pipeline = aacPipeline(options);
            setupPipelineListeners(pipeline, this);
          } else if (!isAac && (!pipeline || pipeline.type !== 'ts')) {
            pipeline = tsPipeline(options);
            setupPipelineListeners(pipeline, this);
          }

          hasFlushed = false;
        }

        pipeline.headOfPipeline.push(bytes);
      };

      this.flush = function () {
        if (!pipeline) {
          return;
        }

        hasFlushed = true;
        pipeline.headOfPipeline.flush();
      };

      this.partialFlush = function () {
        if (!pipeline) {
          return;
        }

        pipeline.headOfPipeline.partialFlush();
      };

      this.endTimeline = function () {
        if (!pipeline) {
          return;
        }

        pipeline.headOfPipeline.endTimeline();
      };

      this.reset = function () {
        if (!pipeline) {
          return;
        }

        pipeline.headOfPipeline.reset();
      };

      this.setBaseMediaDecodeTime = function (baseMediaDecodeTime) {
        if (!options.keepOriginalTimestamps) {
          options.baseMediaDecodeTime = baseMediaDecodeTime;
        }

        if (!pipeline) {
          return;
        }

        if (pipeline.tracks.audio) {
          pipeline.tracks.audio.timelineStartInfo.dts = undefined;
          pipeline.tracks.audio.timelineStartInfo.pts = undefined;
          trackDecodeInfo.clearDtsInfo(pipeline.tracks.audio);

          if (pipeline.audioRollover) {
            pipeline.audioRollover.discontinuity();
          }
        }

        if (pipeline.tracks.video) {
          if (pipeline.videoSegmentStream) {
            pipeline.videoSegmentStream.gopCache_ = [];
          }

          pipeline.tracks.video.timelineStartInfo.dts = undefined;
          pipeline.tracks.video.timelineStartInfo.pts = undefined;
          trackDecodeInfo.clearDtsInfo(pipeline.tracks.video); // pipeline.captionStream.reset();
        }

        if (pipeline.timestampRollover) {
          pipeline.timestampRollover.discontinuity();
        }
      };

      this.setRemux = function (val) {
        options.remux = val;

        if (pipeline && pipeline.coalesceStream) {
          pipeline.coalesceStream.setRemux(val);
        }
      };

      this.setAudioAppendStart = function (audioAppendStart) {
        if (!pipeline || !pipeline.tracks.audio || !pipeline.audioSegmentStream) {
          return;
        }

        pipeline.audioSegmentStream.setAudioAppendStart(audioAppendStart);
      }; // TODO GOP alignment support
      // Support may be a bit trickier than with full segment appends, as GOPs may be split
      // and processed in a more granular fashion


      this.alignGopsWith = function (gopsToAlignWith) {
        return;
      };
    };

    Transmuxer$1.prototype = new stream();
    var transmuxer$1 = Transmuxer$1;
    /**
     * mux.js
     *
     * Copyright (c) Brightcove
     * Licensed Apache-2.0 https://github.com/videojs/mux.js/blob/master/LICENSE
     */

    var toUnsigned = function toUnsigned(value) {
      return value >>> 0;
    };

    var toHexString = function toHexString(value) {
      return ('00' + value.toString(16)).slice(-2);
    };

    var bin = {
      toUnsigned: toUnsigned,
      toHexString: toHexString
    };

    var parseType$1 = function parseType$1(buffer) {
      var result = '';
      result += String.fromCharCode(buffer[0]);
      result += String.fromCharCode(buffer[1]);
      result += String.fromCharCode(buffer[2]);
      result += String.fromCharCode(buffer[3]);
      return result;
    };

    var parseType_1 = parseType$1;
    var toUnsigned$1 = bin.toUnsigned;

    var findBox = function findBox(data, path) {
      var results = [],
          i,
          size,
          type,
          end,
          subresults;

      if (!path.length) {
        // short-circuit the search for empty paths
        return null;
      }

      for (i = 0; i < data.byteLength;) {
        size = toUnsigned$1(data[i] << 24 | data[i + 1] << 16 | data[i + 2] << 8 | data[i + 3]);
        type = parseType_1(data.subarray(i + 4, i + 8));
        end = size > 1 ? i + size : data.byteLength;

        if (type === path[0]) {
          if (path.length === 1) {
            // this is the end of the path and we've found the box we were
            // looking for
            results.push(data.subarray(i + 8, end));
          } else {
            // recursively search for the next box along the path
            subresults = findBox(data.subarray(i + 8, end), path.slice(1));

            if (subresults.length) {
              results = results.concat(subresults);
            }
          }
        }

        i = end;
      } // we've finished searching all of data


      return results;
    };

    var findBox_1 = findBox;
    var toUnsigned$2 = bin.toUnsigned;

    var tfdt = function tfdt(data) {
      var result = {
        version: data[0],
        flags: new Uint8Array(data.subarray(1, 4)),
        baseMediaDecodeTime: toUnsigned$2(data[4] << 24 | data[5] << 16 | data[6] << 8 | data[7])
      };

      if (result.version === 1) {
        result.baseMediaDecodeTime *= Math.pow(2, 32);
        result.baseMediaDecodeTime += toUnsigned$2(data[8] << 24 | data[9] << 16 | data[10] << 8 | data[11]);
      }

      return result;
    };

    var parseTfdt = tfdt;

    var parseSampleFlags = function parseSampleFlags(flags) {
      return {
        isLeading: (flags[0] & 0x0c) >>> 2,
        dependsOn: flags[0] & 0x03,
        isDependedOn: (flags[1] & 0xc0) >>> 6,
        hasRedundancy: (flags[1] & 0x30) >>> 4,
        paddingValue: (flags[1] & 0x0e) >>> 1,
        isNonSyncSample: flags[1] & 0x01,
        degradationPriority: flags[2] << 8 | flags[3]
      };
    };

    var parseSampleFlags_1 = parseSampleFlags;

    var trun$1 = function trun$1(data) {
      var result = {
        version: data[0],
        flags: new Uint8Array(data.subarray(1, 4)),
        samples: []
      },
          view = new DataView(data.buffer, data.byteOffset, data.byteLength),
          // Flag interpretation
      dataOffsetPresent = result.flags[2] & 0x01,
          // compare with 2nd byte of 0x1
      firstSampleFlagsPresent = result.flags[2] & 0x04,
          // compare with 2nd byte of 0x4
      sampleDurationPresent = result.flags[1] & 0x01,
          // compare with 2nd byte of 0x100
      sampleSizePresent = result.flags[1] & 0x02,
          // compare with 2nd byte of 0x200
      sampleFlagsPresent = result.flags[1] & 0x04,
          // compare with 2nd byte of 0x400
      sampleCompositionTimeOffsetPresent = result.flags[1] & 0x08,
          // compare with 2nd byte of 0x800
      sampleCount = view.getUint32(4),
          offset = 8,
          sample;

      if (dataOffsetPresent) {
        // 32 bit signed integer
        result.dataOffset = view.getInt32(offset);
        offset += 4;
      } // Overrides the flags for the first sample only. The order of
      // optional values will be: duration, size, compositionTimeOffset


      if (firstSampleFlagsPresent && sampleCount) {
        sample = {
          flags: parseSampleFlags_1(data.subarray(offset, offset + 4))
        };
        offset += 4;

        if (sampleDurationPresent) {
          sample.duration = view.getUint32(offset);
          offset += 4;
        }

        if (sampleSizePresent) {
          sample.size = view.getUint32(offset);
          offset += 4;
        }

        if (sampleCompositionTimeOffsetPresent) {
          if (result.version === 1) {
            sample.compositionTimeOffset = view.getInt32(offset);
          } else {
            sample.compositionTimeOffset = view.getUint32(offset);
          }

          offset += 4;
        }

        result.samples.push(sample);
        sampleCount--;
      }

      while (sampleCount--) {
        sample = {};

        if (sampleDurationPresent) {
          sample.duration = view.getUint32(offset);
          offset += 4;
        }

        if (sampleSizePresent) {
          sample.size = view.getUint32(offset);
          offset += 4;
        }

        if (sampleFlagsPresent) {
          sample.flags = parseSampleFlags_1(data.subarray(offset, offset + 4));
          offset += 4;
        }

        if (sampleCompositionTimeOffsetPresent) {
          if (result.version === 1) {
            sample.compositionTimeOffset = view.getInt32(offset);
          } else {
            sample.compositionTimeOffset = view.getUint32(offset);
          }

          offset += 4;
        }

        result.samples.push(sample);
      }

      return result;
    };

    var parseTrun = trun$1;

    var tfhd = function tfhd(data) {
      var view = new DataView(data.buffer, data.byteOffset, data.byteLength),
          result = {
        version: data[0],
        flags: new Uint8Array(data.subarray(1, 4)),
        trackId: view.getUint32(4)
      },
          baseDataOffsetPresent = result.flags[2] & 0x01,
          sampleDescriptionIndexPresent = result.flags[2] & 0x02,
          defaultSampleDurationPresent = result.flags[2] & 0x08,
          defaultSampleSizePresent = result.flags[2] & 0x10,
          defaultSampleFlagsPresent = result.flags[2] & 0x20,
          durationIsEmpty = result.flags[0] & 0x010000,
          defaultBaseIsMoof = result.flags[0] & 0x020000,
          i;
      i = 8;

      if (baseDataOffsetPresent) {
        i += 4; // truncate top 4 bytes
        // FIXME: should we read the full 64 bits?

        result.baseDataOffset = view.getUint32(12);
        i += 4;
      }

      if (sampleDescriptionIndexPresent) {
        result.sampleDescriptionIndex = view.getUint32(i);
        i += 4;
      }

      if (defaultSampleDurationPresent) {
        result.defaultSampleDuration = view.getUint32(i);
        i += 4;
      }

      if (defaultSampleSizePresent) {
        result.defaultSampleSize = view.getUint32(i);
        i += 4;
      }

      if (defaultSampleFlagsPresent) {
        result.defaultSampleFlags = view.getUint32(i);
      }

      if (durationIsEmpty) {
        result.durationIsEmpty = true;
      }

      if (!baseDataOffsetPresent && defaultBaseIsMoof) {
        result.baseDataOffsetIsMoof = true;
      }

      return result;
    };

    var parseTfhd = tfhd;
    var discardEmulationPreventionBytes$1 = captionPacketParser.discardEmulationPreventionBytes;
    var CaptionStream$1 = captionStream.CaptionStream;
    /**
      * Maps an offset in the mdat to a sample based on the the size of the samples.
      * Assumes that `parseSamples` has been called first.
      *
      * @param {Number} offset - The offset into the mdat
      * @param {Object[]} samples - An array of samples, parsed using `parseSamples`
      * @return {?Object} The matching sample, or null if no match was found.
      *
      * @see ISO-BMFF-12/2015, Section 8.8.8
     **/

    var mapToSample = function mapToSample(offset, samples) {
      var approximateOffset = offset;

      for (var i = 0; i < samples.length; i++) {
        var sample = samples[i];

        if (approximateOffset < sample.size) {
          return sample;
        }

        approximateOffset -= sample.size;
      }

      return null;
    };
    /**
      * Finds SEI nal units contained in a Media Data Box.
      * Assumes that `parseSamples` has been called first.
      *
      * @param {Uint8Array} avcStream - The bytes of the mdat
      * @param {Object[]} samples - The samples parsed out by `parseSamples`
      * @param {Number} trackId - The trackId of this video track
      * @return {Object[]} seiNals - the parsed SEI NALUs found.
      *   The contents of the seiNal should match what is expected by
      *   CaptionStream.push (nalUnitType, size, data, escapedRBSP, pts, dts)
      *
      * @see ISO-BMFF-12/2015, Section 8.1.1
      * @see Rec. ITU-T H.264, 7.3.2.3.1
     **/


    var findSeiNals = function findSeiNals(avcStream, samples, trackId) {
      var avcView = new DataView(avcStream.buffer, avcStream.byteOffset, avcStream.byteLength),
          result = [],
          seiNal,
          i,
          length,
          lastMatchedSample;

      for (i = 0; i + 4 < avcStream.length; i += length) {
        length = avcView.getUint32(i);
        i += 4; // Bail if this doesn't appear to be an H264 stream

        if (length <= 0) {
          continue;
        }

        switch (avcStream[i] & 0x1F) {
          case 0x06:
            var data = avcStream.subarray(i + 1, i + 1 + length);
            var matchingSample = mapToSample(i, samples);
            seiNal = {
              nalUnitType: 'sei_rbsp',
              size: length,
              data: data,
              escapedRBSP: discardEmulationPreventionBytes$1(data),
              trackId: trackId
            };

            if (matchingSample) {
              seiNal.pts = matchingSample.pts;
              seiNal.dts = matchingSample.dts;
              lastMatchedSample = matchingSample;
            } else if (lastMatchedSample) {
              // If a matching sample cannot be found, use the last
              // sample's values as they should be as close as possible
              seiNal.pts = lastMatchedSample.pts;
              seiNal.dts = lastMatchedSample.dts;
            } else {
              // eslint-disable-next-line no-console
              console.log("We've encountered a nal unit without data. See mux.js#233.");
              break;
            }

            result.push(seiNal);
            break;
        }
      }

      return result;
    };
    /**
      * Parses sample information out of Track Run Boxes and calculates
      * the absolute presentation and decode timestamps of each sample.
      *
      * @param {Array<Uint8Array>} truns - The Trun Run boxes to be parsed
      * @param {Number} baseMediaDecodeTime - base media decode time from tfdt
          @see ISO-BMFF-12/2015, Section 8.8.12
      * @param {Object} tfhd - The parsed Track Fragment Header
      *   @see inspect.parseTfhd
      * @return {Object[]} the parsed samples
      *
      * @see ISO-BMFF-12/2015, Section 8.8.8
     **/


    var parseSamples = function parseSamples(truns, baseMediaDecodeTime, tfhd) {
      var currentDts = baseMediaDecodeTime;
      var defaultSampleDuration = tfhd.defaultSampleDuration || 0;
      var defaultSampleSize = tfhd.defaultSampleSize || 0;
      var trackId = tfhd.trackId;
      var allSamples = [];
      truns.forEach(function (trun) {
        // Note: We currently do not parse the sample table as well
        // as the trun. It's possible some sources will require this.
        // moov > trak > mdia > minf > stbl
        var trackRun = parseTrun(trun);
        var samples = trackRun.samples;
        samples.forEach(function (sample) {
          if (sample.duration === undefined) {
            sample.duration = defaultSampleDuration;
          }

          if (sample.size === undefined) {
            sample.size = defaultSampleSize;
          }

          sample.trackId = trackId;
          sample.dts = currentDts;

          if (sample.compositionTimeOffset === undefined) {
            sample.compositionTimeOffset = 0;
          }

          sample.pts = currentDts + sample.compositionTimeOffset;
          currentDts += sample.duration;
        });
        allSamples = allSamples.concat(samples);
      });
      return allSamples;
    };
    /**
      * Parses out caption nals from an FMP4 segment's video tracks.
      *
      * @param {Uint8Array} segment - The bytes of a single segment
      * @param {Number} videoTrackId - The trackId of a video track in the segment
      * @return {Object.<Number, Object[]>} A mapping of video trackId to
      *   a list of seiNals found in that track
     **/


    var parseCaptionNals = function parseCaptionNals(segment, videoTrackId) {
      // To get the samples
      var trafs = findBox_1(segment, ['moof', 'traf']); // To get SEI NAL units

      var mdats = findBox_1(segment, ['mdat']);
      var captionNals = {};
      var mdatTrafPairs = []; // Pair up each traf with a mdat as moofs and mdats are in pairs

      mdats.forEach(function (mdat, index) {
        var matchingTraf = trafs[index];
        mdatTrafPairs.push({
          mdat: mdat,
          traf: matchingTraf
        });
      });
      mdatTrafPairs.forEach(function (pair) {
        var mdat = pair.mdat;
        var traf = pair.traf;
        var tfhd = findBox_1(traf, ['tfhd']); // Exactly 1 tfhd per traf

        var headerInfo = parseTfhd(tfhd[0]);
        var trackId = headerInfo.trackId;
        var tfdt = findBox_1(traf, ['tfdt']); // Either 0 or 1 tfdt per traf

        var baseMediaDecodeTime = tfdt.length > 0 ? parseTfdt(tfdt[0]).baseMediaDecodeTime : 0;
        var truns = findBox_1(traf, ['trun']);
        var samples;
        var seiNals; // Only parse video data for the chosen video track

        if (videoTrackId === trackId && truns.length > 0) {
          samples = parseSamples(truns, baseMediaDecodeTime, headerInfo);
          seiNals = findSeiNals(mdat, samples, trackId);

          if (!captionNals[trackId]) {
            captionNals[trackId] = [];
          }

          captionNals[trackId] = captionNals[trackId].concat(seiNals);
        }
      });
      return captionNals;
    };
    /**
      * Parses out inband captions from an MP4 container and returns
      * caption objects that can be used by WebVTT and the TextTrack API.
      * @see https://developer.mozilla.org/en-US/docs/Web/API/VTTCue
      * @see https://developer.mozilla.org/en-US/docs/Web/API/TextTrack
      * Assumes that `probe.getVideoTrackIds` and `probe.timescale` have been called first
      *
      * @param {Uint8Array} segment - The fmp4 segment containing embedded captions
      * @param {Number} trackId - The id of the video track to parse
      * @param {Number} timescale - The timescale for the video track from the init segment
      *
      * @return {?Object[]} parsedCaptions - A list of captions or null if no video tracks
      * @return {Number} parsedCaptions[].startTime - The time to show the caption in seconds
      * @return {Number} parsedCaptions[].endTime - The time to stop showing the caption in seconds
      * @return {String} parsedCaptions[].text - The visible content of the caption
     **/


    var parseEmbeddedCaptions = function parseEmbeddedCaptions(segment, trackId, timescale) {
      var seiNals; // the ISO-BMFF spec says that trackId can't be zero, but there's some broken content out there

      if (trackId === null) {
        return null;
      }

      seiNals = parseCaptionNals(segment, trackId);
      return {
        seiNals: seiNals[trackId],
        timescale: timescale
      };
    };
    /**
      * Converts SEI NALUs into captions that can be used by video.js
     **/


    var CaptionParser = function CaptionParser() {
      var isInitialized = false;
      var captionStream; // Stores segments seen before trackId and timescale are set

      var segmentCache; // Stores video track ID of the track being parsed

      var trackId; // Stores the timescale of the track being parsed

      var timescale; // Stores captions parsed so far

      var parsedCaptions; // Stores whether we are receiving partial data or not

      var parsingPartial;
      /**
        * A method to indicate whether a CaptionParser has been initalized
        * @returns {Boolean}
       **/

      this.isInitialized = function () {
        return isInitialized;
      };
      /**
        * Initializes the underlying CaptionStream, SEI NAL parsing
        * and management, and caption collection
       **/


      this.init = function (options) {
        captionStream = new CaptionStream$1();
        isInitialized = true;
        parsingPartial = options ? options.isPartial : false; // Collect dispatched captions

        captionStream.on('data', function (event) {
          // Convert to seconds in the source's timescale
          event.startTime = event.startPts / timescale;
          event.endTime = event.endPts / timescale;
          parsedCaptions.captions.push(event);
          parsedCaptions.captionStreams[event.stream] = true;
        });
      };
      /**
        * Determines if a new video track will be selected
        * or if the timescale changed
        * @return {Boolean}
       **/


      this.isNewInit = function (videoTrackIds, timescales) {
        if (videoTrackIds && videoTrackIds.length === 0 || timescales && typeof timescales === 'object' && Object.keys(timescales).length === 0) {
          return false;
        }

        return trackId !== videoTrackIds[0] || timescale !== timescales[trackId];
      };
      /**
        * Parses out SEI captions and interacts with underlying
        * CaptionStream to return dispatched captions
        *
        * @param {Uint8Array} segment - The fmp4 segment containing embedded captions
        * @param {Number[]} videoTrackIds - A list of video tracks found in the init segment
        * @param {Object.<Number, Number>} timescales - The timescales found in the init segment
        * @see parseEmbeddedCaptions
        * @see m2ts/caption-stream.js
       **/


      this.parse = function (segment, videoTrackIds, timescales) {
        var parsedData;

        if (!this.isInitialized()) {
          return null; // This is not likely to be a video segment
        } else if (!videoTrackIds || !timescales) {
          return null;
        } else if (this.isNewInit(videoTrackIds, timescales)) {
          // Use the first video track only as there is no
          // mechanism to switch to other video tracks
          trackId = videoTrackIds[0];
          timescale = timescales[trackId]; // If an init segment has not been seen yet, hold onto segment
          // data until we have one.
          // the ISO-BMFF spec says that trackId can't be zero, but there's some broken content out there
        } else if (trackId === null || !timescale) {
          segmentCache.push(segment);
          return null;
        } // Now that a timescale and trackId is set, parse cached segments


        while (segmentCache.length > 0) {
          var cachedSegment = segmentCache.shift();
          this.parse(cachedSegment, videoTrackIds, timescales);
        }

        parsedData = parseEmbeddedCaptions(segment, trackId, timescale);

        if (parsedData === null || !parsedData.seiNals) {
          return null;
        }

        this.pushNals(parsedData.seiNals); // Force the parsed captions to be dispatched

        this.flushStream();
        return parsedCaptions;
      };
      /**
        * Pushes SEI NALUs onto CaptionStream
        * @param {Object[]} nals - A list of SEI nals parsed using `parseCaptionNals`
        * Assumes that `parseCaptionNals` has been called first
        * @see m2ts/caption-stream.js
        **/


      this.pushNals = function (nals) {
        if (!this.isInitialized() || !nals || nals.length === 0) {
          return null;
        }

        nals.forEach(function (nal) {
          captionStream.push(nal);
        });
      };
      /**
        * Flushes underlying CaptionStream to dispatch processed, displayable captions
        * @see m2ts/caption-stream.js
       **/


      this.flushStream = function () {
        if (!this.isInitialized()) {
          return null;
        }

        if (!parsingPartial) {
          captionStream.flush();
        } else {
          captionStream.partialFlush();
        }
      };
      /**
        * Reset caption buckets for new data
       **/


      this.clearParsedCaptions = function () {
        parsedCaptions.captions = [];
        parsedCaptions.captionStreams = {};
      };
      /**
        * Resets underlying CaptionStream
        * @see m2ts/caption-stream.js
       **/


      this.resetCaptionStream = function () {
        if (!this.isInitialized()) {
          return null;
        }

        captionStream.reset();
      };
      /**
        * Convenience method to clear all captions flushed from the
        * CaptionStream and still being parsed
        * @see m2ts/caption-stream.js
       **/


      this.clearAllCaptions = function () {
        this.clearParsedCaptions();
        this.resetCaptionStream();
      };
      /**
        * Reset caption parser
       **/


      this.reset = function () {
        segmentCache = [];
        trackId = null;
        timescale = null;

        if (!parsedCaptions) {
          parsedCaptions = {
            captions: [],
            // CC1, CC2, CC3, CC4
            captionStreams: {}
          };
        } else {
          this.clearParsedCaptions();
        }

        this.resetCaptionStream();
      };

      this.reset();
    };

    var captionParser = CaptionParser;
    /* global self */

    var typeFromStreamString = function typeFromStreamString(streamString) {
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


    var wireFullTransmuxerEvents = function wireFullTransmuxerEvents(self, transmuxer) {
      transmuxer.on('data', function (segment) {
        // transfer ownership of the underlying ArrayBuffer
        // instead of doing a copy to save memory
        // ArrayBuffers are transferable but generic TypedArrays are not
        // @link https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API/Using_web_workers#Passing_data_by_transferring_ownership_(transferable_objects)
        var initArray = segment.initSegment;
        segment.initSegment = {
          data: initArray.buffer,
          byteOffset: initArray.byteOffset,
          byteLength: initArray.byteLength
        };
        var typedArray = segment.data;
        segment.data = typedArray.buffer;
        self.postMessage({
          action: 'data',
          segment: segment,
          byteOffset: typedArray.byteOffset,
          byteLength: typedArray.byteLength
        }, [segment.data]);
      });
      transmuxer.on('done', function (data) {
        self.postMessage({
          action: 'done'
        });
      });
      transmuxer.on('gopInfo', function (gopInfo) {
        self.postMessage({
          action: 'gopInfo',
          gopInfo: gopInfo
        });
      });
      transmuxer.on('videoSegmentTimingInfo', function (timingInfo) {
        var videoSegmentTimingInfo = {
          start: {
            decode: clock_4(timingInfo.start.dts),
            presentation: clock_4(timingInfo.start.pts)
          },
          end: {
            decode: clock_4(timingInfo.end.dts),
            presentation: clock_4(timingInfo.end.pts)
          },
          baseMediaDecodeTime: clock_4(timingInfo.baseMediaDecodeTime)
        };

        if (timingInfo.prependedContentDuration) {
          videoSegmentTimingInfo.prependedContentDuration = clock_4(timingInfo.prependedContentDuration);
        }

        self.postMessage({
          action: 'videoSegmentTimingInfo',
          videoSegmentTimingInfo: videoSegmentTimingInfo
        });
      });
      transmuxer.on('id3Frame', function (id3Frame) {
        self.postMessage({
          action: 'id3Frame',
          id3Frame: id3Frame
        });
      });
      transmuxer.on('caption', function (caption) {
        self.postMessage({
          action: 'caption',
          caption: caption
        });
      });
      transmuxer.on('trackinfo', function (trackInfo) {
        self.postMessage({
          action: 'trackinfo',
          trackInfo: trackInfo
        });
      });
      transmuxer.on('audioTimingInfo', function (audioTimingInfo) {
        // convert to video TS since we prioritize video time over audio
        self.postMessage({
          action: 'audioTimingInfo',
          audioTimingInfo: {
            start: clock_4(audioTimingInfo.start),
            end: clock_4(audioTimingInfo.end)
          }
        });
      });
      transmuxer.on('videoTimingInfo', function (videoTimingInfo) {
        self.postMessage({
          action: 'videoTimingInfo',
          videoTimingInfo: {
            start: clock_4(videoTimingInfo.start),
            end: clock_4(videoTimingInfo.end)
          }
        });
      });
    };

    var wirePartialTransmuxerEvents = function wirePartialTransmuxerEvents(self, transmuxer) {
      transmuxer.on('data', function (event) {
        // transfer ownership of the underlying ArrayBuffer
        // instead of doing a copy to save memory
        // ArrayBuffers are transferable but generic TypedArrays are not
        // @link https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API/Using_web_workers#Passing_data_by_transferring_ownership_(transferable_objects)
        var initSegment = {
          data: event.data.track.initSegment.buffer,
          byteOffset: event.data.track.initSegment.byteOffset,
          byteLength: event.data.track.initSegment.byteLength
        };
        var boxes = {
          data: event.data.boxes.buffer,
          byteOffset: event.data.boxes.byteOffset,
          byteLength: event.data.boxes.byteLength
        };
        var segment = {
          boxes: boxes,
          initSegment: initSegment,
          type: event.type,
          sequence: event.data.sequence
        };

        if (typeof event.data.videoFrameDts !== 'undefined') {
          segment.videoFrameDtsTime = clock_4(event.data.videoFrameDts);
        }

        if (typeof event.data.videoFramePts !== 'undefined') {
          segment.videoFramePtsTime = clock_4(event.data.videoFramePts);
        }

        self.postMessage({
          action: 'data',
          segment: segment
        }, [segment.boxes.data, segment.initSegment.data]);
      });
      transmuxer.on('id3Frame', function (id3Frame) {
        self.postMessage({
          action: 'id3Frame',
          id3Frame: id3Frame
        });
      });
      transmuxer.on('caption', function (caption) {
        self.postMessage({
          action: 'caption',
          caption: caption
        });
      });
      transmuxer.on('done', function (data) {
        self.postMessage({
          action: 'done',
          type: typeFromStreamString(data)
        });
      });
      transmuxer.on('partialdone', function (data) {
        self.postMessage({
          action: 'partialdone',
          type: typeFromStreamString(data)
        });
      });
      transmuxer.on('endedsegment', function (data) {
        self.postMessage({
          action: 'endedSegment',
          type: typeFromStreamString(data)
        });
      });
      transmuxer.on('trackinfo', function (trackInfo) {
        self.postMessage({
          action: 'trackinfo',
          trackInfo: trackInfo
        });
      });
      transmuxer.on('audioTimingInfo', function (audioTimingInfo) {
        // This can happen if flush is called when no
        // audio has been processed. This should be an
        // unusual case, but if it does occur should not
        // result in valid data being returned
        if (audioTimingInfo.start === null) {
          self.postMessage({
            action: 'audioTimingInfo',
            audioTimingInfo: audioTimingInfo
          });
          return;
        } // convert to video TS since we prioritize video time over audio


        var timingInfoInSeconds = {
          start: clock_4(audioTimingInfo.start)
        };

        if (audioTimingInfo.end) {
          timingInfoInSeconds.end = clock_4(audioTimingInfo.end);
        }

        self.postMessage({
          action: 'audioTimingInfo',
          audioTimingInfo: timingInfoInSeconds
        });
      });
      transmuxer.on('videoTimingInfo', function (videoTimingInfo) {
        var timingInfoInSeconds = {
          start: clock_4(videoTimingInfo.start)
        };

        if (videoTimingInfo.end) {
          timingInfoInSeconds.end = clock_4(videoTimingInfo.end);
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


    var MessageHandlers = /*#__PURE__*/function () {
      function MessageHandlers(self, options) {
        this.options = options || {};
        this.self = self;
        this.init();
      }
      /**
       * initialize our web worker and wire all the events.
       */


      var _proto = MessageHandlers.prototype;

      _proto.init = function init() {
        if (this.transmuxer) {
          this.transmuxer.dispose();
        }

        this.transmuxer = this.options.handlePartialData ? new transmuxer$1(this.options) : new transmuxer_1(this.options);

        if (this.options.handlePartialData) {
          wirePartialTransmuxerEvents(this.self, this.transmuxer);
        } else {
          wireFullTransmuxerEvents(this.self, this.transmuxer);
        }
      };

      _proto.pushMp4Captions = function pushMp4Captions(data) {
        if (!this.captionParser) {
          this.captionParser = new captionParser();
          this.captionParser.init();
        }

        var segment = new Uint8Array(data.data, data.byteOffset, data.byteLength);
        var parsed = this.captionParser.parse(segment, data.trackIds, data.timescales);
        this.self.postMessage({
          action: 'mp4Captions',
          captions: parsed && parsed.captions || [],
          data: segment.buffer
        }, [segment.buffer]);
      };

      _proto.clearAllMp4Captions = function clearAllMp4Captions() {
        if (this.captionParser) {
          this.captionParser.clearAllCaptions();
        }
      };

      _proto.clearParsedMp4Captions = function clearParsedMp4Captions() {
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
      ;

      _proto.push = function push(data) {
        // Cast array buffer to correct type for transmuxer
        var segment = new Uint8Array(data.data, data.byteOffset, data.byteLength);
        this.transmuxer.push(segment);
      }
      /**
       * Recreate the transmuxer so that the next segment added via `push`
       * start with a fresh transmuxer.
       */
      ;

      _proto.reset = function reset() {
        this.transmuxer.reset();
      }
      /**
       * Set the value that will be used as the `baseMediaDecodeTime` time for the
       * next segment pushed in. Subsequent segments will have their `baseMediaDecodeTime`
       * set relative to the first based on the PTS values.
       *
       * @param {Object} data used to set the timestamp offset in the muxer
       */
      ;

      _proto.setTimestampOffset = function setTimestampOffset(data) {
        var timestampOffset = data.timestampOffset || 0;
        this.transmuxer.setBaseMediaDecodeTime(Math.round(clock_2(timestampOffset)));
      };

      _proto.setAudioAppendStart = function setAudioAppendStart(data) {
        this.transmuxer.setAudioAppendStart(Math.ceil(clock_2(data.appendStart)));
      };

      _proto.setRemux = function setRemux(data) {
        this.transmuxer.setRemux(data.remux);
      }
      /**
       * Forces the pipeline to finish processing the last segment and emit it's
       * results.
       *
       * @param {Object} data event data, not really used
       */
      ;

      _proto.flush = function flush(data) {
        this.transmuxer.flush(); // transmuxed done action is fired after both audio/video pipelines are flushed

        self.postMessage({
          action: 'done',
          type: 'transmuxed'
        });
      };

      _proto.partialFlush = function partialFlush(data) {
        this.transmuxer.partialFlush(); // transmuxed partialdone action is fired after both audio/video pipelines are flushed

        self.postMessage({
          action: 'partialdone',
          type: 'transmuxed'
        });
      };

      _proto.endTimeline = function endTimeline() {
        this.transmuxer.endTimeline(); // transmuxed endedtimeline action is fired after both audio/video pipelines end their
        // timelines

        self.postMessage({
          action: 'endedtimeline',
          type: 'transmuxed'
        });
      };

      _proto.alignGopsWith = function alignGopsWith(data) {
        this.transmuxer.alignGopsWith(data.gopsToAlignWith.slice());
      };

      return MessageHandlers;
    }();
    /**
     * Our web worker interface so that things can talk to mux.js
     * that will be running in a web worker. the scope is passed to this by
     * webworkify.
     *
     * @param {Object} self the scope for the web worker
     */


    var TransmuxerWorker = function TransmuxerWorker(self) {
      self.onmessage = function (event) {
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

    var transmuxerWorker = new TransmuxerWorker(self);
    return transmuxerWorker;
  }();
});

/**
 * @file - codecs.js - Handles tasks regarding codec strings such as translating them to
 * codec strings, or translating codec strings into objects that can be examined.
 */
/**
 * Returns a set of codec strings parsed from the playlist or the default
 * codec strings if no codecs were specified in the playlist
 *
 * @param {Playlist} media the current media playlist
 * @return {Object} an object with the video and audio codecs
 */

var getCodecs = function getCodecs(media) {
  // if the codecs were explicitly specified, use them instead of the
  // defaults
  var mediaAttributes = media.attributes || {};

  if (mediaAttributes.CODECS) {
    return codecs.parseCodecs(mediaAttributes.CODECS);
  }
};

var isMaat = function isMaat(master, media) {
  var mediaAttributes = media.attributes || {};
  return master && master.mediaGroups && master.mediaGroups.AUDIO && mediaAttributes.AUDIO && master.mediaGroups.AUDIO[mediaAttributes.AUDIO];
};
var isMuxed = function isMuxed(master, media) {
  if (!isMaat(master, media)) {
    return true;
  }

  var mediaAttributes = media.attributes || {};
  var audioGroup = master.mediaGroups.AUDIO[mediaAttributes.AUDIO];

  for (var groupId in audioGroup) {
    // If an audio group has a URI (the case for HLS, as HLS will use external playlists),
    // or there are listed playlists (the case for DASH, as the manifest will have already
    // provided all of the details necessary to generate the audio playlist, as opposed to
    // HLS' externally requested playlists), then the content is demuxed.
    if (!audioGroup[groupId].uri && !audioGroup[groupId].playlists) {
      return true;
    }
  }

  return false;
};
/**
 * Calculates the codec strings for a working configuration of
 * SourceBuffers to play variant streams in a master playlist. If
 * there is no possible working configuration, an empty object will be
 * returned.
 *
 * @param master {Object} the m3u8 object for the master playlist
 * @param media {Object} the m3u8 object for the variant playlist
 * @return {Object} the codec strings.
 *
 * @private
 */

var codecsForPlaylist = function codecsForPlaylist(master, media) {
  var mediaAttributes = media.attributes || {};
  var codecInfo = getCodecs(media) || {}; // HLS with multiple-audio tracks must always get an audio codec.
  // Put another way, there is no way to have a video-only multiple-audio HLS!

  if (isMaat(master, media) && !codecInfo.audio) {
    if (!isMuxed(master, media)) {
      // It is possible for codecs to be specified on the audio media group playlist but
      // not on the rendition playlist. This is mostly the case for DASH, where audio and
      // video are always separate (and separately specified).
      var defaultCodecs = codecs.codecsFromDefault(master, mediaAttributes.AUDIO);

      if (defaultCodecs) {
        codecInfo.audio = defaultCodecs.audio;
      }
    }
  }

  var codecs$1 = {};

  if (codecInfo.video) {
    codecs$1.video = codecs.translateLegacyCodec("" + codecInfo.video.type + codecInfo.video.details);
  }

  if (codecInfo.audio) {
    codecs$1.audio = codecs.translateLegacyCodec("" + codecInfo.audio.type + codecInfo.audio.details);
  }

  return codecs$1;
};

var logger = function logger(source) {
  if (videojs$1.log.debug) {
    return videojs$1.log.debug.bind(videojs$1, 'VHS:', source + " >");
  }

  return function () {};
};

var logFn = logger('PlaylistSelector');

var representationToString = function representationToString(representation) {
  if (!representation || !representation.playlist) {
    return;
  }

  var playlist = representation.playlist;
  return JSON.stringify({
    id: playlist.id,
    bandwidth: representation.bandwidth,
    width: representation.width,
    height: representation.height,
    codecs: playlist.attributes && playlist.attributes.CODECS || ''
  });
}; // Utilities

/**
 * Returns the CSS value for the specified property on an element
 * using `getComputedStyle`. Firefox has a long-standing issue where
 * getComputedStyle() may return null when running in an iframe with
 * `display: none`.
 *
 * @see https://bugzilla.mozilla.org/show_bug.cgi?id=548397
 * @param {HTMLElement} el the htmlelement to work on
 * @param {string} the proprety to get the style for
 */


var safeGetComputedStyle = function safeGetComputedStyle(el, property) {
  if (!el) {
    return '';
  }

  var result = window_1.getComputedStyle(el);

  if (!result) {
    return '';
  }

  return result[property];
};
/**
 * Resuable stable sort function
 *
 * @param {Playlists} array
 * @param {Function} sortFn Different comparators
 * @function stableSort
 */


var stableSort = function stableSort(array, sortFn) {
  var newArray = array.slice();
  array.sort(function (left, right) {
    var cmp = sortFn(left, right);

    if (cmp === 0) {
      return newArray.indexOf(left) - newArray.indexOf(right);
    }

    return cmp;
  });
};
/**
 * A comparator function to sort two playlist object by bandwidth.
 *
 * @param {Object} left a media playlist object
 * @param {Object} right a media playlist object
 * @return {number} Greater than zero if the bandwidth attribute of
 * left is greater than the corresponding attribute of right. Less
 * than zero if the bandwidth of right is greater than left and
 * exactly zero if the two are equal.
 */


var comparePlaylistBandwidth = function comparePlaylistBandwidth(left, right) {
  var leftBandwidth;
  var rightBandwidth;

  if (left.attributes.BANDWIDTH) {
    leftBandwidth = left.attributes.BANDWIDTH;
  }

  leftBandwidth = leftBandwidth || window_1.Number.MAX_VALUE;

  if (right.attributes.BANDWIDTH) {
    rightBandwidth = right.attributes.BANDWIDTH;
  }

  rightBandwidth = rightBandwidth || window_1.Number.MAX_VALUE;
  return leftBandwidth - rightBandwidth;
};
/**
 * A comparator function to sort two playlist object by resolution (width).
 *
 * @param {Object} left a media playlist object
 * @param {Object} right a media playlist object
 * @return {number} Greater than zero if the resolution.width attribute of
 * left is greater than the corresponding attribute of right. Less
 * than zero if the resolution.width of right is greater than left and
 * exactly zero if the two are equal.
 */

var comparePlaylistResolution = function comparePlaylistResolution(left, right) {
  var leftWidth;
  var rightWidth;

  if (left.attributes.RESOLUTION && left.attributes.RESOLUTION.width) {
    leftWidth = left.attributes.RESOLUTION.width;
  }

  leftWidth = leftWidth || window_1.Number.MAX_VALUE;

  if (right.attributes.RESOLUTION && right.attributes.RESOLUTION.width) {
    rightWidth = right.attributes.RESOLUTION.width;
  }

  rightWidth = rightWidth || window_1.Number.MAX_VALUE; // NOTE - Fallback to bandwidth sort as appropriate in cases where multiple renditions
  // have the same media dimensions/ resolution

  if (leftWidth === rightWidth && left.attributes.BANDWIDTH && right.attributes.BANDWIDTH) {
    return left.attributes.BANDWIDTH - right.attributes.BANDWIDTH;
  }

  return leftWidth - rightWidth;
};
/**
 * Chooses the appropriate media playlist based on bandwidth and player size
 *
 * @param {Object} master
 *        Object representation of the master manifest
 * @param {number} playerBandwidth
 *        Current calculated bandwidth of the player
 * @param {number} playerWidth
 *        Current width of the player element (should account for the device pixel ratio)
 * @param {number} playerHeight
 *        Current height of the player element (should account for the device pixel ratio)
 * @param {boolean} limitRenditionByPlayerDimensions
 *        True if the player width and height should be used during the selection, false otherwise
 * @return {Playlist} the highest bitrate playlist less than the
 * currently detected bandwidth, accounting for some amount of
 * bandwidth variance
 */

var simpleSelector = function simpleSelector(master, playerBandwidth, playerWidth, playerHeight, limitRenditionByPlayerDimensions) {
  var options = {
    bandwidth: playerBandwidth,
    width: playerWidth,
    height: playerHeight,
    limitRenditionByPlayerDimensions: limitRenditionByPlayerDimensions
  }; // convert the playlists to an intermediary representation to make comparisons easier

  var sortedPlaylistReps = master.playlists.map(function (playlist) {
    var bandwidth;
    var width = playlist.attributes.RESOLUTION && playlist.attributes.RESOLUTION.width;
    var height = playlist.attributes.RESOLUTION && playlist.attributes.RESOLUTION.height;
    bandwidth = playlist.attributes.BANDWIDTH;
    bandwidth = bandwidth || window_1.Number.MAX_VALUE;
    return {
      bandwidth: bandwidth,
      width: width,
      height: height,
      playlist: playlist
    };
  });
  stableSort(sortedPlaylistReps, function (left, right) {
    return left.bandwidth - right.bandwidth;
  }); // filter out any playlists that have been excluded due to
  // incompatible configurations

  sortedPlaylistReps = sortedPlaylistReps.filter(function (rep) {
    return !Playlist.isIncompatible(rep.playlist);
  }); // filter out any playlists that have been disabled manually through the representations
  // api or blacklisted temporarily due to playback errors.

  var enabledPlaylistReps = sortedPlaylistReps.filter(function (rep) {
    return Playlist.isEnabled(rep.playlist);
  });

  if (!enabledPlaylistReps.length) {
    // if there are no enabled playlists, then they have all been blacklisted or disabled
    // by the user through the representations api. In this case, ignore blacklisting and
    // fallback to what the user wants by using playlists the user has not disabled.
    enabledPlaylistReps = sortedPlaylistReps.filter(function (rep) {
      return !Playlist.isDisabled(rep.playlist);
    });
  } // filter out any variant that has greater effective bitrate
  // than the current estimated bandwidth


  var bandwidthPlaylistReps = enabledPlaylistReps.filter(function (rep) {
    return rep.bandwidth * Config.BANDWIDTH_VARIANCE < playerBandwidth;
  });
  var highestRemainingBandwidthRep = bandwidthPlaylistReps[bandwidthPlaylistReps.length - 1]; // get all of the renditions with the same (highest) bandwidth
  // and then taking the very first element

  var bandwidthBestRep = bandwidthPlaylistReps.filter(function (rep) {
    return rep.bandwidth === highestRemainingBandwidthRep.bandwidth;
  })[0]; // if we're not going to limit renditions by player size, make an early decision.

  if (limitRenditionByPlayerDimensions === false) {
    var _chosenRep = bandwidthBestRep || enabledPlaylistReps[0] || sortedPlaylistReps[0];

    if (_chosenRep && _chosenRep.playlist) {
      var type = 'sortedPlaylistReps';

      if (bandwidthBestRep) {
        type = 'bandwidthBestRep';
      }

      if (enabledPlaylistReps[0]) {
        type = 'enabledPlaylistReps';
      }

      logFn("choosing " + representationToString(_chosenRep) + " using " + type + " with options", options);
      return _chosenRep.playlist;
    }

    logFn('could not choose a playlist with options', options);
    return null;
  } // filter out playlists without resolution information


  var haveResolution = bandwidthPlaylistReps.filter(function (rep) {
    return rep.width && rep.height;
  }); // sort variants by resolution

  stableSort(haveResolution, function (left, right) {
    return left.width - right.width;
  }); // if we have the exact resolution as the player use it

  var resolutionBestRepList = haveResolution.filter(function (rep) {
    return rep.width === playerWidth && rep.height === playerHeight;
  });
  highestRemainingBandwidthRep = resolutionBestRepList[resolutionBestRepList.length - 1]; // ensure that we pick the highest bandwidth variant that have exact resolution

  var resolutionBestRep = resolutionBestRepList.filter(function (rep) {
    return rep.bandwidth === highestRemainingBandwidthRep.bandwidth;
  })[0];
  var resolutionPlusOneList;
  var resolutionPlusOneSmallest;
  var resolutionPlusOneRep; // find the smallest variant that is larger than the player
  // if there is no match of exact resolution

  if (!resolutionBestRep) {
    resolutionPlusOneList = haveResolution.filter(function (rep) {
      return rep.width > playerWidth || rep.height > playerHeight;
    }); // find all the variants have the same smallest resolution

    resolutionPlusOneSmallest = resolutionPlusOneList.filter(function (rep) {
      return rep.width === resolutionPlusOneList[0].width && rep.height === resolutionPlusOneList[0].height;
    }); // ensure that we also pick the highest bandwidth variant that
    // is just-larger-than the video player

    highestRemainingBandwidthRep = resolutionPlusOneSmallest[resolutionPlusOneSmallest.length - 1];
    resolutionPlusOneRep = resolutionPlusOneSmallest.filter(function (rep) {
      return rep.bandwidth === highestRemainingBandwidthRep.bandwidth;
    })[0];
  } // fallback chain of variants


  var chosenRep = resolutionPlusOneRep || resolutionBestRep || bandwidthBestRep || enabledPlaylistReps[0] || sortedPlaylistReps[0];

  if (chosenRep && chosenRep.playlist) {
    var _type = 'sortedPlaylistReps';

    if (resolutionPlusOneRep) {
      _type = 'resolutionPlusOneRep';
    } else if (resolutionBestRep) {
      _type = 'resolutionBestRep';
    } else if (bandwidthBestRep) {
      _type = 'bandwidthBestRep';
    } else if (enabledPlaylistReps[0]) {
      _type = 'enabledPlaylistReps';
    }

    logFn("choosing " + representationToString(chosenRep) + " using " + _type + " with options", options);
    return chosenRep.playlist;
  }

  logFn('could not choose a playlist with options', options);
  return null;
}; // Playlist Selectors

/**
 * Chooses the appropriate media playlist based on the most recent
 * bandwidth estimate and the player size.
 *
 * Expects to be called within the context of an instance of VhsHandler
 *
 * @return {Playlist} the highest bitrate playlist less than the
 * currently detected bandwidth, accounting for some amount of
 * bandwidth variance
 */

var lastBandwidthSelector = function lastBandwidthSelector() {
  var pixelRatio = this.useDevicePixelRatio ? window_1.devicePixelRatio || 1 : 1;
  return simpleSelector(this.playlists.master, this.systemBandwidth, parseInt(safeGetComputedStyle(this.tech_.el(), 'width'), 10) * pixelRatio, parseInt(safeGetComputedStyle(this.tech_.el(), 'height'), 10) * pixelRatio, this.limitRenditionByPlayerDimensions);
};
/**
 * Chooses the appropriate media playlist based on the potential to rebuffer
 *
 * @param {Object} settings
 *        Object of information required to use this selector
 * @param {Object} settings.master
 *        Object representation of the master manifest
 * @param {number} settings.currentTime
 *        The current time of the player
 * @param {number} settings.bandwidth
 *        Current measured bandwidth
 * @param {number} settings.duration
 *        Duration of the media
 * @param {number} settings.segmentDuration
 *        Segment duration to be used in round trip time calculations
 * @param {number} settings.timeUntilRebuffer
 *        Time left in seconds until the player has to rebuffer
 * @param {number} settings.currentTimeline
 *        The current timeline segments are being loaded from
 * @param {SyncController} settings.syncController
 *        SyncController for determining if we have a sync point for a given playlist
 * @return {Object|null}
 *         {Object} return.playlist
 *         The highest bandwidth playlist with the least amount of rebuffering
 *         {Number} return.rebufferingImpact
 *         The amount of time in seconds switching to this playlist will rebuffer. A
 *         negative value means that switching will cause zero rebuffering.
 */

var minRebufferMaxBandwidthSelector = function minRebufferMaxBandwidthSelector(settings) {
  var master = settings.master,
      currentTime = settings.currentTime,
      bandwidth = settings.bandwidth,
      duration = settings.duration,
      segmentDuration = settings.segmentDuration,
      timeUntilRebuffer = settings.timeUntilRebuffer,
      currentTimeline = settings.currentTimeline,
      syncController = settings.syncController; // filter out any playlists that have been excluded due to
  // incompatible configurations

  var compatiblePlaylists = master.playlists.filter(function (playlist) {
    return !Playlist.isIncompatible(playlist);
  }); // filter out any playlists that have been disabled manually through the representations
  // api or blacklisted temporarily due to playback errors.

  var enabledPlaylists = compatiblePlaylists.filter(Playlist.isEnabled);

  if (!enabledPlaylists.length) {
    // if there are no enabled playlists, then they have all been blacklisted or disabled
    // by the user through the representations api. In this case, ignore blacklisting and
    // fallback to what the user wants by using playlists the user has not disabled.
    enabledPlaylists = compatiblePlaylists.filter(function (playlist) {
      return !Playlist.isDisabled(playlist);
    });
  }

  var bandwidthPlaylists = enabledPlaylists.filter(Playlist.hasAttribute.bind(null, 'BANDWIDTH'));
  var rebufferingEstimates = bandwidthPlaylists.map(function (playlist) {
    var syncPoint = syncController.getSyncPoint(playlist, duration, currentTimeline, currentTime); // If there is no sync point for this playlist, switching to it will require a
    // sync request first. This will double the request time

    var numRequests = syncPoint ? 1 : 2;
    var requestTimeEstimate = Playlist.estimateSegmentRequestTime(segmentDuration, bandwidth, playlist);
    var rebufferingImpact = requestTimeEstimate * numRequests - timeUntilRebuffer;
    return {
      playlist: playlist,
      rebufferingImpact: rebufferingImpact
    };
  });
  var noRebufferingPlaylists = rebufferingEstimates.filter(function (estimate) {
    return estimate.rebufferingImpact <= 0;
  }); // Sort by bandwidth DESC

  stableSort(noRebufferingPlaylists, function (a, b) {
    return comparePlaylistBandwidth(b.playlist, a.playlist);
  });

  if (noRebufferingPlaylists.length) {
    return noRebufferingPlaylists[0];
  }

  stableSort(rebufferingEstimates, function (a, b) {
    return a.rebufferingImpact - b.rebufferingImpact;
  });
  return rebufferingEstimates[0] || null;
};
/**
 * Chooses the appropriate media playlist, which in this case is the lowest bitrate
 * one with video.  If no renditions with video exist, return the lowest audio rendition.
 *
 * Expects to be called within the context of an instance of VhsHandler
 *
 * @return {Object|null}
 *         {Object} return.playlist
 *         The lowest bitrate playlist that contains a video codec.  If no such rendition
 *         exists pick the lowest audio rendition.
 */

var lowestBitrateCompatibleVariantSelector = function lowestBitrateCompatibleVariantSelector() {
  var _this = this;

  // filter out any playlists that have been excluded due to
  // incompatible configurations or playback errors
  var playlists = this.playlists.master.playlists.filter(Playlist.isEnabled); // Sort ascending by bitrate

  stableSort(playlists, function (a, b) {
    return comparePlaylistBandwidth(a, b);
  }); // Parse and assume that playlists with no video codec have no video
  // (this is not necessarily true, although it is generally true).
  //
  // If an entire manifest has no valid videos everything will get filtered
  // out.

  var playlistsWithVideo = playlists.filter(function (playlist) {
    return !!codecsForPlaylist(_this.playlists.master, playlist).video;
  });
  return playlistsWithVideo[0] || null;
};

/**
 * @file text-tracks.js
 */
/**
 * Create captions text tracks on video.js if they do not exist
 *
 * @param {Object} inbandTextTracks a reference to current inbandTextTracks
 * @param {Object} tech the video.js tech
 * @param {Object} captionStream the caption stream to create
 * @private
 */

var createCaptionsTrackIfNotExists = function createCaptionsTrackIfNotExists(inbandTextTracks, tech, captionStream) {
  if (!inbandTextTracks[captionStream]) {
    tech.trigger({
      type: 'usage',
      name: 'vhs-608'
    });
    tech.trigger({
      type: 'usage',
      name: 'hls-608'
    });
    var track = tech.textTracks().getTrackById(captionStream);

    if (track) {
      // Resuse an existing track with a CC# id because this was
      // very likely created by videojs-contrib-hls from information
      // in the m3u8 for us to use
      inbandTextTracks[captionStream] = track;
    } else {
      // Otherwise, create a track with the default `CC#` label and
      // without a language
      inbandTextTracks[captionStream] = tech.addRemoteTextTrack({
        kind: 'captions',
        id: captionStream,
        label: captionStream
      }, false).track;
    }
  }
};
/**
 * Add caption text track data to a source handler given an array of captions
 *
 * @param {Object}
 *   @param {Object} inbandTextTracks the inband text tracks
 *   @param {number} timestampOffset the timestamp offset of the source buffer
 *   @param {Array} captionArray an array of caption data
 * @private
 */

var addCaptionData = function addCaptionData(_ref) {
  var inbandTextTracks = _ref.inbandTextTracks,
      captionArray = _ref.captionArray,
      timestampOffset = _ref.timestampOffset;

  if (!captionArray) {
    return;
  }

  var Cue = window_1.WebKitDataCue || window_1.VTTCue;
  captionArray.forEach(function (caption) {
    var track = caption.stream;
    inbandTextTracks[track].addCue(new Cue(caption.startTime + timestampOffset, caption.endTime + timestampOffset, caption.text));
  });
};
/**
 * Define properties on a cue for backwards compatability,
 * but warn the user that the way that they are using it
 * is depricated and will be removed at a later date.
 *
 * @param {Cue} cue the cue to add the properties on
 * @private
 */

var deprecateOldCue = function deprecateOldCue(cue) {
  Object.defineProperties(cue.frame, {
    id: {
      get: function get() {
        videojs$1.log.warn('cue.frame.id is deprecated. Use cue.value.key instead.');
        return cue.value.key;
      }
    },
    value: {
      get: function get() {
        videojs$1.log.warn('cue.frame.value is deprecated. Use cue.value.data instead.');
        return cue.value.data;
      }
    },
    privateData: {
      get: function get() {
        videojs$1.log.warn('cue.frame.privateData is deprecated. Use cue.value.data instead.');
        return cue.value.data;
      }
    }
  });
};
/**
 * Add metadata text track data to a source handler given an array of metadata
 *
 * @param {Object}
 *   @param {Object} inbandTextTracks the inband text tracks
 *   @param {Array} metadataArray an array of meta data
 *   @param {number} timestampOffset the timestamp offset of the source buffer
 *   @param {number} videoDuration the duration of the video
 * @private
 */


var addMetadata = function addMetadata(_ref2) {
  var inbandTextTracks = _ref2.inbandTextTracks,
      metadataArray = _ref2.metadataArray,
      timestampOffset = _ref2.timestampOffset,
      videoDuration = _ref2.videoDuration;

  if (!metadataArray) {
    return;
  }

  var Cue = window_1.WebKitDataCue || window_1.VTTCue;
  var metadataTrack = inbandTextTracks.metadataTrack_;

  if (!metadataTrack) {
    return;
  }

  metadataArray.forEach(function (metadata) {
    var time = metadata.cueTime + timestampOffset; // if time isn't a finite number between 0 and Infinity, like NaN,
    // ignore this bit of metadata.
    // This likely occurs when you have an non-timed ID3 tag like TIT2,
    // which is the "Title/Songname/Content description" frame

    if (typeof time !== 'number' || window_1.isNaN(time) || time < 0 || !(time < Infinity)) {
      return;
    }

    metadata.frames.forEach(function (frame) {
      var cue = new Cue(time, time, frame.value || frame.url || frame.data || '');
      cue.frame = frame;
      cue.value = frame;
      deprecateOldCue(cue);
      metadataTrack.addCue(cue);
    });
  });

  if (!metadataTrack.cues || !metadataTrack.cues.length) {
    return;
  } // Updating the metadeta cues so that
  // the endTime of each cue is the startTime of the next cue
  // the endTime of last cue is the duration of the video


  var cues = metadataTrack.cues;
  var cuesArray = []; // Create a copy of the TextTrackCueList...
  // ...disregarding cues with a falsey value

  for (var i = 0; i < cues.length; i++) {
    if (cues[i]) {
      cuesArray.push(cues[i]);
    }
  } // Group cues by their startTime value


  var cuesGroupedByStartTime = cuesArray.reduce(function (obj, cue) {
    var timeSlot = obj[cue.startTime] || [];
    timeSlot.push(cue);
    obj[cue.startTime] = timeSlot;
    return obj;
  }, {}); // Sort startTimes by ascending order

  var sortedStartTimes = Object.keys(cuesGroupedByStartTime).sort(function (a, b) {
    return Number(a) - Number(b);
  }); // Map each cue group's endTime to the next group's startTime

  sortedStartTimes.forEach(function (startTime, idx) {
    var cueGroup = cuesGroupedByStartTime[startTime];
    var nextTime = Number(sortedStartTimes[idx + 1]) || videoDuration; // Map each cue's endTime the next group's startTime

    cueGroup.forEach(function (cue) {
      cue.endTime = nextTime;
    });
  });
};
/**
 * Create metadata text track on video.js if it does not exist
 *
 * @param {Object} inbandTextTracks a reference to current inbandTextTracks
 * @param {string} dispatchType the inband metadata track dispatch type
 * @param {Object} tech the video.js tech
 * @private
 */

var createMetadataTrackIfNotExists = function createMetadataTrackIfNotExists(inbandTextTracks, dispatchType, tech) {
  if (inbandTextTracks.metadataTrack_) {
    return;
  }

  inbandTextTracks.metadataTrack_ = tech.addRemoteTextTrack({
    kind: 'metadata',
    label: 'Timed Metadata'
  }, false).track;
  inbandTextTracks.metadataTrack_.inBandMetadataTrackDispatchType = dispatchType;
};
/**
 * Remove cues from a track on video.js.
 *
 * @param {Double} start start of where we should remove the cue
 * @param {Double} end end of where the we should remove the cue
 * @param {Object} track the text track to remove the cues from
 * @private
 */

var removeCuesFromTrack = function removeCuesFromTrack(start, end, track) {
  var i;
  var cue;

  if (!track) {
    return;
  }

  if (!track.cues) {
    return;
  }

  i = track.cues.length;

  while (i--) {
    cue = track.cues[i]; // Remove any overlapping cue

    if (cue.startTime >= start && cue.endTime <= end) {
      track.removeCue(cue);
    }
  }
};

/**
 * Returns a list of gops in the buffer that have a pts value of 3 seconds or more in
 * front of current time.
 *
 * @param {Array} buffer
 *        The current buffer of gop information
 * @param {number} currentTime
 *        The current time
 * @param {Double} mapping
 *        Offset to map display time to stream presentation time
 * @return {Array}
 *         List of gops considered safe to append over
 */

var gopsSafeToAlignWith = function gopsSafeToAlignWith(buffer, currentTime, mapping) {
  if (typeof currentTime === 'undefined' || currentTime === null || !buffer.length) {
    return [];
  } // pts value for current time + 3 seconds to give a bit more wiggle room


  var currentTimePts = Math.ceil((currentTime - mapping + 3) * clock.ONE_SECOND_IN_TS);
  var i;

  for (i = 0; i < buffer.length; i++) {
    if (buffer[i].pts > currentTimePts) {
      break;
    }
  }

  return buffer.slice(i);
};
/**
 * Appends gop information (timing and byteLength) received by the transmuxer for the
 * gops appended in the last call to appendBuffer
 *
 * @param {Array} buffer
 *        The current buffer of gop information
 * @param {Array} gops
 *        List of new gop information
 * @param {boolean} replace
 *        If true, replace the buffer with the new gop information. If false, append the
 *        new gop information to the buffer in the right location of time.
 * @return {Array}
 *         Updated list of gop information
 */

var updateGopBuffer = function updateGopBuffer(buffer, gops, replace) {
  if (!gops.length) {
    return buffer;
  }

  if (replace) {
    // If we are in safe append mode, then completely overwrite the gop buffer
    // with the most recent appeneded data. This will make sure that when appending
    // future segments, we only try to align with gops that are both ahead of current
    // time and in the last segment appended.
    return gops.slice();
  }

  var start = gops[0].pts;
  var i = 0;

  for (i; i < buffer.length; i++) {
    if (buffer[i].pts >= start) {
      break;
    }
  }

  return buffer.slice(0, i).concat(gops);
};
/**
 * Removes gop information in buffer that overlaps with provided start and end
 *
 * @param {Array} buffer
 *        The current buffer of gop information
 * @param {Double} start
 *        position to start the remove at
 * @param {Double} end
 *        position to end the remove at
 * @param {Double} mapping
 *        Offset to map display time to stream presentation time
 */

var removeGopBuffer = function removeGopBuffer(buffer, start, end, mapping) {
  var startPts = Math.ceil((start - mapping) * clock.ONE_SECOND_IN_TS);
  var endPts = Math.ceil((end - mapping) * clock.ONE_SECOND_IN_TS);
  var updatedBuffer = buffer.slice();
  var i = buffer.length;

  while (i--) {
    if (buffer[i].pts <= endPts) {
      break;
    }
  }

  if (i === -1) {
    // no removal because end of remove range is before start of buffer
    return updatedBuffer;
  }

  var j = i + 1;

  while (j--) {
    if (buffer[j].pts <= startPts) {
      break;
    }
  } // clamp remove range start to 0 index


  j = Math.max(j, 0);
  updatedBuffer.splice(j, i - j + 1);
  return updatedBuffer;
};

var shallowEqual = function shallowEqual(a, b) {
  // if both are undefined
  // or one or the other is undefined
  // they are not equal
  if (!a && !b || !a && b || a && !b) {
    return false;
  } // they are the same object and thus, equal


  if (a === b) {
    return true;
  } // sort keys so we can make sure they have
  // all the same keys later.


  var akeys = Object.keys(a).sort();
  var bkeys = Object.keys(b).sort(); // different number of keys, not equal

  if (akeys.length !== bkeys.length) {
    return false;
  }

  for (var i = 0; i < akeys.length; i++) {
    var key = akeys[i]; // different sorted keys, not equal

    if (key !== bkeys[i]) {
      return false;
    } // different values, not equal


    if (a[key] !== b[key]) {
      return false;
    }
  }

  return true;
};

var CHECK_BUFFER_DELAY = 500;

var finite = function finite(num) {
  return typeof num === 'number' && isFinite(num);
};

var illegalMediaSwitch = function illegalMediaSwitch(loaderType, startingMedia, trackInfo) {
  // Although these checks should most likely cover non 'main' types, for now it narrows
  // the scope of our checks.
  if (loaderType !== 'main' || !startingMedia || !trackInfo) {
    return null;
  }

  if (!trackInfo.hasAudio && !trackInfo.hasVideo) {
    return 'Neither audio nor video found in segment.';
  }

  if (startingMedia.hasVideo && !trackInfo.hasVideo) {
    return 'Only audio found in segment when we expected video.' + ' We can\'t switch to audio only from a stream that had video.' + ' To get rid of this message, please add codec information to the manifest.';
  }

  if (!startingMedia.hasVideo && trackInfo.hasVideo) {
    return 'Video found in segment when we expected only audio.' + ' We can\'t switch to a stream with video from an audio only stream.' + ' To get rid of this message, please add codec information to the manifest.';
  }

  return null;
};
/**
 * Calculates a time value that is safe to remove from the back buffer without interrupting
 * playback.
 *
 * @param {TimeRange} seekable
 *        The current seekable range
 * @param {number} currentTime
 *        The current time of the player
 * @param {number} targetDuration
 *        The target duration of the current playlist
 * @return {number}
 *         Time that is safe to remove from the back buffer without interrupting playback
 */

var safeBackBufferTrimTime = function safeBackBufferTrimTime(seekable, currentTime, targetDuration) {
  // 30 seconds before the playhead provides a safe default for trimming.
  //
  // Choosing a reasonable default is particularly important for high bitrate content and
  // VOD videos/live streams with large windows, as the buffer may end up overfilled and
  // throw an APPEND_BUFFER_ERR.
  var trimTime = currentTime - Config.BACK_BUFFER_LENGTH;

  if (seekable.length) {
    // Some live playlists may have a shorter window of content than the full allowed back
    // buffer. For these playlists, don't save content that's no longer within the window.
    trimTime = Math.max(trimTime, seekable.start(0));
  } // Don't remove within target duration of the current time to avoid the possibility of
  // removing the GOP currently being played, as removing it can cause playback stalls.


  var maxTrimTime = currentTime - targetDuration;
  return Math.min(maxTrimTime, trimTime);
};

var segmentInfoString = function segmentInfoString(segmentInfo) {
  var _segmentInfo$segment = segmentInfo.segment,
      start = _segmentInfo$segment.start,
      end = _segmentInfo$segment.end,
      _segmentInfo$playlist = segmentInfo.playlist,
      seq = _segmentInfo$playlist.mediaSequence,
      id = _segmentInfo$playlist.id,
      _segmentInfo$playlist2 = _segmentInfo$playlist.segments,
      segments = _segmentInfo$playlist2 === void 0 ? [] : _segmentInfo$playlist2,
      index = segmentInfo.mediaIndex,
      timeline = segmentInfo.timeline;
  return ["appending [" + index + "] of [" + seq + ", " + (seq + segments.length) + "] from playlist [" + id + "]", "[" + start + " => " + end + "] in timeline [" + timeline + "]"].join(' ');
};

var timingInfoPropertyForMedia = function timingInfoPropertyForMedia(mediaType) {
  return mediaType + "TimingInfo";
};
/**
 * Returns the timestamp offset to use for the segment.
 *
 * @param {number} segmentTimeline
 *        The timeline of the segment
 * @param {number} currentTimeline
 *        The timeline currently being followed by the loader
 * @param {number} startOfSegment
 *        The estimated segment start
 * @param {TimeRange[]} buffered
 *        The loader's buffer
 * @param {boolean} overrideCheck
 *        If true, no checks are made to see if the timestamp offset value should be set,
 *        but sets it directly to a value.
 *
 * @return {number|null}
 *         Either a number representing a new timestamp offset, or null if the segment is
 *         part of the same timeline
 */


var timestampOffsetForSegment = function timestampOffsetForSegment(_ref) {
  var segmentTimeline = _ref.segmentTimeline,
      currentTimeline = _ref.currentTimeline,
      startOfSegment = _ref.startOfSegment,
      buffered = _ref.buffered,
      overrideCheck = _ref.overrideCheck;

  // Check to see if we are crossing a discontinuity to see if we need to set the
  // timestamp offset on the transmuxer and source buffer.
  //
  // Previously, we changed the timestampOffset if the start of this segment was less than
  // the currently set timestampOffset, but this isn't desirable as it can produce bad
  // behavior, especially around long running live streams.
  if (!overrideCheck && segmentTimeline === currentTimeline) {
    return null;
  } // segmentInfo.startOfSegment used to be used as the timestamp offset, however, that
  // value uses the end of the last segment if it is available. While this value
  // should often be correct, it's better to rely on the buffered end, as the new
  // content post discontinuity should line up with the buffered end as if it were
  // time 0 for the new content.


  return buffered.length ? buffered.end(buffered.length - 1) : startOfSegment;
};
/**
 * Returns whether or not the loader should wait for a timeline change from the timeline
 * change controller before processing the segment.
 *
 * Primary timing in VHS goes by video. This is different from most media players, as
 * audio is more often used as the primary timing source. For the foreseeable future, VHS
 * will continue to use video as the primary timing source, due to the current logic and
 * expectations built around it.

 * Since the timing follows video, in order to maintain sync, the video loader is
 * responsible for setting both audio and video source buffer timestamp offsets.
 *
 * Setting different values for audio and video source buffers could lead to
 * desyncing. The following examples demonstrate some of the situations where this
 * distinction is important. Note that all of these cases involve demuxed content. When
 * content is muxed, the audio and video are packaged together, therefore syncing
 * separate media playlists is not an issue.
 *
 * CASE 1: Audio prepares to load a new timeline before video:
 *
 * Timeline:       0                 1
 * Audio Segments: 0 1 2 3 4 5 DISCO 6 7 8 9
 * Audio Loader:                     ^
 * Video Segments: 0 1 2 3 4 5 DISCO 6 7 8 9
 * Video Loader              ^
 *
 * In the above example, the audio loader is preparing to load the 6th segment, the first
 * after a discontinuity, while the video loader is still loading the 5th segment, before
 * the discontinuity.
 *
 * If the audio loader goes ahead and loads and appends the 6th segment before the video
 * loader crosses the discontinuity, then when appended, the 6th audio segment will use
 * the timestamp offset from timeline 0. This will likely lead to desyncing. In addition,
 * the audio loader must provide the audioAppendStart value to trim the content in the
 * transmuxer, and that value relies on the audio timestamp offset. Since the audio
 * timestamp offset is set by the video (main) loader, the audio loader shouldn't load the
 * segment until that value is provided.
 *
 * CASE 2: Video prepares to load a new timeline before audio:
 *
 * Timeline:       0                 1
 * Audio Segments: 0 1 2 3 4 5 DISCO 6 7 8 9
 * Audio Loader:             ^
 * Video Segments: 0 1 2 3 4 5 DISCO 6 7 8 9
 * Video Loader                      ^
 *
 * In the above example, the video loader is preparing to load the 6th segment, the first
 * after a discontinuity, while the audio loader is still loading the 5th segment, before
 * the discontinuity.
 *
 * If the video loader goes ahead and loads and appends the 6th segment, then once the
 * segment is loaded and processed, both the video and audio timestamp offsets will be
 * set, since video is used as the primary timing source. This is to ensure content lines
 * up appropriately, as any modifications to the video timing are reflected by audio when
 * the video loader sets the audio and video timestamp offsets to the same value. However,
 * setting the timestamp offset for audio before audio has had a chance to change
 * timelines will likely lead to desyncing, as the audio loader will append segment 5 with
 * a timestamp intended to apply to segments from timeline 1 rather than timeline 0.
 *
 * CASE 3: When seeking, audio prepares to load a new timeline before video
 *
 * Timeline:       0                 1
 * Audio Segments: 0 1 2 3 4 5 DISCO 6 7 8 9
 * Audio Loader:           ^
 * Video Segments: 0 1 2 3 4 5 DISCO 6 7 8 9
 * Video Loader            ^
 *
 * In the above example, both audio and video loaders are loading segments from timeline
 * 0, but imagine that the seek originated from timeline 1.
 *
 * When seeking to a new timeline, the timestamp offset will be set based on the expected
 * segment start of the loaded video segment. In order to maintain sync, the audio loader
 * must wait for the video loader to load its segment and update both the audio and video
 * timestamp offsets before it may load and append its own segment. This is the case
 * whether the seek results in a mismatched segment request (e.g., the audio loader
 * chooses to load segment 3 and the video loader chooses to load segment 4) or the
 * loaders choose to load the same segment index from each playlist, as the segments may
 * not be aligned perfectly, even for matching segment indexes.
 *
 * @param {Object} timelinechangeController
 * @param {number} currentTimeline
 *        The timeline currently being followed by the loader
 * @param {number} segmentTimeline
 *        The timeline of the segment being loaded
 * @param {('main'|'audio')} loaderType
 *        The loader type
 * @param {boolean} audioDisabled
 *        Whether the audio is disabled for the loader. This should only be true when the
 *        loader may have muxed audio in its segment, but should not append it, e.g., for
 *        the main loader when an alternate audio playlist is active.
 *
 * @return {boolean}
 *         Whether the loader should wait for a timeline change from the timeline change
 *         controller before processing the segment
 */

var shouldWaitForTimelineChange = function shouldWaitForTimelineChange(_ref2) {
  var timelineChangeController = _ref2.timelineChangeController,
      currentTimeline = _ref2.currentTimeline,
      segmentTimeline = _ref2.segmentTimeline,
      loaderType = _ref2.loaderType,
      audioDisabled = _ref2.audioDisabled;

  if (currentTimeline === segmentTimeline) {
    return false;
  }

  if (loaderType === 'audio') {
    var lastMainTimelineChange = timelineChangeController.lastTimelineChange({
      type: 'main'
    }); // Audio loader should wait if:
    //
    // * main hasn't had a timeline change yet (thus has not loaded its first segment)
    // * main hasn't yet changed to the timeline audio is looking to load

    return !lastMainTimelineChange || lastMainTimelineChange.to !== segmentTimeline;
  } // The main loader only needs to wait for timeline changes if there's demuxed audio.
  // Otherwise, there's nothing to wait for, since audio would be muxed into the main
  // loader's segments (or the content is audio/video only and handled by the main
  // loader).


  if (loaderType === 'main' && audioDisabled) {
    var pendingAudioTimelineChange = timelineChangeController.pendingTimelineChange({
      type: 'audio'
    }); // Main loader should wait for the audio loader if audio is not pending a timeline
    // change to the current timeline.
    //
    // Since the main loader is responsible for setting the timestamp offset for both
    // audio and video, the main loader must wait for audio to be about to change to its
    // timeline before setting the offset, otherwise, if audio is behind in loading,
    // segments from the previous timeline would be adjusted by the new timestamp offset.
    //
    // This requirement means that video will not cross a timeline until the audio is
    // about to cross to it, so that way audio and video will always cross the timeline
    // together.
    //
    // In addition to normal timeline changes, these rules also apply to the start of a
    // stream (going from a non-existent timeline, -1, to timeline 0). It's important
    // that these rules apply to the first timeline change because if they did not, it's
    // possible that the main loader will cross two timelines before the audio loader has
    // crossed one. Logic may be implemented to handle the startup as a special case, but
    // it's easier to simply treat all timeline changes the same.

    if (pendingAudioTimelineChange && pendingAudioTimelineChange.to === segmentTimeline) {
      return false;
    }

    return true;
  }

  return false;
};
/**
 * An object that manages segment loading and appending.
 *
 * @class SegmentLoader
 * @param {Object} options required and optional options
 * @extends videojs.EventTarget
 */

var SegmentLoader = /*#__PURE__*/function (_videojs$EventTarget) {
  inheritsLoose(SegmentLoader, _videojs$EventTarget);

  function SegmentLoader(settings, options) {
    var _this;

    _this = _videojs$EventTarget.call(this) || this; // check pre-conditions

    if (!settings) {
      throw new TypeError('Initialization settings are required');
    }

    if (typeof settings.currentTime !== 'function') {
      throw new TypeError('No currentTime getter specified');
    }

    if (!settings.mediaSource) {
      throw new TypeError('No MediaSource specified');
    } // public properties


    _this.bandwidth = settings.bandwidth;
    _this.throughput = {
      rate: 0,
      count: 0
    };
    _this.roundTrip = NaN;

    _this.resetStats_();

    _this.mediaIndex = null; // private settings

    _this.hasPlayed_ = settings.hasPlayed;
    _this.currentTime_ = settings.currentTime;
    _this.seekable_ = settings.seekable;
    _this.seeking_ = settings.seeking;
    _this.duration_ = settings.duration;
    _this.mediaSource_ = settings.mediaSource;
    _this.vhs_ = settings.vhs;
    _this.loaderType_ = settings.loaderType;
    _this.currentMediaInfo_ = void 0;
    _this.startingMediaInfo_ = void 0;
    _this.segmentMetadataTrack_ = settings.segmentMetadataTrack;
    _this.goalBufferLength_ = settings.goalBufferLength;
    _this.sourceType_ = settings.sourceType;
    _this.sourceUpdater_ = settings.sourceUpdater;
    _this.inbandTextTracks_ = settings.inbandTextTracks;
    _this.state_ = 'INIT';
    _this.handlePartialData_ = settings.handlePartialData;
    _this.timelineChangeController_ = settings.timelineChangeController;
    _this.shouldSaveSegmentTimingInfo_ = true; // private instance variables

    _this.checkBufferTimeout_ = null;
    _this.error_ = void 0;
    _this.currentTimeline_ = -1;
    _this.pendingSegment_ = null;
    _this.xhrOptions_ = null;
    _this.pendingSegments_ = [];
    _this.audioDisabled_ = false;
    _this.isPendingTimestampOffset_ = false; // TODO possibly move gopBuffer and timeMapping info to a separate controller

    _this.gopBuffer_ = [];
    _this.timeMapping_ = 0;
    _this.safeAppend_ = videojs$1.browser.IE_VERSION >= 11;
    _this.appendInitSegment_ = {
      audio: true,
      video: true
    };
    _this.playlistOfLastInitSegment_ = {
      audio: null,
      video: null
    };
    _this.callQueue_ = []; // If the segment loader prepares to load a segment, but does not have enough
    // information yet to start the loading process (e.g., if the audio loader wants to
    // load a segment from the next timeline but the main loader hasn't yet crossed that
    // timeline), then the load call will be added to the queue until it is ready to be
    // processed.

    _this.loadQueue_ = [];
    _this.metadataQueue_ = {
      id3: [],
      caption: []
    }; // Fragmented mp4 playback

    _this.activeInitSegmentId_ = null;
    _this.initSegments_ = {}; // HLSe playback

    _this.cacheEncryptionKeys_ = settings.cacheEncryptionKeys;
    _this.keyCache_ = {};
    _this.decrypter_ = settings.decrypter; // Manages the tracking and generation of sync-points, mappings
    // between a time in the display time and a segment index within
    // a playlist

    _this.syncController_ = settings.syncController;
    _this.syncPoint_ = {
      segmentIndex: 0,
      time: 0
    };
    _this.transmuxer_ = _this.createTransmuxer_();

    _this.triggerSyncInfoUpdate_ = function () {
      return _this.trigger('syncinfoupdate');
    };

    _this.syncController_.on('syncinfoupdate', _this.triggerSyncInfoUpdate_);

    _this.mediaSource_.addEventListener('sourceopen', function () {
      if (!_this.isEndOfStream_()) {
        _this.ended_ = false;
      }
    }); // ...for determining the fetch location


    _this.fetchAtBuffer_ = false;
    _this.logger_ = logger("SegmentLoader[" + _this.loaderType_ + "]");
    Object.defineProperty(assertThisInitialized(_this), 'state', {
      get: function get() {
        return this.state_;
      },
      set: function set(newState) {
        if (newState !== this.state_) {
          this.logger_(this.state_ + " -> " + newState);
          this.state_ = newState;
          this.trigger('statechange');
        }
      }
    });

    _this.sourceUpdater_.on('ready', function () {
      if (_this.hasEnoughInfoToAppend_()) {
        _this.processCallQueue_();
      }
    }); // Only the main loader needs to listen for pending timeline changes, as the main
    // loader should wait for audio to be ready to change its timeline so that both main
    // and audio timelines change together. For more details, see the
    // shouldWaitForTimelineChange function.


    if (_this.loaderType_ === 'main') {
      _this.timelineChangeController_.on('pendingtimelinechange', function () {
        if (_this.hasEnoughInfoToAppend_()) {
          _this.processCallQueue_();
        }
      });
    } // The main loader only listens on pending timeline changes, but the audio loader,
    // since its loads follow main, needs to listen on timeline changes. For more details,
    // see the shouldWaitForTimelineChange function.


    if (_this.loaderType_ === 'audio') {
      _this.timelineChangeController_.on('timelinechange', function () {
        if (_this.hasEnoughInfoToLoad_()) {
          _this.processLoadQueue_();
        }

        if (_this.hasEnoughInfoToAppend_()) {
          _this.processCallQueue_();
        }
      });
    }

    return _this;
  }

  var _proto = SegmentLoader.prototype;

  _proto.createTransmuxer_ = function createTransmuxer_() {
    var transmuxer = new TransmuxWorker();
    transmuxer.postMessage({
      action: 'init',
      options: {
        remux: false,
        alignGopsAtEnd: this.safeAppend_,
        keepOriginalTimestamps: true,
        handlePartialData: this.handlePartialData_
      }
    });
    return transmuxer;
  }
  /**
   * reset all of our media stats
   *
   * @private
   */
  ;

  _proto.resetStats_ = function resetStats_() {
    this.mediaBytesTransferred = 0;
    this.mediaRequests = 0;
    this.mediaRequestsAborted = 0;
    this.mediaRequestsTimedout = 0;
    this.mediaRequestsErrored = 0;
    this.mediaTransferDuration = 0;
    this.mediaSecondsLoaded = 0;
  }
  /**
   * dispose of the SegmentLoader and reset to the default state
   */
  ;

  _proto.dispose = function dispose() {
    this.trigger('dispose');
    this.state = 'DISPOSED';
    this.pause();
    this.abort_();

    if (this.transmuxer_) {
      this.transmuxer_.terminate(); // Although it isn't an instance of a class, the segment transmuxer must still be
      // cleaned up.

      segmentTransmuxer.dispose();
    }

    this.resetStats_();

    if (this.checkBufferTimeout_) {
      window_1.clearTimeout(this.checkBufferTimeout_);
    }

    if (this.syncController_ && this.triggerSyncInfoUpdate_) {
      this.syncController_.off('syncinfoupdate', this.triggerSyncInfoUpdate_);
    }

    this.off();
  };

  _proto.setAudio = function setAudio(enable) {
    this.audioDisabled_ = !enable;

    if (enable) {
      this.appendInitSegment_.audio = true;
    } else {
      // remove current track audio if it gets disabled
      this.sourceUpdater_.removeAudio(0, this.duration_());
    }
  }
  /**
   * abort anything that is currently doing on with the SegmentLoader
   * and reset to a default state
   */
  ;

  _proto.abort = function abort() {
    if (this.state !== 'WAITING') {
      if (this.pendingSegment_) {
        this.pendingSegment_ = null;
      }

      return;
    }

    this.abort_(); // We aborted the requests we were waiting on, so reset the loader's state to READY
    // since we are no longer "waiting" on any requests. XHR callback is not always run
    // when the request is aborted. This will prevent the loader from being stuck in the
    // WAITING state indefinitely.

    this.state = 'READY'; // don't wait for buffer check timeouts to begin fetching the
    // next segment

    if (!this.paused()) {
      this.monitorBuffer_();
    }
  }
  /**
   * abort all pending xhr requests and null any pending segements
   *
   * @private
   */
  ;

  _proto.abort_ = function abort_() {
    if (this.pendingSegment_ && this.pendingSegment_.abortRequests) {
      this.pendingSegment_.abortRequests();
    } // clear out the segment being processed


    this.pendingSegment_ = null;
    this.callQueue_ = [];
    this.loadQueue_ = [];
    this.metadataQueue_.id3 = [];
    this.metadataQueue_.caption = [];
    this.timelineChangeController_.clearPendingTimelineChange(this.loaderType_);
  };

  _proto.checkForAbort_ = function checkForAbort_(requestId) {
    // If the state is APPENDING, then aborts will not modify the state, meaning the first
    // callback that happens should reset the state to READY so that loading can continue.
    if (this.state === 'APPENDING' && !this.pendingSegment_) {
      this.state = 'READY';
      return true;
    }

    if (!this.pendingSegment_ || this.pendingSegment_.requestId !== requestId) {
      return true;
    }

    return false;
  }
  /**
   * set an error on the segment loader and null out any pending segements
   *
   * @param {Error} error the error to set on the SegmentLoader
   * @return {Error} the error that was set or that is currently set
   */
  ;

  _proto.error = function error(_error) {
    if (typeof _error !== 'undefined') {
      this.logger_('error occurred:', _error);
      this.error_ = _error;
    }

    this.pendingSegment_ = null;
    return this.error_;
  };

  _proto.endOfStream = function endOfStream() {
    this.ended_ = true;

    if (this.transmuxer_) {
      // need to clear out any cached data to prepare for the new segment
      segmentTransmuxer.reset(this.transmuxer_);
    }

    this.gopBuffer_.length = 0;
    this.pause();
    this.trigger('ended');
  }
  /**
   * Indicates which time ranges are buffered
   *
   * @return {TimeRange}
   *         TimeRange object representing the current buffered ranges
   */
  ;

  _proto.buffered_ = function buffered_() {
    if (!this.sourceUpdater_ || !this.startingMediaInfo_) {
      return videojs$1.createTimeRanges();
    }

    if (this.loaderType_ === 'main') {
      var _this$startingMediaIn = this.startingMediaInfo_,
          hasAudio = _this$startingMediaIn.hasAudio,
          hasVideo = _this$startingMediaIn.hasVideo,
          isMuxed = _this$startingMediaIn.isMuxed;

      if (hasVideo && hasAudio && !this.audioDisabled_ && !isMuxed) {
        return this.sourceUpdater_.buffered();
      }

      if (hasVideo) {
        return this.sourceUpdater_.videoBuffered();
      }
    } // One case that can be ignored for now is audio only with alt audio,
    // as we don't yet have proper support for that.


    return this.sourceUpdater_.audioBuffered();
  }
  /**
   * Gets and sets init segment for the provided map
   *
   * @param {Object} map
   *        The map object representing the init segment to get or set
   * @param {boolean=} set
   *        If true, the init segment for the provided map should be saved
   * @return {Object}
   *         map object for desired init segment
   */
  ;

  _proto.initSegmentForMap = function initSegmentForMap(map, set) {
    if (set === void 0) {
      set = false;
    }

    if (!map) {
      return null;
    }

    var id = initSegmentId(map);
    var storedMap = this.initSegments_[id];

    if (set && !storedMap && map.bytes) {
      this.initSegments_[id] = storedMap = {
        resolvedUri: map.resolvedUri,
        byterange: map.byterange,
        bytes: map.bytes,
        tracks: map.tracks,
        timescales: map.timescales
      };
    }

    return storedMap || map;
  }
  /**
   * Gets and sets key for the provided key
   *
   * @param {Object} key
   *        The key object representing the key to get or set
   * @param {boolean=} set
   *        If true, the key for the provided key should be saved
   * @return {Object}
   *         Key object for desired key
   */
  ;

  _proto.segmentKey = function segmentKey(key, set) {
    if (set === void 0) {
      set = false;
    }

    if (!key) {
      return null;
    }

    var id = segmentKeyId(key);
    var storedKey = this.keyCache_[id]; // TODO: We should use the HTTP Expires header to invalidate our cache per
    // https://tools.ietf.org/html/draft-pantos-http-live-streaming-23#section-6.2.3

    if (this.cacheEncryptionKeys_ && set && !storedKey && key.bytes) {
      this.keyCache_[id] = storedKey = {
        resolvedUri: key.resolvedUri,
        bytes: key.bytes
      };
    }

    var result = {
      resolvedUri: (storedKey || key).resolvedUri
    };

    if (storedKey) {
      result.bytes = storedKey.bytes;
    }

    return result;
  }
  /**
   * Returns true if all configuration required for loading is present, otherwise false.
   *
   * @return {boolean} True if the all configuration is ready for loading
   * @private
   */
  ;

  _proto.couldBeginLoading_ = function couldBeginLoading_() {
    return this.playlist_ && !this.paused();
  }
  /**
   * load a playlist and start to fill the buffer
   */
  ;

  _proto.load = function load() {
    // un-pause
    this.monitorBuffer_(); // if we don't have a playlist yet, keep waiting for one to be
    // specified

    if (!this.playlist_) {
      return;
    } // not sure if this is the best place for this


    this.syncController_.setDateTimeMapping(this.playlist_); // if all the configuration is ready, initialize and begin loading

    if (this.state === 'INIT' && this.couldBeginLoading_()) {
      return this.init_();
    } // if we're in the middle of processing a segment already, don't
    // kick off an additional segment request


    if (!this.couldBeginLoading_() || this.state !== 'READY' && this.state !== 'INIT') {
      return;
    }

    this.state = 'READY';
  }
  /**
   * Once all the starting parameters have been specified, begin
   * operation. This method should only be invoked from the INIT
   * state.
   *
   * @private
   */
  ;

  _proto.init_ = function init_() {
    this.state = 'READY'; // if this is the audio segment loader, and it hasn't been inited before, then any old
    // audio data from the muxed content should be removed

    this.resetEverything();
    return this.monitorBuffer_();
  }
  /**
   * set a playlist on the segment loader
   *
   * @param {PlaylistLoader} media the playlist to set on the segment loader
   */
  ;

  _proto.playlist = function playlist(newPlaylist, options) {
    if (options === void 0) {
      options = {};
    }

    if (!newPlaylist) {
      return;
    }

    var oldPlaylist = this.playlist_;
    var segmentInfo = this.pendingSegment_;
    this.playlist_ = newPlaylist;
    this.xhrOptions_ = options; // when we haven't started playing yet, the start of a live playlist
    // is always our zero-time so force a sync update each time the playlist
    // is refreshed from the server
    //
    // Use the INIT state to determine if playback has started, as the playlist sync info
    // should be fixed once requests begin (as sync points are generated based on sync
    // info), but not before then.

    if (this.state === 'INIT') {
      newPlaylist.syncInfo = {
        mediaSequence: newPlaylist.mediaSequence,
        time: 0
      };
    }

    var oldId = null;

    if (oldPlaylist) {
      if (oldPlaylist.id) {
        oldId = oldPlaylist.id;
      } else if (oldPlaylist.uri) {
        oldId = oldPlaylist.uri;
      }
    }

    this.logger_("playlist update [" + oldId + " => " + (newPlaylist.id || newPlaylist.uri) + "]"); // in VOD, this is always a rendition switch (or we updated our syncInfo above)
    // in LIVE, we always want to update with new playlists (including refreshes)

    this.trigger('syncinfoupdate'); // if we were unpaused but waiting for a playlist, start
    // buffering now

    if (this.state === 'INIT' && this.couldBeginLoading_()) {
      return this.init_();
    }

    if (!oldPlaylist || oldPlaylist.uri !== newPlaylist.uri) {
      if (this.mediaIndex !== null || this.handlePartialData_) {
        // we must "resync" the segment loader when we switch renditions and
        // the segment loader is already synced to the previous rendition
        //
        // or if we're handling partial data, we need to ensure the transmuxer is cleared
        // out before we start adding more data
        this.resyncLoader();
      }

      this.currentMediaInfo_ = void 0;
      this.trigger('playlistupdate'); // the rest of this function depends on `oldPlaylist` being defined

      return;
    } // we reloaded the same playlist so we are in a live scenario
    // and we will likely need to adjust the mediaIndex


    var mediaSequenceDiff = newPlaylist.mediaSequence - oldPlaylist.mediaSequence;
    this.logger_("live window shift [" + mediaSequenceDiff + "]"); // update the mediaIndex on the SegmentLoader
    // this is important because we can abort a request and this value must be
    // equal to the last appended mediaIndex

    if (this.mediaIndex !== null) {
      this.mediaIndex -= mediaSequenceDiff;
    } // update the mediaIndex on the SegmentInfo object
    // this is important because we will update this.mediaIndex with this value
    // in `handleAppendsDone_` after the segment has been successfully appended


    if (segmentInfo) {
      segmentInfo.mediaIndex -= mediaSequenceDiff; // we need to update the referenced segment so that timing information is
      // saved for the new playlist's segment, however, if the segment fell off the
      // playlist, we can leave the old reference and just lose the timing info

      if (segmentInfo.mediaIndex >= 0) {
        segmentInfo.segment = newPlaylist.segments[segmentInfo.mediaIndex];
      }
    }

    this.syncController_.saveExpiredSegmentInfo(oldPlaylist, newPlaylist);
  }
  /**
   * Prevent the loader from fetching additional segments. If there
   * is a segment request outstanding, it will finish processing
   * before the loader halts. A segment loader can be unpaused by
   * calling load().
   */
  ;

  _proto.pause = function pause() {
    if (this.checkBufferTimeout_) {
      window_1.clearTimeout(this.checkBufferTimeout_);
      this.checkBufferTimeout_ = null;
    }
  }
  /**
   * Returns whether the segment loader is fetching additional
   * segments when given the opportunity. This property can be
   * modified through calls to pause() and load().
   */
  ;

  _proto.paused = function paused() {
    return this.checkBufferTimeout_ === null;
  }
  /**
   * Delete all the buffered data and reset the SegmentLoader
   *
   * @param {Function} [done] an optional callback to be executed when the remove
   * operation is complete
   */
  ;

  _proto.resetEverything = function resetEverything(done) {
    this.ended_ = false;
    this.appendInitSegment_ = {
      audio: true,
      video: true
    };
    this.resetLoader(); // remove from 0, the earliest point, to Infinity, to signify removal of everything.
    // VTT Segment Loader doesn't need to do anything but in the regular SegmentLoader,
    // we then clamp the value to duration if necessary.

    this.remove(0, Infinity, done); // clears fmp4 captions

    if (this.transmuxer_) {
      this.transmuxer_.postMessage({
        action: 'clearAllMp4Captions'
      });
    }
  }
  /**
   * Force the SegmentLoader to resync and start loading around the currentTime instead
   * of starting at the end of the buffer
   *
   * Useful for fast quality changes
   */
  ;

  _proto.resetLoader = function resetLoader() {
    this.fetchAtBuffer_ = false;
    this.resyncLoader();
  }
  /**
   * Force the SegmentLoader to restart synchronization and make a conservative guess
   * before returning to the simple walk-forward method
   */
  ;

  _proto.resyncLoader = function resyncLoader() {
    if (this.transmuxer_) {
      // need to clear out any cached data to prepare for the new segment
      segmentTransmuxer.reset(this.transmuxer_);
    }

    this.mediaIndex = null;
    this.syncPoint_ = null;
    this.isPendingTimestampOffset_ = false;
    this.callQueue_ = [];
    this.loadQueue_ = [];
    this.metadataQueue_.id3 = [];
    this.metadataQueue_.caption = [];
    this.abort();

    if (this.transmuxer_) {
      this.transmuxer_.postMessage({
        action: 'clearParsedMp4Captions'
      });
    }
  }
  /**
   * Remove any data in the source buffer between start and end times
   *
   * @param {number} start - the start time of the region to remove from the buffer
   * @param {number} end - the end time of the region to remove from the buffer
   * @param {Function} [done] - an optional callback to be executed when the remove
   * operation is complete
   */
  ;

  _proto.remove = function remove(start, end, done) {
    if (done === void 0) {
      done = function done() {};
    }

    // clamp end to duration if we need to remove everything.
    // This is due to a browser bug that causes issues if we remove to Infinity.
    // videojs/videojs-contrib-hls#1225
    if (end === Infinity) {
      end = this.duration_();
    }

    if (!this.sourceUpdater_ || !this.currentMediaInfo_) {
      // nothing to remove if we haven't processed any media
      return;
    } // set it to one to complete this function's removes


    var removesRemaining = 1;

    var removeFinished = function removeFinished() {
      removesRemaining--;

      if (removesRemaining === 0) {
        done();
      }
    };

    if (!this.audioDisabled_) {
      removesRemaining++;
      this.sourceUpdater_.removeAudio(start, end, removeFinished);
    }

    if (this.loaderType_ === 'main' && this.currentMediaInfo_ && this.currentMediaInfo_.hasVideo) {
      this.gopBuffer_ = removeGopBuffer(this.gopBuffer_, start, end, this.timeMapping_);
      removesRemaining++;
      this.sourceUpdater_.removeVideo(start, end, removeFinished);
    } // remove any captions and ID3 tags


    for (var track in this.inbandTextTracks_) {
      removeCuesFromTrack(start, end, this.inbandTextTracks_[track]);
    }

    removeCuesFromTrack(start, end, this.segmentMetadataTrack_); // finished this function's removes

    removeFinished();
  }
  /**
   * (re-)schedule monitorBufferTick_ to run as soon as possible
   *
   * @private
   */
  ;

  _proto.monitorBuffer_ = function monitorBuffer_() {
    if (this.checkBufferTimeout_) {
      window_1.clearTimeout(this.checkBufferTimeout_);
    }

    this.checkBufferTimeout_ = window_1.setTimeout(this.monitorBufferTick_.bind(this), 1);
  }
  /**
   * As long as the SegmentLoader is in the READY state, periodically
   * invoke fillBuffer_().
   *
   * @private
   */
  ;

  _proto.monitorBufferTick_ = function monitorBufferTick_() {
    if (this.state === 'READY') {
      this.fillBuffer_();
    }

    if (this.checkBufferTimeout_) {
      window_1.clearTimeout(this.checkBufferTimeout_);
    }

    this.checkBufferTimeout_ = window_1.setTimeout(this.monitorBufferTick_.bind(this), CHECK_BUFFER_DELAY);
  }
  /**
   * fill the buffer with segements unless the sourceBuffers are
   * currently updating
   *
   * Note: this function should only ever be called by monitorBuffer_
   * and never directly
   *
   * @private
   */
  ;

  _proto.fillBuffer_ = function fillBuffer_() {
    // TODO since the source buffer maintains a queue, and we shouldn't call this function
    // except when we're ready for the next segment, this check can most likely be removed
    if (this.sourceUpdater_.updating()) {
      return;
    }

    if (!this.syncPoint_) {
      this.syncPoint_ = this.syncController_.getSyncPoint(this.playlist_, this.duration_(), this.currentTimeline_, this.currentTime_());
    }

    var buffered = this.buffered_(); // see if we need to begin loading immediately

    var segmentInfo = this.checkBuffer_(buffered, this.playlist_, this.mediaIndex, this.hasPlayed_(), this.currentTime_(), this.syncPoint_);

    if (!segmentInfo) {
      return;
    }

    segmentInfo.timestampOffset = timestampOffsetForSegment({
      segmentTimeline: segmentInfo.timeline,
      currentTimeline: this.currentTimeline_,
      startOfSegment: segmentInfo.startOfSegment,
      buffered: buffered,
      overrideCheck: this.isPendingTimestampOffset_
    });
    this.isPendingTimestampOffset_ = false;

    if (typeof segmentInfo.timestampOffset === 'number') {
      this.timelineChangeController_.pendingTimelineChange({
        type: this.loaderType_,
        from: this.currentTimeline_,
        to: segmentInfo.timeline
      });
    }

    this.loadSegment_(segmentInfo);
  }
  /**
   * Determines if we should call endOfStream on the media source based
   * on the state of the buffer or if appened segment was the final
   * segment in the playlist.
   *
   * @param {number} [mediaIndex] the media index of segment we last appended
   * @param {Object} [playlist] a media playlist object
   * @return {boolean} do we need to call endOfStream on the MediaSource
   */
  ;

  _proto.isEndOfStream_ = function isEndOfStream_(mediaIndex, playlist) {
    if (mediaIndex === void 0) {
      mediaIndex = this.mediaIndex;
    }

    if (playlist === void 0) {
      playlist = this.playlist_;
    }

    if (!playlist || !this.mediaSource_) {
      return false;
    } // mediaIndex is zero based but length is 1 based


    var appendedLastSegment = mediaIndex + 1 === playlist.segments.length; // if we've buffered to the end of the video, we need to call endOfStream
    // so that MediaSources can trigger the `ended` event when it runs out of
    // buffered data instead of waiting for me

    return playlist.endList && this.mediaSource_.readyState === 'open' && appendedLastSegment;
  }
  /**
   * Determines what segment request should be made, given current playback
   * state.
   *
   * @param {TimeRanges} buffered - the state of the buffer
   * @param {Object} playlist - the playlist object to fetch segments from
   * @param {number} mediaIndex - the previous mediaIndex fetched or null
   * @param {boolean} hasPlayed - a flag indicating whether we have played or not
   * @param {number} currentTime - the playback position in seconds
   * @param {Object} syncPoint - a segment info object that describes the
   * @return {Object} a segment request object that describes the segment to load
   */
  ;

  _proto.checkBuffer_ = function checkBuffer_(buffered, playlist, currentMediaIndex, hasPlayed, currentTime, syncPoint) {
    var lastBufferedEnd = 0;

    if (buffered.length) {
      lastBufferedEnd = buffered.end(buffered.length - 1);
    }

    var bufferedTime = Math.max(0, lastBufferedEnd - currentTime);

    if (!playlist.segments.length) {
      return null;
    } // if there is plenty of content buffered, and the video has
    // been played before relax for awhile


    if (bufferedTime >= this.goalBufferLength_()) {
      return null;
    } // if the video has not yet played once, and we already have
    // one segment downloaded do nothing


    if (!hasPlayed && bufferedTime >= 1) {
      return null;
    }

    var nextMediaIndex = null;
    var startOfSegment;
    var isSyncRequest = false; // When the syncPoint is null, there is no way of determining a good
    // conservative segment index to fetch from
    // The best thing to do here is to get the kind of sync-point data by
    // making a request

    if (syncPoint === null) {
      nextMediaIndex = this.getSyncSegmentCandidate_(playlist);
      isSyncRequest = true;
    } else if (currentMediaIndex !== null) {
      // Under normal playback conditions fetching is a simple walk forward
      var segment = playlist.segments[currentMediaIndex];

      if (segment && segment.end) {
        startOfSegment = segment.end;
      } else {
        startOfSegment = lastBufferedEnd;
      }

      nextMediaIndex = currentMediaIndex + 1; // There is a sync-point but the lack of a mediaIndex indicates that
      // we need to make a good conservative guess about which segment to
      // fetch
    } else if (this.fetchAtBuffer_) {
      // Find the segment containing the end of the buffer
      var mediaSourceInfo = Playlist.getMediaInfoForTime(playlist, lastBufferedEnd, syncPoint.segmentIndex, syncPoint.time);
      nextMediaIndex = mediaSourceInfo.mediaIndex;
      startOfSegment = mediaSourceInfo.startTime;
    } else {
      // Find the segment containing currentTime
      var _mediaSourceInfo = Playlist.getMediaInfoForTime(playlist, currentTime, syncPoint.segmentIndex, syncPoint.time);

      nextMediaIndex = _mediaSourceInfo.mediaIndex;
      startOfSegment = _mediaSourceInfo.startTime;
    }

    var segmentInfo = this.generateSegmentInfo_(playlist, nextMediaIndex, startOfSegment, isSyncRequest);

    if (!segmentInfo) {
      return;
    } // if this is the last segment in the playlist
    // we are not seeking and end of stream has already been called
    // do not re-request


    if (this.mediaSource_ && this.playlist_ && segmentInfo.mediaIndex === this.playlist_.segments.length - 1 && this.mediaSource_.readyState === 'ended' && !this.seeking_()) {
      return;
    }

    this.logger_("checkBuffer_ returning " + segmentInfo.uri, {
      segmentInfo: segmentInfo,
      playlist: playlist,
      currentMediaIndex: currentMediaIndex,
      nextMediaIndex: nextMediaIndex,
      startOfSegment: startOfSegment,
      isSyncRequest: isSyncRequest
    });
    return segmentInfo;
  }
  /**
   * The segment loader has no recourse except to fetch a segment in the
   * current playlist and use the internal timestamps in that segment to
   * generate a syncPoint. This function returns a good candidate index
   * for that process.
   *
   * @param {Object} playlist - the playlist object to look for a
   * @return {number} An index of a segment from the playlist to load
   */
  ;

  _proto.getSyncSegmentCandidate_ = function getSyncSegmentCandidate_(playlist) {
    var _this2 = this;

    if (this.currentTimeline_ === -1) {
      return 0;
    }

    var segmentIndexArray = playlist.segments.map(function (s, i) {
      return {
        timeline: s.timeline,
        segmentIndex: i
      };
    }).filter(function (s) {
      return s.timeline === _this2.currentTimeline_;
    });

    if (segmentIndexArray.length) {
      return segmentIndexArray[Math.min(segmentIndexArray.length - 1, 1)].segmentIndex;
    }

    return Math.max(playlist.segments.length - 1, 0);
  };

  _proto.generateSegmentInfo_ = function generateSegmentInfo_(playlist, mediaIndex, startOfSegment, isSyncRequest) {
    if (mediaIndex < 0 || mediaIndex >= playlist.segments.length) {
      return null;
    }

    var segment = playlist.segments[mediaIndex];
    var audioBuffered = this.sourceUpdater_.audioBuffered();
    var videoBuffered = this.sourceUpdater_.videoBuffered();
    var audioAppendStart;
    var gopsToAlignWith;

    if (audioBuffered.length) {
      // since the transmuxer is using the actual timing values, but the buffer is
      // adjusted by the timestamp offset, we must adjust the value here
      audioAppendStart = audioBuffered.end(audioBuffered.length - 1) - this.sourceUpdater_.audioTimestampOffset();
    }

    if (videoBuffered.length) {
      gopsToAlignWith = gopsSafeToAlignWith(this.gopBuffer_, // since the transmuxer is using the actual timing values, but the time is
      // adjusted by the timestmap offset, we must adjust the value here
      this.currentTime_() - this.sourceUpdater_.videoTimestampOffset(), this.timeMapping_);
    }

    return {
      requestId: 'segment-loader-' + Math.random(),
      // resolve the segment URL relative to the playlist
      uri: segment.resolvedUri,
      // the segment's mediaIndex at the time it was requested
      mediaIndex: mediaIndex,
      // whether or not to update the SegmentLoader's state with this
      // segment's mediaIndex
      isSyncRequest: isSyncRequest,
      startOfSegment: startOfSegment,
      // the segment's playlist
      playlist: playlist,
      // unencrypted bytes of the segment
      bytes: null,
      // when a key is defined for this segment, the encrypted bytes
      encryptedBytes: null,
      // The target timestampOffset for this segment when we append it
      // to the source buffer
      timestampOffset: null,
      // The timeline that the segment is in
      timeline: segment.timeline,
      // The expected duration of the segment in seconds
      duration: segment.duration,
      // retain the segment in case the playlist updates while doing an async process
      segment: segment,
      byteLength: 0,
      transmuxer: this.transmuxer_,
      audioAppendStart: audioAppendStart,
      gopsToAlignWith: gopsToAlignWith
    };
  }
  /**
   * Determines if the network has enough bandwidth to complete the current segment
   * request in a timely manner. If not, the request will be aborted early and bandwidth
   * updated to trigger a playlist switch.
   *
   * @param {Object} stats
   *        Object containing stats about the request timing and size
   * @return {boolean} True if the request was aborted, false otherwise
   * @private
   */
  ;

  _proto.abortRequestEarly_ = function abortRequestEarly_(stats) {
    if (this.vhs_.tech_.paused() || // Don't abort if the current playlist is on the lowestEnabledRendition
    // TODO: Replace using timeout with a boolean indicating whether this playlist is
    //       the lowestEnabledRendition.
    !this.xhrOptions_.timeout || // Don't abort if we have no bandwidth information to estimate segment sizes
    !this.playlist_.attributes.BANDWIDTH) {
      return false;
    } // Wait at least 1 second since the first byte of data has been received before
    // using the calculated bandwidth from the progress event to allow the bitrate
    // to stabilize


    if (Date.now() - (stats.firstBytesReceivedAt || Date.now()) < 1000) {
      return false;
    }

    var currentTime = this.currentTime_();
    var measuredBandwidth = stats.bandwidth;
    var segmentDuration = this.pendingSegment_.duration;
    var requestTimeRemaining = Playlist.estimateSegmentRequestTime(segmentDuration, measuredBandwidth, this.playlist_, stats.bytesReceived); // Subtract 1 from the timeUntilRebuffer so we still consider an early abort
    // if we are only left with less than 1 second when the request completes.
    // A negative timeUntilRebuffering indicates we are already rebuffering

    var timeUntilRebuffer$1 = timeUntilRebuffer(this.buffered_(), currentTime, this.vhs_.tech_.playbackRate()) - 1; // Only consider aborting early if the estimated time to finish the download
    // is larger than the estimated time until the player runs out of forward buffer

    if (requestTimeRemaining <= timeUntilRebuffer$1) {
      return false;
    }

    var switchCandidate = minRebufferMaxBandwidthSelector({
      master: this.vhs_.playlists.master,
      currentTime: currentTime,
      bandwidth: measuredBandwidth,
      duration: this.duration_(),
      segmentDuration: segmentDuration,
      timeUntilRebuffer: timeUntilRebuffer$1,
      currentTimeline: this.currentTimeline_,
      syncController: this.syncController_
    });

    if (!switchCandidate) {
      return;
    }

    var rebufferingImpact = requestTimeRemaining - timeUntilRebuffer$1;
    var timeSavedBySwitching = rebufferingImpact - switchCandidate.rebufferingImpact;
    var minimumTimeSaving = 0.5; // If we are already rebuffering, increase the amount of variance we add to the
    // potential round trip time of the new request so that we are not too aggressive
    // with switching to a playlist that might save us a fraction of a second.

    if (timeUntilRebuffer$1 <= TIME_FUDGE_FACTOR) {
      minimumTimeSaving = 1;
    }

    if (!switchCandidate.playlist || switchCandidate.playlist.uri === this.playlist_.uri || timeSavedBySwitching < minimumTimeSaving) {
      return false;
    } // set the bandwidth to that of the desired playlist being sure to scale by
    // BANDWIDTH_VARIANCE and add one so the playlist selector does not exclude it
    // don't trigger a bandwidthupdate as the bandwidth is artifial


    this.bandwidth = switchCandidate.playlist.attributes.BANDWIDTH * Config.BANDWIDTH_VARIANCE + 1;
    this.abort();
    this.trigger('earlyabort');
    return true;
  };

  _proto.handleAbort_ = function handleAbort_() {
    this.mediaRequestsAborted += 1;
  }
  /**
   * XHR `progress` event handler
   *
   * @param {Event}
   *        The XHR `progress` event
   * @param {Object} simpleSegment
   *        A simplified segment object copy
   * @private
   */
  ;

  _proto.handleProgress_ = function handleProgress_(event, simpleSegment) {
    if (this.checkForAbort_(simpleSegment.requestId) || this.abortRequestEarly_(simpleSegment.stats)) {
      return;
    }

    this.trigger('progress');
  };

  _proto.handleTrackInfo_ = function handleTrackInfo_(simpleSegment, trackInfo) {
    if (this.checkForAbort_(simpleSegment.requestId) || this.abortRequestEarly_(simpleSegment.stats)) {
      return;
    }

    if (this.checkForIllegalMediaSwitch(trackInfo)) {
      return;
    }

    trackInfo = trackInfo || {}; // When we have track info, determine what media types this loader is dealing with.
    // Guard against cases where we're not getting track info at all until we are
    // certain that all streams will provide it.

    if (!shallowEqual(this.currentMediaInfo_, trackInfo)) {
      this.appendInitSegment_ = {
        audio: true,
        video: true
      };
      this.startingMediaInfo_ = trackInfo;
      this.currentMediaInfo_ = trackInfo;
      this.logger_('trackinfo update', trackInfo);
      this.trigger('trackinfo');
    } // trackinfo may cause an abort if the trackinfo
    // causes a codec change to an unsupported codec.


    if (this.checkForAbort_(simpleSegment.requestId) || this.abortRequestEarly_(simpleSegment.stats)) {
      return;
    } // set trackinfo on the pending segment so that
    // it can append.


    this.pendingSegment_.trackInfo = trackInfo; // check if any calls were waiting on the track info

    if (this.hasEnoughInfoToAppend_()) {
      this.processCallQueue_();
    }
  };

  _proto.handleTimingInfo_ = function handleTimingInfo_(simpleSegment, mediaType, timeType, time) {
    if (this.checkForAbort_(simpleSegment.requestId) || this.abortRequestEarly_(simpleSegment.stats)) {
      return;
    }

    var segmentInfo = this.pendingSegment_;
    var timingInfoProperty = timingInfoPropertyForMedia(mediaType);
    segmentInfo[timingInfoProperty] = segmentInfo[timingInfoProperty] || {};
    segmentInfo[timingInfoProperty][timeType] = time;
    this.logger_("timinginfo: " + mediaType + " - " + timeType + " - " + time); // check if any calls were waiting on the timing info

    if (this.hasEnoughInfoToAppend_()) {
      this.processCallQueue_();
    }
  };

  _proto.handleCaptions_ = function handleCaptions_(simpleSegment, captionData) {
    var _this3 = this;

    if (this.checkForAbort_(simpleSegment.requestId) || this.abortRequestEarly_(simpleSegment.stats)) {
      return;
    } // This could only happen with fmp4 segments, but
    // should still not happen in general


    if (captionData.length === 0) {
      this.logger_('SegmentLoader received no captions from a caption event');
      return;
    }

    var segmentInfo = this.pendingSegment_; // Wait until we have some video data so that caption timing
    // can be adjusted by the timestamp offset

    if (!segmentInfo.hasAppendedData_) {
      this.metadataQueue_.caption.push(this.handleCaptions_.bind(this, simpleSegment, captionData));
      return;
    }

    var timestampOffset = this.sourceUpdater_.videoTimestampOffset() === null ? this.sourceUpdater_.audioTimestampOffset() : this.sourceUpdater_.videoTimestampOffset();
    var captionTracks = {}; // get total start/end and captions for each track/stream

    captionData.forEach(function (caption) {
      // caption.stream is actually a track name...
      // set to the existing values in tracks or default values
      captionTracks[caption.stream] = captionTracks[caption.stream] || {
        // Infinity, as any other value will be less than this
        startTime: Infinity,
        captions: [],
        // 0 as an other value will be more than this
        endTime: 0
      };
      var captionTrack = captionTracks[caption.stream];
      captionTrack.startTime = Math.min(captionTrack.startTime, caption.startTime + timestampOffset);
      captionTrack.endTime = Math.max(captionTrack.endTime, caption.endTime + timestampOffset);
      captionTrack.captions.push(caption);
    });
    Object.keys(captionTracks).forEach(function (trackName) {
      var _captionTracks$trackN = captionTracks[trackName],
          startTime = _captionTracks$trackN.startTime,
          endTime = _captionTracks$trackN.endTime,
          captions = _captionTracks$trackN.captions;
      var inbandTextTracks = _this3.inbandTextTracks_;

      _this3.logger_("adding cues from " + startTime + " -> " + endTime + " for " + trackName);

      createCaptionsTrackIfNotExists(inbandTextTracks, _this3.vhs_.tech_, trackName); // clear out any cues that start and end at the same time period for the same track.
      // We do this because a rendition change that also changes the timescale for captions
      // will result in captions being re-parsed for certain segments. If we add them again
      // without clearing we will have two of the same captions visible.

      removeCuesFromTrack(startTime, endTime, inbandTextTracks[trackName]);
      addCaptionData({
        captionArray: captions,
        inbandTextTracks: inbandTextTracks,
        timestampOffset: timestampOffset
      });
    }); // Reset stored captions since we added parsed
    // captions to a text track at this point

    if (this.transmuxer_) {
      this.transmuxer_.postMessage({
        action: 'clearParsedMp4Captions'
      });
    }
  };

  _proto.handleId3_ = function handleId3_(simpleSegment, id3Frames, dispatchType) {
    if (this.checkForAbort_(simpleSegment.requestId) || this.abortRequestEarly_(simpleSegment.stats)) {
      return;
    }

    var segmentInfo = this.pendingSegment_; // we need to have appended data in order for the timestamp offset to be set

    if (!segmentInfo.hasAppendedData_) {
      this.metadataQueue_.id3.push(this.handleId3_.bind(this, simpleSegment, id3Frames, dispatchType));
      return;
    }

    var timestampOffset = this.sourceUpdater_.videoTimestampOffset() === null ? this.sourceUpdater_.audioTimestampOffset() : this.sourceUpdater_.videoTimestampOffset(); // There's potentially an issue where we could double add metadata if there's a muxed
    // audio/video source with a metadata track, and an alt audio with a metadata track.
    // However, this probably won't happen, and if it does it can be handled then.

    createMetadataTrackIfNotExists(this.inbandTextTracks_, dispatchType, this.vhs_.tech_);
    addMetadata({
      inbandTextTracks: this.inbandTextTracks_,
      metadataArray: id3Frames,
      timestampOffset: timestampOffset,
      videoDuration: this.duration_()
    });
  };

  _proto.processMetadataQueue_ = function processMetadataQueue_() {
    this.metadataQueue_.id3.forEach(function (fn) {
      return fn();
    });
    this.metadataQueue_.caption.forEach(function (fn) {
      return fn();
    });
    this.metadataQueue_.id3 = [];
    this.metadataQueue_.caption = [];
  };

  _proto.processCallQueue_ = function processCallQueue_() {
    var callQueue = this.callQueue_; // Clear out the queue before the queued functions are run, since some of the
    // functions may check the length of the load queue and default to pushing themselves
    // back onto the queue.

    this.callQueue_ = [];
    callQueue.forEach(function (fun) {
      return fun();
    });
  };

  _proto.processLoadQueue_ = function processLoadQueue_() {
    var loadQueue = this.loadQueue_; // Clear out the queue before the queued functions are run, since some of the
    // functions may check the length of the load queue and default to pushing themselves
    // back onto the queue.

    this.loadQueue_ = [];
    loadQueue.forEach(function (fun) {
      return fun();
    });
  }
  /**
   * Determines whether the loader has enough info to load the next segment.
   *
   * @return {boolean}
   *         Whether or not the loader has enough info to load the next segment
   */
  ;

  _proto.hasEnoughInfoToLoad_ = function hasEnoughInfoToLoad_() {
    // Since primary timing goes by video, only the audio loader potentially needs to wait
    // to load.
    if (this.loaderType_ !== 'audio') {
      return true;
    }

    var segmentInfo = this.pendingSegment_; // A fill buffer must have already run to establish a pending segment before there's
    // enough info to load.

    if (!segmentInfo) {
      return false;
    } // The first segment can and should be loaded immediately so that source buffers are
    // created together (before appending). Source buffer creation uses the presence of
    // audio and video data to determine whether to create audio/video source buffers, and
    // uses processed (transmuxed or parsed) media to determine the types required.


    if (!this.currentMediaInfo_) {
      return true;
    }

    if ( // Technically, instead of waiting to load a segment on timeline changes, a segment
    // can be requested and downloaded and only wait before it is transmuxed or parsed.
    // But in practice, there are a few reasons why it is better to wait until a loader
    // is ready to append that segment before requesting and downloading:
    //
    // 1. Because audio and main loaders cross discontinuities together, if this loader
    //    is waiting for the other to catch up, then instead of requesting another
    //    segment and using up more bandwidth, by not yet loading, more bandwidth is
    //    allotted to the loader currently behind.
    // 2. media-segment-request doesn't have to have logic to consider whether a segment
    // is ready to be processed or not, isolating the queueing behavior to the loader.
    // 3. The audio loader bases some of its segment properties on timing information
    //    provided by the main loader, meaning that, if the logic for waiting on
    //    processing was in media-segment-request, then it would also need to know how
    //    to re-generate the segment information after the main loader caught up.
    shouldWaitForTimelineChange({
      timelineChangeController: this.timelineChangeController_,
      currentTimeline: this.currentTimeline_,
      segmentTimeline: segmentInfo.timeline,
      loaderType: this.loaderType_,
      audioDisabled: this.audioDisabled_
    })) {
      return false;
    }

    return true;
  };

  _proto.hasEnoughInfoToAppend_ = function hasEnoughInfoToAppend_() {
    if (!this.sourceUpdater_.ready()) {
      // waiting on one of the segment loaders to get enough data to create source buffers
      return false;
    }

    var segmentInfo = this.pendingSegment_; // no segment to append any data for or
    // we do not have information on this specific
    // segment yet

    if (!segmentInfo || !segmentInfo.trackInfo) {
      return false;
    }

    if (!this.handlePartialData_) {
      var _this$currentMediaInf = this.currentMediaInfo_,
          hasAudio = _this$currentMediaInf.hasAudio,
          hasVideo = _this$currentMediaInf.hasVideo,
          isMuxed = _this$currentMediaInf.isMuxed;

      if (hasVideo && !segmentInfo.videoTimingInfo) {
        return false;
      } // muxed content only relies on video timing information for now.


      if (hasAudio && !this.audioDisabled_ && !isMuxed && !segmentInfo.audioTimingInfo) {
        return false;
      }
    }

    if (shouldWaitForTimelineChange({
      timelineChangeController: this.timelineChangeController_,
      currentTimeline: this.currentTimeline_,
      segmentTimeline: segmentInfo.timeline,
      loaderType: this.loaderType_,
      audioDisabled: this.audioDisabled_
    })) {
      return false;
    }

    return true;
  };

  _proto.handleData_ = function handleData_(simpleSegment, result) {
    if (this.checkForAbort_(simpleSegment.requestId) || this.abortRequestEarly_(simpleSegment.stats)) {
      return;
    } // If there's anything in the call queue, then this data came later and should be
    // executed after the calls currently queued.


    if (this.callQueue_.length || !this.hasEnoughInfoToAppend_()) {
      this.callQueue_.push(this.handleData_.bind(this, simpleSegment, result));
      return;
    }

    var segmentInfo = this.pendingSegment_; // update the time mapping so we can translate from display time to media time

    this.setTimeMapping_(segmentInfo.timeline); // for tracking overall stats

    this.updateMediaSecondsLoaded_(segmentInfo.segment); // Note that the state isn't changed from loading to appending. This is because abort
    // logic may change behavior depending on the state, and changing state too early may
    // inflate our estimates of bandwidth. In the future this should be re-examined to
    // note more granular states.
    // don't process and append data if the mediaSource is closed

    if (this.mediaSource_.readyState === 'closed') {
      return;
    } // if this request included an initialization segment, save that data
    // to the initSegment cache


    if (simpleSegment.map) {
      simpleSegment.map = this.initSegmentForMap(simpleSegment.map, true); // move over init segment properties to media request

      segmentInfo.segment.map = simpleSegment.map;
    } // if this request included a segment key, save that data in the cache


    if (simpleSegment.key) {
      this.segmentKey(simpleSegment.key, true);
    }

    segmentInfo.isFmp4 = simpleSegment.isFmp4;
    segmentInfo.timingInfo = segmentInfo.timingInfo || {};

    if (segmentInfo.isFmp4) {
      this.trigger('fmp4');
      segmentInfo.timingInfo.start = segmentInfo[timingInfoPropertyForMedia(result.type)].start;
    } else {
      var useVideoTimingInfo = this.loaderType_ === 'main' && this.currentMediaInfo_.hasVideo;
      var firstVideoFrameTimeForData;

      if (useVideoTimingInfo) {
        firstVideoFrameTimeForData = this.handlePartialData_ ? result.videoFramePtsTime : segmentInfo.videoTimingInfo.start;
      } // Segment loader knows more about segment timing than the transmuxer (in certain
      // aspects), so make any changes required for a more accurate start time.
      // Don't set the end time yet, as the segment may not be finished processing.


      segmentInfo.timingInfo.start = this.trueSegmentStart_({
        currentStart: segmentInfo.timingInfo.start,
        playlist: segmentInfo.playlist,
        mediaIndex: segmentInfo.mediaIndex,
        currentVideoTimestampOffset: this.sourceUpdater_.videoTimestampOffset(),
        useVideoTimingInfo: useVideoTimingInfo,
        firstVideoFrameTimeForData: firstVideoFrameTimeForData,
        videoTimingInfo: segmentInfo.videoTimingInfo,
        audioTimingInfo: segmentInfo.audioTimingInfo
      });
    } // Init segments for audio and video only need to be appended in certain cases. Now
    // that data is about to be appended, we can check the final cases to determine
    // whether we should append an init segment.


    this.updateAppendInitSegmentStatus(segmentInfo, result.type); // Timestamp offset should be updated once we get new data and have its timing info,
    // as we use the start of the segment to offset the best guess (playlist provided)
    // timestamp offset.

    this.updateSourceBufferTimestampOffset_(segmentInfo); // Save some state so that in the future anything waiting on first append (and/or
    // timestamp offset(s)) can process immediately. While the extra state isn't optimal,
    // we need some notion of whether the timestamp offset or other relevant information
    // has had a chance to be set.

    segmentInfo.hasAppendedData_ = true; // Now that the timestamp offset should be set, we can append any waiting ID3 tags.

    this.processMetadataQueue_();
    this.appendData_(segmentInfo, result);
  };

  _proto.updateAppendInitSegmentStatus = function updateAppendInitSegmentStatus(segmentInfo, type) {
    // alt audio doesn't manage timestamp offset
    if (this.loaderType_ === 'main' && typeof segmentInfo.timestampOffset === 'number' && // in the case that we're handling partial data, we don't want to append an init
    // segment for each chunk
    !segmentInfo.changedTimestampOffset) {
      // if the timestamp offset changed, the timeline may have changed, so we have to re-
      // append init segments
      this.appendInitSegment_ = {
        audio: true,
        video: true
      };
    }

    if (this.playlistOfLastInitSegment_[type] !== segmentInfo.playlist) {
      // make sure we append init segment on playlist changes, in case the media config
      // changed
      this.appendInitSegment_[type] = true;
    }
  };

  _proto.getInitSegmentAndUpdateState_ = function getInitSegmentAndUpdateState_(_ref3) {
    var type = _ref3.type,
        initSegment = _ref3.initSegment,
        map = _ref3.map,
        playlist = _ref3.playlist;

    // "The EXT-X-MAP tag specifies how to obtain the Media Initialization Section
    // (Section 3) required to parse the applicable Media Segments.  It applies to every
    // Media Segment that appears after it in the Playlist until the next EXT-X-MAP tag
    // or until the end of the playlist."
    // https://tools.ietf.org/html/draft-pantos-http-live-streaming-23#section-4.3.2.5
    if (map) {
      var id = initSegmentId(map);

      if (this.activeInitSegmentId_ === id) {
        // don't need to re-append the init segment if the ID matches
        return null;
      } // a map-specified init segment takes priority over any transmuxed (or otherwise
      // obtained) init segment
      //
      // this also caches the init segment for later use


      initSegment = this.initSegmentForMap(map, true).bytes;
      this.activeInitSegmentId_ = id;
    } // We used to always prepend init segments for video, however, that shouldn't be
    // necessary. Instead, we should only append on changes, similar to what we've always
    // done for audio. This is more important (though may not be that important) for
    // frame-by-frame appending for LHLS, simply because of the increased quantity of
    // appends.


    if (initSegment && this.appendInitSegment_[type]) {
      // Make sure we track the playlist that we last used for the init segment, so that
      // we can re-append the init segment in the event that we get data from a new
      // playlist. Discontinuities and track changes are handled in other sections.
      this.playlistOfLastInitSegment_[type] = playlist; // we should only be appending the next init segment if we detect a change, or if
      // the segment has a map

      this.appendInitSegment_[type] = map ? true : false; // we need to clear out the fmp4 active init segment id, since
      // we are appending the muxer init segment

      this.activeInitSegmentId_ = null;
      return initSegment;
    }

    return null;
  };

  _proto.appendToSourceBuffer_ = function appendToSourceBuffer_(_ref4) {
    var _this4 = this;

    var segmentInfo = _ref4.segmentInfo,
        type = _ref4.type,
        initSegment = _ref4.initSegment,
        data = _ref4.data;
    var segments = [data];
    var byteLength = data.byteLength;

    if (initSegment) {
      // if the media initialization segment is changing, append it before the content
      // segment
      segments.unshift(initSegment);
      byteLength += initSegment.byteLength;
    } // Technically we should be OK appending the init segment separately, however, we
    // haven't yet tested that, and prepending is how we have always done things.


    var bytes = concatSegments({
      bytes: byteLength,
      segments: segments
    });
    this.sourceUpdater_.appendBuffer({
      segmentInfo: segmentInfo,
      type: type,
      bytes: bytes
    }, function (error) {
      if (error) {
        _this4.error(type + " append of " + bytes.length + "b failed for segment #" + segmentInfo.mediaIndex + " in playlist " + segmentInfo.playlist.id); // If an append errors, we can't recover.
        // (see https://w3c.github.io/media-source/#sourcebuffer-append-error).
        // Trigger a special error so that it can be handled separately from normal,
        // recoverable errors.


        _this4.trigger('appenderror');
      }
    });
  };

  _proto.handleVideoSegmentTimingInfo_ = function handleVideoSegmentTimingInfo_(requestId, videoSegmentTimingInfo) {
    if (!this.pendingSegment_ || requestId !== this.pendingSegment_.requestId) {
      return;
    }

    var segment = this.pendingSegment_.segment;

    if (!segment.videoTimingInfo) {
      segment.videoTimingInfo = {};
    }

    segment.videoTimingInfo.transmuxerPrependedSeconds = videoSegmentTimingInfo.prependedContentDuration || 0;
    segment.videoTimingInfo.transmuxedPresentationStart = videoSegmentTimingInfo.start.presentation;
    segment.videoTimingInfo.transmuxedPresentationEnd = videoSegmentTimingInfo.end.presentation; // mainly used as a reference for debugging

    segment.videoTimingInfo.baseMediaDecodeTime = videoSegmentTimingInfo.baseMediaDecodeTime;
  };

  _proto.appendData_ = function appendData_(segmentInfo, result) {
    var type = result.type,
        data = result.data;

    if (!data || !data.byteLength) {
      return;
    }

    if (type === 'audio' && this.audioDisabled_) {
      return;
    }

    var initSegment = this.getInitSegmentAndUpdateState_({
      type: type,
      initSegment: result.initSegment,
      playlist: segmentInfo.playlist,
      map: segmentInfo.isFmp4 ? segmentInfo.segment.map : null
    });
    this.appendToSourceBuffer_({
      segmentInfo: segmentInfo,
      type: type,
      initSegment: initSegment,
      data: data
    });
  }
  /**
   * load a specific segment from a request into the buffer
   *
   * @private
   */
  ;

  _proto.loadSegment_ = function loadSegment_(segmentInfo) {
    var _this5 = this;

    this.state = 'WAITING';
    this.pendingSegment_ = segmentInfo;
    this.trimBackBuffer_(segmentInfo);

    if (typeof segmentInfo.timestampOffset === 'number') {
      if (this.transmuxer_) {
        this.transmuxer_.postMessage({
          action: 'clearAllMp4Captions'
        });
      }
    }

    if (!this.hasEnoughInfoToLoad_()) {
      this.loadQueue_.push(function () {
        var buffered = _this5.buffered_();

        if (typeof segmentInfo.timestampOffset === 'number') {
          // The timestamp offset needs to be regenerated, as the buffer most likely
          // changed since the function was added to the queue. This is expected, as the
          // load is usually pending the main loader appending new segments.
          //
          // Note also that the overrideCheck property is set to true. This is because
          // isPendingTimestampOffset is set back to false after the first set of the
          // timestamp offset (before it was added to the queue). But the presence of
          // timestamp offset as a property of segmentInfo serves as enough evidence that
          // it should be regenerated.
          segmentInfo.timestampOffset = timestampOffsetForSegment({
            segmentTimeline: segmentInfo.timeline,
            currentTimeline: _this5.currentTimeline_,
            startOfSegment: segmentInfo.startOfSegment,
            buffered: buffered,
            overrideCheck: true
          });
        }

        delete segmentInfo.audioAppendStart;

        var audioBuffered = _this5.sourceUpdater_.audioBuffered();

        if (audioBuffered.length) {
          // Because the audio timestamp offset may have been changed by the main loader,
          // the audioAppendStart should be regenerated.
          //
          // Since the transmuxer is using the actual timing values, but the buffer is
          // adjusted by the timestamp offset, the value must be adjusted.
          segmentInfo.audioAppendStart = audioBuffered.end(audioBuffered.length - 1) - _this5.sourceUpdater_.audioTimestampOffset();
        }

        _this5.updateTransmuxerAndRequestSegment_(segmentInfo);
      });
      return;
    }

    this.updateTransmuxerAndRequestSegment_(segmentInfo);
  };

  _proto.updateTransmuxerAndRequestSegment_ = function updateTransmuxerAndRequestSegment_(segmentInfo) {
    // We'll update the source buffer's timestamp offset once we have transmuxed data, but
    // the transmuxer still needs to be updated before then.
    //
    // Even though keepOriginalTimestamps is set to true for the transmuxer, timestamp
    // offset must be passed to the transmuxer for stream correcting adjustments.
    if (this.shouldUpdateTransmuxerTimestampOffset_(segmentInfo.timestampOffset)) {
      this.gopBuffer_.length = 0; // gopsToAlignWith was set before the GOP buffer was cleared

      segmentInfo.gopsToAlignWith = [];
      this.timeMapping_ = 0; // reset values in the transmuxer since a discontinuity should start fresh

      this.transmuxer_.postMessage({
        action: 'reset'
      });
      this.transmuxer_.postMessage({
        action: 'setTimestampOffset',
        timestampOffset: segmentInfo.timestampOffset
      });
    }

    var simpleSegment = this.createSimplifiedSegmentObj_(segmentInfo);
    segmentInfo.abortRequests = mediaSegmentRequest({
      xhr: this.vhs_.xhr,
      xhrOptions: this.xhrOptions_,
      decryptionWorker: this.decrypter_,
      segment: simpleSegment,
      handlePartialData: this.handlePartialData_,
      abortFn: this.handleAbort_.bind(this),
      progressFn: this.handleProgress_.bind(this),
      trackInfoFn: this.handleTrackInfo_.bind(this),
      timingInfoFn: this.handleTimingInfo_.bind(this),
      videoSegmentTimingInfoFn: this.handleVideoSegmentTimingInfo_.bind(this, segmentInfo.requestId),
      captionsFn: this.handleCaptions_.bind(this),
      id3Fn: this.handleId3_.bind(this),
      dataFn: this.handleData_.bind(this),
      doneFn: this.segmentRequestFinished_.bind(this)
    });
  }
  /**
   * trim the back buffer so that we don't have too much data
   * in the source buffer
   *
   * @private
   *
   * @param {Object} segmentInfo - the current segment
   */
  ;

  _proto.trimBackBuffer_ = function trimBackBuffer_(segmentInfo) {
    var removeToTime = safeBackBufferTrimTime(this.seekable_(), this.currentTime_(), this.playlist_.targetDuration || 10); // Chrome has a hard limit of 150MB of
    // buffer and a very conservative "garbage collector"
    // We manually clear out the old buffer to ensure
    // we don't trigger the QuotaExceeded error
    // on the source buffer during subsequent appends

    if (removeToTime > 0) {
      this.remove(0, removeToTime);
    }
  }
  /**
   * created a simplified copy of the segment object with just the
   * information necessary to perform the XHR and decryption
   *
   * @private
   *
   * @param {Object} segmentInfo - the current segment
   * @return {Object} a simplified segment object copy
   */
  ;

  _proto.createSimplifiedSegmentObj_ = function createSimplifiedSegmentObj_(segmentInfo) {
    var segment = segmentInfo.segment;
    var simpleSegment = {
      resolvedUri: segment.resolvedUri,
      byterange: segment.byterange,
      requestId: segmentInfo.requestId,
      transmuxer: segmentInfo.transmuxer,
      audioAppendStart: segmentInfo.audioAppendStart,
      gopsToAlignWith: segmentInfo.gopsToAlignWith
    };
    var previousSegment = segmentInfo.playlist.segments[segmentInfo.mediaIndex];

    if (previousSegment && previousSegment.end && previousSegment.timeline === segment.timeline) {
      simpleSegment.baseStartTime = previousSegment.end + segmentInfo.timestampOffset;
    }

    if (segment.key) {
      // if the media sequence is greater than 2^32, the IV will be incorrect
      // assuming 10s segments, that would be about 1300 years
      var iv = segment.key.iv || new Uint32Array([0, 0, 0, segmentInfo.mediaIndex + segmentInfo.playlist.mediaSequence]);
      simpleSegment.key = this.segmentKey(segment.key);
      simpleSegment.key.iv = iv;
    }

    if (segment.map) {
      simpleSegment.map = this.initSegmentForMap(segment.map);
    }

    return simpleSegment;
  };

  _proto.saveTransferStats_ = function saveTransferStats_(stats) {
    // every request counts as a media request even if it has been aborted
    // or canceled due to a timeout
    this.mediaRequests += 1;

    if (stats) {
      this.mediaBytesTransferred += stats.bytesReceived;
      this.mediaTransferDuration += stats.roundTripTime;
    }
  };

  _proto.saveBandwidthRelatedStats_ = function saveBandwidthRelatedStats_(stats) {
    this.bandwidth = stats.bandwidth;
    this.roundTrip = stats.roundTripTime; // byteLength will be used for throughput, and should be based on bytes receieved,
    // which we only know at the end of the request and should reflect total bytes
    // downloaded rather than just bytes processed from components of the segment

    this.pendingSegment_.byteLength = stats.bytesReceived;
  };

  _proto.handleTimeout_ = function handleTimeout_() {
    // although the VTT segment loader bandwidth isn't really used, it's good to
    // maintain functinality between segment loaders
    this.mediaRequestsTimedout += 1;
    this.bandwidth = 1;
    this.roundTrip = NaN;
    this.trigger('bandwidthupdate');
  }
  /**
   * Handle the callback from the segmentRequest function and set the
   * associated SegmentLoader state and errors if necessary
   *
   * @private
   */
  ;

  _proto.segmentRequestFinished_ = function segmentRequestFinished_(error, simpleSegment, result) {
    // TODO handle special cases, e.g., muxed audio/video but only audio in the segment
    // check the call queue directly since this function doesn't need to deal with any
    // data, and can continue even if the source buffers are not set up and we didn't get
    // any data from the segment
    if (this.callQueue_.length) {
      this.callQueue_.push(this.segmentRequestFinished_.bind(this, error, simpleSegment, result));
      return;
    }

    this.saveTransferStats_(simpleSegment.stats); // The request was aborted and the SegmentLoader has already been reset

    if (!this.pendingSegment_) {
      return;
    } // the request was aborted and the SegmentLoader has already started
    // another request. this can happen when the timeout for an aborted
    // request triggers due to a limitation in the XHR library
    // do not count this as any sort of request or we risk double-counting


    if (simpleSegment.requestId !== this.pendingSegment_.requestId) {
      return;
    } // an error occurred from the active pendingSegment_ so reset everything


    if (error) {
      this.pendingSegment_ = null;
      this.state = 'READY'; // aborts are not a true error condition and nothing corrective needs to be done

      if (error.code === REQUEST_ERRORS.ABORTED) {
        return;
      }

      this.pause(); // the error is really just that at least one of the requests timed-out
      // set the bandwidth to a very low value and trigger an ABR switch to
      // take emergency action

      if (error.code === REQUEST_ERRORS.TIMEOUT) {
        this.handleTimeout_();
        return;
      } // if control-flow has arrived here, then the error is real
      // emit an error event to blacklist the current playlist


      this.mediaRequestsErrored += 1;
      this.error(error);
      this.trigger('error');
      return;
    } // the response was a success so set any bandwidth stats the request
    // generated for ABR purposes


    this.saveBandwidthRelatedStats_(simpleSegment.stats);
    var segmentInfo = this.pendingSegment_;
    segmentInfo.endOfAllRequests = simpleSegment.endOfAllRequests;

    if (result.gopInfo) {
      this.gopBuffer_ = updateGopBuffer(this.gopBuffer_, result.gopInfo, this.safeAppend_);
    } // Although we may have already started appending on progress, we shouldn't switch the
    // state away from loading until we are officially done loading the segment data.


    this.state = 'APPENDING';
    var isEndOfStream = this.isEndOfStream_(segmentInfo.mediaIndex, segmentInfo.playlist);
    var isWalkingForward = this.mediaIndex !== null;
    var isDiscontinuity = segmentInfo.timeline !== this.currentTimeline_ && // TODO verify this behavior
    // currentTimeline starts at -1, but we shouldn't end the timeline switching to 0,
    // the first timeline
    segmentInfo.timeline > 0;

    if (!segmentInfo.isFmp4 && (isEndOfStream || isWalkingForward && isDiscontinuity)) {
      segmentTransmuxer.endTimeline(this.transmuxer_);
    } // used for testing


    this.trigger('appending');
    this.waitForAppendsToComplete_(segmentInfo);
  };

  _proto.setTimeMapping_ = function setTimeMapping_(timeline) {
    var timelineMapping = this.syncController_.mappingForTimeline(timeline);

    if (timelineMapping !== null) {
      this.timeMapping_ = timelineMapping;
    }
  };

  _proto.updateMediaSecondsLoaded_ = function updateMediaSecondsLoaded_(segment) {
    if (typeof segment.start === 'number' && typeof segment.end === 'number') {
      this.mediaSecondsLoaded += segment.end - segment.start;
    } else {
      this.mediaSecondsLoaded += segment.duration;
    }
  };

  _proto.shouldUpdateTransmuxerTimestampOffset_ = function shouldUpdateTransmuxerTimestampOffset_(timestampOffset) {
    if (timestampOffset === null) {
      return false;
    } // note that we're potentially using the same timestamp offset for both video and
    // audio


    if (this.loaderType_ === 'main' && timestampOffset !== this.sourceUpdater_.videoTimestampOffset()) {
      return true;
    }

    if (!this.audioDisabled_ && timestampOffset !== this.sourceUpdater_.audioTimestampOffset()) {
      return true;
    }

    return false;
  };

  _proto.trueSegmentStart_ = function trueSegmentStart_(_ref5) {
    var currentStart = _ref5.currentStart,
        playlist = _ref5.playlist,
        mediaIndex = _ref5.mediaIndex,
        firstVideoFrameTimeForData = _ref5.firstVideoFrameTimeForData,
        currentVideoTimestampOffset = _ref5.currentVideoTimestampOffset,
        useVideoTimingInfo = _ref5.useVideoTimingInfo,
        videoTimingInfo = _ref5.videoTimingInfo,
        audioTimingInfo = _ref5.audioTimingInfo;

    if (typeof currentStart !== 'undefined') {
      // if start was set once, keep using it
      return currentStart;
    }

    if (!useVideoTimingInfo) {
      return audioTimingInfo.start;
    }

    var previousSegment = playlist.segments[mediaIndex - 1]; // The start of a segment should be the start of the first full frame contained
    // within that segment. Since the transmuxer maintains a cache of incomplete data
    // from and/or the last frame seen, the start time may reflect a frame that starts
    // in the previous segment. Check for that case and ensure the start time is
    // accurate for the segment.

    if (mediaIndex === 0 || !previousSegment || typeof previousSegment.start === 'undefined' || previousSegment.end !== firstVideoFrameTimeForData + currentVideoTimestampOffset) {
      return firstVideoFrameTimeForData;
    }

    return videoTimingInfo.start;
  };

  _proto.waitForAppendsToComplete_ = function waitForAppendsToComplete_(segmentInfo) {
    if (!this.currentMediaInfo_) {
      this.error({
        message: 'No starting media returned, likely due to an unsupported media format.',
        blacklistDuration: Infinity
      });
      this.trigger('error');
      return;
    } // Although transmuxing is done, appends may not yet be finished. Throw a marker
    // on each queue this loader is responsible for to ensure that the appends are
    // complete.


    var _this$currentMediaInf2 = this.currentMediaInfo_,
        hasAudio = _this$currentMediaInf2.hasAudio,
        hasVideo = _this$currentMediaInf2.hasVideo,
        isMuxed = _this$currentMediaInf2.isMuxed;
    var waitForVideo = this.loaderType_ === 'main' && hasVideo; // TODO: does this break partial support for muxed content?

    var waitForAudio = !this.audioDisabled_ && hasAudio && !isMuxed;
    segmentInfo.waitingOnAppends = 0; // segments with no data

    if (!segmentInfo.hasAppendedData_) {
      if (!segmentInfo.timingInfo && typeof segmentInfo.timestampOffset === 'number') {
        // When there's no audio or video data in the segment, there's no audio or video
        // timing information.
        //
        // If there's no audio or video timing information, then the timestamp offset
        // can't be adjusted to the appropriate value for the transmuxer and source
        // buffers.
        //
        // Therefore, the next segment should be used to set the timestamp offset.
        this.isPendingTimestampOffset_ = true;
      } // override settings for metadata only segments


      segmentInfo.timingInfo = {
        start: 0
      };
      segmentInfo.waitingOnAppends++;

      if (!this.isPendingTimestampOffset_) {
        // update the timestampoffset
        this.updateSourceBufferTimestampOffset_(segmentInfo); // make sure the metadata queue is processed even though we have
        // no video/audio data.

        this.processMetadataQueue_();
      } // append is "done" instantly with no data.


      this.checkAppendsDone_(segmentInfo);
      return;
    } // Since source updater could call back synchronously, do the increments first.


    if (waitForVideo) {
      segmentInfo.waitingOnAppends++;
    }

    if (waitForAudio) {
      segmentInfo.waitingOnAppends++;
    }

    if (waitForVideo) {
      this.sourceUpdater_.videoQueueCallback(this.checkAppendsDone_.bind(this, segmentInfo));
    }

    if (waitForAudio) {
      this.sourceUpdater_.audioQueueCallback(this.checkAppendsDone_.bind(this, segmentInfo));
    }
  };

  _proto.checkAppendsDone_ = function checkAppendsDone_(segmentInfo) {
    if (this.checkForAbort_(segmentInfo.requestId)) {
      return;
    }

    segmentInfo.waitingOnAppends--;

    if (segmentInfo.waitingOnAppends === 0) {
      this.handleAppendsDone_();
    }
  };

  _proto.checkForIllegalMediaSwitch = function checkForIllegalMediaSwitch(trackInfo) {
    var illegalMediaSwitchError = illegalMediaSwitch(this.loaderType_, this.currentMediaInfo_, trackInfo);

    if (illegalMediaSwitchError) {
      this.error({
        message: illegalMediaSwitchError,
        blacklistDuration: Infinity
      });
      this.trigger('error');
      return true;
    }

    return false;
  };

  _proto.updateSourceBufferTimestampOffset_ = function updateSourceBufferTimestampOffset_(segmentInfo) {
    if (segmentInfo.timestampOffset === null || // we don't yet have the start for whatever media type (video or audio) has
    // priority, timing-wise, so we must wait
    typeof segmentInfo.timingInfo.start !== 'number' || // already updated the timestamp offset for this segment
    segmentInfo.changedTimestampOffset || // the alt audio loader should not be responsible for setting the timestamp offset
    this.loaderType_ !== 'main') {
      return;
    }

    var didChange = false; // Primary timing goes by video, and audio is trimmed in the transmuxer, meaning that
    // the timing info here comes from video. In the event that the audio is longer than
    // the video, this will trim the start of the audio.
    // This also trims any offset from 0 at the beginning of the media

    segmentInfo.timestampOffset -= segmentInfo.timingInfo.start; // In the event that there are partial segment downloads, each will try to update the
    // timestamp offset. Retaining this bit of state prevents us from updating in the
    // future (within the same segment), however, there may be a better way to handle it.

    segmentInfo.changedTimestampOffset = true;

    if (segmentInfo.timestampOffset !== this.sourceUpdater_.videoTimestampOffset()) {
      this.sourceUpdater_.videoTimestampOffset(segmentInfo.timestampOffset);
      didChange = true;
    }

    if (segmentInfo.timestampOffset !== this.sourceUpdater_.audioTimestampOffset()) {
      this.sourceUpdater_.audioTimestampOffset(segmentInfo.timestampOffset);
      didChange = true;
    }

    if (didChange) {
      this.trigger('timestampoffset');
    }
  };

  _proto.updateTimingInfoEnd_ = function updateTimingInfoEnd_(segmentInfo) {
    segmentInfo.timingInfo = segmentInfo.timingInfo || {};
    var useVideoTimingInfo = this.loaderType_ === 'main' && this.currentMediaInfo_.hasVideo;
    var prioritizedTimingInfo = useVideoTimingInfo && segmentInfo.videoTimingInfo ? segmentInfo.videoTimingInfo : segmentInfo.audioTimingInfo;

    if (!prioritizedTimingInfo) {
      return;
    }

    segmentInfo.timingInfo.end = typeof prioritizedTimingInfo.end === 'number' ? // End time may not exist in a case where we aren't parsing the full segment (one
    // current example is the case of fmp4), so use the rough duration to calculate an
    // end time.
    prioritizedTimingInfo.end : prioritizedTimingInfo.start + segmentInfo.duration;
  }
  /**
   * callback to run when appendBuffer is finished. detects if we are
   * in a good state to do things with the data we got, or if we need
   * to wait for more
   *
   * @private
   */
  ;

  _proto.handleAppendsDone_ = function handleAppendsDone_() {
    if (!this.pendingSegment_) {
      this.state = 'READY'; // TODO should this move into this.checkForAbort to speed up requests post abort in
      // all appending cases?

      if (!this.paused()) {
        this.monitorBuffer_();
      }

      return;
    }

    this.trigger('appendsdone');
    var segmentInfo = this.pendingSegment_; // Now that the end of the segment has been reached, we can set the end time. It's
    // best to wait until all appends are done so we're sure that the primary media is
    // finished (and we have its end time).

    this.updateTimingInfoEnd_(segmentInfo);

    if (this.shouldSaveSegmentTimingInfo_) {
      // Timeline mappings should only be saved for the main loader. This is for multiple
      // reasons:
      //
      // 1) Only one mapping is saved per timeline, meaning that if both the audio loader
      //    and the main loader try to save the timeline mapping, whichever comes later
      //    will overwrite the first. In theory this is OK, as the mappings should be the
      //    same, however, it breaks for (2)
      // 2) In the event of a live stream, the initial live point will make for a somewhat
      //    arbitrary mapping. If audio and video streams are not perfectly in-sync, then
      //    the mapping will be off for one of the streams, dependent on which one was
      //    first saved (see (1)).
      // 3) Primary timing goes by video in VHS, so the mapping should be video.
      //
      // Since the audio loader will wait for the main loader to load the first segment,
      // the main loader will save the first timeline mapping, and ensure that there won't
      // be a case where audio loads two segments without saving a mapping (thus leading
      // to missing segment timing info).
      this.syncController_.saveSegmentTimingInfo({
        segmentInfo: segmentInfo,
        shouldSaveTimelineMapping: this.loaderType_ === 'main'
      });
    }

    this.logger_(segmentInfoString(segmentInfo));
    this.recordThroughput_(segmentInfo);
    this.pendingSegment_ = null;
    this.state = 'READY'; // TODO minor, but for partial segment downloads, this can be done earlier to save
    // on bandwidth and download time

    if (segmentInfo.isSyncRequest) {
      this.trigger('syncinfoupdate');
      return;
    }

    this.addSegmentMetadataCue_(segmentInfo);
    this.fetchAtBuffer_ = true;

    if (this.currentTimeline_ !== segmentInfo.timeline) {
      this.timelineChangeController_.lastTimelineChange({
        type: this.loaderType_,
        from: this.currentTimeline_,
        to: segmentInfo.timeline
      }); // If audio is not disabled, the main segment loader is responsible for updating
      // the audio timeline as well. If the content is video only, this won't have any
      // impact.

      if (this.loaderType_ === 'main' && !this.audioDisabled_) {
        this.timelineChangeController_.lastTimelineChange({
          type: 'audio',
          from: this.currentTimeline_,
          to: segmentInfo.timeline
        });
      }
    }

    this.currentTimeline_ = segmentInfo.timeline; // We must update the syncinfo to recalculate the seekable range before
    // the following conditional otherwise it may consider this a bad "guess"
    // and attempt to resync when the post-update seekable window and live
    // point would mean that this was the perfect segment to fetch

    this.trigger('syncinfoupdate');
    var segment = segmentInfo.segment; // If we previously appended a segment that ends more than 3 targetDurations before
    // the currentTime_ that means that our conservative guess was too conservative.
    // In that case, reset the loader state so that we try to use any information gained
    // from the previous request to create a new, more accurate, sync-point.

    if (segment.end && this.currentTime_() - segment.end > segmentInfo.playlist.targetDuration * 3) {
      this.resetEverything();
      return;
    }

    var isWalkingForward = this.mediaIndex !== null; // Don't do a rendition switch unless we have enough time to get a sync segment
    // and conservatively guess

    if (isWalkingForward) {
      this.trigger('bandwidthupdate');
    }

    this.trigger('progress');
    this.mediaIndex = segmentInfo.mediaIndex; // any time an update finishes and the last segment is in the
    // buffer, end the stream. this ensures the "ended" event will
    // fire if playback reaches that point.

    if (this.isEndOfStream_(segmentInfo.mediaIndex, segmentInfo.playlist)) {
      this.endOfStream();
    } // used for testing


    this.trigger('appended');

    if (!this.paused()) {
      this.monitorBuffer_();
    }
  }
  /**
   * Records the current throughput of the decrypt, transmux, and append
   * portion of the semgment pipeline. `throughput.rate` is a the cumulative
   * moving average of the throughput. `throughput.count` is the number of
   * data points in the average.
   *
   * @private
   * @param {Object} segmentInfo the object returned by loadSegment
   */
  ;

  _proto.recordThroughput_ = function recordThroughput_(segmentInfo) {
    var rate = this.throughput.rate; // Add one to the time to ensure that we don't accidentally attempt to divide
    // by zero in the case where the throughput is ridiculously high

    var segmentProcessingTime = Date.now() - segmentInfo.endOfAllRequests + 1; // Multiply by 8000 to convert from bytes/millisecond to bits/second

    var segmentProcessingThroughput = Math.floor(segmentInfo.byteLength / segmentProcessingTime * 8 * 1000); // This is just a cumulative moving average calculation:
    //   newAvg = oldAvg + (sample - oldAvg) / (sampleCount + 1)

    this.throughput.rate += (segmentProcessingThroughput - rate) / ++this.throughput.count;
  }
  /**
   * Adds a cue to the segment-metadata track with some metadata information about the
   * segment
   *
   * @private
   * @param {Object} segmentInfo
   *        the object returned by loadSegment
   * @method addSegmentMetadataCue_
   */
  ;

  _proto.addSegmentMetadataCue_ = function addSegmentMetadataCue_(segmentInfo) {
    if (!this.segmentMetadataTrack_) {
      return;
    }

    var segment = segmentInfo.segment;
    var start = segment.start;
    var end = segment.end; // Do not try adding the cue if the start and end times are invalid.

    if (!finite(start) || !finite(end)) {
      return;
    }

    removeCuesFromTrack(start, end, this.segmentMetadataTrack_);
    var Cue = window_1.WebKitDataCue || window_1.VTTCue;
    var value = {
      custom: segment.custom,
      dateTimeObject: segment.dateTimeObject,
      dateTimeString: segment.dateTimeString,
      bandwidth: segmentInfo.playlist.attributes.BANDWIDTH,
      resolution: segmentInfo.playlist.attributes.RESOLUTION,
      codecs: segmentInfo.playlist.attributes.CODECS,
      byteLength: segmentInfo.byteLength,
      uri: segmentInfo.uri,
      timeline: segmentInfo.timeline,
      playlist: segmentInfo.playlist.id,
      start: start,
      end: end
    };
    var data = JSON.stringify(value);
    var cue = new Cue(start, end, data); // Attach the metadata to the value property of the cue to keep consistency between
    // the differences of WebKitDataCue in safari and VTTCue in other browsers

    cue.value = value;
    this.segmentMetadataTrack_.addCue(cue);
  };

  return SegmentLoader;
}(videojs$1.EventTarget);

function noop() {}

var toTitleCase = function toTitleCase(string) {
  if (typeof string !== 'string') {
    return string;
  }

  return string.replace(/./, function (w) {
    return w.toUpperCase();
  });
};

var bufferTypes = ['video', 'audio'];

var _updating = function updating(type, sourceUpdater) {
  var sourceBuffer = sourceUpdater[type + "Buffer"];
  return sourceBuffer && sourceBuffer.updating || sourceUpdater.queuePending[type];
};

var nextQueueIndexOfType = function nextQueueIndexOfType(type, queue) {
  for (var i = 0; i < queue.length; i++) {
    var queueEntry = queue[i];

    if (queueEntry.type === 'mediaSource') {
      // If the next entry is a media source entry (uses multiple source buffers), block
      // processing to allow it to go through first.
      return null;
    }

    if (queueEntry.type === type) {
      return i;
    }
  }

  return null;
};

var shiftQueue = function shiftQueue(type, sourceUpdater) {
  if (sourceUpdater.queue.length === 0) {
    return;
  }

  var queueIndex = 0;
  var queueEntry = sourceUpdater.queue[queueIndex];

  if (queueEntry.type === 'mediaSource') {
    if (!sourceUpdater.updating() && sourceUpdater.mediaSource.readyState !== 'closed') {
      sourceUpdater.queue.shift();
      queueEntry.action(sourceUpdater);

      if (queueEntry.doneFn) {
        queueEntry.doneFn();
      } // Only specific source buffer actions must wait for async updateend events. Media
      // Source actions process synchronously. Therefore, both audio and video source
      // buffers are now clear to process the next queue entries.


      shiftQueue('audio', sourceUpdater);
      shiftQueue('video', sourceUpdater);
    } // Media Source actions require both source buffers, so if the media source action
    // couldn't process yet (because one or both source buffers are busy), block other
    // queue actions until both are available and the media source action can process.


    return;
  }

  if (type === 'mediaSource') {
    // If the queue was shifted by a media source action (this happens when pushing a
    // media source action onto the queue), then it wasn't from an updateend event from an
    // audio or video source buffer, so there's no change from previous state, and no
    // processing should be done.
    return;
  } // Media source queue entries don't need to consider whether the source updater is
  // started (i.e., source buffers are created) as they don't need the source buffers, but
  // source buffer queue entries do.


  if (!sourceUpdater.started_ || sourceUpdater.mediaSource.readyState === 'closed' || _updating(type, sourceUpdater)) {
    return;
  }

  if (queueEntry.type !== type) {
    queueIndex = nextQueueIndexOfType(type, sourceUpdater.queue);

    if (queueIndex === null) {
      // Either there's no queue entry that uses this source buffer type in the queue, or
      // there's a media source queue entry before the next entry of this type, in which
      // case wait for that action to process first.
      return;
    }

    queueEntry = sourceUpdater.queue[queueIndex];
  }

  sourceUpdater.queue.splice(queueIndex, 1);
  queueEntry.action(type, sourceUpdater);

  if (!queueEntry.doneFn) {
    // synchronous operation, process next entry
    shiftQueue(type, sourceUpdater);
    return;
  } // asynchronous operation, so keep a record that this source buffer type is in use


  sourceUpdater.queuePending[type] = queueEntry;
};

var cleanupBuffer = function cleanupBuffer(type, sourceUpdater) {
  var buffer = sourceUpdater[type + "Buffer"];
  var titleType = toTitleCase(type);

  if (!buffer) {
    return;
  }

  buffer.removeEventListener('updateend', sourceUpdater["on" + titleType + "UpdateEnd_"]);
  buffer.removeEventListener('error', sourceUpdater["on" + titleType + "Error_"]);
  sourceUpdater.codecs[type] = null;
  sourceUpdater[type + "Buffer"] = null;
};

var inSourceBuffers = function inSourceBuffers(mediaSource, sourceBuffer) {
  return mediaSource && sourceBuffer && Array.prototype.indexOf.call(mediaSource.sourceBuffers, sourceBuffer) !== -1;
};

var actions = {
  appendBuffer: function appendBuffer(bytes, segmentInfo) {
    return function (type, sourceUpdater) {
      var sourceBuffer = sourceUpdater[type + "Buffer"]; // can't do anything if the media source / source buffer is null
      // or the media source does not contain this source buffer.

      if (!inSourceBuffers(sourceUpdater.mediaSource, sourceBuffer)) {
        return;
      }

      sourceUpdater.logger_("Appending segment " + segmentInfo.mediaIndex + "'s " + bytes.length + " bytes to " + type + "Buffer");
      sourceBuffer.appendBuffer(bytes);
    };
  },
  remove: function remove(start, end) {
    return function (type, sourceUpdater) {
      var sourceBuffer = sourceUpdater[type + "Buffer"]; // can't do anything if the media source / source buffer is null
      // or the media source does not contain this source buffer.

      if (!inSourceBuffers(sourceUpdater.mediaSource, sourceBuffer)) {
        return;
      }

      sourceUpdater.logger_("Removing " + start + " to " + end + " from " + type + "Buffer");
      sourceBuffer.remove(start, end);
    };
  },
  timestampOffset: function timestampOffset(offset) {
    return function (type, sourceUpdater) {
      var sourceBuffer = sourceUpdater[type + "Buffer"]; // can't do anything if the media source / source buffer is null
      // or the media source does not contain this source buffer.

      if (!inSourceBuffers(sourceUpdater.mediaSource, sourceBuffer)) {
        return;
      }

      sourceUpdater.logger_("Setting " + type + "timestampOffset to " + offset);
      sourceBuffer.timestampOffset = offset;
    };
  },
  callback: function callback(_callback) {
    return function (type, sourceUpdater) {
      _callback();
    };
  },
  endOfStream: function endOfStream(error) {
    return function (sourceUpdater) {
      if (sourceUpdater.mediaSource.readyState !== 'open') {
        return;
      }

      sourceUpdater.logger_("Calling mediaSource endOfStream(" + (error || '') + ")");

      try {
        sourceUpdater.mediaSource.endOfStream(error);
      } catch (e) {
        videojs$1.log.warn('Failed to call media source endOfStream', e);
      }
    };
  },
  duration: function duration(_duration) {
    return function (sourceUpdater) {
      sourceUpdater.logger_("Setting mediaSource duration to " + _duration);

      try {
        sourceUpdater.mediaSource.duration = _duration;
      } catch (e) {
        videojs$1.log.warn('Failed to set media source duration', e);
      }
    };
  },
  abort: function abort() {
    return function (type, sourceUpdater) {
      if (sourceUpdater.mediaSource.readyState !== 'open') {
        return;
      }

      var sourceBuffer = sourceUpdater[type + "Buffer"]; // can't do anything if the media source / source buffer is null
      // or the media source does not contain this source buffer.

      if (!inSourceBuffers(sourceUpdater.mediaSource, sourceBuffer)) {
        return;
      }

      sourceUpdater.logger_("calling abort on " + type + "Buffer");

      try {
        sourceBuffer.abort();
      } catch (e) {
        videojs$1.log.warn("Failed to abort on " + type + "Buffer", e);
      }
    };
  },
  addSourceBuffer: function addSourceBuffer(type, codec) {
    return function (sourceUpdater) {
      var titleType = toTitleCase(type);
      var mime = codecs.getMimeForCodec(codec);
      sourceUpdater.logger_("Adding " + type + "Buffer with codec " + codec + " to mediaSource");
      var sourceBuffer = sourceUpdater.mediaSource.addSourceBuffer(mime);
      sourceBuffer.addEventListener('updateend', sourceUpdater["on" + titleType + "UpdateEnd_"]);
      sourceBuffer.addEventListener('error', sourceUpdater["on" + titleType + "Error_"]);
      sourceUpdater.codecs[type] = codec;
      sourceUpdater[type + "Buffer"] = sourceBuffer;
    };
  },
  removeSourceBuffer: function removeSourceBuffer(type) {
    return function (sourceUpdater) {
      var sourceBuffer = sourceUpdater[type + "Buffer"];
      cleanupBuffer(type, sourceUpdater); // can't do anything if the media source / source buffer is null
      // or the media source does not contain this source buffer.

      if (!inSourceBuffers(sourceUpdater.mediaSource, sourceBuffer)) {
        return;
      }

      sourceUpdater.logger_("Removing " + type + "Buffer with codec " + sourceUpdater.codecs[type] + " from mediaSource");

      try {
        sourceUpdater.mediaSource.removeSourceBuffer(sourceBuffer);
      } catch (e) {
        videojs$1.log.warn("Failed to removeSourceBuffer " + type + "Buffer", e);
      }
    };
  },
  changeType: function changeType(codec) {
    return function (type, sourceUpdater) {
      var sourceBuffer = sourceUpdater[type + "Buffer"];
      var mime = codecs.getMimeForCodec(codec); // can't do anything if the media source / source buffer is null
      // or the media source does not contain this source buffer.

      if (!inSourceBuffers(sourceUpdater.mediaSource, sourceBuffer)) {
        return;
      } // do not update codec if we don't need to.


      if (sourceUpdater.codecs[type] === codec) {
        return;
      }

      sourceUpdater.logger_("changing " + type + "Buffer codec from " + sourceUpdater.codecs[type] + " to " + codec);
      sourceBuffer.changeType(mime);
      sourceUpdater.codecs[type] = codec;
    };
  }
};

var pushQueue = function pushQueue(_ref) {
  var type = _ref.type,
      sourceUpdater = _ref.sourceUpdater,
      action = _ref.action,
      doneFn = _ref.doneFn,
      name = _ref.name;
  sourceUpdater.queue.push({
    type: type,
    action: action,
    doneFn: doneFn,
    name: name
  });
  shiftQueue(type, sourceUpdater);
};

var onUpdateend = function onUpdateend(type, sourceUpdater) {
  return function (e) {
    // Although there should, in theory, be a pending action for any updateend receieved,
    // there are some actions that may trigger updateend events without set definitions in
    // the w3c spec. For instance, setting the duration on the media source may trigger
    // updateend events on source buffers. This does not appear to be in the spec. As such,
    // if we encounter an updateend without a corresponding pending action from our queue
    // for that source buffer type, process the next action.
    if (sourceUpdater.queuePending[type]) {
      var doneFn = sourceUpdater.queuePending[type].doneFn;
      sourceUpdater.queuePending[type] = null;

      if (doneFn) {
        // if there's an error, report it
        doneFn(sourceUpdater[type + "Error_"]);
      }
    }

    shiftQueue(type, sourceUpdater);
  };
};
/**
 * A queue of callbacks to be serialized and applied when a
 * MediaSource and its associated SourceBuffers are not in the
 * updating state. It is used by the segment loader to update the
 * underlying SourceBuffers when new data is loaded, for instance.
 *
 * @class SourceUpdater
 * @param {MediaSource} mediaSource the MediaSource to create the SourceBuffer from
 * @param {string} mimeType the desired MIME type of the underlying SourceBuffer
 */


var SourceUpdater = /*#__PURE__*/function (_videojs$EventTarget) {
  inheritsLoose(SourceUpdater, _videojs$EventTarget);

  function SourceUpdater(mediaSource) {
    var _this;

    _this = _videojs$EventTarget.call(this) || this;
    _this.mediaSource = mediaSource;

    _this.sourceopenListener_ = function () {
      return shiftQueue('mediaSource', assertThisInitialized(_this));
    };

    _this.mediaSource.addEventListener('sourceopen', _this.sourceopenListener_);

    _this.logger_ = logger('SourceUpdater'); // initial timestamp offset is 0

    _this.audioTimestampOffset_ = 0;
    _this.videoTimestampOffset_ = 0;
    _this.queue = [];
    _this.queuePending = {
      audio: null,
      video: null
    };
    _this.delayedAudioAppendQueue_ = [];
    _this.videoAppendQueued_ = false;
    _this.codecs = {};
    _this.onVideoUpdateEnd_ = onUpdateend('video', assertThisInitialized(_this));
    _this.onAudioUpdateEnd_ = onUpdateend('audio', assertThisInitialized(_this));

    _this.onVideoError_ = function (e) {
      // used for debugging
      _this.videoError_ = e;
    };

    _this.onAudioError_ = function (e) {
      // used for debugging
      _this.audioError_ = e;
    };

    _this.started_ = false;
    return _this;
  }

  var _proto = SourceUpdater.prototype;

  _proto.ready = function ready() {
    return this.started_;
  };

  _proto.createSourceBuffers = function createSourceBuffers(codecs) {
    if (this.ready()) {
      // already created them before
      return;
    } // the intial addOrChangeSourceBuffers will always be
    // two add buffers.


    this.addOrChangeSourceBuffers(codecs);
    this.started_ = true;
    this.trigger('ready');
  }
  /**
   * Add a type of source buffer to the media source.
   *
   * @param {string} type
   *        The type of source buffer to add.
   *
   * @param {string} codec
   *        The codec to add the source buffer with.
   */
  ;

  _proto.addSourceBuffer = function addSourceBuffer(type, codec) {
    pushQueue({
      type: 'mediaSource',
      sourceUpdater: this,
      action: actions.addSourceBuffer(type, codec),
      name: 'addSourceBuffer'
    });
  }
  /**
   * call abort on a source buffer.
   *
   * @param {string} type
   *        The type of source buffer to call abort on.
   */
  ;

  _proto.abort = function abort(type) {
    pushQueue({
      type: type,
      sourceUpdater: this,
      action: actions.abort(type),
      name: 'abort'
    });
  }
  /**
   * Call removeSourceBuffer and remove a specific type
   * of source buffer on the mediaSource.
   *
   * @param {string} type
   *        The type of source buffer to remove.
   */
  ;

  _proto.removeSourceBuffer = function removeSourceBuffer(type) {
    if (!this.canRemoveSourceBuffer()) {
      videojs$1.log.error('removeSourceBuffer is not supported!');
      return;
    }

    pushQueue({
      type: 'mediaSource',
      sourceUpdater: this,
      action: actions.removeSourceBuffer(type),
      name: 'removeSourceBuffer'
    });
  }
  /**
   * Whether or not the removeSourceBuffer function is supported
   * on the mediaSource.
   *
   * @return {boolean}
   *          if removeSourceBuffer can be called.
   */
  ;

  _proto.canRemoveSourceBuffer = function canRemoveSourceBuffer() {
    // IE reports that it supports removeSourceBuffer, but often throws
    // errors when attempting to use the function. So we report that it
    // does not support removeSourceBuffer.
    return !videojs$1.browser.IE_VERSION && window_1.MediaSource && window_1.MediaSource.prototype && typeof window_1.MediaSource.prototype.removeSourceBuffer === 'function';
  }
  /**
   * Whether or not the changeType function is supported
   * on our SourceBuffers.
   *
   * @return {boolean}
   *         if changeType can be called.
   */
  ;

  SourceUpdater.canChangeType = function canChangeType() {
    return window_1.SourceBuffer && window_1.SourceBuffer.prototype && typeof window_1.SourceBuffer.prototype.changeType === 'function';
  }
  /**
   * Whether or not the changeType function is supported
   * on our SourceBuffers.
   *
   * @return {boolean}
   *         if changeType can be called.
   */
  ;

  _proto.canChangeType = function canChangeType() {
    return this.constructor.canChangeType();
  }
  /**
   * Call the changeType function on a source buffer, given the code and type.
   *
   * @param {string} type
   *        The type of source buffer to call changeType on.
   *
   * @param {string} codec
   *        The codec string to change type with on the source buffer.
   */
  ;

  _proto.changeType = function changeType(type, codec) {
    if (!this.canChangeType()) {
      videojs$1.log.error('changeType is not supported!');
      return;
    }

    pushQueue({
      type: type,
      sourceUpdater: this,
      action: actions.changeType(codec),
      name: 'changeType'
    });
  }
  /**
   * Add source buffers with a codec or, if they are already created,
   * call changeType on source buffers using changeType.
   *
   * @param {Object} codecs
   *        Codecs to switch to
   */
  ;

  _proto.addOrChangeSourceBuffers = function addOrChangeSourceBuffers(codecs) {
    var _this2 = this;

    if (!codecs || typeof codecs !== 'object' || Object.keys(codecs).length === 0) {
      throw new Error('Cannot addOrChangeSourceBuffers to undefined codecs');
    }

    Object.keys(codecs).forEach(function (type) {
      var codec = codecs[type];

      if (!_this2.ready()) {
        return _this2.addSourceBuffer(type, codec);
      }

      if (_this2.canChangeType()) {
        _this2.changeType(type, codec);
      }
    });
  }
  /**
   * Queue an update to append an ArrayBuffer.
   *
   * @param {MediaObject} object containing audioBytes and/or videoBytes
   * @param {Function} done the function to call when done
   * @see http://www.w3.org/TR/media-source/#widl-SourceBuffer-appendBuffer-void-ArrayBuffer-data
   */
  ;

  _proto.appendBuffer = function appendBuffer(options, doneFn) {
    var _this3 = this;

    var segmentInfo = options.segmentInfo,
        type = options.type,
        bytes = options.bytes;
    this.processedAppend_ = true;

    if (type === 'audio' && this.videoBuffer && !this.videoAppendQueued_) {
      this.delayedAudioAppendQueue_.push([options, doneFn]);
      this.logger_("delayed audio append of " + bytes.length + " until video append");
      return;
    }

    pushQueue({
      type: type,
      sourceUpdater: this,
      action: actions.appendBuffer(bytes, segmentInfo || {
        mediaIndex: -1
      }),
      doneFn: doneFn,
      name: 'appendBuffer'
    });

    if (type === 'video') {
      this.videoAppendQueued_ = true;

      if (!this.delayedAudioAppendQueue_.length) {
        return;
      }

      var queue = this.delayedAudioAppendQueue_.slice();
      this.logger_("queuing delayed audio " + queue.length + " appendBuffers");
      this.delayedAudioAppendQueue_.length = 0;
      queue.forEach(function (que) {
        _this3.appendBuffer.apply(_this3, que);
      });
    }
  }
  /**
   * Get the audio buffer's buffered timerange.
   *
   * @return {TimeRange}
   *         The audio buffer's buffered time range
   */
  ;

  _proto.audioBuffered = function audioBuffered() {
    // no media source/source buffer or it isn't in the media sources
    // source buffer list
    if (!inSourceBuffers(this.mediaSource, this.audioBuffer)) {
      return videojs$1.createTimeRange();
    }

    return this.audioBuffer.buffered ? this.audioBuffer.buffered : videojs$1.createTimeRange();
  }
  /**
   * Get the video buffer's buffered timerange.
   *
   * @return {TimeRange}
   *         The video buffer's buffered time range
   */
  ;

  _proto.videoBuffered = function videoBuffered() {
    // no media source/source buffer or it isn't in the media sources
    // source buffer list
    if (!inSourceBuffers(this.mediaSource, this.videoBuffer)) {
      return videojs$1.createTimeRange();
    }

    return this.videoBuffer.buffered ? this.videoBuffer.buffered : videojs$1.createTimeRange();
  }
  /**
   * Get a combined video/audio buffer's buffered timerange.
   *
   * @return {TimeRange}
   *         the combined time range
   */
  ;

  _proto.buffered = function buffered() {
    var video = inSourceBuffers(this.mediaSource, this.videoBuffer) ? this.videoBuffer : null;
    var audio = inSourceBuffers(this.mediaSource, this.audioBuffer) ? this.audioBuffer : null;

    if (audio && !video) {
      return this.audioBuffered();
    }

    if (video && !audio) {
      return this.videoBuffered();
    }

    return bufferIntersection(this.audioBuffered(), this.videoBuffered());
  }
  /**
   * Add a callback to the queue that will set duration on the mediaSource.
   *
   * @param {number} duration
   *        The duration to set
   *
   * @param {Function} [doneFn]
   *        function to run after duration has been set.
   */
  ;

  _proto.setDuration = function setDuration(duration, doneFn) {
    if (doneFn === void 0) {
      doneFn = noop;
    }

    // In order to set the duration on the media source, it's necessary to wait for all
    // source buffers to no longer be updating. "If the updating attribute equals true on
    // any SourceBuffer in sourceBuffers, then throw an InvalidStateError exception and
    // abort these steps." (source: https://www.w3.org/TR/media-source/#attributes).
    pushQueue({
      type: 'mediaSource',
      sourceUpdater: this,
      action: actions.duration(duration),
      name: 'duration',
      doneFn: doneFn
    });
  }
  /**
   * Add a mediaSource endOfStream call to the queue
   *
   * @param {Error} [error]
   *        Call endOfStream with an error
   *
   * @param {Function} [doneFn]
   *        A function that should be called when the
   *        endOfStream call has finished.
   */
  ;

  _proto.endOfStream = function endOfStream(error, doneFn) {
    if (error === void 0) {
      error = null;
    }

    if (doneFn === void 0) {
      doneFn = noop;
    }

    if (typeof error !== 'string') {
      error = undefined;
    } // In order to set the duration on the media source, it's necessary to wait for all
    // source buffers to no longer be updating. "If the updating attribute equals true on
    // any SourceBuffer in sourceBuffers, then throw an InvalidStateError exception and
    // abort these steps." (source: https://www.w3.org/TR/media-source/#attributes).


    pushQueue({
      type: 'mediaSource',
      sourceUpdater: this,
      action: actions.endOfStream(error),
      name: 'endOfStream',
      doneFn: doneFn
    });
  }
  /**
   * Queue an update to remove a time range from the buffer.
   *
   * @param {number} start where to start the removal
   * @param {number} end where to end the removal
   * @param {Function} [done=noop] optional callback to be executed when the remove
   * operation is complete
   * @see http://www.w3.org/TR/media-source/#widl-SourceBuffer-remove-void-double-start-unrestricted-double-end
   */
  ;

  _proto.removeAudio = function removeAudio(start, end, done) {
    if (done === void 0) {
      done = noop;
    }

    if (!this.audioBuffered().length || this.audioBuffered().end(0) === 0) {
      done();
      return;
    }

    pushQueue({
      type: 'audio',
      sourceUpdater: this,
      action: actions.remove(start, end),
      doneFn: done,
      name: 'remove'
    });
  }
  /**
   * Queue an update to remove a time range from the buffer.
   *
   * @param {number} start where to start the removal
   * @param {number} end where to end the removal
   * @param {Function} [done=noop] optional callback to be executed when the remove
   * operation is complete
   * @see http://www.w3.org/TR/media-source/#widl-SourceBuffer-remove-void-double-start-unrestricted-double-end
   */
  ;

  _proto.removeVideo = function removeVideo(start, end, done) {
    if (done === void 0) {
      done = noop;
    }

    if (!this.videoBuffered().length || this.videoBuffered().end(0) === 0) {
      done();
      return;
    }

    pushQueue({
      type: 'video',
      sourceUpdater: this,
      action: actions.remove(start, end),
      doneFn: done,
      name: 'remove'
    });
  }
  /**
   * Whether the underlying sourceBuffer is updating or not
   *
   * @return {boolean} the updating status of the SourceBuffer
   */
  ;

  _proto.updating = function updating() {
    // the audio/video source buffer is updating
    if (_updating('audio', this) || _updating('video', this)) {
      return true;
    }

    return false;
  }
  /**
   * Set/get the timestampoffset on the audio SourceBuffer
   *
   * @return {number} the timestamp offset
   */
  ;

  _proto.audioTimestampOffset = function audioTimestampOffset(offset) {
    if (typeof offset !== 'undefined' && this.audioBuffer && // no point in updating if it's the same
    this.audioTimestampOffset_ !== offset) {
      pushQueue({
        type: 'audio',
        sourceUpdater: this,
        action: actions.timestampOffset(offset),
        name: 'timestampOffset'
      });
      this.audioTimestampOffset_ = offset;
    }

    return this.audioTimestampOffset_;
  }
  /**
   * Set/get the timestampoffset on the video SourceBuffer
   *
   * @return {number} the timestamp offset
   */
  ;

  _proto.videoTimestampOffset = function videoTimestampOffset(offset) {
    if (typeof offset !== 'undefined' && this.videoBuffer && // no point in updating if it's the same
    this.videoTimestampOffset !== offset) {
      pushQueue({
        type: 'video',
        sourceUpdater: this,
        action: actions.timestampOffset(offset),
        name: 'timestampOffset'
      });
      this.videoTimestampOffset_ = offset;
    }

    return this.videoTimestampOffset_;
  }
  /**
   * Add a function to the queue that will be called
   * when it is its turn to run in the audio queue.
   *
   * @param {Function} callback
   *        The callback to queue.
   */
  ;

  _proto.audioQueueCallback = function audioQueueCallback(callback) {
    if (!this.audioBuffer) {
      return;
    }

    pushQueue({
      type: 'audio',
      sourceUpdater: this,
      action: actions.callback(callback),
      name: 'callback'
    });
  }
  /**
   * Add a function to the queue that will be called
   * when it is its turn to run in the video queue.
   *
   * @param {Function} callback
   *        The callback to queue.
   */
  ;

  _proto.videoQueueCallback = function videoQueueCallback(callback) {
    if (!this.videoBuffer) {
      return;
    }

    pushQueue({
      type: 'video',
      sourceUpdater: this,
      action: actions.callback(callback),
      name: 'callback'
    });
  }
  /**
   * dispose of the source updater and the underlying sourceBuffer
   */
  ;

  _proto.dispose = function dispose() {
    var _this4 = this;

    this.trigger('dispose');
    bufferTypes.forEach(function (type) {
      _this4.abort(type);

      if (_this4.canRemoveSourceBuffer()) {
        _this4.removeSourceBuffer(type);
      } else {
        _this4[type + "QueueCallback"](function () {
          return cleanupBuffer(type, _this4);
        });
      }
    });
    this.videoAppendQueued_ = false;
    this.delayedAudioAppendQueue_.length = 0;

    if (this.sourceopenListener_) {
      this.mediaSource.removeEventListener('sourceopen', this.sourceopenListener_);
    }

    this.off();
  };

  return SourceUpdater;
}(videojs$1.EventTarget);

var uint8ToUtf8 = function uint8ToUtf8(uintArray) {
  return decodeURIComponent(escape(String.fromCharCode.apply(null, uintArray)));
};

var VTT_LINE_TERMINATORS = new Uint8Array('\n\n'.split('').map(function (char) {
  return char.charCodeAt(0);
}));
/**
 * An object that manages segment loading and appending.
 *
 * @class VTTSegmentLoader
 * @param {Object} options required and optional options
 * @extends videojs.EventTarget
 */

var VTTSegmentLoader = /*#__PURE__*/function (_SegmentLoader) {
  inheritsLoose(VTTSegmentLoader, _SegmentLoader);

  function VTTSegmentLoader(settings, options) {
    var _this;

    if (options === void 0) {
      options = {};
    }

    _this = _SegmentLoader.call(this, settings, options) || this; // VTT can't handle partial data

    _this.handlePartialData_ = false; // SegmentLoader requires a MediaSource be specified or it will throw an error;
    // however, VTTSegmentLoader has no need of a media source, so delete the reference

    _this.mediaSource_ = null;
    _this.subtitlesTrack_ = null;
    _this.loaderType_ = 'subtitle';
    _this.featuresNativeTextTracks_ = settings.featuresNativeTextTracks; // The VTT segment will have its own time mappings. Saving VTT segment timing info in
    // the sync controller leads to improper behavior.

    _this.shouldSaveSegmentTimingInfo_ = false;
    return _this;
  }

  var _proto = VTTSegmentLoader.prototype;

  _proto.createTransmuxer_ = function createTransmuxer_() {
    // don't need to transmux any subtitles
    return null;
  }
  /**
   * Indicates which time ranges are buffered
   *
   * @return {TimeRange}
   *         TimeRange object representing the current buffered ranges
   */
  ;

  _proto.buffered_ = function buffered_() {
    if (!this.subtitlesTrack_ || !this.subtitlesTrack_.cues.length) {
      return videojs$1.createTimeRanges();
    }

    var cues = this.subtitlesTrack_.cues;
    var start = cues[0].startTime;
    var end = cues[cues.length - 1].startTime;
    return videojs$1.createTimeRanges([[start, end]]);
  }
  /**
   * Gets and sets init segment for the provided map
   *
   * @param {Object} map
   *        The map object representing the init segment to get or set
   * @param {boolean=} set
   *        If true, the init segment for the provided map should be saved
   * @return {Object}
   *         map object for desired init segment
   */
  ;

  _proto.initSegmentForMap = function initSegmentForMap(map, set) {
    if (set === void 0) {
      set = false;
    }

    if (!map) {
      return null;
    }

    var id = initSegmentId(map);
    var storedMap = this.initSegments_[id];

    if (set && !storedMap && map.bytes) {
      // append WebVTT line terminators to the media initialization segment if it exists
      // to follow the WebVTT spec (https://w3c.github.io/webvtt/#file-structure) that
      // requires two or more WebVTT line terminators between the WebVTT header and the
      // rest of the file
      var combinedByteLength = VTT_LINE_TERMINATORS.byteLength + map.bytes.byteLength;
      var combinedSegment = new Uint8Array(combinedByteLength);
      combinedSegment.set(map.bytes);
      combinedSegment.set(VTT_LINE_TERMINATORS, map.bytes.byteLength);
      this.initSegments_[id] = storedMap = {
        resolvedUri: map.resolvedUri,
        byterange: map.byterange,
        bytes: combinedSegment
      };
    }

    return storedMap || map;
  }
  /**
   * Returns true if all configuration required for loading is present, otherwise false.
   *
   * @return {boolean} True if the all configuration is ready for loading
   * @private
   */
  ;

  _proto.couldBeginLoading_ = function couldBeginLoading_() {
    return this.playlist_ && this.subtitlesTrack_ && !this.paused();
  }
  /**
   * Once all the starting parameters have been specified, begin
   * operation. This method should only be invoked from the INIT
   * state.
   *
   * @private
   */
  ;

  _proto.init_ = function init_() {
    this.state = 'READY';
    this.resetEverything();
    return this.monitorBuffer_();
  }
  /**
   * Set a subtitle track on the segment loader to add subtitles to
   *
   * @param {TextTrack=} track
   *        The text track to add loaded subtitles to
   * @return {TextTrack}
   *        Returns the subtitles track
   */
  ;

  _proto.track = function track(_track) {
    if (typeof _track === 'undefined') {
      return this.subtitlesTrack_;
    }

    this.subtitlesTrack_ = _track; // if we were unpaused but waiting for a sourceUpdater, start
    // buffering now

    if (this.state === 'INIT' && this.couldBeginLoading_()) {
      this.init_();
    }

    return this.subtitlesTrack_;
  }
  /**
   * Remove any data in the source buffer between start and end times
   *
   * @param {number} start - the start time of the region to remove from the buffer
   * @param {number} end - the end time of the region to remove from the buffer
   */
  ;

  _proto.remove = function remove(start, end) {
    removeCuesFromTrack(start, end, this.subtitlesTrack_);
  }
  /**
   * fill the buffer with segements unless the sourceBuffers are
   * currently updating
   *
   * Note: this function should only ever be called by monitorBuffer_
   * and never directly
   *
   * @private
   */
  ;

  _proto.fillBuffer_ = function fillBuffer_() {
    var _this2 = this;

    if (!this.syncPoint_) {
      this.syncPoint_ = this.syncController_.getSyncPoint(this.playlist_, this.duration_(), this.currentTimeline_, this.currentTime_());
    } // see if we need to begin loading immediately


    var segmentInfo = this.checkBuffer_(this.buffered_(), this.playlist_, this.mediaIndex, this.hasPlayed_(), this.currentTime_(), this.syncPoint_);
    segmentInfo = this.skipEmptySegments_(segmentInfo);

    if (!segmentInfo) {
      return;
    }

    if (this.syncController_.timestampOffsetForTimeline(segmentInfo.timeline) === null) {
      // We don't have the timestamp offset that we need to sync subtitles.
      // Rerun on a timestamp offset or user interaction.
      var checkTimestampOffset = function checkTimestampOffset() {
        _this2.state = 'READY';

        if (!_this2.paused()) {
          // if not paused, queue a buffer check as soon as possible
          _this2.monitorBuffer_();
        }
      };

      this.syncController_.one('timestampoffset', checkTimestampOffset);
      this.state = 'WAITING_ON_TIMELINE';
      return;
    }

    this.loadSegment_(segmentInfo);
  }
  /**
   * Prevents the segment loader from requesting segments we know contain no subtitles
   * by walking forward until we find the next segment that we don't know whether it is
   * empty or not.
   *
   * @param {Object} segmentInfo
   *        a segment info object that describes the current segment
   * @return {Object}
   *         a segment info object that describes the current segment
   */
  ;

  _proto.skipEmptySegments_ = function skipEmptySegments_(segmentInfo) {
    while (segmentInfo && segmentInfo.segment.empty) {
      segmentInfo = this.generateSegmentInfo_(segmentInfo.playlist, segmentInfo.mediaIndex + 1, segmentInfo.startOfSegment + segmentInfo.duration, segmentInfo.isSyncRequest);
    }

    return segmentInfo;
  };

  _proto.stopForError = function stopForError(error) {
    this.error(error);
    this.state = 'READY';
    this.pause();
    this.trigger('error');
  }
  /**
   * append a decrypted segement to the SourceBuffer through a SourceUpdater
   *
   * @private
   */
  ;

  _proto.segmentRequestFinished_ = function segmentRequestFinished_(error, simpleSegment, result) {
    var _this3 = this;

    if (!this.subtitlesTrack_) {
      this.state = 'READY';
      return;
    }

    this.saveTransferStats_(simpleSegment.stats); // the request was aborted

    if (!this.pendingSegment_) {
      this.state = 'READY';
      this.mediaRequestsAborted += 1;
      return;
    }

    if (error) {
      if (error.code === REQUEST_ERRORS.TIMEOUT) {
        this.handleTimeout_();
      }

      if (error.code === REQUEST_ERRORS.ABORTED) {
        this.mediaRequestsAborted += 1;
      } else {
        this.mediaRequestsErrored += 1;
      }

      this.stopForError(error);
      return;
    } // although the VTT segment loader bandwidth isn't really used, it's good to
    // maintain functionality between segment loaders


    this.saveBandwidthRelatedStats_(simpleSegment.stats);
    this.state = 'APPENDING'; // used for tests

    this.trigger('appending');
    var segmentInfo = this.pendingSegment_;
    var segment = segmentInfo.segment;

    if (segment.map) {
      segment.map.bytes = simpleSegment.map.bytes;
    }

    segmentInfo.bytes = simpleSegment.bytes; // Make sure that vttjs has loaded, otherwise, wait till it finished loading

    if (typeof window_1.WebVTT !== 'function' && this.subtitlesTrack_ && this.subtitlesTrack_.tech_) {
      var loadHandler;

      var errorHandler = function errorHandler() {
        _this3.subtitlesTrack_.tech_.off('vttjsloaded', loadHandler);

        _this3.stopForError({
          message: 'Error loading vtt.js'
        });

        return;
      };

      loadHandler = function loadHandler() {
        _this3.subtitlesTrack_.tech_.off('vttjserror', errorHandler);

        _this3.segmentRequestFinished_(error, simpleSegment, result);
      };

      this.state = 'WAITING_ON_VTTJS';
      this.subtitlesTrack_.tech_.one('vttjsloaded', loadHandler);
      this.subtitlesTrack_.tech_.one('vttjserror', errorHandler);
      return;
    }

    segment.requested = true;

    try {
      this.parseVTTCues_(segmentInfo);
    } catch (e) {
      this.stopForError({
        message: e.message
      });
      return;
    }

    this.updateTimeMapping_(segmentInfo, this.syncController_.timelines[segmentInfo.timeline], this.playlist_);

    if (segmentInfo.cues.length) {
      segmentInfo.timingInfo = {
        start: segmentInfo.cues[0].startTime,
        end: segmentInfo.cues[segmentInfo.cues.length - 1].endTime
      };
    } else {
      segmentInfo.timingInfo = {
        start: segmentInfo.startOfSegment,
        end: segmentInfo.startOfSegment + segmentInfo.duration
      };
    }

    if (segmentInfo.isSyncRequest) {
      this.trigger('syncinfoupdate');
      this.pendingSegment_ = null;
      this.state = 'READY';
      return;
    }

    segmentInfo.byteLength = segmentInfo.bytes.byteLength;
    this.mediaSecondsLoaded += segment.duration;
    segmentInfo.cues.forEach(function (cue) {
      // remove any overlapping cues to prevent doubling
      _this3.remove(cue.startTime, cue.endTime);

      _this3.subtitlesTrack_.addCue(_this3.featuresNativeTextTracks_ ? new window_1.VTTCue(cue.startTime, cue.endTime, cue.text) : cue);
    });
    this.handleAppendsDone_();
  };

  _proto.handleData_ = function handleData_() {// noop as we shouldn't be getting video/audio data captions
    // that we do not support here.
  };

  _proto.updateTimingInfoEnd_ = function updateTimingInfoEnd_() {} // noop

  /**
   * Uses the WebVTT parser to parse the segment response
   *
   * @param {Object} segmentInfo
   *        a segment info object that describes the current segment
   * @private
   */
  ;

  _proto.parseVTTCues_ = function parseVTTCues_(segmentInfo) {
    var decoder;
    var decodeBytesToString = false;

    if (typeof window_1.TextDecoder === 'function') {
      decoder = new window_1.TextDecoder('utf8');
    } else {
      decoder = window_1.WebVTT.StringDecoder();
      decodeBytesToString = true;
    }

    var parser = new window_1.WebVTT.Parser(window_1, window_1.vttjs, decoder);
    segmentInfo.cues = [];
    segmentInfo.timestampmap = {
      MPEGTS: 0,
      LOCAL: 0
    };
    parser.oncue = segmentInfo.cues.push.bind(segmentInfo.cues);

    parser.ontimestampmap = function (map) {
      segmentInfo.timestampmap = map;
    };

    parser.onparsingerror = function (error) {
      videojs$1.log.warn('Error encountered when parsing cues: ' + error.message);
    };

    if (segmentInfo.segment.map) {
      var mapData = segmentInfo.segment.map.bytes;

      if (decodeBytesToString) {
        mapData = uint8ToUtf8(mapData);
      }

      parser.parse(mapData);
    }

    var segmentData = segmentInfo.bytes;

    if (decodeBytesToString) {
      segmentData = uint8ToUtf8(segmentData);
    }

    parser.parse(segmentData);
    parser.flush();
  }
  /**
   * Updates the start and end times of any cues parsed by the WebVTT parser using
   * the information parsed from the X-TIMESTAMP-MAP header and a TS to media time mapping
   * from the SyncController
   *
   * @param {Object} segmentInfo
   *        a segment info object that describes the current segment
   * @param {Object} mappingObj
   *        object containing a mapping from TS to media time
   * @param {Object} playlist
   *        the playlist object containing the segment
   * @private
   */
  ;

  _proto.updateTimeMapping_ = function updateTimeMapping_(segmentInfo, mappingObj, playlist) {
    var segment = segmentInfo.segment;

    if (!mappingObj) {
      // If the sync controller does not have a mapping of TS to Media Time for the
      // timeline, then we don't have enough information to update the cue
      // start/end times
      return;
    }

    if (!segmentInfo.cues.length) {
      // If there are no cues, we also do not have enough information to figure out
      // segment timing. Mark that the segment contains no cues so we don't re-request
      // an empty segment.
      segment.empty = true;
      return;
    }

    var timestampmap = segmentInfo.timestampmap;
    var diff = timestampmap.MPEGTS / clock.ONE_SECOND_IN_TS - timestampmap.LOCAL + mappingObj.mapping;
    segmentInfo.cues.forEach(function (cue) {
      // First convert cue time to TS time using the timestamp-map provided within the vtt
      cue.startTime += diff;
      cue.endTime += diff;
    });

    if (!playlist.syncInfo) {
      var firstStart = segmentInfo.cues[0].startTime;
      var lastStart = segmentInfo.cues[segmentInfo.cues.length - 1].startTime;
      playlist.syncInfo = {
        mediaSequence: playlist.mediaSequence + segmentInfo.mediaIndex,
        time: Math.min(firstStart, lastStart - segment.duration)
      };
    }
  };

  return VTTSegmentLoader;
}(SegmentLoader);

/**
 * @file ad-cue-tags.js
 */
/**
 * Searches for an ad cue that overlaps with the given mediaTime
 */

var findAdCue = function findAdCue(track, mediaTime) {
  var cues = track.cues;

  for (var i = 0; i < cues.length; i++) {
    var cue = cues[i];

    if (mediaTime >= cue.adStartTime && mediaTime <= cue.adEndTime) {
      return cue;
    }
  }

  return null;
};
var updateAdCues = function updateAdCues(media, track, offset) {
  if (offset === void 0) {
    offset = 0;
  }

  if (!media.segments) {
    return;
  }

  var mediaTime = offset;
  var cue;

  for (var i = 0; i < media.segments.length; i++) {
    var segment = media.segments[i];

    if (!cue) {
      // Since the cues will span for at least the segment duration, adding a fudge
      // factor of half segment duration will prevent duplicate cues from being
      // created when timing info is not exact (e.g. cue start time initialized
      // at 10.006677, but next call mediaTime is 10.003332 )
      cue = findAdCue(track, mediaTime + segment.duration / 2);
    }

    if (cue) {
      if ('cueIn' in segment) {
        // Found a CUE-IN so end the cue
        cue.endTime = mediaTime;
        cue.adEndTime = mediaTime;
        mediaTime += segment.duration;
        cue = null;
        continue;
      }

      if (mediaTime < cue.endTime) {
        // Already processed this mediaTime for this cue
        mediaTime += segment.duration;
        continue;
      } // otherwise extend cue until a CUE-IN is found


      cue.endTime += segment.duration;
    } else {
      if ('cueOut' in segment) {
        cue = new window_1.VTTCue(mediaTime, mediaTime + segment.duration, segment.cueOut);
        cue.adStartTime = mediaTime; // Assumes tag format to be
        // #EXT-X-CUE-OUT:30

        cue.adEndTime = mediaTime + parseFloat(segment.cueOut);
        track.addCue(cue);
      }

      if ('cueOutCont' in segment) {
        // Entered into the middle of an ad cue
        // Assumes tag formate to be
        // #EXT-X-CUE-OUT-CONT:10/30
        var _segment$cueOutCont$s = segment.cueOutCont.split('/').map(parseFloat),
            adOffset = _segment$cueOutCont$s[0],
            adTotal = _segment$cueOutCont$s[1];

        cue = new window_1.VTTCue(mediaTime, mediaTime + segment.duration, '');
        cue.adStartTime = mediaTime - adOffset;
        cue.adEndTime = cue.adStartTime + adTotal;
        track.addCue(cue);
      }
    }

    mediaTime += segment.duration;
  }
};

var syncPointStrategies = [// Stategy "VOD": Handle the VOD-case where the sync-point is *always*
//                the equivalence display-time 0 === segment-index 0
{
  name: 'VOD',
  run: function run(syncController, playlist, duration, currentTimeline, currentTime) {
    if (duration !== Infinity) {
      var syncPoint = {
        time: 0,
        segmentIndex: 0
      };
      return syncPoint;
    }

    return null;
  }
}, // Stategy "ProgramDateTime": We have a program-date-time tag in this playlist
{
  name: 'ProgramDateTime',
  run: function run(syncController, playlist, duration, currentTimeline, currentTime) {
    if (!syncController.datetimeToDisplayTime) {
      return null;
    }

    var segments = playlist.segments || [];
    var syncPoint = null;
    var lastDistance = null;
    currentTime = currentTime || 0;

    for (var i = 0; i < segments.length; i++) {
      var segment = segments[i];

      if (segment.dateTimeObject) {
        var segmentTime = segment.dateTimeObject.getTime() / 1000;
        var segmentStart = segmentTime + syncController.datetimeToDisplayTime;
        var distance = Math.abs(currentTime - segmentStart); // Once the distance begins to increase, or if distance is 0, we have passed
        // currentTime and can stop looking for better candidates

        if (lastDistance !== null && (distance === 0 || lastDistance < distance)) {
          break;
        }

        lastDistance = distance;
        syncPoint = {
          time: segmentStart,
          segmentIndex: i
        };
      }
    }

    return syncPoint;
  }
}, // Stategy "Segment": We have a known time mapping for a timeline and a
//                    segment in the current timeline with timing data
{
  name: 'Segment',
  run: function run(syncController, playlist, duration, currentTimeline, currentTime) {
    var segments = playlist.segments || [];
    var syncPoint = null;
    var lastDistance = null;
    currentTime = currentTime || 0;

    for (var i = 0; i < segments.length; i++) {
      var segment = segments[i];

      if (segment.timeline === currentTimeline && typeof segment.start !== 'undefined') {
        var distance = Math.abs(currentTime - segment.start); // Once the distance begins to increase, we have passed
        // currentTime and can stop looking for better candidates

        if (lastDistance !== null && lastDistance < distance) {
          break;
        }

        if (!syncPoint || lastDistance === null || lastDistance >= distance) {
          lastDistance = distance;
          syncPoint = {
            time: segment.start,
            segmentIndex: i
          };
        }
      }
    }

    return syncPoint;
  }
}, // Stategy "Discontinuity": We have a discontinuity with a known
//                          display-time
{
  name: 'Discontinuity',
  run: function run(syncController, playlist, duration, currentTimeline, currentTime) {
    var syncPoint = null;
    currentTime = currentTime || 0;

    if (playlist.discontinuityStarts && playlist.discontinuityStarts.length) {
      var lastDistance = null;

      for (var i = 0; i < playlist.discontinuityStarts.length; i++) {
        var segmentIndex = playlist.discontinuityStarts[i];
        var discontinuity = playlist.discontinuitySequence + i + 1;
        var discontinuitySync = syncController.discontinuities[discontinuity];

        if (discontinuitySync) {
          var distance = Math.abs(currentTime - discontinuitySync.time); // Once the distance begins to increase, we have passed
          // currentTime and can stop looking for better candidates

          if (lastDistance !== null && lastDistance < distance) {
            break;
          }

          if (!syncPoint || lastDistance === null || lastDistance >= distance) {
            lastDistance = distance;
            syncPoint = {
              time: discontinuitySync.time,
              segmentIndex: segmentIndex
            };
          }
        }
      }
    }

    return syncPoint;
  }
}, // Stategy "Playlist": We have a playlist with a known mapping of
//                     segment index to display time
{
  name: 'Playlist',
  run: function run(syncController, playlist, duration, currentTimeline, currentTime) {
    if (playlist.syncInfo) {
      var syncPoint = {
        time: playlist.syncInfo.time,
        segmentIndex: playlist.syncInfo.mediaSequence - playlist.mediaSequence
      };
      return syncPoint;
    }

    return null;
  }
}];

var SyncController = /*#__PURE__*/function (_videojs$EventTarget) {
  inheritsLoose(SyncController, _videojs$EventTarget);

  function SyncController(options) {
    var _this;

    _this = _videojs$EventTarget.call(this) || this; // ...for synching across variants

    _this.timelines = [];
    _this.discontinuities = [];
    _this.datetimeToDisplayTime = null;
    _this.logger_ = logger('SyncController');
    return _this;
  }
  /**
   * Find a sync-point for the playlist specified
   *
   * A sync-point is defined as a known mapping from display-time to
   * a segment-index in the current playlist.
   *
   * @param {Playlist} playlist
   *        The playlist that needs a sync-point
   * @param {number} duration
   *        Duration of the MediaSource (Infinite if playing a live source)
   * @param {number} currentTimeline
   *        The last timeline from which a segment was loaded
   * @return {Object}
   *          A sync-point object
   */


  var _proto = SyncController.prototype;

  _proto.getSyncPoint = function getSyncPoint(playlist, duration, currentTimeline, currentTime) {
    var syncPoints = this.runStrategies_(playlist, duration, currentTimeline, currentTime);

    if (!syncPoints.length) {
      // Signal that we need to attempt to get a sync-point manually
      // by fetching a segment in the playlist and constructing
      // a sync-point from that information
      return null;
    } // Now find the sync-point that is closest to the currentTime because
    // that should result in the most accurate guess about which segment
    // to fetch


    return this.selectSyncPoint_(syncPoints, {
      key: 'time',
      value: currentTime
    });
  }
  /**
   * Calculate the amount of time that has expired off the playlist during playback
   *
   * @param {Playlist} playlist
   *        Playlist object to calculate expired from
   * @param {number} duration
   *        Duration of the MediaSource (Infinity if playling a live source)
   * @return {number|null}
   *          The amount of time that has expired off the playlist during playback. Null
   *          if no sync-points for the playlist can be found.
   */
  ;

  _proto.getExpiredTime = function getExpiredTime(playlist, duration) {
    if (!playlist || !playlist.segments) {
      return null;
    }

    var syncPoints = this.runStrategies_(playlist, duration, playlist.discontinuitySequence, 0); // Without sync-points, there is not enough information to determine the expired time

    if (!syncPoints.length) {
      return null;
    }

    var syncPoint = this.selectSyncPoint_(syncPoints, {
      key: 'segmentIndex',
      value: 0
    }); // If the sync-point is beyond the start of the playlist, we want to subtract the
    // duration from index 0 to syncPoint.segmentIndex instead of adding.

    if (syncPoint.segmentIndex > 0) {
      syncPoint.time *= -1;
    }

    return Math.abs(syncPoint.time + sumDurations(playlist, syncPoint.segmentIndex, 0));
  }
  /**
   * Runs each sync-point strategy and returns a list of sync-points returned by the
   * strategies
   *
   * @private
   * @param {Playlist} playlist
   *        The playlist that needs a sync-point
   * @param {number} duration
   *        Duration of the MediaSource (Infinity if playing a live source)
   * @param {number} currentTimeline
   *        The last timeline from which a segment was loaded
   * @return {Array}
   *          A list of sync-point objects
   */
  ;

  _proto.runStrategies_ = function runStrategies_(playlist, duration, currentTimeline, currentTime) {
    var syncPoints = []; // Try to find a sync-point in by utilizing various strategies...

    for (var i = 0; i < syncPointStrategies.length; i++) {
      var strategy = syncPointStrategies[i];
      var syncPoint = strategy.run(this, playlist, duration, currentTimeline, currentTime);

      if (syncPoint) {
        syncPoint.strategy = strategy.name;
        syncPoints.push({
          strategy: strategy.name,
          syncPoint: syncPoint
        });
      }
    }

    return syncPoints;
  }
  /**
   * Selects the sync-point nearest the specified target
   *
   * @private
   * @param {Array} syncPoints
   *        List of sync-points to select from
   * @param {Object} target
   *        Object specifying the property and value we are targeting
   * @param {string} target.key
   *        Specifies the property to target. Must be either 'time' or 'segmentIndex'
   * @param {number} target.value
   *        The value to target for the specified key.
   * @return {Object}
   *          The sync-point nearest the target
   */
  ;

  _proto.selectSyncPoint_ = function selectSyncPoint_(syncPoints, target) {
    var bestSyncPoint = syncPoints[0].syncPoint;
    var bestDistance = Math.abs(syncPoints[0].syncPoint[target.key] - target.value);
    var bestStrategy = syncPoints[0].strategy;

    for (var i = 1; i < syncPoints.length; i++) {
      var newDistance = Math.abs(syncPoints[i].syncPoint[target.key] - target.value);

      if (newDistance < bestDistance) {
        bestDistance = newDistance;
        bestSyncPoint = syncPoints[i].syncPoint;
        bestStrategy = syncPoints[i].strategy;
      }
    }

    this.logger_("syncPoint for [" + target.key + ": " + target.value + "] chosen with strategy" + (" [" + bestStrategy + "]: [time:" + bestSyncPoint.time + ",") + (" segmentIndex:" + bestSyncPoint.segmentIndex + "]"));
    return bestSyncPoint;
  }
  /**
   * Save any meta-data present on the segments when segments leave
   * the live window to the playlist to allow for synchronization at the
   * playlist level later.
   *
   * @param {Playlist} oldPlaylist - The previous active playlist
   * @param {Playlist} newPlaylist - The updated and most current playlist
   */
  ;

  _proto.saveExpiredSegmentInfo = function saveExpiredSegmentInfo(oldPlaylist, newPlaylist) {
    var mediaSequenceDiff = newPlaylist.mediaSequence - oldPlaylist.mediaSequence; // When a segment expires from the playlist and it has a start time
    // save that information as a possible sync-point reference in future

    for (var i = mediaSequenceDiff - 1; i >= 0; i--) {
      var lastRemovedSegment = oldPlaylist.segments[i];

      if (lastRemovedSegment && typeof lastRemovedSegment.start !== 'undefined') {
        newPlaylist.syncInfo = {
          mediaSequence: oldPlaylist.mediaSequence + i,
          time: lastRemovedSegment.start
        };
        this.logger_("playlist refresh sync: [time:" + newPlaylist.syncInfo.time + "," + (" mediaSequence: " + newPlaylist.syncInfo.mediaSequence + "]"));
        this.trigger('syncinfoupdate');
        break;
      }
    }
  }
  /**
   * Save the mapping from playlist's ProgramDateTime to display. This should
   * only ever happen once at the start of playback.
   *
   * @param {Playlist} playlist - The currently active playlist
   */
  ;

  _proto.setDateTimeMapping = function setDateTimeMapping(playlist) {
    if (!this.datetimeToDisplayTime && playlist.segments && playlist.segments.length && playlist.segments[0].dateTimeObject) {
      var playlistTimestamp = playlist.segments[0].dateTimeObject.getTime() / 1000;
      this.datetimeToDisplayTime = -playlistTimestamp;
    }
  }
  /**
   * Calculates and saves timeline mappings, playlist sync info, and segment timing values
   * based on the latest timing information.
   *
   * @param {Object} options
   *        Options object
   * @param {SegmentInfo} options.segmentInfo
   *        The current active request information
   * @param {boolean} options.shouldSaveTimelineMapping
   *        If there's a timeline change, determines if the timeline mapping should be
   *        saved in timelines.
   */
  ;

  _proto.saveSegmentTimingInfo = function saveSegmentTimingInfo(_ref) {
    var segmentInfo = _ref.segmentInfo,
        shouldSaveTimelineMapping = _ref.shouldSaveTimelineMapping;
    var didCalculateSegmentTimeMapping = this.calculateSegmentTimeMapping_(segmentInfo, segmentInfo.timingInfo, shouldSaveTimelineMapping);

    if (didCalculateSegmentTimeMapping) {
      this.saveDiscontinuitySyncInfo_(segmentInfo); // If the playlist does not have sync information yet, record that information
      // now with segment timing information

      if (!segmentInfo.playlist.syncInfo) {
        segmentInfo.playlist.syncInfo = {
          mediaSequence: segmentInfo.playlist.mediaSequence + segmentInfo.mediaIndex,
          time: segmentInfo.segment.start
        };
      }
    }
  };

  _proto.timestampOffsetForTimeline = function timestampOffsetForTimeline(timeline) {
    if (typeof this.timelines[timeline] === 'undefined') {
      return null;
    }

    return this.timelines[timeline].time;
  };

  _proto.mappingForTimeline = function mappingForTimeline(timeline) {
    if (typeof this.timelines[timeline] === 'undefined') {
      return null;
    }

    return this.timelines[timeline].mapping;
  }
  /**
   * Use the "media time" for a segment to generate a mapping to "display time" and
   * save that display time to the segment.
   *
   * @private
   * @param {SegmentInfo} segmentInfo
   *        The current active request information
   * @param {Object} timingInfo
   *        The start and end time of the current segment in "media time"
   * @param {boolean} shouldSaveTimelineMapping
   *        If there's a timeline change, determines if the timeline mapping should be
   *        saved in timelines.
   * @return {boolean}
   *          Returns false if segment time mapping could not be calculated
   */
  ;

  _proto.calculateSegmentTimeMapping_ = function calculateSegmentTimeMapping_(segmentInfo, timingInfo, shouldSaveTimelineMapping) {
    var segment = segmentInfo.segment;
    var mappingObj = this.timelines[segmentInfo.timeline];

    if (segmentInfo.timestampOffset !== null) {
      mappingObj = {
        time: segmentInfo.startOfSegment,
        mapping: segmentInfo.startOfSegment - timingInfo.start
      };

      if (shouldSaveTimelineMapping) {
        this.timelines[segmentInfo.timeline] = mappingObj;
        this.trigger('timestampoffset');
        this.logger_("time mapping for timeline " + segmentInfo.timeline + ": " + ("[time: " + mappingObj.time + "] [mapping: " + mappingObj.mapping + "]"));
      }

      segment.start = segmentInfo.startOfSegment;
      segment.end = timingInfo.end + mappingObj.mapping;
    } else if (mappingObj) {
      segment.start = timingInfo.start + mappingObj.mapping;
      segment.end = timingInfo.end + mappingObj.mapping;
    } else {
      return false;
    }

    return true;
  }
  /**
   * Each time we have discontinuity in the playlist, attempt to calculate the location
   * in display of the start of the discontinuity and save that. We also save an accuracy
   * value so that we save values with the most accuracy (closest to 0.)
   *
   * @private
   * @param {SegmentInfo} segmentInfo - The current active request information
   */
  ;

  _proto.saveDiscontinuitySyncInfo_ = function saveDiscontinuitySyncInfo_(segmentInfo) {
    var playlist = segmentInfo.playlist;
    var segment = segmentInfo.segment; // If the current segment is a discontinuity then we know exactly where
    // the start of the range and it's accuracy is 0 (greater accuracy values
    // mean more approximation)

    if (segment.discontinuity) {
      this.discontinuities[segment.timeline] = {
        time: segment.start,
        accuracy: 0
      };
    } else if (playlist.discontinuityStarts && playlist.discontinuityStarts.length) {
      // Search for future discontinuities that we can provide better timing
      // information for and save that information for sync purposes
      for (var i = 0; i < playlist.discontinuityStarts.length; i++) {
        var segmentIndex = playlist.discontinuityStarts[i];
        var discontinuity = playlist.discontinuitySequence + i + 1;
        var mediaIndexDiff = segmentIndex - segmentInfo.mediaIndex;
        var accuracy = Math.abs(mediaIndexDiff);

        if (!this.discontinuities[discontinuity] || this.discontinuities[discontinuity].accuracy > accuracy) {
          var time = void 0;

          if (mediaIndexDiff < 0) {
            time = segment.start - sumDurations(playlist, segmentInfo.mediaIndex, segmentIndex);
          } else {
            time = segment.end + sumDurations(playlist, segmentInfo.mediaIndex + 1, segmentIndex);
          }

          this.discontinuities[discontinuity] = {
            time: time,
            accuracy: accuracy
          };
        }
      }
    }
  };

  _proto.dispose = function dispose() {
    this.trigger('dispose');
    this.off();
  };

  return SyncController;
}(videojs$1.EventTarget);

/**
 * The TimelineChangeController acts as a source for segment loaders to listen for and
 * keep track of latest and pending timeline changes. This is useful to ensure proper
 * sync, as each loader may need to make a consideration for what timeline the other
 * loader is on before making changes which could impact the other loader's media.
 *
 * @class TimelineChangeController
 * @extends videojs.EventTarget
 */

var TimelineChangeController = /*#__PURE__*/function (_videojs$EventTarget) {
  inheritsLoose(TimelineChangeController, _videojs$EventTarget);

  function TimelineChangeController() {
    var _this;

    _this = _videojs$EventTarget.call(this) || this;
    _this.pendingTimelineChanges_ = {};
    _this.lastTimelineChanges_ = {};
    return _this;
  }

  var _proto = TimelineChangeController.prototype;

  _proto.clearPendingTimelineChange = function clearPendingTimelineChange(type) {
    this.pendingTimelineChanges_[type] = null;
    this.trigger('pendingtimelinechange');
  };

  _proto.pendingTimelineChange = function pendingTimelineChange(_ref) {
    var type = _ref.type,
        from = _ref.from,
        to = _ref.to;

    if (typeof from === 'number' && typeof to === 'number') {
      this.pendingTimelineChanges_[type] = {
        type: type,
        from: from,
        to: to
      };
      this.trigger('pendingtimelinechange');
    }

    return this.pendingTimelineChanges_[type];
  };

  _proto.lastTimelineChange = function lastTimelineChange(_ref2) {
    var type = _ref2.type,
        from = _ref2.from,
        to = _ref2.to;

    if (typeof from === 'number' && typeof to === 'number') {
      this.lastTimelineChanges_[type] = {
        type: type,
        from: from,
        to: to
      };
      delete this.pendingTimelineChanges_[type];
      this.trigger('timelinechange');
    }

    return this.lastTimelineChanges_[type];
  };

  _proto.dispose = function dispose() {
    this.trigger('dispose');
    this.pendingTimelineChanges_ = {};
    this.lastTimelineChanges_ = {};
    this.off();
  };

  return TimelineChangeController;
}(videojs$1.EventTarget);

var Decrypter = new shimWorker("./decrypter-worker.worker.js", function (window, document) {
  var self = this;
  /*! @name @videojs/http-streaming @version 2.2.0 @license Apache-2.0 */

  var decrypterWorker = function () {

    function _defineProperties(target, props) {
      for (var i = 0; i < props.length; i++) {
        var descriptor = props[i];
        descriptor.enumerable = descriptor.enumerable || false;
        descriptor.configurable = true;
        if ("value" in descriptor) descriptor.writable = true;
        Object.defineProperty(target, descriptor.key, descriptor);
      }
    }

    function _createClass(Constructor, protoProps, staticProps) {
      if (protoProps) _defineProperties(Constructor.prototype, protoProps);
      if (staticProps) _defineProperties(Constructor, staticProps);
      return Constructor;
    }

    var createClass = _createClass;

    function _inheritsLoose(subClass, superClass) {
      subClass.prototype = Object.create(superClass.prototype);
      subClass.prototype.constructor = subClass;
      subClass.__proto__ = superClass;
    }

    var inheritsLoose = _inheritsLoose;
    /*! @name @videojs/vhs-utils @version 2.2.0 @license MIT */

    /**
     * @file stream.js
     */

    /**
     * A lightweight readable stream implemention that handles event dispatching.
     *
     * @class Stream
     */

    var Stream = /*#__PURE__*/function () {
      function Stream() {
        this.listeners = {};
      }
      /**
       * Add a listener for a specified event type.
       *
       * @param {string} type the event name
       * @param {Function} listener the callback to be invoked when an event of
       * the specified type occurs
       */


      var _proto = Stream.prototype;

      _proto.on = function on(type, listener) {
        if (!this.listeners[type]) {
          this.listeners[type] = [];
        }

        this.listeners[type].push(listener);
      }
      /**
       * Remove a listener for a specified event type.
       *
       * @param {string} type the event name
       * @param {Function} listener  a function previously registered for this
       * type of event through `on`
       * @return {boolean} if we could turn it off or not
       */
      ;

      _proto.off = function off(type, listener) {
        if (!this.listeners[type]) {
          return false;
        }

        var index = this.listeners[type].indexOf(listener); // TODO: which is better?
        // In Video.js we slice listener functions
        // on trigger so that it does not mess up the order
        // while we loop through.
        //
        // Here we slice on off so that the loop in trigger
        // can continue using it's old reference to loop without
        // messing up the order.

        this.listeners[type] = this.listeners[type].slice(0);
        this.listeners[type].splice(index, 1);
        return index > -1;
      }
      /**
       * Trigger an event of the specified type on this stream. Any additional
       * arguments to this function are passed as parameters to event listeners.
       *
       * @param {string} type the event name
       */
      ;

      _proto.trigger = function trigger(type) {
        var callbacks = this.listeners[type];

        if (!callbacks) {
          return;
        } // Slicing the arguments on every invocation of this method
        // can add a significant amount of overhead. Avoid the
        // intermediate object creation for the common case of a
        // single callback argument


        if (arguments.length === 2) {
          var length = callbacks.length;

          for (var i = 0; i < length; ++i) {
            callbacks[i].call(this, arguments[1]);
          }
        } else {
          var args = Array.prototype.slice.call(arguments, 1);
          var _length = callbacks.length;

          for (var _i = 0; _i < _length; ++_i) {
            callbacks[_i].apply(this, args);
          }
        }
      }
      /**
       * Destroys the stream and cleans up.
       */
      ;

      _proto.dispose = function dispose() {
        this.listeners = {};
      }
      /**
       * Forwards all `data` events on this stream to the destination stream. The
       * destination stream should provide a method `push` to receive the data
       * events as they arrive.
       *
       * @param {Stream} destination the stream that will receive all `data` events
       * @see http://nodejs.org/api/stream.html#stream_readable_pipe_destination_options
       */
      ;

      _proto.pipe = function pipe(destination) {
        this.on('data', function (data) {
          destination.push(data);
        });
      };

      return Stream;
    }();

    var stream = Stream;
    /*! @name pkcs7 @version 1.0.4 @license Apache-2.0 */

    /**
     * Returns the subarray of a Uint8Array without PKCS#7 padding.
     *
     * @param padded {Uint8Array} unencrypted bytes that have been padded
     * @return {Uint8Array} the unpadded bytes
     * @see http://tools.ietf.org/html/rfc5652
     */

    function unpad(padded) {
      return padded.subarray(0, padded.byteLength - padded[padded.byteLength - 1]);
    }
    /*! @name aes-decrypter @version 3.0.2 @license Apache-2.0 */

    /**
     * @file aes.js
     *
     * This file contains an adaptation of the AES decryption algorithm
     * from the Standford Javascript Cryptography Library. That work is
     * covered by the following copyright and permissions notice:
     *
     * Copyright 2009-2010 Emily Stark, Mike Hamburg, Dan Boneh.
     * All rights reserved.
     *
     * Redistribution and use in source and binary forms, with or without
     * modification, are permitted provided that the following conditions are
     * met:
     *
     * 1. Redistributions of source code must retain the above copyright
     *    notice, this list of conditions and the following disclaimer.
     *
     * 2. Redistributions in binary form must reproduce the above
     *    copyright notice, this list of conditions and the following
     *    disclaimer in the documentation and/or other materials provided
     *    with the distribution.
     *
     * THIS SOFTWARE IS PROVIDED BY THE AUTHORS ``AS IS'' AND ANY EXPRESS OR
     * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
     * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
     * DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> OR CONTRIBUTORS BE
     * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
     * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
     * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR
     * BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
     * WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
     * OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN
     * IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
     *
     * The views and conclusions contained in the software and documentation
     * are those of the authors and should not be interpreted as representing
     * official policies, either expressed or implied, of the authors.
     */

    /**
     * Expand the S-box tables.
     *
     * @private
     */


    var precompute = function precompute() {
      var tables = [[[], [], [], [], []], [[], [], [], [], []]];
      var encTable = tables[0];
      var decTable = tables[1];
      var sbox = encTable[4];
      var sboxInv = decTable[4];
      var i;
      var x;
      var xInv;
      var d = [];
      var th = [];
      var x2;
      var x4;
      var x8;
      var s;
      var tEnc;
      var tDec; // Compute double and third tables

      for (i = 0; i < 256; i++) {
        th[(d[i] = i << 1 ^ (i >> 7) * 283) ^ i] = i;
      }

      for (x = xInv = 0; !sbox[x]; x ^= x2 || 1, xInv = th[xInv] || 1) {
        // Compute sbox
        s = xInv ^ xInv << 1 ^ xInv << 2 ^ xInv << 3 ^ xInv << 4;
        s = s >> 8 ^ s & 255 ^ 99;
        sbox[x] = s;
        sboxInv[s] = x; // Compute MixColumns

        x8 = d[x4 = d[x2 = d[x]]];
        tDec = x8 * 0x1010101 ^ x4 * 0x10001 ^ x2 * 0x101 ^ x * 0x1010100;
        tEnc = d[s] * 0x101 ^ s * 0x1010100;

        for (i = 0; i < 4; i++) {
          encTable[i][x] = tEnc = tEnc << 24 ^ tEnc >>> 8;
          decTable[i][s] = tDec = tDec << 24 ^ tDec >>> 8;
        }
      } // Compactify. Considerable speedup on Firefox.


      for (i = 0; i < 5; i++) {
        encTable[i] = encTable[i].slice(0);
        decTable[i] = decTable[i].slice(0);
      }

      return tables;
    };

    var aesTables = null;
    /**
     * Schedule out an AES key for both encryption and decryption. This
     * is a low-level class. Use a cipher mode to do bulk encryption.
     *
     * @class AES
     * @param key {Array} The key as an array of 4, 6 or 8 words.
     */

    var AES = /*#__PURE__*/function () {
      function AES(key) {
        /**
        * The expanded S-box and inverse S-box tables. These will be computed
        * on the client so that we don't have to send them down the wire.
        *
        * There are two tables, _tables[0] is for encryption and
        * _tables[1] is for decryption.
        *
        * The first 4 sub-tables are the expanded S-box with MixColumns. The
        * last (_tables[01][4]) is the S-box itself.
        *
        * @private
        */
        // if we have yet to precompute the S-box tables
        // do so now
        if (!aesTables) {
          aesTables = precompute();
        } // then make a copy of that object for use


        this._tables = [[aesTables[0][0].slice(), aesTables[0][1].slice(), aesTables[0][2].slice(), aesTables[0][3].slice(), aesTables[0][4].slice()], [aesTables[1][0].slice(), aesTables[1][1].slice(), aesTables[1][2].slice(), aesTables[1][3].slice(), aesTables[1][4].slice()]];
        var i;
        var j;
        var tmp;
        var sbox = this._tables[0][4];
        var decTable = this._tables[1];
        var keyLen = key.length;
        var rcon = 1;

        if (keyLen !== 4 && keyLen !== 6 && keyLen !== 8) {
          throw new Error('Invalid aes key size');
        }

        var encKey = key.slice(0);
        var decKey = [];
        this._key = [encKey, decKey]; // schedule encryption keys

        for (i = keyLen; i < 4 * keyLen + 28; i++) {
          tmp = encKey[i - 1]; // apply sbox

          if (i % keyLen === 0 || keyLen === 8 && i % keyLen === 4) {
            tmp = sbox[tmp >>> 24] << 24 ^ sbox[tmp >> 16 & 255] << 16 ^ sbox[tmp >> 8 & 255] << 8 ^ sbox[tmp & 255]; // shift rows and add rcon

            if (i % keyLen === 0) {
              tmp = tmp << 8 ^ tmp >>> 24 ^ rcon << 24;
              rcon = rcon << 1 ^ (rcon >> 7) * 283;
            }
          }

          encKey[i] = encKey[i - keyLen] ^ tmp;
        } // schedule decryption keys


        for (j = 0; i; j++, i--) {
          tmp = encKey[j & 3 ? i : i - 4];

          if (i <= 4 || j < 4) {
            decKey[j] = tmp;
          } else {
            decKey[j] = decTable[0][sbox[tmp >>> 24]] ^ decTable[1][sbox[tmp >> 16 & 255]] ^ decTable[2][sbox[tmp >> 8 & 255]] ^ decTable[3][sbox[tmp & 255]];
          }
        }
      }
      /**
       * Decrypt 16 bytes, specified as four 32-bit words.
       *
       * @param {number} encrypted0 the first word to decrypt
       * @param {number} encrypted1 the second word to decrypt
       * @param {number} encrypted2 the third word to decrypt
       * @param {number} encrypted3 the fourth word to decrypt
       * @param {Int32Array} out the array to write the decrypted words
       * into
       * @param {number} offset the offset into the output array to start
       * writing results
       * @return {Array} The plaintext.
       */


      var _proto = AES.prototype;

      _proto.decrypt = function decrypt(encrypted0, encrypted1, encrypted2, encrypted3, out, offset) {
        var key = this._key[1]; // state variables a,b,c,d are loaded with pre-whitened data

        var a = encrypted0 ^ key[0];
        var b = encrypted3 ^ key[1];
        var c = encrypted2 ^ key[2];
        var d = encrypted1 ^ key[3];
        var a2;
        var b2;
        var c2; // key.length === 2 ?

        var nInnerRounds = key.length / 4 - 2;
        var i;
        var kIndex = 4;
        var table = this._tables[1]; // load up the tables

        var table0 = table[0];
        var table1 = table[1];
        var table2 = table[2];
        var table3 = table[3];
        var sbox = table[4]; // Inner rounds. Cribbed from OpenSSL.

        for (i = 0; i < nInnerRounds; i++) {
          a2 = table0[a >>> 24] ^ table1[b >> 16 & 255] ^ table2[c >> 8 & 255] ^ table3[d & 255] ^ key[kIndex];
          b2 = table0[b >>> 24] ^ table1[c >> 16 & 255] ^ table2[d >> 8 & 255] ^ table3[a & 255] ^ key[kIndex + 1];
          c2 = table0[c >>> 24] ^ table1[d >> 16 & 255] ^ table2[a >> 8 & 255] ^ table3[b & 255] ^ key[kIndex + 2];
          d = table0[d >>> 24] ^ table1[a >> 16 & 255] ^ table2[b >> 8 & 255] ^ table3[c & 255] ^ key[kIndex + 3];
          kIndex += 4;
          a = a2;
          b = b2;
          c = c2;
        } // Last round.


        for (i = 0; i < 4; i++) {
          out[(3 & -i) + offset] = sbox[a >>> 24] << 24 ^ sbox[b >> 16 & 255] << 16 ^ sbox[c >> 8 & 255] << 8 ^ sbox[d & 255] ^ key[kIndex++];
          a2 = a;
          a = b;
          b = c;
          c = d;
          d = a2;
        }
      };

      return AES;
    }();
    /**
     * A wrapper around the Stream class to use setTimeout
     * and run stream "jobs" Asynchronously
     *
     * @class AsyncStream
     * @extends Stream
     */


    var AsyncStream = /*#__PURE__*/function (_Stream) {
      inheritsLoose(AsyncStream, _Stream);

      function AsyncStream() {
        var _this;

        _this = _Stream.call(this, stream) || this;
        _this.jobs = [];
        _this.delay = 1;
        _this.timeout_ = null;
        return _this;
      }
      /**
       * process an async job
       *
       * @private
       */


      var _proto = AsyncStream.prototype;

      _proto.processJob_ = function processJob_() {
        this.jobs.shift()();

        if (this.jobs.length) {
          this.timeout_ = setTimeout(this.processJob_.bind(this), this.delay);
        } else {
          this.timeout_ = null;
        }
      }
      /**
       * push a job into the stream
       *
       * @param {Function} job the job to push into the stream
       */
      ;

      _proto.push = function push(job) {
        this.jobs.push(job);

        if (!this.timeout_) {
          this.timeout_ = setTimeout(this.processJob_.bind(this), this.delay);
        }
      };

      return AsyncStream;
    }(stream);
    /**
     * Convert network-order (big-endian) bytes into their little-endian
     * representation.
     */


    var ntoh = function ntoh(word) {
      return word << 24 | (word & 0xff00) << 8 | (word & 0xff0000) >> 8 | word >>> 24;
    };
    /**
     * Decrypt bytes using AES-128 with CBC and PKCS#7 padding.
     *
     * @param {Uint8Array} encrypted the encrypted bytes
     * @param {Uint32Array} key the bytes of the decryption key
     * @param {Uint32Array} initVector the initialization vector (IV) to
     * use for the first round of CBC.
     * @return {Uint8Array} the decrypted bytes
     *
     * @see http://en.wikipedia.org/wiki/Advanced_Encryption_Standard
     * @see http://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Cipher_Block_Chaining_.28CBC.29
     * @see https://tools.ietf.org/html/rfc2315
     */


    var decrypt = function decrypt(encrypted, key, initVector) {
      // word-level access to the encrypted bytes
      var encrypted32 = new Int32Array(encrypted.buffer, encrypted.byteOffset, encrypted.byteLength >> 2);
      var decipher = new AES(Array.prototype.slice.call(key)); // byte and word-level access for the decrypted output

      var decrypted = new Uint8Array(encrypted.byteLength);
      var decrypted32 = new Int32Array(decrypted.buffer); // temporary variables for working with the IV, encrypted, and
      // decrypted data

      var init0;
      var init1;
      var init2;
      var init3;
      var encrypted0;
      var encrypted1;
      var encrypted2;
      var encrypted3; // iteration variable

      var wordIx; // pull out the words of the IV to ensure we don't modify the
      // passed-in reference and easier access

      init0 = initVector[0];
      init1 = initVector[1];
      init2 = initVector[2];
      init3 = initVector[3]; // decrypt four word sequences, applying cipher-block chaining (CBC)
      // to each decrypted block

      for (wordIx = 0; wordIx < encrypted32.length; wordIx += 4) {
        // convert big-endian (network order) words into little-endian
        // (javascript order)
        encrypted0 = ntoh(encrypted32[wordIx]);
        encrypted1 = ntoh(encrypted32[wordIx + 1]);
        encrypted2 = ntoh(encrypted32[wordIx + 2]);
        encrypted3 = ntoh(encrypted32[wordIx + 3]); // decrypt the block

        decipher.decrypt(encrypted0, encrypted1, encrypted2, encrypted3, decrypted32, wordIx); // XOR with the IV, and restore network byte-order to obtain the
        // plaintext

        decrypted32[wordIx] = ntoh(decrypted32[wordIx] ^ init0);
        decrypted32[wordIx + 1] = ntoh(decrypted32[wordIx + 1] ^ init1);
        decrypted32[wordIx + 2] = ntoh(decrypted32[wordIx + 2] ^ init2);
        decrypted32[wordIx + 3] = ntoh(decrypted32[wordIx + 3] ^ init3); // setup the IV for the next round

        init0 = encrypted0;
        init1 = encrypted1;
        init2 = encrypted2;
        init3 = encrypted3;
      }

      return decrypted;
    };
    /**
     * The `Decrypter` class that manages decryption of AES
     * data through `AsyncStream` objects and the `decrypt`
     * function
     *
     * @param {Uint8Array} encrypted the encrypted bytes
     * @param {Uint32Array} key the bytes of the decryption key
     * @param {Uint32Array} initVector the initialization vector (IV) to
     * @param {Function} done the function to run when done
     * @class Decrypter
     */


    var Decrypter = /*#__PURE__*/function () {
      function Decrypter(encrypted, key, initVector, done) {
        var step = Decrypter.STEP;
        var encrypted32 = new Int32Array(encrypted.buffer);
        var decrypted = new Uint8Array(encrypted.byteLength);
        var i = 0;
        this.asyncStream_ = new AsyncStream(); // split up the encryption job and do the individual chunks asynchronously

        this.asyncStream_.push(this.decryptChunk_(encrypted32.subarray(i, i + step), key, initVector, decrypted));

        for (i = step; i < encrypted32.length; i += step) {
          initVector = new Uint32Array([ntoh(encrypted32[i - 4]), ntoh(encrypted32[i - 3]), ntoh(encrypted32[i - 2]), ntoh(encrypted32[i - 1])]);
          this.asyncStream_.push(this.decryptChunk_(encrypted32.subarray(i, i + step), key, initVector, decrypted));
        } // invoke the done() callback when everything is finished


        this.asyncStream_.push(function () {
          // remove pkcs#7 padding from the decrypted bytes
          done(null, unpad(decrypted));
        });
      }
      /**
       * a getter for step the maximum number of bytes to process at one time
       *
       * @return {number} the value of step 32000
       */


      var _proto = Decrypter.prototype;
      /**
       * @private
       */

      _proto.decryptChunk_ = function decryptChunk_(encrypted, key, initVector, decrypted) {
        return function () {
          var bytes = decrypt(encrypted, key, initVector);
          decrypted.set(bytes, encrypted.byteOffset);
        };
      };

      createClass(Decrypter, null, [{
        key: "STEP",
        get: function get() {
          // 4 * 8000;
          return 32000;
        }
      }]);
      return Decrypter;
    }();
    /**
     * @file bin-utils.js
     */

    /**
     * Creates an object for sending to a web worker modifying properties that are TypedArrays
     * into a new object with seperated properties for the buffer, byteOffset, and byteLength.
     *
     * @param {Object} message
     *        Object of properties and values to send to the web worker
     * @return {Object}
     *         Modified message with TypedArray values expanded
     * @function createTransferableMessage
     */


    var createTransferableMessage = function createTransferableMessage(message) {
      var transferable = {};
      Object.keys(message).forEach(function (key) {
        var value = message[key];

        if (ArrayBuffer.isView(value)) {
          transferable[key] = {
            bytes: value.buffer,
            byteOffset: value.byteOffset,
            byteLength: value.byteLength
          };
        } else {
          transferable[key] = value;
        }
      });
      return transferable;
    };
    /* global self */

    /**
     * Our web worker interface so that things can talk to aes-decrypter
     * that will be running in a web worker. the scope is passed to this by
     * webworkify.
     *
     * @param {Object} self
     *        the scope for the web worker
     */


    var DecrypterWorker = function DecrypterWorker(self) {
      self.onmessage = function (event) {
        var data = event.data;
        var encrypted = new Uint8Array(data.encrypted.bytes, data.encrypted.byteOffset, data.encrypted.byteLength);
        var key = new Uint32Array(data.key.bytes, data.key.byteOffset, data.key.byteLength / 4);
        var iv = new Uint32Array(data.iv.bytes, data.iv.byteOffset, data.iv.byteLength / 4);
        /* eslint-disable no-new, handle-callback-err */

        new Decrypter(encrypted, key, iv, function (err, bytes) {
          self.postMessage(createTransferableMessage({
            source: data.source,
            decrypted: bytes
          }), [bytes.buffer]);
        });
        /* eslint-enable */
      };
    };

    var decrypterWorker = new DecrypterWorker(self);
    return decrypterWorker;
  }();
});

/**
 * Convert the properties of an HLS track into an audioTrackKind.
 *
 * @private
 */

var audioTrackKind_ = function audioTrackKind_(properties) {
  var kind = properties.default ? 'main' : 'alternative';

  if (properties.characteristics && properties.characteristics.indexOf('public.accessibility.describes-video') >= 0) {
    kind = 'main-desc';
  }

  return kind;
};
/**
 * Pause provided segment loader and playlist loader if active
 *
 * @param {SegmentLoader} segmentLoader
 *        SegmentLoader to pause
 * @param {Object} mediaType
 *        Active media type
 * @function stopLoaders
 */


var stopLoaders = function stopLoaders(segmentLoader, mediaType) {
  segmentLoader.abort();
  segmentLoader.pause();

  if (mediaType && mediaType.activePlaylistLoader) {
    mediaType.activePlaylistLoader.pause();
    mediaType.activePlaylistLoader = null;
  }
};
/**
 * Start loading provided segment loader and playlist loader
 *
 * @param {PlaylistLoader} playlistLoader
 *        PlaylistLoader to start loading
 * @param {Object} mediaType
 *        Active media type
 * @function startLoaders
 */

var startLoaders = function startLoaders(playlistLoader, mediaType) {
  // Segment loader will be started after `loadedmetadata` or `loadedplaylist` from the
  // playlist loader
  mediaType.activePlaylistLoader = playlistLoader;
  playlistLoader.load();
};
/**
 * Returns a function to be called when the media group changes. It performs a
 * non-destructive (preserve the buffer) resync of the SegmentLoader. This is because a
 * change of group is merely a rendition switch of the same content at another encoding,
 * rather than a change of content, such as switching audio from English to Spanish.
 *
 * @param {string} type
 *        MediaGroup type
 * @param {Object} settings
 *        Object containing required information for media groups
 * @return {Function}
 *         Handler for a non-destructive resync of SegmentLoader when the active media
 *         group changes.
 * @function onGroupChanged
 */

var onGroupChanged = function onGroupChanged(type, settings) {
  return function () {
    var _settings$segmentLoad = settings.segmentLoaders,
        segmentLoader = _settings$segmentLoad[type],
        mainSegmentLoader = _settings$segmentLoad.main,
        mediaType = settings.mediaTypes[type];
    var activeTrack = mediaType.activeTrack();
    var activeGroup = mediaType.activeGroup(activeTrack);
    var previousActiveLoader = mediaType.activePlaylistLoader;
    stopLoaders(segmentLoader, mediaType);

    if (!activeGroup) {
      // there is no group active
      return;
    }

    if (!activeGroup.playlistLoader) {
      if (previousActiveLoader) {
        // The previous group had a playlist loader but the new active group does not
        // this means we are switching from demuxed to muxed audio. In this case we want to
        // do a destructive reset of the main segment loader and not restart the audio
        // loaders.
        mainSegmentLoader.resetEverything();
      }

      return;
    } // Non-destructive resync


    segmentLoader.resyncLoader();
    startLoaders(activeGroup.playlistLoader, mediaType);
  };
};
var onGroupChanging = function onGroupChanging(type, settings) {
  return function () {
    var segmentLoader = settings.segmentLoaders[type];
    segmentLoader.abort();
    segmentLoader.pause();
  };
};
/**
 * Returns a function to be called when the media track changes. It performs a
 * destructive reset of the SegmentLoader to ensure we start loading as close to
 * currentTime as possible.
 *
 * @param {string} type
 *        MediaGroup type
 * @param {Object} settings
 *        Object containing required information for media groups
 * @return {Function}
 *         Handler for a destructive reset of SegmentLoader when the active media
 *         track changes.
 * @function onTrackChanged
 */

var onTrackChanged = function onTrackChanged(type, settings) {
  return function () {
    var _settings$segmentLoad2 = settings.segmentLoaders,
        segmentLoader = _settings$segmentLoad2[type],
        mainSegmentLoader = _settings$segmentLoad2.main,
        mediaType = settings.mediaTypes[type];
    var activeTrack = mediaType.activeTrack();
    var activeGroup = mediaType.activeGroup(activeTrack);
    var previousActiveLoader = mediaType.activePlaylistLoader;
    stopLoaders(segmentLoader, mediaType);

    if (!activeGroup) {
      // there is no group active so we do not want to restart loaders
      return;
    }

    if (type === 'AUDIO') {
      if (!activeGroup.playlistLoader) {
        // when switching from demuxed audio/video to muxed audio/video (noted by no
        // playlist loader for the audio group), we want to do a destructive reset of the
        // main segment loader and not restart the audio loaders
        mainSegmentLoader.setAudio(true); // don't have to worry about disabling the audio of the audio segment loader since
        // it should be stopped

        mainSegmentLoader.resetEverything();
        return;
      } // although the segment loader is an audio segment loader, call the setAudio
      // function to ensure it is prepared to re-append the init segment (or handle other
      // config changes)


      segmentLoader.setAudio(true);
      mainSegmentLoader.setAudio(false);
    }

    if (previousActiveLoader === activeGroup.playlistLoader) {
      // Nothing has actually changed. This can happen because track change events can fire
      // multiple times for a "single" change. One for enabling the new active track, and
      // one for disabling the track that was active
      startLoaders(activeGroup.playlistLoader, mediaType);
      return;
    }

    if (segmentLoader.track) {
      // For WebVTT, set the new text track in the segmentloader
      segmentLoader.track(activeTrack);
    } // destructive reset


    segmentLoader.resetEverything();
    startLoaders(activeGroup.playlistLoader, mediaType);
  };
};
var onError = {
  /**
   * Returns a function to be called when a SegmentLoader or PlaylistLoader encounters
   * an error.
   *
   * @param {string} type
   *        MediaGroup type
   * @param {Object} settings
   *        Object containing required information for media groups
   * @return {Function}
   *         Error handler. Logs warning (or error if the playlist is blacklisted) to
   *         console and switches back to default audio track.
   * @function onError.AUDIO
   */
  AUDIO: function AUDIO(type, settings) {
    return function () {
      var segmentLoader = settings.segmentLoaders[type],
          mediaType = settings.mediaTypes[type],
          blacklistCurrentPlaylist = settings.blacklistCurrentPlaylist;
      stopLoaders(segmentLoader, mediaType); // switch back to default audio track

      var activeTrack = mediaType.activeTrack();
      var activeGroup = mediaType.activeGroup();
      var id = (activeGroup.filter(function (group) {
        return group.default;
      })[0] || activeGroup[0]).id;
      var defaultTrack = mediaType.tracks[id];

      if (activeTrack === defaultTrack) {
        // Default track encountered an error. All we can do now is blacklist the current
        // rendition and hope another will switch audio groups
        blacklistCurrentPlaylist({
          message: 'Problem encountered loading the default audio track.'
        });
        return;
      }

      videojs$1.log.warn('Problem encountered loading the alternate audio track.' + 'Switching back to default.');

      for (var trackId in mediaType.tracks) {
        mediaType.tracks[trackId].enabled = mediaType.tracks[trackId] === defaultTrack;
      }

      mediaType.onTrackChanged();
    };
  },

  /**
   * Returns a function to be called when a SegmentLoader or PlaylistLoader encounters
   * an error.
   *
   * @param {string} type
   *        MediaGroup type
   * @param {Object} settings
   *        Object containing required information for media groups
   * @return {Function}
   *         Error handler. Logs warning to console and disables the active subtitle track
   * @function onError.SUBTITLES
   */
  SUBTITLES: function SUBTITLES(type, settings) {
    return function () {
      var segmentLoader = settings.segmentLoaders[type],
          mediaType = settings.mediaTypes[type];
      videojs$1.log.warn('Problem encountered loading the subtitle track.' + 'Disabling subtitle track.');
      stopLoaders(segmentLoader, mediaType);
      var track = mediaType.activeTrack();

      if (track) {
        track.mode = 'disabled';
      }

      mediaType.onTrackChanged();
    };
  }
};
var setupListeners = {
  /**
   * Setup event listeners for audio playlist loader
   *
   * @param {string} type
   *        MediaGroup type
   * @param {PlaylistLoader|null} playlistLoader
   *        PlaylistLoader to register listeners on
   * @param {Object} settings
   *        Object containing required information for media groups
   * @function setupListeners.AUDIO
   */
  AUDIO: function AUDIO(type, playlistLoader, settings) {
    if (!playlistLoader) {
      // no playlist loader means audio will be muxed with the video
      return;
    }

    var tech = settings.tech,
        requestOptions = settings.requestOptions,
        segmentLoader = settings.segmentLoaders[type];
    playlistLoader.on('loadedmetadata', function () {
      var media = playlistLoader.media();
      segmentLoader.playlist(media, requestOptions); // if the video is already playing, or if this isn't a live video and preload
      // permits, start downloading segments

      if (!tech.paused() || media.endList && tech.preload() !== 'none') {
        segmentLoader.load();
      }
    });
    playlistLoader.on('loadedplaylist', function () {
      segmentLoader.playlist(playlistLoader.media(), requestOptions); // If the player isn't paused, ensure that the segment loader is running

      if (!tech.paused()) {
        segmentLoader.load();
      }
    });
    playlistLoader.on('error', onError[type](type, settings));
  },

  /**
   * Setup event listeners for subtitle playlist loader
   *
   * @param {string} type
   *        MediaGroup type
   * @param {PlaylistLoader|null} playlistLoader
   *        PlaylistLoader to register listeners on
   * @param {Object} settings
   *        Object containing required information for media groups
   * @function setupListeners.SUBTITLES
   */
  SUBTITLES: function SUBTITLES(type, playlistLoader, settings) {
    var tech = settings.tech,
        requestOptions = settings.requestOptions,
        segmentLoader = settings.segmentLoaders[type],
        mediaType = settings.mediaTypes[type];
    playlistLoader.on('loadedmetadata', function () {
      var media = playlistLoader.media();
      segmentLoader.playlist(media, requestOptions);
      segmentLoader.track(mediaType.activeTrack()); // if the video is already playing, or if this isn't a live video and preload
      // permits, start downloading segments

      if (!tech.paused() || media.endList && tech.preload() !== 'none') {
        segmentLoader.load();
      }
    });
    playlistLoader.on('loadedplaylist', function () {
      segmentLoader.playlist(playlistLoader.media(), requestOptions); // If the player isn't paused, ensure that the segment loader is running

      if (!tech.paused()) {
        segmentLoader.load();
      }
    });
    playlistLoader.on('error', onError[type](type, settings));
  }
};
var initialize = {
  /**
   * Setup PlaylistLoaders and AudioTracks for the audio groups
   *
   * @param {string} type
   *        MediaGroup type
   * @param {Object} settings
   *        Object containing required information for media groups
   * @function initialize.AUDIO
   */
  'AUDIO': function AUDIO(type, settings) {
    var vhs = settings.vhs,
        sourceType = settings.sourceType,
        segmentLoader = settings.segmentLoaders[type],
        requestOptions = settings.requestOptions,
        _settings$master = settings.master,
        mediaGroups = _settings$master.mediaGroups,
        playlists = _settings$master.playlists,
        _settings$mediaTypes$ = settings.mediaTypes[type],
        groups = _settings$mediaTypes$.groups,
        tracks = _settings$mediaTypes$.tracks,
        masterPlaylistLoader = settings.masterPlaylistLoader; // force a default if we have none

    if (!mediaGroups[type] || Object.keys(mediaGroups[type]).length === 0) {
      mediaGroups[type] = {
        main: {
          default: {
            default: true
          }
        }
      };
    }

    var _loop = function _loop(groupId) {
      if (!groups[groupId]) {
        groups[groupId] = [];
      } // List of playlists that have an AUDIO attribute value matching the current
      // group ID


      var groupPlaylists = playlists.filter(function (playlist) {
        return playlist.attributes[type] === groupId;
      });

      var _loop2 = function _loop2(variantLabel) {
        var properties = mediaGroups[type][groupId][variantLabel]; // List of playlists for the current group ID that have a matching uri with
        // this alternate audio variant

        var matchingPlaylists = groupPlaylists.filter(function (playlist) {
          return playlist.resolvedUri === properties.resolvedUri;
        });

        if (matchingPlaylists.length) {
          // If there is a playlist that has the same uri as this audio variant, assume
          // that the playlist is audio only. We delete the resolvedUri property here
          // to prevent a playlist loader from being created so that we don't have
          // both the main and audio segment loaders loading the same audio segments
          // from the same playlist.
          delete properties.resolvedUri;
        }

        var playlistLoader = void 0; // if vhs-json was provided as the source, and the media playlist was resolved,
        // use the resolved media playlist object

        if (sourceType === 'vhs-json' && properties.playlists) {
          playlistLoader = new PlaylistLoader(properties.playlists[0], vhs, requestOptions);
        } else if (properties.resolvedUri) {
          playlistLoader = new PlaylistLoader(properties.resolvedUri, vhs, requestOptions);
        } else if (properties.playlists && sourceType === 'dash') {
          playlistLoader = new DashPlaylistLoader(properties.playlists[0], vhs, requestOptions, masterPlaylistLoader);
        } else {
          // no resolvedUri means the audio is muxed with the video when using this
          // audio track
          playlistLoader = null;
        }

        properties = videojs$1.mergeOptions({
          id: variantLabel,
          playlistLoader: playlistLoader
        }, properties);
        setupListeners[type](type, properties.playlistLoader, settings);
        groups[groupId].push(properties);

        if (typeof tracks[variantLabel] === 'undefined') {
          var track = new videojs$1.AudioTrack({
            id: variantLabel,
            kind: audioTrackKind_(properties),
            enabled: false,
            language: properties.language,
            default: properties.default,
            label: variantLabel
          });
          tracks[variantLabel] = track;
        }
      };

      for (var variantLabel in mediaGroups[type][groupId]) {
        _loop2(variantLabel);
      }
    };

    for (var groupId in mediaGroups[type]) {
      _loop(groupId);
    } // setup single error event handler for the segment loader


    segmentLoader.on('error', onError[type](type, settings));
  },

  /**
   * Setup PlaylistLoaders and TextTracks for the subtitle groups
   *
   * @param {string} type
   *        MediaGroup type
   * @param {Object} settings
   *        Object containing required information for media groups
   * @function initialize.SUBTITLES
   */
  'SUBTITLES': function SUBTITLES(type, settings) {
    var tech = settings.tech,
        vhs = settings.vhs,
        sourceType = settings.sourceType,
        segmentLoader = settings.segmentLoaders[type],
        requestOptions = settings.requestOptions,
        mediaGroups = settings.master.mediaGroups,
        _settings$mediaTypes$2 = settings.mediaTypes[type],
        groups = _settings$mediaTypes$2.groups,
        tracks = _settings$mediaTypes$2.tracks,
        masterPlaylistLoader = settings.masterPlaylistLoader;

    for (var groupId in mediaGroups[type]) {
      if (!groups[groupId]) {
        groups[groupId] = [];
      }

      for (var variantLabel in mediaGroups[type][groupId]) {
        if (mediaGroups[type][groupId][variantLabel].forced) {
          // Subtitle playlists with the forced attribute are not selectable in Safari.
          // According to Apple's HLS Authoring Specification:
          //   If content has forced subtitles and regular subtitles in a given language,
          //   the regular subtitles track in that language MUST contain both the forced
          //   subtitles and the regular subtitles for that language.
          // Because of this requirement and that Safari does not add forced subtitles,
          // forced subtitles are skipped here to maintain consistent experience across
          // all platforms
          continue;
        }

        var properties = mediaGroups[type][groupId][variantLabel];
        var playlistLoader = void 0;

        if (sourceType === 'hls') {
          playlistLoader = new PlaylistLoader(properties.resolvedUri, vhs, requestOptions);
        } else if (sourceType === 'dash') {
          playlistLoader = new DashPlaylistLoader(properties.playlists[0], vhs, requestOptions, masterPlaylistLoader);
        } else if (sourceType === 'vhs-json') {
          playlistLoader = new PlaylistLoader( // if the vhs-json object included the media playlist, use the media playlist
          // as provided, otherwise use the resolved URI to load the playlist
          properties.playlists ? properties.playlists[0] : properties.resolvedUri, vhs, requestOptions);
        }

        properties = videojs$1.mergeOptions({
          id: variantLabel,
          playlistLoader: playlistLoader
        }, properties);
        setupListeners[type](type, properties.playlistLoader, settings);
        groups[groupId].push(properties);

        if (typeof tracks[variantLabel] === 'undefined') {
          var track = tech.addRemoteTextTrack({
            id: variantLabel,
            kind: 'subtitles',
            default: properties.default && properties.autoselect,
            language: properties.language,
            label: variantLabel
          }, false).track;
          tracks[variantLabel] = track;
        }
      }
    } // setup single error event handler for the segment loader


    segmentLoader.on('error', onError[type](type, settings));
  },

  /**
   * Setup TextTracks for the closed-caption groups
   *
   * @param {String} type
   *        MediaGroup type
   * @param {Object} settings
   *        Object containing required information for media groups
   * @function initialize['CLOSED-CAPTIONS']
   */
  'CLOSED-CAPTIONS': function CLOSEDCAPTIONS(type, settings) {
    var tech = settings.tech,
        mediaGroups = settings.master.mediaGroups,
        _settings$mediaTypes$3 = settings.mediaTypes[type],
        groups = _settings$mediaTypes$3.groups,
        tracks = _settings$mediaTypes$3.tracks;

    for (var groupId in mediaGroups[type]) {
      if (!groups[groupId]) {
        groups[groupId] = [];
      }

      for (var variantLabel in mediaGroups[type][groupId]) {
        var properties = mediaGroups[type][groupId][variantLabel]; // We only support CEA608 captions for now, so ignore anything that
        // doesn't use a CCx INSTREAM-ID

        if (!properties.instreamId.match(/CC\d/)) {
          continue;
        } // No PlaylistLoader is required for Closed-Captions because the captions are
        // embedded within the video stream


        groups[groupId].push(videojs$1.mergeOptions({
          id: variantLabel
        }, properties));

        if (typeof tracks[variantLabel] === 'undefined') {
          var track = tech.addRemoteTextTrack({
            id: properties.instreamId,
            kind: 'captions',
            default: properties.default && properties.autoselect,
            language: properties.language,
            label: variantLabel
          }, false).track;
          tracks[variantLabel] = track;
        }
      }
    }
  }
};
/**
 * Returns a function used to get the active group of the provided type
 *
 * @param {string} type
 *        MediaGroup type
 * @param {Object} settings
 *        Object containing required information for media groups
 * @return {Function}
 *         Function that returns the active media group for the provided type. Takes an
 *         optional parameter {TextTrack} track. If no track is provided, a list of all
 *         variants in the group, otherwise the variant corresponding to the provided
 *         track is returned.
 * @function activeGroup
 */

var activeGroup = function activeGroup(type, settings) {
  return function (track) {
    var masterPlaylistLoader = settings.masterPlaylistLoader,
        groups = settings.mediaTypes[type].groups;
    var media = masterPlaylistLoader.media();

    if (!media) {
      return null;
    }

    var variants = null;

    if (media.attributes[type]) {
      variants = groups[media.attributes[type]];
    }

    variants = variants || groups.main;

    if (typeof track === 'undefined') {
      return variants;
    }

    if (track === null) {
      // An active track was specified so a corresponding group is expected. track === null
      // means no track is currently active so there is no corresponding group
      return null;
    }

    return variants.filter(function (props) {
      return props.id === track.id;
    })[0] || null;
  };
};
var activeTrack = {
  /**
   * Returns a function used to get the active track of type provided
   *
   * @param {string} type
   *        MediaGroup type
   * @param {Object} settings
   *        Object containing required information for media groups
   * @return {Function}
   *         Function that returns the active media track for the provided type. Returns
   *         null if no track is active
   * @function activeTrack.AUDIO
   */
  AUDIO: function AUDIO(type, settings) {
    return function () {
      var tracks = settings.mediaTypes[type].tracks;

      for (var id in tracks) {
        if (tracks[id].enabled) {
          return tracks[id];
        }
      }

      return null;
    };
  },

  /**
   * Returns a function used to get the active track of type provided
   *
   * @param {string} type
   *        MediaGroup type
   * @param {Object} settings
   *        Object containing required information for media groups
   * @return {Function}
   *         Function that returns the active media track for the provided type. Returns
   *         null if no track is active
   * @function activeTrack.SUBTITLES
   */
  SUBTITLES: function SUBTITLES(type, settings) {
    return function () {
      var tracks = settings.mediaTypes[type].tracks;

      for (var id in tracks) {
        if (tracks[id].mode === 'showing' || tracks[id].mode === 'hidden') {
          return tracks[id];
        }
      }

      return null;
    };
  }
};
/**
 * Setup PlaylistLoaders and Tracks for media groups (Audio, Subtitles,
 * Closed-Captions) specified in the master manifest.
 *
 * @param {Object} settings
 *        Object containing required information for setting up the media groups
 * @param {Tech} settings.tech
 *        The tech of the player
 * @param {Object} settings.requestOptions
 *        XHR request options used by the segment loaders
 * @param {PlaylistLoader} settings.masterPlaylistLoader
 *        PlaylistLoader for the master source
 * @param {VhsHandler} settings.vhs
 *        VHS SourceHandler
 * @param {Object} settings.master
 *        The parsed master manifest
 * @param {Object} settings.mediaTypes
 *        Object to store the loaders, tracks, and utility methods for each media type
 * @param {Function} settings.blacklistCurrentPlaylist
 *        Blacklists the current rendition and forces a rendition switch.
 * @function setupMediaGroups
 */

var setupMediaGroups = function setupMediaGroups(settings) {
  ['AUDIO', 'SUBTITLES', 'CLOSED-CAPTIONS'].forEach(function (type) {
    initialize[type](type, settings);
  });
  var mediaTypes = settings.mediaTypes,
      masterPlaylistLoader = settings.masterPlaylistLoader,
      tech = settings.tech,
      vhs = settings.vhs; // setup active group and track getters and change event handlers

  ['AUDIO', 'SUBTITLES'].forEach(function (type) {
    mediaTypes[type].activeGroup = activeGroup(type, settings);
    mediaTypes[type].activeTrack = activeTrack[type](type, settings);
    mediaTypes[type].onGroupChanged = onGroupChanged(type, settings);
    mediaTypes[type].onGroupChanging = onGroupChanging(type, settings);
    mediaTypes[type].onTrackChanged = onTrackChanged(type, settings);
  }); // DO NOT enable the default subtitle or caption track.
  // DO enable the default audio track

  var audioGroup = mediaTypes.AUDIO.activeGroup();

  if (audioGroup) {
    var groupId = (audioGroup.filter(function (group) {
      return group.default;
    })[0] || audioGroup[0]).id;
    mediaTypes.AUDIO.tracks[groupId].enabled = true;
    mediaTypes.AUDIO.onTrackChanged();
  }

  masterPlaylistLoader.on('mediachange', function () {
    ['AUDIO', 'SUBTITLES'].forEach(function (type) {
      return mediaTypes[type].onGroupChanged();
    });
  });
  masterPlaylistLoader.on('mediachanging', function () {
    ['AUDIO', 'SUBTITLES'].forEach(function (type) {
      return mediaTypes[type].onGroupChanging();
    });
  }); // custom audio track change event handler for usage event

  var onAudioTrackChanged = function onAudioTrackChanged() {
    mediaTypes.AUDIO.onTrackChanged();
    tech.trigger({
      type: 'usage',
      name: 'vhs-audio-change'
    });
    tech.trigger({
      type: 'usage',
      name: 'hls-audio-change'
    });
  };

  tech.audioTracks().addEventListener('change', onAudioTrackChanged);
  tech.remoteTextTracks().addEventListener('change', mediaTypes.SUBTITLES.onTrackChanged);
  vhs.on('dispose', function () {
    tech.audioTracks().removeEventListener('change', onAudioTrackChanged);
    tech.remoteTextTracks().removeEventListener('change', mediaTypes.SUBTITLES.onTrackChanged);
  }); // clear existing audio tracks and add the ones we just created

  tech.clearTracks('audio');

  for (var id in mediaTypes.AUDIO.tracks) {
    tech.audioTracks().addTrack(mediaTypes.AUDIO.tracks[id]);
  }
};
/**
 * Creates skeleton object used to store the loaders, tracks, and utility methods for each
 * media type
 *
 * @return {Object}
 *         Object to store the loaders, tracks, and utility methods for each media type
 * @function createMediaTypes
 */

var createMediaTypes = function createMediaTypes() {
  var mediaTypes = {};
  ['AUDIO', 'SUBTITLES', 'CLOSED-CAPTIONS'].forEach(function (type) {
    mediaTypes[type] = {
      groups: {},
      tracks: {},
      activePlaylistLoader: null,
      activeGroup: noop,
      activeTrack: noop,
      onGroupChanged: noop,
      onTrackChanged: noop
    };
  });
  return mediaTypes;
};

var ABORT_EARLY_BLACKLIST_SECONDS = 60 * 2;
var Vhs; // SegmentLoader stats that need to have each loader's
// values summed to calculate the final value

var loaderStats = ['mediaRequests', 'mediaRequestsAborted', 'mediaRequestsTimedout', 'mediaRequestsErrored', 'mediaTransferDuration', 'mediaBytesTransferred'];

var sumLoaderStat = function sumLoaderStat(stat) {
  return this.audioSegmentLoader_[stat] + this.mainSegmentLoader_[stat];
};

var shouldSwitchToMedia = function shouldSwitchToMedia(_ref) {
  var currentPlaylist = _ref.currentPlaylist,
      nextPlaylist = _ref.nextPlaylist,
      forwardBuffer = _ref.forwardBuffer,
      bufferLowWaterLine = _ref.bufferLowWaterLine,
      duration = _ref.duration,
      log = _ref.log;

  // we have no other playlist to switch to
  if (!nextPlaylist) {
    videojs$1.log.warn('We received no playlist to switch to. Please check your stream.');
    return false;
  } // If the playlist is live, then we want to not take low water line into account.
  // This is because in LIVE, the player plays 3 segments from the end of the
  // playlist, and if `BUFFER_LOW_WATER_LINE` is greater than the duration availble
  // in those segments, a viewer will never experience a rendition upswitch.


  if (!currentPlaylist.endList) {
    return true;
  } // For the same reason as LIVE, we ignore the low water line when the VOD
  // duration is below the max potential low water line


  if (duration < Config.MAX_BUFFER_LOW_WATER_LINE) {
    return true;
  } // we want to switch down to lower resolutions quickly to continue playback, but


  if (nextPlaylist.attributes.BANDWIDTH < currentPlaylist.attributes.BANDWIDTH) {
    return true;
  } // ensure we have some buffer before we switch up to prevent us running out of
  // buffer while loading a higher rendition.


  if (forwardBuffer >= bufferLowWaterLine) {
    return true;
  }

  return false;
};
/**
 * the master playlist controller controller all interactons
 * between playlists and segmentloaders. At this time this mainly
 * involves a master playlist and a series of audio playlists
 * if they are available
 *
 * @class MasterPlaylistController
 * @extends videojs.EventTarget
 */


var MasterPlaylistController = /*#__PURE__*/function (_videojs$EventTarget) {
  inheritsLoose(MasterPlaylistController, _videojs$EventTarget);

  function MasterPlaylistController(options) {
    var _this;

    _this = _videojs$EventTarget.call(this) || this;
    var src = options.src,
        handleManifestRedirects = options.handleManifestRedirects,
        withCredentials = options.withCredentials,
        tech = options.tech,
        bandwidth = options.bandwidth,
        externVhs = options.externVhs,
        useCueTags = options.useCueTags,
        blacklistDuration = options.blacklistDuration,
        enableLowInitialPlaylist = options.enableLowInitialPlaylist,
        sourceType = options.sourceType,
        cacheEncryptionKeys = options.cacheEncryptionKeys,
        handlePartialData = options.handlePartialData;

    if (!src) {
      throw new Error('A non-empty playlist URL or JSON manifest string is required');
    }

    Vhs = externVhs;
    _this.withCredentials = withCredentials;
    _this.tech_ = tech;
    _this.vhs_ = tech.vhs;
    _this.sourceType_ = sourceType;
    _this.useCueTags_ = useCueTags;
    _this.blacklistDuration = blacklistDuration;
    _this.enableLowInitialPlaylist = enableLowInitialPlaylist;

    if (_this.useCueTags_) {
      _this.cueTagsTrack_ = _this.tech_.addTextTrack('metadata', 'ad-cues');
      _this.cueTagsTrack_.inBandMetadataTrackDispatchType = '';
    }

    _this.requestOptions_ = {
      withCredentials: withCredentials,
      handleManifestRedirects: handleManifestRedirects,
      timeout: null
    };

    _this.on('error', _this.pauseLoading);

    _this.mediaTypes_ = createMediaTypes();
    _this.mediaSource = new window_1.MediaSource();
    _this.handleDurationChange_ = _this.handleDurationChange_.bind(assertThisInitialized(_this));
    _this.handleSourceOpen_ = _this.handleSourceOpen_.bind(assertThisInitialized(_this));
    _this.handleSourceEnded_ = _this.handleSourceEnded_.bind(assertThisInitialized(_this));

    _this.mediaSource.addEventListener('durationchange', _this.handleDurationChange_); // load the media source into the player


    _this.mediaSource.addEventListener('sourceopen', _this.handleSourceOpen_);

    _this.mediaSource.addEventListener('sourceended', _this.handleSourceEnded_); // we don't have to handle sourceclose since dispose will handle termination of
    // everything, and the MediaSource should not be detached without a proper disposal


    _this.seekable_ = videojs$1.createTimeRanges();
    _this.hasPlayed_ = false;
    _this.syncController_ = new SyncController(options);
    _this.segmentMetadataTrack_ = tech.addRemoteTextTrack({
      kind: 'metadata',
      label: 'segment-metadata'
    }, false).track;
    _this.decrypter_ = new Decrypter();
    _this.sourceUpdater_ = new SourceUpdater(_this.mediaSource);
    _this.inbandTextTracks_ = {};
    _this.timelineChangeController_ = new TimelineChangeController();
    var segmentLoaderSettings = {
      vhs: _this.vhs_,
      mediaSource: _this.mediaSource,
      currentTime: _this.tech_.currentTime.bind(_this.tech_),
      seekable: function seekable() {
        return _this.seekable();
      },
      seeking: function seeking() {
        return _this.tech_.seeking();
      },
      duration: function duration() {
        return _this.duration();
      },
      hasPlayed: function hasPlayed() {
        return _this.hasPlayed_;
      },
      goalBufferLength: function goalBufferLength() {
        return _this.goalBufferLength();
      },
      bandwidth: bandwidth,
      syncController: _this.syncController_,
      decrypter: _this.decrypter_,
      sourceType: _this.sourceType_,
      inbandTextTracks: _this.inbandTextTracks_,
      cacheEncryptionKeys: cacheEncryptionKeys,
      handlePartialData: handlePartialData,
      sourceUpdater: _this.sourceUpdater_,
      timelineChangeController: _this.timelineChangeController_
    }; // The source type check not only determines whether a special DASH playlist loader
    // should be used, but also covers the case where the provided src is a vhs-json
    // manifest object (instead of a URL). In the case of vhs-json, the default
    // PlaylistLoader should be used.

    _this.masterPlaylistLoader_ = _this.sourceType_ === 'dash' ? new DashPlaylistLoader(src, _this.vhs_, _this.requestOptions_) : new PlaylistLoader(src, _this.vhs_, _this.requestOptions_);

    _this.setupMasterPlaylistLoaderListeners_(); // setup segment loaders
    // combined audio/video or just video when alternate audio track is selected


    _this.mainSegmentLoader_ = new SegmentLoader(videojs$1.mergeOptions(segmentLoaderSettings, {
      segmentMetadataTrack: _this.segmentMetadataTrack_,
      loaderType: 'main'
    }), options); // alternate audio track

    _this.audioSegmentLoader_ = new SegmentLoader(videojs$1.mergeOptions(segmentLoaderSettings, {
      loaderType: 'audio'
    }), options);
    _this.subtitleSegmentLoader_ = new VTTSegmentLoader(videojs$1.mergeOptions(segmentLoaderSettings, {
      loaderType: 'vtt',
      featuresNativeTextTracks: _this.tech_.featuresNativeTextTracks
    }), options);

    _this.setupSegmentLoaderListeners_(); // Create SegmentLoader stat-getters
    // mediaRequests_
    // mediaRequestsAborted_
    // mediaRequestsTimedout_
    // mediaRequestsErrored_
    // mediaTransferDuration_
    // mediaBytesTransferred_


    loaderStats.forEach(function (stat) {
      _this[stat + '_'] = sumLoaderStat.bind(assertThisInitialized(_this), stat);
    });
    _this.logger_ = logger('MPC');
    _this.triggeredFmp4Usage = false;

    _this.masterPlaylistLoader_.load();

    return _this;
  }
  /**
   * Register event handlers on the master playlist loader. A helper
   * function for construction time.
   *
   * @private
   */


  var _proto = MasterPlaylistController.prototype;

  _proto.setupMasterPlaylistLoaderListeners_ = function setupMasterPlaylistLoaderListeners_() {
    var _this2 = this;

    this.masterPlaylistLoader_.on('loadedmetadata', function () {
      var media = _this2.masterPlaylistLoader_.media();

      var requestTimeout = media.targetDuration * 1.5 * 1000; // If we don't have any more available playlists, we don't want to
      // timeout the request.

      if (isLowestEnabledRendition(_this2.masterPlaylistLoader_.master, _this2.masterPlaylistLoader_.media())) {
        _this2.requestOptions_.timeout = 0;
      } else {
        _this2.requestOptions_.timeout = requestTimeout;
      } // if this isn't a live video and preload permits, start
      // downloading segments


      if (media.endList && _this2.tech_.preload() !== 'none') {
        _this2.mainSegmentLoader_.playlist(media, _this2.requestOptions_);

        _this2.mainSegmentLoader_.load();
      }

      setupMediaGroups({
        sourceType: _this2.sourceType_,
        segmentLoaders: {
          AUDIO: _this2.audioSegmentLoader_,
          SUBTITLES: _this2.subtitleSegmentLoader_,
          main: _this2.mainSegmentLoader_
        },
        tech: _this2.tech_,
        requestOptions: _this2.requestOptions_,
        masterPlaylistLoader: _this2.masterPlaylistLoader_,
        vhs: _this2.vhs_,
        master: _this2.master(),
        mediaTypes: _this2.mediaTypes_,
        blacklistCurrentPlaylist: _this2.blacklistCurrentPlaylist.bind(_this2)
      });

      _this2.triggerPresenceUsage_(_this2.master(), media);

      _this2.setupFirstPlay();

      if (!_this2.mediaTypes_.AUDIO.activePlaylistLoader || _this2.mediaTypes_.AUDIO.activePlaylistLoader.media()) {
        _this2.trigger('selectedinitialmedia');
      } else {
        // We must wait for the active audio playlist loader to
        // finish setting up before triggering this event so the
        // representations API and EME setup is correct
        _this2.mediaTypes_.AUDIO.activePlaylistLoader.one('loadedmetadata', function () {
          _this2.trigger('selectedinitialmedia');
        });
      }
    });
    this.masterPlaylistLoader_.on('loadedplaylist', function () {
      var updatedPlaylist = _this2.masterPlaylistLoader_.media();

      if (!updatedPlaylist) {
        // exclude any variants that are not supported by the browser before selecting
        // an initial media as the playlist selectors do not consider browser support
        _this2.excludeUnsupportedVariants_();

        var selectedMedia;

        if (_this2.enableLowInitialPlaylist) {
          selectedMedia = _this2.selectInitialPlaylist();
        }

        if (!selectedMedia) {
          selectedMedia = _this2.selectPlaylist();
        }

        _this2.initialMedia_ = selectedMedia;

        _this2.masterPlaylistLoader_.media(_this2.initialMedia_); // Under the standard case where a source URL is provided, loadedplaylist will
        // fire again since the playlist will be requested. In the case of vhs-json
        // (where the manifest object is provided as the source), when the media
        // playlist's `segments` list is already available, a media playlist won't be
        // requested, and loadedplaylist won't fire again, so the playlist handler must be
        // called on its own here.


        var haveJsonSource = _this2.sourceType_ === 'vhs-json' && _this2.initialMedia_.segments;

        if (!haveJsonSource) {
          return;
        }

        updatedPlaylist = _this2.initialMedia_;
      }

      _this2.handleUpdatedMediaPlaylist(updatedPlaylist);
    });
    this.masterPlaylistLoader_.on('error', function () {
      _this2.blacklistCurrentPlaylist(_this2.masterPlaylistLoader_.error);
    });
    this.masterPlaylistLoader_.on('mediachanging', function () {
      _this2.mainSegmentLoader_.abort();

      _this2.mainSegmentLoader_.pause();
    });
    this.masterPlaylistLoader_.on('mediachange', function () {
      var media = _this2.masterPlaylistLoader_.media();

      var requestTimeout = media.targetDuration * 1.5 * 1000; // If we don't have any more available playlists, we don't want to
      // timeout the request.

      if (isLowestEnabledRendition(_this2.masterPlaylistLoader_.master, _this2.masterPlaylistLoader_.media())) {
        _this2.requestOptions_.timeout = 0;
      } else {
        _this2.requestOptions_.timeout = requestTimeout;
      } // TODO: Create a new event on the PlaylistLoader that signals
      // that the segments have changed in some way and use that to
      // update the SegmentLoader instead of doing it twice here and
      // on `loadedplaylist`


      _this2.mainSegmentLoader_.playlist(media, _this2.requestOptions_);

      _this2.mainSegmentLoader_.load();

      _this2.tech_.trigger({
        type: 'mediachange',
        bubbles: true
      });
    });
    this.masterPlaylistLoader_.on('playlistunchanged', function () {
      var updatedPlaylist = _this2.masterPlaylistLoader_.media();

      var playlistOutdated = _this2.stuckAtPlaylistEnd_(updatedPlaylist);

      if (playlistOutdated) {
        // Playlist has stopped updating and we're stuck at its end. Try to
        // blacklist it and switch to another playlist in the hope that that
        // one is updating (and give the player a chance to re-adjust to the
        // safe live point).
        _this2.blacklistCurrentPlaylist({
          message: 'Playlist no longer updating.'
        }); // useful for monitoring QoS


        _this2.tech_.trigger('playliststuck');
      }
    });
    this.masterPlaylistLoader_.on('renditiondisabled', function () {
      _this2.tech_.trigger({
        type: 'usage',
        name: 'vhs-rendition-disabled'
      });

      _this2.tech_.trigger({
        type: 'usage',
        name: 'hls-rendition-disabled'
      });
    });
    this.masterPlaylistLoader_.on('renditionenabled', function () {
      _this2.tech_.trigger({
        type: 'usage',
        name: 'vhs-rendition-enabled'
      });

      _this2.tech_.trigger({
        type: 'usage',
        name: 'hls-rendition-enabled'
      });
    });
  }
  /**
   * Given an updated media playlist (whether it was loaded for the first time, or
   * refreshed for live playlists), update any relevant properties and state to reflect
   * changes in the media that should be accounted for (e.g., cues and duration).
   *
   * @param {Object} updatedPlaylist the updated media playlist object
   *
   * @private
   */
  ;

  _proto.handleUpdatedMediaPlaylist = function handleUpdatedMediaPlaylist(updatedPlaylist) {
    if (this.useCueTags_) {
      this.updateAdCues_(updatedPlaylist);
    } // TODO: Create a new event on the PlaylistLoader that signals
    // that the segments have changed in some way and use that to
    // update the SegmentLoader instead of doing it twice here and
    // on `mediachange`


    this.mainSegmentLoader_.playlist(updatedPlaylist, this.requestOptions_);
    this.updateDuration(!updatedPlaylist.endList); // If the player isn't paused, ensure that the segment loader is running,
    // as it is possible that it was temporarily stopped while waiting for
    // a playlist (e.g., in case the playlist errored and we re-requested it).

    if (!this.tech_.paused()) {
      this.mainSegmentLoader_.load();

      if (this.audioSegmentLoader_) {
        this.audioSegmentLoader_.load();
      }
    }
  }
  /**
   * A helper function for triggerring presence usage events once per source
   *
   * @private
   */
  ;

  _proto.triggerPresenceUsage_ = function triggerPresenceUsage_(master, media) {
    var mediaGroups = master.mediaGroups || {};
    var defaultDemuxed = true;
    var audioGroupKeys = Object.keys(mediaGroups.AUDIO);

    for (var mediaGroup in mediaGroups.AUDIO) {
      for (var label in mediaGroups.AUDIO[mediaGroup]) {
        var properties = mediaGroups.AUDIO[mediaGroup][label];

        if (!properties.uri) {
          defaultDemuxed = false;
        }
      }
    }

    if (defaultDemuxed) {
      this.tech_.trigger({
        type: 'usage',
        name: 'vhs-demuxed'
      });
      this.tech_.trigger({
        type: 'usage',
        name: 'hls-demuxed'
      });
    }

    if (Object.keys(mediaGroups.SUBTITLES).length) {
      this.tech_.trigger({
        type: 'usage',
        name: 'vhs-webvtt'
      });
      this.tech_.trigger({
        type: 'usage',
        name: 'hls-webvtt'
      });
    }

    if (Vhs.Playlist.isAes(media)) {
      this.tech_.trigger({
        type: 'usage',
        name: 'vhs-aes'
      });
      this.tech_.trigger({
        type: 'usage',
        name: 'hls-aes'
      });
    }

    if (audioGroupKeys.length && Object.keys(mediaGroups.AUDIO[audioGroupKeys[0]]).length > 1) {
      this.tech_.trigger({
        type: 'usage',
        name: 'vhs-alternate-audio'
      });
      this.tech_.trigger({
        type: 'usage',
        name: 'hls-alternate-audio'
      });
    }

    if (this.useCueTags_) {
      this.tech_.trigger({
        type: 'usage',
        name: 'vhs-playlist-cue-tags'
      });
      this.tech_.trigger({
        type: 'usage',
        name: 'hls-playlist-cue-tags'
      });
    }
  }
  /**
   * Register event handlers on the segment loaders. A helper function
   * for construction time.
   *
   * @private
   */
  ;

  _proto.setupSegmentLoaderListeners_ = function setupSegmentLoaderListeners_() {
    var _this3 = this;

    this.mainSegmentLoader_.on('bandwidthupdate', function () {
      var nextPlaylist = _this3.selectPlaylist();

      var currentPlaylist = _this3.masterPlaylistLoader_.media();

      var buffered = _this3.tech_.buffered();

      var forwardBuffer = buffered.length ? buffered.end(buffered.length - 1) - _this3.tech_.currentTime() : 0;

      var bufferLowWaterLine = _this3.bufferLowWaterLine();

      if (shouldSwitchToMedia({
        currentPlaylist: currentPlaylist,
        nextPlaylist: nextPlaylist,
        forwardBuffer: forwardBuffer,
        bufferLowWaterLine: bufferLowWaterLine,
        duration: _this3.duration(),
        log: _this3.logger_
      })) {
        _this3.masterPlaylistLoader_.media(nextPlaylist);
      }

      _this3.tech_.trigger('bandwidthupdate');
    });
    this.mainSegmentLoader_.on('progress', function () {
      _this3.trigger('progress');
    });
    this.mainSegmentLoader_.on('error', function () {
      _this3.blacklistCurrentPlaylist(_this3.mainSegmentLoader_.error());
    });
    this.mainSegmentLoader_.on('appenderror', function () {
      _this3.error = _this3.mainSegmentLoader_.error_;

      _this3.trigger('error');
    });
    this.mainSegmentLoader_.on('syncinfoupdate', function () {
      _this3.onSyncInfoUpdate_();
    });
    this.mainSegmentLoader_.on('timestampoffset', function () {
      _this3.tech_.trigger({
        type: 'usage',
        name: 'vhs-timestamp-offset'
      });

      _this3.tech_.trigger({
        type: 'usage',
        name: 'hls-timestamp-offset'
      });
    });
    this.audioSegmentLoader_.on('syncinfoupdate', function () {
      _this3.onSyncInfoUpdate_();
    });
    this.audioSegmentLoader_.on('appenderror', function () {
      _this3.error = _this3.audioSegmentLoader_.error_;

      _this3.trigger('error');
    });
    this.mainSegmentLoader_.on('ended', function () {
      _this3.logger_('main segment loader ended');

      _this3.onEndOfStream();
    });
    this.mainSegmentLoader_.on('earlyabort', function () {
      _this3.blacklistCurrentPlaylist({
        message: 'Aborted early because there isn\'t enough bandwidth to complete the ' + 'request without rebuffering.'
      }, ABORT_EARLY_BLACKLIST_SECONDS);
    });

    var updateCodecs = function updateCodecs() {
      if (!_this3.sourceUpdater_.ready()) {
        return _this3.tryToCreateSourceBuffers_();
      }

      var codecs = _this3.getCodecsOrExclude_(); // no codecs means that the playlist was excluded


      if (!codecs) {
        return;
      }

      _this3.sourceUpdater_.addOrChangeSourceBuffers(codecs);
    };

    this.mainSegmentLoader_.on('trackinfo', updateCodecs);
    this.audioSegmentLoader_.on('trackinfo', updateCodecs);
    this.mainSegmentLoader_.on('fmp4', function () {
      if (!_this3.triggeredFmp4Usage) {
        _this3.tech_.trigger({
          type: 'usage',
          name: 'vhs-fmp4'
        });

        _this3.tech_.trigger({
          type: 'usage',
          name: 'hls-fmp4'
        });

        _this3.triggeredFmp4Usage = true;
      }
    });
    this.audioSegmentLoader_.on('fmp4', function () {
      if (!_this3.triggeredFmp4Usage) {
        _this3.tech_.trigger({
          type: 'usage',
          name: 'vhs-fmp4'
        });

        _this3.tech_.trigger({
          type: 'usage',
          name: 'hls-fmp4'
        });

        _this3.triggeredFmp4Usage = true;
      }
    });
    this.audioSegmentLoader_.on('ended', function () {
      _this3.logger_('audioSegmentLoader ended');

      _this3.onEndOfStream();
    });
  };

  _proto.mediaSecondsLoaded_ = function mediaSecondsLoaded_() {
    return Math.max(this.audioSegmentLoader_.mediaSecondsLoaded + this.mainSegmentLoader_.mediaSecondsLoaded);
  }
  /**
   * Call load on our SegmentLoaders
   */
  ;

  _proto.load = function load() {
    this.mainSegmentLoader_.load();

    if (this.mediaTypes_.AUDIO.activePlaylistLoader) {
      this.audioSegmentLoader_.load();
    }

    if (this.mediaTypes_.SUBTITLES.activePlaylistLoader) {
      this.subtitleSegmentLoader_.load();
    }
  }
  /**
   * Re-tune playback quality level for the current player
   * conditions without performing destructive actions, like
   * removing already buffered content
   *
   * @private
   */
  ;

  _proto.smoothQualityChange_ = function smoothQualityChange_(media) {
    if (media === void 0) {
      media = this.selectPlaylist();
    }

    if (media === this.masterPlaylistLoader_.media()) {
      return;
    }

    this.masterPlaylistLoader_.media(media);
    this.mainSegmentLoader_.resetLoader(); // don't need to reset audio as it is reset when media changes
  }
  /**
   * Re-tune playback quality level for the current player
   * conditions. This method will perform destructive actions like removing
   * already buffered content in order to readjust the currently active
   * playlist quickly. This is good for manual quality changes
   *
   * @private
   */
  ;

  _proto.fastQualityChange_ = function fastQualityChange_(media) {
    var _this4 = this;

    if (media === void 0) {
      media = this.selectPlaylist();
    }

    if (media === this.masterPlaylistLoader_.media()) {
      return;
    }

    this.masterPlaylistLoader_.media(media); // Delete all buffered data to allow an immediate quality switch, then seek to give
    // the browser a kick to remove any cached frames from the previous rendtion (.04 seconds
    // ahead is roughly the minimum that will accomplish this across a variety of content
    // in IE and Edge, but seeking in place is sufficient on all other browsers)
    // Edge/IE bug: https://developer.microsoft.com/en-us/microsoft-edge/platform/issues/14600375/
    // Chrome bug: https://bugs.chromium.org/p/chromium/issues/detail?id=651904

    this.mainSegmentLoader_.resetEverything(function () {
      // Since this is not a typical seek, we avoid the seekTo method which can cause segments
      // from the previously enabled rendition to load before the new playlist has finished loading
      if (videojs$1.browser.IE_VERSION || videojs$1.browser.IS_EDGE) {
        _this4.tech_.setCurrentTime(_this4.tech_.currentTime() + 0.04);
      } else {
        _this4.tech_.setCurrentTime(_this4.tech_.currentTime());
      }
    }); // don't need to reset audio as it is reset when media changes
  }
  /**
   * Begin playback.
   */
  ;

  _proto.play = function play() {
    if (this.setupFirstPlay()) {
      return;
    }

    if (this.tech_.ended()) {
      this.tech_.setCurrentTime(0);
    }

    if (this.hasPlayed_) {
      this.load();
    }

    var seekable = this.tech_.seekable(); // if the viewer has paused and we fell out of the live window,
    // seek forward to the live point

    if (this.tech_.duration() === Infinity) {
      if (this.tech_.currentTime() < seekable.start(0)) {
        return this.tech_.setCurrentTime(seekable.end(seekable.length - 1));
      }
    }
  }
  /**
   * Seek to the latest media position if this is a live video and the
   * player and video are loaded and initialized.
   */
  ;

  _proto.setupFirstPlay = function setupFirstPlay() {
    var _this5 = this;

    var media = this.masterPlaylistLoader_.media(); // Check that everything is ready to begin buffering for the first call to play
    //  If 1) there is no active media
    //     2) the player is paused
    //     3) the first play has already been setup
    // then exit early

    if (!media || this.tech_.paused() || this.hasPlayed_) {
      return false;
    } // when the video is a live stream


    if (!media.endList) {
      var seekable = this.seekable();

      if (!seekable.length) {
        // without a seekable range, the player cannot seek to begin buffering at the live
        // point
        return false;
      }

      if (videojs$1.browser.IE_VERSION && this.tech_.readyState() === 0) {
        // IE11 throws an InvalidStateError if you try to set currentTime while the
        // readyState is 0, so it must be delayed until the tech fires loadedmetadata.
        this.tech_.one('loadedmetadata', function () {
          _this5.trigger('firstplay');

          _this5.tech_.setCurrentTime(seekable.end(0));

          _this5.hasPlayed_ = true;
        });
        return false;
      } // trigger firstplay to inform the source handler to ignore the next seek event


      this.trigger('firstplay'); // seek to the live point

      this.tech_.setCurrentTime(seekable.end(0));
    }

    this.hasPlayed_ = true; // we can begin loading now that everything is ready

    this.load();
    return true;
  }
  /**
   * handle the sourceopen event on the MediaSource
   *
   * @private
   */
  ;

  _proto.handleSourceOpen_ = function handleSourceOpen_() {
    // Only attempt to create the source buffer if none already exist.
    // handleSourceOpen is also called when we are "re-opening" a source buffer
    // after `endOfStream` has been called (in response to a seek for instance)
    this.tryToCreateSourceBuffers_(); // if autoplay is enabled, begin playback. This is duplicative of
    // code in video.js but is required because play() must be invoked
    // *after* the media source has opened.

    if (this.tech_.autoplay()) {
      var playPromise = this.tech_.play(); // Catch/silence error when a pause interrupts a play request
      // on browsers which return a promise

      if (typeof playPromise !== 'undefined' && typeof playPromise.then === 'function') {
        playPromise.then(null, function (e) {});
      }
    }

    this.trigger('sourceopen');
  }
  /**
   * handle the sourceended event on the MediaSource
   *
   * @private
   */
  ;

  _proto.handleSourceEnded_ = function handleSourceEnded_() {
    if (!this.inbandTextTracks_.metadataTrack_) {
      return;
    }

    var cues = this.inbandTextTracks_.metadataTrack_.cues;

    if (!cues || !cues.length) {
      return;
    }

    var duration = this.duration();
    cues[cues.length - 1].endTime = isNaN(duration) || Math.abs(duration) === Infinity ? Number.MAX_VALUE : duration;
  }
  /**
   * handle the durationchange event on the MediaSource
   *
   * @private
   */
  ;

  _proto.handleDurationChange_ = function handleDurationChange_() {
    this.tech_.trigger('durationchange');
  }
  /**
   * Calls endOfStream on the media source when all active stream types have called
   * endOfStream
   *
   * @param {string} streamType
   *        Stream type of the segment loader that called endOfStream
   * @private
   */
  ;

  _proto.onEndOfStream = function onEndOfStream() {
    var isEndOfStream = this.mainSegmentLoader_.ended_;

    if (this.mediaTypes_.AUDIO.activePlaylistLoader) {
      // if the audio playlist loader exists, then alternate audio is active
      if (!this.mainSegmentLoader_.currentMediaInfo_ || this.mainSegmentLoader_.currentMediaInfo_.hasVideo) {
        // if we do not know if the main segment loader contains video yet or if we
        // definitively know the main segment loader contains video, then we need to wait
        // for both main and audio segment loaders to call endOfStream
        isEndOfStream = isEndOfStream && this.audioSegmentLoader_.ended_;
      } else {
        // otherwise just rely on the audio loader
        isEndOfStream = this.audioSegmentLoader_.ended_;
      }
    }

    if (!isEndOfStream) {
      return;
    }

    this.sourceUpdater_.endOfStream();
  }
  /**
   * Check if a playlist has stopped being updated
   *
   * @param {Object} playlist the media playlist object
   * @return {boolean} whether the playlist has stopped being updated or not
   */
  ;

  _proto.stuckAtPlaylistEnd_ = function stuckAtPlaylistEnd_(playlist) {
    var seekable = this.seekable();

    if (!seekable.length) {
      // playlist doesn't have enough information to determine whether we are stuck
      return false;
    }

    var expired = this.syncController_.getExpiredTime(playlist, this.duration());

    if (expired === null) {
      return false;
    } // does not use the safe live end to calculate playlist end, since we
    // don't want to say we are stuck while there is still content


    var absolutePlaylistEnd = Vhs.Playlist.playlistEnd(playlist, expired);
    var currentTime = this.tech_.currentTime();
    var buffered = this.tech_.buffered();

    if (!buffered.length) {
      // return true if the playhead reached the absolute end of the playlist
      return absolutePlaylistEnd - currentTime <= SAFE_TIME_DELTA;
    }

    var bufferedEnd = buffered.end(buffered.length - 1); // return true if there is too little buffer left and buffer has reached absolute
    // end of playlist

    return bufferedEnd - currentTime <= SAFE_TIME_DELTA && absolutePlaylistEnd - bufferedEnd <= SAFE_TIME_DELTA;
  }
  /**
   * Blacklists a playlist when an error occurs for a set amount of time
   * making it unavailable for selection by the rendition selection algorithm
   * and then forces a new playlist (rendition) selection.
   *
   * @param {Object=} error an optional error that may include the playlist
   * to blacklist
   * @param {number=} blacklistDuration an optional number of seconds to blacklist the
   * playlist
   */
  ;

  _proto.blacklistCurrentPlaylist = function blacklistCurrentPlaylist(error, blacklistDuration) {
    if (error === void 0) {
      error = {};
    }

    // If the `error` was generated by the playlist loader, it will contain
    // the playlist we were trying to load (but failed) and that should be
    // blacklisted instead of the currently selected playlist which is likely
    // out-of-date in this scenario
    var currentPlaylist = error.playlist || this.masterPlaylistLoader_.media();
    blacklistDuration = blacklistDuration || error.blacklistDuration || this.blacklistDuration; // If there is no current playlist, then an error occurred while we were
    // trying to load the master OR while we were disposing of the tech

    if (!currentPlaylist) {
      this.error = error;

      if (this.mediaSource.readyState !== 'open') {
        this.trigger('error');
      } else {
        this.sourceUpdater_.endOfStream('network');
      }

      return;
    }

    var playlists = this.masterPlaylistLoader_.master.playlists;
    var enabledPlaylists = playlists.filter(isEnabled);
    var isFinalRendition = enabledPlaylists.length === 1 && enabledPlaylists[0] === currentPlaylist; // Don't blacklist the only playlist unless it was blacklisted
    // forever

    if (playlists.length === 1 && blacklistDuration !== Infinity) {
      videojs$1.log.warn("Problem encountered with playlist " + currentPlaylist.id + ". " + 'Trying again since it is the only playlist.');
      this.tech_.trigger('retryplaylist');
      return this.masterPlaylistLoader_.load(isFinalRendition);
    }

    if (isFinalRendition) {
      // Since we're on the final non-blacklisted playlist, and we're about to blacklist
      // it, instead of erring the player or retrying this playlist, clear out the current
      // blacklist. This allows other playlists to be attempted in case any have been
      // fixed.
      var reincluded = false;
      playlists.forEach(function (playlist) {
        // skip current playlist which is about to be blacklisted
        if (playlist === currentPlaylist) {
          return;
        }

        var excludeUntil = playlist.excludeUntil; // a playlist cannot be reincluded if it wasn't excluded to begin with.

        if (typeof excludeUntil !== 'undefined' && excludeUntil !== Infinity) {
          reincluded = true;
          delete playlist.excludeUntil;
        }
      });

      if (reincluded) {
        videojs$1.log.warn('Removing other playlists from the exclusion list because the last ' + 'rendition is about to be excluded.'); // Technically we are retrying a playlist, in that we are simply retrying a previous
        // playlist. This is needed for users relying on the retryplaylist event to catch a
        // case where the player might be stuck and looping through "dead" playlists.

        this.tech_.trigger('retryplaylist');
      }
    } // Blacklist this playlist


    currentPlaylist.excludeUntil = Date.now() + blacklistDuration * 1000;
    this.tech_.trigger('blacklistplaylist');
    this.tech_.trigger({
      type: 'usage',
      name: 'vhs-rendition-blacklisted'
    });
    this.tech_.trigger({
      type: 'usage',
      name: 'hls-rendition-blacklisted'
    }); // TODO: should we select a new playlist if this blacklist wasn't for the currentPlaylist?
    // Would be something like media().id !=== currentPlaylist.id and we  would need something
    // like `pendingMedia` in playlist loaders to check against that too. This will prevent us
    // from loading a new playlist on any blacklist.
    // Select a new playlist

    var nextPlaylist = this.selectPlaylist();

    if (!nextPlaylist) {
      this.error = 'Playback cannot continue. No available working or supported playlists.';
      this.trigger('error');
      return;
    }

    var logFn = error.internal ? this.logger_ : videojs$1.log.warn;
    var errorMessage = error.message ? ' ' + error.message : '';
    logFn((error.internal ? 'Internal problem' : 'Problem') + " encountered with playlist " + currentPlaylist.id + "." + (errorMessage + " Switching to playlist " + nextPlaylist.id + ".")); // if audio group changed reset audio loaders

    if (nextPlaylist.attributes.AUDIO !== currentPlaylist.attributes.AUDIO) {
      this.delegateLoaders_('audio', ['abort', 'pause']);
    } // if subtitle group changed reset subtitle loaders


    if (nextPlaylist.attributes.SUBTITLES !== currentPlaylist.attributes.SUBTITLES) {
      this.delegateLoaders_('subtitle', ['abort', 'pause']);
    }

    this.delegateLoaders_('main', ['abort', 'pause']);
    return this.masterPlaylistLoader_.media(nextPlaylist, isFinalRendition);
  }
  /**
   * Pause all segment/playlist loaders
   */
  ;

  _proto.pauseLoading = function pauseLoading() {
    this.delegateLoaders_('all', ['abort', 'pause']);
  }
  /**
   * Call a set of functions in order on playlist loaders, segment loaders,
   * or both types of loaders.
   *
   * @param {string} filter
   *        Filter loaders that should call fnNames using a string. Can be:
   *        * all - run on all loaders
   *        * audio - run on all audio loaders
   *        * subtitle - run on all subtitle loaders
   *        * main - run on the main/master loaders
   *
   * @param {Array|string} fnNames
   *        A string or array of function names to call.
   */
  ;

  _proto.delegateLoaders_ = function delegateLoaders_(filter, fnNames) {
    var _this6 = this;

    var loaders = [];
    var dontFilterPlaylist = filter === 'all';

    if (dontFilterPlaylist || filter === 'main') {
      loaders.push(this.masterPlaylistLoader_);
    }

    var mediaTypes = [];

    if (dontFilterPlaylist || filter === 'audio') {
      mediaTypes.push('AUDIO');
    }

    if (dontFilterPlaylist || filter === 'subtitle') {
      mediaTypes.push('CLOSED-CAPTIONS');
      mediaTypes.push('SUBTITLES');
    }

    mediaTypes.forEach(function (mediaType) {
      var loader = _this6.mediaTypes_[mediaType] && _this6.mediaTypes_[mediaType].activePlaylistLoader;

      if (loader) {
        loaders.push(loader);
      }
    });
    ['main', 'audio', 'subtitle'].forEach(function (name) {
      var loader = _this6[name + "SegmentLoader_"];

      if (loader && (filter === name || filter === 'all')) {
        loaders.push(loader);
      }
    });
    loaders.forEach(function (loader) {
      return fnNames.forEach(function (fnName) {
        if (typeof loader[fnName] === 'function') {
          loader[fnName]();
        }
      });
    });
  }
  /**
   * set the current time on all segment loaders
   *
   * @param {TimeRange} currentTime the current time to set
   * @return {TimeRange} the current time
   */
  ;

  _proto.setCurrentTime = function setCurrentTime(currentTime) {
    var buffered = findRange(this.tech_.buffered(), currentTime);

    if (!(this.masterPlaylistLoader_ && this.masterPlaylistLoader_.media())) {
      // return immediately if the metadata is not ready yet
      return 0;
    } // it's clearly an edge-case but don't thrown an error if asked to
    // seek within an empty playlist


    if (!this.masterPlaylistLoader_.media().segments) {
      return 0;
    } // if the seek location is already buffered, continue buffering as usual


    if (buffered && buffered.length) {
      return currentTime;
    } // cancel outstanding requests so we begin buffering at the new
    // location


    this.mainSegmentLoader_.resetEverything();
    this.mainSegmentLoader_.abort();

    if (this.mediaTypes_.AUDIO.activePlaylistLoader) {
      this.audioSegmentLoader_.resetEverything();
      this.audioSegmentLoader_.abort();
    }

    if (this.mediaTypes_.SUBTITLES.activePlaylistLoader) {
      this.subtitleSegmentLoader_.resetEverything();
      this.subtitleSegmentLoader_.abort();
    } // start segment loader loading in case they are paused


    this.load();
  }
  /**
   * get the current duration
   *
   * @return {TimeRange} the duration
   */
  ;

  _proto.duration = function duration() {
    if (!this.masterPlaylistLoader_) {
      return 0;
    }

    var media = this.masterPlaylistLoader_.media();

    if (!media) {
      // no playlists loaded yet, so can't determine a duration
      return 0;
    } // Don't rely on the media source for duration in the case of a live playlist since
    // setting the native MediaSource's duration to infinity ends up with consequences to
    // seekable behavior. See https://github.com/w3c/media-source/issues/5 for details.
    //
    // This is resolved in the spec by https://github.com/w3c/media-source/pull/92,
    // however, few browsers have support for setLiveSeekableRange()
    // https://developer.mozilla.org/en-US/docs/Web/API/MediaSource/setLiveSeekableRange
    //
    // Until a time when the duration of the media source can be set to infinity, and a
    // seekable range specified across browsers, just return Infinity.


    if (!media.endList) {
      return Infinity;
    } // Since this is a VOD video, it is safe to rely on the media source's duration (if
    // available). If it's not available, fall back to a playlist-calculated estimate.


    if (this.mediaSource) {
      return this.mediaSource.duration;
    }

    return Vhs.Playlist.duration(media);
  }
  /**
   * check the seekable range
   *
   * @return {TimeRange} the seekable range
   */
  ;

  _proto.seekable = function seekable() {
    return this.seekable_;
  };

  _proto.onSyncInfoUpdate_ = function onSyncInfoUpdate_() {
    var audioSeekable;

    if (!this.masterPlaylistLoader_) {
      return;
    }

    var media = this.masterPlaylistLoader_.media();

    if (!media) {
      return;
    }

    var expired = this.syncController_.getExpiredTime(media, this.duration());

    if (expired === null) {
      // not enough information to update seekable
      return;
    }

    var suggestedPresentationDelay = this.masterPlaylistLoader_.master.suggestedPresentationDelay;
    var mainSeekable = Vhs.Playlist.seekable(media, expired, suggestedPresentationDelay);

    if (mainSeekable.length === 0) {
      return;
    }

    if (this.mediaTypes_.AUDIO.activePlaylistLoader) {
      media = this.mediaTypes_.AUDIO.activePlaylistLoader.media();
      expired = this.syncController_.getExpiredTime(media, this.duration());

      if (expired === null) {
        return;
      }

      audioSeekable = Vhs.Playlist.seekable(media, expired, suggestedPresentationDelay);

      if (audioSeekable.length === 0) {
        return;
      }
    }

    var oldEnd;
    var oldStart;

    if (this.seekable_ && this.seekable_.length) {
      oldEnd = this.seekable_.end(0);
      oldStart = this.seekable_.start(0);
    }

    if (!audioSeekable) {
      // seekable has been calculated based on buffering video data so it
      // can be returned directly
      this.seekable_ = mainSeekable;
    } else if (audioSeekable.start(0) > mainSeekable.end(0) || mainSeekable.start(0) > audioSeekable.end(0)) {
      // seekables are pretty far off, rely on main
      this.seekable_ = mainSeekable;
    } else {
      this.seekable_ = videojs$1.createTimeRanges([[audioSeekable.start(0) > mainSeekable.start(0) ? audioSeekable.start(0) : mainSeekable.start(0), audioSeekable.end(0) < mainSeekable.end(0) ? audioSeekable.end(0) : mainSeekable.end(0)]]);
    } // seekable is the same as last time


    if (this.seekable_ && this.seekable_.length) {
      if (this.seekable_.end(0) === oldEnd && this.seekable_.start(0) === oldStart) {
        return;
      }
    }

    this.logger_("seekable updated [" + printableRange(this.seekable_) + "]");
    this.tech_.trigger('seekablechanged');
  }
  /**
   * Update the player duration
   */
  ;

  _proto.updateDuration = function updateDuration(isLive) {
    if (this.updateDuration_) {
      this.mediaSource.removeEventListener('sourceopen', this.updateDuration_);
      this.updateDuration_ = null;
    }

    if (this.mediaSource.readyState !== 'open') {
      this.updateDuration_ = this.updateDuration.bind(this, isLive);
      this.mediaSource.addEventListener('sourceopen', this.updateDuration_);
      return;
    }

    if (isLive) {
      var seekable = this.seekable();

      if (!seekable.length) {
        return;
      } // Even in the case of a live playlist, the native MediaSource's duration should not
      // be set to Infinity (even though this would be expected for a live playlist), since
      // setting the native MediaSource's duration to infinity ends up with consequences to
      // seekable behavior. See https://github.com/w3c/media-source/issues/5 for details.
      //
      // This is resolved in the spec by https://github.com/w3c/media-source/pull/92,
      // however, few browsers have support for setLiveSeekableRange()
      // https://developer.mozilla.org/en-US/docs/Web/API/MediaSource/setLiveSeekableRange
      //
      // Until a time when the duration of the media source can be set to infinity, and a
      // seekable range specified across browsers, the duration should be greater than or
      // equal to the last possible seekable value.
      // MediaSource duration starts as NaN
      // It is possible (and probable) that this case will never be reached for many
      // sources, since the MediaSource reports duration as the highest value without
      // accounting for timestamp offset. For example, if the timestamp offset is -100 and
      // we buffered times 0 to 100 with real times of 100 to 200, even though current
      // time will be between 0 and 100, the native media source may report the duration
      // as 200. However, since we report duration separate from the media source (as
      // Infinity), and as long as the native media source duration value is greater than
      // our reported seekable range, seeks will work as expected. The large number as
      // duration for live is actually a strategy used by some players to work around the
      // issue of live seekable ranges cited above.


      if (isNaN(this.mediaSource.duration) || this.mediaSource.duration < seekable.end(seekable.length - 1)) {
        this.sourceUpdater_.setDuration(seekable.end(seekable.length - 1));
      }

      return;
    }

    var buffered = this.tech_.buffered();
    var duration = Vhs.Playlist.duration(this.masterPlaylistLoader_.media());

    if (buffered.length > 0) {
      duration = Math.max(duration, buffered.end(buffered.length - 1));
    }

    if (this.mediaSource.duration !== duration) {
      this.sourceUpdater_.setDuration(duration);
    }
  }
  /**
   * dispose of the MasterPlaylistController and everything
   * that it controls
   */
  ;

  _proto.dispose = function dispose() {
    var _this7 = this;

    this.trigger('dispose');
    this.decrypter_.terminate();
    this.masterPlaylistLoader_.dispose();
    this.mainSegmentLoader_.dispose();
    ['AUDIO', 'SUBTITLES'].forEach(function (type) {
      var groups = _this7.mediaTypes_[type].groups;

      for (var id in groups) {
        groups[id].forEach(function (group) {
          if (group.playlistLoader) {
            group.playlistLoader.dispose();
          }
        });
      }
    });
    this.audioSegmentLoader_.dispose();
    this.subtitleSegmentLoader_.dispose();
    this.sourceUpdater_.dispose();
    this.timelineChangeController_.dispose();

    if (this.updateDuration_) {
      this.mediaSource.removeEventListener('sourceopen', this.updateDuration_);
    }

    this.mediaSource.removeEventListener('durationchange', this.handleDurationChange_); // load the media source into the player

    this.mediaSource.removeEventListener('sourceopen', this.handleSourceOpen_);
    this.mediaSource.removeEventListener('sourceended', this.handleSourceEnded_);
    this.off();
  }
  /**
   * return the master playlist object if we have one
   *
   * @return {Object} the master playlist object that we parsed
   */
  ;

  _proto.master = function master() {
    return this.masterPlaylistLoader_.master;
  }
  /**
   * return the currently selected playlist
   *
   * @return {Object} the currently selected playlist object that we parsed
   */
  ;

  _proto.media = function media() {
    // playlist loader will not return media if it has not been fully loaded
    return this.masterPlaylistLoader_.media() || this.initialMedia_;
  };

  _proto.areMediaTypesKnown_ = function areMediaTypesKnown_() {
    var usingAudioLoader = !!this.mediaTypes_.AUDIO.activePlaylistLoader; // one or both loaders has not loaded sufficently to get codecs

    if (!this.mainSegmentLoader_.currentMediaInfo_ || usingAudioLoader && !this.audioSegmentLoader_.currentMediaInfo_) {
      return false;
    }

    return true;
  };

  _proto.getCodecsOrExclude_ = function getCodecsOrExclude_() {
    var _this8 = this;

    var media = {
      main: this.mainSegmentLoader_.currentMediaInfo_ || {},
      audio: this.audioSegmentLoader_.currentMediaInfo_ || {}
    }; // set "main" media equal to video

    media.video = media.main;
    var playlistCodecs = codecsForPlaylist(this.master(), this.media());
    var codecs$1 = {};
    var usingAudioLoader = !!this.mediaTypes_.AUDIO.activePlaylistLoader;

    if (media.main.hasVideo) {
      codecs$1.video = playlistCodecs.video || media.main.videoCodec || codecs.DEFAULT_VIDEO_CODEC;
    }

    if (media.main.isMuxed) {
      codecs$1.video += "," + (playlistCodecs.audio || media.main.audioCodec || codecs.DEFAULT_AUDIO_CODEC);
    }

    if (media.main.hasAudio && !media.main.isMuxed || media.audio.hasAudio || usingAudioLoader) {
      codecs$1.audio = playlistCodecs.audio || media.main.audioCodec || media.audio.audioCodec || codecs.DEFAULT_AUDIO_CODEC; // set audio isFmp4 so we use the correct "supports" function below

      media.audio.isFmp4 = media.main.hasAudio && !media.main.isMuxed ? media.main.isFmp4 : media.audio.isFmp4;
    } // no codecs, no playback.


    if (!codecs$1.audio && !codecs$1.video) {
      this.blacklistCurrentPlaylist({
        playlist: this.media(),
        message: 'Could not determine codecs for playlist.',
        blacklistDuration: Infinity
      });
      return;
    } // fmp4 relies on browser support, while ts relies on muxer support


    var supportFunction = function supportFunction(isFmp4, codec) {
      return isFmp4 ? codecs.browserSupportsCodec(codec) : codecs.muxerSupportsCodec(codec);
    };

    var unsupportedCodecs = {};
    var unsupportedAudio;
    ['video', 'audio'].forEach(function (type) {
      if (codecs$1.hasOwnProperty(type) && !supportFunction(media[type].isFmp4, codecs$1[type])) {
        var supporter = media[type].isFmp4 ? 'browser' : 'muxer';
        unsupportedCodecs[supporter] = unsupportedCodecs[supporter] || [];
        unsupportedCodecs[supporter].push(codecs$1[type]);

        if (type === 'audio') {
          unsupportedAudio = supporter;
        }
      }
    });

    if (usingAudioLoader && unsupportedAudio && this.media().attributes.AUDIO) {
      var audioGroup = this.media().attributes.AUDIO;
      this.master().playlists.forEach(function (variant) {
        var variantAudioGroup = variant.attributes && variant.attributes.AUDIO;

        if (variantAudioGroup === audioGroup && variant !== _this8.media()) {
          variant.excludeUntil = Infinity;
        }
      });
      this.logger_("excluding audio group " + audioGroup + " as " + unsupportedAudio + " does not support codec(s): \"" + codecs$1.audio + "\"");
    } // if we have any unsupported codecs blacklist this playlist.


    if (Object.keys(unsupportedCodecs).length) {
      var message = Object.keys(unsupportedCodecs).reduce(function (acc, supporter) {
        if (acc) {
          acc += ', ';
        }

        acc += supporter + " does not support codec(s): \"" + unsupportedCodecs[supporter].join(',') + "\"";
        return acc;
      }, '') + '.';
      this.blacklistCurrentPlaylist({
        playlist: this.media(),
        internal: true,
        message: message,
        blacklistDuration: Infinity
      });
      return;
    } // check if codec switching is happening


    if (this.sourceUpdater_.ready() && !this.sourceUpdater_.canChangeType()) {
      var switchMessages = [];
      ['video', 'audio'].forEach(function (type) {
        var newCodec = (codecs.parseCodecs(_this8.sourceUpdater_.codecs[type] || '')[type] || {}).type;
        var oldCodec = (codecs.parseCodecs(codecs$1[type] || '')[type] || {}).type;

        if (newCodec && oldCodec && newCodec.toLowerCase() !== oldCodec.toLowerCase()) {
          switchMessages.push("\"" + _this8.sourceUpdater_.codecs[type] + "\" -> \"" + codecs$1[type] + "\"");
        }
      });

      if (switchMessages.length) {
        this.blacklistCurrentPlaylist({
          playlist: this.media(),
          message: "Codec switching not supported: " + switchMessages.join(', ') + ".",
          blacklistDuration: Infinity,
          internal: true
        });
        return;
      }
    } // TODO: when using the muxer shouldn't we just return
    // the codecs that the muxer outputs?


    return codecs$1;
  }
  /**
   * Create source buffers and exlude any incompatible renditions.
   *
   * @private
   */
  ;

  _proto.tryToCreateSourceBuffers_ = function tryToCreateSourceBuffers_() {
    // media source is not ready yet or sourceBuffers are already
    // created.
    if (this.mediaSource.readyState !== 'open' || this.sourceUpdater_.ready()) {
      return;
    }

    if (!this.areMediaTypesKnown_()) {
      return;
    }

    var codecs = this.getCodecsOrExclude_(); // no codecs means that the playlist was excluded

    if (!codecs) {
      return;
    }

    this.sourceUpdater_.createSourceBuffers(codecs);
    var codecString = [codecs.video, codecs.audio].filter(Boolean).join(',');
    this.excludeIncompatibleVariants_(codecString);
  }
  /**
   * Excludes playlists with codecs that are unsupported by the muxer and browser.
   */
  ;

  _proto.excludeUnsupportedVariants_ = function excludeUnsupportedVariants_() {
    var _this9 = this;

    this.master().playlists.forEach(function (variant) {
      var codecs$1 = codecsForPlaylist(_this9.master, variant);

      if (codecs$1.audio && !codecs.muxerSupportsCodec(codecs$1.audio) && !codecs.browserSupportsCodec(codecs$1.audio)) {
        variant.excludeUntil = Infinity;
      }

      if (codecs$1.video && !codecs.muxerSupportsCodec(codecs$1.video) && !codecs.browserSupportsCodec(codecs$1.video)) {
        variant.excludeUntil = Infinity;
      }
    });
  }
  /**
   * Blacklist playlists that are known to be codec or
   * stream-incompatible with the SourceBuffer configuration. For
   * instance, Media Source Extensions would cause the video element to
   * stall waiting for video data if you switched from a variant with
   * video and audio to an audio-only one.
   *
   * @param {Object} media a media playlist compatible with the current
   * set of SourceBuffers. Variants in the current master playlist that
   * do not appear to have compatible codec or stream configurations
   * will be excluded from the default playlist selection algorithm
   * indefinitely.
   * @private
   */
  ;

  _proto.excludeIncompatibleVariants_ = function excludeIncompatibleVariants_(codecString) {
    var _this10 = this;

    var codecs$1 = codecs.parseCodecs(codecString);
    var codecCount = Object.keys(codecs$1).length;
    this.master().playlists.forEach(function (variant) {
      // skip variants that are already blacklisted forever
      if (variant.excludeUntil === Infinity) {
        return;
      }
      /* TODO: Decide whether two codecs should be assumed here.
       * Right now, for playlists that don't specify codecs, VHS assumes
       * that there are two (one for audio and one for video).
       * Although this is often the case, this may lead to broken behavior
       * if the playlist only has one codec. It may be better in the future
       * to decide at time of segment download how many tracks there are and
       * determine the proper codecs. This will come at a cost of potentially
       * more bandwidth, but will be a more robust approach than the assumption here.
       */


      var variantCodecs = {};
      var variantCodecCount = 2;
      var blacklistReasons = []; // get codecs from the playlist for this variant

      var variantCodecStrings = codecsForPlaylist(_this10.masterPlaylistLoader_.master, variant);

      if (variantCodecStrings.audio || variantCodecStrings.video) {
        var variantCodecString = [variantCodecStrings.video, variantCodecStrings.audio].filter(Boolean).join(',');
        variantCodecs = codecs.parseCodecs(variantCodecString);
        variantCodecCount = Object.keys(variantCodecs).length;
      } // TODO: we can support this by removing the
      // old media source and creating a new one, but it will take some work.
      // The number of streams cannot change


      if (variantCodecCount !== codecCount) {
        blacklistReasons.push("codec count \"" + variantCodecCount + "\" !== \"" + codecCount + "\"");
        variant.excludeUntil = Infinity;
      } // only exclude playlists by codec change, if codecs cannot switch
      // during playback.


      if (!_this10.sourceUpdater_.canChangeType()) {
        // the video codec cannot change
        if (variantCodecs.video && codecs$1.video && variantCodecs.video.type.toLowerCase() !== codecs$1.video.type.toLowerCase()) {
          blacklistReasons.push("video codec \"" + variantCodecs.video.type + "\" !== \"" + codecs$1.video.type + "\"");
          variant.excludeUntil = Infinity;
        } // the audio codec cannot change


        if (variantCodecs.audio && codecs$1.audio && variantCodecs.audio.type.toLowerCase() !== codecs$1.audio.type.toLowerCase()) {
          variant.excludeUntil = Infinity;
          blacklistReasons.push("audio codec \"" + variantCodecs.audio.type + "\" !== \"" + codecs$1.audio.type + "\"");
        }
      }

      if (blacklistReasons.length) {
        _this10.logger_("blacklisting " + variant.id + ": " + blacklistReasons.join(' && '));
      }
    });
  };

  _proto.updateAdCues_ = function updateAdCues_(media) {
    var offset = 0;
    var seekable = this.seekable();

    if (seekable.length) {
      offset = seekable.start(0);
    }

    updateAdCues(media, this.cueTagsTrack_, offset);
  }
  /**
   * Calculates the desired forward buffer length based on current time
   *
   * @return {number} Desired forward buffer length in seconds
   */
  ;

  _proto.goalBufferLength = function goalBufferLength() {
    var currentTime = this.tech_.currentTime();
    var initial = Config.GOAL_BUFFER_LENGTH;
    var rate = Config.GOAL_BUFFER_LENGTH_RATE;
    var max = Math.max(initial, Config.MAX_GOAL_BUFFER_LENGTH);
    return Math.min(initial + currentTime * rate, max);
  }
  /**
   * Calculates the desired buffer low water line based on current time
   *
   * @return {number} Desired buffer low water line in seconds
   */
  ;

  _proto.bufferLowWaterLine = function bufferLowWaterLine() {
    var currentTime = this.tech_.currentTime();
    var initial = Config.BUFFER_LOW_WATER_LINE;
    var rate = Config.BUFFER_LOW_WATER_LINE_RATE;
    var max = Math.max(initial, Config.MAX_BUFFER_LOW_WATER_LINE);
    return Math.min(initial + currentTime * rate, max);
  };

  return MasterPlaylistController;
}(videojs$1.EventTarget);

/**
 * Returns a function that acts as the Enable/disable playlist function.
 *
 * @param {PlaylistLoader} loader - The master playlist loader
 * @param {string} playlistID - id of the playlist
 * @param {Function} changePlaylistFn - A function to be called after a
 * playlist's enabled-state has been changed. Will NOT be called if a
 * playlist's enabled-state is unchanged
 * @param {boolean=} enable - Value to set the playlist enabled-state to
 * or if undefined returns the current enabled-state for the playlist
 * @return {Function} Function for setting/getting enabled
 */

var enableFunction = function enableFunction(loader, playlistID, changePlaylistFn) {
  return function (enable) {
    var playlist = loader.master.playlists[playlistID];
    var incompatible = isIncompatible(playlist);
    var currentlyEnabled = isEnabled(playlist);

    if (typeof enable === 'undefined') {
      return currentlyEnabled;
    }

    if (enable) {
      delete playlist.disabled;
    } else {
      playlist.disabled = true;
    }

    if (enable !== currentlyEnabled && !incompatible) {
      // Ensure the outside world knows about our changes
      changePlaylistFn();

      if (enable) {
        loader.trigger('renditionenabled');
      } else {
        loader.trigger('renditiondisabled');
      }
    }

    return enable;
  };
};
/**
 * The representation object encapsulates the publicly visible information
 * in a media playlist along with a setter/getter-type function (enabled)
 * for changing the enabled-state of a particular playlist entry
 *
 * @class Representation
 */


var Representation = function Representation(vhsHandler, playlist, id) {
  var mpc = vhsHandler.masterPlaylistController_,
      smoothQualityChange = vhsHandler.options_.smoothQualityChange; // Get a reference to a bound version of the quality change function

  var changeType = smoothQualityChange ? 'smooth' : 'fast';
  var qualityChangeFunction = mpc[changeType + "QualityChange_"].bind(mpc); // some playlist attributes are optional

  if (playlist.attributes.RESOLUTION) {
    var resolution = playlist.attributes.RESOLUTION;
    this.width = resolution.width;
    this.height = resolution.height;
  }

  this.bandwidth = playlist.attributes.BANDWIDTH;
  this.codecs = codecsForPlaylist(mpc.master(), playlist);
  this.playlist = playlist; // The id is simply the ordinality of the media playlist
  // within the master playlist

  this.id = id; // Partially-apply the enableFunction to create a playlist-
  // specific variant

  this.enabled = enableFunction(vhsHandler.playlists, playlist.id, qualityChangeFunction);
};
/**
 * A mixin function that adds the `representations` api to an instance
 * of the VhsHandler class
 *
 * @param {VhsHandler} vhsHandler - An instance of VhsHandler to add the
 * representation API into
 */


var renditionSelectionMixin = function renditionSelectionMixin(vhsHandler) {
  var playlists = vhsHandler.playlists; // Add a single API-specific function to the VhsHandler instance

  vhsHandler.representations = function () {
    if (!playlists || !playlists.master || !playlists.master.playlists) {
      return [];
    }

    return playlists.master.playlists.filter(function (media) {
      return !isIncompatible(media);
    }).map(function (e, i) {
      return new Representation(vhsHandler, e, e.id);
    });
  };
};

/**
 * @file playback-watcher.js
 *
 * Playback starts, and now my watch begins. It shall not end until my death. I shall
 * take no wait, hold no uncleared timeouts, father no bad seeks. I shall wear no crowns
 * and win no glory. I shall live and die at my post. I am the corrector of the underflow.
 * I am the watcher of gaps. I am the shield that guards the realms of seekable. I pledge
 * my life and honor to the Playback Watch, for this Player and all the Players to come.
 */

var timerCancelEvents = ['seeking', 'seeked', 'pause', 'playing', 'error'];
/**
 * Returns whether or not the current time should be considered close to buffered content,
 * taking into consideration whether there's enough buffered content for proper playback.
 *
 * @param {Object} options
 *        Options object
 * @param {TimeRange} options.buffered
 *        Current buffer
 * @param {number} options.targetDuration
 *        The active playlist's target duration
 * @param {number} options.currentTime
 *        The current time of the player
 * @return {boolean}
 *         Whether the current time should be considered close to the buffer
 */

var closeToBufferedContent = function closeToBufferedContent(_ref) {
  var buffered = _ref.buffered,
      targetDuration = _ref.targetDuration,
      currentTime = _ref.currentTime;

  if (!buffered.length) {
    return false;
  } // At least two to three segments worth of content should be buffered before there's a
  // full enough buffer to consider taking any actions.


  if (buffered.end(0) - buffered.start(0) < targetDuration * 2) {
    return false;
  } // It's possible that, on seek, a remove hasn't completed and the buffered range is
  // somewhere past the current time. In that event, don't consider the buffered content
  // close.


  if (currentTime > buffered.start(0)) {
    return false;
  } // Since target duration generally represents the max (or close to max) duration of a
  // segment, if the buffer is within a segment of the current time, the gap probably
  // won't be closed, and current time should be considered close to buffered content.


  return buffered.start(0) - currentTime < targetDuration;
};
/**
 * @class PlaybackWatcher
 */

var PlaybackWatcher = /*#__PURE__*/function () {
  /**
   * Represents an PlaybackWatcher object.
   *
   * @class
   * @param {Object} options an object that includes the tech and settings
   */
  function PlaybackWatcher(options) {
    var _this = this;

    this.masterPlaylistController_ = options.masterPlaylistController;
    this.tech_ = options.tech;
    this.seekable = options.seekable;
    this.allowSeeksWithinUnsafeLiveWindow = options.allowSeeksWithinUnsafeLiveWindow;
    this.media = options.media;
    this.consecutiveUpdates = 0;
    this.lastRecordedTime = null;
    this.timer_ = null;
    this.checkCurrentTimeTimeout_ = null;
    this.logger_ = logger('PlaybackWatcher');
    this.logger_('initialize');

    var canPlayHandler = function canPlayHandler() {
      return _this.monitorCurrentTime_();
    };

    var waitingHandler = function waitingHandler() {
      return _this.techWaiting_();
    };

    var cancelTimerHandler = function cancelTimerHandler() {
      return _this.cancelTimer_();
    };

    var fixesBadSeeksHandler = function fixesBadSeeksHandler() {
      return _this.fixesBadSeeks_();
    };

    var mpc = this.masterPlaylistController_;
    var loaderTypes = ['main', 'subtitle', 'audio'];
    var loaderChecks = {};
    loaderTypes.forEach(function (type) {
      loaderChecks[type] = {
        reset: function reset() {
          return _this.resetSegmentDownloads_(type);
        },
        updateend: function updateend() {
          return _this.checkSegmentDownloads_(type);
        }
      };
      mpc[type + "SegmentLoader_"].on('appendsdone', loaderChecks[type].updateend); // If a rendition switch happens during a playback stall where the buffer
      // isn't changing we want to reset. We cannot assume that the new rendition
      // will also be stalled, until after new appends.

      mpc[type + "SegmentLoader_"].on('playlistupdate', loaderChecks[type].reset); // Playback stalls should not be detected right after seeking.
      // This prevents one segment playlists (single vtt or single segment content)
      // from being detected as stalling. As the buffer will not change in those cases, since
      // the buffer is the entire video duration.

      _this.tech_.on(['seeked', 'seeking'], loaderChecks[type].reset);
    });
    this.tech_.on('seekablechanged', fixesBadSeeksHandler);
    this.tech_.on('waiting', waitingHandler);
    this.tech_.on(timerCancelEvents, cancelTimerHandler);
    this.tech_.on('canplay', canPlayHandler); // Define the dispose function to clean up our events

    this.dispose = function () {
      _this.logger_('dispose');

      _this.tech_.off('seekablechanged', fixesBadSeeksHandler);

      _this.tech_.off('waiting', waitingHandler);

      _this.tech_.off(timerCancelEvents, cancelTimerHandler);

      _this.tech_.off('canplay', canPlayHandler);

      loaderTypes.forEach(function (type) {
        mpc[type + "SegmentLoader_"].off('appendsdone', loaderChecks[type].updateend);
        mpc[type + "SegmentLoader_"].off('playlistupdate', loaderChecks[type].reset);

        _this.tech_.off(['seeked', 'seeking'], loaderChecks[type].reset);
      });

      if (_this.checkCurrentTimeTimeout_) {
        window_1.clearTimeout(_this.checkCurrentTimeTimeout_);
      }

      _this.cancelTimer_();
    };
  }
  /**
   * Periodically check current time to see if playback stopped
   *
   * @private
   */


  var _proto = PlaybackWatcher.prototype;

  _proto.monitorCurrentTime_ = function monitorCurrentTime_() {
    this.checkCurrentTime_();

    if (this.checkCurrentTimeTimeout_) {
      window_1.clearTimeout(this.checkCurrentTimeTimeout_);
    } // 42 = 24 fps // 250 is what Webkit uses // FF uses 15


    this.checkCurrentTimeTimeout_ = window_1.setTimeout(this.monitorCurrentTime_.bind(this), 250);
  }
  /**
   * Reset stalled download stats for a specific type of loader
   *
   * @param {string} type
   *        The segment loader type to check.
   *
   * @listens SegmentLoader#playlistupdate
   * @listens Tech#seeking
   * @listens Tech#seeked
   */
  ;

  _proto.resetSegmentDownloads_ = function resetSegmentDownloads_(type) {
    var loader = this.masterPlaylistController_[type + "SegmentLoader_"];

    if (this[type + "StalledDownloads_"] > 0) {
      this.logger_("resetting possible stalled download count for " + type + " loader");
    }

    this[type + "StalledDownloads_"] = 0;
    this[type + "Buffered_"] = loader.buffered_();
  }
  /**
   * Checks on every segment `appendsdone` to see
   * if segment appends are making progress. If they are not
   * and we are still downloading bytes. We blacklist the playlist.
   *
   * @param {string} type
   *        The segment loader type to check.
   *
   * @listens SegmentLoader#appendsdone
   */
  ;

  _proto.checkSegmentDownloads_ = function checkSegmentDownloads_(type) {
    var mpc = this.masterPlaylistController_;
    var loader = mpc[type + "SegmentLoader_"];
    var buffered = loader.buffered_();
    var isBufferedDifferent = isRangeDifferent(this[type + "Buffered_"], buffered);
    this[type + "Buffered_"] = buffered; // if another watcher is going to fix the issue or
    // the buffered value for this loader changed
    // appends are working

    if (isBufferedDifferent) {
      this.resetSegmentDownloads_(type);
      return;
    }

    this[type + "StalledDownloads_"]++;
    this.logger_("found #" + this[type + "StalledDownloads_"] + " " + type + " appends that did not increase buffer (possible stalled download)", {
      playlistId: loader.playlist_ && loader.playlist_.id,
      buffered: timeRangesToArray(buffered)
    }); // after 10 possibly stalled appends with no reset, exclude

    if (this[type + "StalledDownloads_"] < 10) {
      return;
    }

    this.logger_(type + " loader stalled download exclusion");
    this.resetSegmentDownloads_(type);
    this.tech_.trigger({
      type: 'usage',
      name: "vhs-" + type + "-download-exclusion"
    });

    if (type === 'subtitle') {
      // TODO: Is there anything else that we can do here?
      // removing the track and disabling could have accesiblity implications.
      var track = loader.track();
      var label = track.label || track.language || 'Unknown';
      videojs$1.log.warn("Text track \"" + label + "\" is not working correctly. It will be disabled and excluded.");
      track.mode = 'disabled';
      this.tech_.textTracks().removeTrack(track);
      return;
    } // TODO: should we exclude audio tracks rather than main tracks
    // when type is audio?


    mpc.blacklistCurrentPlaylist({
      message: "Excessive " + type + " segment downloading detected."
    }, Infinity);
  }
  /**
   * The purpose of this function is to emulate the "waiting" event on
   * browsers that do not emit it when they are waiting for more
   * data to continue playback
   *
   * @private
   */
  ;

  _proto.checkCurrentTime_ = function checkCurrentTime_() {
    if (this.tech_.seeking() && this.fixesBadSeeks_()) {
      this.consecutiveUpdates = 0;
      this.lastRecordedTime = this.tech_.currentTime();
      return;
    }

    if (this.tech_.paused() || this.tech_.seeking()) {
      return;
    }

    var currentTime = this.tech_.currentTime();
    var buffered = this.tech_.buffered();

    if (this.lastRecordedTime === currentTime && (!buffered.length || currentTime + SAFE_TIME_DELTA >= buffered.end(buffered.length - 1))) {
      // If current time is at the end of the final buffered region, then any playback
      // stall is most likely caused by buffering in a low bandwidth environment. The tech
      // should fire a `waiting` event in this scenario, but due to browser and tech
      // inconsistencies. Calling `techWaiting_` here allows us to simulate
      // responding to a native `waiting` event when the tech fails to emit one.
      return this.techWaiting_();
    }

    if (this.consecutiveUpdates >= 5 && currentTime === this.lastRecordedTime) {
      this.consecutiveUpdates++;
      this.waiting_();
    } else if (currentTime === this.lastRecordedTime) {
      this.consecutiveUpdates++;
    } else {
      this.consecutiveUpdates = 0;
      this.lastRecordedTime = currentTime;
    }
  }
  /**
   * Cancels any pending timers and resets the 'timeupdate' mechanism
   * designed to detect that we are stalled
   *
   * @private
   */
  ;

  _proto.cancelTimer_ = function cancelTimer_() {
    this.consecutiveUpdates = 0;

    if (this.timer_) {
      this.logger_('cancelTimer_');
      clearTimeout(this.timer_);
    }

    this.timer_ = null;
  }
  /**
   * Fixes situations where there's a bad seek
   *
   * @return {boolean} whether an action was taken to fix the seek
   * @private
   */
  ;

  _proto.fixesBadSeeks_ = function fixesBadSeeks_() {
    var seeking = this.tech_.seeking();

    if (!seeking) {
      return false;
    }

    var seekable = this.seekable();
    var currentTime = this.tech_.currentTime();
    var isAfterSeekableRange = this.afterSeekableWindow_(seekable, currentTime, this.media(), this.allowSeeksWithinUnsafeLiveWindow);
    var seekTo;

    if (isAfterSeekableRange) {
      var seekableEnd = seekable.end(seekable.length - 1); // sync to live point (if VOD, our seekable was updated and we're simply adjusting)

      seekTo = seekableEnd;
    }

    if (this.beforeSeekableWindow_(seekable, currentTime)) {
      var seekableStart = seekable.start(0); // sync to the beginning of the live window
      // provide a buffer of .1 seconds to handle rounding/imprecise numbers

      seekTo = seekableStart + ( // if the playlist is too short and the seekable range is an exact time (can
      // happen in live with a 3 segment playlist), then don't use a time delta
      seekableStart === seekable.end(0) ? 0 : SAFE_TIME_DELTA);
    }

    if (typeof seekTo !== 'undefined') {
      this.logger_("Trying to seek outside of seekable at time " + currentTime + " with " + ("seekable range " + printableRange(seekable) + ". Seeking to ") + (seekTo + "."));
      this.tech_.setCurrentTime(seekTo);
      return true;
    }

    var buffered = this.tech_.buffered();

    if (closeToBufferedContent({
      buffered: buffered,
      targetDuration: this.media().targetDuration,
      currentTime: currentTime
    })) {
      seekTo = buffered.start(0) + SAFE_TIME_DELTA;
      this.logger_("Buffered region starts (" + buffered.start(0) + ") " + (" just beyond seek point (" + currentTime + "). Seeking to " + seekTo + "."));
      this.tech_.setCurrentTime(seekTo);
      return true;
    }

    return false;
  }
  /**
   * Handler for situations when we determine the player is waiting.
   *
   * @private
   */
  ;

  _proto.waiting_ = function waiting_() {
    if (this.techWaiting_()) {
      return;
    } // All tech waiting checks failed. Use last resort correction


    var currentTime = this.tech_.currentTime();
    var buffered = this.tech_.buffered();
    var currentRange = findRange(buffered, currentTime); // Sometimes the player can stall for unknown reasons within a contiguous buffered
    // region with no indication that anything is amiss (seen in Firefox). Seeking to
    // currentTime is usually enough to kickstart the player. This checks that the player
    // is currently within a buffered region before attempting a corrective seek.
    // Chrome does not appear to continue `timeupdate` events after a `waiting` event
    // until there is ~ 3 seconds of forward buffer available. PlaybackWatcher should also
    // make sure there is ~3 seconds of forward buffer before taking any corrective action
    // to avoid triggering an `unknownwaiting` event when the network is slow.

    if (currentRange.length && currentTime + 3 <= currentRange.end(0)) {
      this.cancelTimer_();
      this.tech_.setCurrentTime(currentTime);
      this.logger_("Stopped at " + currentTime + " while inside a buffered region " + ("[" + currentRange.start(0) + " -> " + currentRange.end(0) + "]. Attempting to resume ") + 'playback by seeking to the current time.'); // unknown waiting corrections may be useful for monitoring QoS

      this.tech_.trigger({
        type: 'usage',
        name: 'vhs-unknown-waiting'
      });
      this.tech_.trigger({
        type: 'usage',
        name: 'hls-unknown-waiting'
      });
      return;
    }
  }
  /**
   * Handler for situations when the tech fires a `waiting` event
   *
   * @return {boolean}
   *         True if an action (or none) was needed to correct the waiting. False if no
   *         checks passed
   * @private
   */
  ;

  _proto.techWaiting_ = function techWaiting_() {
    var seekable = this.seekable();
    var currentTime = this.tech_.currentTime();

    if (this.tech_.seeking() && this.fixesBadSeeks_()) {
      // Tech is seeking or bad seek fixed, no action needed
      return true;
    }

    if (this.tech_.seeking() || this.timer_ !== null) {
      // Tech is seeking or already waiting on another action, no action needed
      return true;
    }

    if (this.beforeSeekableWindow_(seekable, currentTime)) {
      var livePoint = seekable.end(seekable.length - 1);
      this.logger_("Fell out of live window at time " + currentTime + ". Seeking to " + ("live point (seekable end) " + livePoint));
      this.cancelTimer_();
      this.tech_.setCurrentTime(livePoint); // live window resyncs may be useful for monitoring QoS

      this.tech_.trigger({
        type: 'usage',
        name: 'vhs-live-resync'
      });
      this.tech_.trigger({
        type: 'usage',
        name: 'hls-live-resync'
      });
      return true;
    }

    var sourceUpdater = this.tech_.vhs.masterPlaylistController_.sourceUpdater_;
    var buffered = this.tech_.buffered();
    var videoUnderflow = this.videoUnderflow_({
      audioBuffered: sourceUpdater.audioBuffered(),
      videoBuffered: sourceUpdater.videoBuffered(),
      currentTime: currentTime
    });

    if (videoUnderflow) {
      // Even though the video underflowed and was stuck in a gap, the audio overplayed
      // the gap, leading currentTime into a buffered range. Seeking to currentTime
      // allows the video to catch up to the audio position without losing any audio
      // (only suffering ~3 seconds of frozen video and a pause in audio playback).
      this.cancelTimer_();
      this.tech_.setCurrentTime(currentTime); // video underflow may be useful for monitoring QoS

      this.tech_.trigger({
        type: 'usage',
        name: 'vhs-video-underflow'
      });
      this.tech_.trigger({
        type: 'usage',
        name: 'hls-video-underflow'
      });
      return true;
    }

    var nextRange = findNextRange(buffered, currentTime); // check for gap

    if (nextRange.length > 0) {
      var difference = nextRange.start(0) - currentTime;
      this.logger_("Stopped at " + currentTime + ", setting timer for " + difference + ", seeking " + ("to " + nextRange.start(0)));
      this.cancelTimer_();
      this.timer_ = setTimeout(this.skipTheGap_.bind(this), difference * 1000, currentTime);
      return true;
    } // All checks failed. Returning false to indicate failure to correct waiting


    return false;
  };

  _proto.afterSeekableWindow_ = function afterSeekableWindow_(seekable, currentTime, playlist, allowSeeksWithinUnsafeLiveWindow) {
    if (allowSeeksWithinUnsafeLiveWindow === void 0) {
      allowSeeksWithinUnsafeLiveWindow = false;
    }

    if (!seekable.length) {
      // we can't make a solid case if there's no seekable, default to false
      return false;
    }

    var allowedEnd = seekable.end(seekable.length - 1) + SAFE_TIME_DELTA;
    var isLive = !playlist.endList;

    if (isLive && allowSeeksWithinUnsafeLiveWindow) {
      allowedEnd = seekable.end(seekable.length - 1) + playlist.targetDuration * 3;
    }

    if (currentTime > allowedEnd) {
      return true;
    }

    return false;
  };

  _proto.beforeSeekableWindow_ = function beforeSeekableWindow_(seekable, currentTime) {
    if (seekable.length && // can't fall before 0 and 0 seekable start identifies VOD stream
    seekable.start(0) > 0 && currentTime < seekable.start(0) - SAFE_TIME_DELTA) {
      return true;
    }

    return false;
  };

  _proto.videoUnderflow_ = function videoUnderflow_(_ref2) {
    var videoBuffered = _ref2.videoBuffered,
        audioBuffered = _ref2.audioBuffered,
        currentTime = _ref2.currentTime;

    // audio only content will not have video underflow :)
    if (!videoBuffered) {
      return;
    }

    var gap; // find a gap in demuxed content.

    if (videoBuffered.length && audioBuffered.length) {
      // in Chrome audio will continue to play for ~3s when we run out of video
      // so we have to check that the video buffer did have some buffer in the
      // past.
      var lastVideoRange = findRange(videoBuffered, currentTime - 3);
      var videoRange = findRange(videoBuffered, currentTime);
      var audioRange = findRange(audioBuffered, currentTime);

      if (audioRange.length && !videoRange.length && lastVideoRange.length) {
        gap = {
          start: lastVideoRange.end(0),
          end: audioRange.end(0)
        };
      } // find a gap in muxed content.

    } else {
      var nextRange = findNextRange(videoBuffered, currentTime); // Even if there is no available next range, there is still a possibility we are
      // stuck in a gap due to video underflow.

      if (!nextRange.length) {
        gap = this.gapFromVideoUnderflow_(videoBuffered, currentTime);
      }
    }

    if (gap) {
      this.logger_("Encountered a gap in video from " + gap.start + " to " + gap.end + ". " + ("Seeking to current time " + currentTime));
      return true;
    }

    return false;
  }
  /**
   * Timer callback. If playback still has not proceeded, then we seek
   * to the start of the next buffered region.
   *
   * @private
   */
  ;

  _proto.skipTheGap_ = function skipTheGap_(scheduledCurrentTime) {
    var buffered = this.tech_.buffered();
    var currentTime = this.tech_.currentTime();
    var nextRange = findNextRange(buffered, currentTime);
    this.cancelTimer_();

    if (nextRange.length === 0 || currentTime !== scheduledCurrentTime) {
      return;
    }

    this.logger_('skipTheGap_:', 'currentTime:', currentTime, 'scheduled currentTime:', scheduledCurrentTime, 'nextRange start:', nextRange.start(0)); // only seek if we still have not played

    this.tech_.setCurrentTime(nextRange.start(0) + TIME_FUDGE_FACTOR);
    this.tech_.trigger({
      type: 'usage',
      name: 'vhs-gap-skip'
    });
    this.tech_.trigger({
      type: 'usage',
      name: 'hls-gap-skip'
    });
  };

  _proto.gapFromVideoUnderflow_ = function gapFromVideoUnderflow_(buffered, currentTime) {
    // At least in Chrome, if there is a gap in the video buffer, the audio will continue
    // playing for ~3 seconds after the video gap starts. This is done to account for
    // video buffer underflow/underrun (note that this is not done when there is audio
    // buffer underflow/underrun -- in that case the video will stop as soon as it
    // encounters the gap, as audio stalls are more noticeable/jarring to a user than
    // video stalls). The player's time will reflect the playthrough of audio, so the
    // time will appear as if we are in a buffered region, even if we are stuck in a
    // "gap."
    //
    // Example:
    // video buffer:   0 => 10.1, 10.2 => 20
    // audio buffer:   0 => 20
    // overall buffer: 0 => 10.1, 10.2 => 20
    // current time: 13
    //
    // Chrome's video froze at 10 seconds, where the video buffer encountered the gap,
    // however, the audio continued playing until it reached ~3 seconds past the gap
    // (13 seconds), at which point it stops as well. Since current time is past the
    // gap, findNextRange will return no ranges.
    //
    // To check for this issue, we see if there is a gap that starts somewhere within
    // a 3 second range (3 seconds +/- 1 second) back from our current time.
    var gaps = findGaps(buffered);

    for (var i = 0; i < gaps.length; i++) {
      var start = gaps.start(i);
      var end = gaps.end(i); // gap is starts no more than 4 seconds back

      if (currentTime - start < 4 && currentTime - start > 2) {
        return {
          start: start,
          end: end
        };
      }
    }

    return null;
  };

  return PlaybackWatcher;
}();

var defaultOptions = {
  errorInterval: 30,
  getSource: function getSource(next) {
    var tech = this.tech({
      IWillNotUseThisInPlugins: true
    });
    var sourceObj = tech.currentSource_ || this.currentSource();
    return next(sourceObj);
  }
};
/**
 * Main entry point for the plugin
 *
 * @param {Player} player a reference to a videojs Player instance
 * @param {Object} [options] an object with plugin options
 * @private
 */

var initPlugin = function initPlugin(player, options) {
  var lastCalled = 0;
  var seekTo = 0;
  var localOptions = videojs$1.mergeOptions(defaultOptions, options);
  player.ready(function () {
    player.trigger({
      type: 'usage',
      name: 'vhs-error-reload-initialized'
    });
    player.trigger({
      type: 'usage',
      name: 'hls-error-reload-initialized'
    });
  });
  /**
   * Player modifications to perform that must wait until `loadedmetadata`
   * has been triggered
   *
   * @private
   */

  var loadedMetadataHandler = function loadedMetadataHandler() {
    if (seekTo) {
      player.currentTime(seekTo);
    }
  };
  /**
   * Set the source on the player element, play, and seek if necessary
   *
   * @param {Object} sourceObj An object specifying the source url and mime-type to play
   * @private
   */


  var setSource = function setSource(sourceObj) {
    if (sourceObj === null || sourceObj === undefined) {
      return;
    }

    seekTo = player.duration() !== Infinity && player.currentTime() || 0;
    player.one('loadedmetadata', loadedMetadataHandler);
    player.src(sourceObj);
    player.trigger({
      type: 'usage',
      name: 'vhs-error-reload'
    });
    player.trigger({
      type: 'usage',
      name: 'hls-error-reload'
    });
    player.play();
  };
  /**
   * Attempt to get a source from either the built-in getSource function
   * or a custom function provided via the options
   *
   * @private
   */


  var errorHandler = function errorHandler() {
    // Do not attempt to reload the source if a source-reload occurred before
    // 'errorInterval' time has elapsed since the last source-reload
    if (Date.now() - lastCalled < localOptions.errorInterval * 1000) {
      player.trigger({
        type: 'usage',
        name: 'vhs-error-reload-canceled'
      });
      player.trigger({
        type: 'usage',
        name: 'hls-error-reload-canceled'
      });
      return;
    }

    if (!localOptions.getSource || typeof localOptions.getSource !== 'function') {
      videojs$1.log.error('ERROR: reloadSourceOnError - The option getSource must be a function!');
      return;
    }

    lastCalled = Date.now();
    return localOptions.getSource.call(player, setSource);
  };
  /**
   * Unbind any event handlers that were bound by the plugin
   *
   * @private
   */


  var cleanupEvents = function cleanupEvents() {
    player.off('loadedmetadata', loadedMetadataHandler);
    player.off('error', errorHandler);
    player.off('dispose', cleanupEvents);
  };
  /**
   * Cleanup before re-initializing the plugin
   *
   * @param {Object} [newOptions] an object with plugin options
   * @private
   */


  var reinitPlugin = function reinitPlugin(newOptions) {
    cleanupEvents();
    initPlugin(player, newOptions);
  };

  player.on('error', errorHandler);
  player.on('dispose', cleanupEvents); // Overwrite the plugin function so that we can correctly cleanup before
  // initializing the plugin

  player.reloadSourceOnError = reinitPlugin;
};
/**
 * Reload the source when an error is detected as long as there
 * wasn't an error previously within the last 30 seconds
 *
 * @param {Object} [options] an object with plugin options
 */


var reloadSourceOnError = function reloadSourceOnError(options) {
  initPlugin(this, options);
};

var version$4 = "2.2.0";

var Vhs$1 = {
  PlaylistLoader: PlaylistLoader,
  Playlist: Playlist,
  utils: utils$1,
  STANDARD_PLAYLIST_SELECTOR: lastBandwidthSelector,
  INITIAL_PLAYLIST_SELECTOR: lowestBitrateCompatibleVariantSelector,
  comparePlaylistBandwidth: comparePlaylistBandwidth,
  comparePlaylistResolution: comparePlaylistResolution,
  xhr: xhrFactory()
}; // Define getter/setters for config properties

['GOAL_BUFFER_LENGTH', 'MAX_GOAL_BUFFER_LENGTH', 'BACK_BUFFER_LENGTH', 'GOAL_BUFFER_LENGTH_RATE', 'BUFFER_LOW_WATER_LINE', 'MAX_BUFFER_LOW_WATER_LINE', 'BUFFER_LOW_WATER_LINE_RATE', 'BANDWIDTH_VARIANCE'].forEach(function (prop) {
  Object.defineProperty(Vhs$1, prop, {
    get: function get() {
      videojs$1.log.warn("using Vhs." + prop + " is UNSAFE be sure you know what you are doing");
      return Config[prop];
    },
    set: function set(value) {
      videojs$1.log.warn("using Vhs." + prop + " is UNSAFE be sure you know what you are doing");

      if (typeof value !== 'number' || value < 0) {
        videojs$1.log.warn("value of Vhs." + prop + " must be greater than or equal to 0");
        return;
      }

      Config[prop] = value;
    }
  });
});
var LOCAL_STORAGE_KEY = 'videojs-vhs';
/**
 * Updates the selectedIndex of the QualityLevelList when a mediachange happens in vhs.
 *
 * @param {QualityLevelList} qualityLevels The QualityLevelList to update.
 * @param {PlaylistLoader} playlistLoader PlaylistLoader containing the new media info.
 * @function handleVhsMediaChange
 */

var handleVhsMediaChange = function handleVhsMediaChange(qualityLevels, playlistLoader) {
  var newPlaylist = playlistLoader.media();
  var selectedIndex = -1;

  for (var i = 0; i < qualityLevels.length; i++) {
    if (qualityLevels[i].id === newPlaylist.id) {
      selectedIndex = i;
      break;
    }
  }

  qualityLevels.selectedIndex_ = selectedIndex;
  qualityLevels.trigger({
    selectedIndex: selectedIndex,
    type: 'change'
  });
};
/**
 * Adds quality levels to list once playlist metadata is available
 *
 * @param {QualityLevelList} qualityLevels The QualityLevelList to attach events to.
 * @param {Object} vhs Vhs object to listen to for media events.
 * @function handleVhsLoadedMetadata
 */


var handleVhsLoadedMetadata = function handleVhsLoadedMetadata(qualityLevels, vhs) {
  vhs.representations().forEach(function (rep) {
    qualityLevels.addQualityLevel(rep);
  });
  handleVhsMediaChange(qualityLevels, vhs.playlists);
}; // HLS is a source handler, not a tech. Make sure attempts to use it
// as one do not cause exceptions.


Vhs$1.canPlaySource = function () {
  return videojs$1.log.warn('HLS is no longer a tech. Please remove it from ' + 'your player\'s techOrder.');
};

var emeKeySystems = function emeKeySystems(keySystemOptions, videoPlaylist, audioPlaylist) {
  if (!keySystemOptions) {
    return keySystemOptions;
  }

  var codecs$1 = {
    video: videoPlaylist && videoPlaylist.attributes && videoPlaylist.attributes.CODECS,
    audio: audioPlaylist && audioPlaylist.attributes && audioPlaylist.attributes.CODECS
  };

  if (!codecs$1.audio && codecs$1.video && codecs$1.video.split(',').length > 1) {
    codecs$1.video.split(',').forEach(function (codec) {
      codec = codec.trim();

      if (codecs.isAudioCodec(codec)) {
        codecs$1.audio = codec;
      } else if (codecs.isVideoCodec(codec)) {
        codecs$1.video = codec;
      }
    });
  }

  var videoContentType = codecs$1.video ? "video/mp4;codecs=\"" + codecs$1.video + "\"" : null;
  var audioContentType = codecs$1.audio ? "audio/mp4;codecs=\"" + codecs$1.audio + "\"" : null; // upsert the content types based on the selected playlist

  var keySystemContentTypes = {};

  for (var keySystem in keySystemOptions) {
    keySystemContentTypes[keySystem] = {
      audioContentType: audioContentType,
      videoContentType: videoContentType
    };

    if (videoPlaylist.contentProtection && videoPlaylist.contentProtection[keySystem] && videoPlaylist.contentProtection[keySystem].pssh) {
      keySystemContentTypes[keySystem].pssh = videoPlaylist.contentProtection[keySystem].pssh;
    } // videojs-contrib-eme accepts the option of specifying: 'com.some.cdm': 'url'
    // so we need to prevent overwriting the URL entirely


    if (typeof keySystemOptions[keySystem] === 'string') {
      keySystemContentTypes[keySystem].url = keySystemOptions[keySystem];
    }
  }

  return videojs$1.mergeOptions(keySystemOptions, keySystemContentTypes);
};
/**
 * @typedef {Object} KeySystems
 *
 * keySystems configuration for https://github.com/videojs/videojs-contrib-eme
 * Note: not all options are listed here.
 *
 * @property {Uint8Array} [pssh]
 *           Protection System Specific Header
 */

/**
 * Goes through all the playlists and collects an array of KeySystems options objects
 * containing each playlist's keySystems and their pssh values, if available.
 *
 * @param {Object[]} playlists
 *        The playlists to look through
 * @param {string[]} keySystems
 *        The keySystems to collect pssh values for
 *
 * @return {KeySystems[]}
 *         An array of KeySystems objects containing available key systems and their
 *         pssh values
 */


var getAllPsshKeySystemsOptions = function getAllPsshKeySystemsOptions(playlists, keySystems) {
  return playlists.reduce(function (keySystemsArr, playlist) {
    if (!playlist.contentProtection) {
      return keySystemsArr;
    }

    var keySystemsOptions = keySystems.reduce(function (keySystemsObj, keySystem) {
      var keySystemOptions = playlist.contentProtection[keySystem];

      if (keySystemOptions && keySystemOptions.pssh) {
        keySystemsObj[keySystem] = {
          pssh: keySystemOptions.pssh
        };
      }

      return keySystemsObj;
    }, {});

    if (Object.keys(keySystemsOptions).length) {
      keySystemsArr.push(keySystemsOptions);
    }

    return keySystemsArr;
  }, []);
};
/**
 * If the [eme](https://github.com/videojs/videojs-contrib-eme) plugin is available, and
 * there are keySystems on the source, sets up source options to prepare the source for
 * eme and tries to initialize it early via eme's initializeMediaKeys API (if available).
 *
 * @param {Object} player
 *        The player instance
 * @param {Object[]} sourceKeySystems
 *        The key systems options from the player source
 * @param {Object} media
 *        The active media playlist
 * @param {Object} [audioMedia]
 *        The active audio media playlist (optional)
 * @param {Object[]} mainPlaylists
 *        The playlists found on the master playlist object
 */


var setupEmeOptions = function setupEmeOptions(_ref) {
  var player = _ref.player,
      sourceKeySystems = _ref.sourceKeySystems,
      media = _ref.media,
      audioMedia = _ref.audioMedia,
      mainPlaylists = _ref.mainPlaylists;
  var sourceOptions = emeKeySystems(sourceKeySystems, media, audioMedia);

  if (!sourceOptions) {
    return;
  }

  player.currentSource().keySystems = sourceOptions; // eme handles the rest of the setup, so if it is missing
  // do nothing.

  if (sourceOptions && !player.eme) {
    videojs$1.log.warn('DRM encrypted source cannot be decrypted without a DRM plugin');
    return;
  } // works around https://bugs.chromium.org/p/chromium/issues/detail?id=895449
  // in non-IE11 browsers. In IE11 this is too early to initialize media keys


  if (videojs$1.browser.IE_VERSION === 11 || !player.eme.initializeMediaKeys) {
    return;
  } // TODO should all audio PSSH values be initialized for DRM?
  //
  // All unique video rendition pssh values are initialized for DRM, but here only
  // the initial audio playlist license is initialized. In theory, an encrypted
  // event should be fired if the user switches to an alternative audio playlist
  // where a license is required, but this case hasn't yet been tested. In addition, there
  // may be many alternate audio playlists unlikely to be used (e.g., multiple different
  // languages).


  var playlists = audioMedia ? mainPlaylists.concat([audioMedia]) : mainPlaylists;
  var keySystemsOptionsArr = getAllPsshKeySystemsOptions(playlists, Object.keys(sourceKeySystems)); // Since PSSH values are interpreted as initData, EME will dedupe any duplicates. The
  // only place where it should not be deduped is for ms-prefixed APIs, but the early
  // return for IE11 above, and the existence of modern EME APIs in addition to
  // ms-prefixed APIs on Edge should prevent this from being a concern.
  // initializeMediaKeys also won't use the webkit-prefixed APIs.

  keySystemsOptionsArr.forEach(function (keySystemsOptions) {
    player.eme.initializeMediaKeys({
      keySystems: keySystemsOptions
    });
  });
};

var getVhsLocalStorage = function getVhsLocalStorage() {
  if (!window_1.localStorage) {
    return null;
  }

  var storedObject = window_1.localStorage.getItem(LOCAL_STORAGE_KEY);

  if (!storedObject) {
    return null;
  }

  try {
    return JSON.parse(storedObject);
  } catch (e) {
    // someone may have tampered with the value
    return null;
  }
};

var updateVhsLocalStorage = function updateVhsLocalStorage(options) {
  if (!window_1.localStorage) {
    return false;
  }

  var objectToStore = getVhsLocalStorage();
  objectToStore = objectToStore ? videojs$1.mergeOptions(objectToStore, options) : options;

  try {
    window_1.localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(objectToStore));
  } catch (e) {
    // Throws if storage is full (e.g., always on iOS 5+ Safari private mode, where
    // storage is set to 0).
    // https://developer.mozilla.org/en-US/docs/Web/API/Storage/setItem#Exceptions
    // No need to perform any operation.
    return false;
  }

  return objectToStore;
};
/**
 * Parses VHS-supported media types from data URIs. See
 * https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/Data_URIs
 * for information on data URIs.
 *
 * @param {string} dataUri
 *        The data URI
 *
 * @return {string|Object}
 *         The parsed object/string, or the original string if no supported media type
 *         was found
 */


var expandDataUri = function expandDataUri(dataUri) {
  if (dataUri.toLowerCase().indexOf('data:application/vnd.videojs.vhs+json,') === 0) {
    return JSON.parse(dataUri.substring(dataUri.indexOf(',') + 1));
  } // no known case for this data URI, return the string as-is


  return dataUri;
};
/**
 * Whether the browser has built-in HLS support.
 */


Vhs$1.supportsNativeHls = function () {
  if (!document_1 || !document_1.createElement) {
    return false;
  }

  var video = document_1.createElement('video'); // native HLS is definitely not supported if HTML5 video isn't

  if (!videojs$1.getTech('Html5').isSupported()) {
    return false;
  } // HLS manifests can go by many mime-types


  var canPlay = [// Apple santioned
  'application/vnd.apple.mpegurl', // Apple sanctioned for backwards compatibility
  'audio/mpegurl', // Very common
  'audio/x-mpegurl', // Very common
  'application/x-mpegurl', // Included for completeness
  'video/x-mpegurl', 'video/mpegurl', 'application/mpegurl'];
  return canPlay.some(function (canItPlay) {
    return /maybe|probably/i.test(video.canPlayType(canItPlay));
  });
}();

Vhs$1.supportsNativeDash = function () {
  if (!document_1 || !document_1.createElement || !videojs$1.getTech('Html5').isSupported()) {
    return false;
  }

  return /maybe|probably/i.test(document_1.createElement('video').canPlayType('application/dash+xml'));
}();

Vhs$1.supportsTypeNatively = function (type) {
  if (type === 'hls') {
    return Vhs$1.supportsNativeHls;
  }

  if (type === 'dash') {
    return Vhs$1.supportsNativeDash;
  }

  return false;
};
/**
 * HLS is a source handler, not a tech. Make sure attempts to use it
 * as one do not cause exceptions.
 */


Vhs$1.isSupported = function () {
  return videojs$1.log.warn('HLS is no longer a tech. Please remove it from ' + 'your player\'s techOrder.');
};

var Component = videojs$1.getComponent('Component');
/**
 * The Vhs Handler object, where we orchestrate all of the parts
 * of HLS to interact with video.js
 *
 * @class VhsHandler
 * @extends videojs.Component
 * @param {Object} source the soruce object
 * @param {Tech} tech the parent tech object
 * @param {Object} options optional and required options
 */

var VhsHandler = /*#__PURE__*/function (_Component) {
  inheritsLoose(VhsHandler, _Component);

  function VhsHandler(source, tech, options) {
    var _this;

    _this = _Component.call(this, tech, videojs$1.mergeOptions(options.hls, options.vhs)) || this;

    if (options.hls && Object.keys(options.hls).length) {
      videojs$1.log.warn('Using hls options is deprecated. Use vhs instead.');
    } // tech.player() is deprecated but setup a reference to HLS for
    // backwards-compatibility


    if (tech.options_ && tech.options_.playerId) {
      var _player = videojs$1(tech.options_.playerId);

      if (!_player.hasOwnProperty('hls')) {
        Object.defineProperty(_player, 'hls', {
          get: function get() {
            videojs$1.log.warn('player.hls is deprecated. Use player.tech().vhs instead.');
            tech.trigger({
              type: 'usage',
              name: 'hls-player-access'
            });
            return assertThisInitialized(_this);
          },
          configurable: true
        });
      }

      if (!_player.hasOwnProperty('vhs')) {
        Object.defineProperty(_player, 'vhs', {
          get: function get() {
            videojs$1.log.warn('player.vhs is deprecated. Use player.tech().vhs instead.');
            tech.trigger({
              type: 'usage',
              name: 'vhs-player-access'
            });
            return assertThisInitialized(_this);
          },
          configurable: true
        });
      }

      if (!_player.hasOwnProperty('dash')) {
        Object.defineProperty(_player, 'dash', {
          get: function get() {
            videojs$1.log.warn('player.dash is deprecated. Use player.tech().vhs instead.');
            return assertThisInitialized(_this);
          },
          configurable: true
        });
      }

      _this.player_ = _player;
    }

    _this.tech_ = tech;
    _this.source_ = source;
    _this.stats = {};
    _this.ignoreNextSeekingEvent_ = false;

    _this.setOptions_();

    if (_this.options_.overrideNative && tech.overrideNativeAudioTracks && tech.overrideNativeVideoTracks) {
      tech.overrideNativeAudioTracks(true);
      tech.overrideNativeVideoTracks(true);
    } else if (_this.options_.overrideNative && (tech.featuresNativeVideoTracks || tech.featuresNativeAudioTracks)) {
      // overriding native HLS only works if audio tracks have been emulated
      // error early if we're misconfigured
      throw new Error('Overriding native HLS requires emulated tracks. ' + 'See https://git.io/vMpjB');
    } // listen for fullscreenchange events for this player so that we
    // can adjust our quality selection quickly


    _this.on(document_1, ['fullscreenchange', 'webkitfullscreenchange', 'mozfullscreenchange', 'MSFullscreenChange'], function (event) {
      var fullscreenElement = document_1.fullscreenElement || document_1.webkitFullscreenElement || document_1.mozFullScreenElement || document_1.msFullscreenElement;

      if (fullscreenElement && fullscreenElement.contains(_this.tech_.el())) {
        _this.masterPlaylistController_.smoothQualityChange_();
      }
    });

    _this.on(_this.tech_, 'seeking', function () {
      if (this.ignoreNextSeekingEvent_) {
        this.ignoreNextSeekingEvent_ = false;
        return;
      }

      this.setCurrentTime(this.tech_.currentTime());
    });

    _this.on(_this.tech_, 'error', function () {
      if (this.masterPlaylistController_) {
        this.masterPlaylistController_.pauseLoading();
      }
    });

    _this.on(_this.tech_, 'play', _this.play);

    return _this;
  }

  var _proto = VhsHandler.prototype;

  _proto.setOptions_ = function setOptions_() {
    var _this2 = this;

    // defaults
    this.options_.withCredentials = this.options_.withCredentials || false;
    this.options_.handleManifestRedirects = this.options_.handleManifestRedirects === false ? false : true;
    this.options_.limitRenditionByPlayerDimensions = this.options_.limitRenditionByPlayerDimensions === false ? false : true;
    this.options_.useDevicePixelRatio = this.options_.useDevicePixelRatio || false;
    this.options_.smoothQualityChange = this.options_.smoothQualityChange || false;
    this.options_.useBandwidthFromLocalStorage = typeof this.source_.useBandwidthFromLocalStorage !== 'undefined' ? this.source_.useBandwidthFromLocalStorage : this.options_.useBandwidthFromLocalStorage || false;
    this.options_.customTagParsers = this.options_.customTagParsers || [];
    this.options_.customTagMappers = this.options_.customTagMappers || [];
    this.options_.cacheEncryptionKeys = this.options_.cacheEncryptionKeys || false;
    this.options_.handlePartialData = this.options_.handlePartialData || false;

    if (typeof this.options_.blacklistDuration !== 'number') {
      this.options_.blacklistDuration = 5 * 60;
    }

    if (typeof this.options_.bandwidth !== 'number') {
      if (this.options_.useBandwidthFromLocalStorage) {
        var storedObject = getVhsLocalStorage();

        if (storedObject && storedObject.bandwidth) {
          this.options_.bandwidth = storedObject.bandwidth;
          this.tech_.trigger({
            type: 'usage',
            name: 'vhs-bandwidth-from-local-storage'
          });
          this.tech_.trigger({
            type: 'usage',
            name: 'hls-bandwidth-from-local-storage'
          });
        }

        if (storedObject && storedObject.throughput) {
          this.options_.throughput = storedObject.throughput;
          this.tech_.trigger({
            type: 'usage',
            name: 'vhs-throughput-from-local-storage'
          });
          this.tech_.trigger({
            type: 'usage',
            name: 'hls-throughput-from-local-storage'
          });
        }
      }
    } // if bandwidth was not set by options or pulled from local storage, start playlist
    // selection at a reasonable bandwidth


    if (typeof this.options_.bandwidth !== 'number') {
      this.options_.bandwidth = Config.INITIAL_BANDWIDTH;
    } // If the bandwidth number is unchanged from the initial setting
    // then this takes precedence over the enableLowInitialPlaylist option


    this.options_.enableLowInitialPlaylist = this.options_.enableLowInitialPlaylist && this.options_.bandwidth === Config.INITIAL_BANDWIDTH; // grab options passed to player.src

    ['withCredentials', 'useDevicePixelRatio', 'limitRenditionByPlayerDimensions', 'bandwidth', 'smoothQualityChange', 'customTagParsers', 'customTagMappers', 'handleManifestRedirects', 'cacheEncryptionKeys', 'handlePartialData'].forEach(function (option) {
      if (typeof _this2.source_[option] !== 'undefined') {
        _this2.options_[option] = _this2.source_[option];
      }
    });
    this.limitRenditionByPlayerDimensions = this.options_.limitRenditionByPlayerDimensions;
    this.useDevicePixelRatio = this.options_.useDevicePixelRatio;
  }
  /**
   * called when player.src gets called, handle a new source
   *
   * @param {Object} src the source object to handle
   */
  ;

  _proto.src = function src(_src, type) {
    var _this3 = this;

    // do nothing if the src is falsey
    if (!_src) {
      return;
    }

    this.setOptions_(); // add master playlist controller options

    this.options_.src = expandDataUri(this.source_.src);
    this.options_.tech = this.tech_;
    this.options_.externVhs = Vhs$1;
    this.options_.sourceType = mediaTypes.simpleTypeFromSourceType(type); // Whenever we seek internally, we should update the tech

    this.options_.seekTo = function (time) {
      _this3.tech_.setCurrentTime(time);
    };

    this.masterPlaylistController_ = new MasterPlaylistController(this.options_);
    this.playbackWatcher_ = new PlaybackWatcher(videojs$1.mergeOptions(this.options_, {
      seekable: function seekable() {
        return _this3.seekable();
      },
      media: function media() {
        return _this3.masterPlaylistController_.media();
      },
      masterPlaylistController: this.masterPlaylistController_
    }));
    this.masterPlaylistController_.on('error', function () {
      var player = videojs$1.players[_this3.tech_.options_.playerId];
      var error = _this3.masterPlaylistController_.error;

      if (typeof error === 'object' && !error.code) {
        error.code = 3;
      } else if (typeof error === 'string') {
        error = {
          message: error,
          code: 3
        };
      }

      player.error(error);
    }); // `this` in selectPlaylist should be the VhsHandler for backwards
    // compatibility with < v2

    this.masterPlaylistController_.selectPlaylist = this.selectPlaylist ? this.selectPlaylist.bind(this) : Vhs$1.STANDARD_PLAYLIST_SELECTOR.bind(this);
    this.masterPlaylistController_.selectInitialPlaylist = Vhs$1.INITIAL_PLAYLIST_SELECTOR.bind(this); // re-expose some internal objects for backwards compatibility with < v2

    this.playlists = this.masterPlaylistController_.masterPlaylistLoader_;
    this.mediaSource = this.masterPlaylistController_.mediaSource; // Proxy assignment of some properties to the master playlist
    // controller. Using a custom property for backwards compatibility
    // with < v2

    Object.defineProperties(this, {
      selectPlaylist: {
        get: function get() {
          return this.masterPlaylistController_.selectPlaylist;
        },
        set: function set(selectPlaylist) {
          this.masterPlaylistController_.selectPlaylist = selectPlaylist.bind(this);
        }
      },
      throughput: {
        get: function get() {
          return this.masterPlaylistController_.mainSegmentLoader_.throughput.rate;
        },
        set: function set(throughput) {
          this.masterPlaylistController_.mainSegmentLoader_.throughput.rate = throughput; // By setting `count` to 1 the throughput value becomes the starting value
          // for the cumulative average

          this.masterPlaylistController_.mainSegmentLoader_.throughput.count = 1;
        }
      },
      bandwidth: {
        get: function get() {
          return this.masterPlaylistController_.mainSegmentLoader_.bandwidth;
        },
        set: function set(bandwidth) {
          this.masterPlaylistController_.mainSegmentLoader_.bandwidth = bandwidth; // setting the bandwidth manually resets the throughput counter
          // `count` is set to zero that current value of `rate` isn't included
          // in the cumulative average

          this.masterPlaylistController_.mainSegmentLoader_.throughput = {
            rate: 0,
            count: 0
          };
        }
      },

      /**
       * `systemBandwidth` is a combination of two serial processes bit-rates. The first
       * is the network bitrate provided by `bandwidth` and the second is the bitrate of
       * the entire process after that - decryption, transmuxing, and appending - provided
       * by `throughput`.
       *
       * Since the two process are serial, the overall system bandwidth is given by:
       *   sysBandwidth = 1 / (1 / bandwidth + 1 / throughput)
       */
      systemBandwidth: {
        get: function get() {
          var invBandwidth = 1 / (this.bandwidth || 1);
          var invThroughput;

          if (this.throughput > 0) {
            invThroughput = 1 / this.throughput;
          } else {
            invThroughput = 0;
          }

          var systemBitrate = Math.floor(1 / (invBandwidth + invThroughput));
          return systemBitrate;
        },
        set: function set() {
          videojs$1.log.error('The "systemBandwidth" property is read-only');
        }
      }
    });

    if (this.options_.bandwidth) {
      this.bandwidth = this.options_.bandwidth;
    }

    if (this.options_.throughput) {
      this.throughput = this.options_.throughput;
    }

    Object.defineProperties(this.stats, {
      bandwidth: {
        get: function get() {
          return _this3.bandwidth || 0;
        },
        enumerable: true
      },
      mediaRequests: {
        get: function get() {
          return _this3.masterPlaylistController_.mediaRequests_() || 0;
        },
        enumerable: true
      },
      mediaRequestsAborted: {
        get: function get() {
          return _this3.masterPlaylistController_.mediaRequestsAborted_() || 0;
        },
        enumerable: true
      },
      mediaRequestsTimedout: {
        get: function get() {
          return _this3.masterPlaylistController_.mediaRequestsTimedout_() || 0;
        },
        enumerable: true
      },
      mediaRequestsErrored: {
        get: function get() {
          return _this3.masterPlaylistController_.mediaRequestsErrored_() || 0;
        },
        enumerable: true
      },
      mediaTransferDuration: {
        get: function get() {
          return _this3.masterPlaylistController_.mediaTransferDuration_() || 0;
        },
        enumerable: true
      },
      mediaBytesTransferred: {
        get: function get() {
          return _this3.masterPlaylistController_.mediaBytesTransferred_() || 0;
        },
        enumerable: true
      },
      mediaSecondsLoaded: {
        get: function get() {
          return _this3.masterPlaylistController_.mediaSecondsLoaded_() || 0;
        },
        enumerable: true
      },
      buffered: {
        get: function get() {
          return timeRangesToArray(_this3.tech_.buffered());
        },
        enumerable: true
      },
      currentTime: {
        get: function get() {
          return _this3.tech_.currentTime();
        },
        enumerable: true
      },
      currentSource: {
        get: function get() {
          return _this3.tech_.currentSource_;
        },
        enumerable: true
      },
      currentTech: {
        get: function get() {
          return _this3.tech_.name_;
        },
        enumerable: true
      },
      duration: {
        get: function get() {
          return _this3.tech_.duration();
        },
        enumerable: true
      },
      master: {
        get: function get() {
          return _this3.playlists.master;
        },
        enumerable: true
      },
      playerDimensions: {
        get: function get() {
          return _this3.tech_.currentDimensions();
        },
        enumerable: true
      },
      seekable: {
        get: function get() {
          return timeRangesToArray(_this3.tech_.seekable());
        },
        enumerable: true
      },
      timestamp: {
        get: function get() {
          return Date.now();
        },
        enumerable: true
      },
      videoPlaybackQuality: {
        get: function get() {
          return _this3.tech_.getVideoPlaybackQuality();
        },
        enumerable: true
      }
    });
    this.tech_.one('canplay', this.masterPlaylistController_.setupFirstPlay.bind(this.masterPlaylistController_));
    this.tech_.on('bandwidthupdate', function () {
      if (_this3.options_.useBandwidthFromLocalStorage) {
        updateVhsLocalStorage({
          bandwidth: _this3.bandwidth,
          throughput: Math.round(_this3.throughput)
        });
      }
    });
    this.masterPlaylistController_.on('selectedinitialmedia', function () {
      // Add the manual rendition mix-in to VhsHandler
      renditionSelectionMixin(_this3);
    });
    this.masterPlaylistController_.sourceUpdater_.on('ready', function () {
      var audioPlaylistLoader = _this3.masterPlaylistController_.mediaTypes_.AUDIO.activePlaylistLoader;
      setupEmeOptions({
        player: _this3.player_,
        sourceKeySystems: _this3.source_.keySystems,
        media: _this3.playlists.media(),
        audioMedia: audioPlaylistLoader && audioPlaylistLoader.media(),
        mainPlaylists: _this3.playlists.master.playlists
      });
    }); // the bandwidth of the primary segment loader is our best
    // estimate of overall bandwidth

    this.on(this.masterPlaylistController_, 'progress', function () {
      this.tech_.trigger('progress');
    }); // In the live case, we need to ignore the very first `seeking` event since
    // that will be the result of the seek-to-live behavior

    this.on(this.masterPlaylistController_, 'firstplay', function () {
      this.ignoreNextSeekingEvent_ = true;
    });
    this.setupQualityLevels_(); // do nothing if the tech has been disposed already
    // this can occur if someone sets the src in player.ready(), for instance

    if (!this.tech_.el()) {
      return;
    }

    this.mediaSourceUrl_ = window_1.URL.createObjectURL(this.masterPlaylistController_.mediaSource);
    this.tech_.src(this.mediaSourceUrl_);
  }
  /**
   * Initializes the quality levels and sets listeners to update them.
   *
   * @method setupQualityLevels_
   * @private
   */
  ;

  _proto.setupQualityLevels_ = function setupQualityLevels_() {
    var _this4 = this;

    var player = videojs$1.players[this.tech_.options_.playerId]; // if there isn't a player or there isn't a qualityLevels plugin
    // or qualityLevels_ listeners have already been setup, do nothing.

    if (!player || !player.qualityLevels || this.qualityLevels_) {
      return;
    }

    this.qualityLevels_ = player.qualityLevels();
    this.masterPlaylistController_.on('selectedinitialmedia', function () {
      handleVhsLoadedMetadata(_this4.qualityLevels_, _this4);
    });
    this.playlists.on('mediachange', function () {
      handleVhsMediaChange(_this4.qualityLevels_, _this4.playlists);
    });
  }
  /**
   * return the version
   */
  ;

  VhsHandler.version = function version$5() {
    return {
      '@videojs/http-streaming': version$4,
      'mux.js': version,
      'mpd-parser': version$1,
      'm3u8-parser': version$2,
      'aes-decrypter': version$3
    };
  }
  /**
   * return the version
   */
  ;

  _proto.version = function version() {
    return this.constructor.version();
  };

  _proto.canChangeType = function canChangeType() {
    return SourceUpdater.canChangeType();
  }
  /**
   * Begin playing the video.
   */
  ;

  _proto.play = function play() {
    this.masterPlaylistController_.play();
  }
  /**
   * a wrapper around the function in MasterPlaylistController
   */
  ;

  _proto.setCurrentTime = function setCurrentTime(currentTime) {
    this.masterPlaylistController_.setCurrentTime(currentTime);
  }
  /**
   * a wrapper around the function in MasterPlaylistController
   */
  ;

  _proto.duration = function duration() {
    return this.masterPlaylistController_.duration();
  }
  /**
   * a wrapper around the function in MasterPlaylistController
   */
  ;

  _proto.seekable = function seekable() {
    return this.masterPlaylistController_.seekable();
  }
  /**
   * Abort all outstanding work and cleanup.
   */
  ;

  _proto.dispose = function dispose() {
    if (this.playbackWatcher_) {
      this.playbackWatcher_.dispose();
    }

    if (this.masterPlaylistController_) {
      this.masterPlaylistController_.dispose();
    }

    if (this.qualityLevels_) {
      this.qualityLevels_.dispose();
    }

    if (this.player_) {
      delete this.player_.vhs;
      delete this.player_.dash;
      delete this.player_.hls;
    }

    if (this.tech_ && this.tech_.vhs) {
      delete this.tech_.vhs;
    } // don't check this.tech_.hls as it will log a deprecated warning


    if (this.tech_) {
      delete this.tech_.hls;
    }

    if (this.mediaSourceUrl_ && window_1.URL.revokeObjectURL) {
      window_1.URL.revokeObjectURL(this.mediaSourceUrl_);
      this.mediaSourceUrl_ = null;
    }

    _Component.prototype.dispose.call(this);
  };

  _proto.convertToProgramTime = function convertToProgramTime(time, callback) {
    return getProgramTime({
      playlist: this.masterPlaylistController_.media(),
      time: time,
      callback: callback
    });
  } // the player must be playing before calling this
  ;

  _proto.seekToProgramTime = function seekToProgramTime$1(programTime, callback, pauseAfterSeek, retryCount) {
    if (pauseAfterSeek === void 0) {
      pauseAfterSeek = true;
    }

    if (retryCount === void 0) {
      retryCount = 2;
    }

    return seekToProgramTime({
      programTime: programTime,
      playlist: this.masterPlaylistController_.media(),
      retryCount: retryCount,
      pauseAfterSeek: pauseAfterSeek,
      seekTo: this.options_.seekTo,
      tech: this.options_.tech,
      callback: callback
    });
  };

  return VhsHandler;
}(Component);
/**
 * The Source Handler object, which informs video.js what additional
 * MIME types are supported and sets up playback. It is registered
 * automatically to the appropriate tech based on the capabilities of
 * the browser it is running in. It is not necessary to use or modify
 * this object in normal usage.
 */


var VhsSourceHandler = {
  name: 'videojs-http-streaming',
  VERSION: version$4,
  canHandleSource: function canHandleSource(srcObj, options) {
    if (options === void 0) {
      options = {};
    }

    var localOptions = videojs$1.mergeOptions(videojs$1.options, options);
    return VhsSourceHandler.canPlayType(srcObj.type, localOptions);
  },
  handleSource: function handleSource(source, tech, options) {
    if (options === void 0) {
      options = {};
    }

    var localOptions = videojs$1.mergeOptions(videojs$1.options, options);
    tech.vhs = new VhsHandler(source, tech, localOptions);

    if (!videojs$1.hasOwnProperty('hls')) {
      Object.defineProperty(tech, 'hls', {
        get: function get() {
          videojs$1.log.warn('player.tech().hls is deprecated. Use player.tech().vhs instead.');
          return tech.vhs;
        },
        configurable: true
      });
    }

    tech.vhs.xhr = xhrFactory();
    tech.vhs.src(source.src, source.type);
    return tech.vhs;
  },
  canPlayType: function canPlayType(type, options) {
    if (options === void 0) {
      options = {};
    }

    var _videojs$mergeOptions = videojs$1.mergeOptions(videojs$1.options, options),
        _videojs$mergeOptions2 = _videojs$mergeOptions.vhs.overrideNative,
        overrideNative = _videojs$mergeOptions2 === void 0 ? !videojs$1.browser.IS_ANY_SAFARI : _videojs$mergeOptions2;

    var supportedType = mediaTypes.simpleTypeFromSourceType(type);
    var canUseMsePlayback = supportedType && (!Vhs$1.supportsTypeNatively(supportedType) || overrideNative);
    return canUseMsePlayback ? 'maybe' : '';
  }
};
/**
 * Check to see if the native MediaSource object exists and supports
 * an MP4 container with both H.264 video and AAC-LC audio.
 *
 * @return {boolean} if  native media sources are supported
 */

var supportsNativeMediaSources = function supportsNativeMediaSources() {
  return codecs.browserSupportsCodec('avc1.4d400d,mp4a.40.2');
}; // register source handlers with the appropriate techs


if (supportsNativeMediaSources()) {
  videojs$1.getTech('Html5').registerSourceHandler(VhsSourceHandler, 0);
}

videojs$1.VhsHandler = VhsHandler;
Object.defineProperty(videojs$1, 'HlsHandler', {
  get: function get() {
    videojs$1.log.warn('videojs.HlsHandler is deprecated. Use videojs.VhsHandler instead.');
    return VhsHandler;
  },
  configurable: true
});
videojs$1.VhsSourceHandler = VhsSourceHandler;
Object.defineProperty(videojs$1, 'HlsSourceHandler', {
  get: function get() {
    videojs$1.log.warn('videojs.HlsSourceHandler is deprecated. ' + 'Use videojs.VhsSourceHandler instead.');
    return VhsSourceHandler;
  },
  configurable: true
});
videojs$1.Vhs = Vhs$1;
Object.defineProperty(videojs$1, 'Hls', {
  get: function get() {
    videojs$1.log.warn('videojs.Hls is deprecated. Use videojs.Vhs instead.');
    return Vhs$1;
  },
  configurable: true
});

if (!videojs$1.use) {
  videojs$1.registerComponent('Hls', Vhs$1);
  videojs$1.registerComponent('Vhs', Vhs$1);
}

videojs$1.options.vhs = videojs$1.options.vhs || {};
videojs$1.options.hls = videojs$1.options.hls || {};

if (videojs$1.registerPlugin) {
  videojs$1.registerPlugin('reloadSourceOnError', reloadSourceOnError);
} else {
  videojs$1.plugin('reloadSourceOnError', reloadSourceOnError);
}

var simpleTypeFromSourceType = mediaTypes.simpleTypeFromSourceType;
export { LOCAL_STORAGE_KEY, Vhs$1 as Vhs, VhsHandler, VhsSourceHandler, emeKeySystems, expandDataUri, getAllPsshKeySystemsOptions, setupEmeOptions, simpleTypeFromSourceType };
