function getLocalStorage(key) {
  try {
    return localStorage.getItem(key);
  } catch (e) {
  }
  return null;
}

function setLocalStorage(key, value) {
  try {
    localStorage.setItem(key, value);
    return true;
  } catch (e) {}
  return false;
}

function jumpToBottom(id) {
  const div = document.querySelector(id);
  console.log(div.scrollTop, div.scrollHeight , div.clientHeight)
  div.scrollTo({
    top: div.scrollHeight - div.clientHeight,
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