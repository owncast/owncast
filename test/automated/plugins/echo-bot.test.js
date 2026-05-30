// End-to-end: echo-bot (built from the SDK) receives a chat message via its
// onChatMessage handler and replies with owncast.chat.send, posting under its
// own bot identity. Also verifies the host doesn't loop on the bot's reply.
const { test, expect, beforeAll } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const { registerChat, sendChatMessage } = require('./lib/chat');
const { enableOnly, sleep, ADMIN } = require('./lib/plugins');

const PLUGIN = 'echo-bot';
// Bot display name as declared in the example manifest's `bot.displayName`.
// The chat-bot user posts under this string, independent of the plugin's
// slug or display name.
const BOT_DISPLAY_NAME = 'Example Echo';

beforeAll(async () => {
	await enableOnly(PLUGIN);
});

test('echoes a chat message once, under its own bot identity (no loop)', async () => {
	const { accessToken } = await registerChat();
	const probe = 'echo-probe-' + Date.now();
	sendChatMessage({ body: probe, type: 'CHAT' }, accessToken);

	// user message -> onChatMessage -> chat.send (as the bot) -> persisted.
	await sleep(3000);

	const messages = (
		await request
			.get('/api/admin/chat/messages')
			.auth(...ADMIN)
			.expect(200)
	).body;
	const echoes = messages.filter(
		(m) =>
			m.user &&
			m.user.displayName === BOT_DISPLAY_NAME &&
			typeof m.body === 'string' &&
			m.body.includes(probe),
	);

	// Exactly one: the bot reacted to the user's message, but its own reply
	// (authored by a bot) was not delivered back to it, so it didn't loop.
	expect(echoes.length).toBe(1);
});
