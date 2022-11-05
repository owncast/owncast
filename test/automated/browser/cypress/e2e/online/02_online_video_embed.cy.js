import { setup } from '../../support/setup.js';
setup();

describe(`Online video embed`, () => {
	it('Can visit the page', () => {
		cy.visit('http://localhost:8080');
	});

	it('Should have a play button', () => {
		cy.get('.vjs-big-play-button').should('be.visible');
	});
});
