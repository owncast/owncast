var request = require('supertest');
request = request('http://127.0.0.1:8080');
const fs = require('fs');
const path = require('path');
const randomString = require('./lib/rand').randomString;

const publicPath = path.resolve(__dirname, '../../../data/public');
const filename = randomString() + '.txt';
const fileContent = randomString();

test('random public static file does not exist', async () => {
	await request.get('/public/' + filename).expect(404);
});

test('public directory is writable', async () => {
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
});

test('public static file is accessible', async () => {
	await request
		.get('/public/' + filename)
		.expect(200)
		.then((res) => {
			expect(res.text).toEqual(fileContent);
		});
});

test('public static file is persistent and not locked', async () => {
	fs.unlink(path.join(publicPath, filename), (err) => {
		if (err) {
			throw err;
		}
	});
});

function writeFileToPublic() {
	fs.writeFileSync(path.join(publicPath, filename), fileContent);
}
