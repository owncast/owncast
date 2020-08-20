import { h, Component, Fragment } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
const html = htm.bind(h);


import UsernameForm from './chat/username.js';
import Chat from './chat/chat.js';
import Websocket from './websocket.js';

import { getLocalStorage, generateAvatar, generateUsername, URL_OWNCAST, URL_CONFIG, URL_STATUS, addNewlines } from './utils.js';
import { KEY_USERNAME, KEY_AVATAR,  } from './utils/chat.js';

export default class App extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      websocket: new Websocket(),
      chatEnabled: true, // always true for standalone chat
      username: getLocalStorage(KEY_USERNAME) || generateUsername(),
      userAvatarImage: getLocalStorage(KEY_AVATAR) || generateAvatar(`${this.username}${Date.now()}`),

      streamStatus: null,
      player: null,
      configData: {},
    };

    // timers
    this.playerRestartTimer = null;
    this.offlineTimer = null;
    this.statusTimer = null;
    this.disableChatTimer = null;
    this.streamDurationTimer = null;

    this.handleUsernameChange = this.handleUsernameChange.bind(this);
    this.getConfig = this.getConfig.bind(this);
    this.getStreamStatus = this.getStreamStatus.bind(this);
    this.getExtraUserContent = this.getExtraUserContent.bind(this);

  }

  componentDidMount() {
    this.getConfig();

    // DO LATER..
    // this.player = new OwncastPlayer();
    // this.player.setupPlayerCallbacks({
    //   onReady: this.handlePlayerReady,
    //   onPlaying: this.handlePlayerPlaying,
    //   onEnded: this.handlePlayerEnded,
    //   onError: this.handlePlayerError,
    // });
    // this.player.init();
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
        const descriptionHTML = new showdown.Converter().makeHtml(text);
        this.vueApp.extraUserContent = descriptionHTML;
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
    if (!status) {
      return;
    }
    // update UI
    this.vueApp.viewerCount = status.viewerCount;
    this.vueApp.sessionMaxViewerCount = status.sessionMaxViewerCount;
    this.vueApp.overallMaxViewerCount = status.overallMaxViewerCount;

    this.lastDisconnectTime = status.lastDisconnectTime;

    if (!this.streamStatus) {
      // display offline mode the first time we get status, and it's offline.
      if (!status.online) {
        this.handleOfflineMode();
      } else {
        this.handleOnlineMode();
      }
    } else {
      if (status.online && !this.streamStatus.online) {
        // stream has just come online.
        this.handleOnlineMode();
      } else if (!status.online && this.streamStatus.online) {
        // stream has just flipped offline.
        this.handleOfflineMode();
      }
    }

    // keep a local copy
    this.streamStatus = status;

    if (status.online) {
      // only do this if video is paused, so no unnecessary img fetches
      if (this.player.vjsPlayer && this.player.vjsPlayer.paused()) {
        this.player.setPoster();
      }
    }
  }

  // stop status timer and disable chat after some time.
  handleOfflineMode() {
    this.vueApp.isOnline = false;
    clearInterval(this.streamDurationTimer);
    this.vueApp.streamStatus = MESSAGE_OFFLINE;
    if (this.streamStatus) {
      const remainingChatTime = TIMER_DISABLE_CHAT_AFTER_OFFLINE - (Date.now() - new Date(this.lastDisconnectTime));
      const countdown = (remainingChatTime < 0) ? 0 : remainingChatTime;
      this.disableChatTimer = setTimeout(this.messagingInterface.disableChat, countdown);
    }
  }

  // play video!
  handleOnlineMode() {
    this.vueApp.playerOn = true;
    this.vueApp.isOnline = true;
    this.vueApp.streamStatus = MESSAGE_ONLINE;

    this.player.startPlayer();
    clearTimeout(this.disableChatTimer);
    this.disableChatTimer = null;
    this.messagingInterface.enableChat();

    this.streamDurationTimer =
      setInterval(this.setCurrentStreamDuration, TIMER_STREAM_DURATION_COUNTER);
  }


  handleUsernameChange(newName, newAvatar) {
    this.setState({
      username: newName,
      userAvatarImage: newAvatar,
    });
  }

  handleChatToggle() {
    const { chatEnabled: curChatEnabled } = this.state;
    this.setState({
      chatEnabled: !curChatEnabled,
    });
  }

  handleNetworkingError(error) {
    console.log(`>>> App Error: ${error}`);
  }

  render(props, state) {
    const { username, userAvatarImage, websocket, configData } = state;
    const {
      version: appVersion,
      logo = {},
      socialHandles,
      name: streamnerName,
      summary,
      tags,
      title,
    } = configData;
    const { small: smallLogo, large: largeLogo } = logo;

    const bgLogo = { backgroundImage: `url(${smallLogo})` };
    const bgLogoLarge = { backgroundImage: `url(${largeLogo})` };

    // not needed for standalone, just messages only. remove later.
    return (
      html`
        <div id="app-container" class="flex chat">
          <div id="top-content">
            <header class="flex border-b border-gray-900 border-solid shadow-md">
              <h1 v-cloak class="flex text-gray-400">
                <span
                  id="logo-container"
                  class="rounded-full bg-white px-1 py-1"
                  style=${bgLogo}
                >
                  <img class="logo visually-hidden" src=${smallLogo}>
                </span>
                <span class="instance-title">${title}</span>
              </h1>

              <${UsernameForm}
                username=${username}
                userAvatarImage=${userAvatarImage}
                handleUsernameChange=${this.handleUsernameChange}
                handleChatToggle=${this.handleChatToggle}
              />

            </header>
          </div>

          <main>
            <div
              id="video-container"
              class="flex owncast-video-container bg-black"
              style=${bgLogoLarge}
            >
              <video
                class="video-js vjs-big-play-centered"
                id="video"
                preload="auto"
                controls
                playsinline
              >
              </video>
            </div>


            <section id="stream-info" aria-label="Stream status" class="flex font-mono bg-gray-900 text-indigo-200 shadow-md border-b border-gray-100 border-solid">
              <span>{{ streamStatus }}</span>
              <span v-if="isOnline">{{ viewerCount }} {{ 'viewer' | plural(viewerCount) }}.</span>
              <span v-if="isOnline">Max {{ sessionMaxViewerCount }} {{ 'viewer' | plural(sessionMaxViewerCount) }}.</span>
              <span v-if="isOnline">{{ overallMaxViewerCount }} overall.</span>
            </section>
          </main>

          <section id="user-content" aria-label="User information">
            <user-details
              v-bind:logo="logo"
              v-bind:platforms="socialHandles"
              v-bind:summary="summary"
              v-bind:tags="tags"
            >{{streamerName}}</user-details>

            <div v-html="extraUserContent" class="extra-user-content">{{extraUserContent}}</div>
          </section>

          <footer class="flex">
            <span>
              <a href="${URL_OWNCAST}" target="_blank">About Owncast</a>
            </span>
            <span>Version ${appVersion}</span>
          </footer>
        </div>

        <${Chat}
          websocket=${websocket}
          username=${username}
          userAvatarImage=${userAvatarImage}
          chatEnabled
        />

      </div>
    `);
  }
}

