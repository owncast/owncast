import { extractAPIErrorMessage } from '../utils/apis';

describe('extractAPIErrorMessage', () => {
  test('prefers backend error field when present', () => {
    expect(
      extractAPIErrorMessage(400, {
        error: 'manifest.scripts requires the "http.serve" permission',
        message: 'generic message',
      }),
    ).toBe('manifest.scripts requires the "http.serve" permission');
  });

  test('falls back to backend message when error is absent', () => {
    expect(extractAPIErrorMessage(502, { message: 'registry unavailable' })).toBe(
      'registry unavailable',
    );
  });

  test('falls back to response text when body is not structured json', () => {
    expect(extractAPIErrorMessage(500, null, 'plain text failure')).toBe('plain text failure');
  });

  test('falls back to generic status when no detail exists', () => {
    expect(extractAPIErrorMessage(418, null, '')).toBe('An error has occurred: 418');
  });
});
