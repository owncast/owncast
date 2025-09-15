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
export function trimUnicodeWhitespace(str: string): string {
  // Unicode whitespace regex that matches what Go's strings.TrimSpace() removes
  // Using multiple smaller patterns to avoid ESLint character class warnings
  const patterns = [
    /^[\s\u00A0]+|[\s\u00A0]+$/g, // ASCII whitespace + non-breaking space
    /^[\u1680\u180E]+|[\u1680\u180E]+$/g, // Ogham space mark + Mongolian vowel separator
    /^[\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200A]+|[\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200A]+$/g, // En quad through hair space
    /^[\u200B\u200C\u200D]+|[\u200B\u200C\u200D]+$/g, // Zero-width spaces
    /^[\u2028\u2029]+|[\u2028\u2029]+$/g, // Line separator + paragraph separator
    /^[\u202F\u205F\u3000\uFEFF]+|[\u202F\u205F\u3000\uFEFF]+$/g, // Other Unicode spaces
  ];

  let result = str;
  patterns.forEach(pattern => {
    result = result.replace(pattern, '');
  });

  return result;
}

export function validateDisplayName(
  name: string | undefined,
  currentName: string,
  characterLimit: number = 30,
): DisplayNameValidationResult {
  // Check if name is provided
  if (name === undefined || name === null) {
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
