export function setup() {
	let windowErrorSpy;

	Cypress.on('window:before:load', (win) => {
		windowErrorSpy = cy.spy(win.console, 'error');
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
				expect(windowErrorSpy).to.not.be.called;
			});
		});
	});
}
