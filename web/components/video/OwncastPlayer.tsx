import React from 'react';
import { useRecoilState } from 'recoil';
import VideoJS from './player';
import ViewerPing from './viewer-ping';
import VideoPoster from './VideoPoster';
import { getLocalStorage, setLocalStorage } from '../../utils/helpers';
import { isVideoPlayingAtom } from '../stores/ClientConfigStore';

const PLAYER_VOLUME = 'owncast_volume';

const ping = new ViewerPing();

interface Props {
  source: string;
  online: boolean;
}

export default function OwncastPlayer(props: Props) {
  const playerRef = React.useRef(null);
  const { source, online } = props;
  const [videoPlaying, setVideoPlaying] = useRecoilState<boolean>(isVideoPlayingAtom);

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
        src: source,
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

    player.on('volumechange', handleVolume);
  };

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
