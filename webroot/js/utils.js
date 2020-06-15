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
function mobileVHhack() {
  setVHvar();
  window.addEventListener("orientationchange", setVHvar);
}

function isAndroidMobile() {
  var isAndroid = navigator.userAgent.toLowerCase().indexOf("android") > -1; 
  //&& ua.indexOf("mobile");
  return isAndroid;
}


