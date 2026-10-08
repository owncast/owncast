import { ensureLocalStorage } from '../utils/ensureLocalStorage';

// A plain object stands in for window, so jsdom's own localStorage stays intact.
const fakeWindow = (descriptor: PropertyDescriptor): Window => {
  const win = {} as Window;
  Object.defineProperty(win, 'localStorage', { configurable: true, ...descriptor });
  return win;
};

describe('ensureLocalStorage', () => {
  test('leaves a working localStorage alone', () => {
    window.localStorage.setItem('kept', 'yes');
    expect(ensureLocalStorage(window)).toBe(false);
    expect(window.localStorage.getItem('kept')).toBe('yes');
    window.localStorage.removeItem('kept');
  });

  test('replaces a null localStorage, as a WebView with DOM storage disabled exposes it', () => {
    const win = fakeWindow({ value: null });
    expect(ensureLocalStorage(win)).toBe(true);
    win.localStorage.setItem('lang', 'nl');
    expect(win.localStorage.getItem('lang')).toBe('nl');
  });

  test('replaces a localStorage that throws, as blocked storage does', () => {
    const win = fakeWindow({
      get() {
        throw new DOMException('Access is denied for this document.', 'SecurityError');
      },
    });
    expect(ensureLocalStorage(win)).toBe(true);
    expect(win.localStorage.getItem('missing')).toBeNull();
  });

  test('the stand-in behaves like Storage', () => {
    const win = fakeWindow({ value: null });
    ensureLocalStorage(win);
    const storage = win.localStorage;
    storage.setItem('a', 1 as unknown as string);
    storage.setItem('b', 'two');
    expect(storage.getItem('a')).toBe('1');
    expect(storage.length).toBe(2);
    expect(storage.key(1)).toBe('b');
    expect(storage.key(5)).toBeNull();
    storage.removeItem('a');
    expect(storage.getItem('a')).toBeNull();
    storage.clear();
    expect(storage.length).toBe(0);
  });
});
