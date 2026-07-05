/*
Smoke coverage for the admin web interface. The admin is the most Ant
Design-heavy surface in the project, so these tests exist primarily to catch
UI framework regressions (e.g. antd upgrades) that break rendering or basic
navigation. Deeper federation protocol testing lives in
test/automated/activitypub/.
*/

import { setup } from '../../support/setup.js';
import filterTests from '../../support/filterTests';

setup();

const ADMIN_AUTH = { username: 'admin', password: 'abc123' };

filterTests(['desktop'], () => {
	describe('Admin layout and navigation', () => {
		it('Can visit the admin home page', () => {
			cy.setConfig('serverurl', 'https://testing.biz');
			cy.setConfig('federation/enable', true);
			cy.visit('http://localhost:8080/admin', { auth: ADMIN_AUTH });
			cy.get('#admin-page').should('exist');
		});

		it('Has the main navigation sections', () => {
			[
				'Home',
				'Viewers',
				'Users',
				'Chat',
				'Configuration',
				'Utilities',
				'Integrations',
				'Help',
			].forEach((item) => {
				cy.get('#admin-page .ant6-menu').contains(item).should('exist');
			});
		});

		it('Shows the Fediverse sections when federation is enabled', () => {
			cy.get('#admin-page .ant6-menu').contains('Followers').should('exist');
			cy.get('#admin-page .ant6-menu')
				.contains('Featured Streams')
				.should('exist');
		});

		it('Shows the stream status indicator', () => {
			cy.get('.online-status-indicator').should('contain', 'Offline');
		});
	});

	describe('General configuration page', () => {
		it('Can navigate to the General configuration page', () => {
			cy.get('#admin-page .ant6-menu').contains('Configuration').click();
			cy.get('#admin-page .ant6-menu').contains('General').click();
			cy.url().should('include', '/admin/config/general');
			cy.get('.config-public-details-page').should('exist');
		});

		it('Can switch between the configuration tabs', () => {
			cy.get('.config-public-details-page .ant6-tabs-tab')
				.contains('Appearance')
				.click();
			cy.get('.config-public-details-page .ant6-tabs-tab-active').should(
				'contain',
				'Appearance',
			);
		});
	});

	describe('Federation followers page', () => {
		it('Can view the followers page', () => {
			cy.visit('http://localhost:8080/admin/federation/followers', {
				auth: ADMIN_AUTH,
			});
			cy.get('#admin-page').should('exist');
			cy.get('.ant-table').should('exist');
		});
	});
});
