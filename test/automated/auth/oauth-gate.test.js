// End-to-end OAuth viewer-authentication gate test. Runs against a live Owncast
// instance (started by run.sh) with the `oauth-gate-test` plugin installed and a
// real OpenID Connect provider (node-oidc-provider, started by run.sh) on :9876.
//
// It proves that with the gate enabled: viewer resources are blocked until you
// authenticate, the admin surface stays reachable throughout, completing the
// real OAuth authorization-code flow issues a session that unlocks everything,
// and the resulting identity is recorded on the user.
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const ADMIN = ['admin', 'abc123'];
const SLUG = 'oauth-gate-test';
const SESSION_COOKIE = 'owncast_session';
const LOGIN_PREFIX = `/plugins/${SLUG}/`;
const CALLBACK_PREFIX = `http://localhost:8080/plugins/${SLUG}/callback`;

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

function sessionCookieFrom(res) {
	const setCookie = res.headers['set-cookie'] || [];
	const match = setCookie.find((c) => c.startsWith(`${SESSION_COOKIE}=`));
	return match ? match.split(';')[0] : null;
}

// A cookie-aware redirect follower — the "browser" for the OAuth provider leg.
// node-oidc-provider drives login + consent as a chain of redirects backed by
// session cookies, so plain fetch (no cookie jar) can't complete it. We stop as
// soon as the provider bounces back to the plugin's callback and let supertest
// drive that final hop, so we can observe the Set-Cookie Owncast issues.
async function followToCallback(startUrl) {
	const jar = {};
	let url = startUrl;
	for (let i = 0; i < 25; i++) {
		const cookie = Object.entries(jar)
			.map(([k, v]) => `${k}=${v}`)
			.join('; ');
		const res = await fetch(url, {
			redirect: 'manual',
			headers: cookie ? { cookie } : {},
		});
		for (const c of res.headers.getSetCookie()) {
			const pair = c.split(';')[0];
			const idx = pair.indexOf('=');
			const name = pair.slice(0, idx).trim();
			const val = pair.slice(idx + 1).trim();
			if (val === '') delete jar[name];
			else jar[name] = val;
		}
		const loc = res.headers.get('location');
		if (res.status >= 300 && res.status < 400 && loc) {
			url = new URL(loc, url).toString();
			if (url.startsWith(CALLBACK_PREFIX)) return url;
			continue;
		}
		throw new Error(
			`OAuth provider flow stopped unexpectedly: status=${res.status} url=${url}`,
		);
	}
	throw new Error('OAuth provider flow exceeded redirect limit');
}

beforeAll(async () => {
	await request
		.post(`/api/admin/plugins/${SLUG}/enable`)
		.auth(...ADMIN)
		.expect(200);
	await sleep(1500);
});

afterAll(async () => {
	await request.post(`/api/admin/plugins/${SLUG}/disable`).auth(...ADMIN);
});

describe('before authenticating', () => {
	test('a viewer API endpoint is blocked and redirected to the login screen', async () => {
		const res = await request.get('/api/config').redirects(0).expect(302);
		expect(res.headers.location).toContain(LOGIN_PREFIX);
	});

	test('the home page is blocked and redirected to the login screen', async () => {
		const res = await request.get('/').redirects(0).expect(302);
		expect(res.headers.location).toContain(LOGIN_PREFIX);
	});

	test('the HLS stream playlist is blocked', async () => {
		const res = await request.get('/hls/stream.m3u8').redirects(0);
		expect(res.status).toBe(302);
		expect(res.headers.location).toContain(LOGIN_PREFIX);
	});

	test("the gate plugin's own login screen stays reachable (bootstrap)", async () => {
		await request.get(LOGIN_PREFIX).redirects(0).expect(200);
	});

	test('the admin API stays reachable with admin credentials', async () => {
		await request
			.get('/api/admin/plugins')
			.auth(...ADMIN)
			.redirects(0)
			.expect(200);
	});

	test('the admin web app stays reachable with admin credentials', async () => {
		await request
			.get('/admin/')
			.auth(...ADMIN)
			.redirects(0)
			.expect(200);
	});
});

describe('authenticating via the OAuth provider', () => {
	let sessionCookie;

	test('completing the real OAuth authorization-code flow issues a signed session', async () => {
		// 1. The plugin's login screen links to the provider's authorize endpoint.
		const login = await request.get(LOGIN_PREFIX).redirects(0).expect(200);
		const m = login.text.match(/href="([^"]+\/authorize[^"]*)"/);
		expect(m).toBeTruthy();
		const authorizeUrl = m[1].replace(/&amp;/g, '&');

		// 2. Drive the real provider: authorize -> (login + consent) -> redirect
		//    back to the plugin's callback with a real authorization code.
		const callbackUrl = await followToCallback(authorizeUrl);
		const cb = new URL(callbackUrl);
		expect(cb.searchParams.get('code')).toBeTruthy();

		// 3. Owncast's callback: the plugin exchanges the code for a token, fetches
		//    userinfo, registers the identity, and the host grants the session.
		const res = await request
			.get(cb.pathname + cb.search)
			.redirects(0)
			.expect(302);
		sessionCookie = sessionCookieFrom(res);
		expect(sessionCookie).toBeTruthy();
		// The signed token is "<payload>.<sig>".
		expect(sessionCookie.slice(SESSION_COOKIE.length + 1)).toContain('.');
	});

	test('the session cookie unlocks the viewer API', async () => {
		await request
			.get('/api/config')
			.set('Cookie', sessionCookie)
			.redirects(0)
			.expect(200);
	});

	test('the session cookie unlocks the home page', async () => {
		await request
			.get('/')
			.set('Cookie', sessionCookie)
			.redirects(0)
			.expect(200);
	});

	test('a tampered cookie is rejected and redirected to login', async () => {
		const value = sessionCookie.slice(SESSION_COOKIE.length + 1);
		const flipped = value.slice(0, -1) + (value.endsWith('A') ? 'B' : 'A');
		const res = await request
			.get('/api/config')
			.set('Cookie', `${SESSION_COOKIE}=${flipped}`)
			.redirects(0)
			.expect(302);
		expect(res.headers.location).toContain(LOGIN_PREFIX);
	});

	test('the OAuth identity is recorded on the user (admin users API)', async () => {
		const res = await request
			.get('/api/admin/users?limit=200')
			.auth(...ADMIN)
			.redirects(0)
			.expect(200);
		const match = res.body.results.filter(
			(u) => (u.authProviders || []).includes(SLUG) && u.authenticated,
		);
		expect(match.length).toBeGreaterThanOrEqual(1);
	});

	test('logout clears the session cookie', async () => {
		const res = await request
			.get(`${LOGIN_PREFIX}logout`)
			.set('Cookie', sessionCookie)
			.redirects(0)
			.expect(302);
		const setCookie = res.headers['set-cookie'] || [];
		const cleared = setCookie.find((c) => c.startsWith(`${SESSION_COOKIE}=`));
		expect(cleared).toBeTruthy();
		expect(cleared).toMatch(/owncast_session=;|Max-Age=0/i);
	});
});

describe('after the gate is disabled', () => {
	test('viewer endpoints are reachable again without a session', async () => {
		await request
			.post(`/api/admin/plugins/${SLUG}/disable`)
			.auth(...ADMIN)
			.expect(200);
		await sleep(1000);
		await request.get('/api/config').redirects(0).expect(200);
	});
});
