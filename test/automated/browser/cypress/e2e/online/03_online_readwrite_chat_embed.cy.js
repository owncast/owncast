import { setup } from '../../support/setup.js';
setup();

describe(`Online readwrite chat embed`, () => {
	it('Can visit the page', () => {
		cy.visit('http://localhost:8080/embed/chat/readwrite');
	});

	it('Header should be visible', () => {
		cy.get('header').should('be.visible');
	});

	it('User menu should be visible', () => {
		cy.get('#user-menu').should('be.visible');
	});

	it('Chat join message should exist', () => {
		cy.contains('joined the chat').should('be.visible');
	});

	it('Can send a chat message from the embed', () => {
		cy.get('#chat-input-content-editable').type('embed e2e message{enter}');
		cy.contains('.chat-message_user', 'embed e2e message').should('be.visible');
	});

	it('Click on user menu', () => {
		cy.get('#user-menu').click();
	});

	it('Show change name modal', () => {
		cy.contains('Change name').click();
	});

	it('Close name change modal', () => {
		cy.get('.ant6-modal-close-x').click();
	});
});
