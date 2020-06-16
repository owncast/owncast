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

function jumpToBottom(id) {
  const div = document.querySelector(id);
  div.scrollTo({
    top: div.scrollHeight,// - div.clientHeight,
    left: 0,
    behavior: 'smooth'
  });
}

function uuidv4() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    const r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

function setVHvar() {
  var vh = window.innerHeight * 0.01;
  // Then we set the value in the --vh custom property to the root of the document
  document.documentElement.style.setProperty('--vh', `${vh}px`);
  console.log("== new vh", vh)
}
// delayed
function mobileVHhack() {
  setTimeout(setVHvar, 100);
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

