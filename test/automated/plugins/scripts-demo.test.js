// End-to-end: scripts-demo (built from the SDK) ships a JS file in
// `manifest.scripts`. The host reads the file's bytes from the plugin
// and appends them to the admin's customJavascript at the
// /customjavascript endpoint, so the viewer loads one <script> tag
// covering both sources. Verifies the full pipeline against a
// running Owncast instance.
const { test, expect, beforeAll } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const { enableOnly } = require('./lib/plugins');

const PLUGIN = 'scripts-demo';

beforeAll(async () => {
	await enableOnly(PLUGIN);
});

test('/customjavascript contains the plugin script with a delimiter naming the source', async () => {
	const res = await request.get('/customjavascript').expect(200);
	expect(res.headers['content-type']).toMatch(/javascript/);
	// The host prefixes each plugin contribution with a comment naming
	// the source so devtools can attribute a behavior back to whichever
	// plugin shipped it.
	expect(res.text).toMatch(/\/\/ plugin: scripts-demo /);
	// Plugin's own JS body should be inline.
	expect(res.text).toMatch(/__pluginScriptsDemoLoaded/);
	expect(res.text).toMatch(/\[scripts-demo\] plugin script loaded/);
});

test('disabling the plugin drops its JS from /customjavascript', async () => {
	// Switch to a different plugin so scripts-demo unloads. Pick one
	// without `scripts` so the assertion is unambiguous.
	await enableOnly('echo-bot');
	const res = await request.get('/customjavascript').expect(200);
	expect(res.text).not.toMatch(/plugin: scripts-demo/);
	expect(res.text).not.toMatch(/__pluginScriptsDemoLoaded/);
});
