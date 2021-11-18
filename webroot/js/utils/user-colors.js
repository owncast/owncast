export function messageBubbleColorForHue(hue) {
  // Tweak these to adjust the result of the color
  const saturation = 50;
  const lightness = 50;
  const alpha = 'var(--message-background-alpha)';

  return `hsla(${hue}, ${saturation}%, ${lightness}%, ${alpha})`;
}

export function textColorForHue(hue) {
  // Tweak these to adjust the result of the color
  const saturation = 70;
  const lightness = 80;
  const alpha = 0.85;

  return `hsla(${hue}, ${saturation}%, ${lightness}%, ${alpha})`;
}
