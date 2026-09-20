import { useEffect, useState, type CSSProperties, type FC, type ImgHTMLAttributes } from 'react';
import { CrossfadeImage } from '../../ui/CrossfadeImage/CrossfadeImage';
import styles from './VideoPoster.module.scss';

const REFRESH_INTERVAL = 20_000;

export type VideoPosterProps = {
  initialSrc: string;
  src: string;
  online: boolean;
  refreshInterval?: number;
  objectFit?: CSSProperties['objectFit'];
  width?: string;
  height?: string;
  referrerPolicy?: ImgHTMLAttributes<HTMLImageElement>['referrerPolicy'];
  onError?: ImgHTMLAttributes<HTMLImageElement>['onError'];
};

export const VideoPoster: FC<VideoPosterProps> = ({
  online,
  initialSrc,
  src: base,
  refreshInterval = REFRESH_INTERVAL,
  objectFit = 'contain',
  width = '100%',
  height = 'auto',
  referrerPolicy,
  onError,
}) => {
  const [src, setSrc] = useState(initialSrc);
  const [duration, setDuration] = useState('0s');

  useEffect(() => {
    setSrc(initialSrc);
    setDuration('0s');

    if (!online || refreshInterval <= 0) {
      return undefined;
    }

    const timer = setInterval(() => {
      setDuration(current => (current === '0s' ? '3s' : current));
      setSrc(`${base}${base.includes('?') ? '&' : '?'}cb=${Date.now()}`);
    }, refreshInterval);

    return () => clearInterval(timer);
  }, [base, initialSrc, online, refreshInterval]);

  return (
    <div className={styles.poster}>
      {!online && <img src={initialSrc} alt="logo" />}

      {online && (
        <CrossfadeImage
          src={src}
          duration={duration}
          objectFit={objectFit}
          height={height}
          width={width}
          className={styles.image}
          referrerPolicy={referrerPolicy}
          onError={onError}
        />
      )}
    </div>
  );
};
