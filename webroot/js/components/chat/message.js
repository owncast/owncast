import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import { messageBubbleColorForString } from '../../utils/user-colors.js';
import { convertToText } from '../../utils/chat.js';
import { SOCKET_MESSAGE_TYPES } from '../../utils/websocket.js';

export default class Message extends Component {
  render(props) {
    const { message, username } = props;
    const { type } = message;
    if (type === SOCKET_MESSAGE_TYPES.CHAT) {
      const { author, body, timestamp } = message;
      const formattedMessage = formatMessageText(body, username);
      const formattedTimestamp = formatTimestamp(timestamp);

      const authorColor = messageBubbleColorForString(author);
      const authorTextColor = { color: authorColor };
      return (
        html`
          <div class="message flex flex-row items-start p-3">
            <div class="message-content text-sm break-words w-full">
              <div class="message-author text-white font-bold" style=${authorTextColor}>
                ${author}
              </div>
              <div
                class="message-text text-gray-300 font-normal overflow-y-hidden"
                title=${formattedTimestamp}
                dangerouslySetInnerHTML=${
                  { __html: formattedMessage }
                }
              ></div>
            </div>
        </div>
      `);
    } else if (type === SOCKET_MESSAGE_TYPES.NAME_CHANGE) {
      const { oldName, newName } = message;
      return (
        html`
          <div class="message message-name-change flex items-center justify-start p-3">
            <div class="message-content flex flex-row items-center justify-center text-sm w-full">
              <div class="text-white text-center opacity-50">
                <span class="font-bold">${oldName}</span> is now known as <span class="font-bold">${newName}</span>.
              </div>
            </div>
          </div>
        `
      );
    }
  }
}

export function formatMessageText(message, username) {
  let formattedText = highlightUsername(message, username);
  formattedText = getMessageWithEmbeds(formattedText);
  return convertToMarkup(formattedText);
}

function highlightUsername(message, username) {
  const pattern = new RegExp(
    '@?' + username.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&'),
    'gi'
  );
  return message.replace(
    pattern,
    '<span class="highlighted px-1 rounded font-bold bg-orange-500">$&</span>'
  );
}

function getMessageWithEmbeds(message) {
  var embedText = '';
  // Make a temporary element so we can actually parse the html and pull anchor tags from it.
  // This is a better approach than regex.
  var container = document.createElement('p');
  container.innerHTML = message;

  var anchors = container.getElementsByTagName('a');
  for (var i = 0; i < anchors.length; i++) {
    const url = anchors[i].href;
    if (getYoutubeIdFromURL(url)) {
      const youtubeID = getYoutubeIdFromURL(url);
      embedText += getYoutubeEmbedFromID(youtubeID);
    } else if (url.indexOf('instagram.com/p/') > -1) {
      embedText += getInstagramEmbedFromURL(url);
    } else if (isImage(url)) {
      embedText += getImageForURL(url);
    }
  }

  // If this message only consists of a single embeddable link
  // then only return the embed and strip the link url from the text.
  if (embedText !== '' && anchors.length == 1 && isMessageJustAnchor(message, anchors[0])) {
    return embedText;
  }
  return message + embedText;
}

function getYoutubeIdFromURL(url) {
  try {
    var regExp = /^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|\&v=)([^#\&\?]*).*/;
    var match = url.match(regExp);

    if (match && match[2].length == 11) {
      return match[2];
    } else {
      return null;
    }
  } catch (e) {
    console.log(e);
    return null;
  }
}

function getYoutubeEmbedFromID(id) {
  return `<div class="chat-embed youtube-embed"><lite-youtube videoid="${id}" /></div>`;
}

function getInstagramEmbedFromURL(url) {
  const urlObject = new URL(url.replace(/\/$/, ''));
  urlObject.pathname += '/embed';
  return `<iframe class="chat-embed instagram-embed" src="${urlObject.href}" frameborder="0" allowfullscreen></iframe>`;
}

function isImage(url) {
  const re = /\.(jpe?g|png|gif)$/i;
  return re.test(url);
}

function getImageForURL(url) {
  return `<a target="_blank" href="${url}"><img class="chat-embed embedded-image" src="${url}" /></a>`;
}

function isMessageJustAnchor(message, anchor) {
  return stripTags(message) === stripTags(anchor.innerHTML);
}

function formatTimestamp(sentAt) {
  sentAt = new Date(sentAt);
  if (isNaN(sentAt)) {
    return '';
  }

  let diffInDays = (new Date() - sentAt) / (24 * 3600 * 1000);
  if (diffInDays >= 1) {
    return (
      `Sent at ${sentAt.toLocaleDateString('en-US', {
        dateStyle: 'medium',
      })} at ` + sentAt.toLocaleTimeString()
    );
  }

  return `Sent at ${sentAt.toLocaleTimeString()}`;
}

/*
  You would call this when receiving a plain text
  value back from an API, and before inserting the
  text into the `contenteditable` area on a page.
*/
function convertToMarkup(str = '') {
  return convertToText(str).replace(/\n/g, '<br>');
}

function stripTags(str) {
  return str.replace(/<\/?[^>]+(>|$)/g, '');
}
