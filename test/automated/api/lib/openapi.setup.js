// Jest setupFilesAfterEnv hook: transparently validates every supertest
// response against openapi.yaml. No per-spec changes needed; failures are
// aggregated and thrown in afterAll so floating promises in old tests cannot
// turn into unhandled rejections.
const supertest = require('supertest');
const openapi = require('./openapi');

let routes = null;
const failures = [];
const seen = new Set();

beforeAll(() => {
	routes = openapi.loadRoutes();
});

const origEnd = supertest.Test.prototype.end;
supertest.Test.prototype.end = function end(fn) {
	return origEnd.call(this, (err, res) => {
		if (routes && res) {
			const found = [];
			try {
				openapi.checkResponse(routes, this.method, this.url, res, found);
			} catch (e) {
				found.push(
					`validator error for ${this.method} ${this.url}: ${e.message}`,
				);
			}
			for (const f of found) {
				if (!seen.has(f)) {
					seen.add(f);
					failures.push(f);
				}
			}
		}
		if (fn) fn(err, res);
	});
};

afterAll(() => {
	if (failures.length) {
		throw new Error(
			`OpenAPI response validation failed (${failures.length} problem(s)):\n  ` +
				failures.join('\n  '),
		);
	}
});
