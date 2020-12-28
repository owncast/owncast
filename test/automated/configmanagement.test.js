var request = require('supertest');
request = request('http://127.0.0.1:8080');

const serverName = randomString();
const streamTitle = randomString();
const serverSummary = randomString();
const pageContent = randomString();
const logo = '/img/' + randomString();

test('set server name', async (done) => {
    const res = await sendConfigChangeRequest('name', serverName);
    done();
});

test('set stream title', async (done) => {
    const res = await sendConfigChangeRequest('streamtitle', streamTitle);
    done();
});

test('set server summary', async (done) => {
    const res = await sendConfigChangeRequest('serversummary', serverSummary);
    done();
});

test('set extra page content', async (done) => {
    const res = await sendConfigChangeRequest('pagecontent', pageContent);
    done();
});

test('set logo', async (done) => {
    const res = await sendConfigChangeRequest('logo', logo);
    done();
});


test('verify updated config values', async (done) => {
    const res = await request.get('/api/config');
    expect(res.body.name).toBe(serverName);
    expect(res.body.streamTitle).toBe(streamTitle);
    expect(res.body.summary).toBe(serverSummary);
    expect(res.body.extraPageContent).toBe(pageContent);
    expect(res.body.logo).toBe(logo);

    done();
});

async function sendConfigChangeRequest(endpoint, value) {
    const res = await request.post('/api/admin/config/' + endpoint)
        .auth('admin', 'abc123')
        .send({"value": value}).expect(200);

    expect(res.body.success).toBe(true);

    return res
}

function randomString(length = 20) {
    return Math.random().toString(16).substr(2, length);
}