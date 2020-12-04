import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import ChatMessageView from './chat-message-view.js';

import { messageBubbleColorForString } from '../../utils/user-colors.js';
import { SOCKET_MESSAGE_TYPES } from '../../utils/websocket.js';

export default class Message extends Component {
  render(props) {
    const { message } = props;
    const { type } = message;
    if (type === SOCKET_MESSAGE_TYPES.CHAT || type === SOCKET_MESSAGE_TYPES.SYSTEM) {
      return html`<${ChatMessageView} ...${props} />`;
    } else if (type === SOCKET_MESSAGE_TYPES.NAME_CHANGE) {
      const { oldName, newName } = message;
      return (
        html`
          <div class="message message-name-change flex items-center justify-start p-3">
            <div class="message-content flex flex-row items-center justify-center text-sm w-full">
              <div class="text-white text-center opacity-50 overflow-hidden break-words">
                <span class="font-bold">${oldName}</span> is now known as <span class="font-bold">${newName}</span>.
              </div>
            </div>
          </div>
        `
      );
    } else {
      console.log("Unknown message type:", type);
    }
  }
}

