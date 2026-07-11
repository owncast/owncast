import Bowser from 'bowser';

export function getDiffInDaysFromNow(timestamp) {
  const time = typeof timestamp === 'string' ? new Date(timestamp) : timestamp;
  return (new Date() - time) / (24 * 3600 * 1000);
}

export const isMobileSafariIos = () => {
  try {
    const { browser, platform } = Bowser.parse(navigator.userAgent);
    const isIPadOS = navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1;

    if (!isIPadOS && platform.vendor !== 'Apple') {
      return false;
    }

    if (!isIPadOS && platform.type !== 'mobile' && platform.type !== 'tablet') {
      return false;
    }

    return browser.name === 'Mobile Safari' || browser.name === 'Safari';
  } catch {
    return false;
  }
};

export const isMobileSafariHomeScreenApp = () => {
  if (!isMobileSafariIos()) {
    return false;
  }

  return 'standalone' in window.navigator && window.navigator.standalone;
};
