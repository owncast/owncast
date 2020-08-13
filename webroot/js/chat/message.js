import { h, Component, createRef } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
// Initialize htm with Preact
const html = htm.bind(h);

import { messageBubbleColorForString } from '../utils/user-colors.js';
import { formatMessageText } from '../utils/chat.js';
import { generateAvatar } from '../utils.js';
import SOCKET_MESSAGE_TYPES from '../utils/socket-message-types.js';

export default class Message extends Component {
  render(props) {
    const { message } = props;
    const { type } = message;

    if (type === SOCKET_MESSAGE_TYPES.CHAT) {
      const { image, author, body, type } = message;
      const formattedMessage = formatMessageText(body);
      const avatar = image || generateAvatar(author);
      const avatarBgColor = { backgroundColor: messageBubbleColorForString(author) };
      return (
        html`
          <div class="message flex">
            <div
              class="message-avatar rounded-full flex items-center justify-center"
              style=${avatarBgColor}
            >
              <img src=${avatar} />
            </div>
            <div class="message-content">
              <p class="message-author text-white font-bold">${author}</p>
              <p class="message-text text-gray-400 font-thin">${formattedMessage}</p>
            </div>
        </div>
      `);
    } else if (type === SOCKET_MESSAGE_TYPES.NAME_CHANGE) {
      const { oldName, newName, image } = message;
      return (
        html`
          <div class="message flex">
            <img class="mr-2" src=${image} />
            <div class="text-white text-center">
              <span class="font-bold">${oldName}</span> is now known as <span class="font-bold">${newName}</span>.
            </div>
          </div>
        `
      )
    }
  }
}
