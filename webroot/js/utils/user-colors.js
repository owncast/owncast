export function messageBubbleColorForString(str) {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    // eslint-disable-next-line
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }

  // Tweak these to adjust the result of the color
  const saturation = 25;
  const lightness = 45;
  const alpha = 0.3;
  const hue = parseInt(Math.abs(hash), 16) % 360;

  return `hsla(${hue}, ${saturation}%, ${lightness}%, ${alpha})`;
}

export function textColorForString(str) {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    // eslint-disable-next-line
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }

  // Tweak these to adjust the result of the color
  const saturation = 80;
  const lightness = 80;
  const alpha = 0.8;
  const hue = parseInt(Math.abs(hash), 16) % 360;

  return `hsla(${hue}, ${saturation}%, ${lightness}%, ${alpha})`;
}
