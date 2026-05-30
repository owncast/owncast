// End-to-end: profanity-filter (built from the SDK) filters chat through its
// filterChatMessage handler in a running Owncast instance.
const { test, expect, beforeAll } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const { registerChat, sendChatMessage } = require('./lib/chat');
const { enableOnly, listPlugins, sleep, ADMIN } = require('./lib/plugins');

const PLUGIN = 'profanity-filter';

beforeAll(async () => {
	await enableOnly(PLUGIN);
});

test('the plugin is discovered, enabled, and loaded', async () => {
	// Look up by slug: that's the stable identifier the admin API keys on.
	// The display `name` is now a free-form label ("Profanity Filter").
	const entry = (await listPlugins()).find((p) => p.slug === PLUGIN);
	expect(entry).toBeDefined();
	expect(entry.enabled).toBe(true);
	expect(entry.loaded).toBe(true);
});

test('redacts profanity from chat messages', async () => {
	const { accessToken } = await registerChat();
	sendChatMessage({ body: 'what the hell', type: 'CHAT' }, accessToken);
	await sleep(2000);

	const messages = (
		await request
			.get('/api/admin/chat/messages')
			.auth(...ADMIN)
			.expect(200)
	).body;
	const message = messages.find(
		(m) => typeof m.body === 'string' && m.body.includes('what the'),
	);

	expect(message).toBeDefined();
	// "hell" is redacted by the plugin; the body is the rendered (HTML) message.
	expect(message.body).toBe('<p>what the ****</p>');
});

test('leaves clean messages untouched', async () => {
	const { accessToken } = await registerChat();
	sendChatMessage({ body: 'hello everyone', type: 'CHAT' }, accessToken);
	await sleep(2000);

	const messages = (
		await request
			.get('/api/admin/chat/messages')
			.auth(...ADMIN)
			.expect(200)
	).body;
	const message = messages.find(
		(m) => typeof m.body === 'string' && m.body.includes('hello everyone'),
	);

	expect(message).toBeDefined();
	expect(message.body).toBe('<p>hello everyone</p>');
});
