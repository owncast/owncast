export function formatIPAddress(ipAddress: string): string {
  const ipAddressComponents = ipAddress.split(':')

  // Wipe out the port component
  ipAddressComponents[ipAddressComponents.length - 1] = '';

  let ip = ipAddressComponents.join(':')
  ip = ip.slice(0, ip.length - 1)
  if (ip === '[::1]') {
    return "Localhost"
  }

  return ip;
}