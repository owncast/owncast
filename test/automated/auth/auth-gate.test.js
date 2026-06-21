// End-to-end viewer-authentication gate test. Runs against a live Owncast
// instance (started by run.sh) with the SDK's `basic-auth` example installed.
// It enables the gate, proves endpoints are blocked until you authenticate,
// then proves the issued session cookie unlocks them.
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const ADMIN = ['admin', 'abc123'];
const SLUG = 'basic-auth';
const PASSWORD = 'letmein'; // basic-auth manifest default
const SESSION_COOKIE = 'owncast_session';
const LOGIN_PREFIX = `/plugins/${SLUG}/`;

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

// Pull the "owncast_session=<value>" pair out of a response's Set-Cookie.
function sessionCookieFrom(res) {
	const setCookie = res.headers['set-cookie'] || [];
	const match = setCookie.find((c) => c.startsWith(`${SESSION_COOKIE}=`));
	return match ? match.split(';')[0] : null;
}

beforeAll(async () => {
	// Enabling returns 200 only once the plugin actually loads.
	await request
		.post(`/api/admin/plugins/${SLUG}/enable`)
		.auth(...ADMIN)
		.expect(200);
	await sleep(1500);
});

afterAll(async () => {
	// Always lift the gate so it can't leak into other suites / a left-on state.
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

	test('the gate plugin\'s own login screen stays reachable (bootstrap)', async () => {
		await request.get(LOGIN_PREFIX).redirects(0).expect(200);
	});

	test('the admin API stays reachable with admin credentials', async () => {
		await request
			.get('/api/admin/plugins')
			.auth(...ADMIN)
			.redirects(0)
			.expect(200);
	});
});

describe('authenticating', () => {
	let sessionCookie;

	test('a wrong password issues no session', async () => {
		const res = await request
			.get(`${LOGIN_PREFIX}login?password=wrong-password`)
			.redirects(0)
			.expect(200);
		expect(sessionCookieFrom(res)).toBeNull();
	});

	test('the correct password issues a signed session cookie', async () => {
		const res = await request
			.get(`${LOGIN_PREFIX}login?password=${PASSWORD}&return_to=%2F`)
			.redirects(0)
			.expect(302);
		sessionCookie = sessionCookieFrom(res);
		expect(sessionCookie).toBeTruthy();
		expect(sessionCookie.startsWith(`${SESSION_COOKIE}=`)).toBe(true);
		// The signed token is "<base64url payload>.<base64url sig>".
		const value = sessionCookie.slice(SESSION_COOKIE.length + 1);
		expect(value).toContain('.');
	});

	test('the session cookie unlocks the viewer API', async () => {
		await request
			.get('/api/config')
			.set('Cookie', sessionCookie)
			.redirects(0)
			.expect(200);
	});

	test('the session cookie unlocks the home page', async () => {
		await request.get('/').set('Cookie', sessionCookie).redirects(0).expect(200);
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

	test('logout clears the session cookie', async () => {
		const res = await request
			.get(`${LOGIN_PREFIX}logout`)
			.set('Cookie', sessionCookie)
			.redirects(0)
			.expect(302);
		const setCookie = res.headers['set-cookie'] || [];
		const cleared = setCookie.find((c) => c.startsWith(`${SESSION_COOKIE}=`));
		expect(cleared).toBeTruthy();
		// A clearing cookie has an empty value and an immediate expiry.
		expect(cleared).toMatch(/owncast_session=;|Max-Age=0/i);
	});
});

describe('revocation via onAuthCheck', () => {
	let cookie;

	test('a fresh login establishes a working session', async () => {
		const res = await request
			.get(`${LOGIN_PREFIX}login?password=${PASSWORD}&return_to=%2F`)
			.redirects(0)
			.expect(302);
		cookie = sessionCookieFrom(res);
		expect(cookie).toBeTruthy();
		// Index page allows the session — this also proves the host resolves the
		// viewer identity and passes it to onAuthCheck (which returns ok).
		await request.get('/').set('Cookie', cookie).redirects(0).expect(200);
	});

	test('after an admin revokes, the index page denies the session', async () => {
		await request
			.get(`${LOGIN_PREFIX}revoke`)
			.auth(...ADMIN)
			.redirects(0)
			.expect(200);
		const res = await request.get('/').set('Cookie', cookie).redirects(0).expect(302);
		expect(res.headers.location).toContain(LOGIN_PREFIX);
	});

	test('after the revocation is lifted, the index page allows the session again', async () => {
		await request
			.get(`${LOGIN_PREFIX}unrevoke`)
			.auth(...ADMIN)
			.redirects(0)
			.expect(200);
		await request.get('/').set('Cookie', cookie).redirects(0).expect(200);
	});
});

describe('admin config form (auto-generated from the manifest config block)', () => {
	const CONFIG_URL = `/api/admin/plugins/${SLUG}/config`;

	test('the plugin list exposes the config schema', async () => {
		const res = await request.get('/api/admin/plugins').auth(...ADMIN).redirects(0).expect(200);
		const entry = res.body.find((p) => p.slug === SLUG);
		expect(entry).toBeTruthy();
		expect(entry.config).toBeTruthy();
		expect(entry.config.password.type).toBe('string');
		expect(entry.config.password.default).toBe('letmein');
	});

	test('GET returns current effective values (defaults until overridden)', async () => {
		const res = await request.get(CONFIG_URL).auth(...ADMIN).redirects(0).expect(200);
		expect(res.body.password).toBe('letmein');
	});

	test('POST rejects an unknown key', async () => {
		await request
			.post(CONFIG_URL)
			.auth(...ADMIN)
			.send({ nope: 'x' })
			.redirects(0)
			.expect(400);
	});

	test('saving an override persists and is reflected by GET and by config.get()', async () => {
		await request
			.post(CONFIG_URL)
			.auth(...ADMIN)
			.send({ password: 'config-test-pass' })
			.redirects(0)
			.expect(200);

		// GET now returns the override.
		const after = await request.get(CONFIG_URL).auth(...ADMIN).redirects(0).expect(200);
		expect(after.body.password).toBe('config-test-pass');

		// The plugin reads it via owncast.config.get: the NEW password logs in,
		// the old default no longer does — proving the override is wired through.
		const ok = await request
			.get(`${LOGIN_PREFIX}login?password=config-test-pass&return_to=%2F`)
			.redirects(0)
			.expect(302);
		expect(sessionCookieFrom(ok)).toBeTruthy();

		const bad = await request
			.get(`${LOGIN_PREFIX}login?password=letmein`)
			.redirects(0)
			.expect(200); // re-shows the form, no redirect/cookie
		expect(sessionCookieFrom(bad)).toBeNull();
	});

	test('reset the override back to the default', async () => {
		await request
			.post(CONFIG_URL)
			.auth(...ADMIN)
			.send({ password: 'letmein' })
			.redirects(0)
			.expect(200);
		const res = await request.get(CONFIG_URL).auth(...ADMIN).redirects(0).expect(200);
		expect(res.body.password).toBe('letmein');
	});
});

describe('after the gate is disabled', () => {
	test('viewer endpoints are reachable again without a session', async () => {
		await request.post(`/api/admin/plugins/${SLUG}/disable`).auth(...ADMIN).expect(200);
		await sleep(1000);
		await request.get('/api/config').redirects(0).expect(200);
	});
});
