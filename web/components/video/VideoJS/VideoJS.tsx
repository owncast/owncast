import React, { FC } from 'react';
import videojs from 'video.js';
import type VideoJsPlayer from 'video.js/dist/types/player';

import styles from './VideoJS.module.scss';

require('video.js/dist/video-js.css');

export type VideoJSProps = {
  options: any;
  onReady: (player: VideoJsPlayer, vjsInstance: typeof videojs) => void;
};

export const VideoJS: FC<VideoJSProps> = ({ options, onReady }) => {
  const videoRef = React.useRef<HTMLVideoElement | null>(null);
  const playerRef = React.useRef<VideoJsPlayer | null>(null);

  React.useEffect(() => {
    // Make sure Video.js player is only initialized once
    if (!playerRef.current) {
      const videoElement = videoRef.current;

      // eslint-disable-next-line no-multi-assign
      const player: VideoJsPlayer = (playerRef.current = videojs(videoElement, options, () => {
        console.debug('player is ready');
        return onReady && onReady(player, videojs);
      }));

      player.autoplay(options.autoplay);
      player.src(options.sources);
    }
  }, [options, videoRef]);

  React.useEffect(() => {
    videojs.getPlayer(videoRef.current).on('xhr-hooks-ready', () => {
      const cachebusterRequestHook = o => {
        const { uri } = o;
        let updatedURI = uri;
        if (o.uri.match('m3u8')) {
          const u = uri.startsWith('http')
            ? new URL(uri)
            : new URL(uri, window.location.protocol + window.location.host);
          const cachebuster = Math.random().toString(16).slice(2, 8);
          u.searchParams.append('cachebust', cachebuster);
          updatedURI = u.toString();
        }
        return {
          ...o,
          uri: updatedURI,
        };
      };
      (
        videojs.getPlayer(videoRef.current).tech({ IWillNotUseThisInPlugins: true }) as any
      )?.vhs.xhr.onRequest(cachebusterRequestHook);
    });
  }, []);

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
