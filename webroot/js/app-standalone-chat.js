import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);
import UsernameForm from './components/chat/username.js';
import Chat from './components/chat/chat.js';
import Websocket, {
  CALLBACKS,
  SOCKET_MESSAGE_TYPES,
} from './utils/websocket.js';
import { registerChat } from './chat/register.js';
import { getLocalStorage, setLocalStorage } from './utils/helpers.js';
import {
  CHAT_MAX_MESSAGE_LENGTH,
  EST_SOCKET_PAYLOAD_BUFFER,
  KEY_EMBED_CHAT_ACCESS_TOKEN,
  KEY_ACCESS_TOKEN,
  KEY_USERNAME,
  TIMER_DISABLE_CHAT_AFTER_OFFLINE,
  URL_STATUS,
  URL_CONFIG,
  TIMER_STATUS_UPDATE,
} from './utils/constants.js';
import { URL_WEBSOCKET } from './utils/constants.js';

export default class StandaloneChat extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      websocket: null,
      canChat: false,
      chatEnabled: true, // always true for standalone chat
      chatInputEnabled: false, // chat input box state
      accessToken: null,
      username: null,
      isRegistering: false,
      streamOnline: null, // stream is active/online
      lastDisconnectTime: null,
      configData: {
        loading: true,
      },
    };
    this.disableChatInputTimer = null;
    this.hasConfiguredChat = false;

    this.handleUsernameChange = this.handleUsernameChange.bind(this);
    this.handleOfflineMode = this.handleOfflineMode.bind(this);
    this.handleOnlineMode = this.handleOnlineMode.bind(this);
    this.handleFormFocus = this.handleFormFocus.bind(this);
    this.handleFormBlur = this.handleFormBlur.bind(this);
    this.getStreamStatus = this.getStreamStatus.bind(this);
    this.getConfig = this.getConfig.bind(this);
    this.disableChatInput = this.disableChatInput.bind(this);
    this.setupChatAuth = this.setupChatAuth.bind(this);
    this.disableChat = this.disableChat.bind(this);

    this.socketHostOverride = null;

    // user events
    this.handleWebsocketMessage = this.handleWebsocketMessage.bind(this);

    this.getConfig();

    this.getStreamStatus();
    this.statusTimer = setInterval(this.getStreamStatus, TIMER_STATUS_UPDATE);
  }

  // fetch /config data
  getConfig() {
    fetch(URL_CONFIG)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`Network response was not ok ${response.ok}`);
        }
        return response.json();
      })
      .then((json) => {
        this.setConfigData(json);
      })
      .catch((error) => {
        this.handleNetworkingError(`Fetch config: ${error}`);
      });
  }

  // fetch stream status
  getStreamStatus() {
    fetch(URL_STATUS)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`Network response was not ok ${response.ok}`);
        }
        return response.json();
      })
      .then((json) => {
        this.updateStreamStatus(json);
      })
      .catch((error) => {
        this.handleOfflineMode();
        this.handleNetworkingError(`Stream status: ${error}`);
      });
  }

  setConfigData(data = {}) {
    const { chatDisabled, socketHostOverride } = data;

    // If this is the first time setting the config
    // then setup chat if it's enabled.
    if (!this.hasConfiguredChat && !chatDisabled) {
      this.setupChatAuth();
    }

    this.hasConfiguredChat = true;
    this.socketHostOverride = socketHostOverride;
    this.setState({
      canChat: !chatDisabled,
      configData: {
        ...data,
      },
    });
  }

  // handle UI things from stream status result
  updateStreamStatus(status = {}) {
    const { streamOnline: curStreamOnline } = this.state;

    if (!status) {
      return;
    }
    const { online, lastDisconnectTime } = status;

    this.setState({
      lastDisconnectTime,
      streamOnline: online,
    });

    if (status.online !== curStreamOnline) {
      if (status.online) {
        // stream has just come online.
        this.handleOnlineMode();
      } else {
        // stream has just flipped offline or app just got loaded and stream is offline.
        this.handleOfflineMode(lastDisconnectTime);
      }
    }
  }

  // stop status timer and disable chat after some time.
  handleOfflineMode(lastDisconnectTime) {
    if (lastDisconnectTime) {
      const remainingChatTime =
        TIMER_DISABLE_CHAT_AFTER_OFFLINE -
        (Date.now() - new Date(lastDisconnectTime));
      const countdown = remainingChatTime < 0 ? 0 : remainingChatTime;
      if (countdown > 0) {
        this.setState({
          chatInputEnabled: true,
        });
      }
      this.disableChatInputTimer = setTimeout(this.disableChatInput, countdown);
    }
    this.setState({
      streamOnline: false,
    });
  }

  handleOnlineMode() {
    clearTimeout(this.disableChatInputTimer);
    this.disableChatInputTimer = null;

    this.setState({
      streamOnline: true,
      chatInputEnabled: true,
    });
  }

  handleUsernameChange(newName) {
    this.setState({
      username: newName,
    });
    this.sendUsernameChange(newName);
  }

  disableChatInput() {
    this.setState({
      chatInputEnabled: false,
    });
  }

  handleNetworkingError(error) {
    console.error(`>>> App Error: ${error}`);
  }

  handleWebsocketMessage(e) {
    if (e.type === SOCKET_MESSAGE_TYPES.ERROR_USER_DISABLED) {
      // User has been actively disabled on the backend. Turn off chat for them.
      this.handleBlockedChat();
    } else if (
      e.type === SOCKET_MESSAGE_TYPES.ERROR_NEEDS_REGISTRATION &&
      !this.isRegistering
    ) {
      // User needs an access token, so start the user auth flow.
      this.state.websocket.shutdown();
      this.setState({ websocket: null });
      this.setupChatAuth(true);
    } else if (e.type === SOCKET_MESSAGE_TYPES.ERROR_MAX_CONNECTIONS_EXCEEDED) {
      // Chat server cannot support any more chat clients. Turn off chat for them.
      this.disableChat();
    } else if (e.type === SOCKET_MESSAGE_TYPES.CONNECTED_USER_INFO) {
      // When connected the user will return an event letting us know what our
      // user details are so we can display them properly.
      const { user } = e;
      const { displayName } = user;

      this.setState({ username: displayName });
    }
  }

  handleBlockedChat() {
    setLocalStorage('owncast_chat_blocked', true);
    this.disableChat();
  }

  handleFormFocus() {
    if (this.hasTouchScreen) {
      this.setState({
        touchKeyboardActive: true,
      });
    }
  }

  handleFormBlur() {
    if (this.hasTouchScreen) {
      this.setState({
        touchKeyboardActive: false,
      });
    }
  }

  disableChat() {
    this.state.websocket.shutdown();
    this.setState({ websocket: null, canChat: false });
  }

  async setupChatAuth(force) {
    const { readonly } = this.props;
    var accessToken = readonly
      ? getLocalStorage(KEY_EMBED_CHAT_ACCESS_TOKEN)
      : getLocalStorage(KEY_ACCESS_TOKEN);
    var randomIntArray = new Uint32Array(1);
    window.crypto.getRandomValues(randomIntArray);
    var username = readonly
      ? 'chat-embed-' + randomIntArray[0]
      : getLocalStorage(KEY_USERNAME);

    if (!accessToken || force) {
      try {
        this.isRegistering = true;
        const registration = await registerChat(username);
        accessToken = registration.accessToken;
        username = registration.displayName;

        if (readonly) {
          setLocalStorage(KEY_EMBED_CHAT_ACCESS_TOKEN, accessToken);
        } else {
          setLocalStorage(KEY_ACCESS_TOKEN, accessToken);
          setLocalStorage(KEY_USERNAME, username);
        }

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
    const websocket = new Websocket(
      accessToken,
      this.socketHostOverride || URL_WEBSOCKET
    );
    websocket.addListener(
      CALLBACKS.RAW_WEBSOCKET_MESSAGE_RECEIVED,
      this.handleWebsocketMessage
    );

    this.setState({
      username,
      websocket,
      accessToken,
    });
  }

  sendUsernameChange(newName) {
    const nameChange = {
      type: SOCKET_MESSAGE_TYPES.NAME_CHANGE,
      newName,
    };
    this.state.websocket.send(nameChange);
  }

  render(props, state) {
    const { username, websocket, accessToken, chatInputEnabled, configData } =
      state;

    const { chatDisabled, maxSocketPayloadSize, customStyles, name } =
      configData;

    const { readonly } = props;
    return this.state.websocket
      ? html`${!readonly
            ? html`<style>
                  ${customStyles}
                </style>
                <header
                  class="flex flex-row-reverse fixed z-10 w-full bg-gray-900"
                >
                  <${UsernameForm}
                    username=${username}
                    onUsernameChange=${this.handleUsernameChange}
                    onFocus=${this.handleFormFocus}
                    onBlur=${this.handleFormBlur}
                  />
                </header>`
            : ''}
          <${Chat}
            websocket=${websocket}
            username=${username}
            accessToken=${accessToken}
            readonly=${readonly}
            instanceTitle=${name}
            chatInputEnabled=${chatInputEnabled && !chatDisabled}
            inputMaxBytes=${maxSocketPayloadSize - EST_SOCKET_PAYLOAD_BUFFER ||
            CHAT_MAX_MESSAGE_LENGTH}
          />`
      : null;
  }
}
