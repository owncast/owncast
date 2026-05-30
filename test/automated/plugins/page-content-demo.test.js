// End-to-end: page-content-demo (built from the SDK) ships an HTML
// file in `manifest.extraPageContent`. The host reads the file's
// bytes from the plugin and prepends them to the admin's rendered
// extraPageContent on /api/config, so plugin HTML lands at the top
// of the viewer page's extra-content block. Verifies the full
// pipeline against a running Owncast instance.
const { test, expect, beforeAll } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const { enableOnly } = require('./lib/plugins');

const PLUGIN = 'page-content-demo';

beforeAll(async () => {
	await enableOnly(PLUGIN);
});

test('/api/config extraPageContent contains the plugin HTML with a delimiter naming the source', async () => {
	const res = await request.get('/api/config').expect(200);
	const { extraPageContent } = res.body;
	expect(typeof extraPageContent).toBe('string');
	// Each plugin contribution is wrapped with an HTML comment naming
	// the source so a reader can attribute the markup back.
	expect(extraPageContent).toMatch(/<!-- plugin: page-content-demo /);
	// Plugin's marker id should be inline.
	expect(extraPageContent).toMatch(/plugin-page-content-demo-banner/);
});

test('disabling the plugin drops its HTML from /api/config extraPageContent', async () => {
	// Switch to a different plugin so page-content-demo unloads. Pick
	// one without `extraPageContent` so the assertion is unambiguous.
	await enableOnly('echo-bot');
	const res = await request.get('/api/config').expect(200);
	const { extraPageContent } = res.body;
	expect(typeof extraPageContent).toBe('string');
	expect(extraPageContent).not.toMatch(/plugin: page-content-demo/);
	expect(extraPageContent).not.toMatch(/plugin-page-content-demo-banner/);
});
