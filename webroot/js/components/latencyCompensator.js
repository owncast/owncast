/*
The Owncast Latency Compensator.
It will try to slowly adjust the playback rate to enable the player to get
further into the future, with the goal of being as close to the live edge as
possible, without causing any buffering events.

It will:
  - Determine the start (max) and end (min) latency values.
  - Keep an eye on download speed and stop compensating if it drops too low.
  - Pause the compensation if buffering events occur.
  - Completely give up on all compensation if too many buffering events occur.
*/

const REBUFFER_EVENT_LIMIT = 5; // Max number of buffering events before we stop compensating for latency.
const MIN_BUFFER_DURATION = 300; // Min duration a buffer event must last to be counted.
const MAX_SPEEDUP_RATE = 1.08; // The playback rate when compensating for latency.
const TIMEOUT_DURATION = 20_000; // The amount of time we stop handling latency after certain events.
const CHECK_TIMER_INTERVAL = 5_000; // How often we check if we should be compensating for latency.
const BUFFERING_AMNESTY_DURATION = 4 * 1000 * 60; // How often until a buffering event expires.
const HIGH_LATENCY_ENABLE_THRESHOLD = 20_000; // ms;
const LOW_LATENCY_DISABLE_THRESHOLD = 4500; // ms;
const REQUIRED_BANDWIDTH_RATIO = 2.0; // The player:bitrate ratio required to enable compensating for latency.

class LatencyCompensator {
  constructor(player) {
    this.player = player;
    this.enabled = false;
    this.running = false;
    this.inTimeout = false;
    this.timeoutTimer = 0;
    this.checkTimer = 0;
    this.maxLatencyThreshold = 10000;
    this.minLatencyThreshold = 8000;
    this.bufferingCounter = 0;
    this.bufferingTimer = 0;
    this.playbackRate = 1.0;
    this.player.on('playing', this.handlePlaying.bind(this));
    this.player.on('error', this.handleError.bind(this));
    this.player.on('waiting', this.handleBuffering.bind(this));
    this.player.on('ended', this.handleEnded.bind(this));
    this.player.on('canplaythrough', this.handlePlaying.bind(this));
    this.player.on('canplay', this.handlePlaying.bind(this));
  }

  // This is run on a timer to check if we should be compensating for latency.
  check() {
    console.log(
      'playback rate',
      this.playbackRate,
      'enabled:',
      this.enabled,
      'running: ',
      this.running,
      'timeout: ',
      this.inTimeout,
      'buffer count:',
      this.bufferingCounter
    );

    if (this.inTimeout) {
      return;
    }

    const tech = this.player.tech({ IWillNotUseThisInPlugins: true });

    if (!tech || !tech.vhs) {
      return;
    }

    try {
      // Check the player buffers to make sure there's enough playable content
      // that we can safely play.
      if (tech.vhs.stats.buffered.length === 0) {
        this.timeout();
      }

      let totalBuffered = 0;

      tech.vhs.stats.buffered.forEach((buffer) => {
        totalBuffered += buffer.end - buffer.start;
      });
      console.log('buffered', totalBuffered);

      if (totalBuffered < 18) {
        this.timeout();
      }
    } catch (e) {}

    // Determine how much of the current playlist's bandwidth requirements
    // we're utilizing. If it's too high then we can't afford to push
    // further into the future because we're downloading too slowly.
    const currentPlaylist = tech.vhs.playlists.media();
    const currentPlaylistBandwidth = currentPlaylist.attributes.BANDWIDTH;
    const playerBandwidth = tech.vhs.systemBandwidth;
    const bandwidthRatio = playerBandwidth / currentPlaylistBandwidth;

    if (bandwidthRatio < REQUIRED_BANDWIDTH_RATIO) {
      this.timeout();
      return;
    }

    let proposedPlaybackRate = bandwidthRatio * 0.34;
    console.log('proposed rate', proposedPlaybackRate);

    proposedPlaybackRate = Math.max(
      Math.min(proposedPlaybackRate, MAX_SPEEDUP_RATE),
      1.0
    );

    try {
      const segment = getCurrentlyPlayingSegment(tech);
      if (!segment) {
        return;
      }

      // How far away from live edge do we start the compensator.
      this.maxLatencyThreshold = Math.min(
        segment.duration * 1000 * 4,
        HIGH_LATENCY_ENABLE_THRESHOLD
      );

      // How far away from live edge do we stop the compensator.
      this.minLatencyThreshold = Math.max(
        segment.duration * 1000 * 2,
        LOW_LATENCY_DISABLE_THRESHOLD
      );
      const segmentTime = segment.dateTimeObject.getTime();
      const now = new Date().getTime();
      const latency = now - segmentTime;

      if (latency > this.maxLatencyThreshold) {
        this.start(proposedPlaybackRate);
      } else if (latency < this.minLatencyThreshold) {
        this.stop();
      }
    } catch (err) {
      console.warn(err);
    }
  }

  setPlaybackRate(rate) {
    this.playbackRate = rate;
    this.player.playbackRate(rate);
  }

  start(rate = 1.0) {
    if (this.inTimeout || !this.enabled || rate === this.playbackRate) {
      return;
    }

    this.running = true;
    this.setPlaybackRate(rate);
  }

  stop() {
    this.running = false;
    this.setPlaybackRate(1.0);
  }

  enable() {
    this.enabled = true;
    clearInterval(this.checkTimer);
    clearTimeout(this.bufferingTimer);

    this.checkTimer = setInterval(() => {
      this.check();
    }, CHECK_TIMER_INTERVAL);
  }

  // Disable means we're done for good and should no longer compensate for latency.
  disable() {
    clearInterval(this.checkTimer);
    clearTimeout(this.timeoutTimer);
    this.stop();
    this.enabled = false;
  }

  timeout() {
    this.inTimeout = true;
    this.stop();

    clearTimeout(this.timeoutTimer);
    this.timeoutTimer = setTimeout(() => {
      this.endTimeout();
    }, TIMEOUT_DURATION);
  }

  endTimeout() {
    clearTimeout(this.timeoutTimer);
    this.inTimeout = false;
  }

  handlePlaying() {
    clearTimeout(this.bufferingTimer);

    if (!this.enabled) {
      return;
    }
  }

  handleEnded() {
    if (!this.enabled) {
      return;
    }

    this.disable();
  }

  handleError(e) {
    if (!this.enabled) {
      return;
    }

    console.log('handle error', e);
    this.timeout();
  }

  countBufferingEvent() {
    this.bufferingCounter++;
    if (this.bufferingCounter > REBUFFER_EVENT_LIMIT) {
      this.disable();
      return;
    }

    this.timeout();

    // Allow us to forget about old buffering events if enough time goes by.
    setTimeout(() => {
      if (this.bufferingCounter > 0) {
        this.bufferingCounter--;
      }
    }, BUFFERING_AMNESTY_DURATION);
  }

  handleBuffering() {
    if (!this.enabled) {
      return;
    }

    this.timeout();

    this.bufferingTimer = setTimeout(() => {
      this.countBufferingEvent();
    }, MIN_BUFFER_DURATION);
  }
}

function getCurrentlyPlayingSegment(tech) {
  var target_media = tech.vhs.playlists.media();
  var snapshot_time = tech.currentTime();

  var segment;

  // Itinerate trough available segments and get first within which snapshot_time is
  for (var i = 0, l = target_media.segments.length; i < l; i++) {
    // Note: segment.end may be undefined or is not properly set
    if (snapshot_time < target_media.segments[i].end) {
      segment = target_media.segments[i];
      break;
    }
  }

  if (!segment) {
    segment = target_media.segments[0];
  }

  return segment;
}

export default LatencyCompensator;
