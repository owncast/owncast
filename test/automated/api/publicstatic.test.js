var request = require('supertest');
request = request('http://127.0.0.1:8080');
const fs = require('fs');
const path = require('path');
const randomString = require('./lib/rand').randomString;

const publicPath = path.resolve(__dirname, '../../../data/public');
const filename = randomString() + '.txt';
const fileContent = randomString();

test('random public static file does not exist', async (done) => {
	request.get('/public/' + filename).expect(404);

	done();
});

test('public directory is writable', async (done) => {
	try {
		writeFileToPublic();
	} catch (err) {
		if (err) {
			if (err.code === 'ENOENT') {
				// path does not exist
				fs.mkdirSync(publicPath);
				writeFileToPublic();
			} else {
				throw err;
			}
		}
	}

	done();
});

test('public static file is accessible', async (done) => {
	request
		.get('/public/' + filename)
		.expect(200)
		.then((res) => {
			expect(res.text).toEqual(fileContent);
			done();
		});
});

test('public static file is persistent and not locked', async (done) => {
	fs.unlink(path.join(publicPath, filename), (err) => {
		if (err) {
			throw err;
		}
	});
	done();
});

function writeFileToPublic() {
	fs.writeFileSync(path.join(publicPath, filename), fileContent);
}
