import React, { FC, useEffect } from 'react';
import { useRecoilState, useRecoilValue } from 'recoil';
import { useHotkeys } from 'react-hotkeys-hook';
import { VideoJS } from '../VideoJS/VideoJS';
import ViewerPing from '../viewer-ping';
import { VideoPoster } from '../VideoPoster/VideoPoster';
import { getLocalStorage, setLocalStorage } from '../../../utils/localStorage';
import { isVideoPlayingAtom, clockSkewAtom } from '../../stores/ClientConfigStore';
import PlaybackMetrics from '../metrics/playback';
import createVideoSettingsMenuButton from '../settings-menu';
import LatencyCompensator from '../latencyCompensator';

const VIDEO_CONFIG_URL = '/api/video/variants';
const PLAYER_VOLUME = 'owncast_volume';
const LATENCY_COMPENSATION_ENABLED = 'latencyCompensatorEnabled';

const ping = new ViewerPing();
let playbackMetrics = null;
let latencyCompensator = null;
let latencyCompensatorEnabled = false;

export type OwncastPlayerProps = {
  source: string;
  online: boolean;
};

async function getVideoSettings() {
  let qualities = [];

  try {
    const response = await fetch(VIDEO_CONFIG_URL);
    qualities = await response.json();
  } catch (e) {
    console.error(e);
  }
  return qualities;
}

export const OwncastPlayer: FC<OwncastPlayerProps> = ({ source, online }) => {
  const playerRef = React.useRef(null);
  const [videoPlaying, setVideoPlaying] = useRecoilState<boolean>(isVideoPlayingAtom);
  const clockSkew = useRecoilValue<Number>(clockSkewAtom);

  const setSavedVolume = () => {
    try {
      playerRef.current.volume(getLocalStorage(PLAYER_VOLUME) || 1);
    } catch (err) {
      console.warn(err);
    }
  };

  const handleVolume = () => {
    setLocalStorage(PLAYER_VOLUME, playerRef.current.muted() ? 0 : playerRef.current.volume());
  };

  const togglePlayback = () => {
    if (playerRef.current.paused()) {
      playerRef.current.play();
    } else {
      playerRef.current.pause();
    }
  };

  const toggleMute = () => {
    if (playerRef.current.muted() || playerRef.current.volume() === 0) {
      playerRef.current.volume(0.7);
    } else {
      playerRef.current.volume(0);
    }
  };

  const toggleFullScreen = () => {
    if (playerRef.current.isFullscreen()) {
      playerRef.current.exitFullscreen();
    } else {
      playerRef.current.requestFullscreen();
    }
  };

  const setLatencyCompensatorItemTitle = title => {
    const item = document.querySelector('.latency-toggle-item > .vjs-menu-item-text');
    if (!item) {
      return;
    }

    item.innerHTML = title;
  };

  const startLatencyCompensator = () => {
    if (latencyCompensator) {
      latencyCompensator.stop();
    }

    latencyCompensatorEnabled = true;

    latencyCompensator = new LatencyCompensator(playerRef.current);
    latencyCompensator.setClockSkew(clockSkew);
    latencyCompensator.enable();
    setLocalStorage(LATENCY_COMPENSATION_ENABLED, true);

    setLatencyCompensatorItemTitle('disable minimized latency');
  };

  const stopLatencyCompensator = () => {
    if (latencyCompensator) {
      latencyCompensator.disable();
    }
    latencyCompensator = null;
    latencyCompensatorEnabled = false;
    setLocalStorage(LATENCY_COMPENSATION_ENABLED, false);
    setLatencyCompensatorItemTitle(
      '<span style="font-size: 0.8em">enable minimized latency (experimental)</span>',
    );
  };

  const toggleLatencyCompensator = () => {
    if (latencyCompensatorEnabled) {
      stopLatencyCompensator();
    } else {
      startLatencyCompensator();
    }
  };

  const setupLatencyCompensator = player => {
    const tech = player.tech({ IWillNotUseThisInPlugins: true });

    // VHS is required.
    if (!tech || !tech.vhs) {
      return;
    }

    const latencyCompensatorEnabledSaved = getLocalStorage(LATENCY_COMPENSATION_ENABLED);

    if (latencyCompensatorEnabledSaved === 'true' && tech && tech.vhs) {
      startLatencyCompensator();
    } else {
      stopLatencyCompensator();
    }
  };

  const createSettings = async (player, videojs) => {
    const videoQualities = await getVideoSettings();
    const menuButton = createVideoSettingsMenuButton(
      player,
      videojs,
      videoQualities,
      toggleLatencyCompensator,
    );
    player.controlBar.addChild(
      menuButton,
      {},
      // eslint-disable-next-line no-underscore-dangle
      player.controlBar.children_.length - 2,
    );
    setupLatencyCompensator(player);
  };

  const setupAirplay = (player, videojs) => {
    // eslint-disable-next-line no-prototype-builtins
    if (window.hasOwnProperty('WebKitPlaybackTargetAvailabilityEvent')) {
      const videoJsButtonClass = videojs.getComponent('Button');
      const ConcreteButtonClass = videojs.extend(videoJsButtonClass, {
        // The `init()` method will also work for constructor logic here, but it is
        // deprecated. If you provide an `init()` method, it will override the
        // `constructor()` method!
        constructor() {
          videoJsButtonClass.call(this, player);
        },

        handleClick() {
          try {
            const videoElement = document.getElementsByTagName('video')[0];
            (videoElement as any).webkitShowPlaybackTargetPicker();
          } catch (e) {
            console.error(e);
          }
        },
      });

      const concreteButtonInstance = player.controlBar.addChild(new ConcreteButtonClass());
      concreteButtonInstance.addClass('vjs-airplay');
    }
  };

  // Register keyboard shortcut for the space bar to toggle playback
  useHotkeys('space', e => {
    e.preventDefault();
    togglePlayback();
  });

  // Register keyboard shortcut for f to toggle full screen
  useHotkeys('f', toggleFullScreen, {
    enableOnContentEditable: false,
  });

  // Register keyboard shortcut for the "m" key to toggle mute
  useHotkeys('m', toggleMute, {
    enableOnContentEditable: false,
  });

  useHotkeys('0', () => playerRef.current.volume(playerRef.current.volume() + 0.1), {
    enableOnContentEditable: false,
  });
  useHotkeys('9', () => playerRef.current.volume(playerRef.current.volume() - 0.1), {
    enableOnContentEditable: false,
  });

  const videoJsOptions = {
    autoplay: false,
    controls: true,
    responsive: true,
    fluid: false,
    playsinline: true,
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
    sources: [
      {
        src: source,
        type: 'application/x-mpegURL',
      },
    ],
  };

  const handlePlayerReady = (player, videojs) => {
    playerRef.current = player;
    setSavedVolume();
    setupAirplay(player, videojs);

    // You can handle player events here, for example:
    player.on('waiting', () => {
      player.log('player is waiting');
    });

    player.on('dispose', () => {
      player.log('player will dispose');
      ping.stop();
    });

    player.on('playing', () => {
      player.log('player is playing');
      ping.start();
      setVideoPlaying(true);
    });

    player.on('pause', () => {
      player.log('player is paused');
      ping.stop();
      setVideoPlaying(false);
    });

    player.on('ended', () => {
      player.log('player is ended');
      ping.stop();
      setVideoPlaying(false);
    });

    videojs.hookOnce();

    player.on('volumechange', handleVolume);

    playbackMetrics = new PlaybackMetrics(player, videojs);
    playbackMetrics.setClockSkew(clockSkew);

    createSettings(player, videojs);
  };

  useEffect(() => {
    if (playbackMetrics) {
      playbackMetrics.setClockSkew(clockSkew);
    }
  }, [clockSkew]);

  useEffect(
    () => () => {
      stopLatencyCompensator();
      playbackMetrics.stop();
    },
    [],
  );

  return (
    <div style={{ display: 'grid' }}>
      {online && (
        <div style={{ gridColumn: 1, gridRow: 1 }}>
          <VideoJS options={videoJsOptions} onReady={handlePlayerReady} />
        </div>
      )}
      <div style={{ gridColumn: 1, gridRow: 1 }}>
        {!videoPlaying && (
          <VideoPoster online={online} initialSrc="/thumbnail.jpg" src="/thumbnail.jpg" />
        )}
      </div>
    </div>
  );
};
export default OwncastPlayer;
