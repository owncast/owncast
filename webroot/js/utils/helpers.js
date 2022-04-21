import { ORIENTATION_LANDSCAPE, ORIENTATION_PORTRAIT } from './constants.js';

export function getLocalStorage(key) {
  try {
    return localStorage.getItem(key);
  } catch (e) {}
  return null;
}

export function setLocalStorage(key, value) {
  try {
    if (value !== '' && value !== null) {
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
export function jumpToBottom(element, behavior) {
  if (!element) return;

  if (!behavior) {
    behavior = document.visibilityState === 'visible' ? 'smooth' : 'instant';
  }

  setTimeout(
    () => {
      element.scrollTo({
        top: element.scrollHeight,
        left: 0,
        behavior: behavior,
      });
    },
    50,
    element
  );
}

// convert newlines to <br>s
export function addNewlines(str) {
  return str.replace(/(?:\r\n|\r|\n)/g, '<br />');
}

export function pluralize(string, count) {
  if (count === 1) {
    return string;
  } else {
    return string + 's';
  }
}

// Trying to determine if browser is mobile/tablet.
// Source: https://developer.mozilla.org/en-US/docs/Web/HTTP/Browser_detection_using_the_user_agent
export function hasTouchScreen() {
  let hasTouch = false;
  if ('maxTouchPoints' in navigator) {
    hasTouch = navigator.maxTouchPoints > 0;
  } else if ('msMaxTouchPoints' in navigator) {
    hasTouch = navigator.msMaxTouchPoints > 0;
  } else {
    var mQ = window.matchMedia && matchMedia('(pointer:coarse)');
    if (mQ && mQ.media === '(pointer:coarse)') {
      hasTouch = !!mQ.matches;
    } else if ('orientation' in window) {
      hasTouch = true; // deprecated, but good fallback
    } else {
      // Only as a last resort, fall back to user agent sniffing
      hasTouch = navigator.userAgentData.mobile;
    }
  }
  return hasTouch;
}

export function getOrientation(forTouch = false) {
  // chrome mobile gives misleading matchMedia result when keyboard is up
  if (forTouch && window.screen && window.screen.orientation) {
    return window.screen.orientation.type.match('portrait')
      ? ORIENTATION_PORTRAIT
      : ORIENTATION_LANDSCAPE;
  } else {
    // all other cases
    return window.matchMedia('(orientation: portrait)').matches
      ? ORIENTATION_PORTRAIT
      : ORIENTATION_LANDSCAPE;
  }
}

export function padLeft(text, pad, size) {
  return String(pad.repeat(size) + text).slice(-size);
}

export function parseSecondsToDurationString(seconds = 0) {
  const finiteSeconds = Number.isFinite(+seconds) ? Math.abs(seconds) : 0;

  const days = Math.floor(finiteSeconds / 86400);
  const daysString = days > 0 ? `${days} day${days > 1 ? 's' : ''} ` : '';

  const hours = Math.floor((finiteSeconds / 3600) % 24);
  const hoursString = hours || days ? padLeft(`${hours}:`, '0', 3) : '';

  const mins = Math.floor((finiteSeconds / 60) % 60);
  const minString = padLeft(`${mins}:`, '0', 3);

  const secs = Math.floor(finiteSeconds % 60);
  const secsString = padLeft(`${secs}`, '0', 2);

  return daysString + hoursString + minString + secsString;
}

export function setVHvar() {
  var vh = window.innerHeight * 0.01;
  // Then we set the value in the --vh custom property to the root of the document
  document.documentElement.style.setProperty('--vh', `${vh}px`);
}

export function doesObjectSupportFunction(object, functionName) {
  return typeof object[functionName] === 'function';
}

// return a string of css classes
export function classNames(json) {
  const classes = [];

  Object.entries(json).map(function (item) {
    const [key, value] = item;
    if (value) {
      classes.push(key);
    }
    return null;
  });
  return classes.join(' ');
}

// taken from
// https://medium.com/@TCAS3/debounce-deep-dive-javascript-es6-e6f8d983b7a1
export function debounce(fn, time) {
  let timeout;

  return function () {
    const functionCall = () => fn.apply(this, arguments);

    clearTimeout(timeout);
    timeout = setTimeout(functionCall, time);
  };
}

export function getDiffInDaysFromNow(timestamp) {
  const time = typeof timestamp === 'string' ? new Date(timestamp) : timestamp;
  return (new Date() - time) / (24 * 3600 * 1000);
}

// "Last live today at [time]" or "last live [date]"
export function makeLastOnlineString(timestamp) {
  if (!timestamp) {
    return '';
  }
  let string = '';
  const time = new Date(timestamp);
  const comparisonDate = new Date(time).setHours(0, 0, 0, 0);

  if (comparisonDate == new Date().setHours(0, 0, 0, 0)) {
    const atTime = time.toLocaleTimeString([], {
      hour: '2-digit',
      minute: '2-digit',
    });
    string = `Today ${atTime}`;
  } else {
    string = time.toLocaleDateString();
  }

  return `Last live: ${string}`;
}

// Routing & Tabs
export const ROUTE_RECORDINGS = 'recordings';
export const ROUTE_SCHEDULE = 'schedule';
// looks for `/recording|schedule/id` pattern to determine what to display from the tab view
export function checkUrlPathForDisplay() {
  const pathTest = [ROUTE_RECORDINGS, ROUTE_SCHEDULE];
  const pathParts = window.location.pathname.split('/');

  if (pathParts.length >= 2) {
    const part = pathParts[1].toLowerCase();
    if (pathTest.includes(part)) {
      return {
        section: part,
        sectionId: pathParts[2] || '',
      };
    }
  }
  return null;
}

export function paginateArray(items, page, perPage) {
  const offset = perPage * (page - 1);
  const totalPages = Math.ceil(items.length / perPage);
  const paginatedItems = items.slice(offset, perPage * page);

  return {
    previousPage: page - 1 ? page - 1 : null,
    nextPage: totalPages > page ? page + 1 : null,
    total: items.length,
    totalPages: totalPages,
    items: paginatedItems,
  };
}
