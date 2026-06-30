import {
  parsePendingFediverseAuth,
  PENDING_FEDIVERSE_AUTH_EXPIRY_MS,
} from '../utils/fediverseAuthSession';

describe('parsePendingFediverseAuth', () => {
  const now = 1_700_000_000_000;
  const account = 'streamer@example.com';

  test('returns the account for a fresh entry', () => {
    const raw = JSON.stringify({ account, ts: now });
    expect(parsePendingFediverseAuth(raw, now)).toEqual({ account, ts: now });
  });

  test('keeps an entry right at the expiry boundary', () => {
    const raw = JSON.stringify({ account, ts: now - PENDING_FEDIVERSE_AUTH_EXPIRY_MS });
    expect(parsePendingFediverseAuth(raw, now)).toEqual({
      account,
      ts: now - PENDING_FEDIVERSE_AUTH_EXPIRY_MS,
    });
  });

  test('drops an entry past the OTP lifetime so the viewer is not stranded', () => {
    const raw = JSON.stringify({ account, ts: now - PENDING_FEDIVERSE_AUTH_EXPIRY_MS - 1 });
    expect(parsePendingFediverseAuth(raw, now)).toBeNull();
  });

  test('returns null for missing or malformed data', () => {
    expect(parsePendingFediverseAuth(null, now)).toBeNull();
    expect(parsePendingFediverseAuth('not json', now)).toBeNull();
    expect(parsePendingFediverseAuth(JSON.stringify({ ts: now }), now)).toBeNull();
    expect(parsePendingFediverseAuth(JSON.stringify({ account }), now)).toBeNull();
  });
});
