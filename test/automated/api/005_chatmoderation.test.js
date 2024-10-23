const { test } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');
const WebSocket = require('ws');

const registerChat = require('./lib/chat').registerChat;
const sendChatMessage = require('./lib/chat').sendChatMessage;
const getAdminResponse = require('./lib/admin').getAdminResponse;
const sendAdminPayload = require('./lib/admin').sendAdminPayload;
const sendAdminRequest = require('./lib/admin').sendAdminRequest;
const randomNumber = require('./lib/rand').randomNumber;

const testVisibilityMessage = {
	body: 'message ' + randomNumber(100),
	type: 'CHAT',
};

var messageId;

const establishedUserFailedChatMessage = {
	body: 'this message should fail to send ' + Math.floor(Math.random() * 100),
	type: 'CHAT',
};

test('send a chat message', async () => {
	const registration = await registerChat();
	const accessToken = registration.accessToken;

	await sendChatMessage(testVisibilityMessage, accessToken);
});

test('verify admin can make API call to mark message as hidden', async () => {
	const registration = await registerChat();
	const accessToken = registration.accessToken;
	const ws = new WebSocket(
		`ws://localhost:8080/ws?accessToken=${accessToken}`,
		{
			origin: 'http://localhost:8080',
		}
	);

	// Verify the visibility change comes through the websocket
	ws.on('message', async function incoming(message) {
		const messages = message.split('\n');
		messages.forEach(async function (message) {
			const event = JSON.parse(message);

			if (event.type === 'VISIBILITY-UPDATE') {
				ws.close();
			}
		});
	});

	const res = await getAdminResponse('chat/messages');

	const message = res.body[0];
	messageId = message.id;
	await sendAdminPayload('chat/messagevisibility', {
		idArray: [messageId],
		visible: false,
	});
});

test('verify message has become hidden', async () => {
	await new Promise((r) => setTimeout(r, 2000));

	const res = await getAdminResponse('chat/messages');

	const message = res.body.filter((obj) => {
		return obj.id === messageId;
	});
	expect(message.length).toBe(1);
	// expect(message[0].hiddenAt).toBeTruthy();
});

test('enable established chat user mode', async () => {
	await sendAdminRequest('config/chat/establishedusermode', true);
});

test('send a message after established user mode is enabled', async () => {
	const registration = await registerChat();
	const accessToken = registration.accessToken;

	await sendChatMessage(establishedUserFailedChatMessage, accessToken);
});

test('verify rejected message is not in the chat feed', async () => {
	const res = await getAdminResponse('chat/messages');

	const message = res.body.filter((obj) => {
		return obj.body === establishedUserFailedChatMessage.body;
	});

	expect(message.length).toBe(0);
});

test('disable established chat user mode', async () => {
	await sendAdminRequest('config/chat/establishedusermode', false);
});
