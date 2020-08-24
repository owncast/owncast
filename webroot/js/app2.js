import { h, Component } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
const html = htm.bind(h);


import SocialIcon from './social.js';
import UsernameForm from './chat/username.js';
import Chat from './chat/chat.js';
import Websocket from './websocket.js';
import { OwncastPlayer } from './player.js';

import {
  getLocalStorage,
  setLocalStorage,
  clearLocalStorage,
  generateAvatar,
  generateUsername,
  URL_OWNCAST,
  URL_CONFIG,
  URL_STATUS,
  addNewlines,
  pluralize,
  TIMER_STATUS_UPDATE,
  TIMER_DISABLE_CHAT_AFTER_OFFLINE,
  TIMER_STREAM_DURATION_COUNTER,
  TEMP_IMAGE,
  MESSAGE_OFFLINE,
  MESSAGE_ONLINE,
} from './utils.js';
import { KEY_USERNAME, KEY_AVATAR, KEY_CHAT_DISPLAYED } from './utils/chat.js';

export default class App extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      websocket: new Websocket(),
      displayChat: getLocalStorage(KEY_CHAT_DISPLAYED), // chat panel state
      chatEnabled: false, // chat input box state
      username: getLocalStorage(KEY_USERNAME) || generateUsername(),
      userAvatarImage: getLocalStorage(KEY_AVATAR) || generateAvatar(`${this.username}${Date.now()}`),

      configData: {},
      extraUserContent: '',

      playerActive: false, // player object is active
      streamOnline: false,  // stream is active/online

      //status
      streamStatusMessage: MESSAGE_OFFLINE,
      viewerCount: '',
      sessionMaxViewerCount: '',
      overallMaxViewerCount: '',
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

    this.handleOfflineMode = this.handleOfflineMode.bind(this);
    this.handleOnlineMode = this.handleOnlineMode.bind(this);
    this.disableChatInput = this.disableChatInput.bind(this);

    // player events
    this.handlePlayerReady = this.handlePlayerReady.bind(this);
    this.handlePlayerPlaying = this.handlePlayerPlaying.bind(this);
    this.handlePlayerEnded = this.handlePlayerEnded.bind(this);
    this.handlePlayerError = this.handlePlayerError.bind(this);

    // fetch events
    this.getConfig = this.getConfig.bind(this);
    this.getStreamStatus = this.getStreamStatus.bind(this);
    this.getExtraUserContent = this.getExtraUserContent.bind(this);


  }

  componentDidMount() {
    this.getConfig();

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
  }

  // fetch /config data
  getConfig() {
    fetch(URL_CONFIG)
    .then(response => {
      if (!response.ok) {
        throw new Error(`Network response was not ok ${response.ok}`);
      }
      return response.json();
    })
    .then(json => {
      this.setConfigData(json);
    })
    .catch(error => {
      this.handleNetworkingError(`Fetch config: ${error}`);
    });
  }

  // fetch stream status
  getStreamStatus() {
    fetch(URL_STATUS)
      .then(response => {
        if (!response.ok) {
          throw new Error(`Network response was not ok ${response.ok}`);
        }
        return response.json();
      })
      .then(json => {
        this.updateStreamStatus(json);
      })
      .catch(error => {
        this.handleOfflineMode();
        this.handleNetworkingError(`Stream status: ${error}`);
      });
  }

  // fetch content.md
  getExtraUserContent(path) {
    fetch(path)
      .then(response => {
        if (!response.ok) {
          throw new Error(`Network response was not ok ${response.ok}`);
        }
        return response.text();
      })
      .then(text => {
        this.setState({
          extraUserContent: new showdown.Converter().makeHtml(text),
        });
      })
      .catch(error => {
        this.handleNetworkingError(`Fetch extra content: ${error}`);
      });
  }


  setConfigData(data = {}) {
    const { title, extraUserInfoFileName, summary } = data;

    window.document.title = title;
    if (extraUserInfoFileName) {
      this.getExtraUserContent(extraUserInfoFileName);
    }

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
      sessionMaxViewerCount,
      overallMaxViewerCount,
      online,
    } = status;

    this.lastDisconnectTime = status.lastDisconnectTime;

    if (status.online && !curStreamOnline) {
      // stream has just come online.
      this.handleOnlineMode();
    } else if (!status.online && curStreamOnline) {
      // stream has just flipped offline.
      this.handleOfflineMode();
    }
    if (status.online) {
      // only do this if video is paused, so no unnecessary img fetches
      if (this.player.vjsPlayer && this.player.vjsPlayer.paused()) {
        this.player.setPoster();
      }
    }
    this.setState({
      viewerCount,
      sessionMaxViewerCount,
      overallMaxViewerCount,
      streamOnline: online,
    });
  }

  // when videojs player is ready, start polling for stream
  handlePlayerReady() {
    this.getStreamStatus();
    this.statusTimer = setInterval(this.getStreamStatus, TIMER_STATUS_UPDATE);
  }

  handlePlayerPlaying() {
    // do something?
  }


  // likely called some time after stream status has gone offline.
  // basically hide video and show underlying "poster"
  handlePlayerEnded() {
    this.setState({
      playerActive: false,
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
    const remainingChatTime = TIMER_DISABLE_CHAT_AFTER_OFFLINE - (Date.now() - new Date(this.lastDisconnectTime));
    const countdown = (remainingChatTime < 0) ? 0 : remainingChatTime;
    this.disableChatTimer = setTimeout(this.disableChatInput, countdown);
    this.setState({
      streamOnline: false,
      streamStatusMessage: MESSAGE_OFFLINE,
    });
  }

  // play video!
  handleOnlineMode() {
    this.player.startPlayer();
    clearTimeout(this.disableChatTimer);
    this.disableChatTimer = null;

    this.streamDurationTimer =
      setInterval(this.setCurrentStreamDuration, TIMER_STREAM_DURATION_COUNTER);

    this.setState({
      playerActive: true,
      streamOnline: true,
      chatEnabled: true,
      streamStatusMessage: MESSAGE_ONLINE,
    });
  }


  handleUsernameChange(newName, newAvatar) {
    this.setState({
      username: newName,
      userAvatarImage: newAvatar,
    });
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
      chatEnabled: false,
    });
  }

  handleNetworkingError(error) {
    console.log(`>>> App Error: ${error}`);
  }

  render(props, state) {
    const {
      username,
      userAvatarImage,
      websocket,
      configData,
      extraUserContent,
      displayChat,
      viewerCount,
      sessionMaxViewerCount,
      overallMaxViewerCount,
      playerActive,
      streamOnline,
      streamStatusMessage,
      chatEnabled,
    } = state;

    const {
      version: appVersion,
      logo = {},
      socialHandles = [],
      name: streamerName,
      summary,
      tags = [],
      title,
    } = configData;
    const { small: smallLogo = TEMP_IMAGE, large: largeLogo = TEMP_IMAGE } = logo;

    const bgLogo = { backgroundImage: `url(${smallLogo})` };
    const bgLogoLarge = { backgroundImage: `url(${largeLogo})` };

    const tagList = !tags.length ?
      null :
      tags.map((tag, index) => html`
        <li key="tag${index}" class="tag rounded-sm text-gray-100 bg-gray-700 text-xs uppercase mr-3 p-2 whitespace-no-wrap">${tag}</li>
      `);

    const socialIconsList =
      !socialHandles.length ?
      null :
      socialHandles.map((item, index) => html`
        <li key="social${index}">
          <${SocialIcon} platform=${item.platform} url=${item.url} />
        </li>
      `);


    const chatClass = displayChat ? 'chat' : 'no-chat';
    const mainClass = playerActive ? 'online' : '';
    const streamInfoClass = streamOnline ? 'online' : '';
    return (
      html`
        <div id="app-container" class="flex w-full flex-col justify-start relative ${chatClass}">
          <div id="top-content">
            <header class="flex border-b border-gray-900 border-solid shadow-md fixed z-10 w-full top-0	left-0 flex flex-row justify-between flex-no-wrap">
              <h1 class="flex flex-row items-center justify-start p-2 uppercase text-gray-400 text-xl	font-thin tracking-wider overflow-hidden whitespace-no-wrap">
                <span
                  id="logo-container"
                  class="inline-block	rounded-full bg-white w-8 min-w-8 min-h-8 h-8 p-1 mr-2 bg-no-repeat bg-center"
                  style=${bgLogo}
                >
                  <img class="logo visually-hidden" src=${smallLogo} alt=""/>
                </span>
                <span class="instance-title overflow-hidden truncate">${title}</span>
              </h1>
              <div id="user-options-container" class="flex flex-row justify-end items-center flex-no-wrap">
                <${UsernameForm}
                  username=${username}
                  userAvatarImage=${userAvatarImage}
                  handleUsernameChange=${this.handleUsernameChange}
                />
                <button type="button" id="chat-toggle" onClick=${this.handleChatPanelToggle} class="flex cursor-pointer text-center justify-center items-center min-w-12 h-full bg-gray-800 hover:bg-gray-700">💬</button>
              </div>
            </header>
          </div>

          <main class=${mainClass}>
            <div
              id="video-container"
              class="flex owncast-video-container bg-black w-full bg-center bg-no-repeat flex flex-col items-center justify-start"
              style=${bgLogoLarge}
            >
              <video
                class="video-js vjs-big-play-centered display-block w-full h-full"
                id="video"
                preload="auto"
                controls
                playsinline
              ></video>
            </div>


            <section id="stream-info" aria-label="Stream status" class="flex text-center flex-row justify-between font-mono py-2 px-8 bg-gray-900 text-indigo-200 shadow-md border-b border-gray-100 border-solid ${streamInfoClass}">
              <span>${streamStatusMessage}</span>
              <span>${viewerCount} ${pluralize('viewer', viewerCount)}.</span>
              <span>Max ${pluralize('viewer', sessionMaxViewerCount)}.</span>
              <span>${overallMaxViewerCount} overall.</span>
            </section>
          </main>

          <section id="user-content" aria-label="User information" class="p-8">
            <div class="user-content flex flex-row p-8">
              <div
                class="user-image rounded-full bg-white p-4 mr-8 bg-no-repeat bg-center"
                style=${bgLogoLarge}
              >
                <img
                  class="logo visually-hidden"
                  alt="Logo"
                  src=${largeLogo}/>
              </div>
              <div class="user-content-header border-b border-gray-500 border-solid">
                <h2 class="font-semibold text-5xl">
                  About <span class="streamer-name text-indigo-600">${streamerName}</span>
                </h2>
                <ul id="social-list" class="social-list flex flex-row items-center justify-start flex-wrap">
                  <span class="follow-label text-xs	font-bold mr-2 uppercase">Follow me: </span>
                  ${socialIconsList}
                </ul>
                <div id="stream-summary" class="stream-summary my-4" dangerouslySetInnerHTML=${{ __html: summary }}></div>
                <ul id="tag-list" class="tag-list flex flex-row my-4">
                  ${tagList}
                </ul>
              </div>
            </div>
            <div
              id="extra-user-content"
              class="extra-user-content px-8"
              dangerouslySetInnerHTML=${{ __html: extraUserContent }}
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
          userAvatarImage=${userAvatarImage}
          chatEnabled=${chatEnabled}
        />

      </div>
    `);
  }
}

