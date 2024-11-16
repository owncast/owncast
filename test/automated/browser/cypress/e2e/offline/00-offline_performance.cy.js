// describe('Lighthouse Metrics', () => {
// 	beforeEach(() => {
// 		cy.visit('http://localhost:8080');
// 	});

// 	it('Capture Metrics', () => {
// 		cy.lighthouse({
// 			accessibility: 97,
// 			'best-practices': 90,
// 			seo: 90,
// 			performance: 0, // Once the performance issues are fixed revert this 90,
// 		});
// 	});
// });

import fetchData from '../../support/fetchData.js';

// Alternatively you can use CommonJS syntax:
// require('./commands')

console.log('------------- support/e2e.js');
// Put Owncast in a state where it's ready to be tested.

// Set server URL
fetchData('http://localhost:8080/api/admin/config/serverurl', {
	method: 'POST',
	data: { value: 'https://testing.biz' },
});

// Enable Fediverse features.
fetchData('http://localhost:8080/api/admin/config/federation/enable', {
	method: 'POST',
	data: { value: true },
});
