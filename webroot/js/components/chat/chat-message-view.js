import { h, Component } from '../../web_modules/preact.js';
import htm from '../../web_modules/htm.js';
import Mark from '../../web_modules/markjs/dist/mark.es6.min.js';
const html = htm.bind(h);

import {
  messageBubbleColorForString,
  textColorForString,
} from '../../utils/user-colors.js';
import { convertToText } from '../../utils/chat.js';
import { SOCKET_MESSAGE_TYPES } from '../../utils/websocket.js';

export default class ChatMessageView extends Component {
  constructor(props) {
    super(props);
    this.state = {
      formattedMessage: '',
    };
  }

  shouldComponentUpdate(nextProps, nextState) {
    const { formattedMessage } = this.state;
    const { formattedMessage: nextFormattedMessage } = nextState;

    return (formattedMessage !== nextFormattedMessage);
  }

  async componentDidMount() {
    const { message, username } = this.props;
    if (message && username) {
      const { body } = message;
      const formattedMessage = await formatMessageText(body, username);
      this.setState({
        formattedMessage,
      });
    }
  }


  render() {
    const { message } = this.props;
    const { author, timestamp, visible } = message;

    const { formattedMessage } = this.state;
    if (!formattedMessage) {
      return null;
    }
    const formattedTimestamp = formatTimestamp(timestamp);

    const isSystemMessage = message.type === SOCKET_MESSAGE_TYPES.SYSTEM;

    const authorTextColor = isSystemMessage
      ? { color: '#fff' }
      : { color: textColorForString(author) };
    const backgroundStyle = isSystemMessage
      ? { backgroundColor: '#667eea' }
      : { backgroundColor: messageBubbleColorForString(author) };
    const messageClassString = isSystemMessage
      ? getSystemMessageClassString()
      : getChatMessageClassString();

    return html`
      <div
        style=${backgroundStyle}
        class=${messageClassString}
        title=${formattedTimestamp}
      >
        <div class="message-content break-words w-full">
          <div
            style=${authorTextColor}
            class="message-author font-bold"
          >
            ${author}
          </div>
          <div
            class="message-text text-gray-300 font-normal overflow-y-hidden pt-2"
            dangerouslySetInnerHTML=${{ __html: formattedMessage }}
          ></div>
        </div>
      </div>
    `;
  }
}

function getSystemMessageClassString() {
  return 'message flex flex-row items-start p-4 m-2 rounded-lg shadow-l border-solid border-indigo-700 border-2 border-opacity-60 text-l';
}

function getChatMessageClassString() {
  return 'message flex flex-row items-start p-3 m-3 rounded-lg shadow-s text-sm';
}

export async function formatMessageText(message, username) {
  let formattedText = getMessageWithEmbeds(message);
  formattedText = convertToMarkup(formattedText);
  return await highlightUsername(formattedText, username);
}

function highlightUsername(message, username) {
  // https://github.com/julmot/mark.js/issues/115
  const node = document.createElement('span');
  node.innerHTML = message;
  return new Promise(res => {
    new Mark(node).mark(username, {
      element: 'span',
      className: 'highlighted px-1 rounded font-bold bg-orange-500',
      separateWordSearch: false,
      accuracy: {
        value: 'exactly',
        limiters: [",", ".", "'", '?', '@'],
      },
      done() {
        res(node.innerHTML);
      }
    });
  });
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
  if (
    embedText !== '' &&
    anchors.length == 1 &&
    isMessageJustAnchor(message, anchors[0])
  ) {
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


