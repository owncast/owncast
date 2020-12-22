import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import { OwncastPlayer } from './components/player.js';
import SocialIconsList from './components/platform-logos-list.js';
import UsernameForm from './components/chat/username.js';
import VideoPoster from './components/video-poster.js';
import Chat from './components/chat/chat.js';
import Websocket from './utils/websocket.js';
import { parseSecondsToDurationString, hasTouchScreen, getOrientation } from './utils/helpers.js';

import {
  addNewlines,
  classNames,
  clearLocalStorage,
  debounce,
  generateUsername,
  getLocalStorage,
  pluralize,
  setLocalStorage,
} from './utils/helpers.js';
import {
  HEIGHT_SHORT_WIDE,
  KEY_CHAT_DISPLAYED,
  KEY_USERNAME,
  MESSAGE_OFFLINE,
  MESSAGE_ONLINE,
  ORIENTATION_PORTRAIT,
  OWNCAST_LOGO_LOCAL,
  TEMP_IMAGE,
  TIMER_DISABLE_CHAT_AFTER_OFFLINE,
  TIMER_STATUS_UPDATE,
  TIMER_STREAM_DURATION_COUNTER,
  URL_CONFIG,
  URL_OWNCAST,
  URL_STATUS,
  WIDTH_SINGLE_COL,
} from './utils/constants.js';

export default class App extends Component {
  constructor(props, context) {
    super(props, context);

    const chatStorage = getLocalStorage(KEY_CHAT_DISPLAYED);
    this.hasTouchScreen = hasTouchScreen();

    this.state = {
      websocket: new Websocket(),
      displayChat: chatStorage === null ? true : chatStorage,
      chatInputEnabled: false, // chat input box state
      username: getLocalStorage(KEY_USERNAME) || generateUsername(),
      touchKeyboardActive: false,

      configData: {},
      extraPageContent: '',

      playerActive: false, // player object is active
      streamOnline: false, // stream is active/online
      isPlaying: false, // player is actively playing video

      // status
      streamStatusMessage: MESSAGE_OFFLINE,
      viewerCount: '',

      // dom
      windowWidth: window.innerWidth,
      windowHeight: window.innerHeight,
      orientation: getOrientation(this.hasTouchScreen),
    };

    // timers
    this.playerRestartTimer = null;
    this.offlineTimer = null;
    this.statusTimer = null;
    this.disableChatTimer = null;
    this.streamDurationTimer = null;

    // misc dom events
    this.handleChatPanelToggle = this.handleChatPanelToggle.bind(this);
    this.handleUsernameChange = this.handleUsernameChange.bind(this);
    this.handleFormFocus = this.handleFormFocus.bind(this);
    this.handleFormBlur = this.handleFormBlur.bind(this);
    this.handleWindowResize = debounce(this.handleWindowResize.bind(this), 250);

    this.handleOfflineMode = this.handleOfflineMode.bind(this);
    this.handleOnlineMode = this.handleOnlineMode.bind(this);
    this.disableChatInput = this.disableChatInput.bind(this);
    this.setCurrentStreamDuration = this.setCurrentStreamDuration.bind(this);

    // player events
    this.handlePlayerReady = this.handlePlayerReady.bind(this);
    this.handlePlayerPlaying = this.handlePlayerPlaying.bind(this);
    this.handlePlayerEnded = this.handlePlayerEnded.bind(this);
    this.handlePlayerError = this.handlePlayerError.bind(this);

    // fetch events
    this.getConfig = this.getConfig.bind(this);
    this.getStreamStatus = this.getStreamStatus.bind(this);
  }

  componentDidMount() {
    this.getConfig();
    window.addEventListener('resize', this.handleWindowResize);
    if (this.hasTouchScreen) {
      window.addEventListener('orientationchange', this.handleWindowResize);
    }
    this.player = new OwncastPlayer();
    this.player.setupPlayerCallbacks({
      onReady: this.handlePlayerReady,
      onPlaying: this.handlePlayerPlaying,
      onEnded: this.handlePlayerEnded,
      onError: this.handlePlayerError,
    });
    this.player.init();
  }

  componentWillUnmount() {
    // clear all the timers
    clearInterval(this.playerRestartTimer);
    clearInterval(this.offlineTimer);
    clearInterval(this.statusTimer);
    clearTimeout(this.disableChatTimer);
    clearInterval(this.streamDurationTimer);
    window.removeEventListener('resize', this.handleWindowResize);
    if (this.hasTouchScreen) {
      window.removeEventListener('orientationchange', this.handleWindowResize);
    }
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
    const { title, summary } = data;
    window.document.title = title;

    this.setState({
      configData: {
        ...data,
        summary: summary && addNewlines(summary),
      },
    });
  }

  // handle UI things from stream status result
  updateStreamStatus(status = {}) {
    const { streamOnline: curStreamOnline } = this.state;

    if (!status) {
      return;
    }
    const {
      viewerCount,
      online,
      lastConnectTime,
    } = status;

    this.lastDisconnectTime = status.lastDisconnectTime;

    if (status.online && !curStreamOnline) {
      // stream has just come online.
      this.handleOnlineMode();
    } else if (!status.online && curStreamOnline) {
      // stream has just flipped offline.
      this.handleOfflineMode();
    }

    this.setState({
      viewerCount,
      lastConnectTime,
      streamOnline: online,
    });
  }

  // when videojs player is ready, start polling for stream
  handlePlayerReady() {
    this.getStreamStatus();
    this.statusTimer = setInterval(this.getStreamStatus, TIMER_STATUS_UPDATE);
  }

  handlePlayerPlaying() {
    this.setState({
      isPlaying: true,
    });
  }

  // likely called some time after stream status has gone offline.
  // basically hide video and show underlying "poster"
  handlePlayerEnded() {
    this.setState({
      playerActive: false,
      isPlaying: false,
    });
  }

  handlePlayerError() {
    // do something?
    this.handleOfflineMode();
    this.handlePlayerEnded();
  }

  // stop status timer and disable chat after some time.
  handleOfflineMode() {
    clearInterval(this.streamDurationTimer);
    const remainingChatTime =
      TIMER_DISABLE_CHAT_AFTER_OFFLINE -
      (Date.now() - new Date(this.lastDisconnectTime));
    const countdown = remainingChatTime < 0 ? 0 : remainingChatTime;
    this.disableChatTimer = setTimeout(this.disableChatInput, countdown);
    this.setState({
      streamOnline: false,
      streamStatusMessage: MESSAGE_OFFLINE,
    });

    if (this.player.vjsPlayer && this.player.vjsPlayer.paused()) {
      this.handlePlayerEnded();
    }
  }

  // play video!
  handleOnlineMode() {
    this.player.startPlayer();
    clearTimeout(this.disableChatTimer);
    this.disableChatTimer = null;

    this.streamDurationTimer = setInterval(
      this.setCurrentStreamDuration,
      TIMER_STREAM_DURATION_COUNTER
    );

    this.setState({
      playerActive: true,
      streamOnline: true,
      chatInputEnabled: true,
      streamStatusMessage: MESSAGE_ONLINE,
    });
  }

  setCurrentStreamDuration() {
    let streamDurationString = '';
    if (this.state.lastConnectTime) {
      const diff = (Date.now() - Date.parse(this.state.lastConnectTime)) / 1000;
      streamDurationString = parseSecondsToDurationString(diff);
    }
    this.setState({
      streamStatusMessage: `${MESSAGE_ONLINE} ${streamDurationString}`,
    });

  }

  handleUsernameChange(newName) {
    this.setState({
      username: newName,
    });
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

  handleChatPanelToggle() {
    const { displayChat: curDisplayed } = this.state;

    const displayChat = !curDisplayed;
    if (displayChat) {
      setLocalStorage(KEY_CHAT_DISPLAYED, displayChat);
    } else {
      clearLocalStorage(KEY_CHAT_DISPLAYED);
    }
    this.setState({
      displayChat,
    });
  }

  disableChatInput() {
    this.setState({
      chatInputEnabled: false,
    });
  }

  handleNetworkingError(error) {
    console.log(`>>> App Error: ${error}`);
  }

  handleWindowResize() {
    this.setState({
      windowWidth: window.innerWidth,
      windowHeight: window.innerHeight,
      orientation: getOrientation(this.hasTouchScreen),
    });
  }

  render(props, state) {
    const {
      chatInputEnabled,
      configData,
      displayChat,
      isPlaying,
      orientation,
      playerActive,
      streamOnline,
      streamStatusMessage,
      touchKeyboardActive,
      username,
      viewerCount,
      websocket,
      windowHeight,
      windowWidth,
    } = state;


    const {
      version: appVersion,
      logo = TEMP_IMAGE,
      socialHandles = [],
      name: streamerName,
      summary,
      tags = [],
      title,
      extraPageContent,
    } = configData;

    const bgUserLogo = { backgroundImage: `url(${logo})` };

    const tagList = (tags !== null && tags.length > 0)
      ? tags.map(
          (tag, index) => html`
            <li
              key="tag${index}"
              class="tag rounded-sm text-gray-100 bg-gray-700 text-xs uppercase mr-3 p-2 whitespace-no-wrap"
            >
              ${tag}
            </li>
          `
        )
      : null;

    const mainClass = playerActive ? 'online' : '';
    const streamInfoClass = streamOnline ? 'online' : ''; // need?
    const isPortrait = this.hasTouchScreen && orientation === ORIENTATION_PORTRAIT;
    const shortHeight = windowHeight <= HEIGHT_SHORT_WIDE && !isPortrait;
    const singleColMode = windowWidth <= WIDTH_SINGLE_COL && !shortHeight;

    const extraAppClasses = classNames({
      chat: displayChat,
      'no-chat': !displayChat,
      'single-col': singleColMode,
      'bg-gray-800': singleColMode && displayChat,
      'short-wide': shortHeight && windowWidth > WIDTH_SINGLE_COL,
      'touch-screen': this.hasTouchScreen,
      'touch-keyboard-active': touchKeyboardActive,
    });

    const poster = isPlaying ? null : html`
      <${VideoPoster} offlineImage=${logo} active=${streamOnline} />
    `;

    return html`
      <div
        id="app-container"
        class="flex w-full flex-col justify-start relative ${extraAppClasses}"
      >
        <div id="top-content" class="z-50">
          <header
            class="flex border-b border-gray-900 border-solid shadow-md fixed z-10 w-full top-0	left-0 flex flex-row justify-between flex-no-wrap"
          >
            <h1
              class="flex flex-row items-center justify-start p-2 uppercase text-gray-400 text-xl	font-thin tracking-wider overflow-hidden whitespace-no-wrap"
            >
              <span
                id="logo-container"
                class="inline-block	rounded-full bg-white w-8 min-w-8 min-h-8 h-8 mr-2 bg-no-repeat bg-center"
              >
                <img class="logo visually-hidden" src=${OWNCAST_LOGO_LOCAL} alt="owncast logo" />
              </span>
              <span class="instance-title overflow-hidden truncate"
                >${title}</span
              >
            </h1>
            <div
              id="user-options-container"
              class="flex flex-row justify-end items-center flex-no-wrap"
            >
              <${UsernameForm}
                username=${username}
                onUsernameChange=${this.handleUsernameChange}
                onFocus=${this.handleFormFocus}
                onBlur=${this.handleFormBlur}
              />
              <button
                type="button"
                id="chat-toggle"
                onClick=${this.handleChatPanelToggle}
                class="flex cursor-pointer text-center justify-center items-center min-w-12 h-full bg-gray-800 hover:bg-gray-700"
              >
                💬
              </button>
            </div>
          </header>
        </div>

        <main class=${mainClass}>
          <div
            id="video-container"
            class="flex owncast-video-container bg-black w-full bg-center bg-no-repeat flex flex-col items-center justify-start"
          >
            <video
              class="video-js vjs-big-play-centered display-block w-full h-full"
              id="video"
              preload="auto"
              controls
              playsinline
            ></video>
            ${poster}
          </div>

          <section
            id="stream-info"
            aria-label="Stream status"
            class="flex text-center flex-row justify-between font-mono py-2 px-8 bg-gray-900 text-indigo-200 shadow-md border-b border-gray-100 border-solid ${streamInfoClass}"
          >
            <span>${streamStatusMessage}</span>
            <span>${viewerCount} ${pluralize('viewer', viewerCount)}.</span>
          </section>
        </main>

        <section id="user-content" aria-label="User information" class="p-8">
          <div class="user-content flex flex-row p-8">
            <div
              class="user-image rounded-full bg-white p-4 mr-8 bg-no-repeat bg-center"
              style=${bgUserLogo}
            >
              <img class="logo visually-hidden" alt="" src=${logo} />
            </div>
            <div
              class="user-content-header border-b border-gray-500 border-solid"
            >
              <h2 class="font-semibold text-5xl">
                <span class="streamer-name text-indigo-600"
                  >${streamerName}</span
                >
              </h2>
              <${SocialIconsList} handles=${socialHandles} />
              <div
                id="stream-summary"
                class="stream-summary my-4"
                dangerouslySetInnerHTML=${{ __html: summary }}
              ></div>
              <ul id="tag-list" class="tag-list flex flex-row my-4">
                ${tagList}
              </ul>
            </div>
          </div>
          <div
            id="extra-user-content"
            class="extra-user-content px-8"
            dangerouslySetInnerHTML=${{ __html: extraPageContent }}
          ></div>
        </section>

        <footer class="flex flex-row justify-start p-8 opacity-50 text-xs">
          <span class="mx-1 inline-block">
            <a href="${URL_OWNCAST}" target="_blank">About Owncast</a>
          </span>
          <span class="mx-1 inline-block">Version ${appVersion}</span>
        </footer>

        <${Chat}
          websocket=${websocket}
          username=${username}
          chatInputEnabled=${chatInputEnabled}
        />
      </div>
    `;
  }
}

