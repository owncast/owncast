const en = require('./en.json');
const es = require('./es.json');
const de = require('./de.json');
const fr = require('./fr.json');

const i18n = {
  translations: {
    en,
    es,
    de,
    fr,
  },
  defaultLang: 'en',
  useBrowserDefault: true,
  // optional property, will default to "query" if not set
  languageDataStore: 'query' || 'localStorage',
};

module.exports = i18n;
