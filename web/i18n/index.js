const ar = require('./ar/translation.json');
const bn = require('./bn/translation.json');
const de = require('./de/translation.json');
const el = require('./el/translation.json');
const en = require('./en/translation.json');
const es = require('./es/translation.json');
const fr = require('./fr/translation.json');
const ga = require('./ga/translation.json');
const hi = require('./hi/translation.json');
const hr = require('./hr/translation.json');
const it = require('./it/translation.json');
const ja = require('./ja/translation.json');
const ko = require('./ko/translation.json');
const ms = require('./ms/translation.json');
const nl = require('./nl/translation.json');
const no = require('./no/translation.json');
const pa = require('./pa/translation.json');
const pl = require('./pl/translation.json');
const pt = require('./pt/translation.json');
const ru = require('./ru/translation.json');
const sv = require('./sv/translation.json');
const th = require('./th/translation.json');
const vi = require('./vi/translation.json');
const zh = require('./zh/translation.json');

const i18n = {
  translations: {
    ar,
    bn,
    de,
    el,
    en,
    es,
    fr,
    ga,
    hi,
    hr,
    it,
    ja,
    ko,
    ms,
    nl,
    no,
    pa,
    pl,
    pt,
    ru,
    sv,
    th,
    vi,
    zh,
  },
  defaultLang: 'en',
  useBrowserDefault: true,
  // optional property, will default to "query" if not set
  languageDataStore: 'query' || 'localStorage',
};

module.exports = i18n;
