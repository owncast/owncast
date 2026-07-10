var request = require('supertest');
request = request('http://127.0.0.1:8080');

test('cmcd collector accepts a single event mode report', async () => {
	await request
		.post('/api/metrics/cmcd')
		.set('Content-Type', 'application/json')
		.send({
			v: 2,
			e: 'ps',
			ts: Date.now(),
			sta: 'p',
			ltc: 3200,
			mtp: 5100,
			sid: 'api-test-session',
			sn: 1,
		})
		.expect(200);
});

test('cmcd collector accepts a batched array of reports', async () => {
	await request
		.post('/api/metrics/cmcd')
		.set('Content-Type', 'application/json')
		.send([
			{ v: 2, e: 't', mtp: 5100, sid: 'api-test-session', sn: 2 },
			{ v: 2, e: 'ps', sta: 'w', bs: true, sid: 'api-test-session', sn: 3 },
		])
		.expect(200);
});

test('cmcd collector accepts a report as a query parameter', async () => {
	await request
		.get(
			'/api/metrics/cmcd?CMCD=' +
				encodeURIComponent('e=t,ltc=2900,mtp=4800,sid="api-test-session"'),
		)
		.expect(200);
});

test('cmcd collector answers cross-origin preflight', async () => {
	const res = await request
		.options('/api/metrics/cmcd')
		.set('Origin', 'https://example.com')
		.set('Access-Control-Request-Method', 'POST')
		.set('Access-Control-Request-Headers', 'content-type')
		.expect(204);
	expect(res.headers['access-control-allow-origin']).toBe('*');
	expect(res.headers['access-control-allow-methods']).toContain('POST');
	expect(res.headers['access-control-allow-headers'].toLowerCase()).toContain(
		'content-type',
	);
});

test('cmcd collector rejects malformed json', async () => {
	await request
		.post('/api/metrics/cmcd')
		.set('Content-Type', 'application/json')
		.send('{not json')
		.expect(400);
});

test('legacy playback metrics endpoint still accepts reports', async () => {
	await request
		.post('/api/metrics/playback')
		.set('Content-Type', 'application/json')
		.send({
			bandwidth: 1234,
			latency: 2.5,
			downloadDuration: 0.4,
			errors: 0,
			qualityVariantChanges: 0,
		})
		.expect(200);
});
