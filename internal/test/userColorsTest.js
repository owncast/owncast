// Run on the command line to test colors.
// node userColorsTest.js > colorTest.html

console.log(`<body style="background: #2d3748">`);

for (var i = 0; i < 20; i++) {
  const element = generateElement(randomString(6));
  console.log(element);
}

console.log(`</body>`);

function generateElement(string) {
  const bgColor = messageBubbleColorForString(string);
  const fgColor = textColorForString(string);
  return `<div style="background-color: ${bgColor}; color: ${fgColor}; padding: 30px;">${string}</div>`;
}

function randomString(length) {
  return Math.random().toString(36).substring(length);
}

function messageBubbleColorForString(str) {
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

function textColorForString(str) {
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
