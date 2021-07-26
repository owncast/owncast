import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import Chat from './components/chat/chat.js';
import Websocket from './utils/websocket.js';
import { getLocalStorage, setLocalStorage } from './utils/helpers.js';
import { KEY_EMBED_CHAT_ACCESS_TOKEN } from './utils/constants.js';
import { registerChat } from './chat/register.js';

export default class StandaloneChat extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      chatEnabled: true, // always true for standalone chat
      username: null,
    };

    this.isRegistering = false;
    this.hasConfiguredChat = false;
    this.websocket = null;
    this.handleUsernameChange = this.handleUsernameChange.bind(this);

    // If this is the first time setting the config
    // then setup chat if it's enabled.
    const chatBlocked = getLocalStorage('owncast_chat_blocked');
    if (!chatBlocked && !this.hasConfiguredChat) {
      this.setupChatAuth();
    }

    this.hasConfiguredChat = true;
  }

  handleUsernameChange(newName) {
    this.setState({
      username: newName,
    });
  }

  async setupChatAuth(force) {
    var accessToken = getLocalStorage(KEY_EMBED_CHAT_ACCESS_TOKEN);
    const randomInt = Math.floor(Math.random() * 100) + 1;
    var username = 'chat-embed-' + randomInt;

    if (!accessToken || force) {
      try {
        this.isRegistering = true;
        const registration = await registerChat(username);
        accessToken = registration.accessToken;
        username = registration.displayName;

        setLocalStorage(KEY_EMBED_CHAT_ACCESS_TOKEN, accessToken);

        this.isRegistering = false;
      } catch (e) {
        console.error('registration error:', e);
      }
    }

    if (this.state.websocket) {
      this.state.websocket.shutdown();
      this.setState({
        websocket: null,
      });
    }

    // Without a valid access token he websocket connection will be rejected.
    const websocket = new Websocket(accessToken);

    this.setState({
      username,
      websocket,
      accessToken,
    });
  }

  render(props, state) {
    const { username, websocket, accessToken } = state;
    const { messagesOnly } = props;
    return html`
      <${Chat}
        websocket=${websocket}
        username=${username}
        accessToken=${accessToken}
        messagesOnly=${messagesOnly}
      />
    `;
  }
}
