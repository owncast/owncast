/**
 * Validates a display name for chat usage.
 * @param name The proposed display name
 * @param currentName The user's current display name (to check if it's different)
 * @param characterLimit Maximum allowed character count
 * @returns Object with validation result and error message if invalid
 */
export interface DisplayNameValidationResult {
  isValid: boolean;
  errorMessage?: string;
  trimmedName?: string;
}

/**
 * Trims Unicode whitespace characters, similar to Go's strings.TrimSpace()
 * This includes ASCII whitespace plus Unicode space characters and invisible characters
 */
// Unicode whitespace character class used for trimming (matches Go's strings.TrimSpace)
const UNICODE_WHITESPACE_CLASS =
  '[\\s\\u00A0\\u1680\\u180E\\u2000-\\u200A\\u200B-\\u200D\\u2028\\u2029\\u202F\\u205F\\u3000\\uFEFF]';

export function trimUnicodeWhitespace(str: string): string {
  // Unicode whitespace regex that matches what Go's strings.TrimSpace() removes
  // This pattern matches all relevant Unicode whitespace at start/end of string
  const unicodeWhitespacePattern = new RegExp(
    `^${UNICODE_WHITESPACE_CLASS}+|${UNICODE_WHITESPACE_CLASS}+$`,
    'g',
  );
  return str.replace(unicodeWhitespacePattern, '');
}

export function validateDisplayName(
  name: string | undefined,
  currentName: string,
  characterLimit: number = 30,
): DisplayNameValidationResult {
  // Check if name is provided
  if (name === undefined) {
    return {
      isValid: false,
      errorMessage: 'Display name is required',
    };
  }

  // Trim Unicode whitespace (similar to Go's strings.TrimSpace)
  const trimmedName = trimUnicodeWhitespace(name);

  // Check if trimmed name is empty (was only whitespace or originally empty)
  if (trimmedName.length === 0) {
    return {
      isValid: false,
      errorMessage: 'Display name cannot be empty or contain only whitespace',
    };
  }

  // Check if name is different from current
  if (trimmedName === currentName) {
    return {
      isValid: false,
      errorMessage: 'New name must be different from current name',
    };
  }

  // Check character limit (using Unicode-aware length)
  const characterCount = Array.from(trimmedName).length;
  if (characterCount > characterLimit) {
    return {
      isValid: false,
      errorMessage: `Display name cannot exceed ${characterLimit} characters`,
    };
  }

  return {
    isValid: true,
    trimmedName,
  };
}
