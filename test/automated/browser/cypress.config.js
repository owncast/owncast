const { defineConfig } = require('cypress');

module.exports = defineConfig({
	projectId: 'wwi3xe',
	e2e: {
		setupNodeEvents(on, config) {
			// implement node event listeners here
		},
	},
});
