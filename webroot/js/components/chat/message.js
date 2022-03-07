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
  const { type, body, title, image, link } = message;

  let icon = null;
  switch (type) {
    case SOCKET_MESSAGE_TYPES.FEDIVERSE_ENGAGEMENT_FOLLOW:
      icon = '/img/follow.svg';
      break;
    case SOCKET_MESSAGE_TYPES.FEDIVERSE_ENGAGEMENT_LIKE:
      icon = '/img/like.svg';
      break;
    case SOCKET_MESSAGE_TYPES.FEDIVERSE_ENGAGEMENT_REPOST:
      icon = '/img/repost.svg';
      break;
    default:
      break;
  }

  return html`
    <a
      href=${link}
      target="_blank"
      class="hover:no-underline"
      title="Visit profile"
    >
      <div
        class="federated-action m-2 mt-3 bg-white flex items-center px-2 rounded-xl shadow border"
      >
        <div class="relative" style="top: -6px">
          <img
            src="${image || '/img/logo.svg'}"
            style="max-width: unset"
            class="rounded-full border border-slate-500	w-16"
          />
          <span
            style=${{ backgroundImage: `url(${icon})` }}
            class="absolute h-6 w-6 rounded-full border-2 border-white action-icon"
          ></span>
        </div>
        <div class="px-4 py-2 min-w-0">
          <div class="text-gray-500 text-sm hover:no-underline truncate">
            ${title}
          </div>
          <p
            class=" text-gray-700 w-full text-base leading-6"
            dangerouslySetInnerHTML=${{ __html: body }}
          ></p>
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
    const { displayName, isBot } = user;
    const isAuthorModerator = checkIsModerator(message);
    const messageAuthorFlair = isAuthorModerator
      ? html`<img
          title="Moderator"
          class="inline-block mr-1 w-3 h-3"
          src="/img/moderator-nobackground.svg"
        />`
      : null;

    const contents = html`<div>
      <span class="font-bold">${messageAuthorFlair}${displayName}</span>
      ${' '}joined the chat.
    </div>`;
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
