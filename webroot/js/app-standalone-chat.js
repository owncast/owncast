import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import Chat from './components/chat/chat.js';
import Websocket from './utils/websocket.js';
import { getLocalStorage, generateAvatar, generateUsername } from './utils/helpers.js';
import { KEY_USERNAME, KEY_AVATAR } from './utils/constants.js';

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

  render(props, state) {
    const { username, userAvatarImage, websocket } = state;
    return (
      html`
        <${Chat}
          websocket=${websocket}
          username=${username}
          userAvatarImage=${userAvatarImage}
          messagesOnly
        />
      `
    );
  }
}
