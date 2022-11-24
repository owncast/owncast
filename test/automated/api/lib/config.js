var request = require('supertest');
request = request('http://127.0.0.1:8080');

const defaultStreamKey = 'abc123';

async function sendConfigChangeRequest(endpoint, value) {
    const url = '/api/admin/config/' + endpoint;
    const res = await request
      .post(url)
      .auth('admin', defaultStreamKey)
      .send({ value: value })
      .expect(200);
  
    expect(res.body.success).toBe(true);
    return res;
  }
  
  async function sendConfigChangePayload(endpoint, payload) {
    const url = '/api/admin/config/' + endpoint;
    const res = await request
      .post(url)
      .auth('admin', defaultStreamKey)
      .send(payload)
      .expect(200);
  
    expect(res.body.success).toBe(true);
  
    return res;
  }


module.exports.sendConfigChangeRequest = sendConfigChangeRequest;
module.exports.sendConfigChangePayload = sendConfigChangePayload;