// End-to-end: a Rust module compiled directly to wasm reads its packaged
// manifest from Extism config during register(), receives a real Owncast chat
// event, and calls the owncast_send_chat host function. No language SDK or
// runtime script is involved.
const { test, expect, beforeAll } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const { registerChat, sendChatMessage } = require('./lib/chat');
const { enableOnly, sleep, ADMIN } = require('./lib/plugins');

const PLUGIN = 'native-wasm-test';
const BOT_DISPLAY_NAME = 'Native Wasm Test';

beforeAll(async () => {
	await enableOnly(PLUGIN);
});

test('loads and handles chat as a self-contained wasm module', async () => {
	const { accessToken } = await registerChat();
	const probe = 'native-wasm-probe-' + Date.now();
	sendChatMessage({ body: probe, type: 'CHAT' }, accessToken);

	await sleep(3000);

	const messages = (
		await request
			.get('/api/admin/chat/messages')
			.auth(...ADMIN)
			.expect(200)
	).body;
	const replies = messages.filter(
		(m) =>
			m.user &&
			m.user.displayName === BOT_DISPLAY_NAME &&
			typeof m.body === 'string' &&
			m.body.includes(probe),
	);

	expect(replies.length).toBe(1);
});
