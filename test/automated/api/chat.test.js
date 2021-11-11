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

  const expectedBody = `${testMessage.body}`;
  const message = res.body.filter(function (msg) {
    return msg.body === expectedBody;
  })[0];

  expect(message.body).toBe(expectedBody);
  expect(message.user.displayName).toBe(userDisplayName);
  expect(message.type).toBe(testMessage.type);

  done();
});
