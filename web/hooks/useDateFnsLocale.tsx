import type { Locale } from 'date-fns';
import { useSelectedLanguage } from 'next-export-i18n';
import { getDateFnsLocale } from '../utils/dateFnsLocale';

// Resolves the viewer's currently selected language (?lang=) to the
// matching date-fns Locale, for passing into formatDuration /
// formatDistanceToNow calls whose output feeds a Translation's vars.
export function useDateFnsLocale(): Locale | undefined {
  const { lang } = useSelectedLanguage();
  return getDateFnsLocale(lang);
}
