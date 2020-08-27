import { h, Component } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
const html = htm.bind(h);

import { messageBubbleColorForString } from '../../utils/user-colors.js';
import { formatMessageText } from '../../utils/chat.js';
import { generateAvatar } from '../../utils/helpers.js';
import { SOCKET_MESSAGE_TYPES } from '../../utils/websocket.js';

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
          <div class="message flex flex-row items-start p-3">
            <div
              class="message-avatar rounded-full flex items-center justify-center mr-3"
              style=${avatarBgColor}
            >
              <img src=${avatar} class="p-1" />
            </div>
            <div class="message-content text-sm break-words w-full">
              <div class="message-author text-white font-bold" style=${authorTextColor}>
                ${author}
              </div>
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
