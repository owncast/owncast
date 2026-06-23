// OAuth viewer-auth gate (test fixture).
//
// A standard OAuth2 authorization-code client served under
// /plugins/oauth-gate-test/:
//
//   GET /          -> "sign in" screen linking to the provider's authorize URL
//   GET /callback  -> the provider redirects back with ?code; we exchange it for
//                     a token, fetch the userinfo, register the identity, and
//                     grant the gate session
//   GET /logout    -> clear the session
//
// Everything is hardcoded to the internal OAuth2 provider the test harness
// starts (fake-oauth-server.mjs, a real node-oidc-provider instance on :9876).
// The host owns the signed session cookie end to end via owncast.auth.
const { definePlugin, owncast } = require('@owncast/plugin-sdk');

const ISSUER = 'http://localhost:9876';
const AUTHORIZE_URL = ISSUER + '/authorize';
const TOKEN_URL = ISSUER + '/token';
const USERINFO_URL = ISSUER + '/userinfo';
const CLIENT_ID = 'owncast-oauth-test';
const CLIENT_SECRET = 'owncast-oauth-test-secret';
const REDIRECT_URI = 'http://localhost:8080/plugins/oauth-gate-test/callback';

function formEncode(obj) {
	return Object.keys(obj)
		.map(function (k) {
			return encodeURIComponent(k) + '=' + encodeURIComponent(obj[k]);
		})
		.join('&');
}

module.exports = definePlugin({
	onHttpRequest(req) {
		const query = req.query || {};

		// Login screen: send the viewer to the OAuth provider's authorize endpoint.
		if (req.method === 'GET' && req.path === '/') {
			const returnTo = query.return_to || '/';
			// A real plugin uses a random, per-request state; a fixed one keeps the
			// test deterministic (a single serial login flow).
			const state = 'oauth-test-state';
			owncast.kv.set('oauth_state:' + state, returnTo);
			const authorize =
				AUTHORIZE_URL +
				'?' +
				formEncode({
					response_type: 'code',
					client_id: CLIENT_ID,
					redirect_uri: REDIRECT_URI,
					scope: 'openid profile',
					state: state,
				});
			return {
				status: 200,
				headers: { 'content-type': 'text/html' },
				body:
					'<!doctype html><meta charset=utf-8><title>Sign in</title>' +
					'<h1>This stream is private</h1>' +
					'<a id="login" href="' +
					authorize +
					'">Sign in with OAuth</a>',
			};
		}

		// OAuth callback: exchange the code, look up the user, grant a session.
		if (req.method === 'GET' && req.path === '/callback') {
			const code = query.code || '';
			const state = query.state || '';

			// CSRF: the state must match one we issued.
			const returnTo = owncast.kv.get('oauth_state:' + state);
			if (!returnTo) return { status: 400, body: 'invalid or expired state' };
			owncast.kv.set('oauth_state:' + state, ''); // consume it

			// Exchange the authorization code for an access token (client_secret_post).
			const tokenResp = owncast.http.fetch(TOKEN_URL, {
				method: 'POST',
				headers: {
					'content-type': 'application/x-www-form-urlencoded',
					accept: 'application/json',
				},
				body: formEncode({
					grant_type: 'authorization_code',
					code: code,
					redirect_uri: REDIRECT_URI,
					client_id: CLIENT_ID,
					client_secret: CLIENT_SECRET,
				}),
			});
			const accessToken = JSON.parse(tokenResp.body || '{}').access_token || '';
			if (!accessToken) return { status: 502, body: 'token exchange failed' };

			// Look up the authenticated user behind the token.
			const userResp = owncast.http.fetch(USERINFO_URL, {
				headers: {
					authorization: 'Bearer ' + accessToken,
					accept: 'application/json',
				},
			});
			const profile = JSON.parse(userResp.body || '{}');
			const authId = 'oauth:' + profile.sub;
			const display = profile.name || profile.preferred_username || authId;

			// Identity, then session. register() finds-or-creates the Owncast user for
			// this OAuth identity (the host namespaces authId by our slug).
			const { userId } = owncast.users.register({
				authId: authId,
				displayName: display,
			});
			owncast.auth.grantSession({ userId: userId });

			return { status: 302, headers: { Location: returnTo || '/' } };
		}

		if (req.method === 'GET' && req.path === '/logout') {
			owncast.auth.endSession();
			return { status: 302, headers: { Location: '/' } };
		}

		return { status: 404, body: 'not found' };
	},
});
