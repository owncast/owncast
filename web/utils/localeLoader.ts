import type { NextRouter } from 'next/router';
import i18n from '../i18n';

interface TranslationCatalog {
  [key: string]: string | TranslationCatalog;
}

export const withFallbackTranslations = (
  catalog: TranslationCatalog,
  fallback: TranslationCatalog,
): TranslationCatalog => {
  const merged = { ...catalog };
  Object.entries(fallback).forEach(([key, fallbackValue]) => {
    const value = merged[key];
    if (typeof fallbackValue === 'string') {
      if (typeof value !== 'string' || value.trim() === '') {
        merged[key] = fallbackValue;
      }
      return;
    }

    merged[key] = withFallbackTranslations(typeof value === 'object' ? value : {}, fallbackValue);
  });
  return merged;
};

/**
 * Locales with a catalog under web/i18n/ that we offer to browsers. This
 * mirrors the set that used to be bundled statically in i18n/index.js.
 * ('eu' has a catalog on disk but was never wired up, so it stays out.)
 */
export const AVAILABLE_LOCALES = [
  'ar',
  'bn',
  'de',
  'el',
  'en',
  'es',
  'fr',
  'ga',
  'hi',
  'hr',
  'it',
  'ja',
  'ko',
  'ms',
  'nl',
  'no',
  'pa',
  'pl',
  'pt',
  'ru',
  'sv',
  'th',
  'vi',
  'zh',
  'zh-CN',
];

/**
 * Chinese is the one language where matching the base subtag is not enough:
 * Simplified and Traditional are separate catalogs ('zh-CN' and 'zh') and
 * browsers report the script they want as a region or script subtag.
 *
 * The script subtag wins over the region: 'zh-Hant-CN' asks for Traditional.
 * Bare 'zh' stays Traditional, which is what the catalog Crowdin already
 * fills with.
 */
const chineseLocaleForTag = (tag: string): string => {
  if (tag.includes('-hans')) {
    return 'zh-CN';
  }
  if (tag.includes('-hant')) {
    return 'zh';
  }
  const region = tag.split('-')[1];
  if (region === 'cn' || region === 'sg' || region === 'my') {
    return 'zh-CN';
  }
  return 'zh';
};

/**
 * Maps a browser language tag onto a catalog we ship, or null when we have
 * none. Region and script subtags are dropped only after the language-specific
 * rules above had their say.
 */
const localeForBrowserLanguage = (browserLanguage: string | undefined): string | null => {
  const tag = (browserLanguage || '').toLowerCase();
  if (!tag) {
    return null;
  }
  // Chinese tags are ambiguous, so resolve them before the exact match below.
  if (tag.startsWith('zh-')) {
    return chineseLocaleForTag(tag);
  }
  if (AVAILABLE_LOCALES.includes(tag)) {
    return tag;
  }
  // Region-bearing catalogs keep their canonical casing, so a browser that
  // sends lowercase 'pt-br' still finds 'pt-BR'.
  const exact = AVAILABLE_LOCALES.find(locale => locale.toLowerCase() === tag);
  if (exact) {
    return exact;
  }
  const base = tag.split('-')[0];
  return AVAILABLE_LOCALES.includes(base) ? base : null;
};

/**
 * Decides which locale catalog to fetch. An explicit, known ?lang= value
 * wins, then the browser's primary language. English needs no fetch, and
 * unknown values fall back to the browser language the same way the
 * library's own detection did when every catalog was bundled.
 */
export const pickLocale = (
  queryLang: string | undefined,
  browserLanguage: string | undefined,
): string | null => {
  const fromQuery = queryLang && AVAILABLE_LOCALES.includes(queryLang) ? queryLang : null;
  const fromBrowser = localeForBrowserLanguage(browserLanguage);
  const locale = fromQuery || fromBrowser;
  return locale && locale !== 'en' ? locale : null;
};

let loadStarted = false;

let pendingUrlCleanup: string | null = null;

/**
 * LocaleSync (in _app) calls this on every selected-language change. It
 * returns true exactly once: when the automatic flip for `lang` has landed
 * and the temporary ?lang= parameter we added to trigger it should be
 * removed again, so copied URLs don't pin this viewer's language. Explicit
 * viewer-set ?lang= values never register a cleanup.
 */
export const shouldCleanUrlAfterFlip = (lang: string): boolean => {
  if (!pendingUrlCleanup || pendingUrlCleanup !== lang) {
    return false;
  }
  pendingUrlCleanup = null;
  return true;
};

/**
 * Fetches the viewer's locale catalog and activates it.
 *
 * next-export-i18n only re-evaluates the selected language when the ?lang
 * query value changes: its translations dependency is the same object we
 * mutate here (identity never changes), and its localStorage listener is
 * inert under languageDataStore 'query'. So after patching the catalog in,
 * we change the query parameter to re-run that effect, which now finds the
 * catalog present and switches every t() over. If ?lang was already set we
 * bounce it off and back on; the final effect run always sees the catalog,
 * so the switch cannot be lost to effect timing.
 */
export const loadViewerLocale = async (router: NextRouter): Promise<void> => {
  if (loadStarted || typeof window === 'undefined') {
    return;
  }

  const queryLang = typeof router.query.lang === 'string' ? router.query.lang : undefined;
  const browserLanguage = window.navigator.languages?.[0] || window.navigator.language;
  const locale = pickLocale(queryLang, browserLanguage);
  if (!locale || i18n.translations[locale]) {
    return;
  }
  // Marked only once real work begins: a no-op evaluation (English browser,
  // unknown locale, catalog already present) leaves the loader available for
  // a later call that does name a loadable locale. Everything before this
  // line is synchronous, so a double invocation (React strict mode) can't
  // start the load twice.
  loadStarted = true;

  try {
    // The locale is only known at runtime (browser language or ?lang=), so a
    // static import cannot work here. webpack turns this into one chunk per
    // catalog and loads just the one the viewer needs.
    const importedCatalog = await import(`../i18n/${locale}/translation.json`);
    const catalog = (importedCatalog.default ?? importedCatalog) as TranslationCatalog;
    i18n.translations[locale] = withFallbackTranslations(
      catalog,
      i18n.translations.en as TranslationCatalog,
    );

    const { lang, ...rest } = router.query;
    if (lang) {
      // Explicit ?lang= stays in the URL. Bounce it off and back on so the
      // library re-evaluates it now that the catalog exists.
      await router.replace({ query: rest }, undefined, { shallow: true, scroll: false });
    } else {
      // The parameter is only being added to trigger the switch, so have
      // LocaleSync remove it once the flip is confirmed to have landed.
      pendingUrlCleanup = locale;
    }
    await router.replace({ query: { ...rest, lang: locale } }, undefined, {
      shallow: true,
      scroll: false,
    });
  } catch (e) {
    // A failed catalog fetch (flaky network, content blocker) just leaves
    // the page in English rather than surfacing an error.
    pendingUrlCleanup = null;
    console.error(`unable to load the '${locale}' translation catalog`, e);
  }
};
