import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import Chat from './components/chat/chat.js';
import Websocket from './utils/websocket.js';
import { getLocalStorage, generateUsername } from './utils/helpers.js';
import { KEY_USERNAME } from './utils/constants.js';

export default class StandaloneChat extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      websocket: new Websocket(true), // Send along the "ignoreClient" flag so this isn't counted as a viewer
      chatEnabled: true, // always true for standalone chat
      username: getLocalStorage(KEY_USERNAME) || generateUsername(),
    };

    this.websocket = null;
    this.handleUsernameChange = this.handleUsernameChange.bind(this);
  }

  handleUsernameChange(newName) {
    this.setState({
      username: newName,
    });
  }

  render(props, state) {
    const { username, websocket } = state;
    return (
      html`
        <${Chat}
          websocket=${websocket}
          username=${username}
          messagesOnly
          ignoreClient
        />
      `
    );
  }
}
