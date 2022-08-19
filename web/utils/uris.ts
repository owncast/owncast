export default function isValidURI(uri: string): boolean {
  if (uri.includes(':')) {
    const splittedURI = uri.split(':', 2);
    if (
      splittedURI[0] === ''
    ) {
      return false;
    } else if (
      splittedURI[1] === ''
    ) {
      return false;
    } else if (
      splittedURI[0] === 'http' ||
      splittedURI[0] === 'https'
    ) {
      if (
        (uri.slice(0,7) === 'http://' ||
        uri.slice(0,8) === 'https://') &&
        splittedURI[1].includes('.')
      ) {
        return true;
      }
      return false;
    }
    return true;
  }
  return false;
}
