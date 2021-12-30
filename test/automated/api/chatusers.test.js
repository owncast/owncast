const { test } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');
const WebSocket = require('ws');
const fs = require('fs');

const registerChat = require('./lib/chat').registerChat;
const sendChatMessage = require('./lib/chat').sendChatMessage;

const testVisibilityMessage = {
  body: 'message ' + Math.floor(Math.random() * 100),
  type: 'CHAT',
};

var userId;
var accessToken;
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
  // To allow for visually being able to see the test hiding the
  // message add a short delay.
  await new Promise((r) => setTimeout(r, 1500));

  await request
    .post('/api/admin/chat/users/setenabled')
    .send({ userId: userId, enabled: false })
    .auth('admin', 'abc123')
    .expect(200);
  done();
});

test('verify user is disabled', async (done) => {
  const response = await request
    .get('/api/admin/chat/users/disabled')
    .auth('admin', 'abc123')
    .expect(200);
  const tokenCheck = response.body.filter((user) => user.id === userId);
  expect(tokenCheck).toHaveLength(1);
  done();
});

test('verify messages from user are hidden', async (done) => {
  const response = await request
    .get('/api/admin/chat/messages')
    .auth('admin', 'abc123')
    .expect(200);
  const message = response.body.filter((obj) => {
    return obj.user.id === userId;
  });
  expect(message[0].user.disabledAt).toBeTruthy();
  done();
});

test('can re-enable a user', async (done) => {
  await request
    .post('/api/admin/chat/users/setenabled')
    .send({ userId: userId, enabled: true })
    .auth('admin', 'abc123')
    .expect(200);
  done();
});

test('verify user is enabled', async (done) => {
  const response = await request
    .get('/api/admin/chat/users/disabled')
    .auth('admin', 'abc123')
    .expect(200);
  const tokenCheck = response.body.filter((user) => user.id === userId);
  expect(tokenCheck).toHaveLength(0);

  done();
});

test('can set the user as moderator', async (done) => {
  await request
    .post('/api/admin/chat/users/setmoderator')
    .send({ userId: userId, isModerator: true })
    .auth('admin', 'abc123')
    .expect(200);
  done();
});

test('verify user is a moderator', async (done) => {
  const response = await request
    .get('/api/admin/chat/users/moderators')
    .auth('admin', 'abc123')
    .expect(200);
  const tokenCheck = response.body.filter((user) => user.id === userId);
  expect(tokenCheck).toHaveLength(1);

  done();
});

test('verify user list is populated', async (done) => {
  const ws = new WebSocket(
    `ws://localhost:8080/ws?accessToken=${accessToken}`,
    {
      origin: 'http://localhost:8080',
    }
  );

  ws.on('open', async function open() {
    const response = await request
      .get('/api/admin/chat/clients')
      .auth('admin', 'abc123')
      .expect(200);

    expect(response.body.length).toBeGreaterThan(0);

    // Optionally, if GeoIP is configured, check the location property.
    if (fs.existsSync('../../../data/GeoLite2-City.mmdb')) {
      expect(response.body[0].geo.regionName).toBe('Localhost');
    }

    ws.close();
  });

  ws.on('error', function incoming(data) {
    console.error(data);
    ws.close();
  });

  ws.on('close', function incoming(data) {
    done();
  });
});
