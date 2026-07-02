// Persists an in-progress fediverse OTP verification across a page reload so a
// viewer who backgrounds the tab to copy their code (issue #4887) returns to the
// code-entry step instead of an empty account form.
//
// The window mirrors the backend OTP lifetime (auth/fediverse registrationTimeout
// is 10 minutes). Past that the code the viewer received is dead, so we drop the
// persisted state and send them back to the account step rather than stranding
// them at a code box that can never succeed.

const PENDING_FEDIVERSE_AUTH_KEY = 'pendingFediverseAuth';
export const PENDING_FEDIVERSE_AUTH_EXPIRY_MS = 10 * 60 * 1000;

export type PendingFediverseAuth = {
  account: string;
  ts: number;
};

// Pure decode + expiry check, kept separate from storage I/O so it can be tested
// with a plain string and clock value.
export function parsePendingFediverseAuth(
  raw: string | null,
  now: number,
): PendingFediverseAuth | null {
  if (!raw) {
    return null;
  }

  let parsed: PendingFediverseAuth;
  try {
    parsed = JSON.parse(raw) as PendingFediverseAuth;
  } catch {
    return null;
  }

  if (!parsed || typeof parsed.account !== 'string' || typeof parsed.ts !== 'number') {
    return null;
  }

  // Reject non-finite or future timestamps so a bad value can't pin the modal
  // open forever (Infinity or a future ts never crosses the expiry check).
  if (!Number.isFinite(parsed.ts) || parsed.ts > now) {
    return null;
  }

  if (now - parsed.ts > PENDING_FEDIVERSE_AUTH_EXPIRY_MS) {
    return null;
  }

  return parsed;
}

export function getPendingFediverseAuth(): PendingFediverseAuth | null {
  try {
    const result = parsePendingFediverseAuth(
      sessionStorage.getItem(PENDING_FEDIVERSE_AUTH_KEY),
      Date.now(),
    );
    if (!result) {
      sessionStorage.removeItem(PENDING_FEDIVERSE_AUTH_KEY);
    }
    return result;
  } catch (e) {
    console.error(e);
    return null;
  }
}

export function setPendingFediverseAuth(account: string): void {
  try {
    const value: PendingFediverseAuth = { account, ts: Date.now() };
    sessionStorage.setItem(PENDING_FEDIVERSE_AUTH_KEY, JSON.stringify(value));
  } catch (e) {
    console.error(e);
  }
}

export function clearPendingFediverseAuth(): void {
  try {
    sessionStorage.removeItem(PENDING_FEDIVERSE_AUTH_KEY);
  } catch (e) {
    console.error(e);
  }
}
