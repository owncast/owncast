import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import ChatMessageView from './chat-message-view.js';
import SystemMessageView from './system-message-view.js';

import { SOCKET_MESSAGE_TYPES } from '../../utils/websocket.js';

export default class Message extends Component {
  render(props) {
    const { message } = props;
    const { type } = message;

    if (type === SOCKET_MESSAGE_TYPES.CHAT) {
      return html`<${ChatMessageView} ...${props} />`;
    } else if (type === SOCKET_MESSAGE_TYPES.SYSTEM) {
      return html`<${SystemMessageView} ...${props} />`;
    } else if (type === SOCKET_MESSAGE_TYPES.NAME_CHANGE) {
      const { oldName, newName, image } = message;
      return (
        html`
          <div class="message message-name-change flex items-center justify-start p-3">
            <div class="message-content flex flex-row items-center justify-center text-sm w-full">
              <div
                class="message-avatar rounded-full mr-3 bg-gray-900"
              >
                <img class="mr-2 p-1" src=${image} />
              </div>

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
