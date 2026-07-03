import React from 'react';
import { act, render, screen } from '@testing-library/react';

// These tests pin the next-export-i18n semantics that utils/localeLoader
// depends on: the library only flips language when its query effect re-runs
// AND the catalog for router.query.lang is already present in the shared
// translations object. A library upgrade that changes any of this would
// silently leave non-English viewers on English, so it must fail here first.

// The library reads the shared catalog singleton (web/i18n) at module scope
// and the tests below mutate it, so both modules come fresh out of
// jest.isolateModules per test.
//
// react must NOT get a fresh copy inside that sandbox: react-dom (loaded out
// here by testing-library) wires the hook dispatcher into THIS react
// instance, and an isolated second copy would throw "Invalid hook call" the
// moment the library calls useState. Pinning react to the single actual
// instance keeps every registry, outer and isolated, on the same copy.
jest.mock('react', () => jest.requireActual('react'));

// The library's query effect lists router.query.lang in its deps, so
// useRouter must hand back a FRESH object per render carrying the holder's
// current query. A memoized router object would freeze the dep forever.
const mockRouterState: { query: Record<string, string> } = { query: {} };

jest.mock('next/router', () => ({
  useRouter: () => ({ query: { ...mockRouterState.query } }),
}));

/* eslint-disable global-require, @typescript-eslint/no-var-requires */
const enCatalog = require('../i18n/en/translation.json');
const deCatalog = require('../i18n/de/translation.json');
/* eslint-enable global-require, @typescript-eslint/no-var-requires */

// Assert against the real catalog files, not hardcoded strings, so a
// reworded translation cannot silently strand these tests.
const KEY = 'Admin.AccessTokens.selectAll';
const EN_VALUE: string = enCatalog.Admin.AccessTokens.selectAll;
const DE_VALUE: string = deCatalog.Admin.AccessTokens.selectAll;

interface I18nShape {
  translations: Record<string, Record<string, unknown>>;
}

// Fresh library + fresh catalog singleton + a probe component built against
// exactly that isolated copy of useTranslation.
const freshHarness = (): { i18n: I18nShape; Probe: () => React.ReactElement } => {
  let i18n!: I18nShape;
  let Probe!: () => React.ReactElement;
  jest.isolateModules(() => {
    /* eslint-disable global-require, @typescript-eslint/no-var-requires */
    i18n = require('../i18n');
    const { useTranslation } = require('next-export-i18n');
    /* eslint-enable global-require, @typescript-eslint/no-var-requires */
    Probe = function LocaleProbe() {
      const { t } = useTranslation();
      return <span data-testid="probe">{t(KEY)}</span>;
    };
  });
  return { i18n, Probe };
};

const probeText = () => screen.getByTestId('probe').textContent;

const setQueryLang = (lang?: string) => {
  mockRouterState.query = lang === undefined ? {} : { lang };
};

beforeEach(() => {
  setQueryLang(undefined);
});

describe('next-export-i18n semantics the lazy locale loader depends on', () => {
  beforeAll(() => {
    // If a catalog edit ever made these two values equal, every assertion
    // below would go vacuous. Fail loudly instead.
    expect(typeof EN_VALUE).toBe('string');
    expect(typeof DE_VALUE).toBe('string');
    expect(EN_VALUE).not.toBe(DE_VALUE);
  });

  test('stays English when ?lang=de names a catalog that is not loaded', () => {
    const { Probe } = freshHarness();
    setQueryLang('de');
    render(<Probe />);

    // The query effect ran on mount, saw translations.de missing, and its
    // guard refused the flip. This is what makes the loader's
    // fetch-catalog-BEFORE-touching-the-query ordering mandatory.
    expect(probeText()).toBe(EN_VALUE);
  });

  test('flips to German once the catalog is patched in and ?lang=de arrives', () => {
    const { Probe, i18n } = freshHarness();
    const { rerender } = render(<Probe />);
    expect(probeText()).toBe(EN_VALUE);

    // Exactly what loadViewerLocale does: mutate the shared translations
    // object, then activate through the query channel.
    i18n.translations.de = deCatalog;
    setQueryLang('de');
    act(() => {
      rerender(<Probe />);
    });

    // A plain rerender sufficed here because router.query.lang changed
    // (undefined -> 'de'), which re-runs the effect.
    expect(probeText()).toBe(DE_VALUE);
  });

  test('keeps German after ?lang is removed again, which makes URL cleanup safe', () => {
    const { Probe, i18n } = freshHarness();
    const { rerender } = render(<Probe />);

    i18n.translations.de = deCatalog;
    setQueryLang('de');
    act(() => {
      rerender(<Probe />);
    });
    expect(probeText()).toBe(DE_VALUE);

    // LocaleSync strips ?lang= after observing the flip. The library keeps
    // its selected language: the effect re-runs (dep 'de' -> undefined) but
    // its guard requires a truthy query value to change anything.
    setQueryLang(undefined);
    act(() => {
      rerender(<Probe />);
    });
    expect(probeText()).toBe(DE_VALUE);
  });

  test('a plain rerender cannot flip when ?lang=de predates the catalog, so the loader must bounce the param', () => {
    const { Probe, i18n } = freshHarness();
    setQueryLang('de');
    const { rerender } = render(<Probe />);
    expect(probeText()).toBe(EN_VALUE);

    // FINDING, pinned on purpose: patching the catalog alone changes no
    // effect dep (lang, router.query.lang, and the translations object
    // REFERENCE are all unchanged), so the effect never re-runs and the
    // language stays stuck on English.
    i18n.translations.de = deCatalog;
    act(() => {
      rerender(<Probe />);
    });
    expect(probeText()).toBe(EN_VALUE);

    // The loader recovers by bouncing the query param off and back on,
    // which is why loadViewerLocale does the replace-without-lang /
    // replace-with-lang dance when ?lang= was already set.
    setQueryLang(undefined);
    act(() => {
      rerender(<Probe />);
    });
    setQueryLang('de');
    act(() => {
      rerender(<Probe />);
    });
    expect(probeText()).toBe(DE_VALUE);
  });
});
