// End-to-end: styles-demo (built from the SDK) ships a CSS file in
// `manifest.styles`. The host reads the file's bytes from the plugin
// and appends them to the admin's customStyles on /api/config, so
// the viewer renders one inline <style> block covering both
// sources. Verifies the full pipeline against a running Owncast
// instance.
const { test, expect, beforeAll } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const { enableOnly } = require('./lib/plugins');

const PLUGIN = 'styles-demo';

beforeAll(async () => {
	await enableOnly(PLUGIN);
});

test('/api/config customStyles contains the plugin CSS with a delimiter naming the source', async () => {
	const res = await request.get('/api/config').expect(200);
	const { customStyles } = res.body;
	expect(typeof customStyles).toBe('string');
	// The host prefixes each plugin contribution with a comment naming
	// the source plugin so devtools "view source" can attribute a rule
	// back to whichever plugin shipped it.
	expect(customStyles).toMatch(/\/\* plugin: styles-demo /);
	// Plugin's own CSS body should be inline.
	expect(customStyles).toMatch(/--plugin-styles-demo-loaded/);
	expect(customStyles).toMatch(/body::before/);
});

test('disabling the plugin drops its CSS from /api/config customStyles', async () => {
	// Switch to a different plugin so styles-demo unloads. Any of the
	// other examples works; pick one without `styles` so the assertion
	// is unambiguous.
	await enableOnly('echo-bot');
	const res = await request.get('/api/config').expect(200);
	const { customStyles } = res.body;
	expect(typeof customStyles).toBe('string');
	expect(customStyles).not.toMatch(/plugin: styles-demo/);
	expect(customStyles).not.toMatch(/--plugin-styles-demo-loaded/);
});
