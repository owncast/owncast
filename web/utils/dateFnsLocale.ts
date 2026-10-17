import type { Locale } from 'date-fns';
import {
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
  no: nb,
  zh: zhTW,
};

export function getDateFnsLocale(locale: string): Locale | undefined {
  return DATE_FNS_LOCALE_MAP[locale];
}
