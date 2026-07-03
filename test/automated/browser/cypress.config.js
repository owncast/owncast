const fs = require('node:fs');
const path = require('node:path');
const { defineConfig } = require('cypress');
const remoteFediverse = require('./cypress/support/remote-fediverse-server');

// Append this spec's final failures to cypress/results/failures.json. run.sh
// clears the file at the start of a run and prints it at the end, so one run
// yields one machine-readable list of every failure across every group —
// spec, test, error text and screenshots — that an agent (or human) can act
// on without scrolling back through five groups of output.
function recordFailures(spec, results) {
	if (!results || !results.tests) return;
	const failures = results.tests
		.filter((test) => test.state === 'failed')
		.map((test) => ({
			spec: spec.relative,
			test: test.title.join(' > '),
			error: test.displayError,
			screenshots: (results.screenshots || []).map((s) => s.path),
		}));
	if (!failures.length) return;

	const out = path.join(__dirname, 'cypress', 'results', 'failures.json');
	fs.mkdirSync(path.dirname(out), { recursive: true });
	let existing = [];
	try {
		existing = JSON.parse(fs.readFileSync(out, 'utf8'));
	} catch {
		// first failing spec of the run
	}
	fs.writeFileSync(out, JSON.stringify(existing.concat(failures), null, 2));
}

module.exports = defineConfig({
	projectId: 'wwi3xe',
	e2e: {
		supportFile: 'cypress/support/e2e.js',
		// The specs in this suite visit the page once per describe block and
		// assert across multiple it() blocks, matching how a viewer actually
		// uses the app. Cypress 12+ would otherwise reset to about:blank
		// between tests.
		testIsolation: false,
		setupNodeEvents(on, config) {
			on('task', remoteFediverse.tasks);
			on('after:spec', recordFailures);
			return config;
		},
	},
	video: false,
	retries: 3,
});
