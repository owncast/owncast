// Clip is a saved window of a recorded stream, played back as its own HLS
// stream.
export interface Clip {
  id: string;
  streamId: string;
  title?: string;
  streamTitle?: string;
  clippedBy?: string;
  relativeStartTime: number;
  relativeEndTime: number;
  durationSeconds: number;
  // manifest is the path to the clip's HLS master playlist.
  manifest: string;
  // thumbnail is the clip's poster image, absent when none was generated.
  thumbnail?: string;
  timestamp: string;
}

// Replay is a recorded stream. Clips are taken from replays. Replays are
// managed in the admin and are not listed to viewers.
export interface Replay {
  id: string;
  title?: string;
  startTime: string;
  endTime?: string;
  inProgress?: boolean;
  manifest?: string;
  durationSeconds: number;
  totalBytes: number;
  clipCount: number;
}
