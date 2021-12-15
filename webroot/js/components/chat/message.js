import { h } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import ChatMessageView from './chat-message-view.js';

import { SOCKET_MESSAGE_TYPES } from '../../utils/websocket.js';
import { checkIsModerator } from '../../utils/chat.js';

function SystemMessage(props) {
  const { contents } = props;
  return html`
    <div
      class="message message-name-change flex items-center justify-start p-3"
    >
      <div
        class="message-content flex flex-row items-center justify-center text-sm w-full"
      >
        <div
          class="text-gray-400 w-full text-center opacity-90 overflow-hidden break-words"
        >
          ${contents}
        </div>
      </div>
    </div>
  `;
}

function SingleFederatedUser(props) {
  const { message } = props;
  const { type, title, subtitle, image, link } = message;

  let body = '';
  let icon = null;
  switch (type) {
    case SOCKET_MESSAGE_TYPES.FEDIVERSE_ENGAGEMENT_FOLLOW:
      body = html`<span>just followed this stream.</span>`;
      icon = '/img/follow.svg';
      break;
    case SOCKET_MESSAGE_TYPES.FEDIVERSE_ENGAGEMENT_LIKE:
      body = html`just liked a post.`;
      icon = '/img/like.svg';
      break;
    case SOCKET_MESSAGE_TYPES.FEDIVERSE_ENGAGEMENT_REPOST:
      body = html`just shared a post.`;
      icon = '/img/repost.svg';
      break;
    default:
      body = '';
      break;
  }

  return html`
    <a href=${link} target="_blank">
      <div
        class="federated-action m-2 bg-white flex items-center p-1 px-2 rounded-xl shadow border"
      >
        <div class="relative flex items-center space-x-4">
          <img
            src="${image || '/img/logo.svg'}"
            class="w-12 h-12 rounded-full"
          />
          <span
            style=${{
              backgroundImage: `url(${icon})`,
            }}
            class="absolute h-6 w-6 rounded-full border-2 border-white action-icon"
          ></span>
        </div>
        <div class="p-3 flex-grow">
          <span class="font-semibold text-gray-700 truncate"> ${title} </span>
          ${body}
          <p class=" text-gray-500 truncate">@${subtitle}</p>
        </div>
      </div>
    </a>
  `;
}

export default function Message(props) {
  const { message } = props;
  const { type, oldName, user, body } = message;
  if (
    type === SOCKET_MESSAGE_TYPES.CHAT ||
    type === SOCKET_MESSAGE_TYPES.SYSTEM
  ) {
    return html`<${ChatMessageView} ...${props} />`;
  } else if (type === SOCKET_MESSAGE_TYPES.NAME_CHANGE) {
    // User changed their name
    const { displayName } = user;
    const contents = html`
      <div>
        <span class="font-bold">${oldName}</span> is now known as ${' '}
        <span class="font-bold">${displayName}</span>.
      </div>
    `;
    return html`<${SystemMessage} contents=${contents} />`;
  } else if (type === SOCKET_MESSAGE_TYPES.USER_JOINED) {
    const { displayName } = user;
    const contents = html`<span class="font-bold">${displayName}</span> joined
      the chat.`;
    return html`<${SystemMessage} contents=${contents} />`;
  } else if (type === SOCKET_MESSAGE_TYPES.CHAT_ACTION) {
    const contents = html`<span
      dangerouslySetInnerHTML=${{ __html: body }}
    ></span>`;
    return html`<${SystemMessage} contents=${contents} />`;
  } else if (type === SOCKET_MESSAGE_TYPES.CONNECTED_USER_INFO) {
    // moderator message
    const isModerator = checkIsModerator(message);
    if (isModerator) {
      const contents = html`<div class="rounded-lg bg-gray-700 p-3">
        <img src="/img/moderator.svg" class="moderator-flag" />You are now a
        moderator.
      </div>`;
      return html`<${SystemMessage} contents=${contents} />`;
    }
  } else if (
    type === SOCKET_MESSAGE_TYPES.FEDIVERSE_ENGAGEMENT_FOLLOW ||
    SOCKET_MESSAGE_TYPES.FEDIVERSE_ENGAGEMENT_LIKE ||
    SOCKET_MESSAGE_TYPES.FEDIVERSE_ENGAGEMENT_REPOST
  ) {
    return html` <${SingleFederatedUser} message=${message} /> `;
  } else {
    console.log('Unknown message type:', type);
  }
}
