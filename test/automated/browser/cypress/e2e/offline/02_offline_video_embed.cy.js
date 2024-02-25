import { setup } from '../../support/setup.js';
setup();

describe(`Offline video embed`, () => {
	it('Can visit the page', () => {
		cy.visit('http://localhost:8080/embed/video');
	});

	// Offline banner
	it('Has correct offline embed values', () => {
		cy.contains('This stream is not currently live.').should('be.visible');
	});
});
