import React from 'react';

// This imports the functional component from the previous sample.
import VideoJS from './player';

export default function OwncastPlayer(props) {
  const playerRef = React.useRef(null);
  const { source } = props;

  const videoJsOptions = {
    autoplay: false,
    controls: true,
    responsive: true,
    fluid: true,
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

    // You can handle player events here, for example:
    player.on('waiting', () => {
      player.log('player is waiting');
    });

    player.on('dispose', () => {
      player.log('player will dispose');
    });
  };

  return <VideoJS options={videoJsOptions} onReady={handlePlayerReady} />;
}
