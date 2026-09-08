/*
This test is to verify that the identifiers for specific components are
correctly set. This is to ensure that CSS customizations can be made to the web
UI using these specific IDs and/or class names.
These should be documented so people know how to customize their pages.
If you change one of these identifiers, you must update the documentation.
*/

import filterTests from '../../support/filterTests';

const identifiers = [
	'header', // The entire header component
	'footer', // The entire footer component
	'#global-header-text', // Just the text in the header
	'#offline-banner', // The entire offline banner component
	'#custom-page-content', // The entire custom page content component
	'#notify-button', // The notify button
];

filterTests(['desktop'], () => {
	describe(`Has correct identifiers for overrides`, () => {
		it('Can visit the page', () => {
			cy.visit('http://localhost:8080/');
		});

		// Loop over each identifier and verify it exists.
		identifiers.forEach((identifier) => {
			it(`Has identifier: ${identifier}`, () => {
				cy.get(identifier).should('exist');
			});
		});

		// Modal
		const modalContainer = '#modal-container';
		it(`Has identifier ${modalContainer}`, () => {
			cy.contains('Notify').click();
			cy.get(modalContainer, { timeout: 2000 }).should('be.visible');
		});
	});

	describe('Applies documented CSS variable customizations', () => {
		before(() => {
			cy.setConfig('appearance', {
				'theme-color-action': 'blue',
			});
			cy.setConfig(
				'customstyles',
				`:root {
					--theme-color-action: red;
					--theme-color-components-primary-button-border: green;
					--theme-text-display-font-family: monospace;
					--theme-text-body-font-family: monospace;
				}

				#notify-button {
					border-width: 4px;
				}

				#notify-button:focus {
					background-color: rgb(0, 0, 255);
				}`,
			);
			cy.visit('http://localhost:8080/');
		});

		after(() => {
			cy.setConfig('customstyles', '');
			cy.setConfig('appearance', {});
		});

		it('applies font variables to native and Ant Design elements', () => {
			cy.get('body').should('have.css', 'font-family', 'monospace');
			cy.get('#global-header-text').should(
				'have.css',
				'font-family',
				'monospace',
			);
			cy.get('#notify-button').should('have.css', 'font-family', 'monospace');
		});

		it('lets custom variables override appearance variables', () => {
			cy.document().then((document) => {
				expect(
					getComputedStyle(document.documentElement)
						.getPropertyValue('--theme-color-action')
						.trim(),
				).to.equal('red');
			});
			cy.get('#notify-button').should(
				'have.css',
				'background-color',
				'rgb(255, 0, 0)',
			);
		});

		it('applies custom selectors with higher specificity than component styles', () => {
			cy.get('#notify-button').should('have.css', 'border-width', '4px');
		});

		it('applies custom interactive state selectors', () => {
			cy.get('#notify-button')
				.focus()
				.should('have.css', 'background-color', 'rgb(0, 0, 255)');
		});

		it('retains the cascade after a full page reload', () => {
			cy.reload();
			cy.get('#notify-button').should(
				'have.css',
				'background-color',
				'rgb(255, 0, 0)',
			);
		});
	});
});

filterTests(['mobile'], () => {
	describe('Applies responsive custom CSS', () => {
		before(() => {
			cy.setConfig(
				'customstyles',
				`@media (max-width: 600px) {
					#global-header-text {
						color: rgb(0, 128, 0);
					}
				}`,
			);
			cy.visit('http://localhost:8080/');
		});

		after(() => {
			cy.setConfig('customstyles', '');
		});

		it('applies media-query overrides at the mobile viewport', () => {
			cy.get('#global-header-text').should(
				'have.css',
				'color',
				'rgb(0, 128, 0)',
			);
		});
	});
});
