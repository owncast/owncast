// Only English ships in the main javascript bundle. Every other locale is
// fetched on demand by utils/localeLoader, which webpack splits into one
// small chunk per language, patched into `translations` at runtime. This
// keeps ~500KB of catalogs out of the code that has to download and parse
// before the page becomes interactive.
const en = require('./en/translation.json');

const i18n = {
  translations: {
    en,
  },
  defaultLang: 'en',
  // Browser-language detection is handled by utils/localeLoader so the
  // catalog can be fetched first. The library's own detection only ever
  // sees English until that catalog arrives.
  useBrowserDefault: true,
  languageDataStore: 'query',
};

module.exports = i18n;
