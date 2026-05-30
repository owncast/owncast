// End-to-end: overlay (built from the SDK) serves a static page and a dynamic
// JSON endpoint backed by the owncast.chat.history() host function, exercising
// the plugin HTTP-serving path in a running Owncast instance.
const { test, expect, beforeAll } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const { registerChat, sendChatMessage } = require('./lib/chat');
const { enableOnly, sleep } = require('./lib/plugins');

const PLUGIN = 'overlay';

beforeAll(async () => {
	await enableOnly(PLUGIN);
});

test('serves its static overlay page at /plugins/overlay/', async () => {
	const res = await request.get('/plugins/overlay/').expect(200);
	expect(res.headers['content-type']).toMatch(/html/);
	expect(res.text.length).toBeGreaterThan(0);
});

test('serves a dynamic API backed by chat.history()', async () => {
	// Seed a chat message, then read it back through the plugin's own endpoint.
	const { accessToken } = await registerChat();
	const probe = 'overlay-probe-' + Date.now();
	sendChatMessage({ body: probe, type: 'CHAT' }, accessToken);
	await sleep(2000);

	const res = await request.get('/plugins/overlay/api/messages').expect(200);
	expect(res.headers['content-type']).toMatch(/json/);

	const messages = res.body.messages || [];
	const found = messages.some(
		(m) => typeof m.body === 'string' && m.body.includes(probe),
	);
	expect(found).toBe(true);
});
