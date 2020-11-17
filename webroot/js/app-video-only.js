import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import VideoPoster from './components/video-poster.js';
import { OwncastPlayer } from './components/player.js';

import {
  addNewlines,
  pluralize,
} from './utils/helpers.js';
import {
  URL_CONFIG,
  URL_STATUS,
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
      streamOnline: false,  // stream is active/online

      isPlaying: false,

      //status
      streamStatusMessage: MESSAGE_OFFLINE,
      viewerCount: '',
    };

    // timers
    this.playerRestartTimer = null;
    this.offlineTimer = null;
    this.statusTimer = null;
    this.streamDurationTimer = null;

    this.handleOfflineMode = this.handleOfflineMode.bind(this);
    this.handleOnlineMode = this.handleOnlineMode.bind(this);

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
    this.setState({
      streamOnline: false,
      streamStatusMessage: MESSAGE_OFFLINE,
    });
  }

  // play video!
  handleOnlineMode() {
    this.player.startPlayer();

    this.streamDurationTimer =
      setInterval(this.setCurrentStreamDuration, TIMER_STREAM_DURATION_COUNTER);

    this.setState({
      playerActive: true,
      streamOnline: true,
      streamStatusMessage: MESSAGE_ONLINE,
    });
  }

  handleNetworkingError(error) {
    console.log(`>>> App Error: ${error}`);
  }

  render(props, state) {
    const {
      configData,

      viewerCount,
      playerActive,
      streamOnline,
      streamStatusMessage,
      isPlaying,
    } = state;

    const { logo = TEMP_IMAGE } = configData;
    const streamInfoClass = streamOnline ? 'online' : ''; // need?

    const mainClass = playerActive ? 'online' : '';

    const poster = isPlaying ? null : html`
      <${VideoPoster} offlineImage=${logo} active=${streamOnline} />
    `;
    return (
      html`
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
    `);
  }
}
