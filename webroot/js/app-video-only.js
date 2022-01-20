import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import VideoPoster from './components/video-poster.js';
import { OwncastPlayer } from './components/player.js';

import {
  addNewlines,
  makeLastOnlineString,
  pluralize,
  parseSecondsToDurationString,
} from './utils/helpers.js';
import {
  URL_CONFIG,
  URL_STATUS,
  URL_VIEWER_PING,
  TIMER_STATUS_UPDATE,
  TIMER_STREAM_DURATION_COUNTER,
  TEMP_IMAGE,
  MESSAGE_OFFLINE,
  MESSAGE_ONLINE,
} from './utils/constants.js';

export default class VideoOnly extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      configData: {},

      playerActive: false, // player object is active
      streamOnline: false, // stream is active/online

      isPlaying: false,

      //status
      streamStatusMessage: MESSAGE_OFFLINE,
      viewerCount: '',
      lastDisconnectTime: null,
    };

    // timers
    this.playerRestartTimer = null;
    this.offlineTimer = null;
    this.statusTimer = null;
    this.streamDurationTimer = null;

    this.handleOfflineMode = this.handleOfflineMode.bind(this);
    this.handleOnlineMode = this.handleOnlineMode.bind(this);
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
    clearInterval(this.streamDurationTimer);
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
    const { viewerCount, online, lastConnectTime, lastDisconnectTime } = status;

    if (status.online && !curStreamOnline) {
      // stream has just come online.
      this.handleOnlineMode();
    } else if (!status.online && curStreamOnline) {
      // stream has just flipped offline.
      this.handleOfflineMode();
    }
    this.setState({
      viewerCount,
      streamOnline: online,
      lastDisconnectTime,
      lastConnectTime,
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
    this.setState({
      streamOnline: false,
      streamStatusMessage: MESSAGE_OFFLINE,
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

  // play video!
  handleOnlineMode() {
    this.player.startPlayer();

    this.streamDurationTimer = setInterval(
      this.setCurrentStreamDuration,
      TIMER_STREAM_DURATION_COUNTER
    );

    this.setState({
      playerActive: true,
      streamOnline: true,
      streamStatusMessage: MESSAGE_ONLINE,
    });
  }

  handleNetworkingError(error) {
    console.error(`>>> App Error: ${error}`);
  }

  render(props, state) {
    const {
      configData,

      viewerCount,
      playerActive,
      streamOnline,
      streamStatusMessage,
      lastDisconnectTime,
      isPlaying,
    } = state;

    const { logo = TEMP_IMAGE, customStyles } = configData;

    let viewerCountMessage = '';
    if (streamOnline && viewerCount > 0) {
      viewerCountMessage = html`${viewerCount}
      ${pluralize(' viewer', viewerCount)}`;
    } else if (lastDisconnectTime) {
      viewerCountMessage = makeLastOnlineString(lastDisconnectTime);
    }

    const mainClass = playerActive ? 'online' : '';

    const poster = isPlaying
      ? null
      : html` <${VideoPoster} offlineImage=${logo} active=${streamOnline} /> `;
    return html`
      <main class=${mainClass}>
        <style>
          ${customStyles}
        </style>
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
          class="flex flex-row justify-between font-mono py-2 px-4 bg-gray-900 text-indigo-200 shadow-md border-b border-gray-100 border-solid"
        >
          <span class="text-xs">${streamStatusMessage}</span>
          <span id="stream-viewer-count" class="text-xs text-right"
            >${viewerCountMessage}</span
          >
        </section>
      </main>
    `;
  }
}
