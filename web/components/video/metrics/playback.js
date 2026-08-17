const URL_CMCD_COLLECTOR = `/api/metrics/cmcd`;
const EVENT_REPORT_INTERVAL = 10000;
const MAX_VALID_LATENCY_SECONDS = 100; // Anything > this gets thrown out.
const CMCD_VERSION = 2;

function getCurrentlyPlayingSegment(tech) {
  const targetMedia = tech.vhs.playlists.media();
  const snapshotTime = tech.currentTime();
  let segment;

  // Iterate through available segments and get first within which snapshot_time is
  for (let i = 0, l = targetMedia.segments.length; i < l; i += 1) {
    // Note: segment.end may be undefined or is not properly set
    if (snapshotTime < targetMedia.segments[i].end) {
      segment = targetMedia.segments[i];
      break;
    }
  }
  // No match means currentTime is outside the tracked window (mid-seek, or
  // segment bookkeeping is stale after a live-edge resync). Report no segment
  // and skip the sample: falling back to the oldest segment inflates the
  // reported latency by up to the whole playlist window.

  return segment;
}

function generateSessionID() {
  if (crypto.randomUUID) {
    return crypto.randomUUID();
  }
  // randomUUID is only exposed in secure contexts; build a v4 UUID from
  // getRandomValues, which is available everywhere.
  const bytes = crypto.getRandomValues(new Uint8Array(16));
  // eslint-disable-next-line no-bitwise
  bytes[6] = (bytes[6] & 0x0f) | 0x40;
  // eslint-disable-next-line no-bitwise
  bytes[8] = (bytes[8] & 0x3f) | 0x80;
  const hex = [...bytes].map(b => b.toString(16).padStart(2, '0')).join('');
  return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20)}`;
}

// CMCD keys whose values are tokens, which the spec serializes unquoted
// (ot=v), unlike string values which are quoted (sid="...").
const CMCD_TOKEN_KEYS = new Set(['ot', 'sf', 'st', 'sta', 'e']);

// Serializes CMCD keys as a CMCD dictionary payload: alphabetized,
// booleans as bare keys when true, tokens bare, strings quoted with
// escaping.
function encodeCmcdDictionary(keys) {
  return Object.keys(keys)
    .filter(key => keys[key] !== undefined && keys[key] !== null && keys[key] !== false)
    .sort()
    .map(key => {
      const value = keys[key];
      if (value === true) {
        return key;
      }
      if (typeof value === 'string' && !CMCD_TOKEN_KEYS.has(key)) {
        return `${key}="${value.replace(/\\/g, '\\\\').replace(/"/g, '\\"')}"`;
      }
      return `${key}=${value}`;
    })
    .join(',');
}

function appendCmcdQuery(uri, keys) {
  const payload = encodeCmcdDictionary(keys);
  if (!payload) {
    return uri;
  }
  const separator = uri.includes('?') ? '&' : '?';
  return `${uri}${separator}CMCD=${encodeURIComponent(payload)}`;
}

/*
CMCD v2 (CTA-5004-A) metrics emitter.

Request mode: every playlist and segment request the player makes is
decorated with a CMCD query parameter carrying throughput, buffer, bitrate
and live latency data, which the server harvests off the requests it is
already serving.

Event mode: player state transitions, errors and a periodic time-interval
report are POSTed to the CMCD collector endpoint.
*/
class PlaybackMetrics {
  constructor(player, videojs) {
    this.player = player;
    this.clockSkewMs = 0;
    this.sessionID = generateSessionID();
    this.sequenceNumber = 0;

    this.hasStartedPlaying = false;
    this.isStalled = false;
    this.playRequestedAt = null;
    this.mediaStartDelayMs = null;
    this.mediaStartDelaySent = false;
    this.segmentDownloadTimesMs = [];

    this.handlePlay = this.handlePlay.bind(this);
    this.handlePlaying = this.handlePlaying.bind(this);
    this.handlePause = this.handlePause.bind(this);
    this.handleEnded = this.handleEnded.bind(this);
    this.handleError = this.handleError.bind(this);
    this.handleBuffering = this.handleBuffering.bind(this);
    this.handleNoLongerBuffering = this.handleNoLongerBuffering.bind(this);

    this.player.on('play', this.handlePlay);
    this.player.on('playing', this.handlePlaying);
    this.player.on('pause', this.handlePause);
    this.player.on('ended', this.handleEnded);
    this.player.on('error', this.handleError);
    this.player.on('stalled', this.handleBuffering);
    this.player.on('waiting', this.handleBuffering);
    this.player.on('canplaythrough', this.handleNoLongerBuffering);

    // Decorate every media request with CMCD request-mode data and time
    // segment downloads for ttlb reporting. Keep a reference of the
    // standard vjs xhr function and wrap it.
    const oldVjsXhrCallback = videojs.xhr;
    // eslint-disable-next-line no-param-reassign
    videojs.Vhs.xhr = (...args) => {
      try {
        const request = args[0];
        const isSegment = /\.ts(\?|$)/.test(request.uri);
        const isPlaylist = /\.m3u8(\?|$)/.test(request.uri);
        if (isSegment || isPlaylist) {
          // eslint-disable-next-line no-param-reassign
          request.uri = appendCmcdQuery(request.uri, {
            ...this.commonCmcdKeys(),
            ot: isSegment ? 'av' : 'm',
            sf: 'h',
            st: 'l',
            su: !this.hasStartedPlaying || undefined,
          });
        }

        // Time segment downloads client-side: reported as ttlb, the only
        // download duration measurement that stays honest when a proxy or
        // tunnel sits between the viewer and the server.
        if (isSegment && typeof args[1] === 'function') {
          const start = Date.now();
          const cb = args[1];
          // eslint-disable-next-line no-param-reassign
          args[1] = (request_, error, response) => {
            this.segmentDownloadTimesMs.push(Date.now() - start);
            cb(request_, error, response);
          };
        }
      } catch {
        // Never let metrics decoration break media loading.
      }

      return oldVjsXhrCallback(...args);
    };

    this.eventReportTimer = setInterval(() => {
      // Pauses are reported as their own state transition; a paused player
      // has nothing new to say in between.
      if (!this.player || this.player.paused() || !this.hasStartedPlaying) {
        return;
      }
      this.sendEventReport('t');
    }, EVENT_REPORT_INTERVAL);
  }

  stop() {
    clearInterval(this.eventReportTimer);
    // Component teardown does not reliably emit Video.js's ended event.
    // Preserve the terminal state while its viewer is still active.
    if (this.hasStartedPlaying && !this.player.ended()) {
      this.sendEventReport('ps', { sta: 'e' });
    }
    this.player.off();
  }

  // Keep our client clock in sync with the server clock to determine
  // accurate latency calculations.
  setClockSkew(skewMs) {
    this.clockSkewMs = skewMs;
  }

  handlePlay() {
    if (this.playRequestedAt === null) {
      this.playRequestedAt = Date.now();
    }
  }

  handlePlaying() {
    if (this.mediaStartDelayMs === null && this.playRequestedAt !== null) {
      this.mediaStartDelayMs = Date.now() - this.playRequestedAt;
    }
    const firstPlay = !this.hasStartedPlaying;
    this.hasStartedPlaying = true;
    this.isStalled = false;
    if (firstPlay) {
      this.sendEventReport('ps', { sta: 'p' });
    }
  }

  handlePause() {
    // Tearing down the player fires a pause; don't report it.
    if (!this.player.ended()) {
      this.sendEventReport('ps', { sta: 'a' });
    }
  }

  handleEnded() {
    this.sendEventReport('ps', { sta: 'e' });
  }

  handleError() {
    const error = this.player.error();
    this.sendEventReport('e', { ec: error ? String(error.code) : undefined });
  }

  handleBuffering() {
    // Buffering before the first frame is startup, not an interruption. A
    // waiting event that fires while the player is seeking is the fetch at
    // a seek discontinuity (the viewer scrubbed, or the latency compensator
    // jumped forward), not a playback interruption. Genuine starvation fires
    // with seeking() false, so only report those.
    if (!this.hasStartedPlaying || this.player.seeking() || this.isStalled) {
      return;
    }
    this.isStalled = true;
    this.sendEventReport('ps', { sta: 'w', bs: true });
  }

  handleNoLongerBuffering() {
    this.isStalled = false;
  }

  // CMCD keys shared by request-mode decoration and event reports.
  commonCmcdKeys() {
    const keys = {
      v: CMCD_VERSION,
      sid: this.sessionID,
    };

    const throughput = this.measuredThroughputKbps();
    if (throughput) {
      keys.mtp = throughput;
    }

    const bitrate = this.currentBitrateKbps();
    if (bitrate) {
      keys.br = bitrate;
    }

    const buffer = this.bufferLengthMs();
    if (buffer !== null) {
      keys.bl = buffer;
    }

    const latency = this.liveLatencyMs();
    if (latency !== null) {
      keys.ltc = latency;
    }

    const playbackRate = this.player.playbackRate();
    if (playbackRate !== 1) {
      keys.pr = playbackRate;
    }

    return keys;
  }

  async sendEventReport(event, extraKeys = {}) {
    this.sequenceNumber += 1;

    const report = {
      ...this.commonCmcdKeys(),
      ...extraKeys,
      e: event,
      ts: Date.now(),
      sn: this.sequenceNumber,
    };

    if (report.sta === undefined) {
      report.sta = this.playerStateToken();
    }

    if (!this.mediaStartDelaySent && this.mediaStartDelayMs !== null) {
      this.mediaStartDelaySent = true;
      report.msd = this.mediaStartDelayMs;
    }

    const droppedFrames = this.droppedFrameCount();
    if (droppedFrames) {
      report.df = droppedFrames;
    }

    // Average client-measured segment download time since the last report,
    // reported as ttlb in milliseconds.
    if (this.segmentDownloadTimesMs.length > 0) {
      const sum = this.segmentDownloadTimesMs.reduce((p, c) => p + c, 0);
      report.ttlb = Math.round(sum / this.segmentDownloadTimesMs.length);
      this.segmentDownloadTimesMs = [];
    }

    Object.keys(report).forEach(key => {
      if (report[key] === undefined || report[key] === null) {
        delete report[key];
      }
    });

    try {
      await fetch(URL_CMCD_COLLECTOR, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(report),
      });
    } catch (e) {
      console.error(e);
    }
  }

  playerStateToken() {
    if (this.player.seeking()) {
      return 'k';
    }
    if (this.player.ended()) {
      return 'e';
    }
    if (this.player.paused()) {
      return 'a';
    }
    if (this.isStalled) {
      return 'w';
    }
    return 'p';
  }

  measuredThroughputKbps() {
    const tech = this.player.tech({ IWillNotUseThisInPlugins: true });
    const bandwidth = tech?.vhs?.systemBandwidth;
    if (!bandwidth) {
      return null;
    }
    // The spec asks for throughput rounded to the nearest 100kbps.
    return Math.round(bandwidth / 1000 / 100) * 100;
  }

  currentBitrateKbps() {
    const tech = this.player.tech({ IWillNotUseThisInPlugins: true });
    const bandwidth = tech?.vhs?.playlists?.media()?.attributes?.BANDWIDTH;
    if (!bandwidth) {
      return null;
    }
    return Math.round(bandwidth / 1000);
  }

  bufferLengthMs() {
    try {
      const buffered = this.player.buffered();
      if (!buffered || buffered.length === 0) {
        return null;
      }
      const bufferMs = (buffered.end(buffered.length - 1) - this.player.currentTime()) * 1000;
      if (bufferMs < 0) {
        return null;
      }
      // The spec asks for buffer length rounded to the nearest 100ms.
      return Math.round(bufferMs / 100) * 100;
    } catch {
      return null;
    }
  }

  liveLatencyMs() {
    const tech = this.player.tech({ IWillNotUseThisInPlugins: true });
    if (!tech?.vhs) {
      return null;
    }

    try {
      const segment = getCurrentlyPlayingSegment(tech);
      if (!segment || !segment.dateTimeObject) {
        return null;
      }

      const segmentTime = segment.dateTimeObject.getTime();
      const now = Date.now() + this.clockSkewMs;
      const latency = now - segmentTime;

      // Throw away values that seem invalid.
      if (latency < 0 || latency / 1000 >= MAX_VALID_LATENCY_SECONDS) {
        return null;
      }

      return Math.round(latency);
    } catch {
      return null;
    }
  }

  droppedFrameCount() {
    const tech = this.player.tech({ IWillNotUseThisInPlugins: true });
    const el = tech?.el?.();
    if (!el?.getVideoPlaybackQuality) {
      return null;
    }
    return el.getVideoPlaybackQuality().droppedVideoFrames || null;
  }
}

export default PlaybackMetrics;
