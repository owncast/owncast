import { setup } from '../../support/setup.js';
setup();

describe(`Offline video embed`, () => {
	it('Can visit the page', () => {
		cy.visit('http://localhost:8080/embed/video');
	});

	// Offline banner
	it('Has correct offline banner values', () => {
		cy.contains('This stream is offline. Check back soon!').should(
			'be.visible'
		);
	});
});
