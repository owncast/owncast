import React, { FC, useMemo, useState } from 'react';

type ObjectFit = React.CSSProperties['objectFit'];

export type CrossfadeImageProps = {
  src: string;
  width: string;
  height: string;
  objectFit?: ObjectFit;
  duration?: string;
  className?: string;
};

const imgStyle: React.CSSProperties = {
  position: 'absolute',
  width: `100%`,
  height: `100%`,
};

export const CrossfadeImage: FC<CrossfadeImageProps> = ({
  src = '',
  width,
  height,
  objectFit = 'fill',
  duration = '1s',
  className,
}) => {
  const spanStyle: React.CSSProperties = useMemo(
    () => ({
      display: 'inline-block',
      position: 'relative',
      width,
      height,
    }),
    [width, height],
  );

  const imgStyles = useMemo(
    () => [
      { ...imgStyle, objectFit, opacity: 0, transition: `opacity ${duration}` },
      { ...imgStyle, objectFit, opacity: 1, transition: `opacity ${duration}` },
      { ...imgStyle, objectFit, opacity: 0 },
    ],
    [objectFit, duration],
  );

  const [key, setKey] = useState(0);
  const [srcs, setSrcs] = useState(['', '']);
  const nextSrc = src !== srcs[1] ? src : '';

  const onLoadImg = () => {
    setKey((key + 1) % 3);
    setSrcs([srcs[1], nextSrc]);
  };

  return (
    <span style={spanStyle} className={className}>
      {[...srcs, nextSrc].map(
        (singleSrc, index) =>
          singleSrc !== '' && (
            <img
              key={singleSrc}
              src={singleSrc}
              alt=""
              style={imgStyles[index]}
              onLoad={index === 2 ? onLoadImg : undefined}
            />
          ),
      )}
    </span>
  );
};

CrossfadeImage.defaultProps = {
  objectFit: 'fill',
  duration: '3s',
};
