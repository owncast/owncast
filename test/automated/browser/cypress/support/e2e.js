import './commands';

before(() => {
	// Set server URL. Specs assume this instance is known as testing.biz.
	cy.setConfig('serverurl', 'https://testing.biz');
});

after(() => {
	// When recording, pad the end of each spec: the CDP screencast delivers
	// frames with latency, and ending the spec immediately after its last
	// action drops the in-flight tail (recordings ended seconds before the
	// final tests' UI, e.g. the mobile name-change modal, ever appeared).
	if (Cypress.config('video')) {
		cy.wait(3000);
	}
});
