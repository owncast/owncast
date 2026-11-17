import type { Locale } from 'date-fns';
import {
  ar,
  bn,
  de,
  el,
  enUS,
  es,
  eu,
  fr,
  hi,
  hr,
  it,
  ja,
  ko,
  ms,
  nb,
  nl,
  pl,
  pt,
  ru,
  sv,
  th,
  vi,
  zhTW,
} from 'date-fns/locale';

// date-fns ships no Irish (ga) or Punjabi (pa) locale, so duration text for
// those languages intentionally falls back to English rather than throwing
// or silently rendering undefined.
const DATE_FNS_LOCALE_MAP: Record<string, Locale> = {
  ar,
  bn,
  de,
  el,
  es,
  eu,
  fr,
  hi,
  hr,
  it,
  ja,
  ko,
  ms,
  nl,
  pl,
  pt,
  ru,
  sv,
  th,
  vi,
  ga: enUS,
  pa: enUS,
  no: nb,
  zh: zhTW,
};

export function getDateFnsLocale(locale: string): Locale | undefined {
  return DATE_FNS_LOCALE_MAP[locale];
}
