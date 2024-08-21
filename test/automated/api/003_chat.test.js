const { test } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const registerChat = require('./lib/chat').registerChat;
const sendChatMessage = require('./lib/chat').sendChatMessage;
const getAdminResponse = require('./lib/admin').getAdminResponse;
const randomNumber = require('./lib/rand').randomNumber;

var userDisplayName;
const message = randomNumber(100) + ' test 123';

const testMessage = {
	body: message,
	type: 'CHAT',
};

test('send a chat message', async () => {
	const registration = await registerChat();
	const accessToken = registration.accessToken;
	userDisplayName = registration.displayName;

	sendChatMessage(testMessage, accessToken);
});

test('fetch chat messages by admin', async () => {
	const res = await getAdminResponse('chat/messages');
	const expectedBody = `<p>` + testMessage.body + `</p>`;

	const message = res.body.filter((m) => m.body === expectedBody)[0];
	if (!message) {
		throw new Error('Message not found');
	}

	expect(message.body).toBe(expectedBody);
	expect(message.user.displayName).toBe(userDisplayName);
	expect(message.type).toBe(testMessage.type);
});

test('derive display name from user header', async () => {
	const res = await request
		.post('/api/chat/register')
		.set('X-Forwarded-User', 'test-user')
		.expect(200);

	expect(res.body.displayName).toBe('test-user');
});

test('overwrite user header derived display name with body', async () => {
	const res = await request
		.post('/api/chat/register')
		.send({ displayName: 'TestUserChat' })
		.expect(200);

	expect(res.body.displayName).toBe('TestUserChat');
});
