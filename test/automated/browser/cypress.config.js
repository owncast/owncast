const { defineConfig } = require('cypress');

module.exports = defineConfig({
	projectId: 'wwi3xe',
	e2e: {
		supportFile: 'cypress/support/e2e.js',
		// The specs in this suite visit the page once per describe block and
		// assert across multiple it() blocks, matching how a viewer actually
		// uses the app. Cypress 12+ would otherwise reset to about:blank
		// between tests.
		testIsolation: false,
	},
	video: false,
	retries: 3,
});
