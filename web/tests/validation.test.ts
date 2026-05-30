import { isValidUrl, isValidAccount, isValidFediverseAccount } from '../utils/validators';

describe('test url validation', () => {
  const validURL = 'https://example.com';
  const invalidURL = 'example.jfks';

  test('should succeed', () => {
    expect(isValidUrl(validURL)).toBe(true);
  });

  test('should fail', () => {
    expect(isValidUrl(invalidURL)).toBe(false);
  });
});

describe('test xmpp account validation', () => {
  const validAccount = 'xmpp:something@test.biz';
  const invalidAccount = 'something.invalid@something';

  test('should succeed', () => {
    expect(isValidAccount(validAccount, 'xmpp')).toBe(true);
  });

  test('should fail', () => {
    expect(isValidAccount(invalidAccount, 'xmpp')).toBe(false);
  });
});

describe('test fediverse account validation', () => {
  test('should accept a standard account', () => {
    expect(isValidFediverseAccount('@streamer@example.com')).toBe(true);
  });

  test('should accept a punycode TLD', () => {
    expect(isValidFediverseAccount('retrots3m@live.retrospection.xn--q9jyb4c')).toBe(true);
  });

  test('should accept a unicode TLD', () => {
    expect(isValidFediverseAccount('retrots3m@live.retrospection.みんな')).toBe(true);
  });

  test('should accept unicode domain labels', () => {
    expect(isValidFediverseAccount('person@bière.be')).toBe(true);
    expect(isValidFediverseAccount('person@みんな.みんな')).toBe(true);
  });

  test('should reject malformed accounts', () => {
    expect(isValidFediverseAccount('retrots3m')).toBe(false);
    expect(isValidFediverseAccount('retrots3m@')).toBe(false);
    expect(isValidFediverseAccount('@live.retrospection.みんな')).toBe(false);
    expect(isValidFediverseAccount('retrots3m@live.retrospection.みんな/path')).toBe(false);
    expect(isValidFediverseAccount('retrots3m@localhost')).toBe(false);
    expect(isValidFediverseAccount('retrots3m@owncast')).toBe(false);
  });
});
