import UAParser from 'ua-parser-js';

export function getDiffInDaysFromNow(timestamp) {
  const time = typeof timestamp === 'string' ? new Date(timestamp) : timestamp;
  return (new Date() - time) / (24 * 3600 * 1000);
}

export const isMobileSafariIos = () => {
  try {
    const ua = navigator.userAgent;
    const uaParser = new UAParser(ua);
    const browser = uaParser.getBrowser();
    const device = uaParser.getDevice();

    if (device.vendor !== 'Apple') {
      return false;
    }

    if (device.type !== 'mobile' && device.type !== 'tablet') {
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
