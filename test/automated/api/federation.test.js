var request = require('supertest');
const parseJson = require('parse-json');
const jsonfile = require('jsonfile');
const Ajv = require('ajv-draft-04');
const sendAdminRequest = require('./lib/admin').sendAdminRequest;

request = request('http://127.0.0.1:8080');

var ajv = new Ajv();
var nodeInfoSchema = jsonfile.readFileSync('schema/nodeinfo_2.0.json');

const serverURL = 'owncast.server.test'
const fediUsername = 'streamer'

test('disable federation', async (done) => {
	const res = await sendAdminRequest('config/federation/enable', false);
	done();
});

test('verify responses of /.well-known/webfinger when federation is disabled', async (done) => {
	const res = request.get('/.well-known/webfinger').expect(405);
	done();
});

test('verify responses of /.well-known/host-meta when federation is disabled', async (done) => {
	const res = request.get('/.well-known/host-meta').expect(405);
	done();
});

test('verify responses of /.well-known/nodeinfo when federation is disabled', async (done) => {
	const res = request.get('/.well-known/nodeinfo').expect(405);
	done();
});

test('verify responses of /.well-known/x-nodeinfo2 when federation is disabled', async (done) => {
	const res = request.get('/.well-known/x-nodeinfo2').expect(405);
	done();
});

test('verify responses of /nodeinfo/2.0 when federation is disabled', async (done) => {
	const res = request.get('/nodeinfo/2.0').expect(405);
	done();
});

test('verify responses of /api/v1/instance when federation is disabled', async (done) => {
	const res = request.get('/api/v1/instance').expect(405);
	done();
});

test('verify responses of /federation/user/ when federation is disabled', async (done) => {
	const res = request.get('/federation/user/').expect(405);
	done();
});

test('verify responses of /federation/ when federation is disabled', async (done) => {
	const res = request.get('/federation/').expect(405);
	done();
});

test('set required parameters and enable federation', async (done) => {
	const res1 = await sendAdminRequest(
		'config/serverurl',
		serverURL
	);
	const res2 = await sendAdminRequest(
		'config/federation/username',
		fediUsername
	);
	const res3 = await sendAdminRequest('config/federation/enable', true);
	done();
});

test('verify responses of /.well-known/webfinger when federation is enabled', async (done) => {
	const resNoResource = request.get('/.well-known/webfinger').expect(400);
	const resBadResource = request.get(
		'/.well-known/webfinger?resource=' + fediUsername + '@' + serverURL
	).expect(400);
	const resBadResource2 = request.get(
		'/.well-known/webfinger?resource=notacct:' + fediUsername + '@' + serverURL
	).expect(400);
	const resBadServer = request.get(
		'/.well-known/webfinger?resource=acct:' + fediUsername + '@not' + serverURL
	).expect(404);
	const resBadUser = request.get(
		'/.well-known/webfinger?resource=acct:not' + fediUsername + '@' + serverURL
	).expect(404);
	const resNoAccept = request.get(
		'/.well-known/webfinger?resource=acct:' + fediUsername + '@' + serverURL
	).expect(200)
		.expect('Content-Type', /json/)
		.then((res) => {
			expect(parseJson(res.body)).toBe(true);
			done();
		});
});

test('verify responses of /.well-known/host-meta when federation is enabled', async (done) => {
	const res = request.get('/.well-known/host-meta').expect(200);
	done();
});

test('verify responses of /.well-known/nodeinfo when federation is enabled', async (done) => {
	const res = request.get('/.well-known/nodeinfo').expect(200);
	done();
});

test('verify responses of /.well-known/x-nodeinfo2 when federation is enabled', async (done) => {
	const res = request.get('/.well-known/x-nodeinfo2').expect(200);
	done();
});

test('verify responses of /nodeinfo/2.0 when federation is enabled', async (done) => {
	const res = request
		.get('/nodeinfo/2.0')
		.expect(200)
		.then((res) => {
			expect(ajv.validate(nodeInfoSchema, res.body)).toBe(true);
			done();
		});
});

test('verify responses of /api/v1/instance when federation is enabled', async (done) => {
	const res = request.get('/api/v1/instance').expect(200);
	done();
});

test('verify responses of /federation/user/ when federation is enabled', async (done) => {
	const res = request.get('/federation/user/').expect(200);
	done();
});

test('verify responses of /federation/ when federation is enabled', async (done) => {
	const res = request.get('/federation/').expect(200);
	done();
});
