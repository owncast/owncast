import React, { FC } from 'react';
import videojs from 'video.js';
import type VideoJsPlayer from 'video.js/dist/types/player';
import { useTranslation } from 'next-export-i18n';

import styles from './VideoJS.module.scss';

require('video.js/dist/video-js.css');

const SHORTCUT_SUFFIXES: Record<string, string> = {
  Play: ' (Space)',
  Pause: ' (Space)',
  Mute: ' (m)',
  Unmute: ' (m)',
  Fullscreen: ' (f)',
  'Non-Fullscreen': ' (f)',
  'Picture-in-Picture': ' (i)',
  'Exit Picture-in-Picture': ' (i)',
};

export type VideoJSProps = {
  options: any;
  onReady: (player: VideoJsPlayer, vjsInstance: typeof videojs) => void;
};

export const VideoJS: FC<VideoJSProps> = ({ options, onReady }) => {
  const videoRef = React.useRef<HTMLVideoElement | null>(null);
  const playerRef = React.useRef<VideoJsPlayer | null>(null);
  const { t } = useTranslation();

  const addShortcutsToLanguage = (vjs: typeof videojs, langCode: string) => {
    const updates: Record<string, string> = {};
    Object.keys(SHORTCUT_SUFFIXES).forEach(key => {
      const currentLabel = key;
      const suffix = SHORTCUT_SUFFIXES[key];
      updates[key] = t(`${currentLabel}${suffix}`);
    });
    vjs.addLanguage(langCode, updates);
  };

  React.useEffect(() => {
    // Make sure Video.js player is only initialized once
    if (!playerRef.current) {
      const videoElement = videoRef.current;

      addShortcutsToLanguage(videojs, 'en');
      const finalOptions = {
        ...options,
        noUITitleAttributes: true, // Prevents videojs from adding a title attribute to UI elements, thus preventing "double tooltips".
      };
      // eslint-disable-next-line no-multi-assign
      const player: VideoJsPlayer = (playerRef.current = videojs(videoElement, finalOptions, () => {
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
