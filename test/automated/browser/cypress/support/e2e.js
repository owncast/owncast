import './commands';

before(() => {
	// Set server URL. Specs assume this instance is known as testing.biz.
	cy.setConfig('serverurl', 'https://testing.biz');
});
