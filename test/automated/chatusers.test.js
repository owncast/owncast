const { test } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const registerChat = require('./lib/chat').registerChat;
const sendChatMessage = require('./lib/chat').sendChatMessage;

const testVisibilityMessage = {
    body: "message " + Math.floor(Math.random() * 100),
    type: 'CHAT',
};

var userId
var accessToken
test('can register a user', async (done) => {
    const registration = await registerChat();
    userId = registration.id;
    accessToken = registration.accessToken;
    done();
});

test('can send a chat message', async (done) => {
    sendChatMessage(testVisibilityMessage, accessToken, done);
});

test('can disable a user', async (done) => {
    await request.post('/api/admin/chat/userenabled').send({ "userId": userId, "enabled": false })
    .auth('admin', 'abc123').expect(200);
    done();
});

test('verify user is disabled', async (done) => {
    const response = await request.get('/api/admin/chat/usersdisabled').auth('admin', 'abc123').expect(200);
    expect(response.body[response.body.length - 1].id).toBe(userId);
    done();
});

test('verify messages from user are hidden', async (done) => {
    const response = await request.get('/api/admin/chat/messages')
    .auth('admin', 'abc123')
    .expect(200);
    const message = response.body.filter(obj => {
        return obj.user.id === userId;
    });
    expect(message[0].hiddenAt).toBeTruthy();
    done();
});

test('can re-enable a user', async (done) => {
    await request.post('/api/admin/chat/userenabled').send({ "userId": userId, "enabled": true })
    .auth('admin', 'abc123').expect(200);
    done();
});

test('verify user is enabled', async (done) => {
    const response = await request.get('/api/admin/chat/usersdisabled').auth('admin', 'abc123').expect(200);
    expect(response.body.length).toBe(0);
    done();
});

test('verify messages from user are visible', async (done) => {
    const response = await request.get('/api/admin/chat/messages')
    .auth('admin', 'abc123')
    .expect(200);
    const message = response.body.filter(obj => {
        return obj.user.id === userId;
    });
    expect(message[0].hiddenAt).toBeUndefined();
    done();
});
