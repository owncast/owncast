import { setup } from '../../support/setup.js';
import filterTests from '../../support/filterTests';

setup();

filterTests(['mobile'], () => {
	describe(`Live mobile tests`, () => {
		it('Can visit the page', () => {
			cy.visit('http://localhost:8080');
		});

		it('Mobile chat button should be visible', () => {
			cy.get('#mobile-chat-button').should('be.visible');
		});

		it('Click mobile chat button', () => {
			cy.get('#mobile-chat-button').click();
		});

		it('Mobile chat modal should be visible', () => {
			cy.get('.ant-modal').should('be.visible');
		});

		it('Chat container should be visible', () => {
			cy.get('#chat-container').should('be.visible');
		});

		it('Chat input should be visible', () => {
			cy.get('#chat-input').should('be.visible');
		});

		it('Chat user menu should be visible', () => {
			cy.get('#chat-modal-user-menu').should('be.visible');
		});

		it('Click on user menu', () => {
			cy.get('#chat-modal-user-menu').click();
		});

		it('Show change name modal', () => {
			cy.contains('Change name').click();
		});
	});
});
