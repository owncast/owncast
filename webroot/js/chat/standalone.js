import { h, Component, Fragment } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
const html = htm.bind(h);


import UsernameForm from './username.js';
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
    const { messagesOnly } = props;
    const { username, userAvatarImage, websocket } = state;


    if (messagesOnly) {
      return (
      html`
        <${Chat}
          websocket=${websocket}
          username=${username}
          userAvatarImage=${userAvatarImage}
          messagesOnly
        />
      `);
    }

    // not needed for standalone, just messages only. remove later.
    return (
      html`
        <${Fragment}>
          <${UsernameForm}
            username=${username}
            userAvatarImage=${userAvatarImage}
            handleUsernameChange=${this.handleUsernameChange}
            handleChatToggle=${this.handleChatToggle}
          />

          <${Chat}
            websocket=${websocket}
            username=${username}
            userAvatarImage=${userAvatarImage}
          />
        </${Fragment}>
    `);
  }

}
