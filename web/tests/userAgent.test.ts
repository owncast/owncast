import { formatUAstring } from '../utils/format';
import { isMobileSafariHomeScreenApp, isMobileSafariIos } from '../utils/helpers';

const iphoneSafari =
  'Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1';
const iphoneChrome =
  'Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/126.0.6478.153 Mobile/15E148 Safari/604.1';
const ipadSafari =
  'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1';
const desktopSafari =
  'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15';
const desktopChrome =
  'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36';

function setNavigator(
  userAgent: string,
  { platform = '', maxTouchPoints = 0, standalone = false } = {},
) {
  Object.defineProperties(window.navigator, {
    userAgent: { configurable: true, value: userAgent },
    platform: { configurable: true, value: platform },
    maxTouchPoints: { configurable: true, value: maxTouchPoints },
    standalone: { configurable: true, value: standalone },
  });
}

describe('mobile Safari detection', () => {
  test.each([
    ['iPhone Safari', iphoneSafari, 'iPhone', 5, true],
    ['iPhone Chrome', iphoneChrome, 'iPhone', 5, false],
    ['iPadOS Safari', ipadSafari, 'MacIntel', 5, true],
    ['desktop Safari', desktopSafari, 'MacIntel', 0, false],
  ])('%s', (_name, userAgent, platform, maxTouchPoints, expected) => {
    setNavigator(userAgent, { platform, maxTouchPoints });

    expect(isMobileSafariIos()).toBe(expected);
  });

  test('detects an installed iOS home screen app', () => {
    setNavigator(iphoneSafari, { platform: 'iPhone', maxTouchPoints: 5, standalone: true });

    expect(isMobileSafariHomeScreenApp()).toBe(true);
  });
});

describe('user agent formatting', () => {
  test('preserves the existing browser display format', () => {
    expect(formatUAstring(desktopChrome)).toBe('Chrome 126 on Windows 10\n  ');
    expect(formatUAstring(iphoneSafari)).toBe('Mobile Safari 17 on iOS 17.5\n   (iPhone)');
  });

  test('preserves special and unknown user agents', () => {
    expect(formatUAstring('libmpv')).toBe('mpv media player');
    expect(formatUAstring('curl/8.7.1')).toBe('curl/8.7.1');
  });
});
