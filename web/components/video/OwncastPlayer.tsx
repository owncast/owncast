import React from 'react';
import VideoJS from './player';
import ViewerPing from './viewer-ping';
import { getLocalStorage, setLocalStorage } from '../../utils/helpers';

const PLAYER_VOLUME = 'owncast_volume';

const ping = new ViewerPing();

export default function OwncastPlayer(props) {
  const playerRef = React.useRef(null);
  const { source } = props;

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
        src: `${source}/hls/stream.m3u8`,
        type: 'application/x-mpegURL',
      },
    ],
  };

  const handlePlayerReady = player => {
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
    });

    player.on('pause', () => {
      player.log('player is paused');
      ping.stop();
    });

    player.on('ended', () => {
      player.log('player is ended');
      ping.stop();
    });

    player.on('volumechange', handleVolume);
  };

  return <VideoJS options={videoJsOptions} onReady={handlePlayerReady} />;
}
