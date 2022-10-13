import React, { FC } from 'react';
import videojs from 'video.js';
import styles from './VideoJS.module.scss';

require('video.js/dist/video-js.css');

// TODO: Restore volume that was saved in local storage.
// import { getLocalStorage, setLocalStorage } from '../../utils/helpers.js';
// import { PLAYER_VOLUME, URL_STREAM } from '../../utils/constants.js';
export type VideoJSProps = {
  options: any;
  onReady: (player: videojs.Player, vjsInstance: videojs) => void;
};

export const VideoJS: FC<VideoJSProps> = ({ options, onReady }) => {
  const videoRef = React.useRef(null);
  const playerRef = React.useRef(null);

  React.useEffect(() => {
    // Make sure Video.js player is only initialized once
    if (!playerRef.current) {
      const videoElement = videoRef.current;

      // eslint-disable-next-line no-multi-assign
      const player = (playerRef.current = videojs(videoElement, options, () => {
        player.log('player is ready');
        return onReady && onReady(player, videojs);
      }));

      // TODO: Add airplay support, video settings menu, latency compensator, etc.

      // You can update player in the `else` block here, for example:
      // } else {
      player.autoplay(options.autoplay);
      player.src(options.sources);
    }
  }, [options, videoRef]);

  // Dispose the Video.js player when the functional component unmounts
  React.useEffect(() => {
    const player = playerRef.current;

    return () => {
      if (player) {
        player.dispose();
        playerRef.current = null;
      }
    };
  }, [playerRef]);

  return (
    <div data-vjs-player>
      {/* eslint-disable-next-line jsx-a11y/media-has-caption */}
      <video
        ref={videoRef}
        className={`video-js vjs-big-play-centered vjs-show-big-play-button-on-pause ${styles.player} vjs-owncast`}
      />
    </div>
  );
};

export default VideoJS;

// import PlaybackMetrics from '../metrics/playback.js';
// import LatencyCompensator from './latencyCompensator.js';

// const VIDEO_ID = 'video';
// const LATENCY_COMPENSATION_ENABLED = 'latencyCompensatorEnabled';

// // Video setup
// const VIDEO_SRC = {
//   src: URL_STREAM,
//   type: 'application/x-mpegURL',
// };
// const VIDEO_OPTIONS = {
//   autoplay: false,
//   liveui: true,
//   preload: 'auto',
//   controlBar: {
//     progressControl: {
//       seekBar: false,
//     },
//   },
//   html5: {
//     vhs: {
//       // used to select the lowest bitrate playlist initially. This helps to decrease playback start time. This setting is false by default.
//       enableLowInitialPlaylist: true,
//       experimentalBufferBasedABR: true,
//       useNetworkInformationApi: true,
//       maxPlaylistRetries: 30,
//     },
//   },
//   liveTracker: {
//     trackingThreshold: 0,
//     liveTolerance: 15,
//   },
//   sources: [VIDEO_SRC],
// };

// export const POSTER_DEFAULT = `/img/logo.png`;
// export const POSTER_THUMB = `/thumbnail.jpg`;

// export default class MenuSeparator extends VjsMenuItem {
//   constructor(player, options) {
//     super(player, options);
//   }

//   createEl(tag = 'button', props = {}, attributes = {}) {
//     const el = super.createEl(tag, props, attributes);
//     el.innerHTML = '<hr style="opacity: 0.3; margin-left: 10px; margin-right: 10px;" />';
//     return el;
//   }
// }

// class OwncastPlayer {
//   constructor() {
//     window.VIDEOJS_NO_DYNAMIC_STYLE = true; // style override

//     this.vjsPlayer = null;
//     this.latencyCompensator = null;
//     this.playbackMetrics = null;

//     this.appPlayerReadyCallback = null;
//     this.appPlayerPlayingCallback = null;
//     this.appPlayerEndedCallback = null;

//     this.hasStartedPlayback = false;
//     this.latencyCompensatorEnabled = false;

//     // bind all the things because safari
//     this.startPlayer = this.startPlayer.bind(this);
//     this.handleReady = this.handleReady.bind(this);
//     this.handlePlaying = this.handlePlaying.bind(this);
//     this.handleVolume = this.handleVolume.bind(this);
//     this.handleEnded = this.handleEnded.bind(this);
//     this.handleError = this.handleError.bind(this);
//     this.addQualitySelector = this.addVideoSettingsMenu.bind(this);
//     this.addQualitySelector = this.addVideoSettingsMenu.bind(this);
//     this.toggleLatencyCompensator = this.toggleLatencyCompensator.bind(this);
//     this.startLatencyCompensator = this.startLatencyCompensator.bind(this);
//     this.stopLatencyCompensator = this.stopLatencyCompensator.bind(this);
//     this.qualitySelectionMenu = null;
//     this.latencyCompensatorToggleButton = null;
//   }

//   init() {
//     // this.addAirplay();
//     // this.addVideoSettingsMenu();

//     // Add a cachebuster param to playlist URLs.
//     videojs.Vhs.xhr.beforeRequest = options => {
//       if (options.uri.match('m3u8')) {
//         const cachebuster = Math.random().toString(16).substr(2, 8);
//         options.uri = `${options.uri}?cachebust=${cachebuster}`;
//       }

//       return options;
//     };

//     this.vjsPlayer = videojs(VIDEO_ID, VIDEO_OPTIONS);

//     this.vjsPlayer.ready(this.handleReady);
//   }

//   setupPlayerCallbacks(callbacks) {
//     const { onReady, onPlaying, onEnded, onError } = callbacks;

//     this.appPlayerReadyCallback = onReady;
//     this.appPlayerPlayingCallback = onPlaying;
//     this.appPlayerEndedCallback = onEnded;
//     this.appPlayerErrorCallback = onError;
//   }

//   // play
//   startPlayer() {
//     this.log('Start playing');
//     const source = { ...VIDEO_SRC };

//     try {
//       this.vjsPlayer.volume(getLocalStorage(PLAYER_VOLUME) || 1);
//     } catch (err) {
//       console.warn(err);
//     }
//     this.vjsPlayer.src(source);
//   }

//   setupPlaybackMetrics() {
//     // this.playbackMetrics = new PlaybackMetrics(this.vjsPlayer, videojs);
//   }

//   setupLatencyCompensator() {
//     const tech = this.vjsPlayer.tech({ IWillNotUseThisInPlugins: true });

//     // VHS is required.
//     if (!tech || !tech.vhs) {
//     }

//     // const latencyCompensatorEnabledSaved = getLocalStorage(LATENCY_COMPENSATION_ENABLED);

//     // if (latencyCompensatorEnabledSaved === 'true' && tech && tech.vhs) {
//     //   this.startLatencyCompensator();
//     // } else {
//     //   this.stopLatencyCompensator();
//     // }
//   }

//   startLatencyCompensator() {
//     // this.latencyCompensator = new LatencyCompensator(this.vjsPlayer);
//     // this.latencyCompensator.enable();
//     // this.latencyCompensatorEnabled = true;
//     // this.setLatencyCompensatorItemTitle('disable minimized latency');
//   }

//   stopLatencyCompensator() {
//     // if (this.latencyCompensator) {
//     //   this.latencyCompensator.disable();
//     // }
//     // this.LatencyCompensator = null;
//     // this.latencyCompensatorEnabled = false;
//     // this.setLatencyCompensatorItemTitle(
//     //   '<span style="font-size: 0.8em">enable minimized latency (experimental)</span>',
//     // );
//   }

//   handleReady() {
//     console.log('handleReady');
//     this.vjsPlayer.on('error', this.handleError);
//     this.vjsPlayer.on('playing', this.handlePlaying);
//     this.vjsPlayer.on('volumechange', this.handleVolume);
//     this.vjsPlayer.on('ended', this.handleEnded);

//     this.vjsPlayer.on('loadeddata', () => {
//       this.setupPlaybackMetrics();
//       this.setupLatencyCompensator();
//     });

//     if (this.appPlayerReadyCallback) {
//       // start polling
//       this.appPlayerReadyCallback();
//     }

//     this.vjsPlayer.log.level('debug');
//   }

//   handleVolume() {
//     setLocalStorage(PLAYER_VOLUME, this.vjsPlayer.muted() ? 0 : this.vjsPlayer.volume());
//   }

//   handlePlaying() {
//     if (this.appPlayerPlayingCallback) {
//       // start polling
//       this.appPlayerPlayingCallback();
//     }

//     if (this.latencyCompensator && !this.hasStartedPlayback) {
//       this.latencyCompensator.enable();
//     }

//     this.hasStartedPlayback = true;
//   }

//   handleEnded() {
//     if (this.appPlayerEndedCallback) {
//       this.appPlayerEndedCallback();
//     }

//     this.stopLatencyCompensator();
//   }

//   handleError(e) {
//     this.log(`on Error: ${JSON.stringify(e)}`);
//     if (this.appPlayerEndedCallback) {
//       this.appPlayerEndedCallback();
//     }
//   }

//   toggleLatencyCompensator() {
//     if (this.latencyCompensatorEnabled) {
//       this.stopLatencyCompensator();
//       setLocalStorage(LATENCY_COMPENSATION_ENABLED, false);
//     } else {
//       this.startLatencyCompensator();
//       setLocalStorage(LATENCY_COMPENSATION_ENABLED, true);
//     }
//   }

//   setLatencyCompensatorItemTitle(title) {
//     const item = document.querySelector('.latency-toggle-item > .vjs-menu-item-text');
//     if (!item) {
//       return;
//     }

//     item.innerHTML = title;
//   }

//   log(message) {
//     // console.log(`>>> Player: ${message}`);
//   }

//   render(): JSX.Element {
//     return (
//       <div data-vjs-player>
//         <video ref={videoRef} className='video-js vjs-big-play-centered' />
//       </div>
//     );
//   }

//   async addVideoSettingsMenu() {
//     if (this.qualityMenuButton) {
//       player.controlBar.removeChild(this.qualityMenuButton);
//     }

//     videojs.hookOnce('setup', async player => {
//       let qualities = [];

//       try {
//         const response = await fetch('http://localhost:8080/api/video/variants');
//         qualities = await response.json();
//       } catch (e) {
//         console.error(e);
//       }

//       const MenuItem = videojs.getComponent('MenuItem');
//       const MenuButtonClass = videojs.getComponent('MenuButton');

//       const lowLatencyItem = new MenuItem(player, {
//         selectable: true,
//       });
//       lowLatencyItem.setAttribute('class', 'latency-toggle-item');
//       lowLatencyItem.on('click', () => {
//         this.toggleLatencyCompensator();
//       });
//       this.latencyCompensatorToggleButton = lowLatencyItem;

//       const separator = new MenuSeparator(player, {
//         selectable: false,
//       });

//       const MenuButton = videojs.extend(MenuButtonClass, {
//         // The `init()` method will also work for constructor logic here, but it is
//         // deprecated. If you provide an `init()` method, it will override the
//         // `constructor()` method!
//         constructor() {
//           MenuButtonClass.call(this, player);
//         },

//         createItems() {
//           const tech = this.player_.tech({ IWillNotUseThisInPlugins: true });

//           const defaultAutoItem = new MenuItem(player, {
//             selectable: true,
//             label: 'Auto',
//           });

//           const items = qualities.map(item => {
//             const newMenuItem = new MenuItem(player, {
//               selectable: true,
//               label: item.name,
//             });

//             // Quality selected
//             newMenuItem.on('click', () => {
//               // If for some reason tech doesn't exist, then don't do anything
//               if (!tech) {
//                 console.warn('Invalid attempt to access null player tech');
//                 return;
//               }
//               // Only enable this single, selected representation.
//               tech.vhs.representations().forEach((rep, index) => {
//                 rep.enabled(index === item.index);
//               });
//               newMenuItem.selected(false);
//             });

//             return newMenuItem;
//           });

//           defaultAutoItem.on('click', () => {
//             // Re-enable all representations.
//             tech.vhs.representations().forEach((rep, index) => {
//               rep.enabled(true);
//             });
//             defaultAutoItem.selected(false);
//           });

//           const supportsLatencyCompensator = !!tech && !!tech.vhs;

//           // Only show the quality selector if there is more than one option.
//           if (qualities.length < 2 && supportsLatencyCompensator) {
//             return [lowLatencyItem];
//           }
//           if (qualities.length > 1 && supportsLatencyCompensator) {
//             return [defaultAutoItem, ...items, separator, lowLatencyItem];
//           }
//           if (!supportsLatencyCompensator && qualities.length == 1) {
//             return [];
//           }

//           return [defaultAutoItem, ...items];
//         },
//       });

//       // If none of the settings in this menu are applicable then don't show it.
//       const tech = player.tech({ IWillNotUseThisInPlugins: true });

//       if (qualities.length < 2 && (!tech || !tech.vhs)) {
//         return;
//       }

//       const menuButton = new MenuButton();
//       menuButton.addClass('vjs-quality-selector');
//       player.controlBar.addChild(menuButton, {}, player.controlBar.children_.length - 2);
//       this.qualityMenuButton = menuButton;
//       this.latencyCompensatorToggleButton = lowLatencyItem;
//     });
//   }

//   addAirplay() {
//     videojs.hookOnce('setup', player => {
//       if (window.WebKitPlaybackTargetAvailabilityEvent) {
//         const videoJsButtonClass = videojs.getComponent('Button');
//         const concreteButtonClass = videojs.extend(videoJsButtonClass, {
//           // The `init()` method will also work for constructor logic here, but it is
//           // deprecated. If you provide an `init()` method, it will override the
//           // `constructor()` method!
//           constructor() {
//             videoJsButtonClass.call(this, player);
//           },

//           handleClick() {
//             const videoElement = document.getElementsByTagName('video')[0];
//             videoElement.webkitShowPlaybackTargetPicker();
//           },
//         });

//         const concreteButtonInstance = player.controlBar.addChild(new concreteButtonClass());
//         concreteButtonInstance.addClass('vjs-airplay');
//       }
//     });
//   }

// }

// export { OwncastPlayer };

// const VjsMenuItem = videojs.getComponent('MenuItem');

// VjsMenuItem.registerComponent('MenuSeparator', MenuSeparator);
