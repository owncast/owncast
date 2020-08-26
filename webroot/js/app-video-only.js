import { h, Component } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
const html = htm.bind(h);

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
      sessionMaxViewerCount,
      overallMaxViewerCount,
      playerActive,
      streamOnline,
      streamStatusMessage,
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

    const bgLogoLarge = { backgroundImage: `url(${largeLogo})` };

    const mainClass = playerActive ? 'online' : '';
    return (
      html`
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

          <section id="stream-info" aria-label="Stream status" class="flex text-center flex-row justify-between items-center font-mono py-2 px-8 bg-gray-900 text-indigo-200">
            <span>${streamStatusMessage}</span>
            <span>${viewerCount} ${pluralize('viewer', viewerCount)}.</span>
            <span>Max ${pluralize('viewer', sessionMaxViewerCount)}.</span>
            <span>${overallMaxViewerCount} overall.</span>
          </section>
        </main>
    `);
  }
}
