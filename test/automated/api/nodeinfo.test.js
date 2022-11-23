var request = require('supertest')
const jsonfile = require('jsonfile')
const Ajv = require("ajv-draft-04")

request = request('http://127.0.0.1:8080');

var ajv = new Ajv();

var schema = jsonfile.readFileSync('schema/node_info_2.0_schema.json');

test('nodeinfo 2.0 is valid', (done) => {
    get('/nodeinfo/2.0').expect(200)
        .then((res) => {
            expect(ajv.validate(schema, res.body)).toBe(true);
            done();
        });
});
