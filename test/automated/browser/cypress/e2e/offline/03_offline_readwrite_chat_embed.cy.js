import { setup } from '../../support/setup.js';
setup();

describe(`Offline readwrite chat embed`, () => {
	it('Can visit the page', () => {
		cy.visit('http://localhost:8080/embed/chat/readwrite');
	});
});
