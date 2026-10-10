// localStorage can be null (a WebView with DOM storage disabled) or throw on
// access (blocked storage). next-export-i18n reads it unguarded on every page,
// so an in-memory stand-in takes its place there.

function isUsable(win: Window): boolean {
  try {
    const storage = win.localStorage;
    if (!storage) {
      return false;
    }
    // Read, don't write: a write can fail on a full storage that still works.
    storage.getItem('owncast-storage-probe');
    return true;
  } catch {
    return false;
  }
}

function createMemoryStorage(): Storage {
  const items = new Map<string, string>();
  return {
    get length() {
      return items.size;
    },
    clear: () => items.clear(),
    getItem: key => (items.has(key) ? items.get(key) : null),
    key: index => Array.from(items.keys())[index] ?? null,
    removeItem: key => {
      items.delete(key);
    },
    setItem: (key, value) => {
      items.set(key, String(value));
    },
  };
}

export function ensureLocalStorage(win: Window): boolean {
  if (isUsable(win)) {
    return false;
  }
  try {
    Object.defineProperty(win, 'localStorage', {
      configurable: true,
      value: createMemoryStorage(),
    });
    return true;
  } catch (e) {
    console.error(e);
    return false;
  }
}
