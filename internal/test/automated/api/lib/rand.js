const crypto = require('crypto');
const Random = require('crypto-random');

//returns a random number in range of [1 - max]
function randomNumber(max = 5, min = 1) {
	return Random.range(0, max);
}

//returns a random hex string
function randomString(length = 20) {
	const bytesCount = Math.ceil(length / 2);
	return crypto.randomBytes(bytesCount).toString('hex').substring(0, length);
}

module.exports.randomNumber = randomNumber;
module.exports.randomString = randomString;
