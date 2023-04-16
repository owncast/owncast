import { isValidUrl, isValidAccount, isValidMatrixAccount } from '../utils/validators';

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

describe('test matrix account validation', () => {
  const validMatrixAccount = '@me:matrix.org';
  const validMatrixAccountWithProtocol = 'matrix:@me:matrix.org';
  const invalidMatrixAccount = 'something.invalid@something';

  test('should succeed', () => {
    expect(isValidMatrixAccount(validMatrixAccount)).toBe(true);
    expect(isValidMatrixAccount(validMatrixAccountWithProtocol)).toBe(true);
  });

  test('should fail', () => {
    expect(isValidMatrixAccount(invalidMatrixAccount)).toBe(false);
  });
});
