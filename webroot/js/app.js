import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import { OwncastPlayer } from './components/player.js';
import SocialIconsList from './components/platform-logos-list.js';
import UsernameForm from './components/chat/username.js';
import VideoPoster from './components/video-poster.js';
import Chat from './components/chat/chat.js';
import Websocket from './utils/websocket.js';
import ExternalActionModal, {
  ExternalActionButton,
} from './components/external-action-modal.js';

import {
  addNewlines,
  checkUrlPathForDisplay,
  classNames,
  clearLocalStorage,
  debounce,
  generateUsername,
  getLocalStorage,
  getOrientation,
  hasTouchScreen,
  makeLastOnlineString,
  parseSecondsToDurationString,
  pluralize,
  ROUTE_RECORDINGS,
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
  URL_VIEWER_PING,
  WIDTH_SINGLE_COL,
} from './utils/constants.js';

export default class App extends Component {
  constructor(props, context) {
    super(props, context);

    const chatStorage = getLocalStorage(KEY_CHAT_DISPLAYED);
    this.hasTouchScreen = hasTouchScreen();
    this.windowBlurred = false;

    this.state = {
      websocket: new Websocket(),
      displayChat: chatStorage === null ? true : chatStorage,
      chatInputEnabled: false, // chat input box state
      username: getLocalStorage(KEY_USERNAME) || generateUsername(),
      touchKeyboardActive: false,

      configData: {
        loading: true,
      },
      extraPageContent: '',

      playerActive: false, // player object is active
      streamOnline: false, // stream is active/online
      isPlaying: false, // player is actively playing video

      // status
      streamStatusMessage: MESSAGE_OFFLINE,
      viewerCount: '',
      lastDisconnectTime: null,

      // dom
      windowWidth: window.innerWidth,
      windowHeight: window.innerHeight,
      orientation: getOrientation(this.hasTouchScreen),

      externalAction: null,

      // routing & tabbing
      section: '',
      sectionId: '',
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
    this.handleWindowBlur = this.handleWindowBlur.bind(this);
    this.handleWindowFocus = this.handleWindowFocus.bind(this);
    this.handleWindowResize = debounce(this.handleWindowResize.bind(this), 250);

    this.handleOfflineMode = this.handleOfflineMode.bind(this);
    this.handleOnlineMode = this.handleOnlineMode.bind(this);
    this.disableChatInput = this.disableChatInput.bind(this);
    this.setCurrentStreamDuration = this.setCurrentStreamDuration.bind(this);

    this.handleKeyPressed = this.handleKeyPressed.bind(this);
    this.displayExternalAction = this.displayExternalAction.bind(this);
    this.closeExternalActionModal = this.closeExternalActionModal.bind(this);

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
    if (!this.hasTouchScreen) {
      window.addEventListener('resize', this.handleWindowResize);
    }
    window.addEventListener('blur', this.handleWindowBlur);
    window.addEventListener('focus', this.handleWindowFocus);
    if (this.hasTouchScreen) {
      window.addEventListener('orientationchange', this.handleWindowResize);
    }
    window.addEventListener('keypress', this.handleKeyPressed);
    this.player = new OwncastPlayer();
    this.player.setupPlayerCallbacks({
      onReady: this.handlePlayerReady,
      onPlaying: this.handlePlayerPlaying,
      onEnded: this.handlePlayerEnded,
      onError: this.handlePlayerError,
    });
    this.player.init();

    // check routing
    console.log("==== did mount");
    this.getRoute();
  }

  componentWillUnmount() {
    // clear all the timers
    clearInterval(this.playerRestartTimer);
    clearInterval(this.offlineTimer);
    clearInterval(this.statusTimer);
    clearTimeout(this.disableChatTimer);
    clearInterval(this.streamDurationTimer);
    window.removeEventListener('resize', this.handleWindowResize);
    window.removeEventListener('blur', this.handleWindowBlur);
    window.removeEventListener('focus', this.handleWindowFocus);
    window.removeEventListener('keypress', this.handleKeyPressed);
    if (this.hasTouchScreen) {
      window.removeEventListener('orientationchange', this.handleWindowResize);
    }
  }

  getRoute() {
    const routeInfo = checkUrlPathForDisplay();
    this.setState({
      ...routeInfo,
    });
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

    // Ping the API to let them know we're an active viewer
    fetch(URL_VIEWER_PING).catch((error) => {
      this.handleOfflineMode();
      this.handleNetworkingError(`Viewer PING error: ${error}`);
    });
  }

  setConfigData(data = {}) {
    const { name, summary } = data;
    window.document.title = name;

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
      streamTitle,
      lastDisconnectTime,
    } = status;

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
      streamTitle,
      lastDisconnectTime,
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
      (Date.now() - new Date(this.state.lastDisconnectTime));
    const countdown = remainingChatTime < 0 ? 0 : remainingChatTime;
    this.disableChatTimer = setTimeout(this.disableChatInput, countdown);
    this.setState({
      streamOnline: false,
      streamStatusMessage: MESSAGE_OFFLINE,
    });

    if (this.player.vjsPlayer && this.player.vjsPlayer.paused()) {
      this.handlePlayerEnded();
    }

    if (this.windowBlurred) {
      document.title = ` 🔴 ${
        this.state.configData && this.state.configData.name
      }`;
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
      streamTitle: '',
      streamStatusMessage: MESSAGE_ONLINE,
    });

    if (this.windowBlurred) {
      document.title = ` 🟢 ${
        this.state.configData && this.state.configData.name
      }`;
    }
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

  handleWindowBlur() {
    this.windowBlurred = true;
  }

  handleWindowFocus() {
    this.windowBlurred = false;
    window.document.title = this.state.configData && this.state.configData.name;
  }

  handleSpaceBarPressed(e) {
    e.preventDefault();
    if (this.state.isPlaying) {
      this.setState({
        isPlaying: false,
      });
      this.player.vjsPlayer.pause();
    } else {
      this.setState({
        isPlaying: true,
      });
      this.player.vjsPlayer.play();
    }
  }

  handleMuteKeyPressed() {
    const muted = this.player.vjsPlayer.muted();
    const volume = this.player.vjsPlayer.volume();

    if (volume === 0) {
      this.player.vjsPlayer.volume(0.5);
      this.player.vjsPlayer.muted(false);
    } else {
      this.player.vjsPlayer.muted(!muted);
    }
  }

  handleFullScreenKeyPressed() {
    if (this.player.vjsPlayer.isFullscreen()) {
      this.player.vjsPlayer.exitFullscreen();
    } else {
      this.player.vjsPlayer.requestFullscreen();
    }
  }

  handleVolumeSet(factor) {
    this.player.vjsPlayer.volume(this.player.vjsPlayer.volume() + factor);
  }

  handleKeyPressed(e) {
    if (
      e.target !== document.getElementById('message-input') &&
      e.target !== document.getElementById('username-change-input') &&
      e.target !== document.getElementsByClassName('emoji-picker__search')[0] &&
      this.state.streamOnline
    ) {
      switch (e.code) {
        case 'MediaPlayPause':
        case 'KeyP':
        case 'Space':
          this.handleSpaceBarPressed(e);
          break;
        case 'KeyM':
          this.handleMuteKeyPressed(e);
          break;
        case 'KeyF':
          this.handleFullScreenKeyPressed(e);
          break;
        case 'KeyC':
          this.handleChatPanelToggle();
          break;
        case 'Digit9':
          this.handleVolumeSet(-0.1);
          break;
        case 'Digit0':
          this.handleVolumeSet(0.1);
      }
    }
  }

  displayExternalAction(action) {
    const { username } = this.state;
    if (!action) {
      return;
    }
    const { url: actionUrl, openExternally } = action || {};
    let url = new URL(actionUrl);
    // Append url and username to params so the link knows where we came from and who we are.
    url.searchParams.append('username', username);
    url.searchParams.append('instance', window.location);

    const fullUrl = url.toString();

    if (openExternally) {
      var win = window.open(fullUrl, '_blank');
      win.focus();
      return;
    }
    this.setState({
      externalAction: {
        ...action,
        url: fullUrl,
      },
    });
  }
  closeExternalActionModal() {
    this.setState({
      externalAction: null,
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
      streamTitle,
      touchKeyboardActive,
      username,
      viewerCount,
      websocket,
      windowHeight,
      windowWidth,
      externalAction,
      lastDisconnectTime,
      section,
      sectionId,
    } = state;

    const {
      version: appVersion,
      logo = TEMP_IMAGE,
      socialHandles = [],
      summary,
      tags = [],
      name,
      extraPageContent,
      chatDisabled,
      externalActions,
      customStyles,
    } = configData;

    const bgUserLogo = { backgroundImage: `url(${logo})` };

    const tagList = tags !== null && tags.length > 0 && tags.join(' #');

    let viewerCountMessage = '';
    if (streamOnline && viewerCount > 0) {
      viewerCountMessage = html`${viewerCount}
      ${pluralize(' viewer', viewerCount)}`;
    } else if (lastDisconnectTime) {
      viewerCountMessage = makeLastOnlineString(lastDisconnectTime);
    }

    const mainClass = playerActive ? 'online' : '';
    const isPortrait =
      this.hasTouchScreen && orientation === ORIENTATION_PORTRAIT;
    const shortHeight = windowHeight <= HEIGHT_SHORT_WIDE && !isPortrait;
    const singleColMode = windowWidth <= WIDTH_SINGLE_COL && !shortHeight;

    const noVideoContent = !playerActive || (section === ROUTE_RECORDINGS && sectionId !== '');
    const shouldDisplayChat = displayChat && !chatDisabled && !noVideoContent;
    const usernameStyle = chatDisabled ? 'none' : 'flex';

    const extraAppClasses = classNames({
      'config-loading': configData.loading,

      chat: shouldDisplayChat,
      'no-chat': !shouldDisplayChat,
      'no-video': noVideoContent,
      'single-col': singleColMode,
      'bg-gray-800': singleColMode && shouldDisplayChat,
      'short-wide': shortHeight && windowWidth > WIDTH_SINGLE_COL,
      'touch-screen': this.hasTouchScreen,
      'touch-keyboard-active': touchKeyboardActive,
    });

    const poster = isPlaying
      ? null
      : html` <${VideoPoster} offlineImage=${logo} active=${streamOnline} /> `;

    // modal buttons
    const externalActionButtons =
      externalActions &&
      html`<div
        id="external-actions-container"
        class="flex flex-row justify-end"
      >
        ${externalActions.map(
          function (action) {
            return html`<${ExternalActionButton}
              onClick=${this.displayExternalAction}
              action=${action}
            />`;
          }.bind(this)
        )}
      </div>`;

    // modal component
    const externalActionModal =
      externalAction &&
      html`<${ExternalActionModal}
        action=${externalAction}
        onClose=${this.closeExternalActionModal}
      />`;

    return html`
      <div
        id="app-container"
        class="flex w-full flex-col justify-start relative ${extraAppClasses}"
      >
        <style>
          ${customStyles}
        </style>

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
                <img
                  class="logo visually-hidden"
                  src=${OWNCAST_LOGO_LOCAL}
                  alt="owncast logo"
                />
              </span>
              <span class="instance-title overflow-hidden truncate"
                >${streamOnline && streamTitle ? streamTitle : name}</span
              >
            </h1>
            <div
              id="user-options-container"
              class="flex flex-row justify-end items-center flex-no-wrap"
              style=${{ display: usernameStyle }}
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
            class="flex text-center flex-row justify-between font-mono py-2 px-4 bg-gray-900 text-indigo-200 shadow-md border-b border-gray-100 border-solid"
          >
            <span class="text-xs">${streamStatusMessage}</span>
            <span id="stream-viewer-count" class="text-xs text-right"
              >${viewerCountMessage}</span
            >
          </section>
        </main>

        <section id="user-content" aria-label="User information" class="p-2">

          ${externalActionButtons && html`${externalActionButtons}`}


          <div class="user-content flex flex-row p-8">
            <div class="flex flex-col align-center justify-start mr-8">
              <div
                class="user-image rounded-full bg-white p-4 bg-no-repeat bg-center"
                style=${bgUserLogo}
              >
                <img class="logo visually-hidden" alt="" src=${logo} />
              </div>
              <${SocialIconsList} handles=${socialHandles} />
            </div>

            <div class="user-content-header">
              <h2 class="font-semibold text-5xl">
                <span class="streamer-name text-indigo-600">${name}</span>
              </h2>
              <h3 class="font-semibold text-3xl">
                ${streamOnline && streamTitle}
              </h3>
              <div
                id="stream-summary"
                class="stream-summary my-4"
                dangerouslySetInnerHTML=${{ __html: summary }}
              ></div>
              <div id="tag-list" class="tag-list text-gray-600 mb-3">
                ${tagList && `#${tagList}`}
              </div>
            </div>
          </div>


          <!-- tab bar -->
          <div class="mx-8 border-b border-gray-500 border-solid" id="tab-bar">
            tab bar
          </div>

          <div
            id="extra-user-content"
            class="extra-user-content px-8"
            dangerouslySetInnerHTML=${{ __html: extraPageContent }}
          ></div>
        </section>

        <footer class="flex flex-row justify-start p-8 opacity-50 text-xs">
          <span class="mx-1 inline-block">
            <a href="${URL_OWNCAST}" rel="noopener noreferrer" target="_blank"
              >${appVersion}</a
            >
          </span>
        </footer>

        <${Chat}
          websocket=${websocket}
          username=${username}
          chatInputEnabled=${chatInputEnabled && !chatDisabled}
          instanceTitle=${name}
        />
        ${externalActionModal}
      </div>
    `;
  }
}
