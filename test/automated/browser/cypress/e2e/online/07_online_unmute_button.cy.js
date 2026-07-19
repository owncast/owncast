import { setup } from '../../support/setup.js';
import filterTests from '../../support/filterTests';

setup();

// The big unmute overlay must only appear when the player itself muted
// playback (the muted autoplay fallback), never when the viewer muted on
// purpose. A deliberate mute, in-session or restored from a previous visit,
// keeps the player clean. Desktop-only: these tests click the control bar,
// which stays hidden on mobile viewports until tapped.
filterTests(['desktop'], () => {
	describe('Big unmute button', () => {
		// The autoplay test flips the instance config, so put it back for
		// whoever runs after us against the same datastore.
		after(() => {
			cy.setConfig('autoplay', 'off');
		});

		it('Does not appear when the viewer mutes manually', () => {
			cy.visit('http://localhost:8080', {
				onBeforeLoad(win) {
					win.localStorage.removeItem('owncast_volume');
				},
			});
			cy.get('.vjs-big-play-button').click();
			cy.get('video', { timeout: 15000 }).should(
				($v) => expect($v[0].paused).to.be.false,
			);
			cy.get('.vjs-mute-control').click();
			cy.get('video').should(($v) => expect($v[0].muted).to.be.true);
			cy.get('.vjs-big-unmute-button').should('not.be.visible');
		});

		it('Does not appear when a persisted mute is restored', () => {
			cy.visit('http://localhost:8080', {
				onBeforeLoad(win) {
					win.localStorage.setItem('owncast_volume', '0');
				},
			});
			cy.get('.vjs-big-play-button').click();
			cy.get('video', { timeout: 15000 }).should(
				($v) => expect($v[0].paused).to.be.false,
			);
			cy.get('video').should(($v) => expect($v[0].volume).to.equal(0));
			cy.get('.vjs-big-unmute-button').should('not.be.visible');
		});

		it('Appears when autoplay had to fall back to muted, and unmutes on click', () => {
			cy.setConfig('autoplay', 'always');
			cy.visit('http://localhost:8080', {
				onBeforeLoad(win) {
					win.localStorage.removeItem('owncast_volume');
					// Mimic a browser that blocks audible autoplay so video.js
					// takes its muted fallback path.
					const original = win.HTMLMediaElement.prototype.play;
					win.HTMLMediaElement.prototype.play = function play() {
						if (!this.muted) {
							return Promise.reject(
								new win.DOMException('autoplay blocked', 'NotAllowedError'),
							);
						}
						return original.call(this);
					};
				},
			});
			cy.get('video', { timeout: 15000 }).should(
				($v) => expect($v[0].muted && !$v[0].paused).to.be.true,
			);
			cy.get('.vjs-big-unmute-button').should('be.visible').click();
			cy.get('video').should(($v) => expect($v[0].muted).to.be.false);
			cy.get('.vjs-big-unmute-button').should('not.be.visible');
		});
	});
});
