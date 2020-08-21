import { html, Component } from "https://unpkg.com/htm/preact/index.mjs?module";

import { messageBubbleColorForString } from '../utils/user-colors.js';
import { formatMessageText } from '../utils/chat.js';
import { generateAvatar } from '../utils.js';
import SOCKET_MESSAGE_TYPES from '../utils/socket-message-types.js';

export default class Message extends Component {
  render(props) {
    const { message, username } = props;
    const { type } = message;

    if (type === SOCKET_MESSAGE_TYPES.CHAT) {
      const { image, author, body } = message;
      const formattedMessage = formatMessageText(body, username);
      const avatar = image || generateAvatar(author);

      const authorColor = messageBubbleColorForString(author);
      const avatarBgColor = { backgroundColor: authorColor };
      const authorTextColor = { color: authorColor };
      return (
        html`
          <div class="message flex">
            <div
              class="message-avatar rounded-full flex items-center justify-center"
              style=${avatarBgColor}
            >
              <img src=${avatar} />
            </div>
            <div class="message-content text-sm">
              <p class="message-author text-white font-bold" style=${authorTextColor}>${author}</p>
              <div
                class="message-text text-gray-300 font-normal"
                dangerouslySetInnerHTML=${
                  { __html: formattedMessage }
                }
              ></div>
            </div>
        </div>
      `);
    } else if (type === SOCKET_MESSAGE_TYPES.NAME_CHANGE) {
      const { oldName, newName, image } = message;
      return (
        html`
          <div class="message flex">
            <div class="message-content text-sm">
              <img class="mr-2" src=${image} />
              <div class="text-white text-center">
                <span class="font-bold">${oldName}</span> is now known as <span class="font-bold">${newName}</span>.
              </div>
            </div>
          </div>
        `);
    }
  }
}
