/*
The Owncast Latency Compensator.

Tries to keep playback as close to the live edge as possible without ever
making the experience worse for the viewer. It slowly adjusts the playback
rate (imperceptibly) and, in extreme cases under pristine conditions, jumps
forward — always landing inside already-buffered content.

Safety model (buffer-first):
Near the live edge your playable buffer is physically capped by your latency,
so latency can only be reduced by spending buffer. The compensator will only
spend buffer down to a floor derived from the segment duration. All latency
targets are subordinate to that floor: when the floor is reached, compensation
stops, wherever that leaves us. It never chases a wall-clock number the buffer
can't afford.

How it behaves:
  - Measures latency from segment program-date-time plus the playhead's
    position within the segment, smoothed over several samples.
  - Only acts with bandwidth headroom, a stable player, and buffer above the
    floor. Backs off immediately when any of those degrade.
  - Tapers the speedup toward 1.0 as it approaches the target latency so it
    never overshoots into the thin-buffer zone at full speed.
  - Jumps only when very far behind, only forward, and only within buffered
    content, leaving at least the buffer floor after landing.
  - Remembers recent rebuffering events: they raise the latency floor and
    pause all compensation. Too many events disables it for the session.
*/

// ---- Tunables ----
const REBUFFER_EVENT_LIMIT = 4; // Max rebuffer events before giving up for good.
const MIN_BUFFER_DURATION = 200; // Min ms a buffering event must last to be counted.
const MAX_SPEEDUP_RATE = 1.08; // Max playback rate while compensating.
const MAX_SPEEDUP_RAMP = 0.005; // Max playback rate change applied per ramp tick.
const TIMEOUT_DURATION = 30 * 1000; // How long we back off after risky events.
const CHECK_TIMER_INTERVAL = 3 * 1000; // How often we evaluate conditions.
const SPEED_ADJUSTMENT_INTERVAL = 1 * 1000; // How often micro speed adjustments run.
const BUFFERING_AMNESTY_DURATION = 4 * 60 * 1000; // Rebuffer events expire after this.
const REQUIRED_BANDWIDTH_RATIO = 2.0; // download:bitrate ratio required to start.
const CONTINUE_BANDWIDTH_RATIO = 1.5; // ratio required to keep going once running.
const HIGHEST_LATENCY_SEGMENT_LENGTH_MULTIPLIER = 2.2; // segment length * this = start threshold.
const LOWEST_LATENCY_SEGMENT_LENGTH_MULTIPLIER = 1.7; // segment length * this = stop threshold.
const MIN_LATENCY = 4 * 1000; // Absolute lowest latency we'll compensate toward.
const MAX_LATENCY = 15 * 1000; // Cap on the computed start threshold.
const MAX_JUMP_LATENCY = 5 * 1000; // How far past the start threshold before jumps are considered.
const MAX_JUMP_FREQUENCY = 20 * 1000; // Min ms between jumps.
const MAX_ACTIONABLE_LATENCY = 90 * 1000; // Latency beyond this means bad data; do nothing.
const STARTUP_WAIT_TIME = 20 * 1000; // Leave the player alone this long after startup.
const CONSECUTIVE_STABLE_CHECKS = 3; // Healthy checks required before acting.
const CRITICAL_BUFFER_SECONDS = 2; // Below this playable buffer we hard back off.
const BUFFER_FLOOR_SEGMENT_MULTIPLIER = 1.5; // Buffer floor = segment duration * this...
const MIN_BUFFER_FLOOR_SECONDS = 4; // ...but never below this many seconds.
const MIN_BUFFER_SHRINK_LIMIT = 2; // Base for the shrink-trend limit (scaled by segment duration).
const BUFFER_BOOST_PER_SECOND = 0.02; // Speedup allowance per second of buffer above the floor.
const MIN_ACTIONABLE_RATE = 1.005; // Rates below this aren't worth running.
const SPEED_CHANGE_MIN_INTERVAL = 10 * 1000; // Min ms between starting speed changes.
const JUMP_MIN_WORTHWHILE_SECONDS = 2; // Never jump less than this.
const JUMP_IGNORE_BUFFER_DURATION = 10 * 1000; // Ignore buffer events this long after a seek.
const LATENCY_SAMPLE_COUNT = 5; // Latency samples kept for smoothing.

// ---- Earned trust (rev 9) ----
// After sustained good behavior the safety cushion thins a little: trust is
// earned by 3 continuous minutes of zero rebuffer events AND measured
// bandwidth headroom >= REQUIRED_BANDWIDTH_RATIO. It reverts instantly on
// any rebuffer, timeout, or headroom loss, and must be re-earned in full.
const TRUST_EARN_DURATION = 3 * 60 * 1000; // Clean time required to earn trust.
const MIN_LATENCY_TRUSTED = 3 * 1000; // Trusted absolute latency target floor.
const TRUSTED_FLOOR_JITTER_MARGIN = 1; // Trusted floor = segment duration + this...
const MIN_BUFFER_FLOOR_SECONDS_TRUSTED = 3; // ...but never below this many seconds.

// ---- Clock-skew fallback (rev 9) ----
const SKEW_ESTIMATION_MIN_ERROR = 10 * 1000; // Only correct clock errors bigger than this.
const SKEW_SIGHTING_WINDOW = 10; // New-segment sightings kept for estimation.

const REBUFFER_STORAGE_KEY = 'owncast-latency-compensator.rebuffers';

// The playable-buffer floor the compensator will never spend below.
// Untrusted: 1.5 segments, at least 4s. Trusted: one production cadence plus
// a jitter margin — never a bare multiplier that dips under the sawtooth.
export function bufferFloorSeconds(segmentDurationSec, trusted = false) {
  if (trusted) {
    return Math.max(
      segmentDurationSec + TRUSTED_FLOOR_JITTER_MARGIN,
      MIN_BUFFER_FLOOR_SECONDS_TRUSTED,
    );
  }
  return Math.max(BUFFER_FLOOR_SEGMENT_MULTIPLIER * segmentDurationSec, MIN_BUFFER_FLOOR_SECONDS);
}

// How much the playable buffer may shrink over three checks before we back
// off. Never tighter than the natural per-segment sawtooth.
function bufferShrinkLimit(segmentDurationSec) {
  return Math.max(MIN_BUFFER_SHRINK_LIMIT, segmentDurationSec);
}

// The latency band: compensation starts above maxLatencyMs and stops at or
// below minLatencyMs. Recent rebuffering raises the stop threshold — never
// lowers it.
export function latencyThresholds(segmentDurationSec, worstRebufferLatencyMs, trusted = false) {
  let minLatencyMs = Math.max(
    trusted ? MIN_LATENCY_TRUSTED : MIN_LATENCY,
    segmentDurationSec * 1000 * LOWEST_LATENCY_SEGMENT_LENGTH_MULTIPLIER,
  );

  if (Number.isFinite(worstRebufferLatencyMs)) {
    minLatencyMs = Math.max(minLatencyMs, Math.min(worstRebufferLatencyMs + 1000, MAX_LATENCY));
  }

  let maxLatencyMs = Math.max(
    minLatencyMs * 1.4,
    Math.min(segmentDurationSec * 1000 * HIGHEST_LATENCY_SEGMENT_LENGTH_MULTIPLIER, MAX_LATENCY),
  );

  if (minLatencyMs >= maxLatencyMs) {
    maxLatencyMs = minLatencyMs + 3000;
  }

  return { minLatencyMs, maxLatencyMs };
}

// The band adjusted for what the buffer can actually afford. Near the live
// edge, latency can only be reduced by spending buffer, and the buffer floor
// caps how much is spendable: the lowest reachable latency is
// latency − (buffer − floor). The effective minimum never sits below that,
// so the natural stopping point of a catch-up is always INSIDE the band and
// the resting state is silent instead of chattering on the start threshold.
export function effectiveLatencyBand(
  segmentDurationSec,
  worstRebufferLatencyMs,
  latencyMs,
  playableBufferSec,
  trusted = false,
) {
  let { minLatencyMs, maxLatencyMs } = latencyThresholds(
    segmentDurationSec,
    worstRebufferLatencyMs,
    trusted,
  );

  if (Number.isFinite(latencyMs) && Number.isFinite(playableBufferSec)) {
    const reachableMinMs =
      latencyMs - (playableBufferSec - bufferFloorSeconds(segmentDurationSec, trusted)) * 1000;
    if (reachableMinMs > minLatencyMs) {
      minLatencyMs = reachableMinMs;
      if (minLatencyMs >= maxLatencyMs) {
        maxLatencyMs = minLatencyMs + 3000;
      }
    }
  }

  return { minLatencyMs, maxLatencyMs };
}

/*
The pure decision function: a snapshot of player state in, exactly one action
out. No side effects, no timers, no player access — this is the behavioral
contract the tests pin down.

inputs = {
  latencyMs,              // smoothed latency estimate, or null if unknown
  playableBufferSec,      // mean of the last 4 buffer samples (smoothed)
  playableBufferRawSec,   // instantaneous buffered seconds ahead of the playhead
  bufferTrendSec,         // mean(last 4) − mean(the 4 before); smoothed trend
  bandwidthRatio,         // conservative download:bitrate ratio, or null
  segmentDurationSec,     // current segment duration
  stableChecks,           // consecutive healthy checks
  rebufferEvents,         // count of recent (unexpired) rebuffer events
  worstRebufferLatencyMs, // highest latency among recent rebuffers, or null
  msSinceLastSpeedChange, // ms since the target rate last changed
  msSinceLastJump,        // ms since the last jump/seek we performed
  running,                // whether compensation is currently active
  trusted,                // earned-trust state; thins the floor and thresholds
}

action = { type: "none" | "stop" | "timeout", reason }
       | { type: "speed", rate, reason }
       | { type: "jump", aheadSec, reason }
*/
export function decide(inputs) {
  const {
    latencyMs,
    playableBufferSec,
    playableBufferRawSec,
    bufferTrendSec,
    bandwidthRatio,
    segmentDurationSec,
    stableChecks,
    rebufferEvents,
    worstRebufferLatencyMs,
    msSinceLastSpeedChange,
    msSinceLastJump,
    running,
    trusted,
  } = inputs;

  // Unknown or absurd latency: we can't make safe decisions on bad data.
  if (!Number.isFinite(latencyMs) || Math.abs(latencyMs) > MAX_ACTIONABLE_LATENCY) {
    return { type: 'timeout', reason: 'latency unmeasurable or absurd' };
  }

  // About to run dry: hard back-off regardless of anything else. The
  // emergency check sees the raw instantaneous buffer — smoothing may
  // never delay it.
  const rawSec = Number.isFinite(playableBufferRawSec) ? playableBufferRawSec : playableBufferSec;
  if (Math.min(rawSec, playableBufferSec) < CRITICAL_BUFFER_SECONDS) {
    return { type: 'timeout', reason: 'critically low buffer' };
  }

  const floor = bufferFloorSeconds(segmentDurationSec, trusted);
  const idle = reason => ({ type: running ? 'stop' : 'none', reason });

  // Below the buffer floor we never spend buffer, no matter the latency.
  if (playableBufferSec < floor) {
    return idle('buffer at floor');
  }

  // Buffer shrinking fast: stop spending it before it becomes a stall.
  // The limit scales with segment duration: playable buffer naturally
  // sawtooths by a whole segment as segments land, so a fixed limit
  // false-trips on longer segments and causes engage/stop churn.
  if (bufferTrendSec < -bufferShrinkLimit(segmentDurationSec)) {
    return idle('buffer shrinking');
  }

  const { minLatencyMs, maxLatencyMs } = effectiveLatencyBand(
    segmentDurationSec,
    worstRebufferLatencyMs,
    latencyMs,
    playableBufferSec,
    trusted,
  );

  if (latencyMs <= minLatencyMs) {
    return idle('at target latency');
  }

  // Hysteresis: inside the band we only keep going if already running.
  if (latencyMs <= maxLatencyMs && !running) {
    return { type: 'none', reason: 'inside latency band' };
  }

  // Bandwidth headroom: strict to start, looser to continue.
  const requiredRatio = running ? CONTINUE_BANDWIDTH_RATIO : REQUIRED_BANDWIDTH_RATIO;
  if (!Number.isFinite(bandwidthRatio) || bandwidthRatio < requiredRatio) {
    return idle('insufficient bandwidth headroom');
  }

  // Jump: only when very far behind, under pristine conditions, and always
  // forward, landing inside buffered content with the floor left intact.
  if (
    latencyMs > maxLatencyMs + MAX_JUMP_LATENCY &&
    rebufferEvents === 0 &&
    stableChecks >= CONSECUTIVE_STABLE_CHECKS &&
    msSinceLastJump > MAX_JUMP_FREQUENCY &&
    playableBufferSec >= floor + 2 * segmentDurationSec
  ) {
    const aheadSec = Math.min(
      // Keep at least the buffer floor after landing.
      playableBufferSec - floor,
      // Don't jump past the target latency.
      (latencyMs - minLatencyMs) / 1000 - segmentDurationSec,
    );
    if (aheadSec >= Math.max(segmentDurationSec, JUMP_MIN_WORTHWHILE_SECONDS)) {
      return {
        type: 'jump',
        aheadSec: Math.floor(aheadSec * 100) / 100,
        reason: 'far behind live',
      };
    }
  }

  // Speed-up. Starting requires stability; continuing is re-evaluated on
  // every check so the rate tapers as conditions and distance change.
  if (!running) {
    if (stableChecks < CONSECUTIVE_STABLE_CHECKS) {
      return { type: 'none', reason: 'waiting for stability' };
    }
    if (rebufferEvents > 0) {
      return { type: 'none', reason: 'recent rebuffering' };
    }
    if (msSinceLastSpeedChange < SPEED_CHANGE_MIN_INTERVAL) {
      return { type: 'none', reason: 'rate changed too recently' };
    }
    // Entry hysteresis: starting needs segDur/2 of spendable buffer beyond
    // the floor, continuing only needs the floor — otherwise a buffer
    // hovering at the floor causes engage/stop flapping. (Kept as a plain
    // gate on purpose: folding it into the band geometry bars any catch-up
    // gaining less than the band's degenerate 3s offset.)
    if (playableBufferSec < floor + segmentDurationSec / 2) {
      return { type: 'none', reason: 'buffer too close to floor to start' };
    }
  }

  // Taper toward 1.0 as we approach the target so we never overshoot into
  // the thin-buffer zone at full speed.
  const distance = Math.min((latencyMs - minLatencyMs) / (maxLatencyMs - minLatencyMs), 1);
  const boostCap = Math.min(
    MAX_SPEEDUP_RATE - 1,
    BUFFER_BOOST_PER_SECOND * (playableBufferSec - floor),
  );
  const rate = Math.round((1 + Math.max(boostCap, 0) * distance) * 10000) / 10000;

  if (rate < MIN_ACTIONABLE_RATE) {
    return idle('not enough headroom to be worthwhile');
  }

  return { type: 'speed', rate, reason: 'catching up' };
}

function median(values) {
  const sorted = [...values].sort((a, b) => a - b);
  const mid = Math.floor(sorted.length / 2);
  return sorted.length % 2 ? sorted[mid] : (sorted[mid - 1] + sorted[mid]) / 2;
}

function mean(values) {
  return values.reduce((a, b) => a + b, 0) / values.length;
}

function getCurrentlyPlayingSegment(tech) {
  const playlist = tech.vhs.playlists.media();
  if (!playlist || !playlist.segments || playlist.segments.length === 0) {
    return null;
  }

  const snapshotTime = tech.currentTime();

  // Find the first segment whose end is past the playhead.
  for (let i = 0; i < playlist.segments.length; i += 1) {
    // Note: segment.end may be undefined when not yet buffered.
    if (snapshotTime < playlist.segments[i].end) {
      return playlist.segments[i];
    }
  }

  return playlist.segments[0];
}

// ---- Rebuffer memory persistence ----
// A page reload should not reset learned caution: events are stored in
// sessionStorage and restored (age-pruned) at construction. All storage
// access is failure-tolerant — missing or broken sessionStorage (Node,
// sandboxed iframes, private modes) silently degrades to in-memory only.

function loadPersistedRebufferEvents() {
  try {
    if (typeof window === 'undefined' || !window.sessionStorage) {
      return [];
    }
    const raw = window.sessionStorage.getItem(REBUFFER_STORAGE_KEY);
    if (!raw) {
      return [];
    }
    const events = JSON.parse(raw);
    if (!Array.isArray(events)) {
      return [];
    }
    const cutoff = Date.now() - BUFFERING_AMNESTY_DURATION;
    return events
      .filter(e => e && typeof e.atMs === 'number' && e.atMs > cutoff)
      .map(e => ({
        latencyMs: Number.isFinite(e.latencyMs) ? e.latencyMs : null,
        atMs: e.atMs,
      }));
  } catch {
    return [];
  }
}

function persistRebufferEvents(events) {
  try {
    if (typeof window === 'undefined' || !window.sessionStorage) {
      return;
    }
    window.sessionStorage.setItem(REBUFFER_STORAGE_KEY, JSON.stringify(events));
  } catch {
    // Storage full or forbidden; in-memory behavior is the fallback.
  }
}

class LatencyCompensator {
  constructor(player, options = {}) {
    this.player = player;
    this.playing = false;
    this.enabled = false;
    this.running = false;
    this.inTimeout = false;
    this.jumpingToLiveIgnoreBuffer = false;
    this.timeoutTimer = 0;
    this.timeoutEndingAt = 0;
    this.checkTimer = 0;
    this.bufferingTimer = 0;
    this.speedAdjustmentTimer = 0;
    this.playbackRate = 1.0;
    this.targetPlaybackRate = 1.0;
    this.lastJumpOccurred = 0;
    this.lastSpeedChange = 0;
    this.startupTime = Date.now();
    this.clockSkewMs = 0;
    this.skewWasSet = false;
    this.skewSightings = []; // server-minus-client estimates from new-segment sightings
    this.lastSkewSightingKey = null;
    this.currentLatency = null;
    this.consecutiveStableChecks = 0;
    this.latencySamples = [];
    this.bandwidthHistory = [];
    this.playableBufferHistory = [];

    // Earned trust: 3 clean minutes thin the safety cushion; any rebuffer,
    // timeout, or headroom loss reverts instantly.
    this.trustSince = null;
    this.trusted = false;

    // "lowest" (default) pins to the lowest rendition while compensating;
    // "current" only blocks upswitches.
    this.qualityPinMode = options.qualityPinMode === 'current' ? 'current' : 'lowest';

    // Last decision and the signals that produced it, for stats consumers.
    this.lastAction = null;
    this.lastPlayableBufferSec = null;
    this.lastBufferFloorSec = null;

    // Recent rebuffering events: { latencyMs: number|null, atMs: number }.
    // Pruned by age on access; the count and worst latency derive from it.
    // Restored from sessionStorage so a reload keeps learned caution.
    this.rebufferEvents = loadPersistedRebufferEvents();

    this.player.on('playing', this.handlePlaying.bind(this));
    this.player.on('pause', this.handlePause.bind(this));
    this.player.on('error', this.handleError.bind(this));
    this.player.on('waiting', this.handleBuffering.bind(this));
    this.player.on('stalled', this.handleBuffering.bind(this));
    this.player.on('ended', this.handleEnded.bind(this));
    this.player.on('canplaythrough', this.handlePlaying.bind(this));
    this.player.on('canplay', this.handlePlaying.bind(this));

    this.onStats = null;

    this.check = this.check.bind(this);
    this.smoothSpeedAdjustment = this.smoothSpeedAdjustment.bind(this);
  }

  // To keep our client clock in sync with the server clock to determine
  // accurate latency the clock skew should be set here to be used in
  // the calculation. Otherwise if somebody's client clock is significantly
  // off it will have a very incorrect latency determination and make bad
  // decisions. When this is never called, a coarse self-estimate kicks in
  // for clocks that are more than 10s off (see observeClockSkew).
  setClockSkew(skewMs) {
    this.clockSkewMs = skewMs;
    this.skewWasSet = true;
  }

  // Watch for newly published segments; each first sighting bounds the
  // server-client clock offset: est = (pdt + duration) − clientNow, which is
  // trueSkew minus the publish/poll delay, so est <= trueSkew. The MAX of
  // recent sightings is the least-delayed bound and under-reports latency —
  // the conservative direction.
  observeClockSkew(playlist) {
    for (let i = playlist.segments.length - 1; i >= 0; i -= 1) {
      const segment = playlist.segments[i];
      if (segment.dateTimeObject) {
        const key = segment.dateTimeObject.getTime();
        if (key !== this.lastSkewSightingKey) {
          this.lastSkewSightingKey = key;
          this.skewSightings.push(key + (segment.duration || 0) * 1000 - Date.now());
          if (this.skewSightings.length > SKEW_SIGHTING_WINDOW) {
            this.skewSightings.shift();
          }
        }
        return;
      }
    }
  }

  // The skew actually used in latency math. An explicit setClockSkew always
  // wins. The self-estimate only applies when the apparent clock error is
  // large (> 10s): small errors are harmless under the buffer-first design,
  // but a clock minutes off used to make the compensator permanently inert
  // via the bad-data guard.
  effectiveClockSkewMs() {
    if (this.skewWasSet) {
      return this.clockSkewMs;
    }
    if (this.skewSightings.length === 0) {
      return 0;
    }
    const estimate = Math.max(...this.skewSightings);
    return Math.abs(estimate) > SKEW_ESTIMATION_MIN_ERROR ? estimate : 0;
  }

  enable() {
    this.enabled = true;
    clearInterval(this.checkTimer);
    clearTimeout(this.bufferingTimer);

    this.checkTimer = setInterval(() => {
      this.check();
    }, CHECK_TIMER_INTERVAL);
  }

  // Disable means we're done for good and should no longer compensate.
  disable() {
    clearInterval(this.checkTimer);
    clearInterval(this.speedAdjustmentTimer);
    this.speedAdjustmentTimer = 0;
    clearTimeout(this.timeoutTimer);
    clearTimeout(this.bufferingTimer);
    this.inTimeout = false;
    this.stop();
    this.resetPlaybackRate();
    this.enabled = false;
  }

  // Runs on a timer: gather player state, ask decide() what to do, act on it.
  check() {
    if (!this.enabled) {
      return;
    }

    // Leave the player alone right after startup so it can settle and
    // build up a buffer before we start measuring or acting.
    if (Date.now() - this.startupTime < STARTUP_WAIT_TIME) {
      this.reportIdle('startup grace period');
      return;
    }

    if (this.inTimeout) {
      // Keep the timeout countdown flowing to stats consumers.
      this.reportStats(0, 0);
      return;
    }

    if (this.player.paused()) {
      this.consecutiveStableChecks = 0;
      this.reportIdle('paused');
      return;
    }

    if (this.player.seeking()) {
      this.consecutiveStableChecks = 0;
      this.reportIdle('seeking');
      return;
    }

    const tech = this.player.tech({ IWillNotUseThisInPlugins: true });

    // We need access to the internal tech of VHS to move forward.
    // Native playback (e.g. Safari CoreMedia) doesn't expose it.
    if (!tech || !tech.vhs) {
      this.reportIdle('no VHS tech (native playback?)');
      return;
    }

    // States 0 (NETWORK_EMPTY) and 3 (NETWORK_NO_SOURCE) mean nothing is
    // arriving. State 1 (IDLE) is normal for live HLS between segment
    // fetches and must not reset stability.
    const networkState = this.player.networkState();
    if (networkState === 0 || networkState === 3) {
      this.consecutiveStableChecks = 0;
      this.reportIdle('network not delivering');
      return;
    }

    try {
      const playlist = tech.vhs.playlists.media();
      if (!playlist || !playlist.segments || playlist.segments.length === 0) {
        this.reportIdle('no playlist data yet');
        return;
      }

      const playlistBandwidth = playlist.attributes && playlist.attributes.BANDWIDTH;
      const playerBandwidth = tech.vhs.systemBandwidth;
      if (playlistBandwidth > 0 && playerBandwidth > 0) {
        this.bandwidthHistory.push(playerBandwidth / playlistBandwidth);
        if (this.bandwidthHistory.length > 10) {
          this.bandwidthHistory.shift();
        }
      }

      if (!this.skewWasSet) {
        this.observeClockSkew(playlist);
      }

      const segment = getCurrentlyPlayingSegment(tech);
      if (!segment || !segment.dateTimeObject) {
        this.reportIdle('no segment timing (stream lacks program-date-time?)');
        return;
      }
      const segmentDurationSec = segment.duration || playlist.targetDuration || 4;

      const playableBufferRawSec = this.playableBufferSeconds();
      this.playableBufferHistory.push(playableBufferRawSec);
      if (this.playableBufferHistory.length > 10) {
        this.playableBufferHistory.shift();
      }
      // The raw playable buffer sawtooths by a whole segment as segments
      // land, and the 3s check cadence aliases against 2s/4s segment
      // arrivals into a period-2 alternation that a median passes straight
      // through. An even-length MEAN cancels it: 4 samples span a full beat
      // cycle for both 2s and 4s segments.
      const history = this.playableBufferHistory;
      const playableBufferSec = mean(history.slice(-4));
      const bufferTrendSec =
        history.length >= 8 ? playableBufferSec - mean(history.slice(-8, -4)) : 0;

      // Latency is measured against the segment's program-date-time plus
      // how far into the segment the playhead is — otherwise the estimate
      // saw-tooths by a whole segment duration. Samples are median-smoothed.
      const intraSegmentMs =
        typeof segment.start === 'number'
          ? Math.max((tech.currentTime() - segment.start) * 1000, 0)
          : 0;
      const sample =
        Date.now() +
        this.effectiveClockSkewMs() -
        (segment.dateTimeObject.getTime() + intraSegmentMs);
      this.latencySamples.push(sample);
      if (this.latencySamples.length > LATENCY_SAMPLE_COUNT) {
        this.latencySamples.shift();
      }
      const latencyMs = median(this.latencySamples);
      this.currentLatency = latencyMs;

      // Earned trust: the clock runs only while conditions are clean and
      // resets the moment they aren't. Reverting is instant, re-earning
      // takes the full duration again.
      const bandwidthRatio = this.getAverageBandwidth();
      if (
        this.bufferingCounter > 0 ||
        bandwidthRatio === null ||
        bandwidthRatio < REQUIRED_BANDWIDTH_RATIO
      ) {
        this.trustSince = null;
      } else if (this.trustSince === null) {
        this.trustSince = Date.now();
      }
      this.trusted =
        this.trustSince !== null && Date.now() - this.trustSince >= TRUST_EARN_DURATION;

      const floor = bufferFloorSeconds(segmentDurationSec, this.trusted);
      const healthy =
        playableBufferSec >= floor && bufferTrendSec >= -bufferShrinkLimit(segmentDurationSec);
      this.consecutiveStableChecks = healthy ? this.consecutiveStableChecks + 1 : 0;

      const worstRebufferLatencyMs = this.worstRebufferLatency();
      const action = decide({
        latencyMs,
        playableBufferSec,
        playableBufferRawSec,
        bufferTrendSec,
        bandwidthRatio,
        segmentDurationSec,
        stableChecks: this.consecutiveStableChecks,
        rebufferEvents: this.bufferingCounter,
        worstRebufferLatencyMs,
        msSinceLastSpeedChange: Date.now() - this.lastSpeedChange,
        msSinceLastJump: this.lastJumpOccurred ? Date.now() - this.lastJumpOccurred : Infinity,
        running: this.running,
        trusted: this.trusted,
      });

      this.lastAction = action;
      this.lastPlayableBufferSec = playableBufferSec;
      this.lastBufferFloorSec = floor;

      const { minLatencyMs, maxLatencyMs } = effectiveLatencyBand(
        segmentDurationSec,
        worstRebufferLatencyMs,
        latencyMs,
        playableBufferSec,
        this.trusted,
      );
      this.reportStats(minLatencyMs, maxLatencyMs);

      switch (action.type) {
        case 'timeout':
          this.timeout();
          break;
        case 'stop':
          this.stop();
          break;
        case 'jump':
          this.jump(this.player.currentTime() + action.aheadSec);
          break;
        case 'speed':
          this.startSpeed(action.rate);
          break;
        default:
          break;
      }

      console.info(
        'latency',
        latencyMs / 1000,
        'min',
        minLatencyMs / 1000,
        'max',
        maxLatencyMs / 1000,
        'playable',
        playableBufferSec,
        'action',
        action.type,
        `(${action.reason})`,
        'rate',
        this.playbackRate,
        'target rate',
        this.targetPlaybackRate,
        'rebuffer events',
        this.bufferingCounter,
        'stable checks',
        this.consecutiveStableChecks,
      );
    } catch (err) {
      console.trace(err);
    }
  }

  // Seconds of buffered content ahead of the playhead, using the buffered
  // range that actually contains the playhead (not just the first range).
  playableBufferSeconds() {
    const buffered = this.player.buffered();
    const currentTime = this.player.currentTime();
    for (let i = 0; i < buffered.length; i++) {
      if (currentTime >= buffered.start(i) - 0.5 && currentTime <= buffered.end(i)) {
        return buffered.end(i) - currentTime;
      }
    }
    return 0;
  }

  getAverageBandwidth() {
    if (this.bandwidthHistory.length < 5) {
      return null;
    }
    // Use the 30th percentile for a conservative estimate.
    const sorted = [...this.bandwidthHistory].sort((a, b) => a - b);
    return sorted[Math.floor(sorted.length * 0.3)];
  }

  getAveragePlayableBuffer() {
    if (this.playableBufferHistory.length === 0) {
      return 0;
    }
    return (
      this.playableBufferHistory.reduce((a, b) => a + b, 0) / this.playableBufferHistory.length
    );
  }

  // ---- Rebuffer memory ----

  pruneRebufferEvents() {
    const cutoff = Date.now() - BUFFERING_AMNESTY_DURATION;
    this.rebufferEvents = this.rebufferEvents.filter(e => e.atMs > cutoff);
  }

  get bufferingCounter() {
    this.pruneRebufferEvents();
    return this.rebufferEvents.length;
  }

  worstRebufferLatency() {
    this.pruneRebufferEvents();
    const latencies = this.rebufferEvents.map(e => e.latencyMs).filter(l => Number.isFinite(l));
    return latencies.length ? Math.max(...latencies) : null;
  }

  countBufferingEvent() {
    this.pruneRebufferEvents();
    this.rebufferEvents.push({
      latencyMs: Number.isFinite(this.currentLatency) ? this.currentLatency : null,
      atMs: Date.now(),
    });
    persistRebufferEvents(this.rebufferEvents);
    this.consecutiveStableChecks = 0;
    // Rebuffering instantly revokes earned trust.
    this.trustSince = null;
    this.trusted = false;

    if (this.rebufferEvents.length > REBUFFER_EVENT_LIMIT) {
      this.disable();
      return;
    }

    console.log(
      'latency compensation timeout due to buffering:',
      this.rebufferEvents.length,
      'buffering events of',
      REBUFFER_EVENT_LIMIT,
    );
  }

  // ---- Actions ----

  startSpeed(rate) {
    if (this.inTimeout || !this.enabled) {
      return;
    }

    if (!this.running) {
      this.running = true;
      this.enableOnlyLowQualityPlayback();
    }

    if (rate !== this.targetPlaybackRate) {
      this.targetPlaybackRate = rate;
      this.lastSpeedChange = Date.now();
    }

    this.ensureRampTimer();
  }

  stop() {
    if (this.running) {
      console.log('stopping latency compensator...');
    }
    this.running = false;
    this.enableAllQualityPlayback();

    // Ramp gently back to 1.0 — a stop is not an emergency.
    this.targetPlaybackRate = 1.0;
    if (this.playbackRate !== 1.0) {
      this.ensureRampTimer();
    }
  }

  timeout() {
    if (this.jumpingToLiveIgnoreBuffer) {
      return;
    }

    this.inTimeout = true;
    this.stop();
    // A timeout means risk was detected: return to normal speed immediately
    // and revoke earned trust.
    this.resetPlaybackRate();
    this.consecutiveStableChecks = 0;
    this.trustSince = null;
    this.trusted = false;

    clearTimeout(this.timeoutTimer);
    this.timeoutTimer = setTimeout(() => {
      this.endTimeout();
    }, TIMEOUT_DURATION);
    if (this.timeoutEndingAt === 0) {
      this.timeoutEndingAt = Date.now() + TIMEOUT_DURATION;
    }
    this.reportStats(0, 0);
  }

  endTimeout() {
    clearTimeout(this.timeoutTimer);
    this.inTimeout = false;
    this.timeoutEndingAt = 0;
  }

  jump(seekPosition) {
    // Safety invariant: we never seek backwards.
    if (seekPosition <= this.player.currentTime()) {
      return;
    }

    this.jumpingToLiveIgnoreBuffer = true;
    this.lastJumpOccurred = Date.now();
    this.consecutiveStableChecks = 0;
    this.latencySamples = [];
    this.resetPlaybackRate();

    console.info('jumping from', this.player.currentTime(), 'to', seekPosition);
    this.player.currentTime(seekPosition);

    setTimeout(() => {
      this.jumpingToLiveIgnoreBuffer = false;
    }, JUMP_IGNORE_BUFFER_DURATION);
  }

  // ---- Playback rate plumbing ----

  applyPlaybackRate(rate) {
    this.playbackRate = rate;
    this.player.playbackRate(rate);
  }

  resetPlaybackRate() {
    this.targetPlaybackRate = 1.0;
    this.applyPlaybackRate(1.0);
  }

  ensureRampTimer() {
    if (!this.speedAdjustmentTimer) {
      this.speedAdjustmentTimer = setInterval(
        this.smoothSpeedAdjustment,
        SPEED_ADJUSTMENT_INTERVAL,
      );
    }
  }

  // Moves the actual playback rate toward the target in imperceptible steps,
  // in both directions. Shuts its own timer down once settled at 1.0.
  smoothSpeedAdjustment() {
    if (!this.enabled || this.inTimeout) {
      this.resetPlaybackRate();
    }

    const diff = this.targetPlaybackRate - this.playbackRate;
    if (Math.abs(diff) <= MAX_SPEEDUP_RAMP) {
      if (this.playbackRate !== this.targetPlaybackRate) {
        this.applyPlaybackRate(this.targetPlaybackRate);
      }
      if (this.targetPlaybackRate === 1.0) {
        clearInterval(this.speedAdjustmentTimer);
        this.speedAdjustmentTimer = 0;
      }
      return;
    }

    this.applyPlaybackRate(this.playbackRate + Math.sign(diff) * MAX_SPEEDUP_RAMP);
  }

  // ---- Event handlers ----

  handlePlaying() {
    const wasPreviouslyPlaying = this.playing;
    this.playing = true;

    clearTimeout(this.bufferingTimer);
    if (!this.enabled) {
      return;
    }

    // Coming back from a stall: don't add seeking to the mix, just play.
    if (wasPreviouslyPlaying) {
      return;
    }

    // Cold start or resuming from a pause: catch up the accumulated gap by
    // seeking near — never onto — the live edge.
    this.seekNearLiveEdge();
  }

  // Seek forward to (live edge - buffer floor - one segment): the viewer
  // skips the gap but lands with breathing room instead of on the edge
  // itself, where the buffer is zero by definition.
  seekNearLiveEdge() {
    if (this.bufferingCounter > 0) {
      return;
    }
    if (Date.now() - this.lastJumpOccurred < MAX_JUMP_FREQUENCY) {
      return;
    }

    try {
      const seekable = this.player.seekable();
      if (!seekable || seekable.length === 0) {
        return;
      }
      const liveEdge = seekable.end(seekable.length - 1);

      const tech = this.player.tech({ IWillNotUseThisInPlugins: true });
      const playlist = tech && tech.vhs && tech.vhs.playlists.media();
      const segmentDurationSec = (playlist && playlist.targetDuration) || 4;

      const target = liveEdge - (bufferFloorSeconds(segmentDurationSec) + segmentDurationSec);

      // Only forward, and only when the gap is worth skipping.
      if (target - this.player.currentTime() < 2 * segmentDurationSec) {
        return;
      }

      this.jumpingToLiveIgnoreBuffer = true;
      this.lastJumpOccurred = Date.now();
      this.consecutiveStableChecks = 0;
      this.latencySamples = [];

      console.info('seeking near live edge:', target);
      this.player.currentTime(target);

      setTimeout(() => {
        this.jumpingToLiveIgnoreBuffer = false;
      }, JUMP_IGNORE_BUFFER_DURATION);
    } catch (err) {
      console.trace(err);
    }
  }

  handlePause() {
    this.playing = false;
    this.consecutiveStableChecks = 0;
  }

  handleEnded() {
    if (!this.enabled) {
      return;
    }

    this.disable();
  }

  handleError() {
    if (!this.enabled) {
      return;
    }

    this.timeout();
  }

  handleBuffering() {
    if (Date.now() - this.startupTime < STARTUP_WAIT_TIME) {
      return;
    }

    if (this.jumpingToLiveIgnoreBuffer) {
      return;
    }

    // A stall while the player is seeking is the fetch at a seek
    // discontinuity (e.g. the viewer scrubbed), not playback starvation —
    // genuine starvation fires with seeking() false. Not our fault: don't
    // count it against the compensator or enter a timeout for it.
    if (this.player.seeking()) {
      return;
    }

    // Stop spending buffer immediately.
    if (this.running) {
      this.resetPlaybackRate();
    }

    this.timeout();

    // Only count events that last long enough to matter.
    clearTimeout(this.bufferingTimer);
    this.bufferingTimer = setTimeout(() => {
      this.countBufferingEvent();
    }, MIN_BUFFER_DURATION);
  }

  // ---- Quality pinning ----
  // While actively compensating, renditions are restricted per
  // options.qualityPinMode: "lowest" (default) keeps only the cheapest
  // rendition so the freed headroom goes toward catching up; "current"
  // only blocks upswitches. All renditions are restored the moment
  // compensation stops.

  enableOnlyLowQualityPlayback() {
    this.setQualityPinned(true);
  }

  enableAllQualityPlayback() {
    this.setQualityPinned(false);
  }

  setQualityPinned(pinned) {
    try {
      const tech = this.player.tech({ IWillNotUseThisInPlugins: true });
      const representations =
        tech && tech.vhs && typeof tech.vhs.representations === 'function'
          ? tech.vhs.representations()
          : null;
      if (!representations || representations.length < 2) {
        return;
      }

      if (pinned && this.qualityPinMode === 'current') {
        const playlist = tech.vhs.playlists.media();
        const currentBandwidth = playlist && playlist.attributes && playlist.attributes.BANDWIDTH;
        if (!(currentBandwidth > 0)) {
          return; // can't tell what "current" is; leave renditions alone
        }
        representations.forEach(rep => {
          rep.enabled((rep.bandwidth || 0) <= currentBandwidth);
        });
        return;
      }

      const lowest = representations.reduce((a, b) =>
        (a.bandwidth || 0) <= (b.bandwidth || 0) ? a : b,
      );
      representations.forEach(rep => {
        rep.enabled(pinned ? rep === lowest : true);
      });
    } catch {
      // Renditions API unavailable; quality pinning is best-effort.
    }
  }

  // For inert states (startup grace, paused, no tech, ...) surface WHY the
  // compensator is doing nothing, so stats consumers never go stale exactly
  // when an explanation is most needed.
  reportIdle(reason) {
    this.lastAction = { type: 'idle', reason };
    this.reportStats(0, 0);
  }

  reportStats(minEndingLatency, maxStartingLatency) {
    if (!this.onStats) {
      return;
    }

    const stats = {
      latency: this.currentLatency,
      playbackRate: this.playbackRate,
      targetPlaybackRate: this.player.playbackRate(),
      enabled: this.enabled,
      running: this.running,
      bufferingEvents: this.bufferingCounter,
      inTimeout: this.inTimeout,
      timeoutRemaining: this.inTimeout ? this.timeoutEndingAt - Date.now() : 0,
      minEndingLatency,
      maxStartingLatency,
      averageBufferSeconds: this.getAveragePlayableBuffer(),
      averageBandwidthRatio: this.getAverageBandwidth(),
      consecutiveStableChecks: this.consecutiveStableChecks,
      // The last decision and the signals that produced it.
      action: this.lastAction,
      targetRate: this.targetPlaybackRate,
      playableBufferSeconds: this.lastPlayableBufferSec,
      bufferFloorSeconds: this.lastBufferFloorSec,
      trusted: this.trusted,
    };

    this.onStats(stats);
  }
}

export default LatencyCompensator;
