// Viewer-facing autoplay behaviors. Mirrors the backend config enum
// (off / always / sound-only) served in the client config.
const AUTOPLAY_VALUES = ['off', 'always', 'sound-only'];

// video.js autoplay option values:
//   false  -> no autoplay
//   'any'  -> try to play with sound, fall back to muted if the browser blocks it
//   'play' -> play with sound when the browser allows it, otherwise stay paused
type AutoplayMode = false | 'any' | 'play';

// Map an autoplay setting to the video.js autoplay value. A viewer who has
// signaled they would rather not autoplay (data-saver on, or reduced-motion
// preferred) gets no autoplay regardless of the configured setting.
export function autoplayModeForSetting(
  setting: string,
  preferences: {
    prefersReducedMotion?: boolean;
    saveData?: boolean;
    startsInaudible?: boolean;
  } = {},
): AutoplayMode {
  if (preferences.prefersReducedMotion || preferences.saveData) {
    return false;
  }
  if (setting === 'always') {
    return 'any';
  }
  if (setting === 'sound-only') {
    // "Only with sound" must never start silently. If the player would begin
    // inaudible (an embed asking to start muted, or a persisted volume of 0),
    // a play attempt can succeed silently anyway: Firefox permits autoplay of
    // muted or volume-0 media. So don't attempt autoplay at all.
    return preferences.startsInaudible ? false : 'play';
  }
  return false;
}

// Resolve the effective autoplay setting for an embed: a valid ?autoplay=
// override wins, otherwise fall back to the instance config value. Anything
// unrecognized falls back to the config value.
export function resolveAutoplaySetting(
  queryValue: string | undefined,
  configValue: string,
): string {
  if (queryValue && AUTOPLAY_VALUES.includes(queryValue)) {
    return queryValue;
  }
  return configValue;
}
