import { AutoplaySetting, autoplayModeForSetting, resolveAutoplaySetting } from '../utils/autoplay';

describe('autoplayModeForSetting', () => {
  test('off maps to no autoplay', () => {
    expect(autoplayModeForSetting(AutoplaySetting.Off)).toBe(false);
  });

  test('always maps to the muted-fallback mode', () => {
    expect(autoplayModeForSetting(AutoplaySetting.Always)).toBe('any');
  });

  test('sound-only maps to the sound-or-paused mode', () => {
    expect(autoplayModeForSetting(AutoplaySetting.SoundOnly)).toBe('play');
  });

  test('unrecognized values map to no autoplay', () => {
    expect(autoplayModeForSetting('bogus')).toBe(false);
    expect(autoplayModeForSetting('')).toBe(false);
  });

  test('reduced-motion forces autoplay off even when set to always', () => {
    expect(autoplayModeForSetting(AutoplaySetting.Always, { prefersReducedMotion: true })).toBe(
      false,
    );
  });

  test('data-saver forces autoplay off even when set to sound-only', () => {
    expect(autoplayModeForSetting(AutoplaySetting.SoundOnly, { saveData: true })).toBe(false);
  });

  test('an inaudible start suppresses sound-only autoplay', () => {
    // A persisted mute (volume 0) or an embed starting muted would let the
    // browser autoplay silently, violating "never starts silently".
    expect(autoplayModeForSetting(AutoplaySetting.SoundOnly, { startsInaudible: true })).toBe(
      false,
    );
  });

  test('an inaudible start does not affect always', () => {
    expect(autoplayModeForSetting(AutoplaySetting.Always, { startsInaudible: true })).toBe('any');
  });

  test('an audible start keeps sound-only autoplay', () => {
    expect(autoplayModeForSetting(AutoplaySetting.SoundOnly, { startsInaudible: false })).toBe(
      'play',
    );
  });

  test('preferences left unset do not change the mapping', () => {
    expect(
      autoplayModeForSetting(AutoplaySetting.Always, {
        prefersReducedMotion: false,
        saveData: false,
      }),
    ).toBe('any');
  });
});

describe('resolveAutoplaySetting', () => {
  test('a valid query value overrides the config value', () => {
    expect(resolveAutoplaySetting(AutoplaySetting.Off, AutoplaySetting.Always)).toBe(
      AutoplaySetting.Off,
    );
    expect(resolveAutoplaySetting(AutoplaySetting.Always, AutoplaySetting.Off)).toBe(
      AutoplaySetting.Always,
    );
    expect(resolveAutoplaySetting(AutoplaySetting.SoundOnly, AutoplaySetting.Off)).toBe(
      AutoplaySetting.SoundOnly,
    );
  });

  test('a missing query value falls back to the config value', () => {
    expect(resolveAutoplaySetting(undefined, AutoplaySetting.Always)).toBe(AutoplaySetting.Always);
  });

  test('an unrecognized query value falls back to the config value', () => {
    expect(resolveAutoplaySetting('bogus', AutoplaySetting.SoundOnly)).toBe(
      AutoplaySetting.SoundOnly,
    );
    expect(resolveAutoplaySetting('', AutoplaySetting.Off)).toBe(AutoplaySetting.Off);
  });
});
