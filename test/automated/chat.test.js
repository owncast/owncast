const { test } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const registerChat = require('./lib/chat').registerChat;
const sendChatMessage = require('./lib/chat').sendChatMessage;

var userDisplayName;
const message = Math.floor(Math.random() * 100) + ' test 123';
// const messageRaw = message + ' *and some markdown too*';
// const messageMarkdown = '<p>' + message + ' <em>and some markdown too</em></p>'

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

test('can fetch chat messages', (done) => {
  request
    .get('/api/admin/chat/messages')
    .auth('admin', 'abc123')
    .expect(200)
    .then((res) => {
      const message = res.body.filter(function (msg) {
        return (msg.body = testMessage.body);
      })[0];

      expect(message.body).toBe(message.body);
      expect(message.user.displayName).toBe(userDisplayName);
      expect(message.type).toBe(testMessage.type);

      done();
    });
});
