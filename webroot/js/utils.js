export const URL_STATUS = `/status`;
export const URL_CHAT_HISTORY = `/chat`;
export const URL_CUSTOM_EMOJIS = `/emoji`;
export const URL_CONFIG = `/config`;

// TODO: This directory is customizable in the config.  So we should expose this via the config API.
export const URL_STREAM = `/hls/stream.m3u8`;
export const URL_WEBSOCKET = `${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/entry`;

export const TIMER_STATUS_UPDATE = 5000; // ms
export const TIMER_DISABLE_CHAT_AFTER_OFFLINE = 5 * 60 * 1000; // 5 mins
export const TIMER_STREAM_DURATION_COUNTER = 1000;
export const TEMP_IMAGE = 'data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7';

export const MESSAGE_OFFLINE = 'Stream is offline.';
export const MESSAGE_ONLINE = 'Stream is online';

export const URL_OWNCAST = 'https://github.com/gabek/owncast'; // used in footer

export function getLocalStorage(key) {
  try {
    return localStorage.getItem(key);
  } catch (e) {
  }
  return null;
}

export function setLocalStorage(key, value) {
  try {
    if (value !== "" && value !== null) {
      localStorage.setItem(key, value);
    } else {
      localStorage.removeItem(key);
    }
    return true;
  } catch (e) {}
  return false;
}

export function clearLocalStorage(key) {
  localStorage.removeItem(key);
}

// jump down to the max height of a div, with a slight delay
export function jumpToBottom(element) {
  if (!element) return;

  setTimeout(() => {
    element.scrollTo({
      top: element.scrollHeight,
      left: 0,
      behavior: 'smooth'
    });
  }, 50, element);
}

// convert newlines to <br>s
export function addNewlines(str) {
  return str.replace(/(?:\r\n|\r|\n)/g, '<br />');
}

export function pluralize(string, count) {
  if (count === 1) {
    return string;
  } else {
    return string + "s";
  }
}


// Trying to determine if browser is mobile/tablet.
// Source: https://developer.mozilla.org/en-US/docs/Web/HTTP/Browser_detection_using_the_user_agent
export function hasTouchScreen() {
  let hasTouch = false;
  if ("maxTouchPoints" in navigator) {
    hasTouch = navigator.maxTouchPoints > 0;
  } else if ("msMaxTouchPoints" in navigator) {
    hasTouch = navigator.msMaxTouchPoints > 0;
  } else {
      var mQ = window.matchMedia && matchMedia("(pointer:coarse)");
      if (mQ && mQ.media === "(pointer:coarse)") {
        hasTouch = !!mQ.matches;
      } else if ('orientation' in window) {
        hasTouch = true; // deprecated, but good fallback
      } else {
          // Only as a last resort, fall back to user agent sniffing
          var UA = navigator.userAgent;
          hasTouch = (
              /\b(BlackBerry|webOS|iPhone|IEMobile)\b/i.test(UA) ||
              /\b(Android|Windows Phone|iPad|iPod)\b/i.test(UA)
          );
      }
  }
  return hasTouch;
}

// generate random avatar from https://robohash.org
export function generateAvatar(hash) {
  const avatarSource = 'https://robohash.org/';
  const optionSize = '?size=80x80';
  const optionSet = '&set=set3';
  const optionBg = ''; // or &bgset=bg1 or bg2

  return avatarSource + hash + optionSize + optionSet + optionBg;
}

export function generateUsername() {
  return `User ${(Math.floor(Math.random() * 42) + 1)}`;
}

export function secondsToHMMSS(seconds = 0) {
  const finiteSeconds = Number.isFinite(+seconds) ? Math.abs(seconds) : 0;

  const hours = Math.floor(finiteSeconds / 3600);
  const hoursString = hours ? `${hours}:` : '';

  const mins = Math.floor((finiteSeconds / 60) % 60);
  const minString = mins < 10 ? `0${mins}:` : `${mins}:`;

  const secs = Math.floor(finiteSeconds % 60);
  const secsString = secs < 10 ? `0${secs}` : `${secs}`;

  return hoursString + minString + secsString;
}

export function setVHvar() {
  var vh = window.innerHeight * 0.01;
  // Then we set the value in the --vh custom property to the root of the document
  document.documentElement.style.setProperty('--vh', `${vh}px`);
  console.log("== new vh", vh)
}

export function doesObjectSupportFunction(object, functionName) {
  return typeof object[functionName] === "function";
}
