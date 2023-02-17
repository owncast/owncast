import { convertToText } from '../chat';
import { getDiffInDaysFromNow } from '../../../utils/helpers';

const convertToMarkup = (str = '') => convertToText(str).replace(/\n/g, '<p></p>');

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

/*
  You would call this when receiving a plain text
  value back from an API, and before inserting the
  text into the `contenteditable` area on a page.
*/

export function formatMessageText(message: string) {
  const formattedText = convertToMarkup(message);
  return formattedText;
}
