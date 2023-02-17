import PropTypes from 'prop-types';
import React from 'react';
import { FC } from 'react';

export type ColorProps = {
  color: any; // TODO specify better type
};

export const Color: FC<ColorProps> = ({ color }) => {
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
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    textShadow: '0 0 15px black',
    height: '70%',
    width: '100%',
    backgroundColor: resolvedColor,
  };

  const colorTextStyle = {
    color: 'white',
    alignText: 'center',
  };

  const colorDescriptionStyle = {
    margin: '5px',
    color: 'gray',
    fontSize: '0.95vw',
    textAlign: 'center' as 'center',
    lineHeight: 1.0,
  };

  return (
    <figure style={containerStyle}>
      <div style={colorBlockStyle}>
        <div style={colorTextStyle}>{resolvedColor}</div>
      </div>
      <figcaption style={colorDescriptionStyle}>{color}</figcaption>
    </figure>
  );
};

Color.propTypes = {
  color: PropTypes.string.isRequired,
};

const rowStyle = {
  display: 'flex',
  flexDirection: 'row' as 'row',
  flexWrap: 'wrap' as 'wrap',
  alignItems: 'center',
};

export const ColorRow = props => {
  const { colors } = props;

  return (
    <div style={rowStyle}>
      {colors.map(color => (
        <Color key={color} color={color} />
      ))}
    </div>
  );
};

ColorRow.propTypes = {
  colors: PropTypes.arrayOf(PropTypes.string).isRequired,
};

export const getColor = color => {
  const colorValue = getComputedStyle(document.documentElement).getPropertyValue(`--${color}`);
  return { [color]: colorValue };
  // return { [color]: colorValue };
};
