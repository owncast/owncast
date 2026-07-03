import { setup } from '../../support/setup.js';
import filterTests from '../../support/filterTests';
setup();

describe('Fediverse tests', () => {
	// Enable Fediverse features.
	it('Can visit the page', () => {
		cy.setConfig('serverurl', 'https://testing.biz');
		cy.setConfig('federation/enable', true);

		cy.visit('http://localhost:8080/');
	});

	// Offline banner
	it('Has correct offline banner values', () => {
		cy.contains(
			'This stream is offline. You can be notified the next time New Owncast Server goes live or follow streamer@testing.biz on the Fediverse.',
		).should('exist');
	});

	// Followers
	const followersCollection = '#followers-collection';
	it(`Has identifier: ${followersCollection}`, () => {
		cy.contains('Followers').click();
		cy.get(followersCollection).should('be.visible');
	});

	it(`Has identifier: #follow-button`, () => {
		cy.get('#follow-button').should('exist');
	});

	it('Can change to Followers tab', () => {
		cy.contains('Followers').click();
	});
});

// Desktop-only: at the mobile viewport the follow button fails Cypress's
// visibility check (covered by the mobile layout), so exercise the modal on
// desktop where the interaction is stable.
filterTests(['desktop'], () => {
	describe('Follow modal', () => {
		it('Can open the follow modal', () => {
			cy.get('#follow-button').click();
			cy.get('#follow-modal').should('be.visible');
			// The modal shows this server's fediverse account to follow.
			cy.get('#follow-modal').contains('streamer@testing.biz').should('exist');
		});

		it('Follow button is disabled until a valid account is entered', () => {
			cy.get('#follow-modal')
				.contains('button', 'Follow')
				.should('be.disabled');
			cy.get('#follow-modal input').type('not-a-valid-account');
			cy.get('#follow-modal')
				.contains('button', 'Follow')
				.should('be.disabled');
		});

		it('Shows an error when the remote account cannot be reached', () => {
			// .invalid is reserved (RFC 2606) and never resolves, so the server's
			// WebFinger lookup fails deterministically without external network.
			cy.get('#follow-modal input').clear();
			cy.get('#follow-modal input').type('someone@remote.invalid');
			cy.get('#follow-modal')
				.contains('button', 'Follow')
				.should('not.be.disabled');
			cy.get('#follow-modal').contains('button', 'Follow').click();
			cy.get('#follow-modal').contains('Follow Error').should('be.visible');
		});
	});
});
