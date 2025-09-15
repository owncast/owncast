import { validateDisplayName } from '../utils/displayNameValidation';

describe('Display Name Validation', () => {
  const currentName = 'CurrentUser';
  const characterLimit = 30;

  describe('Valid names', () => {
    test('should accept valid name', () => {
      const result = validateDisplayName('NewUser', currentName, characterLimit);
      expect(result.isValid).toBe(true);
      expect(result.trimmedName).toBe('NewUser');
      expect(result.errorMessage).toBeUndefined();
    });

    test('should trim whitespace and accept valid name', () => {
      const result = validateDisplayName('  NewUser  ', currentName, characterLimit);
      expect(result.isValid).toBe(true);
      expect(result.trimmedName).toBe('NewUser');
    });

    test('should accept name with Unicode characters', () => {
      const result = validateDisplayName('用户名', currentName, characterLimit);
      expect(result.isValid).toBe(true);
      expect(result.trimmedName).toBe('用户名');
    });

    test('should accept name at character limit', () => {
      const longName = 'A'.repeat(characterLimit);
      const result = validateDisplayName(longName, currentName, characterLimit);
      expect(result.isValid).toBe(true);
      expect(result.trimmedName).toBe(longName);
    });
  });

  describe('Invalid names - whitespace and empty', () => {
    test('should reject undefined name', () => {
      const result = validateDisplayName(undefined, currentName, characterLimit);
      expect(result.isValid).toBe(false);
      expect(result.errorMessage).toBe('Display name is required');
    });

    test('should reject empty string', () => {
      const result = validateDisplayName('', currentName, characterLimit);
      expect(result.isValid).toBe(false);
      expect(result.errorMessage).toBe('Display name cannot be empty or contain only whitespace');
    });

    test('should reject spaces only', () => {
      const result = validateDisplayName('   ', currentName, characterLimit);
      expect(result.isValid).toBe(false);
      expect(result.errorMessage).toBe('Display name cannot be empty or contain only whitespace');
    });

    test('should reject tabs only', () => {
      const result = validateDisplayName('\t\t\t', currentName, characterLimit);
      expect(result.isValid).toBe(false);
      expect(result.errorMessage).toBe('Display name cannot be empty or contain only whitespace');
    });

    test('should reject newlines only', () => {
      const result = validateDisplayName('\n\n', currentName, characterLimit);
      expect(result.isValid).toBe(false);
      expect(result.errorMessage).toBe('Display name cannot be empty or contain only whitespace');
    });

    test('should reject mixed whitespace', () => {
      const result = validateDisplayName(' \t\n\r ', currentName, characterLimit);
      expect(result.isValid).toBe(false);
      expect(result.errorMessage).toBe('Display name cannot be empty or contain only whitespace');
    });

    test('should reject Unicode whitespace', () => {
      const result = validateDisplayName('\u00A0\u2000\u2001\u2002', currentName, characterLimit);
      expect(result.isValid).toBe(false);
      expect(result.errorMessage).toBe('Display name cannot be empty or contain only whitespace');
    });
  });

  describe('Invalid names - other reasons', () => {
    test('should reject name same as current', () => {
      const result = validateDisplayName(currentName, currentName, characterLimit);
      expect(result.isValid).toBe(false);
      expect(result.errorMessage).toBe('New name must be different from current name');
    });

    test('should reject name same as current after trimming', () => {
      const result = validateDisplayName(`  ${currentName}  `, currentName, characterLimit);
      expect(result.isValid).toBe(false);
      expect(result.errorMessage).toBe('New name must be different from current name');
    });

    test('should reject name exceeding character limit', () => {
      const longName = 'A'.repeat(characterLimit + 1);
      const result = validateDisplayName(longName, currentName, characterLimit);
      expect(result.isValid).toBe(false);
      expect(result.errorMessage).toBe(`Display name cannot exceed ${characterLimit} characters`);
    });

    test('should handle Unicode characters correctly for length', () => {
      // Emoji characters count as 2 code units but should be treated as 1 character
      const emojiName = '😀'.repeat(characterLimit + 1);
      const result = validateDisplayName(emojiName, currentName, characterLimit);
      expect(result.isValid).toBe(false);
      expect(result.errorMessage).toBe(`Display name cannot exceed ${characterLimit} characters`);
    });
  });

  describe('Edge cases', () => {
    test('should handle zero character limit', () => {
      const result = validateDisplayName('A', currentName, 0);
      expect(result.isValid).toBe(false);
      expect(result.errorMessage).toBe('Display name cannot exceed 0 characters');
    });

    test('should handle very long current name', () => {
      const veryLongCurrentName = 'A'.repeat(100);
      const result = validateDisplayName('NewUser', veryLongCurrentName, characterLimit);
      expect(result.isValid).toBe(true);
      expect(result.trimmedName).toBe('NewUser');
    });

    test('should handle special characters', () => {
      const specialName = '!@#$%^&*()_+{}|:"<>?[]\\;\',./-=`~';
      const result = validateDisplayName(specialName, currentName, 50);
      expect(result.isValid).toBe(true);
      expect(result.trimmedName).toBe(specialName);
    });
  });
});

describe('Display Name Validation - Real-world test cases', () => {
  const currentName = 'TestUser';

  const testCases = [
    { input: 'John', expected: true, description: 'Simple name' },
    { input: 'John Doe', expected: true, description: 'Name with space' },
    { input: '  John  ', expected: true, description: 'Name with padding spaces' },
    { input: '', expected: false, description: 'Empty string' },
    { input: '   ', expected: false, description: 'Only spaces' },
    { input: '\t', expected: false, description: 'Only tab' },
    { input: '\n', expected: false, description: 'Only newline' },
    { input: 'TestUser', expected: false, description: 'Same as current' },
    { input: '  TestUser  ', expected: false, description: 'Same as current with spaces' },
    { input: 'A'.repeat(31), expected: false, description: 'Too long (31 chars)' },
    { input: 'A'.repeat(30), expected: true, description: 'At limit (30 chars)' },
    { input: '用户', expected: true, description: 'Chinese characters' },
    { input: '😀🎉', expected: true, description: 'Emoji characters' },
  ];

  test.each(testCases)('$description', ({ input, expected }) => {
    const result = validateDisplayName(input, currentName, 30);
    expect(result.isValid).toBe(expected);
    if (expected && result.isValid) {
      expect(result.trimmedName).toBe(input.trim());
    }
  });
});
