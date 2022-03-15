import { URL_PLAYBACK_METRICS } from '../utils/constants.js';
const METRICS_SEND_INTERVAL = 10000;

class PlaybackMetrics {
  constructor() {
    this.hasPerformedInitialVariantChange = false;

    this.segmentDownloadTime = [];
    this.bandwidthTracking = [];
    this.latencyTracking = [];
    this.errors = 0;
    this.qualityVariantChanges = 0;
    this.isBuffering = false;

    setInterval(() => {
      this.send();
    }, METRICS_SEND_INTERVAL);
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

  trackSegmentDownloadTime(seconds) {
    this.segmentDownloadTime.push(seconds);
  }

  trackBandwidth(bps) {
    this.bandwidthTracking.push(bps);
  }

  trackLatency(latency) {
    this.latencyTracking.push(latency);
  }

  async send() {
    if (
      this.segmentDownloadTime.length < 4 ||
      this.bandwidthTracking.length < 4
    ) {
      return;
    }

    const errorCount = this.errors;
    const average = (arr) => arr.reduce((p, c) => p + c, 0) / arr.length;

    const averageDownloadDuration = average(this.segmentDownloadTime) / 1000;
    const roundedAverageDownloadDuration =
      Math.round(averageDownloadDuration * 1000) / 1000;

    const averageBandwidth = average(this.bandwidthTracking) / 1000;
    const roundedAverageBandwidth = Math.round(averageBandwidth * 1000) / 1000;

    const averageLatency = average(this.latencyTracking) / 1000;
    const roundedAverageLatency = Math.round(averageLatency * 1000) / 1000;

    const data = {
      bandwidth: roundedAverageBandwidth,
      latency: roundedAverageLatency,
      downloadDuration: roundedAverageDownloadDuration,
      errors: errorCount + this.isBuffering ? 1 : 0,
      qualityVariantChanges: this.qualityVariantChanges,
    };
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
      fetch(URL_PLAYBACK_METRICS, options);
    } catch (e) {
      console.error(e);
    }

    console.log(data);
  }
}

export default PlaybackMetrics;
