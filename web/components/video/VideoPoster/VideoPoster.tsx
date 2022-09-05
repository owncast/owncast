import { useEffect, useState } from 'react';
import CrossfadeImage from '../../ui/CrossfadeImage/CrossfadeImage';
import s from './VideoPoster.module.scss';

const REFRESH_INTERVAL = 20_000;

interface Props {
  initialSrc: string;
  src: string;
  online: boolean;
}

export const VideoPoster = (props: Props) => {
  const { online, initialSrc, src: base } = props;

  let timer: ReturnType<typeof setInterval>;
  const [src, setSrc] = useState(initialSrc);
  const [duration, setDuration] = useState('0s');

  useEffect(() => {
    clearInterval(timer);
    timer = setInterval(() => {
      if (duration === '0s') {
        setDuration('3s');
      }
      setSrc(`${base}?${Date.now()}`);
    }, REFRESH_INTERVAL);
  }, []);

  return (
    <div className={s.poster}>
      {!online && <img src={initialSrc} alt="logo" />}

      {online && (
        <CrossfadeImage
          src={src}
          duration={duration}
          objectFit="cover"
          width="100%"
          height="100%"
        />
      )}
    </div>
  );
};
export default VideoPoster;
