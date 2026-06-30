import { autoplayModeForSetting, resolveAutoplaySetting } from '../utils/autoplay';

describe('autoplayModeForSetting', () => {
  test('off maps to no autoplay', () => {
    expect(autoplayModeForSetting('off')).toBe(false);
  });

  test('always maps to the muted-fallback mode', () => {
    expect(autoplayModeForSetting('always')).toBe('any');
  });

  test('sound-only maps to the sound-or-paused mode', () => {
    expect(autoplayModeForSetting('sound-only')).toBe('play');
  });

  test('unrecognized values map to no autoplay', () => {
    expect(autoplayModeForSetting('bogus')).toBe(false);
    expect(autoplayModeForSetting('')).toBe(false);
  });

  test('reduced-motion forces autoplay off even when set to always', () => {
    expect(autoplayModeForSetting('always', { prefersReducedMotion: true })).toBe(false);
  });

  test('data-saver forces autoplay off even when set to sound-only', () => {
    expect(autoplayModeForSetting('sound-only', { saveData: true })).toBe(false);
  });

  test('preferences left unset do not change the mapping', () => {
    expect(autoplayModeForSetting('always', { prefersReducedMotion: false, saveData: false })).toBe(
      'any',
    );
  });
});

describe('resolveAutoplaySetting', () => {
  test('a valid query value overrides the config value', () => {
    expect(resolveAutoplaySetting('off', 'always')).toBe('off');
    expect(resolveAutoplaySetting('always', 'off')).toBe('always');
    expect(resolveAutoplaySetting('sound-only', 'off')).toBe('sound-only');
  });

  test('a missing query value falls back to the config value', () => {
    expect(resolveAutoplaySetting(undefined, 'always')).toBe('always');
  });

  test('an unrecognized query value falls back to the config value', () => {
    expect(resolveAutoplaySetting('bogus', 'sound-only')).toBe('sound-only');
    expect(resolveAutoplaySetting('', 'off')).toBe('off');
  });
});
