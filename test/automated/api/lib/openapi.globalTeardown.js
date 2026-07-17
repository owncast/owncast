// Writes results/openapi-coverage.json: which documented path+method pairs
// the suite exercised. Informational only; must never fail the run.
const fs = require('fs');
const path = require('path');
const openapi = require('./openapi');

module.exports = async () => {
	try {
		const routes = openapi.loadRoutes();
		const documented = routes.map((r) => `${r.method} ${r.path}`).sort();
		let lines = [];
		try {
			lines = fs
				.readFileSync(openapi.EXERCISED_FILE, 'utf8')
				.split('\n')
				.filter(Boolean);
		} catch (err) {
			// no requests recorded
		}
		const exercisedSet = new Set(lines);
		const exercised = documented.filter((d) => exercisedSet.has(d));
		const missed = documented.filter((d) => !exercisedSet.has(d));
		const out = {
			documentedCount: documented.length,
			exercisedCount: exercised.length,
			missedCount: missed.length,
			documented,
			exercised,
			missed,
		};
		fs.mkdirSync(openapi.RESULTS_DIR, { recursive: true });
		fs.writeFileSync(
			path.join(openapi.RESULTS_DIR, 'openapi-coverage.json'),
			JSON.stringify(out, null, 2),
		);
		console.log(
			`OpenAPI coverage: ${exercised.length}/${documented.length} documented endpoint+method pairs exercised (${missed.length} missed). See results/openapi-coverage.json`,
		);
	} catch (err) {
		console.warn(`OpenAPI coverage report skipped: ${err.message}`);
	}
};
