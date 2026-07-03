// Custom Cypress commands shared by all specs.

const ADMIN_USERNAME = 'admin';
const ADMIN_STREAMKEY = 'abc123';

// Set an Owncast admin config value and wait for the server to acknowledge
// it. Unlike a fire-and-forget fetch, this is properly chained into the
// Cypress command queue, so later visits observe the new config.
Cypress.Commands.add('setConfig', (path, value) => {
	cy.request({
		method: 'POST',
		url: `http://localhost:8080/api/admin/config/${path}`,
		auth: { username: ADMIN_USERNAME, password: ADMIN_STREAMKEY },
		body: { value },
	});
});

// Make an authenticated request to an Owncast admin API endpoint.
Cypress.Commands.add('adminRequest', (method, url, body) => {
	cy.request({
		method,
		url,
		auth: { username: ADMIN_USERNAME, password: ADMIN_STREAMKEY },
		body,
	});
});
