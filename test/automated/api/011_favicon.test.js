var request = require('supertest');

request = request('http://127.0.0.1:8080');

const defaultAdminPassword = 'abc123';

// Create base64 data URLs for test images.
// Minimal valid PNG (1x1 pixel).
const minimalPNG = Buffer.from([
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49,
	0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x02,
	0x00, 0x00, 0x00, 0x90, 0x77, 0x53, 0xde, 0x00, 0x00, 0x00, 0x0c, 0x49, 0x44,
	0x41, 0x54, 0x08, 0xd7, 0x63, 0xf8, 0xff, 0xff, 0x3f, 0x00, 0x05, 0xfe, 0x02,
	0xfe, 0xdc, 0xcc, 0x59, 0xe7, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44,
	0xae, 0x42, 0x60, 0x82,
]);

// Minimal valid ICO (1x1 pixel).
const minimalICO = Buffer.from([
	0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x18,
	0x00, 0x30, 0x00, 0x00, 0x00, 0x16, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00,
	0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x18, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff,
	0x00, 0x00, 0x00, 0x00, 0x00,
]);

const pngBase64 = `data:image/png;base64,${minimalPNG.toString('base64')}`;
const icoBase64 = `data:image/x-icon;base64,${minimalICO.toString('base64')}`;
const jpegBase64 = `data:image/jpeg;base64,${minimalPNG.toString('base64')}`;

test('upload PNG favicon successfully', async () => {
	const res = await request
		.post('/api/admin/config/favicon')
		.auth('admin', defaultAdminPassword)
		.send({ value: pngBase64 })
		.expect(200);

	expect(res.body.success).toBe(true);
	expect(res.body.message).toBe('favicon updated');
});

test('verify favicon.ico endpoint returns the uploaded favicon', async () => {
	const res = await request.get('/favicon.ico').expect(200);

	expect(res.headers['content-type']).toMatch(/image\/(png|x-icon)/);
	expect(res.body).toBeDefined();
});

test('upload ICO favicon successfully', async () => {
	const res = await request
		.post('/api/admin/config/favicon')
		.auth('admin', defaultAdminPassword)
		.send({ value: icoBase64 })
		.expect(200);

	expect(res.body.success).toBe(true);
	expect(res.body.message).toBe('favicon updated');
});

test('reject unsupported file type (JPEG)', async () => {
	const res = await request
		.post('/api/admin/config/favicon')
		.auth('admin', defaultAdminPassword)
		.send({ value: jpegBase64 })
		.expect(400);

	expect(res.body.success).toBe(false);
	expect(res.body.message).toBe('favicon must be PNG or ICO format');
});

test('reject oversized favicon (>200KB)', async () => {
	// 210KB decoded exceeds the 200KB server limit but stays under the
	// MaxBytesReader limit so the handler's own size check is exercised.
	const oversized = Buffer.alloc(210 * 1024, 0x42);
	const oversizedBase64 = `data:image/png;base64,${oversized.toString('base64')}`;

	const res = await request
		.post('/api/admin/config/favicon')
		.auth('admin', defaultAdminPassword)
		.send({ value: oversizedBase64 })
		.expect(400);

	expect(res.body.success).toBe(false);
	expect(res.body.message).toBe('file too large, max 200KB');
});

test('reject request without image data', async () => {
	const res = await request
		.post('/api/admin/config/favicon')
		.auth('admin', defaultAdminPassword)
		.send({})
		.expect(400);

	expect(res.body.success).toBe(false);
});

test('reject unauthenticated request', async () => {
	await request.post('/api/admin/config/favicon').expect(401);
});

test('reset favicon to default successfully', async () => {
	// First upload a favicon to ensure we have one to reset
	await request
		.post('/api/admin/config/favicon')
		.auth('admin', defaultAdminPassword)
		.send({ value: pngBase64 })
		.expect(200);

	// Now reset to default
	const res = await request
		.delete('/api/admin/config/favicon')
		.auth('admin', defaultAdminPassword)
		.expect(200);

	expect(res.body.success).toBe(true);
	expect(res.body.message).toBe('favicon reset to default');
});

test('verify favicon.ico returns default after reset', async () => {
	const res = await request.get('/favicon.ico').expect(200);

	expect(res.headers['content-type']).toMatch(/image\/(png|x-icon)/);
	expect(res.body).toBeDefined();
});

test('reject unauthenticated reset request', async () => {
	await request.delete('/api/admin/config/favicon').expect(401);
});
