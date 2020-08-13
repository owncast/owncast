
const URL_STATUS = `/status`;
const URL_CHAT_HISTORY = `/chat`;
// TODO: This directory is customizable in the config.  So we should expose this via the config API.
const URL_STREAM = `/hls/stream.m3u8`;
const URL_WEBSOCKET = `${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/entry`;

const POSTER_DEFAULT = `/img/logo.png`;
const POSTER_THUMB = `/thumbnail.jpg`;

const URL_OWNCAST = 'https://github.com/gabek/owncast'; // used in footer

function getLocalStorage(key) {
  try {
    return localStorage.getItem(key);
  } catch (e) {
  }
  return null;
}

function setLocalStorage(key, value) {
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

function clearLocalStorage(key) {
  localStorage.removeItem(key);
}

// jump down to the max height of a div, with a slight delay
function jumpToBottom(element) {
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
function addNewlines(str) {
  return str.replace(/(?:\r\n|\r|\n)/g, '<br />');
}

function pluralize(string, count) {
  if (count === 1) {
    return string;
  } else {
    return string + "s";
  }
}


// Trying to determine if browser is mobile/tablet.
// Source: https://developer.mozilla.org/en-US/docs/Web/HTTP/Browser_detection_using_the_user_agent
function hasTouchScreen() {
  var hasTouchScreen = false;
  if ("maxTouchPoints" in navigator) { 
      hasTouchScreen = navigator.maxTouchPoints > 0;
  } else if ("msMaxTouchPoints" in navigator) {
      hasTouchScreen = navigator.msMaxTouchPoints > 0; 
  } else {
      var mQ = window.matchMedia && matchMedia("(pointer:coarse)");
      if (mQ && mQ.media === "(pointer:coarse)") {
          hasTouchScreen = !!mQ.matches;
      } else if ('orientation' in window) {
          hasTouchScreen = true; // deprecated, but good fallback
      } else {
          // Only as a last resort, fall back to user agent sniffing
          var UA = navigator.userAgent;
          hasTouchScreen = (
              /\b(BlackBerry|webOS|iPhone|IEMobile)\b/i.test(UA) ||
              /\b(Android|Windows Phone|iPad|iPod)\b/i.test(UA)
          );
      }
  }
  return hasTouchScreen;
}

// generate random avatar from https://robohash.org
function generateAvatar(hash) {
  const avatarSource = 'https://robohash.org/';
  const optionSize = '?size=80x80';
  const optionSet = '&set=set3'; 
  const optionBg = ''; // or &bgset=bg1 or bg2

  return avatarSource + hash + optionSize + optionSet + optionBg;
}

function generateUsername() {
  return `User ${(Math.floor(Math.random() * 42) + 1)}`;
}

function secondsToHMMSS(seconds = 0) {
  const finiteSeconds = Number.isFinite(+seconds) ? Math.abs(seconds) : 0;

  const hours = Math.floor(finiteSeconds / 3600);
  const hoursString = hours ? `${hours}:` : '';

  const mins = Math.floor((finiteSeconds / 60) % 60);
  const minString = mins < 10 ? `0${mins}:` : `${mins}:`;

  const secs = Math.floor(finiteSeconds % 60);
  const secsString = secs < 10 ? `0${secs}` : `${secs}`;

  return hoursString + minString + secsString;
}

function setVHvar() {
  var vh = window.innerHeight * 0.01;
  // Then we set the value in the --vh custom property to the root of the document
  document.documentElement.style.setProperty('--vh', `${vh}px`);
  console.log("== new vh", vh)
}

function doesObjectSupportFunction(object, functionName) {
  return typeof object[functionName] === "function";
}