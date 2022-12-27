export function setup() {
	let windowErrorSpy;

	Cypress.on('window:before:load', (win) => {
		windowErrorSpy = cy.spy(win.console, 'error');
	});

	Cypress.on(
		'uncaught:exception',
		(err) => !err.message.includes('ResizeObserver loop limit exceeded')
	);

	describe('Listen for errors', () => {
		afterEach(() => {
			cy.wait(1000).then(() => {
				expect(windowErrorSpy).to.not.be.called;
			});
		});
	});
}
