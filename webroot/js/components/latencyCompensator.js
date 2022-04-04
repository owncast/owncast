/*
The Owncast Latency Compensator.
It will try to slowly adjust the playback rate to enable the player to get
further into the future, with the goal of being as close to the live edge as
possible, without causing any buffering events.

It will:
  - Determine the start (max) and end (min) latency values.
  - Keep an eye on download speed and stop compensating if it drops too low.
  - Limit the playback speedup so it doesn't sound weird by jumping rates.
  - Dynamically calculate the speedup rate based on network speed.
  - Pause the compensation if buffering events occur.
  - Completely give up on all compensation if too many buffering events occur.
*/

const REBUFFER_EVENT_LIMIT = 5; // Max number of buffering events before we stop compensating for latency.
const MIN_BUFFER_DURATION = 500; // Min duration a buffer event must last to be counted.
const MAX_SPEEDUP_RATE = 1.07; // The playback rate when compensating for latency.
const MAX_SPEEDUP_RAMP = 0.005; // The max amount we will increase the playback rate at once.
const TIMEOUT_DURATION = 30 * 1000; // The amount of time we stop handling latency after certain events.
const CHECK_TIMER_INTERVAL = 5_000; // How often we check if we should be compensating for latency.
const BUFFERING_AMNESTY_DURATION = 3 * 1000 * 60; // How often until a buffering event expires.
const REQUIRED_BANDWIDTH_RATIO = 2.0; // The player:bitrate ratio required to enable compensating for latency.
const HIGHEST_LATENCY_SEGMENT_LENGTH_MULTIPLIER = 2.5; // Segment length * this value is when we start compensating.
const LOWEST_LATENCY_SEGMENT_LENGTH_MULTIPLIER = 1.7; // Segment length * this value is when we stop compensating.
const MIN_LATENCY = 4 * 1000; // The lowest we'll continue compensation to be running at.
const MAX_LATENCY = 8 * 1000; // The highest we'll allow a target latency to be before we start compensating.
class LatencyCompensator {
  constructor(player) {
    this.player = player;
    this.enabled = false;
    this.running = false;
    this.inTimeout = false;
    this.timeoutTimer = 0;
    this.checkTimer = 0;
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

    // If we don't think we have the bandwidth to play faster, then don't do it.
    if (bandwidthRatio < REQUIRED_BANDWIDTH_RATIO) {
      this.timeout();
      return;
    }

    // Using our bandwidth ratio determine a wide guess at how fast we can play.
    let proposedPlaybackRate = bandwidthRatio * 0.33;
    console.log('proposed rate', proposedPlaybackRate);

    // But limit the playback rate to a max value.
    proposedPlaybackRate = Math.max(
      Math.min(proposedPlaybackRate, MAX_SPEEDUP_RATE),
      1.0
    );

    // If this proposed speed is substantially faster than the current rate,
    // then allow us to ramp up by using a slower value for now.
    if (proposedPlaybackRate > this.playbackRate + MAX_SPEEDUP_RAMP) {
      proposedPlaybackRate = this.playbackRate + MAX_SPEEDUP_RAMP;
    }

    proposedPlaybackRate =
      Math.round(proposedPlaybackRate * Math.pow(10, 3)) / Math.pow(10, 3);

    try {
      const segment = getCurrentlyPlayingSegment(tech);
      if (!segment) {
        return;
      }

      // How far away from live edge do we start the compensator.
      const maxLatencyThreshold = Math.min(
        MAX_LATENCY,
        segment.duration * 1000 * HIGHEST_LATENCY_SEGMENT_LENGTH_MULTIPLIER
      );

      // How far away from live edge do we stop the compensator.
      const minLatencyThreshold = Math.max(
        MIN_LATENCY,
        segment.duration * 1000 * LOWEST_LATENCY_SEGMENT_LENGTH_MULTIPLIER
      );

      const segmentTime = segment.dateTimeObject.getTime();
      const now = new Date().getTime();
      const latency = now - segmentTime;

      console.log(
        'latency',
        latency / 1000,
        'min',
        minLatencyThreshold / 1000,
        'max',
        maxLatencyThreshold / 1000
      );

      if (latency > maxLatencyThreshold) {
        this.start(proposedPlaybackRate);
      } else if (latency < minLatencyThreshold) {
        this.stop();
      }
    } catch (err) {
      // console.warn(err);
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
