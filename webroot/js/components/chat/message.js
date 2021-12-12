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
  } else {
    console.log('Unknown message type:', type);
  }
}
