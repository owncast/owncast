function getHashFromString(string) {
  let hash = 1;
  for (let i = 0; i < string.length; i++) {
    const codepoint = string.charCodeAt(i);
    hash *= codepoint;
  }

  return Math.abs(hash);
}

function digitsFromNumber(number) {
  const numberString = number.toString();
  let digits = [];

  for (let i = 0, len = numberString.length; i < len; i += 1) {
    digits.push(numberString.charAt(i));
  }

  return digits;
}

// function avatarFromString(string) {
//   const hash = getHashFromString(string);
//   const digits = digitsFromNumber(hash);
//   // eslint-disable-next-line
//   const sum = digits.reduce(function (total, number) {
//     return total + number;
//   });
//   const sumDigits = digitsFromNumber(sum);
//   const first = sumDigits[0];
//   const second = sumDigits[1];
//   let filename = '/avatars/';

//   // eslint-disable-next-line
//   if (first == 1 || first == 2) {
//     filename += '1' + second.toString();
//     // eslint-disable-next-line
//   } else if (first == 3 || first == 4) {
//     filename += '2' + second.toString();
//     // eslint-disable-next-line
//   } else if (first == 5 || first == 6) {
//     filename += '3' + second.toString();
//     // eslint-disable-next-line
//   } else if (first == 7 || first == 8) {
//     filename += '4' + second.toString();
//   } else {
//     filename += '5';
//   }

//   return filename + '.svg';
// }

function colorForString(str) {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    // eslint-disable-next-line
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }
  let colour = '#';
  for (let i = 0; i < 3; i++) {
    // eslint-disable-next-line
    let value = (hash >> (i * 8)) & 0xff;
    colour += ('00' + value.toString(16)).substr(-2);
  }
  return colour;
}

function messageBubbleColorForString(str) {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    // eslint-disable-next-line
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }
  let color = '#';
  for (let i = 0; i < 3; i++) {
    // eslint-disable-next-line
    let value = (hash >> (i * 8)) & 0xff;
    color += ('00' + value.toString(16)).substr(-2);
  }
  // Convert to RGBA
  let result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(color);
  let rgb = result ? {
    r: parseInt(result[1], 16),
    g: parseInt(result[2], 16),
    b: parseInt(result[3], 16),
  } : null;
  return 'rgba(' + rgb.r + ',' + rgb.g + ',' + rgb.b + ', 0.4)';
}