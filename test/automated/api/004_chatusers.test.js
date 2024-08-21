const { test } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');
const WebSocket = require('ws');
const fs = require('fs');

const registerChat = require('./lib/chat').registerChat;
const sendChatMessage = require('./lib/chat').sendChatMessage;
const sendAdminRequest = require('./lib/admin').sendAdminRequest;
const sendAdminPayload = require('./lib/admin').sendAdminPayload;
const getAdminResponse = require('./lib/admin').getAdminResponse;
const randomNumber = require('./lib/rand').randomNumber;

const localIPAddressV4 = '127.0.0.1';
const localIPAddressV6 = '::1';

const testVisibilityMessage = {
	body: 'message ' + randomNumber(100),
	type: 'CHAT',
};

var userId;
var accessToken;
test('register a user', async () => {
	const registration = await registerChat();
	userId = registration.id;
	accessToken = registration.accessToken;
});

test('send a chat message', async () => {
	await sendChatMessage(testVisibilityMessage, accessToken);
});

test('set the user as moderator', async () => {
	const res = await sendAdminPayload('chat/users/setmoderator', {
		userId: userId,
		isModerator: true,
	});
});

test('verify user is a moderator', async () => {
	const response = await getAdminResponse('chat/users/moderators');
	const tokenCheck = response.body.filter((user) => user.id === userId);
	expect(tokenCheck).toHaveLength(1);
});

test('verify user list is populated', async () => {
	const ws = new WebSocket(
		`ws://localhost:8080/ws?accessToken=${accessToken}`,
		{
			origin: 'http://localhost:8080',
		}
	);

	ws.on('open', async function open() {
		const response = await getAdminResponse('chat/clients');

		expect(response.body.length).toBeGreaterThan(0);

		// Optionally, if GeoIP is configured, check the location property.
		if (fs.existsSync('../../../data/GeoLite2-City.mmdb')) {
			expect(response.body[0].geo.regionName).toBe('Localhost');
		}

		ws.close();
	});

	ws.on('error', function incoming(data) {
		console.error(data);
		ws.close();
	});

	ws.on('close', function incoming(data) {});
});

test('disable a user by admin', async () => {
	// To allow for visually being able to see the test hiding the
	// message add a short delay.
	await new Promise((r) => setTimeout(r, 1500));

	const ws = new WebSocket(
		`ws://localhost:8080/ws?accessToken=${accessToken}`,
		{
			origin: 'http://localhost:8080',
		}
	);

	const res = await sendAdminPayload('chat/users/setenabled', {
		userId: userId,
		enabled: false,
	});

	await new Promise((r) => setTimeout(r, 1500));
});

test('verify user is disabled', async () => {
	const response = await getAdminResponse('chat/users/disabled');
	const tokenCheck = response.body.filter((user) => user.id === userId);
	expect(tokenCheck).toHaveLength(1);
});

test('verify messages from user are hidden', async () => {
	const response = await getAdminResponse('chat/messages');
	const message = response.body.filter((obj) => {
		return obj.user.id === userId;
	});
	expect(message[0].user.disabledAt).toBeTruthy();
});

test('re-enable a user by admin', async () => {
	const res = await sendAdminPayload('chat/users/setenabled', {
		userId: userId,
		enabled: true,
	});
});

test('verify user is enabled', async () => {
	const response = await getAdminResponse('chat/users/disabled');
	const tokenCheck = response.body.filter((user) => user.id === userId);
	expect(tokenCheck).toHaveLength(0);
});

test('ban an ip address by admin', async () => {
	const resIPv4 = await sendAdminRequest(
		'chat/users/ipbans/create',
		localIPAddressV4
	);
	const resIPv6 = await sendAdminRequest(
		'chat/users/ipbans/create',
		localIPAddressV6
	);
});

test('verify IP address is blocked from the ban', async () => {
	const response = await getAdminResponse('chat/users/ipbans');

	expect(response.body).toHaveLength(2);
	expect(onlyLocalIPAddress(response.body)).toBe(true);
});

test('verify access is denied', async () => {
	await request.get(`/api/chat?accessToken=${accessToken}`).expect(401);
});

test('remove an ip address ban by admin', async () => {
	const resIPv4 = await sendAdminRequest(
		'chat/users/ipbans/remove',
		localIPAddressV4
	);
	const resIPv6 = await sendAdminRequest(
		'chat/users/ipbans/remove',
		localIPAddressV6
	);
});

test('verify IP address is no longer banned', async () => {
	const response = await getAdminResponse('chat/users/ipbans');

	expect(response.body).toHaveLength(0);
});

test('verify access is allowed after unban', async () => {
	await request.get(`/api/chat?accessToken=${accessToken}`).expect(200);
});

// This function expects the local address to be localIPAddressV4 & localIPAddressV6
function onlyLocalIPAddress(banInfo) {
	for (let i = 0; i < banInfo.length; i++) {
		if (
			banInfo[i].ipAddress != localIPAddressV4 &&
			banInfo[i].ipAddress != localIPAddressV6
		) {
			return false;
		}
	}
	return true;
}
