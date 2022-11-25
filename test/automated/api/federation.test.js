var request = require('supertest')
const jsonfile = require('jsonfile')
const Ajv = require("ajv-draft-04")
const sendConfigChangeRequest = require('./lib/config').sendConfigChangeRequest;

request = request('http://127.0.0.1:8080');

var ajv = new Ajv();
var nodeInfoSchema = jsonfile.readFileSync('schema/node_info_2.0_schema.json');

test('disable federation', async (done) => {
  const res = await sendConfigChangeRequest('federation/enable', false);
  done();
});

test('verify responses of /.well-known/webfinger when federation is disabled', async (done) => {
  const res = request
  .get('/.well-known/webfinger')
  .expect(405);
  done();
});

test('verify responses of /.well-known/host-meta when federation is disabled', async (done) => {
  const res = request
  .get('/.well-known/host-meta')
  .expect(405);
  done();
});

test('verify responses of /.well-known/nodeinfo when federation is disabled', async (done) => {
  const res = request
  .get('/.well-known/nodeinfo')
  .expect(405);
  done();
});

test('verify responses of /.well-known/x-nodeinfo2 when federation is disabled', async (done) => {
  const res = request
  .get('/.well-known/x-nodeinfo2')
  .expect(405);
  done();
});

test('verify responses of /nodeinfo/2.0 when federation is disabled', async (done) => { 
  const res = request
  .get('/nodeinfo/2.0')
  .expect(405);
  done();
});

test('verify responses of /api/v1/instance when federation is disabled', async (done) => {
  const res = request
  .get('/api/v1/instance')
  .expect(405);
  done();
});

test('verify responses of /federation/user/ when federation is disabled', async (done) => {
  const res = request
  .get('/federation/user/')
  .expect(405);
  done();
});

test('verify responses of /federation/ when federation is disabled', async (done) => {
  const res = request
  .get('/federation/')
  .expect(405);
  done();
});

test('enable federation', async (done) => {
  const res = await sendConfigChangeRequest('federation/enable', true);
  done();
});

test('verify responses of /.well-known/webfinger when federation is enabled', async (done) => {
  const res = request
  .get('/.well-known/webfinger')
  .expect(200);
  done();
});

test('verify responses of /.well-known/host-meta when federation is enabled', async (done) => {
  const res = request
  .get('/.well-known/host-meta')
  .expect(200);
  done();
});

test('verify responses of /.well-known/nodeinfo when federation is enabled', async (done) => {
  const res = request
  .get('/.well-known/nodeinfo')
  .expect(200);
  done();
});

test('verify responses of /.well-known/x-nodeinfo2 when federation is enabled', async (done) => {
  const res = request
  .get('/.well-known/x-nodeinfo2')
  .expect(200);
  done();
});

test('verify responses of /nodeinfo/2.0 when federation is enabled', async (done) => { 
  const res = request
  .get('/nodeinfo/2.0')
  .expect(200);
  done();
});

test('verify responses of /api/v1/instance when federation is enabled', async (done) => {
  const res = request
  .get('/api/v1/instance')
  .expect(200);
  done();
});

test('verify responses of /federation/user/ when federation is enabled', async (done) => {
  const res = request
  .get('/federation/user/')
  .expect(200);
  done();
});

test('verify responses of /federation/ when federation is enabled', async (done) => {
  const res = request
  .get('/federation/')
  .expect(200);
  done();
});

test('verify nodeinfo 2.0 is valid', (done) => {
  request.get('/nodeinfo/2.0').expect(200)
    .then((res) => {
      expect(ajv.validate(nodeInfoSchema, res.body)).toBe(true);
      done();
    });
});


