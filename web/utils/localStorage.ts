export const LOCAL_STORAGE_KEYS = {
  username: 'username',
  hasDisplayedNotificationModal: 'HAS_DISPLAYED_NOTIFICATION_MODAL',
  userVisitCount: 'USER_VISIT_COUNT',
};

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
        behavior,
      });
    },
    50,
    element,
  );
}
