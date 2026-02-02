import { generateRndKey } from '../components/admin/config/server/StreamKeys';
import { REGEX_STREAM_KEY } from '../utils/config-constants';

describe('generateRndKey', () => {
  test('should generate a key that matches the stream key regex', () => {
    const key = generateRndKey();
    // Use the same regex constant that the implementation uses
    expect(REGEX_STREAM_KEY.test(key)).toBe(true);
  });

  test('should only contain alphanumeric characters (no special characters for broadcasting software compatibility)', () => {
    const key = generateRndKey();
    // Keys should be purely alphanumeric for compatibility with broadcasting software
    expect(key).toMatch(/^[a-zA-Z0-9]+$/);
  });

  test('should not contain dashes', () => {
    const key = generateRndKey();
    // Dashes are explicitly forbidden as they can break RTMP URL parsing
    expect(key).not.toContain('-');
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

  test('should contain at least one uppercase letter, one lowercase letter, and one digit', () => {
    const key = generateRndKey();
    expect(key).toMatch(/[A-Z]/); // has uppercase
    expect(key).toMatch(/[a-z]/); // has lowercase
    expect(key).toMatch(/[0-9]/); // has digit
  });

  test('should consistently generate valid keys across multiple invocations', () => {
    // Generate multiple keys to ensure consistency
    for (let i = 0; i < 10; i++) {
      const key = generateRndKey();
      expect(REGEX_STREAM_KEY.test(key)).toBe(true);
      expect(key).toMatch(/^[a-zA-Z0-9]+$/);
    }
  });
});
