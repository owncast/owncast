import React, { useEffect } from 'react';
import { useRecoilState, useRecoilValue } from 'recoil';
import { useHotkeys } from 'react-hotkeys-hook';
import VideoJS from './player';
import ViewerPing from './viewer-ping';
import VideoPoster from './VideoPoster';
import { getLocalStorage, setLocalStorage } from '../../utils/localStorage';
import { isVideoPlayingAtom, clockSkewAtom } from '../stores/ClientConfigStore';
import PlaybackMetrics from './metrics/playback';
import createVideoSettingsMenuButton from './settings-menu';

const VIDEO_CONFIG_URL = '/api/video/variants';
const PLAYER_VOLUME = 'owncast_volume';

const ping = new ViewerPing();
let playbackMetrics = null;
interface Props {
  source: string;
  online: boolean;
}

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

export default function OwncastPlayer(props: Props) {
  const playerRef = React.useRef(null);
  const { source, online } = props;
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

  // Register keyboard shortcut for the space bar to toggle playback
  useHotkeys('space', togglePlayback, {
    enableOnContentEditable: false,
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
    playsInline: true,
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

    const createSettings = async () => {
      const videoQualities = await getVideoSettings();
      const menuButton = createVideoSettingsMenuButton(player, videojs, videoQualities);
      player.controlBar.addChild(
        menuButton,
        {},
        // eslint-disable-next-line no-underscore-dangle
        player.controlBar.children_.length - 2,
      );
      // this.latencyCompensatorToggleButton = lowLatencyItem;
    };

    createSettings();
  };

  useEffect(() => {
    if (playbackMetrics) {
      playbackMetrics.setClockSkew(clockSkew);
    }
  }, [clockSkew]);

  return (
    <div style={{ display: 'grid' }}>
      {online && (
        <div style={{ gridColumn: 1, gridRow: 1 }}>
          <VideoJS options={videoJsOptions} onReady={handlePlayerReady} />
        </div>
      )}
      <div style={{ gridColumn: 1, gridRow: 1 }}>
        {!videoPlaying && <VideoPoster online={online} initialSrc="/logo" src="/thumbnail.jpg" />}
      </div>
    </div>
  );
}
