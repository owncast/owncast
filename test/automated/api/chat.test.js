const { test } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const registerChat = require('./lib/chat').registerChat;
const sendChatMessage = require('./lib/chat').sendChatMessage;

var userDisplayName;
const message = Math.floor(Math.random() * 100) + ' test 123';

const testMessage = {
  body: message,
  type: 'CHAT',
};

test('can send a chat message', async (done) => {
  const registration = await registerChat();
  const accessToken = registration.accessToken;
  userDisplayName = registration.displayName;

  sendChatMessage(testMessage, accessToken, done);
});

test('can fetch chat messages', async (done) => {
  const res = await request
    .get('/api/admin/chat/messages')
    .auth('admin', 'abc123')
    .expect(200);

  const message = res.body.filter((m) => m.body === testMessage.body)[0];
  if (!message) {
    throw new Error('Message not found');
  }

  const expectedBody = testMessage.body;

  expect(message.body).toBe(expectedBody);
  expect(message.user.displayName).toBe(userDisplayName);
  expect(message.type).toBe(testMessage.type);

  done();
});

test('can derive display name from user header', async (done) => {
  const res = await request
    .post('/api/chat/register')
    .set('X-Forwarded-User', 'test-user')
    .expect(200);

  expect(res.body.displayName).toBe('test-user');
  done();
});

test('can overwrite user header derived display name with body', async (done) => {
  const res = await request
    .post('/api/chat/register')
    .send({ displayName: 'TestUserChat' })
    .set('X-Forwarded-User', 'test-user')
    .expect(200);

  expect(res.body.displayName).toBe('TestUserChat');
  done();
});
