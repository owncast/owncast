// Helpers for driving the plugin admin API during the integration tests.
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const ADMIN = ['admin', 'abc123'];

// Every example plugin this suite installs. The tests share one Owncast
// instance, so enableOnly() unloads the others to keep each test isolated.
const ALL = [
	'profanity-filter',
	'echo-bot',
	'overlay',
	'styles-demo',
	'scripts-demo',
	'viewer-gate',
	'page-content-demo',
	'tabs-demo',
];

const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

// enableOnly loads `name`, unloads every other example, and waits for the
// manager to instantiate the wasm runtime. The enable call returns 200 only
// when the plugin actually loaded, so this also asserts a successful load.
async function enableOnly(name) {
	for (const plugin of ALL) {
		if (plugin === name) {
			continue;
		}
		await request
			.post(`/api/admin/plugins/${plugin}/disable`)
			.auth(...ADMIN)
			.expect(200);
	}
	await request
		.post(`/api/admin/plugins/${name}/enable`)
		.auth(...ADMIN)
		.expect(200);
	await sleep(1500);
}

async function listPlugins() {
	const res = await request
		.get('/api/admin/plugins')
		.auth(...ADMIN)
		.expect(200);
	return res.body;
}

module.exports = { enableOnly, listPlugins, sleep, ADMIN };
