
const LOCAL_TEST = true;


const MESSAGE_OFFLINE = 'Stream is offline.';
const MESSAGE_ONLINE = 'Stream is online.';

const URL_PREFIX = LOCAL_TEST ? 'https://goth.land' : ''; 
const URL_STATUS = `${URL_PREFIX}/status`;
const URL_STREAM = `${URL_PREFIX}/hls/stream.m3u8`;

const URL_WEBSOCKET = LOCAL_TEST 
  ? 'wss://goth.land/entry'
  : `${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/entry`;

const POSTER_DEFAULT = `${URL_PREFIX}/img/logo.png`;
const POSTER_THUMB = `${URL_PREFIX}/thumbnail.jpg`;

const URL_CONFIG = `${URL_PREFIX}/config`;

const SOCKET_MESSAGE_TYPES = {
	CHAT: 'CHAT',
	PING: 'PING'
}

const KEY_USERNAME = 'owncast_username';
const KEY_AVATAR = 'owncast_avatar';
const KEY_CHAT_DISPLAYED = 'owncast_chat';

const VIDEO_ID = 'video';

const URL_OWNCAST = 'https://github.com/gabek/owncast';

const TIMER_STATUS_UPDATE = 3000; // ms
const TIMER_WEBSOCKET_RECONNECT = 3000; // ms

const TEMP_IMAGE = 'data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7';

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

