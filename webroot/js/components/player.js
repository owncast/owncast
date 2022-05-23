// https://docs.videojs.com/player

import videojs from '/js/web_modules/videojs/dist/video.min.js';
import { getLocalStorage, setLocalStorage } from '../utils/helpers.js';
import { PLAYER_VOLUME, URL_STREAM } from '../utils/constants.js';
import PlaybackMetrics from '../metrics/playback.js';
import LatencyCompensator from './latencyCompensator.js';

const VIDEO_ID = 'video';
const LATENCY_COMPENSATION_ENABLED = 'latencyCompensatorEnabled';

// Video setup
const VIDEO_SRC = {
  src: URL_STREAM,
  type: 'application/x-mpegURL',
};
const VIDEO_OPTIONS = {
  autoplay: false,
  liveui: true,
  preload: 'auto',
  controlBar: {
    progressControl: {
      seekBar: false,
    },
  },
  html5: {
    vhs: {
      // used to select the lowest bitrate playlist initially. This helps to decrease playback start time. This setting is false by default.
      enableLowInitialPlaylist: true,
      experimentalBufferBasedABR: true,
      useNetworkInformationApi: true,
      maxPlaylistRetries: 30,
    },
  },
  liveTracker: {
    trackingThreshold: 0,
    liveTolerance: 15,
  },
  sources: [VIDEO_SRC],
};

export const POSTER_DEFAULT = `/img/logo.png`;
export const POSTER_THUMB = `/thumbnail.jpg`;

class OwncastPlayer {
  constructor() {
    window.VIDEOJS_NO_DYNAMIC_STYLE = true; // style override

    this.vjsPlayer = null;
    this.latencyCompensator = null;
    this.playbackMetrics = null;

    this.appPlayerReadyCallback = null;
    this.appPlayerPlayingCallback = null;
    this.appPlayerEndedCallback = null;

    this.hasStartedPlayback = false;
    this.latencyCompensatorEnabled = false;

    this.clockSkewMs = 0;

    // bind all the things because safari
    this.startPlayer = this.startPlayer.bind(this);
    this.handleReady = this.handleReady.bind(this);
    this.handlePlaying = this.handlePlaying.bind(this);
    this.handleVolume = this.handleVolume.bind(this);
    this.handleEnded = this.handleEnded.bind(this);
    this.handleError = this.handleError.bind(this);
    this.addQualitySelector = this.addVideoSettingsMenu.bind(this);
    this.addQualitySelector = this.addVideoSettingsMenu.bind(this);
    this.toggleLatencyCompensator = this.toggleLatencyCompensator.bind(this);
    this.startLatencyCompensator = this.startLatencyCompensator.bind(this);
    this.stopLatencyCompensator = this.stopLatencyCompensator.bind(this);
    this.qualitySelectionMenu = null;
    this.latencyCompensatorToggleButton = null;
  }

  init() {
    this.addAirplay();
    this.addVideoSettingsMenu();

    // Add a cachebuster param to playlist URLs.
    videojs.Vhs.xhr.beforeRequest = (options) => {
      if (options.uri.match('m3u8')) {
        const cachebuster = Math.random().toString(16).substr(2, 8);
        options.uri = `${options.uri}?cachebust=${cachebuster}`;
      }

      return options;
    };

    this.vjsPlayer = videojs(VIDEO_ID, VIDEO_OPTIONS);

    this.vjsPlayer.ready(this.handleReady);
  }

  setClockSkew(skewMs) {
    this.clockSkewMs = skewMs;

    if (this.playbackMetrics) {
      this.playbackMetrics.setClockSkew(skewMs);
    }

    if (this.latencyCompensator) {
      this.latencyCompensator.setClockSkew(skewMs);
    }
  }

  setupPlayerCallbacks(callbacks) {
    const { onReady, onPlaying, onEnded, onError } = callbacks;

    this.appPlayerReadyCallback = onReady;
    this.appPlayerPlayingCallback = onPlaying;
    this.appPlayerEndedCallback = onEnded;
    this.appPlayerErrorCallback = onError;
  }

  // play
  startPlayer() {
    this.log('Start playing');
    const source = { ...VIDEO_SRC };

    try {
      this.vjsPlayer.volume(getLocalStorage(PLAYER_VOLUME) || 1);
    } catch (err) {
      console.warn(err);
    }
    this.vjsPlayer.src(source);
  }

  setupPlaybackMetrics() {
    this.playbackMetrics = new PlaybackMetrics(this.vjsPlayer, videojs);
    this.playbackMetrics.setClockSkew(this.clockSkewMs);
  }

  setupLatencyCompensator() {
    const tech = this.vjsPlayer.tech({ IWillNotUseThisInPlugins: true });

    // VHS is required.
    if (!tech || !tech.vhs) {
      return;
    }

    const latencyCompensatorEnabledSaved = getLocalStorage(
      LATENCY_COMPENSATION_ENABLED
    );

    if (latencyCompensatorEnabledSaved === 'true' && tech && tech.vhs) {
      this.startLatencyCompensator();
    } else {
      this.stopLatencyCompensator();
    }
  }

  startLatencyCompensator() {
    this.latencyCompensator = new LatencyCompensator(this.vjsPlayer);
    this.playbackMetrics.setClockSkew(this.clockSkewMs);
    this.latencyCompensator.enable();
    this.latencyCompensatorEnabled = true;
    this.setLatencyCompensatorItemTitle('disable minimized latency');
  }

  stopLatencyCompensator() {
    if (this.latencyCompensator) {
      this.latencyCompensator.disable();
    }
    this.LatencyCompensator = null;
    this.latencyCompensatorEnabled = false;
    this.setLatencyCompensatorItemTitle(
      '<span style="font-size: 0.8em">enable minimized latency (experimental)</span>'
    );
  }

  handleReady() {
    console.log('handleReady');
    this.vjsPlayer.on('error', this.handleError);
    this.vjsPlayer.on('playing', this.handlePlaying);
    this.vjsPlayer.on('volumechange', this.handleVolume);
    this.vjsPlayer.on('ended', this.handleEnded);

    this.vjsPlayer.on('loadeddata', () => {
      this.setupPlaybackMetrics();
      this.setupLatencyCompensator();
    });

    if (this.appPlayerReadyCallback) {
      // start polling
      this.appPlayerReadyCallback();
    }

    this.vjsPlayer.log.level('debug');
  }

  handleVolume() {
    setLocalStorage(
      PLAYER_VOLUME,
      this.vjsPlayer.muted() ? 0 : this.vjsPlayer.volume()
    );
  }

  handlePlaying() {
    if (this.appPlayerPlayingCallback) {
      // start polling
      this.appPlayerPlayingCallback();
    }

    if (this.latencyCompensator && !this.hasStartedPlayback) {
      this.latencyCompensator.enable();
    }

    this.hasStartedPlayback = true;
  }

  handleEnded() {
    if (this.appPlayerEndedCallback) {
      this.appPlayerEndedCallback();
    }

    this.stopLatencyCompensator();
  }

  handleError(e) {
    this.log(`on Error: ${JSON.stringify(e)}`);
    if (this.appPlayerEndedCallback) {
      this.appPlayerEndedCallback();
    }
  }

  toggleLatencyCompensator() {
    if (this.latencyCompensatorEnabled) {
      this.stopLatencyCompensator();
      setLocalStorage(LATENCY_COMPENSATION_ENABLED, false);
    } else {
      this.startLatencyCompensator();
      setLocalStorage(LATENCY_COMPENSATION_ENABLED, true);
    }
  }

  setLatencyCompensatorItemTitle(title) {
    const item = document.querySelector(
      '.latency-toggle-item > .vjs-menu-item-text'
    );
    if (!item) {
      return;
    }

    item.innerHTML = title;
  }

  log(message) {
    // console.log(`>>> Player: ${message}`);
  }

  async addVideoSettingsMenu() {
    if (this.qualityMenuButton) {
      player.controlBar.removeChild(this.qualityMenuButton);
    }

    videojs.hookOnce(
      'setup',
      async function (player) {
        var qualities = [];

        try {
          const response = await fetch('/api/video/variants');
          qualities = await response.json();
        } catch (e) {
          console.error(e);
        }

        var MenuItem = videojs.getComponent('MenuItem');
        var MenuButtonClass = videojs.getComponent('MenuButton');

        const lowLatencyItem = new MenuItem(player, {
          selectable: true,
        });
        lowLatencyItem.setAttribute('class', 'latency-toggle-item');
        lowLatencyItem.on('click', () => {
          this.toggleLatencyCompensator();
        });
        this.latencyCompensatorToggleButton = lowLatencyItem;

        const separator = new MenuSeparator(player, {
          selectable: false,
        });

        var MenuButton = videojs.extend(MenuButtonClass, {
          // The `init()` method will also work for constructor logic here, but it is
          // deprecated. If you provide an `init()` method, it will override the
          // `constructor()` method!
          constructor: function () {
            MenuButtonClass.call(this, player);
          },

          createItems: function () {
            const tech = this.player_.tech({ IWillNotUseThisInPlugins: true });

            const defaultAutoItem = new MenuItem(player, {
              selectable: true,
              label: 'Auto',
            });

            const items = qualities.map(function (item) {
              var newMenuItem = new MenuItem(player, {
                selectable: true,
                label: item.name,
              });

              // Quality selected
              newMenuItem.on('click', function () {
                // If for some reason tech doesn't exist, then don't do anything
                if (!tech) {
                  console.warn('Invalid attempt to access null player tech');
                  return;
                }
                // Only enable this single, selected representation.
                tech.vhs.representations().forEach(function (rep, index) {
                  rep.enabled(index === item.index);
                });
                newMenuItem.selected(false);
              });

              return newMenuItem;
            });

            defaultAutoItem.on('click', function () {
              // Re-enable all representations.
              tech.vhs.representations().forEach(function (rep, index) {
                rep.enabled(true);
              });
              defaultAutoItem.selected(false);
            });

            const supportsLatencyCompensator = !!tech && !!tech.vhs;

            // Only show the quality selector if there is more than one option.
            if (qualities.length < 2 && supportsLatencyCompensator) {
              return [lowLatencyItem];
            } else if (qualities.length > 1 && supportsLatencyCompensator) {
              return [defaultAutoItem, ...items, separator, lowLatencyItem];
            } else if (!supportsLatencyCompensator && qualities.length == 1) {
              return [];
            }

            return [defaultAutoItem, ...items];
          },
        });

        // If none of the settings in this menu are applicable then don't show it.
        const tech = player.tech({ IWillNotUseThisInPlugins: true });

        if (qualities.length < 2 && (!tech || !tech.vhs)) {
          return;
        }

        var menuButton = new MenuButton();
        menuButton.addClass('vjs-quality-selector');
        player.controlBar.addChild(
          menuButton,
          {},
          player.controlBar.children_.length - 2
        );
        this.qualityMenuButton = menuButton;
        this.latencyCompensatorToggleButton = lowLatencyItem;
      }.bind(this)
    );
  }

  addAirplay() {
    videojs.hookOnce('setup', function (player) {
      if (window.WebKitPlaybackTargetAvailabilityEvent) {
        var videoJsButtonClass = videojs.getComponent('Button');
        var concreteButtonClass = videojs.extend(videoJsButtonClass, {
          // The `init()` method will also work for constructor logic here, but it is
          // deprecated. If you provide an `init()` method, it will override the
          // `constructor()` method!
          constructor: function () {
            videoJsButtonClass.call(this, player);
          },

          handleClick: function () {
            const videoElement = document.getElementsByTagName('video')[0];
            videoElement.webkitShowPlaybackTargetPicker();
          },
        });

        var concreteButtonInstance = player.controlBar.addChild(
          new concreteButtonClass()
        );
        concreteButtonInstance.addClass('vjs-airplay');
      }
    });
  }
}

export { OwncastPlayer };

const VjsMenuItem = videojs.getComponent('MenuItem');

export default class MenuSeparator extends VjsMenuItem {
  constructor(player, options) {
    super(player, options);
  }

  createEl(tag = 'button', props = {}, attributes = {}) {
    let el = super.createEl(tag, props, attributes);
    el.innerHTML =
      '<hr style="opacity: 0.3; margin-left: 10px; margin-right: 10px;" />';
    return el;
  }
}

VjsMenuItem.registerComponent('MenuSeparator', MenuSeparator);
