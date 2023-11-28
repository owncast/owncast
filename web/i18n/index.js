const en = require('./translations.en.json');
const de = require('./translations.de.json');
const fr = require('./translations.fr.json');
const sp = require('./translations.sp.json');

const i18n = {
  translations: {
    en,
    de,
    fr,
    sp,
  },
  defaultLang: 'en',
  useBrowserDefault: true,
  // optional property, will default to "query" if not set
  languageDataStore: 'query' || 'localStorage',
};

module.exports = i18n;
