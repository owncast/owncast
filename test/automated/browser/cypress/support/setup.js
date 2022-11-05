export function setup() {
	Cypress.on(
		'uncaught:exception',
		(err) => !err.message.includes('ResizeObserver loop limit exceeded')
	);
}
