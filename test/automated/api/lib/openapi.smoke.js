// Server-free self-check for the OpenAPI response validator.
// Run: node lib/openapi.smoke.js
const assert = require('assert');
const supertest = require('supertest');
const openapi = require('./openapi');

function fakeRes(status, body) {
	return { status, type: 'application/json', body };
}

(async () => {
	assert(
		typeof supertest.Test.prototype.end === 'function',
		'supertest.Test missing',
	);

	const routes = openapi.loadRoutes();
	assert(routes.length > 100, `expected >100 routes, got ${routes.length}`);

	// exact path match
	const status = openapi.matchRoute(routes, 'GET', '/api/status');
	assert(status && status.path === '/status', 'GET /api/status did not match');

	// templated path match
	const tmpl = openapi.matchRoute(
		routes,
		'DELETE',
		'/api/admin/federation/servers/123',
	);
	assert(
		tmpl && tmpl.path === '/admin/federation/servers/{id}',
		'templated path did not match',
	);

	// undocumented path ignored
	assert.strictEqual(openapi.matchRoute(routes, 'GET', '/robots.txt'), null);

	// valid canned /status body passes
	let failures = [];
	openapi.checkResponse(
		routes,
		'GET',
		'http://127.0.0.1:8080/api/status',
		fakeRes(200, {
			online: true,
			viewerCount: 3,
			serverTime: '2026-07-16T00:00:00Z',
			versionNumber: '0.2.4',
			streamTitle: '',
			lastConnectTime: '2026-07-16T00:00:00Z',
			lastDisconnectTime: '2026-07-16T00:00:00Z',
		}),
		failures,
	);
	assert.deepStrictEqual(failures, [], `valid body flagged: ${failures}`);

	// invalid body (wrong type) fails
	failures = [];
	openapi.checkResponse(
		routes,
		'GET',
		'http://127.0.0.1:8080/api/status',
		fakeRes(200, { online: 'yes', viewerCount: 'many' }),
		failures,
	);
	assert(
		failures.length === 1 && /does not match schema/.test(failures[0]),
		`expected schema failure, got: ${failures}`,
	);

	// undocumented status flagged
	failures = [];
	openapi.checkResponse(
		routes,
		'GET',
		'http://127.0.0.1:8080/api/status',
		fakeRes(418, {}),
		failures,
	);
	assert(
		failures.length === 1 && /status not documented/.test(failures[0]),
		`expected status failure, got: ${failures}`,
	);

	console.log(
		`openapi.smoke: OK (${routes.length} documented method+path pairs)`,
	);
})().catch((err) => {
	console.error('openapi.smoke: FAIL', err);
	process.exit(1);
});
