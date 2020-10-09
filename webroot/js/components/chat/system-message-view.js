import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import { messageBubbleColorForString } from '../../utils/user-colors.js';
import { formatMessageText, formatTimestamp } from '../../utils/chat.js';
import { generateAvatar } from '../../utils/helpers.js';

export default class SystemMessageView extends Component {
  render() {
    const { message, username } = this.props;
    const { image, author, body, timestamp } = message;
    const formattedMessage = formatMessageText(body, username);
    const avatar = image || generateAvatar(author);
    const formattedTimestamp = formatTimestamp(timestamp);

    const authorColor = messageBubbleColorForString(author);
    const avatarBgColor = { backgroundColor: authorColor };
    const authorTextColor = { color: authorColor };
    return html`
      <div
        class="message flex flex-row items-start p-3 m-2 bg-orange-700 bg-opacity-25 rounded-lg shadow-xs"
      >
        <div class="message-content break-words w-full">
          <div
            class="message-author text-white font-bold text-lg"
            style=${authorTextColor}
          >
            <img src=${avatar} class="h-6 w-6 inline" /> ${author}
          </div>
          <div
            class="message-text text-gray-300 font-normal overflow-y-hidden"
            title=${formattedTimestamp}
            dangerouslySetInnerHTML=${{ __html: formattedMessage }}
          ></div>
        </div>
      </div>
    `;
  }
}
