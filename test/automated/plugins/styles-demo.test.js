// End-to-end: styles-demo (built from the SDK) ships a CSS file in
// `manifest.styles`. The host reads the file's bytes from the plugin
// and exposes them in a dedicated `pluginStyles` field on /api/config
// (kept separate from the admin's own `customStyles` so the viewer can
// render plugin CSS as a baseline before the admin's styling wins on
// overlap). Verifies the full pipeline against a running Owncast
// instance.
const { test, expect, beforeAll } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const { enableOnly } = require('./lib/plugins');

const PLUGIN = 'styles-demo';

beforeAll(async () => {
	await enableOnly(PLUGIN);
});

test('/api/config pluginStyles contains the plugin CSS with a delimiter naming the source', async () => {
	const res = await request.get('/api/config').expect(200);
	const { pluginStyles } = res.body;
	expect(typeof pluginStyles).toBe('string');
	// The host prefixes each plugin contribution with a comment naming
	// the source plugin so devtools "view source" can attribute a rule
	// back to whichever plugin shipped it.
	expect(pluginStyles).toMatch(/\/\* plugin: styles-demo /);
	// Plugin's own CSS body should be inline.
	expect(pluginStyles).toMatch(/--plugin-styles-demo-loaded/);
	expect(pluginStyles).toMatch(/body::before/);
});

test('disabling the plugin drops its CSS from /api/config pluginStyles', async () => {
	// Switch to a different plugin so styles-demo unloads. Any of the
	// other examples works; pick one without `styles` so the assertion
	// is unambiguous.
	await enableOnly('echo-bot');
	const res = await request.get('/api/config').expect(200);
	const { pluginStyles } = res.body;
	expect(typeof pluginStyles).toBe('string');
	expect(pluginStyles).not.toMatch(/plugin: styles-demo/);
	expect(pluginStyles).not.toMatch(/--plugin-styles-demo-loaded/);
});
