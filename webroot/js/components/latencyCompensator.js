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

const BUFFER_LIMIT = 10; // Max number of buffering events before we stop compensating for latency.
const MIN_BUFFER_DURATION = 300; // Min duration a buffer event must last to be counted.
const SPEEDUP_RATE = 1.06; // The playback rate when compensating for latency.
const TIMEOUT_DURATION = 20_000; // The amount of time we stop handling latency after certain events.
const CHECK_TIMER_INTERVAL = 5_000; // How often we check if we should be compensating for latency.
const BUFFERING_AMNESTY_DURATION = 2 * 1000 * 60; // How often until a buffering event expires.
const HIGH_LATENCY_ENABLE_THRESHOLD = 20_000; // ms;
const LOW_LATENCY_DISABLE_THRESHOLD = 4500; // ms;
const REQUIRED_BANDWIDTH_RATIO = 2.0; // The player:bitrate ratio required to enable compensating for latency.

class LatencyCompensator {
  constructor(player) {
    this.player = player;
    this.enabled = true;
    this.running = false;
    this.inTimeout = false;
    this.timeoutTimer = 0;
    this.checkTimer = 0;
    this.maxLatencyThreshold = 10000;
    this.minLatencyThreshold = 8000;
    this.bufferingCounter = 0;
    this.bufferingTimer = 0;
    this.bufferStartedTimestamp = 0;
    this.player.on('playing', this.handlePlaying.bind(this));
    this.player.on('error', this.handleError.bind(this));
    this.player.on('waiting', this.handleBuffering.bind(this));
    this.player.on('ended', this.handleEnded.bind(this));
    this.player.on('canplaythrough', this.handlePlaying.bind(this));
    this.player.on('canplay', this.handlePlaying.bind(this));
  }

  // This is run on a timer to check if we should be compensating for latency.
  check() {
    if (this.inTimeout) {
      return;
    }

    const tech = this.player.tech({ IWillNotUseThisInPlugins: true });

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
        this.start();
      } else if (latency < this.minLatencyThreshold) {
        this.stop();
      }
    } catch (err) {
      console.warn(err);
    }
  }

  start() {
    if (this.running || this.disabled) {
      return;
    }

    this.running = true;
    this.player.playbackRate(SPEEDUP_RATE);
  }

  stop() {
    this.running = false;
    this.player.playbackRate(1);
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
    this.start();
  }

  handlePlaying() {
    if (!this.enabled) {
      return;
    }

    if (this.bufferStartedTimestamp !== 0) {
      const bufferingDuration =
        new Date().getTime() - this.bufferStartedTimestamp;
      this.bufferStartedTimestamp = 0;

      // If the buffering event lasted long enough then we will stay in
      // a timeout and count it as a real buffering event. Otherwise
      // we will ignore it.
      if (bufferingDuration > MIN_BUFFER_DURATION) {
        this.countBufferingEvent();
      } else {
        this.endTimeout();
      }
    }
    clearInterval(this.checkTimer);
    clearTimeout(this.bufferingTimer);

    this.checkTimer = setInterval(() => {
      this.check();
    }, CHECK_TIMER_INTERVAL);
  }

  handleEnded() {
    this.disable();
  }

  handleError() {
    this.timeout();
  }

  countBufferingEvent() {
    this.bufferingCounter++;
    if (this.bufferingCounter > BUFFER_LIMIT) {
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
    this.bufferStartedTimestamp = new Date().getTime();
    this.timeout();
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
