import fs from 'fs';
import path from 'path';
import type { NextRouter } from 'next/router';
import type { loadViewerLocale, shouldCleanUrlAfterFlip } from '../utils/localeLoader';
import { AVAILABLE_LOCALES, pickLocale } from '../utils/localeLoader';

type LoadViewerLocale = typeof loadViewerLocale;
type ShouldCleanUrlAfterFlip = typeof shouldCleanUrlAfterFlip;

interface I18nShape {
  translations: Record<string, Record<string, unknown>>;
}

// loadViewerLocale has a module-level once-guard, so every test that calls it
// needs a fresh copy of the module. The loader mutates the i18n singleton it
// imported, so both modules must come out of the same isolated registry or
// the test would inspect a different object than the one being mutated.
const freshModules = (): {
  load: LoadViewerLocale;
  shouldClean: ShouldCleanUrlAfterFlip;
  i18n: I18nShape;
} => {
  let load!: LoadViewerLocale;
  let shouldClean!: ShouldCleanUrlAfterFlip;
  let i18n!: I18nShape;
  jest.isolateModules(() => {
    /* eslint-disable global-require, @typescript-eslint/no-var-requires */
    load = require('../utils/localeLoader').loadViewerLocale;
    shouldClean = require('../utils/localeLoader').shouldCleanUrlAfterFlip;
    i18n = require('../i18n');
    /* eslint-enable global-require, @typescript-eslint/no-var-requires */
  });
  return { load, shouldClean, i18n };
};

const setBrowserLanguage = (language: string) => {
  Object.defineProperty(window.navigator, 'languages', {
    value: [language],
    configurable: true,
  });
};

afterEach(() => {
  // Remove the instance override so the next test (or file) sees jsdom's own
  // navigator.languages again.
  Reflect.deleteProperty(window.navigator, 'languages');
});

const makeRouter = (query: Record<string, string>) => {
  const replace = jest.fn(async () => true);
  const router = { query, replace } as unknown as NextRouter;
  return { router, replace };
};

describe('pickLocale', () => {
  test.each([
    ['valid query wins over browser language', 'de', 'fr-FR', 'de'],
    ['unknown query falls back to browser language', 'xx', 'de-DE', 'de'],
    ['absent query uses browser language', undefined, 'de-DE', 'de'],
    ['regional browser tag maps to its base locale', undefined, 'de-AT', 'de'],
    ['browser tag is case-insensitive', undefined, 'DE-at', 'de'],
    ['query is exact-match only, no case folding', 'DE', 'zz-ZZ', null],
    ['english from the query needs no catalog', 'en', 'de-DE', null],
    ['english from the browser needs no catalog', undefined, 'en-US', null],
    ['unknown query and unknown browser yield nothing', 'xx', 'zz-ZZ', null],
    ['both inputs missing yields nothing', undefined, undefined, null],
  ])('%s', (_name, queryLang, browserLanguage, expected) => {
    expect(pickLocale(queryLang, browserLanguage)).toBe(expected);
  });
});

describe('AVAILABLE_LOCALES', () => {
  test('offers en and de but deliberately not eu', () => {
    expect(AVAILABLE_LOCALES).toContain('en');
    expect(AVAILABLE_LOCALES).toContain('de');
    expect(AVAILABLE_LOCALES).not.toContain('eu');
  });

  test('has no duplicate entries', () => {
    expect(new Set(AVAILABLE_LOCALES).size).toBe(AVAILABLE_LOCALES.length);
  });

  test('every offered locale has a catalog file on disk', () => {
    const missing = AVAILABLE_LOCALES.filter(
      locale => !fs.existsSync(path.join(__dirname, '..', 'i18n', locale, 'translation.json')),
    );
    expect(missing).toEqual([]);
  });
});

describe('loadViewerLocale', () => {
  test('German browser without ?lang loads the catalog and sets lang once', async () => {
    setBrowserLanguage('de-DE');
    const { load, i18n } = freshModules();
    const { router, replace } = makeRouter({ existing: 'param' });

    await load(router);

    // The real catalog was fetched and patched into the shared i18n object.
    const de = i18n.translations.de as {
      Admin: { AccessTokens: { createNewAccessToken: string } };
    };
    expect(typeof de.Admin.AccessTokens.createNewAccessToken).toBe('string');

    // One replace that adds lang while keeping existing query params.
    expect(replace).toHaveBeenCalledTimes(1);
    expect(replace).toHaveBeenCalledWith({ query: { existing: 'param', lang: 'de' } }, undefined, {
      shallow: true,
      scroll: false,
    });
  });

  test('explicit ?lang=de bounces the query param off and back on', async () => {
    setBrowserLanguage('en-US');
    const { load, i18n } = freshModules();
    const { router, replace } = makeRouter({ lang: 'de', foo: 'bar' });

    await load(router);

    expect(i18n.translations.de).toBeDefined();
    expect(replace).toHaveBeenCalledTimes(2);
    expect(replace).toHaveBeenNthCalledWith(1, { query: { foo: 'bar' } }, undefined, {
      shallow: true,
      scroll: false,
    });
    expect(replace).toHaveBeenNthCalledWith(2, { query: { foo: 'bar', lang: 'de' } }, undefined, {
      shallow: true,
      scroll: false,
    });
  });

  test('English browser without ?lang changes nothing', async () => {
    setBrowserLanguage('en-US');
    const { load, i18n } = freshModules();
    const { router, replace } = makeRouter({});

    await load(router);

    expect(replace).not.toHaveBeenCalled();
    expect(Object.keys(i18n.translations)).toEqual(['en']);
  });

  test('second call is a no-op even with different inputs', async () => {
    setBrowserLanguage('de-DE');
    const { load, i18n } = freshModules();
    const first = makeRouter({});
    await load(first.router);
    expect(first.replace).toHaveBeenCalledTimes(1);

    // Different browser language and an explicit ?lang: the once-guard must
    // still win, so no French catalog and no navigation.
    setBrowserLanguage('fr-FR');
    const second = makeRouter({ lang: 'fr' });
    await load(second.router);

    expect(second.replace).not.toHaveBeenCalled();
    expect(i18n.translations.fr).toBeUndefined();
  });
});

describe('shouldCleanUrlAfterFlip', () => {
  test('browser-default load registers a consume-once cleanup for its locale', async () => {
    setBrowserLanguage('de-DE');
    const { load, shouldClean } = freshModules();
    const { router, replace } = makeRouter({});

    await load(router);

    // The loader itself never performs the cleanup: still exactly one
    // replace, the one that adds ?lang=de.
    expect(replace).toHaveBeenCalledTimes(1);

    // A non-matching locale neither matches nor consumes the pending cleanup.
    expect(shouldClean('fr')).toBe(false);

    // The matching locale consumes it exactly once.
    expect(shouldClean('de')).toBe(true);
    expect(shouldClean('de')).toBe(false);
    expect(shouldClean('fr')).toBe(false);
  });

  test('explicit ?lang load never registers a cleanup', async () => {
    setBrowserLanguage('en-US');
    const { load, shouldClean } = freshModules();
    const { router, replace } = makeRouter({ lang: 'de' });

    await load(router);

    expect(replace).toHaveBeenCalledTimes(2);
    expect(shouldClean('de')).toBe(false);
  });

  test('returns false before any load and after a no-op English load', async () => {
    setBrowserLanguage('en-US');
    const { load, shouldClean } = freshModules();

    expect(shouldClean('de')).toBe(false);

    const { router, replace } = makeRouter({});
    await load(router);

    expect(replace).not.toHaveBeenCalled();
    expect(shouldClean('de')).toBe(false);
    expect(shouldClean('en')).toBe(false);
  });
});
