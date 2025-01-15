const en = require('./en/translation.json');
// const es = require('./es/translation.json');
// const de = require('./de/translation.json');
// const fr = require('./fr/translation.json');

const i18n = {
  translations: {
    en,
    // es,
    // de,
    // fr,
  },
  defaultLang: 'en',
  useBrowserDefault: true,
  // optional property, will default to "query" if not set
  languageDataStore: 'query' || 'localStorage',
};

module.exports = i18n;
