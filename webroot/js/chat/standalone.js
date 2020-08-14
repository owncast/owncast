import { html, Component } from "https://unpkg.com/htm/preact/index.mjs?module";
import UserInfo from './user-info.js';
import Chat from './chat.js';
import Websocket from '../websocket.js';

import { getLocalStorage, generateAvatar, generateUsername } from '../utils.js';
import { KEY_USERNAME, KEY_AVATAR } from '../utils/chat.js';

export default class StandaloneChat extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      websocket: new Websocket(),
      chatEnabled: true, // always true for standalone chat
      username: getLocalStorage(KEY_USERNAME) || generateUsername(),
      userAvatarImage: getLocalStorage(KEY_AVATAR) || generateAvatar(`${this.username}${Date.now()}`),
    };

    this.websocket = null;
    this.handleUsernameChange = this.handleUsernameChange.bind(this);
  }

  handleUsernameChange(newName, newAvatar) {
    this.setState({
      username: newName,
      userAvatarImage: newAvatar,
    });
  }

  handleChatToggle() {
    return;
  }

  render(props, state) {
    const { username, userAvatarImage, websocket } = state;
    return (
      html`
        <div class="flex">
          <${UserInfo}
            username=${username}
            userAvatarImage=${userAvatarImage}
            handleUsernameChange=${this.handleUsernameChange}
            handleChatToggle=${this.handleChatToggle}
          />
          <${Chat}
            websocket=${websocket}
            username=${username}
            userAvatarImage=${userAvatarImage}
            chatEnabled />
        </div>
    `);
  }

}
