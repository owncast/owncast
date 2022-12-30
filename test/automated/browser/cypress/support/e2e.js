// ***********************************************************
// This example support/e2e.js is processed and
// loaded automatically before your test files.
//
// This is a great place to put global configuration and
// behavior that modifies Cypress.
//
// You can change the location of this file or turn off
// automatically serving support files with the
// 'supportFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/configuration
// ***********************************************************

// Import commands.js using ES2015 syntax:
import './commands';
import fetchData from './fetchData.js';

// Alternatively you can use CommonJS syntax:
// require('./commands')

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
