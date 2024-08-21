var request = require('supertest');
const parseJson = require('parse-json');
const jsonfile = require('jsonfile');
const Ajv = require('ajv-draft-04');
const sendAdminRequest = require('./lib/admin').sendAdminRequest;

request = request('http://127.0.0.1:8080');

var ajv = new Ajv();
var nodeInfoSchema = jsonfile.readFileSync('schema/nodeinfo_2.0.json');

const serverName = 'owncast.server.test';
const serverURL = 'http://' + serverName;
const fediUsername = 'streamer';

test('disable federation', async () => {
	await sendAdminRequest('config/federation/enable', false);
});

test('verify responses of /.well-known/webfinger when federation is disabled', async () => {
	await request.get('/.well-known/webfinger').expect(405);
});

test('verify responses of /.well-known/host-meta when federation is disabled', async () => {
	await request.get('/.well-known/host-meta').expect(405);
});

test('verify responses of /.well-known/nodeinfo when federation is disabled', async () => {
	await request.get('/.well-known/nodeinfo').expect(405);
});

test('verify responses of /.well-known/x-nodeinfo2 when federation is disabled', async () => {
	await request.get('/.well-known/x-nodeinfo2').expect(405);
});

test('verify responses of /nodeinfo/2.0 when federation is disabled', async () => {
	await request.get('/nodeinfo/2.0').expect(405);
});

test('verify responses of /api/v1/instance when federation is disabled', async () => {
	await request.get('/api/v1/instance').expect(405);
});

test('verify responses of /federation/user/ when federation is disabled', async () => {
	await request.get('/federation/user/').expect(405);
});

test('verify responses of /federation/ when federation is disabled', async () => {
	await request.get('/federation/').expect(405);
});

test('set required parameters and enable federation', async () => {
	await sendAdminRequest('config/serverurl', serverURL);
	await sendAdminRequest('config/federation/username', fediUsername);
	await sendAdminRequest('config/federation/enable', true);
});

test('verify responses of /.well-known/webfinger when federation is enabled', async () => {
	const resNoResource = request.get('/.well-known/webfinger').expect(400);
	const resBadResource = request
		.get('/.well-known/webfinger?resource=' + fediUsername + '@' + serverName)
		.expect(400);
	const resBadResource2 = request
		.get(
			'/.well-known/webfinger?resource=notacct:' +
				fediUsername +
				'@' +
				serverName
		)
		.expect(400);
	const resBadServer = request
		.get(
			'/.well-known/webfinger?resource=acct:' +
				fediUsername +
				'@not' +
				serverName
		)
		.expect(404);
	const resBadUser = request
		.get(
			'/.well-known/webfinger?resource=acct:not' +
				fediUsername +
				'@' +
				serverName
		)
		.expect(404);
	const resNoAccept = request
		.get(
			'/.well-known/webfinger?resource=acct:' + fediUsername + '@' + serverName
		)
		.expect(200)
		.expect('Content-Type', /json/)
		.then((res) => {
			parseJson(res.text);
		});
	const resWithAccept = request
		.get(
			'/.well-known/webfinger?resource=acct:' + fediUsername + '@' + serverName
		)
		.expect(200)
		.set('Accept', 'application/json')
		.expect('Content-Type', /json/)
		.then((res) => {
			parseJson(res.text);
		});
});

test('verify responses of /.well-known/host-meta when federation is enabled', async () => {
	request
		.get('/.well-known/host-meta')
		.expect(200)
		.expect('Content-Type', /xml/);
});

test('verify responses of /.well-known/nodeinfo when federation is enabled', async () => {
	await request
		.get('/.well-known/nodeinfo')
		.expect(200)
		.expect('Content-Type', /json/)
		.then((res) => {
			parseJson(res.text);
		});
});

test('verify responses of /.well-known/x-nodeinfo2 when federation is enabled', async () => {
	await request
		.get('/.well-known/x-nodeinfo2')
		.expect(200)
		.expect('Content-Type', /json/)
		.then((res) => {
			parseJson(res.text);
		});
});

test('verify responses of /nodeinfo/2.0 when federation is enabled', async () => {
	await request
		.get('/nodeinfo/2.0')
		.expect(200)
		.expect('Content-Type', /json/)
		.then((res) => {
			parseJson(res.text);
			expect(ajv.validate(nodeInfoSchema, res.body)).toBe(true);
		});
});

test('verify responses of /api/v1/instance when federation is enabled', async () => {
	await request
		.get('/api/v1/instance')
		.expect(200)
		.expect('Content-Type', /json/)
		.then((res) => {
			parseJson(res.text);
		});
});

test('verify responses of /federation/user/ when federation is enabled', async () => {
	const resNoAccept = request.get('/federation/user/').expect(307);
	const resWithAccept = request
		.get('/federation/user/')
		.set('Accept', 'application/json')
		.expect(404);
	const resWithAcceptWrongUsername = request
		.get('/federation/user/not' + fediUsername)
		.set('Accept', 'application/json')
		.expect(404);
	const resWithAcceptUsername = request
		.get('/federation/user/' + fediUsername)
		.set('Accept', 'application/json')
		.expect(200)
		.expect('Content-Type', /json/)
		.then((res) => {
			parseJson(res.text);
		});
});

test('verify responses of /federation/ when federation is enabled', async () => {
	const resNoAccept = request.get('/federation/').expect(307);
	const resWithAccept = request
		.get('/federation/')
		.set('Accept', 'application/json')
		.expect(404);
});
