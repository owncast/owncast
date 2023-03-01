import { setup } from '../../support/setup.js';
setup();

describe(`Offline readwrite chat embed`, () => {
	it('Can visit the page', () => {
		cy.visit('http://localhost:8080/embed/chat/readwrite');
	});

	it('Header should be visible', () => {
		cy.get('header').should('be.visible');
	});

	it('Chat should be visible', () => {
		cy.get('#chat-container').should('be.visible');
	});
});
