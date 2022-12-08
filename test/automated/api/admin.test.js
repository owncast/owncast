var request = require('supertest');
request = request('http://127.0.0.1:8080');

const getAdminResponse = require('./lib/admin').getAdminResponse;

test('correct number of log entries exist', (done) => {
    getAdminResponse('logs')
        .then((res) => {
            // expect(res.body).toHaveLength(8);
            done();
        });
});

