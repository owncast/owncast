import React from 'react';
import { FC } from 'react';

export type ImageAssetProps = {
  name: string;
  src: string;
};

export const ImageAsset: FC<ImageAssetProps> = ({ name, src }) => {
  const containerStyle = {
    borderRadius: '20px',
    width: '12vw',
    height: '12vw',
    minWidth: '100px',
    minHeight: '100px',
    borderWidth: '1.5px',
    borderStyle: 'solid',
    borderColor: 'lightgray',
    overflow: 'hidden',
    margin: '0.3vw',
  };

  const colorDescriptionStyle = {
    textAlign: 'center' as 'center',
    color: 'gray',
    fontSize: '0.8em',
  };

  const imageStyle = {
    width: '100%',
    height: '80%',
    backgroundRepeat: 'no-repeat',
    backgroundSize: 'contain',
    backgroundPosition: 'center',
    marginTop: '5px',
    backgroundImage: `url(${src})`,
  };

  return (
    <figure style={containerStyle}>
      <a href={src} target="_blank" rel="noopener noreferrer">
        <div style={imageStyle} />
        <figcaption style={colorDescriptionStyle}>{name}</figcaption>
      </a>
    </figure>
  );
};

const rowStyle = {
  display: 'flex',
  flexDirection: 'row' as 'row',
  flexWrap: 'wrap' as 'wrap',
  // justifyContent: 'space-around',
  alignItems: 'center',
};

export const ImageRow = (props: ImageRowProps) => {
  const { images } = props;

  return (
    <div style={rowStyle}>
      {images.map(image => (
        <ImageAsset key={image.src} src={image.src} name={image.name} />
      ))}
    </div>
  );
};

interface ImageRowProps {
  images: ImageAssetProps[];
}
