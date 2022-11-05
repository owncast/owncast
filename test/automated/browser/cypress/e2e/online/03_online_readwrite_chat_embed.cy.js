import { setup } from '../../support/setup.js';
setup();

describe(`Online readwrite chat embed`, () => {
	it('Can visit the page', () => {
		cy.visit('http://localhost:8080/embed/chat/readwrite');
	});

	// it('Chat should be visible', () => {
	// 	cy.get('#chat-container').should('be.visible');
	// });

	// it('User menu should be visible', () => {
	// 	cy.get('#user-menu').should('be.visible');
	// });

	// it('Chat join message should exist', () => {
	// 	cy.contains('joined the chat').should('be.visible');
	// });

	// it('User menu should be visible', () => {
	// 	cy.get('#user-menu').should('be.visible');
	// });

	// it('Click on user menu', () => {
	// 	cy.get('#user-menu').click();
	// });

	// it('Show change name modal', () => {
	// 	cy.contains('Change name').click();
	// });
});
