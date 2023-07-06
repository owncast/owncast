/* eslint-disable no-plusplus */
const URL_PLAYBACK_METRICS = `/api/metrics/playback`;
const METRICS_SEND_INTERVAL = 10000;
const MAX_VALID_LATENCY_SECONDS = 100; // Anything > this gets thrown out.

function getCurrentlyPlayingSegment(tech) {
  const targetMedia = tech.vhs.playlists.media();
  const snapshotTime = tech.currentTime();
  let segment;

  // Iterate trough available segments and get first within which snapshot_time is
  // eslint-disable-next-line no-plusplus
  for (let i = 0, l = targetMedia.segments.length; i < l; i++) {
    // Note: segment.end may be undefined or is not properly set
    if (snapshotTime < targetMedia.segments[i].end) {
      segment = targetMedia.segments[i];
      break;
    }
  }

  if (!segment) {
    [segment] = targetMedia.segments;
  }

  return segment;
}

class PlaybackMetrics {
  constructor(player, videojs) {
    this.player = player;
    this.supportsDetailedMetrics = false;
    this.hasPerformedInitialVariantChange = false;
    this.clockSkewMs = 0;

    this.segmentDownloadTime = [];
    this.bandwidthTracking = [];
    this.latencyTracking = [];
    this.errors = 0;
    this.qualityVariantChanges = 0;
    this.isBuffering = false;
    this.bufferingDurationTimer = 0;
    this.collectPlaybackMetricsTimer = 0;

    this.videoJSReady = this.videoJSReady.bind(this);
    this.handlePlaying = this.handlePlaying.bind(this);
    this.handleBuffering = this.handleBuffering.bind(this);
    this.handleEnded = this.handleEnded.bind(this);
    this.handleError = this.handleError.bind(this);
    this.send = this.send.bind(this);
    this.collectPlaybackMetrics = this.collectPlaybackMetrics.bind(this);
    this.handleNoLongerBuffering = this.handleNoLongerBuffering.bind(this);
    this.sendMetricsTimer = 0;

    this.player.on('canplaythrough', this.handleNoLongerBuffering);
    this.player.on('error', this.handleError);
    this.player.on('stalled', this.handleBuffering);
    this.player.on('waiting', this.handleBuffering);
    this.player.on('playing', this.handlePlaying);
    this.player.on('ended', this.handleEnded);

    // Keep a reference of the standard vjs xhr function.
    const oldVjsXhrCallback = videojs.xhr;

    // Override the xhr function to track segment download time.
    // eslint-disable-next-line no-param-reassign
    videojs.Vhs.xhr = (...args) => {
      if (args[0].uri.match('.ts')) {
        const start = new Date();

        const cb = args[1];
        // eslint-disable-next-line no-param-reassign
        args[1] = (request, error, response) => {
          const end = new Date();
          const delta = end.getTime() - start.getTime();
          this.trackSegmentDownloadTime(delta);
          cb(request, error, response);
        };
      }

      return oldVjsXhrCallback(...args);
    };

    this.videoJSReady();

    this.sendMetricsTimer = setInterval(() => {
      this.send();
    }, METRICS_SEND_INTERVAL);
  }

  stop() {
    clearInterval(this.sendMetricsTimer);
    this.player.off();
  }

  // Keep our client clock in sync with the server clock to determine
  // accurate latency calculations.
  setClockSkew(skewMs) {
    this.clockSkewMs = skewMs;
  }

  videoJSReady() {
    const tech = this.player.tech({ IWillNotUseThisInPlugins: true });
    this.supportsDetailedMetrics = !!tech;

    tech?.on('usage', e => {
      if (e.name === 'vhs-unknown-waiting') {
        this.setIsBuffering(true);
      }

      if (e.name === 'vhs-rendition-change-abr') {
        // Quality variant has changed
        this.incrementQualityVariantChanges();
      }
    });

    // Variant changed
    const trackElements = this.player.textTracks();
    trackElements.addEventListener('cuechange', () => {
      this.incrementQualityVariantChanges();
    });
  }

  handlePlaying() {
    clearInterval(this.collectPlaybackMetricsTimer);
    this.collectPlaybackMetricsTimer = setInterval(() => {
      this.collectPlaybackMetrics();
    }, 5000);
  }

  handleEnded() {
    clearInterval(this.collectPlaybackMetricsTimer);
  }

  handleBuffering() {
    this.incrementErrorCount(1);
    this.setIsBuffering(true);
  }

  handleNoLongerBuffering() {
    this.setIsBuffering(false);
  }

  handleError() {
    this.incrementErrorCount(1);
  }

  incrementErrorCount(count) {
    this.errors += count;
  }

  incrementQualityVariantChanges() {
    // We always start the player at the lowest quality, so let's just not
    // count the first change.
    if (!this.hasPerformedInitialVariantChange) {
      this.hasPerformedInitialVariantChange = true;
      return;
    }
    this.qualityVariantChanges++;
  }

  setIsBuffering(isBuffering) {
    this.isBuffering = isBuffering;

    if (!isBuffering) {
      clearTimeout(this.bufferingDurationTimer);
      return;
    }

    this.bufferingDurationTimer = setTimeout(() => {
      this.incrementErrorCount(1);
    }, 500);
  }

  trackSegmentDownloadTime(seconds) {
    this.segmentDownloadTime.push(seconds);
  }

  trackBandwidth(bps) {
    this.bandwidthTracking.push(bps);
  }

  trackLatency(latency) {
    this.latencyTracking.push(latency);
  }

  collectPlaybackMetrics() {
    const tech = this.player.tech({ IWillNotUseThisInPlugins: true });
    if (!tech || !tech.vhs) {
      return;
    }

    // If we're paused then do nothing.
    if (this.player.paused()) {
      return;
    }

    // Network state 2 means we're actively using the network.
    // We only care about reporting metrics with network activity stats
    // if it's actively being used, so don't report otherwise.
    const networkState = this.player.networkState();
    if (networkState !== 2) {
      return;
    }

    const bandwidth = tech.vhs.systemBandwidth;
    this.trackBandwidth(bandwidth);

    try {
      const segment = getCurrentlyPlayingSegment(tech);
      if (!segment || !segment.dateTimeObject) {
        return;
      }

      const segmentTime = segment.dateTimeObject.getTime();
      const now = new Date().getTime() + this.clockSkewMs;
      const latency = now - segmentTime;

      // Throw away values that seem invalid.
      if (latency < 0 || latency / 1000 >= MAX_VALID_LATENCY_SECONDS) {
        return;
      }

      this.trackLatency(latency);
    } catch (err) {
      console.warn(err);
    }
  }

  async send() {
    if (this.segmentDownloadTime.length === 0) {
      return;
    }

    // If we're paused then do nothing.
    if (!this.player || this.player.paused()) {
      return;
    }

    const errorCount = this.errors;

    let data;
    if (this.supportsDetailedMetrics) {
      const average = arr => arr.reduce((p, c) => p + c, 0) / arr.length;

      const averageDownloadDuration = average(this.segmentDownloadTime) / 1000;
      const roundedAverageDownloadDuration = Math.round(averageDownloadDuration * 1000) / 1000;

      const averageBandwidth = average(this.bandwidthTracking) / 1000;
      const roundedAverageBandwidth = Math.round(averageBandwidth * 1000) / 1000;

      const averageLatency = average(this.latencyTracking) / 1000;
      const roundedAverageLatency = Math.round(averageLatency * 1000) / 1000;

      data = {
        bandwidth: roundedAverageBandwidth,
        latency: roundedAverageLatency,
        downloadDuration: roundedAverageDownloadDuration,
        errors: errorCount + (this.isBuffering ? 1 : 0),
        qualityVariantChanges: this.qualityVariantChanges,
      };
    } else {
      data = {
        errors: errorCount + (this.isBuffering ? 1 : 0),
      };
    }

    this.errors = 0;
    this.qualityVariantChanges = 0;
    this.segmentDownloadTime = [];
    this.bandwidthTracking = [];
    this.latencyTracking = [];

    const options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    };

    try {
      await fetch(URL_PLAYBACK_METRICS, options);
    } catch (e) {
      console.error(e);
    }
  }
}

export default PlaybackMetrics;
