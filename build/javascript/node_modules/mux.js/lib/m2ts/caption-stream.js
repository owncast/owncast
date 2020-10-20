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

'use strict';

// -----------------
// Link To Transport
// -----------------

var Stream = require('../utils/stream');
var cea708Parser = require('../tools/caption-packet-parser');

var CaptionStream = function() {

  CaptionStream.prototype.init.call(this);

  this.captionPackets_ = [];

  this.ccStreams_ = [
    new Cea608Stream(0, 0), // eslint-disable-line no-use-before-define
    new Cea608Stream(0, 1), // eslint-disable-line no-use-before-define
    new Cea608Stream(1, 0), // eslint-disable-line no-use-before-define
    new Cea608Stream(1, 1) // eslint-disable-line no-use-before-define
  ];

  this.reset();

  // forward data and done events from CCs to this CaptionStream
  this.ccStreams_.forEach(function(cc) {
    cc.on('data', this.trigger.bind(this, 'data'));
    cc.on('partialdone', this.trigger.bind(this, 'partialdone'));
    cc.on('done', this.trigger.bind(this, 'done'));
  }, this);

};

CaptionStream.prototype = new Stream();
CaptionStream.prototype.push = function(event) {
  var sei, userData, newCaptionPackets;

  // only examine SEI NALs
  if (event.nalUnitType !== 'sei_rbsp') {
    return;
  }

  // parse the sei
  sei = cea708Parser.parseSei(event.escapedRBSP);

  // ignore everything but user_data_registered_itu_t_t35
  if (sei.payloadType !== cea708Parser.USER_DATA_REGISTERED_ITU_T_T35) {
    return;
  }

  // parse out the user data payload
  userData = cea708Parser.parseUserData(sei);

  // ignore unrecognized userData
  if (!userData) {
    return;
  }

  // Sometimes, the same segment # will be downloaded twice. To stop the
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
  } else if ((event.dts === this.latestDts_) && (this.ignoreNextEqualDts_)) {
    this.numSameDts_--;
    if (!this.numSameDts_) {
      // We've received the last duplicate packet, time to start processing again
      this.ignoreNextEqualDts_ = false;
    }
    return;
  }

  // parse out CC data packets and save them for later
  newCaptionPackets = cea708Parser.parseCaptionPackets(event.pts, userData);
  this.captionPackets_ = this.captionPackets_.concat(newCaptionPackets);
  if (this.latestDts_ !== event.dts) {
    this.numSameDts_ = 0;
  }
  this.numSameDts_++;
  this.latestDts_ = event.dts;
};

CaptionStream.prototype.flushCCStreams = function(flushType) {
  this.ccStreams_.forEach(function(cc) {
    return flushType === 'flush' ? cc.flush() : cc.partialFlush();
  }, this);
};

CaptionStream.prototype.flushStream = function(flushType) {
  // make sure we actually parsed captions before proceeding
  if (!this.captionPackets_.length) {
    this.flushCCStreams(flushType);
    return;
  }

  // In Chrome, the Array#sort function is not stable so add a
  // presortIndex that we can use to ensure we get a stable-sort
  this.captionPackets_.forEach(function(elem, idx) {
    elem.presortIndex = idx;
  });

  // sort caption byte-pairs based on their PTS values
  this.captionPackets_.sort(function(a, b) {
    if (a.pts === b.pts) {
      return a.presortIndex - b.presortIndex;
    }
    return a.pts - b.pts;
  });

  this.captionPackets_.forEach(function(packet) {
    if (packet.type < 2) {
      // Dispatch packet to the right Cea608Stream
      this.dispatchCea608Packet(packet);
    }
    // this is where an 'else' would go for a dispatching packets
    // to a theoretical Cea708Stream that handles SERVICEn data
  }, this);

  this.captionPackets_.length = 0;
  this.flushCCStreams(flushType);
};

CaptionStream.prototype.flush = function() {
  return this.flushStream('flush');
};

// Only called if handling partial data
CaptionStream.prototype.partialFlush = function() {
  return this.flushStream('partialFlush');
};

CaptionStream.prototype.reset = function() {
  this.latestDts_ = null;
  this.ignoreNextEqualDts_ = false;
  this.numSameDts_ = 0;
  this.activeCea608Channel_ = [null, null];
  this.ccStreams_.forEach(function(ccStream) {
    ccStream.reset();
  });
};

// From the CEA-608 spec:
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
CaptionStream.prototype.dispatchCea608Packet = function(packet) {
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

CaptionStream.prototype.setsChannel1Active = function(packet) {
  return ((packet.ccData & 0x7800) === 0x1000);
};
CaptionStream.prototype.setsChannel2Active = function(packet) {
  return ((packet.ccData & 0x7800) === 0x1800);
};
CaptionStream.prototype.setsTextOrXDSActive = function(packet) {
  return ((packet.ccData & 0x7100) === 0x0100) ||
    ((packet.ccData & 0x78fe) === 0x102a) ||
    ((packet.ccData & 0x78fe) === 0x182a);
};

// ----------------------
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
  0x2a: 0xe1,     // á
  0x5c: 0xe9,     // é
  0x5e: 0xed,     // í
  0x5f: 0xf3,     // ó
  0x60: 0xfa,     // ú
  0x7b: 0xe7,     // ç
  0x7c: 0xf7,     // ÷
  0x7d: 0xd1,     // Ñ
  0x7e: 0xf1,     // ñ
  0x7f: 0x2588,   // █
  0x0130: 0xae,   // ®
  0x0131: 0xb0,   // °
  0x0132: 0xbd,   // ½
  0x0133: 0xbf,   // ¿
  0x0134: 0x2122, // ™
  0x0135: 0xa2,   // ¢
  0x0136: 0xa3,   // £
  0x0137: 0x266a, // ♪
  0x0138: 0xe0,   // à
  0x0139: 0xa0,   //
  0x013a: 0xe8,   // è
  0x013b: 0xe2,   // â
  0x013c: 0xea,   // ê
  0x013d: 0xee,   // î
  0x013e: 0xf4,   // ô
  0x013f: 0xfb,   // û
  0x0220: 0xc1,   // Á
  0x0221: 0xc9,   // É
  0x0222: 0xd3,   // Ó
  0x0223: 0xda,   // Ú
  0x0224: 0xdc,   // Ü
  0x0225: 0xfc,   // ü
  0x0226: 0x2018, // ‘
  0x0227: 0xa1,   // ¡
  0x0228: 0x2a,   // *
  0x0229: 0x27,   // '
  0x022a: 0x2014, // —
  0x022b: 0xa9,   // ©
  0x022c: 0x2120, // ℠
  0x022d: 0x2022, // •
  0x022e: 0x201c, // “
  0x022f: 0x201d, // ”
  0x0230: 0xc0,   // À
  0x0231: 0xc2,   // Â
  0x0232: 0xc7,   // Ç
  0x0233: 0xc8,   // È
  0x0234: 0xca,   // Ê
  0x0235: 0xcb,   // Ë
  0x0236: 0xeb,   // ë
  0x0237: 0xce,   // Î
  0x0238: 0xcf,   // Ï
  0x0239: 0xef,   // ï
  0x023a: 0xd4,   // Ô
  0x023b: 0xd9,   // Ù
  0x023c: 0xf9,   // ù
  0x023d: 0xdb,   // Û
  0x023e: 0xab,   // «
  0x023f: 0xbb,   // »
  0x0320: 0xc3,   // Ã
  0x0321: 0xe3,   // ã
  0x0322: 0xcd,   // Í
  0x0323: 0xcc,   // Ì
  0x0324: 0xec,   // ì
  0x0325: 0xd2,   // Ò
  0x0326: 0xf2,   // ò
  0x0327: 0xd5,   // Õ
  0x0328: 0xf5,   // õ
  0x0329: 0x7b,   // {
  0x032a: 0x7d,   // }
  0x032b: 0x5c,   // \
  0x032c: 0x5e,   // ^
  0x032d: 0x5f,   // _
  0x032e: 0x7c,   // |
  0x032f: 0x7e,   // ~
  0x0330: 0xc4,   // Ä
  0x0331: 0xe4,   // ä
  0x0332: 0xd6,   // Ö
  0x0333: 0xf6,   // ö
  0x0334: 0xdf,   // ß
  0x0335: 0xa5,   // ¥
  0x0336: 0xa4,   // ¤
  0x0337: 0x2502, // │
  0x0338: 0xc5,   // Å
  0x0339: 0xe5,   // å
  0x033a: 0xd8,   // Ø
  0x033b: 0xf8,   // ø
  0x033c: 0x250c, // ┌
  0x033d: 0x2510, // ┐
  0x033e: 0x2514, // └
  0x033f: 0x2518  // ┘
};

var getCharFromCode = function(code) {
  if (code === null) {
    return '';
  }
  code = CHARACTER_TRANSLATION[code] || code;
  return String.fromCharCode(code);
};

// the index of the last row in a CEA-608 display buffer
var BOTTOM_ROW = 14;

// This array is used for mapping PACs -> row #, since there's no way of
// getting it through bit logic.
var ROWS = [0x1100, 0x1120, 0x1200, 0x1220, 0x1500, 0x1520, 0x1600, 0x1620,
            0x1700, 0x1720, 0x1000, 0x1300, 0x1320, 0x1400, 0x1420];

// CEA-608 captions are rendered onto a 34x15 matrix of character
// cells. The "bottom" row is the last element in the outer array.
var createDisplayBuffer = function() {
  var result = [], i = BOTTOM_ROW + 1;
  while (i--) {
    result.push('');
  }
  return result;
};

var Cea608Stream = function(field, dataChannel) {
  Cea608Stream.prototype.init.call(this);

  this.field_ = field || 0;
  this.dataChannel_ = dataChannel || 0;

  this.name_ = 'CC' + (((this.field_ << 1) | this.dataChannel_) + 1);

  this.setConstants();
  this.reset();

  this.push = function(packet) {
    var data, swap, char0, char1, text;
    // remove the parity bits
    data = packet.ccData & 0x7f7f;

    // ignore duplicate control codes; the spec demands they're sent twice
    if (data === this.lastControlCode_) {
      this.lastControlCode_ = null;
      return;
    }

    // Store control codes
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
      this.clearFormatting(packet.pts);
      // if a caption was being displayed, it's gone now
      this.flushDisplayed(packet.pts);

      // flip memory
      swap = this.displayed_;
      this.displayed_ = this.nonDisplayed_;
      this.nonDisplayed_ = swap;

      // start measuring the time to display the caption
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
      this.startPts_ = packet.pts;

    // Append special characters to caption text
    } else if (this.isSpecialCharacter(char0, char1)) {
      // Bitmask char0 so that we can apply character transformations
      // regardless of field and data channel.
      // Then byte-shift to the left and OR with char1 so we can pass the
      // entire character code to `getCharFromCode`.
      char0 = (char0 & 0x03) << 8;
      text = getCharFromCode(char0 | char1);
      this[this.mode_](packet.pts, text);
      this.column_++;

    // Append extended characters to caption text
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
      }

      // Bitmask char0 so that we can apply character transformations
      // regardless of field and data channel.
      // Then byte-shift to the left and OR with char1 so we can pass the
      // entire character code to `getCharFromCode`.
      char0 = (char0 & 0x03) << 8;
      text = getCharFromCode(char0 | char1);
      this[this.mode_](packet.pts, text);
      this.column_++;

    // Process mid-row codes
    } else if (this.isMidRowCode(char0, char1)) {
      // Attributes are not additive, so clear all formatting
      this.clearFormatting(packet.pts);

      // According to the standard, mid-row codes
      // should be replaced with spaces, so add one now
      this[this.mode_](packet.pts, ' ');
      this.column_++;

      if ((char1 & 0xe) === 0xe) {
        this.addFormatting(packet.pts, ['i']);
      }

      if ((char1 & 0x1) === 0x1) {
        this.addFormatting(packet.pts, ['u']);
      }

    // Detect offset control codes and adjust cursor
    } else if (this.isOffsetControlCode(char0, char1)) {
      // Cursor position is set by indent PAC (see below) in 4-column
      // increments, with an additional offset code of 1-3 to reach any
      // of the 32 columns specified by CEA-608. So all we need to do
      // here is increment the column cursor by the given offset.
      this.column_ += (char1 & 0x03);

    // Detect PACs (Preamble Address Codes)
    } else if (this.isPAC(char0, char1)) {

      // There's no logic for PAC -> row mapping, so we have to just
      // find the row code in an array and use its index :(
      var row = ROWS.indexOf(data & 0x1f20);

      // Configure the caption window if we're in roll-up mode
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
      }
      // All PACs can apply underline, so detect and apply
      // (All odd-numbered second bytes set underline)
      if ((char1 & 0x1) && (this.formatting_.indexOf('u') === -1)) {
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
      }

    // We have a normal character in char0, and possibly one in char1
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
Cea608Stream.prototype = new Stream();
// Trigger a cue point that captures the current state of the
// display buffer
Cea608Stream.prototype.flushDisplayed = function(pts) {
  var content = this.displayed_
    // remove spaces from the start and end of the string
    .map(function(row) {
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
    })
    // combine all text rows to display in one cue
    .join('\n')
    // and remove blank rows from the start and end, but not the middle
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
Cea608Stream.prototype.reset = function() {
  this.mode_ = 'popOn';
  // When in roll-up mode, the index of the last row that will
  // actually display captions. If a caption is shifted to a row
  // with a lower index than this, it is cleared from the display
  // buffer
  this.topRow_ = 0;
  this.startPts_ = 0;
  this.displayed_ = createDisplayBuffer();
  this.nonDisplayed_ = createDisplayBuffer();
  this.lastControlCode_ = null;

  // Track row and column for proper line-breaking and spacing
  this.column_ = 0;
  this.row_ = BOTTOM_ROW;
  this.rollUpRows_ = 2;

  // This variable holds currently-applied formatting
  this.formatting_ = [];
};

/**
 * Sets up control code and related constants for this instance
 */
Cea608Stream.prototype.setConstants = function() {
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
    this.BASE_     = 0x10;
    this.EXT_      = 0x11;
    this.CONTROL_  = (0x14 | this.field_) << 8;
    this.OFFSET_   = 0x17;
  } else if (this.dataChannel_ === 1) {
    this.BASE_     = 0x18;
    this.EXT_      = 0x19;
    this.CONTROL_  = (0x1c | this.field_) << 8;
    this.OFFSET_   = 0x1f;
  }

  // Constants for the LSByte command codes recognized by Cea608Stream. This
  // list is not exhaustive. For a more comprehensive listing and semantics see
  // http://www.gpo.gov/fdsys/pkg/CFR-2010-title47-vol1/pdf/CFR-2010-title47-vol1-sec15-119.pdf
  // Padding
  this.PADDING_                    = 0x0000;
  // Pop-on Mode
  this.RESUME_CAPTION_LOADING_     = this.CONTROL_ | 0x20;
  this.END_OF_CAPTION_             = this.CONTROL_ | 0x2f;
  // Roll-up Mode
  this.ROLL_UP_2_ROWS_             = this.CONTROL_ | 0x25;
  this.ROLL_UP_3_ROWS_             = this.CONTROL_ | 0x26;
  this.ROLL_UP_4_ROWS_             = this.CONTROL_ | 0x27;
  this.CARRIAGE_RETURN_            = this.CONTROL_ | 0x2d;
  // paint-on mode
  this.RESUME_DIRECT_CAPTIONING_   = this.CONTROL_ | 0x29;
  // Erasure
  this.BACKSPACE_                  = this.CONTROL_ | 0x21;
  this.ERASE_DISPLAYED_MEMORY_     = this.CONTROL_ | 0x2c;
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
Cea608Stream.prototype.isSpecialCharacter = function(char0, char1) {
  return (char0 === this.EXT_ && char1 >= 0x30 && char1 <= 0x3f);
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
Cea608Stream.prototype.isExtCharacter = function(char0, char1) {
  return ((char0 === (this.EXT_ + 1) || char0 === (this.EXT_ + 2)) &&
    (char1 >= 0x20 && char1 <= 0x3f));
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
Cea608Stream.prototype.isMidRowCode = function(char0, char1) {
  return (char0 === this.EXT_ && (char1 >= 0x20 && char1 <= 0x2f));
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
Cea608Stream.prototype.isOffsetControlCode = function(char0, char1) {
  return (char0 === this.OFFSET_ && (char1 >= 0x21 && char1 <= 0x23));
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
Cea608Stream.prototype.isPAC = function(char0, char1) {
  return (char0 >= this.BASE_ && char0 < (this.BASE_ + 8) &&
    (char1 >= 0x40 && char1 <= 0x7f));
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
Cea608Stream.prototype.isColorPAC = function(char1) {
  return ((char1 >= 0x40 && char1 <= 0x4f) || (char1 >= 0x60 && char1 <= 0x7f));
};

/**
 * Detects if a single byte is in the range of a normal character
 *
 * Normal text bytes are in the range 0x20 to 0x7f.
 *
 * @param  {Integer} char  The byte
 * @return {Boolean}       Whether the byte is a normal character
 */
Cea608Stream.prototype.isNormalChar = function(char) {
  return (char >= 0x20 && char <= 0x7f);
};

/**
 * Configures roll-up
 *
 * @param  {Integer} pts         Current PTS
 * @param  {Integer} newBaseRow  Used by PACs to slide the current window to
 *                               a new position
 */
Cea608Stream.prototype.setRollUp = function(pts, newBaseRow) {
  // Reset the base row to the bottom row when switching modes
  if (this.mode_ !== 'rollUp') {
    this.row_ = BOTTOM_ROW;
    this.mode_ = 'rollUp';
    // Spec says to wipe memories when switching to roll-up
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
};

// Adds the opening HTML tag for the passed character to the caption text,
// and keeps track of it for later closing
Cea608Stream.prototype.addFormatting = function(pts, format) {
  this.formatting_ = this.formatting_.concat(format);
  var text = format.reduce(function(text, format) {
    return text + '<' + format + '>';
  }, '');
  this[this.mode_](pts, text);
};

// Adds HTML closing tags for current formatting to caption text and
// clears remembered formatting
Cea608Stream.prototype.clearFormatting = function(pts) {
  if (!this.formatting_.length) {
    return;
  }
  var text = this.formatting_.reverse().reduce(function(text, format) {
    return text + '</' + format + '>';
  }, '');
  this.formatting_ = [];
  this[this.mode_](pts, text);
};

// Mode Implementations
Cea608Stream.prototype.popOn = function(pts, text) {
  var baseRow = this.nonDisplayed_[this.row_];

  // buffer characters
  baseRow += text;
  this.nonDisplayed_[this.row_] = baseRow;
};

Cea608Stream.prototype.rollUp = function(pts, text) {
  var baseRow = this.displayed_[this.row_];

  baseRow += text;
  this.displayed_[this.row_] = baseRow;

};

Cea608Stream.prototype.shiftRowsUp_ = function() {
  var i;
  // clear out inactive rows
  for (i = 0; i < this.topRow_; i++) {
    this.displayed_[i] = '';
  }
  for (i = this.row_ + 1; i < BOTTOM_ROW + 1; i++) {
    this.displayed_[i] = '';
  }
  // shift displayed rows up
  for (i = this.topRow_; i < this.row_; i++) {
    this.displayed_[i] = this.displayed_[i + 1];
  }
  // clear out the bottom row
  this.displayed_[this.row_] = '';
};

Cea608Stream.prototype.paintOn = function(pts, text) {
  var baseRow = this.displayed_[this.row_];

  baseRow += text;
  this.displayed_[this.row_] = baseRow;
};

// exports
module.exports = {
  CaptionStream: CaptionStream,
  Cea608Stream: Cea608Stream
};
