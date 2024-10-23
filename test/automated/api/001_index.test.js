var request = require('supertest');
request = request('http://127.0.0.1:8080');

test('service is online', () => {
	request
		.get('/api/status')
		.expect(200)
		.then((res) => {
			expect(res.body.online).toBe(true);
		});
});
