import { setup } from '../../support/setup.js';
import filterTests from '../../support/filterTests';

setup();

filterTests(['mobile'], () => {
	describe(`Live mobile tests`, () => {
		it('Can visit the page', () => {
			cy.visit('http://localhost:8080');
		});

		it('Mobile chat button should not be visible', () => {
			cy.get('#mobile-chat-button').should('not.exist');
		});
	});
});
