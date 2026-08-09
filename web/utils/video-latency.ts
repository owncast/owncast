// The live player exposes its media coordinates only when a clip action is
// pressed. Keeping a supplier avoids rendering on every playback tick.
export interface VideoPosition {
  latencySeconds: number | null;
  playheadSeconds: number | null;
}

let supplier: (() => VideoPosition | null) | null = null;

export function setVideoPositionSupplier(fn: (() => VideoPosition | null) | null) {
  supplier = fn;
}

export function getVideoPosition(): VideoPosition | null {
  if (!supplier) {
    return null;
  }
  try {
    const position = supplier();
    if (!position) {
      return null;
    }
    const { latencySeconds, playheadSeconds } = position;
    return {
      latencySeconds:
        typeof latencySeconds === 'number' && Number.isFinite(latencySeconds) && latencySeconds >= 0
          ? latencySeconds
          : null,
      playheadSeconds:
        typeof playheadSeconds === 'number' &&
        Number.isFinite(playheadSeconds) &&
        playheadSeconds >= 0
          ? playheadSeconds
          : null,
    };
  } catch {
    return null;
  }
}
