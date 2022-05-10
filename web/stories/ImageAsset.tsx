import PropTypes from 'prop-types';

export function ImageAsset(props) {
  const { name, src } = props;

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

  const colorBlockStyle = {
    height: '70%',
    width: '100%',
    backgroundColor: 'white',
  };

  const colorDescriptionStyle = {
    textAlign: 'center',
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
}

Image.propTypes = {
  name: PropTypes.string.isRequired,
};

const rowStyle = {
  display: 'flex',
  flexDirection: 'row',
  flexWrap: 'wrap',
  // justifyContent: 'space-around',
  alignItems: 'center',
};

export function ImageRow(props) {
  const { images } = props;

  return (
    <div style={rowStyle}>
      {images.map(image => (
        <ImageAsset key={image.src} src={image.src} name={image.name} />
      ))}
    </div>
  );
}

ImageRow.propTypes = {
  images: PropTypes.arrayOf(PropTypes.object).isRequired,
};
