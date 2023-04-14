import { setup } from '../../support/setup.js';
setup();

describe(`Offline read-only chat embed`, () => {
	it('Can visit the page', () => {
		cy.visit('http://localhost:8080/embed/chat/readonly');
	});

	it('Chat should be visible', () => {
		cy.get('#chat-container').should('be.visible');
	});
});
