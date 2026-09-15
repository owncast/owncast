import './commands';

let ranTest = false;

Cypress.on('test:after:run', () => {
	ranTest = true;
});

before(() => {
	// Set server URL. Specs assume this instance is known as testing.biz.
	cy.setConfig('serverurl', 'https://testing.biz');
});

after(() => {
	// The screenshot is taken from the final tested state, before Cypress
	// tears down the page. A frame chosen from the recording can be blank.
	if (Cypress.config('video') && ranTest) {
		cy.screenshot('preview');

		// Pad the recording so its in-flight tail is retained for playback.
		cy.wait(3000);
	}
});
