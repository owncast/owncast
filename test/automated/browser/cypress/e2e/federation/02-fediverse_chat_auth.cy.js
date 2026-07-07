/*
Full end-to-end fediverse chat authentication (OTP over ActivityPub), using
the fake remote fediverse server in cypress/support/remote-fediverse-server.js.

The whole login loop runs with no mocks inside Owncast: the viewer asks to
link a fediverse account, Owncast performs a real WebFinger lookup against
the remote server, resolves the actor, and delivers a signed direct message
containing a one-time code to the remote inbox. The test reads the code out
of the harness inbox exactly like a human reading their DMs, submits it, and
verifies the chat user is visibly authenticated.

The flow is one test on purpose: a retry re-runs the whole journey with a
fresh remote account (Owncast will not re-send a DM for a still-pending
account, and codes are single-use).
*/

import { setup } from '../../support/setup.js';
import filterTests from '../../support/filterTests';

setup();

filterTests(['desktop'], () => {
	describe('Fediverse chat authentication', () => {
		it('Authenticates chat with a one-time code sent over ActivityPub', () => {
			cy.setConfig('federation/enable', true);

			cy.task('fediverse:createActor').then((actor) => {
				cy.visit('http://localhost:8080/');

				// Open the auth modal from the user menu and pick FediAuth. The
				// auth modal runs on Ant Design v6 (ant- class prefix).
				cy.get('#user-menu').click();
				cy.contains('Authenticate').click();
				cy.contains('.ant-tabs-tab', 'FediAuth').click();

				// Request a code for the remote account. Scope to the active tab
				// pane: antd keeps the inactive IndieAuth pane mounted, and it
				// contains its own search input and button.
				cy.get('.ant-modal .ant-tabs-content-active').within(() => {
					cy.get('input[placeholder="youraccount@yourserver.com"]').type(
						actor.account,
					);
					cy.get('.ant-input-search-btn').click();
				});

				// The code travels over real federation: WebFinger, actor
				// resolution, then a signed DM into the harness inbox.
				cy.task('fediverse:waitForOTPCode', { actorName: actor.name }).then(
					(code) => {
						cy.get('.ant-modal .ant-tabs-content-active').within(() => {
							cy.get('input[placeholder="123456"]').type(code);
							cy.contains('button', 'Verify Code').click();
						});
					},
				);

				// Successful verification reloads the page; the auth modal from
				// the old document is gone once the new one is active.
				cy.get('.ant-modal', { timeout: 15000 }).should('not.exist');
				cy.get('#user-menu').should('be.visible');

				// Visual proof: messages from this user now carry the
				// authenticated badge in chat.
				cy.get('#chat-input-content-editable').type(
					'hello from the fediverse{enter}',
				);
				cy.contains('.chat-message_user', 'hello from the fediverse')
					.find('.chat-user-badge')
					.should('have.attr', 'title', 'Authenticated');
			});
		});
	});
});
