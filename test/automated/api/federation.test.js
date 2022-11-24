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

test('check disabled federation responses of /.well-known/webfinger', async (done) => {
  const res = request
  .get('/.well-known/webfinger')
  .expect(405);
  done();
});

test('check disabled federation responses of /.well-known/host-meta', async (done) => {
  const res = request
  .get('/.well-known/host-meta')
  .expect(405);
  done();
});

test('check disabled federation responses of /.well-known/nodeinfo', async (done) => {
  const res = request
  .get('/.well-known/nodeinfo')
  .expect(405);
  done();
});

test('check disabled federation responses of /.well-known/x-nodeinfo2', async (done) => {
  const res = request
  .get('/.well-known/x-nodeinfo2')
  .expect(405);
  done();
});

test('check disabled federation responses of /nodeinfo/2.0', async (done) => { 
  const res = request
  .get('/nodeinfo/2.0')
  .expect(405);
  done();
});

test('check disabled federation responses of /api/v1/instance', async (done) => {
  const res = request
  .get('/api/v1/instance')
  .expect(405);
  done();
});

test('check disabled federation responses of /federation/user/', async (done) => {
  const res = request
  .get('/federation/user/')
  .expect(405);
  done();
});

test('check disabled federation responses of /federation/', async (done) => {
  const res = request
  .get('/federation/')
  .expect(405);
  done();
});

test('enable federation', async (done) => {
  const res = await sendConfigChangeRequest('federation/enable', true);
  done();
});

test('check enabled federation responses of /.well-known/webfinger', async (done) => {
  const res = request
  .get('/.well-known/webfinger')
  .expect(200);
  done();
});

test('check enabled federation responses of /.well-known/host-meta', async (done) => {
  const res = request
  .get('/.well-known/host-meta')
  .expect(200);
  done();
});

test('check enabled federation responses of /.well-known/nodeinfo', async (done) => {
  const res = request
  .get('/.well-known/nodeinfo')
  .expect(200);
  done();
});

test('check enabled federation responses of /.well-known/x-nodeinfo2', async (done) => {
  const res = request
  .get('/.well-known/x-nodeinfo2')
  .expect(200);
  done();
});

test('check enabled federation responses of /nodeinfo/2.0', async (done) => { 
  const res = request
  .get('/nodeinfo/2.0')
  .expect(200);
  done();
});

test('check enabled federation responses of /api/v1/instance', async (done) => {
  const res = request
  .get('/api/v1/instance')
  .expect(200);
  done();
});

test('check enabled federation responses of /federation/user/', async (done) => {
  const res = request
  .get('/federation/user/')
  .expect(200);
  done();
});

test('check enabled federation responses of /federation/', async (done) => {
  const res = request
  .get('/federation/')
  .expect(200);
  done();
});

test('nodeinfo 2.0 is valid', (done) => {
  request.get('/nodeinfo/2.0').expect(200)
    .then((res) => {
      expect(ajv.validate(nodeInfoSchema, res.body)).toBe(true);
      done();
    });
});


