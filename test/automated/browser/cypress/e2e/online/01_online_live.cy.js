import { setup } from '../../support/setup.js';
import fetchData from '../../support/fetchData.js';
import filterTests from '../../support/filterTests';

setup();

describe(`Live tests`, () => {
	it('Can visit the page', () => {
		cy.visit('http://localhost:8080');
	});

	it('Should have a play button', () => {
		cy.get('.vjs-big-play-button').should('be.visible');
	});

	it('User menu should be visible', () => {
		cy.get('#user-menu').should('be.visible');
	});
});

filterTests(['desktop'], () => {
	describe(`Live desktop tests`, () => {
		it('Click on user menu', () => {
			cy.get('#user-menu').click();
		});
		it('Can toggle chat off', () => {
			cy.contains('Hide Chat').click();
		});

		it('Chat should not be visible', () => {
			cy.get('#chat-container').should('not.exist');
		});

		it('Click on user menu', () => {
			cy.get('#user-menu').click();
		});

		it('Can toggle chat on', () => {
			cy.contains('Show Chat').click();
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
			cy.wait(1500);
			// cy.contains('is now known as').should('be.visible');
		});

		it('Should change to custom websocket host', () => {
			fetchData('http://localhost:8080/api/admin/config/sockethostoverride', {
				method: 'POST',
				data: { value: 'ws://localhost:8080' },
			});
			cy.wait(1500);
		});

		it('Refresh page with new socket host', () => {
			cy.visit('http://localhost:8080');
		});
	});
});
