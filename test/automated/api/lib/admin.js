var request = require('supertest');
request = request('http://127.0.0.1:8080');

const defaultAdminPassword = 'abc123';

async function getAdminConfig(adminPassword = defaultAdminPassword) {
	const res = request
		.get('/api/admin/serverconfig')
		.auth('admin', adminPassword)
		.expect(200);

	return res;
}

async function getAdminDisabledChatUsers(adminPassword = defaultAdminPassword) {
	const res = request
		.get('/api/admin/chat/users/disabled')
		.auth('admin', adminPassword)
		.expect(200);

	return res;
}

async function getAdminBlockedChatIPs(adminPassword = defaultAdminPassword) {
	const res = request
		.get('/api/admin/chat/users/ipbans')
		.auth('admin', adminPassword)
		.expect(200);

	return res;
}

async function getAdminStatus(adminPassword = defaultAdminPassword) {
	const res = request
		.get('/api/admin/status')
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

	expect(res.body.success).toBe(true);

	return res;
}

module.exports.getAdminConfig = getAdminConfig;
module.exports.getAdminStatus = getAdminStatus;
module.exports.getAdminDisabledChatUsers = getAdminDisabledChatUsers;
module.exports.getAdminBlockedChatIPs = getAdminBlockedChatIPs;
module.exports.sendAdminRequest = sendAdminRequest;
module.exports.sendAdminPayload = sendAdminPayload;
