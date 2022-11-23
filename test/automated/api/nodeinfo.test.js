var request = require('supertest') 
const Ajv = require("ajv")
var readTextFile = require('read-text-file');

request = request('http://127.0.0.1:8080');
const ajv = new Ajv();
 
var schema = readTextFile.readSync('schema/node_info_2.0_schema.json');

test('nodeinfo 2.0 is valid', (done) => {
    request.get('/nodeinfo/2.0').expect(200)
        .then((res) => {
            expect(ajv.validate(schema, res.body)).toBe(true);
            done();
        });
});
