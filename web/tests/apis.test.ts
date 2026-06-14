import { extractAPIErrorMessage, fetchData } from '../utils/apis';

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

describe('fetchData', () => {
  const originalFetch = global.fetch;

  afterEach(() => {
    global.fetch = originalFetch;
    jest.restoreAllMocks();
  });

  test('throws when a successful response contains invalid json', async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: true,
      status: 200,
      text: async () => '<html>not json</html>',
    } as Response);

    await expect(fetchData('/api/admin/plugins/registry/install')).rejects.toThrow(
      'Invalid JSON response from server',
    );
  });
});
