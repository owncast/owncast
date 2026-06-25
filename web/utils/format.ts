import { format } from 'date-fns';
import UAParser from 'ua-parser-js';

// formatDisplayDate renders a user timestamp, including the year only when it
// is not the current year. Shared by the admin user views.
export function formatDisplayDate(date: string | Date) {
  const d = new Date(date);
  if (d.getFullYear() !== new Date().getFullYear()) {
    return format(new Date(date), 'MMM d, yyyy H:mma');
  }

  return format(new Date(date), 'MMM d H:mma');
}

// formatDateOnly renders just the calendar date (no time of day), e.g.
// "Jun 25, 2026". Used for the admin user table's "Created" column.
export function formatDateOnly(date: string | Date) {
  return format(new Date(date), 'MMM d, yyyy');
}

export function formatIPAddress(ipAddress: string): string {
  const ipAddressComponents = ipAddress.split(':');

  // Wipe out the port component
  ipAddressComponents[ipAddressComponents.length - 1] = '';

  let ip = ipAddressComponents.join(':');
  ip = ip.slice(0, ip.length - 1);
  if (ip === '[::1]' || ip === '127.0.0.1') {
    return 'Localhost';
  }

  return ip;
}

// check if obj is {}
export function isEmptyObject(obj) {
  return !obj || (Object.keys(obj).length === 0 && obj.constructor === Object);
}

function padLeft(text, pad, size) {
  return String(pad.repeat(size) + text).slice(-size);
}

export function parseSecondsToDurationString(seconds = 0) {
  const finiteSeconds = Number.isFinite(+seconds) ? Math.abs(seconds) : 0;

  const days = Math.floor(finiteSeconds / 86400);
  const daysString = days > 0 ? `${days} day${days > 1 ? 's' : ''} ` : '';

  const hours = Math.floor((finiteSeconds / 3600) % 24);
  const hoursString = hours || days ? padLeft(`${hours}:`, '0', 3) : '';

  const mins = Math.floor((finiteSeconds / 60) % 60);
  const minString = padLeft(`${mins}:`, '0', 3);

  const secs = Math.floor(finiteSeconds % 60);
  const secsString = padLeft(`${secs}`, '0', 2);

  return daysString + hoursString + minString + secsString;
}

export function formatUAstring(uaString: string) {
  const parser = UAParser(uaString);
  const { device, os, browser } = parser;
  const { major: browserVersion, name } = browser;
  const { version: osVersion, name: osName } = os;
  const { model, type } = device;

  if (uaString === 'libmpv') {
    return 'mpv media player';
  }
  // Fallback to just displaying the raw agent string.
  if (!name || !browserVersion || !osName) {
    return uaString;
  }

  const deviceString = model || type ? ` (${model || type})` : '';
  return `${name} ${browserVersion} on ${osName} ${osVersion}
  ${deviceString}`;
}
