import PropTypes from 'prop-types';

export function Color(props) {
  const { color } = props;
  const resolvedColor = getComputedStyle(document.documentElement).getPropertyValue(`--${color}`);

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
    backgroundColor: resolvedColor,
  };

  const colorDescriptionStyle = {
    margin: '5px',
    color: 'gray',
    fontSize: '0.8em',
  };

  return (
    <figure style={containerStyle}>
      <div style={colorBlockStyle} />
      <figcaption>
        <span style={colorDescriptionStyle}>{color}</span>
      </figcaption>
    </figure>
  );
}

Color.propTypes = {
  color: PropTypes.string.isRequired,
};

const rowStyle = {
  display: 'flex',
  flexDirection: 'row' as 'row',
  flexWrap: 'wrap' as 'wrap',
  // justifyContent: 'space-around',
  alignItems: 'center',
};

export function ColorRow(props) {
  const { colors } = props;

  return (
    <div style={rowStyle}>
      {colors.map(color => (
        <Color key={color} color={color} />
      ))}
    </div>
  );
}

ColorRow.propTypes = {
  colors: PropTypes.arrayOf(PropTypes.string).isRequired,
};
