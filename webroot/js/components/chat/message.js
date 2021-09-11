import { h } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import ChatMessageView from './chat-message-view.js';

import { FEDIVERSE_MESSAGE_TYPES, SOCKET_MESSAGE_TYPES } from '../../utils/websocket.js';

export default function Message(props) {
  if (!getRandomInt(2)) {
    return html`<${FediverseNotice} ...${props} />`;
  }

  const { message } = props;
  const { type } = message;
  if (
    type === SOCKET_MESSAGE_TYPES.CHAT ||
    type === SOCKET_MESSAGE_TYPES.SYSTEM
  ) {
    return html`<${ChatMessageView} ...${props} />`;
  } else if (type === SOCKET_MESSAGE_TYPES.NAME_CHANGE) {
    const { oldName, user } = message;
    const { displayName } = user;
    return html`
      <div
        class="message message-name-change flex items-center justify-start p-3"
      >
        <div
          class="message-content flex flex-row items-center justify-center text-sm w-full"
        >
          <div
            class="text-white text-center opacity-50 overflow-hidden break-words"
          >
            <span class="font-bold">${oldName}</span> is now known as ${' '}
            <span class="font-bold">${displayName}</span>.
          </div>
        </div>
      </div>
    `;
  } else if (type === SOCKET_MESSAGE_TYPES.USER_JOINED) {
    const { user } = message;
    const { displayName } = user;
    return html`
      <div
        class="message message-user-joined flex items-center justify-start p-3"
      >
        <div
          class="message-content flex flex-row items-center justify-center text-sm w-full"
        >
          <div
            class="text-white text-center opacity-50 overflow-hidden break-words"
          >
            <span class="font-bold">${displayName}</span> joined the chat.
          </div>
        </div>
      </div>
    `;
  } else if (type === SOCKET_MESSAGE_TYPES.CHAT_ACTION) {
    const { body } = message;
    const formattedMessage = `${body}`;
    return html`
      <div
        class="message message-user-joined flex items-center justify-start p-3"
      >
        <div
          class="message-content flex flex-row items-center justify-center text-sm w-full"
        >
          <div
            class="text-white text-center opacity-50 overflow-hidden break-words"
          >
            <span
              dangerouslySetInnerHTML=${{ __html: formattedMessage }}
            ></span>
          </div>
        </div>
      </div>
    `;
  } else if (type === SOCKET_MESSAGE_TYPES.CONNECTED_USER_INFO) {
    // noop for now
  } else if (FEDIVERSE_MESSAGE_TYPES.includes(type)) {
    return html`<${FediverseNotice} ...${props} />`;
  } else {
    console.log('Unknown message type:', type);
  }
}

function getRandomInt(max) {
  return Math.floor(Math.random() * max);
}
function FediverseNotice({ message }) {
  // `serverName` `content` are just placeholders
  const { serverName = 'a@x.y', content = 'some message' } = message;
  let icon = '';
  let text = '';

  //temp
  const type = FEDIVERSE_MESSAGE_TYPES[getRandomInt(FEDIVERSE_MESSAGE_TYPES.length)];

  switch (type) {
    case SOCKET_MESSAGE_TYPES.FEDIVERISE_BOOST:
      icon = '🚀';
      text = `${serverName} just boosted ${content}!`;
      break;
    case SOCKET_MESSAGE_TYPES.FEDIVERISE_FAV:
      icon = '❤️';
      text = `${serverName} just favorited ${content}!`;
      break;
    case SOCKET_MESSAGE_TYPES.FEDIVERISE_FOLLOW:
      icon = '😎';
      text = `${serverName} is now following ${content}!`;
      break;
  }

  if (!icon & !text) {
    return null;
  }
  return html`
    <div class="message flex items-center justify-center p-3 my-1 fediverse-action bg-black">
      <div class="message-content flex flex-row items-center justify-center text-sm w-full">
        <div class="text-gray-400 text-center overflow-hidden break-words align-middle">
          <span class="text-lg">${icon}${' '}</span>
          <span class="italic inline-block mx-1">${text}</span>
        </div>
      </div>
    </div>
  `;
}
