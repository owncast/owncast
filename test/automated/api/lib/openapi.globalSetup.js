const fs = require('fs');
const { RESULTS_DIR, EXERCISED_FILE } = require('./openapi');

module.exports = async () => {
	fs.mkdirSync(RESULTS_DIR, { recursive: true });
	fs.writeFileSync(EXERCISED_FILE, '');
};
