import { isNewerVersion, isPluginUpdateAvailable } from '../utils/apis';

describe('isNewerVersion', () => {
  test('candidate strictly newer is true', () => {
    expect(isNewerVersion('1.2.4', '1.2.3')).toBe(true);
    expect(isNewerVersion('1.3.0', '1.2.9')).toBe(true);
    expect(isNewerVersion('2.0.0', '1.9.9')).toBe(true);
  });

  test('equal or older candidate is false', () => {
    expect(isNewerVersion('1.2.3', '1.2.3')).toBe(false);
    expect(isNewerVersion('1.2.2', '1.2.3')).toBe(false);
    expect(isNewerVersion('0.2.1', '0.3.0')).toBe(false); // the dev-build/downgrade case
  });

  test('missing versions are false (no crash)', () => {
    expect(isNewerVersion(undefined, '1.2.3')).toBe(false);
    expect(isNewerVersion('1.2.3', undefined)).toBe(false);
    expect(isNewerVersion('', '')).toBe(false);
  });

  test('non-semver strings are false rather than throwing', () => {
    expect(isNewerVersion('latest', '1.2.3')).toBe(false);
    expect(isNewerVersion('1.2.3', 'dev')).toBe(false);
  });
});

describe('isPluginUpdateAvailable', () => {
  describe('an update is offered only when the registry version is strictly newer', () => {
    test('newer patch in the registry', () => {
      expect(isPluginUpdateAvailable('0.2.1', '0.2.2')).toBe(true);
    });

    test('newer minor in the registry', () => {
      expect(isPluginUpdateAvailable('0.2.1', '0.3.0')).toBe(true);
    });

    test('newer major in the registry', () => {
      expect(isPluginUpdateAvailable('0.2.1', '1.0.0')).toBe(true);
    });

    test('same version is not an update', () => {
      expect(isPluginUpdateAvailable('0.2.1', '0.2.1')).toBe(false);
    });
  });

  // The regression this guards: the UI used to flag any version *difference*
  // as an update, so a newer local/dev build looked like it needed a
  // "downgrade" to the registry version.
  describe('a newer locally-installed (dev) build is NOT flagged as an update', () => {
    test('installed 0.3.0 vs registry 0.2.1', () => {
      expect(isPluginUpdateAvailable('0.3.0', '0.2.1')).toBe(false);
    });

    test('installed major ahead of registry', () => {
      expect(isPluginUpdateAvailable('2.0.0', '1.9.9')).toBe(false);
    });
  });

  describe('missing versions never offer an update', () => {
    test('no installed version', () => {
      expect(isPluginUpdateAvailable(undefined, '0.2.1')).toBe(false);
    });

    test('no registry version', () => {
      expect(isPluginUpdateAvailable('0.2.1', undefined)).toBe(false);
    });

    test('neither version', () => {
      expect(isPluginUpdateAvailable(undefined, undefined)).toBe(false);
    });

    test('empty strings', () => {
      expect(isPluginUpdateAvailable('', '')).toBe(false);
    });
  });

  describe('un-orderable (non-semver) versions are treated conservatively as no update', () => {
    test('non-semver installed version', () => {
      expect(isPluginUpdateAvailable('dev', '0.2.1')).toBe(false);
    });

    test('non-semver registry version', () => {
      expect(isPluginUpdateAvailable('0.2.1', 'latest')).toBe(false);
    });

    test('both non-semver', () => {
      expect(isPluginUpdateAvailable('foo', 'bar')).toBe(false);
    });
  });

  describe('prerelease ordering follows semver', () => {
    test('a prerelease is older than its release', () => {
      expect(isPluginUpdateAvailable('1.0.0-rc.1', '1.0.0')).toBe(true);
    });

    test('a release is not "updated" by an older prerelease', () => {
      expect(isPluginUpdateAvailable('1.0.0', '1.0.0-rc.1')).toBe(false);
    });
  });
});
