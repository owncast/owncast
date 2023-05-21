export const LOCAL_STORAGE_KEYS = {
  username: 'username',
  hasDisplayedNotificationModal: 'HAS_DISPLAYED_NOTIFICATION_MODAL',
  userVisitCount: 'USER_VISIT_COUNT',
};

export function getLocalStorage(key) {
  try {
    return localStorage.getItem(key);
  } catch (e) {
    console.error(e);
  }
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
  } catch (e) {
    console.error(e);
  }
  return false;
}
