export function messageBubbleColorForHue(hue) {
  // Tweak these to adjust the result of the color
  const saturation = 45;
  const lightness = 50;
  const alpha = 'var(--message-background-alpha)';

  return `hsla(${hue}, ${saturation}%, ${lightness}%, ${alpha})`;
}

export function textColorForHue(hue) {
  // Tweak these to adjust the result of the color
  const saturation = 80;
  const lightness = 80;
  const alpha = 0.8;

  return `hsla(${hue}, ${saturation}%, ${lightness}%, ${alpha})`;
}
