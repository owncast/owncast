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

test('check disabled federation responses', async (done) => {
    const res1 = request
    .get('/.well-known/webfinger')
    .expect(405);

    const res2 = request
    .get('/.well-known/host-meta')
    .expect(405);
  
    const res3 = request
    .get('/.well-known/nodeinfo')
    .expect(405);

    const res4 = request
    .get('/.well-known/x-nodeinfo2')
    .expect(405);
    
    const res5 = request
    .get('/nodeinfo/2.0')
    .expect(405);
    
    const res6 = request
    .get('/api/v1/instance')
    .expect(405);

    const res7 = request
    .get('/federation/user/')
    .expect(405);

    const res8 = request
    .get('/federation/')
    .expect(405);

    done();
  });

  test('enable federation', async (done) => {
    const res = await sendConfigChangeRequest('federation/enable', true);
    done();
  });

test('nodeinfo 2.0 is valid', (done) => {
    request.get('/nodeinfo/2.0').expect(200)
        .then((res) => {
            expect(ajv.validate(nodeInfoSchema, res.body)).toBe(true);
            done();
        });
});


