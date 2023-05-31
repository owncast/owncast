import { generateRndKey } from '../components/admin/config/server/StreamKeys';

describe('generateRndKey', () => {
  test('should generate a key that matches the regular expression', () => {
    const key = generateRndKey();
    const regex = /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[!@#$^&*]).{8,192}$/;
    expect(regex.test(key)).toBe(true);
  });

  test('returns a string', () => {
    const result = generateRndKey();
    expect(typeof result).toBe('string');
  });

  test('should generate a key of length between 8 and 192 characters', () => {
    const key = generateRndKey();
    expect(key.length).toBeGreaterThanOrEqual(8);
    expect(key.length).toBeLessThanOrEqual(192);
  });

  test('should generate a unique key on each invocation', () => {
    const key1 = generateRndKey();
    const key2 = generateRndKey();
    expect(key1).not.toBe(key2);
  });
});
