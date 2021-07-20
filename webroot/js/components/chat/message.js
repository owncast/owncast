import { h } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import ChatMessageView from './chat-message-view.js';

import { SOCKET_MESSAGE_TYPES } from '../../utils/websocket.js';

export default function Message(props) {
  const { message } = props;
  const { type } = message;
  if (type === SOCKET_MESSAGE_TYPES.CHAT || type === SOCKET_MESSAGE_TYPES.SYSTEM) {
    return html`<${ChatMessageView} ...${props} />`;
  } else if (type === SOCKET_MESSAGE_TYPES.NAME_CHANGE) {
    const { oldName, user } = message;
    const { displayName } = user;  
    return (
      html`
        <div class="message message-name-change flex items-center justify-start p-3">
          <div class="message-content flex flex-row items-center justify-center text-sm w-full">
            <div class="text-white text-center opacity-50 overflow-hidden break-words">
              <span class="font-bold">${oldName}</span> is now known as <span class="font-bold">${displayName}</span>.
            </div>
          </div>
        </div>
      `
    );
  } else if (type === SOCKET_MESSAGE_TYPES.USER_JOINED) {
    const { user } = message
    const { displayName } = user;
    return (
      html`
          <div class="message message-user-joined flex items-center justify-start p-3">
            <div class="message-content flex flex-row items-center justify-center text-sm w-full">
              <div class="text-white text-center opacity-50 overflow-hidden break-words">
                <span class="font-bold">${displayName}</span> joined the chat.
              </div>
            </div>
          </div>
        `
    );
  } else if (type === SOCKET_MESSAGE_TYPES.CHAT_ACTION) {
    const { body } = message;
    const formattedMessage = `${body}`
    return (
      html`
          <div class="message message-user-joined flex items-center justify-start p-3">
            <div class="message-content flex flex-row items-center justify-center text-sm w-full">
              <div class="text-white text-center opacity-50 overflow-hidden break-words">
                <span dangerouslySetInnerHTML=${{ __html: formattedMessage }}></span>
              </div>
            </div>
          </div>
        `
    );
  } else if (type === SOCKET_MESSAGE_TYPES.CONNECTED_USER_INFO) {
    // noop for now
  } else {
    console.log("Unknown message type:", type);
  }
}

