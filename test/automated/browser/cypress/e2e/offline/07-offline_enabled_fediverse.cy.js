import fetchData from '../../support/fetchData.js';
import { setup } from '../../support/setup.js';
setup();

describe('Fediverse tests', () => {
	// Enable Fediverse features.
	it('Can visit the page', () => {
		fetchData('http://localhost:8080/api/admin/config/serverurl', {
			method: 'POST',
			data: { value: 'https://testing.biz' },
		});
		fetchData('http://localhost:8080/api/admin/config/federation/enable', {
			method: 'POST',
			data: { value: true },
		});
		cy.wait(1500);

		cy.visit('http://localhost:8080/');
		// cy.reload(true, { timeout: 10000 });
	});

	// Offline banner
	it('Has correct offline banner values', () => {
		cy.contains(
			'This stream is offline. You can be notified the next time New Owncast Server goes live or follow streamer@testing.biz on the Fediverse.'
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
