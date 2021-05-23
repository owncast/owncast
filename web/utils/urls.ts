// to use with <input type="url"> fields, as the default pattern only checks for `:`,
export const DEFAULT_TEXTFIELD_URL_PATTERN = 'https?://.*';

export default function isValidUrl(url: string): boolean {
  const validProtocols = ['http:', 'https:'];

  try {
   const validationObject = new URL(url);
   if (validationObject.protocol === '' || validationObject.hostname === '' || !validProtocols.includes(validationObject.protocol)) {
     return false;
   }
  } catch(e) {
    return false;
  }

  return true;
}
