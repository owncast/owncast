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
  return !obj || Object.keys(obj).length === 0;
}
