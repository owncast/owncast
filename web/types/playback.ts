export interface GeoDetails {
  countryCode: string;
  regionName: string;
  timeZone: string;
}

/**
 * Where a playback client's measurements came from. `client` means the
 * player reported them itself, `server` means the server observed them
 * while serving video to a player that reports nothing.
 */
export type PlaybackSource = 'client' | 'server';

/**
 * The latest playback health of one client. A null measurement is unknown
 * for that client rather than measured as zero, since not every player
 * reports every value.
 */
export interface PlaybackClientHealth {
  source: PlaybackSource;
  lastUpdate: string;
  playerState?: string;
  measurementStatus?: string;
  bandwidthKbps: number | null;
  latencySeconds: number | null;
  downloadSeconds: number | null;
  bitrateKbps: number | null;
  errorCount: number | null;
  qualityVariantChanges: number | null;
}

/**
 * A single client currently playing back video. `playback` is null when no
 * measurements exist for it at all.
 */
export interface PlaybackClient {
  clientID: string;
  viewerID: string;
  firstSeen: string;
  geo: GeoDetails | null;
  userAgent: string;
  ipAddress: string;
  playback: PlaybackClientHealth | null;
}
