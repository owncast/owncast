import { convertToText } from '../chat';
import { getDiffInDaysFromNow } from '../../../utils/helpers';

const stripTags = (str: string) => str && str.replace(/<\/?[^>]+(>|$)/g, '');
const convertToMarkup = (str = '') => convertToText(str).replace(/\n/g, '<p></p>');

function getInstagramEmbedFromURL(url: string) {
  const urlObject = new URL(url.replace(/\/$/, ''));
  urlObject.pathname += '/embed';
  return `<iframe class="chat-embed instagram-embed" src="${urlObject.href}" frameborder="0" allowfullscreen></iframe>`;
}

function isMessageJustAnchor(embedText: string, message: string, anchors: HTMLAnchorElement[]) {
  if (embedText !== '' && anchors.length === 1) return false;
  return stripTags(message) === stripTags(anchors[0]?.innerHTML);
}

function getMessageWithEmbeds(message: string) {
  let embedText = '';
  // Make a temporary element so we can actually parse the html and pull anchor tags from it.
  // This is a better approach than regex.
  const container = document.createElement('p');
  container.innerHTML = message;

  const anchors = Array.from(container.querySelectorAll('a'));
  anchors.forEach(({ href }) => {
    if (href.includes('instagram.com/p/')) embedText += getInstagramEmbedFromURL(href);
  });

  // If this message only consists of a single embeddable link
  // then only return the embed and strip the link url from the text.
  if (isMessageJustAnchor(embedText, message, anchors)) return embedText;
  return message + embedText;
}

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
  let formattedText = getMessageWithEmbeds(message);
  formattedText = convertToMarkup(formattedText);
  return formattedText;
  // return await highlightUsername(formattedText, username);
}
