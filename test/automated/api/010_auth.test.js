var request = require('supertest');
request = request('http://127.0.0.1:8080');

test('main page requires no auth', async () => {
	await request.get('/').expect(200);
});

test('admin without trailing slash redirects', async () => {
	await request.get('/admin').expect(301);
});

test('admin with trailing slash requires auth', async () => {
	await request.get('/admin/').expect(401);
});

const paths = [
	'/admin/config/general/',
	'/admin/config/server/',
	'/admin/config-video',
	'/admin/config-chat/',
	'/admin/config-federation/',
	'/admin/config-notify',
	'/admin/federation/followers/',
	'/admin/chat/messages',
	'/admin/viewer-info/',
	'/admin/chat/users/',
	'/admin/stream-health',
	'/admin/hardware-info/',
	// Some APIs too
	'/api/admin/status',
	'/api/admin/serverconfig',
	'/api/admin/chat/clients',
	'/api/admin/chat/messages',
	'/api/admin/followers',
	'/api/admin/prometheus',
];

// Test a bunch of paths to make sure random different pages don't slip by for some reason.
// Technically this shouldn't be possible but it's a sanity check anyway.
paths.forEach((path) => {
	test(`admin path ${path} requires auth and should fail`, async () => {
		await request.get(path).expect(401);
	});
});

// Try them again with auth. Some with trailing slashes some without.
// Allow redirects.
paths.forEach((path) => {
	test(`admin path ${path} requires auth and should pass`, async () => {
		const r = await request.get(path).auth('admin', 'abc123');
		expect([200, 301]).toContain(r.status);
	});
});
