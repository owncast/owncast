export function messageBubbleColorForString(str) {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    // eslint-disable-next-line
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }

  // Tweak these to adjust the result of the color
  const saturation = 75;
  const lightness = 65;
  const alpha = 1.0;
  const hue = parseInt(Math.abs(hash), 16) % 300;

  return `hsla(${hue}, ${saturation}%, ${lightness}%, ${alpha})`;
}
