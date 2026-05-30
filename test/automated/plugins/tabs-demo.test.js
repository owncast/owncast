// End-to-end: tabs-demo (built from the SDK) declares two viewer-page
// tabs in manifest.tabs. The host reads each tab's content file from
// the plugin's assets/ directory and inlines the HTML into the
// pluginTabs[] array on /api/config. Verifies the full pipeline
// against a running Owncast instance.
const { test, expect, beforeAll } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const { enableOnly } = require('./lib/plugins');

const PLUGIN = 'tabs-demo';

beforeAll(async () => {
	await enableOnly(PLUGIN);
});

test('/api/config pluginTabs advertises both tabs with their inlined HTML', async () => {
	const res = await request.get('/api/config').expect(200);
	const { pluginTabs } = res.body;
	expect(Array.isArray(pluginTabs)).toBe(true);
	const tabs = pluginTabs.filter((t) => t.slug === PLUGIN);
	expect(tabs.length).toBe(2);

	const music = tabs.find((t) => t.title === 'Music');
	expect(music).toBeDefined();
	expect(music.html).toMatch(/plugin-tabs-demo-music/);
	expect(music.html).toMatch(/What we're listening to/);

	const schedule = tabs.find((t) => t.title === 'Schedule');
	expect(schedule).toBeDefined();
	expect(schedule.html).toMatch(/plugin-tabs-demo-schedule/);
	expect(schedule.html).toMatch(/Stream schedule/);
});

test('disabling the plugin drops its tabs from /api/config', async () => {
	// Switch to a different plugin so tabs-demo unloads. Pick one
	// without manifest.tabs so the assertion is unambiguous.
	await enableOnly('echo-bot');
	const res = await request.get('/api/config').expect(200);
	const { pluginTabs } = res.body;
	expect(Array.isArray(pluginTabs)).toBe(true);
	expect(pluginTabs.find((t) => t.slug === PLUGIN)).toBeUndefined();
});
