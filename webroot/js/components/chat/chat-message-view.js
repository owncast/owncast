import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
import Mark from '/js/web_modules/markjs/dist/mark.es6.min.js';
const html = htm.bind(h);

import {
  messageBubbleColorForHue,
  textColorForHue,
} from '../../utils/user-colors.js';
import { convertToText, checkIsModerator } from '../../utils/chat.js';
import { SOCKET_MESSAGE_TYPES } from '../../utils/websocket.js';
import { getDiffInDaysFromNow } from '../../utils/helpers.js';
import ModeratorActions from './moderator-actions.js';

export default class ChatMessageView extends Component {
  constructor(props) {
    super(props);
    this.state = {
      formattedMessage: '',
      moderatorMenuOpen: false,
    };
  }

  shouldComponentUpdate(nextProps, nextState) {
    const { formattedMessage } = this.state;
    const { formattedMessage: nextFormattedMessage } = nextState;

    return (
      formattedMessage !== nextFormattedMessage ||
      (!this.props.isModerator && nextProps.isModerator)
    );
  }

  async componentDidMount() {
    const { message, username } = this.props;
    const { body } = message;

    if (message && username) {
      const formattedMessage = await formatMessageText(body, username);
      this.setState({
        formattedMessage,
      });
    }
  }
  render() {
    const { message, isModerator, accessToken } = this.props;
    const { user, timestamp } = message;

    if (!user) {
      return null;
    }

    const { displayName, displayColor, createdAt, isBot, authenticated } = user;
    const isAuthorModerator = checkIsModerator(message);

    const isMessageModeratable =
      isModerator && message.type === SOCKET_MESSAGE_TYPES.CHAT;

    const { formattedMessage } = this.state;
    if (!formattedMessage) {
      return null;
    }
    const formattedTimestamp = `Sent at ${formatTimestamp(timestamp)}`;
    const userMetadata = createdAt
      ? `${displayName} first joined ${formatTimestamp(createdAt)}`
      : null;

    const isSystemMessage = message.type === SOCKET_MESSAGE_TYPES.SYSTEM;

    const authorTextColor = isSystemMessage
      ? { color: '#fff' }
      : { color: textColorForHue(displayColor) };
    const backgroundStyle = isSystemMessage
      ? { backgroundColor: '#667eea' }
      : { backgroundColor: messageBubbleColorForHue(displayColor) };
    const messageClassString = isSystemMessage
      ? 'message flex flex-row items-start p-4 m-2 rounded-lg shadow-l border-solid border-indigo-700 border-2 border-opacity-60 text-l'
      : `message relative flex flex-row items-start p-3 m-3 rounded-lg shadow-s text-sm ${
          isMessageModeratable ? 'moderatable' : ''
        }`;

    const isModeratorFlair = isAuthorModerator
      ? html`<img
          class="flair"
          title="Moderator"
          src="/img/moderator-nobackground.svg"
        />`
      : null;

    const isBotFlair = isBot
      ? html`<img
          title="Bot"
          class="inline-block mr-1 w-4 h-4 relative"
          style=${{ bottom: '1px' }}
          src="/img/bot.svg"
        />`
      : null;

    const authorAuthenticatedFlair = authenticated
      ? html`<img
          class="flair"
          title="Authenticated"
          src="/img/authenticated.svg"
        />`
      : null;

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
            title=${userMetadata}
          >
            ${isBotFlair} ${authorAuthenticatedFlair} ${isModeratorFlair}
            ${displayName}
          </div>
          ${isMessageModeratable &&
          html`<${ModeratorActions}
            message=${message}
            accessToken=${accessToken}
          />`}
          <div
            class="message-text text-gray-300 font-normal overflow-y-hidden pt-2"
            dangerouslySetInnerHTML=${{ __html: formattedMessage }}
          ></div>
        </div>
      </div>
    `;
  }
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
  return new Promise((res) => {
    new Mark(node).mark(username, {
      element: 'span',
      className: 'highlighted px-1 rounded font-bold bg-orange-500',
      separateWordSearch: false,
      accuracy: {
        value: 'exactly',
        limiters: [',', '.', "'", '?', '@'],
      },
      done() {
        res(node.innerHTML);
      },
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
    if (url.indexOf('instagram.com/p/') > -1) {
      embedText += getInstagramEmbedFromURL(url);
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

function getInstagramEmbedFromURL(url) {
  const urlObject = new URL(url.replace(/\/$/, ''));
  urlObject.pathname += '/embed';
  return `<iframe class="chat-embed instagram-embed" src="${urlObject.href}" frameborder="0" allowfullscreen></iframe>`;
}

function isMessageJustAnchor(message, anchor) {
  return stripTags(message) === stripTags(anchor.innerHTML);
}

function formatTimestamp(sentAt) {
  sentAt = new Date(sentAt);
  if (isNaN(sentAt)) {
    return '';
  }

  let diffInDays = getDiffInDaysFromNow(sentAt);
  if (diffInDays >= 1) {
    return (
      `at ${sentAt.toLocaleDateString('en-US', {
        dateStyle: 'medium',
      })} at ` + sentAt.toLocaleTimeString()
    );
  }

  return `${sentAt.toLocaleTimeString()}`;
}

/*
  You would call this when receiving a plain text
  value back from an API, and before inserting the
  text into the `contenteditable` area on a page.
*/
function convertToMarkup(str = '') {
  return convertToText(str).replace(/\n/g, '<p></p>');
}

function stripTags(str) {
  return str.replace(/<\/?[^>]+(>|$)/g, '');
}
