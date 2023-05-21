import { getDiffInDaysFromNow } from '../../../utils/helpers';

export function formatTimestamp(sentAt: Date) {
  const now = new Date(sentAt);
  if (Number.isNaN(now)) return '';

  const diffInDays = getDiffInDaysFromNow(sentAt);

  if (diffInDays >= 1) {
    const localeDate = now.toLocaleDateString('en-US', {
      dateStyle: 'medium',
    });
    return `${localeDate} at ${now.toLocaleTimeString()}`;
  }

  return `${now.toLocaleTimeString()}`;
}
