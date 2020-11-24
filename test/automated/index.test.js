var request = require('supertest');
request = request('http://127.0.0.1:8080');

test('service is online', (done) => {
    request.get('/api/status').expect(200)
        .then((res) => {
            expect(res.body.online).toBe(true);
            done();
        });
});

test('frontend configuration is correct', (done) => {
    request.get('/api/config').expect(200)
        .then((res) => {
            expect(res.body.title).toBe('Owncast');
            expect(res.body.logo).toBe('/img/logo.svg');
            expect(res.body.socialHandles[0].platform).toBe('github');
            expect(res.body.socialHandles[0].url).toBe('http://github.com/owncast/owncast');
            done();
        });
});
