var request = require('supertest');
request = request('http://127.0.0.1:8080');

const defaultAdminPassword = 'abc123';

async function getAdminResponse(endpoint, adminPassword = defaultAdminPassword) {
	const url = '/api/admin/' + endpoint;
	const res = request
		.get(url)
		.auth('admin', adminPassword)
		.expect(200);

	return res;
}

async function sendAdminRequest(
	endpoint,
	value,
	adminPassword = defaultAdminPassword
) {
	const url = '/api/admin/' + endpoint;
	const res = await request
		.post(url)
		.auth('admin', adminPassword)
		.send({ value: value })
		.expect(200);

	expect(res.body.success).toBe(true);
	return res;
}

async function sendAdminPayload(
	endpoint,
	payload,
	adminPassword = defaultAdminPassword
) {
	const url = '/api/admin/' + endpoint;
	const res = await request
		.post(url)
		.auth('admin', adminPassword)
		.send(payload)
		.expect(200);

	expect(res.body.success).not.toBe(false);

	return res;
}

async function failAdminRequest(
	endpoint,
	value,
	adminPassword = defaultAdminPassword,
	responseCode = 400
) {
	const url = '/api/admin/' + endpoint;
	const res = await request
		.post(url)
		.auth('admin', adminPassword)
		.send({ value: value })
		.expect(responseCode);

	expect(res.body.success).toBe(false);
	return res;
}

module.exports.getAdminResponse = getAdminResponse;
module.exports.sendAdminRequest = sendAdminRequest;
module.exports.sendAdminPayload = sendAdminPayload;
module.exports.failAdminRequest = failAdminRequest;

