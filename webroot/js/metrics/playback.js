const MAX_METRICS_COUNT = 20;
const METRICS_SEND_INTERVAL = 3000;

class PlaybackMetrics {
  constructor() {
    this.segmentDownloadTime = [];
    this.bandwidthTracking = [];
    this.latancyTracking = [];
    this.errors;

    setInterval(() => {
      this.send();
    }, METRICS_SEND_INTERVAL);
  }

  trackError(error) {
    console.log(error);
  }

  trackSegmentDownloadTime(seconds) {
    if (this.segmentDownloadTime.length > MAX_METRICS_COUNT) {
      this.segmentDownloadTime.shift();
    }

    this.segmentDownloadTime.push(seconds);
  }

  trackBandwidth(bps) {
    if (this.bandwidthTracking.length > MAX_METRICS_COUNT) {
      this.bandwidthTracking.shift();
    }

    this.bandwidthTracking.push(bps);
  }

  trackLatancy(latency) {
    if (this.latancyTracking.length > MAX_METRICS_COUNT) {
      this.latancyTracking.shift();
    }

    this.latancyTracking.push(latency);
  }

  send() {
    if (
      this.segmentDownloadTime.length < 4 ||
      this.bandwidthTracking.length < 4
    ) {
      return;
    }

    const average = (arr) => arr.reduce((p, c) => p + c, 0) / arr.length;

    const averageDownloadDuration = average(this.segmentDownloadTime) / 1000;
    const roundedAverageDownloadDuration =
      Math.round(averageDownloadDuration * 1000) / 1000;

    const averageBandwidth = average(this.bandwidthTracking) / 1000;
    const roundedAverageBandwidth = Math.round(averageBandwidth * 1000) / 1000;

    const averageLatancy = average(this.latancyTracking) / 1000;
    const roundedAverageLatancy = Math.round(averageLatancy * 1000) / 1000;

    const data = {
      bandwidth: roundedAverageBandwidth,
      latancy: roundedAverageLatancy,
      downloadDuration: roundedAverageDownloadDuration,
    };

    console.log(data);
  }
}

export default PlaybackMetrics;
