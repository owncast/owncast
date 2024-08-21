var request = require('supertest');
request = request('http://127.0.0.1:8080');

const getAdminResponse = require('./lib/admin').getAdminResponse;
const sendAdminPayload = require('./lib/admin').sendAdminPayload;

var accessToken = '';
var webhookID;

const webhook = 'https://super.duper.cool.thing.biz/owncast';
const events = ['CHAT'];

test('create webhook', async () => {
	const res = await sendAdminPayload('webhooks/create', {
		url: webhook,
		events: events,
	});

	expect(res.body.url).toBe(webhook);
	expect(res.body.timestamp).toBeTruthy();
	expect(res.body.events).toStrictEqual(events);
});

test('check webhooks', (done) => {
	getAdminResponse('webhooks').then((res) => {
		expect(res.body).toHaveLength(1);
		expect(res.body[0].url).toBe(webhook);
		expect(res.body[0].events).toStrictEqual(events);
		webhookID = res.body[0].id;
		done();
	});
});

test('delete webhook', async () => {
	const res = await sendAdminPayload('webhooks/delete', {
		id: webhookID,
	});
	expect(res.body.success).toBe(true);
});

test('check that webhook was deleted', (done) => {
	getAdminResponse('webhooks').then((res) => {
		expect(res.body).toHaveLength(0);
		done();
	});
});

test('create access token', async () => {
	const name = 'Automated integration test';
	const scopes = [
		'CAN_SEND_SYSTEM_MESSAGES',
		'CAN_SEND_MESSAGES',
		'HAS_ADMIN_ACCESS',
	];
	const res = await sendAdminPayload('accesstokens/create', {
		name: name,
		scopes: scopes,
	});

	expect(res.body.accessToken).toBeTruthy();
	expect(res.body.createdAt).toBeTruthy();
	expect(res.body.displayName).toBe(name);
	expect(res.body.scopes).toStrictEqual(scopes);
	accessToken = res.body.accessToken;
});

test('check access tokens', async () => {
	const res = await getAdminResponse('accesstokens');
	const tokenCheck = res.body.filter(
		(token) => token.accessToken === accessToken
	);
	expect(tokenCheck).toHaveLength(1);
});

test('send a system message using access token', async () => {
	const payload = {
		body: 'This is a test system message from the automated integration test',
	};
	await request
		.post('/api/integrations/chat/system')
		.set('Authorization', 'bearer ' + accessToken)
		.send(payload)
		.expect(200);
});

test('send an external integration message using access token', async () => {
	const payload = {
		body: 'This is a test external message from the automated integration test',
	};
	await request
		.post('/api/integrations/chat/send')
		.set('Authorization', 'Bearer ' + accessToken)
		.send(payload)
		.expect(200);
});

test('send an external integration action using access token', async () => {
	const payload = {
		body: 'This is a test external action from the automated integration test',
	};
	await request
		.post('/api/integrations/chat/action')
		.set('Authorization', 'Bearer ' + accessToken)
		.send(payload)
		.expect(200);
});

test('test fetch chat history using access token', async () => {
	await request
		.get('/api/integrations/chat')
		.set('Authorization', 'Bearer ' + accessToken)
		.expect(200);
});

test('test fetch chat history failure using invalid access token', async () => {
	await request
		.get('/api/integrations/chat')
		.set('Authorization', 'Bearer ' + 'invalidToken')
		.expect(401);
});

test('test fetch chat history OPTIONS request', async () => {
	const res = await request
		.options('/api/integrations/chat')
		.set('Authorization', 'Bearer ' + accessToken)
		.expect(204);
});

test('delete access token', async () => {
	const res = await sendAdminPayload('accesstokens/delete', {
		token: accessToken,
	});
	expect(res.body.success).toBe(true);
});

test('check token delete was successful', async () => {
	const res = await getAdminResponse('accesstokens');
	const tokenCheck = res.body.filter(
		(token) => token.accessToken === accessToken
	);
	expect(tokenCheck).toHaveLength(0);
});

test('send an external integration action using access token expecting failure', async () => {
	const payload = {
		body: 'This is a test external action from the automated integration test',
	};
	await request
		.post('/api/integrations/chat/action')
		.set('Authorization', 'Bearer ' + accessToken)
		.send(payload)
		.expect(401);
});
