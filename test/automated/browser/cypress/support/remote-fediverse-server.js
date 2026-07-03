// A fake remote fediverse server for full end-to-end federation tests.
//
// Runs inside the Cypress Node process (wired up as cy.task()s in
// cypress.config.js) and speaks just enough real ActivityPub for Owncast to
// treat it as a genuine remote server:
//
//   - WebFinger  (GET /.well-known/webfinger)   account -> actor IRI
//   - Actor      (GET /actor/<name>)            Person doc with a public key
//   - Inbox      (POST /inbox/<name>)           captures deliveries (DMs,
//                                               Accepts) for assertions
//
// It can also SEND signed activities (e.g. Follow) to Owncast's inbox using
// draft-cavage HTTP signatures, the same scheme real fediverse servers use.
//
// Owncast must run with OWNCAST_ALLOW_INTERNAL_FEDERATION=true and
// OWNCAST_INSECURE_SKIP_VERIFY=true (set in run.sh) so its SSRF guard and TLS
// verification accept this loopback, self-signed server. The server listens
// on https://127.0.0.1:9443: Owncast requires https actor/key IRIs, and the
// web UI's account validator requires a dot in the hostname, which 127.0.0.1
// satisfies and "localhost" does not.

const crypto = require('node:crypto');
const https = require('node:https');
const { execFileSync } = require('node:child_process');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');

const HOST = '127.0.0.1';
const PORT = 9443;
const BASE_URL = `https://${HOST}:${PORT}`;

let server = null;

// name -> { name, account, iri, publicKeyPem, privateKey, inbox: [] }
const actors = new Map();

// The https server needs a TLS certificate; Node cannot mint one natively, so
// shell out to openssl (present on dev machines and CI runners).
function makeSelfSignedCert() {
	const dir = fs.mkdtempSync(path.join(os.tmpdir(), 'owncast-fedi-harness-'));
	const keyPath = path.join(dir, 'key.pem');
	const certPath = path.join(dir, 'cert.pem');
	execFileSync(
		'openssl',
		[
			'req',
			'-x509',
			'-newkey',
			'rsa:2048',
			'-keyout',
			keyPath,
			'-out',
			certPath,
			'-days',
			'2',
			'-nodes',
			'-subj',
			`/CN=${HOST}`,
		],
		{ stdio: 'ignore' },
	);
	return { key: fs.readFileSync(keyPath), cert: fs.readFileSync(certPath) };
}

function actorDocument(actor) {
	return {
		'@context': [
			'https://www.w3.org/ns/activitystreams',
			'https://w3id.org/security/v1',
		],
		type: 'Person',
		id: actor.iri,
		preferredUsername: actor.name,
		name: actor.name,
		inbox: `${BASE_URL}/inbox/${actor.name}`,
		publicKey: {
			id: `${actor.iri}#main-key`,
			owner: actor.iri,
			publicKeyPem: actor.publicKeyPem,
		},
	};
}

function handleRequest(req, res) {
	const url = new URL(req.url, BASE_URL);
	const respondJSON = (code, body) => {
		res.writeHead(code, { 'Content-Type': 'application/activity+json' });
		res.end(JSON.stringify(body));
	};

	if (url.pathname === '/.well-known/webfinger') {
		// resource=acct:name@127.0.0.1:9443
		const resource = url.searchParams.get('resource') || '';
		const name = resource.replace(/^acct:/, '').split('@')[0];
		const actor = actors.get(name);
		if (!actor) return respondJSON(404, { error: `no such actor ${name}` });
		return respondJSON(200, {
			subject: `acct:${actor.account}`,
			links: [
				{ rel: 'self', type: 'application/activity+json', href: actor.iri },
			],
		});
	}

	if (url.pathname.startsWith('/actor/')) {
		const actor = actors.get(url.pathname.split('/')[2]);
		if (!actor) return respondJSON(404, { error: 'no such actor' });
		return respondJSON(200, actorDocument(actor));
	}

	if (req.method === 'POST' && url.pathname.startsWith('/inbox/')) {
		const actor = actors.get(url.pathname.split('/')[2]);
		if (!actor) return respondJSON(404, { error: 'no such actor' });
		let body = '';
		req.on('data', (chunk) => {
			body += chunk;
		});
		req.on('end', () => {
			try {
				actor.inbox.push(JSON.parse(body));
			} catch {
				actor.inbox.push({ unparseable: body });
			}
			res.writeHead(202);
			res.end();
		});
		return undefined;
	}

	return respondJSON(404, { error: `unhandled ${req.method} ${url.pathname}` });
}

function ensureStarted() {
	if (server) return Promise.resolve(BASE_URL);
	return new Promise((resolve, reject) => {
		server = https.createServer(makeSelfSignedCert(), handleRequest);
		server.on('error', (err) => {
			server = null;
			reject(
				new Error(
					`fediverse harness failed to listen on ${BASE_URL}: ${err.message}`,
				),
			);
		});
		server.listen(PORT, HOST, () => resolve(BASE_URL));
	});
}

// Create a fresh actor (new name, new keypair, empty inbox). Flow specs
// create a new actor per test attempt so Cypress retries start clean:
// Owncast deduplicates repeat follows and OTP requests per account.
function createActor() {
	const name = `feditest${crypto.randomBytes(4).toString('hex')}`;
	const { publicKey, privateKey } = crypto.generateKeyPairSync('rsa', {
		modulusLength: 2048,
	});
	const actor = {
		name,
		account: `${name}@${HOST}:${PORT}`,
		iri: `${BASE_URL}/actor/${name}`,
		publicKeyPem: publicKey.export({ type: 'spki', format: 'pem' }),
		privateKey,
		inbox: [],
	};
	actors.set(name, actor);
	return { name: actor.name, account: actor.account, iri: actor.iri };
}

// Sign and deliver an activity to an Owncast inbox with a draft-cavage HTTP
// signature over (request-target), date and digest — what Owncast verifies.
async function deliverSignedActivity(actor, inboxUrl, activity) {
	const body = JSON.stringify(activity);
	const url = new URL(inboxUrl);
	const digest = `SHA-256=${crypto.createHash('sha256').update(body).digest('base64')}`;
	const date = new Date().toUTCString();
	const signingString = [
		`(request-target): post ${url.pathname}`,
		`date: ${date}`,
		`digest: ${digest}`,
	].join('\n');
	const signature = crypto
		.sign('sha256', Buffer.from(signingString), actor.privateKey)
		.toString('base64');

	let response;
	try {
		response = await fetch(inboxUrl, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/activity+json',
				Date: date,
				Digest: digest,
				Signature: `keyId="${actor.iri}#main-key",algorithm="rsa-sha256",headers="(request-target) date digest",signature="${signature}"`,
			},
			body,
		});
	} catch (err) {
		throw new Error(
			`Could not reach the Owncast inbox at ${inboxUrl}: ${err.cause?.message || err.message}. Is Owncast running?`,
		);
	}
	if (response.status >= 400) {
		throw new Error(
			`Owncast inbox POST ${inboxUrl} responded ${response.status}: ${await response.text()}`,
		);
	}
	return response.status;
}

// Wait for an activity matching `predicate` to arrive in the actor's inbox.
// On timeout, throw an error describing everything that DID arrive so a
// failing test explains itself.
async function waitForInboxActivity(actor, description, predicate, timeoutMs) {
	const deadline = Date.now() + timeoutMs;
	for (;;) {
		const match = actor.inbox.find(predicate);
		if (match) return match;
		if (Date.now() > deadline) {
			const seen =
				actor.inbox.map((a) => a.type || 'unknown').join(', ') || 'nothing';
			throw new Error(
				`Timed out after ${timeoutMs}ms waiting for ${description} in ${actor.account}'s inbox. ` +
					`Activities received instead: ${seen}. ` +
					`Check that Owncast is running with OWNCAST_ALLOW_INTERNAL_FEDERATION=true ` +
					`and OWNCAST_INSECURE_SKIP_VERIFY=true, and see the owncast log for delivery errors.`,
			);
		}
		// eslint-disable-next-line no-await-in-loop
		await new Promise((r) => {
			setTimeout(r, 250);
		});
	}
}

function getActorOrThrow(name) {
	const actor = actors.get(name);
	if (!actor) {
		throw new Error(
			`unknown harness actor "${name}" — call fediverse:createActor first`,
		);
	}
	return actor;
}

// cy.task() handlers. Tasks return JSON-serializable values; thrown errors
// fail the test with the error message.
const tasks = {
	// Start the harness (idempotent) and create a fresh remote actor.
	// Returns { name, account, iri }.
	async 'fediverse:createActor'() {
		await ensureStarted();
		return createActor();
	},

	// Send a signed Follow from the named actor to Owncast and wait for the
	// Accept to be federated back to the actor's inbox. Proves the entire
	// round trip: signature verification, actor resolution over webfinger-less
	// direct fetch, follower persistence, and outbound signed delivery.
	// Returns the Accept activity.
	async 'fediverse:sendFollow'({
		actorName,
		owncastInbox,
		owncastActorIri,
		timeoutMs = 15000,
	}) {
		const actor = getActorOrThrow(actorName);
		const follow = {
			'@context': 'https://www.w3.org/ns/activitystreams',
			id: `${BASE_URL}/activities/${crypto.randomUUID()}`,
			type: 'Follow',
			actor: actor.iri,
			object: owncastActorIri,
		};
		await deliverSignedActivity(actor, owncastInbox, follow);
		return waitForInboxActivity(
			actor,
			'an Accept of the Follow',
			(a) => a.type === 'Accept',
			timeoutMs,
		);
	},

	// Wait for a direct message (Create/Note) to arrive for the named actor
	// and extract the 6-digit one-time code from it. Returns the code string.
	async 'fediverse:waitForOTPCode'({ actorName, timeoutMs = 20000 }) {
		const actor = getActorOrThrow(actorName);
		const create = await waitForInboxActivity(
			actor,
			'the one-time-code direct message',
			(a) =>
				a.type === 'Create' && /\d{6}/.test(JSON.stringify(a.object || '')),
			timeoutMs,
		);
		const content =
			typeof create.object === 'object'
				? create.object.content
				: String(create.object);
		return content.match(/(\d{6})/)[1];
	},
};

module.exports = { tasks };
