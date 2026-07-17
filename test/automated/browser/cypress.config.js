const fs = require('node:fs');
const path = require('node:path');
const { defineConfig } = require('cypress');
const remoteFediverse = require('./cypress/support/remote-fediverse-server');

// Append this spec's final failures to cypress/results/failures.json, and
// tests that only passed on a Cypress retry to cypress/results/flaky.json.
// run.sh clears both files at the start of a run and prints them at the end,
// so one run yields one machine-readable list of every failure and every
// swallowed flake across every group — spec, test, error text and
// screenshots — that an agent (or human) can act on without scrolling back
// through five groups of output.
function appendResults(file, records) {
	if (!records.length) return;
	const out = path.join(__dirname, 'cypress', 'results', file);
	fs.mkdirSync(path.dirname(out), { recursive: true });
	let existing = [];
	try {
		existing = JSON.parse(fs.readFileSync(out, 'utf8'));
	} catch {
		// first spec of the run to record
	}
	fs.writeFileSync(out, JSON.stringify(existing.concat(records), null, 2));
}

function recordResults(spec, results) {
	if (!results || !results.tests) return;
	appendResults(
		'failures.json',
		results.tests
			.filter((test) => test.state === 'failed')
			.map((test) => ({
				spec: spec.relative,
				test: test.title.join(' > '),
				error: test.displayError,
				screenshots: (results.screenshots || []).map((s) => s.path),
			})),
	);
	// A test that passed only on a Cypress retry (retries: 3 below) is a
	// swallowed flake: the suite exits 0 and nothing else surfaces it.
	appendResults(
		'flaky.json',
		results.tests
			.filter(
				(test) => test.state === 'passed' && (test.attempts || []).length > 1,
			)
			.map((test) => ({
				spec: spec.relative,
				test: test.title.join(' > '),
				attempts: test.attempts.length,
			})),
	);
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
			on('after:spec', recordResults);
			return config;
		},
	},
	// Off by default so local runs stay fast; CI enables recording through
	// cypress's native env override (CYPRESS_VIDEO=true in the workflow).
	// run.sh gives each group its own videosFolder so desktop and mobile
	// runs of the same specs don't overwrite each other.
	video: false,
	retries: 3,
});
