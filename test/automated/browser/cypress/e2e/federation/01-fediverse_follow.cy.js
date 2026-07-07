/*
Full end-to-end inbound follow flow, using the fake remote fediverse server
in cypress/support/remote-fediverse-server.js. Nothing is mocked inside
Owncast: the harness delivers a real signed Follow activity to the inbox,
Owncast verifies the HTTP signature, resolves the remote actor over HTTPS,
persists the follower, federates a signed Accept back to the harness, and
broadcasts an engagement event to chat.

Assertions cover both layers:
	- protocol: the Accept arrives back at the remote server's inbox
	- UI: the follow appears live in chat, in the Followers tab, and in the
		admin followers table

Each test attempt creates a fresh remote actor because Owncast deduplicates
repeat follows per actor (a re-follow does not re-fire the chat event).
*/

import { setup } from '../../support/setup.js';
import filterTests from '../../support/filterTests';

setup();

const OWNCAST_INBOX = 'http://localhost:8080/federation/user/streamer/inbox';
// Actor IRIs are derived from the configured server URL (set to
// https://testing.biz for the whole suite in support/e2e.js).
const OWNCAST_ACTOR_IRI = 'https://testing.biz/federation/user/streamer';

filterTests(['desktop'], () => {
	describe('Inbound fediverse follow', () => {
		it('A signed Follow is Accepted and shows live in chat', () => {
			cy.setConfig('federation/enable', true);
			cy.setConfig('federation/private', false);
			cy.setConfig('federation/showengagement', true);

			// Chat must be connected before the follow arrives: engagement
			// events are broadcast live over the websocket and are not part of
			// the chat history backlog.
			cy.visit('http://localhost:8080/');
			cy.get('#chat-container').should('be.visible');

			cy.task('fediverse:createActor').then((actor) => {
				cy.task('fediverse:sendFollow', {
					actorName: actor.name,
					owncastInbox: OWNCAST_INBOX,
					owncastActorIri: OWNCAST_ACTOR_IRI,
				}).then((accept) => {
					// Protocol-level proof: Owncast federated a signed Accept
					// back to the remote server.
					expect(accept.type).to.equal('Accept');
					expect(accept.actor).to.equal(OWNCAST_ACTOR_IRI);
				});

				// Visual proof: the follow shows up in the chat feed as a
				// social event card.
				cy.get('.chat-message_social', { timeout: 15000 })
					.should('contain', actor.name)
					.and('contain', 'followed this stream');

				// The new follower is listed in the Followers tab.
				cy.contains('Followers').click();
				cy.get('#followers-collection').should('contain', actor.name);
			});
		});

		it('The follower appears in the admin followers table', () => {
			cy.visit('http://localhost:8080/admin/federation/followers', {
				auth: { username: 'admin', password: 'abc123' },
			});
			// The table lists the follower's name and actor URL. Any follower
			// from the harness domain proves persistence without coupling this
			// test to the actor created by the previous one.
			cy.get('.ant6-table').should('contain', '127.0.0.1:9443/actor/');
		});
	});
});
