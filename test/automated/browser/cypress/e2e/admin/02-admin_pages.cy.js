/*
Renders every admin page and fails on any console error (see support/setup.js:
the error text is included in the failure). The admin is the most Ant
Design-heavy surface in the project, so this spec is the early-warning system
for antd upgrades: a page that crashes, renders the Next.js client-exception
screen, or logs a React/antd error fails its own named test, telling you
exactly which page broke.

One it() per page keeps failures independent: a broken Video config page does
not hide a broken Webhooks page.
*/

import { setup } from '../../support/setup.js';
import filterTests from '../../support/filterTests';

setup();

const ADMIN_AUTH = { username: 'admin', password: 'abc123' };

// Every admin page reachable without extra route params.
// Deliberately excluded: /admin/upgrade (fetches releases from the GitHub
// API, an external dependency that would make this suite flaky offline).
const pages = [
	'/admin',
	'/admin/viewer-info',
	'/admin/users',
	'/admin/chat/messages',
	'/admin/chat/emojis',
	'/admin/access-tokens',
	'/admin/actions',
	'/admin/webhooks',
	'/admin/plugins',
	'/admin/config/general',
	'/admin/config/server',
	'/admin/config-video',
	'/admin/config-chat',
	'/admin/config-featured',
	'/admin/config-federation',
	'/admin/config-notify',
	'/admin/config-social-items',
	'/admin/hardware-info',
	'/admin/stream-health',
	'/admin/logs',
	'/admin/help',
	'/admin/federation/followers',
	'/admin/federation/actions',
];

filterTests(['desktop'], () => {
	describe('Every admin page renders', () => {
		pages.forEach((page) => {
			it(`Renders ${page}`, () => {
				cy.visit(`http://localhost:8080${page}`, { auth: ADMIN_AUTH });
				cy.get('#admin-page').should('exist');
				// The Next.js static export shows this text when a page throws
				// during client-side render.
				cy.contains('Application error').should('not.exist');
			});
		});
	});
});
