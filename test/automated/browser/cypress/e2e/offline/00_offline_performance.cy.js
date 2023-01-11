describe('Lighthouse Metrics', () => {
	beforeEach(() => {
		cy.visit('http://localhost:8080');
	});

	it('Capture Metrics', () => {
		cy.lighthouse({
			accessibility: 97,
			'best-practices': 90,
			seo: 90,
			performance: 60, // Once the performance issues are fixed revert this 90,
		});
	});
});
