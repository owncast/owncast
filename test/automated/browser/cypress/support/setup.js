// Shared per-spec setup: fail any test whose page logged a console error.
//
// The captured error text is included in the assertion failure, so a failing
// run tells you exactly what the application logged, not just that
// "a spy was called".
export function setup() {
	let consoleErrors = [];

	Cypress.on('window:before:load', (win) => {
		const original = win.console.error;
		// eslint-disable-next-line no-param-reassign
		win.console.error = (...args) => {
			consoleErrors.push(
				args
					.map((a) => {
						if (a instanceof win.Error) return a.stack || a.message;
						if (typeof a === 'object') {
							try {
								return JSON.stringify(a);
							} catch {
								return String(a);
							}
						}
						return String(a);
					})
					.join(' '),
			);
			original.apply(win.console, args);
		};
	});

	Cypress.on(
		'uncaught:exception',
		(err) => !err.message.includes('ResizeObserver loop limit exceeded'),
	);

	describe('Listen for errors', () => {
		afterEach(() => {
			// Give any in-flight async work a beat to surface console errors
			// before asserting. Kept short: this runs after every single test.
			cy.wait(300).then(() => {
				const errors = consoleErrors;
				consoleErrors = [];
				expect(
					errors,
					`the page logged console.error:\n${errors.join('\n---\n')}`,
				).to.be.empty;
			});
		});
	});
}
