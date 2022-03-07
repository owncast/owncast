const { test } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');
const WebSocket = require('ws');

const registerChat = require('./lib/chat').registerChat;
const sendChatMessage = require('./lib/chat').sendChatMessage;
const listenForEvent = require('./lib/chat').listenForEvent;

const testVisibilityMessage = {
  body: 'message ' + Math.floor(Math.random() * 100),
  type: 'CHAT',
};

test('can send a chat message', async (done) => {
  const registration = await registerChat();
  const accessToken = registration.accessToken;

  sendChatMessage(testVisibilityMessage, accessToken, done);
});

test('verify we can make API call to mark message as hidden', async (done) => {
  const registration = await registerChat();
  const accessToken = registration.accessToken;
  const ws = new WebSocket(
    `ws://localhost:8080/ws?accessToken=${accessToken}`,
    {
      origin: 'http://localhost:8080',
    }
  );

  // Verify the visibility change comes through the websocket
  ws.on('message', function incoming(message) {
    const messages = message.split('\n');
    messages.forEach(function (message) {
      const event = JSON.parse(message);

      if (event.type === 'VISIBILITY-UPDATE') {
        done();
        ws.close();
      }
    });
  });

  const res = await request
    .get('/api/admin/chat/messages')
    .auth('admin', 'abc123')
    .expect(200);

  const message = res.body[0];
  const messageId = message.id;
  await request
    .post('/api/admin/chat/updatemessagevisibility')
    .auth('admin', 'abc123')
    .send({ idArray: [messageId], visible: false })
    .expect(200);
});

test('verify message has become hidden', async (done) => {
  const res = await request
    .get('/api/admin/chat/messages')
    .expect(200)
    .auth('admin', 'abc123');

  const message = res.body.filter((obj) => {
    return obj.body === `${testVisibilityMessage.body}`;
  });
  expect(message.length).toBe(1);
  // expect(message[0].hiddenAt).toBeTruthy();
  done();
});
