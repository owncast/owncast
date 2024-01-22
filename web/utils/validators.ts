// to use with <input type="url"> fields, as the default pattern only checks for `:`,
export const DEFAULT_TEXTFIELD_URL_PATTERN = 'https?://.*';

/**
 * Determines if a URL is valid
 * @param {string} url - A URL to validate.
 * @param {string[]} validProtocols - An array of valid protocols. Defaults to web.
 * @returns {boolean} - True if the URI is valid, false otherwise.
 */
export function isValidUrl(url: string, validProtocols: string[] = ['http:', 'https:']): boolean {
  try {
    const validationObject = new URL(url);

    if (
      validationObject.protocol === '' ||
      validationObject.hostname === '' ||
      !validProtocols.includes(validationObject.protocol)
    ) {
      return false;
    }
  } catch (e) {
    return false;
  }

  return true;
}

/**
 * Determines if an account is valid by checking for a protocol, username
 * and server, delimited by a colon. For example: @username:example.com
 * @param {string} account - An account to validate.
 * @param {string} protocol - The protocol we expect the account to be using.
 * @returns {boolean} - True if the account is valid, false otherwise.
 */
export function isValidAccount(account: string, protocol: string): boolean {
  if (account.startsWith('@')) {
    // eslint-disable-next-line no-param-reassign
    account = account.slice(1);
  }

  const components = account.split(/:|@/);
  const [service, user, host] = components;

  if (service !== protocol) {
    return false;
  }

  if (components.length !== 3 || !service || !user || !host) {
    return false;
  }

  return true;
}

/**
 * Determines if a fediverse account is valid.
 * For example: @username@example.com
 * @param {string} account - An account to validate.
 * @returns {boolean} - True if the account is valid, false otherwise.
 */
export function isValidFediverseAccount(account: string): boolean {
  const sanitized = account.replace(/^@+/, '');
  const regex =
    /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return regex.test(String(sanitized).toLowerCase());
}
