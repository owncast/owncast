const en = require('./en-US.json');
const es = require('./es-ES.json');
const de = require('./de-DE.json');
const fr = require('./fr-FR.json');

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
