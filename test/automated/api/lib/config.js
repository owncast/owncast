var request = require('supertest');
request = request('http://127.0.0.1:8080');

const defaultAdminPassword = 'abc123';

async function sendConfigChangeRequest(endpoint, value, adminPassword = defaultAdminPassword) {
  const url = '/api/admin/config/' + endpoint;
  const res = await request
    .post(url)
    .auth('admin', adminPassword)
    .send({ value: value })
    .expect(200);

  expect(res.body.success).toBe(true);
  return res;
}

async function sendConfigChangePayload(endpoint, payload, adminPassword = defaultAdminPassword) {
  const url = '/api/admin/config/' + endpoint;
  const res = await request
    .post(url)
    .auth('admin', adminPassword)
    .send(payload)
    .expect(200);

  expect(res.body.success).toBe(true);

  return res;
}


module.exports.sendConfigChangeRequest = sendConfigChangeRequest;
module.exports.sendConfigChangePayload = sendConfigChangePayload;
