import { FC, useMemo } from 'react';
import { VideoJS } from '../VideoJS/VideoJS';
import styles from './ClipPlayer.module.scss';

export type ClipPlayerProps = {
  // source is the clip's HLS master playlist.
  source: string;
  // poster is the clip's thumbnail, shown before playback starts.
  poster?: string;
  autoplay?: boolean;
};

// ClipPlayer plays a clip as ordinary VOD. The live player's machinery
// (viewer pings, latency compensation, playback metrics) does not apply to a
// finished clip, so this uses video.js directly.
export const ClipPlayer: FC<ClipPlayerProps> = ({ source, poster, autoplay = false }) => {
  const options = useMemo(
    () => ({
      autoplay,
      controls: true,
      responsive: true,
      fluid: true,
      poster,
      sources: [
        {
          src: source,
          type: 'application/x-mpegURL',
        },
      ],
    }),
    [source, poster, autoplay],
  );

  return (
    <div className={styles.player}>
      <VideoJS options={options} onReady={() => {}} />
    </div>
  );
};
