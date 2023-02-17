/*
This test is to verify that the identifiers for specific components are
correctly set. This is to ensure that CSS customizations can be made to the web
UI using these specific IDs and/or class names.
These should be documented so people know how to customize their pages.
If you change one of these identifiers, you must update the documentation.
*/

import filterTests from '../../support/filterTests';

const identifiers = [
	'header', // The entire header component
	'footer', // The entire footer component
	'#global-header-text', // Just the text in the header
	'#offline-banner', // The entire offline banner component
	'#custom-page-content', // The entire custom page content component
	'#notify-button', // The notify button
	'#follow-button', // The follow button
];

filterTests(['desktop'], () => {
	describe(`Has correct identifiers for overrides`, () => {
		it('Can visit the page', () => {
			cy.visit('http://localhost:8080/');
		});

		// Loop over each identifier and verify it exists.
		identifiers.forEach((identifier) => {
			it(`Has identifier: ${identifier}`, () => {
				cy.get(identifier).should('exist');
			});
		});

		// Followers
		const followersCollection = '#followers-collection';
		it(`Has identifier: ${followersCollection}`, () => {
			cy.contains('Followers').click();
			cy.get(followersCollection).should('be.visible');
		});

		// Modal
		const modalContainer = '#modal-container';
		it(`Has identifier ${modalContainer}`, () => {
			cy.contains('Notify').click();
			cy.get(modalContainer, { timeout: 2000 }).should('be.visible');
		});
	});
});
