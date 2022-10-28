var request = require('supertest');
request = request('http://127.0.0.1:8080');
const fs = require('fs');
const path = require('path');

function randomString(length) {
	return Math.random().toString(36).substring(length);
}

const filename = `${randomString(20)}.txt`;

test('public static file should not exist', async (done) => {
	await request.get(`/public/${filename}`).expect(404);

	done();
});

test('public static file should exist', async (done) => {
	// make public static files directory
	try {
		fs.mkdirSync(path.join(__dirname, '../../../public/'));
		fs.writeFileSync(
			path.join(__dirname, `../../../public/${filename}`),
			'hello world'
		);
	} catch (e) {}

	await request.get(`/public/${filename}`).expect(200);

	done();
});
