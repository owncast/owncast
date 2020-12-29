export function formatIPAddress(ipAddress: string): string {
  const ipAddressComponents = ipAddress.split(':')

  // Wipe out the port component
  ipAddressComponents[ipAddressComponents.length - 1] = '';

  let ip = ipAddressComponents.join(':')
  ip = ip.slice(0, ip.length - 1)
  if (ip === '[::1]' || ip === '127.0.0.1') {
    return "Localhost"
  }

  return ip;
}

// check if obj is {}
export function isEmptyObject(obj) {
  return !obj || (Object.keys(obj).length === 0 && obj.constructor === Object);
}

export function padLeft(text, pad, size) {
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
