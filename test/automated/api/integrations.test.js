var request = require('supertest');
request = request('http://127.0.0.1:8080');

var accessToken = '';
var webhookID;

const webhook = 'https://super.duper.cool.thing.biz/owncast';
const events = ['CHAT'];

test('create webhook', async (done) => {
  const res = await sendIntegrationsChangePayload('webhooks/create', {
    url: webhook,
    events: events,
  });

  expect(res.body.url).toBe(webhook);
  expect(res.body.timestamp).toBeTruthy();
  expect(res.body.events).toStrictEqual(events);
  done();
});

test('check webhooks', (done) => {
  request
    .get('/api/admin/webhooks')
    .auth('admin', 'abc123')
    .expect(200)
    .then((res) => {
      expect(res.body).toHaveLength(1);
      expect(res.body[0].url).toBe(webhook);
      expect(res.body[0].events).toStrictEqual(events);
      webhookID = res.body[0].id;
      done();
    });
});

test('delete webhook', async (done) => {
  const res = await sendIntegrationsChangePayload('webhooks/delete', {
    id: webhookID,
  });
  expect(res.body.success).toBe(true);
  done();
});

test('check that webhook was deleted', (done) => {
  request
    .get('/api/admin/webhooks')
    .auth('admin', 'abc123')
    .expect(200)
    .then((res) => {
      expect(res.body).toHaveLength(0);
      done();
    });
});

test('create access token', async (done) => {
  const name = 'Automated integration test';
  const scopes = [
    'CAN_SEND_SYSTEM_MESSAGES',
    'CAN_SEND_MESSAGES',
    'HAS_ADMIN_ACCESS',
  ];
  const res = await sendIntegrationsChangePayload('accesstokens/create', {
    name: name,
    scopes: scopes,
  });

  expect(res.body.accessToken).toBeTruthy();
  expect(res.body.createdAt).toBeTruthy();
  expect(res.body.displayName).toBe(name);
  expect(res.body.scopes).toStrictEqual(scopes);
  accessToken = res.body.accessToken;

  done();
});

test('check access tokens', async (done) => {
  const res = await request
    .get('/api/admin/accesstokens')
    .auth('admin', 'abc123')
    .expect(200);
  const tokenCheck = res.body.filter(
    (token) => token.accessToken === accessToken
  );
  expect(tokenCheck).toHaveLength(1);
  done();
});

test('send a system message using access token', async (done) => {
  const payload = {
    body: 'This is a test system message from the automated integration test',
  };
  const res = await request
    .post('/api/integrations/chat/system')
    .set('Authorization', 'Bearer ' + accessToken)
    .send(payload)
    .expect(200);
  done();
});

test('send an external integration message using access token', async (done) => {
  const payload = {
    body: 'This is a test external message from the automated integration test',
  };
  const res = await request
    .post('/api/integrations/chat/send')
    .set('Authorization', 'Bearer ' + accessToken)
    .send(payload)
    .expect(200);
  done();
});

test('send an external integration action using access token', async (done) => {
  const payload = {
    body: 'This is a test external action from the automated integration test',
  };
  const res = await request
    .post('/api/integrations/chat/action')
    .set('Authorization', 'Bearer ' + accessToken)
    .send(payload)
    .expect(200);
  done();
});

test('test fetch chat history using access token', async (done) => {
  const res = await request
    .get('/api/integrations/chat')
    .set('Authorization', 'Bearer ' + accessToken)
    .expect(200);
  done();
});

test('test fetch chat history failure using invalid access token', async (done) => {
  const res = await request
    .get('/api/integrations/chat')
    .set('Authorization', 'Bearer ' + 'invalidToken')
    .expect(401);
  done();
});

test('test fetch chat history OPTIONS request', async (done) => {
  const res = await request
    .options('/api/integrations/chat')
    .set('Authorization', 'Bearer ' + accessToken)
    .expect(204);
  done();
});

test('delete access token', async (done) => {
  const res = await sendIntegrationsChangePayload('accesstokens/delete', {
    token: accessToken,
  });
  expect(res.body.success).toBe(true);
  done();
});

test('check token delete was successful', async (done) => {
  const res = await request
    .get('/api/admin/accesstokens')
    .auth('admin', 'abc123')
    .expect(200);
  const tokenCheck = res.body.filter(
    (token) => token.accessToken === accessToken
  );
  expect(tokenCheck).toHaveLength(0);
  done();
});

async function sendIntegrationsChangePayload(endpoint, payload) {
  const url = '/api/admin/' + endpoint;
  const res = await request
    .post(url)
    .auth('admin', 'abc123')
    .send(payload)
    .expect(200);

  return res;
}
