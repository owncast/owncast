Cypress.on(
	'uncaught:exception',
	(err) => !err.message.includes('ResizeObserver loop limit exceeded')
);

describe('Live tests', () => {
	it('Can visit the page', () => {
		cy.visit('http://localhost:8080/');
	});

	it('Should have a play button', () => {
		cy.get('.vjs-big-play-button').should('be.visible');
	});

	it('Chat should be visible', () => {
		cy.get('#chat-container').should('be.visible');
	});

	it('User menu should be visible', () => {
		cy.get('#user-menu').should('be.visible');
	});

	it('Chat join message should exist', () => {
		cy.contains('joined the chat').should('be.visible');
	});

	it('User menu should be visible', () => {
		cy.get('#user-menu').should('be.visible');
	});

	it('Click on user menu', () => {
		cy.get('#user-menu').click();
	});

	it('Can toggle chat off', () => {
		cy.contains('Toggle chat').click();
	});

	it('Chat should not be visible', () => {
		cy.get('#chat-container').should('not.exist');
	});

	it('Click on user menu', () => {
		cy.get('#user-menu').click();
	});

	it('Can toggle chat on', () => {
		cy.contains('Toggle chat').click();
	});

	it('Chat should be visible', () => {
		cy.get('#chat-container').should('be.visible');
	});

	it('Click on user menu', () => {
		cy.get('#user-menu').click();
	});

	it('Show change name modal', () => {
		cy.contains('Change name').click();
	});

	it('Should change name', () => {
		cy.get('#name-change-field').focus();
		cy.get('#name-change-field').type('{ctrl+a}');
		cy.get('#name-change-field').type('my-new-name');
		cy.get('#name-change-submit').click();
		cy.get('.ant-modal-close-x').click();
		cy.contains('is now known as').should('be.visible');
	});
});
