/*
This test is to verify that the identifiers for specific components are
correctly set. This is to ensure that CSS customizations can be made to the web
UI using these specific IDs and/or class names.
These should be documented so people know how to customize their pages.
If you change one of these identifiers, you must update the documentation.
*/

import filterTests from '../../support/filterTests';

const identifiers = [
	'#chat-container', // The entire chat container component
];

filterTests(['desktop'], () => {
	describe(`Has correct identifiers for overrides`, () => {
		it('Can visit the page', () => {
			cy.visit('http://localhost:8080/');
		});

		// Loop over each identifier and verify it exists.
		identifiers.forEach((identifier) => {
			it(`Has identifier: ${identifier}`, () => {
				cy.get(identifier).should('be.visible');
			});
		});
	});
});
