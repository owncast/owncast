// https://docs.videojs.com/player

import videojs from '/js/web_modules/videojs/dist/video.min.js';
import { getLocalStorage, setLocalStorage } from '../utils/helpers.js';
import { PLAYER_VOLUME, URL_STREAM } from '../utils/constants.js';
import PlaybackMetrics from '../metrics/playback.js';
import LatencyCompensator from './latencyCompensator.js';

const VIDEO_ID = 'video';

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
      maxPlaylistRetries: 30,
    },
  },
  liveTracker: {
    trackingThreshold: 0,
  },
  sources: [VIDEO_SRC],
};

export const POSTER_DEFAULT = `/img/logo.png`;
export const POSTER_THUMB = `/thumbnail.jpg`;

function getCurrentlyPlayingSegment(tech, old_segment = null) {
  var target_media = tech.vhs.playlists.media();
  var snapshot_time = tech.currentTime();

  var segment;
  var segment_time;

  // Itinerate trough available segments and get first within which snapshot_time is
  for (var i = 0, l = target_media.segments.length; i < l; i++) {
    // Note: segment.end may be undefined or is not properly set
    if (snapshot_time < target_media.segments[i].end) {
      segment = target_media.segments[i];
      break;
    }
  }

  // Null segment_time in case it's lower then 0.
  if (segment) {
    segment_time = Math.max(
      0,
      snapshot_time - (segment.end - segment.duration)
    );
    // Because early segments don't have end property
  } else {
    segment = target_media.segments[0];
    segment_time = 0;
  }

  return segment;
}

class OwncastPlayer {
  constructor() {
    window.VIDEOJS_NO_DYNAMIC_STYLE = true; // style override

    this.playbackMetrics = new PlaybackMetrics();

    this.vjsPlayer = null;
    this.latencyCompensator = null;

    this.appPlayerReadyCallback = null;
    this.appPlayerPlayingCallback = null;
    this.appPlayerEndedCallback = null;
    this.hasStartedPlayback = false;

    // bind all the things because safari
    this.startPlayer = this.startPlayer.bind(this);
    this.handleReady = this.handleReady.bind(this);
    this.handlePlaying = this.handlePlaying.bind(this);
    this.handleVolume = this.handleVolume.bind(this);
    this.handleEnded = this.handleEnded.bind(this);
    this.handleError = this.handleError.bind(this);
    this.handleWaiting = this.handleWaiting.bind(this);
    this.handleNoLongerBuffering = this.handleNoLongerBuffering.bind(this);
    this.addQualitySelector = this.addQualitySelector.bind(this);
    this.qualitySelectionMenu = null;
  }

  init() {
    this.addAirplay();
    this.addQualitySelector();

    // Keep a reference of the standard vjs xhr function.
    const oldVjsXhrCallback = videojs.xhr;

    // Override the xhr function to track segment download time.
    videojs.Vhs.xhr = (...args) => {
      if (args[0].uri.match('.ts')) {
        const start = new Date();

        const cb = args[1];
        args[1] = (request, error, response) => {
          const end = new Date();
          const delta = end.getTime() - start.getTime();
          this.playbackMetrics.trackSegmentDownloadTime(delta);
          cb(request, error, response);
        };
      }

      return oldVjsXhrCallback(...args);
    };

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
    // this.vjsPlayer.play();
  }

  handleReady() {
    this.vjsPlayer.on('error', this.handleError);
    this.vjsPlayer.on('playing', this.handlePlaying);
    this.vjsPlayer.on('waiting', this.handleWaiting);
    this.vjsPlayer.on('canplaythrough', this.handleNoLongerBuffering);
    this.vjsPlayer.on('volumechange', this.handleVolume);
    this.vjsPlayer.on('ended', this.handleEnded);

    this.vjsPlayer.on('ready', () => {
      const tech = this.vjsPlayer.tech({ IWillNotUseThisInPlugins: true });
      tech.on('usage', (e) => {
        if (e.name === 'vhs-unknown-waiting') {
          this.playbackMetrics.incrementErrorCount(1);
        }

        if (e.name === 'vhs-rendition-change-abr') {
          // Quality variant has changed
          this.playbackMetrics.incrementQualityVariantChanges();
        }
      });

      // Variant changed
      const trackElements = this.vjsPlayer.textTracks();
      trackElements.addEventListener('cuechange', function (c) {
        console.log(c);
      });

      this.latencyCompensator = new LatencyCompensator(this.vjsPlayer);
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
    this.log('on Playing');
    if (this.appPlayerPlayingCallback) {
      // start polling
      this.appPlayerPlayingCallback();
    }

    if (this.latencyCompensator && !this.hasStartedPlayback) {
      this.latencyCompensator.enable();
    }

    this.hasStartedPlayback = true;

    setInterval(() => {
      this.collectPlaybackMetrics();
    }, 5000);
  }

  collectPlaybackMetrics() {
    const tech = this.vjsPlayer.tech({ IWillNotUseThisInPlugins: true });
    if (!tech || !tech.vhs) {
      return;
    }

    const bandwidth = tech.vhs.systemBandwidth;
    this.playbackMetrics.trackBandwidth(bandwidth);

    try {
      const segment = getCurrentlyPlayingSegment(tech);
      const segmentTime = segment.dateTimeObject.getTime();
      const now = new Date().getTime();
      const latency = now - segmentTime;
      this.playbackMetrics.trackLatency(latency);
    } catch (err) {
      console.warn(err);
    }
  }

  handleEnded() {
    this.log('on Ended');
    if (this.appPlayerEndedCallback) {
      this.appPlayerEndedCallback();
    }

    this.latencyCompensator.disable();
  }

  handleError(e) {
    this.log(`on Error: ${JSON.stringify(e)}`);
    if (this.appPlayerEndedCallback) {
      this.appPlayerEndedCallback();
    }

    this.playbackMetrics.incrementErrorCount(1);
  }

  handleWaiting(e) {
    // this.playbackMetrics.incrementErrorCount(1);
    this.playbackMetrics.isBuffering = true;
  }

  handleNoLongerBuffering() {
    this.playbackMetrics.isBuffering = false;
  }

  log(message) {
    // console.log(`>>> Player: ${message}`);
  }

  async addQualitySelector() {
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
        var MenuButton = videojs.extend(MenuButtonClass, {
          // The `init()` method will also work for constructor logic here, but it is
          // deprecated. If you provide an `init()` method, it will override the
          // `constructor()` method!
          constructor: function () {
            MenuButtonClass.call(this, player);
          },

          createItems: function () {
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
                // Only enable this single, selected representation.
                player
                  .tech({ IWillNotUseThisInPlugins: true })
                  .vhs.representations()
                  .forEach(function (rep, index) {
                    rep.enabled(index === item.index);
                  });
                newMenuItem.selected(false);
              });

              return newMenuItem;
            });

            defaultAutoItem.on('click', function () {
              // Re-enable all representations.
              player
                .tech({ IWillNotUseThisInPlugins: true })
                .vhs.representations()
                .forEach(function (rep, index) {
                  rep.enabled(true);
                });
              defaultAutoItem.selected(false);
            });

            return [defaultAutoItem, ...items];
          },
        });

        // Only show the quality selector if there is more than one option.
        if (qualities.length < 2) {
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
