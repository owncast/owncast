import { html, Component } from "https://unpkg.com/htm/preact/index.mjs?module";
import UserInfo from './user-info.js';
import Chat from './chat.js';

import { getLocalStorage, generateAvatar, generateUsername } from '../utils.js';
import { KEY_USERNAME, KEY_AVATAR } from '../utils/chat.js';

export class StandaloneChat extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      chatEnabled: true, // always true for standalone chat
      username: getLocalStorage(KEY_USERNAME) || generateUsername(),
      userAvatarImage: getLocalStorage(KEY_AVATAR) || generateAvatar(`${this.username}${Date.now()}`),
    };

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
    const { username, userAvatarImage } = state;
    return (
      html`
        <div class="flex">
          <${UserInfo}
            username=${username}
            userAvatarImage=${userAvatarImage}
            handleUsernameChange=${this.handleUsernameChange}
            handleChatToggle=${this.handleChatToggle}
          />
          <${Chat} username=${username} userAvatarImage=${userAvatarImage} chatEnabled />
        </div>
    `);
  }

}
